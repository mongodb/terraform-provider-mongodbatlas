resource "mongodbatlas_ai_model_rate_limit" "this" {
  project_id                 = var.project_id
  model_group_name           = "embed_large"
  requests_per_minute_limit  = 100
  tokens_per_minute_limit    = 10000
}

data "mongodbatlas_ai_model_rate_limit" "this" {
  project_id       = var.project_id
  model_group_name = mongodbatlas_ai_model_rate_limit.this.model_group_name
}

data "mongodbatlas_ai_model_rate_limits" "this" {
  project_id = var.project_id
}

output "ai_model_rate_limit_requests_per_minute" {
  description = "The requests per minute limit for the AI Model Rate Limit."
  value       = mongodbatlas_ai_model_rate_limit.this.requests_per_minute_limit
}

output "ai_model_rate_limit_tokens_per_minute" {
  description = "The tokens per minute limit for the AI Model Rate Limit."
  value       = mongodbatlas_ai_model_rate_limit.this.tokens_per_minute_limit
}

output "ai_model_rate_limit_model_names" {
  description = "The model names included in this model group from the data source."
  value       = data.mongodbatlas_ai_model_rate_limit.this.model_names
}

output "ai_model_rate_limits_results" {
  description = "All AI Model Rate Limits in the project."
  value       = data.mongodbatlas_ai_model_rate_limits.this.results
}
