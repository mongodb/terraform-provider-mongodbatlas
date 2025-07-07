package organization

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"go.mongodb.org/atlas-sdk/v20250312005/admin"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/dsschema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func FlattenUsers(users []admin.OrgUserResponse) []map[string]any {
	ret := make([]map[string]any, len(users))
	for i := range users {
		user := &users[i]
		ret[i] = map[string]any{
			"id":                    user.GetId(),
			"org_membership_status": user.GetOrgMembershipStatus(),
			"roles":                 FlattenUserRoles(user.GetRoles()),
			"team_ids":              user.GetTeamIds(),
			"username":              user.GetUsername(),
			"invitation_created_at": user.GetInvitationCreatedAt().Format(time.RFC3339),
			"invitation_expires_at": user.GetInvitationExpiresAt().Format(time.RFC3339),
			"inviter_username":      user.GetInviterUsername(),
			"country":               user.GetCountry(),
			"created_at":            user.GetCreatedAt().Format(time.RFC3339),
			"first_name":            user.GetFirstName(),
			"last_auth":             user.GetLastAuth().Format(time.RFC3339),
			"last_name":             user.GetLastName(),
			"mobile_number":         user.GetMobileNumber(),
		}
	}
	return ret
}

func FlattenUserRoles(roles admin.OrgUserRolesResponse) []map[string]any {
	ret := []map[string]any{}
	roleMap := map[string]any{
		"org_roles":     []string{},
		"project_roles": []map[string]any{},
	}
	if roles.HasOrgRoles() {
		roleMap["org_roles"] = roles.GetOrgRoles()
	}
	if roles.HasGroupRoleAssignments() {
		roleMap["project_roles"] = FlattenGroupRolesAssignments(roles.GetGroupRoleAssignments())
	}
	ret = append(ret, roleMap)
	return ret
}

func FlattenGroupRolesAssignments(assignments []admin.GroupRoleAssignment) []map[string]any {
	ret := make([]map[string]any, len(assignments))
	for i, assignment := range assignments {
		ret[i] = map[string]any{
			"group_id":    assignment.GetGroupId(),
			"group_roles": assignment.GetGroupRoles(),
		}
	}
	return ret
}

var (
	DSOrgUsersSchema = schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"org_membership_status": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"roles": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"org_roles": {
								Type:     schema.TypeList,
								Computed: true,
								Elem:     &schema.Schema{Type: schema.TypeString},
							},
							"project_roles": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"group_id": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"group_roles": {
											Type:     schema.TypeList,
											Computed: true,
											Elem:     &schema.Schema{Type: schema.TypeString},
										},
									},
								},
							},
						},
					},
				},
				"team_ids": {
					Type:     schema.TypeList,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"username": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"invitation_created_at": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"invitation_expires_at": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"inviter_username": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"country": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"created_at": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"first_name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"last_auth": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"last_name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"mobile_number": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_deleted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rel": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"users": &DSOrgUsersSchema,
			"api_access_list_required": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"multi_factor_auth_required": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"restrict_employee_access": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"gen_ai_features_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"skip_default_alerts_settings": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"security_contact": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2
	orgID := d.Get("org_id").(string)

	organization, _, err := conn.OrganizationsApi.GetOrganization(ctx, orgID).Execute()

	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting organizations information: %s", err))
	}

	if err := d.Set("name", organization.GetName()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `name`: %s", err))
	}

	if err := d.Set("skip_default_alerts_settings", organization.GetSkipDefaultAlertsSettings()); err != nil {
		return diag.Errorf("error setting `skip_default_alerts_settings`: %s", err)
	}

	if err := d.Set("is_deleted", organization.GetIsDeleted()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `is_deleted`: %s", err))
	}

	if err := d.Set("links", conversion.FlattenLinks(organization.GetLinks())); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `is_deleted`: %s", err))
	}

	users, err := ListAllOrganizationUsers(ctx, orgID, conn)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting organization users: %s", err))
	}
	if err := d.Set("users", FlattenUsers(users)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `users`: %s", err))
	}

	settings, _, err := conn.OrganizationsApi.GetOrganizationSettings(ctx, orgID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting organization settings: %s", err))
	}
	if err := d.Set("api_access_list_required", settings.ApiAccessListRequired); err != nil {
		return diag.Errorf("error setting `api_access_list_required` for organization (%s): %s", orgID, err)
	}
	if err := d.Set("multi_factor_auth_required", settings.MultiFactorAuthRequired); err != nil {
		return diag.Errorf("error setting `multi_factor_auth_required` for organization (%s): %s", orgID, err)
	}
	if err := d.Set("restrict_employee_access", settings.RestrictEmployeeAccess); err != nil {
		return diag.Errorf("error setting `restrict_employee_access` for organization (%s): %s", orgID, err)
	}
	if err := d.Set("gen_ai_features_enabled", settings.GenAIFeaturesEnabled); err != nil {
		return diag.Errorf("error setting `gen_ai_features_enabled` for organization (%s): %s", orgID, err)
	}
	if err := d.Set("security_contact", settings.SecurityContact); err != nil {
		return diag.Errorf("error setting `security_contact` for organization (%s): %s", orgID, err)
	}

	d.SetId(organization.GetId())

	return nil
}

func ListAllOrganizationUsers(ctx context.Context, orgID string, conn *admin.APIClient) ([]admin.OrgUserResponse, error) {
	return dsschema.AllPages(ctx, func(ctx context.Context, pageNum int) (dsschema.PaginateResponse[admin.OrgUserResponse], *http.Response, error) {
		request := conn.MongoDBCloudUsersApi.ListOrganizationUsers(ctx, orgID)
		request = request.PageNum(pageNum)
		return request.Execute()
	})
}
