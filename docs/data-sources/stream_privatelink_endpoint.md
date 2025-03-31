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

### Azure Privatelink
```terraform
resource "azurerm_resource_group" "rg" {
  name     = var.azure_resource_group
  location = var.azure_region
}

resource "azurerm_virtual_network" "vnet" {
  name                = var.vnet_name
  address_space       = var.vnet_address_space
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
}

resource "azurerm_subnet" "subnet" {
  name                 = var.subnet_name
  resource_group_name  = azurerm_resource_group.rg.name
  virtual_network_name = azurerm_virtual_network.vnet.name
  address_prefixes     = var.subnet_address_prefix
}

resource "azurerm_eventhub_namespace" "eventhub_ns" {
  name = var.eventhub_namespace_name
  location = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
  sku = "Standard" # Minimum SKU for Private Link
  capacity = 1
}

resource "azurerm_eventhub" "eventhub" {
  name                = var.eventhub_name
  namespace_name = azurerm_eventhub_namespace.eventhub_ns.name
  resource_group_name = azurerm_resource_group.rg.name
  partition_count     = 1
  message_retention   = 1
}

resource "azurerm_private_dns_zone" "dns_zone" {
  name                = "privatelink.servicebus.windows.net" # should always be "privatelink.servicebus.windows.net"
  resource_group_name = azurerm_resource_group.rg.name
}

resource "azurerm_private_dns_zone_virtual_network_link" "dns_zone_link" {
  name                  = "${var.vnet_name}-dns-link"
  resource_group_name   = azurerm_resource_group.rg.name
  private_dns_zone_name = azurerm_private_dns_zone.dns_zone.name
  virtual_network_id    = azurerm_virtual_network.vnet.id
}

resource "azurerm_private_endpoint" "eventhub_endpoint" {
 name = "pe-${var.eventhub_namespace_name}"
    location = azurerm_resource_group.rg.location
    resource_group_name = azurerm_resource_group.rg.name
    subnet_id = azurerm_subnet.subnet.id

    private_service_connection {
        name = "psc-${var.eventhub_namespace_name}"
        is_manual_connection = false
        private_connection_resource_id = azurerm_eventhub_namespace.eventhub_ns.id
        subresource_names = ["namespace"]
    }

    private_dns_zone_group {
        name = "default-dns-group"
        private_dns_zone_ids = [azurerm_private_dns_zone.dns_zone.id]
    }

    depends_on = [azurerm_private_dns_zone_virtual_network_link.dns_zone_link]
}

data "azurerm_client_config" "current" {}

resource "mongodbatlas_stream_privatelink_endpoint" "test-stream-privatelink" {
  project_id          = var.project_id
  # dns_domain comes from the hostname of the Event Hub Namespace in Azure.
  dns_domain          = "${var.eventhub_namespace_name}.servicebus.windows.net"
  provider_name       = "AZURE"
  region              = var.atlas_region
  vendor              = "GENERIC"
  # The service endpoint ID is generated as follows: /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/{namespaceName}
  service_endpoint_id = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${var.azure_resource_group}/providers/Microsoft.EventHub/namespaces/${var.eventhub_namespace_name}"
  depends_on = [azurerm_private_endpoint.eventhub_endpoint]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) The ID of the Private Link connection.
- `project_id` (String) Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.

**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group or project id remains the same. The resource and corresponding endpoints use the term groups.

### Read-Only

- `arn` (String) Amazon Resource Name (ARN).
- `dns_domain` (String) Domain name of Privatelink connected cluster.
- `dns_sub_domain` (List of String) Sub-Domain name of Confluent cluster. These are typically your availability zones.
- `error_message` (String) Error message if the connection is in a failed state.
- `interface_endpoint_id` (String) Interface endpoint ID that is created from the specified service endpoint ID.
- `interface_endpoint_name` (String) Name of interface endpoint that is created from the specified service endpoint ID.
- `provider_account_id` (String) Account ID from the cloud provider.
- `provider_name` (String) Provider where the Kafka cluster is deployed.
- `region` (String) When the vendor is `CONFLUENT`, this is the domain name of Confluent cluster. When the vendor is `MSK`, this is computed by the API from the provided `arn`.
- `service_endpoint_id` (String) Service Endpoint ID.
- `state` (String) Status of the connection.
- `vendor` (String) Vendor who manages the Kafka cluster.

For more information see: [MongoDB Atlas API - Streams Privatelink](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Streams/operation/createPrivateLinkConnection) Documentation.
