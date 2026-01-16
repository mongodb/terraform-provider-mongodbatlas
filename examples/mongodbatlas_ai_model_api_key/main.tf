resource "mongodbatlas_ai_model_api_key" "this" {
  project_id = var.project_id
  name       = "example-ai-model-key"
}

data "mongodbatlas_ai_model_api_key" "this" {
  project_id = var.project_id
  api_key_id = mongodbatlas_ai_model_api_key.this.api_key_id
}

data "mongodbatlas_ai_model_api_keys" "this" {
  project_id = var.project_id
}

output "ai_model_api_key_id" {
  description = "The ID of the AI Model API key."
  value       = mongodbatlas_ai_model_api_key.this.api_key_id
}

output "ai_model_api_key_secret" {
  description = "The secret value of the AI Model API key, null if the resource was imported."
  value       = mongodbatlas_ai_model_api_key.this.secret
  sensitive   = true
}

output "ai_model_api_key_name" {
  description = "The name of the AI Model API key from the data source."
  value       = data.mongodbatlas_ai_model_api_key.this.name
}

output "ai_model_api_keys_results" {
  description = "All AI Model API keys in the project."
  value       = data.mongodbatlas_ai_model_api_keys.this.results
}
