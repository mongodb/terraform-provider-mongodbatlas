package cloudprovideraccess

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"go.mongodb.org/atlas-sdk/v20250312008/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/cleanup"
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
	errorCreate                  = "error creating cloud provider access %s"
	errorUpdate                  = "error updating cloud provider access %s"
	errorDelete                  = "error deleting cloud provider access %s"
	errorImporter                = "error importing cloud provider access %s"
	ErrorGetRead                 = "error reading cloud provider access %s"
	defaultTimeout time.Duration = 20 * time.Minute
)

func ResourceSetup() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout:   resourceCloudProviderAccessSetupRead,
		CreateWithoutTimeout: resourceCloudProviderAccessSetupCreate,
		UpdateWithoutTimeout: resourceCloudProviderAccessAuthorizationPlaceHolder,
		DeleteWithoutTimeout: resourceCloudProviderAccessSetupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceCloudProviderAccessSetupImportState,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(defaultTimeout),
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
			"delete_on_create_timeout": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether to delete the resource being created if a timeout is reached when waiting for completion. When set to `true` and timeout occurs, it triggers the deletion and returns immediately without waiting for deletion to complete. When set to `false`, the timeout will not trigger resource deletion. If you suspect a transient error when the value is `true`, wait before retrying to allow resource deletion to finish. Default is `true`.",
			},
		},
	}
}

func resourceCloudProviderAccessSetupRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	roleID := ids["id"]

	role, resp, err := conn.CloudProviderAccessApi.GetCloudProviderAccess(context.Background(), projectID, roleID).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf(ErrorGetRead, err))
	}

	roleSchema, err := roleToSchemaSetup(role)
	if err != nil {
		return diag.Errorf(errorCreate, err)
	}
	for key, val := range roleSchema {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(fmt.Errorf(ErrorGetRead, err))
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

	role, _, err := conn.CloudProviderAccessApi.CreateCloudProviderAccess(ctx, projectID, requestParameters).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorCreate, err))
	}

	resourceID := role.GetRoleId()
	if role.ProviderName == constant.AZURE {
		resourceID = role.GetId() // For Azure, the unique identifier is in the "id" field, not "roleId".
	}

	if role.ProviderName == constant.GCP {
		deleteOnCreateTimeout := true // default value when not set
		if v, ok := d.GetOkExists("delete_on_create_timeout"); ok {
			deleteOnCreateTimeout = v.(bool)
		}
		r, err := waitForGCPProviderAccessCompletion(ctx, projectID, resourceID, conn, d.Timeout(schema.TimeoutCreate), deleteOnCreateTimeout)
		if err != nil {
			return diag.FromErr(err)
		}
		role = r
	}

	// once multiple providers enable here do a switch, select for provider type
	roleSchema, err := roleToSchemaSetup(role)
	if err != nil {
		return diag.Errorf(errorCreate, err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"id":            resourceID,
		"project_id":    projectID,
		"provider_name": role.GetProviderName(),
	}))

	for key, val := range roleSchema {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(fmt.Errorf(errorCreate, err))
		}
	}

	return nil
}

func waitForGCPProviderAccessCompletion(ctx context.Context, projectID, resourceID string, conn *admin.APIClient, timeout time.Duration, deleteOnCreateTimeout bool) (*admin.CloudProviderAccessRole, error) {
	requestParams := &admin.GetCloudProviderAccessApiParams{
		RoleId:  resourceID,
		GroupId: projectID,
	}

	stateConf := retry.StateChangeConf{
		Pending:    []string{"IN_PROGRESS", "NOT_INITIATED"},
		Target:     []string{"COMPLETE", "FAILED"},
		Refresh:    resourceRefreshFunc(ctx, requestParams, conn),
		Timeout:    timeout,
		MinTimeout: 60 * time.Second,
		Delay:      30 * time.Second,
	}

	finalResponse, err := stateConf.WaitForStateContext(ctx)
	err = cleanup.HandleCreateTimeout(deleteOnCreateTimeout, err, func(ctxCleanup context.Context) error {
		_, errCleanup := conn.CloudProviderAccessApi.DeauthorizeProviderAccessRole(ctxCleanup, projectID, constant.GCP, resourceID).Execute()
		return errCleanup
	})
	if err != nil {
		return nil, err
	}

	r, ok := finalResponse.(*admin.CloudProviderAccessRole)
	if !ok {
		return nil, fmt.Errorf("unexpected type for result: %T", finalResponse)
	}
	return r, nil
}

