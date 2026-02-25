package organization

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strings"

	"go.mongodb.org/atlas-sdk/v20250312014/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var (
	attrsCreateRequired    = []string{"org_owner_id"}              // name not included as it's already required in the schema.
	attrsCreateRequiredPAK = []string{"description", "role_names"} // only required when creating a PAK (no service_account block).
	attrsCreateOnly        = []string{"org_owner_id", "description", "role_names", "federation_settings_id", "service_account"}
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate,
		ReadContext:   resourceRead,
		UpdateContext: resourceUpdate,
		DeleteContext: resourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceImport,
		},
		Schema: map[string]*schema.Schema{
			"org_owner_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"org_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"private_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"role_names": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"federation_settings_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"api_access_list_required": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"multi_factor_auth_required": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"restrict_employee_access": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"gen_ai_features_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			// skip_default_alerts_settings defaults to `true` to prevent Atlas from automatically creating organization-level alerts not explicitly managed through Terraform.
			// Note that this deviates from the API default of `false` for this attribute.
			"skip_default_alerts_settings": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"security_contact": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_account": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Required: true,
						},
						"roles": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"secret_expires_after_hours": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"client_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"secrets": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"created_at": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"expires_at": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"secret_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"secret": {
										Type:      schema.TypeString,
										Computed:  true,
										Sensitive: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	for _, attr := range attrsCreateRequired {
		if _, ok := d.GetOk(attr); !ok {
			return diag.FromErr(fmt.Errorf("%s is required during organization creation", attr))
		}
	}
	_, usingSA := d.GetOk("service_account")
	if !usingSA {
		for _, attr := range attrsCreateRequiredPAK {
			if _, ok := d.GetOk(attr); !ok {
				return diag.FromErr(fmt.Errorf("%s is required during organization creation when not using service_account", attr))
			}
		}
		if err := ValidateAPIKeyIsOrgOwner(conversion.ExpandStringList(d.Get("role_names").(*schema.Set).List())); err != nil {
			return diag.FromErr(err)
		}
	}
	conn := getAtlasV2Connection(ctx, d, meta) // Using provider credentials.
	organization, resp, err := conn.OrganizationsApi.CreateOrg(ctx, newCreateOrganizationRequest(d)).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) && !strings.Contains(err.Error(), "USER_NOT_FOUND") {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error creating Organization: %s", err))
	}
	if usingSA {
		sa, saOk := organization.GetServiceAccountOk()
		if !saOk {
			return diag.FromErr(fmt.Errorf("service account was not returned by the API"))
		}
		if err := setServiceAccountState(d, sa); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := d.Set("private_key", organization.ApiKey.GetPrivateKey()); err != nil {
			return diag.FromErr(fmt.Errorf("error setting `private_key`: %s", err))
		}
		if err := d.Set("public_key", organization.ApiKey.GetPublicKey()); err != nil {
			return diag.FromErr(fmt.Errorf("error setting `public_key`: %s", err))
		}
	}
	conn = getAtlasV2Connection(ctx, d, meta) // Using new credentials from the created organization.
	orgID := organization.Organization.GetId()
	_, _, errUpdate := conn.OrganizationsApi.UpdateOrgSettings(ctx, orgID, newOrganizationSettings(d)).Execute()
	if errUpdate != nil {
		if _, err := conn.OrganizationsApi.DeleteOrg(ctx, orgID).Execute(); err != nil {
			d.SetId("")
			return diag.FromErr(fmt.Errorf("an error occurred when updating Organization settings: %s.\n Unable to delete organization, there may be dangling resources: %s", errUpdate.Error(), err.Error()))
		}
		d.SetId("")
		return diag.FromErr(fmt.Errorf("an error occurred when updating Organization settings: %s", err))
	}
	d.SetId(conversion.EncodeStateID(map[string]string{
		"org_id": orgID,
	}))
	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := getAtlasV2Connection(ctx, d, meta)
	ids := conversion.DecodeStateID(d.Id())
	orgID := ids["org_id"]

	organization, resp, err := conn.OrganizationsApi.GetOrg(ctx, orgID).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) {
			log.Printf("warning Organization deleted will recreate: %s \n", err.Error())
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error reading organization information: %s", err))
	}

	if err := d.Set("name", organization.Name); err != nil {
		return diag.Errorf("error setting `name` for organization (%s): %s", *organization.Id, err)
	}
	if err := d.Set("skip_default_alerts_settings", organization.SkipDefaultAlertsSettings); err != nil {
		return diag.Errorf("error setting `skip_default_alerts_settings` for organization (%s): %s", orgID, err)
	}
	if err := d.Set("org_id", orgID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `org_id`: %s", err))
	}

	settings, _, err := conn.OrganizationsApi.GetOrgSettings(ctx, orgID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading organization settings: %s", err))
	}

	if err := d.Set("api_access_list_required", settings.ApiAccessListRequired); err != nil {
		return diag.Errorf("error setting `api_access_list_required` for organization (%s): %s", orgID, err)
	}
	if err := d.Set("multi_factor_auth_required", settings.MultiFactorAuthRequired); err != nil {
		return diag.Errorf("error setting `multi_factor_auth_required` for organization (%s): %s", orgID, err)
	}
	if err := d.Set("restrict_employee_access", settings.RestrictEmployeeAccess); err != nil {
		return diag.Errorf("error setting `restrict_employee_access` for organization (%s): %s", orgID, err)
	}
	if err := d.Set("gen_ai_features_enabled", settings.GenAIFeaturesEnabled); err != nil {
		return diag.Errorf("error setting `gen_ai_features_enabled` for organization (%s): %s", orgID, err)
	}
	if err := d.Set("security_contact", settings.SecurityContact); err != nil {
		return diag.Errorf("error setting `security_contact` for organization (%s): %s", orgID, err)
	}
	return nil
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := getAtlasV2Connection(ctx, d, meta)
	ids := conversion.DecodeStateID(d.Id())
	orgID := ids["org_id"]
	for _, attr := range attrsCreateOnly {
		if d.HasChange(attr) {
			return diag.Errorf("%s cannot be changed after creation", attr)
		}
	}
	if d.HasChange("name") ||
		d.HasChange("skip_default_alerts_settings") {
		updateRequest := &admin.AtlasOrganization{
			Name:                      d.Get("name").(string),
			SkipDefaultAlertsSettings: new(d.Get("skip_default_alerts_settings").(bool)),
		}
		if _, _, err := conn.OrganizationsApi.UpdateOrg(ctx, orgID, updateRequest).Execute(); err != nil {
			return diag.FromErr(fmt.Errorf("error updating Organization name: %s", err))
		}
	}

	if d.HasChange("api_access_list_required") ||
		d.HasChange("multi_factor_auth_required") ||
		d.HasChange("restrict_employee_access") ||
		d.HasChange("gen_ai_features_enabled") ||
		d.HasChange("security_contact") {
		if _, _, err := conn.OrganizationsApi.UpdateOrgSettings(ctx, orgID, newOrganizationSettings(d)).Execute(); err != nil {
			return diag.FromErr(fmt.Errorf("error updating Organization settings: %s", err))
		}
	}

	return resourceRead(ctx, d, meta)
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := getAtlasV2Connection(ctx, d, meta)
	ids := conversion.DecodeStateID(d.Id())
	orgID := ids["org_id"]

	if _, err := conn.OrganizationsApi.DeleteOrg(ctx, orgID).Execute(); err != nil {
		return diag.FromErr(fmt.Errorf("error deleting Organization: %s", err))
	}
	return nil
}

func resourceImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	d.SetId(conversion.EncodeStateID(map[string]string{
		"org_id": d.Id(),
	}))
	return []*schema.ResourceData{d}, nil
}

func newCreateOrganizationRequest(d *schema.ResourceData) *admin.CreateOrganizationRequest {
	// skip_default_alerts_settings defaults to `true` to prevent Atlas from automatically creating organization-level alerts not explicitly managed through Terraform.
	// Note that this deviates from the API default of `false` for this attribute.
	skipDefaultAlertsSettings := true

	if v, ok := d.GetOkExists("skip_default_alerts_settings"); ok {
		skipDefaultAlertsSettings = v.(bool)
	}

	createRequest := &admin.CreateOrganizationRequest{
		Name:                      d.Get("name").(string),
		OrgOwnerId:                new(d.Get("org_owner_id").(string)),
		SkipDefaultAlertsSettings: new(skipDefaultAlertsSettings),
	}

	if federationSettingsID, ok := d.Get("federation_settings_id").(string); ok && federationSettingsID != "" {
		createRequest.FederationSettingsId = &federationSettingsID
	}

	// API does not allow both apiKey and serviceAccount in the same request.
	if v, ok := d.GetOk("service_account"); ok {
		saList := v.([]any)
		if len(saList) > 0 {
			saMap := saList[0].(map[string]any)
			createRequest.ServiceAccount = &admin.OrgServiceAccountRequest{
				Name:                    saMap["name"].(string),
				Description:             saMap["description"].(string),
				Roles:                   conversion.ExpandStringList(saMap["roles"].(*schema.Set).List()),
				SecretExpiresAfterHours: saMap["secret_expires_after_hours"].(int),
			}
		}
	} else {
		createRequest.ApiKey = &admin.CreateAtlasOrganizationApiKey{
			Roles: conversion.ExpandStringList(d.Get("role_names").(*schema.Set).List()),
			Desc:  d.Get("description").(string),
		}
	}

	return createRequest
}

