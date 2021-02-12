package mongodbatlas

import (
	"encoding/base64"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/spf13/cast"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

// Provider returns the provider to be use by the code.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"public_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MONGODB_ATLAS_PUBLIC_KEY", ""),
				Description: "MongoDB Atlas Programmatic Public Key",
			},
			"private_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MONGODB_ATLAS_PRIVATE_KEY", ""),
				Description: "MongoDB Atlas Programmatic Private Key",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"mongodbatlas_custom_db_role":                        dataSourceMongoDBAtlasCustomDBRole(),
			"mongodbatlas_custom_db_roles":                       dataSourceMongoDBAtlasCustomDBRoles(),
			"mongodbatlas_database_user":                         dataSourceMongoDBAtlasDatabaseUser(),
			"mongodbatlas_database_users":                        dataSourceMongoDBAtlasDatabaseUsers(),
			"mongodbatlas_project":                               dataSourceMongoDBAtlasProject(),
			"mongodbatlas_projects":                              dataSourceMongoDBAtlasProjects(),
			"mongodbatlas_project_ip_whitelist":                  dataSourceMongoDBAtlasProjectIPWhitelist(),
			"mongodbatlas_cluster":                               dataSourceMongoDBAtlasCluster(),
			"mongodbatlas_clusters":                              dataSourceMongoDBAtlasClusters(),
			"mongodbatlas_cloud_provider_snapshot":               dataSourceMongoDBAtlasCloudProviderSnapshot(),
			"mongodbatlas_cloud_provider_snapshots":              dataSourceMongoDBAtlasCloudProviderSnapshots(),
			"mongodbatlas_network_container":                     dataSourceMongoDBAtlasNetworkContainer(),
			"mongodbatlas_network_containers":                    dataSourceMongoDBAtlasNetworkContainers(),
			"mongodbatlas_network_peering":                       dataSourceMongoDBAtlasNetworkPeering(),
			"mongodbatlas_network_peerings":                      dataSourceMongoDBAtlasNetworkPeerings(),
			"mongodbatlas_cloud_provider_snapshot_restore_job":   dataSourceMongoDBAtlasCloudProviderSnapshotRestoreJob(),
			"mongodbatlas_cloud_provider_snapshot_restore_jobs":  dataSourceMongoDBAtlasCloudProviderSnapshotRestoreJobs(),
			"mongodbatlas_maintenance_window":                    dataSourceMongoDBAtlasMaintenanceWindow(),
			"mongodbatlas_auditing":                              dataSourceMongoDBAtlasAuditing(),
			"mongodbatlas_team":                                  dataSourceMongoDBAtlasTeam(),
			"mongodbatlas_teams":                                 dataSourceMongoDBAtlasTeam(),
			"mongodbatlas_global_cluster_config":                 dataSourceMongoDBAtlasGlobalCluster(),
			"mongodbatlas_alert_configuration":                   dataSourceMongoDBAtlasAlertConfiguration(),
			"mongodbatlas_x509_authentication_database_user":     dataSourceMongoDBAtlasX509AuthDBUser(),
			"mongodbatlas_privatelink_endpoint":                  dataSourceMongoDBAtlasPrivateLinkEndpoint(),
			"mongodbatlas_privatelink_endpoint_service":          dataSourceMongoDBAtlasPrivateEndpointServiceLink(),
			"mongodbatlas_private_endpoint":                      dataSourceMongoDBAtlasPrivateEndpoint(),
			"mongodbatlas_private_endpoint_interface_link":       dataSourceMongoDBAtlasPrivateEndpointInterfaceLink(),
			"mongodbatlas_cloud_provider_snapshot_backup_policy": dataSourceMongoDBAtlasCloudProviderSnapshotBackupPolicy(),
			"mongodbatlas_third_party_integrations":              dataSourceMongoDBAtlasThirdPartyIntegrations(),
			"mongodbatlas_third_party_integration":               dataSourceMongoDBAtlasThirdPartyIntegration(),
			"mongodbatlas_project_ip_access_list":                dataSourceMongoDBAtlasProjectIPAccessList(),
			"mongodbatlas_cloud_provider_access":                 dataSourceMongoDBAtlasCloudProviderAccessList(),
			"mongodbatlas_custom_dns_configuration_cluster_aws":  dataSourceMongoDBAtlasCustomDNSConfigurationAWS(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"mongodbatlas_custom_db_role":                        resourceMongoDBAtlasCustomDBRole(),
			"mongodbatlas_database_user":                         resourceMongoDBAtlasDatabaseUser(),
			"mongodbatlas_project_ip_whitelist":                  resourceMongoDBAtlasProjectIPWhitelist(),
			"mongodbatlas_project":                               resourceMongoDBAtlasProject(),
			"mongodbatlas_cluster":                               resourceMongoDBAtlasCluster(),
			"mongodbatlas_cloud_provider_snapshot":               resourceMongoDBAtlasCloudProviderSnapshot(),
			"mongodbatlas_network_container":                     resourceMongoDBAtlasNetworkContainer(),
			"mongodbatlas_cloud_provider_snapshot_restore_job":   resourceMongoDBAtlasCloudProviderSnapshotRestoreJob(),
			"mongodbatlas_network_peering":                       resourceMongoDBAtlasNetworkPeering(),
			"mongodbatlas_encryption_at_rest":                    resourceMongoDBAtlasEncryptionAtRest(),
			"mongodbatlas_private_ip_mode":                       resourceMongoDBAtlasPrivateIPMode(),
			"mongodbatlas_maintenance_window":                    resourceMongoDBAtlasMaintenanceWindow(),
			"mongodbatlas_auditing":                              resourceMongoDBAtlasAuditing(),
			"mongodbatlas_team":                                  resourceMongoDBAtlasTeam(),
			"mongodbatlas_teams":                                 resourceMongoDBAtlasTeam(),
			"mongodbatlas_global_cluster_config":                 resourceMongoDBAtlasGlobalCluster(),
			"mongodbatlas_alert_configuration":                   resourceMongoDBAtlasAlertConfiguration(),
			"mongodbatlas_x509_authentication_database_user":     resourceMongoDBAtlasX509AuthDBUser(),
			"mongodbatlas_privatelink_endpoint":                  resourceMongoDBAtlasPrivateLinkEndpoint(),
			"mongodbatlas_privatelink_endpoint_service":          resourceMongoDBAtlasPrivateEndpointServiceLink(),
			"mongodbatlas_private_endpoint":                      resourceMongoDBAtlasPrivateEndpoint(),
			"mongodbatlas_private_endpoint_interface_link":       resourceMongoDBAtlasPrivateEndpointInterfaceLink(),
			"mongodbatlas_cloud_provider_snapshot_backup_policy": resourceMongoDBAtlasCloudProviderSnapshotBackupPolicy(),
			"mongodbatlas_third_party_integration":               resourceMongoDBAtlasThirdPartyIntegration(),
			"mongodbatlas_project_ip_access_list":                resourceMongoDBAtlasProjectIPAccessList(),
			"mongodbatlas_cloud_provider_access":                 resourceMongoDBAtlasCloudProviderAccess(),
			"mongodbatlas_custom_dns_configuration_cluster_aws":  resourceMongoDBAtlasCustomDNSConfiguration(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		PublicKey:  d.Get("public_key").(string),
		PrivateKey: d.Get("private_key").(string),
	}

	return config.NewClient(), nil
}

