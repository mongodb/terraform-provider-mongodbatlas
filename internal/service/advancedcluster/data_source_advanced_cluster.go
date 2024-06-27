package advancedcluster

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"use_independent_shards": {
				Type:     schema.TypeBool,
				Optional: true,
			},
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
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"encryption_at_rest_provider": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"labels": {
				Type:       schema.TypeSet,
				Computed:   true,
				Deprecated: fmt.Sprintf(constant.DeprecationParamByDateWithReplacement, "September 2024", "tags"),
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
				Required: true,
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
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": { // new API only
							Type:     schema.TypeString,
							Computed: true,
						},
						"external_id": { // new API only
							Type:     schema.TypeString,
							Computed: true,
						},
						"num_shards": {
							Type:     schema.TypeInt,
							Computed: true,
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
		},
	}
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	connLatest := meta.(*config.MongoDBClient).AtlasV2Preview

	projectID := d.Get("project_id").(string)
	clusterName := d.Get("name").(string)
	// useIndependentShards := false

	// if v, ok := d.GetOk("use_independent_shards"); ok {
	// 	useIndependentShards = v.(bool)
	// }

	// if !useIndependentShards {
	cluster, resp, err := connV2.ClustersApi.GetCluster(ctx, projectID, clusterName).Execute()

	// } else {
	//  cluster, resp, err := connPreview.ClustersApi.GetCluster(ctx, projectID, clusterName).Execute() //var cluster *adminPreview.ClusterDescription20240710

	// }

	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil
		}
		return diag.FromErr(fmt.Errorf(errorRead, clusterName, err))
	}

	if err := d.Set("backup_enabled", cluster.GetBackupEnabled()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorClusterAdvancedSetting, "backup_enabled", clusterName, err))
	}

	if err := d.Set("bi_connector_config", flattenBiConnectorConfig(convertBiConnectToLatest(cluster.GetBiConnector()))); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorClusterAdvancedSetting, "bi_connector_config", clusterName, err))
	}

	if err := d.Set("cluster_type", cluster.GetClusterType()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorClusterAdvancedSetting, "cluster_type", clusterName, err))
	}

	if err := d.Set("connection_strings", flattenConnectionStrings(convertConnectionStringToLatest(cluster.GetConnectionStrings()))); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorClusterAdvancedSetting, "connection_strings", clusterName, err))
	}

	if err := d.Set("create_date", conversion.TimePtrToStringPtr(cluster.CreateDate)); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorClusterAdvancedSetting, "create_date", clusterName, err))
	}

	if err := d.Set("disk_size_gb", cluster.GetDiskSizeGB()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorClusterAdvancedSetting, "disk_size_gb", clusterName, err))
	}

	if err := d.Set("encryption_at_rest_provider", cluster.GetEncryptionAtRestProvider()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorClusterAdvancedSetting, "encryption_at_rest_provider", clusterName, err))
	}

	if err := d.Set("labels", flattenLabels(convertLabelsToLatest(cluster.GetLabels()))); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorClusterAdvancedSetting, "labels", clusterName, err))
	}

	if err := d.Set("tags", flattenTags(*convertTagsToLatest(cluster.Tags))); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorClusterAdvancedSetting, "tags", clusterName, err))
	}

	if err := d.Set("mongo_db_major_version", cluster.GetMongoDBMajorVersion()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorClusterAdvancedSetting, "mongo_db_major_version", clusterName, err))
	}

	if err := d.Set("mongo_db_version", cluster.GetMongoDBVersion()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorClusterAdvancedSetting, "mongo_db_version", clusterName, err))
	}

	if err := d.Set("name", cluster.GetName()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorClusterAdvancedSetting, "name", clusterName, err))
	}

	if err := d.Set("paused", cluster.GetPaused()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorClusterAdvancedSetting, "paused", clusterName, err))
	}

	if err := d.Set("pit_enabled", cluster.GetPitEnabled()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorClusterAdvancedSetting, "pit_enabled", clusterName, err))
	}

	replicationSpecs, err := FlattenAdvancedReplicationSpecsOldSDK(ctx, cluster.GetReplicationSpecs(), d.Get("replication_specs").([]any), d, connLatest)
	if err != nil {
		return diag.FromErr(fmt.Errorf(ErrorClusterAdvancedSetting, "replication_specs", clusterName, err))
	}

	if err := d.Set("replication_specs", replicationSpecs); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorClusterAdvancedSetting, "replication_specs", clusterName, err))
	}

	if err := d.Set("root_cert_type", cluster.GetRootCertType()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorClusterAdvancedSetting, "state_name", clusterName, err))
	}

	if err := d.Set("state_name", cluster.GetStateName()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorClusterAdvancedSetting, "state_name", clusterName, err))
	}
	if err := d.Set("termination_protection_enabled", cluster.GetTerminationProtectionEnabled()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorClusterAdvancedSetting, "termination_protection_enabled", clusterName, err))
	}
	if err := d.Set("version_release_system", cluster.GetVersionReleaseSystem()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorClusterAdvancedSetting, "version_release_system", clusterName, err))
	}
	if err := d.Set("global_cluster_self_managed_sharding", cluster.GetGlobalClusterSelfManagedSharding()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorClusterAdvancedSetting, "global_cluster_self_managed_sharding", clusterName, err))
	}

	// TODO: update to use connLatest to call below API
	processArgs, _, err := connV2.ClustersApi.GetClusterAdvancedConfiguration(ctx, projectID, clusterName).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(ErrorAdvancedConfRead, clusterName, err))
	}

	if err := d.Set("advanced_configuration", flattenProcessArgs(processArgs)); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorClusterAdvancedSetting, "advanced_configuration", clusterName, err))
	}

	d.SetId(cluster.GetId())
	return nil
}
