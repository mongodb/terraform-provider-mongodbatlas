package streamprivatelinkendpoint

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"go.mongodb.org/atlas-sdk/v20241113003/admin"
)

func NewTFModel(ctx context.Context, projectID string, apiResp *admin.StreamsPrivateLinkConnection) (*TFModel, diag.Diagnostics) {
	result := &TFModel{
		Id:                  types.StringPointerValue(apiResp.Id),
		DnsDomain:           types.StringPointerValue(apiResp.DnsDomain),
		ProjectId:           types.StringPointerValue(&projectID),
		InterfaceEndpointId: types.StringPointerValue(apiResp.InterfaceEndpointId),
		Provider:            types.StringPointerValue(apiResp.Provider),
		Region:              types.StringPointerValue(apiResp.Region),
		ServiceEndpointId:   types.StringPointerValue(apiResp.ServiceEndpointId),
		State:               types.StringPointerValue(apiResp.State),
		Vendor:              types.StringPointerValue(apiResp.Vendor),
	}
	if apiResp.DnsSubDomain != nil {
		subdomain, diagn := types.ListValueFrom(ctx, types.StringType, apiResp.GetDnsSubDomain())
		if diagn.HasError() {
			return nil, diagn
		}
		result.DnsSubDomain = subdomain
	}

	return result, nil
}

func NewAtlasReq(ctx context.Context, plan *TFModel) (*admin.StreamsPrivateLinkConnection, diag.Diagnostics) {
	result := &admin.StreamsPrivateLinkConnection{
		DnsDomain:         plan.DnsDomain.ValueStringPointer(),
		DnsSubDomain:      &[]string{},
		Provider:          plan.Provider.ValueStringPointer(),
		Region:            plan.Region.ValueStringPointer(),
		ServiceEndpointId: plan.ServiceEndpointId.ValueStringPointer(),
		State:             plan.State.ValueStringPointer(),
		Vendor:            plan.Vendor.ValueStringPointer(),
	}

	if !plan.DnsSubDomain.IsNull() {
		var dnsSubdomains []string
		diags := plan.DnsSubDomain.ElementsAs(ctx, &dnsSubdomains, false)
		if diags.HasError() {
			return nil, diags
		}
		result.DnsSubDomain = &dnsSubdomains
	}
	return result, nil
}

func NewTFModelPluralDS(ctx context.Context, projectID string, sdkResults []admin.StreamsPrivateLinkConnection) (*TFModelDSP, diag.Diagnostics) {
	diags := &diag.Diagnostics{}
	tfModels := make([]TFModel, len(sdkResults))
	for i := range sdkResults {
		tfModel, diagsLocal := NewTFModel(ctx, projectID, &sdkResults[i])
		diags.Append(diagsLocal...)
		if tfModel != nil {
			tfModels[i] = *tfModel
		}
	}
	if diags.HasError() {
		return nil, *diags
	}
	return &TFModelDSP{
		ProjectId: types.StringValue(projectID),
		Results:   tfModels,
	}, *diags
}
