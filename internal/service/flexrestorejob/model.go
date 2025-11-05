package flexrestorejob

import (
	"go.mongodb.org/atlas-sdk/v20250312009/admin"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

func NewTFModel(apiResp *admin.FlexBackupRestoreJob20241113) *TFModel {
	return &TFModel{
		ProjectID:                types.StringPointerValue(apiResp.ProjectId),
		Name:                     types.StringPointerValue(apiResp.InstanceName),
		DeliveryType:             types.StringPointerValue(apiResp.DeliveryType),
		ExpirationDate:           types.StringPointerValue(conversion.TimePtrToStringPtr(apiResp.ExpirationDate)),
		RestoreJobID:             types.StringPointerValue(apiResp.Id),
		RestoreFinishedDate:      types.StringPointerValue(conversion.TimePtrToStringPtr(apiResp.RestoreFinishedDate)),
		RestoreScheduledDate:     types.StringPointerValue(conversion.TimePtrToStringPtr(apiResp.RestoreScheduledDate)),
		SnapshotFinishedDate:     types.StringPointerValue(conversion.TimePtrToStringPtr(apiResp.SnapshotFinishedDate)),
		SnapshotID:               types.StringPointerValue(apiResp.SnapshotId),
		SnapshotUrl:              types.StringPointerValue(apiResp.SnapshotUrl),
		Status:                   types.StringPointerValue(apiResp.Status),
		TargetDeploymentItemName: types.StringPointerValue(apiResp.TargetDeploymentItemName),
		TargetProjectID:          types.StringPointerValue(apiResp.TargetProjectId),
	}
}

func NewTFModelPluralDS(projectID, name string, apiResp *[]admin.FlexBackupRestoreJob20241113) *TFFlexRestoreJobsDSModel {
	if apiResp == nil {
		return nil
	}
	var results []TFModel
	for _, job := range *apiResp {
		results = append(results, *NewTFModel(&job))
	}
	return &TFFlexRestoreJobsDSModel{
		ProjectID: types.StringValue(projectID),
		Name:      types.StringValue(name),
		Results:   results,
	}
}
