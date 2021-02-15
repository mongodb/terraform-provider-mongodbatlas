package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasDataLake_basic(t *testing.T) {
	SkipTestExtCred(t)
	var (
		resourceName        = "mongodbatlas_data_lake.basic_ds"
		orgID               = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName         = acctest.RandomWithPrefix("test-acc")
		name                = acctest.RandomWithPrefix("test-acc")
		roleID              = os.Getenv("AWS_ROLE_ID")
		testS3Bucket        = os.Getenv("AWS_S3_BUCKET")
		testS3BucketUpdated = os.Getenv("AWS_S3_BUCKET")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasDataLakeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataLakeConfig(projectName, orgID, name, roleID, testS3Bucket),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			{
				Config: testAccMongoDBAtlasDataLakeConfig(projectName, orgID, name, roleID, testS3BucketUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasDataLake_importBasic(t *testing.T) {
	SkipTestExtCred(t)
	var (
		resourceName = "mongodbatlas_data_lake.basic_ds"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		name         = acctest.RandomWithPrefix("test-acc")
		roleID       = os.Getenv("AWS_ROLE_ID")
		testS3Bucket = os.Getenv("AWS_S3_BUCKET")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasDataLakeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataLakeConfig(projectName, orgID, name, roleID, testS3Bucket),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasDataLakeImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMongoDBAtlasDataLakeImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := decodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s-%s", ids["project_id"], ids["name"]), nil
	}
}

func testAccCheckMongoDBAtlasDataLakeExists(resourceName string, dbUser *matlas.DataLake) resource.TestCheckFunc {
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

		if dbUserResp, _, err := conn.DataLakes.Get(context.Background(), ids["project_id"], ids["name"]); err == nil {
			*dbUser = *dbUserResp
			return nil
		}

		return fmt.Errorf("database user(%s) does not exist", ids["project_id"])
	}
}

func testAccCheckMongoDBAtlasDataLakeAttributes(dataLake *matlas.DataLake, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Printf("[DEBUG] difference dbUser.Username: %s , username : %s", dataLake.Name, name)
		if dataLake.Name != name {
			return fmt.Errorf("bad username: %s", dataLake.Name)
		}

		return nil
	}
}

func testAccCheckMongoDBAtlasDataLakeDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*matlas.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_data_lake" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)
		// Try to find the database user
		_, _, err := conn.DataLakes.Get(context.Background(), ids["project_id"], ids["name"])
		if err == nil {
			return fmt.Errorf("database user (%s) still exists", ids["project_id"])
		}
	}

	return nil
}

func testAccMongoDBAtlasDataLakeConfig(projectName, orgID, name, roleID, testS3Bucket string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = "%s"
			org_id = "%s"
		}
		resource "mongodbatlas_cloud_provider_access" "test" {
			project_id = mongodbatlas_project.test.id
			provider_name = "AWS"
			iam_assumed_role_arn = "%s"
		}

		resource "mongodbatlas_data_lake" "basic_ds" {
			project_id         = mongodbatlas_project.test.id
			name = "%s"
			aws{
				role_id = mongodbatlas_cloud_provider_access.test.role_id
				test_s3_bucket = "%s"
			}
		}
	`, projectName, orgID, roleID, name, testS3Bucket)
}
