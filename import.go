package got

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strconv"

	"golang.org/x/tools/go/ast/astutil"
)

func appendImport(dataDir, importPath string) error {
	f, fset, node, err := parseToolsFile(dataDir)
	if err != nil {
		return nil
	}
	defer f.Close()

	astutil.AddNamedImport(fset, node, "_", importPath)
	return printer.Fprint(f, fset, node)
}

func deleteImport(dataDir, importPath string) error {
	f, fset, node, err := parseToolsFile(dataDir)
	if err != nil {
		return nil
	}
	defer f.Close()

	astutil.DeleteNamedImport(fset, node, "_", importPath)
	return printer.Fprint(f, fset, node)
}

func getImports(dataDir string) ([]string, error) {
	f, fset, node, err := parseToolsFile(dataDir)
	if err != nil {
		return nil, nil
	}
	defer f.Close()

	iSpecs := astutil.Imports(fset, node)
	imports := make([]string, len(iSpecs[0]))
	for i, iSpec := range iSpecs[0] {
		var err error
		imports[i], err = strconv.Unquote(iSpec.Path.Value)
		if err != nil {
			return nil, err
		}
	}
	return imports, nil
}

func parseToolsFile(dataDir string) (*os.File, *token.FileSet, *ast.File, error) {
	f, toolsPath, err := ensureToolsFile(dataDir)
	if err != nil {
		return nil, nil, nil, err
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, toolsPath, nil, parser.Mode(0))
	if err != nil {
		return nil, nil, nil, err
	}
	return f, fset, node, nil
}

func ensureToolsFile(dataDir string) (*os.File, string, error) {
	path := filepath.Join(dataDir, "tools.go")
	if _, err := os.Stat(path); err != nil {
		f, err := os.Create(path)
		if err != nil {
			return nil, "", err
		}
		f.Write([]byte("package main\n"))
		return f, path, nil
	}
	f, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return nil, "", err
	}
	return f, path, nil
}
