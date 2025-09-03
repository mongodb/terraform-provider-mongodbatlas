package federateddatabaseinstance

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasFederatedDatabaseInstanceRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hostnames": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cloud_provider_config": cloudProviderConfig(true),
			"data_process_region": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloud_provider": {
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
			"storage_databases": schemaFederatedDatabaseInstanceDatabasesDataSource(),
			"storage_stores":    schemaFederatedDatabaseInstanceStoresDataSource(),
		},
	}
}

func schemaFederatedDatabaseInstanceDatabasesDataSource() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"collections": {
					Type:     schema.TypeSet,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"name": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"data_sources": {
								Type:     schema.TypeSet,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"store_name": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"dataset_name": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"default_format": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"path": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"allow_insecure": {
											Type:     schema.TypeBool,
											Computed: true,
										},
										"database": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"database_regex": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"collection": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"collection_regex": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"provenance_field_name": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"urls": {
											Type:     schema.TypeList,
											Computed: true,
											Elem: &schema.Schema{
												Type: schema.TypeString,
											},
										},
									},
								},
							},
						},
					},
				},
				"views": {
					Type:     schema.TypeSet,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"name": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"source": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"pipeline": {
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},
				"max_wildcard_collections": {
					Type:     schema.TypeInt,
					Computed: true,
				},
			},
		},
	}
}

func schemaFederatedDatabaseInstanceStoresDataSource() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"provider": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"region": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"bucket": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"cluster_name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"project_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"prefix": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"delimiter": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"include_tags": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"allow_insecure": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"additional_storage_classes": {
					Type:     schema.TypeList,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"public": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"default_format": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"urls": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"read_preference": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"mode": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"max_staleness_seconds": {
								Type:     schema.TypeInt,
								Computed: true,
							},
							"tag_sets": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"tags": {
											Type:     schema.TypeList,
											Computed: true,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"name": {
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

func dataSourceMongoDBAtlasFederatedDatabaseInstanceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	projectID := d.Get("project_id").(string)
	name := d.Get("name").(string)

	dataFederationInstance, _, err := connV2.DataFederationApi.GetDataFederation(ctx, projectID, name).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("couldn't import data lake(%s) for project (%s), error: %s", name, projectID, err))
	}

	if err := d.Set("project_id", projectID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `project_id` for data lakes (%s): %s", d.Id(), err))
	}

	if err := d.Set("name", dataFederationInstance.Name); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `name` for data lakes (%s): %s", d.Id(), err))
	}

	if cloudProviderField := flattenCloudProviderConfig(d, dataFederationInstance.CloudProviderConfig); cloudProviderField != nil {
		if err = d.Set("cloud_provider_config", cloudProviderField); err != nil {
			return diag.FromErr(fmt.Errorf(errorFederatedDatabaseInstanceSetting, "cloud_provider_config", name, err))
		}
	}

	if err := d.Set("data_process_region", flattenDataProcessRegion(dataFederationInstance.DataProcessRegion)); err != nil {
		return diag.FromErr(fmt.Errorf(errorFederatedDatabaseInstanceSetting, "data_process_region", name, err))
	}

	if storageDatabaseField := flattenDataFederationDatabase(dataFederationInstance.Storage.GetDatabases()); storageDatabaseField != nil {
		if err := d.Set("storage_databases", storageDatabaseField); err != nil {
			return diag.FromErr(fmt.Errorf(errorFederatedDatabaseInstanceSetting, "storage_databases", name, err))
		}
	}

	if err := d.Set("storage_stores", flattenDataFederationStores(dataFederationInstance.Storage.GetStores())); err != nil {
		return diag.FromErr(fmt.Errorf(errorFederatedDatabaseInstanceSetting, "storage_stores", name, err))
	}

	if err := d.Set("state", dataFederationInstance.State); err != nil {
		return diag.FromErr(fmt.Errorf(errorFederatedDatabaseInstanceSetting, "state", name, err))
	}

	if err := d.Set("hostnames", dataFederationInstance.Hostnames); err != nil {
		return diag.FromErr(fmt.Errorf(errorFederatedDatabaseInstanceSetting, "hostnames", name, err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id": projectID,
		"name":       dataFederationInstance.GetName(),
	}))

	return nil
}
