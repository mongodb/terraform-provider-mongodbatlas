package team

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20240530002/admin"
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
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func LegacyTeamsResource() *schema.Resource {
	res := Resource()
	res.DeprecationMessage = fmt.Sprintf(constant.DeprecationResourceByDateWithReplacement, "November 2024", "mongodbatlas_team")
	return res
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	orgID := d.Get("org_id").(string)

	usernames := conversion.ExpandStringListFromSetSchema(d.Get("usernames").(*schema.Set))
	teamsResp, _, err := connV2.TeamsApi.CreateTeam(ctx, orgID,
		&admin.Team{
			Name:      d.Get("name").(string),
			Usernames: &usernames,
		}).Execute()
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
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	ids := conversion.DecodeStateID(d.Id())
	orgID := ids["org_id"]
	teamID := ids["id"]

	team, resp, err := connV2.TeamsApi.GetTeamById(context.Background(), orgID, teamID).Execute()

	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
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

	users, _, err := connV2.TeamsApi.ListTeamUsers(ctx, orgID, teamID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorTeamRead, err))
	}

	usernames := []string{}
	for i := range users.GetResults() {
		usernames = append(usernames, users.GetResults()[i].GetUsername())
	}

	if err := d.Set("usernames", usernames); err != nil {
		return diag.FromErr(fmt.Errorf(errorTeamSetting, "usernames", teamID, err))
	}

	return nil
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	ids := conversion.DecodeStateID(d.Id())
	orgID := ids["org_id"]
	teamID := ids["id"]

	if d.HasChange("name") {
		_, _, err := connV2.TeamsApi.RenameTeam(ctx, orgID, teamID,
			&admin.TeamUpdate{Name: d.Get("name").(string)},
		).Execute()
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorTeamUpdate, err))
		}
	}

	if d.HasChange("usernames") {
		users, _, err := connV2.TeamsApi.ListTeamUsers(ctx, orgID, teamID).Execute()

		if err != nil {
			return diag.FromErr(fmt.Errorf(errorTeamRead, err))
		}

		index := make(map[string]admin.CloudAppUser)
		for i := range users.GetResults() {
			index[users.GetResults()[i].GetUsername()] = users.GetResults()[i]
		}

		cleanUsers := func() error {
			for i := range users.GetResults() {
				_, err := connV2.TeamsApi.RemoveTeamUser(ctx, orgID, teamID, users.GetResults()[i].GetId()).Execute()
				if err != nil {
					return fmt.Errorf("error deleting Atlas User (%s) information: %s", teamID, err)
				}
			}
			return nil
		}

		var newUsers []admin.AddUserToTeam

		for _, username := range d.Get("usernames").(*schema.Set).List() {
			user, _, err := connV2.MongoDBCloudUsersApi.GetUserByUsername(ctx, username.(string)).Execute()

			updatedUserData := user

			if err != nil {
				if !strings.Contains(err.Error(), "401") {
					return diag.FromErr(fmt.Errorf("error getting Atlas User (%s) information: %s", username, err))
				}

				log.Printf("[WARN] error fetching information user for (%s): %s\n", username, err)
				if user == nil {
					log.Printf("[WARN] there is no runtime information to fetch, checking in the existing users")

					cached, ok := index[username.(string)]

					if !ok {
						log.Printf("[WARN] no information in cached for (%s)", username)
						return diag.FromErr(fmt.Errorf("error getting Atlas User (%s) information: %s", username, err))
					}
					updatedUserData = &cached
				}
			}
			newUsers = append(newUsers, admin.AddUserToTeam{Id: updatedUserData.GetId()})
		}

		err = cleanUsers()
		if err != nil {
			return diag.FromErr(err)
		}

		_, _, err = connV2.TeamsApi.AddTeamUser(ctx, orgID, teamID, &newUsers).Execute()
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorTeamAddUsers, err))
		}
	}

	return resourceRead(ctx, d, meta)
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	orgID := ids["org_id"]
	id := ids["id"]

	err := retry.RetryContext(ctx, 1*time.Hour, func() *retry.RetryError {
		_, _, err := connV2.TeamsApi.DeleteTeam(ctx, orgID, id).Execute()
		if err != nil {
			if admin.IsErrorCode(err, "CANNOT_DELETE_TEAM_ASSIGNED_TO_PROJECT") {
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

func getProjectIDByTeamID(ctx context.Context, connV2 *admin.APIClient, teamID string) (string, error) {
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
