data "aws_partition" "current" {}

data "aws_region" "current" {}


locals {
  UsingDefaultBucket = var.qss3_bucket_name == "aws-quickstart"
  ActivateResources = var.activate_mongo_db_resources == "Yes"
}
module atlas-basic{
    source = "/Users/sowbaranikat/Documents/CFN/terraform-provider-mongodbatlas/modules/atlas-basic"

    private_key = var.private_key
    public_key = var.public_key
    password = var.password
    database_name = var.database_name
    atlas_org_id = var.atlas_org_id
    region = var.region
}

resource "aws_ecs_service" "client_service" {
  name = "client_service"
  cluster = aws_ecs_cluster.cluster.id
  desired_count = 1
  launch_type = "FARGATE"
  platform_version = "1.4.0"
  propagate_tags = "SERVICE"
  scheduling_strategy = "REPLICA"
  network_configuration {
    subnets = [aws_subnet.subnet_east_a.id, aws_subnet.subnet_east_b.id]  # Replace with your subnet IDs
    security_groups = [aws_security_group.default_network.id]           # Replace with your security group ID
  }
  task_definition = aws_ecs_task_definition.client_task_definition.arn
  tags = {
    environment_name = var.environmentId
    Project = "MongoDbTerraformProvider"
    created_by = "aws-farget"
    creation_date = timestamp()
  }
  lifecycle {
    ignore_changes = [tags]
  }
}

resource "aws_service_discovery_service" "client_service_discovery_entry" {
  description = "Client service discovery entry in Cloud Map"
  name = "client_service_discovery_entry"
  namespace_id    = aws_service_discovery_private_dns_namespace.cloud_map.id
  dns_config {
    namespace_id    = aws_service_discovery_private_dns_namespace.cloud_map.id
    routing_policy  = "MULTIVALUE"  # Specify your routing policy
    dns_records {
      ttl  = 10
      type = "A"
    }
  }
}

resource "aws_load_balancer_listener_policy" "client_tcp8080_listener" {
  load_balancer_name = aws_load_balancer_listener_policy.load_balancer.id
  load_balancer_port = 8080
  policy_names = [
    aws_load_balancer_policy.wu-tang-ssl-tls-1-1.policy_name,
  ]
}

resource "aws_ecs_task_definition" "client_task_definition" {
  container_definitions = jsonencode([
    {
      Image = "docker/ecs-searchdomain-sidecar:1.0"
      Name = "Client_ResolvConf_InitContainer"
    },
    {
      Image = var.client_service_ecr_image_uri
      Name = "client"
    }
  ])
  cpu = "256"
  execution_role_arn = aws_iam_role.client_task_execution_role.arn
  family = "partner-meanstack-atlas-fargate-client"
  memory = "512"
  network_mode = "awsvpc"
  requires_compatibilities = [
    "FARGATE"
  ]

}

resource "aws_iam_role" "client_task_execution_role" {
  assume_role_policy = jsonencode({
    Statement = [
      {
        Action = "sts:AssumeRole"
        Sid = ""
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      },
    ]
    Version = "2012-10-17"
  })
  managed_policy_arns = [
    "arn:${data.aws_partition.current.partition}:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy",
    "arn:${data.aws_partition.current.partition}:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
  ]
  tags = {
    environment_name = var.environmentId
    Project = "MongoDbTerraformProvider"
    created_by = "aws-farget"
    creation_date = timestamp()
  }
  lifecycle {
    ignore_changes = [tags]
  }
}

resource "aws_service_discovery_private_dns_namespace" "cloud_map" {
  description = "Service Map for Docker Compose project partner-meanstack-atlas-fargate"
  name = "partner-meanstack-atlas-fargate.local"
  vpc = aws_vpc.vpc_east.id
}

resource "aws_ecs_cluster" "cluster" {
  name = "partner-meanstack-atlas-fargate"
  tags = {
    environment_name = var.environmentId
    Project = "MongoDbTerraformProvider"
    created_by = "aws-farget"
    creation_date = timestamp()
  }
  lifecycle {
    ignore_changes = [tags]
  }
}

resource "aws_vpc_security_group_ingress_rule" "default5200_ingress" {
  security_group_id = aws_security_group.default_network.id
  cidr_ipv4   = "10.0.0.0/16"
  description = "server:5200/tcp on default network"
  from_port = 5200
  ip_protocol = "TCP"
  to_port = 5200
}

resource "aws_vpc_security_group_ingress_rule" "default8080_ingress" {
  security_group_id = aws_security_group.default_network.id
  description = "client:8080/tcp on default network"
  from_port = 8080
  referenced_security_group_id = aws_security_group.default_network.id
  ip_protocol = "TCP"
  to_port = 8080
}

resource "aws_security_group" "default_network" {
  description = "partner-meanstack-atlas-fargate Security Group for default network"
  vpc_id = aws_vpc.vpc_east.id
  tags = {
    environment_name = var.environmentId
    Project = "MongoDbTerraformProvider"
    created_by = "aws-farget"
    creation_date = timestamp()
  }
  lifecycle {
    ignore_changes = [tags]
  }
}

