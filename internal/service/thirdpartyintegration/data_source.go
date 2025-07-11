package thirdpartyintegration

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func DataSource() *schema.Resource {
	integration := thirdPartyIntegrationSchema()
	integration.Schema["project_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}

	integration.Schema["type"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Third-party service integration identifier",
	}

	integration.ReadContext = dataSourceRead

	return integration
}

func thirdPartyIntegrationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"api_key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"team_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"channel_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"routing_key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"microsoft_teams_webhook_url": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"user_name": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"service_discovery": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"send_collection_latency_metrics": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"send_database_metrics": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"send_user_provided_resource_tags": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	projectID := d.Get("project_id").(string)
	queryType := d.Get("type").(string)

	connV2 := meta.(*config.MongoDBClient).AtlasV2

	integration, _, err := connV2.ThirdPartyIntegrationsApi.GetThirdPartyIntegration(ctx, projectID, queryType).Execute()

	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting third party integration for type %s %w", queryType, err))
	}

	fieldMap := integrationToSchema(d, integration)

	for property, value := range fieldMap {
		if err = d.Set(property, value); err != nil {
			return diag.FromErr(fmt.Errorf("error setting %s for third party integration %w", property, err))
		}
	}

	d.SetId(integration.GetId())
	return nil
}
