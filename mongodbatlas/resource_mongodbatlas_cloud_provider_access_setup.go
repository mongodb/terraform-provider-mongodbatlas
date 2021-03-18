package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

/*
	mongodb_atlas_cloud_provider_access_setup
	-> Creates the the information from the mongodbatlas side
	-> The delete deletes and deauthorize the role
*/

func resourceMongoDBAtlasCloudProviderAccessSetup() *schema.Resource {
	return &schema.Resource{
		Read:   resourceMongoDBAtlasCloudProviderAccessSetupRead,
		Create: resourceMongoDBAtlasCloudProviderAccessSetupCreate,
		Update: resourceMongoDBAtlasCloudProviderAccessAuthorizationPlaceHolder,
		Delete: resourceMongoDBAtlasCloudProviderAccessSetupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMongoDBAtlasCloudProviderAccessSetupImportState,
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

func resourceMongoDBAtlasCloudProviderAccessSetupRead(d *schema.ResourceData, meta interface{}) error {
	// sadly there is no just get API
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	providerName := ids["provider_name"]

	roles, _, err := conn.CloudProviderAccess.ListRoles(context.Background(), projectID)

	if err != nil {
		return fmt.Errorf(errorGetRead, err)
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
				return fmt.Errorf(errorGetRead, err)
			}
		}
	} else {
		// planning for the future multiple providers
		return fmt.Errorf(errorGetRead,
			fmt.Sprintf("unsopported provider type %s", providerName))
	}

	return nil
}

func resourceMongoDBAtlasCloudProviderAccessSetupCreate(d *schema.ResourceData, meta interface{}) error {
	projectID := d.Get("project_id").(string)

	conn := meta.(*matlas.Client)

	requestParameters := &matlas.CloudProviderAccessRoleRequest{
		ProviderName: d.Get("provider_name").(string),
	}

	role, _, err := conn.CloudProviderAccess.CreateRole(context.Background(), projectID, requestParameters)

	if err != nil {
		return fmt.Errorf(errorCloudProviderAccessCreate, err)
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
			return fmt.Errorf(errorCloudProviderAccessCreate, err)
		}
	}

	return nil
}

func resourceMongoDBAtlasCloudProviderAccessSetupDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())

	projectID := ids["project_id"]
	roleID := ids["id"]
	providerName := ids["provider_name"]

	req := &matlas.CloudProviderDeauthorizationRequest{
		ProviderName: providerName,
		RoleID:       roleID,
		GroupID:      projectID,
	}

	_, err := conn.CloudProviderAccess.DeauthorizeRole(context.Background(), req)

	if err != nil {
		return fmt.Errorf(errorCloudProviderAccessDelete, err)
	}

	return nil
}

func roleToSchemaSetup(role *matlas.AWSIAMRole) map[string]interface{} {
	out := map[string]interface{}{
		"provider_name": role.ProviderName,
		"aws": map[string]interface{}{
			"atlas_aws_account_arn":          role.AtlasAWSAccountARN,
			"atlas_assumed_role_external_id": role.AtlasAssumedRoleExternalID,
		},
		"created_date": role.CreatedDate,
		"role_id":      role.RoleID,
	}

	return out
}

func resourceMongoDBAtlasCloudProviderAccessSetupImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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

	err = resourceMongoDBAtlasCloudProviderAccessSetupRead(d, meta)

	if err != nil {
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
