package federatedquerylimit

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312008/admin"
)

const (
	errorFederatedDatabaseQueryLimitCreate = "error creating MongoDB Atlas Federated Database Query Limit (%s): %s"
	errorFederatedDatabaseQueryLimitRead   = "error reading MongoDB Atlas Federated Database Query Limit (%s): %s"
	errorFederatedDatabaseQueryLimitDelete = "error deleting MongoDB Atlas Federated Database Query Limit (%s): %s"
	errorFederatedDatabaseQueryLimitUpdate = "error updating MongoDB Atlas Federated Database Query Limit (%s): %s"
	errorFederatedDatabaseQueryLimit       = "error setting `%s` for Atlas Federated Database Query Limit (%s): %s"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate,
		ReadContext:   resourceRead,
		UpdateContext: resourceUpdate,
		DeleteContext: resourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceImportState,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"limit_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"overrun_policy": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"current_usage": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"last_modified_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"maximum_limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return createOrUpdate(ctx, meta, d, errorFederatedDatabaseQueryLimitCreate)
}

func createOrUpdate(ctx context.Context, meta any, d *schema.ResourceData, errorTemplate string) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get(string("project_id")).(string)
	tenantName := d.Get("tenant_name").(string)
	limitName := d.Get("limit_name").(string)

	requestBody := &admin.DataFederationTenantQueryLimit{
		OverrunPolicy: conversion.StringPtr(d.Get("overrun_policy").(string)),
		Value:         int64(d.Get("value").(int)),
	}

	federatedDatabaseQueryLimit, _, err := conn.DataFederationApi.SetDataFederationLimit(ctx, projectID, tenantName, limitName, requestBody).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorTemplate, limitName, err))
	}
	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":  projectID,
		"tenant_name": federatedDatabaseQueryLimit.GetTenantName(),
		"limit_name":  federatedDatabaseQueryLimit.Name,
	}))

	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	tenantName := ids["tenant_name"]
	limitName := ids["limit_name"]

	queryLimit, resp, err := conn.DataFederationApi.GetDataFederationLimit(ctx, projectID, tenantName, limitName).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf(errorFederatedDatabaseQueryLimitRead, limitName, err))
	}

	if err = setResourceFieldsFromFederatedDatabaseQueryLimit(d, projectID, queryLimit); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":  projectID,
		"tenant_name": queryLimit.GetTenantName(),
		"limit_name":  queryLimit.Name,
	}))

	return nil
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return createOrUpdate(ctx, meta, d, errorFederatedDatabaseQueryLimitUpdate)
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2

	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	tenantName := ids["tenant_name"]
	limitName := ids["limit_name"]

	if _, err := conn.DataFederationApi.DeleteDataFederationLimit(ctx, projectID, tenantName, limitName).Execute(); err != nil {
		return diag.FromErr(fmt.Errorf(errorFederatedDatabaseQueryLimitDelete, limitName, err))
	}

	return nil
}

func resourceImportState(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	conn := meta.(*config.MongoDBClient).AtlasV2
	parts := strings.Split(d.Id(), "--")

	var projectID, tenantName, limitName string

	if len(parts) != 3 {
		return nil, errors.New("import format error: to import a MongoDB Atlas Federated Database Query Limit, use the format {project_id}-{tenant_name}-{limit_name}")
	}
	projectID, tenantName, limitName = parts[0], parts[1], parts[2]

	queryLimit, _, err := conn.DataFederationApi.GetDataFederationLimit(ctx, projectID, tenantName, limitName).Execute()

	if err != nil {
		return nil, fmt.Errorf("couldn't import federated database query limit(%s) for project (%s), tenant (%s), error: %s", limitName, projectID, tenantName, err)
	}

	if err := setResourceFieldsFromFederatedDatabaseQueryLimit(d, projectID, queryLimit); err != nil {
		return nil, err
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":  projectID,
		"tenant_name": queryLimit.GetTenantName(),
		"limit_name":  queryLimit.Name,
	}))

	return []*schema.ResourceData{d}, nil
}

func setResourceFieldsFromFederatedDatabaseQueryLimit(d *schema.ResourceData, projectID string, queryLimit *admin.DataFederationTenantQueryLimit) error {
	if err := d.Set("project_id", projectID); err != nil {
		return fmt.Errorf(errorFederatedDatabaseQueryLimit, "project_id", d.Id(), err)
	}

	if err := d.Set("limit_name", queryLimit.Name); err != nil {
		return fmt.Errorf(errorFederatedDatabaseQueryLimit, "limit_name", d.Id(), err)
	}

	if err := d.Set("tenant_name", queryLimit.GetTenantName()); err != nil {
		return fmt.Errorf(errorFederatedDatabaseQueryLimit, "tenant_name", d.Id(), err)
	}

	if err := d.Set("overrun_policy", queryLimit.GetOverrunPolicy()); err != nil {
		return fmt.Errorf(errorFederatedDatabaseQueryLimit, "overrun_policy", d.Id(), err)
	}

	if err := d.Set("value", queryLimit.Value); err != nil {
		return fmt.Errorf(errorFederatedDatabaseQueryLimit, "value", d.Id(), err)
	}

	if err := d.Set("current_usage", queryLimit.GetCurrentUsage()); err != nil {
		return fmt.Errorf(errorFederatedDatabaseQueryLimit, "current_usage", d.Id(), err)
	}

	if err := d.Set("default_limit", queryLimit.GetDefaultLimit()); err != nil {
		return fmt.Errorf(errorFederatedDatabaseQueryLimit, "default_limit", d.Id(), err)
	}

	if err := d.Set("last_modified_date", conversion.TimeToString(queryLimit.GetLastModifiedDate())); err != nil {
		return fmt.Errorf(errorFederatedDatabaseQueryLimit, "last_modified_date", d.Id(), err)
	}

	if err := d.Set("maximum_limit", queryLimit.GetMaximumLimit()); err != nil {
		return fmt.Errorf(errorFederatedDatabaseQueryLimit, "maximum_limit", d.Id(), err)
	}

	return nil
}
