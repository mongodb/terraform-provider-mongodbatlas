# MongoDB Atlas Provider -- AI Model API Key

This example shows how to create an AI Model API Key in MongoDB Atlas. AI Model API Keys are used to authenticate requests to AI embedding and reranking services.

## Important Notes

When an AI Model API Key is created, a secret is automatically generated and stored in Terraform state. The secret can be retrieved from the resource output at any time.

You can retrieve it using (**warning**: this prints the secret to your terminal):

```bash
terraform output -raw ai_model_api_key_secret
```

**Note on Import:** When importing an existing AI Model API Key, the `secret` attribute will be null because the secret is only returned by the API at creation time. If you need the secret value, you must create a new API key through Terraform.

## Variables Required to be set:
- `atlas_client_id`: MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret
- `project_id`: Project ID where the AI Model API Key will be created

## Outputs
- `ai_model_api_key_id`: The ID of the created AI Model API Key
- `ai_model_api_key_secret` (sensitive): The secret value (null if resource was imported)
- `ai_model_api_key_name`: The name of the AI Model API key from the data source
- `ai_model_api_keys_results`: All AI Model API keys in the project
