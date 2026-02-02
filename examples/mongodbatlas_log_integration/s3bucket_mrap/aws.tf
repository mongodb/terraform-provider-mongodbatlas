# Generate a random suffix for globally unique S3 bucket names
resource "random_id" "bucket_suffix" {
  byte_length = 4
}

# S3 buckets (one per region for MRAP)
resource "aws_s3_bucket" "mrap_bucket_us_east_1" {
  provider      = aws.us_east_1
  bucket        = "${var.name_prefix}-${random_id.bucket_suffix.hex}-us-east-1"
  force_destroy = true

  tags = {
    Purpose = "Atlas Log Integration MRAP"
    Region  = "us-east-1"
  }
}

resource "aws_s3_bucket" "mrap_bucket_us_west_2" {
  provider      = aws.us_west_2
  bucket        = "${var.name_prefix}-${random_id.bucket_suffix.hex}-us-west-2"
  force_destroy = true

  tags = {
    Purpose = "Atlas Log Integration MRAP"
    Region  = "us-west-2"
  }
}

# S3 Multi-Region Access Point
resource "aws_s3control_multi_region_access_point" "atlas_logs" {
  details {
    name = "${var.name_prefix}-mrap-${random_id.bucket_suffix.hex}"

    region {
      bucket = aws_s3_bucket.mrap_bucket_us_east_1.id
    }

    region {
      bucket = aws_s3_bucket.mrap_bucket_us_west_2.id
    }
  }
}

# IAM role for Atlas to assume
resource "aws_iam_role" "atlas_role" {
  name                 = "${var.name_prefix}-role-${random_id.bucket_suffix.hex}"
  max_session_duration = 43200

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          AWS = mongodbatlas_cloud_provider_access_setup.setup.aws_config[0].atlas_aws_account_arn
        }
        Action = "sts:AssumeRole"
        Condition = {
          StringEquals = {
            "sts:ExternalId" = mongodbatlas_cloud_provider_access_setup.setup.aws_config[0].atlas_assumed_role_external_id
          }
        }
      }
    ]
  })

  tags = {
    Purpose = "Atlas Log Integration MRAP"
  }
}

# IAM policy for MRAP access
resource "aws_iam_role_policy" "mrap_policy" {
  name = "${var.name_prefix}-mrap-policy"
  role = aws_iam_role.atlas_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "MRAPAccess"
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:ListBucket",
          "s3:GetBucketLocation",
          "s3:DeleteObject"
        ]
        # MRAP ARN format: arn:aws:s3::ACCOUNT_ID:accesspoint/MRAP_ALIAS
        Resource = [
          aws_s3control_multi_region_access_point.atlas_logs.arn,
          "${aws_s3control_multi_region_access_point.atlas_logs.arn}/object/*"
        ]
      },
      {
        Sid    = "MRAPMetadata"
        Effect = "Allow"
        Action = [
          "s3:GetAccessPoint",
          "s3:ListAccessPoints"
        ]
        Resource = "*"
      },
      {
        # Direct bucket access is also needed for some operations
        Sid    = "BackingBucketAccess"
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:ListBucket",
          "s3:GetBucketLocation",
          "s3:DeleteObject"
        ]
        Resource = [
          aws_s3_bucket.mrap_bucket_us_east_1.arn,
          "${aws_s3_bucket.mrap_bucket_us_east_1.arn}/*",
          aws_s3_bucket.mrap_bucket_us_west_2.arn,
          "${aws_s3_bucket.mrap_bucket_us_west_2.arn}/*"
        ]
      }
    ]
  })
}
