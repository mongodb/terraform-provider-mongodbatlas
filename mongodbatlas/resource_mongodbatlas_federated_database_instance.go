package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"go.mongodb.org/atlas-sdk/v20230201005/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	errorFederatedDatabaseInstanceCreate  = "error creating MongoDB Atlas Federated Database Instace: %s"
	errorFederatedDatabaseInstanceRead    = "error reading MongoDB Atlas Federated Database Instace (%s): %s"
	errorFederatedDatabaseInstanceDelete  = "error deleting MongoDB Atlas Federated Database Instace (%s): %s"
	errorFederatedDatabaseInstanceUpdate  = "error updating MongoDB Atlas Federated Database Instace (%s): %s"
	errorFederatedDatabaseInstanceSetting = "error setting `%s` for MongoDB Atlas Federated Database Instace (%s): %s"
)

func resourceMongoDBAtlasFederatedDatabaseInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBFederatedDatabaseInstanceCreate,
		ReadContext:   resourceMongoDBAFederatedDatabaseInstanceRead,
		UpdateContext: resourceMongoDBFederatedDatabaseInstanceUpdate,
		DeleteContext: resourceMongoDBAtlasFederatedDatabaseInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasFederatedDatabaseInstanceImportState,
		},
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
			"cloud_provider_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"aws": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Required: true,
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
					},
				},
			},
			"data_process_region": {
				Type:     schema.TypeList,
				MaxItems: 1,
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
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"collections": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"name": {
								Type:     schema.TypeString,
								Optional: true,
							},
							"data_sources": {
								Type:     schema.TypeSet,
								Optional: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"store_name": {
											Type:     schema.TypeString,
											Optional: true,
										},
										"default_format": {
											Type:     schema.TypeString,
											Optional: true,
										},
										"path": {
											Type:     schema.TypeString,
											Optional: true,
										},
										"allow_insecure": {
											Type:     schema.TypeBool,
											Optional: true,
										},
										"database": {
											Type:     schema.TypeString,
											Optional: true,
										},
										"database_regex": {
											Type:     schema.TypeString,
											Optional: true,
										},
										"collection": {
											Type:     schema.TypeString,
											Optional: true,
										},
										"collection_regex": {
											Type:     schema.TypeString,
											Optional: true,
										},
										"provenance_field_name": {
											Type:     schema.TypeString,
											Optional: true,
										},
										"urls": {
											Type:     schema.TypeList,
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
					Optional: true,
				},
				"provider": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"region": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"bucket": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"cluster_name": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"project_id": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"prefix": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"delimiter": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"include_tags": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"allow_insecure": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"additional_storage_classes": {
					Type:     schema.TypeList,
					Optional: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"public": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"default_format": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"urls": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"read_preference": {
					Type:     schema.TypeList,
					MaxItems: 1,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"mode": {
								Type:     schema.TypeString,
								Optional: true,
							},
							"max_staleness_seconds": {
								Type:     schema.TypeInt,
								Optional: true,
							},
						},
					},
				},
			},
		},
	}
}

func resourceMongoDBFederatedDatabaseInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	connV2 := meta.(*MongoDBClient).AtlasV2

	projectID := d.Get("project_id").(string)
	name := d.Get("name").(string)

	if _, _, err := connV2.DataFederationApi.CreateFederatedDatabase(ctx, projectID, &admin.DataLakeTenant{
		Name:                stringPtr(name),
		CloudProviderConfig: newCloudProviderConfig(d),
		DataProcessRegion:   newDataProcessRegion(d),
		Storage:             newDataFederationStorage(d),
	}).Execute(); err != nil {
		return diag.FromErr(fmt.Errorf(errorFederatedDatabaseInstanceCreate, err))
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"name":       name,
	}))

	return resourceMongoDBAFederatedDatabaseInstanceRead(ctx, d, meta)
}

func resourceMongoDBAFederatedDatabaseInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	connV2 := meta.(*MongoDBClient).AtlasV2
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	name := ids["name"]

	dataFederationInstance, resp, err := connV2.DataFederationApi.GetFederatedDatabase(ctx, projectID, name).Execute()
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorFederatedDatabaseInstanceRead, name, err))
	}

	if val, ok := dataFederationInstance.GetCloudProviderConfigOk(); ok {
		if cloudProviderField := flattenCloudProviderConfig(d, val); cloudProviderField != nil {
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

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"name":       name,
	}))

	return nil
}

func resourceMongoDBFederatedDatabaseInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	connV2 := meta.(*MongoDBClient).AtlasV2

	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	name := ids["name"]

	if _, _, err := connV2.DataFederationApi.UpdateFederatedDatabase(ctx, projectID, name, &admin.DataLakeTenant{
		Name:                stringPtr(name),
		CloudProviderConfig: newCloudProviderConfig(d),
		DataProcessRegion:   newDataProcessRegion(d),
		Storage:             newDataFederationStorage(d),
	}).Execute(); err != nil {
		return diag.FromErr(fmt.Errorf(errorFederatedDatabaseInstanceUpdate, name, err))
	}

	return resourceMongoDBAFederatedDatabaseInstanceRead(ctx, d, meta)
}

func resourceMongoDBAtlasFederatedDatabaseInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	connV2 := meta.(*MongoDBClient).AtlasV2

	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	name := ids["name"]

	if _, _, err := connV2.DataFederationApi.DeleteFederatedDatabase(ctx, projectID, name).Execute(); err != nil {
		return diag.FromErr(fmt.Errorf(errorFederatedDatabaseInstanceDelete, name, err))
	}

	return nil
}

func resourceMongoDBAtlasFederatedDatabaseInstanceImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	connV2 := meta.(*MongoDBClient).AtlasV2

	projectID, name, s3Bucket, err := splitDataFederatedInstanceImportID(d.Id())
	if err != nil {
		return nil, err
	}

	// test_s3_bucket is not part of the API response
	if s3Bucket != "" {
		cloudProviderConfig := []map[string][]map[string]interface{}{
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
		if cloudProviderField := flattenCloudProviderConfig(d, val); cloudProviderField != nil {
			if err = d.Set("cloud_provider_config", cloudProviderField); err != nil {
				return nil, fmt.Errorf(errorFederatedDatabaseInstanceSetting, "cloud_provider_config", name, err)
			}
		}
	}

	if storage, ok := dataFederationInstance.GetStorageOk(); ok {
		if databases, ok := storage.GetDatabasesOk(); ok {
			if storageDatabaseField := flattenDataFederationDatabase(databases); storageDatabaseField != nil {
				if err := d.Set("storage_databases", storageDatabaseField); err != nil {
					return nil, fmt.Errorf(errorFederatedDatabaseInstanceSetting, "storage_databases", name, err)
				}
			}
		}

		if stores, ok := storage.GetStoresOk(); ok {
			if err := d.Set("storage_stores", flattenDataFederationStores(stores)); err != nil {
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

	d.SetId(encodeStateID(map[string]string{
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

func newStores(d *schema.ResourceData) []admin.DataLakeStoreSettings {
	storesFromConf := d.Get("storage_stores").(*schema.Set).List()
	if len(storesFromConf) == 0 {
		return nil
	}

	stores := make([]admin.DataLakeStoreSettings, len(storesFromConf))
	for i, storeFromConf := range storesFromConf {
		storeFromConfMap := storeFromConf.(map[string]interface{})
		stores[i] = admin.DataLakeStoreSettings{
			Name:                     stringPtr(storeFromConfMap["name"].(string)),
			Provider:                 storeFromConfMap["provider"].(string),
			Region:                   stringPtr(storeFromConfMap["region"].(string)),
			ProjectId:                pointer(storeFromConfMap["project_id"].(string)),
			Bucket:                   stringPtr(storeFromConfMap["bucket"].(string)),
			ClusterName:              stringPtr(storeFromConfMap["cluster_name"].(string)),
			Prefix:                   stringPtr(storeFromConfMap["prefix"].(string)),
			Delimiter:                stringPtr(storeFromConfMap["delimiter"].(string)),
			IncludeTags:              pointer(storeFromConfMap["include_tags"].(bool)),
			AdditionalStorageClasses: newAdditionalStorageClasses(storeFromConfMap["additional_storage_classes"].([]interface{})),
			ReadPreference:           newReadPreference(storeFromConfMap),
		}
	}

	return stores
}

func newAdditionalStorageClasses(additionalStorageClassesFromConfig []interface{}) []string {
	if len(additionalStorageClassesFromConfig) == 0 {
		return nil
	}

	additionalStorageClasses := make([]string, len(additionalStorageClassesFromConfig))
	for i, additionalStorageClassFromConfig := range additionalStorageClassesFromConfig {
		additionalStorageClasses[i] = additionalStorageClassFromConfig.(string)
	}

	return additionalStorageClasses
}

func newReadPreference(storeFromConfMap map[string]interface{}) *admin.DataLakeAtlasStoreReadPreference {
	readPreferenceFromConf, ok := storeFromConfMap["read_preference"].([]interface{})
	if !ok || len(readPreferenceFromConf) == 0 {
		return nil
	}
	readPreferenceFromConfMap := readPreferenceFromConf[0].(map[string]interface{})
	return &admin.DataLakeAtlasStoreReadPreference{
		Mode:                stringPtr(readPreferenceFromConfMap["mode"].(string)),
		MaxStalenessSeconds: intPtr(readPreferenceFromConfMap["max_staleness_seconds"].(int)),
	}

}

func newDataFederationDatabase(d *schema.ResourceData) []admin.DataLakeDatabaseInstance {
	storageDBsFromConf := d.Get("storage_databases").(*schema.Set).List()
	if len(storageDBsFromConf) == 0 {
		return nil
	}

	dbs := make([]admin.DataLakeDatabaseInstance, len(storageDBsFromConf))
	for i, storageDBFromConf := range storageDBsFromConf {
		storageDBFromConfMap := storageDBFromConf.(map[string]interface{})

		dbs[i] = admin.DataLakeDatabaseInstance{
			Name:                   stringPtr(storageDBFromConfMap["name"].(string)),
			MaxWildcardCollections: intPtr(storageDBFromConfMap["max_wildcard_collections"].(int)),
			Collections:            newDataFederationCollections(storageDBFromConfMap),
		}
	}

	return dbs
}

func newDataFederationCollections(storageDBFromConfMap map[string]interface{}) []admin.DataLakeDatabaseCollection {
	collectionsFromConf := storageDBFromConfMap["collections"].(*schema.Set).List()
	if len(collectionsFromConf) == 0 {
		return nil
	}

	collections := make([]admin.DataLakeDatabaseCollection, len(collectionsFromConf))
	for i, collectionFromConf := range collectionsFromConf {
		collections[i] = admin.DataLakeDatabaseCollection{
			Name:        stringPtr(collectionFromConf.(map[string]interface{})["name"].(string)),
			DataSources: newDataFederationDataSource(collectionFromConf.(map[string]interface{})),
		}
	}

	return collections
}

func newDataFederationDataSource(collectionFromConf map[string]interface{}) []admin.DataLakeDatabaseDataSourceSettings {
	dataSourcesFromConf := collectionFromConf["data_sources"].(*schema.Set).List()
	if len(dataSourcesFromConf) == 0 {
		return nil
	}
	dataSources := make([]admin.DataLakeDatabaseDataSourceSettings, len(dataSourcesFromConf))
	for i, dataSourceFromConf := range dataSourcesFromConf {
		dataSourceFromConfMap := dataSourceFromConf.(map[string]interface{})

		dataSources[i] = admin.DataLakeDatabaseDataSourceSettings{
			AllowInsecure:       pointer(dataSourceFromConfMap["allow_insecure"].(bool)),
			Database:            stringPtr(dataSourceFromConfMap["database"].(string)),
			Collection:          stringPtr(dataSourceFromConfMap["collection"].(string)),
			CollectionRegex:     stringPtr(dataSourceFromConfMap["collection_regex"].(string)),
			DefaultFormat:       stringPtr(dataSourceFromConfMap["default_format"].(string)),
			Path:                stringPtr(dataSourceFromConfMap["path"].(string)),
			ProvenanceFieldName: stringPtr(dataSourceFromConfMap["provenance_field_name"].(string)),
			StoreName:           stringPtr(dataSourceFromConfMap["store_name"].(string)),
			Urls:                newUrls(dataSourceFromConfMap["urls"].([]interface{})),
		}
	}

	return dataSources
}

func newUrls(urlsFromConfig []interface{}) []string {
	if len(urlsFromConfig) == 0 {
		return nil
	}

	urls := make([]string, len(urlsFromConfig))
	for i, urlFromConfig := range urlsFromConfig {
		urls[i] = urlFromConfig.(string)
	}

	return urls
}

func newCloudProviderConfig(d *schema.ResourceData) *admin.DataLakeCloudProviderConfig {
	if cloudProvider, ok := d.Get("cloud_provider_config").([]interface{}); ok && len(cloudProvider) == 1 {
		return admin.NewDataLakeCloudProviderConfig(*newAWSConfig(cloudProvider))
	}

	return nil
}

func newAWSConfig(cloudProvider []interface{}) *admin.DataLakeAWSCloudProviderConfig {
	if aws, ok := cloudProvider[0].(map[string]interface{})["aws"].([]interface{}); ok && len(aws) == 1 {
		awsSchema := aws[0].(map[string]interface{})
		return admin.NewDataLakeAWSCloudProviderConfig(awsSchema["role_id"].(string), awsSchema["test_s3_bucket"].(string))
	}

	return nil
}

func newDataProcessRegion(d *schema.ResourceData) *admin.DataLakeDataProcessRegion {
	if dataProcessRegion, ok := d.Get("data_process_region").([]interface{}); ok && len(dataProcessRegion) == 1 {
		return &admin.DataLakeDataProcessRegion{
			CloudProvider: dataProcessRegion[0].(map[string]interface{})["cloud_provider"].(string),
			Region:        dataProcessRegion[0].(map[string]interface{})["region"].(string),
		}
	}

	return nil
}

func flattenCloudProviderConfig(d *schema.ResourceData, cloudProviderConfig *admin.DataLakeCloudProviderConfig) []map[string]interface{} {
	if cloudProviderConfig == nil {
		return nil
	}

	aws := cloudProviderConfig.GetAws()

	awsOut := []map[string]interface{}{
		{
			"role_id":              aws.GetRoleId(),
			"iam_assumed_role_arn": aws.GetIamAssumedRoleARN(),
			"iam_user_arn":         aws.GetIamUserARN(),
			"external_id":          aws.GetExternalId(),
		},
	}

	currentCloudProviderConfig, ok := d.Get("cloud_provider_config").([]interface{})
	if !ok || len(currentCloudProviderConfig) == 0 {
		return []map[string]interface{}{
			{
				"aws": &awsOut,
			},
		}
	}
	// test_s3_bucket is not part of the API response
	if currentAWS, ok := currentCloudProviderConfig[0].(map[string]interface{})["aws"].([]interface{}); ok {
		if testS3Bucket, ok := currentAWS[0].(map[string]interface{})["test_s3_bucket"].(string); ok {
			awsOut[0]["test_s3_bucket"] = testS3Bucket
			return []map[string]interface{}{
				{
					"aws": &awsOut,
				},
			}
		}
	}

	return awsOut
}

func flattenDataProcessRegion(processRegion *admin.DataLakeDataProcessRegion) []map[string]interface{} {
	if processRegion == nil || (processRegion.Region != "" && processRegion.CloudProvider != "") {
		return nil
	}

	return []map[string]interface{}{
		{
			"cloud_provider": processRegion.GetCloudProvider(),
			"region":         processRegion.GetRegion(),
		},
	}
}

func flattenDataFederationDatabase(atlasDatabases []admin.DataLakeDatabaseInstance) []map[string]interface{} {
	dbs := make([]map[string]interface{}, len(atlasDatabases))

	for i, atlasDatabase := range atlasDatabases {
		dbs[i] = map[string]interface{}{
			"name":                     atlasDatabase.GetName(),
			"max_wildcard_collections": atlasDatabase.GetMaxWildcardCollections(),
			"collections":              flattenDataFederationCollections(atlasDatabase.Collections),
			"views":                    flattenDataFederationDatabaseViews(atlasDatabase.Views),
		}
	}

	return dbs
}

func flattenDataFederationDatabaseViews(atlasViews []admin.DataLakeApiBase) []map[string]interface{} {
	views := make([]map[string]interface{}, len(atlasViews))

	for i, atlasView := range atlasViews {
		views[i] = map[string]interface{}{
			"name":     atlasView.GetName(),
			"source":   atlasView.GetSource(),
			"pipeline": atlasView.GetPipeline(),
		}
	}

	return views
}

func flattenDataFederationCollections(atlasCollections []admin.DataLakeDatabaseCollection) []map[string]interface{} {
	colls := make([]map[string]interface{}, len(atlasCollections))

	for i, atlasCollection := range atlasCollections {
		colls[i] = map[string]interface{}{
			"name":         atlasCollection.GetName(),
			"data_sources": flattenDataFederationDataSources(atlasCollection.DataSources),
		}
	}

	return colls
}

func flattenDataFederationDataSources(atlasDataSources []admin.DataLakeDatabaseDataSourceSettings) []map[string]interface{} {
	out := make([]map[string]interface{}, len(atlasDataSources))

	for i, AtlasDataSource := range atlasDataSources {
		out[i] = map[string]interface{}{
			"allow_insecure":        AtlasDataSource.GetAllowInsecure(),
			"collection":            AtlasDataSource.GetCollection(),
			"collection_regex":      AtlasDataSource.GetCollectionRegex(),
			"database":              AtlasDataSource.GetDatabase(),
			"database_regex":        AtlasDataSource.GetDatabaseRegex(),
			"default_format":        AtlasDataSource.GetDefaultFormat(),
			"path":                  AtlasDataSource.GetPath(),
			"provenance_field_name": AtlasDataSource.GetProvenanceFieldName(),
			"store_name":            AtlasDataSource.GetStoreName(),
			"urls":                  AtlasDataSource.GetUrls(),
		}
	}

	return out
}

func flattenDataFederationStores(stores []admin.DataLakeStoreSettings) []map[string]interface{} {
	store := make([]map[string]interface{}, 0)

	for i := range stores {
		store = append(store, map[string]interface{}{
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

func newReadPreferenceField(atlasReadPreference *admin.DataLakeAtlasStoreReadPreference) []map[string]interface{} {
	if atlasReadPreference == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"mode":                  atlasReadPreference.GetMode(),
			"max_staleness_seconds": atlasReadPreference.GetMaxStalenessSeconds(),
		},
	}
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
