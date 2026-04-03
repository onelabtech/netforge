package cmd

import (
	"context"
	"fmt"
	"netforge/internal/checker"
	"netforge/internal/util"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var monitorCmd = &cobra.Command{
	Use:   "monitor [url]",
	Short: "Repeatedly monitor a target endpoint",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := args[0]
		target, _ := util.NormalizeTarget(input)
		interval, _ := cmd.Flags().GetDuration("interval")

		fmt.Printf("⏱ Monitoring %s every %s... (Ctrl+C to stop)\n", target.GetURL(), interval)

		c := checker.NewHTTPChecker()
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-sigChan:
				fmt.Println("\nStopped monitoring.")
				return nil
			case <-ticker.C:
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				res, err := c.Check(ctx, target)
				cancel()

				if err != nil {
					fmt.Printf("[%s] ERROR: %v\n", time.Now().Format("15:04:05"), err)
				} else {
					statusIcon := "✅"
					if res.StatusCode >= 400 {
						statusIcon = "❌"
					}
					fmt.Printf("[%s] %s Status: %d, Latency: %s\n", time.Now().Format("15:04:05"), statusIcon, res.StatusCode, res.Duration)
				}
			}
		}
	},
}

func init() {
	monitorCmd.Flags().DurationP("interval", "i", 30*time.Second, "Monitoring interval (e.g. 5s, 1m)")
	rootCmd.AddCommand(monitorCmd)
}
