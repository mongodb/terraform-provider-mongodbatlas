package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorPrivateLinkEndpointsCreate  = "error creating MongoDB Private Endpoints Connection: %s"
	errorPrivateLinkEndpointsRead    = "error reading MongoDB Private Endpoints Connection(%s): %s"
	errorPrivateLinkEndpointsDelete  = "error deleting MongoDB Private Endpoints Connection(%s): %s"
	errorPrivateLinkEndpointsSetting = "error setting `%s` for MongoDB Private Endpoints Connection(%s): %s"
)

func resourceMongoDBAtlasPrivateLinkEndpoint() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBAtlasPrivateLinkEndpointCreate,
		ReadContext:   resourceMongoDBAtlasPrivateLinkEndpointRead,
		DeleteContext: resourceMongoDBAtlasPrivateLinkEndpointDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasPrivateLinkEndpointImportState,
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
	}
}

func resourceMongoDBAtlasPrivateLinkEndpointCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)
	providerName := d.Get("provider_name").(string)
	region := d.Get("region").(string)

	request := &matlas.PrivateEndpointConnection{
		ProviderName: providerName,
		Region:       region,
	}

	privateEndpointConn, _, err := conn.PrivateEndpoints.Create(ctx, projectID, request)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsCreate, err))
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"INITIATING", "DELETING"},
		Target:     []string{"WAITING_FOR_USER", "FAILED", "DELETED", "AVAILABLE"},
		Refresh:    resourcePrivateLinkEndpointRefreshFunc(ctx, conn, projectID, providerName, privateEndpointConn.ID),
		Timeout:    1 * time.Hour,
		MinTimeout: 5 * time.Second,
		Delay:      3 * time.Second,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsCreate, err))
	}

	d.SetId(encodeStateID(map[string]string{
		"private_link_id": privateEndpointConn.ID,
		"project_id":      projectID,
		"provider_name":   providerName,
		"region":          region,
	}))

	return resourceMongoDBAtlasPrivateLinkEndpointRead(ctx, d, meta)
}

func resourceMongoDBAtlasPrivateLinkEndpointRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	privateLinkID := ids["private_link_id"]
	providerName := ids["provider_name"]
	region := ids["region"]

	privateEndpoint, resp, err := conn.PrivateEndpoints.Get(context.Background(), projectID, providerName, privateLinkID)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsRead, privateLinkID, err))
	}

	if err := d.Set("private_link_id", privateEndpoint.ID); err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsSetting, "private_link_id", privateLinkID, err))
	}

	if err := d.Set("endpoint_service_name", privateEndpoint.EndpointServiceName); err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsSetting, "endpoint_service_name", privateLinkID, err))
	}

	if err := d.Set("error_message", privateEndpoint.ErrorMessage); err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsSetting, "error_message", privateLinkID, err))
	}

	if err := d.Set("interface_endpoints", privateEndpoint.InterfaceEndpoints); err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsSetting, "interface_endpoints", privateLinkID, err))
	}

	if err := d.Set("private_endpoints", privateEndpoint.PrivateEndpoints); err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsSetting, "interface_endpoints", privateLinkID, err))
	}

	if err := d.Set("private_link_service_name", privateEndpoint.PrivateLinkServiceName); err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsSetting, "interface_endpoints", privateLinkID, err))
	}

	if err := d.Set("private_link_service_resource_id", privateEndpoint.PrivateLinkServiceResourceID); err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsSetting, "interface_endpoints", privateLinkID, err))
	}

	if err := d.Set("status", privateEndpoint.Status); err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsSetting, "status", privateLinkID, err))
	}

	if err := d.Set("provider_name", providerName); err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsSetting, "provider_name", privateLinkID, err))
	}

	if err := d.Set("region", region); err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsSetting, "region", privateLinkID, err))
	}

	if err := d.Set("endpoint_group_names", privateEndpoint.EndpointGroupNames); err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsSetting, "endpoint_group_names", privateLinkID, err))
	}

	if err := d.Set("region_name", privateEndpoint.RegionName); err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsSetting, "region_name", privateLinkID, err))
	}

	if err := d.Set("service_attachment_names", privateEndpoint.ServiceAttachmentNames); err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsSetting, "service_attachment_names", privateLinkID, err))
	}

	return nil
}

func resourceMongoDBAtlasPrivateLinkEndpointDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	ids := decodeStateID(d.Id())
	privateLinkID := ids["private_link_id"]
	projectID := ids["project_id"]
	providerName := ids["provider_name"]

	resp, err := conn.PrivateEndpoints.Delete(ctx, projectID, providerName, privateLinkID)
	if err != nil {
		if resp.Response.StatusCode == 404 {
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsDelete, privateLinkID, err))
	}

	log.Println("[INFO] Waiting for MongoDB Private Endpoints Connection to be destroyed")

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"DELETING"},
		Target:     []string{"DELETED", "FAILED"},
		Refresh:    resourcePrivateLinkEndpointRefreshFunc(ctx, conn, projectID, providerName, privateLinkID),
		Timeout:    1 * time.Hour,
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

func resourceMongoDBAtlasPrivateLinkEndpointImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*MongoDBClient).Atlas

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
		log.Printf(errorPrivateLinkEndpointsSetting, "project_id", privateLinkID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"private_link_id": privateEndpoint.ID,
		"project_id":      projectID,
		"provider_name":   providerName,
		"region":          region,
	}))

	return []*schema.ResourceData{d}, nil
}

func resourcePrivateLinkEndpointRefreshFunc(ctx context.Context, client *matlas.Client, projectID, providerName, privateLinkID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		p, resp, err := client.PrivateEndpoints.Get(ctx, projectID, providerName, privateLinkID)
		if err != nil {
			if resp.Response.StatusCode == 404 {
				return "", "DELETED", nil
			}

			return nil, "REJECTED", err
		}

		if p.Status != "WAITING_FOR_USER" {
			return "", p.Status, nil
		}

		return p, p.Status, nil
	}
}
