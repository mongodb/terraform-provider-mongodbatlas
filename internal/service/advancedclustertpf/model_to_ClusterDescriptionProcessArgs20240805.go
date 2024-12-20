package advancedclustertpf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20241113003/admin"
)

func NewAtlasReqAdvancedConfiguration(ctx context.Context, input *ObjectValueOf[TFAdvancedConfigurationModel], diags *diag.Diagnostics) *admin.ClusterDescriptionProcessArgs20240805 {
	var out admin.ClusterDescriptionProcessArgs20240805
	if input == nil || input.IsUnknown() || input.IsNull() {
		return &out
	}
	tf, localDiags := input.ToPtr(ctx)
	diags.Append(localDiags...)
	if diags.HasError() {
		return &out
	}
	return &admin.ClusterDescriptionProcessArgs20240805{
		ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds: conversion.NilForUnknown(tf.ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds, conversion.Int64PtrToIntPtr(tf.ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds.ValueInt64Pointer())),
		DefaultWriteConcern:              conversion.NilForUnknown(tf.DefaultWriteConcern, tf.DefaultWriteConcern.ValueStringPointer()),
		JavascriptEnabled:                conversion.NilForUnknown(tf.JavascriptEnabled, tf.JavascriptEnabled.ValueBoolPointer()),
		MinimumEnabledTlsProtocol:        conversion.NilForUnknown(tf.MinimumEnabledTlsProtocol, tf.MinimumEnabledTlsProtocol.ValueStringPointer()),
		NoTableScan:                      conversion.NilForUnknown(tf.NoTableScan, tf.NoTableScan.ValueBoolPointer()),
		OplogMinRetentionHours:           conversion.NilForUnknown(tf.OplogMinRetentionHours, tf.OplogMinRetentionHours.ValueFloat64Pointer()),
		OplogSizeMB:                      conversion.NilForUnknown(tf.OplogSizeMb, conversion.Int64PtrToIntPtr(tf.OplogSizeMb.ValueInt64Pointer())),
		SampleRefreshIntervalBIConnector: conversion.NilForUnknown(tf.SampleRefreshIntervalBiconnector, conversion.Int64PtrToIntPtr(tf.SampleRefreshIntervalBiconnector.ValueInt64Pointer())),
		SampleSizeBIConnector:            conversion.NilForUnknown(tf.SampleSizeBiconnector, conversion.Int64PtrToIntPtr(tf.SampleSizeBiconnector.ValueInt64Pointer())),
		TransactionLifetimeLimitSeconds:  conversion.NilForUnknown(tf.TransactionLifetimeLimitSeconds, tf.TransactionLifetimeLimitSeconds.ValueInt64Pointer()),
	}
}
