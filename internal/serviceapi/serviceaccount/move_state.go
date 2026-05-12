package serviceaccount

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// apiResourceSourceTypeName is the resource type name of the generic resource
// we accept state from.
const apiResourceSourceTypeName = "mongodbatlas_api_resource"

// idPattern matches the api_resource state.id for an org-scoped service
// account: /api/atlas/v2/orgs/<orgId>/serviceAccounts/<clientId>[/]
var idPattern = regexp.MustCompile(`^/api/atlas/v2/orgs/([^/]+)/serviceAccounts/([^/]+)/?$`)

// parseAPIResourceID extracts org_id and client_id from a mongodbatlas_api_resource
// state.id targeting an org-scoped service account endpoint.
func parseAPIResourceID(id string) (orgID, clientID string, err error) {
	m := idPattern.FindStringSubmatch(id)
	if m == nil {
		return "", "", fmt.Errorf("source id %q does not match an org-scoped service account URL", id)
	}
	return m[1], m[2], nil
}

// Compile-time check that the typed resource implements ResourceWithMoveState.
var _ resource.ResourceWithMoveState = &rs{}

// MoveState declares that mongodbatlas_service_account can accept state moved
// from mongodbatlas_api_resource (the generic resource). Used with `moved {}`
// blocks to upgrade from early-adopter generic usage to the typed resource.
func (r *rs) MoveState(context.Context) []resource.StateMover {
	return []resource.StateMover{{StateMover: stateMover}}
}

var sourceAttrTypes = map[string]tftypes.Type{
	"id": tftypes.String,
}

// stateMover translates state from mongodbatlas_api_resource into the typed
// service_account schema. It only sets identity (org_id, client_id) — Read
// refills computed attributes and the new config supplies required attributes.
func stateMover(ctx context.Context, req resource.MoveStateRequest, resp *resource.MoveStateResponse) {
	if req.SourceTypeName != apiResourceSourceTypeName {
		return
	}
	if !strings.HasSuffix(req.SourceProviderAddress, "/mongodbatlas") {
		return
	}

	rawStateValue, err := req.SourceRawState.UnmarshalWithOpts(
		tftypes.Object{AttributeTypes: sourceAttrTypes},
		tfprotov6.UnmarshalOpts{ValueFromJSONOpts: tftypes.ValueFromJSONOpts{IgnoreUndefinedAttributes: true}},
	)
	if err != nil {
		resp.Diagnostics.AddError("Unable to unmarshal source state", err.Error())
		return
	}
	var stateObj map[string]tftypes.Value
	if err := rawStateValue.As(&stateObj); err != nil {
		resp.Diagnostics.AddError("Unable to parse source state", err.Error())
		return
	}
	var idPtr *string
	_ = stateObj["id"].As(&idPtr)
	if idPtr == nil || *idPtr == "" {
		resp.Diagnostics.AddError(
			"Missing id in source state",
			"mongodbatlas_api_resource state does not contain a non-empty id; cannot derive org_id and client_id.",
		)
		return
	}

	orgID, clientID, err := parseAPIResourceID(*idPtr)
	if err != nil {
		resp.Diagnostics.AddError("Unsupported source id for service_account migration", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.TargetState.SetAttribute(ctx, path.Root("org_id"), types.StringValue(orgID))...)
	resp.Diagnostics.Append(resp.TargetState.SetAttribute(ctx, path.Root("client_id"), types.StringValue(clientID))...)
}
