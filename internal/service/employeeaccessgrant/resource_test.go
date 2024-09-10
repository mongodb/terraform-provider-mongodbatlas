package employeeaccessgrant_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName   = "mongodbatlas_employee_access_grant.test"
	dataSourceName = "data.mongodbatlas_employee_access_grant.test"
)

func TestAccEmployeeAccessGrant_basic(t *testing.T) {
	resource.Test(t, *basicTestCase(t))
}

func basicTestCase(tb testing.TB) *resource.TestCase {
	tb.Helper()

	var (
		projectID, clusterName = acc.ClusterNameExecution(tb)
	)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(tb) },
		ExternalProviders:        acc.ExternalProvidersOnlyAWS(),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, clusterName, "CLUSTER_INFRASTRUCTURE", "2025-08-01T12:00:00Z"),
				Check:  checkBasic(projectID, clusterName, "CLUSTER_INFRASTRUCTURE", "2025-08-01T12:00:00Z"),
			},
		},
	}
}

func configBasic(projectID, clusterName, grantType, expirationTime string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_employee_access_grant" "test" {
			project_id 			= %[1]q
			cluster_name 		= %[2]q
			grant_type 			= %[3]q
			expiration_time = %[4]q
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
		if rs.Type == "mongodbatlas_employee_access_grant" {
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
