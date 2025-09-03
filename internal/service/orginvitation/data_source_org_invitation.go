package orginvitation

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: fmt.Sprintf(constant.DeprecationNextMajorWithReplacementGuide, "data source", "mongodbatlas_cloud_user_org_assignment", "https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/org-invitation-to-cloud-user-org-assignment-migration-guide"),
		ReadContext:        dataSourceRead,
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"invitation_id": {
				Type:     schema.TypeString,
				Required: true,
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
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"roles": {
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
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	orgID := d.Get("org_id").(string)
	username := d.Get("username").(string)
	invitationID := d.Get("invitation_id").(string)

	orgInvitation, _, err := connV2.OrganizationsApi.GetOrganizationInvitation(ctx, orgID, invitationID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting Organization Invitation information: %w", err))
	}

	if err := d.Set("username", orgInvitation.GetUsername()); err != nil {
		return diag.FromErr(fmt.Errorf("error getting `username` for Organization Invitation (%s): %w", d.Id(), err))
	}

	if err := d.Set("org_id", orgInvitation.GetOrgId()); err != nil {
		return diag.FromErr(fmt.Errorf("error getting `username` for Organization Invitation (%s): %w", d.Id(), err))
	}

	if err := d.Set("invitation_id", orgInvitation.GetId()); err != nil {
		return diag.FromErr(fmt.Errorf("error getting `invitation_id` for Organization Invitation (%s): %w", d.Id(), err))
	}

	if err := d.Set("expires_at", conversion.TimePtrToStringPtr(orgInvitation.ExpiresAt)); err != nil {
		return diag.FromErr(fmt.Errorf("error getting `expires_at` for Organization Invitation (%s): %w", d.Id(), err))
	}

	if err := d.Set("created_at", conversion.TimePtrToStringPtr(orgInvitation.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("error getting `created_at` for Organization Invitation (%s): %w", d.Id(), err))
	}

	if err := d.Set("inviter_username", orgInvitation.GetInviterUsername()); err != nil {
		return diag.FromErr(fmt.Errorf("error getting `inviter_username` for Organization Invitation (%s): %w", d.Id(), err))
	}

	if err := d.Set("teams_ids", orgInvitation.GetTeamIds()); err != nil {
		return diag.FromErr(fmt.Errorf("error getting `teams_ids` for Organization Invitation (%s): %w", d.Id(), err))
	}

	if err := d.Set("roles", orgInvitation.GetRoles()); err != nil {
		return diag.FromErr(fmt.Errorf("error getting `roles` for Organization Invitation (%s): %w", d.Id(), err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"username":      username,
		"org_id":        orgID,
		"invitation_id": invitationID,
	}))

	return nil
}
