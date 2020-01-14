package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"

	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

const (
	errorTeamCreate        = "error creating Team information: %s"
	errorTeamAddUsers      = "error adding users to the Team information: %s"
	errorTeamRead          = "error getting Team information: %s"
	errorTeamUpdate        = "error updating Team information: %s"
	errorTeamUpdatingRoles = "error updating Team Roles information: %s"
	errorTeamDelete        = "error deleting Team (%s): %s"
	errorTeamSetting       = "error setting `%s` for Team (%s): %s"
)

func resourceMongoDBAtlasTeam() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoDBAtlasTeamCreate,
		Read:   resourceMongoDBAtlasTeamRead,
		Update: resourceMongoDBAtlasTeamUpdate,
		Delete: resourceMongoDBAtlasTeamDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMongoDBAtlasTeamImportState,
		},
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"team_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"usernames": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"team_roles": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceMongoDBAtlasTeamCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)
	orgID := d.Get("org_id").(string)
	projectID := d.Get("project_id").(string)

	// Creating the team
	teamsResp, _, err := conn.Teams.Create(context.Background(), orgID,
		&matlas.Team{
			Name:      d.Get("name").(string),
			Usernames: expandStringListFromSetSchema(d.Get("usernames").(*schema.Set)),
		})
	if err != nil {
		return fmt.Errorf(errorTeamCreate, err)
	}

	// Linking the team inside of one Project and adding roles if roles isn't empty
	_, _, err = conn.Projects.AddTeamsToProject(context.Background(), projectID,
		[]*matlas.ProjectTeam{{
			TeamID:    teamsResp.ID,
			RoleNames: expandStringListFromSetSchema(d.Get("team_roles").(*schema.Set)),
		}},
	)
	if err != nil {
		return fmt.Errorf("error linking the team(%s) to the Project(%s) information: %s", teamsResp.ID, projectID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"org_id":     orgID,
		"id":         teamsResp.ID,
		"project_id": projectID,
	}))

	return resourceMongoDBAtlasTeamRead(d, meta)
}

func resourceMongoDBAtlasTeamRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	ids := decodeStateID(d.Id())
	orgID := ids["org_id"]
	teamID := ids["id"]
	projectID := ids["project_id"]

	team, _, err := conn.Teams.Get(context.Background(), orgID, teamID)
	if err != nil {
		return fmt.Errorf(errorTeamRead, err)
	}

	if err := d.Set("name", team.Name); err != nil {
		return fmt.Errorf(errorTeamSetting, "name", teamID, err)
	}
	if err := d.Set("team_id", team.ID); err != nil {
		return fmt.Errorf(errorTeamSetting, "team_id", teamID, err)
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
		return fmt.Errorf(errorTeamSetting, "usernames", teamID, err)
	}

	// Get all the teams assiged to a Project
	teams, _, err := conn.Projects.GetProjectTeamsAssigned(context.Background(), projectID)
	if err != nil {
		return fmt.Errorf("error getting Teams from a Project(%s) information: %s", orgID, err)
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

	return nil
}

func resourceMongoDBAtlasTeamUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	ids := decodeStateID(d.Id())
	orgID := ids["org_id"]
	teamID := ids["id"]
	projectID := ids["project_id"]

	if d.HasChange("name") {
		_, _, err := conn.Teams.Rename(context.Background(), orgID, teamID, d.Get("name").(string))
		if err != nil {
			return fmt.Errorf(errorTeamUpdate, err)
		}
	}

	if d.HasChange("usernames") {
		// First, we need to remove the current users of the team and later add the new users
		// Get the current team's users
		users, _, err := conn.Teams.GetTeamUsersAssigned(context.Background(), orgID, teamID)
		if err != nil {
			return fmt.Errorf(errorTeamRead, err)
		}
		// Removing each user
		for _, u := range users {
			_, err := conn.Teams.RemoveUserToTeam(context.Background(), orgID, teamID, u.ID)
			if err != nil {
				return fmt.Errorf("error deleting Atlas User (%s) information: %s", teamID, err)
			}
		}

		// Verify if the gave users exists
		var newUsers []string
		for _, username := range d.Get("usernames").(*schema.Set).List() {
			user, _, err := conn.AtlasUsers.GetByName(context.Background(), username.(string))
			if err != nil {
				return fmt.Errorf("error getting Atlas User (%s) information: %s", username, err)
			}
			// if the user exists, we will storage its teamID
			newUsers = append(newUsers, user.ID)
		}

		// Adding the new existing users by teamID
		_, _, err = conn.Teams.AddUsersToTeam(context.Background(), orgID, teamID, newUsers)
		if err != nil {
			return fmt.Errorf(errorTeamAddUsers, err)
		}
	}

	if d.HasChange("team_roles") {
		_, _, err := conn.Teams.UpdateTeamRoles(context.Background(), projectID, teamID,
			&matlas.TeamUpdateRoles{
				RoleNames: expandStringListFromSetSchema(d.Get("team_roles").(*schema.Set)),
			})
		if err != nil {
			return fmt.Errorf(errorTeamUpdatingRoles, err)
		}
	}

	return resourceMongoDBAtlasTeamRead(d, meta)
}

func resourceMongoDBAtlasTeamDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	orgID := ids["org_id"]
	id := ids["id"]
	projectID := ids["project_id"]

	_, err := conn.Teams.RemoveTeamFromProject(context.Background(), projectID, id)
	if err != nil {
		return fmt.Errorf(errorTeamDelete, id, err)
	}

	_, err = conn.Teams.RemoveTeamFromOrganization(context.Background(), orgID, id)
	if err != nil {
		return fmt.Errorf(errorTeamDelete, id, err)
	}

	return nil
}

func resourceMongoDBAtlasTeamImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*matlas.Client)

	parts := strings.SplitN(d.Id(), "-", 3)
	if len(parts) != 3 {
		return nil, errors.New("import format error: to import a team, use the format {group_id}-{team_id}-{project_id}")
	}

	orgID := parts[0]
	teamID := parts[1]
	projectID := parts[2]

	u, _, err := conn.Teams.Get(context.Background(), orgID, teamID)
	if err != nil {
		return nil, fmt.Errorf("couldn't import team (%s) in project (%s), error: %s", teamID, orgID, err)
	}

	if err := d.Set("org_id", orgID); err != nil {
		log.Printf("[WARN] Error setting org_id for (%s): %s", teamID, err)
	}
	if err := d.Set("team_id", teamID); err != nil {
		log.Printf("[WARN] Error setting team_id for (%s): %s", teamID, err)
	}
	if err := d.Set("project_id", projectID); err != nil {
		log.Printf("[WARN] Error setting project_id for (%s): %s", teamID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"org_id":     orgID,
		"id":         u.ID,
		"project_id": projectID,
	}))

	return []*schema.ResourceData{d}, nil
}

func expandStringListFromSetSchema(list *schema.Set) []string {
	res := make([]string, list.Len())
	for i, v := range list.List() {
		res[i] = v.(string)
	}
	return res
}
