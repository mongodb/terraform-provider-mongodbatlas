package openapi

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ghodss/yaml"
)

func ParseAtlasAdminAPI(urlPath string) (*openapi3.T, error) {
	specFilePath := fmt.Sprintf("open-api-spec.yml")
	if err := downloadOpenAPISpec(urlPath, specFilePath); err != nil {
		return nil, err
	}

	openAPISpecFileYaml, err := os.ReadFile(specFilePath)
	if err != nil {
		return nil, err
	}

	specYaml, err := yaml.YAMLToJSON(openAPISpecFileYaml)
	if err != nil {
		fmt.Printf("err: %v\n", err)

		return nil, err
	}
	doc, err := openapi3.NewLoader().LoadFromData(specYaml)
	if err != nil {
		return nil, err
	}

	if doc == nil {
		fmt.Println("empty document found")
		os.Exit(1)
	}

	return doc, nil
}

func downloadOpenAPISpec(url, specFilePath string) (err error) {
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
