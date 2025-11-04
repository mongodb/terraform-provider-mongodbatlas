package openapi

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

func ParseAtlasAdminAPI(filePath string) (*libopenapi.DocumentModel[v3.Document], error) {
	atlasAPISpec, _ := os.ReadFile(filePath)
	document, err := libopenapi.NewDocument(atlasAPISpec)
	if err != nil {
		return nil, fmt.Errorf("cannot create new document: %e", err)
	}
	docModel, err := document.BuildV3Model()
	if err != nil {
		return nil, fmt.Errorf("cannot create v3 model from document: %w", err)
	}

	return docModel, nil
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
