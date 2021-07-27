package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasThirdPartyIntegrations() *schema.Resource {
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

func dataSourceMongoDBAtlasThirdPartyIntegrationsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	projectID := d.Get("project_id").(string)
	integrations, _, err := conn.Integrations.List(ctx, projectID)

	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting third party integration list: %s", err))
	}

	if err = d.Set("results", flattenIntegrations(integrations, projectID)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting results for third party integrations %s", err))
	}

	d.SetId(resource.UniqueId())

	return nil
}

func flattenIntegrations(integrations *matlas.ThirdPartyIntegrations, projectID string) (list []map[string]interface{}) {
	if len(integrations.Results) == 0 {
		return
	}

	list = make([]map[string]interface{}, 0, len(integrations.Results))

	for _, integration := range integrations.Results {
		service := integrationToSchema(integration)
		service["project_id"] = projectID
		list = append(list, service)
	}

	return
}

func integrationToSchema(integration *matlas.ThirdPartyIntegration) map[string]interface{} {
	out := map[string]interface{}{
		"type":         integration.Type,
		"license_key":  integration.LicenseKey,
		"account_id":   integration.AccountID,
		"write_token":  integration.WriteToken,
		"read_token":   integration.ReadToken,
		"api_key":      integration.APIKey,
		"region":       integration.Region,
		"service_key":  integration.ServiceKey,
		"api_token":    integration.APIToken,
		"team_name":    integration.TeamName,
		"channel_name": integration.ChannelName,
		"routing_key":  integration.RoutingKey,
		"flow_name":    integration.FlowName,
		"org_name":     integration.OrgName,
		"url":          integration.URL,
		"secret":       integration.Secret,
	}

	// removing optional empty values, terraform complains about unexpected values even though they're empty
	optionals := []string{"license_key", "account_id", "write_token",
		"read_token", "api_key", "region", "service_key", "api_token",
		"team_name", "channel_name", "flow_name", "org_name", "url", "secret"}

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

	if licenseKey, ok := in.GetOk("license_key"); ok {
		out.LicenseKey = licenseKey.(string)
	}

	if accountID, ok := in.GetOk("account_id"); ok {
		out.AccountID = accountID.(string)
	}

	if writeToken, ok := in.GetOk("write_token"); ok {
		out.WriteToken = writeToken.(string)
	}

	if readToken, ok := in.GetOk("read_token"); ok {
		out.ReadToken = readToken.(string)
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

	if apiToken, ok := in.GetOk("api_token"); ok {
		out.APIToken = apiToken.(string)
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

	if flowName, ok := in.GetOk("flow_name"); ok {
		out.FlowName = flowName.(string)
	}

	if orgName, ok := in.GetOk("org_name"); ok {
		out.OrgName = orgName.(string)
	}

	if url, ok := in.GetOk("url"); ok {
		out.URL = url.(string)
	}

	if secret, ok := in.GetOk("secret"); ok {
		out.Secret = secret.(string)
	}

	return out
}

func updateIntegrationFromSchema(d *schema.ResourceData, integration *matlas.ThirdPartyIntegration) {
	if d.HasChange("license_key") {
		integration.LicenseKey = d.Get("license_key").(string)
	}

	if d.HasChange("account_id") {
		integration.AccountID = d.Get("account_id").(string)
	}

	if d.HasChange("write_token") {
		integration.WriteToken = d.Get("write_token").(string)
	}

	if d.HasChange("read_token") {
		integration.ReadToken = d.Get("read_token").(string)
	}

	if d.HasChange("api_key") {
		integration.APIKey = d.Get("api_key").(string)
	}

	if d.HasChange("region") {
		integration.Region = d.Get("region").(string)
	}

	if d.HasChange("service_key") {
		integration.ServiceKey = d.Get("service_key").(string)
	}

	if d.HasChange("api_token") {
		integration.APIToken = d.Get("api_token").(string)
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

	if d.HasChange("flow_name") {
		integration.FlowName = d.Get("flow_name").(string)
	}

	if d.HasChange("org_name") {
		integration.OrgName = d.Get("org_name").(string)
	}

	if d.HasChange("url") {
		integration.URL = d.Get("url").(string)
	}

	if d.HasChange("secret") {
		integration.Secret = d.Get("secret").(string)
	}
}
