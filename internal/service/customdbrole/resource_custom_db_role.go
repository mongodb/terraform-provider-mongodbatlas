package customdbrole

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/spf13/cast"
	"go.mongodb.org/atlas-sdk/v20250312007/admin"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate,
		ReadContext:   resourceRead,
		UpdateContext: resourceUpdate,
		DeleteContext: resourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceImport,
		},
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
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`[\w-]+`), "`role_name` can contain only letters, digits, underscores, and dashes"),
					func(v any, k string) (ws []string, es []error) {
						value := v.(string)
						if strings.HasPrefix(value, "x-gen") {
							es = append(es, fmt.Errorf("`role_name` cannot start with 'xgen-'"))
						}
						return
					},
				),
			},
			"actions": {
				Type:     schema.TypeList,
				Optional: true,
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
				Type:     schema.TypeSet,
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

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)

	customDBRoleReq := &admin.UserCustomDBRole{
		RoleName:       d.Get("role_name").(string),
		Actions:        expandActions(d),
		InheritedRoles: expandInheritedRoles(d),
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"created", "failed"},
		Refresh: func() (any, string, error) {
			// Atlas Create is called inside refresh because the endpoint doesn't support concurrent POST requests so it's retried if it fails because of that.
			customDBRoleRes, _, err := connV2.CustomDatabaseRolesApi.CreateCustomDbRole(ctx, projectID, customDBRoleReq).Execute()
			if err != nil {
				if strings.Contains(err.Error(), "Unexpected error") ||
					strings.Contains(err.Error(), "UNEXPECTED_ERROR") ||
					strings.Contains(err.Error(), "500") ||
					strings.Contains(err.Error(), "404") ||
					strings.Contains(err.Error(), "ATLAS_CUSTOM_ROLE_NOT_FOUND") {
					return nil, "pending", nil
				}
				return nil, "failed", err
			}

			return customDBRoleRes, "created", nil
		},
		Timeout:    10 * time.Minute,
		Delay:      3 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	// Wait, catching any errors
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating custom db role: %s", err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id": projectID,
		"role_name":  customDBRoleReq.RoleName,
	}))

	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	roleName := ids["role_name"]

	customDBRole, resp, err := connV2.CustomDatabaseRolesApi.GetCustomDbRole(ctx, projectID, roleName).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf("error getting custom db role information: %s", err))
	}

	if err := d.Set("role_name", customDBRole.GetRoleName()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `role_name` for custom db role (%s): %s", d.Id(), err))
	}

	if err := d.Set("actions", flattenActions(customDBRole.GetActions())); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `actions` for custom db role (%s): %s", d.Id(), err))
	}

	if err := d.Set("inherited_roles", flattenInheritedRoles(customDBRole.GetInheritedRoles())); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `inherited_roles` for custom db role (%s): %s", d.Id(), err))
	}

	return nil
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	roleName := ids["role_name"]

	if d.HasChange("actions") || d.HasChange("inherited_roles") {
		updateParams := &admin.UpdateCustomDBRole{
			Actions:        expandActions(d),
			InheritedRoles: expandInheritedRoles(d),
		}
		_, _, err := connV2.CustomDatabaseRolesApi.UpdateCustomDbRole(ctx, projectID, roleName, updateParams).Execute()
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating custom db role (%s): %s", roleName, err))
		}
	}
	return resourceRead(ctx, d, meta)
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	roleName := ids["role_name"]

	stateConf := &retry.StateChangeConf{
		Pending: []string{"deleting"},
		Target:  []string{"deleted", "failed"},
		Refresh: func() (any, string, error) {
			_, _, err := connV2.CustomDatabaseRolesApi.GetCustomDbRole(ctx, projectID, roleName).Execute()
			if err != nil {
				if strings.Contains(err.Error(), "404") ||
					strings.Contains(err.Error(), "ATLAS_CUSTOM_ROLE_NOT_FOUND") {
					return "", "deleted", nil
				}
				return nil, "failed", err
			}

			_, err = connV2.CustomDatabaseRolesApi.DeleteCustomDbRole(ctx, projectID, roleName).Execute()
			if err != nil {
				return nil, "failed", fmt.Errorf("error deleting custom db role (%s): %s", roleName, err)
			}

			return nil, "deleting", nil
		},
		Timeout:    10 * time.Minute,
		Delay:      3 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	// Wait, catching any errors
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 2 {
		return nil, errors.New("import format error: to import a custom db role use the format {project_id}-{role_name}")
	}

	projectID := parts[0]
	roleName := parts[1]

	r, _, err := connV2.CustomDatabaseRolesApi.GetCustomDbRole(ctx, projectID, roleName).Execute()
	if err != nil {
		return nil, fmt.Errorf("couldn't import custom db role %s in project %s, error: %s", roleName, projectID, err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id": projectID,
		"role_name":  r.GetRoleName(),
	}))

	if err := d.Set("project_id", projectID); err != nil {
		log.Printf("[WARN] Error setting project_id for (%s): %s", d.Id(), err)
	}

	return []*schema.ResourceData{d}, nil
}

