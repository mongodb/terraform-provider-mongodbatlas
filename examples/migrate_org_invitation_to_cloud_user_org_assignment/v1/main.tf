############################################################
# v1: Initial State
# - One pending invitation managed via mongodbatlas_org_invitation (with teams)
# - One active user present in org (no invitation resource in state)
############################################################

# Pending invitation (with teams)
resource "mongodbatlas_org_invitation" "pending" {
  org_id    = var.org_id
  username  = var.pending_username
  roles     = var.roles
  teams_ids = var.pending_team_ids
}

# Active user is represented only for reference via data source
data "mongodbatlas_organization" "org" {
  org_id = var.org_id
}

locals {
  active_users = {
    for u in data.mongodbatlas_organization.org.users :
    u.username => u if u.org_membership_status == "ACTIVE" && u.username == var.active_username
  }
}

output "active_users" {
  value = local.active_users
}
