package privatelinkendpointservicedatafederationonlinearchive

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
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
	endopointID := d.Get("endpoint_id").(string)

	privateEndpoint, _, err := connV2.DataFederationApi.GetDataFederationPrivateEndpoint(ctx, projectID, endopointID).Execute()
	if err != nil {
		return diag.Errorf(errorPrivateEndpointServiceDataFederationOnlineArchiveRead, endopointID, projectID, err)
	}
	err = populateResourceData(d, privateEndpoint)
	if err != nil {
		return diag.Errorf(errorPrivateEndpointServiceDataFederationOnlineArchiveRead, endopointID, projectID, err)
	}
	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":  projectID,
		"endpoint_id": endopointID,
	}))

	return nil
}
