resource "mongodbatlas_third_party_integration" "test_msteams" {
  project_id                  = mongodbatlas_project.project.id
  type                        = "MICROSOFT_TEAMS"
  microsoft_teams_webhook_url = var.microsoft_teams_webhook_url
}

resource "mongodbatlas_third_party_integration" "test_prometheus" {
  project_id        = mongodbatlas_project.project.id
  type              = "PROMETHEUS"
  user_name         = var.user_name
  password          = var.password
  service_discovery = "file"
  scheme            = "https"
  enabled           = true
}

data "template_file" "init" {
  template = file("prom-conf-yml.tpl")

  vars = {
    job_name = "test-job"
    group_id = mongodbatlas_project.project.id
    password = mongodbatlas_third_party_integration.test_prometheus.password
  }
}