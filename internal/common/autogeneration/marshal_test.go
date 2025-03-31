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
	require.Equal(t, "value_date", model.CreateDate.ValueString())
}

type tfModelTest struct {
	BucketName types.String `tfsdk:"bucket_name"`
	CreateDate types.String `tfsdk:"create_date"`
}

const tfModelJSON = `
{
	"bucket_name": "value_bucket",
	"create_date": "value_date",
}
`
