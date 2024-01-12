package organization

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20231115003/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mwielbut/pointy"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func PluralDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasOrganizationsRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"include_deleted_orgs": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: fmt.Sprintf(constant.DeprecationParamByDate, "January 2025"),
			},
			"page_num": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"items_per_page": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_deleted": {
							Type:     schema.TypeBool,
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
				},
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceMongoDBAtlasOrganizationsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*config.MongoDBClient).AtlasV2

	organizationOptions := &admin.ListOrganizationsApiParams{
		PageNum:      pointy.Int(d.Get("page_num").(int)),
		ItemsPerPage: pointy.Int(d.Get("items_per_page").(int)),
		Name:         pointy.String(d.Get("name").(string)),
	}

	organizations, _, err := conn.OrganizationsApi.ListOrganizationsWithParams(ctx, organizationOptions).Execute()

	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting organization information: %s", err))
	}

	if err := d.Set("results", flattenOrganizations(organizations.Results)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `results`: %s", err))
	}

	if err := d.Set("total_count", organizations.TotalCount); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `total_count`: %s", err))
	}

	d.SetId(id.UniqueId())

	return nil
}

func flattenOrganizationLinks(links []admin.Link) []map[string]any {
	linksList := make([]map[string]any, 0)

	for _, link := range links {
		mLink := map[string]any{
			"href": link.Href,
			"rel":  link.Rel,
		}
		linksList = append(linksList, mLink)
	}

	return linksList
}

func flattenOrganizations(organizations []admin.AtlasOrganization) []map[string]any {
	var results []map[string]any

	if len(organizations) == 0 {
		return results
	}

	results = make([]map[string]any, len(organizations))

	for k, organization := range organizations {
		results[k] = map[string]any{
			"id":         organization.Id,
			"name":       organization.Name,
			"is_deleted": organization.IsDeleted,
			"links":      flattenOrganizationLinks(organization.Links),
		}
	}

	return results
}
