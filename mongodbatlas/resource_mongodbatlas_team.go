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
	errorTeamCreate   = "error creating Team information: %s"
	errorTeamAddUsers = "error adding users to the Team information: %s"
	errorTeamRead     = "error getting Team information: %s"
	errorTeamUpdate   = "error updating Team information: %s"
	errorTeamDelete   = "error deleting Team (%s): %s"
	errorTeamSetting  = "error setting `%s` for Team (%s): %s"
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
		},
	}
}

func resourceMongoDBAtlasTeamCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)
	orgID := d.Get("org_id").(string)

	teamReq := &matlas.Team{
		Name:      d.Get("name").(string),
		Usernames: expandStringListFromSetSchema(d.Get("usernames").(*schema.Set)),
	}

	teamsResp, _, err := conn.Teams.Create(context.Background(), orgID, teamReq)
	if err != nil {
		return fmt.Errorf(errorTeamCreate, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"org_id": orgID,
		"id":     teamsResp.ID,
	}))

	return resourceMongoDBAtlasTeamRead(d, meta)
}

func resourceMongoDBAtlasTeamRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	ids := decodeStateID(d.Id())
	orgID := ids["org_id"]
	id := ids["id"]

	team, _, err := conn.Teams.Get(context.Background(), orgID, id)
	if err != nil {
		return fmt.Errorf(errorTeamRead, err)
	}

	if err := d.Set("name", team.Name); err != nil {
		return fmt.Errorf(errorTeamSetting, "name", d.Id(), err)
	}
	if err := d.Set("team_id", team.ID); err != nil {
		return fmt.Errorf(errorTeamSetting, "team_id", d.Id(), err)
	}

	//Set Usernames
	users, _, err := conn.Teams.GetTeamUsersAssigned(context.Background(), orgID, id)
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

	return nil
}

func resourceMongoDBAtlasTeamUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	ids := decodeStateID(d.Id())
	orgID := ids["org_id"]
	id := ids["id"]

	if d.HasChange("name") {
		_, _, err := conn.Teams.Rename(context.Background(), orgID, id, d.Get("name").(string))
		if err != nil {
			return fmt.Errorf(errorTeamUpdate, err)
		}
	}

	if d.HasChange("usernames") {
		// First, we need to remove the current users of the team and later add the new users
		// Get the current team's users
		users, _, err := conn.Teams.GetTeamUsersAssigned(context.Background(), orgID, id)
		if err != nil {
			return fmt.Errorf(errorTeamRead, err)
		}
		// Removing each user
		for _, u := range users {
			_, err := conn.Teams.RemoveUserToTeam(context.Background(), orgID, id, u.ID)
			if err != nil {
				return fmt.Errorf("error deleting Atlas User (%s) information: %s", id, err)
			}
		}

		// Verify if the gave users exists
		var newUsers []string
		for _, username := range d.Get("usernames").(*schema.Set).List() {
			user, _, err := conn.AtlasUsers.GetByName(context.Background(), username.(string))
			if err != nil {
				return fmt.Errorf("error getting Atlas User (%s) information: %s", username, err)
			}
			// if the user exists, we will storage its id
			newUsers = append(newUsers, user.ID)
		}

		// Adding the new existing users by id
		_, _, err = conn.Teams.AddUsersToTeam(context.Background(), orgID, id, newUsers)
		if err != nil {
			return fmt.Errorf(errorTeamAddUsers, err)
		}
	}

	return resourceMongoDBAtlasTeamRead(d, meta)
}

func resourceMongoDBAtlasTeamDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	orgID := ids["org_id"]
	id := ids["id"]

	_, err := conn.Teams.RemoveTeamFromOrganization(context.Background(), orgID, id)
	if err != nil {
		return fmt.Errorf(errorTeamDelete, id, err)
	}

	return nil
}

func resourceMongoDBAtlasTeamImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*matlas.Client)

	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 2 {
		return nil, errors.New("import format error: to import a team, use the format {group_id}-{team_id}")
	}

	orgID := parts[0]
	teamID := parts[1]

	u, _, err := conn.Teams.Get(context.Background(), orgID, teamID)
	if err != nil {
		return nil, fmt.Errorf("couldn't import user (%s) in project (%s), error: %s", teamID, orgID, err)
	}

	if err := d.Set("org_id", orgID); err != nil {
		log.Printf("[WARN] Error setting org_id for (%s): %s", d.Id(), err)
	}

	d.SetId(encodeStateID(map[string]string{
		"org_id": orgID,
		"id":     u.ID,
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
