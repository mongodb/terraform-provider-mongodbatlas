package cloudbackupcollectionrestorejobcollection_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/serviceapi/cloudbackupcollectionrestorejobcollection"
)

func TestPluralDSSchema_listFiltersAreComputedOnly(t *testing.T) {
	ds := cloudbackupcollectionrestorejobcollection.PluralDataSource()
	var resp datasource.SchemaResponse
	ds.Schema(context.Background(), datasource.SchemaRequest{}, &resp)
	require.False(t, resp.Diagnostics.HasError())

	for _, name := range []string{"source_namespace", "state", "target_namespace"} {
		t.Run(name, func(t *testing.T) {
			attr, ok := resp.Schema.Attributes[name]
			require.True(t, ok, "schema must keep %s so TFPluralDSModel conversion succeeds", name)
			str, ok := attr.(dsschema.StringAttribute)
			require.True(t, ok)
			assert.True(t, str.Computed)
			assert.False(t, str.Optional)
			assert.False(t, str.Required)
		})
	}
}
