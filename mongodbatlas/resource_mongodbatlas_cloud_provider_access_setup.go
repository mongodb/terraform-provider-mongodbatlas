package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func resourceMongoDBAtlasCloudProviderAccessSetup() *schema.Resource {
	return &schema.Resource{
		Read: resourceMongoDBAtlasCloudProviderAccessSetupRead,
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
