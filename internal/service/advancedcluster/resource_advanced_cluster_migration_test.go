package advancedcluster_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigAdvancedCluster_singleAWSProvider(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		clusterName = acc.RandomClusterName()
		config      = configSingleProvider(orgID, projectName, clusterName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "retain_backups_enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestMigAdvancedCluster_multiCloud(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		clusterName = acc.RandomClusterName()
		config      = configMultiCloud(orgID, projectName, clusterName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "retain_backups_enabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestMigAdvancedCluster_partialAdvancedConf(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		rName       = acc.RandomClusterName()
		processArgs = `advanced_configuration  {
			fail_index_key_too_long              = false
			javascript_enabled                   = true
			minimum_enabled_tls_protocol         = "TLS1_1"
			no_table_scan                        = false
		  }`
		biConnector = `bi_connector_config {
			enabled = true
		  }`
		processArgsUpdated = `advanced_configuration  {
			fail_index_key_too_long              = false
			javascript_enabled                   = true
			minimum_enabled_tls_protocol         = "TLS1_1"
			no_table_scan                        = false
			default_read_concern                 = "available"
			sample_size_bi_connector			 = 110
            sample_refresh_interval_bi_connector = 310
		  }`
		biConnectorUpdated = `bi_connector_config {
			enabled = false
			read_preference = "secondary"
		  }`
		config        = configPartialAdvancedConfig(orgID, projectName, rName, processArgs, biConnector)
		configUpdated = configPartialAdvancedConfig(orgID, projectName, rName, processArgsUpdated, biConnectorUpdated)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.fail_index_key_too_long", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.javascript_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.minimum_enabled_tls_protocol", "TLS1_1"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.no_table_scan", "false"),
					resource.TestCheckResourceAttr(resourceName, "bi_connector_config.0.enabled", "true"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configUpdated,
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.fail_index_key_too_long", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.javascript_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.minimum_enabled_tls_protocol", "TLS1_1"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.no_table_scan", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.sample_refresh_interval_bi_connector", "310"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.sample_size_bi_connector", "110"),
					resource.TestCheckResourceAttr(resourceName, "bi_connector_config.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "bi_connector_config.0.read_preference", "secondary"),
				),
			},
		},
	})
}

func configPartialAdvancedConfig(orgID, projectName, name, processArgs, biConnector string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "cluster_project" {
	name   = %[2]q
	org_id = %[1]q
}	
resource "mongodbatlas_advanced_cluster" "test" {
  project_id             = mongodbatlas_project.cluster_project.id
  name                   = %[3]q
  cluster_type           = "REPLICASET"

  %[5]s

   replication_specs {
    region_configs {
      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }
      analytics_specs {
        instance_size = "M10"
        node_count    = 1
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "EU_WEST_1"
    }
  }

  %[4]s
    
}

	`, orgID, projectName, name, processArgs, biConnector)
}