func resourceRefreshFunc(ctx context.Context, requestParams *admin.GetCloudProviderAccessApiParams, conn *admin.APIClient) retry.StateRefreshFunc {
	return func() (any, string, error) {
		role, resp, err := conn.CloudProviderAccessApi.GetCloudProviderAccessWithParams(ctx, requestParams).Execute()
		if err != nil {
			if validate.StatusNotFound(resp) {
				return nil, "", fmt.Errorf("cloud provider access role %q not found in project %q", requestParams.RoleId, requestParams.GroupId)
			}
			return nil, "", err
		}

		status := role.GetStatus()
		switch status {
		case "IN_PROGRESS", "NOT_INITIATED", "COMPLETE":
			return role, status, nil
		default:
			return nil, status, fmt.Errorf("cloud provider access setup failed %q for role %q", status, requestParams.RoleId)
		}
	}
}

func resourceCloudProviderAccessSetupDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())

	projectID := ids["project_id"]
	roleID := ids["id"]
	providerName := ids["provider_name"]

	req := &admin.DeauthorizeProviderAccessRoleApiParams{
		CloudProvider: providerName,
		RoleId:        roleID,
		GroupId:       projectID,
	}

	_, err := conn.CloudProviderAccessApi.DeauthorizeProviderAccessRoleWithParams(ctx, req).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorDelete, err))
	}

	d.SetId("")
	return nil
}

func roleToSchemaSetup(role *admin.CloudProviderAccessRole) (map[string]any, error) {
	switch role.ProviderName {
	case constant.AWS:
		return map[string]any{
			"provider_name": role.GetProviderName(),
			"aws_config": []any{map[string]any{
				"atlas_aws_account_arn":          role.GetAtlasAWSAccountArn(),
				"atlas_assumed_role_external_id": role.GetAtlasAssumedRoleExternalId(),
			}},
			"gcp_config":   []any{map[string]any{}},
			"created_date": conversion.TimeToString(role.GetCreatedDate()),
			"role_id":      role.GetRoleId(), // For AWS, the unique identifier is in the "roleId" field
		}, nil
	case constant.AZURE:
		return map[string]any{
			"provider_name": role.ProviderName,
			"azure_config": []any{map[string]any{
				"atlas_azure_app_id":   role.GetAtlasAzureAppId(),
				"service_principal_id": role.GetServicePrincipalId(),
				"tenant_id":            role.GetTenantId(),
			}},
			"aws_config":        []any{map[string]any{}},
			"gcp_config":        []any{map[string]any{}},
			"created_date":      conversion.TimeToString(role.GetCreatedDate()),
			"last_updated_date": conversion.TimeToString(role.GetLastUpdatedDate()),
			"role_id":           role.GetId(), // For Azure, the unique identifier is in the "id" field, not "roleId".
		}, nil
	case constant.GCP:
		return map[string]any{
			"provider_name": role.GetProviderName(),
			"gcp_config": []any{map[string]any{
				"status":                    role.GetStatus(),
				"service_account_for_atlas": role.GetGcpServiceAccountForAtlas(),
			}},
			"aws_config":   []any{map[string]any{}},
			"role_id":      role.GetRoleId(), // For GCP, the unique identifier is in the "roleId"
			"created_date": conversion.TimeToString(role.GetCreatedDate()),
		}, nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", role.GetProviderName())
	}
}

func resourceCloudProviderAccessSetupImportState(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	projectID, providerName, roleID, err := splitCloudProviderAccessID(d.Id())

	if err != nil {
		return nil, fmt.Errorf(errorImporter, err)
	}

	// searching id in internal format
	d.SetId(conversion.EncodeStateID(map[string]string{
		"id":            roleID,
		"project_id":    projectID,
		"provider_name": providerName,
	}))

	err2 := resourceCloudProviderAccessSetupRead(ctx, d, meta)

	if err2 != nil {
		return nil, fmt.Errorf(errorImporter, err)
	}

	// case of not found
	if d.Id() == "" {
		return nil, fmt.Errorf(errorImporter, " Resource not found at the cloud please check your id")
	}

	// params syncing
	if err = d.Set("project_id", projectID); err != nil {
		return nil, fmt.Errorf(errorImporter, err)
	}

	return []*schema.ResourceData{d}, nil
}

// format  {project_id}-{provider-name}-{role-id}
func splitCloudProviderAccessID(id string) (projectID, providerName, roleID string, err error) {
	var re = regexp.MustCompile(`(?s)^([0-9a-fA-F]{24})-(.*)-(.*)$`)
	parts := re.FindStringSubmatch(id)

	if len(parts) != 4 {
		err = fmt.Errorf(errorImporter, "format please use {project_id}-{provider-name}-{role-id}")
		return
	}

	projectID, providerName, roleID = parts[1], parts[2], parts[3]

	return
}
