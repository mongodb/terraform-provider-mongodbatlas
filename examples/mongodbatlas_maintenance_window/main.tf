resource "mongodbatlas_project" "example" {
  name   = "project-name"
  org_id = var.org_id
}

resource "mongodbatlas_maintenance_window" "example" {
  project_id              = mongodbatlas_project.example.id
  auto_defer_once_enabled = true
  hour_of_day             = 23
  start_asap              = true
  day_of_week             = 1
  protected_hours {
    start_hour_of_day = 9
    end_hour_of_day   = 17
  }
}

data "mongodbatlas_maintenance_window" "example" {
  project_id   = mongodbatlas_maintenance_window.example.project_id
}

