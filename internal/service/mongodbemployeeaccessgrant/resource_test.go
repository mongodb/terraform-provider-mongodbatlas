package mongodbemployeeaccessgrant_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
	projectID, clusterName := acc.ClusterNameExecution(tb)
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
		},
	}
}

func TestAccMongoDBEmployeeAccessGrant_invalidExpirationTime(t *testing.T) {
	projectID, clusterName := acc.ClusterNameExecution(t)
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
	projectID, clusterName := acc.ClusterNameExecution(t)
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
	projectID, clusterName := acc.ClusterNameExecution(t)
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
	projectID, clusterName := acc.ClusterNameExecution(t)
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
	cluster, _, _ := acc.ConnV2().ClustersApi.GetCluster(context.Background(), projectID, clusterName).Execute()
	resp, _ := cluster.GetMongoDBEmployeeAccessGrantOk()
	return resp != nil
}
