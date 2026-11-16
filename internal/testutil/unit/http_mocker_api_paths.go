package unit

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

type OpenapiSchema struct {
	Paths map[string]map[string]any `yaml:"paths"`
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

	downloadTimeout    = 30 * time.Second
	downloadAttempts   = 3
	downloadRetryDelay = time.Second
)

func DownloadOpenAPISpec(url, specFilePath string) error {
	var lastErr error
	for attempt := range downloadAttempts {
		if attempt > 0 {
			time.Sleep(downloadRetryDelay)
		}
		lastErr = downloadOpenAPISpecAttempt(url, specFilePath)
		if lastErr == nil {
			return nil
		}
	}
	return fmt.Errorf("download failed after %d attempts: %w", downloadAttempts, lastErr)
}

func downloadOpenAPISpecAttempt(url, specFilePath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), downloadTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return err
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return os.WriteFile(specFilePath, body, 0o600)
}

var (
	apiSpecPaths     map[string][]APISpecPath
	apiSpecPathsOnce sync.Once
)

// ReadAPISpecPaths returns the parsed OpenAPI spec paths, loading them on first use.
func ReadAPISpecPaths() map[string][]APISpecPath {
	InitializeAPISpecPaths()
	return apiSpecPaths
}

// InitializeAPISpecPaths loads the OpenAPI spec paths once per test binary, downloading the spec when it is not already present.
// It runs on first use so test binaries that never read the paths, such as replay-only or migration runs, do not download the spec.
func InitializeAPISpecPaths() {
	apiSpecPathsOnce.Do(func() {
		specPath := RepoPath(specFileRelPath)
		if !FileExist(specPath) {
			if err := DownloadOpenAPISpec(atlasAdminAPISpecURL, specPath); err != nil {
				panic(fmt.Sprintf("error downloading OpenAPI spec: %s", err))
			}
		}
		var err error
		apiSpecPaths, err = parseAPISpecPaths(specPath)
		if err != nil {
			panic(fmt.Sprintf("error parsing OpenAPI spec: %s", err))
		}
	})
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
