# MongoDB Atlas Provider -- AI Model Rate Limit

This example shows how to configure AI Model Rate Limits in MongoDB Atlas. AI Model Rate Limits control the number of requests and tokens per minute allowed for AI embedding and reranking model groups.

## Important Notes

Rate limits are configured per model group (e.g., `embed_large`, `embed_small`). Each model group has its own separate limits for requests per minute and tokens per minute.

Project-level rate limits cannot exceed the organization-level limits, which are determined by your organization's tier. You can retrieve organization-level limits using the `mongodbatlas_ai_model_org_rate_limits` data source.

## Variables Required to be set:
- `atlas_client_id`: MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret
- `project_id`: Project ID where the AI Model Rate Limit will be configured

## Outputs
- `ai_model_rate_limit_requests_per_minute`: The requests per minute limit
- `ai_model_rate_limit_tokens_per_minute`: The tokens per minute limit
- `ai_model_rate_limit_model_names`: The model names included in the model group
- `ai_model_rate_limits_results`: All AI Model Rate Limits in the project
