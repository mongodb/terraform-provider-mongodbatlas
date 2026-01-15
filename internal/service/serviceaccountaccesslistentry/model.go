package serviceaccountaccesslistentry

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20250312012/admin"
)

func NewMongoDBServiceAccountAccessListEntry(model *TFServiceAccountAccessListEntryModel) *[]admin.ServiceAccountIPAccessListEntry {
	return &[]admin.ServiceAccountIPAccessListEntry{
		{
			CidrBlock: conversion.StringPtr(model.CIDRBlock.ValueString()),
			IpAddress: conversion.StringPtr(model.IPAddress.ValueString()),
		},
	}
}

func NewTFServiceAccountAccessListModel(orgID, clientID string, entry *admin.ServiceAccountIPAccessListEntry) *TFServiceAccountAccessListEntryModel {
	return &TFServiceAccountAccessListEntryModel{
		OrgID:           types.StringValue(orgID),
		ClientID:        types.StringValue(clientID),
		CIDRBlock:       types.StringValue(entry.GetCidrBlock()),
		IPAddress:       types.StringValue(entry.GetIpAddress()),
		CreatedAt:       types.StringPointerValue(conversion.TimePtrToStringPtr(entry.CreatedAt)),
		LastUsedAddress: types.StringPointerValue(entry.LastUsedAddress),
		LastUsedAt:      types.StringPointerValue(conversion.TimePtrToStringPtr(entry.LastUsedAt)),
		RequestCount:    types.Int64PointerValue(conversion.IntPtrToInt64Ptr(entry.RequestCount)),
	}
}

func NewTFServiceAccountAccessListEntriesPluralDSModel(orgID, clientID string, entries []admin.ServiceAccountIPAccessListEntry) *TFServiceAccountAccessListEntriesPluralDSModel {
	results := make([]*TFServiceAccountAccessListEntryModel, len(entries))
	for i := range entries {
		results[i] = NewTFServiceAccountAccessListModel(orgID, clientID, &entries[i])
	}
	return &TFServiceAccountAccessListEntriesPluralDSModel{
		OrgID:    types.StringValue(orgID),
		ClientID: types.StringValue(clientID),
		Results:  results,
	}
}
