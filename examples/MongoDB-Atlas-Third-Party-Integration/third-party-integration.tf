resource "mongodbatlas_third_party_integration" "test_msteams" {
  project_id                  = mongodbatlas_project.project.id
  type                        = "MICROSOFT_TEAMS"
  microsoft_teams_webhook_url = "https://mongodb0.webhook.office.com/webhookb2/zfd-15e8-47de-9a7a-355183e89a68@thi-841b-4ef9-af16-33548de0c958/IncomingWebhook/xyz"
}

resource "mongodbatlas_third_party_integration" "test_prometheus" {
  project_id        = mongodbatlas_project.project.id
  type              = "PROMETHEUS"
  user_name         = "prom_user_621952567c87684fd69b0101"
  password          = "KeQvcbkBhrNeuhVE"
  service_discovery = "file"
  scheme            = "https"
  enabled           = true
}
