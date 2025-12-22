# AWS resources for log export
resource "aws_s3_bucket" "log_bucket" {
  bucket        = var.s3_bucket_name
  force_destroy = true
}

resource "aws_iam_role" "atlas_role" {
  name = var.aws_iam_role_name

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          AWS = mongodbatlas_cloud_provider_access_setup.setup_only.aws_config[0].atlas_aws_account_arn
        }
        Action = "sts:AssumeRole"
        Condition = {
          StringEquals = {
            "sts:ExternalId" = mongodbatlas_cloud_provider_access_setup.setup_only.aws_config[0].atlas_assumed_role_external_id
          }
        }
      }
    ]
  })
}

resource "aws_iam_role_policy" "atlas_role_policy" {
  name = "${var.aws_iam_role_name}-policy"
  role = aws_iam_role.atlas_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:PutObject",
          "s3:GetBucketLocation",
          "s3:ListBucket"
        ]
        Resource = [
          aws_s3_bucket.log_bucket.arn,
          "${aws_s3_bucket.log_bucket.arn}/*"
        ]
      }
    ]
  })
}

