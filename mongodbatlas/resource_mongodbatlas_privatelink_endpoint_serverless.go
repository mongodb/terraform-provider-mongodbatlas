package mongodbatlas

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
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorServerlessEndpointAdd    = "error adding MongoDB Serverless PrivateLink Endpoint Connection(%s): %s"
	errorServerlessEndpointDelete = "error deleting MongoDB Serverless PrivateLink Endpoint Connection(%s): %s"
)

func resourceMongoDBAtlasPrivateLinkEndpointServerless() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBAtlasPrivateLinkEndpointServerlessCreate,
		ReadContext:   resourceMongoDBAtlasPrivateLinkEndpointServerlessRead,
		DeleteContext: resourceMongoDBAtlasPrivateLinkEndpointServerlessDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasPrivateLinkEndpointServerlessImportState,
		},
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
				Computed: true,
			},
			"provider_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"AWS", "AZURE"}, false),
			},
			"endpoint_service_name": {
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
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Hour),
			Delete: schema.DefaultTimeout(2 * time.Hour),
		},
	}
}

func resourceMongoDBAtlasPrivateLinkEndpointServerlessCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)
	instanceName := d.Get("instance_name").(string)

	privateLinkRequest := &matlas.ServerlessPrivateEndpointConnection{
		Comment: "create",
	}

	endPoint, _, err := conn.ServerlessPrivateEndpoints.Create(ctx, projectID, instanceName, privateLinkRequest)
	if err != nil {
		return diag.Errorf(errorServerlessServiceEndpointAdd, privateLinkRequest.CloudProviderEndpointID, err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"RESERVATION_REQUESTED", "INITIATING", "DELETING"},
		Target:     []string{"RESERVED", "FAILED", "DELETED", "AVAILABLE"},
		Refresh:    resourcePrivateLinkEndpointServerlessRefreshFunc(ctx, conn, projectID, instanceName, endPoint.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 5 * time.Second,
		Delay:      3 * time.Second,
	}
	// RESERVATION_REQUESTED, RESERVED, INITIATING, AVAILABLE, FAILED, DELETING.
	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorServerlessEndpointAdd, err, endPoint.ID))
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":    projectID,
		"instance_name": instanceName,
		"endpoint_id":   endPoint.ID,
	}))

	return resourceMongoDBAtlasPrivateLinkEndpointServerlessRead(ctx, d, meta)
}

func resourceMongoDBAtlasPrivateLinkEndpointServerlessRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	instanceName := ids["instance_name"]
	endpointID := ids["endpoint_id"]

	privateLinkResponse, _, err := conn.ServerlessPrivateEndpoints.Get(ctx, projectID, instanceName, endpointID)
	if err != nil {
		// case 404/ 400
		// deleted in the backend case
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "400") {
			d.SetId("")
			return nil
		}

		return diag.Errorf("error getting Serverless private link endpoint  information: %s", err)
	}

	if err := d.Set("endpoint_id", privateLinkResponse.ID); err != nil {
		return diag.Errorf("error setting `endpoint_id` for endpoint_id (%s): %s", d.Id(), err)
	}

	if err := d.Set("instance_name", instanceName); err != nil {
		return diag.Errorf("error setting `instance Name` for endpoint_id (%s): %s", d.Id(), err)
	}

	if err := d.Set("endpoint_service_name", privateLinkResponse.EndpointServiceName); err != nil {
		return diag.Errorf("error setting `endpoint_service_name Name` for endpoint_id (%s): %s", d.Id(), err)
	}

	if err := d.Set("private_link_service_resource_id", privateLinkResponse.PrivateLinkServiceResourceID); err != nil {
		return diag.Errorf("error setting `private_link_service_resource_id Name` for endpoint_id (%s): %s", d.Id(), err)
	}

	if err := d.Set("status", privateLinkResponse.Status); err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsSetting, "status", d.Id(), err))
	}

	return nil
}

func resourceMongoDBAtlasPrivateLinkEndpointServerlessDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	instanceName := ids["instance_name"]
	endpointID := ids["endpoint_id"]

	_, _, err := conn.ServerlessPrivateEndpoints.Get(ctx, projectID, instanceName, endpointID)
	if err != nil {
		// case 404
		// deleted in the backend case
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "400") {
			d.SetId("")
			return nil
		}

		return diag.Errorf("error getting Serverless private link endpoint  information: %s", err)
	}

	_, err = conn.ServerlessPrivateEndpoints.Delete(ctx, projectID, instanceName, endpointID)
	if err != nil {
		return diag.Errorf("error deleting serverless private link endpoint(%s): %s", endpointID, err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"DELETING"},
		Target:     []string{"DELETED", "FAILED"},
		Refresh:    resourcePrivateLinkEndpointServerlessRefreshFunc(ctx, conn, projectID, instanceName, endpointID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		MinTimeout: 5 * time.Second,
		Delay:      3 * time.Second,
	}
	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorServerlessEndpointDelete, endpointID, err))
	}

	return nil
}

func resourceMongoDBAtlasPrivateLinkEndpointServerlessImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*MongoDBClient).Atlas

	parts := strings.SplitN(d.Id(), "--", 3)
	if len(parts) != 3 {
		return nil, errors.New("import format error: to import a search index, use the format {project_id}--{instance_name}--{endpoint_id}")
	}

	projectID := parts[0]
	instanceName := parts[1]
	endpointID := parts[2]

	privateLinkResponse, _, err := conn.ServerlessPrivateEndpoints.Get(ctx, projectID, instanceName, endpointID)
	if err != nil {
		return nil, fmt.Errorf("couldn't import serverless private link endpoint (%s) in projectID (%s) , error: %s", endpointID, projectID, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		log.Printf("[WARN] Error setting project_id for (%s): %s", projectID, err)
	}

	if err := d.Set("endpoint_id", endpointID); err != nil {
		log.Printf("[WARN] Error setting endpoint_id for (%s): %s", endpointID, err)
	}
	if err := d.Set("instance_name", instanceName); err != nil {
		log.Printf("[WARN] Error setting instance_name for (%s): %s", endpointID, err)
	}

	if err := d.Set("endpoint_service_name", privateLinkResponse.EndpointServiceName); err != nil {
		log.Printf("[WARN] Error setting endpoint_service_name for (%s): %s", endpointID, err)
	}

	if privateLinkResponse.PrivateLinkServiceResourceID != "" {
		if err := d.Set("provider_name", "AZURE"); err != nil {
			log.Printf("[WARN] Error setting provider_name for (%s): %s", endpointID, err)
		}
	} else {
		if err := d.Set("provider_name", "AWS"); err != nil {
			log.Printf("[WARN] Error setting provider_name for (%s): %s", endpointID, err)
		}
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":    projectID,
		"instance_name": instanceName,
		"endpoint_id":   endpointID,
	}))

	return []*schema.ResourceData{d}, nil
}

func resourcePrivateLinkEndpointServerlessRefreshFunc(ctx context.Context, client *matlas.Client, projectID, instanceName, privateLinkID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		p, resp, err := client.ServerlessPrivateEndpoints.Get(ctx, projectID, instanceName, privateLinkID)
		if err != nil {
			if resp.Response.StatusCode == 404 || resp.Response.StatusCode == 400 {
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
