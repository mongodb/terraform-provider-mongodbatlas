package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMongoDBAtlasFederatedSettingsOrganizationRoleMapping() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBAtlasFederatedSettingsOrganizationRoleMappingCreate,
		ReadContext:   resourceMongoDBAtlasFederatedSettingsOrganizationRoleMappingRead,
		UpdateContext: resourceMongoDBAtlasFederatedSettingsOrganizationRoleMappingUpdate,
		DeleteContext: resourceMongoDBAtlasFederatedSettingsOrganizationRoleMappingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasFederatedSettingsOrganizationRoleMappingImportState,
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
			"external_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"organization_roles": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_roles": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceMongoDBAtlasFederatedSettingsOrganizationRoleMappingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas

	ids := decodeStateID(d.Id())
	federationSettingsID := ids["federation_settings_id"]
	orgID := ids["org_id"]
	roleMappingID := ids["role_mapping_id"]

	federatedSettingsOrganizationRoleMapping, resp, err := conn.FederatedSettings.GetRoleMapping(context.Background(), federationSettingsID, orgID, roleMappingID)
	if err != nil {
		// case 404
		// deleted in the backend case
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf("error getting federated settings organization config: %s", err))
	}

	if err := d.Set("external_group_name", federatedSettingsOrganizationRoleMapping.ExternalGroupName); err != nil {
		return diag.FromErr(fmt.Errorf("error setting external group name (%s): %s", d.Id(), err))
	}

	orgRoles := []string{}
	groupRoles := []string{}

	for i := range federatedSettingsOrganizationRoleMapping.RoleAssignments {
		if federatedSettingsOrganizationRoleMapping.RoleAssignments[i].GroupID == "" {
			orgRoles = append(orgRoles, federatedSettingsOrganizationRoleMapping.RoleAssignments[i].Role)
		}

		if federatedSettingsOrganizationRoleMapping.RoleAssignments[i].OrgID == "" {
			groupRoles = append(groupRoles, federatedSettingsOrganizationRoleMapping.RoleAssignments[i].Role)
		}
	}

	if err := d.Set("organization_roles", orgRoles); err != nil {
		return diag.FromErr(fmt.Errorf("error setting org roles (%s): %s", d.Id(), err))
	}

	if err := d.Set("group_roles", groupRoles); err != nil {
		return diag.FromErr(fmt.Errorf("error setting group roles (%s): %s", d.Id(), err))
	}

	d.SetId(encodeStateID(map[string]string{
		"federation_settings_id": federationSettingsID,
		"org_id":                 orgID,
		"role_mapping_id":        roleMappingID,
	}))

	return nil
}

func resourceMongoDBAtlasFederatedSettingsOrganizationRoleMappingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	federationSettingsID := d.Get("federation_settings_id").(string)
	orgID := d.Get("org_id").(string)
	groupID := d.Get("group_id").(string)

	externalGroupName := d.Get("external_group_name").(string)

	body := &mongodbatlas.FederatedSettingsOrganizationRoleMapping{}
	body.ExternalGroupName = externalGroupName

	for _, role := range d.Get("organization_roles").(*schema.Set).List() {
		roleAssignment := mongodbatlas.RoleAssignments{}
		roleAssignment.Role = role.(string)

		roleAssignment.OrgID = orgID
		roleAssignment.GroupID = ""
		if roleAssignment.Role != "" {
			body.RoleAssignments = append(body.RoleAssignments, &roleAssignment)
		}
	}

	for _, role := range d.Get("group_roles").(*schema.Set).List() {
		roleAssignment := mongodbatlas.RoleAssignments{}
		roleAssignment.Role = role.(string)

		roleAssignment.OrgID = ""
		roleAssignment.GroupID = groupID
		if roleAssignment.Role != "" {
			body.RoleAssignments = append(body.RoleAssignments, &roleAssignment)
		}
	}

	federatedSettingsOrganizationRoleMapping, resp, err := conn.FederatedSettings.CreateRoleMapping(context.Background(), federationSettingsID, orgID, body)
	if err != nil {
		// case 404
		// deleted in the backend case
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf("error getting federated settings organization config: %s", err))
	}

	d.SetId(encodeStateID(map[string]string{
		"federation_settings_id": federationSettingsID,
		"org_id":                 orgID,
		"role_mapping_id":        federatedSettingsOrganizationRoleMapping.ID,
	}))

	return resourceMongoDBAtlasFederatedSettingsOrganizationRoleMappingRead(ctx, d, meta)
}

