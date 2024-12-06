package unit

import (
	"fmt"
	"strings"
)

type APISpecPath struct {
	Path string
}

func (a *APISpecPath) Variables(path string) map[string]string {
	variables := make(map[string]string)
	expectedParts := strings.Split(a.Path, "/")
	actualParts := strings.Split(path, "/")
	for i, part := range expectedParts {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			variables[part[1:len(part)-1]] = actualParts[i]
		}
	}
	return variables
}

func (a *APISpecPath) Match(path string) bool {
	expectedParts := strings.Split(a.Path, "/")
	actualParts := strings.Split(path, "/")
	if len(expectedParts) != len(actualParts) {
		return false
	}
	for i, expected := range expectedParts {
		actual := actualParts[i]
		if expected == actual {
			continue
		}
		if strings.HasPrefix(expected, "{") && strings.HasSuffix(expected, "}") {
			continue
		}
		return false
	}
	return true
}

func FindNormalizedPath(path string, apiSpecPaths *[]APISpecPath) (APISpecPath, error) {
	if strings.Contains(path, "?") {
		path = strings.Split(path, "?")[0]
	}
	path = strings.TrimRight(path, "/")
	for _, apiSpecPath := range *apiSpecPaths {
		if apiSpecPath.Match(path) {
			return apiSpecPath, nil
		}
	}
	return APISpecPath{}, fmt.Errorf("could not find path: %s", path)
}
