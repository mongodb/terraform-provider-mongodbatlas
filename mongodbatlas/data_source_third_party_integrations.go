package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func DataSourceThirdPartyIntegrations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasThirdPartyIntegrationsRead,
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

func dataSourceMongoDBAtlasThirdPartyIntegrationsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).Atlas

	projectID := d.Get("project_id").(string)
	integrations, _, err := conn.Integrations.List(ctx, projectID)

	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting third party integration list: %s", err))
	}

	if err = d.Set("results", flattenIntegrations(d, integrations, projectID)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting results for third party integrations %s", err))
	}

	d.SetId(id.UniqueId())

	return nil
}

func flattenIntegrations(d *schema.ResourceData, integrations *matlas.ThirdPartyIntegrations, projectID string) (list []map[string]any) {
	if len(integrations.Results) == 0 {
		return
	}

	list = make([]map[string]any, 0, len(integrations.Results))

	for _, integration := range integrations.Results {
		service := integrationToSchema(d, integration)
		service["project_id"] = projectID
		list = append(list, service)
	}

	return
}

func integrationToSchema(d *schema.ResourceData, integration *matlas.ThirdPartyIntegration) map[string]any {
	integrationSchema := schemaToIntegration(d)
	if integrationSchema.APIKey == "" {
		integrationSchema.APIKey = integration.APIKey
	}
	if integrationSchema.ServiceKey == "" {
		integrationSchema.APIKey = integration.ServiceKey
	}
	if integrationSchema.APIToken == "" {
		integrationSchema.APIKey = integration.APIToken
	}
	if integrationSchema.RoutingKey == "" {
		integrationSchema.APIKey = integration.RoutingKey
	}
	if integrationSchema.Secret == "" {
		integrationSchema.APIKey = integration.Secret
	}
	if integrationSchema.Password == "" {
		integrationSchema.APIKey = integration.Password
	}
	if integrationSchema.UserName == "" {
		integrationSchema.APIKey = integration.UserName
	}
	if integrationSchema.URL == "" {
		integrationSchema.URL = integration.URL
	}

	out := map[string]any{
		"type":                        integration.Type,
		"api_key":                     integrationSchema.APIKey,
		"region":                      integration.Region,
		"service_key":                 integrationSchema.ServiceKey,
		"team_name":                   integration.TeamName,
		"channel_name":                integration.ChannelName,
		"routing_key":                 integrationSchema.RoutingKey,
		"url":                         integrationSchema.URL,
		"secret":                      integrationSchema.Secret,
		"microsoft_teams_webhook_url": integrationSchema.MicrosoftTeamsWebhookURL,
		"user_name":                   integrationSchema.UserName,
		"password":                    integrationSchema.Password,
		"service_discovery":           integration.ServiceDiscovery,
		"scheme":                      integration.Scheme,
		"enabled":                     integration.Enabled,
	}

	// removing optional empty values, terraform complains about unexpected values even though they're empty
	optionals := []string{"api_key", "region", "service_key",
		"team_name", "channel_name", "url", "secret", "password"}

	for _, attr := range optionals {
		if val, ok := out[attr]; ok {
			strval, okT := val.(string)
			if okT && strval == "" {
				delete(out, attr)
			}
		}
	}

	return out
}

func schemaToIntegration(in *schema.ResourceData) (out *matlas.ThirdPartyIntegration) {
	out = &matlas.ThirdPartyIntegration{}

	if _type, ok := in.GetOk("type"); ok {
		out.Type = _type.(string)
	}

	if apiKey, ok := in.GetOk("api_key"); ok {
		out.APIKey = apiKey.(string)
	}

	if region, ok := in.GetOk("region"); ok {
		out.Region = region.(string)
	}

	if serviceKey, ok := in.GetOk("service_key"); ok {
		out.ServiceKey = serviceKey.(string)
	}

	if teamName, ok := in.GetOk("team_name"); ok {
		out.TeamName = teamName.(string)
	}

	if channelName, ok := in.GetOk("channel_name"); ok {
		out.ChannelName = channelName.(string)
	}

	if routingKey, ok := in.GetOk("routing_key"); ok {
		out.RoutingKey = routingKey.(string)
	}

	if url, ok := in.GetOk("url"); ok {
		out.URL = url.(string)
	}

	if secret, ok := in.GetOk("secret"); ok {
		out.Secret = secret.(string)
	}

	if microsoftTeamsWebhookURL, ok := in.GetOk("microsoft_teams_webhook_url"); ok {
		out.MicrosoftTeamsWebhookURL = microsoftTeamsWebhookURL.(string)
	}

	if userName, ok := in.GetOk("user_name"); ok {
		out.UserName = userName.(string)
	}

	if password, ok := in.GetOk("password"); ok {
		out.Password = password.(string)
	}

	if serviceDiscovery, ok := in.GetOk("service_discovery"); ok {
		out.ServiceDiscovery = serviceDiscovery.(string)
	}

	if scheme, ok := in.GetOk("scheme"); ok {
		out.Scheme = scheme.(string)
	}

	if enabled, ok := in.GetOk("enabled"); ok {
		out.Enabled = enabled.(bool)
	}

	return out
}

func updateIntegrationFromSchema(d *schema.ResourceData, integration *matlas.ThirdPartyIntegration) {
	integration.APIKey = d.Get("api_key").(string)

	if d.HasChange("region") {
		integration.Region = d.Get("region").(string)
	}

	if d.HasChange("service_key") {
		integration.ServiceKey = d.Get("service_key").(string)
	}

	if d.HasChange("team_name") {
		integration.TeamName = d.Get("team_name").(string)
	}

	if d.HasChange("channel_name") {
		integration.ChannelName = d.Get("channel_name").(string)
	}

	if d.HasChange("routing_key") {
		integration.RoutingKey = d.Get("routing_key").(string)
	}

	if d.HasChange("url") {
		integration.URL = d.Get("url").(string)
	}

	if d.HasChange("secret") {
		integration.Secret = d.Get("secret").(string)
	}

	if d.HasChange("microsoft_teams_webhook_url") {
		integration.MicrosoftTeamsWebhookURL = d.Get("microsoft_teams_webhook_url").(string)
	}

	if d.HasChange("user_name") {
		integration.UserName = d.Get("user_name").(string)
	}

	if d.HasChange("password") {
		integration.Password = d.Get("password").(string)
	}

	if d.HasChange("service_discovery") {
		integration.ServiceDiscovery = d.Get("service_discovery").(string)
	}

	if d.HasChange("scheme") {
		integration.Scheme = d.Get("scheme").(string)
	}

	if d.HasChange("enabled") {
		integration.Enabled = d.Get("enabled").(bool)
	}
}
