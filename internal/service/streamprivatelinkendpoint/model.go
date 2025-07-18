package streamprivatelinkendpoint

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"go.mongodb.org/atlas-sdk/v20250312005/admin"
)

const (
	VendorConfluent = "CONFLUENT"
	VendorMSK       = "MSK"
)

func NewTFModel(ctx context.Context, projectID string, apiResp *admin.StreamsPrivateLinkConnection) (*TFModel, diag.Diagnostics) {
	result := &TFModel{
		Id:                    types.StringPointerValue(apiResp.Id),
		DnsDomain:             types.StringPointerValue(apiResp.DnsDomain),
		ErrorMessage:          types.StringPointerValue(apiResp.ErrorMessage),
		ProjectId:             types.StringPointerValue(&projectID),
		InterfaceEndpointId:   types.StringPointerValue(apiResp.InterfaceEndpointId),
		InterfaceEndpointName: types.StringPointerValue(apiResp.InterfaceEndpointName),
		Provider:              types.StringValue(apiResp.Provider),
		ProviderAccountId:     types.StringPointerValue(apiResp.ProviderAccountId),
		Region:                types.StringPointerValue(apiResp.Region),
		ServiceEndpointId:     types.StringPointerValue(apiResp.ServiceEndpointId),
		State:                 types.StringPointerValue(apiResp.State),
		Vendor:                types.StringPointerValue(apiResp.Vendor),
		Arn:                   types.StringPointerValue(apiResp.Arn),
	}

	subdomain, diags := types.ListValueFrom(ctx, types.StringType, apiResp.GetDnsSubDomain())
	if diags.HasError() {
		return nil, diags
	}
	result.DnsSubDomain = subdomain

	return result, nil
}

func NewAtlasReq(ctx context.Context, plan *TFModel) (*admin.StreamsPrivateLinkConnection, diag.Diagnostics) {
	if plan.Vendor.ValueString() == VendorConfluent {
		diags := diag.Diagnostics{}
		if plan.ServiceEndpointId.IsNull() {
			diags.AddError(fmt.Sprintf("service_endpoint_id is required for vendor %s", VendorConfluent), "")
		}
		if plan.DnsDomain.IsNull() {
			diags.AddError(fmt.Sprintf("dns_domain is required for vendor %s", VendorConfluent), "")
		}
		if plan.Region.IsNull() {
			diags.AddError(fmt.Sprintf("region is required for vendor %s", VendorConfluent), "")
		}
		if diags.HasError() {
			return nil, diags
		}
	}

	if plan.Vendor.ValueString() == VendorMSK {
		diags := diag.Diagnostics{}
		if plan.Arn.IsNull() {
			diags.AddError(fmt.Sprintf("arn is required for vendor %s", VendorMSK), "")
		}
		if plan.Region.ValueString() != "" {
			diags.AddError(fmt.Sprintf("region cannot be set for vendor %s", VendorMSK), "")
		}
		if diags.HasError() {
			return nil, diags
		}
	}

	result := &admin.StreamsPrivateLinkConnection{
		DnsDomain:         plan.DnsDomain.ValueStringPointer(),
		Provider:          plan.Provider.ValueString(),
		Region:            plan.Region.ValueStringPointer(),
		ServiceEndpointId: plan.ServiceEndpointId.ValueStringPointer(),
		State:             plan.State.ValueStringPointer(),
		Vendor:            plan.Vendor.ValueStringPointer(),
		Arn:               plan.Arn.ValueStringPointer(),
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
