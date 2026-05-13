package projectipaccesslist_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312018/admin"
	"go.mongodb.org/atlas-sdk/v20250312018/mockadmin"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/serviceapi/projectipaccesslist"
)

const testProjectID = "507f1f77bcf86cd799439011"

func TestFetchAllEntries_Success(t *testing.T) {
	api := mockadmin.NewProjectIPAccessListApi(t)
	entries := []admin.NetworkPermissionEntry{
		{IpAddress: new("1.2.3.4"), Comment: new("test")},
		{CidrBlock: new("10.0.0.0/8")},
	}

	api.EXPECT().ListAccessListEntries(mock.Anything, testProjectID).
		Return(admin.ListAccessListEntriesApiRequest{ApiService: api})
	api.EXPECT().ListAccessListEntriesExecute(mock.Anything).
		Return(&admin.PaginatedNetworkAccess{Results: entries}, &http.Response{StatusCode: http.StatusOK}, nil)

	result, err := projectipaccesslist.FetchAllEntries(context.Background(), api, testProjectID)
	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "1.2.3.4", result[0].GetIpAddress())
	assert.Equal(t, "10.0.0.0/8", result[1].GetCidrBlock())
}

func TestFetchAllEntries_Pagination(t *testing.T) {
	api := mockadmin.NewProjectIPAccessListApi(t)

	page1 := make([]admin.NetworkPermissionEntry, 500)
	for i := range page1 {
		page1[i] = admin.NetworkPermissionEntry{IpAddress: new("1.2.3.4")}
	}
	page2 := []admin.NetworkPermissionEntry{
		{IpAddress: new("5.6.7.8")},
	}

	api.EXPECT().ListAccessListEntries(mock.Anything, testProjectID).
		Return(admin.ListAccessListEntriesApiRequest{ApiService: api}).Times(2)
	api.EXPECT().ListAccessListEntriesExecute(mock.Anything).
		Return(&admin.PaginatedNetworkAccess{Results: page1}, &http.Response{StatusCode: http.StatusOK}, nil).Once()
	api.EXPECT().ListAccessListEntriesExecute(mock.Anything).
		Return(&admin.PaginatedNetworkAccess{Results: page2}, &http.Response{StatusCode: http.StatusOK}, nil).Once()

	result, err := projectipaccesslist.FetchAllEntries(context.Background(), api, testProjectID)
	require.NoError(t, err)
	assert.Len(t, result, 501)
}

func TestFetchAllEntries_APIError(t *testing.T) {
	api := mockadmin.NewProjectIPAccessListApi(t)

	api.EXPECT().ListAccessListEntries(mock.Anything, testProjectID).
		Return(admin.ListAccessListEntriesApiRequest{ApiService: api})
	api.EXPECT().ListAccessListEntriesExecute(mock.Anything).
		Return(nil, &http.Response{StatusCode: http.StatusInternalServerError}, errors.New("internal server error"))

	_, err := projectipaccesslist.FetchAllEntries(context.Background(), api, testProjectID)
	require.Error(t, err)
}

func TestAccessListEntryValue(t *testing.T) {
	tests := []struct {
		name     string
		entry    admin.NetworkPermissionEntry
		expected string
	}{
		{
			name:     "ip address",
			entry:    admin.NetworkPermissionEntry{IpAddress: new("1.2.3.4")},
			expected: "1.2.3.4",
		},
		{
			name:     "cidr block",
			entry:    admin.NetworkPermissionEntry{CidrBlock: new("10.0.0.0/8")},
			expected: "10.0.0.0/8",
		},
		{
			name:     "aws security group",
			entry:    admin.NetworkPermissionEntry{AwsSecurityGroup: new("sg-12345")},
			expected: "sg-12345",
		},
		{
			name:     "ip address takes precedence over cidr",
			entry:    admin.NetworkPermissionEntry{IpAddress: new("1.2.3.4"), CidrBlock: new("10.0.0.0/8")},
			expected: "1.2.3.4",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := projectipaccesslist.AccessListEntryValue(&tc.entry)
			assert.Equal(t, tc.expected, result)
		})
	}
}
