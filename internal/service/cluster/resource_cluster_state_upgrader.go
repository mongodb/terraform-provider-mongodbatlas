package cluster

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
)

func ResourceClusterResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_scaling_disk_gb_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"backup_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"bi_connector": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeString,
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
				Optional: true,
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
					},
				},
			},
			"disk_size_gb": {
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
			},
			"encryption_at_rest_provider": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mongo_db_major_version": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				StateFunc: formatMongoDBMajorVersion,
			},
			"num_shards": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"provider_backup_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"provider_instance_size_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"provider_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"pit_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"backing_provider_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"provider_disk_iops": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"provider_disk_type_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"provider_encrypt_ebs_volume": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"provider_region_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"provider_volume_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"replication_factor": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"replication_specs": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"num_shards": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"regions_config": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region_name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"electable_nodes": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"priority": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"read_only_nodes": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  0,
									},
									"analytics_nodes": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  0,
									},
								},
							},
						},
						"zone_name": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "ZoneName managed by Terraform",
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
			"advanced_configuration": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"fail_index_key_too_long": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"javascript_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"minimum_enabled_tls_protocol": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"no_table_scan": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"oplog_size_mb": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"sample_size_bi_connector": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"sample_refresh_interval_bi_connector": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"labels": {
				Type:     schema.TypeSet,
				Optional: true,
				Set:      advancedcluster.HashFunctionForKeyValuePair,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"snapshot_backup_policy": computedCloudProviderSnapshotBackupPolicySchema(),
		},
	}
}

func resourceMongoDBAtlasClusterStateUpgradeV0(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	log.Println("[INFO] Found MongoDB Cluser state v0; migrating to v1")

	return MigrateAdvancedConfiguration(rawState), nil
}

func MigrateAdvancedConfiguration(rawState map[string]any) map[string]any {
	rawState["advanced_configuration"] = []any{}

	return rawState
}
