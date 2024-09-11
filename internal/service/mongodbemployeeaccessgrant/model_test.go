package mongodbemployeeaccessgrant_test

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/mongodbemployeeaccessgrant"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20240805003/admin"
)

func TestNewTFModel(t *testing.T) {
	testCases := map[string]struct {
		apiResp             *admin.EmployeeAccessGrant
		expectedTFModel     *mongodbemployeeaccessgrant.TFModel
		expectedErrContains string
		projectID           string
		clusterName         string
	}{
		"valid": {
			projectID:   "123456789012345678901234",
			clusterName: "clusterName",
			apiResp: &admin.EmployeeAccessGrant{
				GrantType:      "grantType",
				ExpirationTime: time.Date(2024, 10, 13, 0, 0, 0, 0, time.UTC),
			},
			expectedTFModel: &mongodbemployeeaccessgrant.TFModel{
				ID:             types.StringValue("123456789012345678901234-clusterName"),
				ProjectID:      types.StringValue("123456789012345678901234"),
				ClusterName:    types.StringValue("clusterName"),
				GrantType:      types.StringValue("grantType"),
				ExpirationTime: types.StringValue("2024-10-13T00:00:00Z"),
			},
		},
		"invalid project_id": {
			projectID:           "invalid_project_id",
			clusterName:         "clusterName",
			expectedErrContains: "invalid_project_id",
		},
		"invalid cluster_name": {
			projectID:           "123456789012345678901234",
			clusterName:         "_invalid_cluster_name",
			expectedErrContains: "_invalid_cluster_name",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tfModel, err := mongodbemployeeaccessgrant.NewTFModel(tc.projectID, tc.clusterName, tc.apiResp)
			assert.Equal(t, tc.expectedErrContains == "", err == nil)
			if err == nil {
				assert.Equal(t, tc.expectedTFModel, tfModel)
			} else {
				assert.Contains(t, err.Error(), tc.expectedErrContains)
			}
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
