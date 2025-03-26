package flexcluster_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var (
	resourceName         = "mongodbatlas_flex_cluster.test"
	dataSourceName       = "data.mongodbatlas_flex_cluster.test"
	dataSourcePluralName = "data.mongodbatlas_flex_clusters.test"
)

func TestAccFlexClusterRS_basic(t *testing.T) {
	tc := basicTestCase(t)
	// Tests include testing of plural data source and so cannot be run in parallel
	resource.Test(t, *tc)
}

func TestAccFlexClusterRS_failedUpdate(t *testing.T) {
	tc := failedUpdateTestCase(t)
	resource.Test(t, *tc)
}

func basicTestCase(t *testing.T) *resource.TestCase {
	t.Helper()
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomName()
		provider    = "AWS"
		region      = "US_EAST_1"
	)
	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyFlexCluster,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, clusterName, provider, region, true, false),
				Check:  checksFlexCluster(projectID, clusterName, true, false),
			},
			{
				Config: configBasic(projectID, clusterName, provider, region, false, true),
				Check:  checksFlexCluster(projectID, clusterName, false, true),
			},
			{
				Config:            configBasic(projectID, clusterName, provider, region, true, true),
				ResourceName:      resourceName,
				ImportStateIdFunc: acc.ImportStateIDFuncProjectIDClusterName(resourceName, "project_id", "name"),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	}
}

func failedUpdateTestCase(t *testing.T) *resource.TestCase {
	t.Helper()
	var (
		projectID          = acc.ProjectIDExecution(t)
		projectIDUpdated   = projectID + "-updated"
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
		CheckDestroy:             acc.CheckDestroyFlexCluster,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, clusterName, provider, region, false, false),
				Check:  checksFlexCluster(projectID, clusterName, false, false),
			},
			{
				Config:      configBasic(projectID, clusterNameUpdated, provider, region, false, false),
				ExpectError: regexp.MustCompile("name cannot be updated"),
			},
			{
				Config:      configBasic(projectIDUpdated, clusterName, provider, region, false, false),
				ExpectError: regexp.MustCompile("project_id cannot be updated"),
			},
			{
				Config:      configBasic(projectID, clusterName, providerUpdated, region, false, false),
				ExpectError: regexp.MustCompile("provider_settings.backing_provider_name cannot be updated"),
			},
			{
				Config:      configBasic(projectID, clusterName, provider, regionUpdated, false, false),
				ExpectError: regexp.MustCompile("provider_settings.region_name cannot be updated"),
			},
		},
	}
}

func configBasic(projectID, clusterName, provider, region string, terminationProtectionEnabled, tags bool) string {
	tagsConfig := ""
	if tags {
		tagsConfig = `
			tags = {
				testKey = "testValue"
			}`
	}
	return fmt.Sprintf(`
		resource "mongodbatlas_flex_cluster" "test" {
			project_id = %[1]q
			name       = %[2]q
			provider_settings = {
				backing_provider_name = %[3]q
				region_name           = %[4]q
			}
			termination_protection_enabled = %[5]t
			%[6]s
		}
		%[7]s
		`, projectID, clusterName, provider, region, terminationProtectionEnabled, tagsConfig, acc.FlexDataSource)
}

func checksFlexCluster(projectID, clusterName string, terminationProtectionEnabled, tagsCheck bool) resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{acc.CheckExistsFlexCluster()}
	attrMap := map[string]string{
		"project_id":                     projectID,
		"name":                           clusterName,
		"termination_protection_enabled": fmt.Sprintf("%v", terminationProtectionEnabled),
	}
	if tagsCheck {
		attrMap["tags.testKey"] = "testValue"
	}
	pluralMap := map[string]string{
		"project_id": projectID,
		"results.#":  "1",
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
	checks = acc.AddAttrChecks(dataSourcePluralName, checks, pluralMap)
	return acc.CheckRSAndDS(resourceName, &dataSourceName, &dataSourcePluralName, attrSet, attrMap, checks...)
}
