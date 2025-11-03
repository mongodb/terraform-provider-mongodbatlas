package federateddatabaseinstance

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20250312008/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func PluralDataSource() *schema.Resource {
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
						"cloud_provider_config": cloudProviderConfig(true),
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

func dataSourceMongoDBAtlasFederatedDatabaseInstancesRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	projectID := d.Get("project_id").(string)

	federatedDatabaseInstances, _, err := connV2.DataFederationApi.ListDataFederation(ctx, projectID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting MongoDB Atlas Federated Database Instances information: %s", err))
	}

	if results := flattenFederatedDatabaseInstances(d, projectID, federatedDatabaseInstances); results != nil {
		if err := d.Set("results", results); err != nil {
			return diag.FromErr(fmt.Errorf(errorFederatedDatabaseInstanceSetting, "results", projectID, err))
		}
	}

	d.SetId(id.UniqueId())

	return nil
}

func flattenFederatedDatabaseInstances(d *schema.ResourceData, projectID string, federatedDatabaseInstances []admin.DataLakeTenant) []map[string]any {
	var federatedDatabaseInstancesMap []map[string]any

	if len(federatedDatabaseInstances) > 0 {
		federatedDatabaseInstancesMap = make([]map[string]any, len(federatedDatabaseInstances))

		for i := range federatedDatabaseInstances {
			federatedDatabaseInstancesMap[i] = map[string]any{
				"project_id":            projectID,
				"name":                  federatedDatabaseInstances[i].GetName(),
				"state":                 federatedDatabaseInstances[i].GetState(),
				"hostnames":             federatedDatabaseInstances[i].GetHostnames(),
				"cloud_provider_config": flattenCloudProviderConfig(d, federatedDatabaseInstances[i].CloudProviderConfig),
				"data_process_region":   flattenDataProcessRegion(federatedDatabaseInstances[i].DataProcessRegion),
				"storage_databases":     flattenDataFederationDatabase(federatedDatabaseInstances[i].Storage.GetDatabases()),
				"storage_stores":        flattenDataFederationStores(federatedDatabaseInstances[i].Storage.GetStores()),
			}
		}
	}

	return federatedDatabaseInstancesMap
}
