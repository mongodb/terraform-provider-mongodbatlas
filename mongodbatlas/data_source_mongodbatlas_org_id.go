package mongodbatlas

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
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

	for idx, role := range apiKeyOrgList.APIKey.Roles {
		if strings.HasPrefix(role.RoleName, "ORG_") {
			if err := d.Set("org_id", apiKeyOrgList.APIKey.Roles[idx].OrgID); err != nil {
				return diag.Errorf(errorProjectSetting, `org_id`, root.APIKey.ID, err)
			}
			d.SetId(apiKeyOrgList.APIKey.Roles[idx].OrgID)
			return nil
		}
	}

	d.SetId(id.UniqueId())

	return nil
}
