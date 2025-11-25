package openapi

import (
	"fmt"
	"os"

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
