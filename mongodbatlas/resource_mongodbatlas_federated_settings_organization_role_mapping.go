package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
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

	if err := d.Set("role_assignments", flattenRoleAssignmentsSpecial(federatedSettingsOrganizationRoleMapping.RoleAssignments)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting role_assignments (%s): %s", d.Id(), err))
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
	federationSettingsID, federationSettingsIDOk := d.GetOk("federation_settings_id")

	if !federationSettingsIDOk {
		return diag.FromErr(errors.New("federation_settings_id must be configured"))
	}

	orgID, orgIDOk := d.GetOk("org_id")

	if !orgIDOk {
		return diag.FromErr(errors.New("org_id must be configured"))
	}

	externalGroupName := d.Get("external_group_name").(string)

	body := &matlas.FederatedSettingsOrganizationRoleMapping{}

	ra := []*matlas.RoleAssignments{}

	body.ExternalGroupName = externalGroupName
	roleAssignments := expandRoleAssignments(d)

	for i := range roleAssignments {
		ra = append(ra, &roleAssignments[i])
	}

	body.RoleAssignments = ra

	federatedSettingsOrganizationRoleMapping, resp, err := conn.FederatedSettings.CreateRoleMapping(context.Background(), federationSettingsID.(string), orgID.(string), body)
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
		"federation_settings_id": federationSettingsID.(string),
		"org_id":                 orgID.(string),
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

	federatedSettingsOrganizationRoleMappingUpdate, _, err := conn.FederatedSettings.GetRoleMapping(context.Background(), federationSettingsID, orgID, roleMappingID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retreiving federation settings connected organization (%s): %s", federationSettingsID, err))
	}

	if d.HasChange("external_group_name") {
		externalGroupName := d.Get("external_group_name").(string)
		federatedSettingsOrganizationRoleMappingUpdate.ExternalGroupName = externalGroupName
	}

	if d.HasChange("role_assignments") {
		federatedSettingsOrganizationRoleMappingUpdate.RoleAssignments = nil

		ra := []*matlas.RoleAssignments{}

		roleAssignments := expandRoleAssignments(d)

		for i := range roleAssignments {
			ra = append(ra, &roleAssignments[i])
		}

		federatedSettingsOrganizationRoleMappingUpdate.RoleAssignments = ra
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

	if err := d.Set("org_id", *orgID); err != nil {
		return nil, fmt.Errorf("error setting role mapping in Federation settings (%s): %s", d.Id(), err)
	}

	if err := d.Set("role_assignments", flattenRoleAssignmentsSpecial(federatedSettingsOrganizationRoleMapping.RoleAssignments)); err != nil {
		return nil, fmt.Errorf("error setting role_assignments (%s): %s", d.Id(), err)
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

type roleAssignmentsByFields []matlas.RoleAssignments

func (ra roleAssignmentsByFields) Len() int      { return len(ra) }
func (ra roleAssignmentsByFields) Swap(i, j int) { ra[i], ra[j] = ra[j], ra[i] }

func (ra roleAssignmentsByFields) Less(i, j int) bool {
	compareVal := strings.Compare(ra[i].OrgID, ra[j].OrgID)

	if compareVal != 0 {
		return compareVal < 0
	}

	compareVal = strings.Compare(ra[i].GroupID, ra[j].GroupID)

	if compareVal != 0 {
		return compareVal < 0
	}

	return ra[i].Role < ra[j].Role
}

type roleAssignmentRefsByFields []*matlas.RoleAssignments

func (ra roleAssignmentRefsByFields) Len() int      { return len(ra) }
func (ra roleAssignmentRefsByFields) Swap(i, j int) { ra[i], ra[j] = ra[j], ra[i] }

func (ra roleAssignmentRefsByFields) Less(i, j int) bool {
	compareVal := strings.Compare(ra[i].OrgID, ra[j].OrgID)

	if compareVal != 0 {
		return compareVal < 0
	}

	compareVal = strings.Compare(ra[i].GroupID, ra[j].GroupID)

	if compareVal != 0 {
		return compareVal < 0
	}

	return ra[i].Role < ra[j].Role
}

func expandRoleAssignments(d *schema.ResourceData) []matlas.RoleAssignments {
	var roleAssignmentsReturn []matlas.RoleAssignments

	if v, ok := d.GetOk("role_assignments"); ok {
		if rs := v.(*schema.Set); rs.Len() > 0 {
			roleAssignments := []matlas.RoleAssignments{}
			roleAssignment := matlas.RoleAssignments{}

			for _, r := range rs.List() {
				roleMap := r.(map[string]interface{})

				for _, role := range roleMap["roles"].(*schema.Set).List() {
					roleAssignment.OrgID = roleMap["org_id"].(string)
					roleAssignment.GroupID = roleMap["group_id"].(string)
					roleAssignment.Role = role.(string)
					roleAssignments = append(roleAssignments, roleAssignment)
				}
				roleAssignmentsReturn = roleAssignments
			}
		}
	}

	sort.Sort(roleAssignmentsByFields(roleAssignmentsReturn))

	return roleAssignmentsReturn
}

func flattenRoleAssignmentsSpecial(roleAssignments []*matlas.RoleAssignments) []map[string]interface{} {
	if len(roleAssignments) == 0 {
		return nil
	}

	sort.Sort(roleAssignmentRefsByFields(roleAssignments))

	var flattenedRoleAssignments []map[string]interface{}
	var roleAssignment = map[string]interface{}{
		"group_id": roleAssignments[0].GroupID,
		"org_id":   roleAssignments[0].OrgID,
		"roles":    []string{},
	}

	for _, row := range roleAssignments {
		if (roleAssignment["org_id"] != "" && roleAssignment["org_id"] != row.OrgID) ||
			(roleAssignment["group_id"] != "" && roleAssignment["group_id"] != row.GroupID) {
			flattenedRoleAssignments = append(flattenedRoleAssignments, roleAssignment)

			roleAssignment = map[string]interface{}{
				"group_id": row.GroupID,
				"org_id":   row.OrgID,
				"roles":    []string{},
			}
		}

		roleAssignment["roles"] = append(roleAssignment["roles"].([]string), row.Role)
	}

	flattenedRoleAssignments = append(flattenedRoleAssignments, roleAssignment)

	return flattenedRoleAssignments
}
