package linter

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
)

// Run runs all analyzers.
func Run() {
	var analyzers []*analysis.Analyzer
	analyzers = append(analyzers, GetXPassesAnalyzers()...)
	analyzers = append(analyzers, GetStaticCheckAnalyzers()...)
	analyzers = append(analyzers, GetStyleCheckAnalyzers()...)
	analyzers = append(analyzers, GetExternalAnalyzers()...)
	analyzers = append(analyzers, NoExitInMainAnalyzer)

	multichecker.Main(analyzers...)
}
