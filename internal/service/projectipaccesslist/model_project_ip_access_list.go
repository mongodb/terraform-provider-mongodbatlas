package projectipaccesslist

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20250312011/admin"
)

func getAccessListEntry(model *TfProjectIPAccessListModel) string {
	if model.CIDRBlock.ValueString() != "" {
		return model.CIDRBlock.ValueString()
	}
	if model.IPAddress.ValueString() != "" {
		return model.IPAddress.ValueString()
	}
	if model.AWSSecurityGroup.ValueString() != "" {
		return model.AWSSecurityGroup.ValueString()
	}
	return ""
}

func getAccessListEntryFromAPI(entry *admin.NetworkPermissionEntry) string {
	if entry.GetCidrBlock() != "" {
		return entry.GetCidrBlock()
	}
	if entry.GetIpAddress() != "" {
		return entry.GetIpAddress()
	}
	if entry.GetAwsSecurityGroup() != "" {
		return entry.GetAwsSecurityGroup()
	}
	return ""
}

func NewMongoDBProjectIPAccessList(projectIPAccessListModel *TfProjectIPAccessListModel) *[]admin.NetworkPermissionEntry {
	return &[]admin.NetworkPermissionEntry{
		{
			AwsSecurityGroup: conversion.StringPtr(projectIPAccessListModel.AWSSecurityGroup.ValueString()),
			CidrBlock:        conversion.StringPtr(projectIPAccessListModel.CIDRBlock.ValueString()),
			IpAddress:        conversion.StringPtr(projectIPAccessListModel.IPAddress.ValueString()),
			Comment:          conversion.StringPtr(projectIPAccessListModel.Comment.ValueString()),
		},
	}
}

func NewTfProjectIPAccessListModel(projectIPAccessListModel *TfProjectIPAccessListModel, projectIPAccessList *admin.NetworkPermissionEntry) *TfProjectIPAccessListModel {
	entry := getAccessListEntryFromAPI(projectIPAccessList)

	id := conversion.EncodeStateID(map[string]string{
		"entry":      entry,
		"project_id": projectIPAccessList.GetGroupId(),
	})

	return &TfProjectIPAccessListModel{
		ID:               types.StringValue(id),
		ProjectID:        types.StringValue(projectIPAccessList.GetGroupId()),
		CIDRBlock:        conversion.StringNullIfEmpty(projectIPAccessList.GetCidrBlock()),
		IPAddress:        conversion.StringNullIfEmpty(projectIPAccessList.GetIpAddress()),
		AWSSecurityGroup: conversion.StringNullIfEmpty(projectIPAccessList.GetAwsSecurityGroup()),
		Comment:          conversion.StringNullIfEmpty(projectIPAccessList.GetComment()),
		Timeouts:         projectIPAccessListModel.Timeouts,
	}
}

func NewTfProjectIPAccessListDSModel(ctx context.Context, accessList *admin.NetworkPermissionEntry) (*TfProjectIPAccessListDSModel, diag.Diagnostics) {
	databaseUserModel := &TfProjectIPAccessListDSModel{
		ProjectID:        types.StringValue(accessList.GetGroupId()),
		Comment:          conversion.StringNullIfEmpty(accessList.GetComment()),
		CIDRBlock:        conversion.StringNullIfEmpty(accessList.GetCidrBlock()),
		IPAddress:        conversion.StringNullIfEmpty(accessList.GetIpAddress()),
		AWSSecurityGroup: conversion.StringNullIfEmpty(accessList.GetAwsSecurityGroup()),
	}

	entry := getAccessListEntryFromAPI(accessList)

	id := conversion.EncodeStateID(map[string]string{
		"entry":      entry,
		"project_id": accessList.GetGroupId(),
	})

	databaseUserModel.ID = types.StringValue(id)
	return databaseUserModel, nil
}
