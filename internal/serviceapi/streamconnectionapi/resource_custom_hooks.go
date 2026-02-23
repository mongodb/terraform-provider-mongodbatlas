package streamconnectionapi

import (
	"fmt"
	"regexp"
	"strings"
)

func (r *rs) PreImport(id string) (string, error) {
	if strings.Contains(id, "/") {
		return id, nil
	}

	normalizedID, err := parseLegacyImportID(id)
	if err == nil {
		return normalizedID, nil
	}

	return "", fmt.Errorf("use one of the formats: {project_id}/{workspace_name}/{connection_name} or {workspace_name}-{project_id}-{connection_name}")
}

func parseLegacyImportID(id string) (string, error) {
	re := regexp.MustCompile(`^(.*)-([0-9a-fA-F]{24})-(.*)$`)
	m := re.FindStringSubmatch(id)
	if len(m) != 4 || m[1] == "" || m[3] == "" {
		return "", fmt.Errorf("invalid legacy import format")
	}

	workspaceName := m[1]
	projectID := m[2]
	connectionName := m[3]
	// Normalize to default format: {project_id}/{workspace_name}/{connection_name}
	return fmt.Sprintf("%s/%s/%s", projectID, workspaceName, connectionName), nil
}
