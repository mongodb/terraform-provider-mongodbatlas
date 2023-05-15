package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMongoDBAtlasFederatedDatabaseQueryLimit() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasFederatedDatabaseQueryLimitRead,
		Schema:      schemaMongoDBAtlasFederatedDatabaseQueryLimitDataSource(),
	}
}

func schemaMongoDBAtlasFederatedDatabaseQueryLimitDataSource() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"project_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		// check if required or not
		"limit_name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"tenant_name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"overrun_policy": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"value": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"current_usage": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"default_limit": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"last_modified_date": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"maximum_limit": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

func dataSourceMongoDBAtlasFederatedDatabaseQueryLimitRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	projectID := d.Get("project_id").(string)
	tenantName := d.Get("tenant_name").(string)
	limitName := d.Get("limit_name").(string)

	queryLimit, _, err := conn.DataFederation.GetQueryLimit(ctx, projectID, tenantName, limitName)

	if err != nil {
		return diag.FromErr(fmt.Errorf("couldn't import federated database query limit(%s) for project (%s), tenant (%s), error: %s", limitName, projectID, tenantName, err))
	}

	err = setResourceFieldsFromFederatedDatabaseQueryLimit(d, projectID, queryLimit)
	if err != nil {
		return diag.FromErr(err)
	}

	// if err := d.Set("project_id", projectID); err != nil {
	// 	return diag.FromErr(fmt.Errorf(errorFederatedDatabaseQueryLimit, "project_id", d.Id(), err))
	// }

	// if err := d.Set("limit_name", queryLimit.Name); err != nil {
	// 	return diag.FromErr(fmt.Errorf(errorFederatedDatabaseQueryLimit, "limit_name", d.Id(), err))
	// }

	// if err := d.Set("tenant_name", queryLimit.TenantName); err != nil {
	// 	return diag.FromErr(fmt.Errorf(errorFederatedDatabaseQueryLimit, "tenant_name", d.Id(), err))
	// }

	// if err := d.Set("overrun_policy", queryLimit.OverrunPolicy); err != nil {
	// 	return diag.FromErr(fmt.Errorf(errorFederatedDatabaseQueryLimit, "overrun_policy", d.Id(), err))
	// }

	// if err := d.Set("value", queryLimit.Value); err != nil {
	// 	return diag.FromErr(fmt.Errorf(errorFederatedDatabaseQueryLimit, "value", d.Id(), err))
	// }

	// if err := d.Set("current_usage", queryLimit.CurrentUsage); err != nil {
	// 	return diag.FromErr(fmt.Errorf(errorFederatedDatabaseQueryLimit, "current_usage", d.Id(), err))
	// }

	// if err := d.Set("default_limit", queryLimit.DefaultLimit); err != nil {
	// 	return diag.FromErr(fmt.Errorf(errorFederatedDatabaseQueryLimit, "default_limit", d.Id(), err))
	// }

	// if err := d.Set("last_modified_date", queryLimit.LastModifiedDate); err != nil {
	// 	return diag.FromErr(fmt.Errorf(errorFederatedDatabaseQueryLimit, "last_modified_date", d.Id(), err))
	// }

	// if err := d.Set("maximum_limit", queryLimit.MaximumLimit); err != nil {
	// 	return diag.FromErr(fmt.Errorf(errorFederatedDatabaseQueryLimit, "maximum_limit", d.Id(), err))
	// }

	d.SetId(encodeStateID(map[string]string{
		"project_id":  projectID,
		"tenant_name": queryLimit.TenantName,
		"limit_name":  queryLimit.Name,
	}))

	return nil
}
