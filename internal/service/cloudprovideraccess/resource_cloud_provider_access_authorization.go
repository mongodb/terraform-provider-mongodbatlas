package cloudprovideraccess

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312006/admin"
)

/*
	A cloud provider access authorization
*/

func ResourceAuthorization() *schema.Resource {
	// HELP-71400 provides context as to why import is not implemented
	return &schema.Resource{
		ReadContext:   resourceCloudProviderAccessAuthorizationRead,
		CreateContext: resourceCloudProviderAccessAuthorizationCreate,
		UpdateContext: resourceCloudProviderAccessAuthorizationUpdate,
		DeleteContext: resourceCloudProviderAccessAuthorizationPlaceHolder,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"aws": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"iam_assumed_role_arn": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"azure": {
				Type:     schema.TypeList,
				MaxItems: 1,
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
			"gcp": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_account_for_atlas": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"feature_usages": {
				Type:     schema.TypeList,
				Elem:     featureUsagesSchema(),
				Computed: true,
			},
			"authorized_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceCloudProviderAccessAuthorizationResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceCloudProviderAccessAuthorizationStateUpgradeV0,
				Version: 0,
			},
		},
	}
}

func resourceCloudProviderAccessAuthorizationRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// sadly there is no just get API
	conn := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())

	roleID := ids["id"] // atlas ID
	projectID := ids["project_id"]

	targetRole, err := FindRole(ctx, conn, projectID, roleID)
	if err != nil {
		reset := strings.Contains(err.Error(), "404") && !d.IsNewResource()
		if reset {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	if targetRole == nil {
		return diag.FromErr(fmt.Errorf(ErrorCloudProviderGetRead, "cloud provider access role not found in mongodbatlas, please create it first"))
	}

	roleSchema := roleToSchemaAuthorization(targetRole)
	for key, val := range roleSchema {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(fmt.Errorf(ErrorCloudProviderGetRead, err))
		}
	}

	// If not authorize , then request the authorization
	if targetRole.ProviderName == constant.AWS && conversion.TimeToString(targetRole.GetAuthorizedDate()) == "" && !d.IsNewResource() {
		d.SetId("")
		return nil
	}

	return nil
}

func resourceCloudProviderAccessAuthorizationCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2

	projectID := d.Get("project_id").(string)
	roleID := d.Get("role_id").(string)

	// validation
	targetRole, err := FindRole(ctx, conn, projectID, roleID)

	if err != nil {
		return diag.FromErr(err)
	}

	if targetRole == nil {
		return diag.FromErr(fmt.Errorf(ErrorCloudProviderGetRead, "cloud provider access role not found in mongodbatlas, please create it first"))
	}

	return authorizeRole(ctx, conn, d, projectID, targetRole)
}

func resourceCloudProviderAccessAuthorizationUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())

	roleID := ids["id"]
	projectID := ids["project_id"]

	targetRole, err := FindRole(ctx, conn, projectID, roleID)

	if err != nil {
		return diag.FromErr(err)
	}

	if targetRole == nil {
		return diag.FromErr(fmt.Errorf(ErrorCloudProviderGetRead, "cloud provider access role not found in mongodbatlas, please create it first"))
	}

	if d.HasChange("aws") || d.HasChange("azure") {
		return authorizeRole(ctx, conn, d, projectID, targetRole)
	}

	return nil
}

func resourceCloudProviderAccessAuthorizationPlaceHolder(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	d.SetId("")
	return nil
}

func roleToSchemaAuthorization(role *admin.CloudProviderAccessRole) map[string]any {
	out := map[string]any{
		"role_id": role.GetRoleId(),
		"aws": []any{map[string]any{
			"iam_assumed_role_arn": role.GetIamAssumedRoleArn(),
		}},
		"authorized_date": conversion.TimeToString(role.GetAuthorizedDate()),
		"gcp":             []any{map[string]any{}},
	}

	if role.ProviderName == "AZURE" {
		out = map[string]any{
			"role_id": role.GetId(),
			"azure": []any{map[string]any{
				"atlas_azure_app_id":   role.GetAtlasAzureAppId(),
				"service_principal_id": role.GetServicePrincipalId(),
				"tenant_id":            role.GetTenantId(),
			}},
			"authorized_date": conversion.TimeToString(role.GetAuthorizedDate()),
			"gcp":             []any{map[string]any{}},
		}
	}
	if role.ProviderName == "GCP" {
		out = map[string]any{
			"role_id": role.GetRoleId(),
			"gcp": []any{map[string]any{
				"service_account_for_atlas": role.GetGcpServiceAccountForAtlas(),
			}},
		}
	}

	features := make([]map[string]any, 0, len(role.GetFeatureUsages()))
	for _, featureUsage := range role.GetFeatureUsages() {
		features = append(features, featureToSchema(featureUsage))
	}

	out["feature_usages"] = features
	return out
}

