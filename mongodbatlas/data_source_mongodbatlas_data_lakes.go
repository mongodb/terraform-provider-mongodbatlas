package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasDataLakes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasDataLakesRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
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
				},
			},
		},
	}
}

func dataSourceMongoDBAtlasDataLakesRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*matlas.Client)

	projectID := d.Get("project_id").(string)

	dataLakes, _, err := conn.DataLakes.List(context.Background(), projectID)
	if err != nil {
		return fmt.Errorf("error getting database users information: %s", err)
	}

	if err := d.Set("results", flattenDataLakes(dataLakes)); err != nil {
		return fmt.Errorf(errorDataLakeSetting, "results", projectID, err)
	}

	d.SetId(resource.UniqueId())

	return nil
}

func flattenDataLakes(dataLakes []matlas.DataLake) []map[string]interface{} {
	var dataLakesMap []map[string]interface{}

	if len(dataLakes) > 0 {
		dataLakesMap = make([]map[string]interface{}, len(dataLakes))

		for i := range dataLakes {
			dataLakesMap[i] = map[string]interface{}{
				"project_id":               dataLakes[i].GroupID,
				"name":                     dataLakes[i].Name,
				"aws_role_id":              dataLakes[i].CloudProviderConfig.AWSConfig.RoleID,
				"aws_iam_assumed_role_arn": dataLakes[i].CloudProviderConfig.AWSConfig.IAMAssumedRoleARN,
				"aws_iam_user_arn":         dataLakes[i].CloudProviderConfig.AWSConfig.IAMUserARN,
				"aws_external_id":          dataLakes[i].CloudProviderConfig.AWSConfig.ExternalID,
				"data_process_region":      flattenDataLakeProcessRegion(&dataLakes[i].DataProcessRegion),
				"hostnames":                dataLakes[i].Hostnames,
				"state":                    dataLakes[i].State,
				"storage_databases":        flattenDataLakeStorageDatabases(dataLakes[i].Storage.Databases),
				"storage_stores":           flattenDataLakeStorageStores(dataLakes[i].Storage.Stores),
			}
		}
	}

	return dataLakesMap
}
