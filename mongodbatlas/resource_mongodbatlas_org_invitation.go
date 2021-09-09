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

func resourceMongoDBAtlasOrgInvitation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBAtlasOrgInvitationCreate,
		ReadContext:   resourceMongoDBAtlasOrgInvitationRead,
		DeleteContext: resourceMongoDBAtlasOrgInvitationDelete,
		UpdateContext: resourceMongoDBAtlasOrgInvitationUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasOrgInvitationImportState,
		},
		Schema: map[string]*schema.Schema{
			"org_id": {
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

func resourceMongoDBAtlasOrgInvitationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	orgID := ids["org_id"]
	username := ids["username"]
	invitationID := ids["invitation_id"]

	orgInvitation, resp, err := conn.Organizations.Invitation(context.Background(), orgID, invitationID)
	if err != nil {
		// case 404
		// deleted in the backend case
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf("error getting Organisation Invitation information: %s", err))
	}

	if err := d.Set("username", orgInvitation.Username); err != nil {
		return diag.FromErr(fmt.Errorf("error getting `username` for Organisation Invitation (%s): %s", d.Id(), err))
	}

	if err := d.Set("org_id", orgInvitation.GroupID); err != nil {
		return diag.FromErr(fmt.Errorf("error getting `username` for Organisation Invitation (%s): %s", d.Id(), err))
	}

	if err := d.Set("invitation_id", orgInvitation.ID); err != nil {
		return diag.FromErr(fmt.Errorf("error getting `invitation_id` for Organisation Invitation (%s): %s", d.Id(), err))
	}

	if err := d.Set("expires_at", orgInvitation.ExpiresAt); err != nil {
		return diag.FromErr(fmt.Errorf("error getting `expires_at` for Organisation Invitation (%s): %s", d.Id(), err))
	}

	if err := d.Set("created_at", orgInvitation.CreatedAt); err != nil {
		return diag.FromErr(fmt.Errorf("error getting `created_at` for Organisation Invitation (%s): %s", d.Id(), err))
	}

	if err := d.Set("roles", orgInvitation.Roles); err != nil {
		return diag.FromErr(fmt.Errorf("error getting `roles` for Organisation Invitation (%s): %s", d.Id(), err))
	}

	d.SetId(encodeStateID(map[string]string{
		"username":      username,
		"org_id":        orgID,
		"invitation_id": invitationID,
	}))

	return nil
}

func resourceMongoDBAtlasOrgInvitationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	orgID := d.Get("org_id").(string)

	err := validateOrgRoles(d.Get("roles").(*schema.Set))
	if err != nil {
		return diag.FromErr(err)
	}

	InvitationReq := &matlas.Invitation{
		Roles:    createOrgStringListFromSetSchema(d.Get("roles").(*schema.Set)),
		Username: d.Get("username").(string),
	}

	InvitationRes, resp, err := conn.Organizations.InviteUser(ctx, orgID, InvitationReq)
	if err != nil {
		// case 404
		// deleted in the backend case
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error creating Organisation invitation: %s", err))
	}

	d.SetId(encodeStateID(map[string]string{
		"username":      InvitationRes.Username,
		"org_id":        InvitationRes.OrgID,
		"invitation_id": InvitationRes.ID,
	}))

	return resourceMongoDBAtlasOrgInvitationRead(ctx, d, meta)
}

func resourceMongoDBAtlasOrgInvitationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	orgID := ids["org_id"]
	username := ids["username"]
	invitationID := ids["invitation_id"]

	_, err := conn.Organizations.DeleteInvitation(ctx, orgID, invitationID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating Organisation invitation: %s for %s", username, err))
	}

	return nil
}

func resourceMongoDBAtlasOrgInvitationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	orgID := ids["org_id"]
	username := ids["username"]
	invitationID := ids["invitation_id"]

	InvitationReq := &matlas.Invitation{
		Roles: createOrgStringListFromSetSchema(d.Get("roles").(*schema.Set)),
	}

	_, _, err := conn.Organizations.UpdateInvitationByID(ctx, orgID, invitationID, InvitationReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating Organisation invitation: %s for %s", username, err))
	}

	return resourceMongoDBAtlasOrgInvitationRead(ctx, d, meta)
}

func resourceMongoDBAtlasOrgInvitationImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*MongoDBClient).Atlas
	orgID, username, err := splitProjectInvitationImportID(d.Id())
	if err != nil {
		return nil, err
	}

	orgInvitations, _, err := conn.Organizations.Invitations(ctx, orgID, nil)
	if err != nil {
		return nil, fmt.Errorf("couldn't import Organisation invitations, error: %s", err)
	}

	for _, orgInvitation := range orgInvitations {
		if orgInvitation.Username == username {

			if err := d.Set("username", orgInvitation.Username); err != nil {
				return nil, fmt.Errorf("error getting `username` for Organisation Invitation (%s): %s", username, err)
			}

			if err := d.Set("org_id", orgInvitation.GroupID); err != nil {
				return nil, fmt.Errorf("error getting `org_id` for Organisation Invitation (%s): %s", username, err)
			}

			if err := d.Set("invitation_id", orgInvitation.ID); err != nil {
				return nil, fmt.Errorf("error getting `invitation_id` for Organisation Invitation (%s): %s", username, err)
			}

			d.SetId(encodeStateID(map[string]string{
				"username":      username,
				"org_id":        orgID,
				"invitation_id": orgInvitation.ID,
			}))

			return []*schema.ResourceData{d}, nil
		}
	}

	return nil, fmt.Errorf("Could not import Organisation Invitation for %s", d.Id())
}

func splitOrgInvitationImportID(id string) (orgID, username string, err error) {
	var re = regexp.MustCompile(`(?s)^([0-9a-fA-F]{24})-(.*)$`)
	parts := re.FindStringSubmatch(id)

	if len(parts) != 3 {
		err = fmt.Errorf("import format error: to import a Organisation Invitation, use the format {project_id}-{username}")
		return
	}

	orgID = parts[1]
	username = parts[2]

	return
}

func validateOrgRoles(list *schema.Set) error {
	if rs := list.List(); list.Len() > 0 {
		for _, role := range rs {
			if validateOrgRole(role.(string)) == false {
				return fmt.Errorf("error creating an invite: %s is an invalid role for a Organisation", role)
			}
		}
	}

	return nil
}

func validateOrgRole(str string) bool {
	org_roles := []string{
		"ORG_OWNER",
		"ORG_GROUP_CREATOR",
		"ORG_BILLING_ADMIN",
		"ORG_READ_ONLY",
		"ORG_MEMBER",
	}

	for _, valid_role := range org_roles {
		if valid_role == str {
			return true
		}
	}

	return false
}

func createOrgStringListFromSetSchema(list *schema.Set) []string {
	res := make([]string, list.Len())
	for i, v := range list.List() {
		res[i] = v.(string)
	}

	return res
}
