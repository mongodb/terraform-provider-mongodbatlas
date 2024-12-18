package advancedclustertpf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
)

func NewAtlasReqAdvancedConfigurationLegacy(ctx context.Context, input *types.List, diags *diag.Diagnostics) *admin20240530.ClusterDescriptionProcessArgs {
	var resp *admin20240530.ClusterDescriptionProcessArgs
	if input == nil || input.IsUnknown() || input.IsNull() || len(input.Elements()) == 0 {
		return resp
	}
	elements := make([]TFAdvancedConfigurationModel, len(input.Elements()))
	diags.Append(input.ElementsAs(ctx, &elements, false)...)
	if diags.HasError() {
		return nil
	}
	item := elements[0]

	// Choosing to only handle legacy fields in the old API
	return &admin20240530.ClusterDescriptionProcessArgs{
		DefaultReadConcern:  conversion.NilForUnknown(item.DefaultReadConcern, item.DefaultReadConcern.ValueStringPointer()),
		FailIndexKeyTooLong: conversion.NilForUnknown(item.FailIndexKeyTooLong, item.FailIndexKeyTooLong.ValueBoolPointer()),
	}
}
