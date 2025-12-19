package main

import (
	"context"
	"fmt"
	"go/format"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	providerPath = "internal/provider/provider.go"
	serviceapi   = "internal/serviceapi"
)

func main() {
	content, err := os.ReadFile(providerPath)
	if err != nil {
		log.Fatalf("Read failed: %v", err)
	}
	s := string(content)

	entries, _ := os.ReadDir(serviceapi)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		pkg := entry.Name()
		path := filepath.Join(serviceapi, pkg)

		importLine := fmt.Sprintf("\"github.com/mongodb/terraform-provider-mongodbatlas/%s/%s\"", serviceapi, pkg)
		if !strings.Contains(s, importLine) {
			s = strings.Replace(s, "import (", "import (\n\t"+importLine, 1)
		}

		if fileExists(path, "resource.go") {
			s = inject(s, "project.Resource,", pkg+".Resource,")
		}
		if fileExists(path, "data_source.go") {
			s = inject(s, "project.DataSource,", pkg+".DataSource,")
		}
		if fileExists(path, "plural_data_source.go") {
			s = inject(s, "project.DataSource,", pkg+".PluralDataSource,")
		}
	}

	res, err := format.Source([]byte(s))
	if err != nil {
		log.Fatalf("Formatting failed: %v", err)
	}

	if err := os.WriteFile(providerPath, res, 0o600); err != nil {
		log.Fatalf("Write failed: %v", err)
	}

	log.Println("Syncing dependencies with go mod tidy...")
	if err := tidy(); err != nil {
		log.Printf("Warning: go mod tidy failed: %v", err)
	}

	log.Println("Successfully updated provider.go and dependencies.")
}

func inject(content, anchor, newLine string) string {
	if strings.Contains(content, newLine) {
		return content
	}
	return strings.Replace(content, anchor, newLine+"\n\t\t"+anchor, 1)
}

func fileExists(path, name string) bool {
	_, err := os.Stat(filepath.Join(path, name))
	return err == nil
}

func tidy() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
