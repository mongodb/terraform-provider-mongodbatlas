package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

/*
	mongodb_atlas_cloud_provider_access_setup
	-> Creates the the information from the mongodbatlas side
	-> The delete deletes and deauthorize the role
*/

func resourceMongoDBAtlasCloudProviderAccessSetup() *schema.Resource {
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
			// Note: when new providers will be added, this will trigger a recreate
			"provider_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"AWS"}, false),
			},
			"aws": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Deprecated: "use aws_config instead",
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
			"created_date": {
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

func resourceMongoDBAtlasCloudProviderAccessSetupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// sadly there is no just get API
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	providerName := ids["provider_name"]

	roles, resp, err := conn.CloudProviderAccess.ListRoles(context.Background(), projectID)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorGetRead, err))
	}

	// aws specific
	if providerName == "AWS" {
		var targetRole matlas.AWSIAMRole
		// searching in roles
		for i := range roles.AWSIAMRoles {
			role := &(roles.AWSIAMRoles[i])
			if role.RoleID == ids["id"] && role.ProviderName == ids["provider_name"] {
				targetRole = *role
			}
		}
		// Not Found
		if targetRole.RoleID == "" && !d.IsNewResource() {
			d.SetId("")
			return nil
		}

		roleSchema := roleToSchemaSetup(&targetRole)

		for key, val := range roleSchema {
			if err := d.Set(key, val); err != nil {
				return diag.FromErr(fmt.Errorf(errorGetRead, err))
			}
		}
	} else {
		// planning for the future multiple providers
		return diag.FromErr(fmt.Errorf(errorGetRead,
			fmt.Sprintf("unsopported provider type %s", providerName)))
	}

	return nil
}

func resourceMongoDBAtlasCloudProviderAccessSetupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectID := d.Get("project_id").(string)

	conn := meta.(*MongoDBClient).Atlas

	requestParameters := &matlas.CloudProviderAccessRoleRequest{
		ProviderName: d.Get("provider_name").(string),
	}

	role, _, err := conn.CloudProviderAccess.CreateRole(ctx, projectID, requestParameters)

	if err != nil {
		return diag.FromErr(fmt.Errorf(errorCloudProviderAccessCreate, err))
	}

	// once multiple providers enable here do a switch, select for provider type
	roleSchema := roleToSchemaSetup(role)

	d.SetId(encodeStateID(map[string]string{
		"id":            role.RoleID,
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

func resourceMongoDBAtlasCloudProviderAccessSetupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func roleToSchemaSetup(role *matlas.AWSIAMRole) map[string]interface{} {
	out := map[string]interface{}{
		"provider_name": role.ProviderName,
		"aws": map[string]interface{}{ // Deprecated, will be deleted later
			"atlas_aws_account_arn":          role.AtlasAWSAccountARN,
			"atlas_assumed_role_external_id": role.AtlasAssumedRoleExternalID,
		},
		"aws_config": []interface{}{map[string]interface{}{
			"atlas_aws_account_arn":          role.AtlasAWSAccountARN,
			"atlas_assumed_role_external_id": role.AtlasAssumedRoleExternalID,
		}},
		"created_date": role.CreatedDate,
		"role_id":      role.RoleID,
	}

	return out
}

func resourceMongoDBAtlasCloudProviderAccessSetupImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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
