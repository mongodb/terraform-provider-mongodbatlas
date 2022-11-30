package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccBackupRSDataLake_basic(t *testing.T) {
	SkipTestExtCred(t)
	var (
		resourceName        = "mongodbatlas_data_lake.test"
		orgID               = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName         = acctest.RandomWithPrefix("test-acc")
		name                = acctest.RandomWithPrefix("test-acc")
		policyName          = acctest.RandomWithPrefix("test-acc")
		roleName            = acctest.RandomWithPrefix("test-acc")
		testS3Bucket        = os.Getenv("AWS_S3_BUCKET")
		testS3BucketUpdated = os.Getenv("AWS_S3_BUCKET_UPDATED")
		dataLakeRegion      = "VIRGINIA_USA"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasDataLakeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataLakeConfig(policyName, roleName, projectName, orgID, name, testS3Bucket, dataLakeRegion, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			{
				Config: testAccMongoDBAtlasDataLakeConfig(policyName, roleName, projectName, orgID, name, testS3BucketUpdated, dataLakeRegion, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
		},
	})
}

func TestAccBackupRSDataLake_importBasic(t *testing.T) {
	SkipTestExtCred(t)
	var (
		resourceName = "mongodbatlas_data_lake.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		name         = acctest.RandomWithPrefix("test-acc")
		policyName   = acctest.RandomWithPrefix("test-acc")
		roleName     = acctest.RandomWithPrefix("test-acc")
		testS3Bucket = os.Getenv("AWS_S3_BUCKET")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasDataLakeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataLakeConfig(policyName, roleName, projectName, orgID, name, testS3Bucket, "", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasDataLakeImportStateIDFunc(resourceName, testS3Bucket),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMongoDBAtlasDataLakeImportStateIDFunc(resourceName, s3Bucket string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := decodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s--%s--%s", ids["project_id"], ids["name"], s3Bucket), nil
	}
}

func testAccCheckMongoDBAtlasDataLakeExists(resourceName string, dataLake *matlas.DataLake) resource.TestCheckFunc {
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

		if dataLakeResp, _, err := conn.DataLakes.Get(context.Background(), ids["project_id"], ids["name"]); err == nil {
			*dataLake = *dataLakeResp
			return nil
		}

		return fmt.Errorf("datalake (%s) does not exist", ids["project_id"])
	}
}

func testAccCheckMongoDBAtlasDataLakeAttributes(dataLake *matlas.DataLake, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Printf("[DEBUG] difference dataLake.Name: %s , username : %s", dataLake.Name, name)
		if dataLake.Name != name {
			return fmt.Errorf("bad datalake name: %s", dataLake.Name)
		}

		return nil
	}
}

func testAccCheckMongoDBAtlasDataLakeDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_data_lake" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)
		// Try to find the database user
		_, _, err := conn.DataLakes.Get(context.Background(), ids["project_id"], ids["name"])
		if err == nil {
			return fmt.Errorf("datalake (%s) still exists", ids["project_id"])
		}
	}

	return nil
}

func testAccMongoDBAtlasDataLakeConfig(policyName, roleName, projectName, orgID, name, testS3Bucket, dataLakeRegion string, isUpdate bool) string {
	stepDataLakeConfig := testAccMongoDBAtlasDataLakeConfigFirstStep(name, testS3Bucket)
	if isUpdate {
		stepDataLakeConfig = testAccMongoDBAtlasDataLakeConfigSecondStep(name, testS3Bucket, dataLakeRegion)
	}
	return fmt.Sprintf(`
resource "aws_iam_role_policy" "test_policy" {
  name = %[1]q
  role = aws_iam_role.test_role.id

  policy = <<-EOF
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
		"Action": "*",
		"Resource": "*"
      }
    ]
  }
  EOF
}

resource "aws_iam_role" "test_role" {
  name = %[2]q

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": "${mongodbatlas_cloud_provider_access_setup.setup_only.aws_config.0.atlas_aws_account_arn}"
      },
      "Action": "sts:AssumeRole",
      "Condition": {
        "StringEquals": {
          "sts:ExternalId": "${mongodbatlas_cloud_provider_access_setup.setup_only.aws_config.0.atlas_assumed_role_external_id}"
        }
      }
    }
  ]
}
EOF

}

resource "mongodbatlas_project" "test" {
   name   = %[3]q
   org_id = %[4]q
}


resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
   project_id = mongodbatlas_project.test.id
   provider_name = "AWS"
}

resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
   project_id = mongodbatlas_project.test.id
   role_id =  mongodbatlas_cloud_provider_access_setup.setup_only.role_id

   aws {
      iam_assumed_role_arn = aws_iam_role.test_role.arn
   }
}

%s
	`, policyName, roleName, projectName, orgID, stepDataLakeConfig)
}
func testAccMongoDBAtlasDataLakeConfigFirstStep(name, testS3Bucket string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_data_lake" "test" {
   project_id         = mongodbatlas_project.test.id
   name = %[1]q
   aws {
     role_id = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
     test_s3_bucket = %[2]q
   }
}
	`, name, testS3Bucket)
}
func testAccMongoDBAtlasDataLakeConfigSecondStep(name, testS3Bucket, dataLakeRegion string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_data_lake" "test" {
   project_id         = mongodbatlas_project.test.id
   name = %[1]q
   aws {
     role_id = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
     test_s3_bucket = %[2]q
   }
   data_process_region {
      cloud_provider = "AWS"
      region = %[3]q
   }
}
	`, name, testS3Bucket, dataLakeRegion)
}
