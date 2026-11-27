package mongodbemployeeaccessgrant_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312028/admin"
	"go.mongodb.org/atlas-sdk/v20250312028/mockadmin"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/mongodbemployeeaccessgrant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName          = "mongodbatlas_mongodb_employee_access_grant.test"
	dataSourceName        = "data." + resourceName
	grantType             = "CLUSTER_INFRASTRUCTURE"
	grantTypeUpdated      = "CLUSTER_DATABASE_LOGS"
	grantTypeInvalid      = "invalid_grant_type"
	expirationTime        = "2025-08-01T12:00:00Z"
	expirationTimeUpdated = "2025-09-01T12:00:00Z"
	expirationTimeInvalid = "invalid_time"
)

func TestAccMongoDBEmployeeAccessGrant_basic(t *testing.T) {
	resource.Test(t, *basicTestCase(t))
}

func basicTestCase(tb testing.TB) *resource.TestCase {
	tb.Helper()
	projectID, clusterName := acc.ClusterNameExecution(tb, false)
	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, clusterName, grantType, expirationTime),
				Check:  checkBasic(projectID, clusterName, grantType, expirationTime),
			},
			{
				Config: configBasic(projectID, clusterName, grantTypeUpdated, expirationTime),
				Check:  checkBasic(projectID, clusterName, grantTypeUpdated, expirationTime),
			},
			{
				Config: configBasic(projectID, clusterName, grantTypeUpdated, expirationTimeUpdated),
				Check:  checkBasic(projectID, clusterName, grantTypeUpdated, expirationTimeUpdated),
			},
			{
				Config: configBasic(projectID, clusterName, grantType, expirationTime),
				Check:  checkBasic(projectID, clusterName, grantType, expirationTime),
			},
			{
				Config:                               configBasic(projectID, clusterName, grantType, expirationTime),
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    acc.ImportStateIDFuncProjectIDClusterName(resourceName, "project_id", "cluster_name"),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "cluster_name",
			},
		},
	}
}

func TestAccMongoDBEmployeeAccessGrant_invalidExpirationTime(t *testing.T) {
	projectID, clusterName := acc.ClusterNameExecution(t, false)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config:      configBasic(projectID, clusterName, grantType, expirationTimeInvalid),
				ExpectError: regexp.MustCompile("expiration_time format is incorrect.*" + expirationTimeInvalid),
			},
		},
	})
}

func TestAccMongoDBEmployeeAccessGrant_invalidExpirationTimeUpdate(t *testing.T) {
	projectID, clusterName := acc.ClusterNameExecution(t, false)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, clusterName, grantType, expirationTime),
				Check:  checkBasic(projectID, clusterName, grantType, expirationTime),
			},
			{
				Config:      configBasic(projectID, clusterName, grantType, expirationTimeInvalid),
				ExpectError: regexp.MustCompile("expiration_time format is incorrect.*" + expirationTimeInvalid),
			},
		},
	})
}

func TestAccMongoDBEmployeeAccessGrant_invalidGrantType(t *testing.T) {
	projectID, clusterName := acc.ClusterNameExecution(t, false)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config:      configBasic(projectID, clusterName, grantTypeInvalid, expirationTime),
				ExpectError: regexp.MustCompile("invalid enumeration value.*" + grantTypeInvalid),
			},
		},
	})
}

func TestAccMongoDBEmployeeAccessGrant_invalidGrantTypeUpdate(t *testing.T) {
	projectID, clusterName := acc.ClusterNameExecution(t, false)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, clusterName, grantType, expirationTime),
				Check:  checkBasic(projectID, clusterName, grantType, expirationTime),
			},
			{
				Config:      configBasic(projectID, clusterName, grantTypeInvalid, expirationTime),
				ExpectError: regexp.MustCompile("invalid enumeration value.*" + grantTypeInvalid),
			},
		},
	})
}

