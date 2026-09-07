package mcpconfig

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var (
	_ autogen.ResourceSchemaHook   = (*rs)(nil)
	_ autogen.PreUpdateAPICallHook = (*rs)(nil)
)

// ResourceSchema adds validators to reject ip_access_list entries that set both ip_address and cidr_block.
func (r *rs) ResourceSchema(ctx context.Context, s schema.Schema) schema.Schema {
	ipAccessList, ok := s.Attributes["ip_access_list"].(schema.SetNestedAttribute)
	if !ok {
		return s
	}
	if ipAddress, ok := ipAccessList.NestedObject.Attributes["ip_address"].(schema.StringAttribute); ok {
		ipAddress.Validators = append(ipAddress.Validators, stringvalidator.ExactlyOneOf(
			path.MatchRelative().AtParent().AtName("cidr_block"),
		))
		ipAccessList.NestedObject.Attributes["ip_address"] = ipAddress
	}
	s.Attributes["ip_access_list"] = ipAccessList
	return s
}

// PreUpdateAPICall
// The API returns both cidr_block and ip_address when they are equivalent, regardless of which one was sent, so both may be present in the state.
// This is not a problem for an entry being modified since the access_list attribute is a set, so a modification to an entry is planned as a set removal + addition.
// However, modifications to other attributes would plan both ip_address + cidr_block to be sent in the PATCH for all entries that have them set in the state.
// So here we drop the cidr_block if ip_address is already set. This is a safe assumption since the only way for both to get to this point is because
// the entry comes from the state and was not modified in the config, and sending ip_address or cidr_block is equivalent.
func (r *rs) PreUpdateAPICall(callParams config.APICallParams, bodyReq []byte) (modifiedParams config.APICallParams, modifiedBody []byte) {
	var body map[string]any
	if err := json.Unmarshal(bodyReq, &body); err != nil {
		return callParams, bodyReq
	}
	entries, ok := body["ipAccessList"].([]any)
	if !ok {
		return callParams, bodyReq
	}
	for _, e := range entries {
		entry, ok := e.(map[string]any)
		if !ok {
			continue
		}
		ip, hasIP := entry["ipAddress"].(string)
		_, hasCIDR := entry["cidrBlock"].(string)
		if hasIP && hasCIDR && ip != "" {
			delete(entry, "cidrBlock")
		}
	}
	updated, err := json.Marshal(body)
	if err != nil {
		return callParams, bodyReq
	}
	return callParams, updated
}
