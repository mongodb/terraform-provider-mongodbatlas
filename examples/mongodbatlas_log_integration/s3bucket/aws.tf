# S3 bucket for storing Atlas logs
resource "aws_s3_bucket" "log_bucket" {
  bucket        = var.s3_bucket_name
  force_destroy = true
}

# IAM role for Atlas to assume
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

# IAM policy for S3 access
resource "aws_iam_role_policy" "atlas_s3_policy" {
  name = "${var.aws_iam_role_name}-s3-policy"
  role = aws_iam_role.atlas_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:PutObject",
          "s3:GetObject",
          "s3:ListBucket",
          "s3:DeleteObject",
          "s3:GetBucketLocation"
        ]
        Resource = [
          aws_s3_bucket.log_bucket.arn,
          "${aws_s3_bucket.log_bucket.arn}/*"
        ]
      }
    ]
  })
}
