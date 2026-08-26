package projectmcpconfig

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
	if cidrBlock, ok := ipAccessList.NestedObject.Attributes["cidr_block"].(schema.StringAttribute); ok {
		cidrBlock.Validators = append(cidrBlock.Validators, stringvalidator.ExactlyOneOf(
			path.MatchRelative().AtParent().AtName("ip_address"),
		))
		ipAccessList.NestedObject.Attributes["cidr_block"] = cidrBlock
	}
	s.Attributes["ip_access_list"] = ipAccessList
	return s
}

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
		if hasIP && hasCIDR && ip != "" && cidr == ip+"/32" {
			delete(entry, "cidrBlock")
		}
	}
	updated, err := json.Marshal(body)
	if err != nil {
		return bodyReq
	}
	return updated
}