func newOrganizationSettings(d *schema.ResourceData) *admin.OrganizationSettings {
	return &admin.OrganizationSettings{
		ApiAccessListRequired:   new(d.Get("api_access_list_required").(bool)),
		MultiFactorAuthRequired: new(d.Get("multi_factor_auth_required").(bool)),
		RestrictEmployeeAccess:  new(d.Get("restrict_employee_access").(bool)),
		GenAIFeaturesEnabled:    new(d.Get("gen_ai_features_enabled").(bool)),
		SecurityContact:         new(d.Get("security_contact").(string)),
	}
}

func ValidateAPIKeyIsOrgOwner(roles []string) error {
	if slices.Contains(roles, constant.OrgOwner) {
		return nil
	}

	return fmt.Errorf("`role_names` for new API Key must have the ORG_OWNER role to use this resource")
}

// setServiceAccountState merges the SA API response into the existing service_account block in state.
func setServiceAccountState(d *schema.ResourceData, sa *admin.OrgServiceAccount) error {
	// Preserve user-configured input values from the existing block, only adding computed outputs.
	existing := d.Get("service_account").([]any)
	if len(existing) == 0 {
		return fmt.Errorf("service account was returned by the API but service_account block is not configured")
	}
	saMap := existing[0].(map[string]any)
	saMap["client_id"] = sa.GetClientId()
	saMap["created_at"] = sa.GetCreatedAt().String()

	var secretsList []map[string]any
	for _, s := range sa.GetSecrets() {
		secretsList = append(secretsList, map[string]any{
			"created_at": s.GetCreatedAt().String(),
			"expires_at": s.GetExpiresAt().String(),
			"secret_id":  s.GetId(),
			"secret":     s.GetSecret(),
		})
	}
	saMap["secrets"] = secretsList
	return d.Set("service_account", []any{saMap})
}

// getAtlasV2Connection uses the created credentials for the organization if they exist.
// It tries PAK credentials first, then SA credentials from the service_account block,
// and falls back to provider credentials (e.g. if the resource was imported).
// For SA credentials, if they are present but no longer valid (e.g. expired secret or
// insufficient access), the function falls back to provider credentials.
func getAtlasV2Connection(ctx context.Context, d *schema.ResourceData, meta any) *admin.APIClient {
	currentClient := meta.(*config.MongoDBClient)

	// Try PAK credentials
	publicKey := d.Get("public_key").(string)
	privateKey := d.Get("private_key").(string)
	if publicKey != "" && privateKey != "" {
		c := &config.Credentials{
			PublicKey:  publicKey,
			PrivateKey: privateKey,
			BaseURL:    currentClient.BaseURL,
		}
		if newClient, err := config.NewClient(c, currentClient.TerraformVersion); err == nil {
			return newClient.AtlasV2
		}
	}

	// Try SA credentials from service_account block.
	// Falls back to provider credentials if SA creds are no longer valid.
	if v, ok := d.GetOk("service_account"); ok {
		saList := v.([]any)
		if len(saList) > 0 {
			saMap := saList[0].(map[string]any)
			clientID, _ := saMap["client_id"].(string)
			secretValue := ""
			if secretsList, ok := saMap["secrets"].([]any); ok && len(secretsList) > 0 {
				if secretMap, ok := secretsList[0].(map[string]any); ok {
					secretValue, _ = secretMap["secret"].(string)
				}
			}
			if clientID != "" && secretValue != "" {
				if saClient := newSAClient(ctx, d, clientID, secretValue, currentClient); saClient != nil {
					return saClient
				}
			}
		}
	}

	return currentClient.AtlasV2
}

// newSAClient creates an API client using SA credentials and verifies access to the organization
// with an explicit GetOrg call. Returns nil if the credentials are invalid or lack access.
func newSAClient(ctx context.Context, d *schema.ResourceData, clientID, secretValue string, currentClient *config.MongoDBClient) *admin.APIClient {
	c := &config.Credentials{
		ClientID:     clientID,
		ClientSecret: secretValue,
		BaseURL:      currentClient.BaseURL,
	}
	newClient, err := config.NewClient(c, currentClient.TerraformVersion)
	if err != nil {
		return nil
	}
	if d.Id() != "" {
		ids := conversion.DecodeStateID(d.Id())
		if orgID := ids["org_id"]; orgID != "" {
			if _, _, err := newClient.AtlasV2.OrganizationsApi.GetOrg(ctx, orgID).Execute(); err != nil {
				return nil
			}
		}
	}
	return newClient.AtlasV2
}
