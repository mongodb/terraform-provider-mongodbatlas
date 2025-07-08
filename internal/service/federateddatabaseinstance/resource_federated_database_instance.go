package federateddatabaseinstance

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/atlas-sdk/v20250312004/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	errorFederatedDatabaseInstanceCreate  = "error creating MongoDB Atlas Federated Database Instace: %s"
	errorFederatedDatabaseInstanceRead    = "error reading MongoDB Atlas Federated Database Instace (%s): %s"
	errorFederatedDatabaseInstanceDelete  = "error deleting MongoDB Atlas Federated Database Instace (%s): %s"
	errorFederatedDatabaseInstanceUpdate  = "error updating MongoDB Atlas Federated Database Instace (%s): %s"
	errorFederatedDatabaseInstanceSetting = "error setting `%s` for MongoDB Atlas Federated Database Instace (%s): %s"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate,
		ReadContext:   resourceRead,
		UpdateContext: resourceUpdate,
		DeleteContext: resourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceImport,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"cloud_provider_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"aws": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"role_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"test_s3_bucket": {
										Type:     schema.TypeString,
										Required: true,
									},
									"iam_assumed_role_arn": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"iam_user_arn": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"external_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"azure": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"role_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"atlas_app_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"service_principal_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"tenant_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"data_process_region": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloud_provider": {
							Type:     schema.TypeString,
							Required: true,
						},
						"region": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"storage_databases": schemaFederatedDatabaseInstanceDatabases(),
			"storage_stores":    schemaFederatedDatabaseInstanceStores(),
		},
	}
}

