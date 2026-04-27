package acc

import (
	"fmt"
	"strings"
)

func ConfigDataConfluentDedicatedCluster(networkID, privatelinkAccessID string) string {
	return fmt.Sprintf(`
	%[1]s

	data "confluent_environment" "default_environment" {
  		display_name = "default"
	}

	data "confluent_network" "private-link" {
		id     = %[2]q
		environment {
			id = data.confluent_environment.default_environment.id
		}
	}

	data "confluent_private_link_access" "aws" {
		id = %[3]q
		environment {
			id = data.confluent_environment.default_environment.id
		}
	}`, ConfigConfluentProvider(), networkID, privatelinkAccessID)
}

func configBasic(projectID, provider, region, vendor string, withDNSSubdomains bool) string {
	dnsSubDomainConfig := ""
	if withDNSSubdomains {
		dnsSubDomainConfig = `dns_sub_domain = local.dns_sub_domain_entries`
	}

	return fmt.Sprintf(`
	locals {
		dns_sub_domain_entries = [
    		for zone in confluent_network.private-link.zones :
    		"${zone}.${confluent_network.private-link.dns_domain}"
  		]
	}	

	resource "mongodbatlas_stream_privatelink_endpoint" "test" {
		project_id          = %[1]q
		dns_domain          = confluent_network.private-link.dns_domain
		provider_name       = %[2]q
		region              = %[3]q
		vendor              = %[4]q
		service_endpoint_id = confluent_network.private-link.aws[0].private_link_endpoint_service
		%[5]s
		depends_on = [
			confluent_kafka_cluster.dedicated
		]
	}

	data "mongodbatlas_stream_privatelink_endpoint" "test" {
		project_id = %[1]q
		id         = mongodbatlas_stream_privatelink_endpoint.test.id
		depends_on = [
    		mongodbatlas_stream_privatelink_endpoint.test
  		]
	}

	data "mongodbatlas_stream_privatelink_endpoints" "test" {
		project_id = %[1]q
		depends_on = [
    		mongodbatlas_stream_privatelink_endpoint.test
  		]
	}`, projectID, provider, region, vendor, dnsSubDomainConfig)
}

func configNewConfluentDedicatedCluster(provider, region, awsAccountID string) string {
	return fmt.Sprintf(`
	%[1]s

	data "confluent_environment" "default_environment" {
  		display_name = "default"
	}

	resource "confluent_network" "private-link" {
		display_name     = "terraform-test-private-link-network"
		cloud            = %[2]q
		region           = %[3]q
		connection_types = ["PRIVATELINK"]
		zones            = ["use1-az1", "use1-az2", "use1-az4"]
		environment {
			id = data.confluent_environment.default_environment.id
		}
		dns_config {
			resolution = "PRIVATE"
		}
	}

	resource "confluent_private_link_access" "aws" {
		display_name = "terraform-test-aws-private-link-access"
		aws {
			account = %[4]q
		}
		environment {
			id = data.confluent_environment.default_environment.id
		}
		network {
			id = confluent_network.private-link.id
		}
	}

	resource "confluent_kafka_cluster" "dedicated" {
		display_name = "terraform-test-cluster"
		availability = "MULTI_ZONE"
		cloud        = confluent_network.private-link.cloud
		region       = confluent_network.private-link.region
		dedicated {
			cku = 2
		}
		environment {
			id = data.confluent_environment.default_environment.id
		}
		network {
			id = confluent_network.private-link.id
		}
	}`, ConfigConfluentProvider(), provider, region, awsAccountID)
}

func GetCompleteConfluentConfig(usesExistingConfluentCluster, withDNSSubdomains bool, projectID, provider, region, vendor, awsAccountID, networkID, privatelinkAccessID string) string {
	if usesExistingConfluentCluster {
		configBasicUsingDatasources := strings.ReplaceAll(configBasic(projectID, provider, region, vendor, withDNSSubdomains), "confluent_network.private-link", "data.confluent_network.private-link")
		configBasicUsingDatasourcesWithoutDependsOnCluster := strings.ReplaceAll(configBasicUsingDatasources, "confluent_kafka_cluster.dedicated", "")
		return ConfigDataConfluentDedicatedCluster(networkID, privatelinkAccessID) + configBasicUsingDatasourcesWithoutDependsOnCluster
	}
	return configNewConfluentDedicatedCluster(provider, region, awsAccountID) + configBasic(projectID, provider, region, vendor, true)
}

