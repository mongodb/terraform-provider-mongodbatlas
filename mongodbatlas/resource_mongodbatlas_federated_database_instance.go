package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
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
				"additional_storage_classes": {
					Type:     schema.TypeList,
					Optional: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
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
							"tags": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"name": {
											Type:     schema.TypeString,
											Optional: true,
										},
										"value": {
											Type:     schema.TypeString,
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
	}
}

func resourceMongoDBFederatedDatabaseInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	projectID := d.Get("project_id").(string)
	name := d.Get("name").(string)

	requestBody := &matlas.DataFederationInstance{
		Name:                name,
		CloudProviderConfig: newCloudProviderConfig(d),
		DataProcessRegion:   newDataProcessRegion(d),
		Storage:             newDataFederationStorage(d),
	}

	_, _, err := conn.DataFederation.Create(ctx, projectID, requestBody)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorFederatedDatabaseInstanceCreate, err))
	}
	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"name":       requestBody.Name,
	}))

	return resourceMongoDBAFederatedDatabaseInstanceRead(ctx, d, meta)
}

func resourceMongoDBAFederatedDatabaseInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	name := ids["name"]

	dataFederationInstance, resp, err := conn.DataFederation.Get(ctx, projectID, name)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorFederatedDatabaseInstanceRead, name, err))
	}

	if awsField := flattenCloudProviderConfig(d, dataFederationInstance.CloudProviderConfig); awsField != nil {
		if err = d.Set("aws", awsField); err != nil {
			return diag.FromErr(fmt.Errorf(errorFederatedDatabaseInstanceSetting, "aws", name, err))
		}
	}

	if dataProcessRegionField := flattenDataProcessRegion(dataFederationInstance.DataProcessRegion); dataProcessRegionField != nil {
		if err := d.Set("data_process_region", dataProcessRegionField); err != nil {
			return diag.FromErr(fmt.Errorf(errorFederatedDatabaseInstanceSetting, "data_process_region", name, err))
		}
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"name":       name,
	}))

	return nil
}

func resourceMongoDBFederatedDatabaseInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	name := ids["name"]

	requestBody := &matlas.DataFederationInstance{
		Name:                name,
		CloudProviderConfig: newCloudProviderConfig(d),
		DataProcessRegion:   newDataProcessRegion(d),
		Storage:             newDataFederationStorage(d),
	}
	_, _, err := conn.DataFederation.Update(ctx, projectID, name, requestBody, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorFederatedDatabaseInstanceUpdate, name, err))
	}

	return resourceMongoDBAFederatedDatabaseInstanceRead(ctx, d, meta)
}

func resourceMongoDBAtlasFederatedDatabaseInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	name := ids["name"]

	_, err := conn.DataFederation.Delete(ctx, projectID, name)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorFederatedDatabaseInstanceDelete, name, err))
	}

	return nil
}

func resourceMongoDBAtlasFederatedDatabaseInstanceImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*MongoDBClient).Atlas

	projectID, name, s3Bucket, err := splitDataLakeImportID(d.Id())
	if err != nil {
		return nil, err
	}

	dataFederationInstance, _, err := conn.DataFederation.Get(ctx, projectID, name)
	if err != nil {
		return nil, fmt.Errorf("couldn't import data lake(%s) for project (%s), error: %s", name, projectID, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		return nil, fmt.Errorf("error setting `project_id` for data lakes (%s): %s", d.Id(), err)
	}

	if err := d.Set("name", dataFederationInstance.Name); err != nil {
		return nil, fmt.Errorf("error setting `name` for data lakes (%s): %s", d.Id(), err)
	}
	mapAws := make([]map[string]interface{}, 0)

	mapAws = append(mapAws, map[string]interface{}{
		"test_s3_bucket": s3Bucket,
	})

	if err := d.Set("aws", mapAws); err != nil {
		return nil, fmt.Errorf("error setting `aws` for data lakes (%s): %s", d.Id(), err)
	}

	if storageDatabaseField := flattenDataFederationDatabase(dataFederationInstance.Storage.Databases); storageDatabaseField != nil {
		if err := d.Set("storage_databases", storageDatabaseField); err != nil {
			return nil, fmt.Errorf(errorFederatedDatabaseInstanceSetting, "storage_databases", name, err)
		}
	}

	if err := d.Set("storage_stores", flattenDataFederationStores(dataFederationInstance.Storage.Stores)); err != nil {
		return nil, fmt.Errorf(errorFederatedDatabaseInstanceSetting, "storage_stores", name, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"name":       dataFederationInstance.Name,
	}))

	return []*schema.ResourceData{d}, nil
}

