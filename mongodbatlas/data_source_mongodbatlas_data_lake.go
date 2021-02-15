package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasDataLake() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasDataLakeRead,
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
				Computed: true,
			},
			"aws_test_s3_bucket": {
				Type:     schema.TypeString,
				Computed: true,
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

func dataSourceMongoDBAtlasDataLakeRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)
	name := d.Get("name").(string)

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
