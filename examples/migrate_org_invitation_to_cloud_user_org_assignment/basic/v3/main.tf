############################################################
# v3: Cleaned Up Configuration
# - Remove org_invitation and moved/import blocks after v2 is applied
# - Keep only cloud_user_org_assignment and team assignments
############################################################

# Pending user is now managed directly by cloud_user_org_assignment
resource "mongodbatlas_cloud_user_org_assignment" "pending" {
  org_id   = var.org_id
  username = var.pending_username
  roles    = { org_roles = var.roles }
}

# Active user is already imported; now managed directly
resource "mongodbatlas_cloud_user_org_assignment" "active" {
  org_id   = var.org_id
  username = var.active_username
  roles    = { org_roles = ["ORG_MEMBER"] }
}

# Team assignments are managed directly; no import blocks required
resource "mongodbatlas_cloud_user_team_assignment" "teams" {
  for_each = var.pending_team_ids

  org_id  = var.org_id
  team_id = each.key
  user_id = mongodbatlas_cloud_user_org_assignment.pending.user_id
}
