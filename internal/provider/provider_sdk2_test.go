package provider_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/provider"
)

func TestSdkV2Provider(t *testing.T) {
	if err := provider.NewSdkV2Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
