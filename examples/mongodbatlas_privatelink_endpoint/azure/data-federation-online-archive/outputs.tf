output "resource_endpoint_id" {
  value = mongodbatlas_privatelink_endpoint_service_data_federation_online_archive.this.endpoint_id
}

output "resource_provider_name" {
  value = mongodbatlas_privatelink_endpoint_service_data_federation_online_archive.this.provider_name
}

output "single_ds_endpoint_id" {
  value = data.mongodbatlas_privatelink_endpoint_service_data_federation_online_archive.single.endpoint_id
}

output "single_ds_provider_name" {
  value = data.mongodbatlas_privatelink_endpoint_service_data_federation_online_archive.single.provider_name
}

output "list_ds_results_count" {
  value = length(data.mongodbatlas_privatelink_endpoint_service_data_federation_online_archives.list.results)
}
