package globalclusterconfig_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccClusterRSGlobalCluster_basic(t *testing.T) {
	var (
		clusterInfo = acc.GetClusterInfo(t, &acc.ClusterRequest{Geosharded: true})
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(&clusterInfo, false, false),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "custom_zone_mappings.#"),
					resource.TestCheckResourceAttrSet(resourceName, "custom_zone_mapping.%"),
					resource.TestCheckResourceAttrSet(resourceName, "custom_zone_mapping.CA"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterInfo.ClusterName),
					resource.TestCheckResourceAttr(resourceName, "managed_namespaces.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "managed_namespaces.0.is_custom_shard_key_hashed", "false"),
					resource.TestCheckResourceAttr(resourceName, "managed_namespaces.0.is_shard_key_unique", "false"),
				),
			},
			{
				Config: configBasic(&clusterInfo, true, false),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "managed_namespaces.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "managed_namespaces.0.is_custom_shard_key_hashed", "true"),
					resource.TestCheckResourceAttr(resourceName, "managed_namespaces.0.is_shard_key_unique", "false"),
				),
			},
			{
				Config: configBasic(&clusterInfo, false, true),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "managed_namespaces.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "managed_namespaces.0.is_custom_shard_key_hashed", "false"),
					resource.TestCheckResourceAttr(resourceName, "managed_namespaces.0.is_shard_key_unique", "true"),
				),
			},
		},
	})
}

func TestAccClusterRSGlobalCluster_withAWSAndBackup(t *testing.T) {
	var (
		clusterInfo = acc.GetClusterInfo(t, &acc.ClusterRequest{Geosharded: true, CloudBackup: true})
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(&clusterInfo, false, false),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "managed_namespaces.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "custom_zone_mappings.#"),
					resource.TestCheckResourceAttrSet(resourceName, "custom_zone_mapping.%"),
					resource.TestCheckResourceAttrSet(resourceName, "custom_zone_mapping.CA"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterInfo.ClusterName),
				),
			},
		},
	})
}

func TestAccClusterRSGlobalCluster_importBasic(t *testing.T) {
	var (
		clusterInfo = acc.GetClusterInfo(t, &acc.ClusterRequest{Geosharded: true})
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(&clusterInfo, false, false),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       importStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"custom_zone_mappings"},
			},
		},
	})
}

func TestAccClusterRSGlobalCluster_database(t *testing.T) {
	var (
		clusterInfo = acc.GetClusterInfo(t, &acc.ClusterRequest{Geosharded: true, ExtraConfig: zonesStr})
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithDBConfig(&clusterInfo, customZone),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "managed_namespaces.#", "5"),
					resource.TestCheckResourceAttrSet(resourceName, "custom_zone_mappings.#"),
					resource.TestCheckResourceAttrSet(resourceName, "custom_zone_mapping.%"),
					resource.TestCheckResourceAttrSet(resourceName, "custom_zone_mapping.US"),
					resource.TestCheckResourceAttrSet(resourceName, "custom_zone_mapping.IE"),
					resource.TestCheckResourceAttrSet(resourceName, "custom_zone_mapping.DE"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterInfo.ClusterName),
				),
			},
			{
				Config: configWithDBConfig(&clusterInfo, customZoneUpdated),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "managed_namespaces.#", "5"),
					resource.TestCheckResourceAttrSet(resourceName, "custom_zone_mappings.#"),
					resource.TestCheckResourceAttrSet(resourceName, "custom_zone_mapping.%"),
					resource.TestCheckResourceAttrSet(resourceName, "custom_zone_mapping.US"),
					resource.TestCheckResourceAttrSet(resourceName, "custom_zone_mapping.IE"),
					resource.TestCheckResourceAttrSet(resourceName, "custom_zone_mapping.DE"),
					resource.TestCheckResourceAttrSet(resourceName, "custom_zone_mapping.JP"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterInfo.ClusterName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       importStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"custom_zone_mappings"},
			},
		},
	})
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		config, _, err := acc.ConnV2().GlobalClustersApi.GetManagedNamespace(context.Background(), ids["project_id"], ids["cluster_name"]).Execute()
		if err == nil {
			if len(config.GetCustomZoneMapping()) > 0 || len(config.GetManagedNamespaces()) > 0 {
				return nil
			}
		}
		return fmt.Errorf("global config for cluster(%s) does not exist", ids["cluster_name"])
	}
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		return fmt.Sprintf("%s-%s", ids["project_id"], ids["cluster_name"]), nil
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_global_cluster_config" {
			continue
		}

		globalConfig, _, err := acc.ConnV2().GlobalClustersApi.GetManagedNamespace(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.Attributes["cluster_name"]).Execute()
		if err != nil {
			if strings.Contains(err.Error(), "NOT_FOUND") {
				return nil
			}
			return err
		}

		if len(globalConfig.GetCustomZoneMapping()) > 0 || len(globalConfig.GetManagedNamespaces()) > 0 {
			return fmt.Errorf("global cluster configuration for cluster(%s) still exists", rs.Primary.Attributes["cluster_name"])
		}
	}
	return nil
}