func newDataFederationStorage(d *schema.ResourceData) *matlas.DataFederationStorage {
	return &matlas.DataFederationStorage{
		Databases: newDataFederationDatabase(d),
		Stores:    newStores(d),
	}
}

func newStores(d *schema.ResourceData) []*matlas.DataFederationStore {
	storesFromConf := d.Get("storage_stores").(*schema.Set).List()
	if len(storesFromConf) == 0 {
		return nil
	}

	stores := make([]*matlas.DataFederationStore, len(storesFromConf))
	for i, storeFromConf := range storesFromConf {
		storeFromConfMap := storeFromConf.(map[string]interface{})
		stores[i] = &matlas.DataFederationStore{
			Name:                     storeFromConfMap["name"].(string),
			Provider:                 storeFromConfMap["provider"].(string),
			Region:                   storeFromConfMap["region"].(string),
			ProjectID:                storeFromConfMap["project_id"].(string),
			Bucket:                   storeFromConfMap["bucket"].(string),
			ClusterName:              storeFromConfMap["cluster_name"].(string),
			Prefix:                   storeFromConfMap["prefix"].(string),
			Delimiter:                storeFromConfMap["delimiter"].(string),
			IncludeTags:              pointer(storeFromConfMap["include_tags"].(bool)),
			AdditionalStorageClasses: newAdditionalStorageClasses(storeFromConfMap["additional_storage_classes"].([]interface{})),
			ReadPreference:           newReadPreference(storeFromConfMap),
		}
	}

	return stores
}

func newAdditionalStorageClasses(additionalStorageClassesFromConfig []interface{}) []*string {
	if len(additionalStorageClassesFromConfig) == 0 {
		return nil
	}

	additionalStorageClasses := make([]*string, len(additionalStorageClassesFromConfig))
	for i, additionalStorageClassFromConfig := range additionalStorageClassesFromConfig {
		additionalStorageClasses[i] = pointer(additionalStorageClassFromConfig.(string))
	}

	return additionalStorageClasses
}

func newReadPreference(storeFromConfMap map[string]interface{}) *matlas.ReadPreference {
	readPreferenceFromConf, ok := storeFromConfMap["read_preference"].([]interface{})
	if !ok || len(readPreferenceFromConf) == 0 {
		return nil
	}
	readPreferenceFromConfMap := readPreferenceFromConf[0].(map[string]interface{})
	return &matlas.ReadPreference{
		Mode:                readPreferenceFromConfMap["mode"].(string),
		MaxStalenessSeconds: int32(readPreferenceFromConfMap["max_staleness_seconds"].(int)),
		TagSets:             newTagSets(readPreferenceFromConfMap),
	}
}

func newTagSets(readPreferenceFromConfMap map[string]interface{}) []*matlas.TagSet {
	tagSetsFromConf, ok := readPreferenceFromConfMap["tag_sets"].([]interface{})
	if !ok || len(tagSetsFromConf) == 0 {
		return nil
	}

	tagSets := make([]*matlas.TagSet, len(tagSetsFromConf))
	for i, tagSetFromConf := range tagSetsFromConf {
		storeFromConfMap := tagSetFromConf.(map[string]interface{})
		tagSets[i] = &matlas.TagSet{
			Name:  storeFromConfMap["name"].(string),
			Value: storeFromConfMap["value"].(string),
		}
	}

	return tagSets
}

func newDataFederationDatabase(d *schema.ResourceData) []*matlas.DataFederationDatabase {
	storageDBsFromConf := d.Get("storage_databases").(*schema.Set).List()
	if len(storageDBsFromConf) == 0 {
		return nil
	}

	dbs := make([]*matlas.DataFederationDatabase, len(storageDBsFromConf))
	for i, storageDBFromConf := range storageDBsFromConf {
		storageDBFromConfMap := storageDBFromConf.(map[string]interface{})
		dbs[i] = &matlas.DataFederationDatabase{
			Name:                   storageDBFromConfMap["name"].(string),
			MaxWildcardCollections: int32(storageDBFromConfMap["max_wildcard_collections"].(int)),
			Collections:            newDataFederationCollections(storageDBFromConfMap),
		}
	}

	return dbs
}

