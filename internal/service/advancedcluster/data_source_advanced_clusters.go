package advancedcluster

import (
	"context"
	"fmt"
	"log"

	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
	"go.mongodb.org/atlas-sdk/v20250312002/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/flexcluster"
)

const (
	errorListRead = "error reading advanced cluster list for project(%s): %s"
)

func PluralDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePluralRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"use_replication_spec_per_shard": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"advanced_configuration": SchemaAdvancedConfigDS(),
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
										Optional: true,
										Computed: true,
									},
									"read_preference": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"cluster_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_strings": SchemaConnectionStrings(),
						"create_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_size_gb": {
							Type:       schema.TypeFloat,
							Computed:   true,
							Deprecated: DeprecationMsgOldSchema,
						},
						"encryption_at_rest_provider": {
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
						"tags": &DSTagsSchema,
						"mongo_db_major_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mongo_db_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"paused": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"pit_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"replication_specs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:       schema.TypeString,
										Computed:   true,
										Deprecated: DeprecationMsgOldSchema,
									},
									"zone_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"external_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"num_shards": {
										Type:       schema.TypeInt,
										Computed:   true,
										Deprecated: DeprecationMsgOldSchema,
									},
									"region_configs": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"analytics_specs": schemaSpecs(),
												"auto_scaling": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"disk_gb_enabled": {
																Type:     schema.TypeBool,
																Computed: true,
															},
															"compute_enabled": {
																Type:     schema.TypeBool,
																Computed: true,
															},
															"compute_scale_down_enabled": {
																Type:     schema.TypeBool,
																Computed: true,
															},
															"compute_min_instance_size": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"compute_max_instance_size": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
												"analytics_auto_scaling": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"disk_gb_enabled": {
																Type:     schema.TypeBool,
																Computed: true,
															},
															"compute_enabled": {
																Type:     schema.TypeBool,
																Computed: true,
															},
															"compute_scale_down_enabled": {
																Type:     schema.TypeBool,
																Computed: true,
															},
															"compute_min_instance_size": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"compute_max_instance_size": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
												"backing_provider_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"electable_specs": schemaSpecs(),
												"priority": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"provider_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"read_only_specs": schemaSpecs(),
												"region_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"container_id": {
										Type: schema.TypeMap,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed: true,
									},
									"zone_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"root_cert_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state_name": {
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
						"global_cluster_self_managed_sharding": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"replica_set_scaling_strategy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"redact_client_log_data": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"config_server_management_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"config_server_type": {
							Type:     schema.TypeString,
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
	connV220240530 := meta.(*config.MongoDBClient).AtlasV220240530
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)
	useReplicationSpecPerShard := false

	d.SetId(id.UniqueId())

	if v, ok := d.GetOk("use_replication_spec_per_shard"); ok {
		useReplicationSpecPerShard = v.(bool)
	}

	list, resp, err := connV2.ClustersApi.ListClusters(ctx, projectID).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) {
			return nil
		}
		return diag.FromErr(fmt.Errorf(errorListRead, projectID, err))
	}
	results, diags := flattenAdvancedClusters(ctx, connV220240530, connV2, list.GetResults(), d, useReplicationSpecPerShard)
	if len(diags) > 0 {
		return diags
	}

	listFlexClusters, err := flexcluster.ListFlexClusters(ctx, projectID, connV2.FlexClustersApi)
	if err != nil {
		if validate.StatusNotFound(resp) {
			return nil
		}
		return diag.FromErr(fmt.Errorf(errorListRead, projectID, err))
	}
	results = append(results, flexcluster.FlattenFlexClustersToAdvancedClusters(listFlexClusters)...)

	if err := d.Set("results", results); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorClusterAdvancedSetting, "results", d.Id(), err))
	}

	return nil
}

