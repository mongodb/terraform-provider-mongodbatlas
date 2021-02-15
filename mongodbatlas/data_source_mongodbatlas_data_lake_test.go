package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccDataSourceMongoDBAtlasDataLake_basic(t *testing.T) {
	SkipTestExtCred(t)
	var (
		resourceName = "mongodbatlas_data_lake.basic_ds"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		name         = acctest.RandomWithPrefix("test-acc")
		roleID       = os.Getenv("AWS_ROLE_ID")
		testS3Bucket = os.Getenv("AWS_S3_BUCKET")
		dataLake     = matlas.DataLake{}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataLakeDataSourceConfig(projectName, orgID, name, roleID, testS3Bucket),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasDataLakeExists(resourceName, &dataLake),
					testAccCheckMongoDBAtlasDataLakeAttributes(&dataLake, name),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
		},
	})
}

func testAccMongoDBAtlasDataLakeDataSourceConfig(projectName, orgID, name, roleID, testS3Bucket string) string {
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

		data "mongodbatlas_data_lake" "test" {
		  project_id           = mongodbatlas_data_lake.test.project_id
		  name = mongodbatlas_data_lake.test.name
		}
	`, projectName, orgID, roleID, name, testS3Bucket)
}
