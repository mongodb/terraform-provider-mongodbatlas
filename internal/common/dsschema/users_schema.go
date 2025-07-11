package dsschema

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DSOrgUsersSchema() *schema.Schema {
	return &schema.Schema{
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
								Type:     schema.TypeSet,
								Computed: true,
								Elem:     &schema.Schema{Type: schema.TypeString},
							},
							"project_roles_assignments": {
								Type:     schema.TypeSet,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"project_id": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"project_roles": {
											Type:     schema.TypeSet,
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
}
