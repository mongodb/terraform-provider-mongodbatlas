package team

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	admin20241113 "go.mongodb.org/atlas-sdk/v20241113005/admin"
	"go.mongodb.org/atlas-sdk/v20250312005/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/dsschema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
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
			"users": &dsschema.DSOrgUsersSchema,
		},
	}
}

func LegacyTeamsDataSource() *schema.Resource {
	res := DataSource()
	res.DeprecationMessage = fmt.Sprintf(constant.DeprecationDataSourceByDateWithReplacement, "November 2024", "mongodbatlas_team")
	return res
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		connV2           = meta.(*config.MongoDBClient).AtlasV220241113
		conn             = meta.(*config.MongoDBClient).AtlasV2
		orgID            = d.Get("org_id").(string)
		teamID, teamIDOk = d.GetOk("team_id")
		name, nameOk     = d.GetOk("name")

		err  error
		team *admin20241113.TeamResponse
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
		return diag.FromErr(fmt.Errorf(errorTeamSetting, "team_id", d.Id(), err))
	}

	if err := d.Set("name", team.GetName()); err != nil {
		return diag.FromErr(fmt.Errorf(errorTeamSetting, "name", d.Id(), err))
	}

	teamUsers, err := listAllTeamUsersDS(ctx, conn, orgID, team.GetId())

	if err != nil {
		return diag.FromErr(fmt.Errorf(errorTeamRead, err))
	}

	usernames := []string{}
	for i := range teamUsers {
		usernames = append(usernames, teamUsers[i].GetUsername())
	}

	if err := d.Set("usernames", usernames); err != nil {
		return diag.FromErr(fmt.Errorf(errorTeamSetting, "usernames", d.Id(), err))
	}

	if err := d.Set("users", dsschema.FlattenUsers(teamUsers)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `users`: %s", err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"org_id": orgID,
		"id":     team.GetId(),
	}))

	return nil
}

func listAllTeamUsersDS(ctx context.Context, conn *admin.APIClient, orgID, teamID string) ([]admin.OrgUserResponse, error) {
	return dsschema.AllPages(ctx, func(ctx context.Context, pageNum int) (dsschema.PaginateResponse[admin.OrgUserResponse], *http.Response, error) {
		request := conn.MongoDBCloudUsersApi.ListTeamUsers(ctx, orgID, teamID)
		request = request.PageNum(pageNum)
		return request.Execute()
	})
}
