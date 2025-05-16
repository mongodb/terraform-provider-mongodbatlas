package thirdpartyintegration

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20250312003/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func PluralDataSource() *schema.Resource {
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
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	projectID := d.Get("project_id").(string)
	integrations, _, err := connV2.ThirdPartyIntegrationsApi.ListThirdPartyIntegrations(ctx, projectID).Execute()

	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting third party integration list: %s", err))
	}

	if err = d.Set("results", flattenIntegrations(d, integrations, projectID)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting results for third party integrations %s", err))
	}

	d.SetId(id.UniqueId())

	return nil
}

func flattenIntegrations(d *schema.ResourceData, integrations *admin.PaginatedIntegration, projectID string) (list []map[string]any) {
	results := integrations.GetResults()
	if len(results) == 0 {
		return
	}

	list = make([]map[string]any, 0, len(results))

	for i := range results {
		service := integrationToSchema(d, &results[i])
		service["project_id"] = projectID
		list = append(list, service)
	}

	return
}

func integrationToSchema(d *schema.ResourceData, integration *admin.ThirdPartyIntegration) map[string]any {
	integrationSchema := schemaToIntegration(d)

	if integrationSchema.ApiKey == nil {
		integrationSchema.ApiKey = integration.ApiKey
	}
	if integrationSchema.ServiceKey == nil {
		integrationSchema.ServiceKey = integration.ServiceKey
	}
	if integrationSchema.ApiToken == nil {
		integrationSchema.ApiToken = integration.ApiToken
	}
	if integrationSchema.RoutingKey == nil {
		integrationSchema.RoutingKey = integration.RoutingKey
	}
	if integrationSchema.Secret == nil {
		integrationSchema.Secret = integration.Secret
	}
	if integrationSchema.Password == nil {
		integrationSchema.Password = integration.Password
	}
	if integrationSchema.Url == nil {
		integrationSchema.Url = integration.Url
	}
	if integrationSchema.SendCollectionLatencyMetrics == nil {
		integrationSchema.SendCollectionLatencyMetrics = integration.SendCollectionLatencyMetrics
	}
	if integrationSchema.SendDatabaseMetrics == nil {
		integrationSchema.SendDatabaseMetrics = integration.SendDatabaseMetrics
	}

	out := map[string]any{
		"id":                              integration.Id,
		"type":                            integration.Type,
		"api_key":                         integrationSchema.ApiKey,
		"region":                          integration.Region,
		"service_key":                     integrationSchema.ServiceKey,
		"team_name":                       integration.TeamName,
		"channel_name":                    integration.ChannelName,
		"routing_key":                     integration.RoutingKey,
		"url":                             integrationSchema.Url,
		"secret":                          integrationSchema.Secret,
		"microsoft_teams_webhook_url":     integrationSchema.MicrosoftTeamsWebhookUrl,
		"user_name":                       integration.Username,
		"password":                        integrationSchema.Password,
		"service_discovery":               integration.ServiceDiscovery,
		"enabled":                         integration.Enabled,
		"send_collection_latency_metrics": integration.SendCollectionLatencyMetrics,
		"send_database_metrics":           integration.SendDatabaseMetrics,
	}

	// removing optional empty values, terraform complains about unexpected values even though they're empty
	optionals := []string{"api_key", "region", "service_key",
		"team_name", "channel_name", "url", "secret", "password"}

	for _, attr := range optionals {
		if val, ok := out[attr]; ok {
			stringPtr, okT := val.(*string)
			if okT && !conversion.IsStringPresent(stringPtr) {
				delete(out, attr)
			}
		}
	}

	return out
}

