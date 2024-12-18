package advancedclustertpf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20241113003/admin"
)

func NewAtlasReqAdvancedConfiguration(ctx context.Context, input *types.List, diags *diag.Diagnostics) *admin.ClusterDescriptionProcessArgs20240805 {
	var resp *admin.ClusterDescriptionProcessArgs20240805
	if input == nil || input.IsUnknown() || input.IsNull() || len(input.Elements()) == 0 {
		return resp
	}
	elements := make([]TFAdvancedConfigurationModel, len(input.Elements()))
	diags.Append(input.ElementsAs(ctx, &elements, false)...)
	if diags.HasError() {
		return nil
	}
	item := elements[0]
	return &admin.ClusterDescriptionProcessArgs20240805{
		ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds: conversion.NilForUnknown(item.ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds, conversion.Int64PtrToIntPtr(item.ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds.ValueInt64Pointer())),
		DefaultWriteConcern:              conversion.NilForUnknown(item.DefaultWriteConcern, item.DefaultWriteConcern.ValueStringPointer()),
		JavascriptEnabled:                conversion.NilForUnknown(item.JavascriptEnabled, item.JavascriptEnabled.ValueBoolPointer()),
		MinimumEnabledTlsProtocol:        conversion.NilForUnknown(item.MinimumEnabledTlsProtocol, item.MinimumEnabledTlsProtocol.ValueStringPointer()),
		NoTableScan:                      conversion.NilForUnknown(item.NoTableScan, item.NoTableScan.ValueBoolPointer()),
		OplogMinRetentionHours:           conversion.NilForUnknown(item.OplogMinRetentionHours, item.OplogMinRetentionHours.ValueFloat64Pointer()),
		OplogSizeMB:                      conversion.NilForUnknown(item.OplogSizeMb, conversion.Int64PtrToIntPtr(item.OplogSizeMb.ValueInt64Pointer())),
		SampleRefreshIntervalBIConnector: conversion.NilForUnknown(item.SampleRefreshIntervalBiconnector, conversion.Int64PtrToIntPtr(item.SampleRefreshIntervalBiconnector.ValueInt64Pointer())),
		SampleSizeBIConnector:            conversion.NilForUnknown(item.SampleSizeBiconnector, conversion.Int64PtrToIntPtr(item.SampleSizeBiconnector.ValueInt64Pointer())),
		TransactionLifetimeLimitSeconds:  conversion.NilForUnknown(item.TransactionLifetimeLimitSeconds, item.TransactionLifetimeLimitSeconds.ValueInt64Pointer()),
	}
}
