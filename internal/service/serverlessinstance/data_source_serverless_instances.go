package serverlessinstance

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	admin20241113 "go.mongodb.org/atlas-sdk/v20241113005/admin"
)

func PluralDataSource() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: fmt.Sprintf(constant.DeprecationDataSourceByDateWithExternalLink, constant.ServerlessSharedEOLDate, "https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/serverless-shared-migration-guide"),
		ReadContext:        dataSourcePluralRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: dataSourceSchema(),
				},
			},
		},
	}
}

func dataSourcePluralRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV220241113 := meta.(*config.MongoDBClient).AtlasV220241113
	projectIDValue, projectIDOK := d.GetOk("project_id")
	if !(projectIDOK) {
		return diag.Errorf("project_id must be configured")
	}
	projectID := projectIDValue.(string)
	options := &admin20241113.ListServerlessInstancesApiParams{
		ItemsPerPage: conversion.IntPtr(500),
		IncludeCount: new(true),
		GroupId:      projectID,
	}

	serverlessInstances, err := getServerlessList(ctx, connV220241113, options, 0)
	if err != nil {
		return diag.Errorf("error getting serverless instances information: %s", err)
	}

	autoIndexingList := make([]bool, len(serverlessInstances))
	for i := range serverlessInstances {
		resp, _, _ := connV220241113.PerformanceAdvisorApi.GetServerlessAutoIndexing(ctx, projectID, serverlessInstances[i].GetName()).Execute()
		autoIndexingList[i] = resp
	}

	flatServerlessInstances := flattenServerlessInstances(serverlessInstances, autoIndexingList)
	if err := d.Set("results", flatServerlessInstances); err != nil {
		return diag.Errorf("error setting `results` for serverless instances: %s", err)
	}

	d.SetId(id.UniqueId())
	return nil
}

func getServerlessList(ctx context.Context, connV220241113 *admin20241113.APIClient, options *admin20241113.ListServerlessInstancesApiParams, obtainedItemsCount int) ([]admin20241113.ServerlessInstanceDescription, error) {
	if options.PageNum == nil {
		options.PageNum = conversion.IntPtr(1)
	} else {
		*options.PageNum++
	}
	var list []admin20241113.ServerlessInstanceDescription
	serverlessInstances, _, err := connV220241113.ServerlessInstancesApi.ListServerlessInstancesWithParams(ctx, options).Execute()
	if err != nil {
		return list, fmt.Errorf("error getting serverless instances information: %s", err)
	}

	list = append(list, serverlessInstances.GetResults()...)
	obtainedItemsCount += len(serverlessInstances.GetResults())

	if serverlessInstances.GetTotalCount() > *options.ItemsPerPage && obtainedItemsCount < *serverlessInstances.TotalCount {
		instances, err := getServerlessList(ctx, connV220241113, options, obtainedItemsCount)
		if err != nil {
			return list, fmt.Errorf("error getting serverless instances information: %s", err)
		}
		list = append(list, instances...)
	}
	return list, nil
}

func flattenServerlessInstances(serverlessInstances []admin20241113.ServerlessInstanceDescription, autoIndexingList []bool) []map[string]any {
	var serverlessInstancesMap []map[string]any
	if len(serverlessInstances) == 0 {
		return nil
	}
	serverlessInstancesMap = make([]map[string]any, len(serverlessInstances))

	for i := range serverlessInstances {
		serverlessInstancesMap[i] = map[string]any{
			"connection_strings_standard_srv": serverlessInstances[i].ConnectionStrings.GetStandardSrv(),
			"create_date":                     conversion.TimePtrToStringPtr(serverlessInstances[i].CreateDate),
			"id":                              serverlessInstances[i].GetId(),
			"links":                           flattenLinks(serverlessInstances[i].GetLinks()),
			"mongo_db_version":                serverlessInstances[i].GetMongoDBVersion(),
			"name":                            serverlessInstances[i].GetName(),
			"provider_settings_backing_provider_name": serverlessInstances[i].ProviderSettings.GetBackingProviderName(),
			"provider_settings_region_name":           serverlessInstances[i].ProviderSettings.GetRegionName(),
			"provider_settings_provider_name":         serverlessInstances[i].ProviderSettings.GetProviderName(),
			"state_name":                              serverlessInstances[i].GetStateName(),
			"termination_protection_enabled":          serverlessInstances[i].GetTerminationProtectionEnabled(),
			"continuous_backup_enabled":               serverlessInstances[i].ServerlessBackupOptions.GetServerlessContinuousBackupEnabled(),
			"tags":                                    flattenTags(serverlessInstances[i].GetTags()),
			"auto_indexing":                           autoIndexingList[i],
		}
	}
	return serverlessInstancesMap
}
