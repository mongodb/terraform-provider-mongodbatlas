package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// covers both resource creation and datasource creation for setup

const (
	dataSourceCPASProviderConfig = `
	resource "mongodbatlas_project" "test" {
		name   = %[3]q
		org_id = %[2]q
	}
	resource "mongodbatlas_cloud_provider_access_setup" "%[1]s" {
		project_id = mongodbatlas_project.test.id
		provider_name = %[4]q
	 }
	 
	 data "mongodbatlas_cloud_provider_access_setup" "%[5]s" {
		project_id = mongodbatlas_cloud_provider_access_setup.%[1]s.project_id
		provider_name = %[4]q
		role_id =  mongodbatlas_cloud_provider_access_setup.%[1]s.role_id
	 }
	 `
)

func TestAccConfigDSCloudProviderAccessSetup_aws_basic(t *testing.T) {
	var (
		suffix      = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
		name        = "cpas" + suffix
		dataSCName  = "ds_cpas" + suffix
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acctest.RandomWithPrefix("test-acc")

		// resources fqdn name
		resourceName = "mongodbatlas_cloud_provider_access_setup." + name
		dsName       = "data.mongodbatlas_cloud_provider_access_setup." + dataSCName
	)

	config := fmt.Sprintf(dataSourceCPASProviderConfig, name, orgID, projectName, "AWS", dataSCName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "created_date"),
					resource.TestCheckResourceAttrSet(resourceName, "role_id"),
					resource.TestCheckResourceAttrSet(dsName, "aws.atlas_assumed_role_external_id"),
				),
			},
		},
	})
}
