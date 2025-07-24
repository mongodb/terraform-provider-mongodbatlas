module "option_2_alternative" {
  source = "../.."

  name       = "single-region-with-analytics"
  project_id = var.project_id
  regions = [
    {
      name                 = "US_EAST_1"
      node_count           = 3
      provider_name        = "AWS"
      node_count_analytics = 1
    }
  ]
  auto_scaling = {
    compute_enabled            = true
    compute_max_instance_size  = "M60"
    compute_min_instance_size  = "M30"
    compute_scale_down_enabled = true
    disk_gb_enabled            = true
  }
  auto_scaling_analytics = {
    compute_enabled            = true
    compute_max_instance_size  = "M30"
    compute_min_instance_size  = "M10"
    compute_scale_down_enabled = true
    disk_gb_enabled            = true
  }
}
