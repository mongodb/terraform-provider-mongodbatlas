resource "mongodbatlas_org_maintenance_settings" "example" {
  org_id               = var.org_id
  wave_assignment_mode = "ENV_TAG_MAPPING"
}

resource "mongodbatlas_project" "dev" {
  name   = var.dev_project_name
  org_id = var.org_id

  tags = {
    environment = "development"
  }
}

resource "mongodbatlas_project" "prod" {
  name   = var.prod_project_name
  org_id = var.org_id

  tags = {
    environment = "production"
  }
}

data "mongodbatlas_org_maintenance_settings" "example" {
  org_id = mongodbatlas_org_maintenance_settings.example.org_id
}

output "org_maintenance_settings" {
  value = {
    wave_assignment_mode           = data.mongodbatlas_org_maintenance_settings.example.wave_assignment_mode
    effective_wave_assignment_mode = data.mongodbatlas_org_maintenance_settings.example.effective_wave_assignment_mode
  }
}
