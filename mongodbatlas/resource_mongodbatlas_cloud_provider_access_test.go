package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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

func TestAccResourceMongoDBAtlasCloudProviderAccess_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_cloud_provider_access.test_basic_" + acctest.RandString(10)
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		targetRole   = matlas.AWSIAMRole{}
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasProviderAccessDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(createProviderAccessRole, resourceName, projectID, "AWS"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProviderAccessExists(resourceName, &targetRole),
					resource.TestCheckResourceAttr(resourceName, "atlas_assumed_role_external_id", targetRole.AtlasAssumedRoleExternalID),
					resource.TestCheckResourceAttr(resourceName, "atlas_aws_account_arn", targetRole.AtlasAWSAccountARN),
				),
			},
		},
	},
	)
}

func testAccCheckMongoDBAtlasProviderAccessDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*matlas.Client)
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
		conn := testAccProvider.Meta().(*matlas.Client)

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
			role := &(roles.AWSIAMRoles[i])

			if role.RoleID == ids["id"] && role.ProviderName == ids["provider_name"] {
				*targetRole = *role
				return nil
			}
		}

		return fmt.Errorf("error cloud Provider Access (%s) does not exist", ids["project_id"])
	}
}
