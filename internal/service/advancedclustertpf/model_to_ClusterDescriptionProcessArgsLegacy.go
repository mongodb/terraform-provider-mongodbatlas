package advancedclustertpf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
)

func NewAtlasReqAdvancedConfigurationLegacy(ctx context.Context, input *types.List, diags *diag.Diagnostics) *admin20240530.ClusterDescriptionProcessArgs {
	return conversion.SingleListTFToSDK(ctx, diags, input, func(tf TFAdvancedConfigurationModel) *admin20240530.ClusterDescriptionProcessArgs {
		// Choosing to only handle legacy fields in the old API
		return &admin20240530.ClusterDescriptionProcessArgs{
			DefaultReadConcern:  conversion.NilForUnknown(tf.DefaultReadConcern, tf.DefaultReadConcern.ValueStringPointer()),
			FailIndexKeyTooLong: conversion.NilForUnknown(tf.FailIndexKeyTooLong, tf.FailIndexKeyTooLong.ValueBoolPointer()),
		}
	})
}
