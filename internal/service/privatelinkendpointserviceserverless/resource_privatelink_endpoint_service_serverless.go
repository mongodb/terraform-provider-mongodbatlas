package privatelinkendpointserviceserverless

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/atlas-sdk/v20250312001/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	ErrorServerlessServiceEndpointAdd    = "error adding MongoDB Serverless PrivateLink Endpoint Connection(%s): %s"
	errorServerlessServiceEndpointDelete = "error deleting MongoDB Serverless PrivateLink Endpoint Connection(%s): %s"
	errorServerlessInstanceListStatus    = "error awaiting serverless instance list status IDLE: %s"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceCreate,
		ReadWithoutTimeout:   resourceRead,
		DeleteWithoutTimeout: resourceDelete,
		UpdateWithoutTimeout: resourceUpdate,
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
				Required: true,
				ForceNew: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
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

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)
	instanceName := d.Get("instance_name").(string)
	endpointID := d.Get("endpoint_id").(string)

	_, _, err := connV2.ServerlessPrivateEndpointsApi.GetServerlessPrivateEndpoint(ctx, projectID, instanceName, endpointID).Execute()
	if err != nil {
		return diag.Errorf("error getting Serverless PrivateLink Endpoint Information: %s", err)
	}

	updateRequest := admin.ServerlessTenantEndpointUpdate{
		Comment:                  conversion.StringPtr(d.Get("comment").(string)),
		ProviderName:             d.Get("provider_name").(string),
		CloudProviderEndpointId:  conversion.StringPtr(d.Get("cloud_provider_endpoint_id").(string)),
		PrivateEndpointIpAddress: conversion.StringPtr(d.Get("private_endpoint_ip_address").(string)),
	}

	endPoint, _, err := connV2.ServerlessPrivateEndpointsApi.UpdateServerlessPrivateEndpoint(ctx, projectID, instanceName, endpointID, &updateRequest).Execute()
	if err != nil {
		return diag.Errorf(ErrorServerlessServiceEndpointAdd, endpointID, err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"RESERVATION_REQUESTED", "INITIATING", "DELETING"},
		Target:     []string{"RESERVED", "FAILED", "DELETED", "AVAILABLE"},
		Refresh:    resourceRefreshFunc(ctx, connV2, projectID, instanceName, endpointID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 5 * time.Second,
		Delay:      5 * time.Minute,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf(ErrorServerlessServiceEndpointAdd, endpointID, err))
	}

	clusterConf := &retry.StateChangeConf{
		Pending:    []string{"REPEATING", "PENDING"},
		Target:     []string{"IDLE", "DELETED"},
		Refresh:    resourceListRefreshFunc(ctx, projectID, connV2),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 5 * time.Second,
		Delay:      5 * time.Minute,
	}

	if _, err = clusterConf.WaitForStateContext(ctx); err != nil {
		// error awaiting serverless instances to IDLE should not result in failure to apply changes to this resource
		log.Printf(errorServerlessInstanceListStatus, err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":    projectID,
		"instance_name": instanceName,
		"endpoint_id":   endPoint.GetId(),
	}))

	return resourceRead(ctx, d, meta)
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	if !d.HasChange("comment") {
		return resourceRead(ctx, d, meta)
	}

	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)
	instanceName := d.Get("instance_name").(string)
	endpointID := d.Get("endpoint_id").(string)

	// only "comment" attribute update is supported, updating other attributes forces replacement of this resource
	updateRequest := admin.ServerlessTenantEndpointUpdate{
		Comment:      conversion.StringPtr(d.Get("comment").(string)),
		ProviderName: d.Get("provider_name").(string),
	}

	endPoint, _, err := connV2.ServerlessPrivateEndpointsApi.UpdateServerlessPrivateEndpoint(ctx, projectID, instanceName, endpointID, &updateRequest).Execute()
	if err != nil {
		return diag.Errorf(ErrorServerlessServiceEndpointAdd, endpointID, err)
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
		// case 404: deleted in the backend case
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "400") {
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting Serverless private link endpoint  information: %s", err)
	}

	privateLinkResponse.ProviderName = conversion.StringPtr(d.Get("provider_name").(string))

	if err := d.Set("endpoint_id", privateLinkResponse.GetId()); err != nil {
		return diag.Errorf("error setting `endpoint_id` for endpoint_id (%s): %s", d.Id(), err)
	}

	if err := d.Set("instance_name", instanceName); err != nil {
		return diag.Errorf("error setting `instance Name` for endpoint_id (%s): %s", d.Id(), err)
	}

	if err := d.Set("comment", privateLinkResponse.GetComment()); err != nil {
		return diag.Errorf("error setting `comment` for endpoint_id (%s): %s", d.Id(), err)
	}

	if err := d.Set("provider_name", privateLinkResponse.GetProviderName()); err != nil {
		return diag.Errorf("error setting `provider_name` for endpoint_id (%s): %s", d.Id(), err)
	}

	if err := d.Set("status", privateLinkResponse.GetStatus()); err != nil {
		return diag.Errorf("error setting `status` for endpoint_id (%s): %s", d.Id(), err)
	}

	if err := d.Set("cloud_provider_endpoint_id", privateLinkResponse.GetCloudProviderEndpointId()); err != nil {
		return diag.Errorf("error setting `cloud_provider_endpoint_id` for endpoint_id (%s): %s", d.Id(), err)
	}

	if err := d.Set("private_link_service_resource_id", privateLinkResponse.GetPrivateLinkServiceResourceId()); err != nil {
		return diag.Errorf("error setting `private_link_service_resource_id` for endpoint_id (%s): %s", d.Id(), err)
	}

	if err := d.Set("private_endpoint_ip_address", privateLinkResponse.GetPrivateEndpointIpAddress()); err != nil {
		return diag.Errorf("error setting `private_endpoint_ip_address` for endpoint_id (%s): %s", d.Id(), err)
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
		// case 404: deleted in the backend case
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "400") {
			d.SetId("")
			return nil
		}
		return diag.Errorf(errorServerlessServiceEndpointDelete, endpointID, err)
	}

	d.SetId("") // Set to null as linked resource will delete servless endpoint

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