func GetCompleteMskConfig(projectID, clusterArn string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_stream_privatelink_endpoint" "test" {
		project_id          = %[1]q
		provider_name       = "AWS"
		vendor              = "MSK"
		arn                 = %[2]q
	}

	data "mongodbatlas_stream_privatelink_endpoint" "test" {
		project_id = %[1]q
		id         = mongodbatlas_stream_privatelink_endpoint.test.id
		depends_on = [
			mongodbatlas_stream_privatelink_endpoint.test
		]
	}

	data "mongodbatlas_stream_privatelink_endpoints" "test" {
		project_id = %[1]q
		depends_on = [
			mongodbatlas_stream_privatelink_endpoint.test
		]
	}`, projectID, clusterArn)
}

func GetCompleteAzureBlobStorageConfig(projectID, clusterName, subscriptionID, clientID, clientSecret, tenantID, resourceGroupName, storageAccountName string) string {
	return fmt.Sprintf(`
	%[1]s

	resource "mongodbatlas_advanced_cluster" "test" {
		project_id   = %[2]q
		name         = %[5]q
		cluster_type = "REPLICASET"
		replication_specs = [{
			region_configs = [{
				priority      = 7
				provider_name = "AZURE"
				region_name   = "US_EAST_2"
				electable_specs = {
					instance_size = "M10"
					node_count    = 3
				}
			}]
		}]
	}

	resource "azurerm_resource_group" "blob_pl_rg" {
		name     = %[3]q
		location = "East US 2"
	}

	resource "azurerm_storage_account" "blob_pl_storage" {
		name                     = %[4]q
		resource_group_name      = azurerm_resource_group.blob_pl_rg.name
		location                 = azurerm_resource_group.blob_pl_rg.location
		account_tier             = "Standard"
		account_replication_type = "LRS"
	}

	resource "mongodbatlas_stream_privatelink_endpoint" "test" {
		project_id          = %[2]q
		provider_name       = "AZURE"
		vendor              = "AZURE_BLOB_STORAGE"
		region              = "eastus2"
		service_endpoint_id = azurerm_storage_account.blob_pl_storage.id
		dns_domain          = "${azurerm_storage_account.blob_pl_storage.name}.blob.core.windows.net"
		depends_on          = [mongodbatlas_advanced_cluster.test]
	}

	data "mongodbatlas_stream_privatelink_endpoint" "test" {
		project_id = %[2]q
		id         = mongodbatlas_stream_privatelink_endpoint.test.id
	}

	data "mongodbatlas_stream_privatelink_endpoints" "test" {
		project_id = %[2]q
		depends_on = [
			mongodbatlas_stream_privatelink_endpoint.test
		]
	}`, ConfigAzurermProvider(subscriptionID, clientID, clientSecret, tenantID),
		projectID, resourceGroupName, storageAccountName, clusterName)
}

func GetCompleteS3Config(projectID, region string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_stream_privatelink_endpoint" "test" {
		project_id          = %[1]q
		provider_name       = "AWS"
		vendor              = "S3"
		region              = %[2]q
		service_endpoint_id = "com.amazonaws.%[2]s.s3"
	}

	data "mongodbatlas_stream_privatelink_endpoint" "test" {
		project_id = %[1]q
		id         = mongodbatlas_stream_privatelink_endpoint.test.id
	}

	data "mongodbatlas_stream_privatelink_endpoints" "test" {
		project_id = %[1]q
		depends_on = [
			mongodbatlas_stream_privatelink_endpoint.test
		]
	}`, projectID, region)
}

func GetCompletePubSubConfig(projectID, clusterName, region string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_advanced_cluster" "test" {
		project_id   = %[1]q
		name         = %[3]q
		cluster_type = "REPLICASET"
		replication_specs = [{
			region_configs = [{
				priority      = 7
				provider_name = "GCP"
				region_name   = "US_EAST_4"
				electable_specs = {
					instance_size = "M10"
					node_count    = 3
				}
			}]
		}]
	}

	resource "mongodbatlas_stream_privatelink_endpoint" "test" {
		project_id    = %[1]q
		provider_name = "GCP"
		vendor        = "PUBSUB"
		region        = %[2]q
		depends_on    = [mongodbatlas_advanced_cluster.test]
	}

	data "mongodbatlas_stream_privatelink_endpoint" "test" {
		project_id = %[1]q
		id         = mongodbatlas_stream_privatelink_endpoint.test.id
	}

	data "mongodbatlas_stream_privatelink_endpoints" "test" {
		project_id = %[1]q
		depends_on = [
			mongodbatlas_stream_privatelink_endpoint.test
		]
	}`, projectID, region, clusterName)
}
