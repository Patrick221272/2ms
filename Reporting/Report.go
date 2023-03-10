package Reporting

import (
	"fmt"
	"github.com/zricethezav/gitleaks/v8/report"
	"strings"
)

func ShowReport(report Report) {
	fmt.Println("Summary:")
	fmt.Printf("- Total items scanned: %d\n", 1)
	fmt.Printf("- Total items with secrets: %d\n", 1)
	fmt.Println("Detailed Report:")
	generateResultsReport(report.Results)

}

func generateResultsReport(results map[string][]Secret) {
	for filepath, secrets := range results {
		itemLink := getItemId(filepath)
		fmt.Printf("- Item ID: %s\n", itemLink)
		fmt.Printf(" - Item Link: %s\n", filepath)
		fmt.Println("  - Secrets:")
		for _, secret := range secrets {
			fmt.Printf("   - Type: %s\n", secret.Description)
			fmt.Printf("    - Location: %d-%d\n", secret.StartLine, secret.EndLine)
			fmt.Printf("    - Value: %.20s\n", secret.Value)
		}
	}
}

func AddSecretToFile(report Report, value report.Finding, secret Secret) Report {
	report.Results[value.File] = append(report.Results[value.File], secret)
	return report
}

func CreateNewResult(report Report, value report.Finding, secret Secret) Report {
	results := make(map[string][]Secret)
	results[value.File] = append(results[value.File], secret)
	report.Results = results
	return report
}

func getItemId(fullPath string) string {
	itemLinkStrings := strings.Split(fullPath, "\\")
	itemLink := itemLinkStrings[len(itemLinkStrings)-1]
	return itemLink
}

type Report struct {
	Results map[string][]Secret
}

type Secret struct {
	Description string
	StartLine   int
	EndLine     int
	StartColumn int
	EndColumn   int
	Value       string
}
