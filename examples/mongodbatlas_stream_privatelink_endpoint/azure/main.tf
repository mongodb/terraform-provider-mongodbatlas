resource "mongodbatlas_stream_privatelink_endpoint" "test-stream-privatelink" {
  project_id = var.project_id
  # dns_domain comes from the hostname of the Event Hub Namespace in Azure.
  dns_domain    = "${var.eventhub_namespace_name}.servicebus.windows.net"
  provider_name = "AZURE"
  region        = var.atlas_region
  vendor        = "EVENTHUB"
  # The service endpoint ID is generated as follows: /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/{namespaceName}
  service_endpoint_id = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${var.azure_resource_group}/providers/Microsoft.EventHub/namespaces/${var.eventhub_namespace_name}"
  depends_on          = [azurerm_private_endpoint.eventhub_endpoint]
}