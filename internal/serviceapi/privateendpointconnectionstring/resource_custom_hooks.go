package privateendpointconnectionstring

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/customplanmodifier"
)

var _ autogen.ResourceSchemaHook = (*rs)(nil)

// ResourceSchema makes optimized_mode_enabled create-only. The PATCH body does not accept optimizedModeEnabled, but
// codegen keeps the attribute updatable because the spec default (false) makes it computed, and codegen never adds the
// create-only plan modifier to computed attributes. Without this hook, changing the value plans an update that Atlas
// ignores, and the apply fails with an inconsistent result. CreateOnlyBoolWithDefault puts the default in the plan on
// create, keeps the state value on later plans, and rejects changes with an error that points to the state value.
func (r *rs) ResourceSchema(ctx context.Context, baseSchema schema.Schema) schema.Schema {
	attr := baseSchema.Attributes["optimized_mode_enabled"].(schema.BoolAttribute)
	attr.PlanModifiers = []planmodifier.Bool{customplanmodifier.CreateOnlyBoolWithDefault(false)}
	attr.MarkdownDescription += " This value cannot be changed after creation. To use a different mode, create a new connection string, move your applications to it, and then delete the old one."
	baseSchema.Attributes["optimized_mode_enabled"] = attr
	return baseSchema
}
