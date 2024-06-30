package analizer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var OSExitAnalyzer = &analysis.Analyzer{
	Name: "osExit",
	Doc:  "check for usage of os.Exit in main function",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		// Check for main package
		if file.Name.Name != "main" {
			continue
		}
		// Inspect nodes of AST
		ast.Inspect(file, func(n ast.Node) bool {
			// Check for main function
			if fn, ok := n.(*ast.FuncDecl); ok {
				if fn.Name.Name == "main" && fn.Recv == nil {
					// Inspect the body of the main function for os.Exit calls
					ast.Inspect(fn.Body, func(n ast.Node) bool {
						if expr, ok := n.(*ast.CallExpr); ok {
							if fun, ok := expr.Fun.(*ast.SelectorExpr); ok {
								if pkgIdent, ok := fun.X.(*ast.Ident); ok {
									if pkgIdent.Name == "os" && fun.Sel.Name == "Exit" {
										pos := pass.Fset.Position(pkgIdent.NamePos)
										pass.Reportf(pkgIdent.NamePos, "usage of os.Exit in main at %d:%d", pos.Line, pos.Column)
										return false // Stop inspecting as we found os.Exit
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
