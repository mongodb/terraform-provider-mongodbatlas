package autogeneration_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogeneration"
	"github.com/stretchr/testify/require"
)

func TestUnmarshal(t *testing.T) {
	var model tfModelTest
	require.NoError(t, autogeneration.Unmarshal([]byte(tfModelJSON), &model))
	require.Equal(t, "value_bucket", model.BucketName.ValueString())
	require.Equal(t, "2023-10-01T00:00:00Z", model.CreateDate.ValueString())
	require.Equal(t, "value_group", model.GroupId.ValueString())
	require.Equal(t, "value_role", model.IamRoleId.ValueString())
	require.Equal(t, "value_prefix", model.PrefixPath.ValueString())
	require.Equal(t, "state", model.State.ValueString())
}

type tfModelTest struct {
	BucketName types.String `tfsdk:"bucket_name"`
	CreateDate types.String `tfsdk:"create_date"`
	GroupId    types.String `tfsdk:"group_id"`
	IamRoleId  types.String `tfsdk:"iam_role_id"`
	PrefixPath types.String `tfsdk:"prefix_path"`
	State      types.String `tfsdk:"state"`
}

const tfModelJSON = `
{
	"bucket_name": "value_bucket",
	"create_date": "2023-10-01T00:00:00Z",
	"group_id": "value_group",
	"iam_role_id": "value_role",
	"prefix_path": "value_prefix",
	"state": "state"
}
`
