/*
Package analizer provides a static analysis tool to detect the usage of os.Exit in the main function
of the main package.

The OSExitAnalyzer checks Go source files for calls to os.Exit within the main function of the main
package, and reports the line and column of each occurrence. This can help developers identify and
refactor code that might be prematurely terminating their programs.

Example Usage:

To use this analyzer in your project, you can run it as part of your go vet toolchain or integrate
it with other analysis tools. Here's an example of how to register and run the analyzer:

	package main

	import (
		"golang.org/x/tools/go/analysis/singlechecker"
		"github.com/yourusername/yourrepo/analizer"
	)

	func main() {
		singlechecker.Main(analizer.OSExitAnalyzer)
	}

This will run the osExit analyzer on your Go source files and report any occurrences of os.Exit in
the main function of the main package.
*/
package analizer