resource "aws_vpc_security_group_ingress_rule" "default_network_ingress" {
  description = "Allow communication within network default."
  referenced_security_group_id = aws_security_group.default_network.id
  ip_protocol = "-1"
  security_group_id = aws_security_group.default_network.id
}

resource "aws_load_balancer_listener_policy" "load_balancer" {
  load_balancer_name = "MeanStackApp"
  load_balancer_port = 443
  policy_names = [
     aws_load_balancer_policy.wu-tang-ssl-tls-1-1.policy_name
  ]
}

resource "aws_ecs_service" "server_service" {
  name = "server_service"
  cluster = aws_ecs_cluster.cluster.id
  // CF Property(DeploymentConfiguration) = {
  //   MaximumPercent = 200
  //   MinimumHealthyPercent = 100
  // }
  desired_count = 1
  launch_type = "FARGATE"
  platform_version = "1.4.0"
  propagate_tags = "SERVICE"
  scheduling_strategy = "REPLICA"

  network_configuration {
    subnets = [aws_subnet.subnet_east_a.id, aws_subnet.subnet_east_b.id]  # Replace with your subnet IDs
    security_groups = [aws_security_group.default_network.id]         # Replace with your security group ID
  }
  task_definition = aws_ecs_task_definition.server_task_definition.arn
  tags = {
    environment_name = var.environmentId
    Project = "MongoDbTerraformProvider"
    created_by = "aws-farget"
    creation_date = timestamp()
  }
  lifecycle {
    ignore_changes = [tags]
  }
}

resource "aws_service_discovery_service" "server_service_discovery_entry" {
  description = "Server service discovery entry in Cloud Map"
  namespace_id    = aws_service_discovery_private_dns_namespace.cloud_map.id
  name = "server"
  dns_config {
    namespace_id    = aws_service_discovery_private_dns_namespace.cloud_map.id
    routing_policy  = "MULTIVALUE"  # Specify your routing policy
    dns_records {
      ttl  = 10
      type = "A"
    }
  }
}

resource "aws_load_balancer_listener_policy" "server_tcp5200_listener" {
  // CF Property(DefaultActions) = [
  //   {
  //     ForwardConfig = {
  //       TargetGroups = [
  //         {
  //           TargetGroupArn = aws_lb_target_group_attachment.server_tcp5200_target_group.id
  //         }
  //       ]
  //     }
  //     Type = "forward"
  //   }
  // ]
  load_balancer_name = aws_load_balancer_listener_policy.load_balancer.id
  load_balancer_port = 5200
  // CF Property(Protocol) = "TCP"
}

resource "aws_lb_target_group_attachment" "server_tcp5200_target_group" {
  target_group_arn = aws_lb_target_group.log_group.arn
  target_id = aws_instance.aws-instance.id
  port = 8090
}

data "aws_ami" "aws-ami" {
  most_recent = true
  owners      = ["amazon"]
  filter {
    name   = "architecture"
    values = ["arm64"]
  }
  filter {
    name   = "name"
    values = ["al2023-ami-2023*"]
  }
}

resource "aws_instance" "aws-instance" {
  ami = data.aws_ami.aws-ami.id
  instance_type = "t4g.micro"  
  vpc_security_group_ids = [aws_security_group.default_network.id]  # Replace with your security group ID(s)
  subnet_id = aws_subnet.subnet_east_a.id           
}

resource "aws_ecs_task_definition" "server_task_definition" {
  container_definitions = jsonencode([
    {
      Image = "docker/ecs-searchdomain-sidecar:1.0"
      Name = "Server_ResolvConf_InitContainer"
    },
    {
      Image = var.server_service_ecr_image_uri
      Name = "server"
    }
  ])
  cpu = "256"
  execution_role_arn = aws_iam_role.client_task_execution_role.arn
  family = "partner-meanstack-atlas-fargate-server"
  memory = "512"
  network_mode = "awsvpc"
  requires_compatibilities = [
    "FARGATE"
  ]
}

resource "aws_iam_role" "server_task_execution_role" {
  assume_role_policy = jsonencode({
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      },
    ]
    Version = "2012-10-17"
  })
  managed_policy_arns = [
    "arn:${data.aws_partition.current.partition}:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy",
    "arn:${data.aws_partition.current.partition}:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
  ]
  tags = {
    environment_name = var.environmentId
    Project = "MongoDbTerraformProvider"
    created_by = "aws-farget"
    creation_date = timestamp()
  }
  lifecycle {
    ignore_changes = [tags]
  }
}


output "atlas_database_user" {
  description = "Atlas database user, configured for AWS IAM role access."
  value = module.atlas-basic.dbuser
}

output "atlas_project" {
  description = "Information about your Atlas deployment."
  value = module.atlas-basic.project
}

output "atlas_project_ip_access_list" {
  description = "Atlas project IP access list."
  value = module.atlas-basic.projectIpAccessList
}

output "atlas_cluster" {
  description = "Information about your Atlas cluster."
  value = module.atlas-basic.cluster
}

output "cluster_srv_address" {
  description = "Hostname for mongodb+srv:// connection string."
  value = module.atlas-basic.cluster.connection_strings
}

output "client_url" {
  description = "Load balancer URL for client application."
  value = "http://${aws_load_balancer_listener_policy.load_balancer.load_balancer_name}:8080"
}
