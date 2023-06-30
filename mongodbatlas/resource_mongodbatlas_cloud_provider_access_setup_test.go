package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	createProviderAccessSetupRole = `
	resource "mongodbatlas_project" "test" {
		name   = %[3]q
		org_id = %[2]q
	}
	resource "mongodbatlas_cloud_provider_access_setup" "%[1]s" {
		project_id = mongodbatlas_project.test.id
		provider_name = %[4]q
	 }

	`
)

func TestAccConfigRSCloudProviderAccessSetup_basic(t *testing.T) {
	var (
		name         = "test_basic" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
		resourceName = "mongodbatlas_cloud_provider_access_setup." + name
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		targetRole   = matlas.AWSIAMRole{}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		// same as regular cloud provider access resource
		CheckDestroy: testAccCheckMongoDBAtlasProviderAccessDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(createProviderAccessSetupRole, name, orgID, projectName, "AWS"),
				Check: resource.ComposeTestCheckFunc(
					// same as regular cloud resource
					testAccCheckMongoDBAtlasProviderAccessExists(resourceName, &targetRole),
					resource.TestCheckResourceAttrSet(resourceName, "aws.atlas_assumed_role_external_id"),
					resource.TestCheckResourceAttrSet(resourceName, "aws.atlas_aws_account_arn"),
				),
			},
		},
	},
	)
}

func TestAccConfigRSCloudProviderAccessSetup_importBasic(t *testing.T) {
	var (
		name         = "test_basic" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
		resourceName = "mongodbatlas_cloud_provider_access_setup." + name
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		targetRole   = matlas.AWSIAMRole{}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasProviderAccessDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(createProviderAccessSetupRole, name, orgID, projectName, "AWS"),
				Check: resource.ComposeTestCheckFunc(
					// same as regular cloud provider because we are just checking in the api
					testAccCheckMongoDBAtlasProviderAccessExists(resourceName, &targetRole),
					resource.TestCheckResourceAttrSet(resourceName, "aws.atlas_assumed_role_external_id"),
					resource.TestCheckResourceAttrSet(resourceName, "aws.atlas_aws_account_arn"),
				),
			},
			{
				ResourceName: resourceName,
				// ID remains the same project-id, provider-name and id for consistency
				ImportStateIdFunc: testAccCheckMongoDBAtlasCloudProviderAccessImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	},
	)
}
