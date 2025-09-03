package projectapikey

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"api_key_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_assignment": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role_names": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	projectID := d.Get("project_id").(string)
	apiKeyID := d.Get("api_key_id").(string)
	projectAPIKeys, _, err := connV2.ProgrammaticAPIKeysApi.ListGroupApiKeys(ctx, projectID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting api key information: %s", err))
	}

	for _, val := range projectAPIKeys.GetResults() {
		if val.GetId() != apiKeyID {
			continue
		}

		if err := d.Set("description", val.GetDesc()); err != nil {
			return diag.FromErr(fmt.Errorf("error setting `description`: %s", err))
		}

		if err := d.Set("public_key", val.GetPublicKey()); err != nil {
			return diag.FromErr(fmt.Errorf("error setting `public_key`: %s", err))
		}

		if err := d.Set("private_key", val.GetPrivateKey()); err != nil {
			return diag.FromErr(fmt.Errorf("error setting `private_key`: %s", err))
		}

		details, _, err := getKeyDetails(ctx, connV2, apiKeyID)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error getting api key information: %s", err))
		}
		if err := d.Set("project_assignment", flattenProjectAssignments(details.GetRoles())); err != nil {
			return diag.FromErr(fmt.Errorf("error setting `project_assignment`: %s", err))
		}
	}

	d.SetId(id.UniqueId())

	return nil
}
