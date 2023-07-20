package cmd

import (
	"os"

	"github.com/mona-actions/gh-migration-monitor/pkg/monitor"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "migration-monitor",
	Short: "gh cli extension to monitor migration status",
	Long:  `gh cli extension to monitor migration status`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },

	Run: func(cmd *cobra.Command, args []string) {
		orgName := cmd.Flag("organization").Value.String()
		token := cmd.Flag("github-token").Value.String()

		// Set the GitHub Organization
		os.Setenv("GHMM_GITHUB_ORGANIZATION", orgName)
		viper.BindEnv("GITHUB_ORGANIZATION")

		// Set the GitHub Token
		os.Setenv("GHMM_GITHUB_TOKEN", token)
		viper.BindEnv("GITHUB_TOKEN")

		// Call the monitor
		monitor.Organization()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gh-migration-monitor.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringP("organization", "o", "", "Organization to monitor")
	rootCmd.MarkFlagRequired("organization")

	// Not required because we can use the github token from the environment
	rootCmd.Flags().StringP("github-token", "t", "", "Github token to use")
}

func initConfig() {
	// Set env prefix
	viper.SetEnvPrefix("GHMM")

	// Read in environment variables that match
	// Specifically we are looking for GHMM_GITHUB_TOKEN
	viper.AutomaticEnv()
}
