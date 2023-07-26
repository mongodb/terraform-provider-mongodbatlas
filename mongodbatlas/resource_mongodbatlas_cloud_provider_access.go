package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorCloudProviderAccessCreate   = "error creating cloud provider access %s"
	errorCloudProviderAccessUpdate   = "error updating cloud provider access %s"
	errorCloudProviderAccessDelete   = "error deleting cloud provider access %s"
	errorCloudProviderAccessImporter = "error importing cloud provider access %s"
)

func resourceMongoDBAtlasCloudProviderAccess() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBAtlasCloudProviderAccessCreate,
		ReadContext:   resourceMongoDBAtlasCloudProviderAccessRead,
		UpdateContext: resourceMongoDBAtlasCloudProviderAccessUpdate,
		DeleteContext: resourceMongoDBAtlasCloudProviderAccessDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasCloudProviderAccessImportState,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"provider_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{AWS, AZURE}, false),
			},
			"atlas_aws_account_arn": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"atlas_azure_app_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"service_principal_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"atlas_assumed_role_external_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"authorized_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_update_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"iam_assumed_role_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"role_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"feature_usages": {
				Type:     schema.TypeList,
				Elem:     featureUsagesSchema(),
				Computed: true,
			},
		},
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceMongoDBAtlasCloudProviderAccessV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceMongoDBAtlasCloudProviderAccessStateUpgradeV0,
				Version: 0,
			},
		},
	}
}

func resourceMongoDBAtlasCloudProviderAccessCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectID := d.Get("project_id").(string)

	conn := meta.(*MongoDBClient).Atlas

	req := &matlas.CloudProviderAccessRoleRequest{
		ProviderName: d.Get("provider_name").(string),
	}

	if req.ProviderName == AZURE {
		if value, ok := d.GetOk("atlas_azure_app_id"); ok {
			req.AtlasAzureAppID = pointer(value.(string))
		}

		if value, ok := d.GetOk("service_principal_id"); ok {
			req.AzureServicePrincipalID = pointer(value.(string))
		}

		if value, ok := d.GetOk("tenant_id"); ok {
			req.AzureTenantID = pointer(value.(string))
		}
	}

	role, _, err := conn.CloudProviderAccess.CreateRole(ctx, projectID, req)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorCloudProviderAccessCreate, err))
	}

	roleSchema := roleToSchema(role)
	for key, val := range roleSchema {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(fmt.Errorf(errorCloudProviderAccessCreate, err))
		}
	}

	d.SetId(encodeStateID(map[string]string{
		"id":            role.RoleID,
		"project_id":    projectID,
		"provider_name": role.ProviderName,
	}))

	return nil
}

func resourceMongoDBAtlasCloudProviderAccessRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// sadly there is no just get API
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	roleID := ids["id"]

	role, resp, err := conn.CloudProviderAccess.GetRole(context.Background(), projectID, roleID)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorGetRead, err))
	}

	roleSchema := roleToSchema(role)
	for key, val := range roleSchema {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(fmt.Errorf(errorGetRead, err))
		}
	}

	return nil
}

func resourceMongoDBAtlasCloudProviderAccessUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	if !d.HasChanges("provider_name", "iam_assumed_role_arn", "tenant_id", "service_principal_id", "atlas_azure_app_id") {
		return nil
	}

	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	roleID := ids["id"]

	req := &matlas.CloudProviderAccessRoleRequest{
		ProviderName: d.Get("provider_name").(string),
	}

	if value, ok := d.GetOk("iam_assumed_role_arn"); ok {
		req.IAMAssumedRoleARN = pointer(value.(string))
	}

	if value, ok := d.GetOk("atlas_azure_app_id"); ok {
		req.AtlasAzureAppID = pointer(value.(string))
	}

	if value, ok := d.GetOk("service_principal_id"); ok {
		req.AzureServicePrincipalID = pointer(value.(string))
	}

	if value, ok := d.GetOk("tenant_id"); ok {
		req.AzureTenantID = pointer(value.(string))
	}

	role, _, err := conn.CloudProviderAccess.AuthorizeRole(ctx, projectID, roleID, req)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorCloudProviderAccessUpdate, err))
	}

	roleSchema := roleToSchema(role)
	for key, val := range roleSchema {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(fmt.Errorf(errorGetRead, err))
		}
	}

	return nil
}

