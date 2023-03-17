package cmd

import (
	"2ms/Reporting"
	"2ms/plugins"
	"2ms/wrapper"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "2ms",
	Short: "2ms Secrets Detection",
	Run:   runDetection,
}

func init() {
	cobra.OnInitialize(initLog)
	rootCmd.Flags().BoolP("all", "a", true, "scan all plugins")
	rootCmd.Flags().StringP("confluence", "c", "", "scan confluence url")
	rootCmd.Flags().StringP("confluence-user", "", "", "confluence username or email")
	rootCmd.Flags().StringP("confluence-token", "", "", "confluence token")
	rootCmd.Flags().StringSlice("rules", []string{"all"}, "select rules to be applied")
	rootCmd.PersistentFlags().StringP("log-level", "l", "info", "log level (trace, debug, info, warn, error, fatal)")
}

func initLog() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	ll, err := rootCmd.Flags().GetString("log-level")
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	switch strings.ToLower(ll) {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "err", "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Msg(err.Error())
	}
}

func isValidFilter(rulesFilter []string) bool {
	for _, filter := range rulesFilter {
		if strings.EqualFold(filter, "all") || strings.EqualFold(filter, "token") || strings.EqualFold(filter, "key") || strings.EqualFold(filter, "id") {
			return true
		}
	}
	return false
}

func runDetection(cmd *cobra.Command, args []string) {
	rulesFilter, err := cmd.Flags().GetStringSlice("rules")
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	if !isValidFilter(rulesFilter) {
		log.Fatal().Msg(`rules filter allowed: "all", "token", "id", "key"`)
	}

	// Get desired plugins content
	plugins := plugins.NewPlugins()
	//allPlugins, _ := cmd.Flags().GetBool("all")

	confluence, _ := cmd.Flags().GetString("confluence")
	confluenceUser, _ := cmd.Flags().GetString("confluence-user")
	confluenceToken, _ := cmd.Flags().GetString("confluence-token")

	if confluence != "" {
		plugins.AddPlugin("confluence", confluence, confluenceUser, confluenceToken)
	}

	contents := plugins.RunPlugins()

	report := Reporting.Report{}
	report.Results = make(map[string][]Reporting.Secret)

	// Run with default configuration
	wrap := wrapper.NewWrapper(rulesFilter)

	for _, c := range contents {
		secrets := wrap.Detect(c.Content)
		for _, secret := range secrets {
			report.Results[c.OriginalUrl] = append(report.Results[c.OriginalUrl], secret)
		}
	}
	report.TotalItemsScanned = len(contents)
	Reporting.ShowReport(report)
}
