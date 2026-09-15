output "adaptive_settings_overrides" {
  description = "Customer-specified Adaptive Settings overrides."
  value       = data.mongodbatlas_cluster_adaptive_settings.this.adaptive_settings_overrides
}

output "effective_adaptive_settings" {
  description = "Effective Adaptive Settings after Atlas applies overrides and managed defaults."
  value       = data.mongodbatlas_cluster_adaptive_settings.this.effective_adaptive_settings
}