func resourceMongoDBAtlasCloudProviderAccessDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())

	projectID := ids["project_id"]
	roleID := ids["id"]
	providerName := ids["provider_name"]

	req := &matlas.CloudProviderDeauthorizationRequest{
		ProviderName: providerName,
		RoleID:       roleID,
		GroupID:      projectID,
	}

	_, err := conn.CloudProviderAccess.DeauthorizeRole(ctx, req)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorCloudProviderAccessDelete, err))
	}

	return nil
}

func resourceMongoDBAtlasCloudProviderAccessImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	projectID, providerName, roleID, err := splitCloudProviderAccessID(d.Id())

	if err != nil {
		return nil, fmt.Errorf(errorCloudProviderAccessImporter, err)
	}

	// searching id in internal format
	d.SetId(encodeStateID(map[string]string{
		"id":            roleID,
		"project_id":    projectID,
		"provider_name": providerName,
	}))

	err2 := resourceMongoDBAtlasCloudProviderAccessRead(ctx, d, meta)
	if err2 != nil {
		return nil, fmt.Errorf(errorCloudProviderAccessImporter, err)
	}

	// case of not found
	if d.Id() == "" {
		return nil, fmt.Errorf(errorCloudProviderAccessImporter, " Resource not found at the cloud please check your id")
	}

	// params syncing
	if err = d.Set("project_id", projectID); err != nil {
		return nil, fmt.Errorf(errorCloudProviderAccessImporter, err)
	}

	return []*schema.ResourceData{d}, nil
}

// format  {project_id}-{provider-name}-{role-id}
func splitCloudProviderAccessID(id string) (projectID, providerName, roleID string, err error) {
	var re = regexp.MustCompile(`(?s)^([0-9a-fA-F]{24})-(.*)-(.*)$`)
	parts := re.FindStringSubmatch(id)

	if len(parts) != 4 {
		err = fmt.Errorf(errorCloudProviderAccessImporter, "format please use {project_id}-{provider-name}-{role-id}")
		return
	}

	projectID, providerName, roleID = parts[1], parts[2], parts[3]

	return
}

func resourceMongoDBAtlasCloudProviderAccessV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"provider_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"AWS"}, false),
			},
			"atlas_aws_account_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"atlas_assumed_role_external_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"authorized_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"iam_assumed_role_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"role_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"feature_usages": {
				Type:     schema.TypeList,
				Elem:     featureUsagesSchemaV0(),
				Computed: true,
			},
		},
	}
}

func resourceMongoDBAtlasCloudProviderAccessStateUpgradeV0(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	rawState["feature_usages"] = []interface{}{map[string]interface{}{}}
	return rawState, nil
}

func roleToSchema(role *matlas.CloudProviderAccessRole) map[string]interface{} {
	out := map[string]interface{}{
		"atlas_aws_account_arn":          role.AtlasAWSAccountARN,
		"atlas_assumed_role_external_id": role.AtlasAssumedRoleExternalID,
		"authorized_date":                role.AuthorizedDate,
		"created_date":                   role.CreatedDate,
		"last_update_date":               role.LastUpdatedDate,
		"iam_assumed_role_arn":           role.IAMAssumedRoleARN,
		"provider_name":                  role.ProviderName,
		"role_id":                        role.RoleID,
		"tenant_id":                      role.AzureTenantID,
		"service_principal_id":           role.AzureServicePrincipalID,
		"atlas_azure_app_id":             role.AtlasAzureAppID,
	}

	if role.ProviderName == AZURE {
		out["role_id"] = role.AzureID
	}

	features := make([]map[string]interface{}, 0, len(role.FeatureUsages))
	for _, featureUsage := range role.FeatureUsages {
		features = append(features, featureToSchema(featureUsage))
	}
	out["feature_usages"] = features
	return out
}