func resourceMongoDBAtlasFederatedSettingsOrganizationRoleMappingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	federationSettingsID := ids["federation_settings_id"]
	orgID := ids["org_id"]
	roleMappingID := ids["role_mapping_id"]
	groupID := d.Get("group_id").(string)

	federatedSettingsOrganizationRoleMappingUpdate, _, err := conn.FederatedSettings.GetRoleMapping(context.Background(), federationSettingsID, orgID, roleMappingID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retreiving federation settings connected organization (%s): %s", federationSettingsID, err))
	}

	if d.HasChange("external_group_name") {
		externalGroupName := d.Get("external_group_name").(string)
		federatedSettingsOrganizationRoleMappingUpdate.ExternalGroupName = externalGroupName
	}

	federatedSettingsOrganizationRoleMappingUpdate.RoleAssignments = nil

	for _, role := range d.Get("organization_roles").(*schema.Set).List() {
		roleAssignment := mongodbatlas.RoleAssignments{}
		roleAssignment.Role = role.(string)

		roleAssignment.OrgID = orgID
		roleAssignment.GroupID = ""
		if roleAssignment.Role != "" {
			federatedSettingsOrganizationRoleMappingUpdate.RoleAssignments = append(federatedSettingsOrganizationRoleMappingUpdate.RoleAssignments, &roleAssignment)
		}
	}

	for _, role := range d.Get("group_roles").(*schema.Set).List() {
		roleAssignment := mongodbatlas.RoleAssignments{}
		roleAssignment.Role = role.(string)

		roleAssignment.OrgID = ""
		roleAssignment.GroupID = groupID
		if roleAssignment.Role != "" {
			federatedSettingsOrganizationRoleMappingUpdate.RoleAssignments = append(federatedSettingsOrganizationRoleMappingUpdate.RoleAssignments, &roleAssignment)
		}
	}

	_, _, err = conn.FederatedSettings.UpdateRoleMapping(ctx, federationSettingsID, orgID, roleMappingID, federatedSettingsOrganizationRoleMappingUpdate)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating federation settings connected organization (%s): %s", federationSettingsID, err))
	}

	return resourceMongoDBAtlasFederatedSettingsOrganizationRoleMappingRead(ctx, d, meta)
}

func resourceMongoDBAtlasFederatedSettingsOrganizationRoleMappingDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	federationSettingsID := ids["federation_settings_id"]
	orgID := ids["org_id"]
	roleMappingID := ids["role_mapping_id"]

	_, err := conn.FederatedSettings.DeleteRoleMapping(ctx, federationSettingsID, orgID, roleMappingID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting federation settings connected organization (%s): %s", federationSettingsID, err))
	}

	return nil
}

func resourceMongoDBAtlasFederatedSettingsOrganizationRoleMappingImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*MongoDBClient).Atlas

	federationSettingsID, orgID, roleMappingID, err := splitFederatedSettingsOrganizationRoleMappingImportID(d.Id())
	if err != nil {
		return nil, err
	}

	federatedSettingsOrganizationRoleMapping, _, err := conn.FederatedSettings.GetRoleMapping(context.Background(), *federationSettingsID, *orgID, *roleMappingID)
	if err != nil {
		return nil, fmt.Errorf("couldn't import Role Mappings (%s) in Federation settings (%s), error: %s", *roleMappingID, *federationSettingsID, err)
	}

	if err := d.Set("federation_settings_id", *federationSettingsID); err != nil {
		return nil, fmt.Errorf("error setting role mapping in Federation settings (%s): %s", d.Id(), err)
	}

	orgRoles := []string{}
	groupRoles := []string{}

	for i := range federatedSettingsOrganizationRoleMapping.RoleAssignments {
		if federatedSettingsOrganizationRoleMapping.RoleAssignments[i].GroupID == "" {
			if err := d.Set("org_id", federatedSettingsOrganizationRoleMapping.RoleAssignments[i].OrgID); err != nil {
				return nil, fmt.Errorf("error setting org id (%s): %s", d.Id(), err)
			}
			orgRoles = append(orgRoles, federatedSettingsOrganizationRoleMapping.RoleAssignments[i].Role)
		}

		if federatedSettingsOrganizationRoleMapping.RoleAssignments[i].OrgID == "" {
			if err := d.Set("group_id", federatedSettingsOrganizationRoleMapping.RoleAssignments[i].GroupID); err != nil {
				return nil, fmt.Errorf("error setting group id  (%s): %s", d.Id(), err)
			}
			groupRoles = append(groupRoles, federatedSettingsOrganizationRoleMapping.RoleAssignments[i].Role)
		}
	}

	if err := d.Set("organization_roles", orgRoles); err != nil {
		return nil, fmt.Errorf("error setting org roles (%s): %s", d.Id(), err)
	}

	if err := d.Set("group_roles", groupRoles); err != nil {
		return nil, fmt.Errorf("error setting group roles (%s): %s", d.Id(), err)
	}

	d.SetId(encodeStateID(map[string]string{
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
