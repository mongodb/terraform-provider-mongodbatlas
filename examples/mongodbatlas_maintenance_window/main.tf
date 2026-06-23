resource "mongodbatlas_project" "example" {
  name   = "project-name"
  org_id = var.org_id
}

resource "mongodbatlas_maintenance_window" "example" {
  project_id              = mongodbatlas_project.example.id
  auto_defer_once_enabled = true
  hour_of_day             = 23
  day_of_week             = 1
  wave_assignment         = 1
  protected_hours {
    start_hour_of_day = 9
    end_hour_of_day   = 17
  }
}

data "mongodbatlas_maintenance_window" "example" {
  project_id = mongodbatlas_maintenance_window.example.project_id
}

output "maintenance_window" {
  value = {
    time_zone_id              = data.mongodbatlas_maintenance_window.example.time_zone_id
    wave_assignment           = data.mongodbatlas_maintenance_window.example.wave_assignment
    effective_wave_assignment = data.mongodbatlas_maintenance_window.example.effective_wave_assignment
  }
}