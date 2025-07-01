package mongodbemployeeaccessgrant_test

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/mongodbemployeeaccessgrant"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20250312005/admin"
)

func TestNewTFModel(t *testing.T) {
	testCases := map[string]struct {
		apiResp         *admin.EmployeeAccessGrant
		expectedTFModel *mongodbemployeeaccessgrant.TFModel
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
			expectedTFModel: &mongodbemployeeaccessgrant.TFModel{
				ProjectID:      types.StringValue("projectID"),
				ClusterName:    types.StringValue("clusterName"),
				GrantType:      types.StringValue("grantType"),
				ExpirationTime: types.StringValue("2024-10-13T00:00:00Z"),
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tfModel := mongodbemployeeaccessgrant.NewTFModel(tc.projectID, tc.clusterName, tc.apiResp)
			assert.Equal(t, tc.expectedTFModel, tfModel)
		})
	}
}

func TestNewAtlasReq(t *testing.T) {
	testCases := map[string]struct {
		tfModel             *mongodbemployeeaccessgrant.TFModel
		expectedReq         *admin.EmployeeAccessGrant
		expectedErrContains string
	}{
		"valid": {
			tfModel: &mongodbemployeeaccessgrant.TFModel{
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
			tfModel: &mongodbemployeeaccessgrant.TFModel{
				ProjectID:      types.StringValue("projectID"),
				ClusterName:    types.StringValue("clusterName"),
				GrantType:      types.StringValue("grantType"),
				ExpirationTime: types.StringValue("invalid_time"),
			},
			expectedErrContains: "invalid_time",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			req, err := mongodbemployeeaccessgrant.NewAtlasReq(tc.tfModel)
			assert.Equal(t, tc.expectedErrContains == "", err == nil)
			if err == nil {
				assert.Equal(t, tc.expectedReq, req)
			} else {
				assert.Contains(t, err.Error(), tc.expectedErrContains)
			}
		})
	}
}
