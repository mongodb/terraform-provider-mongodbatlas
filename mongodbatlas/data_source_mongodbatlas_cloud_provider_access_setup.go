package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasCloudProviderAccessSetup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasCloudProviderAccessSetupRead,
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
			"role_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"aws": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
		},
	}
}

func dataSourceMongoDBAtlasCloudProviderAccessSetupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)
	providerName := d.Get("provider_name").(string)
	roleID := d.Get("role_id").(string)

	roles, _, err := conn.CloudProviderAccess.ListRoles(ctx, projectID)

	if err != nil {
		return diag.FromErr(fmt.Errorf(errorGetRead, err))
	}

	// aws specific
	if providerName == "AWS" {
		var targetRole matlas.AWSIAMRole
		// searching in roles
		for i := range roles.AWSIAMRoles {
			role := &(roles.AWSIAMRoles[i])
			if role.RoleID == roleID && role.ProviderName == providerName {
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
			fmt.Sprintf("unsupported provider type %s", providerName)))
	}

	d.SetId(resource.UniqueId())

	return nil
}
