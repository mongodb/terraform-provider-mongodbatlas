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
			"project_id": {
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
			"usernames": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"team_roles": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceMongoDBAtlasTeamRead(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	orgID := d.Get("org_id").(string)
	teamID := d.Get("team_id").(string)
	projectID := d.Get("project_id").(string)

	team, _, err := conn.Teams.Get(context.Background(), orgID, teamID)
	if err != nil {
		return fmt.Errorf(errorTeamRead, err)
	}

	if err := d.Set("name", team.Name); err != nil {
		return fmt.Errorf(errorTeamSetting, "name", d.Id(), err)
	}

	//Set Usernames
	users, _, err := conn.Teams.GetTeamUsersAssigned(context.Background(), orgID, teamID)
	if err != nil {
		return fmt.Errorf(errorTeamRead, err)
	}

	var usernames []string
	for _, u := range users {
		usernames = append(usernames, u.Username)
	}

	if err := d.Set("usernames", usernames); err != nil {
		return fmt.Errorf(errorTeamSetting, "usernames", d.Id(), err)
	}

	// Get all the teams assiged to a Project
	teams, _, err := conn.Projects.GetProjectTeamsAssigned(context.Background(), projectID)
	if err != nil {
		return fmt.Errorf(errorTeamRead, err)
	}

	var teamRoles []string
	for _, team := range teams.Results { // looking for our current team
		// If the team exists then save the roles
		if team.TeamID == teamID {
			teamRoles = team.RoleNames
			break
		}
	}

	// Set the roles
	if err := d.Set("team_roles", teamRoles); err != nil {
		return fmt.Errorf(errorTeamSetting, "team_roles", teamID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"org_id":     orgID,
		"id":         team.ID,
		"project_id": projectID,
	}))

	return nil
}