func configBasic(info *acc.ClusterInfo, isCustomShard, isShardKeyUnique bool) string {
	return info.ClusterTerraformStr + fmt.Sprintf(`
		resource "mongodbatlas_global_cluster_config" "config" {
			cluster_name     = %[1]s
			project_id       = %[2]s

			managed_namespaces {
				db               		   = "mydata"
				collection       		   = "publishers"
				custom_shard_key		   = "city"
				is_custom_shard_key_hashed = %[3]t
				is_shard_key_unique 	   = %[4]t
			}

			custom_zone_mappings {
				location = "CA"
				zone     = "Zone 1"
			}
		}

		data "mongodbatlas_global_cluster_config" "config" {
			cluster_name     = %[1]s
			project_id       = %[2]s
		}	
	`, info.ClusterNameStr, info.ProjectIDStr, isCustomShard, isShardKeyUnique)
}

func configWithDBConfig(info *acc.ClusterInfo, zones string) string {
	return info.ClusterTerraformStr + fmt.Sprintf(`
		resource "mongodbatlas_global_cluster_config" "config" {
			cluster_name     = %[1]s
			project_id       = %[2]s

			managed_namespaces {
				db               = "horizonv2-sg"
				collection       = "entitlements.entitlement"
				custom_shard_key = "orgId"
			}
			managed_namespaces {
				db               = "horizonv2-sg"
				collection       = "entitlements.homesitemapping"
				custom_shard_key = "orgId"
			}
			managed_namespaces {
				db               = "horizonv2-sg"
				collection       = "entitlements.site"
				custom_shard_key = "orgId"
			}
			managed_namespaces {
				db               = "horizonv2-sg"
				collection       = "entitlements.userDesktopMapping"
				custom_shard_key = "orgId"
			}
			managed_namespaces {
				db               = "horizonv2-sg"
				collection       = "session"
				custom_shard_key = "orgId"
			}
			%[3]s
		}
	`, info.ClusterNameStr, info.ProjectIDStr, zones)
}

const (
	resourceName   = "mongodbatlas_global_cluster_config.config"
	dataSourceName = "data.mongodbatlas_global_cluster_config.config"

	customZone = `
		custom_zone_mappings {
			location = "US"
			zone     = "US"
		}
		custom_zone_mappings {
			location = "IE"
			zone     = "EU"
		}
		custom_zone_mappings {
			location = "DE"
			zone     = "DE"
		}
	`
	customZoneUpdated = `
		custom_zone_mappings {
			location = "US"
			zone     = "US"
		}
		custom_zone_mappings {
			location = "IE"
			zone     = "EU"
		}
		custom_zone_mappings {
			location = "DE"
			zone     = "DE"
		}
		custom_zone_mappings {
			location = "JP"
			zone     = "JP"
		}
	`

	zonesStr = `
		replication_specs {
			zone_name  = "US"
			num_shards = 1
			regions_config {
				region_name     = "US_EAST_1"
				electable_nodes = 3
				priority        = 7
				read_only_nodes = 0
			}
		}
		replication_specs {
			zone_name  = "EU"
			num_shards = 1
			regions_config {
				region_name     = "EU_WEST_1"
				electable_nodes = 3
				priority        = 7
				read_only_nodes = 0
			}
		}
		replication_specs {
			zone_name  = "DE"
			num_shards = 1
			regions_config {
				region_name     = "EU_NORTH_1"
				electable_nodes = 3
				priority        = 7
				read_only_nodes = 0
			}
		}
		replication_specs {
			zone_name  = "JP"
			num_shards = 1
			regions_config {
				region_name     = "AP_NORTHEAST_1"
				electable_nodes = 3
				priority        = 7
				read_only_nodes = 0
			}
		}
	`
)
