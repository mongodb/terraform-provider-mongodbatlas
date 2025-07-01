package privatelinkendpointservice

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
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"go.mongodb.org/atlas-sdk/v20250312004/admin"
)

const (
	errorServiceEndpointAdd  = "error adding MongoDB Private Service Endpoint Connection(%s) to a Private Endpoint (%s): %s"
	ErrorServiceEndpointRead = "error reading MongoDB Private Service Endpoint Connection(%s): %s"
	errorEndpointDelete      = "error deleting MongoDB Private Service Endpoint Connection(%s): %s"
	ErrorEndpointSetting     = "error setting `%s` for MongoDB Private Service Endpoint Connection(%s): %s"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext:      resourceCreate,
		ReadWithoutTimeout: resourceRead,
		DeleteContext:      resourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceImportState,
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
			"private_link_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"endpoint_service_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"interface_endpoint_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_endpoint_connection_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_endpoint_ip_address": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"gcp_project_id", "endpoints"},
			},
			"private_endpoint_resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"delete_requested": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"error_message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"aws_connection_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"azure_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gcp_project_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"private_endpoint_ip_address"},
			},
			"endpoints": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"private_endpoint_ip_address"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoint_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"gcp_status": {
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
	privateLinkID := conversion.GetEncodedID(d.Get("private_link_id").(string), "private_link_id")
	providerName := d.Get("provider_name").(string)
	endpointServiceID := d.Get("endpoint_service_id").(string)
	pEIA, pEIAOk := d.GetOk("private_endpoint_ip_address")
	gPI, gPIOk := d.GetOk("gcp_project_id")
	e, eOk := d.GetOk("endpoints")

	createEndpointRequest := &admin.CreateEndpointRequest{}

	switch providerName {
	case "AWS":
		createEndpointRequest.Id = &endpointServiceID
	case "AZURE":
		if !pEIAOk {
			return diag.FromErr(errors.New("`private_endpoint_ip_address` must be set when `provider_name` is `AZURE`"))
		}
		createEndpointRequest.Id = &endpointServiceID
		createEndpointRequest.PrivateEndpointIPAddress = conversion.Pointer(pEIA.(string))
	case "GCP":
		if !gPIOk || !eOk {
			return diag.FromErr(errors.New("`gcp_project_id`, `endpoints` must be set when `provider_name` is `GCP`"))
		}
		createEndpointRequest.EndpointGroupName = &endpointServiceID
		createEndpointRequest.GcpProjectId = conversion.Pointer(gPI.(string))
		createEndpointRequest.Endpoints = expandGCPEndpoints(e.([]any))
	}

	_, _, err := connV2.PrivateEndpointServicesApi.CreatePrivateEndpoint(ctx, projectID, providerName, privateLinkID, createEndpointRequest).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorServiceEndpointAdd, providerName, privateLinkID, err))
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"NONE", "INITIATING", "PENDING_ACCEPTANCE", "PENDING", "DELETING", "VERIFIED"},
		Target:     []string{"AVAILABLE", "REJECTED", "DELETED", "FAILED"},
		Refresh:    resourceRefreshFunc(ctx, connV2, projectID, providerName, privateLinkID, endpointServiceID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 5 * time.Second,
		Delay:      1 * time.Minute,
	}
	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorServiceEndpointAdd, endpointServiceID, privateLinkID, err))
	}

	clusterConf := &retry.StateChangeConf{
		Pending:    []string{"REPEATING", "PENDING"},
		Target:     []string{"IDLE", "DELETED"},
		Refresh:    advancedcluster.ResourceClusterListAdvancedRefreshFunc(ctx, projectID, connV2.ClustersApi),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 5 * time.Second,
		Delay:      1 * time.Minute,
	}

	if _, err = clusterConf.WaitForStateContext(ctx); err != nil {
		// error awaiting advanced clusters IDLE should not result in failure to apply changes to this resource
		log.Printf(advancedcluster.ErrorAdvancedClusterListStatus, err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":          projectID,
		"private_link_id":     privateLinkID,
		"endpoint_service_id": endpointServiceID,
		"provider_name":       providerName,
	}))

	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	privateLinkID := ids["private_link_id"]
	endpointServiceID := ids["endpoint_service_id"]
	providerName := ids["provider_name"]

	privateEndpoint, resp, err := connV2.PrivateEndpointServicesApi.GetPrivateEndpoint(context.Background(), projectID, providerName, endpointServiceID, privateLinkID).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf(ErrorServiceEndpointRead, endpointServiceID, err))
	}

	if err := d.Set("delete_requested", privateEndpoint.GetDeleteRequested()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorEndpointSetting, "delete_requested", endpointServiceID, err))
	}

	if err := d.Set("error_message", privateEndpoint.GetErrorMessage()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorEndpointSetting, "error_message", endpointServiceID, err))
	}

	if err := d.Set("aws_connection_status", privateEndpoint.GetConnectionStatus()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorEndpointSetting, "aws_connection_status", endpointServiceID, err))
	}

	if providerName == "AZURE" {
		if err := d.Set("azure_status", privateEndpoint.GetStatus()); err != nil {
			return diag.FromErr(fmt.Errorf(ErrorEndpointSetting, "azure_status", endpointServiceID, err))
		}
	}

	if err := d.Set("interface_endpoint_id", privateEndpoint.GetInterfaceEndpointId()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorEndpointSetting, "interface_endpoint_id", endpointServiceID, err))
	}

	if err := d.Set("private_endpoint_connection_name", privateEndpoint.GetPrivateEndpointConnectionName()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorEndpointSetting, "private_endpoint_connection_name", endpointServiceID, err))
	}

	if err := d.Set("private_endpoint_ip_address", privateEndpoint.GetPrivateEndpointIPAddress()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorEndpointSetting, "private_endpoint_ip_address", endpointServiceID, err))
	}

	if err := d.Set("private_endpoint_resource_id", privateEndpoint.GetPrivateEndpointResourceId()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorEndpointSetting, "private_endpoint_resource_id", endpointServiceID, err))
	}

	if err := d.Set("endpoint_service_id", endpointServiceID); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorEndpointSetting, "endpoint_service_id", endpointServiceID, err))
	}

	if err := d.Set("endpoints", flattenGCPEndpoints(privateEndpoint.Endpoints)); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorEndpointSetting, "endpoints", endpointServiceID, err))
	}

	if providerName == "GCP" {
		if err := d.Set("gcp_status", privateEndpoint.GetStatus()); err != nil {
			return diag.FromErr(fmt.Errorf(ErrorEndpointSetting, "gcp_status", endpointServiceID, err))
		}
	}

	if privateEndpoint.GetErrorMessage() != "" {
		return diag.FromErr(fmt.Errorf("privatelink endpoint service is in a failed state: %s", privateEndpoint.GetErrorMessage()))
	}
	return nil
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	privateLinkID := ids["private_link_id"]
	endpointServiceID := ids["endpoint_service_id"]
	providerName := ids["provider_name"]

	if endpointServiceID != "" {
		_, err := connV2.PrivateEndpointServicesApi.DeletePrivateEndpoint(ctx, projectID, providerName, endpointServiceID, privateLinkID).Execute()
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorEndpointDelete, endpointServiceID, err))
		}

		stateConf := &retry.StateChangeConf{
			Pending:    []string{"NONE", "PENDING_ACCEPTANCE", "PENDING", "DELETING", "INITIATING"},
			Target:     []string{"REJECTED", "DELETED", "FAILED"},
			Refresh:    resourceRefreshFunc(ctx, connV2, projectID, providerName, privateLinkID, endpointServiceID),
			Timeout:    d.Timeout(schema.TimeoutDelete),
			MinTimeout: 5 * time.Second,
			Delay:      3 * time.Second,
		}

		// Wait, catching any errors
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorEndpointDelete, endpointServiceID, err))
		}

		clusterConf := &retry.StateChangeConf{
			Pending:    []string{"REPEATING", "PENDING"},
			Target:     []string{"IDLE", "DELETED"},
			Refresh:    advancedcluster.ResourceClusterListAdvancedRefreshFunc(ctx, projectID, connV2.ClustersApi),
			Timeout:    d.Timeout(schema.TimeoutDelete),
			MinTimeout: 5 * time.Second,
			Delay:      1 * time.Minute,
		}

		if _, err = clusterConf.WaitForStateContext(ctx); err != nil {
			// error awaiting advanced clusters IDLE should not result in failure to apply changes to this resource
			log.Printf(advancedcluster.ErrorAdvancedClusterListStatus, err)
		}
	}

	return nil
}

