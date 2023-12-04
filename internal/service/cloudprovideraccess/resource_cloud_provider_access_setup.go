package cloudprovideraccess

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

/*
	mongodb_atlas_cloud_provider_access_setup
	-> Creates the the information from the mongodbatlas side
	-> The delete deletes and deauthorize the role
*/

func ResourceSetup() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceMongoDBAtlasCloudProviderAccessSetupRead,
		CreateContext: resourceMongoDBAtlasCloudProviderAccessSetupCreate,
		UpdateContext: resourceMongoDBAtlasCloudProviderAccessAuthorizationPlaceHolder,
		DeleteContext: resourceMongoDBAtlasCloudProviderAccessSetupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasCloudProviderAccessSetupImportState,
		},

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"provider_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{constant.AWS, constant.AZURE}, false),
				ForceNew:     true,
			},
			"aws_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"atlas_aws_account_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"atlas_assumed_role_external_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"azure_config": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"atlas_azure_app_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"service_principal_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"created_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_updated_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceMongoDBAtlasCloudProviderAccessSetupRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).Atlas
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	roleID := ids["id"]

	role, resp, err := conn.CloudProviderAccess.GetRole(context.Background(), projectID, roleID)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf(ErrorCloudProviderGetRead, err))
	}

	roleSchema := roleToSchemaSetup(role)
	for key, val := range roleSchema {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(fmt.Errorf(ErrorCloudProviderGetRead, err))
		}
	}

	return nil
}

func resourceMongoDBAtlasCloudProviderAccessSetupCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	projectID := d.Get("project_id").(string)

	conn := meta.(*config.MongoDBClient).Atlas

	requestParameters := &matlas.CloudProviderAccessRoleRequest{
		ProviderName: d.Get("provider_name").(string),
	}

	if value, ok := d.GetOk("azure_config.0.atlas_azure_app_id"); ok {
		requestParameters.AtlasAzureAppID = conversion.Pointer(value.(string))
	}

	if value, ok := d.GetOk("azure_config.0.service_principal_id"); ok {
		requestParameters.AzureServicePrincipalID = conversion.Pointer(value.(string))
	}

	if value, ok := d.GetOk("azure_config.0.tenant_id"); ok {
		requestParameters.AzureTenantID = conversion.Pointer(value.(string))
	}

	if value, ok := d.GetOk("azure_config.0.atlas_azure_app_id"); ok {
		requestParameters.AtlasAzureAppID = conversion.Pointer(value.(string))
	}

	if value, ok := d.GetOk("azure_config.0.service_principal_id"); ok {
		requestParameters.AzureServicePrincipalID = conversion.Pointer(value.(string))
	}

	if value, ok := d.GetOk("azure_config.0.tenant_id"); ok {
		requestParameters.AzureTenantID = conversion.Pointer(value.(string))
	}

	role, _, err := conn.CloudProviderAccess.CreateRole(ctx, projectID, requestParameters)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorCloudProviderAccessCreate, err))
	}

	// once multiple providers enable here do a switch, select for provider type
	roleSchema := roleToSchemaSetup(role)

	resourceID := role.RoleID
	if role.ProviderName == constant.AZURE {
		resourceID = *role.AzureID
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"id":            resourceID,
		"project_id":    projectID,
		"provider_name": role.ProviderName,
	}))

	for key, val := range roleSchema {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(fmt.Errorf(errorCloudProviderAccessCreate, err))
		}
	}

	return nil
}

func resourceMongoDBAtlasCloudProviderAccessSetupDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).Atlas
	ids := conversion.DecodeStateID(d.Id())

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

	d.SetId("")
	d.SetId("")
	return nil
}

func roleToSchemaSetup(role *matlas.CloudProviderAccessRole) map[string]any {
	if role.ProviderName == "AWS" {
		out := map[string]any{
			"provider_name": role.ProviderName,
			"aws_config": []any{map[string]any{
				"atlas_aws_account_arn":          role.AtlasAWSAccountARN,
				"atlas_assumed_role_external_id": role.AtlasAssumedRoleExternalID,
			}},
			"created_date": role.CreatedDate,
			"role_id":      role.RoleID,
		}
		return out
	}

	out := map[string]any{
		"provider_name": role.ProviderName,
		"azure_config": []any{map[string]any{
			"atlas_azure_app_id":   role.AtlasAzureAppID,
			"service_principal_id": role.AzureServicePrincipalID,
			"tenant_id":            role.AzureTenantID,
		}},
		"aws_config":        []any{map[string]any{}},
		"created_date":      role.CreatedDate,
		"last_updated_date": role.LastUpdatedDate,
		"role_id":           role.AzureID,
	}

	return out
}

func resourceMongoDBAtlasCloudProviderAccessSetupImportState(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	projectID, providerName, roleID, err := splitCloudProviderAccessID(d.Id())

	if err != nil {
		return nil, fmt.Errorf(errorCloudProviderAccessImporter, err)
	}

	// searching id in internal format
	d.SetId(conversion.EncodeStateID(map[string]string{
		"id":            roleID,
		"project_id":    projectID,
		"provider_name": providerName,
	}))

	err2 := resourceMongoDBAtlasCloudProviderAccessSetupRead(ctx, d, meta)

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
