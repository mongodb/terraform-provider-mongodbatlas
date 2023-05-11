package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasFederatedDatabaseInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasFederatedDatabaseInstancesRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
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
						"storage_databases": schemaFederatedDatabaseInstanceDatabasesDataSource(),
						"storage_stores":    schemaFederatedDatabaseInstanceStoresDataSource(),
					},
				},
			},
		},
	}
}

func dataSourceMongoDBAtlasFederatedDatabaseInstancesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	projectID := d.Get("project_id").(string)

	federatedDatabaseInstances, _, err := conn.DataFederation.List(ctx, projectID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting MongoDB Atlas Federated Database Instances information: %s", err))
	}

	if results := flattenFederatedDatabaseInstances(projectID, federatedDatabaseInstances); results != nil {
		if err := d.Set("results", results); err != nil {
			return diag.FromErr(fmt.Errorf(errorFederatedDatabaseInstanceSetting, "results", projectID, err))
		}
	}

	d.SetId(id.UniqueId())

	return nil
}

func flattenFederatedDatabaseInstances(projectID string, federatedDatabaseInstances []*matlas.DataFederationInstance) []map[string]interface{} {
	var federatedDatabaseInstancesMap []map[string]interface{}

	if len(federatedDatabaseInstances) > 0 {
		federatedDatabaseInstancesMap = make([]map[string]interface{}, len(federatedDatabaseInstances))

		for i := range federatedDatabaseInstances {
			federatedDatabaseInstancesMap[i] = map[string]interface{}{
				"project_id":          projectID,
				"name":                federatedDatabaseInstances[i].Name,
				"aws":                 flattenAWSBlock(federatedDatabaseInstances[i].CloudProviderConfig),
				"data_process_region": flattenDataProcessRegion(federatedDatabaseInstances[i].DataProcessRegion),
				"storage_databases":   flattenDataFederationDatabase(federatedDatabaseInstances[i].Storage.Databases),
				"storage_stores":      flattenDataFederationStores(federatedDatabaseInstances[i].Storage.Stores),
			}
		}
	}

	return federatedDatabaseInstancesMap
}
