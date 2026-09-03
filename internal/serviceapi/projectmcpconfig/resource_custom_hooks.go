package projectmcpconfig

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var (
	_ autogen.ResourceSchemaHook   = (*rs)(nil)
	_ autogen.PreUpdateAPICallHook = (*rs)(nil)
)

// ResourceSchema rejects ip_access_list entries that set both ip_address and cidr_block
// (or neither) at plan time, matching the API's server-side validation.
func (r *rs) ResourceSchema(ctx context.Context, s schema.Schema) schema.Schema {
	ipAccessList, ok := s.Attributes["ip_access_list"].(schema.ListNestedAttribute)
	if !ok {
		return s
	}
	if ipAddress, ok := ipAccessList.NestedObject.Attributes["ip_address"].(schema.StringAttribute); ok {
		ipAddress.Validators = append(ipAddress.Validators, stringvalidator.ExactlyOneOf(
			path.MatchRelative().AtParent().AtName("cidr_block"),
		))
		ipAccessList.NestedObject.Attributes["ip_address"] = ipAddress
	}
	if cidrBlock, ok := ipAccessList.NestedObject.Attributes["cidr_block"].(schema.StringAttribute); ok {
		cidrBlock.Validators = append(cidrBlock.Validators, stringvalidator.ExactlyOneOf(
			path.MatchRelative().AtParent().AtName("ip_address"),
		))
		ipAccessList.NestedObject.Attributes["cidr_block"] = cidrBlock
	}
	ipAccessList.PlanModifiers = append(ipAccessList.PlanModifiers, ipAccessListPlanModifier{})
	s.Attributes["ip_access_list"] = ipAccessList
	return s
}

// PreUpdateAPICall strips the API-echoed cidr_block from update payloads. The API returns
// both ip_address and the derived cidr_block (/32 or /128) for single-address entries but
// rejects write requests containing both fields, and the update payload is built from state.
func (r *rs) PreUpdateAPICall(callParams config.APICallParams, bodyReq []byte) (modifiedParams config.APICallParams, modifiedBody []byte) {
	return callParams, stripEchoedIPAccessListFields(bodyReq)
}

func stripEchoedIPAccessListFields(bodyReq []byte) []byte {
	var body map[string]any
	if err := json.Unmarshal(bodyReq, &body); err != nil {
		return bodyReq
	}
	entries, ok := body["ipAccessList"].([]any)
	if !ok {
		return bodyReq
	}
	for _, e := range entries {
		entry, ok := e.(map[string]any)
		if !ok {
			continue
		}
		ip, hasIP := entry["ipAddress"].(string)
		cidr, hasCIDR := entry["cidrBlock"].(string)
		if hasIP && hasCIDR && ip != "" && (cidr == ip+"/32" || cidr == ip+"/128") {
			delete(entry, "cidrBlock")
		}
	}
	updated, err := json.Marshal(body)
	if err != nil {
		return bodyReq
	}
	return updated
}

// ipAccessListPlanModifier marks the server-computed ip_access_list sub-fields as unknown on
// updates. The API recreates all entries on every update, regenerating these fields, which
// otherwise fails Terraform's plan/result consistency check.
type ipAccessListPlanModifier struct{}

func (m ipAccessListPlanModifier) Description(context.Context) string {
	return "Marks computed ip_access_list sub-fields as unknown on updates, since the API recreates all entries on every update."
}

func (m ipAccessListPlanModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m ipAccessListPlanModifier) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() || req.StateValue.IsNull() {
		return
	}
	if req.Plan.Raw.Equal(req.State.Raw) {
		return
	}

	var planEntries []TFIpAccessListModel
	if diags := req.PlanValue.ElementsAs(ctx, &planEntries, false); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	changed := false
	for i := range planEntries {
		if !planEntries[i].CreatedAt.IsUnknown() {
			planEntries[i].CreatedAt = types.StringUnknown()
			changed = true
		}
		if !planEntries[i].LastUsedAddress.IsUnknown() {
			planEntries[i].LastUsedAddress = types.StringUnknown()
			changed = true
		}
		if !planEntries[i].LastUsedAt.IsUnknown() {
			planEntries[i].LastUsedAt = types.StringUnknown()
			changed = true
		}
		if !planEntries[i].RequestCount.IsUnknown() {
			planEntries[i].RequestCount = types.Int64Unknown()
			changed = true
		}
	}
	if !changed {
		return
	}

	newList, diags := types.ListValueFrom(ctx, req.PlanValue.ElementType(ctx), planEntries)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
	resp.PlanValue = newList
}
