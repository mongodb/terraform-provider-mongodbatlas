package mongodbatlas

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func resourceMongoDBAtlasCustomDBRole() *schema.Resource {
	return &schema.Resource{
		Create:   resourceMongoDBAtlasCustomDBRoleCreate,
		Read:     resourceMongoDBAtlasCustomDBRoleRead,
		Update:   nil,
		Delete:   nil,
		Importer: nil,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"actions": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:     schema.TypeString,
							Required: true,
						},
						"resources": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"collection_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"database_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"cluster": {
										Type:     schema.TypeBool,
										Optional: true,
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
							Required: true,
						},
						"role_name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceMongoDBAtlasCustomDBRoleCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)

	customDBRoleReq := &matlas.CustomDBRole{
		RoleName:       d.Get("role_name").(string),
		Actions:        expandActions(d),
		InheritedRoles: expandInheritedRoles(d),
	}

	customDBRoleRes, _, err := conn.CustomDBRoles.Create(context.Background(), projectID, customDBRoleReq)

	if err != nil {
		return fmt.Errorf("error creating custom db role: %s", err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"role_name":  customDBRoleRes.RoleName,
	}))

	return resourceMongoDBAtlasCustomDBRoleRead(d, meta)
}

func resourceMongoDBAtlasCustomDBRoleRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	roleName := ids["role_name"]

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

	return nil
}

func expandActions(d *schema.ResourceData) []matlas.Action {
	var actions []matlas.Action
	if v, ok := d.GetOk("actions"); ok {
		if rs := v.([]interface{}); len(rs) > 0 {
			actions = make([]matlas.Action, len(rs))
			for k, a := range rs {
				actionMap := a.(map[string]interface{})
				actions[k] = matlas.Action{
					Action:    actionMap["action"].(string),
					Resources: expandActionResources(actionMap["resources"].([]interface{})),
				}
			}
		}
	}
	return actions
}

func expandActionResources(resources []interface{}) []matlas.Resource {
	actionResources := make([]matlas.Resource, len(resources))
	for k, v := range resources {
		resourceMap := v.(map[string]interface{})
		if cluster := resourceMap["cluster"]; cluster.(bool) {
			actionResources[k] = matlas.Resource{
				Cluster: resourceMap["cluster"].(bool),
			}
		} else {
			actionResources[k] = matlas.Resource{
				Db:         resourceMap["database_name"].(string),
				Collection: resourceMap["collection_name"].(string),
			}
		}
	}
	return actionResources
}

func flattenActions(actions []matlas.Action) []map[string]interface{} {
	actionList := make([]map[string]interface{}, 0)
	for _, v := range actions {
		actionList = append(actionList, map[string]interface{}{
			"action":    v.Action,
			"resources": flattenActionResources(v.Resources),
		})
	}
	return actionList
}

func flattenActionResources(resources []matlas.Resource) []map[string]interface{} {
	actionResourceList := make([]map[string]interface{}, 0)
	for _, v := range resources {
		if cluster := v.Cluster; cluster {
			actionResourceList = append(actionResourceList, map[string]interface{}{
				"cluster": v.Cluster,
			})
		} else {
			actionResourceList = append(actionResourceList, map[string]interface{}{
				"database_name":   v.Db,
				"collection_name": v.Collection,
			})
		}
	}
	return actionResourceList
}

func expandInheritedRoles(d *schema.ResourceData) []matlas.InheritedRole {
	var inheritedRoles []matlas.InheritedRole
	if v, ok := d.GetOk("inherited_roles"); ok {
		if rs := v.([]interface{}); len(rs) > 0 {
			inheritedRoles = make([]matlas.InheritedRole, len(rs))
			for k, r := range rs {
				roleMap := r.(map[string]interface{})
				inheritedRoles[k] = matlas.InheritedRole{
					Db:   roleMap["database_name"].(string),
					Role: roleMap["role_name"].(string),
				}
			}
		}
	}
	return inheritedRoles
}

func flattenInheritedRoles(roles []matlas.InheritedRole) []map[string]interface{} {
	inheritedRoleList := make([]map[string]interface{}, 0)
	for _, v := range roles {
		inheritedRoleList = append(inheritedRoleList, map[string]interface{}{
			"database_name": v.Db,
			"role_name":     v.Role,
		})
	}
	return inheritedRoleList
}
