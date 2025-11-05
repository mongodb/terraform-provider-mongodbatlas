package federatedsettingsorgrolemapping

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312009/admin"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate,
		ReadContext:   resourceRead,
		UpdateContext: resourceUpdate,
		DeleteContext: resourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceImportState,
		},
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
				Computed: true,
			},
			"external_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role_assignments": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"org_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"roles": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2

	ids := conversion.DecodeStateID(d.Id())
	federationSettingsID := ids["federation_settings_id"]
	orgID := ids["org_id"]
	roleMappingID := ids["role_mapping_id"]

	federatedSettingsOrganizationRoleMapping, resp, err := conn.FederatedAuthenticationApi.GetRoleMapping(context.Background(), federationSettingsID, roleMappingID, orgID).Execute()

	if err != nil {
		if validate.StatusNotFound(resp) {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf("error getting federated settings organization config: %s", err))
	}

	if err := d.Set("role_mapping_id", federatedSettingsOrganizationRoleMapping.GetId()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting role_mapping_id (%s): %s", d.Id(), err))
	}

	if err := d.Set("external_group_name", federatedSettingsOrganizationRoleMapping.GetExternalGroupName()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting external group name (%s): %s", d.Id(), err))
	}

	if err := d.Set("role_assignments", flattenRoleAssignmentsResource(federatedSettingsOrganizationRoleMapping.GetRoleAssignments())); err != nil {
		return diag.FromErr(fmt.Errorf("error setting role_assignments (%s): %s", d.Id(), err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"federation_settings_id": federationSettingsID,
		"org_id":                 orgID,
		"role_mapping_id":        roleMappingID,
	}))

	return nil
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2
	federationSettingsID, federationSettingsIDOk := d.GetOk("federation_settings_id")

	if !federationSettingsIDOk {
		return diag.FromErr(errors.New("federation_settings_id must be configured"))
	}

	orgID, orgIDOk := d.GetOk("org_id")

	if !orgIDOk {
		return diag.FromErr(errors.New("org_id must be configured"))
	}

	externalGroupName := d.Get("external_group_name").(string)

	roleMapping := admin.AuthFederationRoleMapping{
		ExternalGroupName: externalGroupName,
		RoleAssignments:   expandRoleAssignments(d),
	}
	federatedSettingsOrganizationRoleMapping, resp, err := conn.FederatedAuthenticationApi.CreateRoleMapping(ctx, federationSettingsID.(string), orgID.(string), &roleMapping).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf("error getting federated settings organization config: %s", err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"federation_settings_id": federationSettingsID.(string),
		"org_id":                 orgID.(string),
		"role_mapping_id":        federatedSettingsOrganizationRoleMapping.GetId(),
	}))

	return resourceRead(ctx, d, meta)
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	federationSettingsID := ids["federation_settings_id"]
	orgID := ids["org_id"]
	roleMappingID := ids["role_mapping_id"]

	federatedSettingsOrganizationRoleMappingUpdate, _, err := conn.FederatedAuthenticationApi.GetRoleMapping(context.Background(), federationSettingsID, roleMappingID, orgID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retreiving federation settings connected organization (%s): %s", federationSettingsID, err))
	}

	if d.HasChange("external_group_name") {
		externalGroupName := d.Get("external_group_name").(string)
		federatedSettingsOrganizationRoleMappingUpdate.ExternalGroupName = externalGroupName
	}

	if d.HasChange("role_assignments") {
		federatedSettingsOrganizationRoleMappingUpdate.RoleAssignments = expandRoleAssignments(d)
	}
	_, _, err = conn.FederatedAuthenticationApi.UpdateRoleMapping(ctx, federationSettingsID, roleMappingID, orgID, federatedSettingsOrganizationRoleMappingUpdate).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating federation settings connected organization (%s): %s", federationSettingsID, err))
	}

	return resourceRead(ctx, d, meta)
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	federationSettingsID := ids["federation_settings_id"]
	orgID := ids["org_id"]
	roleMappingID := ids["role_mapping_id"]

	_, err := conn.FederatedAuthenticationApi.DeleteRoleMapping(ctx, federationSettingsID, roleMappingID, orgID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting federation settings connected organization (%s): %s", federationSettingsID, err))
	}

	return nil
}

func resourceImportState(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	conn := meta.(*config.MongoDBClient).AtlasV2

	federationSettingsID, orgID, roleMappingID, err := splitFederatedSettingsOrganizationRoleMappingImportID(d.Id())
	if err != nil {
		return nil, err
	}

	_, _, err = conn.FederatedAuthenticationApi.GetRoleMapping(context.Background(), *federationSettingsID, *roleMappingID, *orgID).Execute()
	if err != nil {
		return nil, fmt.Errorf("couldn't import Role Mappings (%s) in Federation settings (%s), error: %s", *roleMappingID, *federationSettingsID, err)
	}

	if err := d.Set("federation_settings_id", *federationSettingsID); err != nil {
		return nil, fmt.Errorf("error setting federation_settings_id for role mapping in Federation settings (%s): %s", d.Id(), err)
	}

	if err := d.Set("org_id", *orgID); err != nil {
		return nil, fmt.Errorf("error setting org_id for role mapping in Federation settings (%s): %s", d.Id(), err)
	}
	d.SetId(conversion.EncodeStateID(map[string]string{
		"federation_settings_id": *federationSettingsID,
		"org_id":                 *orgID,
		"role_mapping_id":        *roleMappingID,
	}))

	return []*schema.ResourceData{d}, nil
}

func splitFederatedSettingsOrganizationRoleMappingImportID(id string) (federationSettingsID, orgID, roleMappingID *string, err error) {
	var re = regexp.MustCompile(`(?s)^(.*)-(.*)-(.*)$`)
	parts := re.FindStringSubmatch(id)

	if len(parts) != 4 {
		err = errors.New("import format error: to import a Federated Settings Role Mappings, use the format {federation_settings_id}-{org_id}-{role_mapping_id}")
		return
	}

	federationSettingsID = &parts[1]
	orgID = &parts[2]
	roleMappingID = &parts[3]

	return
}
