package advancedclustertpf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
)

func NewAtlasReqAdvancedConfigurationLegacy(ctx context.Context, input *ObjectValueOf[TFAdvancedConfigurationModel], diags *diag.Diagnostics) *admin20240530.ClusterDescriptionProcessArgs {
	var out admin20240530.ClusterDescriptionProcessArgs
	if input == nil || input.IsUnknown() || input.IsNull() {
		return &out
	}
	tf, localDiags := input.ToPtr(ctx)
	diags.Append(localDiags...)
	if diags.HasError() {
		return &out
	}
	return &admin20240530.ClusterDescriptionProcessArgs{
		DefaultReadConcern:  conversion.NilForUnknown(tf.DefaultReadConcern, tf.DefaultReadConcern.ValueStringPointer()),
		FailIndexKeyTooLong: conversion.NilForUnknown(tf.FailIndexKeyTooLong, tf.FailIndexKeyTooLong.ValueBoolPointer()),
	}
}
