package analizer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// OSExitAnalyzer is an analyzer that checks for the usage of os.Exit in the main function of the main package.
var OSExitAnalyzer = &analysis.Analyzer{
	Name: "osExit",
	Doc:  "check for usage of os.Exit in main function",
	Run:  run,
}

// run is the function that implements the analysis logic.
func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		// Check if the file belongs to the main package.
		if file.Name.Name != "main" {
			continue
		}
		// Inspect nodes of the AST.
		ast.Inspect(file, func(n ast.Node) bool {
			// Check if the node is a function declaration.
			if fn, ok := n.(*ast.FuncDecl); ok {
				// Check if the function is the main function with no receiver.
				if fn.Name.Name == "main" && fn.Recv == nil {
					// Inspect the body of the main function for os.Exit calls.
					ast.Inspect(fn.Body, func(n ast.Node) bool {
						// Check if the node is a call expression.
						if expr, ok := n.(*ast.CallExpr); ok {
							// Check if the call expression is a selector expression.
							if fun, ok := expr.Fun.(*ast.SelectorExpr); ok {
								// Check if the selector expression's X is an identifier named "os" and the selector is "Exit".
								if pkgIdent, ok := fun.X.(*ast.Ident); ok {
									if pkgIdent.Name == "os" && fun.Sel.Name == "Exit" {
										// Report the position of the os.Exit call.
										pos := pass.Fset.Position(pkgIdent.NamePos)
										pass.Reportf(pkgIdent.NamePos, "usage of os.Exit in main at %d:%d", pos.Line, pos.Column)
										return false // Stop inspecting as we found os.Exit.
									}
								}
							}
						}
						return true
					})
				}
			}
			return true
		})
	}
	return nil, nil
}
