package flexcluster_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var (
	resourceType   = "mongodbatlas_flex_cluster"
	resourceName   = "mongodbatlas_flex_cluster.test"
	dataSourceName = "data.mongodbatlas_flex_cluster.test"
)

func TestAccFlexClusterRS_basic(t *testing.T) {
	tc := basicTestCase(t)
	resource.ParallelTest(t, *tc)
}

func TestAccFlexClusterRS_failedUpdate(t *testing.T) {
	tc := failedUpdateTestCase(t)
	resource.ParallelTest(t, *tc)
}

func basicTestCase(t *testing.T) *resource.TestCase {
	t.Helper()
	var (
		projectID   = os.Getenv("MONGODB_ATLAS_FLEX_PROJECT_ID")
		clusterName = acc.RandomName()
		provider    = "AWS"
		region      = "US_EAST_1"
	)
	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, clusterName, provider, region, true),
				Check:  checksFlexCluster(projectID, clusterName, true),
			},
			{
				Config: configBasic(projectID, clusterName, provider, region, false),
				Check:  checksFlexCluster(projectID, clusterName, false),
			},
			{
				Config:            configBasic(projectID, clusterName, provider, region, true),
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	}
}

func failedUpdateTestCase(t *testing.T) *resource.TestCase {
	t.Helper()
	var (
		projectID          = os.Getenv("MONGODB_ATLAS_FLEX_PROJECT_ID")
		projectIDUpdated   = os.Getenv("MONGODB_ATLAS_FLEX_PROJECT_ID") + "-updated"
		clusterName        = acc.RandomName()
		clusterNameUpdated = clusterName + "-updated"
		provider           = "AWS"
		providerUpdated    = "GCP"
		region             = "US_EAST_1"
		regionUpdated      = "US_EAST_2"
	)
	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, clusterName, provider, region, false),
				Check:  checksFlexCluster(projectID, clusterName, false),
			},
			{
				Config:      configBasic(projectID, clusterNameUpdated, provider, region, false),
				ExpectError: regexp.MustCompile("name cannot be updated"),
			},
			{
				Config:      configBasic(projectIDUpdated, clusterName, provider, region, false),
				ExpectError: regexp.MustCompile("project_id cannot be updated"),
			},
			{
				Config:      configBasic(projectID, clusterName, providerUpdated, region, false),
				ExpectError: regexp.MustCompile("provider_settings.backing_provider_name cannot be updated"),
			},
			{
				Config:      configBasic(projectID, clusterName, provider, regionUpdated, false),
				ExpectError: regexp.MustCompile("provider_settings.region_name cannot be updated"),
			},
		},
	}
}

func configBasic(projectID, clusterName, provider, region string, terminationProtectionEnabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_flex_cluster" "test" {
			project_id = %[1]q
			name       = %[2]q
			provider_settings = {
				backing_provider_name = %[3]q
				region_name           = %[4]q
			}
			termination_protection_enabled = %[5]t
		}
		data "mongodbatlas_flex_cluster" "test" {
			project_id = mongodbatlas_flex_cluster.test.project_id
			name       = mongodbatlas_flex_cluster.test.name
		}`, projectID, clusterName, provider, region, terminationProtectionEnabled)
}

func checksFlexCluster(projectID, clusterName string, terminationProtectionEnabled bool) resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{checkExists()}
	attrMap := map[string]string{
		"project_id":                     projectID,
		"name":                           clusterName,
		"termination_protection_enabled": fmt.Sprintf("%v", terminationProtectionEnabled),
	}
	attrSet := []string{
		"backup_settings.enabled",
		"cluster_type",
		"connection_strings.standard",
		"create_date",
		"id",
		"mongo_db_version",
		"state_name",
		"version_release_system",
		"provider_settings.provider_name",
	}
	checks = acc.AddAttrChecks(resourceName, checks, attrMap)
	checks = acc.AddAttrChecks(dataSourceName, checks, attrMap)
	checks = acc.AddAttrSetChecks(resourceName, checks, attrSet...)
	checks = acc.AddAttrSetChecks(dataSourceName, checks, attrSet...)
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func checkExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type == resourceType {
				projectID := rs.Primary.Attributes["project_id"]
				name := rs.Primary.Attributes["name"]
				_, _, err := acc.ConnV2().FlexClustersApi.GetFlexCluster(context.Background(), projectID, name).Execute()
				if err != nil {
					return fmt.Errorf("flex cluster (%s:%s) not found", projectID, name)
				}
			}
		}
		return nil
	}
}

func checkDestroy(state *terraform.State) error {
	for _, rs := range state.RootModule().Resources {
		if rs.Type == resourceType {
			projectID := rs.Primary.Attributes["project_id"]
			name := rs.Primary.Attributes["name"]
			_, _, err := acc.ConnV2().FlexClustersApi.GetFlexCluster(context.Background(), projectID, name).Execute()
			if err == nil {
				return fmt.Errorf("flex cluster (%s:%s) still exists", projectID, name)
			}
		}
	}
	return nil
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return fmt.Sprintf("%s-%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["name"]), nil
	}
}
