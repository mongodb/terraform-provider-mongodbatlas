package pushbasedlogexport

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TFPushBasedLogExportDSModel struct {
	BucketName types.String `tfsdk:"bucket_name"`
	CreateDate types.String `tfsdk:"create_date"`
	ProjectID  types.String `tfsdk:"project_id"`
	IamRoleID  types.String `tfsdk:"iam_role_id"`
	PrefixPath types.String `tfsdk:"prefix_path"`
	State      types.String `tfsdk:"state"`
}
