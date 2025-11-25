package thirdpartyintegration

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func PluralDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: pluralDataSourceRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     thirdPartyIntegrationSchema(),
			},
		},
	}
}

func pluralDataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	projectID := d.Get("project_id").(string)
	integrations, _, err := connV2.ThirdPartyIntegrationsApi.ListGroupIntegrations(ctx, projectID).Execute()

	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting third party integration list: %s", err))
	}

	if err = d.Set("results", flattenIntegrations(d, integrations, projectID)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting results for third party integrations %s", err))
	}

	d.SetId(id.UniqueId())

	return nil
}
