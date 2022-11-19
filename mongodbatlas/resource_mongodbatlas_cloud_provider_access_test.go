package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	createProviderAccessRole = `
	resource "mongodbatlas_cloud_provider_access" "%[1]s" {
		project_id = "%[2]s"
		provider_name = "%[3]s"
	 }

	`
)

func TestAccConfigRSCloudProviderAccess_basic(t *testing.T) {
	var (
		name         = "test_basic" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
		resourceName = "mongodbatlas_cloud_provider_access." + name
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		targetRole   = matlas.AWSIAMRole{}
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasProviderAccessDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(createProviderAccessRole, name, projectID, "AWS"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProviderAccessExists(resourceName, &targetRole),
					resource.TestCheckResourceAttrSet(resourceName, "atlas_assumed_role_external_id"),
					resource.TestCheckResourceAttrSet(resourceName, "atlas_aws_account_arn"),
				),
			},
		},
	},
	)
}

func TestAccConfigRSCloudProviderAccess_importBasic(t *testing.T) {
	var (
		name         = "test_basic" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
		resourceName = "mongodbatlas_cloud_provider_access." + name
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		targetRole   = matlas.AWSIAMRole{}
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasProviderAccessDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(createProviderAccessRole, name, projectID, "AWS"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProviderAccessExists(resourceName, &targetRole),
					resource.TestCheckResourceAttrSet(resourceName, "atlas_assumed_role_external_id"),
					resource.TestCheckResourceAttrSet(resourceName, "atlas_aws_account_arn"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasCloudProviderAccessImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	},
	)
}

func testAccCheckMongoDBAtlasCloudProviderAccessImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := decodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s-%s-%s", ids["project_id"], ids["provider_name"], ids["id"]), nil
	}
}

func testAccCheckMongoDBAtlasProviderAccessDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_cloud_provider_access" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)

		roles, _, err := conn.CloudProviderAccess.ListRoles(context.Background(), ids["project_id"])

		if err != nil {
			return fmt.Errorf(errorGetRead, err)
		}

		var targetRole matlas.AWSIAMRole

		// searching in roles
		for i := range roles.AWSIAMRoles {
			role := &(roles.AWSIAMRoles[i])

			if role.RoleID == ids["id"] && role.ProviderName == ids["provider_name"] {
				targetRole = *role
			}
		}

		//  Found !!
		if targetRole.RoleID != "" {
			return fmt.Errorf("error cloud Provider Access Role (%s) still exists", ids["id"])
		}
	}

	return nil
}

func testAccCheckMongoDBAtlasProviderAccessExists(resourceName string, targetRole *matlas.AWSIAMRole) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.Attributes["project_id"] == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := decodeStateID(rs.Primary.ID)

		roles, _, err := conn.CloudProviderAccess.ListRoles(context.Background(), ids["project_id"])

		if err != nil {
			return fmt.Errorf(errorGetRead, err)
		}

		// searching in roles
		for i := range roles.AWSIAMRoles {
			if roles.AWSIAMRoles[i].RoleID == ids["id"] && roles.AWSIAMRoles[i].ProviderName == ids["provider_name"] {
				*targetRole = roles.AWSIAMRoles[i]
				return nil
			}
		}

		return fmt.Errorf("error cloud Provider Access (%s) does not exist", ids["project_id"])
	}
}
