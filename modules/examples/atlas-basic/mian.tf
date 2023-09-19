module "atlas-basic"  {
  source = "../../atlas-basic"

  public_key = "<publicKey>"
  password = ["<password>","<password>"]
  private_key = "<privateKey>"
  database_name = ["test1","test2"]
  atlas_org_id = "<orgId>"
  region = "US_EAST_1"
}