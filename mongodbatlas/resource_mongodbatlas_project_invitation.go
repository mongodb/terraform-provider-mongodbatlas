package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func resourceMongoDBAtlasProjectInvitation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBAtlasProjectInvitationCreate,
		ReadContext:   resourceMongoDBAtlasProjectInvitationRead,
		DeleteContext: resourceMongoDBAtlasProjectInvitationDelete,
		UpdateContext: resourceMongoDBAtlasProjectInvitationUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasProjectInvitationImportState,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"invitation_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expires_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"roles": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceMongoDBAtlasProjectInvitationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	username := ids["username"]
	invitationID := ids["invitation_id"]

	projectInvitation, resp, err := conn.Projects.Invitation(context.Background(), projectID, invitationID)
	if err != nil {
		// case 404
		// deleted in the backend case
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf("error getting Project Invitation information: %s", err))
	}

	if err := d.Set("username", projectInvitation.Username); err != nil {
		return diag.FromErr(fmt.Errorf("error getting `username` for Project Invitation (%s): %s", d.Id(), err))
	}

	if err := d.Set("project_id", projectInvitation.GroupID); err != nil {
		return diag.FromErr(fmt.Errorf("error getting `project_id` for Project Invitation (%s): %s", d.Id(), err))
	}

	if err := d.Set("invitation_id", projectInvitation.ID); err != nil {
		return diag.FromErr(fmt.Errorf("error getting `invitation_id` for Project Invitation (%s): %s", d.Id(), err))
	}

	if err := d.Set("expires_at", projectInvitation.ExpiresAt); err != nil {
		return diag.FromErr(fmt.Errorf("error getting `expires_at` for Project Invitation (%s): %s", d.Id(), err))
	}

	if err := d.Set("created_at", projectInvitation.CreatedAt); err != nil {
		return diag.FromErr(fmt.Errorf("error getting `created_at` for Project Invitation (%s): %s", d.Id(), err))
	}

	if err := d.Set("roles", projectInvitation.Roles); err != nil {
		return diag.FromErr(fmt.Errorf("error getting `roles` for Project Invitation (%s): %s", d.Id(), err))
	}

	d.SetId(encodeStateID(map[string]string{
		"username":      username,
		"project_id":    projectID,
		"invitation_id": invitationID,
	}))

	return nil
}

func resourceMongoDBAtlasProjectInvitationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)

	err := validateProjectRoles(d.Get("roles").(*schema.Set))
	if err != nil {
		return diag.FromErr(err)
	}

	InvitationReq := &matlas.Invitation{
		Roles:    createProjectStringListFromSetSchema(d.Get("roles").(*schema.Set)),
		Username: d.Get("username").(string),
	}

	InvitationRes, resp, err := conn.Projects.InviteUser(ctx, projectID, InvitationReq)
	if err != nil {
		// case 404
		// deleted in the backend case
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error creating Atlas user: %s", err))
	}

	d.SetId(encodeStateID(map[string]string{
		"username":      InvitationRes.Username,
		"project_id":    InvitationRes.GroupID,
		"invitation_id": InvitationRes.ID,
	}))

	return resourceMongoDBAtlasProjectInvitationRead(ctx, d, meta)
}

func resourceMongoDBAtlasProjectInvitationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	username := ids["username"]
	invitationID := ids["invitation_id"]

	_, err := conn.Projects.DeleteInvitation(ctx, projectID, invitationID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating Project invitation: %s for %s", username, err))
	}

	return nil
}

func resourceMongoDBAtlasProjectInvitationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	username := ids["username"]
	invitationID := ids["invitation_id"]

	InvitationReq := &matlas.Invitation{
		Roles: createProjectStringListFromSetSchema(d.Get("roles").(*schema.Set)),
	}

	_, _, err := conn.Projects.UpdateInvitationByID(ctx, projectID, invitationID, InvitationReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating Project invitation: %s for %s", username, err))
	}

	return resourceMongoDBAtlasProjectInvitationRead(ctx, d, meta)
}

func resourceMongoDBAtlasProjectInvitationImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*MongoDBClient).Atlas
	projectID, username, err := splitProjectInvitationImportID(d.Id())
	if err != nil {
		return nil, err
	}

	projectInvitations, _, err := conn.Projects.Invitations(ctx, projectID, nil)
	if err != nil {
		return nil, fmt.Errorf("couldn't import Project invitations, error: %s", err)
	}

	for _, projectInvitation := range projectInvitations {
		if projectInvitation.Username == username {

			if err := d.Set("username", projectInvitation.Username); err != nil {
				return nil, fmt.Errorf("error getting `username` for Project Invitation (%s): %s", username, err)
			}

			if err := d.Set("project_id", projectInvitation.GroupID); err != nil {
				return nil, fmt.Errorf("error getting `project_id` for Project Invitation (%s): %s", username, err)
			}

			if err := d.Set("invitation_id", projectInvitation.ID); err != nil {
				return nil, fmt.Errorf("error getting `invitation_id` for Project Invitation (%s): %s", username, err)
			}

			d.SetId(encodeStateID(map[string]string{
				"username":      username,
				"project_id":    projectID,
				"invitation_id": projectInvitation.ID,
			}))

			return []*schema.ResourceData{d}, nil
		}
	}

	return nil, fmt.Errorf("Could not import Project Invitation for %s", d.Id())
}

func splitProjectInvitationImportID(id string) (projectID, username string, err error) {
	var re = regexp.MustCompile(`(?s)^([0-9a-fA-F]{24})-(.*)$`)
	parts := re.FindStringSubmatch(id)

	if len(parts) != 3 {
		err = fmt.Errorf("import format error: to import a Project Invitation, use the format {project_id}-{username}")
		return
	}

	projectID = parts[1]
	username = parts[2]

	return
}

func validateProjectRoles(list *schema.Set) error {
	if rs := list.List(); list.Len() > 0 {
		for _, role := range rs {
			if validateProjectRole(role.(string)) == false {
				return fmt.Errorf("error creating an invite: %s is an invalid role for a Project", role)
			}
		}
	}

	return nil
}

func validateProjectRole(str string) bool {
	proj_roles := []string{
		"GROUP_OWNER",
		"GROUP_CLUSTER_MANAGER",
		"GROUP_READ_ONLY",
		"GROUP_DATA_ACCESS_ADMIN",
		"GROUP_DATA_ACCESS_READ_WRITE",
		"GROUP_DATA_ACCESS_READ_ONLY",
	}

	for _, valid_role := range proj_roles {
		if valid_role == str {
			return true
		}
	}

	return false
}

func createProjectStringListFromSetSchema(list *schema.Set) []string {
	res := make([]string, list.Len())
	for i, v := range list.List() {
		res[i] = v.(string)
	}

	return res
}
