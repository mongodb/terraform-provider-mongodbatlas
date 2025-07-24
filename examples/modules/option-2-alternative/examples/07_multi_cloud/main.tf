module "option-2-alternative" {
  source = "../.."

  name       = "multi-cloud"
  project_id = var.project_id
  regions = [
    {
      name          = "US_WEST_2"
      node_count    = 2
      shard_index   = 0
      provider_name = "AZURE"
      }, {
      name                 = "US_EAST_2"
      node_count           = 1
      shard_index          = 1
      provider_name        = "AWS"
      node_count_read_only = 2
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