func FindRole(ctx context.Context, conn *admin.APIClient, projectID, roleID string) (*admin.CloudProviderAccessRole, error) {
	role, _, err := conn.CloudProviderAccessApi.GetCloudProviderAccessRole(ctx, projectID, roleID).Execute()
	if err != nil {
		return nil, fmt.Errorf(ErrorCloudProviderGetRead, err)
	}

	return role, nil
}

func resourceCloudProviderAccessAuthorizationResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"aws": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"iam_assumed_role_arn": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"feature_usages": {
				Type:     schema.TypeList,
				Elem:     featureUsagesSchema(),
				Computed: true,
			},
			"authorized_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCloudProviderAccessAuthorizationStateUpgradeV0(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	rawState["aws"] = []any{}

	return rawState, nil
}

func authorizeRole(ctx context.Context, client *admin.APIClient, d *schema.ResourceData, projectID string, targetRole *admin.CloudProviderAccessRole) diag.Diagnostics {
	req := &admin.CloudProviderAccessRoleRequestUpdate{
		ProviderName: targetRole.ProviderName,
	}

	roleID := targetRole.GetRoleId()
	if targetRole.ProviderName == constant.AWS {
		roleAWS, ok := d.GetOk("aws")
		if !ok {
			return diag.FromErr(fmt.Errorf("error CloudProviderAccessAuthorization missing iam_assumed_role_arn"))
		}

		req.SetIamAssumedRoleArn(roleAWS.([]any)[0].(map[string]any)["iam_assumed_role_arn"].(string))
	}

	if targetRole.ProviderName == constant.AZURE {
		req.SetAtlasAzureAppId(targetRole.GetAtlasAzureAppId())
		req.SetTenantId(targetRole.GetTenantId())
		req.SetServicePrincipalId(targetRole.GetServicePrincipalId())
		roleID = targetRole.GetId()
	}
	// No specific GCP config is needed, only providerName and roleID are needed

	var role *admin.CloudProviderAccessRole
	var err error

	for range 3 {
		role, _, err = client.CloudProviderAccessApi.AuthorizeCloudProviderAccessRole(ctx, projectID, roleID, req).Execute()
		if err != nil && strings.Contains(err.Error(), "CANNOT_ASSUME_ROLE") { // aws takes time to update , in case of single path
			log.Printf("warning issue performing authorize: %s \n", err.Error())
			log.Println("retrying")
			time.Sleep(10 * time.Second)
			continue
		}
		if err != nil {
			log.Printf("MISSED ERROR %s", err.Error())
		}
		break
	}

	if err != nil {
		return diag.FromErr(fmt.Errorf("error cloud provider access authorization %s", err))
	}

	authSchema := roleToSchemaAuthorization(role)

	resourceID := role.GetRoleId()
	if role.ProviderName == constant.AZURE {
		resourceID = role.GetId()
	}
	d.SetId(conversion.EncodeStateID(map[string]string{
		"id":         resourceID,
		"project_id": projectID,
	}))

	for key, val := range authSchema {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(fmt.Errorf(errorCloudProviderAccessCreate, err))
		}
	}

	return nil
}

func featureUsagesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"feature_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"feature_id": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func featureToSchema(feature admin.CloudProviderAccessFeatureUsage) map[string]any {
	featureID := feature.GetFeatureId()
	featureIDMap := map[string]any{
		"project_id":  featureID.GetGroupId(),
		"bucket_name": featureID.GetBucketName(),
	}
	return map[string]any{
		"feature_type": feature.GetFeatureType(),
		"feature_id":   featureIDMap,
	}
}
