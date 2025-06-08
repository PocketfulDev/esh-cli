package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var version = "dev"

// SetVersion sets the version for the CLI
func SetVersion(v string) {
	version = v
	rootCmd.Version = v
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "esh-cli",
	Short: "ESH CLI tool for managing git tags and deployments",
	Long: `ESH CLI tool for adding and pushing hot fix tags.
Tag format is env_major.minor.patch-release
In some projects this triggers deployment in CircleCI.`,
}

// NewRootCmd creates a new isolated root command instance for testing
// This prevents race conditions when running tests concurrently
func NewRootCmd(cmdVersion string) *cobra.Command {
	// Create isolated config file variable to avoid race conditions
	var isolatedCfgFile string

	cmd := &cobra.Command{
		Use:   "esh-cli",
		Short: "ESH CLI tool for managing git tags and deployments",
		Long: `ESH CLI tool for adding and pushing hot fix tags.
Tag format is env_major.minor.patch-release
In some projects this triggers deployment in CircleCI.`,
		Version: cmdVersion,
	}

	// Add persistent flags with isolated variable
	cmd.PersistentFlags().StringVar(&isolatedCfgFile, "config", "", "config file (default is $HOME/.esh-cli.yaml)")

	// Add all subcommands to the new instance
	addSubcommands(cmd)

	return cmd
}

// addSubcommands adds all subcommands to the given root command
func addSubcommands(cmd *cobra.Command) {
	// Add all the subcommands that are normally added via init() functions
	cmd.AddCommand(addTagCmd)
	cmd.AddCommand(lastTagCmd)
	cmd.AddCommand(initCmd)
	cmd.AddCommand(projectsCmd)
	cmd.AddCommand(branchVersionCmd)
	cmd.AddCommand(bumpVersionCmd)
	cmd.AddCommand(changelogCmd)
	cmd.AddCommand(versionDiffCmd)
	cmd.AddCommand(versionListCmd)
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.esh-cli.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".esh-cli")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {
		// Config file not found - check if we should auto-initialize
		// Only auto-initialize if this is not the init command itself
		if shouldAutoInitialize() {
			fmt.Fprintln(os.Stderr, "🤖 No configuration found. Consider running 'esh-cli init' for AI project discovery.")
		}
	}
}

// shouldAutoInitialize checks if we should show auto-initialization message
func shouldAutoInitialize() bool {
	// Don't show message if running init command
	if len(os.Args) > 1 {
		command := strings.ToLower(os.Args[1])
		if command == "init" || command == "help" || command == "--help" || command == "-h" {
			return false
		}
	}
	return true
}
