provider "mongodbatlas" {
  client_id     = var.atlas_client_id
  client_secret = var.atlas_client_secret
}

provider "datadog" {
  api_key = var.datadog_api_key
  app_key = var.datadog_app_key
}
