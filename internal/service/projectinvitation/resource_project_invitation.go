package projectinvitation

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"go.mongodb.org/atlas-sdk/v20250312005/admin"
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
			"inviter_username": {
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

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)

	roles := createProjectStringListFromSetSchema(d.Get("roles").(*schema.Set))
	invitationReq := &admin.GroupInvitationRequest{
		Roles:    &roles,
		Username: conversion.StringPtr(d.Get("username").(string)),
	}

	invitationRes, _, err := connV2.ProjectsApi.CreateProjectInvitation(ctx, projectID, invitationReq).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating Project invitation for user %s: %w", d.Get("username").(string), err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"username":      invitationRes.GetUsername(),
		"project_id":    invitationRes.GetGroupId(),
		"invitation_id": invitationRes.GetId(),
	}))

	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	username := ids["username"]
	invitationID := ids["invitation_id"]

	projectInvitation, resp, err := connV2.ProjectsApi.GetProjectInvitation(ctx, projectID, invitationID).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) { // case 404: deleted in the backend case
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf("error getting Project Invitation information: %w", err))
	}

	if err := d.Set("username", projectInvitation.GetUsername()); err != nil {
		return diag.FromErr(fmt.Errorf("error getting `username` for Project Invitation (%s): %w", d.Id(), err))
	}

	if err := d.Set("project_id", projectInvitation.GetGroupId()); err != nil {
		return diag.FromErr(fmt.Errorf("error getting `project_id` for Project Invitation (%s): %w", d.Id(), err))
	}

	if err := d.Set("invitation_id", projectInvitation.GetId()); err != nil {
		return diag.FromErr(fmt.Errorf("error getting `invitation_id` for Project Invitation (%s): %w", d.Id(), err))
	}

	if err := d.Set("expires_at", conversion.TimePtrToStringPtr(projectInvitation.ExpiresAt)); err != nil {
		return diag.FromErr(fmt.Errorf("error getting `expires_at` for Project Invitation (%s): %w", d.Id(), err))
	}

	if err := d.Set("created_at", conversion.TimePtrToStringPtr(projectInvitation.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("error getting `created_at` for Project Invitation (%s): %w", d.Id(), err))
	}

	if err := d.Set("inviter_username", projectInvitation.GetInviterUsername()); err != nil {
		return diag.FromErr(fmt.Errorf("error getting `inviter_username` for Project Invitation (%s): %w", d.Id(), err))
	}

	if err := d.Set("roles", projectInvitation.GetRoles()); err != nil {
		return diag.FromErr(fmt.Errorf("error getting `roles` for Project Invitation (%s): %w", d.Id(), err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"username":      username,
		"project_id":    projectID,
		"invitation_id": invitationID,
	}))

	return nil
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	username := ids["username"]
	invitationID := ids["invitation_id"]

	roles := conversion.ExpandStringListFromSetSchema(d.Get("roles").(*schema.Set))
	invitationReq := &admin.GroupInvitationUpdateRequest{
		Roles: &roles,
	}
	_, _, err := connV2.ProjectsApi.UpdateProjectInvitationById(ctx, projectID, invitationID, invitationReq).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating Project invitation for user %s: %w", username, err))
	}
	return resourceRead(ctx, d, meta)
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	username := ids["username"]
	invitationID := ids["invitation_id"]
	_, err := connV2.ProjectsApi.DeleteProjectInvitation(ctx, projectID, invitationID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting Project invitation for user %s: %w", username, err))
	}
	return nil
}

func resourceImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID, username, err := splitProjectInvitationImportID(d.Id())
	if err != nil {
		return nil, err
	}

	projectInvitations, _, err := connV2.ProjectsApi.ListProjectInvitations(ctx, projectID).Execute()
	if err != nil {
		return nil, fmt.Errorf("couldn't import Project invitations, error: %s", err)
	}

	for _, projectInvitation := range projectInvitations {
		if conversion.SafeString(projectInvitation.Username) != username {
			continue
		}

		if err := d.Set("username", projectInvitation.GetUsername()); err != nil {
			return nil, fmt.Errorf("error getting `username` for Project Invitation (%s): %w", username, err)
		}
		if err := d.Set("project_id", projectInvitation.GetGroupId()); err != nil {
			return nil, fmt.Errorf("error getting `project_id` for Project Invitation (%s): %w", username, err)
		}
		if err := d.Set("invitation_id", projectInvitation.GetId()); err != nil {
			return nil, fmt.Errorf("error getting `invitation_id` for Project Invitation (%s): %w", username, err)
		}
		d.SetId(conversion.EncodeStateID(map[string]string{
			"username":      username,
			"project_id":    projectID,
			"invitation_id": projectInvitation.GetId(),
		}))
		return []*schema.ResourceData{d}, nil
	}

	return nil, fmt.Errorf("could not import Project Invitation for %s", d.Id())
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

func createProjectStringListFromSetSchema(list *schema.Set) []string {
	res := make([]string, list.Len())
	for i, v := range list.List() {
		res[i] = v.(string)
	}

	return res
}
