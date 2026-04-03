package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"netforge/internal/analyzer"
	"netforge/internal/checker"
	"netforge/internal/model"
	"netforge/internal/output"
	"netforge/internal/util"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// doctorCmd represents the doctor command
var doctorCmd = &cobra.Command{
	Use:   "doctor [target]",
	Short: "Run comprehensive intelligent diagnostics on a target",
	Long: `Doctor orchestrates multiple network checks (DNS, TCP, TLS, HTTP) in parallel 
and analyzes the results for probable causes and suggested fixes.

Example:
  netforge doctor google.com
  netforge doctor https://api.example.com`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := args[0]
		target, err := util.NormalizeTarget(input)
		if err != nil {
			return fmt.Errorf("invalid target: %v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(viper.GetInt("timeout"))*time.Second)
		defer cancel()

		fmt.Printf("🩺 Running diagnostics for %s...\n", target.Host)

		report := &model.SummaryReport{
			Target:    target.Host,
			Timestamp: time.Now(),
		}

		// Initialize checkers
		dnsChecker := checker.NewDNSChecker()
		tcpChecker := checker.NewTCPChecker()
		tlsChecker := checker.NewTLSChecker()
		httpChecker := checker.NewHTTPChecker()

		var wg sync.WaitGroup
		var mu sync.Mutex

		// Parallel Checks
		wg.Add(4)

		go func() {
			defer wg.Done()
			dnsRes, _ := dnsChecker.Check(ctx, target)
			mu.Lock()
			report.DNS = dnsRes
			mu.Unlock()
		}()

		go func() {
			defer wg.Done()
			tcpRes, _ := tcpChecker.Check(ctx, target)
			mu.Lock()
			report.TCP = tcpRes
			mu.Unlock()
		}()

		go func() {
			defer wg.Done()
			tlsRes, _ := tlsChecker.Check(ctx, target)
			mu.Lock()
			report.TLS = tlsRes
			mu.Unlock()
		}()

		go func() {
			defer wg.Done()
			httpRes, _ := httpChecker.Check(ctx, target)
			mu.Lock()
			report.HTTP = httpRes
			mu.Unlock()
		}()

		wg.Wait()

		// Analyze Results
		engine := analyzer.NewAnalyzer()
		report.Issues = engine.Analyze(report)

		// Determine Overall Health
		report.OverallHealth = model.StatusSuccess
		for _, issue := range report.Issues {
			if issue.Severity == model.StatusFailure {
				report.OverallHealth = model.StatusFailure
				break
			}
			if issue.Severity == model.StatusWarning {
				report.OverallHealth = model.StatusWarning
			}
		}

		// Output Results
		if viper.GetBool("json") {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(report)
		}

		renderer := output.NewRenderer()
		renderer.RenderDoctorReport(report)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
