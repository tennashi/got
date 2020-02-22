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
	tools, err := ensureToolsFile(dataDir)
	if err != nil {
		return err
	}
	fset, node, err := parseToolsFile(tools)
	if err != nil {
		return err
	}

	astutil.AddNamedImport(fset, node, "_", importPath)
	return writeToolsFile(tools, fset, node)
}

func removeImport(dataDir, importPath string) error {
	tools, err := ensureToolsFile(dataDir)
	if err != nil {
		return err
	}
	fset, node, err := parseToolsFile(tools)
	if err != nil {
		return err
	}

	astutil.DeleteNamedImport(fset, node, "_", importPath)
	return writeToolsFile(tools, fset, node)
}

func getImports(dataDir string) ([]string, error) {
	tools, err := ensureToolsFile(dataDir)
	if err != nil {
		return nil, err
	}
	fset, node, err := parseToolsFile(tools)
	if err != nil {
		return nil, err
	}

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

func parseToolsFile(path string) (*token.FileSet, *ast.File, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.Mode(0))
	if err != nil {
		return nil, nil, err
	}
	return fset, node, nil
}

func ensureToolsFile(dataDir string) (string, error) {
	path := filepath.Join(dataDir, "tools.go")
	if _, err := os.Stat(path); err != nil {
		f, err := os.Create(path)
		if err != nil {
			return "", err
		}
		defer f.Close()
		f.WriteString("package main\n")
	}
	return path, nil
}

func writeToolsFile(path string, fset *token.FileSet, node *ast.File) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	return printer.Fprint(f, fset, node)
}
