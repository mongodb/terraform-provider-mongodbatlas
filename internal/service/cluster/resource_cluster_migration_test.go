package cluster_test

import (
	"fmt"
	"os"
	"testing"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

// TestAccClusterRSCluster_withDefaultBiConnectorAndAdvancedConfiguration_backwardCompatibility In mongodbatlas provider v1.12.3
// 'schema.SchemaConfigModeAttr' was added to collection attributes to
// support Terraform versions <1.1.5. This setting is not recommended by Hashicorp and breaks CDKTF.
// This test ensures removing this setting does not have a user impact (only Terraform versions >=1.2.x can be supported)
func TestAccMigrationClusterRSCluster_withDefaultBiConnectorAndAdvancedConfiguration_backwardCompatibility(t *testing.T) {
	var (
		cluster      matlas.Cluster
		resourceName = "mongodbatlas_cluster.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		clusterName  = acc.RandomClusterName()
		cfg          = testAccMongoDBAtlasClusterConfigAWS(orgID, projectName, clusterName, true, true)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviders("1.12.3"),
				Config:            cfg,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasClusterAttributes(&cluster, clusterName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
				),
			},
			mig.TestStepCheckEmptyPlan(cfg),
		},
	})
}

// TestAccMigrationClusterRSCluster_basic_PartialAdvancedConf_backwardCompatibility In mongodbatlas provider v1.12.3
// 'schema.SchemaConfigModeAttr' was added to collection attributes to
// support Terraform versions <1.1.5. This setting is not recommended by Hashicorp and breaks CDKTF.
// This test ensures removing this setting does not have a user impact (only Terraform versions >=1.2.x can be supported)
func TestAccMigrationClusterRSCluster_basic_PartialAdvancedConf_backwardCompatibility(t *testing.T) {
	var (
		cluster      matlas.Cluster
		resourceName = "mongodbatlas_cluster.advance_conf"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		name         = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
		cfgPartial   = testAccMongoDBAtlasClusterConfigAdvancedConfPartial(orgID, projectName, name, "false", &matlas.ProcessArgs{
			MinimumEnabledTLSProtocol: "TLS1_2",
		})
		cfgPartialUpdated = testAccMongoDBAtlasClusterConfigAdvancedConfPartialUpdated(orgID, projectName, name, "false", &matlas.ProcessArgs{
			MinimumEnabledTLSProtocol: "TLS1_2",
			SampleSizeBIConnector:     conversion.Pointer[int64](110),
		})
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviders("1.12.3"),
				Config:            cfgPartial,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.minimum_enabled_tls_protocol", "TLS1_2"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   cfgPartial,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						acc.DebugPlan(),
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   cfgPartialUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.minimum_enabled_tls_protocol", "TLS1_2"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.sample_size_bi_connector", "110"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasClusterConfigAdvancedConfPartialUpdated(orgID, projectName, name, autoscalingEnabled string, p *matlas.ProcessArgs) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "cluster_project" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_cluster" "advance_conf" {
			project_id   = mongodbatlas_project.cluster_project.id
			name         = %[3]q
			disk_size_gb = 10

            cluster_type = "REPLICASET"
		    replication_specs {
			  num_shards = 1
			  regions_config {
			     region_name     = "EU_CENTRAL_1"
			     electable_nodes = 3
			     priority        = 7
                 read_only_nodes = 0
		       }
		    }

			backup_enabled               = false
			auto_scaling_disk_gb_enabled =  %[4]s

			// Provider Settings "block"
			provider_name               = "AWS"
			provider_instance_size_name = "M10"
			provider_region_name        = "EU_CENTRAL_1"

			advanced_configuration {
				minimum_enabled_tls_protocol         = %[5]q
				sample_size_bi_connector			 = %[6]d
			}
		}
	`, orgID, projectName, name, autoscalingEnabled, p.MinimumEnabledTLSProtocol, *p.SampleSizeBIConnector)
}
