package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// covers both resource creation and datasource creation for setup

const (
	dataSourceCPASProviderConfig = `
	resource "mongodbatlas_cloud_provider_access_setup" "%[1]s" {
		project_id = "%[2]s"
		provider_name = "%[3]s"
	 }
	 
	 data "mongodbatlas_cloud_provider_access_setup" "%[4]s" {
		project_id = mongodbatlas_cloud_provider_access_setup.%[1]s.project_id
		provider_name = "%[3]s"
		role_id =  mongodbatlas_cloud_provider_access_setup.%[1]s.role_id
	 }
	 `
)

func TestAccdataSourceMongoDBAtlasCloudProviderAccessSetup_aws_basic(t *testing.T) {
	var (
		suffix     = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
		name       = "cpas" + suffix
		dataSCName = "ds_cpas" + suffix
		projectID  = os.Getenv("MONGODB_ATLAS_PROJECT_ID")

		// resources fqdn name
		resourceName = "mongodbatlas_cloud_provider_access_setup." + name
		dsName       = "data.mongodbatlas_cloud_provider_access_setup." + dataSCName
	)

	config := fmt.Sprintf(dataSourceCPASProviderConfig, name, projectID, "AWS", dataSCName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "created_date"),
					resource.TestCheckResourceAttrSet(resourceName, "role_id"),
					resource.TestCheckResourceAttr(dsName, "project_id", projectID),
					resource.TestCheckResourceAttrSet(dsName, "aws.atlas_assumed_role_external_id"),
				),
			},
		},
	})
}
