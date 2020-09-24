package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

var integrationTypes = []string{
	"PAGER_DUTY",
	"SLACK",
	"DATADOG",
	"NEW_RELIC",
	"OPS_GENIE",
	"VICTOR_OPS",
	"FLOWDOCK",
	"WEBHOOK",
}

var requiredPerType = map[string][]string{
	"PAGER_DUTY": {"service_key"},
	"SLACK":      {"api_token", "team_name"},
	"DATADOG":    {"api_key", "region"},
	"NEW_RELIC":  {"license_key", "account_id", "write_token", "read_token"},
	"OPS_GENIE":  {"api_key", "region"},
	"VICTOR_OPS": {"api_key"},
	"FLOWDOCK":   {"flow_name", "api_token", "org_name"},
	"WEBHOOK":    {"url"},
}

func resourceMongoDBAtlasThirdPartyIntegration() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoDBAtlasThirdPartyIntegrationCreate,
		Read:   resourceMongoDBAtlasThirdPartyIntegrationRead,
		Update: resourceMongoDBAtlasThirdPartyIntegrationUpdate,
		Delete: resourceMongoDBAtlasThirdPartyIntegrationDelete,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(integrationTypes, false),
			},
			"license_key": {
				Type:      schema.TypeString,
				Sensitive: true,
			},
			"account_id": {
				Type: schema.TypeString,
			},
			"write_token": {
				Type:      schema.TypeString,
				Sensitive: true,
			},
			"read_token": {
				Type:      schema.TypeString,
				Sensitive: true,
			},
			"api_key": {
				Type:      schema.TypeString,
				Sensitive: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_key": {
				Type:      schema.TypeString,
				Sensitive: true,
			},
			"api_token": {
				Type:      schema.TypeString,
				Sensitive: true,
			},
			"team_name": {
				Type: schema.TypeString,
			},
			"channel_name": {
				Type: schema.TypeString,
			},
			"routing_key": {
				Type:      schema.TypeString,
				Sensitive: true,
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
			},
		},
	}
}

func resourceMongoDBAtlasThirdPartyIntegrationCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)
	integrationType := d.Get("type").(string)

	// checking per type fields
	if requiredSet, ok := requiredPerType[integrationType]; ok {
		for _, key := range requiredSet {
			_, valid := d.GetOk(key)

			if !valid {
				return fmt.Errorf("error attribute for third party integration %s. please set it", key)
			}
		}
	}

	requestBody := schemaToIntegration(d)

	_, _, err := conn.Integrations.Create(context.Background(), projectID, integrationType, requestBody)

	if err != nil {
		return fmt.Errorf("error creating third party integration %s", err)
	}

	// ID is equal to project_id+type need to ask
	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"type":       integrationType,
	}))

	return resourceMongoDBAtlasThirdPartyIntegrationRead(d, meta)
}

func resourceMongoDBAtlasThirdPartyIntegrationRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())

	projectID := ids["project_id"]
	integrationType := ids["type"]

	integration, _, err := conn.Integrations.Get(context.Background(), projectID, integrationType)

	if err != nil {
		return fmt.Errorf("error getting third party integration resource info %s", integration)
	}

	integrationMap := integrationToSchema(integration)

	for key, val := range integrationMap {
		if err := d.Set(key, val); err != nil {
			return fmt.Errorf("error setting `%s` for third party integration (%s): %s", key, d.Id(), err)
		}
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"type":       integrationType,
	}))

	return nil
}

func resourceMongoDBAtlasThirdPartyIntegrationUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceMongoDBAtlasThirdPartyIntegrationDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
