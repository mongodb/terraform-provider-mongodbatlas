package cloudbackupsnapshot

import (
	"errors"
	"regexp"

	"go.mongodb.org/atlas-sdk/v20250312028/admin"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

// SplitSnapshotImportID accepts the new "/"-separated format first and falls back to the
// legacy "-"-separated format. legacy is true only when the legacy format parsed successfully.
func SplitSnapshotImportID(id string) (params *admin.GetClusterBackupSnapshotApiParams, legacy bool, err error) {
	if ok, parts := conversion.ImportSplit(id, 3); ok {
		return &admin.GetClusterBackupSnapshotApiParams{
			GroupId:     parts[0],
			ClusterName: parts[1],
			SnapshotId:  parts[2],
		}, false, nil
	}

	legacy = true
	var re = regexp.MustCompile(`(?s)^([0-9a-fA-F]{24})-(.*)-([0-9a-fA-F]{24})$`)
	parts := re.FindStringSubmatch(id)
	if len(parts) != 4 {
		legacy = false
		err = errors.New("import format error: to import a snapshot, use the format {project_id}/{cluster_name}/{snapshot_id} or the legacy format {project_id}-{cluster_name}-{snapshot_id}")
		return
	}
	params = &admin.GetClusterBackupSnapshotApiParams{
		GroupId:     parts[1],
		ClusterName: parts[2],
		SnapshotId:  parts[3],
	}
	return
}
