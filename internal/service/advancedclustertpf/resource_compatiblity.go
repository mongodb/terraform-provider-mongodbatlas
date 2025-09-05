package advancedclustertpf

import (
	"context"
	"fmt"
	"strconv"

	"go.mongodb.org/atlas-sdk/v20250312007/admin"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

type MajorVersionOperator int

const (
	EqualOrHigher MajorVersionOperator = iota
	Higher
	EqualOrLower
	Lower
)

func MajorVersionCompatible(input *string, version float64, operator MajorVersionOperator) *bool {
	if !conversion.IsStringPresent(input) {
		return nil
	}
	value, err := strconv.ParseFloat(*input, 64)
	if err != nil {
		return nil
	}
	var result bool
	switch operator {
	case EqualOrHigher:
		result = value >= version
	case Higher:
		result = value > version
	case EqualOrLower:
		result = value <= version
	case Lower:
		result = value < version
	default:
		return nil
	}
	return &result
}

func containerIDKey(providerName, regionName string) string {
	return fmt.Sprintf("%s:%s", providerName, regionName)
}

// based on flattenAdvancedReplicationSpecRegionConfigs in model_advanced_cluster.go
func resolveContainerIDs(ctx context.Context, projectID string, cluster *admin.ClusterDescription20240805, api admin.NetworkPeeringApi) (map[string]string, error) {
	containerIDs := map[string]string{}
	responseCache := map[string]*admin.PaginatedCloudProviderContainer{}
	for _, spec := range cluster.GetReplicationSpecs() {
		for _, regionConfig := range spec.GetRegionConfigs() {
			providerName := regionConfig.GetProviderName()
			if providerName == constant.TENANT {
				continue
			}
			params := &admin.ListGroupContainersApiParams{
				GroupId:      projectID,
				ProviderName: &providerName,
			}
			key := containerIDKey(providerName, regionConfig.GetRegionName())
			if _, ok := containerIDs[key]; ok {
				continue
			}
			var containersResponse *admin.PaginatedCloudProviderContainer
			var err error
			if response, ok := responseCache[providerName]; ok {
				containersResponse = response
			} else {
				containersResponse, _, err = api.ListGroupContainersWithParams(ctx, params).Execute()
				if err != nil {
					return nil, err
				}
				responseCache[providerName] = containersResponse
			}
			if results := GetAdvancedClusterContainerID(containersResponse.GetResults(), &regionConfig); results != "" {
				containerIDs[key] = results
			} else {
				return nil, fmt.Errorf("container id not found for %s", key)
			}
		}
	}
	return containerIDs, nil
}

func overrideAttributesWithPrevStateValue(modelIn, modelOut *TFModel) {
	beforeVersion := conversion.NilForUnknown(modelIn.MongoDBMajorVersion, modelIn.MongoDBMajorVersion.ValueStringPointer())
	if beforeVersion != nil && !modelIn.MongoDBMajorVersion.Equal(modelOut.MongoDBMajorVersion) {
		modelOut.MongoDBMajorVersion = types.StringPointerValue(beforeVersion)
	}
	retainBackups := conversion.NilForUnknown(modelIn.RetainBackupsEnabled, modelIn.RetainBackupsEnabled.ValueBoolPointer())
	if retainBackups != nil && !modelIn.RetainBackupsEnabled.Equal(modelOut.RetainBackupsEnabled) {
		modelOut.RetainBackupsEnabled = types.BoolPointerValue(retainBackups)
	}
	if modelIn.DeleteOnCreateTimeout.ValueBoolPointer() != nil {
		modelOut.DeleteOnCreateTimeout = modelIn.DeleteOnCreateTimeout
	}
	overrideMapStringWithPrevStateValue(&modelIn.Labels, &modelOut.Labels)
	overrideMapStringWithPrevStateValue(&modelIn.Tags, &modelOut.Tags)
}

func overrideMapStringWithPrevStateValue(mapIn, mapOut *types.Map) {
	if mapIn == nil || mapOut == nil || len(mapOut.Elements()) > 0 {
		return
	}
	if mapIn.IsNull() {
		*mapOut = types.MapNull(types.StringType)
	} else {
		*mapOut = types.MapValueMust(types.StringType, nil)
	}
}
