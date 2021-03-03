resource "mongodbatlas_cloud_provider_access" "test" {
  project_id           = var.project_id
  provider_name        = "AWS"
  
  #(Optional) Since we update the `iam_assumed_role_arn` resource using an HTTP call and not by the `mongodbatlas_cloud_provider_access` resource argument, 
  #the lifecycle argument was added so that terraform would ignore changes of the `iam_assumed_role_arn` argument in future terraform applies.
  lifecycle {
    ignore_changes = [
      iam_assumed_role_arn
    ]
  }
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

# The null resource updates the `mongodbatlas_cloud_provider_access` resource with the correct IAM role ARN using an API HTTP PATCH request.
# sleep 10 - Waits ten seconds to make sure that all AWS servers are updated with the new IAM Role.
resource "null_resource" "link_role_arn_to_cloud_provider_access" {
  provisioner "local-exec" {
      command = <<EOT
      sleep 10;
      curl --user "${var.public_key}:${var.private_key}" -X PATCH --digest \
            --header "Accept: application/json" \
            --header "Content-Type: application/json" \
            "https://cloud.mongodb.com/api/atlas/v1.0/groups/${var.project_id}/cloudProviderAccess/${mongodbatlas_cloud_provider_access.test.role_id}?pretty=true" \
            --data '{ "providerName": "AWS", "iamAssumedRoleArn" : "${aws_iam_role.test_role.arn}" }'

EOT
  }
}


output "cpa_role_id" {
  value = mongodbatlas_cloud_provider_access.test.role_id
}
