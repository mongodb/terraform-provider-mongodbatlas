# Old module usage
module "user_team_assignment" {  
  source     = "../../module_maintainer/v1" 
  org_id     = var.org_id  
  team_name  = var.team_name  
  usernames  = var.usernames  
}