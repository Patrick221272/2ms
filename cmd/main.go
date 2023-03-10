package cmd

import (
	"2ms/plugins"
	"2ms/wrapper"
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "2ms",
	Short: "2ms Secrets Detection",
	Run:   runDetection,
}

func init() {
	rootCmd.Flags().BoolP("all", "a", true, "scan all plugins")
	rootCmd.Flags().StringP("confluence", "c", "", "scan confluence url")
	rootCmd.Flags().StringP("confluence-user", "", "", "confluence username or email")
	rootCmd.Flags().StringP("confluence-token", "", "", "confluence token")
	rootCmd.Flags().BoolP("all-rules", "r", true, "use all rules")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func runDetection(cmd *cobra.Command, args []string) {
	allRules, err := cmd.Flags().GetBool("all-rules")
	if err != nil {
		panic(err)
	}

	// Get desired plugins content
	plugins := plugins.NewPlugins()
	allPlugins, _ := cmd.Flags().GetBool("all")

	confluence, _ := cmd.Flags().GetString("confluence")
	confluenceUser, _ := cmd.Flags().GetString("confluence-user")
	confluenceToken, _ := cmd.Flags().GetString("confluence-token")

	if confluence != "" || allPlugins {
		plugins.AddPlugin("confluence", confluence, confluenceUser, confluenceToken)
	}

	// Run with default configuration
	if allRules {
		wrap := wrapper.NewWrapper()

		for find := range wrap.Detect("sfafaf") {
			fmt.Println(find)
		}

	}
}