func resourceRefreshFunc(ctx context.Context, client *admin.APIClient, projectID, instanceName, endpointServiceID string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		serverlessTenantEndpoint, resp, err := client.ServerlessPrivateEndpointsApi.GetServerlessPrivateEndpoint(ctx, projectID, instanceName, endpointServiceID).Execute()
		if err != nil {
			if validate.StatusNotFound(resp) || validate.StatusBadRequest(resp) {
				return "", "DELETED", nil
			}

			return nil, "", err
		}

		if serverlessTenantEndpoint.GetStatus() != "AVAILABLE" {
			return "", serverlessTenantEndpoint.GetStatus(), nil
		}
		resultStatus := serverlessTenantEndpoint.GetStatus()

		return serverlessTenantEndpoint, resultStatus, nil
	}
}

func resourceListRefreshFunc(ctx context.Context, projectID string, client *admin.APIClient) retry.StateRefreshFunc {
	return func() (any, string, error) {
		serverlessInstances, resp, err := client.ServerlessInstancesApi.ListServerlessInstances(ctx, projectID).Execute()

		if err != nil && strings.Contains(err.Error(), "reset by peer") {
			return nil, "REPEATING", nil
		}

		if err != nil && serverlessInstances == nil && resp == nil {
			return nil, "", err
		} else if err != nil {
			if validate.StatusNotFound(resp) {
				return "", "DELETED", nil
			}
			if validate.StatusServiceUnavailable(resp) {
				return "", "PENDING", nil
			}
			return nil, "", err
		}

		for i := range serverlessInstances.GetResults() {
			if serverlessInstances.GetResults()[i].GetStateName() != "IDLE" {
				return serverlessInstances.GetResults()[i], "PENDING", nil
			}
		}

		return serverlessInstances, "IDLE", nil
	}
}
