package autogen

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	rawStateUnmarshalOpts = tfprotov6.UnmarshalOpts{ValueFromJSONOpts: tftypes.ValueFromJSONOpts{IgnoreUndefinedAttributes: true}}
)

// HandleMove migrates state from a supported source resource to a target resource.
// It extracts the specified id attributes from the source state and sets them in the target state.
func HandleMove(ctx context.Context, supportedSources, idAttrs []string, req resource.MoveStateRequest, resp *resource.MoveStateResponse) {
	if !slices.Contains(supportedSources, req.SourceTypeName) || !strings.HasSuffix(req.SourceProviderAddress, "/mongodbatlas") {
		return
	}

	attrTypes := make(map[string]tftypes.Type)
	for _, attrName := range idAttrs {
		attrTypes[attrName] = tftypes.String // Assuming id attrs are always a String
	}

	rawStateValue, err := req.SourceRawState.UnmarshalWithOpts(tftypes.Object{AttributeTypes: attrTypes}, rawStateUnmarshalOpts)
	if err != nil {
		resp.Diagnostics.AddError("Unable to unmarshal source state", err.Error())
		return
	}

	var stateObj map[string]tftypes.Value
	if err := rawStateValue.As(&stateObj); err != nil {
		resp.Diagnostics.AddError("Unable to parse source state", err.Error())
		return
	}

	for _, attrName := range idAttrs {
		var value *string
		_ = stateObj[attrName].As(&value)

		if value == nil || *value == "" {
			resp.Diagnostics.AddError("Unable to read attribute from state", fmt.Sprintf("Ensure the moved block references a valid source state containing %s.", attrName))
			return
		}

		resp.Diagnostics.Append(resp.TargetState.SetAttribute(ctx, path.Root(attrName), types.StringValue(*value))...)
	}
}