func resourceImportState(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	parts := strings.SplitN(d.Id(), "--", 4)
	if len(parts) != 4 {
		return nil, errors.New("import format error: to import a MongoDB Private Endpoint, use the format {project_id}--{private_link_id}--{endpoint_service_id}--{provider_name}")
	}

	projectID := parts[0]
	privateLinkID := parts[1]
	endpointServiceID := parts[2]
	providerName := parts[3]

	_, _, err := connV2.PrivateEndpointServicesApi.GetPrivateEndpoint(ctx, projectID, providerName, endpointServiceID, privateLinkID).Execute()
	if err != nil {
		return nil, fmt.Errorf(ErrorServiceEndpointRead, endpointServiceID, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		return nil, fmt.Errorf(ErrorEndpointSetting, "project_id", privateLinkID, err)
	}

	if err := d.Set("private_link_id", privateLinkID); err != nil {
		return nil, fmt.Errorf(ErrorEndpointSetting, "private_link_id", privateLinkID, err)
	}

	if err := d.Set("endpoint_service_id", endpointServiceID); err != nil {
		return nil, fmt.Errorf(ErrorEndpointSetting, "endpoint_service_id", privateLinkID, err)
	}

	if err := d.Set("provider_name", providerName); err != nil {
		return nil, fmt.Errorf(ErrorEndpointSetting, "provider_name", privateLinkID, err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":          projectID,
		"private_link_id":     privateLinkID,
		"endpoint_service_id": endpointServiceID,
		"provider_name":       providerName,
	}))

	return []*schema.ResourceData{d}, nil
}

