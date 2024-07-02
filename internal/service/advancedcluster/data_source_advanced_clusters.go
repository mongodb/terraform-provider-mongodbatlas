package advancedcluster

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20240530002/admin"
)

func PluralDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePluralRead,
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
				},
			},
		},
	}
}

func dataSourcePluralRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)
	d.SetId(id.UniqueId())

	list, resp, err := connV2.ClustersApi.ListClusters(ctx, projectID).Execute()
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil
		}
		return diag.FromErr(fmt.Errorf("error reading advanced cluster list for project(%s): %s", projectID, err))
	}
	if err := d.Set("results", flattenAdvancedClusters(ctx, connV2, list.GetResults(), d)); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorClusterAdvancedSetting, "results", d.Id(), err))
	}

	return nil
}

func flattenAdvancedClusters(ctx context.Context, connV2 *admin.APIClient, clusters []admin.AdvancedClusterDescription, d *schema.ResourceData) []map[string]any {
	results := make([]map[string]any, 0, len(clusters))
	for i := range clusters {
		cluster := &clusters[i]
		processArgs, _, err := connV2.ClustersApi.GetClusterAdvancedConfiguration(ctx, cluster.GetGroupId(), cluster.GetName()).Execute()
		if err != nil {
			log.Printf("[WARN] Error setting `advanced_configuration` for the cluster(%s): %s", cluster.GetId(), err)
		}
		replicationSpecs, err := FlattenAdvancedReplicationSpecs(ctx, cluster.GetReplicationSpecs(), nil, d, connV2)
		if err != nil {
			log.Printf("[WARN] Error setting `replication_specs` for the cluster(%s): %s", cluster.GetId(), err)
		}

		result := map[string]any{
			"advanced_configuration":               flattenProcessArgs(processArgs),
			"backup_enabled":                       cluster.GetBackupEnabled(),
			"bi_connector_config":                  flattenBiConnectorConfig(cluster.GetBiConnector()),
			"cluster_type":                         cluster.GetClusterType(),
			"create_date":                          conversion.TimePtrToStringPtr(cluster.CreateDate),
			"connection_strings":                   flattenConnectionStrings(cluster.GetConnectionStrings()),
			"disk_size_gb":                         cluster.GetDiskSizeGB(),
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
		}
		results = append(results, result)
	}
	return results
}
