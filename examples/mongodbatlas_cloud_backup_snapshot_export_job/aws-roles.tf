resource "aws_iam_role_policy" "test_policy" {
  name = "mongo_setup_policy"
  role = aws_iam_role.test_role.id

  policy = <<-EOF
  {
      "Version": "2012-10-17",
      "Statement": [
      {
          "Effect": "Allow",
          "Action": "s3:GetBucketLocation",
          "Resource": "arn:aws:s3:::${aws_s3_bucket.test_bucket.bucket}"
      },
      {
          "Effect": "Allow",
          "Action": "s3:PutObject",
          "Resource": "arn:aws:s3:::${aws_s3_bucket.test_bucket.bucket}/*"
      }]
  }
  EOF
}

resource "aws_iam_role" "test_role" {
  name = "mongo_setup_test_role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": "${mongodbatlas_cloud_provider_access_setup.setup_only.aws_config[0].atlas_aws_account_arn}"
      },
      "Action": "sts:AssumeRole",
      "Condition": {
        "StringEquals": {
          "sts:ExternalId": "${mongodbatlas_cloud_provider_access_setup.setup_only.aws_config[0].atlas_assumed_role_external_id}"
        }
      }
    }
  ]
}
EOF
}
