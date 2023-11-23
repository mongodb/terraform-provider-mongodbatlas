package mongodbatlas_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccConfigRSCloudProviderAccessAWS_basic(t *testing.T) {
	var (
		resourceName   = "mongodbatlas_cloud_provider_access.test"
		dataSourceName = "data.mongodbatlas_cloud_provider_access.test"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acctest.RandomWithPrefix("test-acc")
		targetRole     = matlas.CloudProviderAccessRole{}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasic(t) },
		ProtoV6ProviderFactories: testAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudProviderAccessAWS(orgID, projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProviderAccessExists(resourceName, &targetRole),
					resource.TestCheckResourceAttrSet(resourceName, "atlas_assumed_role_external_id"),
					resource.TestCheckResourceAttrSet(resourceName, "atlas_aws_account_arn"),
					resource.TestCheckResourceAttrSet(dataSourceName, "aws_iam_roles.0.atlas_assumed_role_external_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "aws_iam_roles.0.atlas_aws_account_arn"),
					resource.TestCheckResourceAttrSet(dataSourceName, "aws_iam_roles.0.created_date"),
					resource.TestCheckResourceAttrSet(dataSourceName, "aws_iam_roles.0.provider_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "aws_iam_roles.0.role_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "aws_iam_roles.0.created_date"),
				),
			},
		},
	},
	)
}

func TestAccConfigRSCloudProviderAccess_importBasic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_cloud_provider_access.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		targetRole   = matlas.CloudProviderAccessRole{}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasic(t) },
		ProtoV6ProviderFactories: testAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudProviderAccessAWS(orgID, projectName),
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
	conn := testAccProviderSdkV2.Meta().(*MongoDBClient).Atlas
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_cloud_provider_access" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)

		roles, _, err := conn.CloudProviderAccess.ListRoles(context.Background(), ids["project_id"])

		if err != nil {
			return fmt.Errorf(errorGetRead, err)
		}

		var targetRole matlas.CloudProviderAccessRole

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

func testAccCheckMongoDBAtlasProviderAccessExists(resourceName string, targetRole *matlas.CloudProviderAccessRole) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProviderSdkV2.Meta().(*MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.Attributes["project_id"] == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := decodeStateID(rs.Primary.ID)
		providerName := ids["provider_name"]
		id := ids["id"]

		roles, _, err := conn.CloudProviderAccess.ListRoles(context.Background(), ids["project_id"])

		if err != nil {
			return fmt.Errorf(errorGetRead, err)
		}

		if providerName == "AWS" {
			for i := range roles.AWSIAMRoles {
				if roles.AWSIAMRoles[i].RoleID == id && roles.AWSIAMRoles[i].ProviderName == providerName {
					*targetRole = roles.AWSIAMRoles[i]
					return nil
				}
			}
		}

		if providerName == "AZURE" {
			for i := range roles.AzureServicePrincipals {
				if *roles.AzureServicePrincipals[i].AzureID == id && roles.AzureServicePrincipals[i].ProviderName == providerName {
					*targetRole = roles.AzureServicePrincipals[i]
					return nil
				}
			}
		}

		return fmt.Errorf("error cloud Provider Access (%s) does not exist", ids["project_id"])
	}
}

func testAccMongoDBAtlasCloudProviderAccessAWS(orgID, projectName string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_project" "test" {
		name   = %[2]q
		org_id = %[1]q
	}
	resource "mongodbatlas_cloud_provider_access" "test" {
		project_id = mongodbatlas_project.test.id
		provider_name = "AWS"
	 }

	 data "mongodbatlas_cloud_provider_access" "test" {
		project_id = mongodbatlas_cloud_provider_access.test.project_id
	 }	 
	`, orgID, projectName)
}
