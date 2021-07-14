package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasDataLakes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasDataLakesRead,
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
						"aws": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"role_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"test_s3_bucket": {
										Type:     schema.TypeString,
										Computed: true,
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

func dataSourceMongoDBAtlasDataLakesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas

	projectID := d.Get("project_id").(string)

	dataLakes, _, err := conn.DataLakes.List(ctx, projectID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting MongoDB Atlas Data Lakes information: %s", err))
	}

	if err := d.Set("results", flattenDataLakes(dataLakes)); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakeSetting, "results", projectID, err))
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
				"project_id":          dataLakes[i].GroupID,
				"name":                dataLakes[i].Name,
				"aws":                 flattenAWSBlock(&dataLakes[i].CloudProviderConfig),
				"data_process_region": flattenDataLakeProcessRegion(&dataLakes[i].DataProcessRegion),
				"hostnames":           dataLakes[i].Hostnames,
				"state":               dataLakes[i].State,
				"storage_databases":   flattenDataLakeStorageDatabases(dataLakes[i].Storage.Databases),
				"storage_stores":      flattenDataLakeStorageStores(dataLakes[i].Storage.Stores),
			}
		}
	}

	return dataLakesMap
}
