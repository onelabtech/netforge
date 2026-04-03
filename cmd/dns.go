package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"netforge/internal/checker"
	"netforge/internal/util"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var dnsCmd = &cobra.Command{
	Use:   "dns [domain]",
	Short: "Advanced DNS resolution lookup",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target, _ := util.NormalizeTarget(args[0])
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(viper.GetInt("timeout"))*time.Second)
		defer cancel()

		c := checker.NewDNSChecker()
		res, err := c.Check(ctx, target)
		if err != nil {
			return err
		}

		if viper.GetBool("json") {
			json.NewEncoder(os.Stdout).Encode(res)
			return nil
		}

		fmt.Printf("DNS Analysis for %s:\n", target.Host)
		fmt.Printf("  IPs (A):    %v\n", res.IPs)
		fmt.Printf("  IPv6s (AAAA): %v\n", res.IPv6s)
		fmt.Printf("  CNAMEs:     %v\n", res.CNAMEs)
		fmt.Printf("  MX:         %v\n", res.MX)
		fmt.Printf("  NS:         %v\n", res.NS)
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(dnsCmd)
}
