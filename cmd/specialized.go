package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"netforge/internal/checker"
	"netforge/internal/util"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var httpCmd = &cobra.Command{
	Use:   "http [url]",
	Short: "Advanced HTTP request and performance inspection",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target, _ := util.NormalizeTarget(args[0])
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(viper.GetInt("timeout"))*time.Second)
		defer cancel()

		c := checker.NewHTTPChecker()
		res, err := c.Check(ctx, target)
		if err != nil {
			return err
		}

		if viper.GetBool("json") {
			json.NewEncoder(os.Stdout).Encode(res)
			return nil
		}

		fmt.Printf("HTTP Analysis for %s:\n", target.GetURL())
		fmt.Printf("  Status:      %d %s\n", res.StatusCode, res.Status)
		fmt.Printf("  Protocol:    %s\n", res.Proto)
		fmt.Printf("  Duration:    %s\n", res.Duration)
		fmt.Printf("  TTFB:        %s\n", res.TTFB)
		fmt.Printf("  Content-Len: %d\n", res.ContentLength)
		fmt.Printf("  Compression: %s\n", res.Compression)
		
		if len(res.Headers) > 0 {
			fmt.Println("  Headers:")
			for k, v := range res.Headers {
				fmt.Printf("    %s: %s\n", k, v)
			}
		}

		return nil
	},
}

var tlsCmd = &cobra.Command{
	Use:   "tls [host:port]",
	Short: "Advanced TLS/SSL handshake and certificate inspection",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := args[0]
		if !strings.Contains(input, ":") {
			input = input + ":443"
		}
		target, _ := util.NormalizeTarget(input)
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(viper.GetInt("timeout"))*time.Second)
		defer cancel()

		c := checker.NewTLSChecker()
		res, err := c.Check(ctx, target)
		if err != nil {
			return err
		}

		if viper.GetBool("json") {
			json.NewEncoder(os.Stdout).Encode(res)
			return nil
		}

		fmt.Printf("TLS Analysis for %s:%d:\n", target.Host, target.Port)
		fmt.Printf("  Version:      %s\n", res.Version)
		fmt.Printf("  Cipher:       %s\n", res.CipherSuite)
		fmt.Printf("  Subject:      %s\n", res.Subject)
		fmt.Printf("  Issuer:       %s\n", res.Issuer)
		fmt.Printf("  Expires:      %s (%d days left)\n", res.Expiry.Format("2006-01-02"), res.DaysUntilExpiry)
		fmt.Printf("  SANs:         %v\n", res.SANs)
		fmt.Printf("  Chain Valid:  %v\n", res.ChainValidated)

		return nil
	},
}

var scanCmd = &cobra.Command{
	Use:   "scan [target]",
	Short: "Safe TCP-connect port scan",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target, _ := util.NormalizeTarget(args[0])
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		commonPorts := []int{21, 22, 23, 25, 53, 80, 110, 143, 443, 445, 993, 995, 3306, 3389, 5432, 8080, 8443}
		
		fmt.Printf("🔍 Scanning common ports on %s...\n", target.Host)

		c := checker.NewTCPChecker()
		res, err := c.ScanPorts(ctx, target.Host, commonPorts)
		if err != nil {
			return err
		}

		if viper.GetBool("json") {
			json.NewEncoder(os.Stdout).Encode(res)
			return nil
		}

		fmt.Println("Scanned Ports:")
		for _, p := range res.Ports {
			status := "CLOSED"
			if p.Status == "open" {
				status = "OPEN"
			}
			fmt.Printf("  Port %d:\t%s\n", p.Port, status)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
	rootCmd.AddCommand(tlsCmd)
	rootCmd.AddCommand(scanCmd)
}
