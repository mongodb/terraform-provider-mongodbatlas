data "mongodbatlas_alert_configurations" "import" {
  project_id = var.project_id

  output_type = ["resource_hcl", "resource_import"]
}

locals {
  alerts = data.mongodbatlas_alert_configurations.import.results

  alert_resources = compact([
    for i, alert in local.alerts :
    alert.output == null ? null :
    length(alert.output) < 1 == null ? null : alert.output[0].value
  ])

  alert_imports = compact([
    for i, alert in local.alerts :
    alert.output == null ? null :
    length(alert.output) < 2 == null ? null : alert.output[1].value
  ])
}

output "alert_resources" {
  value = join("\n", local.alert_resources)
}

output "alert_imports" {
  value = join("", local.alert_imports)
}
