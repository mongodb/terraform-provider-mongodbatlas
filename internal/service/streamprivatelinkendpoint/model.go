package streamprivatelinkendpoint

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"go.mongodb.org/atlas-sdk/v20250312008/admin"
)

const (
	VendorConfluent = "CONFLUENT"
	VendorMSK       = "MSK"
	VendorS3        = "S3"
	ProviderGCP     = "GCP"
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

	if len(apiResp.GetGcpServiceAttachmentUris()) > 0 {
		serviceAttachmentUris, diagsServiceAttachment := types.ListValueFrom(ctx, types.StringType, apiResp.GetGcpServiceAttachmentUris())
		if diagsServiceAttachment.HasError() {
			return nil, diagsServiceAttachment
		}
		result.ServiceAttachmentUris = serviceAttachmentUris
	} else {
		result.ServiceAttachmentUris = types.ListNull(types.StringType)
	}

	return result, nil
}

func NewAtlasReq(ctx context.Context, plan *TFModel) (*admin.StreamsPrivateLinkConnection, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if plan.Vendor.ValueString() == VendorConfluent {
		// Validate that exactly one of service_endpoint_id or service_attachment_uris is provided
		hasServiceEndpointID := !plan.ServiceEndpointId.IsNull() && plan.ServiceEndpointId.ValueString() != ""
		hasServiceAttachmentUris := !plan.ServiceAttachmentUris.IsNull() && len(plan.ServiceAttachmentUris.Elements()) > 0

		if !hasServiceEndpointID && !hasServiceAttachmentUris {
			diags.AddError(fmt.Sprintf("Either service_endpoint_id or service_attachment_uris must be provided for vendor %s", VendorConfluent), "")
		}
		if hasServiceEndpointID && hasServiceAttachmentUris {
			diags.AddError("Only one of service_endpoint_id or service_attachment_uris can be provided", "")
		}
		if plan.DnsDomain.IsNull() {
			diags.AddError(fmt.Sprintf("dns_domain is required for vendor %s", VendorConfluent), "")
		}
		if plan.Region.IsNull() {
			diags.AddError(fmt.Sprintf("region is required for vendor %s", VendorConfluent), "")
		}
	}

	if plan.Vendor.ValueString() == VendorMSK {
		if plan.Arn.IsNull() {
			diags.AddError(fmt.Sprintf("arn is required for vendor %s", VendorMSK), "")
		}
		if plan.Region.ValueString() != "" {
			diags.AddError(fmt.Sprintf("region cannot be set for vendor %s", VendorMSK), "")
		}
	}

	if plan.Vendor.ValueString() == VendorS3 {
		if plan.Region.IsNull() {
			diags.AddError(fmt.Sprintf("region is required for vendor %s", VendorS3), "")
		}
		if plan.ServiceEndpointId.IsNull() {
			diags.AddError(fmt.Sprintf("service_endpoint_id is required for vendor %s", VendorS3), "It should follow the format 'com.amazonaws.<region>.s3', for example 'com.amazonaws.us-east-1.s3'")
		}
	}

	if diags.HasError() {
		return nil, diags
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

	if !plan.ServiceAttachmentUris.IsNull() {
		var serviceAttachmentUris []string
		diags := plan.ServiceAttachmentUris.ElementsAs(ctx, &serviceAttachmentUris, false)
		if diags.HasError() {
			return nil, diags
		}
		result.GcpServiceAttachmentUris = &serviceAttachmentUris
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
