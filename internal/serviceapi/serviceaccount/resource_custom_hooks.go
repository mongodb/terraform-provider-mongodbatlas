package serviceaccount

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
)

// PathParamsProvider is an interface for models that can provide path parameters.
// This allows both TFModel and TFModelWithID to be used interchangeably.
type PathParamsProvider interface {
	PathParams() map[string]string
}

// PathParams returns the path parameters needed for API calls.
func (m *TFModel) PathParams() map[string]string {
	return map[string]string{
		"orgId":    m.OrgId.ValueString(),
		"clientId": m.ClientId.ValueString(),
	}
}

// TFModelWithID embeds TFModel and adds the computed id attribute.
// This is used for state operations to include the id field that isn't in the auto-generated model.
type TFModelWithID struct {
	TFModel
	Id types.String `tfsdk:"id"`
}

// Compile-time interface checks
var _ autogen.SchemaExtensionHook = (*rs)(nil)
var _ autogen.StateModelHook = (*rs)(nil)
var _ PathParamsProvider = (*TFModel)(nil)
var _ PathParamsProvider = (*TFModelWithID)(nil)

// ExtendSchema adds the computed 'id' attribute to the auto-generated schema.
func (r *rs) ExtendSchema(ctx context.Context, baseSchema schema.Schema) schema.Schema {
	baseSchema.Attributes["id"] = schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "Unique identifier for this resource, composed of org_id and client_id.",
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
	return baseSchema
}

// NewStateModel returns a pointer to a new TFModelWithID instance for reading state.
func (r *rs) NewStateModel() any {
	return &TFModelWithID{}
}

// PrepareForState prepares the extended model before setting state.
// It computes the id from org_id and client_id.
func (r *rs) PrepareForState(model any) any {
	// If already an extended model, just update the Id field
	if m, ok := model.(*TFModelWithID); ok {
		m.Id = types.StringValue(fmt.Sprintf("%s-%s", m.OrgId.ValueString(), m.ClientId.ValueString()))
		return m
	}
	// Otherwise wrap the base model
	m := model.(*TFModel)
	return &TFModelWithID{
		TFModel: *m,
		Id:      types.StringValue(fmt.Sprintf("%s-%s", m.OrgId.ValueString(), m.ClientId.ValueString())),
	}
}
