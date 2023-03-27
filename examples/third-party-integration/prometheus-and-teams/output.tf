output "prometheus_config" {
  value = data.template_file.init.rendered
}