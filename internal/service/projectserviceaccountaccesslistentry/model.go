package projectserviceaccountaccesslistentry

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20250312012/admin"
)

func NewMongoDBProjectServiceAccountAccessListEntry(model *TFProjectServiceAccountAccessListEntryModel) *[]admin.ServiceAccountIPAccessListEntry {
	return &[]admin.ServiceAccountIPAccessListEntry{
		{
			CidrBlock: conversion.StringPtr(model.CIDRBlock.ValueString()),
			IpAddress: conversion.StringPtr(model.IPAddress.ValueString()),
		},
	}
}

func NewTFProjectServiceAccountAccessListModel(projectID, clientID string, entry *admin.ServiceAccountIPAccessListEntry) *TFProjectServiceAccountAccessListEntryModel {
	return &TFProjectServiceAccountAccessListEntryModel{
		ProjectID:       types.StringValue(projectID),
		ClientID:        types.StringValue(clientID),
		CIDRBlock:       types.StringValue(entry.GetCidrBlock()),
		IPAddress:       types.StringValue(entry.GetIpAddress()),
		CreatedAt:       types.StringPointerValue(conversion.TimePtrToStringPtr(entry.CreatedAt)),
		LastUsedAddress: types.StringPointerValue(entry.LastUsedAddress),
		LastUsedAt:      types.StringPointerValue(conversion.TimePtrToStringPtr(entry.LastUsedAt)),
		RequestCount:    types.Int64PointerValue(conversion.IntPtrToInt64Ptr(entry.RequestCount)),
	}
}

func NewTFProjectServiceAccountAccessListEntriesPluralDSModel(projectID, clientID string, entries []admin.ServiceAccountIPAccessListEntry) *TFProjectServiceAccountAccessListEntriesPluralDSModel {
	results := make([]*TFProjectServiceAccountAccessListEntryModel, len(entries))
	for i := range entries {
		results[i] = NewTFProjectServiceAccountAccessListModel(projectID, clientID, &entries[i])
	}
	return &TFProjectServiceAccountAccessListEntriesPluralDSModel{
		ProjectID: types.StringValue(projectID),
		ClientID:  types.StringValue(clientID),
		Results:   results,
	}
}