func resourceRefreshFunc(ctx context.Context, client *admin.APIClient, projectID, providerName, privateLinkID, endpointServiceID string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		i, resp, err := client.PrivateEndpointServicesApi.GetPrivateEndpoint(ctx, projectID, providerName, endpointServiceID, privateLinkID).Execute()
		if err != nil {
			if validate.StatusNotFound(resp) {
				return "", "DELETED", nil
			}

			return nil, "", err
		}

		if strings.EqualFold(providerName, "azure") || strings.EqualFold(providerName, "gcp") {
			if i.GetStatus() != "AVAILABLE" {
				return "", i.GetStatus(), nil
			}
			return i, *i.Status, nil
		}
		if i.GetConnectionStatus() != "AVAILABLE" {
			return "", i.GetConnectionStatus(), nil
		}

		return i, *i.ConnectionStatus, nil
	}
}

func expandGCPEndpoint(tfMap map[string]any) admin.CreateGCPForwardingRuleRequest {
	apiObject := admin.CreateGCPForwardingRuleRequest{}

	if v, ok := tfMap["endpoint_name"]; ok {
		apiObject.EndpointName = conversion.Pointer(v.(string))
	}
	if v, ok := tfMap["ip_address"]; ok {
		apiObject.IpAddress = conversion.Pointer(v.(string))
	}

	return apiObject
}

func expandGCPEndpoints(tfList []any) *[]admin.CreateGCPForwardingRuleRequest {
	if len(tfList) == 0 {
		return nil
	}

	var apiObjects []admin.CreateGCPForwardingRuleRequest

	for _, tfMapRaw := range tfList {
		if tfMap, ok := tfMapRaw.(map[string]any); ok {
			if tfMap != nil {
				apiObject := expandGCPEndpoint(tfMap)
				apiObjects = append(apiObjects, apiObject)
			}
		}
	}

	return &apiObjects
}

func flattenGCPEndpoint(apiObject admin.GCPConsumerForwardingRule) map[string]any {
	tfMap := map[string]any{}

	log.Printf("[DEBIG] apiObject : %+v", apiObject)

	tfMap["endpoint_name"] = apiObject.GetEndpointName()
	tfMap["ip_address"] = apiObject.GetIpAddress()
	tfMap["status"] = apiObject.GetStatus()

	return tfMap
}

func flattenGCPEndpoints(apiObjects *[]admin.GCPConsumerForwardingRule) []any {
	if apiObjects == nil || len(*apiObjects) == 0 {
		return nil
	}

	var tfList []any

	for _, apiObject := range *apiObjects {
		tfList = append(tfList, flattenGCPEndpoint(apiObject))
	}

	return tfList
}