func TestDeleteSurfacesNonNotFoundErrors(t *testing.T) {
	ctx := t.Context()
	const (
		projectID   = "650000000000000000000001"
		clusterName = "cluster-test"
	)
	apiErr := errors.New("atlas api error")
	testCases := map[string]struct {
		httpResp    *http.Response
		err         error
		expectError bool
	}{
		"success":                         {httpResp: &http.Response{StatusCode: http.StatusNoContent}},
		"not found is treated as success": {httpResp: &http.Response{StatusCode: http.StatusNotFound}, err: apiErr},
		"bad request is surfaced":         {httpResp: &http.Response{StatusCode: http.StatusBadRequest}, err: apiErr, expectError: true},
		"unauthorized is surfaced":        {httpResp: &http.Response{StatusCode: http.StatusUnauthorized}, err: apiErr, expectError: true},
		"forbidden is surfaced":           {httpResp: &http.Response{StatusCode: http.StatusForbidden}, err: apiErr, expectError: true},
		"conflict is surfaced":            {httpResp: &http.Response{StatusCode: http.StatusConflict}, err: apiErr, expectError: true},
		"too many requests is surfaced":   {httpResp: &http.Response{StatusCode: http.StatusTooManyRequests}, err: apiErr, expectError: true},
		"server error is surfaced":        {httpResp: &http.Response{StatusCode: http.StatusInternalServerError}, err: apiErr, expectError: true},
		"transport error is surfaced":     {httpResp: nil, err: apiErr, expectError: true},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			api := mockadmin.NewClustersAPI(t)
			api.EXPECT().RevokeMongoEmployeeAccess(mock.Anything, projectID, clusterName).Return(admin.RevokeMongoEmployeeAccessApiRequest{ApiService: api}).Once()
			api.EXPECT().RevokeMongoEmployeeAccessExecute(mock.Anything).Return(tc.httpResp, tc.err).Once()

			r := mongodbemployeeaccessgrant.Resource()
			r.(config.ImplementedResource).SetClient(&config.MongoDBClient{AtlasV2: &admin.APIClient{ClustersAPI: api}})
			var schemaResp fwresource.SchemaResponse
			r.Schema(ctx, fwresource.SchemaRequest{}, &schemaResp)

			state := tfsdk.State{
				Schema: schemaResp.Schema,
				Raw: tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"project_id":      tftypes.NewValue(tftypes.String, projectID),
					"cluster_name":    tftypes.NewValue(tftypes.String, clusterName),
					"grant_type":      tftypes.NewValue(tftypes.String, grantType),
					"expiration_time": tftypes.NewValue(tftypes.String, expirationTime),
				}),
			}
			var resp fwresource.DeleteResponse
			r.Delete(ctx, fwresource.DeleteRequest{State: state}, &resp)

			if tc.expectError {
				require.Len(t, resp.Diagnostics.Errors(), 1)
				require.Equal(t, "Error deleting resource mongodbatlas_mongodb_employee_access_grant", resp.Diagnostics.Errors()[0].Summary())
				require.Contains(t, resp.Diagnostics.Errors()[0].Detail(), apiErr.Error())
				return
			}
			require.False(t, resp.Diagnostics.HasError(), "unexpected diagnostics: %#v", resp.Diagnostics.Errors())
		})
	}
}

func configBasic(projectID, clusterName, grantType, expirationTime string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_mongodb_employee_access_grant" "test" {
			project_id 			= %[1]q
			cluster_name 		= %[2]q
			grant_type 			= %[3]q
			expiration_time = %[4]q
		}

		data "mongodbatlas_mongodb_employee_access_grant" "test" {
			project_id 			= mongodbatlas_mongodb_employee_access_grant.test.project_id
			cluster_name 		= mongodbatlas_mongodb_employee_access_grant.test.cluster_name
			
			depends_on 			= [mongodbatlas_mongodb_employee_access_grant.test]
		}
	`, projectID, clusterName, grantType, expirationTime)
}

func checkBasic(projectID, clusterName, grantType, expirationTime string) resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{checkExists(resourceName)}
	attrsMap := map[string]string{
		"project_id":      projectID,
		"cluster_name":    clusterName,
		"grant_type":      grantType,
		"expiration_time": expirationTime,
	}
	checks = acc.AddAttrChecks(resourceName, checks, attrsMap)
	checks = acc.AddAttrChecks(dataSourceName, checks, attrsMap)
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if !exists(rs) {
			return fmt.Errorf("employee access grant (%s) does not exist", resourceName)
		}
		return nil
	}
}

func checkDestroy(state *terraform.State) error {
	for _, rs := range state.RootModule().Resources {
		if rs.Type == "mongodbatlas_mongodb_employee_access_grant" {
			if exists(rs) {
				return fmt.Errorf("employee access grant still exists")
			}
		}
	}
	return nil
}

func exists(rs *terraform.ResourceState) bool {
	projectID := rs.Primary.Attributes["project_id"]
	clusterName := rs.Primary.Attributes["cluster_name"]
	cluster, _, _ := acc.ConnV2().ClustersAPI.GetCluster(context.Background(), projectID, clusterName).Execute()
	resp, _ := cluster.GetMongoDBEmployeeAccessGrantOk()
	return resp != nil
}
