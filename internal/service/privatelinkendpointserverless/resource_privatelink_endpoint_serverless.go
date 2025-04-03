package privatelinkendpointserverless

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/atlas-sdk/v20250312002/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/privatelinkendpoint"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/privatelinkendpointserviceserverless"
)

const (
	errorServerlessEndpointAdd    = "error adding MongoDB Serverless PrivateLink Endpoint Connection(%s): %s"
	errorServerlessEndpointDelete = "error deleting MongoDB Serverless PrivateLink Endpoint Connection(%s): %s"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate,
		ReadContext:   resourceRead,
		DeleteContext: resourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceImport,
		},
		DeprecationMessage: fmt.Sprintf(constant.DeprecationResourceByDateWithExternalLink, "March 2025", "https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/serverless-shared-migration-guide"),
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

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)
	instanceName := d.Get("instance_name").(string)

	endPoint, _, err := connV2.ServerlessPrivateEndpointsApi.CreateServerlessPrivateEndpoint(ctx, projectID, instanceName, &admin.ServerlessTenantCreateRequest{}).Execute()
	if err != nil {
		return diag.Errorf(privatelinkendpointserviceserverless.ErrorServerlessServiceEndpointAdd, endPoint.GetCloudProviderEndpointId(), err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"RESERVATION_REQUESTED", "INITIATING", "DELETING"},
		Target:     []string{"RESERVED", "FAILED", "DELETED", "AVAILABLE"},
		Refresh:    resourceRefreshFunc(ctx, connV2, projectID, instanceName, endPoint.GetId()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 5 * time.Second,
		Delay:      5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorServerlessEndpointAdd, err, endPoint.GetId()))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":    projectID,
		"instance_name": instanceName,
		"endpoint_id":   endPoint.GetId(),
	}))

	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	instanceName := ids["instance_name"]
	endpointID := ids["endpoint_id"]

	privateLinkResponse, _, err := connV2.ServerlessPrivateEndpointsApi.GetServerlessPrivateEndpoint(ctx, projectID, instanceName, endpointID).Execute()
	if err != nil {
		// case 404/400: deleted in the backend case
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "400") {
			d.SetId("")
			return nil
		}

		return diag.Errorf("error getting Serverless private link endpoint  information: %s", err)
	}

	if err := d.Set("endpoint_id", privateLinkResponse.GetId()); err != nil {
		return diag.Errorf("error setting `endpoint_id` for endpoint_id (%s): %s", d.Id(), err)
	}

	if err := d.Set("instance_name", instanceName); err != nil {
		return diag.Errorf("error setting `instance Name` for endpoint_id (%s): %s", d.Id(), err)
	}

	if err := d.Set("endpoint_service_name", privateLinkResponse.GetEndpointServiceName()); err != nil {
		return diag.Errorf("error setting `endpoint_service_name Name` for endpoint_id (%s): %s", d.Id(), err)
	}

	if err := d.Set("private_link_service_resource_id", privateLinkResponse.GetPrivateLinkServiceResourceId()); err != nil {
		return diag.Errorf("error setting `private_link_service_resource_id Name` for endpoint_id (%s): %s", d.Id(), err)
	}

	if err := d.Set("status", privateLinkResponse.GetStatus()); err != nil {
		return diag.FromErr(fmt.Errorf(privatelinkendpoint.ErrorPrivateLinkEndpointsSetting, "status", d.Id(), err))
	}

	return nil
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	instanceName := ids["instance_name"]
	endpointID := ids["endpoint_id"]

	_, _, err := connV2.ServerlessPrivateEndpointsApi.GetServerlessPrivateEndpoint(ctx, projectID, instanceName, endpointID).Execute()
	if err != nil {
		// case 404/400: deleted in the backend case
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "400") {
			d.SetId("")
			return nil
		}

		return diag.Errorf("error getting Serverless private link endpoint  information: %s", err)
	}

	_, _, err = connV2.ServerlessPrivateEndpointsApi.DeleteServerlessPrivateEndpoint(ctx, projectID, instanceName, endpointID).Execute()
	if err != nil {
		return diag.Errorf("error deleting serverless private link endpoint(%s): %s", endpointID, err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"DELETING"},
		Target:     []string{"DELETED", "FAILED"},
		Refresh:    resourceRefreshFunc(ctx, connV2, projectID, instanceName, endpointID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		MinTimeout: 5 * time.Second,
		Delay:      5 * time.Second,
	}
	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorServerlessEndpointDelete, endpointID, err))
	}

	return nil
}

func resourceImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	parts := strings.SplitN(d.Id(), "--", 3)
	if len(parts) != 3 {
		return nil, errors.New("import format error: to import a search index, use the format {project_id}--{instance_name}--{endpoint_id}")
	}

	projectID := parts[0]
	instanceName := parts[1]
	endpointID := parts[2]

	privateLinkResponse, _, err := connV2.ServerlessPrivateEndpointsApi.GetServerlessPrivateEndpoint(ctx, projectID, instanceName, endpointID).Execute()
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

	if err := d.Set("endpoint_service_name", privateLinkResponse.GetEndpointServiceName()); err != nil {
		log.Printf("[WARN] Error setting endpoint_service_name for (%s): %s", endpointID, err)
	}

	if privateLinkResponse.GetPrivateLinkServiceResourceId() != "" {
		if err := d.Set("provider_name", "AZURE"); err != nil {
			log.Printf("[WARN] Error setting provider_name for (%s): %s", endpointID, err)
		}
	} else {
		if err := d.Set("provider_name", "AWS"); err != nil {
			log.Printf("[WARN] Error setting provider_name for (%s): %s", endpointID, err)
		}
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":    projectID,
		"instance_name": instanceName,
		"endpoint_id":   endpointID,
	}))

	return []*schema.ResourceData{d}, nil
}

func resourceRefreshFunc(ctx context.Context, client *admin.APIClient, projectID, instanceName, privateLinkID string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		p, resp, err := client.ServerlessPrivateEndpointsApi.GetServerlessPrivateEndpoint(ctx, projectID, instanceName, privateLinkID).Execute()
		if err != nil {
			if validate.StatusNotFound(resp) || validate.StatusBadRequest(resp) {
				return "", "DELETED", nil
			}
			return nil, "REJECTED", err
		}

		status := p.GetStatus()

		if status != "WAITING_FOR_USER" {
			return "", status, nil
		}

		return p, status, nil
	}
}