func newDataFederationCollections(storageDBFromConfMap map[string]interface{}) []*matlas.DataFederationCollection {
	collectionsFromConf := storageDBFromConfMap["collections"].(*schema.Set).List()
	if len(collectionsFromConf) == 0 {
		return nil
	}

	collections := make([]*matlas.DataFederationCollection, len(collectionsFromConf))
	for i, collectionFromConf := range collectionsFromConf {
		collections[i] = &matlas.DataFederationCollection{
			Name:        collectionFromConf.(map[string]interface{})["name"].(string),
			DataSources: newDataFederationDataSource(collectionFromConf.(map[string]interface{})),
		}
	}

	return collections
}

func newDataFederationDataSource(collectionFromConf map[string]interface{}) []*matlas.DataFederationDataSource {
	dataSourcesFromConf := collectionFromConf["data_sources"].(*schema.Set).List()
	if len(dataSourcesFromConf) == 0 {
		return nil
	}
	dataSources := make([]*matlas.DataFederationDataSource, len(dataSourcesFromConf))
	for i, dataSourceFromConf := range dataSourcesFromConf {
		dataSourceFromConfMap := dataSourceFromConf.(map[string]interface{})
		dataSources[i] = &matlas.DataFederationDataSource{
			AllowInsecure:       pointer(dataSourceFromConfMap["allow_insecure"].(bool)),
			Database:            dataSourceFromConfMap["database"].(string),
			Collection:          dataSourceFromConfMap["collection"].(string),
			CollectionRegex:     dataSourceFromConfMap["collection_regex"].(string),
			DefaultFormat:       dataSourceFromConfMap["default_format"].(string),
			Path:                dataSourceFromConfMap["path"].(string),
			ProvenanceFieldName: dataSourceFromConfMap["provenance_field_name"].(string),
			StoreName:           dataSourceFromConfMap["store_name"].(string),
			Urls:                newUrls(dataSourceFromConfMap["urls"].([]interface{})),
		}
	}

	return dataSources
}

func newUrls(urlsFromConfig []interface{}) []*string {
	if len(urlsFromConfig) == 0 {
		return nil
	}

	urls := make([]*string, len(urlsFromConfig))
	for i, urlFromConfig := range urlsFromConfig {
		urls[i] = pointer(urlFromConfig.(string))
	}

	return urls
}

func newCloudProviderConfig(d *schema.ResourceData) *matlas.CloudProviderConfig {
	if aws, ok := d.Get("aws").([]interface{}); ok && len(aws) == 1 {
		return &matlas.CloudProviderConfig{
			AWSConfig: *newAwsCloudProviderConfig(aws[0].(map[string]interface{})),
		}
	}

	return nil
}

func newAwsCloudProviderConfig(awsSchema map[string]interface{}) *matlas.AwsCloudProviderConfig {
	return &matlas.AwsCloudProviderConfig{
		RoleID:       awsSchema["role_id"].(string),
		TestS3Bucket: awsSchema["test_s3_bucket"].(string),
	}
}

func newDataProcessRegion(d *schema.ResourceData) *matlas.DataProcessRegion {
	if dataProcessRegion, ok := d.Get("data_process_region").([]interface{}); ok && len(dataProcessRegion) == 1 {
		return &matlas.DataProcessRegion{
			CloudProvider: dataProcessRegion[0].(map[string]interface{})["cloud_provider"].(string),
			Region:        dataProcessRegion[0].(map[string]interface{})["region"].(string),
		}
	}

	return nil
}

func flattenCloudProviderConfig(d *schema.ResourceData, aws *matlas.CloudProviderConfig) []map[string]interface{} {
	if aws == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"role_id":              aws.AWSConfig.RoleID,
			"iam_assumed_role_arn": aws.AWSConfig.IAMAssumedRoleARN,
			"iam_user_arn":         aws.AWSConfig.IAMUserARN,
			"external_id":          aws.AWSConfig.ExternalID,
			"test_s3_bucket":       d.Get("aws").([]interface{})[0].(map[string]interface{})["test_s3_bucket"].(string), // test_s3_bucket is not part of the API response
		},
	}
}

