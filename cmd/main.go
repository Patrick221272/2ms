package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "2ms",
	Short: "2ms Secrets Detection",
}

func init() {
	cobra.OnInitialize(initLog)
	rootCmd.PersistentFlags().BoolP("all", "a", true, "scan all plugins")
	rootCmd.PersistentFlags().BoolP("confluence", "c", false, "scan confluence")
}

func initLog() {

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
