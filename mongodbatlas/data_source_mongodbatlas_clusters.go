package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasClustersRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auto_scaling_disk_gb_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"backup_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"bi_connector": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeString, //Convert to Bool
										Computed: true,
									},
									"read_preference": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"cluster_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_size_gb": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"encryption_at_rest_provider": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mongo_db_major_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"num_shards": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"provider_backup_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"provider_instance_size_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"provider_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backing_provider_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"provider_disk_iops": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"provider_disk_type_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"provider_encrypt_ebs_volume": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"provider_region_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"provider_volume_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"replication_factor": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"replication_specs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"num_shards": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"regions_config": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"region_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"electable_nodes": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"priority": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"read_only_nodes": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"analytics_nodes": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
									"zone_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"mongo_db_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mongo_uri": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mongo_uri_updated": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mongo_uri_with_options": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"paused": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"srv_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceMongoDBAtlasClustersRead(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)

	clusters, resp, err := conn.Clusters.List(context.Background(), projectID, nil)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("error reading cluster list for project(%s): %s", projectID, err)
	}

	if err := d.Set("results", flattenClusters(clusters)); err != nil {
		return fmt.Errorf("error setting cluster list %s", err)
	}

	d.SetId(resource.UniqueId())

	return nil
}

func flattenClusters(clusters []matlas.Cluster) []map[string]interface{} {
	results := make([]map[string]interface{}, 0)

	for _, cluster := range clusters {
		result := map[string]interface{}{
			"auto_scaling_disk_gb_enabled": cluster.BackupEnabled,
			"backup_enabled":               cluster.BackupEnabled,
			"provider_backup_enabled":      cluster.ProviderBackupEnabled,
			"cluster_type":                 cluster.ClusterType,
			"disk_size_gb":                 cluster.DiskSizeGB,
			"encryption_at_rest_provider":  cluster.EncryptionAtRestProvider,
			"mongo_db_major_version":       cluster.MongoDBMajorVersion,
			"name":                         cluster.Name,
			"num_shards":                   cluster.NumShards,
			"mongo_db_version":             cluster.MongoDBVersion,
			"mongo_uri":                    cluster.MongoURI,
			"mongo_uri_updated":            cluster.MongoURIUpdated,
			"mongo_uri_with_options":       cluster.MongoURIWithOptions,
			"paused":                       cluster.Paused,
			"srv_address":                  cluster.SrvAddress,
			"state_name":                   cluster.StateName,
			"replication_factor":           cluster.ReplicationFactor,
			"backing_provider_name":        cluster.ProviderSettings.BackingProviderName,
			"provider_disk_iops":           cluster.ProviderSettings.DiskIOPS,
			"provider_disk_type_name":      cluster.ProviderSettings.DiskTypeName,
			"provider_encrypt_ebs_volume":  cluster.ProviderSettings.EncryptEBSVolume,
			"provider_instance_size_name":  cluster.ProviderSettings.InstanceSizeName,
			"provider_name":                cluster.ProviderSettings.ProviderName,
			"provider_region_name":         cluster.ProviderSettings.RegionName,
			"bi_connector":                 flattenBiConnector(cluster.BiConnector),
			"replication_specs":            flattenReplicationSpecs(cluster.ReplicationSpecs),
		}
		results = append(results, result)
	}
	return results
}