func schemaFederatedDatabaseInstanceDatabases() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Computed: true,
					Optional: true,
				},
				"collections": {
					Type:     schema.TypeSet,
					Computed: true,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"name": {
								Type:     schema.TypeString,
								Computed: true,
								Optional: true,
							},
							"data_sources": {
								Type:     schema.TypeSet,
								Optional: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"store_name": {
											Type:     schema.TypeString,
											Computed: true,
											Optional: true,
										},
										"dataset_name": {
											Type:     schema.TypeString,
											Computed: true,
											Optional: true,
										},
										"default_format": {
											Type:     schema.TypeString,
											Computed: true,
											Optional: true,
										},
										"path": {
											Type:     schema.TypeString,
											Computed: true,
											Optional: true,
										},
										"allow_insecure": {
											Type:     schema.TypeBool,
											Computed: true,
											Optional: true,
										},
										"database": {
											Type:     schema.TypeString,
											Computed: true,
											Optional: true,
										},
										"database_regex": {
											Type:     schema.TypeString,
											Computed: true,
											Optional: true,
										},
										"collection": {
											Type:     schema.TypeString,
											Computed: true,
											Optional: true,
										},
										"collection_regex": {
											Type:     schema.TypeString,
											Computed: true,
											Optional: true,
										},
										"provenance_field_name": {
											Type:     schema.TypeString,
											Computed: true,
											Optional: true,
										},
										"urls": {
											Type:     schema.TypeList,
											Computed: true,
											Optional: true,
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
					Optional: true,
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

func schemaFederatedDatabaseInstanceStores() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Computed: true,
					Optional: true,
				},
				"provider": {
					Type:     schema.TypeString,
					Computed: true,
					Optional: true,
				},
				"region": {
					Type:     schema.TypeString,
					Computed: true,
					Optional: true,
				},
				"bucket": {
					Type:     schema.TypeString,
					Computed: true,
					Optional: true,
				},
				"cluster_name": {
					Type:     schema.TypeString,
					Computed: true,
					Optional: true,
				},
				"project_id": {
					Type:     schema.TypeString,
					Computed: true,
					Optional: true,
				},
				"prefix": {
					Type:     schema.TypeString,
					Computed: true,
					Optional: true,
				},
				"delimiter": {
					Type:     schema.TypeString,
					Computed: true,
					Optional: true,
				},
				"include_tags": {
					Type:     schema.TypeBool,
					Computed: true,
					Optional: true,
				},
				"allow_insecure": {
					Type:     schema.TypeBool,
					Computed: true,
					Optional: true,
				},
				"additional_storage_classes": {
					Type:     schema.TypeList,
					Computed: true,
					Optional: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"public": {
					Type:     schema.TypeString,
					Computed: true,
					Optional: true,
				},
				"default_format": {
					Type:     schema.TypeString,
					Computed: true,
					Optional: true,
				},
				"urls": {
					Type:     schema.TypeList,
					Computed: true,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"read_preference": {
					Type:     schema.TypeList,
					MaxItems: 1,
					Computed: true,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"mode": {
								Type:     schema.TypeString,
								Computed: true,
								Optional: true,
							},
							"max_staleness_seconds": {
								Type:     schema.TypeInt,
								Computed: true,
								Optional: true,
							},
							"tag_sets": {
								Type:     schema.TypeList,
								Computed: true,
								Optional: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"tags": {
											Type:     schema.TypeList,
											Required: true,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"name": {
														Type:     schema.TypeString,
														Computed: true,
														Optional: true,
													},
													"value": {
														Type:     schema.TypeString,
														Computed: true,
														Optional: true,
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

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	projectID := d.Get("project_id").(string)
	name := d.Get("name").(string)

	if _, _, err := connV2.DataFederationApi.CreateFederatedDatabase(ctx, projectID, &admin.DataLakeTenant{
		Name:                conversion.StringPtr(name),
		CloudProviderConfig: newCloudProviderConfig(d),
		DataProcessRegion:   newDataProcessRegion(d),
		Storage:             newDataFederationStorage(d),
	}).Execute(); err != nil {
		return diag.FromErr(fmt.Errorf(errorFederatedDatabaseInstanceCreate, err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id": projectID,
		"name":       name,
	}))

	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	name := ids["name"]

	dataFederationInstance, resp, err := connV2.DataFederationApi.GetFederatedDatabase(ctx, projectID, name).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorFederatedDatabaseInstanceRead, name, err))
	}

	if val, ok := dataFederationInstance.GetCloudProviderConfigOk(); ok {
		if cloudProviderField := flattenCloudProviderConfig(val); cloudProviderField != nil {
			if err = d.Set("cloud_provider_config", cloudProviderField); err != nil {
				return diag.FromErr(fmt.Errorf(errorFederatedDatabaseInstanceSetting, "cloud_provider_config", name, err))
			}
		}
	}

	if val, ok := dataFederationInstance.GetDataProcessRegionOk(); ok {
		if dataProcessRegionField := flattenDataProcessRegion(val); dataProcessRegionField != nil {
			if err := d.Set("data_process_region", dataProcessRegionField); err != nil {
				return diag.FromErr(fmt.Errorf(errorFederatedDatabaseInstanceSetting, "data_process_region", name, err))
			}
		}
	}

	if err := d.Set("state", dataFederationInstance.GetState()); err != nil {
		return diag.FromErr(fmt.Errorf(errorFederatedDatabaseInstanceSetting, "state", name, err))
	}

	if err := d.Set("hostnames", dataFederationInstance.GetHostnames()); err != nil {
		return diag.FromErr(fmt.Errorf(errorFederatedDatabaseInstanceSetting, "hostnames", name, err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id": projectID,
		"name":       name,
	}))

	return nil
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	name := ids["name"]

	dataLakeTenant := &admin.DataLakeTenant{
		Name:                conversion.StringPtr(name),
		CloudProviderConfig: newCloudProviderConfig(d),
		DataProcessRegion:   newDataProcessRegion(d),
		Storage:             newDataFederationStorage(d),
	}

	if _, _, err := connV2.DataFederationApi.UpdateFederatedDatabaseWithParams(ctx, &admin.UpdateFederatedDatabaseApiParams{
		GroupId:            projectID,
		TenantName:         name,
		SkipRoleValidation: admin.PtrBool(false),
		DataLakeTenant:     dataLakeTenant,
	}).Execute(); err != nil {
		return diag.FromErr(fmt.Errorf(errorFederatedDatabaseInstanceUpdate, name, err))
	}

	return resourceRead(ctx, d, meta)
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	name := ids["name"]

	if _, err := connV2.DataFederationApi.DeleteFederatedDatabase(ctx, projectID, name).Execute(); err != nil {
		return diag.FromErr(fmt.Errorf(errorFederatedDatabaseInstanceDelete, name, err))
	}

	return nil
}

func resourceImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	projectID, name, s3Bucket, err := splitDataFederatedInstanceImportID(d.Id())
	if err != nil {
		return nil, err
	}

	// test_s3_bucket is not part of the API response
	if s3Bucket != "" {
		cloudProviderConfig := []map[string][]map[string]any{
			{
				"aws": {
					{
						"test_s3_bucket": s3Bucket,
					},
				},
			},
		}
		if err = d.Set("cloud_provider_config", cloudProviderConfig); err != nil {
			return nil, fmt.Errorf(errorFederatedDatabaseInstanceSetting, "cloud_provider_config", name, err)
		}
	}

	dataFederationInstance, _, err := connV2.DataFederationApi.GetFederatedDatabase(ctx, projectID, name).Execute()
	if err != nil {
		return nil, fmt.Errorf("couldn't import data federated instance (%s) for project (%s), error: %s", name, projectID, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		return nil, fmt.Errorf("error setting `project_id` for data federated instance (%s): %s", d.Id(), err)
	}

	if err := d.Set("name", dataFederationInstance.GetName()); err != nil {
		return nil, fmt.Errorf("error setting `name` for data federated instance (%s): %s", d.Id(), err)
	}

	if val, ok := dataFederationInstance.GetCloudProviderConfigOk(); ok {
		if cloudProviderField := flattenCloudProviderConfig(val); cloudProviderField != nil {
			if err = d.Set("cloud_provider_config", cloudProviderField); err != nil {
				return nil, fmt.Errorf(errorFederatedDatabaseInstanceSetting, "cloud_provider_config", name, err)
			}
		}
	}

	if storage, ok := dataFederationInstance.GetStorageOk(); ok {
		if databases, ok := storage.GetDatabasesOk(); ok {
			if storageDatabaseField := flattenDataFederationDatabase(*databases); storageDatabaseField != nil {
				if err := d.Set("storage_databases", storageDatabaseField); err != nil {
					return nil, fmt.Errorf(errorFederatedDatabaseInstanceSetting, "storage_databases", name, err)
				}
			}
		}

		if stores, ok := storage.GetStoresOk(); ok {
			if err := d.Set("storage_stores", flattenDataFederationStores(*stores)); err != nil {
				return nil, fmt.Errorf(errorFederatedDatabaseInstanceSetting, "storage_stores", name, err)
			}
		}
	}

	if err := d.Set("state", dataFederationInstance.GetState()); err != nil {
		return nil, fmt.Errorf(errorFederatedDatabaseInstanceSetting, "state", name, err)
	}

	if err := d.Set("hostnames", dataFederationInstance.GetHostnames()); err != nil {
		return nil, fmt.Errorf(errorFederatedDatabaseInstanceSetting, "hostnames", name, err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id": projectID,
		"name":       *dataFederationInstance.Name,
	}))

	return []*schema.ResourceData{d}, nil
}

func newDataFederationStorage(d *schema.ResourceData) *admin.DataLakeStorage {
	return &admin.DataLakeStorage{
		Databases: newDataFederationDatabase(d),
		Stores:    newStores(d),
	}
}

func newStores(d *schema.ResourceData) *[]admin.DataLakeStoreSettings {
	storesFromConf := d.Get("storage_stores").(*schema.Set).List()
	if len(storesFromConf) == 0 {
		return new([]admin.DataLakeStoreSettings)
	}
	stores := make([]admin.DataLakeStoreSettings, len(storesFromConf))
	for i, storeFromConf := range storesFromConf {
		storeFromConfMap := storeFromConf.(map[string]any)
		stores[i] = admin.DataLakeStoreSettings{
			Name:                     conversion.StringPtr(storeFromConfMap["name"].(string)),
			Provider:                 storeFromConfMap["provider"].(string),
			Region:                   conversion.StringPtr(storeFromConfMap["region"].(string)),
			ProjectId:                conversion.StringPtr(storeFromConfMap["project_id"].(string)),
			Bucket:                   conversion.StringPtr(storeFromConfMap["bucket"].(string)),
			ClusterName:              conversion.StringPtr(storeFromConfMap["cluster_name"].(string)),
			Prefix:                   conversion.StringPtr(storeFromConfMap["prefix"].(string)),
			Delimiter:                conversion.StringPtr(storeFromConfMap["delimiter"].(string)),
			IncludeTags:              conversion.Pointer(storeFromConfMap["include_tags"].(bool)),
			AdditionalStorageClasses: newAdditionalStorageClasses(storeFromConfMap["additional_storage_classes"].([]any)),
			ReadPreference:           newReadPreference(storeFromConfMap),
		}
	}
	return &stores
}

func newAdditionalStorageClasses(additionalStorageClassesFromConfig []any) *[]string {
	if len(additionalStorageClassesFromConfig) == 0 {
		return new([]string)
	}
	additionalStorageClasses := make([]string, len(additionalStorageClassesFromConfig))
	for i, additionalStorageClassFromConfig := range additionalStorageClassesFromConfig {
		additionalStorageClasses[i] = additionalStorageClassFromConfig.(string)
	}
	return &additionalStorageClasses
}

func newReadPreference(storeFromConfMap map[string]any) *admin.DataLakeAtlasStoreReadPreference {
	readPreferenceFromConf, ok := storeFromConfMap["read_preference"].([]any)
	if !ok || len(readPreferenceFromConf) == 0 {
		return nil
	}
	readPreferenceFromConfMap := readPreferenceFromConf[0].(map[string]any)
	return &admin.DataLakeAtlasStoreReadPreference{
		Mode:                conversion.StringPtr(readPreferenceFromConfMap["mode"].(string)),
		MaxStalenessSeconds: conversion.IntPtr(readPreferenceFromConfMap["max_staleness_seconds"].(int)),
		TagSets:             newTagSets(readPreferenceFromConfMap),
	}
}

func newTagSets(readPreferenceFromConfMap map[string]any) *[][]admin.DataLakeAtlasStoreReadPreferenceTag {
	tagSetsFromConf, ok := readPreferenceFromConfMap["tag_sets"].([]any)
	if !ok || len(tagSetsFromConf) == 0 {
		return new([][]admin.DataLakeAtlasStoreReadPreferenceTag)
	}
	var res [][]admin.DataLakeAtlasStoreReadPreferenceTag
	for ts := 0; ts < len(tagSetsFromConf); ts++ {
		tagSetFromConfMap := tagSetsFromConf[ts].(map[string]any)
		tagsFromConfigMap := tagSetFromConfMap["tags"].([]any)
		var atlastags []admin.DataLakeAtlasStoreReadPreferenceTag
		for t := 0; t < len(tagsFromConfigMap); t++ {
			tagFromConfMap := tagsFromConfigMap[t].(map[string]any)
			atlastags = append(atlastags, admin.DataLakeAtlasStoreReadPreferenceTag{
				Name:  conversion.StringPtr(tagFromConfMap["name"].(string)),
				Value: conversion.StringPtr(tagFromConfMap["value"].(string)),
			})
		}
		res = append(res, atlastags)
	}
	return &res
}

func newDataFederationDatabase(d *schema.ResourceData) *[]admin.DataLakeDatabaseInstance {
	storageDBsFromConf := d.Get("storage_databases").(*schema.Set).List()
	if len(storageDBsFromConf) == 0 {
		return new([]admin.DataLakeDatabaseInstance)
	}
	dbs := make([]admin.DataLakeDatabaseInstance, len(storageDBsFromConf))
	for i, storageDBFromConf := range storageDBsFromConf {
		storageDBFromConfMap := storageDBFromConf.(map[string]any)
		dbs[i] = admin.DataLakeDatabaseInstance{
			Name:                   conversion.StringPtr(storageDBFromConfMap["name"].(string)),
			MaxWildcardCollections: conversion.IntPtr(storageDBFromConfMap["max_wildcard_collections"].(int)),
			Collections:            newDataFederationCollections(storageDBFromConfMap),
		}
	}
	return &dbs
}

func newDataFederationCollections(storageDBFromConfMap map[string]any) *[]admin.DataLakeDatabaseCollection {
	collectionsFromConf := storageDBFromConfMap["collections"].(*schema.Set).List()
	if len(collectionsFromConf) == 0 {
		return new([]admin.DataLakeDatabaseCollection)
	}
	collections := make([]admin.DataLakeDatabaseCollection, len(collectionsFromConf))
	for i, collectionFromConf := range collectionsFromConf {
		collections[i] = admin.DataLakeDatabaseCollection{
			Name:        conversion.StringPtr(collectionFromConf.(map[string]any)["name"].(string)),
			DataSources: newDataFederationDataSource(collectionFromConf.(map[string]any)),
		}
	}
	return &collections
}

func newDataFederationDataSource(collectionFromConf map[string]any) *[]admin.DataLakeDatabaseDataSourceSettings {
	dataSourcesFromConf := collectionFromConf["data_sources"].(*schema.Set).List()
	if len(dataSourcesFromConf) == 0 {
		return new([]admin.DataLakeDatabaseDataSourceSettings)
	}
	dataSources := make([]admin.DataLakeDatabaseDataSourceSettings, len(dataSourcesFromConf))
	for i, dataSourceFromConf := range dataSourcesFromConf {
		dataSourceFromConfMap := dataSourceFromConf.(map[string]any)
		dataSources[i] = admin.DataLakeDatabaseDataSourceSettings{
			AllowInsecure:       conversion.Pointer(dataSourceFromConfMap["allow_insecure"].(bool)),
			Database:            conversion.StringPtr(dataSourceFromConfMap["database"].(string)),
			Collection:          conversion.StringPtr(dataSourceFromConfMap["collection"].(string)),
			CollectionRegex:     conversion.StringPtr(dataSourceFromConfMap["collection_regex"].(string)),
			DatabaseRegex:       conversion.StringPtr(dataSourceFromConfMap["database_regex"].(string)),
			DefaultFormat:       conversion.StringPtr(dataSourceFromConfMap["default_format"].(string)),
			Path:                conversion.StringPtr(dataSourceFromConfMap["path"].(string)),
			ProvenanceFieldName: conversion.StringPtr(dataSourceFromConfMap["provenance_field_name"].(string)),
			StoreName:           conversion.StringPtr(dataSourceFromConfMap["store_name"].(string)),
			DatasetName:         conversion.StringPtr(dataSourceFromConfMap["dataset_name"].(string)),
			Urls:                newUrls(dataSourceFromConfMap["urls"].([]any)),
		}
	}
	return &dataSources
}

func newUrls(urlsFromConfig []any) *[]string {
	if len(urlsFromConfig) == 0 {
		return new([]string)
	}
	urls := make([]string, len(urlsFromConfig))
	for i, urlFromConfig := range urlsFromConfig {
		urls[i] = urlFromConfig.(string)
	}
	return &urls
}

func newCloudProviderConfig(d *schema.ResourceData) *admin.DataLakeCloudProviderConfig {
	if cloudProvider, ok := d.Get("cloud_provider_config").([]any); ok && len(cloudProvider) == 1 {
		return &admin.DataLakeCloudProviderConfig{
			Aws:   newAWSConfig(cloudProvider),
			Azure: newAzureConfig(cloudProvider),
		}
	}

	return nil
}

func newAWSConfig(cloudProvider []any) *admin.DataLakeAWSCloudProviderConfig {
	if aws, ok := cloudProvider[0].(map[string]any)["aws"].([]any); ok && len(aws) == 1 {
		awsSchema := aws[0].(map[string]any)
		return admin.NewDataLakeAWSCloudProviderConfig(awsSchema["role_id"].(string), awsSchema["test_s3_bucket"].(string))
	}

	return nil
}

func newAzureConfig(cloudProvider []any) *admin.DataFederationAzureCloudProviderConfig {
	if azure, ok := cloudProvider[0].(map[string]any)["azure"].([]any); ok && len(azure) == 1 {
		azureSchema := azure[0].(map[string]any)
		return admin.NewDataFederationAzureCloudProviderConfig(azureSchema["role_id"].(string))
	}

	return nil
}

func newDataProcessRegion(d *schema.ResourceData) *admin.DataLakeDataProcessRegion {
	if dataProcessRegion, ok := d.Get("data_process_region").([]any); ok && len(dataProcessRegion) == 1 {
		return &admin.DataLakeDataProcessRegion{
			CloudProvider: dataProcessRegion[0].(map[string]any)["cloud_provider"].(string),
			Region:        dataProcessRegion[0].(map[string]any)["region"].(string),
		}
	}

	return nil
}

func flattenCloudProviderConfig(cloudProviderConfig *admin.DataLakeCloudProviderConfig) []map[string]any {
	if cloudProviderConfig == nil {
		return nil
	}

	return []map[string]any{
		{
			"aws":   flattenAWSCloudProviderConfig(cloudProviderConfig.Aws),
			"azure": flattenAzureCloudProviderConfig(cloudProviderConfig.Azure),
		},
	}
}

func flattenAWSCloudProviderConfig(aws *admin.DataLakeAWSCloudProviderConfig) []map[string]any {
	if aws == nil {
		return nil
	}

	// test_s3_bucket is not part of the API response

	return []map[string]any{
		{
			"role_id":              aws.GetRoleId(),
			"test_s3_bucket":       aws.GetTestS3Bucket(),
			"iam_assumed_role_arn": aws.GetIamAssumedRoleARN(),
			"iam_user_arn":         aws.GetIamUserARN(),
			"external_id":          aws.GetExternalId(),
		},
	}
}

func flattenAzureCloudProviderConfig(azure *admin.DataFederationAzureCloudProviderConfig) []map[string]any {
	if azure == nil {
		return nil
	}

	return []map[string]any{
		{
			"role_id":              azure.GetRoleId(),
			"atlas_app_id":         azure.GetAtlasAppId(),
			"service_principal_id": azure.GetServicePrincipalId(),
			"tenant_id":            azure.GetTenantId(),
		},
	}
}

func flattenDataProcessRegion(processRegion *admin.DataLakeDataProcessRegion) []map[string]any {
	if processRegion == nil || (processRegion.Region == "" && processRegion.CloudProvider == "") {
		return nil
	}

	return []map[string]any{
		{
			"cloud_provider": processRegion.GetCloudProvider(),
			"region":         processRegion.GetRegion(),
		},
	}
}

func flattenDataFederationDatabase(atlasDatabases []admin.DataLakeDatabaseInstance) []map[string]any {
	dbs := make([]map[string]any, len(atlasDatabases))

	for i, atlasDatabase := range atlasDatabases {
		dbs[i] = map[string]any{
			"name":                     atlasDatabase.GetName(),
			"max_wildcard_collections": atlasDatabase.GetMaxWildcardCollections(),
			"collections":              flattenDataFederationCollections(atlasDatabase.GetCollections()),
			"views":                    flattenDataFederationDatabaseViews(atlasDatabase.GetViews()),
		}
	}

	return dbs
}

func flattenDataFederationDatabaseViews(atlasViews []admin.DataLakeApiBase) []map[string]any {
	views := make([]map[string]any, len(atlasViews))

	for i, atlasView := range atlasViews {
		views[i] = map[string]any{
			"name":     atlasView.GetName(),
			"source":   atlasView.GetSource(),
			"pipeline": atlasView.GetPipeline(),
		}
	}

	return views
}

func flattenDataFederationCollections(atlasCollections []admin.DataLakeDatabaseCollection) []map[string]any {
	colls := make([]map[string]any, len(atlasCollections))

	for i, atlasCollection := range atlasCollections {
		colls[i] = map[string]any{
			"name":         atlasCollection.GetName(),
			"data_sources": flattenDataFederationDataSources(atlasCollection.GetDataSources()),
		}
	}

	return colls
}

func flattenDataFederationDataSources(atlasDataSources []admin.DataLakeDatabaseDataSourceSettings) []map[string]any {
	out := make([]map[string]any, len(atlasDataSources))

	for i, AtlasDataSource := range atlasDataSources {
		out[i] = map[string]any{
			"allow_insecure":        AtlasDataSource.GetAllowInsecure(),
			"collection":            AtlasDataSource.GetCollection(),
			"collection_regex":      AtlasDataSource.GetCollectionRegex(),
			"database":              AtlasDataSource.GetDatabase(),
			"database_regex":        AtlasDataSource.GetDatabaseRegex(),
			"default_format":        AtlasDataSource.GetDefaultFormat(),
			"path":                  AtlasDataSource.GetPath(),
			"provenance_field_name": AtlasDataSource.GetProvenanceFieldName(),
			"store_name":            AtlasDataSource.GetStoreName(),
			"dataset_name":          AtlasDataSource.GetDatasetName(),
			"urls":                  AtlasDataSource.GetUrls(),
		}
	}

	return out
}

func flattenDataFederationStores(stores []admin.DataLakeStoreSettings) []map[string]any {
	store := make([]map[string]any, 0)

	for i := range stores {
		store = append(store, map[string]any{
			"name":                       stores[i].GetName(),
			"provider":                   stores[i].GetProvider(),
			"region":                     stores[i].GetRegion(),
			"project_id":                 stores[i].GetProjectId(),
			"bucket":                     stores[i].GetBucket(),
			"cluster_name":               stores[i].GetClusterName(),
			"prefix":                     stores[i].GetPrefix(),
			"delimiter":                  stores[i].GetDelimiter(),
			"include_tags":               stores[i].GetIncludeTags(),
			"additional_storage_classes": stores[i].GetAdditionalStorageClasses(),
			"allow_insecure":             stores[i].GetAllowInsecure(),
			"public":                     stores[i].Public,
			"default_format":             stores[i].GetDefaultFormat(),
			"urls":                       stores[i].GetUrls(),
			"read_preference":            newReadPreferenceField(stores[i].ReadPreference),
		})
	}

	return store
}

func newReadPreferenceField(atlasReadPreference *admin.DataLakeAtlasStoreReadPreference) []map[string]any {
	if atlasReadPreference == nil {
		return nil
	}

	return []map[string]any{
		{
			"mode":                  atlasReadPreference.GetMode(),
			"max_staleness_seconds": atlasReadPreference.GetMaxStalenessSeconds(),
			"tag_sets":              flattenReadPreferenceTagSets(atlasReadPreference.GetTagSets()),
		},
	}
}

func flattenReadPreferenceTagSets(tagSets [][]admin.DataLakeAtlasStoreReadPreferenceTag) []map[string]any {
	tfTagSets := make([]map[string]any, 0)

	for i := range tagSets {
		tfTagSets = append(tfTagSets, map[string]any{
			"tags": flattenReadPreferenceTags(tagSets[i]),
		})
	}

	return tfTagSets
}

func flattenReadPreferenceTags(tags []admin.DataLakeAtlasStoreReadPreferenceTag) []map[string]any {
	tfTags := make([]map[string]any, 0)

	for i := range tags {
		tfTags = append(tfTags, map[string]any{
			"name":  tags[i].Name,
			"value": tags[i].Value,
		})
	}

	return tfTags
}

func splitDataFederatedInstanceImportID(id string) (projectID, name, s3Bucket string, err error) {
	var parts = strings.Split(id, "--")

	if len(parts) == 2 {
		projectID = parts[0]
		name = parts[1]
		return
	}

	if len(parts) == 3 {
		projectID = parts[0]
		name = parts[1]
		s3Bucket = parts[2]
		return
	}

	err = errors.New("import format error: to import a Data Federated instance, use the format {project_id}--{name} or {project_id}--{name}--{test_s3_bucket}")
	return
}
