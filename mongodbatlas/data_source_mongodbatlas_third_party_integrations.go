package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasThirdPartyIntegrations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasThirdPartyIntegrationsRead,
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

func dataSourceMongoDBAtlasThirdPartyIntegrationsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	projectID := d.Get("project_id").(string)
	integrations, _, err := conn.Integrations.List(context.Background(), projectID)

	if err != nil {
		return fmt.Errorf("error getting third party integration list: %s", err)
	}

	if err = d.Set("results", flattenIntegrations(integrations)); err != nil {
		return fmt.Errorf("error setting results for third party integrations %s", err)
	}

	d.SetId(resource.UniqueId())

	return nil
}

func flattenIntegrations(integrations *matlas.ThirdPartyIntegrations) (list []map[string]interface{}) {
	if len(integrations.Results) == 0 {
		return
	}

	list = make([]map[string]interface{}, 0, len(integrations.Results))

	for _, integration := range integrations.Results {
		list = append(list, integrationToSchema(integration))
	}

	return
}

func integrationToSchema(integration *matlas.ThirdPartyIntegration) map[string]interface{} {
	return map[string]interface{}{
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
