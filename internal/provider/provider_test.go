package provider_test

import (
	"testing"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/go-test/deep"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/provider"
)

func TestSdkV2Provider(t *testing.T) {
	if err := provider.NewSdkV2Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

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

	got := config.RemoveLabel(labels, toRemove)

	if diff := deep.Equal(expected, got); diff != nil {
		t.Fatalf("Bad removeLabel return \n got = %#v\nwant = %#v \ndiff = %#v", got, expected, diff)
	}
}
