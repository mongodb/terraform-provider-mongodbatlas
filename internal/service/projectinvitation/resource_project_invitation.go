package projectinvitation

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
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
	// Get client connection.
	conn := meta.(*config.MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)

	invitationReq := &matlas.Invitation{
		Roles:    createProjectStringListFromSetSchema(d.Get("roles").(*schema.Set)),
		Username: d.Get("username").(string),
	}

	invitationRes, _, err := conn.Projects.InviteUser(ctx, projectID, invitationReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating Project invitation for user %s: %w", d.Get("username").(string), err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"username":      invitationRes.Username,
		"project_id":    invitationRes.GroupID,
		"invitation_id": invitationRes.ID,
	}))

	return resourceRead(ctx, d, meta)
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).Atlas
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	username := ids["username"]
	invitationID := ids["invitation_id"]

	_, err := conn.Projects.DeleteInvitation(ctx, projectID, invitationID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting Project invitation for user %s: %w", username, err))
	}

	return nil
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*config.MongoDBClient).Atlas
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	username := ids["username"]
	invitationID := ids["invitation_id"]

	projectInvitation, resp, err := conn.Projects.Invitation(ctx, projectID, invitationID)
	if err != nil {
		// case 404
		// deleted in the backend case
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf("error getting Project Invitation information: %w", err))
	}

	if err := d.Set("username", projectInvitation.Username); err != nil {
		return diag.FromErr(fmt.Errorf("error getting `username` for Project Invitation (%s): %w", d.Id(), err))
	}

	if err := d.Set("project_id", projectInvitation.GroupID); err != nil {
		return diag.FromErr(fmt.Errorf("error getting `project_id` for Project Invitation (%s): %w", d.Id(), err))
	}

	if err := d.Set("invitation_id", projectInvitation.ID); err != nil {
		return diag.FromErr(fmt.Errorf("error getting `invitation_id` for Project Invitation (%s): %w", d.Id(), err))
	}

	if err := d.Set("expires_at", projectInvitation.ExpiresAt); err != nil {
		return diag.FromErr(fmt.Errorf("error getting `expires_at` for Project Invitation (%s): %w", d.Id(), err))
	}

	if err := d.Set("created_at", projectInvitation.CreatedAt); err != nil {
		return diag.FromErr(fmt.Errorf("error getting `created_at` for Project Invitation (%s): %w", d.Id(), err))
	}

	if err := d.Set("inviter_username", projectInvitation.InviterUsername); err != nil {
		return diag.FromErr(fmt.Errorf("error getting `inviter_username` for Project Invitation (%s): %w", d.Id(), err))
	}

	if err := d.Set("roles", projectInvitation.Roles); err != nil {
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
	conn := meta.(*config.MongoDBClient).Atlas
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	username := ids["username"]
	invitationID := ids["invitation_id"]

	invitationReq := &matlas.Invitation{
		Roles: conversion.ExpandStringListFromSetSchema(d.Get("roles").(*schema.Set)),
	}

	_, _, err := conn.Projects.UpdateInvitationByID(ctx, projectID, invitationID, invitationReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating Project invitation for user %s: %w", username, err))
	}

	return resourceRead(ctx, d, meta)
}

func resourceImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	conn := meta.(*config.MongoDBClient).Atlas
	projectID, username, err := splitProjectInvitationImportID(d.Id())
	if err != nil {
		return nil, err
	}

	projectInvitations, _, err := conn.Projects.Invitations(ctx, projectID, nil)
	if err != nil {
		return nil, fmt.Errorf("couldn't import Project invitations, error: %s", err)
	}

	for _, projectInvitation := range projectInvitations {
		if projectInvitation.Username != username {
			continue
		}

		if err := d.Set("username", projectInvitation.Username); err != nil {
			return nil, fmt.Errorf("error getting `username` for Project Invitation (%s): %w", username, err)
		}
		if err := d.Set("project_id", projectInvitation.GroupID); err != nil {
			return nil, fmt.Errorf("error getting `project_id` for Project Invitation (%s): %w", username, err)
		}
		if err := d.Set("invitation_id", projectInvitation.ID); err != nil {
			return nil, fmt.Errorf("error getting `invitation_id` for Project Invitation (%s): %w", username, err)
		}
		d.SetId(conversion.EncodeStateID(map[string]string{
			"username":      username,
			"project_id":    projectID,
			"invitation_id": projectInvitation.ID,
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
