package pushbasedlogexport

import (
	"context"

	"go.mongodb.org/atlas-sdk/v20250312004/admin"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

func NewTFPushBasedLogExport(ctx context.Context, projectID string, apiResp *admin.PushBasedLogExportProject, timeout *timeouts.Value) (*TFPushBasedLogExportRSModel, diag.Diagnostics) {
	tfModel := &TFPushBasedLogExportRSModel{
		ProjectID:  types.StringPointerValue(&projectID),
		BucketName: types.StringPointerValue(apiResp.BucketName),
		IamRoleID:  types.StringPointerValue(apiResp.IamRoleId),
		PrefixPath: types.StringPointerValue(apiResp.PrefixPath),
		CreateDate: types.StringPointerValue(conversion.TimePtrToStringPtr(apiResp.CreateDate)),
		State:      types.StringPointerValue(apiResp.State),
	}

	if timeout != nil {
		tfModel.Timeouts = *timeout
	}
	return tfModel, nil
}

func NewPushBasedLogExportCreateReq(plan *TFPushBasedLogExportRSModel) *admin.CreatePushBasedLogExportProjectRequest {
	return &admin.CreatePushBasedLogExportProjectRequest{
		BucketName: plan.BucketName.ValueString(),
		IamRoleId:  plan.IamRoleID.ValueString(),
		PrefixPath: plan.PrefixPath.ValueString(),
	}
}

func NewPushBasedLogExportUpdateReq(plan *TFPushBasedLogExportRSModel) *admin.PushBasedLogExportProject {
	return &admin.PushBasedLogExportProject{
		BucketName: plan.BucketName.ValueStringPointer(),
		IamRoleId:  plan.IamRoleID.ValueStringPointer(),
		PrefixPath: plan.PrefixPath.ValueStringPointer(),
	}
}
