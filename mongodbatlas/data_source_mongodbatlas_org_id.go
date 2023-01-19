package mongodbatlas

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasOrgID() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasOrgIDRead,
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceMongoDBAtlasOrgIDRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas

	var (
		err  error
		root *matlas.Root
	)

	options := &matlas.ListOptions{}
	apiKeyOrgList, _, err := conn.Root.List(ctx, options)
	if err != nil {
		return diag.Errorf("error getting API Key's org assigned (%s): ", err)
	}

	if err := d.Set("org_id", apiKeyOrgList.APIKey.Roles[0].OrgID); err != nil {
		return diag.Errorf(errorProjectSetting, `org_id`, root.APIKey.ID, err)
	}

	for _, role := range apiKeyOrgList.APIKey.Roles {
		if strings.HasPrefix(role.RoleName, "ORG_") {
			d.SetId(apiKeyOrgList.APIKey.Roles[0].OrgID)
		}
	}

	return nil
}
