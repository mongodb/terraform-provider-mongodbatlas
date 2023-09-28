module "atlas-basic"  {
  source = "../../atlas-basic"

  public_key = "<publicKey>"
  private_key = "<privateKey>"
  atlas_org_id = "<orgId>"

  database_name = ["test1","test2"]
  db_users = ["user1","user2"]
  db_passwords = ["<password>","<password>"]
  database_names = ["test-db1","test-db2"]
  region = "US_EAST_1"

  aws_vpc_cidr_block = "1.0.0.0/16"
  aws_vpc_egress = "0.0.0.0/0"
  aws_vpc_ingress = "0.0.0.0/0"
  aws_subnet_cidr_block1 = "1.0.1.0/24"
  aws_subnet_cidr_block2 = "1.0.2.0/24"

  cidr_block = ["10.1.0.0/16","12.2.0.0/16"]
  ip_address = ["208.169.90.207","63.167.210.250"]

}