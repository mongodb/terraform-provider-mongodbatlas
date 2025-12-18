package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/ast/astutil"
)

const (
	serviceAPIDir    = "internal/serviceapi"
	providerFilePath = "internal/provider/provider.go"
	importPrefix     = "github.com/mongodb/terraform-provider-mongodbatlas/internal/serviceapi/"
)

type AutogenComponents map[string][]string

const (
	ComponentResource           = "resource"
	ComponentSingularDatasource = "singularDatasource"
	ComponentPluralDatasource   = "pluralDatasource"
)

func main() {
	components, err := discoverAutogenComponents(serviceAPIDir)
	if err != nil {
		log.Fatalf("[ERROR] Discovery failed: %v", err)
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, providerFilePath, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("[ERROR] Failed to parse AST: %v", err)
	}

	// Move file operations to a separate function to ensure defer f.Close() runs
	if err := runTransformAndSave(fset, node, components); err != nil {
		log.Fatalf("[ERROR] %v", err)
	}

	log.Println("[INFO] Provider updated successfully.")
}

func runTransformAndSave(fset *token.FileSet, node *ast.File, components AutogenComponents) error {
	// 1. Update Imports
	for pkg := range components {
		astutil.AddImport(fset, node, importPrefix+pkg)
	}

	// 2. Inject Registrations (Targeting specific variable names)
	updateRegistrations(node, "Resources", "resources", "Resource", ComponentResource, components)
	updateRegistrations(node, "DataSources", "dataSources", "DataSource", ComponentSingularDatasource, components)
	updateRegistrations(node, "DataSources", "dataSources", "PluralDataSource", ComponentPluralDatasource, components)

	// 3. Save Changes using modern octal literal 0o644
	f, err := os.OpenFile(providerFilePath, os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open provider file: %w", err)
	}
	defer f.Close()

	if err := format.Node(f, fset, node); err != nil {
		return fmt.Errorf("failed to format and write: %w", err)
	}
	return nil
}

func updateRegistrations(node *ast.File, funcName, varName, suffix, componentType string, components AutogenComponents) {
	for _, decl := range node.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name.Name != funcName {
			continue
		}

		ast.Inspect(fn.Body, func(n ast.Node) bool {
			assign, ok := n.(*ast.AssignStmt)
			if !ok || len(assign.Lhs) == 0 {
				return true
			}

			// Check if the variable name on the Left Hand Side matches our target (e.g., "dataSources")
			ident, ok := assign.Lhs[0].(*ast.Ident)
			if !ok || ident.Name != varName {
				return true
			}

			for _, rhs := range assign.Rhs {
				if cl, ok := rhs.(*ast.CompositeLit); ok {
					for pkg, types := range components {
						if hasComponent(types, componentType) {
							cl.Elts = append(cl.Elts, &ast.SelectorExpr{
								X:   ast.NewIdent(pkg),
								Sel: ast.NewIdent(suffix),
							})
						}
					}
				}
			}
			return false // Stop once we've handled the target variable
		})
	}
}

func discoverAutogenComponents(dir string) (AutogenComponents, error) {
	components := make(AutogenComponents)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		pkg := entry.Name()
		path := filepath.Join(dir, pkg)

		var list []string
		if _, err := os.Stat(filepath.Join(path, "resource.go")); err == nil {
			list = append(list, ComponentResource)
		}
		if _, err := os.Stat(filepath.Join(path, "data_source.go")); err == nil {
			list = append(list, ComponentSingularDatasource)
		}
		if _, err := os.Stat(filepath.Join(path, "plural_data_source.go")); err == nil {
			list = append(list, ComponentPluralDatasource)
		}

		if len(list) > 0 {
			components[pkg] = list
		}
	}
	return components, nil
}

func hasComponent(list []string, target string) bool {
	for _, s := range list {
		if s == target {
			return true
		}
	}
	return false
}
