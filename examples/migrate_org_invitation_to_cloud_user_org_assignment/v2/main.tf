############################################################
# v2: Migration
# - Pending invitation → cloud_user_org_assignment via moved block
# - Demonstrate import path for ACTIVE users and team assignments
############################################################

# New resource + moved block (recommended)
resource "mongodbatlas_cloud_user_org_assignment" "pending" {
  org_id   = var.org_id
  username = var.pending_username
  roles    = { org_roles = var.roles }
}

moved {
  from = mongodbatlas_org_invitation.pending
  to   = mongodbatlas_cloud_user_org_assignment.pending
}

# Import ACTIVE users discovered via data source
data "mongodbatlas_organization" "org" {
  org_id = var.org_id
}

locals {
  active_users = {
    for u in data.mongodbatlas_organization.org.users :
    u.id => u if u.org_membership_status == "ACTIVE" && u.username == var.active_username
  }
}

resource "mongodbatlas_cloud_user_org_assignment" "active" {
  for_each = local.active_users

  org_id   = var.org_id
  username = each.value.username
  roles    = { org_roles = each.value.roles[0].org_roles }
}

import {
  for_each = local.active_users
  to       = mongodbatlas_cloud_user_org_assignment.active[each.key]
  id       = "${var.org_id}/${each.key}"
}

# Team assignments for the pending user (after moved/import)
resource "mongodbatlas_cloud_user_team_assignment" "teams" {
  for_each = var.pending_team_ids

  org_id  = var.org_id
  team_id = each.key
  user_id = mongodbatlas_cloud_user_org_assignment.pending.user_id
}

import {
  for_each = var.pending_team_ids
  to       = mongodbatlas_cloud_user_team_assignment.teams[each.key]
  id       = "${var.org_id}/${each.key}/${var.pending_username}"
}
