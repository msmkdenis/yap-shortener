package linter

import (
	"golang.org/x/tools/go/analysis"
	"honnef.co/go/tools/staticcheck"
)

// GetStaticCheckAnalyzers returns a list of staticcheck analyzers.
func GetStaticCheckAnalyzers() []*analysis.Analyzer {
	var response []*analysis.Analyzer

	for _, check := range staticcheck.Analyzers {
		response = append(response, check.Analyzer)
	}

	return response
}
