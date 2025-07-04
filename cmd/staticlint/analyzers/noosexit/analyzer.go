package noosexit

import (
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "noosexit",
	Doc:  "запрещает использование os.Exit в функции main пакета main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}

	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Name.Name != "main" || fn.Body == nil {
				continue
			}

			ast.Inspect(fn.Body, func(n ast.Node) bool {
				callExpr, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}
				if sel, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
					if ident, ok := sel.X.(*ast.Ident); ok &&
						ident.Name == "os" && sel.Sel.Name == "Exit" {

						obj := pass.TypesInfo.Uses[ident]
						if pkgName, ok := obj.(*types.PkgName); ok && pkgName.Imported().Path() == "os" {
							pass.Reportf(callExpr.Pos(), "нельзя использовать os.Exit в функции main")
						}
					}
				}
				return true
			})
		}
	}
	return nil, nil
}
