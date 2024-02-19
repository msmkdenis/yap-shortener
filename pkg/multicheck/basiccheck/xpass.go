// Package basicchecks defines a list of standard static analyzers.
// Includes the following list of checks:
//
// Includes the following list of checks:
// `asmdecl` checks for mismatches between assembly files and Go declarations.
// `assign` checks for useless assignments.
// `atomic` checks for common mistakes using the sync/atomic package.
// `atomicalign` checks for non-64-bit-aligned arguments to sync/atomic functions.
// `bools` checks for common mistakes involving boolean expressions.
// `cgocall` checks for calls to C code.
// `composites` checks for composite literals that can be simplified.
// `composite` checks for unkeyed composite literals.
// `copylocks` checks for locks erroneously passed by value.
// `deepequalerrors` checks for the use of reflect.DeepEqual with error values.
// `defers` checks for common mistakes in defer statements.
// `directive` checks known Go toolchain directives.
// `errorsas` checks that the second argument to errors.As is a pointer to a type implementing error.
// `httpresponse` checks for mistakes using HTTP responses.
// `ifaceassert` detect impossible interface-to-interface type assertions
// `loopclosure` checks for references to enclosing loop variables from within nested functions.
// `lostcancel` check cancel func returned by context.WithCancel is called
// `nilfunc` checks for useless comparisons between functions and nil.
// `nilness` check for redundant or impossible nil comparisons
// `printf` check consistency of Printf format strings and arguments
// `reflectvaluecompare` check for comparing reflect.Value values with == or reflect.DeepEqual
// `shadow` checks for shadowed variables.
// `sigchanyzer` check for unbuffered channel of os.Signal
// `slog` check for invalid structured logging calls
// `sortslice` checks for calls to sort.Slice that do not use a slice type as first argument.
// `stdmethods` check signature of methods of well-known interfaces
// `stringintconv` check for string(int) conversions
// `structtag` checks struct field tags are well formed.
// `testinggoroutine` report calls to (*testing.T).Fatal from goroutines started by a test.
// `tests` checks for common mistaken usages of tests and examples.
// `timeformat` check for calls of (time.Time).Format or time.Parse with 2006-02-01
// `unmarshal` checks for passing non-pointer or non-interface types to unmarshal and decode functions.
// `unreachable` checks for unreachable code.
// `unsafeptr` check for invalid conversions of uintptr to unsafe.Pointer
// `unusedresult` checks for unused results of calls to certain pure functions.
// `unusedwrite` checks for unused writes
package basiccheck

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/defers"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/slog"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"honnef.co/go/tools/analysis/facts/nilness"
)

// GetXPassesAnalyzers returns a list of x/tools/go/analysis/passes analyzers.
func GetXPassesAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		atomicalign.Analyzer,
		bools.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		deepequalerrors.Analyzer,
		defers.Analyzer,
		directive.Analyzer,
		errorsas.Analyzer,
		httpresponse.Analyzer,
		ifaceassert.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		nilness.Analysis,
		printf.Analyzer,
		reflectvaluecompare.Analyzer,
		shadow.Analyzer,
		sigchanyzer.Analyzer,
		slog.Analyzer,
		sortslice.Analyzer,
		stdmethods.Analyzer,
		stringintconv.Analyzer,
		structtag.Analyzer,
		testinggoroutine.Analyzer,
		tests.Analyzer,
		timeformat.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
		unusedwrite.Analyzer,
	}
}
