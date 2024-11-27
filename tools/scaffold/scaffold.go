package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	ResourceCmd         = "resource"
	DataSourceCmd       = "data-source"
	PluralDataSourceCmd = "plural-data-source"
)

// struct which is applied to go template files
type ScaffoldParams struct {
	GenerationType    string
	NamePascalCase    string
	NameCamelCase     string
	NameSnakeCase     string
	NameLowerNoSpaces string
}

type FileGeneration struct {
	TemplatePath string
	OutputPath   string
}

func main() {
	nameCamelCase := os.Args[1]
	generationType := os.Args[2]

	params := ScaffoldParams{
		GenerationType:    generationType,
		NamePascalCase:    ToPascalCase(nameCamelCase),
		NameCamelCase:     nameCamelCase,
		NameSnakeCase:     ToSnakeCase(nameCamelCase),
		NameLowerNoSpaces: strings.ToLower(nameCamelCase),
	}

	files, err := filesToGenerate(&params)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if err := generateFileFromTemplate(file, &params); err != nil {
			panic(err)
		}
	}
}

func filesToGenerate(params *ScaffoldParams) ([]FileGeneration, error) {
	folderPath := fmt.Sprintf("internal/service/%s", params.NameLowerNoSpaces)

	switch params.GenerationType {
	case ResourceCmd:
		return []FileGeneration{
			{
				TemplatePath: "tools/scaffold/template/resource.tmpl",
				OutputPath:   fmt.Sprintf("%s/resource.go", folderPath),
			},
			{
				TemplatePath: "tools/scaffold/template/acc_test.tmpl",
				OutputPath:   fmt.Sprintf("%s/resource_test.go", folderPath),
			},
			{
				TemplatePath: "tools/scaffold/template/model.tmpl",
				OutputPath:   fmt.Sprintf("%s/model.go", folderPath),
			},
			{
				TemplatePath: "tools/scaffold/template/model_test.tmpl",
				OutputPath:   fmt.Sprintf("%s/model_test.go", folderPath),
			},
			{
				TemplatePath: "tools/scaffold/template/generator_config.tmpl",
				OutputPath:   fmt.Sprintf("%s/tfplugingen/generator_config.yml", folderPath),
			},
			{
				TemplatePath: "tools/scaffold/template/main_test.tmpl",
				OutputPath:   fmt.Sprintf("%s/main_test.go", folderPath),
			},
		}, nil
	case DataSourceCmd:
		return []FileGeneration{
			{
				TemplatePath: "tools/scaffold/template/datasource.tmpl",
				OutputPath:   fmt.Sprintf("%s/data_source.go", folderPath),
			},
			{
				TemplatePath: "tools/scaffold/template/acc_test.tmpl",
				OutputPath:   fmt.Sprintf("%s/resource_test.go", folderPath),
			},
			{
				TemplatePath: "tools/scaffold/template/model.tmpl",
				OutputPath:   fmt.Sprintf("%s/model.go", folderPath),
			},
			{
				TemplatePath: "tools/scaffold/template/model_test.tmpl",
				OutputPath:   fmt.Sprintf("%s/model_test.go", folderPath),
			},
			{
				TemplatePath: "tools/scaffold/template/generator_config.tmpl",
				OutputPath:   fmt.Sprintf("%s/tfplugingen/generator_config.yml", folderPath),
			},
			{
				TemplatePath: "tools/scaffold/template/main_test.tmpl",
				OutputPath:   fmt.Sprintf("%s/main_test.go", folderPath),
			},
		}, nil
	case PluralDataSourceCmd:
		return []FileGeneration{
			{
				TemplatePath: "tools/scaffold/template/pluraldatasource.tmpl",
				OutputPath:   fmt.Sprintf("%s/plural_data_source.go", folderPath),
			},
			{
				TemplatePath: "tools/scaffold/template/acc_test.tmpl",
				OutputPath:   fmt.Sprintf("%s/resource_test.go", folderPath),
			},
			{
				TemplatePath: "tools/scaffold/template/model.tmpl",
				OutputPath:   fmt.Sprintf("%s/model.go", folderPath),
			},
			{
				TemplatePath: "tools/scaffold/template/model_test.tmpl",
				OutputPath:   fmt.Sprintf("%s/model_test.go", folderPath),
			},
			{
				TemplatePath: "tools/scaffold/template/generator_config.tmpl",
				OutputPath:   fmt.Sprintf("%s/tfplugingen/generator_config.yml", folderPath),
			},
			{
				TemplatePath: "tools/scaffold/template/main_test.tmpl",
				OutputPath:   fmt.Sprintf("%s/main_test.go", folderPath),
			},
		}, nil
	default:
		return nil, errors.New("unknown generation type provided")
	}
}

func generateFileFromTemplate(generation FileGeneration, params *ScaffoldParams) error {
	tmpl, err := template.ParseFiles(generation.TemplatePath)
	if err != nil {
		return err
	}

	// ensure content of existing files is not overwritten
	if _, err := os.Stat(generation.OutputPath); err == nil {
		log.Printf("File already exists: %s", generation.OutputPath)
		return nil
	}
	file := createDirsAndFile(generation.OutputPath)

	if err := tmpl.Execute(file, params); err != nil {
		return err
	}
	file.Close()
	return nil
}

func createDirsAndFile(path string) *os.File {
	dirPath := filepath.Dir(path)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		log.Fatalf("Failed to create directories: %s", err)
	}

	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("Failed to create file: %s", err)
	}
	return file
}
