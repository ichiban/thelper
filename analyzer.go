package unmarkedhelper

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer reports unmarked test helpers.
var Analyzer = &analysis.Analyzer{
	Name:     "unmarkedhelper",
	Doc:      `report unmarked test helpers`,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	var t *ast.Ident
	inspect.Nodes([]ast.Node{
		(*ast.File)(nil),
		(*ast.FuncDecl)(nil),
		(*ast.CallExpr)(nil),
	}, func(n ast.Node, push bool) bool {
		if !push {
			if _, ok := n.(*ast.FuncDecl); ok && t != nil {
				pass.Reportf(n.Pos(), "unmarked test helper: should call %s.Helper()", t.Name)
			}

			return false
		}

		switch n := n.(type) {
		case *ast.File:
			f := pass.Fset.File(n.Pos())
			return strings.HasSuffix(f.Name(), "_test.go")
		case *ast.FuncDecl:
			if strings.HasPrefix(n.Name.Name, "Test") {
				return false
			}

			t = containsT(n.Type.Params, pass.TypesInfo)
			return true
		case *ast.CallExpr:
			if helper(n.Fun, pass.TypesInfo) {
				t = nil
				return false
			}
			return true
		default:
			panic(n)
		}
	})

	return nil, nil
}

func containsT(n ast.Node, t *types.Info) *ast.Ident {
	switch n := n.(type) {
	case *ast.FieldList:
		for _, f := range n.List {
			if i := containsT(f, t); i != nil {
				return i
			}
		}
		return nil
	case *ast.Field:
		for _, n := range n.Names {
			if i := containsT(n, t); i != nil {
				return i
			}
		}
		return nil
	case *ast.Ident:
		t := t.TypeOf(n)
		if t.String() != "*testing.T" {
			return nil
		}
		return n
	default:
		panic(n)
	}
}

func helper(n ast.Node, t *types.Info) bool {
	switch n := n.(type) {
	case *ast.SelectorExpr:
		return helper(n.Sel, t)
	case *ast.Ident:
		o := t.ObjectOf(n)
		f, ok := o.(*types.Func)
		if !ok {
			return false
		}
		name := f.FullName()
		return name == "(*testing.common).Helper"
	default:
		return false
	}
}
