package cluster

import (
	"context"
	"fmt"
	"log"
	"net/http"

	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
	"go.mongodb.org/atlas-sdk/v20250312007/admin"
	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
)

func PluralDataSource() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: fmt.Sprintf(constant.DeprecationNextMajorWithReplacementGuide, "datasource", "mongodbatlas_advanced_clusters", clusterToAdvancedClusterGuide),
		ReadContext:        dataSourcePluralRead,
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
						"advanced_configuration": advancedcluster.SchemaAdvancedConfigDS(),
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
						"bi_connector_config": {
							Type:     schema.TypeList,
							Computed: true,
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
												"srv_shard_optimized_connection_string": {
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
						"tags":                   &advancedcluster.DSTagsSchema,
						"snapshot_backup_policy": computedCloudProviderSnapshotBackupPolicySchema(),
						"container_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"termination_protection_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"version_release_system": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"redact_client_log_data": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"pinned_fcv": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"expiration_date": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourcePluralRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).Atlas
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	connV220240530 := meta.(*config.MongoDBClient).AtlasV220240530
	projectID := d.Get("project_id").(string)
	d.SetId(id.UniqueId())

	clusters, resp, err := conn.Clusters.List(ctx, projectID, nil)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil
		}
		return diag.FromErr(fmt.Errorf("error reading cluster list for project(%s): %s", projectID, err))
	}

	latestClusterModels, err := newAtlasList(ctx, connV2, projectID)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil
		}
		return diag.FromErr(fmt.Errorf("error reading new cluster list for project(%s): %s", projectID, err))
	}

	if err := d.Set("results", flattenClusters(ctx, d, conn, connV2, connV220240530, clusters, latestClusterModels)); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "results", d.Id(), err))
	}

	return nil
}

func flattenClusters(ctx context.Context, d *schema.ResourceData, conn *matlas.Client, connV2 *admin.APIClient, connV220240530 *admin20240530.APIClient, clusters []matlas.Cluster, latestClusterModels map[string]*admin.ClusterDescription20240805) []map[string]any {
	results := make([]map[string]any, 0)

	for i := range clusters {
		snapshotBackupPolicy, err := flattenCloudProviderSnapshotBackupPolicy(ctx, d, conn, clusters[i].GroupID, clusters[i].Name)
		if err != nil {
			log.Printf("[WARN] Error setting `snapshot_backup_policy` for the cluster(%s): %s", clusters[i].ID, err)
		}

		processArgs20240530, _, err := connV220240530.ClustersApi.GetClusterAdvancedConfiguration(ctx, clusters[i].GroupID, clusters[i].Name).Execute()
		if err != nil {
			log.Printf("[WARN] Error setting `advanced_configuration` for the cluster(%s): %s", clusters[i].ID, err)
		}
		processArgs, _, err := connV2.ClustersApi.GetProcessArgs(ctx, clusters[i].GroupID, clusters[i].Name).Execute()
		if err != nil {
			log.Printf("[WARN] Error setting `advanced_configuration` for the cluster(%s): %s", clusters[i].ID, err)
		}

		var containerID string
		if clusters[i].ProviderSettings != nil && clusters[i].ProviderSettings.ProviderName != "TENANT" {
			containers, _, err := conn.Containers.List(ctx, clusters[i].GroupID,
				&matlas.ContainersListOptions{ProviderName: clusters[i].ProviderSettings.ProviderName})
			if err != nil {
				log.Printf(errorClusterRead, clusters[i].Name, err)
			}

			containerID = getContainerID(containers, &clusters[i])
		}
		p := &ProcessArgs{
			argsDefault:           processArgs,
			argsLegacy:            processArgs20240530,
			clusterAdvancedConfig: clusters[i].AdvancedConfiguration,
		}
		result := map[string]any{
			"advanced_configuration":                  flattenProcessArgs(p),
			"auto_scaling_compute_enabled":            clusters[i].AutoScaling.Compute.Enabled,
			"auto_scaling_compute_scale_down_enabled": clusters[i].AutoScaling.Compute.ScaleDownEnabled,
			"auto_scaling_disk_gb_enabled":            clusters[i].AutoScaling.DiskGBEnabled,
			"backup_enabled":                          clusters[i].BackupEnabled,
			"provider_backup_enabled":                 clusters[i].ProviderBackupEnabled,
			"cluster_type":                            clusters[i].ClusterType,
			"connection_strings":                      flattenConnectionStrings(clusters[i].ConnectionStrings),
			"disk_size_gb":                            clusters[i].DiskSizeGB,
			"encryption_at_rest_provider":             clusters[i].EncryptionAtRestProvider,
			"mongo_db_major_version":                  latestClusterModels[clusters[i].Name].MongoDBMajorVersion, // uses 2024-08-05 or above as it has fix for correct value when FCV is active
			"name":                                    clusters[i].Name,
			"num_shards":                              clusters[i].NumShards,
			"mongo_db_version":                        clusters[i].MongoDBVersion,
			"mongo_uri":                               clusters[i].MongoURI,
			"mongo_uri_updated":                       clusters[i].MongoURIUpdated,
			"mongo_uri_with_options":                  clusters[i].MongoURIWithOptions,
			"pit_enabled":                             clusters[i].PitEnabled,
			"paused":                                  clusters[i].Paused,
			"srv_address":                             clusters[i].SrvAddress,
			"state_name":                              clusters[i].StateName,
			"replication_factor":                      clusters[i].ReplicationFactor,
			"provider_auto_scaling_compute_min_instance_size": clusters[i].ProviderSettings.AutoScaling.Compute.MinInstanceSize,
			"provider_auto_scaling_compute_max_instance_size": clusters[i].ProviderSettings.AutoScaling.Compute.MaxInstanceSize,
			"backing_provider_name":                           clusters[i].ProviderSettings.BackingProviderName,
			"provider_disk_iops":                              clusters[i].ProviderSettings.DiskIOPS,
			"provider_disk_type_name":                         clusters[i].ProviderSettings.DiskTypeName,
			"provider_encrypt_ebs_volume":                     clusters[i].ProviderSettings.EncryptEBSVolume,
			"provider_instance_size_name":                     clusters[i].ProviderSettings.InstanceSizeName,
			"provider_name":                                   clusters[i].ProviderSettings.ProviderName,
			"provider_region_name":                            clusters[i].ProviderSettings.RegionName,
			"bi_connector_config":                             flattenBiConnectorConfig(clusters[i].BiConnector),
			"replication_specs":                               flattenReplicationSpecs(clusters[i].ReplicationSpecs),
			"labels":                                          flattenLabels(clusters[i].Labels),
			"tags":                                            flattenTags(clusters[i].Tags),
			"snapshot_backup_policy":                          snapshotBackupPolicy,
			"termination_protection_enabled":                  clusters[i].TerminationProtectionEnabled,
			"version_release_system":                          clusters[i].VersionReleaseSystem,
			"container_id":                                    containerID,
			"redact_client_log_data":                          latestClusterModels[clusters[i].Name].GetRedactClientLogData(),
			"pinned_fcv":                                      advancedcluster.FlattenPinnedFCV(latestClusterModels[clusters[i].Name]),
		}
		results = append(results, result)
	}

	return results
}
