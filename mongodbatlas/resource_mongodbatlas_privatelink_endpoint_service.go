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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spf13/cast"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorServiceEndpointAdd  = "error adding MongoDB Private Service Endpoint Connection(%s) to a Private Endpoint (%s): %s"
	errorServiceEndpointRead = "error reading MongoDB Private Service Endpoint Connection(%s): %s"
	errorEndpointDelete      = "error deleting MongoDB Private Service Endpoint Connection(%s): %s"
	errorEndpointSetting     = "error setting `%s` for MongoDB Private Service Endpoint Connection(%s): %s"
)

func resourceMongoDBAtlasPrivateEndpointServiceLink() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceMongoDBAtlasPrivateEndpointServiceLinkCreate,
		ReadWithoutTimeout:   resourceMongoDBAtlasPrivateEndpointServiceLinkRead,
		DeleteWithoutTimeout: resourceMongoDBAtlasPrivateEndpointServiceLinkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasPrivateEndpointServiceLinkImportState,
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
						"service_attachment_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"gcp_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceMongoDBAtlasPrivateEndpointServiceLinkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)
	privateLinkID := getEncodedID(d.Get("private_link_id").(string), "private_link_id")
	providerName := d.Get("provider_name").(string)
	endpointServiceID := d.Get("endpoint_service_id").(string)
	pEIA, pEIAOk := d.GetOk("private_endpoint_ip_address")
	gPI, gPIOk := d.GetOk("gcp_project_id")
	e, eOk := d.GetOk("endpoints")

	request := &matlas.InterfaceEndpointConnection{}

	switch providerName {
	case "AWS":
		request.ID = endpointServiceID
	case "AZURE":
		if !pEIAOk {
			return diag.FromErr(errors.New("`private_endpoint_ip_address` must be set when `provider_name` is `AZURE`"))
		}
		request.ID = endpointServiceID
		request.PrivateEndpointIPAddress = pEIA.(string)
	case "GCP":
		if !gPIOk || !eOk {
			return diag.FromErr(errors.New("`gcp_project_id`, `endpoints` must be set when `provider_name` is `GCP`"))
		}
		request.EndpointGroupName = endpointServiceID
		request.GCPProjectID = gPI.(string)
		request.Endpoints = expandGCPEndpoints(e.([]interface{}))
	}

	_, _, err := conn.PrivateEndpoints.AddOnePrivateEndpoint(ctx, projectID, providerName, privateLinkID, request)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorServiceEndpointAdd, providerName, privateLinkID, err))
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"NONE", "INITIATING", "PENDING_ACCEPTANCE", "PENDING", "DELETING", "VERIFIED"},
		Target:     []string{"AVAILABLE", "REJECTED", "DELETED", "FAILED"},
		Refresh:    resourceServiceEndpointRefreshFunc(ctx, conn, projectID, providerName, privateLinkID, endpointServiceID),
		Timeout:    1 * time.Hour,
		MinTimeout: 5 * time.Second,
		Delay:      5 * time.Minute,
	}
	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorServiceEndpointAdd, endpointServiceID, privateLinkID, err))
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":          projectID,
		"private_link_id":     privateLinkID,
		"endpoint_service_id": endpointServiceID,
		"provider_name":       providerName,
	}))

	return resourceMongoDBAtlasPrivateEndpointServiceLinkRead(ctx, d, meta)
}

