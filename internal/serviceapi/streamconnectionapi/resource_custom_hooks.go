package streamconnectionapi

import (
	"fmt"
	"regexp"
)

func (r *rs) PreImport(id string) (string, error) {
	normalizedID, err := parseLegacyImportID(id)
	if err == nil {
		return normalizedID, nil
	}

	return "", fmt.Errorf("use the format {workspace_name}-{project_id}-{connection_name}")
}

func parseLegacyImportID(id string) (string, error) {
	re := regexp.MustCompile(`^(.*)-([0-9a-fA-F]{24})-(.*)$`)
	m := re.FindStringSubmatch(id)
	if len(m) != 4 || m[1] == "" || m[3] == "" {
		return "", fmt.Errorf("use the format {workspace_name}-{project_id}-{connection_name}")
	}

	workspaceName := m[1]
	projectID := m[2]
	connectionName := m[3]
	return fmt.Sprintf("%s/%s/%s", projectID, workspaceName, connectionName), nil
}
