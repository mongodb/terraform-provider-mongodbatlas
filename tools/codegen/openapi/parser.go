package openapi

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

func ParseAtlasAdminAPI(urlPath string) (*libopenapi.DocumentModel[v3.Document], error) {
	err := fetchOpenAPISpec(context.Background(), urlPath)
	if err != nil {
		return nil, err
	}

	atlasAPISpec, _ := os.ReadFile("open-api-spec.yml")
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

func fetchOpenAPISpec(ctx context.Context, urlPath string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", urlPath, http.NoBody)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download Atlas API spec: %s", resp.Status)
	}

	file, err := os.Create("open-api-spec.yml")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	println("Downloaded Atlas API spec")
	return nil
}
