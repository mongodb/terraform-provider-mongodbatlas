resource "mongodbatlas_org_maintenance_settings" "example" {
  org_id               = var.org_id
  wave_assignment_mode = "ENV_TAG_MAPPING"
}

resource "mongodbatlas_project" "example" {
  name   = "project-name"
  org_id = var.org_id

  tags = {
    Environment = "dev"
  }
}

data "mongodbatlas_org_maintenance_settings" "example" {
  org_id = mongodbatlas_org_maintenance_settings.example.org_id
}

data "mongodbatlas_maintenance_window" "example" {
  project_id = mongodbatlas_project.example.id
}

output "org_maintenance_settings" {
  value = {
    wave_assignment_mode           = data.mongodbatlas_org_maintenance_settings.example.wave_assignment_mode
    effective_wave_assignment_mode = data.mongodbatlas_org_maintenance_settings.example.effective_wave_assignment_mode
  }
}

output "project_maintenance_window" {
  value = {
    effective_wave_assignment = data.mongodbatlas_maintenance_window.example.effective_wave_assignment
  }
}
