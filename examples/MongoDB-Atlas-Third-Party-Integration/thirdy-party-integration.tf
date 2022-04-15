resource "mongodbatlas_third_party_integration" "test_msteams" {
  project_id                  = mongodbatlas_project.project.id
  type                        = "MICROSOFT_TEAMS"
  microsoft_teams_webhook_url = "https://mongodb0.webhook.office.com/webhookb2/97142c7c-15e8-47de-9a7a-355183e89a68@c96563a8-841b-4ef9-af16-33548de0c958/IncomingWebhook/c0337c212d1f41149d754e9f4e0229b6/fa477bb6-9b05-473a-be3e-1157d9633ee8"
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