func resourceMongoDBAtlasPrivateEndpointServiceLinkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	privateLinkID := ids["private_link_id"]
	endpointServiceID := ids["endpoint_service_id"]
	providerName := ids["provider_name"]

	privateEndpoint, resp, err := conn.PrivateEndpoints.GetOnePrivateEndpoint(context.Background(), projectID, providerName, privateLinkID, endpointServiceID)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorServiceEndpointRead, endpointServiceID, err))
	}

	if err := d.Set("delete_requested", cast.ToBool(privateEndpoint.DeleteRequested)); err != nil {
		return diag.FromErr(fmt.Errorf(errorEndpointSetting, "delete_requested", endpointServiceID, err))
	}

	if err := d.Set("error_message", privateEndpoint.ErrorMessage); err != nil {
		return diag.FromErr(fmt.Errorf(errorEndpointSetting, "error_message", endpointServiceID, err))
	}

	if err := d.Set("aws_connection_status", privateEndpoint.AWSConnectionStatus); err != nil {
		return diag.FromErr(fmt.Errorf(errorEndpointSetting, "aws_connection_status", endpointServiceID, err))
	}

	if providerName == "AZURE" {
		if err := d.Set("azure_status", privateEndpoint.Status); err != nil {
			return diag.FromErr(fmt.Errorf(errorEndpointSetting, "azure_status", endpointServiceID, err))
		}
	}

	if err := d.Set("interface_endpoint_id", privateEndpoint.InterfaceEndpointID); err != nil {
		return diag.FromErr(fmt.Errorf(errorEndpointSetting, "connection_status", endpointServiceID, err))
	}

	if err := d.Set("private_endpoint_connection_name", privateEndpoint.PrivateEndpointConnectionName); err != nil {
		return diag.FromErr(fmt.Errorf(errorEndpointSetting, "connection_status", endpointServiceID, err))
	}

	if err := d.Set("private_endpoint_ip_address", privateEndpoint.PrivateEndpointIPAddress); err != nil {
		return diag.FromErr(fmt.Errorf(errorEndpointSetting, "connection_status", endpointServiceID, err))
	}

	if err := d.Set("private_endpoint_resource_id", privateEndpoint.PrivateEndpointResourceID); err != nil {
		return diag.FromErr(fmt.Errorf(errorEndpointSetting, "connection_status", endpointServiceID, err))
	}

	if err := d.Set("endpoint_service_id", endpointServiceID); err != nil {
		return diag.FromErr(fmt.Errorf(errorEndpointSetting, "endpoint_service_id", endpointServiceID, err))
	}

	if err := d.Set("endpoints", flattenGCPEndpoints(privateEndpoint.Endpoints)); err != nil {
		return diag.FromErr(fmt.Errorf(errorEndpointSetting, "endpoints", endpointServiceID, err))
	}

	if providerName == "GCP" {
		if err := d.Set("gcp_status", privateEndpoint.Status); err != nil {
			return diag.FromErr(fmt.Errorf(errorEndpointSetting, "gcp_status", endpointServiceID, err))
		}
	}

	return nil
}

func resourceMongoDBAtlasPrivateEndpointServiceLinkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	privateLinkID := ids["private_link_id"]
	endpointServiceID := ids["endpoint_service_id"]
	providerName := ids["provider_name"]

	if endpointServiceID != "" {
		_, err := conn.PrivateEndpoints.DeleteOnePrivateEndpoint(ctx, projectID, providerName, privateLinkID, endpointServiceID)
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorEndpointDelete, endpointServiceID, err))
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"NONE", "PENDING_ACCEPTANCE", "PENDING", "DELETING", "INITIATING"},
			Target:     []string{"REJECTED", "DELETED", "FAILED"},
			Refresh:    resourceServiceEndpointRefreshFunc(ctx, conn, projectID, providerName, privateLinkID, endpointServiceID),
			Timeout:    1 * time.Hour,
			MinTimeout: 5 * time.Second,
			Delay:      3 * time.Second,
		}

		// Wait, catching any errors
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorEndpointDelete, endpointServiceID, err))
		}
	}

	return nil
}