func flattenDataProcessRegion(processRegion *matlas.DataProcessRegion) []map[string]interface{} {
	if processRegion == nil || (processRegion.Region != "" && processRegion.CloudProvider != "") {
		return nil
	}

	return []map[string]interface{}{
		{
			"cloud_provider": processRegion.CloudProvider,
			"region":         processRegion.Region,
		},
	}
}

func flattenDataFederationDatabase(atlasDatabases []*matlas.DataFederationDatabase) []map[string]interface{} {
	dbs := make([]map[string]interface{}, len(atlasDatabases))

	for i, atlasDatabase := range atlasDatabases {
		dbs[i] = map[string]interface{}{
			"name":                     atlasDatabase.Name,
			"max_wildcard_collections": atlasDatabase.MaxWildcardCollections,
			"collections":              flattenDataFederationCollections(atlasDatabase.Collections),
			"views":                    flattenDataFederationDatabaseViews(atlasDatabase.Views),
		}
	}

	return dbs
}

func flattenDataFederationDatabaseViews(atlasViews []*matlas.DataFederationDatabaseView) []map[string]interface{} {
	views := make([]map[string]interface{}, len(atlasViews))

	for i, atlasView := range atlasViews {
		views[i] = map[string]interface{}{
			"name":     atlasView.Name,
			"source":   atlasView.Source,
			"pipeline": atlasView.Pipeline,
		}
	}

	return views
}

func flattenDataFederationCollections(atlasCollections []*matlas.DataFederationCollection) []map[string]interface{} {
	colls := make([]map[string]interface{}, len(atlasCollections))

	for i, atlasCollection := range atlasCollections {
		colls[i] = map[string]interface{}{
			"name":         atlasCollection.Name,
			"data_sources": flattenDataFederationDataSources(atlasCollection.DataSources),
		}
	}

	return colls
}

func flattenDataFederationDataSources(atlasDataSources []*matlas.DataFederationDataSource) []map[string]interface{} {
	out := make([]map[string]interface{}, len(atlasDataSources))

	for i, AtlasDataSource := range atlasDataSources {
		out[i] = map[string]interface{}{
			"allow_insecure":        AtlasDataSource.AllowInsecure,
			"collection":            AtlasDataSource.Collection,
			"collection_regex":      AtlasDataSource.CollectionRegex,
			"database":              AtlasDataSource.Database,
			"database_regex":        AtlasDataSource.DatabaseRegex,
			"default_format":        AtlasDataSource.DefaultFormat,
			"path":                  AtlasDataSource.Path,
			"provenance_field_name": AtlasDataSource.ProvenanceFieldName,
			"store_name":            AtlasDataSource.StoreName,
			"urls":                  AtlasDataSource.Urls,
		}
	}

	return out
}

func flattenDataFederationStores(stores []*matlas.DataFederationStore) []map[string]interface{} {
	store := make([]map[string]interface{}, 0)

	for i := range stores {
		store = append(store, map[string]interface{}{
			"name":                       stores[i].Name,
			"provider":                   stores[i].Provider,
			"region":                     stores[i].Region,
			"project_id":                 stores[i].ProjectID,
			"bucket":                     stores[i].Bucket,
			"cluster_name":               stores[i].ClusterName,
			"prefix":                     stores[i].Prefix,
			"delimiter":                  stores[i].Delimiter,
			"include_tags":               stores[i].IncludeTags,
			"additional_storage_classes": stores[i].AdditionalStorageClasses,
			"read_preference":            newReadPreferenceField(stores[i].ReadPreference),
		})
	}

	return store
}

func newReadPreferenceField(atlasReadPreference *matlas.ReadPreference) []map[string]interface{} {
	if atlasReadPreference == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"mode":                  atlasReadPreference.Mode,
			"max_staleness_seconds": atlasReadPreference.MaxStalenessSeconds,
			"tags":                  atlasReadPreference.TagSets,
		},
	}
}
