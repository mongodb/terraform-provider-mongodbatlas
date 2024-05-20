package federatedsettingsorgrolemapping

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"federation_settings_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"org_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role_mapping_id": {
				Type:     schema.TypeString,
				Required: true,
			},

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
	}
}
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2

	federationSettingsID, federationSettingsIDOk := d.GetOk("federation_settings_id")

	if !federationSettingsIDOk {
		return diag.FromErr(errors.New("federation_settings_id must be configured"))
	}

	orgID, orgIDOk := d.GetOk("org_id")

	if !orgIDOk {
		return diag.FromErr(errors.New("org_id must be configured"))
	}

	roleMappingID, roleMappingOk := d.GetOk("role_mapping_id")

	if !roleMappingOk {
		return diag.FromErr(errors.New("role_mapping_id must be configured"))
	}

	federatedSettingsOrganizationRoleMapping, _, err := conn.FederatedAuthenticationApi.GetRoleMapping(ctx, federationSettingsID.(string), roleMappingID.(string), orgID.(string)).Execute()
	if err != nil {
		return diag.Errorf("error getting federatedSettings Role Mapping assigned (%s): %s", federationSettingsID, err)
	}

	if err := d.Set("external_group_name", federatedSettingsOrganizationRoleMapping.GetExternalGroupName()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `result` for federatedSettings Role Mapping: %s", err))
	}

	if err := d.Set("role_assignments", FlattenRoleAssignments(federatedSettingsOrganizationRoleMapping.GetRoleAssignments())); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `result` for federatedSettings Role Mapping: %s", err))
	}

	d.SetId(federatedSettingsOrganizationRoleMapping.GetId())

	return nil
}
