package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	matlas "go.mongodb.org/atlas/mongodbatlas"
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
						"auto_scaling_compute_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"auto_scaling_compute_scale_down_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"backup_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"bi_connector": {
							Type:       schema.TypeMap,
							Computed:   true,
							Deprecated: "use bi_connector_config instead",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeString, // Convert to Bool
										Computed: true,
									},
									"read_preference": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"bi_connector_config": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
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
						"connection_strings": {
							Type:     schema.TypeList,
							MinItems: 1,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"standard": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"standard_srv": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"aws_private_link": {
										Type:     schema.TypeMap,
										Computed: true,
									},
									"aws_private_link_srv": {
										Type:     schema.TypeMap,
										Computed: true,
									},
									"private": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"private_srv": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"private_endpoint": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"connection_string": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"endpoints": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"endpoint_id": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"provider_name": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"region": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
												"srv_connection_string": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"type": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
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
						"provider_auto_scaling_compute_min_instance_size": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"provider_auto_scaling_compute_max_instance_size": {
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
						"pit_enabled": {
							Type:     schema.TypeBool,
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
						"labels": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"snapshot_backup_policy": computedCloudProviderSnapshotBackupPolicySchema(),
						"container_id": {
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
	// Get client connection.
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)
	d.SetId(resource.UniqueId())

	clusters, resp, err := conn.Clusters.List(context.Background(), projectID, nil)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil
		}

		return fmt.Errorf("error reading cluster list for project(%s): %s", projectID, err)
	}

	if err := d.Set("results", flattenClusters(d, conn, clusters)); err != nil {
		return fmt.Errorf(errorClusterSetting, "results", d.Id(), err)
	}

	return nil
}

func flattenClusters(d *schema.ResourceData, conn *matlas.Client, clusters []matlas.Cluster) []map[string]interface{} {
	results := make([]map[string]interface{}, 0)

	for i := range clusters {
		snapshotBackupPolicy, err := flattenCloudProviderSnapshotBackupPolicy(d, conn, clusters[i].GroupID, clusters[i].Name)
		if err != nil {
			log.Printf("[WARN] Error setting `snapshot_backup_policy` for the cluster(%s): %s", clusters[i].ID, err)
		}

		result := map[string]interface{}{
			"auto_scaling_disk_gb_enabled": clusters[i].BackupEnabled,
			"backup_enabled":               clusters[i].BackupEnabled,
			"provider_backup_enabled":      clusters[i].ProviderBackupEnabled,
			"cluster_type":                 clusters[i].ClusterType,
			"connection_strings":           flattenConnectionStrings(clusters[i].ConnectionStrings),
			"disk_size_gb":                 clusters[i].DiskSizeGB,
			"encryption_at_rest_provider":  clusters[i].EncryptionAtRestProvider,
			"mongo_db_major_version":       clusters[i].MongoDBMajorVersion,
			"name":                         clusters[i].Name,
			"num_shards":                   clusters[i].NumShards,
			"mongo_db_version":             clusters[i].MongoDBVersion,
			"mongo_uri":                    clusters[i].MongoURI,
			"mongo_uri_updated":            clusters[i].MongoURIUpdated,
			"mongo_uri_with_options":       clusters[i].MongoURIWithOptions,
			"pit_enabled":                  clusters[i].PitEnabled,
			"paused":                       clusters[i].Paused,
			"srv_address":                  clusters[i].SrvAddress,
			"state_name":                   clusters[i].StateName,
			"replication_factor":           clusters[i].ReplicationFactor,
			"backing_provider_name":        clusters[i].ProviderSettings.BackingProviderName,
			"provider_disk_iops":           clusters[i].ProviderSettings.DiskIOPS,
			"provider_disk_type_name":      clusters[i].ProviderSettings.DiskTypeName,
			"provider_encrypt_ebs_volume":  clusters[i].ProviderSettings.EncryptEBSVolume,
			"provider_instance_size_name":  clusters[i].ProviderSettings.InstanceSizeName,
			"provider_name":                clusters[i].ProviderSettings.ProviderName,
			"provider_region_name":         clusters[i].ProviderSettings.RegionName,
			"bi_connector":                 flattenBiConnector(clusters[i].BiConnector),
			"bi_connector_config":          flattenBiConnectorConfig(clusters[i].BiConnector),
			"replication_specs":            flattenReplicationSpecs(clusters[i].ReplicationSpecs),
			"labels":                       flattenLabels(clusters[i].Labels),
			"snapshot_backup_policy":       snapshotBackupPolicy,
		}
		results = append(results, result)
	}

	return results
}
