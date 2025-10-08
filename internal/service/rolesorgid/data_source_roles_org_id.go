package rolesorgid

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312007/admin"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	orgID, err := GetCurrentOrgID(ctx, connV2)
	if err != nil {
		return diag.Errorf("error getting current organization ID: %v", err)
	}
	if err := d.Set("org_id", orgID); err != nil {
		return diag.Errorf("error setting `org_id`: %v", err)
	}
	d.SetId(orgID)
	return nil
}

// GetCurrentOrgID returns the current organization ID for the SA or Programmatic API key (PAK) from the authenticated user.
func GetCurrentOrgID(ctx context.Context, connV2 *admin.APIClient) (string, error) {
	resp, _, err := connV2.OrganizationsApi.ListOrgs(ctx).Execute()
	if err != nil {
		return "", err
	}
	orgIDs := resp.GetResults()
	if len(orgIDs) == 0 {
		return "", fmt.Errorf("no organizations found")
	}

	// At present a PAK or SA belongs to exactly one organization. If this changes in the future, this logic will need to be updated.
	return orgIDs[0].GetId(), nil
}
