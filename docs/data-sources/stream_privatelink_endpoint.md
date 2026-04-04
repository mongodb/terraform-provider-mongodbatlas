---
subcategory: "Streams"
---

# Data Source: mongodbatlas_stream_privatelink_endpoint

`mongodbatlas_stream_privatelink_endpoint` describes a Privatelink Endpoint for Streams.

## Example Usages

### AWS Confluent Privatelink
```terraform
resource "confluent_environment" "staging" {
  display_name = "Staging"
}

resource "confluent_network" "private_link" {
  display_name     = "terraform-test-private-link-network-manual"
  cloud            = "AWS"
  region           = var.aws_region
  connection_types = ["PRIVATELINK"]
  zones            = keys(var.subnets_to_privatelink)
  environment {
    id = confluent_environment.staging.id
  }
  dns_config {
    resolution = "PRIVATE"
  }
}

resource "confluent_private_link_access" "aws" {
  display_name = "example-private-link-access"
  aws {
    account = var.aws_account_id
  }
  environment {
    id = confluent_environment.staging.id
  }
  network {
    id = confluent_network.private_link.id
  }
}

resource "confluent_kafka_cluster" "dedicated" {
  display_name = "example-dedicated-cluster"
  availability = "MULTI_ZONE"
  cloud        = confluent_network.private_link.cloud
  region       = confluent_network.private_link.region
  dedicated {
    cku = 2
  }
  environment {
    id = confluent_environment.staging.id
  }
  network {
    id = confluent_network.private_link.id
  }
}

resource "mongodbatlas_stream_privatelink_endpoint" "test" {
  project_id          = var.project_id
  dns_domain          = confluent_network.private_link.dns_domain
  provider_name       = "AWS"
  region              = var.aws_region
  vendor              = "CONFLUENT"
  service_endpoint_id = confluent_network.private_link.aws[0].private_link_endpoint_service
  dns_sub_domain      = confluent_network.private_link.zonal_subdomains
}

data "mongodbatlas_stream_privatelink_endpoint" "singular_datasource" {
  project_id = var.project_id
  id         = mongodbatlas_stream_privatelink_endpoint.test.id
}

data "mongodbatlas_stream_privatelink_endpoints" "plural_datasource" {
  project_id = var.project_id
}

output "interface_endpoint_id" {
  value = data.mongodbatlas_stream_privatelink_endpoint.singular_datasource.interface_endpoint_id
}

output "interface_endpoint_ids" {
  value = data.mongodbatlas_stream_privatelink_endpoints.plural_datasource.results[*].interface_endpoint_id
}
```

### AWS MSK Privatelink
```terraform
resource "aws_vpc" "vpc" {
  cidr_block = "192.168.0.0/22"
}

data "aws_availability_zones" "azs" {
  state = "available"
}

resource "aws_subnet" "subnet_az1" {
  availability_zone = data.aws_availability_zones.azs.names[0]
  cidr_block        = "192.168.0.0/24"
  vpc_id            = aws_vpc.vpc.id
}

resource "aws_subnet" "subnet_az2" {
  availability_zone = data.aws_availability_zones.azs.names[1]
  cidr_block        = "192.168.1.0/24"
  vpc_id            = aws_vpc.vpc.id
}

resource "aws_security_group" "sg" {
  vpc_id = aws_vpc.vpc.id
}

resource "aws_msk_cluster" "example" {
  cluster_name           = var.msk_cluster_name
  kafka_version          = "3.6.0"
  number_of_broker_nodes = 2

  broker_node_group_info {
    instance_type = "kafka.m5.large"
    client_subnets = [
      aws_subnet.subnet_az1.id,
      aws_subnet.subnet_az2.id,
    ]
    security_groups = [aws_security_group.sg.id]

    connectivity_info {
      vpc_connectivity {
        client_authentication {
          sasl {
            scram = true
          }
        }
      }
    }
  }

  client_authentication {
    sasl {
      scram = true
    }
  }

  configuration_info {
    arn      = aws_msk_configuration.example.arn
    revision = aws_msk_configuration.example.latest_revision
  }
}

resource "aws_msk_cluster_policy" "example" {
  cluster_arn = aws_msk_cluster.example.arn

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Effect = "Allow"
      Principal = {
        "AWS" = "arn:aws:iam::${var.aws_account_id}:root"
      }
      Action = [
        "kafka:CreateVpcConnection",
        "kafka:GetBootstrapBrokers",
        "kafka:DescribeCluster",
        "kafka:DescribeClusterV2"
      ]
      Resource = aws_msk_cluster.example.arn
    }]
  })
}

resource "aws_msk_single_scram_secret_association" "example" {
  cluster_arn = aws_msk_cluster.example.arn
  secret_arn  = var.aws_secret_arn
}

resource "aws_msk_configuration" "example" {
  name = "${var.msk_cluster_name}-msk-configuration"

  # Default ASW MSK configuration with "allow.everyone.if.no.acl.found=false" added
  server_properties = <<PROPERTIES
auto.create.topics.enable=false
default.replication.factor=3
min.insync.replicas=2
num.io.threads=8
num.network.threads=5
num.partitions=1
num.replica.fetchers=2
replica.lag.time.max.ms=30000
socket.receive.buffer.bytes=102400
socket.request.max.bytes=104857600
socket.send.buffer.bytes=102400
unclean.leader.election.enable=true
allow.everyone.if.no.acl.found=false
PROPERTIES
}

resource "mongodbatlas_stream_privatelink_endpoint" "test" {
  project_id    = var.project_id
  provider_name = "AWS"
  vendor        = "MSK"
  arn           = aws_msk_cluster.example.arn
}

data "mongodbatlas_stream_privatelink_endpoint" "singular_datasource" {
  project_id = var.project_id
  id         = mongodbatlas_stream_privatelink_endpoint.test.id
}

output "privatelink_endpoint_id" {
  value = data.mongodbatlas_stream_privatelink_endpoint.singular_datasource.id
}
```

