package globalclusterconfig_test

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName     = "mongodbatlas_global_cluster_config.config"
	dataSourceName   = "data.mongodbatlas_global_cluster_config.config"
	dataSourceConfig = `
	data "mongodbatlas_global_cluster_config" "config" {
		project_id   = mongodbatlas_global_cluster_config.config.project_id
		cluster_name = mongodbatlas_global_cluster_config.config.cluster_name
		depends_on   = [mongodbatlas_global_cluster_config.config]
	}
	`
)

func TestAccGlobalClusterConfig_basic(t *testing.T) {
	resource.ParallelTest(t, *basicTestCase(t, true, false))
}

func TestAccGlobalClusterConfig_withBackup(t *testing.T) {
	resource.ParallelTest(t, *basicTestCase(t, true, true))
}

func TestAccGlobalClusterConfig_iss(t *testing.T) {
	const (
		zone1 = "Zone 1"
		zone2 = "Zone 2"
	)
	var (
		replicationSpec1 = acc.ReplicationSpecRequest{
			ZoneName:     zone1,
			Region:       "US_EAST_1",
			InstanceSize: "M30",
		}
		replicationSpec2 = acc.ReplicationSpecRequest{
			ZoneName:     zone2,
			Region:       "US_EAST_2",
			InstanceSize: "M10",
		}
		clusterInfo = acc.GetClusterInfo(t, &acc.ClusterRequest{Geosharded: true, ReplicationSpecs: []acc.ReplicationSpecRequest{replicationSpec1, replicationSpec2}})
		attrsMap    = map[string]string{
			"cluster_name":         clusterInfo.Name,
			"managed_namespaces.#": "1",
			"managed_namespaces.0.is_custom_shard_key_hashed": "false",
			"managed_namespaces.0.is_shard_key_unique":        "false",
			"custom_zone_mapping.%":                           "0",
			"custom_zone_mapping_zone_id.%":                   "2",
		}
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, &clusterInfo, "", ""),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configISS(&clusterInfo, false, false, zone1, zone2),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					acc.CheckRSAndDS(resourceName, conversion.Pointer(dataSourceName), nil, []string{"project_id"}, attrsMap),
				),
			},
		},
	})
}

func basicTestCase(tb testing.TB, checkZoneID, withBackup bool) *resource.TestCase {
	tb.Helper()
	clusterInfo := acc.GetClusterInfo(tb, &acc.ClusterRequest{Geosharded: true, CloudBackup: withBackup})
	attrsMap := map[string]string{
		"cluster_name":         clusterInfo.Name,
		"managed_namespaces.#": "1",
		"managed_namespaces.0.is_custom_shard_key_hashed": "false",
		"managed_namespaces.0.is_shard_key_unique":        "false",
		"custom_zone_mapping.%":                           "1",
	}
	if checkZoneID {
		attrsMap["custom_zone_mapping_zone_id.%"] = "1"
	}

	return &resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(tb, &clusterInfo, "", ""),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(&clusterInfo, false, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					checkZone(0, "CA", clusterInfo.ResourceName, checkZoneID),
					acc.CheckRSAndDS(resourceName, conversion.Pointer(dataSourceName), nil, []string{"project_id"}, attrsMap)),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       importStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"custom_zone_mappings"},
			},
			{
				Config:      configBasic(&clusterInfo, true, false),
				ExpectError: regexp.MustCompile("managed namespace for collection 'publishers' in db 'mydata' cannot be modified"),
			},
		},
	}
}

