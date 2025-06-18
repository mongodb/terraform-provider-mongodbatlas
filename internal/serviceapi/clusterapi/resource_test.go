package clusterapi_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const resourceName = "mongodbatlas_cluster_api.test"

func TestAccClusterAPI_basic(t *testing.T) {
	var (
		groupID     = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(groupID, clusterName, "M10"),
				Check:  checkBasic(groupID, clusterName, "M10"),
			},
			{
				Config: configBasic(groupID, clusterName, "M30"),
				Check:  checkBasic(groupID, clusterName, "M30"),
			},
			{
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    importStateIDFunc(resourceName),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore: []string{
					"retain_backups_enabled",   // This field is TF specific and not returned by Atlas, so Import can't fill it in.
					"mongo_db_major_version",   // Risks plan change of 8 --> 8.0 (always normalized to `major.minor`)
					"state_name",               // Cluster state can change from IDLE to UPDATING and risks making the test flaky
					"delete_on_create_timeout", // This field is TF specific and not returned by Atlas, so Import can't fill it in.
				},
			},
		},
	})
}

func configBasic(groupID, clusterName, instanceSize string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster_api" "test" {
			group_id     = %[1]q
			name         = %[2]q
			cluster_type = "SHARDED"
			replication_specs = [{
				region_configs = [{
					provider_name = "AWS"
					region_name   = "US_EAST_1"
					priority      = 7
					electable_specs = {
						node_count    = 3
						instance_size = %[3]q
					}
				}]
			}]
		}
	`, groupID, clusterName, instanceSize)
}

func checkBasic(groupID, clusterName, instanceSize string) resource.TestCheckFunc {
	mapChecks := map[string]string{
		"group_id": groupID,
		"name":     clusterName,
		"replication_specs.0.region_configs.0.electable_specs.instance_size": instanceSize,
	}
	checks := acc.AddAttrChecks(resourceName, nil, mapChecks)
	checks = append(checks, checkExists(resourceName))
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		groupID := rs.Primary.Attributes["group_id"]
		name := rs.Primary.Attributes["name"]
		if groupID == "" || name == "" {
			return "", fmt.Errorf("import, attributes not found for: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", groupID, name), nil
	}
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		groupID := rs.Primary.Attributes["group_id"]
		clusterName := rs.Primary.Attributes["name"]
		if groupID == "" || clusterName == "" {
			return fmt.Errorf("groupID or clusterName is empty: %s, %s", groupID, clusterName)
		}
		if _, _, err := acc.ConnV2().ClustersApi.GetCluster(context.Background(), groupID, clusterName).Execute(); err != nil {
			return fmt.Errorf("cluster(%s:%s) does not exist: %w", groupID, clusterName, err)
		}
		return nil
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_cluster_api" {
			continue
		}
		groupID := rs.Primary.Attributes["group_id"]
		clusterName := rs.Primary.Attributes["name"]
		if groupID == "" || clusterName == "" {
			return fmt.Errorf("groupID or clusterName is empty: %s, %s", groupID, clusterName)
		}
		resp, _, _ := acc.ConnV2().ClustersApi.GetCluster(context.Background(), groupID, clusterName).Execute()
		if resp.GetId() != "" {
			return fmt.Errorf("cluster (%s:%s) still exists", clusterName, rs.Primary.ID)
		}
	}
	return nil
}
