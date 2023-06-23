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
	resource "mongodbatlas_project" "test" {
		name   = %[3]q
		org_id = %[2]q
	}
	resource "mongodbatlas_cloud_provider_access" %[1]q {
		project_id = mongodbatlas_project.test.id
		provider_name = %[4]q
	 }
	 
	 data "mongodbatlas_cloud_provider_access" %[5]q {
		project_id = mongodbatlas_cloud_provider_access.%[1]s.project_id
	 }	 
	`
)

func TestAccConfigDSCloudProviderAccess_basic(t *testing.T) {
	var (
		suffix       = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
		name         = "datasource_test_role_basic" + suffix
		dataSCName   = "datasource_test_all_roles" + suffix
		resourceName = "mongodbatlas_cloud_provider_access." + name
		dsName       = "data.mongodbatlas_cloud_provider_access." + dataSCName
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	config := fmt.Sprintf(dataSourceProviderConfig, name, orgID, projectName, "AWS", dataSCName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "atlas_assumed_role_external_id"),
					resource.TestCheckResourceAttrSet(resourceName, "atlas_aws_account_arn"),
					resource.TestCheckResourceAttrSet(dsName, "aws_iam_roles.0.atlas_assumed_role_external_id"),
				),
			},
		},
	},
	)
}
