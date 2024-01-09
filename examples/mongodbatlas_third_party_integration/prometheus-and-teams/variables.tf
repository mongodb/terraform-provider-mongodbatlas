variable "public_key" {
  type        = string
  description = "Public Programmatic API key to authenticate to Atlas"
}
variable "private_key" {
  type        = string
  description = "Private Programmatic API key to authenticate to Atlas"
}
variable "org_id" {
  type        = string
  description = "MongoDB Organization ID"
}
variable "project_name" {
  type        = string
  description = "The MongoDB Atlas Project Name"
}
variable "user_name" {
  type        = string
  description = "The Prometheus User Name"
  default     = "puser"
}
variable "password" {
  type        = string
  description = "The Prometheus Password"
  default     = "ppassword"
}
variable "microsoft_teams_webhook_url" {
  type        = string
  description = "The Microsoft Teams Webhook URL"
  default     = "https://yourcompany.webhook.office.com/webhookb2/zzz@yyy/IncomingWebhook/xyz"
}
