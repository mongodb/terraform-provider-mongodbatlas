package mongodbatlas

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	matlas "go.mongodb.org/atlas/mongodbatlas"
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
				Type:          schema.TypeString,
				Computed:      true,
				Optional:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Computed:      true,
				Optional:      true,
				ConflictsWith: []string{"team_id"},
			},
			"usernames": {
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
	var (
		conn             = meta.(*matlas.Client)
		orgID            = d.Get("org_id").(string)
		teamID, teamIDOk = d.GetOk("team_id")
		name, nameOk     = d.GetOk("name")

		err  error
		team *matlas.Team
	)

	if !teamIDOk && !nameOk {
		return errors.New("either team_id or name must be configured")
	}

	if teamIDOk {
		team, _, err = conn.Teams.Get(context.Background(), orgID, teamID.(string))
	} else {
		team, _, err = conn.Teams.GetOneTeamByName(context.Background(), orgID, name.(string))
	}

	if err != nil {
		return fmt.Errorf(errorTeamRead, err)
	}

	if err := d.Set("team_id", team.ID); err != nil {
		return fmt.Errorf(errorTeamSetting, "name", d.Id(), err)
	}

	if err := d.Set("name", team.Name); err != nil {
		return fmt.Errorf(errorTeamSetting, "name", d.Id(), err)
	}

	// Set Usernames
	users, _, err := conn.Teams.GetTeamUsersAssigned(context.Background(), orgID, team.ID)
	if err != nil {
		return fmt.Errorf(errorTeamRead, err)
	}

	usernames := []string{}
	for i := range users {
		usernames = append(usernames, users[i].Username)
	}

	if err := d.Set("usernames", usernames); err != nil {
		return fmt.Errorf(errorTeamSetting, "usernames", d.Id(), err)
	}

	d.SetId(encodeStateID(map[string]string{
		"org_id": orgID,
		"id":     team.ID,
	}))

	return nil
}
