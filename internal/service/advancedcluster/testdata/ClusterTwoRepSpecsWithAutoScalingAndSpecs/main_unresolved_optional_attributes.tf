resource "terraform_data" "unresolved_map" {
  input = { owner = "team-atlas" }
}

resource "terraform_data" "unresolved_bool" {
  input = true
}

resource "terraform_data" "unresolved_string" {
  input = "AWS"
}

resource "terraform_data" "unresolved_date" {
  input = "2099-01-01T00:00:00Z"
}

resource "terraform_data" "unresolved_pinned_fcv" {
  input = { version = "8.0", expiration_date = "2099-01-01T00:00:00Z" }
}

resource "terraform_data" "unresolved_timeouts" {
  input = { create = "60m", update = "90m", delete = "120m" }
}

resource "mongodbatlas_advanced_cluster" "test" {
  project_id   = "111111111111111111111111"
  name         = "mocked-cluster"
  cluster_type = "GEOSHARDED"

  tags                   = terraform_data.unresolved_map.output
  labels                 = terraform_data.unresolved_map.output
  retain_backups_enabled = terraform_data.unresolved_bool.output
  adaptive_capacity      = terraform_data.unresolved_string.output
  use_effective_fields   = terraform_data.unresolved_bool.output
  pinned_fcv             = terraform_data.unresolved_pinned_fcv.output
  timeouts               = terraform_data.unresolved_timeouts.output

  accept_data_risks_and_force_replica_set_reconfig = terraform_data.unresolved_date.output

  replication_specs = [{
    region_configs = [{
      electable_specs = {
        instance_size = "M10"
        node_count    = 5
      }
      priority      = 7
      provider_name = "AWS"
      backing_provider_name = terraform_data.unresolved_string.output
      region_name   = "US_EAST_1"
    }]
    zone_name = "Zone 1"
    }, {
    region_configs = [{
      electable_specs = {
        instance_size = "M20"
        node_count    = 3
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "US_WEST_2"
    }]
    zone_name = "Zone 2"
  }]
}
