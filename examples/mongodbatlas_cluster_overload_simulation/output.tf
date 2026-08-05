output "simulation_id" {
  description = "Unique identifier of the overload protection simulation."
  value       = mongodbatlas_cluster_overload_simulation.this.simulation_id
}

output "simulation_state" {
  description = "Current state of the overload protection simulation."
  value       = data.mongodbatlas_cluster_overload_simulation.this.state
}

output "simulation_expires_at" {
  description = "Expiration timestamp of the overload protection simulation."
  value       = data.mongodbatlas_cluster_overload_simulation.this.expires_at
}
