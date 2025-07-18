package customdbrole

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"actions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resources": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"collection_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"database_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cluster": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"inherited_roles": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)
	roleName := d.Get("role_name").(string)

	customDBRole, _, err := connV2.CustomDatabaseRolesApi.GetCustomDatabaseRole(ctx, projectID, roleName).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting custom db role information: %s", err))
	}

	if err := d.Set("role_name", customDBRole.RoleName); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `role_name` for custom db role (%s): %s", d.Id(), err))
	}

	if err := d.Set("actions", flattenActions(customDBRole.GetActions())); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `actions` for custom db role (%s): %s", d.Id(), err))
	}

	if err := d.Set("inherited_roles", flattenInheritedRoles(customDBRole.GetInheritedRoles())); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `inherited_roles` for custom db role (%s): %s", d.Id(), err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id": projectID,
		"role_name":  customDBRole.RoleName,
	}))

	return nil
}
