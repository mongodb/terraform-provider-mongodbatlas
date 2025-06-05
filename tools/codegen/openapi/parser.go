package openapi

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/config"
	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

const (
	specFileExtension = ".yaml"
)

func ParseAtlasAdminAPI(filePath string) (*libopenapi.DocumentModel[v3.Document], error) {
	atlasAPISpec, _ := os.ReadFile(filePath)
	document, err := libopenapi.NewDocument(atlasAPISpec)
	if err != nil {
		return nil, fmt.Errorf("cannot create new document: %e", err)
	}
	docModel, errors := document.BuildV3Model()
	if len(errors) > 0 {
		for i := range errors {
			fmt.Printf("error: %e\n", errors[i])
		}
		return nil, fmt.Errorf("cannot create v3 model from document: %d errors reported", len(errors))
	}

	return docModel, nil
}

func DownloadOpenAPISpecs(specs []config.APISpec, specDirPath string) error {
	if err := os.MkdirAll(specDirPath, 0o755); err != nil {
		return fmt.Errorf("failed to create spec directory: %w", err)
	}
	for _, spec := range specs {
		specFilePath := SpecFilePath(specDirPath, spec.Name)
		if err := DownloadOpenAPISpec(spec.URL, specFilePath); err != nil {
			return err
		}
	}
	return nil
}

func ParseAPISpecs(specDirPath string, names []string) map[string]*libopenapi.DocumentModel[v3.Document] {
	apiSpecs := ReadAPISpecs(specDirPath)
	apiSpecsParsed := make(map[string]*libopenapi.DocumentModel[v3.Document])
	for _, specName := range names {
		specFilePath, ok := apiSpecs[specName]
		if !ok {
			log.Fatalf("API spec file for %s not found in directory %s", specName, specDirPath)
		}
		specParsed, err := ParseAtlasAdminAPI(specFilePath)
		if err != nil {
			log.Fatalf("an error occurred when parsing Atlas Admin API spec @ %s: %v", specFilePath, err)
		}
		apiSpecsParsed[specName] = specParsed
	}
	return apiSpecsParsed
}

func ReadAPISpecs(specDirPath string) map[string]string {
	specs := make(map[string]string)
	files, err := os.ReadDir(specDirPath)
	if err != nil {
		return specs
	}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), specFileExtension) {
			specName := strings.TrimSuffix(file.Name(), specFileExtension)
			specs[specName] = SpecFilePath(specDirPath, specName)
		}
	}
	return specs
}

func DownloadOpenAPISpec(url, specFilePath string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return err
	}

	client := http.Client{}
	res, getErr := client.Do(req)
	if getErr != nil {
		return getErr
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return readErr
	}

	err = os.WriteFile(specFilePath, body, 0o600)
	return err
}

func SpecFilePath(specDirPath, specName string) string {
	return fmt.Sprintf("%s/%s", specDirPath, specName) + specFileExtension
}
