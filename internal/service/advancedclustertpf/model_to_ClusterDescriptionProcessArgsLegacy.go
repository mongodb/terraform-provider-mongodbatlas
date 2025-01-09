package advancedclustertpf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
)

func NewAtlasReqAdvancedConfigurationLegacy(ctx context.Context, objInput *types.Object, diags *diag.Diagnostics) *admin20240530.ClusterDescriptionProcessArgs {
	var resp *admin20240530.ClusterDescriptionProcessArgs
	if objInput == nil || objInput.IsUnknown() || objInput.IsNull() {
		return resp
	}
	input := &TFAdvancedConfigurationModel{}
	if localDiags := objInput.As(ctx, input, basetypes.ObjectAsOptions{}); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return resp
	}
	// Choosing to only handle legacy fields in the old API
	return &admin20240530.ClusterDescriptionProcessArgs{
		DefaultReadConcern:  conversion.NilForUnknown(input.DefaultReadConcern, input.DefaultReadConcern.ValueStringPointer()),
		FailIndexKeyTooLong: conversion.NilForUnknown(input.FailIndexKeyTooLong, input.FailIndexKeyTooLong.ValueBoolPointer()),
	}
}
