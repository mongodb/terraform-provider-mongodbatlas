package team

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	admin20241113 "go.mongodb.org/atlas-sdk/v20241113005/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/dsschema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	errorTeamCreate   = "error creating Team information: %s"
	errorTeamAddUsers = "error adding users to the Team information: %s"
	errorTeamRead     = "error getting Team information: %s"
	errorTeamUpdate   = "error updating Team information: %s"
	errorTeamDelete   = "error deleting Team (%s): %s"
	errorTeamSetting  = "error setting `%s` for Team (%s): %s"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate,
		ReadContext:   resourceRead,
		UpdateContext: resourceUpdate,
		DeleteContext: resourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceImport,
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
				Type:       schema.TypeSet,
				Optional:   true,
				Computed:   true,
				Deprecated: fmt.Sprintf(constant.DeprecationNextMajorWithReplacementGuide, "parameter", "mongodbatlas_cloud_user_team_assignment", "<link-to-migration-guide>"),
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV220241113
	orgID := d.Get("org_id").(string)

	usernames := conversion.ExpandStringListFromSetSchema(d.Get("usernames").(*schema.Set))
	createTeamReq := &admin20241113.Team{
		Name: d.Get("name").(string),
	}

	if len(usernames) > 0 {
		createTeamReq.Usernames = usernames
	}

	teamsResp, _, err := connV2.TeamsApi.CreateTeam(ctx, orgID, createTeamReq).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorTeamCreate, err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"org_id": orgID,
		"id":     teamsResp.GetId(),
	}))

	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV220241113

	ids := conversion.DecodeStateID(d.Id())
	orgID := ids["org_id"]
	teamID := ids["id"]

	team, resp, err := connV2.TeamsApi.GetTeamById(context.Background(), orgID, teamID).Execute()

	if err != nil {
		if validate.StatusNotFound(resp) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf(errorTeamRead, err))
	}

	if err = d.Set("name", team.GetName()); err != nil {
		return diag.FromErr(fmt.Errorf(errorTeamSetting, "name", teamID, err))
	}

	if err = d.Set("team_id", team.GetId()); err != nil {
		return diag.FromErr(fmt.Errorf(errorTeamSetting, "team_id", teamID, err))
	}

	teamUsers, err := listAllTeamUsers(ctx, connV2, orgID, team.GetId())

	if err != nil {
		return diag.FromErr(fmt.Errorf(errorTeamRead, err))
	}

	usernames := []string{}
	for i := range teamUsers {
		usernames = append(usernames, teamUsers[i].GetUsername())
	}

	if err := d.Set("usernames", usernames); err != nil {
		return diag.FromErr(fmt.Errorf(errorTeamSetting, "usernames", teamID, err))
	}

	return nil
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV220241113

	ids := conversion.DecodeStateID(d.Id())
	orgID := ids["org_id"]
	teamID := ids["id"]

	if d.HasChange("name") {
		_, _, err := connV2.TeamsApi.RenameTeam(ctx, orgID, teamID,
			&admin20241113.TeamUpdate{Name: d.Get("name").(string)},
		).Execute()
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorTeamUpdate, err))
		}
	}

	if d.HasChange("usernames") {
		existingUsers, err := listAllTeamUsers(ctx, connV2, orgID, teamID)

		if err != nil {
			return diag.FromErr(fmt.Errorf(errorTeamRead, err))
		}
		newUsernames := conversion.ExpandStringList(d.Get("usernames").(*schema.Set).List())

		err = UpdateTeamUsers(connV2.TeamsApi, connV2.MongoDBCloudUsersApi, existingUsers, newUsernames, orgID, teamID)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error when updating usernames in team: %s", err))
		}
	}

	return resourceRead(ctx, d, meta)
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV220241113
	ids := conversion.DecodeStateID(d.Id())
	orgID := ids["org_id"]
	id := ids["id"]

	err := retry.RetryContext(ctx, 1*time.Hour, func() *retry.RetryError {
		_, _, err := connV2.TeamsApi.DeleteTeam(ctx, orgID, id).Execute()
		if err != nil {
			if admin20241113.IsErrorCode(err, "CANNOT_DELETE_TEAM_ASSIGNED_TO_PROJECT") {
				projectID, err := getProjectIDByTeamID(ctx, connV2, id)
				if err != nil {
					return retry.NonRetryableError(err)
				}

				_, err = connV2.TeamsApi.RemoveProjectTeam(ctx, projectID, id).Execute()
				if err != nil {
					return retry.NonRetryableError(fmt.Errorf(errorTeamDelete, id, err))
				}
				return retry.RetryableError(fmt.Errorf("will retry again"))
			}
		}
		return nil
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorTeamDelete, id, err))
	}

	return nil
}

func resourceImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 2 {
		return nil, errors.New("import format error: to import a team, use the format {group_id}-{team_id}")
	}

	orgID := parts[0]
	teamID := parts[1]

	team, _, err := connV2.TeamsApi.GetTeamById(ctx, orgID, teamID).Execute()
	if err != nil {
		return nil, fmt.Errorf("couldn't import team (%s) in organization(%s), error: %s", teamID, orgID, err)
	}

	if err := d.Set("org_id", orgID); err != nil {
		log.Printf("[WARN] Error setting org_id for (%s): %s", teamID, err)
	}

	if err := d.Set("team_id", teamID); err != nil {
		log.Printf("[WARN] Error setting team_id for (%s): %s", teamID, err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"org_id": orgID,
		"id":     team.GetId(),
	}))

	return []*schema.ResourceData{d}, nil
}

func getProjectIDByTeamID(ctx context.Context, connV2 *admin20241113.APIClient, teamID string) (string, error) {
	projects, _, err := connV2.ProjectsApi.ListProjects(ctx).Execute()
	if err != nil {
		return "", fmt.Errorf("error getting projects information: %s", err)
	}

	for _, project := range projects.GetResults() {
		teams, _, err := connV2.TeamsApi.ListProjectTeams(ctx, project.GetId()).Execute()
		if err != nil {
			return "", fmt.Errorf("error getting teams from project information: %s", err)
		}

		for _, team := range teams.GetResults() {
			if team.GetTeamId() == teamID {
				return project.GetId(), nil
			}
		}
	}

	return "", nil
}

func listAllTeamUsers(ctx context.Context, connV2 *admin20241113.APIClient, orgID, teamID string) ([]admin20241113.CloudAppUser, error) {
	return dsschema.AllPages(ctx, func(ctx context.Context, pageNum int) (dsschema.PaginateResponse[admin20241113.CloudAppUser], *http.Response, error) {
		request := connV2.TeamsApi.ListTeamUsers(ctx, orgID, teamID)
		request = request.PageNum(pageNum)
		return request.Execute()
	})
}
