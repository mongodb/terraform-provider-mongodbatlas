package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	dataSourceProviderConfig = `
	resource "mongodbatlas_cloud_provider_access" "%[1]s" {
		project_id = "%[2]s"
		provider_name = "%[3]s"
	 }
	 
	 data "mongodbatlas_cloud_provider_access" "%[4]s" {
		project_id = mongodbatlas_cloud_provider_access.%[1]s.project_id
	 }	 
	`
)

func TestAccdataSourceMongoDBAtlasCloudProviderAccess_basic(t *testing.T) {
	var (
		suffix       = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
		name         = "datasource_test_role_basic" + suffix
		dataSCName   = "datasource_test_all_roles" + suffix
		resourceName = "mongodbatlas_cloud_provider_access." + name
		dsName       = "data.mongodbatlas_cloud_provider_access." + dataSCName
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	)

	config := fmt.Sprintf(dataSourceProviderConfig, name, projectID, "AWS", dataSCName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "atlas_assumed_role_external_id"),
					resource.TestCheckResourceAttrSet(resourceName, "atlas_aws_account_arn"),
					resource.TestCheckResourceAttr(dsName, "project_id", projectID),
					resource.TestCheckResourceAttrSet(dsName, "aws_iam_roles.0.atlas_assumed_role_external_id"),
				),
			},
		},
	},
	)
}
