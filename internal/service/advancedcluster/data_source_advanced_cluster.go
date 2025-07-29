package advancedcluster

import (
	"context"
	"fmt"

	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedclustertpf"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"use_replication_spec_per_shard": {
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
	}
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.MongoDBClient)
	connV220240530 := client.AtlasV220240530
	connV2 := client.AtlasV2

	projectID := d.Get("project_id").(string)
	clusterName := d.Get("name").(string)
	useReplicationSpecPerShard := false
	var replicationSpecs []map[string]any

	if v, ok := d.GetOk("use_replication_spec_per_shard"); ok {
		useReplicationSpecPerShard = v.(bool)
	}
	clusterDesc, flexClusterResp, diags := GetClusterDetails(ctx, client, projectID, clusterName)
	if diags.HasError() {
		return diags
	}
	if clusterDesc == nil && flexClusterResp == nil {
		return nil
	}
	if flexClusterResp != nil {
		diags := setFlexFields(d, flexClusterResp)
		if diags.HasError() {
			return diags
		}
		d.SetId(flexClusterResp.GetId())
		return nil
	}

	zoneNameToOldReplicationSpecMeta, err := GetReplicationSpecAttributesFromOldAPI(ctx, projectID, clusterName, connV220240530.ClustersApi)
	if err != nil {
		if apiError, ok := admin20240530.AsError(err); !ok {
			return diag.FromErr(err)
		} else if apiError.GetErrorCode() == "ASYMMETRIC_SHARD_UNSUPPORTED" && !useReplicationSpecPerShard {
			return diag.FromErr(fmt.Errorf("please add `use_replication_spec_per_shard = true` to your data source configuration to enable asymmetric shard support. Refer to documentation for more details. %s", err))
		} else if apiError.GetErrorCode() != "ASYMMETRIC_SHARD_UNSUPPORTED" {
			return diag.FromErr(fmt.Errorf(errorRead, clusterName, err))
		}
	}
	diags = setRootFields(d, clusterDesc, false)
	if diags.HasError() {
		return diags
	}

	if !useReplicationSpecPerShard {
		replicationSpecs, err = FlattenAdvancedReplicationSpecsOldShardingConfig(ctx, clusterDesc.GetReplicationSpecs(), zoneNameToOldReplicationSpecMeta, d.Get("replication_specs").([]any), d, connV2)
		if err != nil {
			return diag.FromErr(fmt.Errorf(ErrorClusterAdvancedSetting, "replication_specs", clusterName, err))
		}
	} else {
		replicationSpecs, err = flattenAdvancedReplicationSpecsDS(ctx, clusterDesc.GetReplicationSpecs(), zoneNameToOldReplicationSpecMeta, d, connV2)
		if err != nil {
			return diag.FromErr(fmt.Errorf(ErrorClusterAdvancedSetting, "replication_specs", clusterName, err))
		}
	}

	if err := d.Set("replication_specs", replicationSpecs); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorClusterAdvancedSetting, "replication_specs", clusterName, err))
	}

	processArgs, _, err := connV2.ClustersApi.GetClusterAdvancedConfiguration(ctx, projectID, clusterName).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(ErrorAdvancedConfRead, "", clusterName, err))
	}
	advConfigAttr := flattenProcessArgs(&advancedclustertpf.ProcessArgs{
		ArgsDefault:           processArgs,
		ClusterAdvancedConfig: clusterDesc.AdvancedConfiguration,
	})
	if err := d.Set("advanced_configuration", advConfigAttr); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorClusterAdvancedSetting, "advanced_configuration", clusterName, err))
	}

	d.SetId(clusterDesc.GetId())
	return nil
}
