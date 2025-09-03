package flexsnapshot

import (
	"go.mongodb.org/atlas-sdk/v20250312007/admin"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

func NewTFModel(projectID, name string, apiResp *admin.FlexBackupSnapshot20241113) *TFModel {
	return &TFModel{
		ProjectId:      types.StringValue(projectID),
		Name:           types.StringValue(name),
		SnapshotId:     types.StringPointerValue(apiResp.Id),
		Expiration:     types.StringPointerValue(conversion.TimePtrToStringPtr(apiResp.Expiration)),
		FinishTime:     types.StringPointerValue(conversion.TimePtrToStringPtr(apiResp.FinishTime)),
		MongoDBVersion: types.StringPointerValue(apiResp.MongoDBVersion),
		ScheduledTime:  types.StringPointerValue(conversion.TimePtrToStringPtr(apiResp.ScheduledTime)),
		StartTime:      types.StringPointerValue(conversion.TimePtrToStringPtr(apiResp.StartTime)),
		Status:         types.StringPointerValue(apiResp.Status),
	}
}

func NewTFModelPluralDS(projectID, name string, apiResp *[]admin.FlexBackupSnapshot20241113) *TFFlexSnapshotsDSModel {
	if apiResp == nil {
		return nil
	}
	var results []TFModel
	for _, snapshot := range *apiResp {
		results = append(results, *NewTFModel(projectID, name, &snapshot))
	}
	return &TFFlexSnapshotsDSModel{
		ProjectId: types.StringValue(projectID),
		Name:      types.StringValue(name),
		Results:   results,
	}
}
