package projectipaccesslist

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20250312001/admin"
)

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
	entry := projectIPAccessList.GetIpAddress()
	if projectIPAccessList.GetCidrBlock() != "" {
		entry = projectIPAccessList.GetCidrBlock()
	} else if projectIPAccessList.GetAwsSecurityGroup() != "" {
		entry = projectIPAccessList.GetAwsSecurityGroup()
	}

	id := conversion.EncodeStateID(map[string]string{
		"entry":      entry,
		"project_id": projectIPAccessList.GetGroupId(),
	})

	return &TfProjectIPAccessListModel{
		ID:               types.StringValue(id),
		ProjectID:        types.StringValue(projectIPAccessList.GetGroupId()),
		CIDRBlock:        types.StringValue(projectIPAccessList.GetCidrBlock()),
		IPAddress:        types.StringValue(projectIPAccessList.GetIpAddress()),
		AWSSecurityGroup: types.StringValue(projectIPAccessList.GetAwsSecurityGroup()),
		Comment:          types.StringValue(projectIPAccessList.GetComment()),
		Timeouts:         projectIPAccessListModel.Timeouts,
	}
}

func NewTfProjectIPAccessListDSModel(ctx context.Context, accessList *admin.NetworkPermissionEntry) (*TfProjectIPAccessListDSModel, diag.Diagnostics) {
	databaseUserModel := &TfProjectIPAccessListDSModel{
		ProjectID:        types.StringValue(accessList.GetGroupId()),
		Comment:          types.StringValue(accessList.GetComment()),
		CIDRBlock:        types.StringValue(accessList.GetCidrBlock()),
		IPAddress:        types.StringValue(accessList.GetIpAddress()),
		AWSSecurityGroup: types.StringValue(accessList.GetAwsSecurityGroup()),
	}

	entry := accessList.GetCidrBlock()
	if accessList.GetIpAddress() != "" {
		entry = accessList.GetIpAddress()
	} else if accessList.GetAwsSecurityGroup() != "" {
		entry = accessList.GetAwsSecurityGroup()
	}

	id := conversion.EncodeStateID(map[string]string{
		"entry":      entry,
		"project_id": accessList.GetGroupId(),
	})

	databaseUserModel.ID = types.StringValue(id)
	return databaseUserModel, nil
}
