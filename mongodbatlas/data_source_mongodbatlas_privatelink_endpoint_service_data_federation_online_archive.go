package mongodbatlas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMongoDBAtlasPrivatelinkEndpointServiceDataFederationOnlineArchive() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasPrivatelinkEndpointServiceDataFederationOnlineArchiveRead,
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
			"comment": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceMongoDBAtlasPrivatelinkEndpointServiceDataFederationOnlineArchiveRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)
	endopointID := d.Get("endpoint_id").(string)

	privateEndpoint, _, err := conn.DataLakes.GetPrivateLinkEndpoint(context.Background(), projectID, endopointID)
	if err != nil {
		return diag.Errorf(errorPrivateEndpointServiceDataFederationOnlineArchiveRead, endopointID, projectID, err)
	}

	if err := d.Set("comment", privateEndpoint.Comment); err != nil {
		return diag.Errorf(errorPrivateEndpointServiceDataFederationOnlineArchiveRead, endopointID, projectID, err)
	}

	if err := d.Set("provider_name", privateEndpoint.Provider); err != nil {
		return diag.Errorf(errorPrivateEndpointServiceDataFederationOnlineArchiveRead, endopointID, projectID, err)
	}

	if err := d.Set("type", privateEndpoint.Type); err != nil {
		return diag.Errorf(errorPrivateEndpointServiceDataFederationOnlineArchiveRead, endopointID, projectID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":  projectID,
		"endpoint_id": endopointID,
	}))

	return nil
}
