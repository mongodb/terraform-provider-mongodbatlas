package teams

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
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

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		connV2           = meta.(*config.MongoDBClient).AtlasV2
		orgID            = d.Get("org_id").(string)
		teamID, teamIDOk = d.GetOk("team_id")
		name, nameOk     = d.GetOk("name")

		err  error
		team *admin.TeamResponse
	)

	if !teamIDOk && !nameOk {
		return diag.FromErr(errors.New("either team_id or name must be configured"))
	}

	if teamIDOk {
		team, _, err = connV2.TeamsApi.GetTeamById(ctx, orgID, teamID.(string)).Execute()
	} else {
		team, _, err = connV2.TeamsApi.GetTeamByName(ctx, orgID, name.(string)).Execute()
	}

	if err != nil {
		return diag.FromErr(fmt.Errorf(errorTeamRead, err))
	}

	if err := d.Set("team_id", team.GetId()); err != nil {
		return diag.FromErr(fmt.Errorf(errorTeamSetting, "name", d.Id(), err))
	}

	if err := d.Set("name", team.GetName()); err != nil {
		return diag.FromErr(fmt.Errorf(errorTeamSetting, "name", d.Id(), err))
	}

	users, _, err := connV2.TeamsApi.ListTeamUsers(ctx, orgID, team.GetId()).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorTeamRead, err))
	}

	usernames := []string{}
	for i := range users.GetResults() {
		usernames = append(usernames, users.GetResults()[i].GetUsername())
	}

	if err := d.Set("usernames", usernames); err != nil {
		return diag.FromErr(fmt.Errorf(errorTeamSetting, "usernames", d.Id(), err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"org_id": orgID,
		"id":     team.GetId(),
	}))

	return nil
}
