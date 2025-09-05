package cluster

import (
	"context"
	"fmt"
	"net/http"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
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
			"pit_enabled": {
				Type:     schema.TypeBool,
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
			"provider_encrypt_ebs_volume_flag": {
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
	}
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).Atlas
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	connV220240530 := meta.(*config.MongoDBClient).AtlasV220240530
	projectID := d.Get("project_id").(string)
	clusterName := d.Get("name").(string)

	cluster, resp, err := conn.Clusters.Get(ctx, projectID, clusterName)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorClusterRead, clusterName, err))
	}

	if err := d.Set("auto_scaling_compute_enabled", cluster.AutoScaling.Compute.Enabled); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "auto_scaling_compute_enabled", clusterName, err))
	}

	if err := d.Set("auto_scaling_compute_scale_down_enabled", cluster.AutoScaling.Compute.ScaleDownEnabled); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "auto_scaling_compute_scale_down_enabled", clusterName, err))
	}

	if err := d.Set("provider_auto_scaling_compute_min_instance_size", cluster.ProviderSettings.AutoScaling.Compute.MinInstanceSize); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "provider_auto_scaling_compute_min_instance_size", clusterName, err))
	}

	if err := d.Set("provider_auto_scaling_compute_max_instance_size", cluster.ProviderSettings.AutoScaling.Compute.MaxInstanceSize); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "provider_auto_scaling_compute_max_instance_size", clusterName, err))
	}

	if err := d.Set("auto_scaling_disk_gb_enabled", cluster.AutoScaling.DiskGBEnabled); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "auto_scaling_disk_gb_enabled", clusterName, err))
	}

	if err := d.Set("backup_enabled", cluster.BackupEnabled); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "backup_enabled", clusterName, err))
	}

	if err := d.Set("pit_enabled", cluster.PitEnabled); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "pit_enabled", clusterName, err))
	}

	if err := d.Set("provider_backup_enabled", cluster.ProviderBackupEnabled); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "provider_backup_enabled", clusterName, err))
	}

	if err := d.Set("cluster_type", cluster.ClusterType); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "cluster_type", clusterName, err))
	}

	if err := d.Set("connection_strings", flattenConnectionStrings(cluster.ConnectionStrings)); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "connection_strings", clusterName, err))
	}

	if err := d.Set("disk_size_gb", cluster.DiskSizeGB); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "disk_size_gb", clusterName, err))
	}

	if err := d.Set("encryption_at_rest_provider", cluster.EncryptionAtRestProvider); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "encryption_at_rest_provider", clusterName, err))
	}

	// Avoid Global Cluster issues. (NumShards is not present in Global Clusters)
	if cluster.NumShards != nil {
		if err := d.Set("num_shards", cluster.NumShards); err != nil {
			return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "num_shards", clusterName, err))
		}
	}

	if err := d.Set("mongo_db_version", cluster.MongoDBVersion); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "mongo_db_version", clusterName, err))
	}

	if err := d.Set("mongo_uri", cluster.MongoURI); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "mongo_uri", clusterName, err))
	}

	if err := d.Set("mongo_uri_updated", cluster.MongoURIUpdated); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "mongo_uri_updated", clusterName, err))
	}

	if err := d.Set("mongo_uri_with_options", cluster.MongoURIWithOptions); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "mongo_uri_with_options", clusterName, err))
	}

	if err := d.Set("paused", cluster.Paused); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "paused", clusterName, err))
	}

	if err := d.Set("srv_address", cluster.SrvAddress); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "srv_address", clusterName, err))
	}

	if err := d.Set("state_name", cluster.StateName); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "state_name", clusterName, err))
	}

	if err := d.Set("bi_connector_config", flattenBiConnectorConfig(cluster.BiConnector)); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "bi_connector_config", clusterName, err))
	}

	if cluster.ProviderSettings != nil {
		flattenProviderSettings(d, cluster.ProviderSettings, clusterName)
	}

	if err := d.Set("replication_specs", flattenReplicationSpecs(cluster.ReplicationSpecs)); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "replication_specs", clusterName, err))
	}

	if err := d.Set("replication_factor", cluster.ReplicationFactor); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "replication_factor", clusterName, err))
	}

	if err := d.Set("labels", flattenLabels(cluster.Labels)); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "labels", clusterName, err))
	}

	if err := d.Set("tags", flattenTags(cluster.Tags)); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "tags", clusterName, err))
	}

	if err := d.Set("termination_protection_enabled", cluster.TerminationProtectionEnabled); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "termination_protection_enabled", clusterName, err))
	}

	if err := d.Set("version_release_system", cluster.VersionReleaseSystem); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "version_release_system", clusterName, err))
	}

	if cluster.ProviderSettings != nil && cluster.ProviderSettings.ProviderName != "TENANT" {
		containers, _, err := conn.Containers.List(ctx, projectID,
			&matlas.ContainersListOptions{ProviderName: cluster.ProviderSettings.ProviderName})
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorClusterRead, clusterName, err))
		}

		if err := d.Set("container_id", getContainerID(containers, cluster)); err != nil {
			return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "container_id", clusterName, err))
		}
	}

	/*
		Get the advaced configuration options and set up to the terraform state
	*/
	processArgs20240530, _, err := connV220240530.ClustersApi.GetClusterAdvancedConfiguration(ctx, projectID, clusterName).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorAdvancedConfRead, advancedcluster.V20240530, clusterName, err))
	}
	processArgs, _, err := connV2.ClustersApi.GetProcessArgs(ctx, projectID, clusterName).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorAdvancedConfRead, "", clusterName, err))
	}

	p := &ProcessArgs{
		argsDefault:           processArgs,
		argsLegacy:            processArgs20240530,
		clusterAdvancedConfig: cluster.AdvancedConfiguration,
	}

	if err := d.Set("advanced_configuration", flattenProcessArgs(p)); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "advanced_configuration", clusterName, err))
	}

	// Get the snapshot policy and set the data
	snapshotBackupPolicy, err := flattenCloudProviderSnapshotBackupPolicy(ctx, d, conn, projectID, clusterName)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("snapshot_backup_policy", snapshotBackupPolicy); err != nil {
		return diag.FromErr(err)
	}

	latestClusterModel, err := newAtlasGet(ctx, connV2, projectID, clusterName)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterRead, clusterName, err))
	}

	if err := d.Set("mongo_db_major_version", latestClusterModel.MongoDBMajorVersion); err != nil { // uses 2024-08-05 or above as it has fix for correct value when FCV is active
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "mongo_db_major_version", clusterName, err))
	}

	if err := d.Set("redact_client_log_data", latestClusterModel.GetRedactClientLogData()); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "redact_client_log_data", clusterName, err))
	}

	if err := d.Set("pinned_fcv", advancedcluster.FlattenPinnedFCV(latestClusterModel)); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "pinned_fcv", clusterName, err))
	}

	d.SetId(cluster.ID)

	return nil
}
