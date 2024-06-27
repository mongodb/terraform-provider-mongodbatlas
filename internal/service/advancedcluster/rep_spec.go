package advancedcluster

import (
	"context"

	"go.mongodb.org/atlas-sdk/v20240530001/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// func FlattenAdvancedReplicationSpecsDS0710(ctx context.Context, apiRepSpecs []admin.ReplicationSpec20240710, connV2 *admin.APIClient) ([]map[string]any, error) {
// 	if len(apiRepSpecs) == 0 {
// 		return nil, nil
// 	}

// 	tfList := make([]map[string]any, len(apiRepSpecs))

// 	for i, apiRepSpec := range apiRepSpecs {
// 		tfListObj := map[string]any{
// 			"external_id": apiRepSpec.GetId(),
// 			"zone_id":     apiRepSpec.GetZoneId(),
// 			"zone_name":   apiRepSpec.GetZoneName(),
// 			// "region_configs":apiRepSpec.RegionConfigs,
// 			// "container_id":      flattenEndpoints(endpoint.GetEndpoints()),
// 		}
// 		object, containerIDs, err := flattenAdvancedReplicationSpecRegionConfigsDS0710(ctx, apiRepSpec.GetRegionConfigs(), connV2)
// 		if err != nil {
// 			return nil, err
// 		}

// 		tfListObj["region_configs"] = object
// 		tfListObj["container_id"] = containerIDs

// 		tfList[i] = tfListObj
// 	}
// 	return tfList, nil
// }

func flattenAdvancedReplicationSpec(ctx context.Context, apiObject *admin.ReplicationSpec20240710, tfMapObject map[string]any,
	d *schema.ResourceData, connV2 *admin.APIClient) (map[string]any, error) {
	if apiObject == nil {
		return nil, nil
	}

	tfMap := map[string]any{}
	tfMap["external_id"] = apiObject.GetId()
	if tfMapObject != nil {
		object, containerIDs, err := flattenAdvancedReplicationSpecRegionConfigs(ctx, apiObject.GetRegionConfigs(), tfMapObject["region_configs"].([]any), d, connV2)
		if err != nil {
			return nil, err
		}
		tfMap["region_configs"] = object
		tfMap["container_id"] = containerIDs
	} else {
		object, containerIDs, err := flattenAdvancedReplicationSpecRegionConfigs(ctx, apiObject.GetRegionConfigs(), nil, d, connV2)
		if err != nil {
			return nil, err
		}
		tfMap["region_configs"] = object
		tfMap["container_id"] = containerIDs
	}
	tfMap["zone_name"] = apiObject.GetZoneName()
	tfMap["zone_id"] = apiObject.GetZoneId()

	return tfMap, nil
}
