package projectipaccesslist

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20231115002/admin"
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
