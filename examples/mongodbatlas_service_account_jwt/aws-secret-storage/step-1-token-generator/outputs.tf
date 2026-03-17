output "aws_secret_id" {
  description = "ARN of the AWS Secrets Manager secret containing the Atlas JWT."
  value       = aws_secretsmanager_secret.atlas_token.arn
}
