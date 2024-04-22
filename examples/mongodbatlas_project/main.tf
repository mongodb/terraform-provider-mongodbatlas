resource "mongodbatlas_project" "test" {
  name   = "project-name"
  org_id = var.org_id

  limits {
    name  = "atlas.project.deployment.clusters"
    value = 2
  }

  limits {
    name  = "atlas.project.deployment.nodesPerPrivateLinkRegion"
    value = 3
  }
  tags = {
    Owner       = "Terraform"
    Environment = "Example"
    Team        = "tf-experts"
    CurrentDRI  = "unset"
  }
  lifecycle {
    ignore_changes = [
      tags["CostCenter"]
    ]
  }
}
