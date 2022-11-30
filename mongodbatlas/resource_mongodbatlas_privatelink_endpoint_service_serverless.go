package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorServerlessServiceEndpointAdd    = "error adding MongoDB Serverless PrivateLink Endpoint Connection(%s): %s"
	errorServerlessServiceEndpointDelete = "error deleting MongoDB Serverless PrivateLink Endpoint Connection(%s): %s"
)

func resourceMongoDBAtlasPrivateLinkEndpointServiceServerless() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceMongoDBAtlasPrivateLinkEndpointServiceServerlessCreate,
		ReadWithoutTimeout:   resourceMongoDBAtlasPrivateLinkEndpointServiceServerlessRead,
		DeleteWithoutTimeout: resourceMongoDBAtlasPrivateLinkEndpointServiceServerlessDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasPrivateLinkEndpointServiceServerlessImportState,
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
				Required: true,
				ForceNew: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"provider_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"AWS", "AZURE"}, false),
				ForceNew:     true,
			},
			"cloud_provider_endpoint_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"private_link_service_resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_endpoint_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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

func resourceMongoDBAtlasPrivateLinkEndpointServiceServerlessCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)
	instanceName := d.Get("instance_name").(string)
	endpointID := d.Get("endpoint_id").(string)

	privateLink, _, err := conn.ServerlessPrivateEndpoints.Get(ctx, projectID, instanceName, endpointID)
	if err != nil {
		return diag.Errorf("error getting Serverless PrivateLink Endpoint Information: %s", err)
	}

	privateLink.Comment = d.Get("comment").(string)
	privateLink.CloudProviderEndpointID = d.Get("cloud_provider_endpoint_id").(string)
	privateLink.ProviderName = d.Get("provider_name").(string)
	privateLink.PrivateLinkServiceResourceID = ""
	privateLink.PrivateEndpointIPAddress = d.Get("private_endpoint_ip_address").(string)
	privateLink.ID = ""
	privateLink.Status = ""
	privateLink.EndpointServiceName = ""

	endPoint, _, err := conn.ServerlessPrivateEndpoints.Update(ctx, projectID, instanceName, endpointID, privateLink)
	if err != nil {
		return diag.Errorf(errorServerlessServiceEndpointAdd, endpointID, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"RESERVATION_REQUESTED", "INITIATING", "DELETING"},
		Target:     []string{"RESERVED", "FAILED", "DELETED", "AVAILABLE"},
		Refresh:    resourceServiceEndpointServerlessRefreshFunc(ctx, conn, projectID, instanceName, endpointID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 5 * time.Second,
		Delay:      5 * time.Minute,
	}
	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorServerlessServiceEndpointAdd, endpointID, err))
	}

	clusterConf := &resource.StateChangeConf{
		Pending:    []string{"REPEATING", "PENDING"},
		Target:     []string{"IDLE", "DELETED"},
		Refresh:    resourceServerlessInstanceListRefreshFunc(ctx, projectID, conn),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 5 * time.Second,
		Delay:      5 * time.Minute,
	}

	_, err = clusterConf.WaitForStateContext(ctx)

	if err != nil {
		// error awaiting advanced clusters IDLE should not result in failure to apply changes to this resource
		log.Printf(errorAdvancedClusterListStatus, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":    projectID,
		"instance_name": instanceName,
		"endpoint_id":   endPoint.ID,
	}))

	return resourceMongoDBAtlasPrivateLinkEndpointServiceServerlessRead(ctx, d, meta)
}

func resourceMongoDBAtlasPrivateLinkEndpointServiceServerlessRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	instanceName := ids["instance_name"]
	endpointID := ids["endpoint_id"]

	privateLinkResponse, _, err := conn.ServerlessPrivateEndpoints.Get(ctx, projectID, instanceName, endpointID)
	if err != nil {
		// case 404
		// deleted in the backend case
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "400") {
			d.SetId("")
			return nil
		}

		return diag.Errorf("error getting Serverless private link endpoint  information: %s", err)
	}

	privateLinkResponse.ProviderName = d.Get("provider_name").(string)

	if err := d.Set("endpoint_id", privateLinkResponse.ID); err != nil {
		return diag.Errorf("error setting `endpoint_id` for endpoint_id (%s): %s", d.Id(), err)
	}

	if err := d.Set("instance_name", instanceName); err != nil {
		return diag.Errorf("error setting `instance Name` for endpoint_id (%s): %s", d.Id(), err)
	}

	if err := d.Set("comment", privateLinkResponse.Comment); err != nil {
		return diag.Errorf("error setting `comment` for endpoint_id (%s): %s", d.Id(), err)
	}

	if err := d.Set("provider_name", privateLinkResponse.ProviderName); err != nil {
		return diag.Errorf("error setting `provider_name` for endpoint_id (%s): %s", d.Id(), err)
	}

	if err := d.Set("status", privateLinkResponse.Status); err != nil {
		return diag.Errorf("error setting `status` for endpoint_id (%s): %s", d.Id(), err)
	}

	if err := d.Set("cloud_provider_endpoint_id", privateLinkResponse.CloudProviderEndpointID); err != nil {
		return diag.Errorf("error setting `cloud_provider_endpoint_id` for endpoint_id (%s): %s", d.Id(), err)
	}

	if err := d.Set("private_link_service_resource_id", privateLinkResponse.PrivateLinkServiceResourceID); err != nil {
		return diag.Errorf("error setting `private_link_service_resource_id` for endpoint_id (%s): %s", d.Id(), err)
	}

	if err := d.Set("private_endpoint_ip_address", privateLinkResponse.PrivateEndpointIPAddress); err != nil {
		return diag.Errorf("error setting `private_endpoint_ip_address` for endpoint_id (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceMongoDBAtlasPrivateLinkEndpointServiceServerlessDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

		return diag.Errorf(errorServerlessServiceEndpointDelete, endpointID, err)
	}

	d.SetId("") // Set to null as linked resource will delete servless endpoint

	return nil
}

func resourceMongoDBAtlasPrivateLinkEndpointServiceServerlessImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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

func resourceServiceEndpointServerlessRefreshFunc(ctx context.Context, client *matlas.Client, projectID, instanceName, endpointServiceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		i, resp, err := client.ServerlessPrivateEndpoints.Get(ctx, projectID, instanceName, endpointServiceID)
		if err != nil {
			if resp != nil && resp.StatusCode == 404 || resp.Response.StatusCode == 400 {
				return "", "DELETED", nil
			}

			return nil, "", err
		}

		if i.Status != "AVAILABLE" {
			return "", i.Status, nil
		}

		return i, i.Status, nil
	}
}
