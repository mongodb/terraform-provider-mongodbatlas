package projectinvitation

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This data source is deprecated and will be removed in the next major release. Please transition to mongodbatlas_cloud_user_project_assignment. For more details, see [Migration Guide: Project Invitation to Cloud User Project Assignment](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/project-invitation-to-cloud-user-project-assignment-migration-guide)",
		ReadContext:        dataSourceRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
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
	projectID := d.Get("project_id").(string)
	username := d.Get("username").(string)
	invitationID := d.Get("invitation_id").(string)

	projectInvitation, _, err := connV2.ProjectsApi.GetProjectInvitation(ctx, projectID, invitationID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting Project Invitation information: %w", err))
	}

	if err := d.Set("username", projectInvitation.GetUsername()); err != nil {
		return diag.FromErr(fmt.Errorf("error getting `username` for Project Invitation (%s): %w", d.Id(), err))
	}

	if err := d.Set("project_id", projectInvitation.GetGroupId()); err != nil {
		return diag.FromErr(fmt.Errorf("error getting `username` for Project Invitation (%s): %w", d.Id(), err))
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
		return diag.FromErr(fmt.Errorf("error getting `roles` for Project Invitation (%s): %s", d.Id(), err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"username":      username,
		"project_id":    projectID,
		"invitation_id": invitationID,
	}))

	return nil
}
