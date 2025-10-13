package customdbrole

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312008/admin"
)

func PluralDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePluralRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_name": {
							Type:     schema.TypeString,
							Computed: true,
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
										Computed: true,
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
							Computed: true,
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
				},
			},
		},
	}
}

func dataSourcePluralRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)

	customDBRoles, _, err := connV2.CustomDatabaseRolesApi.ListCustomDbRoles(ctx, projectID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting custom db roles information: %s", err))
	}

	if err := d.Set("results", flattenCustomDBRoles(customDBRoles)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `results for custom db roles: %s", err))
	}

	d.SetId(id.UniqueId())

	return nil
}

func flattenCustomDBRoles(customDBRoles []admin.UserCustomDBRole) []map[string]any {
	customDBRolesMap := make([]map[string]any, len(customDBRoles))
	for k, customDBRole := range customDBRoles {
		customDBRolesMap[k] = map[string]any{
			"role_name":       customDBRole.RoleName,
			"actions":         flattenActions(customDBRole.GetActions()),
			"inherited_roles": flattenInheritedRoles(customDBRole.GetInheritedRoles()),
		}
	}
	return customDBRolesMap
}
