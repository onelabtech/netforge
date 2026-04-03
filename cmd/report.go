package cmd

import (
	"encoding/json"
	"fmt"
	"netforge/internal/model"
	"os"

	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:   "report [json-file]",
	Short: "Generate a human-readable summary from a JSON report",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file, err := os.Open(args[0])
		if err != nil {
			return err
		}
		defer file.Close()

		var report model.SummaryReport
		if err := json.NewDecoder(file).Decode(&report); err != nil {
			return err
		}

		fmt.Printf("# NetForge Diagnostic Report for %s\n", report.Target)
		fmt.Printf("## Generated at: %s\n", report.Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Printf("## Overall Health: %s\n\n", report.OverallHealth)

		if len(report.Issues) > 0 {
			fmt.Println("### Identified Issues")
			for i, issue := range report.Issues {
				fmt.Printf("%d. [%s] %s\n", i+1, issue.Severity, issue.Title)
				fmt.Printf("   - Probable Cause: %s\n", issue.ProbableCause)
				fmt.Printf("   - Suggested Fix: %s\n", issue.SuggestedFix)
				fmt.Printf("   - Evidence: %s\n\n", issue.Evidence)
			}
		} else {
			fmt.Println("✓ No issues found.")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
}
