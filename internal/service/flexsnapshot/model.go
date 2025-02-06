package flexsnapshot

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"

	"go.mongodb.org/atlas-sdk/v20241113004/admin"
)

func NewTFModel(projectId, name string, apiResp *admin.FlexBackupSnapshot20241113) *TFModel {
	return &TFModel{
		ProjectId:      types.StringValue(projectId),
		Name:           types.StringValue(name),
		SnapshotId:     types.StringPointerValue(apiResp.Id),
		Expiration:     types.StringPointerValue(conversion.TimePtrToStringPtr(apiResp.Expiration)),
		FinishTime:     types.StringPointerValue(conversion.TimePtrToStringPtr(apiResp.FinishTime)),
		MongoDbversion: types.StringPointerValue(apiResp.MongoDBVersion),
		ScheduledTime:  types.StringPointerValue(conversion.TimePtrToStringPtr(apiResp.ScheduledTime)),
		StartTime:      types.StringPointerValue(conversion.TimePtrToStringPtr(apiResp.StartTime)),
		Status:         types.StringPointerValue(apiResp.Status),
	}
}
