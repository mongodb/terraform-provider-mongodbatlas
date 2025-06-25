package federatedsettingsorgrolemapping

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312004/admin"
)

func PluralDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePluralRead,
		Schema: map[string]*schema.Schema{
			"federation_settings_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"org_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"page_num": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"items_per_page": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"external_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role_assignments": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"org_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"role": {
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
	conn := meta.(*config.MongoDBClient).AtlasV2

	federationSettingsID, federationSettingsIDOk := d.GetOk("federation_settings_id")

	if !federationSettingsIDOk {
		return diag.FromErr(errors.New("federation_settings_id must be configured"))
	}

	orgID, orgIDOk := d.GetOk("org_id")

	if !orgIDOk {
		return diag.FromErr(errors.New("org_id must be configured"))
	}
	federatedSettingsOrganizationRoleMappings, _, err := conn.FederatedAuthenticationApi.ListRoleMappings(ctx, federationSettingsID.(string), orgID.(string)).Execute()
	if err != nil {
		return diag.Errorf("error getting federatedSettings Role Mapping: assigned (%s): %s", federationSettingsID, err)
	}

	if err := d.Set("results", flattenFederatedSettingsOrganizationRoleMappings(federatedSettingsOrganizationRoleMappings)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `result` for federatedSettings Role Mapping:: %s", err))
	}

	d.SetId(federationSettingsID.(string))

	return nil
}

func flattenFederatedSettingsOrganizationRoleMappings(federatedSettingsOrganizationRoleMapping *admin.PaginatedRoleMapping) []map[string]any {
	var results []map[string]any
	mappings := federatedSettingsOrganizationRoleMapping.GetResults()
	for i := range mappings {
		mapping := &mappings[i]
		results = append(results, map[string]any{
			"external_group_name": mapping.GetExternalGroupName(),
			"id":                  mapping.GetId(),
			"role_assignments":    FlattenRoleAssignments(mapping.GetRoleAssignments()),
		})
	}
	return results
}
