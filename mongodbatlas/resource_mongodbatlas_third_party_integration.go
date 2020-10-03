package mongodbatlas

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

var integrationTypes = []string{
	"PAGER_DUTY",
	"DATADOG",
	"NEW_RELIC",
	"OPS_GENIE",
	"VICTOR_OPS",
	"FLOWDOCK",
	"WEBHOOK",
}

var requiredPerType = map[string][]string{
	"PAGER_DUTY": {"service_key"},
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
		Importer: &schema.ResourceImporter{
			State: resourceMongoDBAtlasThirdPartyIntegrationImportState,
		},
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
				Optional:  true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"write_token": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
			},
			"read_token": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
			},
			"api_key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
			},
			"api_token": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
			},
			"team_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"channel_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"routing_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"flow_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"org_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"secret": {
				Type:      schema.TypeString,
				Optional:  true,
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
		return fmt.Errorf("error getting third party integration resource info %s %w", integrationType, err)
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
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())

	projectID := ids["project_id"]
	integrationType := ids["type"]

	integration, _, err := conn.Integrations.Get(context.Background(), projectID, integrationType)

	if err != nil {
		return fmt.Errorf("error getting third party integration resource info %s %w", integrationType, err)
	}

	// check for changed attributes per type

	updateIntegrationFromSchema(d, integration)

	_, _, err = conn.Integrations.Replace(context.Background(), projectID, integrationType, integration)

	if err != nil {
		return fmt.Errorf("error updating third party integration type `%s` (%s): %w", integrationType, d.Id(), err)
	}

	return resourceMongoDBAtlasThirdPartyIntegrationRead(d, meta)
}

func resourceMongoDBAtlasThirdPartyIntegrationDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())

	projectID := ids["project_id"]
	integrationType := ids["type"]

	_, err := conn.Integrations.Delete(context.Background(), projectID, integrationType)

	if err != nil {
		return fmt.Errorf("error deleting third party integration type `%s` (%s): %w", integrationType, d.Id(), err)
	}

	return nil
}

func resourceMongoDBAtlasThirdPartyIntegrationImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*matlas.Client)

	projectID, integrationType, err := splitIntegrationTypeID(d.Id())

	if err != nil {
		return nil, err
	}

	integration, _, err := conn.Integrations.Get(context.Background(), projectID, integrationType)

	if err != nil {
		return nil, fmt.Errorf("couldn't import third party integration (%s) in project(%s), error: %w", integrationType, projectID, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		return nil, fmt.Errorf("error setting `project_id` for third party integration (%s): %w", d.Id(), err)
	}

	if err := d.Set("type", integration.Type); err != nil {
		return nil, fmt.Errorf("error setting `type` for third party integration (%s): %w", d.Id(), err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"type":       integrationType,
	}))

	return []*schema.ResourceData{d}, nil
}

// format {project_id}-{integration_type}
func splitIntegrationTypeID(id string) (projectID, integrationType string, err error) {
	var re = regexp.MustCompile(`(?s)^([0-9a-fA-F]{24})-(.*)$`)
	parts := re.FindStringSubmatch(id)

	if len(parts) != 3 {
		err = fmt.Errorf("import format error: to import a third party integration, use the format {project_id}-{integration_type} %s, %+v", id, parts)
		return
	}

	projectID, integrationType = parts[1], parts[2]

	return
}
