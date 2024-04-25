# Create IAM role & policy to authorize with Atlas
resource "aws_iam_role_policy" "test_policy" {
  name = var.aws_iam_role_policy_name
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
  name                 = var.aws_iam_role_name
  max_session_duration = 43200

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

# Create S3 buckets
resource "aws_s3_bucket" "log_bucket" {
  bucket        = var.s3_bucket_name
  force_destroy = true # required for destroying as Atlas may create a test folder in the bucket when push-based log export is set up 
}

# Add authorization policy to existing IAM role
resource "aws_iam_role_policy" "s3_bucket_policy" {
  name = var.s3_bucket_policy_name
  role = aws_iam_role.test_role.id

  policy = <<-EOF
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Action": [
          "s3:ListBucket",
          "s3:PutObject",
          "s3:GetObject",
          "s3:GetBucketLocation"
        ],
        "Resource": [
          "${aws_s3_bucket.log_bucket.arn}",
          "${aws_s3_bucket.log_bucket.arn}/*"
        ]
      }
    ]
  }
  EOF
}
