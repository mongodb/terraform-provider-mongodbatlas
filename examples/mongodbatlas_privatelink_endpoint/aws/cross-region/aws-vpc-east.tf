resource "aws_vpc_endpoint" "vpce_east" {
  vpc_id             = aws_vpc.vpc_east.id
  service_name       = mongodbatlas_privatelink_endpoint.pe_east.endpoint_service_name
  vpc_endpoint_type  = "Interface"
  subnet_ids         = [aws_subnet.subnet_east.id]
  security_group_ids = [aws_security_group.sg_east.id]
}

resource "aws_vpc" "vpc_east" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true
}

resource "aws_internet_gateway" "ig_east" {
  vpc_id = aws_vpc.vpc_east.id
}

resource "aws_route" "route_east" {
  route_table_id         = aws_vpc.vpc_east.main_route_table_id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.ig_east.id
}

resource "aws_subnet" "subnet_east" {
  vpc_id                  = aws_vpc.vpc_east.id
  cidr_block              = "10.0.2.0/24"
  map_public_ip_on_launch = true
  availability_zone       = "${var.aws_region_east}a"
}

resource "aws_security_group" "sg_east" {
  name_prefix = "default-"
  description = "Default security group for all instances in ${aws_vpc.vpc_east.id}"
  vpc_id      = aws_vpc.vpc_east.id
  ingress {
    from_port = 80
    to_port   = 80
    protocol  = "tcp"
    cidr_blocks = [
      aws_vpc.vpc_east.cidr_block,
    ]
  }
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
