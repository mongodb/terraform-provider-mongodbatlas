package cluster_test

import (
	"fmt"
	"testing"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigClusterRSCluster_withDefaultBiConnectorAndAdvancedConfiguration_backwardCompatibility(t *testing.T) {
	var (
		projectID   = mig.ProjectIDGlobal(t)
		clusterName = acc.RandomClusterName()
		cfg         = configAWS(projectID, clusterName, true, true)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            cfg,
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
				),
			},
			mig.TestStepCheckEmptyPlan(cfg),
		},
	})
}

func TestMigClusterRSCluster_basic_PartialAdvancedConf_backwardCompatibility(t *testing.T) {
	var (
		projectID   = mig.ProjectIDGlobal(t)
		clusterName = acc.RandomClusterName()
		cfgPartial  = configAdvancedConfPartial(projectID, clusterName, "false", &matlas.ProcessArgs{
			MinimumEnabledTLSProtocol: "TLS1_2",
		})
		cfgPartialUpdated = configAdvancedConfPartialUpdated(projectID, clusterName, "false", &matlas.ProcessArgs{
			MinimumEnabledTLSProtocol: "TLS1_2",
			SampleSizeBIConnector:     conversion.Pointer[int64](110),
		})
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            cfgPartial,
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.minimum_enabled_tls_protocol", "TLS1_2"),
				),
			},
			mig.TestStepCheckEmptyPlan(cfgPartial),
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

func configAdvancedConfPartialUpdated(projectID, name, autoscalingEnabled string, p *matlas.ProcessArgs) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "test" {
			project_id   = %[1]q
			name         = %[2]q
			disk_size_gb = 10

            cluster_type = "REPLICASET"
		    replication_specs {
			  num_shards = 1
			  regions_config {
			     region_name     = "US_WEST_2"
			     electable_nodes = 3
			     priority        = 7
                 read_only_nodes = 0
		       }
		    }

			backup_enabled               = false
			auto_scaling_disk_gb_enabled =  %[3]s

			// Provider Settings "block"
			provider_name               = "AWS"
			provider_instance_size_name = "M10"
			provider_region_name        = "US_WEST_2"

			advanced_configuration {
				minimum_enabled_tls_protocol         = %[4]q
				sample_size_bi_connector			 = %[5]d
			}
		}
	`, projectID, name, autoscalingEnabled, p.MinimumEnabledTLSProtocol, *p.SampleSizeBIConnector)
}
