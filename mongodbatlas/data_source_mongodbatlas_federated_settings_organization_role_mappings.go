package mongodbatlas

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasFederatedSettingsOrganizationRoleMappings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasFederatedSettingsOrganizationRoleMappingsRead,
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
func dataSourceMongoDBAtlasFederatedSettingsOrganizationRoleMappingsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas

	federationSettingsID, federationSettingsIDOk := d.GetOk("federation_settings_id")

	if !federationSettingsIDOk {
		return diag.FromErr(errors.New("federation_settings_id must be configured"))
	}

	orgID, orgIDOk := d.GetOk("org_id")

	if !orgIDOk {
		return diag.FromErr(errors.New("org_id must be configured"))
	}

	options := &matlas.ListOptions{
		PageNum:      d.Get("page_num").(int),
		ItemsPerPage: d.Get("items_per_page").(int),
	}

	federatedSettingsOrganizationRoleMappings, _, err := conn.FederatedSettings.ListRoleMappings(ctx, federationSettingsID.(string), orgID.(string), options)
	if err != nil {
		return diag.Errorf("error getting federatedSettings IdentityProviders assigned (%s): %s", federationSettingsID, err)
	}

	if err := d.Set("results", flattenFederatedSettingsOrganizationRoleMappings(federatedSettingsOrganizationRoleMappings)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `result` for federatedSettings IdentityProviders: %s", err))
	}

	d.SetId(federationSettingsID.(string))

	return nil
}

func flattenFederatedSettingsOrganizationRoleMappings(federatedSettingsOrganizationRoleMapping *matlas.FederatedSettingsOrganizationRoleMappings) []map[string]interface{} {
	var federatedSettingsOrganizationRoleMappingMap []map[string]interface{}

	if federatedSettingsOrganizationRoleMapping.TotalCount > 0 {
		federatedSettingsOrganizationRoleMappingMap = make([]map[string]interface{}, federatedSettingsOrganizationRoleMapping.TotalCount)

		for i := range federatedSettingsOrganizationRoleMapping.Results {
			federatedSettingsOrganizationRoleMappingMap[i] = map[string]interface{}{
				"external_group_name": federatedSettingsOrganizationRoleMapping.Results[i].ExternalGroupName,
				"id":                  federatedSettingsOrganizationRoleMapping.Results[i].ID,
				"role_assignments":    flattenRoleAssignments(federatedSettingsOrganizationRoleMapping.Results[i].RoleAssignments),
			}
		}
	}

	return federatedSettingsOrganizationRoleMappingMap
}
