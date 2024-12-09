package unit

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
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

func replaceVars(text string, vars map[string]string) string {
	for key, value := range vars {
		text = strings.ReplaceAll(text, fmt.Sprintf("{%s}", key), value)
	}
	return text
}

var versionDatePattern = regexp.MustCompile(`(\d{4}-\d{2}-\d{2})`)

func ExtractVersion(contentType string) (string, error) {
	match := versionDatePattern.FindStringSubmatch(contentType)
	if len(match) > 1 {
		return match[1], nil
	}
	return "", fmt.Errorf("could not extract version from %s header", contentType)
}

func ExtractVersionRequestResponse(headerValueRequest, headerValueResponse string) string {
	found := versionDatePattern.FindString(headerValueRequest)
	if found != "" {
		return found
	}
	return versionDatePattern.FindString(headerValueResponse)
}

func extractAndNormalizePayload(body io.Reader) (originalPayload, normalizedPayload string, err error) {
	if body != nil {
		payloadBytes, err := io.ReadAll(body)
		if err != nil {
			return "", "", err
		}
		originalPayload = string(payloadBytes)
	}
	normalizedPayload, err = normalizePayload(originalPayload)
	if err != nil {
		return "", "", err
	}
	return originalPayload, normalizedPayload, nil
}

func normalizePayload(payload string) (string, error) {
	if payload == "" {
		return "", nil
	}
	var tempHolder any
	err := json.Unmarshal([]byte(payload), &tempHolder)
	if err != nil {
		return "", err
	}
	sb := strings.Builder{}
	encoder := json.NewEncoder(&sb)
	encoder.SetIndent("", " ")
	err = encoder.Encode(tempHolder)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(sb.String(), "\n"), nil
}
