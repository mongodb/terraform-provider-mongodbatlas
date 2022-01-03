package mongodbatlas

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMongoDBAtlasPrivateLinkEndpointServiceADL() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasPrivateLinkEndpointServiceADLRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"endpoint_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"provider_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceMongoDBAtlasPrivateLinkEndpointServiceADLRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)
	endpointID := d.Get("endpoint_id").(string)

	privateLinkResponse, _, err := conn.DataLakes.GetPrivateLinkEndpoint(ctx, projectID, endpointID)
	if err != nil {
		// case 404
		// deleted in the backend case
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
			return nil
		}

		return diag.Errorf("error getting ADL PrivateLink Endpoint Information: %s", err)
	}

	if err := d.Set("endpoint_id", privateLinkResponse.EndpointID); err != nil {
		return diag.Errorf("error setting `endpoint_id` for endpoint_id (%s): %s", d.Id(), err)
	}

	if err := d.Set("type", privateLinkResponse.Type); err != nil {
		return diag.Errorf("error setting `type` for endpoint_id (%s): %s", d.Id(), err)
	}

	if err := d.Set("comment", privateLinkResponse.Comment); err != nil {
		return diag.Errorf("error setting `comment` for endpoint_id (%s): %s", d.Id(), err)
	}

	if err := d.Set("provider_name", privateLinkResponse.Provider); err != nil {
		return diag.Errorf("error setting `provider_name` for endpoint_id (%s): %s", d.Id(), err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":  projectID,
		"endpoint_id": privateLinkResponse.EndpointID,
	}))

	return nil
}
