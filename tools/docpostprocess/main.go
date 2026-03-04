package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	fileFlag := flag.String("file", "", "Single markdown file to process")
	dirFlag := flag.String("dir", "", "Directory to recursively process all .md files")
	flag.Parse()

	if *fileFlag == "" && *dirFlag == "" {
		log.Fatal("either --file or --dir must be specified")
	}
	if *fileFlag != "" && *dirFlag != "" {
		log.Fatal("--file and --dir are mutually exclusive")
	}

	var files []string
	if *fileFlag != "" {
		files = []string{*fileFlag}
	} else {
		var err error
		files, err = collectMarkdownFiles(*dirFlag)
		if err != nil {
			log.Fatalf("failed to collect markdown files: %v", err)
		}
	}

	processed := 0
	for _, f := range files {
		changed, err := processFile(f)
		if err != nil {
			log.Fatalf("failed to process %s: %v", f, err)
		}
		if changed {
			processed++
			fmt.Printf("Post-processed: %s\n", f)
		}
	}
	if processed == 0 {
		fmt.Println("No files required polymorphic restructuring.")
	}
}

func collectMarkdownFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".md") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func processFile(path string) (bool, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return false, fmt.Errorf("reading file: %w", err)
	}

	original := string(content)
	result := RestructurePolymorphicDocs(original)

	if result == original {
		return false, nil
	}

	if err := os.WriteFile(path, []byte(result), 0o600); err != nil {
		return false, fmt.Errorf("writing file: %w", err)
	}
	return true, nil
}