func TestAccGlobalClusterConfig_database(t *testing.T) {
	const (
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
	)

	var (
		specUS      = acc.ReplicationSpecRequest{ZoneName: "US", Region: "US_EAST_1"}
		specEU      = acc.ReplicationSpecRequest{ZoneName: "EU", Region: "EU_WEST_1"}
		specDE      = acc.ReplicationSpecRequest{ZoneName: "DE", Region: "EU_NORTH_1"}
		specJP      = acc.ReplicationSpecRequest{ZoneName: "JP", Region: "AP_NORTHEAST_1"}
		clusterInfo = acc.GetClusterInfo(t, &acc.ClusterRequest{Geosharded: true, ReplicationSpecs: []acc.ReplicationSpecRequest{specUS, specEU, specDE, specJP}})
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, &clusterInfo, "", ""),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithDBConfig(&clusterInfo, customZone),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					checkZone(0, "US", clusterInfo.ResourceName, true),
					checkZone(1, "IE", clusterInfo.ResourceName, true),
					checkZone(2, "DE", clusterInfo.ResourceName, true),
					acc.CheckRSAndDS(resourceName, conversion.Pointer(dataSourceName), nil,
						[]string{"project_id"},
						map[string]string{
							"cluster_name":         clusterInfo.Name,
							"managed_namespaces.#": "5",
							"managed_namespaces.0.is_custom_shard_key_hashed": "false",
							"managed_namespaces.0.is_shard_key_unique":        "false",
							"custom_zone_mapping_zone_id.%":                   "3",
							"custom_zone_mapping.%":                           "3",
						}),
				),
			},
			{
				Config: configWithDBConfig(&clusterInfo, customZoneUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					checkZone(0, "US", clusterInfo.ResourceName, true),
					checkZone(1, "IE", clusterInfo.ResourceName, true),
					checkZone(2, "DE", clusterInfo.ResourceName, true),
					checkZone(3, "JP", clusterInfo.ResourceName, true),
					acc.CheckRSAndDS(resourceName, conversion.Pointer(dataSourceName), nil,
						[]string{"project_id"},
						map[string]string{
							"cluster_name":         clusterInfo.Name,
							"managed_namespaces.#": "5",
							"managed_namespaces.0.is_custom_shard_key_hashed": "false",
							"managed_namespaces.0.is_shard_key_unique":        "false",
							"custom_zone_mapping_zone_id.%":                   "4",
							"custom_zone_mapping.%":                           "4",
						}),
				),
			},
			{
				Config:      configWithDBConfig(&clusterInfo, customZone),
				ExpectError: regexp.MustCompile("partial deletion of custom_zone_mappings is not allowed; remove either all mappings or none"),
			},
			{
				Config: configWithDBConfig(&clusterInfo, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					acc.CheckRSAndDS(resourceName, conversion.Pointer(dataSourceName), nil,
						[]string{"project_id"},
						map[string]string{
							"cluster_name":         clusterInfo.Name,
							"managed_namespaces.#": "5",
							"managed_namespaces.0.is_custom_shard_key_hashed": "false",
							"managed_namespaces.0.is_shard_key_unique":        "false",
							"custom_zone_mapping_zone_id.%":                   "0",
							"custom_zone_mapping.%":                           "0",
						}),
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

func checkZone(pos int, zone, clusterName string, checkZoneID bool) resource.TestCheckFunc {
	firstID := fmt.Sprintf("custom_zone_mapping.%s", zone)
	secondID := fmt.Sprintf("replication_specs.%d.id", pos)
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrPair(resourceName, firstID, clusterName, secondID),
		resource.TestCheckResourceAttrPair(dataSourceName, firstID, clusterName, secondID),
	}
	if checkZoneID {
		firstZoneID := fmt.Sprintf("custom_zone_mapping_zone_id.%s", zone)
		secondZoneID := fmt.Sprintf("replication_specs.%d.zone_id", pos)
		checks = append(checks,
			resource.TestCheckResourceAttrPair(resourceName, firstZoneID, clusterName, secondZoneID),
			resource.TestCheckResourceAttrPair(dataSourceName, firstZoneID, clusterName, secondZoneID),
		)
	}
	return resource.ComposeAggregateTestCheckFunc(checks...)
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
	return info.TerraformStr + fmt.Sprintf(`
		resource "mongodbatlas_global_cluster_config" "config" {
			cluster_name     = %[1]s
			project_id       = %[2]q

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

			depends_on = [%[5]s]
		}
	`, info.TerraformNameRef, info.ProjectID, isCustomShard, isShardKeyUnique, info.ResourceName) + dataSourceConfig
}

func configWithDBConfig(info *acc.ClusterInfo, zones string) string {
	return info.TerraformStr + fmt.Sprintf(`
		resource "mongodbatlas_global_cluster_config" "config" {
			cluster_name     = %[1]s
			project_id       = %[2]q

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

			depends_on = [%[4]s]
		}
	`, info.TerraformNameRef, info.ProjectID, zones, info.ResourceName) + dataSourceConfig
}

func configISS(info *acc.ClusterInfo, isCustomShard, isShardKeyUnique bool, zone1, zone2 string) string {
	return info.TerraformStr + fmt.Sprintf(`
		resource "mongodbatlas_global_cluster_config" "config" {
			cluster_name     = %[1]s
			project_id       = %[2]q

			managed_namespaces {
				db               		   = "mydata"
				collection       		   = "publishers"
				custom_shard_key		   = "city"
				is_custom_shard_key_hashed = %[3]t
				is_shard_key_unique 	   = %[4]t
			}
			depends_on = [%[5]s]

			custom_zone_mappings {
				location = "US"
				zone     = %[6]q
			}
			custom_zone_mappings {
				location = "IE"
				zone     = %[7]q
			}
		}
	`, info.TerraformNameRef, info.ProjectID, isCustomShard, isShardKeyUnique, info.ResourceName, zone1, zone2) + dataSourceConfig
}
