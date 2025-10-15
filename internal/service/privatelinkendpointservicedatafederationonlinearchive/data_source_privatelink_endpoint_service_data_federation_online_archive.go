package privatelinkendpointservicedatafederationonlinearchive

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"customer_endpoint_dns_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)
	endpointID := d.Get("endpoint_id").(string)

	privateEndpoint, _, err := connV2.DataFederationApi.GetPrivateEndpointId(ctx, projectID, endpointID).Execute()
	if err != nil {
		return diag.Errorf(errorPrivateEndpointServiceDataFederationOnlineArchiveRead, endpointID, projectID, err)
	}
	err = populateResourceData(d, privateEndpoint, projectID, endpointID)
	if err != nil {
		return diag.Errorf(errorPrivateEndpointServiceDataFederationOnlineArchiveRead, endpointID, projectID, err)
	}
	return nil
}
