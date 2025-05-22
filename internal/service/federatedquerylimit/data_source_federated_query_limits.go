package federatedquerylimit

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312003/admin"
)

func PluralDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcesRead,
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
					Schema: schemaDataSource(),
				},
			},
		},
	}
}

func dataSourcesRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2

	projectID := d.Get("project_id").(string)
	tenantName := d.Get("tenant_name").(string)

	queryLimits, _, err := conn.DataFederationApi.ReturnFederatedDatabaseQueryLimits(ctx, projectID, tenantName).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting federated database query limits for project (%s), tenant (%s), error: %s", projectID, tenantName, err))
	}

	if results := flattenFederatedDatabaseQueryLimits(projectID, queryLimits); results != nil {
		if err := d.Set("results", results); err != nil {
			return diag.FromErr(fmt.Errorf(errorFederatedDatabaseQueryLimit, "results", projectID, err))
		}
	}

	d.SetId(id.UniqueId())

	return nil
}

func flattenFederatedDatabaseQueryLimits(projectID string, queryLimits []admin.DataFederationTenantQueryLimit) []map[string]any {
	var federatedDatabaseQueryLimitMap []map[string]any
	if len(queryLimits) == 0 {
		return federatedDatabaseQueryLimitMap
	}

	federatedDatabaseQueryLimitMap = make([]map[string]any, len(queryLimits))
	for i := range queryLimits {
		federatedDatabaseQueryLimitMap[i] = map[string]any{
			"project_id":         projectID,
			"tenant_name":        queryLimits[i].GetTenantName(),
			"limit_name":         queryLimits[i].Name,
			"overrun_policy":     queryLimits[i].GetOverrunPolicy(),
			"value":              queryLimits[i].Value,
			"current_usage":      queryLimits[i].GetCurrentUsage(),
			"default_limit":      queryLimits[i].GetDefaultLimit(),
			"last_modified_date": conversion.TimeToString(queryLimits[i].GetLastModifiedDate()),
			"maximum_limit":      queryLimits[i].GetMaximumLimit(),
		}
	}

	return federatedDatabaseQueryLimitMap
}
