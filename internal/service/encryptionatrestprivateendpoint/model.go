package encryptionatrestprivateendpoint

import (
	"go.mongodb.org/atlas-sdk/v20250312003/admin"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

func NewTFEarPrivateEndpoint(apiResp admin.EARPrivateEndpoint, projectID string) TFEarPrivateEndpointModel {
	return TFEarPrivateEndpointModel{
		ProjectID:                     types.StringValue(projectID),
		CloudProvider:                 conversion.StringNullIfEmpty(apiResp.GetCloudProvider()),
		ErrorMessage:                  conversion.StringNullIfEmpty(apiResp.GetErrorMessage()),
		ID:                            conversion.StringNullIfEmpty(apiResp.GetId()),
		RegionName:                    conversion.StringNullIfEmpty(apiResp.GetRegionName()),
		Status:                        conversion.StringNullIfEmpty(apiResp.GetStatus()),
		PrivateEndpointConnectionName: conversion.StringNullIfEmpty(apiResp.GetPrivateEndpointConnectionName()),
	}
}

func NewEarPrivateEndpointReq(tfPlan *TFEarPrivateEndpointModel) *admin.EARPrivateEndpoint {
	if tfPlan == nil {
		return nil
	}
	return &admin.EARPrivateEndpoint{
		CloudProvider: tfPlan.CloudProvider.ValueStringPointer(),
		RegionName:    tfPlan.RegionName.ValueStringPointer(),
	}
}

func NewTFEarPrivateEndpoints(projectID, cloudProvider string, sdkResults []admin.EARPrivateEndpoint) *TFEncryptionAtRestPrivateEndpointsDSModel {
	results := make([]TFEarPrivateEndpointModel, len(sdkResults))
	for i := range sdkResults {
		result := NewTFEarPrivateEndpoint(sdkResults[i], projectID)
		results[i] = result
	}
	return &TFEncryptionAtRestPrivateEndpointsDSModel{
		ProjectID:     types.StringValue(projectID),
		CloudProvider: types.StringValue(cloudProvider),
		Results:       results,
	}
}
