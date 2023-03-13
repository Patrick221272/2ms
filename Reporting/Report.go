package Reporting

import (
	"fmt"
	"strings"
)

func ShowReport(report Report) {
	fmt.Println("Summary:")
	fmt.Printf("- Total items scanned: %d\n", 1)
	fmt.Printf("- Total items with secrets: %d\n", len(report.Results))
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

func (R *Report) AddSecret(source string, secret Secret) {
	if len(R.Results) > 0 {
		_, fileExist := R.Results[source]
		if fileExist {
			R.AddSecretToFile(source, secret)
		} else {
			R.CreateNewResult(source, secret)
		}
	}
	R.Results[source] = append(R.Results[source], secret)
}

func (R *Report) AddSecretToFile(source string, secret Secret) {
	R.Results[source] = append(R.Results[source], secret)
}

func (R *Report) CreateNewResult(source string, secret Secret) {
	results := make(map[string][]Secret)
	results[source] = append(results[source], secret)
	R.Results = results
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
