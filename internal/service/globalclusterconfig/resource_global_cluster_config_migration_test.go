package globalclusterconfig_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
	"go.mongodb.org/atlas-sdk/v20231115006/admin"
)

func TestAccMigrationClusterRSGlobalCluster_basic(t *testing.T) {
	var (
		globalConfig admin.GeoSharding
		clusterInfo  = acc.GetClusterInfo(&acc.ClusterRequest{Geosharded: true})
		config       = configBasic(&clusterInfo, false, false)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, &globalConfig),
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
			mig.TestStep(config),
		},
	})
}
