package encryptionatrestprivateendpoint

import (
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20240805001/admin"
)

func NewTFEarPrivateEndpoint(apiResp *admin.EARPrivateEndpoint) *TFEarPrivateEndpointModel {
	if apiResp == nil {
		return nil
	}
	return &TFEarPrivateEndpointModel{
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
		CloudProvider:                 tfPlan.CloudProvider.ValueStringPointer(),
		ErrorMessage:                  tfPlan.ErrorMessage.ValueStringPointer(),
		Id:                            tfPlan.ID.ValueStringPointer(),
		RegionName:                    tfPlan.RegionName.ValueStringPointer(),
		Status:                        tfPlan.Status.ValueStringPointer(),
		PrivateEndpointConnectionName: tfPlan.PrivateEndpointConnectionName.ValueStringPointer(),
	}
}
