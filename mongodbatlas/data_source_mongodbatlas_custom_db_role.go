package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasCustomDBRole() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasCustomDBRoleRead,
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

func dataSourceMongoDBAtlasCustomDBRoleRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)
	roleName := d.Get("role_name").(string)

	customDBRole, _, err := conn.CustomDBRoles.Get(context.Background(), projectID, roleName)
	if err != nil {
		return fmt.Errorf("error getting custom db role information: %s", err)
	}

	if err := d.Set("role_name", customDBRole.RoleName); err != nil {
		return fmt.Errorf("error setting `role_name` for custom db role (%s): %s", d.Id(), err)
	}
	if err := d.Set("actions", flattenActions(customDBRole.Actions)); err != nil {
		return fmt.Errorf("error setting `actions` for custom db role (%s): %s", d.Id(), err)
	}
	if err := d.Set("inherited_roles", flattenInheritedRoles(customDBRole.InheritedRoles)); err != nil {
		return fmt.Errorf("error setting `inherited_roles` for custom db role (%s): %s", d.Id(), err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"role_name":  customDBRole.RoleName,
	}))

	return nil
}
