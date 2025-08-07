# S3 bucket for stream data
resource "aws_s3_bucket" "stream_bucket" {
  provider      = aws.s3_region
  bucket        = var.s3_bucket_name
  force_destroy = true
}

resource "aws_s3_bucket_versioning" "stream_bucket_versioning" {
  provider = aws.s3_region
  bucket   = aws_s3_bucket.stream_bucket.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "stream_bucket_encryption" {
  provider = aws.s3_region
  bucket   = aws_s3_bucket.stream_bucket.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

# PrivateLink for S3
resource "mongodbatlas_stream_privatelink_endpoint" "this" {
  project_id          = var.project_id
  provider_name       = "AWS"
  vendor              = "S3"
  region              = var.region
  service_endpoint_id = var.service_endpoint_id
}

output "privatelink_endpoint_id" {
  value = mongodbatlas_stream_privatelink_endpoint.this.id
}
