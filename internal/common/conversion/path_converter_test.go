package conversion_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

func TestConvertingAttributePathsToTPFPaths(t *testing.T) {
	var testCases = map[string]struct {
		in  *tftypes.AttributePath
		out path.Path
	}{
		"root attribute": {
			in: tftypes.NewAttributePathWithSteps(
				[]tftypes.AttributePathStep{tftypes.AttributeName("root")},
			),
			out: path.Root("root"),
		},
		"root attribute with string key": {
			in: tftypes.NewAttributePathWithSteps(
				[]tftypes.AttributePathStep{
					tftypes.AttributeName("tags"),
					tftypes.ElementKeyString("Name"),
				},
			),
			out: path.Root("tags").AtMapKey("Name"),
		},
		"root attribute with int key": {
			in: tftypes.NewAttributePathWithSteps(
				[]tftypes.AttributePathStep{
					tftypes.AttributeName("specs"),
					tftypes.ElementKeyInt(0),
				},
			),
			out: path.Root("specs").AtListIndex(0),
		},
		"root attribute with value key": {
			in: tftypes.NewAttributePathWithSteps(
				[]tftypes.AttributePathStep{
					tftypes.AttributeName("specs"),
				},
			).WithElementKeyValue(
				tftypes.NewValue(tftypes.String, "value"),
			),
			out: path.Root("specs").AtSetValue(types.StringValue("value")),
		},
		"nested attribute": {
			in: tftypes.NewAttributePathWithSteps(
				[]tftypes.AttributePathStep{
					tftypes.AttributeName("specs"),
					tftypes.AttributeName("instance_size"),
				},
			),
			out: path.Root("specs").AtName("instance_size"),
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			out, err := conversion.ConvertAttributePath(*testCase.in)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if !out.Equal(testCase.out) {
				t.Fatalf("expected %s, got %s", testCase.out.String(), out.String())
			}
		})
	}
}
