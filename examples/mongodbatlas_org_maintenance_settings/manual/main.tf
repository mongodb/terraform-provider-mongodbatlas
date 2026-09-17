resource "mongodbatlas_org_maintenance_settings" "example" {
  org_id               = var.org_id
  wave_assignment_mode = "MANUAL"
}

resource "mongodbatlas_project" "dev" {
  name   = var.dev_project_name
  org_id = var.org_id
}

resource "mongodbatlas_project" "prod" {
  name   = var.prod_project_name
  org_id = var.org_id
}

resource "mongodbatlas_maintenance_window" "dev" {
  project_id      = mongodbatlas_project.dev.id
  hour_of_day     = 23
  day_of_week     = 1
  wave_assignment = 1
}

resource "mongodbatlas_maintenance_window" "prod" {
  project_id      = mongodbatlas_project.prod.id
  hour_of_day     = 23
  day_of_week     = 1
  wave_assignment = 3
}

data "mongodbatlas_org_maintenance_settings" "example" {
  org_id = mongodbatlas_org_maintenance_settings.example.org_id
}

data "mongodbatlas_maintenance_window" "dev" {
  project_id = mongodbatlas_maintenance_window.dev.project_id
}

data "mongodbatlas_maintenance_window" "prod" {
  project_id = mongodbatlas_maintenance_window.prod.project_id
}

output "org_maintenance_settings" {
  value = {
    wave_assignment_mode = data.mongodbatlas_org_maintenance_settings.example.wave_assignment_mode
  }
}

output "dev_project_maintenance_window" {
  value = {
    wave_assignment           = data.mongodbatlas_maintenance_window.dev.wave_assignment
    effective_wave_assignment = data.mongodbatlas_maintenance_window.dev.effective_wave_assignment
  }
}

output "prod_project_maintenance_window" {
  value = {
    wave_assignment           = data.mongodbatlas_maintenance_window.prod.wave_assignment
    effective_wave_assignment = data.mongodbatlas_maintenance_window.prod.effective_wave_assignment
  }
}
