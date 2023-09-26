module "atlas-basic"  {
  source = "../../atlas-basic"

  public_key = "<publicKey>"
  password = ["<password>","<password>"]
  private_key = "<privateKey>"
  database_name = ["test1","test2"]
  atlas_org_id = "<orgId>"
  region = "US_EAST_1"
  aws_vpc_egress = "0.0.0.0/0"
  aws_vpc_ingress = "0.0.0.0/0"
  aws_vpc_cidr_block = "10.0.0.0/16"
  cidr_block = ["10.1.0.0/16","12.2.0.0/16"]
  ip_address =["208.169.90.207","0.0.0.0"]
}