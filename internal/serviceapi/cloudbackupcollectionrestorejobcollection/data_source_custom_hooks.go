package cloudbackupcollectionrestorejobcollection

import (
	"context"

	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
)

var _ autogen.DataSourceSchemaHook = (*pluralDS)(nil)

var pluralListFilterAttrs = []string{"source_namespace", "state", "target_namespace"}

// DataSourceSchema keeps list query filters on the schema (required for TFPluralDSModel conversion)
// but marks them computed-only. Autogen emits them as optional, yet HandleDataSourceReadList
// replaces QueryParams with pageNum only (CLOUDP-433211). config.yml ignores cannot drop these names:
// they also match the singular data source path param and computed result fields.
func (d *pluralDS) DataSourceSchema(_ context.Context, baseSchema datasourceschema.Schema) datasourceschema.Schema {
	for _, name := range pluralListFilterAttrs {
		attr, ok := baseSchema.Attributes[name].(datasourceschema.StringAttribute)
		if !ok {
			continue
		}
		attr.Optional = false
		attr.Required = false
		attr.Computed = true
		baseSchema.Attributes[name] = attr
	}
	return baseSchema
}
