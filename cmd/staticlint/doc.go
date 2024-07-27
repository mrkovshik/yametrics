// Package main is a static analysis tool that runs various analyzers on Go code.
//
// The tool uses a variety of analyzers from the "golang.org/x/tools/go/analysis"
// package as well as additional analyzers from the "honnef.co/go/tools" package.
//
// This package provides the main entry point for the tool, which configures and runs
// the analyzers using the main package.
//
// Available analyzers include:
// - asmdecl: Reports mismatches between assembly and Go declarations.
// - assign: Checks for useless assignments.
// - atomic: Detects common mistakes with the sync/atomic package.
// - bools: Checks for common mistakes involving booleans.
// - buildtag: Ensures that go:build tags are correctly formed.
// - cgocall: Detects violations of the cgo pointer passing rules.
// - composite: Checks for unkeyed composite literals.
// - copylock: Reports locks erroneously passed by value.
// - httpresponse: Checks for mistakes using net/http's ResponseWriter.
// - loopclosure: Checks for references to loop variables from within nested functions.
// - lostcancel: Checks for failure to call a context cancellation function.
// - nilfunc: Checks for useless comparisons between functions and nil.
// - pkgfact: Checks facts about packages.
// - printf: Checks consistency of Printf format strings and arguments.
// - shadow: Checks for shadowed variables.
// - shift: Checks for shifts that exceed the width of an integer.
// - stdmethods: Checks for misspellings of method names such as Write and Read.
// - stringintconv: Checks for conversions of integers to strings.
// - structtag: Checks for struct tags not conforming to reflect.StructTag.Get.
// - testinggoroutine: Checks for goroutine leaks in tests.
// - tests: Checks for common mistaken usages of the testing package.
// - unmarshal: Checks for unmarshal errors.
// - unsafeptr: Checks for misuse of unsafe.Pointer.
//
// Additionally, it includes analyzers from staticcheck, simplecheck, and stylecheck packages.
//
// The OSExitAnalyzer checks for usage of os.Exit in the main function of the main package.
// It reports any instances of os.Exit usage found in the main function, as it is often
// considered a bad practice due to its abrupt termination of the program.
//
// Usage:
//
//	go run main.go <packages>
//
// The tool is configured by listing the analyzers in the main function and passing them
// to the multichecker.Main function.
package main
