package advancedcluster_test

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestRemoveLabel(t *testing.T) {
	toRemove := matlas.Label{Key: "To Remove", Value: "To remove value"}

	expected := []matlas.Label{
		{Key: "Name", Value: "Test"},
		{Key: "Version", Value: "1.0"},
		{Key: "Type", Value: "testing"},
	}

	labels := []matlas.Label{
		{Key: "Name", Value: "Test"},
		{Key: "Version", Value: "1.0"},
		{Key: "To Remove", Value: "To remove value"},
		{Key: "Type", Value: "testing"},
	}

	got := advancedcluster.RemoveLabel(labels, toRemove)

	if diff := deep.Equal(expected, got); diff != nil {
		t.Fatalf("Bad removeLabel return \n got = %#v\nwant = %#v \ndiff = %#v", got, expected, diff)
	}
}
