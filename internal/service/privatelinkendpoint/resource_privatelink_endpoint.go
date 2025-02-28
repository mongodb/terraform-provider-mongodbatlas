package privatelinkendpoint

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

const (
	errorPrivateLinkEndpointsCreate  = "error creating MongoDB Private Endpoints Connection: %s"
	errorPrivateLinkEndpointsRead    = "error reading MongoDB Private Endpoints Connection(%s): %s"
	errorPrivateLinkEndpointsDelete  = "error deleting MongoDB Private Endpoints Connection(%s): %s"
	ErrorPrivateLinkEndpointsSetting = "error setting `%s` for MongoDB Private Endpoints Connection(%s): %s"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate,
		ReadContext:   resourceRead,
		DeleteContext: resourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceImport,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"provider_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"AWS", "AZURE", "GCP"}, false),
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"private_link_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint_service_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"interface_endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"private_endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"private_link_service_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_link_service_resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint_group_names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"region_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_attachment_names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Hour),
			Delete: schema.DefaultTimeout(1 * time.Hour),
		},
	}
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)
	providerName := d.Get("provider_name").(string)
	region := d.Get("region").(string)

	request := &admin.CloudProviderEndpointServiceRequest{
		ProviderName: providerName,
		Region:       region,
	}

	privateEndpoint, _, err := connV2.PrivateEndpointServicesApi.CreatePrivateEndpointService(ctx, projectID, request).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsCreate, err))
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"INITIATING", "DELETING"},
		Target:     []string{"WAITING_FOR_USER", "FAILED", "DELETED", "AVAILABLE"},
		Refresh:    refreshFunc(ctx, connV2, projectID, providerName, privateEndpoint.GetId()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 5 * time.Second,
		Delay:      3 * time.Second,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsCreate, err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"private_link_id": privateEndpoint.GetId(),
		"project_id":      projectID,
		"provider_name":   providerName,
		"region":          region,
	}))

	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	privateLinkID := ids["private_link_id"]
	providerName := ids["provider_name"]
	region := ids["region"]

	privateEndpoint, resp, err := connV2.PrivateEndpointServicesApi.GetPrivateEndpointService(context.Background(), projectID, providerName, privateLinkID).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsRead, privateLinkID, err))
	}

	if err := d.Set("private_link_id", privateEndpoint.GetId()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorPrivateLinkEndpointsSetting, "private_link_id", privateLinkID, err))
	}

	if err := d.Set("endpoint_service_name", privateEndpoint.GetEndpointServiceName()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorPrivateLinkEndpointsSetting, "endpoint_service_name", privateLinkID, err))
	}

	if err := d.Set("error_message", privateEndpoint.GetErrorMessage()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorPrivateLinkEndpointsSetting, "error_message", privateLinkID, err))
	}

	if err := d.Set("interface_endpoints", privateEndpoint.GetInterfaceEndpoints()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorPrivateLinkEndpointsSetting, "interface_endpoints", privateLinkID, err))
	}

	if err := d.Set("private_endpoints", privateEndpoint.GetPrivateEndpoints()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorPrivateLinkEndpointsSetting, "interface_endpoints", privateLinkID, err))
	}

	if err := d.Set("private_link_service_name", privateEndpoint.GetPrivateLinkServiceName()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorPrivateLinkEndpointsSetting, "interface_endpoints", privateLinkID, err))
	}

	if err := d.Set("private_link_service_resource_id", privateEndpoint.GetPrivateLinkServiceResourceId()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorPrivateLinkEndpointsSetting, "interface_endpoints", privateLinkID, err))
	}

	if err := d.Set("status", privateEndpoint.GetStatus()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorPrivateLinkEndpointsSetting, "status", privateLinkID, err))
	}

	if err := d.Set("provider_name", providerName); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorPrivateLinkEndpointsSetting, "provider_name", privateLinkID, err))
	}

	if err := d.Set("region", region); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorPrivateLinkEndpointsSetting, "region", privateLinkID, err))
	}

	if err := d.Set("endpoint_group_names", privateEndpoint.GetEndpointGroupNames()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorPrivateLinkEndpointsSetting, "endpoint_group_names", privateLinkID, err))
	}

	if err := d.Set("region_name", privateEndpoint.GetRegionName()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorPrivateLinkEndpointsSetting, "region_name", privateLinkID, err))
	}

	if err := d.Set("service_attachment_names", privateEndpoint.GetServiceAttachmentNames()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorPrivateLinkEndpointsSetting, "service_attachment_names", privateLinkID, err))
	}

	if privateEndpoint.GetErrorMessage() != "" {
		return diag.FromErr(fmt.Errorf("privatelink endpoint is in a failed state: %s", privateEndpoint.GetErrorMessage()))
	}
	return nil
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	ids := conversion.DecodeStateID(d.Id())
	privateLinkID := ids["private_link_id"]
	projectID := ids["project_id"]
	providerName := ids["provider_name"]

	_, resp, err := connV2.PrivateEndpointServicesApi.DeletePrivateEndpointService(ctx, projectID, providerName, privateLinkID).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) {
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsDelete, privateLinkID, err))
	}

	log.Println("[INFO] Waiting for MongoDB Private Endpoints Connection to be destroyed")

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"DELETING"},
		Target:     []string{"DELETED", "FAILED"},
		Refresh:    refreshFunc(ctx, connV2, projectID, providerName, privateLinkID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		MinTimeout: 5 * time.Second,
		Delay:      3 * time.Second,
	}
	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsDelete, privateLinkID, err))
	}

	return nil
}

func resourceImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	conn := meta.(*config.MongoDBClient).Atlas

	parts := strings.Split(d.Id(), "-")

	if len(parts) != 6 && len(parts) != 5 && len(parts) != 4 {
		return nil, errors.New("import format error: to import a MongoDB Private Endpoint, use the format {project_id}-{private_link_id}-{provider_name}-{region}")
	}

	projectID := parts[0]
	privateLinkID := parts[1]
	providerName := parts[2]
	region := parts[3] // If region it's azure or Atlas format like US_EAST_1
	if len(parts) == 5 {
		region = fmt.Sprintf("%s-%s", parts[3], parts[4])
	}
	if len(parts) == 6 {
		region = fmt.Sprintf("%s-%s-%s", parts[3], parts[4], parts[5])
	}

	privateEndpoint, _, err := conn.PrivateEndpoints.Get(ctx, projectID, providerName, privateLinkID)
	if err != nil {
		return nil, fmt.Errorf("couldn't import peer %s in project %s with cloud provider name %s, error: %s", privateLinkID, projectID, providerName, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		log.Printf(ErrorPrivateLinkEndpointsSetting, "project_id", privateLinkID, err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"private_link_id": privateEndpoint.ID,
		"project_id":      projectID,
		"provider_name":   providerName,
		"region":          region,
	}))

	return []*schema.ResourceData{d}, nil
}

func refreshFunc(ctx context.Context, client *admin.APIClient, projectID, providerName, privateLinkID string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		p, resp, err := client.PrivateEndpointServicesApi.GetPrivateEndpointService(ctx, projectID, providerName, privateLinkID).Execute()
		if err != nil {
			if validate.StatusNotFound(resp) {
				return "", "DELETED", nil
			}

			return nil, "REJECTED", err
		}

		status := ""
		if _, ok := p.GetStatusOk(); ok {
			status = p.GetStatus()
		}

		if status != "WAITING_FOR_USER" {
			return "", status, nil
		}

		return p, status, nil
	}
}
