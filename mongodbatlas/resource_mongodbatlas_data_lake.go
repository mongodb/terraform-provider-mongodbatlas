package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/spf13/cast"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorDataLakeCreate  = "error creating MongoDB Atlas DataLake: %s"
	errorDataLakeRead    = "error reading MongoDB Atlas DataLake (%s): %s"
	errorDataLakeDelete  = "error deleting MongoDB Atlas DataLake (%s): %s"
	errorDataLakeUpdate  = "error updating MongoDB Atlas DataLake (%s): %s"
	errorDataLakeSetting = "error setting `%s` for MongoDB Atlas DataLake (%s): %s"
)

func resourceMongoDBAtlasDataLake() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBAtlasDataLakeCreate,
		ReadContext:   resourceMongoDBAtlasDataLakeRead,
		UpdateContext: resourceMongoDBAtlasDataLakeUpdate,
		DeleteContext: resourceMongoDBAtlasDataLakeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasDataLakeImportState,
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
				Computed: true,
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
			"hostnames": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_databases": schemaDataLakesDatabases(),
			"storage_stores":    schemaDataLakesStores(),
		},
	}
}

func schemaDataLakesDatabases() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"collections": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"name": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"data_sources": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"store_name": {
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
									},
								},
							},
						},
					},
				},
				"views": {
					Type:     schema.TypeList,
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

func schemaDataLakesStores() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
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
				"additional_storage_classes": {
					Type:     schema.TypeSet,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}

func resourceMongoDBAtlasDataLakeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)
	name := d.Get("name").(string)

	cloudConfig := &matlas.CloudProviderConfig{
		AWSConfig: expandDataLakeAwsBlock(d),
	}

	dataLakeReq := &matlas.DataLakeCreateRequest{
		CloudProviderConfig: cloudConfig,
		Name:                name,
	}

	dataLake, _, err := conn.DataLakes.Create(ctx, projectID, dataLakeReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakeCreate, err))
	}

	dataLake.CloudProviderConfig.AWSConfig.TestS3Bucket = cloudConfig.AWSConfig.TestS3Bucket

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"name":       dataLake.Name,
	}))

	return resourceMongoDBAtlasDataLakeRead(ctx, d, meta)
}

func resourceMongoDBAtlasDataLakeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	name := ids["name"]

	dataLake, resp, err := conn.DataLakes.Get(ctx, projectID, name)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorDataLakeRead, name, err))
	}

	values := flattenAWSBlock(&dataLake.CloudProviderConfig)
	if len(values) != 0 {
		if !counterEmptyValues(values[0]) {
			if value, ok := d.GetOk("aws"); ok {
				v := value.([]interface{})
				if len(v) != 0 {
					v1 := v[0].(map[string]interface{})
					values[0]["test_s3_bucket"] = cast.ToString(v1["test_s3_bucket"])
				}
			}

			if err = d.Set("aws", values); err != nil {
				return diag.FromErr(fmt.Errorf(errorDataLakeSetting, "aws", name, err))
			}
		}
	}

	if err := d.Set("data_process_region", flattenDataLakeProcessRegion(&dataLake.DataProcessRegion)); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakeSetting, "data_process_region", name, err))
	}

	if err := d.Set("hostnames", dataLake.Hostnames); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakeSetting, "hostnames", name, err))
	}

	if err := d.Set("state", dataLake.State); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakeSetting, "state", name, err))
	}

	if err := d.Set("storage_databases", flattenDataLakeStorageDatabases(dataLake.Storage.Databases)); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakeSetting, "storage_databases", name, err))
	}

	if err := d.Set("storage_stores", flattenDataLakeStorageStores(dataLake.Storage.Stores)); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakeSetting, "storage_stores", name, err))
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"name":       name,
	}))

	return nil
}

func resourceMongoDBAtlasDataLakeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	name := ids["name"]

	dataProcess := &matlas.DataProcessRegion{}
	awsConfig := matlas.AwsCloudProviderConfig{}

	if d.HasChange("aws_role_id") {
		awsConfig.RoleID = cast.ToString(d.Get("aws_role_id"))
	}

	if d.HasChange("aws_test_s3_bucket") {
		awsConfig.TestS3Bucket = cast.ToString(d.Get("aws_test_s3_bucket"))
	}

	if d.HasChange("data_process_region") {
		dataProcess = expandDataLakeDataProcessRegion(d)
	}

	dataLakeReq := &matlas.DataLakeUpdateRequest{
		CloudProviderConfig: &matlas.CloudProviderConfig{AWSConfig: awsConfig},
		DataProcessRegion:   dataProcess,
	}
	_, _, err := conn.DataLakes.Update(ctx, projectID, name, dataLakeReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakeUpdate, name, err))
	}

	return resourceMongoDBAtlasDataLakeRead(ctx, d, meta)
}

func resourceMongoDBAtlasDataLakeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	name := ids["name"]

	_, err := conn.DataLakes.Delete(ctx, projectID, name)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakeDelete, name, err))
	}

	return nil
}

func resourceMongoDBAtlasDataLakeImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*MongoDBClient).Atlas

	projectID, name, s3Bucket, err := splitDataLakeImportID(d.Id())
	if err != nil {
		return nil, err
	}

	u, _, err := conn.DataLakes.Get(ctx, projectID, name)
	if err != nil {
		return nil, fmt.Errorf("couldn't import data lake(%s) for project (%s), error: %s", name, projectID, err)
	}

	if err := d.Set("project_id", u.GroupID); err != nil {
		return nil, fmt.Errorf("error setting `project_id` for data lakes (%s): %s", d.Id(), err)
	}

	if err := d.Set("name", u.Name); err != nil {
		return nil, fmt.Errorf("error setting `name` for data lakes (%s): %s", d.Id(), err)
	}
	mapAws := make([]map[string]interface{}, 0)

	mapAws = append(mapAws, map[string]interface{}{
		"test_s3_bucket": s3Bucket,
	})

	if err := d.Set("aws", mapAws); err != nil {
		return nil, fmt.Errorf("error setting `aws` for data lakes (%s): %s", d.Id(), err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"name":       u.Name,
	}))

	return []*schema.ResourceData{d}, nil
}

func splitDataLakeImportID(id string) (projectID, name, s3Bucket string, err error) {
	var parts = strings.Split(id, "--")

	if len(parts) != 3 {
		err = errors.New("import format error: to import a Data Lake, use the format {project_id}--{name}--{test_s3_bucket}")
		return
	}

	projectID = parts[0]
	name = parts[1]
	s3Bucket = parts[2]

	return
}

func flattenAWSBlock(aws *matlas.CloudProviderConfig) []map[string]interface{} {
	if aws == nil {
		return nil
	}

	database := make([]map[string]interface{}, 0)

	database = append(database, map[string]interface{}{
		"role_id":              aws.AWSConfig.RoleID,
		"iam_assumed_role_arn": aws.AWSConfig.IAMAssumedRoleARN,
		"iam_user_arn":         aws.AWSConfig.IAMUserARN,
		"external_id":          aws.AWSConfig.ExternalID,
	})

	return database
}

func flattenDataLakeProcessRegion(processRegion *matlas.DataProcessRegion) []interface{} {
	if processRegion != nil && (processRegion.Region != "" || processRegion.CloudProvider != "") {
		return []interface{}{map[string]interface{}{
			"cloud_provider": processRegion.CloudProvider,
			"region":         processRegion.Region,
		}}
	}

	return []interface{}{}
}

func flattenDataLakeStorageDatabases(databases []matlas.DataLakeDatabase) []map[string]interface{} {
	database := make([]map[string]interface{}, 0)

	for _, db := range databases {
		database = append(database, map[string]interface{}{
			"name":                     db.Name,
			"collections":              flattenDataLakeStorageDatabaseCollections(db.Collections),
			"views":                    flattenDataLakeStorageDatabaseViews(db.Views),
			"max_wildcard_collections": db.MaxWildcardCollections,
		})
	}

	return database
}

func flattenDataLakeStorageDatabaseCollections(collections []matlas.DataLakeCollection) []map[string]interface{} {
	database := make([]map[string]interface{}, 0)

	for _, db := range collections {
		database = append(database, map[string]interface{}{
			"name":         db.Name,
			"data_sources": flattenDataLakeStorageDatabaseCollectionsDataSources(db.DataSources),
		})
	}

	return database
}

func flattenDataLakeStorageDatabaseCollectionsDataSources(dataSources []matlas.DataLakeDataSource) []map[string]interface{} {
	database := make([]map[string]interface{}, 0)

	for _, db := range dataSources {
		database = append(database, map[string]interface{}{
			"store_name":     db.StoreName,
			"default_format": db.DefaultFormat,
			"path":           db.Path,
		})
	}

	return database
}

func flattenDataLakeStorageDatabaseViews(views []matlas.DataLakeDatabaseView) []map[string]interface{} {
	view := make([]map[string]interface{}, 0)

	for _, db := range views {
		view = append(view, map[string]interface{}{
			"name":     db.Name,
			"source":   db.Source,
			"pipeline": db.Pipeline,
		})
	}

	return view
}

func flattenDataLakeStorageStores(stores []matlas.DataLakeStore) []map[string]interface{} {
	store := make([]map[string]interface{}, 0)

	for i := range stores {
		store = append(store, map[string]interface{}{
			"name":                       stores[i].Name,
			"provider":                   stores[i].Provider,
			"region":                     stores[i].Region,
			"bucket":                     stores[i].Bucket,
			"prefix":                     stores[i].Prefix,
			"delimiter":                  stores[i].Delimiter,
			"include_tags":               stores[i].IncludeTags,
			"additional_storage_classes": stores[i].AdditionalStorageClasses,
		})
	}

	return store
}

func expandDataLakeAwsBlock(d *schema.ResourceData) matlas.AwsCloudProviderConfig {
	aws := matlas.AwsCloudProviderConfig{}
	if value, ok := d.GetOk("aws"); ok {
		v := value.([]interface{})
		if len(v) != 0 {
			v1 := v[0].(map[string]interface{})

			aws.RoleID = cast.ToString(v1["role_id"])
			aws.TestS3Bucket = cast.ToString(v1["test_s3_bucket"])
		}
	}
	return aws
}

func expandDataLakeDataProcessRegion(d *schema.ResourceData) *matlas.DataProcessRegion {
	if value, ok := d.GetOk("data_process_region"); ok {
		vL := value.([]interface{})

		if len(vL) != 0 {
			v := vL[0].(map[string]interface{})

			return &matlas.DataProcessRegion{
				CloudProvider: cast.ToString(v["cloud_provider"]),
				Region:        cast.ToString(v["region"]),
			}
		}
	}
	return nil
}
