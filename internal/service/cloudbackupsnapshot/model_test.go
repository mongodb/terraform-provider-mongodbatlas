package cloudbackupsnapshot_test

import (
	"reflect"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/cloudbackupsnapshot"
	"go.mongodb.org/atlas-sdk/v20250312028/admin"
)

func TestSplitSnapshotImportID(t *testing.T) {
	expected := &admin.GetClusterBackupSnapshotApiParams{
		GroupId:     "5cf5a45a9ccf6400e60981b6",
		ClusterName: "projectname-environment-mongo-global-cluster",
		SnapshotId:  "5cf5a45a9ccf6400e60981b7",
	}

	got, legacy, err := cloudbackupsnapshot.SplitSnapshotImportID("5cf5a45a9ccf6400e60981b6/projectname-environment-mongo-global-cluster/5cf5a45a9ccf6400e60981b7")
	if err != nil {
		t.Errorf("splitSnapshotImportID returned error(%s), expected error=nil", err)
	}
	if legacy {
		t.Error("splitSnapshotImportID returned legacy=true for the new import format, expected false")
	}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Bad splitSnapshotImportID return \n got = %#v\nwant = %#v", expected, *got)
	}

	gotLegacy, legacy, err := cloudbackupsnapshot.SplitSnapshotImportID("5cf5a45a9ccf6400e60981b6-projectname-environment-mongo-global-cluster-5cf5a45a9ccf6400e60981b7")
	if err != nil {
		t.Errorf("splitSnapshotImportID returned error(%s), expected error=nil", err)
	}
	if !legacy {
		t.Error("splitSnapshotImportID returned legacy=false for the legacy import format, expected true")
	}
	if !reflect.DeepEqual(expected, gotLegacy) {
		t.Errorf("Bad splitSnapshotImportID return \n got = %#v\nwant = %#v", expected, *gotLegacy)
	}

	if _, _, err := cloudbackupsnapshot.SplitSnapshotImportID("5cf5a45a9ccf6400e60981b6projectname-environment-mongo-global-cluster5cf5a45a9ccf6400e60981b7"); err == nil {
		t.Error("splitSnapshotImportID expected to have error")
	}
}
