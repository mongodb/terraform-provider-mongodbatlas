module "option4-autoscaling-poc" {
  source = "../.."

  name       = "single-region"
  project_id = var.project_id
  regions = [
    {
      name          = "US_EAST_1"
      node_count    = 3
      provider_name = "AWS"
    }
  ]
  auto_scaling = {
    compute_enabled            = true
    compute_max_instance_size  = "M60"
    compute_min_instance_size  = "M30"
    compute_scale_down_enabled = true
    disk_gb_enabled            = true
  }
}
