package mongodbatlas

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasTeam() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasTeamRead,
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

func dataSourceMongoDBAtlasTeamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conn             = meta.(*MongoDBClient).Atlas
		orgID            = d.Get("org_id").(string)
		teamID, teamIDOk = d.GetOk("team_id")
		name, nameOk     = d.GetOk("name")

		err  error
		team *matlas.Team
	)

	if !teamIDOk && !nameOk {
		return diag.FromErr(errors.New("either team_id or name must be configured"))
	}

	if teamIDOk {
		team, _, err = conn.Teams.Get(ctx, orgID, teamID.(string))
	} else {
		team, _, err = conn.Teams.GetOneTeamByName(ctx, orgID, name.(string))
	}

	if err != nil {
		return diag.FromErr(fmt.Errorf(errorTeamRead, err))
	}

	if err := d.Set("team_id", team.ID); err != nil {
		return diag.FromErr(fmt.Errorf(errorTeamSetting, "name", d.Id(), err))
	}

	if err := d.Set("name", team.Name); err != nil {
		return diag.FromErr(fmt.Errorf(errorTeamSetting, "name", d.Id(), err))
	}

	// Set Usernames
	users, _, err := conn.Teams.GetTeamUsersAssigned(ctx, orgID, team.ID)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorTeamRead, err))
	}

	usernames := []string{}
	for i := range users {
		usernames = append(usernames, users[i].Username)
	}

	if err := d.Set("usernames", usernames); err != nil {
		return diag.FromErr(fmt.Errorf(errorTeamSetting, "usernames", d.Id(), err))
	}

	d.SetId(encodeStateID(map[string]string{
		"org_id": orgID,
		"id":     team.ID,
	}))

	return nil
}
