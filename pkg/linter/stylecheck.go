package linter

import (
	"golang.org/x/tools/go/analysis"
	"honnef.co/go/tools/stylecheck"
)

// GetStyleCheckAnalyzers returns a list of stylecheck analyzers.
func GetStyleCheckAnalyzers() []*analysis.Analyzer {
	var response []*analysis.Analyzer

	var exclude = map[string]bool{
		"ST1000": true,
		"ST1020": true,
	}

	for _, check := range stylecheck.Analyzers {
		if !exclude[check.Analyzer.Name] {
			response = append(response, check.Analyzer)
		}
	}

	return response
}
