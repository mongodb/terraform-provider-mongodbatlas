package cloudbackupsnapshot

import (
	"errors"
	"regexp"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func SplitSnapshotImportID(id string) (*matlas.SnapshotReqPathParameters, error) {
	var re = regexp.MustCompile(`(?s)^([0-9a-fA-F]{24})-(.*)-([0-9a-fA-F]{24})$`)

	parts := re.FindStringSubmatch(id)

	if len(parts) != 4 {
		return nil, errors.New("import format error: to import a snapshot, use the format {project_id}-{cluster_name}-{snapshot_id}")
	}

	return &matlas.SnapshotReqPathParameters{
		GroupID:     parts[1],
		ClusterName: parts[2],
		SnapshotID:  parts[3],
	}, nil
}
