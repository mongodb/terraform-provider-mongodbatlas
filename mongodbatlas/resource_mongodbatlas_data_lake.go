package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/cast"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorDataLakeCreate  = "error creating MongoDB DataLake: %s"
	errorDataLakeRead    = "error reading MongoDB DataLake (%s): %s"
	errorDataLakeDelete  = "error deleting MongoDB DataLake (%s): %s"
	errorDataLakeUpdate  = "error updating MongoDB DataLake (%s): %s"
	errorDataLakeSetting = "error setting `%s` for MongoDB DataLake (%s): %s"
)

func resourceMongoDBAtlasDataLake() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoDBAtlasDataLakeCreate,
		Read:   resourceMongoDBAtlasDataLakeRead,
		Update: resourceMongoDBAtlasDataLakeUpdate,
		Delete: resourceMongoDBAtlasDataLakeDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMongoDBAtlasDataLakeImportState,
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
			"aws_role_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"aws_test_s3_bucket": {
				Type:     schema.TypeString,
				Required: true,
			},
			"aws_iam_assumed_role_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"aws_iam_user_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"aws_external_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_process_region": {
				Type:     schema.TypeMap,
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
			},
		},
	}
}

func resourceMongoDBAtlasDataLakeCreate(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)
	name := d.Get("name").(string)

	cloudConfig := &matlas.CloudProviderConfig{
		AWSConfig: matlas.AwsCloudProviderConfig{RoleID: d.Get("aws_role_id").(string), TestS3Bucket: d.Get("aws_test_s3_bucket").(string)},
	}

	dataLakeReq := &matlas.DataLakeCreateRequest{
		CloudProviderConfig: cloudConfig,
		Name:                name,
	}

	dataLake, _, err := conn.DataLakes.Create(context.Background(), projectID, dataLakeReq)
	if err != nil {
		return fmt.Errorf(errorDataLakeCreate, err)
	}

	dataLake.CloudProviderConfig.AWSConfig.TestS3Bucket = cloudConfig.AWSConfig.TestS3Bucket

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"name":       dataLake.Name,
	}))

	return resourceMongoDBAtlasDataLakeRead(d, meta)
}

func resourceMongoDBAtlasDataLakeRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	name := ids["name"]

	dataLake, resp, err := conn.DataLakes.Get(context.Background(), projectID, name)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return fmt.Errorf(errorDataLakeRead, name, err)
	}

	if err := d.Set("aws_role_id", dataLake.CloudProviderConfig.AWSConfig.RoleID); err != nil {
		return fmt.Errorf(errorDataLakeSetting, "aws_role_id", name, err)
	}

	if err := d.Set("aws_iam_assumed_role_arn", dataLake.CloudProviderConfig.AWSConfig.IAMAssumedRoleARN); err != nil {
		return fmt.Errorf(errorDataLakeSetting, "aws_iam_assumed_role_arn", name, err)
	}

	if err := d.Set("aws_iam_user_arn", dataLake.CloudProviderConfig.AWSConfig.IAMUserARN); err != nil {
		return fmt.Errorf(errorDataLakeSetting, "aws_iam_user_arn", name, err)
	}

	if err := d.Set("aws_external_id", dataLake.CloudProviderConfig.AWSConfig.ExternalID); err != nil {
		return fmt.Errorf(errorDataLakeSetting, "aws_external_id", name, err)
	}

	if err := d.Set("data_process_region", flattenDataLakeProcessRegion(&dataLake.DataProcessRegion)); err != nil {
		return fmt.Errorf(errorDataLakeSetting, "data_process_region", name, err)
	}

	if err := d.Set("hostnames", dataLake.Hostnames); err != nil {
		return fmt.Errorf(errorDataLakeSetting, "hostnames", name, err)
	}

	if err := d.Set("state", dataLake.State); err != nil {
		return fmt.Errorf(errorDataLakeSetting, "state", name, err)
	}

	if err := d.Set("storage_databases", flattenDataLakeStorageDatabases(dataLake.Storage.Databases)); err != nil {
		return fmt.Errorf(errorDataLakeSetting, "storage_databases", name, err)
	}

	if err := d.Set("storage_stores", flattenDataLakeStorageStores(dataLake.Storage.Stores)); err != nil {
		return fmt.Errorf(errorDataLakeSetting, "storage_stores", name, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"name":       name,
	}))

	return nil
}

func resourceMongoDBAtlasDataLakeUpdate(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*matlas.Client)
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
	_, _, err := conn.DataLakes.Update(context.Background(), projectID, name, dataLakeReq)
	if err != nil {
		return fmt.Errorf(errorDataLakeUpdate, name, err)
	}

	return resourceMongoDBAtlasDataLakeRead(d, meta)
}

func resourceMongoDBAtlasDataLakeDelete(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	name := ids["name"]

	_, err := conn.DataLakes.Delete(context.Background(), projectID, name)
	if err != nil {
		return fmt.Errorf(errorDataLakeDelete, name, err)
	}

	return nil
}

func resourceMongoDBAtlasDataLakeImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*matlas.Client)

	projectID, name, s3Bucket, err := splitDataLakeImportID(d.Id())
	if err != nil {
		return nil, err
	}

	u, _, err := conn.DataLakes.Get(context.Background(), projectID, name)
	if err != nil {
		return nil, fmt.Errorf("couldn't import user(%s) in data lake(%s), error: %s", name, projectID, err)
	}

	if err := d.Set("project_id", u.GroupID); err != nil {
		return nil, fmt.Errorf("error setting `project_id` for data lakes (%s): %s", d.Id(), err)
	}

	if err := d.Set("name", u.Name); err != nil {
		return nil, fmt.Errorf("error setting `name` for data lakes (%s): %s", d.Id(), err)
	}

	if err := d.Set("aws_test_s3_bucket", s3Bucket); err != nil {
		return nil, fmt.Errorf("error setting `aws_test_s3_bucket` for data lakes (%s): %s", d.Id(), err)
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

func flattenDataLakeProcessRegion(processRegion *matlas.DataProcessRegion) map[string]interface{} {
	if processRegion != nil && (processRegion.Region != "" || processRegion.CloudProvider != "") {
		return map[string]interface{}{
			"cloud_provider": processRegion.CloudProvider,
			"region":         processRegion.Region,
		}
	}

	return map[string]interface{}{}
}

func flattenDataLakeStorageDatabases(databases []matlas.DataLakeDatabase) []map[string]interface{} {
	database := make([]map[string]interface{}, 0)

	for _, db := range databases {
		database = append(database, map[string]interface{}{
			"name":        db.Name,
			"collections": flattenDataLakeStorageDatabaseCollections(db.Collections),
			"views":       flattenDataLakeStorageDatabaseViews(db.Views),
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

	for _, db := range stores {
		store = append(store, map[string]interface{}{
			"name":         db.Name,
			"provider":     db.Provider,
			"region":       db.Region,
			"bucket":       db.Bucket,
			"prefix":       db.Prefix,
			"delimiter":    db.Delimiter,
			"include_tags": db.IncludeTags,
		})
	}

	return store
}

func expandDataLakeDataProcessRegion(d *schema.ResourceData) *matlas.DataProcessRegion {
	if value, ok := d.GetOk("data_process_region"); ok {
		v := value.(map[string]interface{})

		return &matlas.DataProcessRegion{
			CloudProvider: cast.ToString(v["cloud_provider"]),
			Region:        cast.ToString(v["region"]),
		}
	}
	return nil
}
