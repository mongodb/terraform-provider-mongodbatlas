package streaminstance_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccStreamDSStreamInstance_basic(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_stream_instance.test"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acc.RandomProjectName()
		instanceName   = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckPreviewFlag(t); acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyStreamInstance,
		Steps: []resource.TestStep{
			{
				Config: streamInstanceDataSourceConfig(orgID, projectName, instanceName, region, cloudProvider),
				Check: resource.ComposeTestCheckFunc(
					streamInstanceAttributeChecks(dataSourceName, instanceName, region, cloudProvider),
					resource.TestCheckResourceAttr(dataSourceName, "stream_config.tier", "SP30"),
				),
			},
		},
	})
}

func streamInstanceDataSourceConfig(orgID, projectName, instanceName, region, cloudProvider string) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_stream_instance" "test" {
			project_id = mongodbatlas_stream_instance.test.project_id
			instance_name = mongodbatlas_stream_instance.test.instance_name
		}
	`, acc.StreamInstanceConfig(orgID, projectName, instanceName, region, cloudProvider))
}
