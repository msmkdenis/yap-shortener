// Package stylecheck contains analyzes that enforce style rules.
// Includes the following list of checks:
//
//	ST1001	Dot imports are discouraged
//	ST1003	Poorly chosen identifier
//	ST1005	Incorrectly formatted error string
//	ST1006	Poorly chosen receiver name
//	ST1008	A function's error value should be its last return value
//	ST1011	Poorly chosen name for variable of type time.Duration
//	ST1012	Poorly chosen name for error variable
//	ST1013	Should use constants for HTTP error codes, not magic numbers
//	ST1015	A switch's default case should be the first or last case
//	ST1016	Use consistent method receiver names
//	ST1017	Don't use Yoda conditions
//	ST1018	Avoid zero-width and control characters in string literals
//	ST1019	Importing the same package multiple times
//	ST1021	The documentation of an exported type should start with type's name
//	ST1022	The documentation of an exported variable or constant should start with variable's name
//	ST1023	Redundant type in variable declaration
package stylecheck

import (
	"golang.org/x/tools/go/analysis"
	"honnef.co/go/tools/stylecheck"
)

// GetStyleCheckAnalyzers returns a list of stylecheck analyzers.
func GetStyleCheckAnalyzers() []*analysis.Analyzer {
	var response []*analysis.Analyzer

	exclude := map[string]bool{
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
