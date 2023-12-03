package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func DataSourceFederatedDatabaseQueryLimit() *schema.Resource {
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

func dataSourceMongoDBAtlasFederatedDatabaseQueryLimitRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).Atlas

	projectID := d.Get("project_id").(string)
	tenantName := d.Get("tenant_name").(string)
	limitName := d.Get("limit_name").(string)

	queryLimit, _, err := conn.DataFederation.GetQueryLimit(ctx, projectID, tenantName, limitName)

	if err != nil {
		return diag.FromErr(fmt.Errorf("couldn't import federated database query limit(%s) for project (%s), tenant (%s), error: %s", limitName, projectID, tenantName, err))
	}

	if err = setResourceFieldsFromFederatedDatabaseQueryLimit(d, projectID, queryLimit); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":  projectID,
		"tenant_name": queryLimit.TenantName,
		"limit_name":  queryLimit.Name,
	}))

	return nil
}
