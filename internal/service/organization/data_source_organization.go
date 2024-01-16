package organization

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasOrganizationRead,
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_deleted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rel": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceMongoDBAtlasOrganizationRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*config.MongoDBClient).AtlasV2
	orgID := d.Get("org_id").(string)

	organization, _, err := conn.OrganizationsApi.GetOrganization(ctx, orgID).Execute()

	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting organizations information: %s", err))
	}

	if err := d.Set("name", organization.Name); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `name`: %s", err))
	}

	if err := d.Set("is_deleted", organization.IsDeleted); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `is_deleted`: %s", err))
	}

	if err := d.Set("links", flattenOrganizationLinks(organization.GetLinks())); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `is_deleted`: %s", err))
	}

	d.SetId(*organization.Id)

	return nil
}
