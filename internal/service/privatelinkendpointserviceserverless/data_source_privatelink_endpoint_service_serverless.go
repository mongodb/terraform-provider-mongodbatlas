package privatelinkendpointserviceserverless

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/privatelinkendpointservice"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext:        dataSourceRead,
		DeprecationMessage: fmt.Sprintf(constant.DeprecationDataSourceByDateWithExternalLink, "March 2025", "placeholder-serverless-deprecation-url"),
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"endpoint_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint_service_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloud_provider_endpoint_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_link_service_resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_endpoint_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	projectID := d.Get("project_id").(string)
	instanceName := d.Get("instance_name").(string)
	endpointID := d.Get("endpoint_id").(string)

	serviceEndpoint, _, err := connV2.ServerlessPrivateEndpointsApi.GetServerlessPrivateEndpoint(ctx, projectID, instanceName, endpointID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(privatelinkendpointservice.ErrorServiceEndpointRead, endpointID, err))
	}

	if err := d.Set("error_message", serviceEndpoint.GetErrorMessage()); err != nil {
		return diag.FromErr(fmt.Errorf(privatelinkendpointservice.ErrorEndpointSetting, "error_message", endpointID, err))
	}

	if err := d.Set("status", serviceEndpoint.GetStatus()); err != nil {
		return diag.FromErr(fmt.Errorf(privatelinkendpointservice.ErrorEndpointSetting, "status", endpointID, err))
	}

	if err := d.Set("comment", serviceEndpoint.GetComment()); err != nil {
		return diag.FromErr(fmt.Errorf(privatelinkendpointservice.ErrorEndpointSetting, "comment", endpointID, err))
	}

	if err := d.Set("endpoint_service_name", serviceEndpoint.GetEndpointServiceName()); err != nil {
		return diag.FromErr(fmt.Errorf(privatelinkendpointservice.ErrorEndpointSetting, "endpoint_service_name", endpointID, err))
	}

	if err := d.Set("cloud_provider_endpoint_id", serviceEndpoint.GetCloudProviderEndpointId()); err != nil {
		return diag.FromErr(fmt.Errorf(privatelinkendpointservice.ErrorEndpointSetting, "cloud_provider_endpoint_id", endpointID, err))
	}

	if err := d.Set("private_link_service_resource_id", serviceEndpoint.GetPrivateLinkServiceResourceId()); err != nil {
		return diag.FromErr(fmt.Errorf(privatelinkendpointservice.ErrorEndpointSetting, "private_link_service_resource_id", endpointID, err))
	}

	if err := d.Set("private_endpoint_ip_address", serviceEndpoint.GetPrivateEndpointIpAddress()); err != nil {
		return diag.FromErr(fmt.Errorf(privatelinkendpointservice.ErrorEndpointSetting, "private_endpoint_ip_address", endpointID, err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":    projectID,
		"instance_name": instanceName,
		"endpoint_id":   endpointID,
	}))

	return nil
}
