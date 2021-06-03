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
		resourceName = "mongodbatlas_data_lake.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		name         = acctest.RandomWithPrefix("test-acc")
		policyName   = acctest.RandomWithPrefix("test-acc")
		roleName     = acctest.RandomWithPrefix("test-acc")
		testS3Bucket = os.Getenv("AWS_S3_BUCKET")
		dataLake     = matlas.DataLake{}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataLakeDataSourceConfig(policyName, roleName, projectName, orgID, name, testS3Bucket),
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

func testAccMongoDBAtlasDataLakeDataSourceConfig(policyName, roleName, projectName, orgID, name, testS3Bucket string) string {
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
        "AWS": "${mongodbatlas_cloud_provider_access_setup.setup_only.aws.atlas_aws_account_arn}"
      },
      "Action": "sts:AssumeRole",
      "Condition": {
        "StringEquals": {
          "sts:ExternalId": "${mongodbatlas_cloud_provider_access_setup.setup_only.aws.atlas_assumed_role_external_id}"
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

   aws = {
      iam_assumed_role_arn = aws_iam_role.test_role.arn
   }
}

resource "mongodbatlas_data_lake" "test" {
   project_id         = mongodbatlas_project.test.id
   name = %[5]q
   aws_role_id = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
   aws_test_s3_bucket = %[6]q
}

data "mongodbatlas_data_lake" "test" {
  project_id           = mongodbatlas_data_lake.test.project_id
  name = mongodbatlas_data_lake.test.name
}
	`, policyName, roleName, projectName, orgID, name, testS3Bucket)
}
