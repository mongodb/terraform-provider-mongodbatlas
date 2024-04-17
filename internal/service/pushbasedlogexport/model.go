package pushbasedlogexport

import (
	"context"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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

	links, diagnostics := types.ListValueFrom(ctx, linkObjectType, NewTFLinksModel(apiResp.GetLinks()))
	if diagnostics.HasError() {
		return nil, diagnostics
	}

	tfModel.Links = links

	if timeout != nil {
		tfModel.Timeouts = *timeout
	}
	return tfModel, nil
}

func NewPushBasedLogExportReq(plan *TFPushBasedLogExportRSModel) *admin.PushBasedLogExportProject {
	return &admin.PushBasedLogExportProject{
		BucketName: plan.BucketName.ValueStringPointer(),
		IamRoleId:  plan.IamRoleID.ValueStringPointer(),
		PrefixPath: plan.PrefixPath.ValueStringPointer(),
	}
}

func NewTFLinksModel(links []admin.Link) []TFLinkModel {
	result := make([]TFLinkModel, len(links))
	for i, v := range links {
		result[i] = TFLinkModel{
			Href: types.StringPointerValue(v.Href),
			Rel:  types.StringPointerValue(v.Rel),
		}
	}

	return result
}

type TFLinkModel struct {
	Href types.String `tfsdk:"href"`
	Rel  types.String `tfsdk:"rel"`
}

var linkObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"href": types.StringType,
	"rel":  types.StringType,
}}
