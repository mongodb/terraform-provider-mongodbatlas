package cloudbackupcollectionrestorejobcollection

import (
	"context"

	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
)

var _ autogen.DataSourceSchemaHook = (*pluralDS)(nil)

// DataSourceSchema drops list query filters. Autogen codegen emits them, but HandleDataSourceReadList
// replaces QueryParams with pageNum only (CLOUDP-433211). config.yml ignores cannot drop these names:
// they also match the singular data source path param and computed result fields.
func (d *pluralDS) DataSourceSchema(_ context.Context, baseSchema datasourceschema.Schema) datasourceschema.Schema {
	delete(baseSchema.Attributes, "source_namespace")
	delete(baseSchema.Attributes, "state")
	delete(baseSchema.Attributes, "target_namespace")
	return baseSchema
}
