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
    arn = aws_msk_configuration.example.arn
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
  name           = "${var.msk_cluster_name}-msk-configuration"

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
  project_id          = var.project_id
  provider_name       = "AWS"
  vendor              = "MSK"
  arn                 = aws_msk_cluster.example.arn
}

data "mongodbatlas_stream_privatelink_endpoint" "singular_datasource" {
  project_id          = var.project_id
  id                  = mongodbatlas_stream_privatelink_endpoint.test.id
}

output "privatelink_endpoint_id" {
  value = data.mongodbatlas_stream_privatelink_endpoint.singular_datasource.id
}
