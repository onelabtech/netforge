package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "netforge",
	Short: "NetForge: Next-gen intelligent networking diagnostics CLI",
	Long: `NetForge is a world-class production-grade networking tool designed to 
replace fragmented networking utilities such as ping, dig, curl, and traceroute.
It provides unified, developer-first diagnostics with structured, actionable 
root-cause analysis.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.netforge.yaml)")

	// Global output flags
	rootCmd.PersistentFlags().BoolP("json", "j", false, "Output results in JSON format")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Only show critical output")
	rootCmd.PersistentFlags().IntP("timeout", "t", 10, "Default timeout in seconds")

	// Bind viper to pflags
	viper.BindPFlag("json", rootCmd.PersistentFlags().Lookup("json"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("quiet", rootCmd.PersistentFlags().Lookup("quiet"))
	viper.BindPFlag("timeout", rootCmd.PersistentFlags().Lookup("timeout"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".netforge" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".netforge")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("verbose") {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}
}

// GetConfigPath returns the default config path
func GetConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".netforge.yaml")
}

// IsJSON returns true if JSON output is requested
func IsJSON() bool {
	return viper.GetBool("json")
}

// IsVerbose returns true if verbose output is requested
func IsVerbose() bool {
	return viper.GetBool("verbose")
}

// GetTimeout returns the default timeout as a duration
func GetTimeout() int {
	return viper.GetInt("timeout")
}

// FormatOutput format strings for terminal
func FormatOutput(s string) string {
	// Basic string normalization for terminal output if needed
	return strings.TrimSpace(s)
}
