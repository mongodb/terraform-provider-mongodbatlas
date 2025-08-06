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

func TestAccFlexClusterRS_createTimeoutWithDeleteOnCreateFlex(t *testing.T) {
	var (
		projectID             = acc.ProjectIDExecution(t)
		clusterName           = acc.RandomName()
		provider              = "AWS"
		region                = "US_EAST_1"
		createTimeout         = "1s"
		deleteOnCreateTimeout = true
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config:      configBasic(projectID, clusterName, provider, region, acc.TimeoutConfig(&createTimeout, nil, nil, true), true, false, &deleteOnCreateTimeout),
				ExpectError: regexp.MustCompile("will run cleanup because delete_on_create_timeout is true"),
			},
		},
	})
}

func TestAccFlexClusterRS_updateDeleteTimeout(t *testing.T) {
	var (
		projectID     = acc.ProjectIDExecution(t)
		clusterName   = acc.RandomName()
		provider      = "AWS"
		region        = "US_EAST_1"
		updateTimeout = "1s"
		deleteTimeout = "1s"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, clusterName, provider, region, acc.TimeoutConfig(nil, &updateTimeout, &deleteTimeout, true), false, false, nil),
			},
			{
				Config:      configBasic(projectID, clusterName, provider, region, acc.TimeoutConfig(nil, &updateTimeout, &deleteTimeout, true), false, true, nil),
				ExpectError: regexp.MustCompile("timeout while waiting for state to become 'IDLE'"),
			},
			{
				Config:      acc.ConfigEmpty(), // triggers delete and because delete timeout is 1s, it times out
				ExpectError: regexp.MustCompile("timeout while waiting for state to become 'DELETED'"),
			},
			{
				// deletion of the flex cluster has been triggered, but has timed out in previous step, so this is needed in order to avoid "Error running post-test destroy, there may be dangling resource [...] Cluster already requested to be deleted"
				Config: acc.ConfigRemove(resourceName),
			},
		},
	})
}

func basicTestCase(t *testing.T) *resource.TestCase {
	t.Helper()
	var (
		projectID          = acc.ProjectIDExecution(t)
		clusterName        = acc.RandomName()
		provider           = "AWS"
		region             = "US_EAST_1"
		emptyTimeoutConfig = ""
	)
	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyFlexCluster,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, clusterName, provider, region, emptyTimeoutConfig, true, false, nil),
				Check:  checksFlexCluster(projectID, clusterName, true, false),
			},
			{
				Config: configBasic(projectID, clusterName, provider, region, emptyTimeoutConfig, false, true, nil),
				Check:  checksFlexCluster(projectID, clusterName, false, true),
			},
			{
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
		emptyTimeoutConfig = ""
	)
	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyFlexCluster,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, clusterName, provider, region, emptyTimeoutConfig, false, false, nil),
				Check:  checksFlexCluster(projectID, clusterName, false, false),
			},
			{
				Config:      configBasic(projectID, clusterNameUpdated, provider, region, emptyTimeoutConfig, false, false, nil),
				ExpectError: regexp.MustCompile("name cannot be updated"),
			},
			{
				Config:      configBasic(projectIDUpdated, clusterName, provider, region, emptyTimeoutConfig, false, false, nil),
				ExpectError: regexp.MustCompile("project_id cannot be updated"),
			},
			{
				Config:      configBasic(projectID, clusterName, providerUpdated, region, emptyTimeoutConfig, false, false, nil),
				ExpectError: regexp.MustCompile("provider_settings.backing_provider_name cannot be updated"),
			},
			{
				Config:      configBasic(projectID, clusterName, provider, regionUpdated, emptyTimeoutConfig, false, false, nil),
				ExpectError: regexp.MustCompile("provider_settings.region_name cannot be updated"),
			},
		},
	}
}

func configBasic(projectID, clusterName, provider, region, timeoutConfig string, terminationProtectionEnabled, tags bool, deleteOnCreateTimeout *bool) string {
	tagsConfig := ""
	if tags {
		tagsConfig = `
			tags = {
				testKey = "testValue"
			}`
	}
	deleteOnCreateTimeoutConfig := ""
	if deleteOnCreateTimeout != nil {
		deleteOnCreateTimeoutConfig = fmt.Sprintf(`
			delete_on_create_timeout = %[1]t
		`, *deleteOnCreateTimeout)
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
			%[7]s
			%[8]s
		}
		%[9]s
		`, projectID, clusterName, provider, region, terminationProtectionEnabled, deleteOnCreateTimeoutConfig, tagsConfig, timeoutConfig, acc.FlexDataSource)
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
