package advancedclustertpf

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

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
