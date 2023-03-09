package cmd

import (
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
	rootCmd.Flags().BoolP("confluence", "c", false, "scan confluence")
	rootCmd.Flags().BoolP("all-rules", "r", false, "use all rules")
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

	// Run with default configuration
	if allRules {
		wrap := wrapper.NewWrapper()

		for find := range wrap.Detect("sfafaf") {
			fmt.Println(find)
		}

	}
}
