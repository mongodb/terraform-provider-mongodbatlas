module "atlas-basic"  {
  source = "../../atlas-basic"

  public_key = ""
  password = []
  private_key = ""
  database_name = ["test1","test2"]
  atlas_org_id = ""
  region = "US_EAST_1"
}