package projectipaddresses_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/projectipaddresses"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20240805002/admin"
)

const (
	dummyProjectID = "111111111111111111111111"
)

type sdkToTFModelTestCase struct {
	SDKResp         *admin.GroupIPAddresses
	expectedTFModel *projectipaddresses.ProjectIpAddressesModel
}

func TestProjectIPAddressesSDKToTFModel(t *testing.T) {
	testCases := map[string]sdkToTFModelTestCase{
		"Complete SDK response": {
			SDKResp: &admin.GroupIPAddresses{
				GroupId: admin.PtrString(dummyProjectID),
				Services: &admin.GroupService{
					Clusters: &[]admin.ClusterIPAddresses{
						{
							ClusterName: admin.PtrString("cluster1"),
							Inbound:     &[]string{"192.168.1.1", "192.168.1.2"},
							Outbound:    &[]string{"10.0.0.1", "10.0.0.2"},
						},
						{
							ClusterName: admin.PtrString("cluster2"),
							Inbound:     &[]string{"192.168.2.1"},
							Outbound:    &[]string{"10.0.1.1"},
						},
					},
				},
			},
			expectedTFModel: &projectipaddresses.ProjectIpAddressesModel{
				ProjectId: types.StringValue(dummyProjectID),
				Services: types.ObjectValueMust(projectipaddresses.ServicesObjectType.AttrTypes, map[string]attr.Value{
					"clusters": types.ListValueMust(types.ObjectType{AttrTypes: projectipaddresses.ClusterIPsObjectType.AttrTypes}, []attr.Value{
						types.ObjectValueMust(projectipaddresses.ClusterIPsObjectType.AttrTypes, map[string]attr.Value{
							"cluster_name": types.StringValue("cluster1"),
							"inbound":      toTFList(t, []string{"192.168.1.1", "192.168.1.2"}),
							"outbound":     toTFList(t, []string{"10.0.0.1", "10.0.0.2"}),
						}),
						types.ObjectValueMust(projectipaddresses.ClusterIPsObjectType.AttrTypes, map[string]attr.Value{
							"cluster_name": types.StringValue("cluster2"),
							"inbound":      toTFList(t, []string{"192.168.2.1"}),
							"outbound":     toTFList(t, []string{"10.0.1.1"}),
						}),
					}),
				}),
			},
		},
		"Empty Services": {
			SDKResp: &admin.GroupIPAddresses{
				GroupId:  admin.PtrString(dummyProjectID),
				Services: &admin.GroupService{},
			},
			expectedTFModel: &projectipaddresses.ProjectIpAddressesModel{
				ProjectId: types.StringValue(dummyProjectID),
				Services: types.ObjectValueMust(projectipaddresses.ServicesObjectType.AttrTypes, map[string]attr.Value{
					"clusters": types.ListValueMust(types.ObjectType{AttrTypes: projectipaddresses.ClusterIPsObjectType.AttrTypes}, []attr.Value{}),
				}),
			},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			resultModel, diags := projectipaddresses.NewTFProjectIPAddresses(context.Background(), tc.SDKResp)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedTFModel, resultModel, "created terraform model did not match expected output")
		})
	}
}

func toTFList(t *testing.T, values []string) attr.Value {
	t.Helper()
	list, diags := types.ListValue(types.StringType, convertToAttrValues(values))
	if diags.HasError() {
		t.Errorf("unexpected errors found when creating test cases: %s", diags.Errors()[0].Summary())
	}
	return list
}

func convertToAttrValues(values []string) []attr.Value {
	attrValues := make([]attr.Value, len(values))
	for i, v := range values {
		attrValues[i] = types.StringValue(v)
	}
	return attrValues
}
