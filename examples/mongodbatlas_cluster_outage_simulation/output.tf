output "simulation_id" {
  value = mongodbatlas_cluster_outage_simulation.outage_simulation.simulation_id
}

output "simulation_state" {
  value = mongodbatlas_cluster_outage_simulation.outage_simulation.state
}

output "simulation_start_date" {
  value = mongodbatlas_cluster_outage_simulation.outage_simulation.start_request_date
}
