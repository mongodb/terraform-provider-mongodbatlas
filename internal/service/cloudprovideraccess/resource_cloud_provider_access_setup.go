package cloudprovideraccess

import (
	"context"
	"fmt"
	"regexp"

	"go.mongodb.org/atlas-sdk/v20250312006/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

/*
	mongodb_atlas_cloud_provider_access_setup
	-> Creates the the information from the mongodbatlas side
	-> The delete deletes and deauthorize the role
*/

const (
	errorCloudProviderAccessCreate   = "error creating cloud provider access %s"
	errorCloudProviderAccessUpdate   = "error updating cloud provider access %s"
	errorCloudProviderAccessDelete   = "error deleting cloud provider access %s"
	errorCloudProviderAccessImporter = "error importing cloud provider access %s"
	ErrorCloudProviderGetRead        = "error reading cloud provider access %s"
)

func ResourceSetup() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceCloudProviderAccessSetupRead,
		CreateContext: resourceCloudProviderAccessSetupCreate,
		UpdateContext: resourceCloudProviderAccessAuthorizationPlaceHolder,
		DeleteContext: resourceCloudProviderAccessSetupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceCloudProviderAccessSetupImportState,
		},

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"provider_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"gcp_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_account_for_atlas": {
							Type:     schema.TypeString,
							Computed: true,
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

func resourceCloudProviderAccessSetupRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	roleID := ids["id"]

	role, resp, err := conn.CloudProviderAccessApi.GetCloudProviderAccessRole(context.Background(), projectID, roleID).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) {
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

func resourceCloudProviderAccessSetupCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	projectID := d.Get("project_id").(string)

	conn := meta.(*config.MongoDBClient).AtlasV2

	requestParameters := &admin.CloudProviderAccessRoleRequest{
		ProviderName: d.Get("provider_name").(string),
	}

	if value, ok := d.GetOk("azure_config.0.atlas_azure_app_id"); ok {
		requestParameters.SetAtlasAzureAppId(value.(string))
	}

	if value, ok := d.GetOk("azure_config.0.service_principal_id"); ok {
		requestParameters.SetServicePrincipalId(value.(string))
	}

	if value, ok := d.GetOk("azure_config.0.tenant_id"); ok {
		requestParameters.SetTenantId(value.(string))
	}

	role, _, err := conn.CloudProviderAccessApi.CreateCloudProviderAccessRole(ctx, projectID, requestParameters).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorCloudProviderAccessCreate, err))
	}

	// once multiple providers enable here do a switch, select for provider type
	roleSchema := roleToSchemaSetup(role)

	resourceID := role.GetRoleId()
	if role.ProviderName == constant.AZURE {
		resourceID = role.GetId()
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"id":            resourceID,
		"project_id":    projectID,
		"provider_name": role.GetProviderName(),
	}))

	for key, val := range roleSchema {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(fmt.Errorf(errorCloudProviderAccessCreate, err))
		}
	}

	return nil
}

func resourceCloudProviderAccessSetupDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())

	projectID := ids["project_id"]
	roleID := ids["id"]
	providerName := ids["provider_name"]

	req := &admin.DeauthorizeCloudProviderAccessRoleApiParams{
		CloudProvider: providerName,
		RoleId:        roleID,
		GroupId:       projectID,
	}

	_, err := conn.CloudProviderAccessApi.DeauthorizeCloudProviderAccessRoleWithParams(ctx, req).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorCloudProviderAccessDelete, err))
	}

	d.SetId("")
	d.SetId("")
	return nil
}

func roleToSchemaSetup(role *admin.CloudProviderAccessRole) map[string]any {
	out := map[string]any{}

	if role.ProviderName == "AWS" {
		out["provider_name"] = role.GetProviderName()
		out["aws_config"] = []any{map[string]any{
			"atlas_aws_account_arn":          role.GetAtlasAWSAccountArn(),
			"atlas_assumed_role_external_id": role.GetAtlasAssumedRoleExternalId(),
		}}
		out["gcp_config"] = []any{map[string]any{}}
		out["created_date"] = conversion.TimeToString(role.GetCreatedDate())
		out["role_id"] = role.GetRoleId()
	} else if role.ProviderName == "AZURE" {
		out["provider_name"] = role.GetProviderName()
		out["azure_config"] = []any{map[string]any{
			"atlas_azure_app_id":   role.GetAtlasAzureAppId(),
			"service_principal_id": role.GetServicePrincipalId(),
			"tenant_id":            role.GetTenantId(),
		}}
		out["aws_config"] = []any{map[string]any{}}
		out["gcp_config"] = []any{map[string]any{}}
		out["created_date"] = conversion.TimeToString(role.GetCreatedDate())
		out["last_updated_date"] = conversion.TimeToString(role.GetLastUpdatedDate())
		out["role_id"] = role.GetId()
	} else if role.ProviderName == "GCP" {
		out["provider_name"] = role.GetProviderName()
		out["gcp_config"] = []any{map[string]any{
			"status":                    role.GetStatus(),
			"service_account_for_atlas": role.GetGcpServiceAccountForAtlas(),
		}}
		out["role_id"] = role.GetId()
		out["aws_config"] = []any{map[string]any{}}
		out["created_date"] = conversion.TimeToString(role.GetCreatedDate())
	}
	return out
}

func resourceCloudProviderAccessSetupImportState(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
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

	err2 := resourceCloudProviderAccessSetupRead(ctx, d, meta)

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
