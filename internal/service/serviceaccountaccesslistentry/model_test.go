package serviceaccountaccesslistentry_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/serviceaccountaccesslistentry"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20250312011/admin"
)

var (
	modelOrgID           = "000000000000000000000000"
	modelClientID        = "mdb_sa_id_000000000000000000000000"
	modelCIDRBlock       = "192.168.1.0/24"
	modelIPAddress       = "192.168.1.1"
	modelCreatedAt       = "2025-01-01T01:01:00Z"
	modelLastUsedAddress = "192.168.1.2"
	modelLastUsedAt      = "2025-02-02T02:02:00Z"
	modelRequestCount    = 42

	modelCreatedAtTime, _  = conversion.StringToTime(modelCreatedAt)
	modelLastUsedAtTime, _ = conversion.StringToTime(modelLastUsedAt)
)

func TestNewTFServiceAccountAccessListEntryModel(t *testing.T) {
	testCases := map[string]struct {
		sdkEntry        *admin.ServiceAccountIPAccessListEntry
		expectedTFModel *serviceaccountaccesslistentry.TFServiceAccountAccessListEntryModel
		orgID           string
		clientID        string
	}{
		"Complete SDK response": {
			orgID:    modelOrgID,
			clientID: modelClientID,
			sdkEntry: &admin.ServiceAccountIPAccessListEntry{
				CidrBlock:       &modelCIDRBlock,
				IpAddress:       &modelIPAddress,
				CreatedAt:       &modelCreatedAtTime,
				LastUsedAddress: &modelLastUsedAddress,
				LastUsedAt:      &modelLastUsedAtTime,
				RequestCount:    &modelRequestCount,
			},
			expectedTFModel: &serviceaccountaccesslistentry.TFServiceAccountAccessListEntryModel{
				OrgID:           types.StringValue(modelOrgID),
				ClientID:        types.StringValue(modelClientID),
				CIDRBlock:       types.StringValue(modelCIDRBlock),
				IPAddress:       types.StringValue(modelIPAddress),
				CreatedAt:       types.StringValue(modelCreatedAt),
				LastUsedAddress: types.StringValue(modelLastUsedAddress),
				LastUsedAt:      types.StringValue(modelLastUsedAt),
				RequestCount:    types.Int64Value(int64(modelRequestCount)),
			},
		},
		"Only CIDR block": {
			orgID:    modelOrgID,
			clientID: modelClientID,
			sdkEntry: &admin.ServiceAccountIPAccessListEntry{
				CidrBlock: &modelCIDRBlock,
			},
			expectedTFModel: &serviceaccountaccesslistentry.TFServiceAccountAccessListEntryModel{
				OrgID:           types.StringValue(modelOrgID),
				ClientID:        types.StringValue(modelClientID),
				CIDRBlock:       types.StringValue(modelCIDRBlock),
				IPAddress:       types.StringValue(""),
				CreatedAt:       types.StringNull(),
				LastUsedAddress: types.StringNull(),
				LastUsedAt:      types.StringNull(),
				RequestCount:    types.Int64Null(),
			},
		},
		"Only IP address": {
			orgID:    modelOrgID,
			clientID: modelClientID,
			sdkEntry: &admin.ServiceAccountIPAccessListEntry{
				IpAddress: &modelIPAddress,
			},
			expectedTFModel: &serviceaccountaccesslistentry.TFServiceAccountAccessListEntryModel{
				OrgID:           types.StringValue(modelOrgID),
				ClientID:        types.StringValue(modelClientID),
				CIDRBlock:       types.StringValue(""),
				IPAddress:       types.StringValue(modelIPAddress),
				CreatedAt:       types.StringNull(),
				LastUsedAddress: types.StringNull(),
				LastUsedAt:      types.StringNull(),
				RequestCount:    types.Int64Null(),
			},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			resultModel := serviceaccountaccesslistentry.NewTFServiceAccountAccessListModel(tc.orgID, tc.clientID, tc.sdkEntry)
			assert.Equal(t, tc.expectedTFModel, resultModel, "created terraform model did not match expected output")
		})
	}
}

func TestNewMongoDBServiceAccountAccessListEntry(t *testing.T) {
	testCases := map[string]struct {
		tfModel        *serviceaccountaccesslistentry.TFServiceAccountAccessListEntryModel
		expectedSDKReq *[]admin.ServiceAccountIPAccessListEntry
	}{
		"With CIDR block": {
			tfModel: &serviceaccountaccesslistentry.TFServiceAccountAccessListEntryModel{
				OrgID:     types.StringValue(modelOrgID),
				ClientID:  types.StringValue(modelClientID),
				CIDRBlock: types.StringValue(modelCIDRBlock),
				IPAddress: types.StringNull(),
			},
			expectedSDKReq: &[]admin.ServiceAccountIPAccessListEntry{
				{
					CidrBlock: &modelCIDRBlock,
					IpAddress: nil,
				},
			},
		},
		"With IP address": {
			tfModel: &serviceaccountaccesslistentry.TFServiceAccountAccessListEntryModel{
				OrgID:     types.StringValue(modelOrgID),
				ClientID:  types.StringValue(modelClientID),
				CIDRBlock: types.StringNull(),
				IPAddress: types.StringValue(modelIPAddress),
			},
			expectedSDKReq: &[]admin.ServiceAccountIPAccessListEntry{
				{
					CidrBlock: nil,
					IpAddress: &modelIPAddress,
				},
			},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			apiReqResult := serviceaccountaccesslistentry.NewMongoDBServiceAccountAccessListEntry(tc.tfModel)
			assert.Equal(t, tc.expectedSDKReq, apiReqResult, "created sdk model did not match expected output")
		})
	}
}

