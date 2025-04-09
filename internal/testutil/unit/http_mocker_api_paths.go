package unit

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type OpenapiSchema struct {
	Paths map[string]map[string]interface{} `yaml:"paths"`
}

func parseModel(apiSpecPath string) (OpenapiSchema, error) {
	data, err := os.ReadFile(apiSpecPath)
	if err != nil {
		return OpenapiSchema{}, err
	}

	var model OpenapiSchema
	err = yaml.Unmarshal(data, &model)
	if err != nil {
		return OpenapiSchema{}, err
	}

	return model, nil
}

func parseAPISpecPaths(apiSpecPath string) (map[string][]APISpecPath, error) {
	model, err := parseModel(apiSpecPath)
	if err != nil {
		return nil, err
	}
	paths := make(map[string][]APISpecPath)
	for path, pathDict := range model.Paths {
		for method := range pathDict {
			methodUpper := strings.ToUpper(method)
			paths[methodUpper] = append(paths[methodUpper], APISpecPath{Path: path})
		}
	}
	return paths, nil
}

// copied from tools/codegen/openapi/parser.go
const (
	atlasAdminAPISpecURL = "https://raw.githubusercontent.com/mongodb/atlas-sdk-go/main/openapi/atlas-api-transformed.yaml"
	specFileRelPath      = "tools/codegen/open-api-spec.yml"
)

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

var apiSpecPaths map[string][]APISpecPath

func ReadAPISpecPaths() map[string][]APISpecPath {
	return apiSpecPaths
}

func FileExist(fullPath string) bool {
	_, err := os.Stat(fullPath)
	if err == nil {
		return true
	}
	return !os.IsNotExist(err)
}

func RepoPath(relPath string) string {
	workDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("error getting working directory: %s", err))
	}
	workdDirParts := strings.Split(workDir, "/")
	workdDirParts[0] = "/" + workdDirParts[0]
	for i := range workdDirParts {
		parentCandidate := workdDirParts[:len(workdDirParts)-i]
		parentCandidate = append(parentCandidate, ".git")
		gitDir := path.Join(parentCandidate...)
		if FileExist(gitDir) {
			repoPath, _ := strings.CutSuffix(gitDir, ".git")
			return path.Join(repoPath, relPath)
		}
	}
	panic("could not find repo root")
}

func PackagePath(name string) string {
	return RepoPath(path.Join("internal/service", name))
}

func init() {
	InitializeAPISpecPaths()
}

func InitializeAPISpecPaths() {
	specPath := RepoPath(specFileRelPath)
	var err error
	if !FileExist(specPath) {
		err = DownloadOpenAPISpec(atlasAdminAPISpecURL, specPath)
		if err != nil {
			panic(fmt.Sprintf("error downloading OpenAPI spec: %s", err))
		}
	}
	apiSpecPaths, err = parseAPISpecPaths(specPath)
	if err != nil {
		panic(fmt.Sprintf("error parsing OpenAPI spec: %s", err))
	}
}
