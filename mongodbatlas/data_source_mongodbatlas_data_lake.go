package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMongoDBAtlasDataLake() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasDataLakeRead,
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
	}
}

func dataSourceMongoDBAtlasDataLakeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)
	name := d.Get("name").(string)

	dataLake, resp, err := conn.DataLakes.Get(ctx, projectID, name)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorDataLakeRead, name, err))
	}

	if err := d.Set("aws", flattenAWSBlock(&dataLake.CloudProviderConfig)); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakeSetting, "aws", name, err))
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
