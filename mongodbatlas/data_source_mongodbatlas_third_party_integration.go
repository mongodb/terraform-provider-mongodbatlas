package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMongoDBAtlasThirdPartyIntegration() *schema.Resource {
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

	integration.ReadContext = dataSourceMongoDBAtlasThirdPartyIntegrationRead

	return integration
}

func thirdPartyIntegrationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"license_key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"write_token": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"read_token": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
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
			"api_token": {
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
			"flow_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"org_name": {
				Type:     schema.TypeString,
				Computed: true,
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
		},
	}
}

func dataSourceMongoDBAtlasThirdPartyIntegrationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectID := d.Get("project_id").(string)
	queryType := d.Get("type").(string)

	conn := meta.(*MongoDBClient).Atlas

	integration, _, err := conn.Integrations.Get(ctx, projectID, queryType)

	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting third party integration for type %s %w", queryType, err))
	}

	fieldMap := integrationToSchema(integration)

	for property, value := range fieldMap {
		if err = d.Set(property, value); err != nil {
			return diag.FromErr(fmt.Errorf("error setting %s for third party integration %w", property, err))
		}
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"type":       queryType,
	}))

	return nil
}
