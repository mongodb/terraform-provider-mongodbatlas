package cloudbackupsnapshot

import (
	"errors"
	"regexp"

	"go.mongodb.org/atlas-sdk/v20250312008/admin"
)

func SplitSnapshotImportID(id string) (*admin.GetClusterBackupSnapshotApiParams, error) {
	var re = regexp.MustCompile(`(?s)^([0-9a-fA-F]{24})-(.*)-([0-9a-fA-F]{24})$`)
	parts := re.FindStringSubmatch(id)
	if len(parts) != 4 {
		return nil, errors.New("import format error: to import a snapshot, use the format {project_id}-{cluster_name}-{snapshot_id}")
	}
	return &admin.GetClusterBackupSnapshotApiParams{
		GroupId:     parts[1],
		ClusterName: parts[2],
		SnapshotId:  parts[3],
	}, nil
}