func flattenAdvancedClusters(ctx context.Context, connV220240530 *admin20240530.APIClient, connV2 *admin.APIClient, clusters []admin.ClusterDescription20240805, d *schema.ResourceData, useReplicationSpecPerShard bool) ([]map[string]any, diag.Diagnostics) {
	results := make([]map[string]any, 0, len(clusters))
	for i := range clusters {
		cluster := &clusters[i]
		processArgs20240530, _, err := connV220240530.ClustersApi.GetClusterAdvancedConfiguration(ctx, cluster.GetGroupId(), cluster.GetName()).Execute()
		if err != nil {
			log.Printf("[WARN] Error setting `advanced_configuration` for the cluster(%s): %s", cluster.GetId(), err)
		}
		processArgs, _, err := connV2.ClustersApi.GetClusterAdvancedConfiguration(ctx, cluster.GetGroupId(), cluster.GetName()).Execute()
		if err != nil {
			log.Printf("[WARN] Error setting `advanced_configuration` for the cluster(%s): %s", cluster.GetId(), err)
		}

		zoneNameToOldReplicationSpecMeta, err := GetReplicationSpecAttributesFromOldAPI(ctx, cluster.GetGroupId(), cluster.GetName(), connV220240530.ClustersApi)
		if err != nil {
			errNotFound := admin20240530.IsErrorCode(err, "CLUSTER_NOT_FOUND")
			errAsymmetricUnsupported := admin20240530.IsErrorCode(err, "ASYMMETRIC_SHARD_UNSUPPORTED")
			if errNotFound || (errAsymmetricUnsupported && !useReplicationSpecPerShard) {
				continue
			}
			if !errAsymmetricUnsupported {
				return nil, diag.FromErr(err)
			}
		}

		var replicationSpecs []map[string]any

		if !useReplicationSpecPerShard {
			replicationSpecs, err = FlattenAdvancedReplicationSpecsOldShardingConfig(ctx, cluster.GetReplicationSpecs(), zoneNameToOldReplicationSpecMeta, nil, d, connV2)
			if err != nil {
				log.Printf("[WARN] Error setting `replication_specs` for the cluster(%s): %s", cluster.GetId(), err)
			}
		} else {
			replicationSpecs, err = flattenAdvancedReplicationSpecsDS(ctx, cluster.GetReplicationSpecs(), zoneNameToOldReplicationSpecMeta, d, connV2)
			if err != nil {
				log.Printf("[WARN] Error setting `replication_specs` for the cluster(%s): %s", cluster.GetId(), err)
			}
		}

		result := map[string]any{
			"advanced_configuration":               flattenProcessArgs(processArgs20240530, processArgs, cluster.AdvancedConfiguration),
			"backup_enabled":                       cluster.GetBackupEnabled(),
			"bi_connector_config":                  flattenBiConnectorConfig(cluster.BiConnector),
			"cluster_type":                         cluster.GetClusterType(),
			"create_date":                          conversion.TimePtrToStringPtr(cluster.CreateDate),
			"connection_strings":                   flattenConnectionStrings(cluster.GetConnectionStrings()),
			"disk_size_gb":                         GetDiskSizeGBFromReplicationSpec(cluster),
			"encryption_at_rest_provider":          cluster.GetEncryptionAtRestProvider(),
			"labels":                               flattenLabels(cluster.GetLabels()),
			"tags":                                 conversion.FlattenTags(cluster.GetTags()),
			"mongo_db_major_version":               cluster.GetMongoDBMajorVersion(),
			"mongo_db_version":                     cluster.GetMongoDBVersion(),
			"name":                                 cluster.GetName(),
			"paused":                               cluster.GetPaused(),
			"pit_enabled":                          cluster.GetPitEnabled(),
			"replication_specs":                    replicationSpecs,
			"root_cert_type":                       cluster.GetRootCertType(),
			"state_name":                           cluster.GetStateName(),
			"termination_protection_enabled":       cluster.GetTerminationProtectionEnabled(),
			"version_release_system":               cluster.GetVersionReleaseSystem(),
			"global_cluster_self_managed_sharding": cluster.GetGlobalClusterSelfManagedSharding(),
			"replica_set_scaling_strategy":         cluster.GetReplicaSetScalingStrategy(),
			"redact_client_log_data":               cluster.GetRedactClientLogData(),
			"config_server_management_mode":        cluster.GetConfigServerManagementMode(),
			"config_server_type":                   cluster.GetConfigServerType(),
			"pinned_fcv":                           FlattenPinnedFCV(cluster),
		}
		results = append(results, result)
	}
	return results, nil
}
