package orginvitation

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate,
		ReadContext:   resourceRead,
		DeleteContext: resourceDelete,
		UpdateContext: resourceUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: resourceImport,
		},
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:     schema.TypeString,
				Required: true,
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
			"teams_ids": {
				Type:       schema.TypeSet,
				Optional:   true,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	orgID := ids["org_id"]
	username := ids["username"]
	invitationID := ids["invitation_id"]

	if orgID != invitationID {
		orgInvitation, _, err := connV2.OrganizationsApi.GetOrganizationInvitation(ctx, orgID, invitationID).Execute()
		if err != nil {
			if strings.Contains(err.Error(), "404") { // case 404: deleted in the backend case
				if validateOrgInvitationAlreadyAccepted(ctx, meta.(*config.MongoDBClient), username, orgID) {
					d.SetId("")
					return nil
				}
				return nil
			}

			return diag.Errorf("error getting Organization Invitation information: %s", err)
		}

		if err := d.Set("username", orgInvitation.GetUsername()); err != nil {
			return diag.Errorf("error getting `username` for Organization Invitation (%s): %s", d.Id(), err)
		}

		if err := d.Set("org_id", orgInvitation.GetOrgId()); err != nil {
			return diag.Errorf("error getting `username` for Organization Invitation (%s): %s", d.Id(), err)
		}

		if err := d.Set("invitation_id", orgInvitation.GetId()); err != nil {
			return diag.Errorf("error getting `invitation_id` for Organization Invitation (%s): %s", d.Id(), err)
		}

		if err := d.Set("expires_at", conversion.TimePtrToStringPtr(orgInvitation.ExpiresAt)); err != nil {
			return diag.Errorf("error getting `expires_at` for Organization Invitation (%s): %s", d.Id(), err)
		}

		if err := d.Set("created_at", conversion.TimePtrToStringPtr(orgInvitation.CreatedAt)); err != nil {
			return diag.Errorf("error getting `created_at` for Organization Invitation (%s): %s", d.Id(), err)
		}

		if err := d.Set("inviter_username", orgInvitation.GetInviterUsername()); err != nil {
			return diag.Errorf("error getting `inviter_username` for Organization Invitation (%s): %s", d.Id(), err)
		}

		if err := d.Set("teams_ids", orgInvitation.GetTeamIds()); err != nil {
			return diag.Errorf("error getting `teams_ids` for Organization Invitation (%s): %s", d.Id(), err)
		}

		if err := d.Set("roles", orgInvitation.GetRoles()); err != nil {
			return diag.Errorf("error getting `roles` for Organization Invitation (%s): %s", d.Id(), err)
		}
	}
	d.SetId(conversion.EncodeStateID(map[string]string{
		"username":      username,
		"org_id":        orgID,
		"invitation_id": invitationID,
	}))

	return nil
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*config.MongoDBClient).Atlas
	orgID := d.Get("org_id").(string)

	invitationReq := &matlas.Invitation{
		Roles:    conversion.ExpandStringListFromSetSchema(d.Get("roles").(*schema.Set)),
		TeamIDs:  conversion.ExpandStringListFromSetSchema(d.Get("teams_ids").(*schema.Set)),
		Username: d.Get("username").(string),
	}

	if validateOrgInvitationAlreadyAccepted(ctx, meta.(*config.MongoDBClient), invitationReq.Username, orgID) {
		d.SetId(conversion.EncodeStateID(map[string]string{
			"username":      invitationReq.Username,
			"org_id":        orgID,
			"invitation_id": orgID,
		}))
	} else {
		invitationRes, _, err := conn.Organizations.InviteUser(ctx, orgID, invitationReq)
		if err != nil {
			return diag.Errorf("error creating Organization invitation for user %s: %s", d.Get("username").(string), err)
		}

		d.SetId(conversion.EncodeStateID(map[string]string{
			"username":      invitationRes.Username,
			"org_id":        invitationRes.OrgID,
			"invitation_id": invitationRes.ID,
		}))
	}
	return resourceRead(ctx, d, meta)
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	orgID := ids["org_id"]
	username := ids["username"]
	invitationID := ids["invitation_id"]

	_, _, err := connV2.OrganizationsApi.GetOrganizationInvitation(ctx, orgID, invitationID).Execute()
	if err != nil {
		if strings.Contains(err.Error(), "404") { // case 404: deleted in the backend case
			if validateOrgInvitationAlreadyAccepted(ctx, meta.(*config.MongoDBClient), username, orgID) {
				d.SetId("")
				return nil
			}
			return nil
		}
	}
	_, _, err = connV2.OrganizationsApi.DeleteOrganizationInvitation(ctx, orgID, invitationID).Execute()
	if err != nil {
		return diag.Errorf("error deleting Organization invitation for user %s: %s", username, err)
	}
	d.SetId("")
	return nil
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).Atlas
	ids := conversion.DecodeStateID(d.Id())
	orgID := ids["org_id"]
	username := ids["username"]
	invitationID := ids["invitation_id"]

	invitationReq := &matlas.Invitation{
		Roles: conversion.ExpandStringListFromSetSchema(d.Get("roles").(*schema.Set)),
	}

	_, _, err := conn.Organizations.UpdateInvitationByID(ctx, orgID, invitationID, invitationReq)
	if err != nil {
		return diag.Errorf("error updating Organization invitation for user %s: for %s", username, err)
	}

	return resourceRead(ctx, d, meta)
}

func resourceImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	conn := meta.(*config.MongoDBClient).Atlas
	orgID, username, err := splitOrgInvitationImportID(d.Id())
	if err != nil {
		return nil, err
	}

	orgInvitations, _, err := conn.Organizations.Invitations(ctx, orgID, nil)
	if err != nil {
		return nil, fmt.Errorf("couldn't import Organization invitations, error: %s", err)
	}

	for _, orgInvitation := range orgInvitations {
		if orgInvitation.Username != username {
			continue
		}

		if err := d.Set("username", orgInvitation.Username); err != nil {
			return nil, fmt.Errorf("error getting `username` for Organization Invitation (%s): %s", username, err)
		}
		if err := d.Set("org_id", orgInvitation.GroupID); err != nil {
			return nil, fmt.Errorf("error getting `org_id` for Organization Invitation (%s): %s", username, err)
		}
		if err := d.Set("invitation_id", orgInvitation.ID); err != nil {
			return nil, fmt.Errorf("error getting `invitation_id` for Organization Invitation (%s): %s", username, err)
		}
		d.SetId(conversion.EncodeStateID(map[string]string{
			"username":      username,
			"org_id":        orgID,
			"invitation_id": orgInvitation.ID,
		}))
		return []*schema.ResourceData{d}, nil
	}

	return nil, fmt.Errorf("could not import Organization Invitation for %s", d.Id())
}

func splitOrgInvitationImportID(id string) (orgID, username string, err error) {
	var re = regexp.MustCompile(`(?s)^([0-9a-fA-F]{24})-(.*)$`)
	parts := re.FindStringSubmatch(id)

	if len(parts) != 3 {
		err = fmt.Errorf("import format error: to import a Organization Invitation, use the format {org_id}-{username}")
		return
	}

	orgID = parts[1]
	username = parts[2]

	return
}

func validateOrgInvitationAlreadyAccepted(ctx context.Context, conn *config.MongoDBClient, username, orgID string) bool {
	user, _, err := conn.Atlas.AtlasUsers.GetByName(ctx, username)
	if err != nil {
		return false
	}

	for _, role := range user.Roles {
		if role.OrgID == orgID {
			return true
		}
	}

	return false
}
