resource "aws_vpc_endpoint" "vpce_west" {
  provider           = aws.west
  vpc_id             = aws_vpc.vpc_west.id
  service_name       = mongodbatlas_privatelink_endpoint.pe_west.endpoint_service_name
  vpc_endpoint_type  = "Interface"
  subnet_ids         = [aws_subnet.subnet_west.id]
  security_group_ids = [aws_security_group.sg_west.id]
}

resource "aws_vpc" "vpc_west" {
  provider             = aws.west
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true
}

resource "aws_internet_gateway" "ig_west" {
  provider = aws.west
  vpc_id   = aws_vpc.vpc_west.id
}

resource "aws_route" "route_west" {
  provider               = aws.west
  route_table_id         = aws_vpc.vpc_west.main_route_table_id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.ig_west.id
}

resource "aws_subnet" "subnet_west" {
  provider                = aws.west
  vpc_id                  = aws_vpc.vpc_west.id
  cidr_block              = "10.0.1.0/24"
  map_public_ip_on_launch = true
  availability_zone       = "${var.aws_region_west}b"
}

resource "aws_security_group" "sg_west" {
  provider    = aws.west
  name_prefix = "default-"
  description = "Default security group for all instances in ${aws_vpc.vpc_west.id}"
  vpc_id      = aws_vpc.vpc_west.id
  ingress {
    from_port = 80
    to_port   = 80
    protocol  = "tcp"
    cidr_blocks = [
      aws_vpc.vpc_west.cidr_block,
    ]
  }
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
