# Create Primary VPC
resource "aws_vpc" "west" {
  provider             = aws.west
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true
}

# Create IGW
resource "aws_internet_gateway" "west" {
  provider = aws.west
  vpc_id   = aws_vpc.west.id
}

# Route Table
resource "aws_route" "west-internet_access" {
  provider               = aws.west
  route_table_id         = aws_vpc.west.main_route_table_id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.west.id
}

# Subnet-B
resource "aws_subnet" "west" {
  provider                = aws.west
  vpc_id                  = aws_vpc.west.id
  cidr_block              = "10.0.1.0/24"
  map_public_ip_on_launch = true
  availability_zone       = "${var.aws_region_west}b"
}

/*Security-Group
Ingress - Port 80 -- limited to instance
          Port 22 -- Open to ssh without limitations
Egress  - Open to All*/

resource "aws_security_group" "west" {
  provider    = aws.west
  name_prefix = "default-"
  description = "Default security group for all instances in ${aws_vpc.west.id}"
  vpc_id      = aws_vpc.west.id
  ingress {
    from_port = 80
    to_port   = 80
    protocol  = "tcp"
    cidr_blocks = [
      aws_vpc.primary.cidr_block,
    ]
  }
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