func encodeStateID(values map[string]string) string {
	encode := func(e string) string { return base64.StdEncoding.EncodeToString([]byte(e)) }
	encodedValues := make([]string, 0)

	// sort to make sure the same encoding is returned in case of same input
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		encodedValues = append(encodedValues, fmt.Sprintf("%s:%s", encode(key), encode(values[key])))
	}

	return strings.Join(encodedValues, "-")
}

func decodeStateID(stateID string) map[string]string {
	decode := func(d string) string {
		decodedString, err := base64.StdEncoding.DecodeString(d)
		if err != nil {
			log.Printf("[WARN] error decoding state ID: %s", err)
		}

		return string(decodedString)
	}
	decodedValues := make(map[string]string)
	encodedValues := strings.Split(stateID, "-")

	for _, value := range encodedValues {
		keyValue := strings.Split(value, ":")
		decodedValues[decode(keyValue[0])] = decode(keyValue[1])
	}

	return decodedValues
}

func valRegion(reg interface{}, opt ...string) (string, error) {
	region, err := cast.ToStringE(reg)
	if err != nil {
		return "", err
	}

	if region == "" {
		return "", fmt.Errorf("region must be set")
	}

	/*
		We need to check if the option will be similar to network_peering word
		 (this comes in from the same resource) because network_pering resource
		 has not the standard region name pattern "US_EAST_1",
		 instead it needs the following one: "us-east-1".
	*/
	if len(opt) > 0 && strings.EqualFold("network_peering", opt[0]) {
		return strings.ToLower(strings.ReplaceAll(region, "_", "-")), nil
	}

	return strings.ReplaceAll(region, "-", "_"), nil
}

func flattenLabels(l []matlas.Label) []map[string]interface{} {
	labels := make([]map[string]interface{}, len(l))
	for i, v := range l {
		labels[i] = map[string]interface{}{
			"key":   v.Key,
			"value": v.Value,
		}
	}

	return labels
}

func expandLabelSliceFromSetSchema(d *schema.ResourceData) []matlas.Label {
	list := d.Get("labels").(*schema.Set)
	res := make([]matlas.Label, list.Len())

	for i, val := range list.List() {
		v := val.(map[string]interface{})
		res[i] = matlas.Label{
			Key:   v["key"].(string),
			Value: v["value"].(string),
		}
	}

	return res
}

func containsLabelOrKey(list []matlas.Label, item matlas.Label) bool {
	for _, v := range list {
		if reflect.DeepEqual(v, item) || v.Key == item.Key {
			return true
		}
	}

	return false
}

func removeLabel(list []matlas.Label, item matlas.Label) []matlas.Label {
	var pos int

	for _, v := range list {
		if reflect.DeepEqual(v, item) {
			list = append(list[:pos], list[pos+1:]...)

			if pos > 0 {
				pos--
			}

			continue
		}
		pos++
	}

	return list
}

func expandStringList(list []interface{}) (res []string) {
	for _, v := range list {
		res = append(res, v.(string))
	}

	return
}
