package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

/*
	A cloud provider access authorization
*/

func resourceMongoDBAtlasCloudProviderAccessAuthorization() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceMongoDBAtlasCloudProviderAccessAuthorizationRead,
		CreateContext: resourceMongoDBAtlasCloudProviderAccessAuthorizationCreate,
		UpdateContext: resourceMongoDBAtlasCloudProviderAccessAuthorizationUpdate,
		DeleteContext: resourceMongoDBAtlasCloudProviderAccessAuthorizationPlaceHolder,

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
				Type:    resourceMongoDBAtlasCloudProviderAccessAuthorizationResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceMongoDBAtlasCloudProviderAccessAuthorizationStateUpgradeV0,
				Version: 0,
			},
		},
	}
}

func resourceMongoDBAtlasCloudProviderAccessAuthorizationRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// sadly there is no just get API
	conn := meta.(*config.MongoDBClient).Atlas
	ids := config.DecodeStateID(d.Id())

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
		return diag.FromErr(fmt.Errorf(config.ErrorGetRead, "cloud provider access role not found in mongodbatlas, please create it first"))
	}

	roleSchema := roleToSchemaAuthorization(targetRole)
	for key, val := range roleSchema {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(fmt.Errorf(config.ErrorGetRead, err))
		}
	}

	// If not authorize , then request the authorization
	if targetRole.ProviderName == config.AWS && targetRole.AuthorizedDate == "" && !d.IsNewResource() {
		d.SetId("")
		return nil
	}

	return nil
}

func resourceMongoDBAtlasCloudProviderAccessAuthorizationCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).Atlas

	projectID := d.Get("project_id").(string)
	roleID := d.Get("role_id").(string)

	// validation
	targetRole, err := FindRole(ctx, conn, projectID, roleID)

	if err != nil {
		return diag.FromErr(err)
	}

	if targetRole == nil {
		return diag.FromErr(fmt.Errorf(config.ErrorGetRead, "cloud provider access role not found in mongodbatlas, please create it first"))
	}

	return authorizeRole(ctx, conn, d, projectID, targetRole)
}

func resourceMongoDBAtlasCloudProviderAccessAuthorizationUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).Atlas
	ids := config.DecodeStateID(d.Id())

	roleID := ids["id"]
	projectID := ids["project_id"]

	targetRole, err := FindRole(ctx, conn, projectID, roleID)

	if err != nil {
		return diag.FromErr(err)
	}

	if targetRole == nil {
		return diag.FromErr(fmt.Errorf(config.ErrorGetRead, "cloud provider access role not found in mongodbatlas, please create it first"))
	}

	if d.HasChange("aws") || d.HasChange("azure") {
		return authorizeRole(ctx, conn, d, projectID, targetRole)
	}

	return nil
}

func resourceMongoDBAtlasCloudProviderAccessAuthorizationPlaceHolder(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	d.SetId("")
	return nil
}

func roleToSchemaAuthorization(role *matlas.CloudProviderAccessRole) map[string]any {
	out := map[string]any{
		"role_id": role.RoleID,
		"aws": []any{map[string]any{
			"iam_assumed_role_arn": role.IAMAssumedRoleARN,
		}},
		"authorized_date": role.AuthorizedDate,
	}

	if role.ProviderName == "AZURE" {
		out = map[string]any{
			"role_id": role.AzureID,
			"azure": []any{map[string]any{
				"atlas_azure_app_id":   role.AtlasAzureAppID,
				"service_principal_id": role.AzureServicePrincipalID,
				"tenant_id":            role.AzureTenantID,
			}},
			"authorized_date": role.AuthorizedDate,
		}
	}

	features := make([]map[string]any, 0, len(role.FeatureUsages))
	for _, featureUsage := range role.FeatureUsages {
		features = append(features, featureToSchema(featureUsage))
	}

	out["feature_usages"] = features
	return out
}

func FindRole(ctx context.Context, conn *matlas.Client, projectID, roleID string) (*matlas.CloudProviderAccessRole, error) {
	role, _, err := conn.CloudProviderAccess.GetRole(ctx, projectID, roleID)
	if err != nil {
		return nil, fmt.Errorf(config.ErrorGetRead, err)
	}

	return role, nil
}

func resourceMongoDBAtlasCloudProviderAccessAuthorizationResourceV0() *schema.Resource {
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

func resourceMongoDBAtlasCloudProviderAccessAuthorizationStateUpgradeV0(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	rawState["aws"] = []any{}

	return rawState, nil
}

func authorizeRole(ctx context.Context, client *matlas.Client, d *schema.ResourceData, projectID string, targetRole *matlas.CloudProviderAccessRole) diag.Diagnostics {
	req := &matlas.CloudProviderAccessRoleRequest{
		ProviderName: targetRole.ProviderName,
	}

	roleID := targetRole.RoleID
	if targetRole.ProviderName == config.AWS {
		roleAWS, ok := d.GetOk("aws")
		if !ok {
			return diag.FromErr(fmt.Errorf("error CloudProviderAccessAuthorization missing iam_assumed_role_arn"))
		}

		req.IAMAssumedRoleARN = pointer(roleAWS.([]any)[0].(map[string]any)["iam_assumed_role_arn"].(string))
	}

	if targetRole.ProviderName == config.AZURE {
		req.AtlasAzureAppID = targetRole.AtlasAzureAppID
		req.AzureTenantID = targetRole.AzureTenantID
		req.AzureServicePrincipalID = targetRole.AzureServicePrincipalID
		roleID = *targetRole.AzureID
	}

	var role *matlas.CloudProviderAccessRole
	var err error

	for i := 0; i < 3; i++ {
		role, _, err = client.CloudProviderAccess.AuthorizeRole(ctx, projectID, roleID, req)
		if err != nil && strings.Contains(err.Error(), "CANNOT_ASSUME_ROLE") { // aws takes time to update , in case of single path
			log.Printf("warning issue performing authorize: %s \n", err.Error())
			log.Println("retrying")
			time.Sleep(10 * time.Second)
			continue
		}
		if err != nil {
			log.Printf("MISSED ERRROR %s", err.Error())
		}
		break
	}

	if err != nil {
		return diag.FromErr(fmt.Errorf("error cloud provider access authorization %s", err))
	}

	authSchema := roleToSchemaAuthorization(role)

	resourceID := role.RoleID
	if role.ProviderName == config.AZURE {
		resourceID = *role.AzureID
	}
	d.SetId(config.EncodeStateID(map[string]string{
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