### AWS S3 Privatelink
```terraform
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
```

### GCP Confluent Privatelink
```terraform
resource "mongodbatlas_stream_privatelink_endpoint" "gcp_confluent" {
  project_id = var.project_id

  provider_name = "GCP"
  vendor        = "CONFLUENT"
  region        = var.gcp_region

  dns_domain     = var.confluent_dns_domain
  dns_sub_domain = var.confluent_dns_subdomains

  service_attachment_uris = [
    "projects/my-project/regions/us-west1/serviceAttachments/confluent-attachment-1",
    "projects/my-project/regions/us-west1/serviceAttachments/confluent-attachment-2"
  ]
}

data "mongodbatlas_stream_privatelink_endpoint" "gcp_confluent" {
  project_id = var.project_id
  id         = mongodbatlas_stream_privatelink_endpoint.gcp_confluent.id
}

output "privatelink_endpoint_id" {
  description = "The ID of the MongoDB Atlas Stream Private Link Endpoint"
  value       = mongodbatlas_stream_privatelink_endpoint.gcp_confluent.id
}

output "privatelink_endpoint_state" {
  description = "The state of the MongoDB Atlas Stream Private Link Endpoint"
  value       = data.mongodbatlas_stream_privatelink_endpoint.gcp_confluent.state
}

output "service_attachment_uris" {
  description = "The GCP service attachment URIs used for the private link"
  value       = mongodbatlas_stream_privatelink_endpoint.gcp_confluent.service_attachment_uris
}
```

### GCP Pub/Sub Private Service Connect
```terraform
resource "mongodbatlas_advanced_cluster" "cluster" {
  project_id   = var.project_id
  name         = var.cluster_name
  cluster_type = "REPLICASET"
  replication_specs = [{
    region_configs = [{
      priority      = 7
      provider_name = "GCP"
      region_name   = "US_EAST_4"
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
    }]
  }]
}

resource "mongodbatlas_stream_privatelink_endpoint" "gcp_pubsub" {
  project_id = var.project_id

  provider_name = "GCP"
  vendor        = "PUBSUB"
  region        = var.gcp_region
  depends_on    = [mongodbatlas_advanced_cluster.cluster]
}

data "mongodbatlas_stream_privatelink_endpoint" "gcp_pubsub" {
  project_id = var.project_id
  id         = mongodbatlas_stream_privatelink_endpoint.gcp_pubsub.id
}

output "privatelink_endpoint_id" {
  description = "The ID of the MongoDB Atlas Stream Private Link Endpoint"
  value       = mongodbatlas_stream_privatelink_endpoint.gcp_pubsub.id
}

output "privatelink_endpoint_state" {
  description = "The state of the MongoDB Atlas Stream Private Link Endpoint"
  value       = data.mongodbatlas_stream_privatelink_endpoint.gcp_pubsub.state
}

output "dns_domain" {
  description = "The DNS domain computed by the API for the GCP Pub/Sub private link"
  value       = mongodbatlas_stream_privatelink_endpoint.gcp_pubsub.dns_domain
}
```

