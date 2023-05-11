resource "aws_iam_role_policy" "test_policy" {
  name = var.policy_name
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

resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
   project_id = var.project_id
   provider_name = "AWS"
}

resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
   project_id = var.project_id
   role_id =  mongodbatlas_cloud_provider_access_setup.setup_only.role_id
   aws {
      iam_assumed_role_arn = aws_iam_role.test_role.arn
   }
}


resource "aws_iam_role" "test_role" {
  name = var.role_name
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


resource "mongodbatlas_federated_database_instance" "test" {
  project_id         = mongodbatlas_project.test.id
  name = var.name
  aws {
    role_id = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
    test_s3_bucket = var.test_s3_bucket
  }
  storage_databases {
    name = "VirtualDatabase0"
    collections {
      name = "VirtualCollection0"
      data_sources {
          collection = var.collection
          database = var.database
          store_name =  var.cluster_name
      }
      data_sources {
          store_name = var.test_s3_bucket
          path = var.path
      }
    }
  }

  storage_stores {
    name = var.cluster_name
    cluster_name = var.cluster_name
    project_id = mongodbatlas_project.test.id
    provider = "atlas"
    read_preference {
    mode = "secondary"
    }
  }

  storage_stores {
    bucket = var.test_s3_bucket
    delimiter = "/"
    name = var.test_s3_bucket
    prefix = var.prefix
    provider = "s3"
    region = var.aws_region
  }

  storage_stores {
    name = "dataStore0"
    cluster_name = var.cluster_name
    project_id = mongodbatlas_project.test.id
    provider = "atlas"
    read_preference {
      mode = "secondary"
    }
  }
}