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
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceMongoDBAtlasTeamCreate(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	orgID := d.Get("org_id").(string)

	teamReq := &matlas.Team{
		Name: d.Get("name").(string),
	}

	if v, ok := d.GetOk("usernames"); ok {
		teamReq.Usernames = expandStringList(v.([]interface{}))
	}

	teamsResp, _, err := conn.Teams.Create(context.Background(), orgID, teamReq)

	if err != nil {
		return fmt.Errorf("error creating team: %s", err)
	}

	if err := d.Set("team_id", teamsResp.ID); err != nil {
		return fmt.Errorf("error setting `id` for team (%s): %s", d.Id(), err)
	}

	d.SetId(encodeStateID(map[string]string{
		"org_id": orgID,
		"id":     teamsResp.ID,
	}))

	return resourceMongoDBAtlasTeamRead(d, meta)
}

func resourceMongoDBAtlasTeamRead(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	orgID := ids["org_id"]
	id := ids["id"]

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

	var usernames []string
	for _, v := range users {
		usernames = append(usernames, v.Username)
	}

	if err := d.Set("usernames", usernames); err != nil {
		return fmt.Errorf("error setting `usernames` for team (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceMongoDBAtlasTeamUpdate(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	orgID := ids["org_id"]
	id := ids["id"]

	if d.HasChange("name") {
		_, _, err := conn.Teams.Rename(context.Background(), orgID, id, d.Get("name").(string))

		if err != nil {
			return fmt.Errorf("error updating team(%s): %s", id, err)
		}
	}

	return resourceMongoDBAtlasTeamRead(d, meta)
}

func resourceMongoDBAtlasTeamDelete(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	orgID := ids["org_id"]
	id := ids["id"]

	_, err := conn.Teams.RemoveTeamFromOrganization(context.Background(), orgID, id)

	if err != nil {
		return fmt.Errorf("error deleting team (%s): %s", id, err)
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
	team_id := parts[1]

	u, _, err := conn.Teams.Get(context.Background(), orgID, team_id)
	if err != nil {
		return nil, fmt.Errorf("couldn't import user %s in project %s, error: %s", team_id, orgID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"org_id": orgID,
		"id":     u.ID,
	}))

	if err := d.Set("team_id", u.ID); err != nil {
		log.Printf("[WARN] Error setting team_id for (%s): %s", d.Id(), err)
	}

	if err := d.Set("org_id", orgID); err != nil {
		log.Printf("[WARN] Error setting org_id for (%s): %s", d.Id(), err)
	}

	return []*schema.ResourceData{d}, nil
}

func expandStringList(list []interface{}) []string {
	var sList []string

	for _, v := range list {
		sList = append(sList, v.(string))
	}

	return sList
}
