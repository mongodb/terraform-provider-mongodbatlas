package thirdpartyintegration

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var integrationTypes = []string{
	"PAGER_DUTY",
	"DATADOG",
	"OPS_GENIE",
	"VICTOR_OPS",
	"WEBHOOK",
	"MICROSOFT_TEAMS",
	"PROMETHEUS",
}

var requiredPerType = map[string][]string{
	"PAGER_DUTY":      {"service_key"},
	"DATADOG":         {"api_key", "region"},
	"OPS_GENIE":       {"api_key", "region"},
	"VICTOR_OPS":      {"api_key"},
	"WEBHOOK":         {"url"},
	"MICROSOFT_TEAMS": {"microsoft_teams_webhook_url"},
	"PROMETHEUS":      {"user_name", "password", "service_discovery", "enabled"},
}

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate,
		ReadContext:   resourceRead,
		UpdateContext: resourceUpdate,
		DeleteContext: resourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceImportState,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateIntegrationType(),
			},
			"api_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"service_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"team_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"channel_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"routing_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"url": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"microsoft_teams_webhook_url": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"user_name": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"service_discovery": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"send_collection_latency_metrics": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"send_database_metrics": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"send_user_provided_resource_tags": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)
	integrationType := d.Get("type").(string)

	// checking per type fields
	if requiredSet, ok := requiredPerType[integrationType]; ok {
		for _, key := range requiredSet {
			_, valid := d.GetOk(key)

			if !valid {
				return diag.FromErr(fmt.Errorf("error attribute for third party integration %s. please set it", key))
			}
		}
	}

	requestBody := schemaToIntegration(d)

	_, _, err := connV2.ThirdPartyIntegrationsApi.CreateThirdPartyIntegration(ctx, integrationType, projectID, requestBody).Execute()

	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating third party integration %s", err))
	}

	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	projectID := d.Get("project_id").(string)
	integrationType := d.Get("type").(string)

	integration, resp, err := connV2.ThirdPartyIntegrationsApi.GetThirdPartyIntegration(ctx, projectID, integrationType).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf("error getting third party integration resource info %s %w", integrationType, err))
	}

	integrationMap := integrationToSchema(d, integration)

	for key, val := range integrationMap {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(fmt.Errorf("error setting `%s` for third party integration (%s): %s", key, d.Id(), err))
		}
	}

	d.SetId(integration.GetId())
	return nil
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	projectID := d.Get("project_id").(string)
	integrationType := d.Get("type").(string)

	integration, _, err := connV2.ThirdPartyIntegrationsApi.GetThirdPartyIntegration(ctx, projectID, integrationType).Execute()

	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting third party integration resource info %s %w", integrationType, err))
	}

	// check for changed attributes per type

	updateIntegrationFromSchema(d, integration)

	_, _, err = connV2.ThirdPartyIntegrationsApi.UpdateThirdPartyIntegration(ctx, integrationType, projectID, integration).Execute()

	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating third party integration type `%s` (%s): %w", integrationType, d.Id(), err))
	}

	return resourceRead(ctx, d, meta)
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).Atlas

	projectID := d.Get("project_id").(string)
	integrationType := d.Get("type").(string)

	_, err := conn.Integrations.Delete(ctx, projectID, integrationType)

	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting third party integration type `%s` (%s): %w", integrationType, d.Id(), err))
	}

	return nil
}

func resourceImportState(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	projectID, integrationType, err := splitIntegrationTypeID(d.Id())

	if err != nil {
		return nil, err
	}

	_, _, err = connV2.ThirdPartyIntegrationsApi.GetThirdPartyIntegration(ctx, projectID, integrationType).Execute()

	if err != nil {
		return nil, fmt.Errorf("couldn't import third party integration (%s) in project(%s), error: %w", integrationType, projectID, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		return nil, fmt.Errorf("error setting `project_id` for third party integration (%s): %w", d.Id(), err)
	}

	if err := d.Set("type", integrationType); err != nil {
		return nil, fmt.Errorf("error setting `type` for third party integration (%s): %w", d.Id(), err)
	}

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

func validateIntegrationType() schema.SchemaValidateDiagFunc {
	return func(v any, p cty.Path) diag.Diagnostics {
		value := v.(string)
		var diags diag.Diagnostics
		if !isElementExist(integrationTypes, value) {
			diagError := diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid Third Party Integration type",
				Detail:   fmt.Sprintf("Third Party integration type %q is not a valid value. Possible values are: %q.", value, integrationTypes),
			}
			diags = append(diags, diagError)
		}
		return diags
	}
}

func isElementExist(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
