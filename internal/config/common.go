package config

import (
	"errors"
	"regexp"
)

func SplitDatabaseUserImportID(id string) (projectID, username, authDatabaseName string, err error) {
	var re = regexp.MustCompile(`(?s)^([0-9a-fA-F]{24})-(.*)-([$a-z]{1,15})$`)
	parts := re.FindStringSubmatch(id)

	if len(parts) != 4 {
		err = errors.New("import format error: to import a Database User, use the format {project_id}-{username}-{auth_database_name}")
		return
	}

	projectID = parts[1]
	username = parts[2]
	authDatabaseName = parts[3]

	return
}
