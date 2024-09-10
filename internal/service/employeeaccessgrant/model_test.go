package employeeaccessgrant_test

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/employeeaccessgrant"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20240805003/admin"
)

func TestNewTFModel(t *testing.T) {
	testCases := map[string]struct {
		apiResp         *admin.EmployeeAccessGrant
		expectedTFModel *employeeaccessgrant.TFModel
		projectID       string
		clusterName     string
	}{
		"valid": {
			projectID:   "projectID",
			clusterName: "clusterName",
			apiResp: &admin.EmployeeAccessGrant{
				GrantType:      "grantType",
				ExpirationTime: time.Date(2024, 10, 13, 0, 0, 0, 0, time.UTC),
			},
			expectedTFModel: &employeeaccessgrant.TFModel{
				ProjectID:      types.StringValue("projectID"),
				ClusterName:    types.StringValue("clusterName"),
				GrantType:      types.StringValue("grantType"),
				ExpirationTime: types.StringValue("2024-10-13T00:00:00Z"),
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tfModel := employeeaccessgrant.NewTFModel(tc.projectID, tc.clusterName, tc.apiResp)
			assert.Equal(t, tc.expectedTFModel, tfModel)
		})
	}
}

func TestNewAtlasReq(t *testing.T) {
	testCases := map[string]struct {
		tfModel        *employeeaccessgrant.TFModel
		expectedReq    *admin.EmployeeAccessGrant
		expectedHasErr bool
	}{
		"valid": {
			tfModel: &employeeaccessgrant.TFModel{
				ProjectID:      types.StringValue("projectID"),
				ClusterName:    types.StringValue("clusterName"),
				GrantType:      types.StringValue("grantType"),
				ExpirationTime: types.StringValue("2024-10-13T00:00:00Z"),
			},
			expectedReq: &admin.EmployeeAccessGrant{
				GrantType:      "grantType",
				ExpirationTime: time.Date(2024, 10, 13, 0, 0, 0, 0, time.UTC),
			},
		},
		"invalid expiration time": {
			tfModel: &employeeaccessgrant.TFModel{
				ProjectID:      types.StringValue("projectID"),
				ClusterName:    types.StringValue("clusterName"),
				GrantType:      types.StringValue("grantType"),
				ExpirationTime: types.StringValue("invalid_time"),
			},
			expectedHasErr: true,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			req, err := employeeaccessgrant.NewAtlasReq(tc.tfModel)
			assert.Equal(t, tc.expectedHasErr, err != nil)
			if err == nil {
				assert.Equal(t, tc.expectedReq, req)
			}
		})
	}
}
