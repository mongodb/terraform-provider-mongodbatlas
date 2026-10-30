package flexcluster

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/customplanmodifier"
)

var _ autogen.ResourceSchemaHook = (*rs)(nil)

// ResourceSchema moves the create-only check of provider_settings from the object to its user-set
// attributes, as the handwritten resource does. Codegen puts CreateOnly on the whole object because
// provider_settings is not in the PATCH body, but the object also holds computed attributes
// (disk_size_gb, provider_name) that are still unknown when the object's plan modifiers run, so the
// object-level check fails on every update (e.g. changing only tags).
func (r *rs) ResourceSchema(ctx context.Context, baseSchema schema.Schema) schema.Schema {
	providerSettings, ok := baseSchema.Attributes["provider_settings"].(schema.SingleNestedAttribute)
	if !ok {
		return baseSchema
	}
	providerSettings.PlanModifiers = nil
	for _, name := range []string{"backing_provider_name", "region_name"} {
		if attr, ok := providerSettings.Attributes[name].(schema.StringAttribute); ok {
			attr.PlanModifiers = append(attr.PlanModifiers, customplanmodifier.CreateOnly())
			providerSettings.Attributes[name] = attr
		}
	}
	baseSchema.Attributes["provider_settings"] = providerSettings
	return baseSchema
}