func TestNewTFServiceAccountAccessListEntriesPluralDSModel(t *testing.T) {
	testCases := map[string]struct {
		expectedTFModel *serviceaccountaccesslistentry.TFServiceAccountAccessListEntriesPluralDSModel
		orgID           string
		clientID        string
		entries         []admin.ServiceAccountIPAccessListEntry
	}{
		"Single entry": {
			orgID:    modelOrgID,
			clientID: modelClientID,
			entries: []admin.ServiceAccountIPAccessListEntry{
				{
					CidrBlock:       &modelCIDRBlock,
					IpAddress:       &modelIPAddress,
					CreatedAt:       &modelCreatedAtTime,
					LastUsedAddress: &modelLastUsedAddress,
					LastUsedAt:      &modelLastUsedAtTime,
					RequestCount:    &modelRequestCount,
				},
			},
			expectedTFModel: &serviceaccountaccesslistentry.TFServiceAccountAccessListEntriesPluralDSModel{
				OrgID:    types.StringValue(modelOrgID),
				ClientID: types.StringValue(modelClientID),
				Results: []*serviceaccountaccesslistentry.TFServiceAccountAccessListEntryModel{
					{
						OrgID:           types.StringValue(modelOrgID),
						ClientID:        types.StringValue(modelClientID),
						CIDRBlock:       types.StringValue(modelCIDRBlock),
						IPAddress:       types.StringValue(modelIPAddress),
						CreatedAt:       types.StringValue(modelCreatedAt),
						LastUsedAddress: types.StringValue(modelLastUsedAddress),
						LastUsedAt:      types.StringValue(modelLastUsedAt),
						RequestCount:    types.Int64Value(int64(modelRequestCount)),
					},
				},
			},
		},
		"Multiple entries": {
			orgID:    modelOrgID,
			clientID: modelClientID,
			entries: []admin.ServiceAccountIPAccessListEntry{
				{
					CidrBlock: &modelCIDRBlock,
				},
				{
					IpAddress:       &modelIPAddress,
					CreatedAt:       &modelCreatedAtTime,
					LastUsedAddress: &modelLastUsedAddress,
					LastUsedAt:      &modelLastUsedAtTime,
					RequestCount:    &modelRequestCount,
				},
			},
			expectedTFModel: &serviceaccountaccesslistentry.TFServiceAccountAccessListEntriesPluralDSModel{
				OrgID:    types.StringValue(modelOrgID),
				ClientID: types.StringValue(modelClientID),
				Results: []*serviceaccountaccesslistentry.TFServiceAccountAccessListEntryModel{
					{
						OrgID:           types.StringValue(modelOrgID),
						ClientID:        types.StringValue(modelClientID),
						CIDRBlock:       types.StringValue(modelCIDRBlock),
						IPAddress:       types.StringValue(""),
						CreatedAt:       types.StringNull(),
						LastUsedAddress: types.StringNull(),
						LastUsedAt:      types.StringNull(),
						RequestCount:    types.Int64Null(),
					},
					{
						OrgID:           types.StringValue(modelOrgID),
						ClientID:        types.StringValue(modelClientID),
						CIDRBlock:       types.StringValue(""),
						IPAddress:       types.StringValue(modelIPAddress),
						CreatedAt:       types.StringValue(modelCreatedAt),
						LastUsedAddress: types.StringValue(modelLastUsedAddress),
						LastUsedAt:      types.StringValue(modelLastUsedAt),
						RequestCount:    types.Int64Value(int64(modelRequestCount)),
					},
				},
			},
		},
		"Empty entries": {
			orgID:    modelOrgID,
			clientID: modelClientID,
			entries:  []admin.ServiceAccountIPAccessListEntry{},
			expectedTFModel: &serviceaccountaccesslistentry.TFServiceAccountAccessListEntriesPluralDSModel{
				OrgID:    types.StringValue(modelOrgID),
				ClientID: types.StringValue(modelClientID),
				Results:  []*serviceaccountaccesslistentry.TFServiceAccountAccessListEntryModel{},
			},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			resultModel := serviceaccountaccesslistentry.NewTFServiceAccountAccessListEntriesPluralDSModel(tc.orgID, tc.clientID, tc.entries)
			assert.Equal(t, tc.expectedTFModel, resultModel, "created terraform model did not match expected output")
		})
	}
}
