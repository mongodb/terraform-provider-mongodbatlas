package projectserviceaccountaccesslistentry_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/projectserviceaccountaccesslistentry"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20250312013/admin"
)

var (
	modelProjectID       = "000000000000000000000000"
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

func TestNewTFProjectServiceAccountAccessListEntryModel(t *testing.T) {
	testCases := map[string]struct {
		sdkEntry        *admin.ServiceAccountIPAccessListEntry
		expectedTFModel *projectserviceaccountaccesslistentry.TFProjectServiceAccountAccessListEntryModel
		projectID       string
		clientID        string
	}{
		"Complete SDK response": {
			projectID: modelProjectID,
			clientID:  modelClientID,
			sdkEntry: &admin.ServiceAccountIPAccessListEntry{
				CidrBlock:       &modelCIDRBlock,
				IpAddress:       &modelIPAddress,
				CreatedAt:       &modelCreatedAtTime,
				LastUsedAddress: &modelLastUsedAddress,
				LastUsedAt:      &modelLastUsedAtTime,
				RequestCount:    &modelRequestCount,
			},
			expectedTFModel: &projectserviceaccountaccesslistentry.TFProjectServiceAccountAccessListEntryModel{
				ProjectID:       types.StringValue(modelProjectID),
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
			projectID: modelProjectID,
			clientID:  modelClientID,
			sdkEntry: &admin.ServiceAccountIPAccessListEntry{
				CidrBlock: &modelCIDRBlock,
			},
			expectedTFModel: &projectserviceaccountaccesslistentry.TFProjectServiceAccountAccessListEntryModel{
				ProjectID:       types.StringValue(modelProjectID),
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
			projectID: modelProjectID,
			clientID:  modelClientID,
			sdkEntry: &admin.ServiceAccountIPAccessListEntry{
				IpAddress: &modelIPAddress,
			},
			expectedTFModel: &projectserviceaccountaccesslistentry.TFProjectServiceAccountAccessListEntryModel{
				ProjectID:       types.StringValue(modelProjectID),
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
			resultModel := projectserviceaccountaccesslistentry.NewTFProjectServiceAccountAccessListModel(tc.projectID, tc.clientID, tc.sdkEntry)
			assert.Equal(t, tc.expectedTFModel, resultModel, "created terraform model did not match expected output")
		})
	}
}

func TestNewMongoDBProjectServiceAccountAccessListEntry(t *testing.T) {
	testCases := map[string]struct {
		tfModel        *projectserviceaccountaccesslistentry.TFProjectServiceAccountAccessListEntryModel
		expectedSDKReq *[]admin.ServiceAccountIPAccessListEntry
	}{
		"With CIDR block": {
			tfModel: &projectserviceaccountaccesslistentry.TFProjectServiceAccountAccessListEntryModel{
				ProjectID: types.StringValue(modelProjectID),
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
			tfModel: &projectserviceaccountaccesslistentry.TFProjectServiceAccountAccessListEntryModel{
				ProjectID: types.StringValue(modelProjectID),
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
			apiReqResult := projectserviceaccountaccesslistentry.NewMongoDBProjectServiceAccountAccessListEntry(tc.tfModel)
			assert.Equal(t, tc.expectedSDKReq, apiReqResult, "created sdk model did not match expected output")
		})
	}
}

func TestNewTFProjectServiceAccountAccessListEntriesPluralDSModel(t *testing.T) {
	testCases := map[string]struct {
		expectedTFModel *projectserviceaccountaccesslistentry.TFProjectServiceAccountAccessListEntriesPluralDSModel
		projectID       string
		clientID        string
		entries         []admin.ServiceAccountIPAccessListEntry
	}{
		"Single entry": {
			projectID: modelProjectID,
			clientID:  modelClientID,
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
			expectedTFModel: &projectserviceaccountaccesslistentry.TFProjectServiceAccountAccessListEntriesPluralDSModel{
				ProjectID: types.StringValue(modelProjectID),
				ClientID:  types.StringValue(modelClientID),
				Results: []*projectserviceaccountaccesslistentry.TFProjectServiceAccountAccessListEntryModel{
					{
						ProjectID:       types.StringValue(modelProjectID),
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
			projectID: modelProjectID,
			clientID:  modelClientID,
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
			expectedTFModel: &projectserviceaccountaccesslistentry.TFProjectServiceAccountAccessListEntriesPluralDSModel{
				ProjectID: types.StringValue(modelProjectID),
				ClientID:  types.StringValue(modelClientID),
				Results: []*projectserviceaccountaccesslistentry.TFProjectServiceAccountAccessListEntryModel{
					{
						ProjectID:       types.StringValue(modelProjectID),
						ClientID:        types.StringValue(modelClientID),
						CIDRBlock:       types.StringValue(modelCIDRBlock),
						IPAddress:       types.StringValue(""),
						CreatedAt:       types.StringNull(),
						LastUsedAddress: types.StringNull(),
						LastUsedAt:      types.StringNull(),
						RequestCount:    types.Int64Null(),
					},
					{
						ProjectID:       types.StringValue(modelProjectID),
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
			projectID: modelProjectID,
			clientID:  modelClientID,
			entries:   []admin.ServiceAccountIPAccessListEntry{},
			expectedTFModel: &projectserviceaccountaccesslistentry.TFProjectServiceAccountAccessListEntriesPluralDSModel{
				ProjectID: types.StringValue(modelProjectID),
				ClientID:  types.StringValue(modelClientID),
				Results:   []*projectserviceaccountaccesslistentry.TFProjectServiceAccountAccessListEntryModel{},
			},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			resultModel := projectserviceaccountaccesslistentry.NewTFProjectServiceAccountAccessListEntriesPluralDSModel(tc.projectID, tc.clientID, tc.entries)
			assert.Equal(t, tc.expectedTFModel, resultModel, "created terraform model did not match expected output")
		})
	}
}
