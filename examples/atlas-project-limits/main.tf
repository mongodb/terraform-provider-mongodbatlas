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
}