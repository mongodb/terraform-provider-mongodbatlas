module "option_2_alternative" {
  source = "../.."

  name       = "single-region-sharded"
  project_id = var.project_id
  regions = [
    {
      name          = "US_EAST_1"
      node_count    = 3
      shard_index   = 0
      instance_size = "M40"
      }, {
      name          = "US_EAST_1"
      node_count    = 3
      shard_index   = 1
      instance_size = "M30"
    }
  ]
  provider_name = "AWS"
}
