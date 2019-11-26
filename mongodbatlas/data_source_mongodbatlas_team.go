package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"

	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasTeam() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasTeamRead,
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"team_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"users": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"first_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"last_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"roles": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"project_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"org_id": {
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
						"team_ids": {
							Type:     schema.TypeList,
							Computed: true,
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

func dataSourceMongoDBAtlasTeamRead(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	orgID := d.Get("org_id").(string)
	id := d.Get("team_id").(string)

	team, _, err := conn.Teams.Get(context.Background(), orgID, id)

	if err != nil {
		return fmt.Errorf("error getting team information: %s", err)
	}

	if err := d.Set("name", team.Name); err != nil {
		return fmt.Errorf("error setting `name` for team (%s): %s", d.Id(), err)
	}

	//Set Usernames
	users, _, err := conn.Teams.GetTeamUsersAssigned(context.Background(), orgID, id)

	if err != nil {
		return fmt.Errorf("error getting team user assigned information: %s", err)
	}

	if err := d.Set("users", flattenAtlasUsers(users)); err != nil {
		return fmt.Errorf("error setting `users` for team (%s): %s", d.Id(), err)
	}

	d.SetId(team.ID)

	return nil
}

func flattenAtlasUsers(users []matlas.AtlasUser) []map[string]interface{} {
	var atlasUsers []map[string]interface{}

	for _, user := range users {
		atlasUser := map[string]interface{}{
			"id":            user.ID,
			"email_address": user.EmailAddress,
			"first_name":    user.FirstName,
			"last_name":     user.LastName,
			"roles":         flattenAtlasRoles(user.Roles),
			"team_ids":      user.TeamIds,
			"username":      user.Username,
		}

		atlasUsers = append(atlasUsers, atlasUser)
	}

	return atlasUsers
}

func flattenAtlasRoles(roles []matlas.AtlasRole) []map[string]interface{} {
	var atlasRoles []map[string]interface{}

	for _, role := range roles {
		atlasRole := map[string]interface{}{
			"project_id": role.GroupID,
			"org_id":     role.OrgID,
			"role_name":  role.RoleName,
		}
		atlasRoles = append(atlasRoles, atlasRole)
	}
	return atlasRoles
}
