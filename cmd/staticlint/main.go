package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"

	"github.com/msmkdenis/yap-shortener/pkg/linter"
)

func main() {
	var analyzers []*analysis.Analyzer

	analyzers = append(analyzers, linter.GetXPassesAnalyzers()...)
	analyzers = append(analyzers, linter.GetStaticCheckAnalyzers()...)
	analyzers = append(analyzers, linter.GetStyleCheckAnalyzers()...)
	analyzers = append(analyzers, linter.GetExternalAnalyzers()...)
	analyzers = append(analyzers, linter.NoExitInMainAnalyzer)

	multichecker.Main(analyzers...)
}
