package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mwielbut/pointy"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasOrganizations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasOrganizationsRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"include_deleted_orgs": {
				Type:     schema.TypeBool,
				Optional: true,
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

func dataSourceMongoDBAtlasOrganizationsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas

	options := &matlas.ListOptions{
		PageNum:      d.Get("page_num").(int),
		ItemsPerPage: d.Get("items_per_page").(int),
	}

	organizationOptions := &matlas.OrganizationsListOptions{
		Name:               d.Get("name").(string),
		IncludeDeletedOrgs: pointy.Bool(d.Get("include_deleted_orgs").(bool)),
		ListOptions:        *options,
	}

	organizations, _, err := conn.Organizations.List(ctx, organizationOptions)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting organization information: %s", err))
	}

	if err := d.Set("results", flattenOrganizations(organizations.Results)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `results`: %s", err))
	}

	if err := d.Set("total_count", organizations.TotalCount); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `total_count`: %s", err))
	}

	d.SetId(resource.UniqueId())

	return nil
}

func flattenOrganizationLinks(links []*matlas.Link) []map[string]interface{} {
	linksList := make([]map[string]interface{}, 0)

	for _, link := range links {
		mLink := map[string]interface{}{
			"href": link.Href,
			"rel":  link.Rel,
		}
		linksList = append(linksList, mLink)
	}

	return linksList
}

func flattenOrganizations(organizations []*matlas.Organization) []map[string]interface{} {
	var results []map[string]interface{}

	if len(organizations) == 0 {
		return results
	}

	results = make([]map[string]interface{}, len(organizations))

	for k, organization := range organizations {
		results[k] = map[string]interface{}{
			"id":         organization.ID,
			"name":       organization.Name,
			"is_deleted": organization.IsDeleted,
			"links":      flattenOrganizationLinks(organization.Links),
		}
	}

	return results
}
