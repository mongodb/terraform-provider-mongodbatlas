resource "aws_vpc_endpoint" "vpce_east" {
  vpc_id             = aws_vpc.vpc_east.id
  service_name       = mongodbatlas_privatelink_endpoint.pe_east.endpoint_service_name
  vpc_endpoint_type  = "Interface"
  subnet_ids         = [aws_subnet.subnet_east_a.id, aws_subnet.subnet_east_b.id]
  security_group_ids = [aws_security_group.sg_east.id]
}

resource "aws_vpc" "vpc_east" {
  cidr_block           = var.aws_vpc_cidr_block
  enable_dns_hostnames = true
  enable_dns_support   = true
}

resource "aws_internet_gateway" "ig_east" {
  vpc_id = aws_vpc.vpc_east.id
}

resource "aws_route" "route_east" {
  route_table_id         = aws_vpc.vpc_east.main_route_table_id
  destination_cidr_block = var.aws_route_table_cidr_block
  gateway_id             = aws_internet_gateway.ig_east.id
}

resource "aws_subnet" "subnet_east_a" {
  vpc_id                  = aws_vpc.vpc_east.id
  cidr_block              = var.aws_subnet_cidr_block1
  map_public_ip_on_launch = true
  availability_zone       = var.aws_subnet_availability_zone1
}

resource "aws_subnet" "subnet_east_b" {
  vpc_id                  = aws_vpc.vpc_east.id
  cidr_block              = var.aws_subnet_cidr_block2
  map_public_ip_on_launch = false
  availability_zone       = var.aws_subnet_availability_zone2
}

resource "aws_security_group" "sg_east" {
  name_prefix = "default-"
  description = "Default security group for all instances in vpc"
  vpc_id      = aws_vpc.vpc_east.id
  ingress {
    from_port = var.aws_sg_ingress_from_port
    to_port   = var.aws_sg_ingress_to_port
    protocol  = var.aws_sg_ingress_protocol
    cidr_blocks = [
      var.aws_vpc_cidr_block,
    ]
  }
  egress {
    from_port   = var.aws_sg_egress_from_port
    to_port     = var.aws_sg_egress_to_port
    protocol    = var.aws_sg_egress_protocol
    cidr_blocks = [
      var.aws_vpc_cidr_block
    ]
  }
}
