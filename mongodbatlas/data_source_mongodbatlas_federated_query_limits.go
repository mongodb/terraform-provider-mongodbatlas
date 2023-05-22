package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasFederatedDatabaseQueryLimits() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasFederatedDatabaseQueryLimitsRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: schemaMongoDBAtlasFederatedDatabaseQueryLimitDataSource(),
				},
			},
		},
	}
}

func dataSourceMongoDBAtlasFederatedDatabaseQueryLimitsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	projectID := d.Get("project_id").(string)
	tenantName := d.Get("tenant_name").(string)

	queryLimits, _, err := conn.DataFederation.ListQueryLimits(ctx, projectID, tenantName)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting federated database query limits for project (%s), tenant (%s), error: %s", projectID, tenantName, err))
	}

	if results := flattenFederatedDatabaseQueryLimits(projectID, tenantName, queryLimits); results != nil {
		if err := d.Set("results", results); err != nil {
			return diag.FromErr(fmt.Errorf(errorFederatedDatabaseQueryLimit, "results", projectID, err))
		}
	}

	d.SetId(id.UniqueId())

	return nil
}

func flattenFederatedDatabaseQueryLimits(projectID, tenantName string, queryLimits []*matlas.DataFederationQueryLimit) []map[string]interface{} {
	var federatedDatabaseQueryLimitMap []map[string]interface{}

	if len(queryLimits) > 0 {
		federatedDatabaseQueryLimitMap = make([]map[string]interface{}, len(queryLimits))

		for i := range queryLimits {
			federatedDatabaseQueryLimitMap[i] = map[string]interface{}{
				"project_id":         projectID,
				"tenant_name":        queryLimits[i].TenantName,
				"limit_name":         queryLimits[i].Name,
				"overrun_policy":     queryLimits[i].OverrunPolicy,
				"value":              queryLimits[i].Value,
				"current_usage":      queryLimits[i].CurrentUsage,
				"default_limit":      queryLimits[i].DefaultLimit,
				"last_modified_date": queryLimits[i].LastModifiedDate,
				"maximum_limit":      queryLimits[i].MaximumLimit,
			}
		}
	}

	return federatedDatabaseQueryLimitMap
}
