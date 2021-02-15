resource "mongodbatlas_project" "test" {
  name   = var.project_name
  org_id = var.org_id
}

resource "mongodbatlas_cloud_provider_access" "test" {
  project_id           = mongodbatlas_project.test.id
  provider_name        = "AWS"
  iam_assumed_role_arn = var.aws_iam_role_arn
}

resource "aws_iam_role_policy" "test_policy" {
  name = "test_policy"
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
  name = "test_role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": "${mongodbatlas_cloud_provider_access.test.atlas_aws_account_arn}"
      },
      "Action": "sts:AssumeRole",
      "Condition": {
        "StringEquals": {
          "sts:ExternalId": "${mongodbatlas_cloud_provider_access.test.atlas_assumed_role_external_id}"
        }
      }
    }
  ]
}
EOF


}

output "aws_iam_role_arn" {
  value = aws_iam_role.test_role.arn
}
output "cpa_role_id" {
  value = mongodbatlas_cloud_provider_access.test.role_id
}
output "project_id" {
  value = mongodbatlas_project.test.id
}