func schemaToIntegration(in *schema.ResourceData) (out *admin.ThirdPartyIntegration) {
	out = &admin.ThirdPartyIntegration{}

	if _type, ok := in.GetOk("type"); ok {
		out.Type = admin.PtrString(_type.(string))
	}

	if apiKey, ok := in.GetOk("api_key"); ok {
		out.ApiKey = admin.PtrString(apiKey.(string))
	}

	if region, ok := in.GetOk("region"); ok {
		out.Region = admin.PtrString(region.(string))
	}

	if serviceKey, ok := in.GetOk("service_key"); ok {
		out.ServiceKey = admin.PtrString(serviceKey.(string))
	}

	if teamName, ok := in.GetOk("team_name"); ok {
		out.TeamName = admin.PtrString(teamName.(string))
	}

	if channelName, ok := in.GetOk("channel_name"); ok {
		out.ChannelName = admin.PtrString(channelName.(string))
	}

	if routingKey, ok := in.GetOk("routing_key"); ok {
		out.RoutingKey = admin.PtrString(routingKey.(string))
	}

	if url, ok := in.GetOk("url"); ok {
		out.Url = admin.PtrString(url.(string))
	}

	if secret, ok := in.GetOk("secret"); ok {
		out.Secret = admin.PtrString(secret.(string))
	}

	if microsoftTeamsWebhookURL, ok := in.GetOk("microsoft_teams_webhook_url"); ok {
		out.MicrosoftTeamsWebhookUrl = admin.PtrString(microsoftTeamsWebhookURL.(string))
	}

	if userName, ok := in.GetOk("user_name"); ok {
		out.Username = admin.PtrString(userName.(string))
	}

	if password, ok := in.GetOk("password"); ok {
		out.Password = admin.PtrString(password.(string))
	}

	if serviceDiscovery, ok := in.GetOk("service_discovery"); ok {
		out.ServiceDiscovery = admin.PtrString(serviceDiscovery.(string))
	}

	if enabled, ok := in.GetOk("enabled"); ok {
		out.Enabled = admin.PtrBool(enabled.(bool))
	}

	if sendCollectionLatencyMetrics, ok := in.GetOk("send_collection_latency_metrics"); ok {
		out.SendCollectionLatencyMetrics = admin.PtrBool(sendCollectionLatencyMetrics.(bool))
	}

	if sendDatabaseMetrics, ok := in.GetOk("send_database_metrics"); ok {
		out.SendDatabaseMetrics = admin.PtrBool(sendDatabaseMetrics.(bool))
	}

	return out
}

func updateIntegrationFromSchema(d *schema.ResourceData, integration *admin.ThirdPartyIntegration) {
	integration.ApiKey = conversion.StringPtr(d.Get("api_key").(string))

	if d.HasChange("region") {
		integration.Region = conversion.StringPtr(d.Get("region").(string))
	}

	if d.HasChange("service_key") {
		integration.ServiceKey = conversion.StringPtr(d.Get("service_key").(string))
	}

	if d.HasChange("team_name") {
		integration.TeamName = conversion.StringPtr(d.Get("team_name").(string))
	}

	if d.HasChange("channel_name") {
		integration.ChannelName = conversion.StringPtr(d.Get("channel_name").(string))
	}

	if d.HasChange("routing_key") {
		integration.RoutingKey = conversion.StringPtr(d.Get("routing_key").(string))
	}

	if d.HasChange("url") {
		integration.Url = conversion.StringPtr(d.Get("url").(string))
	}

	if d.HasChange("secret") {
		integration.Secret = conversion.StringPtr(d.Get("secret").(string))
	}

	if d.HasChange("microsoft_teams_webhook_url") {
		integration.MicrosoftTeamsWebhookUrl = conversion.StringPtr(d.Get("microsoft_teams_webhook_url").(string))
	}

	if d.HasChange("user_name") {
		integration.Username = conversion.StringPtr(d.Get("user_name").(string))
	}

	if d.HasChange("password") {
		integration.Password = conversion.StringPtr(d.Get("password").(string))
	}

	if d.HasChange("service_discovery") {
		integration.ServiceDiscovery = conversion.StringPtr(d.Get("service_discovery").(string))
	}

	if d.HasChange("enabled") {
		integration.Enabled = admin.PtrBool(d.Get("enabled").(bool))
	}

	if d.HasChange("send_collection_latency_metrics") {
		integration.SendCollectionLatencyMetrics = admin.PtrBool(d.Get("send_collection_latency_metrics").(bool))
	}

	if d.HasChange("send_database_metrics") {
		integration.SendDatabaseMetrics = admin.PtrBool(d.Get("send_database_metrics").(bool))
	}
}