func resourceMongoDBAtlasPrivateEndpointServiceLinkImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*MongoDBClient).Atlas

	parts := strings.SplitN(d.Id(), "--", 4)
	if len(parts) != 4 {
		return nil, errors.New("import format error: to import a MongoDB Private Endpoint, use the format {project_id}--{private_link_id}--{endpoint_service_id}--{provider_name}")
	}

	projectID := parts[0]
	privateLinkID := parts[1]
	endpointServiceID := parts[2]
	providerName := parts[3]

	_, _, err := conn.PrivateEndpoints.GetOnePrivateEndpoint(ctx, projectID, providerName, privateLinkID, endpointServiceID)
	if err != nil {
		return nil, fmt.Errorf(errorServiceEndpointRead, endpointServiceID, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		return nil, fmt.Errorf(errorEndpointSetting, "project_id", privateLinkID, err)
	}

	if err := d.Set("private_link_id", privateLinkID); err != nil {
		return nil, fmt.Errorf(errorEndpointSetting, "private_link_id", privateLinkID, err)
	}

	if err := d.Set("endpoint_service_id", endpointServiceID); err != nil {
		return nil, fmt.Errorf(errorEndpointSetting, "endpoint_service_id", privateLinkID, err)
	}

	if err := d.Set("provider_name", providerName); err != nil {
		return nil, fmt.Errorf(errorEndpointSetting, "provider_name", privateLinkID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":          projectID,
		"private_link_id":     privateLinkID,
		"endpoint_service_id": endpointServiceID,
		"provider_name":       providerName,
	}))

	return []*schema.ResourceData{d}, nil
}

func resourceServiceEndpointRefreshFunc(ctx context.Context, client *matlas.Client, projectID, providerName, privateLinkID, endpointServiceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		i, resp, err := client.PrivateEndpoints.GetOnePrivateEndpoint(ctx, projectID, providerName, privateLinkID, endpointServiceID)
		if err != nil {
			if resp != nil && resp.StatusCode == 404 {
				return "", "DELETED", nil
			}

			return nil, "", err
		}

		if strings.EqualFold(providerName, "azure") || strings.EqualFold(providerName, "gcp") {
			if i.Status != "AVAILABLE" {
				return "", i.Status, nil
			}
			return i, i.Status, nil
		}
		if i.AWSConnectionStatus != "AVAILABLE" {
			return "", i.AWSConnectionStatus, nil
		}

		return i, i.AWSConnectionStatus, nil
	}
}

func expandGCPEndpoint(tfMap map[string]interface{}) *matlas.GCPEndpoint {
	if tfMap == nil {
		return nil
	}

	apiObject := &matlas.GCPEndpoint{}

	if v, ok := tfMap["endpoint_name"]; ok {
		apiObject.EndpointName = cast.ToString(v)
	}
	if v, ok := tfMap["ip_address"]; ok {
		apiObject.IPAddress = cast.ToString(v)
	}

	return apiObject
}

func expandGCPEndpoints(tfList []interface{}) []*matlas.GCPEndpoint {
	if len(tfList) == 0 {
		return nil
	}

	var apiObjects []*matlas.GCPEndpoint

	for _, tfMapRaw := range tfList {
		if tfMap, ok := tfMapRaw.(map[string]interface{}); ok {
			apiObject := expandGCPEndpoint(tfMap)
			if apiObject == nil {
				continue
			}
			apiObjects = append(apiObjects, apiObject)
		}
	}

	return apiObjects
}

func flattenGCPEndpoint(apiObject *matlas.GCPEndpoint) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	log.Printf("[DEBIG] apiObject : %+v", apiObject)

	tfMap["endpoint_name"] = apiObject.EndpointName
	tfMap["ip_address"] = apiObject.IPAddress
	tfMap["status"] = apiObject.Status
	tfMap["service_attachment_name"] = apiObject.ServiceAttachmentName

	return tfMap
}

func flattenGCPEndpoints(apiObjects []*matlas.GCPEndpoint) []interface{} {
	if len(apiObjects) == 0 {
		return nil
	}

	var tfList []interface{}

	for _, apiObject := range apiObjects {
		if apiObject == nil {
			continue
		}

		tfList = append(tfList, flattenGCPEndpoint(apiObject))
	}

	return tfList
}
