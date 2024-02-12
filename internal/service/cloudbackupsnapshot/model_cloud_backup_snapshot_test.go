package cloudbackupsnapshot_test

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/cloudbackupsnapshot"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestSplitSnapshotImportID(t *testing.T) {
	got, err := cloudbackupsnapshot.SplitSnapshotImportID("5cf5a45a9ccf6400e60981b6-projectname-environment-mongo-global-cluster-5cf5a45a9ccf6400e60981b7")
	if err != nil {
		t.Errorf("splitSnapshotImportID returned error(%s), expected error=nil", err)
	}

	expected := &matlas.SnapshotReqPathParameters{
		GroupID:     "5cf5a45a9ccf6400e60981b6",
		ClusterName: "projectname-environment-mongo-global-cluster",
		SnapshotID:  "5cf5a45a9ccf6400e60981b7",
	}

	if diff := deep.Equal(expected, got); diff != nil {
		t.Errorf("Bad splitSnapshotImportID return \n got = %#v\nwant = %#v \ndiff = %#v", expected, *got, diff)
	}

	if _, err := cloudbackupsnapshot.SplitSnapshotImportID("5cf5a45a9ccf6400e60981b6projectname-environment-mongo-global-cluster5cf5a45a9ccf6400e60981b7"); err == nil {
		t.Error("splitSnapshotImportID expected to have error")
	}
}