func expandActions(d *schema.ResourceData) *[]admin.DatabasePrivilegeAction {
	actions := make([]admin.DatabasePrivilegeAction, len(d.Get("actions").([]any)))
	for k, v := range d.Get("actions").([]any) {
		a := v.(map[string]any)
		actions[k] = admin.DatabasePrivilegeAction{
			Action:    a["action"].(string),
			Resources: expandActionResources(a["resources"].(*schema.Set)),
		}
	}
	return &actions
}

func flattenActions(actions []admin.DatabasePrivilegeAction) []map[string]any {
	actionList := make([]map[string]any, 0)
	for _, v := range actions {
		actionList = append(actionList, map[string]any{
			"action":    v.Action,
			"resources": flattenActionResources(v.GetResources()),
		})
	}
	return actionList
}

func expandActionResources(resources *schema.Set) *[]admin.DatabasePermittedNamespaceResource {
	actionResources := make([]admin.DatabasePermittedNamespaceResource, resources.Len())
	for k, v := range resources.List() {
		resourceMap := v.(map[string]any)
		actionResources[k] = admin.DatabasePermittedNamespaceResource{
			Db:         resourceMap["database_name"].(string),
			Collection: resourceMap["collection_name"].(string),
			Cluster:    cast.ToBool(resourceMap["cluster"]),
		}
	}
	return &actionResources
}

func flattenActionResources(resources []admin.DatabasePermittedNamespaceResource) []map[string]any {
	actionResourceList := make([]map[string]any, 0)
	for _, v := range resources {
		if v.Cluster {
			actionResourceList = append(actionResourceList, map[string]any{
				"cluster": v.Cluster,
			})
		} else {
			actionResourceList = append(actionResourceList, map[string]any{
				"database_name":   cast.ToString(v.GetDb()),
				"collection_name": cast.ToString(v.GetCollection()),
			})
		}
	}
	return actionResourceList
}

func expandInheritedRoles(d *schema.ResourceData) *[]admin.DatabaseInheritedRole {
	vIR := d.Get("inherited_roles").(*schema.Set).List()
	ir := make([]admin.DatabaseInheritedRole, len(vIR))
	if len(vIR) != 0 {
		for i := range vIR {
			r := vIR[i].(map[string]any)
			ir[i] = admin.DatabaseInheritedRole{
				Db:   r["database_name"].(string),
				Role: r["role_name"].(string),
			}
		}
	}
	return &ir
}

func flattenInheritedRoles(roles []admin.DatabaseInheritedRole) []map[string]any {
	inheritedRoleList := make([]map[string]any, 0)
	for _, v := range roles {
		inheritedRoleList = append(inheritedRoleList, map[string]any{
			"database_name": v.GetDb(),
			"role_name":     v.GetRole(),
		})
	}
	return inheritedRoleList
}
