package main

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"

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

	// 1. Update Imports
	for pkg := range components {
		astutil.AddImport(fset, node, importPrefix+pkg)
	}

	// 2. Inject Registrations
	// Maps the function name in provider.go to the file-suffix/component-type logic
	updateRegistrations(node, components, "Resources", "Resource", ComponentResource)
	updateRegistrations(node, components, "DataSources", "DataSource", ComponentSingularDatasource)
	updateRegistrations(node, components, "DataSources", "PluralDataSource", ComponentPluralDatasource)

	// 3. Save Changes
	f, err := os.OpenFile(providerFilePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("[ERROR] Failed to open provider file: %v", err)
	}
	defer f.Close()

	if err := format.Node(f, fset, node); err != nil {
		log.Fatalf("[ERROR] Failed to format and write: %v", err)
	}

	log.Println("[INFO] Provider updated successfully.")
}

func updateRegistrations(node *ast.File, components AutogenComponents, funcName, suffix, componentType string) {
	for _, decl := range node.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name.Name != funcName {
			continue
		}

		// Search for the specific variable assignment: dataSources := []...{...}
		// or resources := []...{...}
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			assign, ok := n.(*ast.AssignStmt)
			if !ok || len(assign.Lhs) == 0 {
				return true
			}

			// Identify the variable name (e.g., "dataSources" or "resources")
			ident, ok := assign.Lhs[0].(*ast.Ident)
			if !ok {
				return true
			}

			// Targeting the primary slices specifically.
			// In Provider.go these are usually named "resources" and "dataSources"
			targetName := strings.ToLower(funcName[:1]) + funcName[1:]
			if ident.Name != targetName {
				return true // Skip things like "analyticsDataSources"
			}

			// Now find the CompositeLit within this specific assignment
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
			return false
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
		if fileExists(filepath.Join(path, "resource.go")) {
			list = append(list, ComponentResource)
		}
		if fileExists(filepath.Join(path, "data_source.go")) {
			list = append(list, ComponentSingularDatasource)
		}
		if fileExists(filepath.Join(path, "plural_data_source.go")) {
			list = append(list, ComponentPluralDatasource)
		}

		if len(list) > 0 {
			components[pkg] = list
		}
	}
	return components, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func hasComponent(list []string, target string) bool {
	for _, s := range list {
		if s == target {
			return true
		}
	}
	return false
}