### Azure Blob Storage Privatelink
```terraform
resource "mongodbatlas_stream_privatelink_endpoint" "azure_blob" {
  project_id          = var.project_id
  provider_name       = "AZURE"
  vendor              = "AZURE_BLOB_STORAGE"
  region              = var.azure_region
  service_endpoint_id = "/subscriptions/${var.azure_subscription_id}/resourceGroups/${var.azure_resource_group}/providers/Microsoft.Storage/storageAccounts/${var.azure_storage_account_name}"
  dns_domain          = "${var.azure_storage_account_name}.blob.core.windows.net"
}

data "mongodbatlas_stream_privatelink_endpoint" "azure_blob" {
  project_id = var.project_id
  id         = mongodbatlas_stream_privatelink_endpoint.azure_blob.id
}

output "privatelink_endpoint_id" {
  value = mongodbatlas_stream_privatelink_endpoint.azure_blob.id
}

output "privatelink_endpoint_state" {
  value = data.mongodbatlas_stream_privatelink_endpoint.azure_blob.state
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) The ID of the Private Link connection.
- `project_id` (String) Unique 24-hexadecimal digit string that identifies your project, also known as `groupId` in the official documentation.

### Read-Only

- `arn` (String) Amazon Resource Name (ARN). Required for AWS Provider and MSK vendor.
- `dns_domain` (String) The domain hostname. Required for the following provider and vendor combinations:

	* AWS provider with CONFLUENT vendor.

	* AZURE provider with EVENTHUB, CONFLUENT, or AZURE_BLOB_STORAGE vendor. For AZURE_BLOB_STORAGE, this should follow the format '{storageAccount}.blob.core.windows.net'.

	* For GCP provider with PUBSUB vendor, the API computes this process.
- `dns_sub_domain` (List of String) Sub-Domain name of Confluent cluster. These are typically your availability zones. Required for AWS Provider and CONFLUENT vendor. If your AWS CONFLUENT cluster doesn't use subdomains, you must set this to the empty array [].
- `error_message` (String) Error message if the connection is in a failed state.
- `interface_endpoint_id` (String) Interface endpoint ID that is created from the specified service endpoint ID.
- `interface_endpoint_name` (String) Name of interface endpoint that is created from the specified service endpoint ID.
- `provider_account_id` (String) Account ID from the cloud provider.
- `provider_name` (String) Provider where the endpoint is deployed. Valid values are AWS, AZURE, and GCP.
- `region` (String) The region of the Provider’s cluster. See [AZURE](https://www.mongodb.com/docs/atlas/reference/microsoft-azure/#stream-processing-instances) and [AWS](https://www.mongodb.com/docs/atlas/reference/amazon-aws/#stream-processing-instances) supported regions. When the vendor is `CONFLUENT`, this is the domain name of Confluent cluster. When the vendor is `MSK`, this is computed by the API from the provided `arn`.
- `service_attachment_uris` (List of String) List of GCP service attachment URIs for Confluent vendor. Required for GCP provider with CONFLUENT vendor.
- `service_endpoint_id` (String) For AZURE EVENTHUB, this is the [namespace endpoint ID](https://learn.microsoft.com/en-us/rest/api/eventhub/namespaces/get). For AWS CONFLUENT cluster, this is the [VPC Endpoint service name](https://docs.confluent.io/cloud/current/networking/private-links/aws-privatelink.html). For AZURE AZURE_BLOB_STORAGE, this is the Azure Resource Manager path of the storage account in the format '/subscriptions/{subscriptionId}/resourceGroups/{resourceGroup}/providers/Microsoft.Storage/storageAccounts/{storageAccount}'.
- `state` (String) Status of the connection.
- `vendor` (String) Vendor that manages the endpoint. The following are the vendor values per provider:

	* **AWS**: MSK, CONFLUENT, and S3

	* **Azure**: EVENTHUB, CONFLUENT, and AZURE_BLOB_STORAGE

	* **GCP**: CONFLUENT and PUBSUB

For more information see: [MongoDB Atlas API - Streams Privatelink](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/operation/operation-createprivatelinkconnection) Documentation.
