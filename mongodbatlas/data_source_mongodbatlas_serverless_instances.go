package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasServerlessInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasServerlessInstancesRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: returnServerlessInstanceDSSchema(),
				},
			},
		},
	}
}

func dataSourceMongoDBAtlasServerlessInstancesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectID, projectIDOK := d.GetOk("project_id")

	if !(projectIDOK) {
		return diag.Errorf("project_id must be configured")
	}

	options := &matlas.ListOptions{
		ItemsPerPage: 500,
		IncludeCount: true,
	}

	serverlessInstances, err := getServerlessList(ctx, meta, projectID.(string), options, 0)
	if err != nil {
		return diag.Errorf("error getting serverless instances information: %s", err)
	}

	flatServerlessInstances := flattenServerlessInstances(serverlessInstances)

	if err := d.Set("results", flatServerlessInstances); err != nil {
		return diag.Errorf("error setting `results` for serverless instances: %s", err)
	}

	d.SetId(resource.UniqueId())

	return nil
}

func getServerlessList(ctx context.Context, meta interface{}, projectID string, options *matlas.ListOptions, obtainedItemsCount int) ([]*matlas.Cluster, error) {
	// Get client connection.
	var list []*matlas.Cluster
	options.PageNum++
	conn := meta.(*MongoDBClient).Atlas

	serverlessInstances, _, err := conn.ServerlessInstances.List(ctx, projectID, options)
	if err != nil {
		return list, fmt.Errorf("error getting serverless instances information: %s", err)
	}

	list = append(list, serverlessInstances.Results...)
	obtainedItemsCount += len(serverlessInstances.Results)

	if serverlessInstances.TotalCount > options.ItemsPerPage && obtainedItemsCount < serverlessInstances.TotalCount {
		instances, err := getServerlessList(ctx, meta, projectID, options, obtainedItemsCount)
		if err != nil {
			return list, fmt.Errorf("error getting serverless instances information: %s", err)
		}
		list = append(list, instances...)
	}

	return list, nil
}

func flattenServerlessInstances(serverlessInstances []*matlas.Cluster) []map[string]interface{} {
	var serverlessInstancesMap []map[string]interface{}

	if len(serverlessInstances) == 0 {
		return nil
	}
	serverlessInstancesMap = make([]map[string]interface{}, len(serverlessInstances))

	for i := range serverlessInstances {
		serverlessInstancesMap[i] = map[string]interface{}{
			"connection_strings_standard_srv": serverlessInstances[i].ConnectionStrings.StandardSrv,
			"create_date":                     serverlessInstances[i].CreateDate,
			"id":                              serverlessInstances[i].ID,
			"links":                           flattenLinks(serverlessInstances[i].Links),
			"mongo_db_version":                serverlessInstances[i].MongoDBVersion,
			"name":                            serverlessInstances[i].Name,
			"provider_settings_backing_provider_name": serverlessInstances[i].ProviderSettings.BackingProviderName,
			"provider_settings_region_name":           serverlessInstances[i].ProviderSettings.RegionName,
			"provider_settings_provider_name":         serverlessInstances[i].ProviderSettings.ProviderName,
			"state_name":                              serverlessInstances[i].StateName,
		}
	}

	return serverlessInstancesMap
}
