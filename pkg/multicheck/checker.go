package multicheck

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"

	"github.com/msmkdenis/yap-shortener/pkg/multicheck/basiccheck"
	"github.com/msmkdenis/yap-shortener/pkg/multicheck/customcheck/noexitmaincheck"
	"github.com/msmkdenis/yap-shortener/pkg/multicheck/staticcheck"
	"github.com/msmkdenis/yap-shortener/pkg/multicheck/stylecheck"
)

// Run runs all analyzers.
func Run() {
	var analyzers []*analysis.Analyzer
	analyzers = append(analyzers, basiccheck.GetXPassesAnalyzers()...)
	analyzers = append(analyzers, staticcheck.GetStaticCheckAnalyzers()...)
	analyzers = append(analyzers, stylecheck.GetStyleCheckAnalyzers()...)
	analyzers = append(analyzers, noexitmaincheck.NoExitInMainAnalyzer)

	multichecker.Main(analyzers...)
}
