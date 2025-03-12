package projectipaddresses_test

import (
	"testing"

	"go.mongodb.org/atlas-sdk/v20250219001/admin"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/projectipaddresses"
)

const (
	dummyProjectID = "111111111111111111111111"
)

type sdkToTFModelTestCase struct {
	SDKResp         *admin.GroupIPAddresses
	expectedTFModel *projectipaddresses.TFProjectIpAddressesModel
}

func TestProjectIPAddressesSDKToTFModel(t *testing.T) {
	testCases := map[string]sdkToTFModelTestCase{
		"Complete SDK response": {
			SDKResp: &admin.GroupIPAddresses{
				GroupId: admin.PtrString(dummyProjectID),
				Services: &admin.GroupService{
					Clusters: &[]admin.ClusterIPAddresses{
						{
							ClusterName:    admin.PtrString("cluster1"),
							Inbound:        &[]string{"192.168.1.1", "192.168.1.2"},
							Outbound:       &[]string{"10.0.0.1", "10.0.0.2"},
							FutureInbound:  &[]string{"192.168.1.1", "192.168.1.2"},
							FutureOutbound: &[]string{"10.0.0.1", "10.0.0.2"},
						},
						{
							ClusterName:    admin.PtrString("cluster2"),
							Inbound:        &[]string{"192.168.2.1"},
							Outbound:       &[]string{"10.0.1.1"},
							FutureInbound:  &[]string{"192.168.2.1"},
							FutureOutbound: &[]string{"10.0.1.1"},
						},
					},
				},
			},
			expectedTFModel: &projectipaddresses.TFProjectIpAddressesModel{
				ProjectId: types.StringValue(dummyProjectID),
				Services: createExpectedServices(t, []projectipaddresses.TFClusterValueModel{
					{
						ClusterName: types.StringValue("cluster1"),
						Inbound: types.ListValueMust(types.StringType, []attr.Value{
							types.StringValue("192.168.1.1"),
							types.StringValue("192.168.1.2"),
						}),
						Outbound: types.ListValueMust(types.StringType, []attr.Value{
							types.StringValue("10.0.0.1"),
							types.StringValue("10.0.0.2"),
						}),
						FutureInbound: types.ListValueMust(types.StringType, []attr.Value{
							types.StringValue("192.168.1.1"),
							types.StringValue("192.168.1.2"),
						}),
						FutureOutbound: types.ListValueMust(types.StringType, []attr.Value{
							types.StringValue("10.0.0.1"),
							types.StringValue("10.0.0.2"),
						}),
					},
					{
						ClusterName: types.StringValue("cluster2"),
						Inbound: types.ListValueMust(types.StringType, []attr.Value{
							types.StringValue("192.168.2.1"),
						}),
						Outbound: types.ListValueMust(types.StringType, []attr.Value{
							types.StringValue("10.0.1.1"),
						}),
						FutureInbound: types.ListValueMust(types.StringType, []attr.Value{
							types.StringValue("192.168.2.1"),
						}),
						FutureOutbound: types.ListValueMust(types.StringType, []attr.Value{
							types.StringValue("10.0.1.1"),
						}),
					},
				}),
			},
		},
		"Single Cluster with no IP addresses": { // case when cluster is being created
			SDKResp: &admin.GroupIPAddresses{
				GroupId: admin.PtrString(dummyProjectID),
				Services: &admin.GroupService{
					Clusters: &[]admin.ClusterIPAddresses{
						{
							ClusterName:    admin.PtrString("cluster1"),
							Inbound:        &[]string{},
							Outbound:       &[]string{},
							FutureInbound:  &[]string{},
							FutureOutbound: &[]string{},
						},
					},
				},
			},
			expectedTFModel: &projectipaddresses.TFProjectIpAddressesModel{
				ProjectId: types.StringValue(dummyProjectID),
				Services: createExpectedServices(t, []projectipaddresses.TFClusterValueModel{
					{
						ClusterName:    types.StringValue("cluster1"),
						Inbound:        types.ListValueMust(types.StringType, []attr.Value{}),
						Outbound:       types.ListValueMust(types.StringType, []attr.Value{}),
						FutureInbound:  types.ListValueMust(types.StringType, []attr.Value{}),
						FutureOutbound: types.ListValueMust(types.StringType, []attr.Value{}),
					},
				}),
			},
		},
		"Empty Services": {
			SDKResp: &admin.GroupIPAddresses{
				GroupId:  admin.PtrString(dummyProjectID),
				Services: &admin.GroupService{},
			},
			expectedTFModel: &projectipaddresses.TFProjectIpAddressesModel{
				ProjectId: types.StringValue(dummyProjectID),
				Services:  createExpectedServices(t, []projectipaddresses.TFClusterValueModel{}),
			},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			resultModel, diags := projectipaddresses.NewTFProjectIPAddresses(t.Context(), tc.SDKResp)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedTFModel, resultModel, "created terraform model did not match expected output")
		})
	}
}

func createExpectedServices(t *testing.T, clusters []projectipaddresses.TFClusterValueModel) types.Object {
	t.Helper()
	servicesValue := projectipaddresses.TFServicesModel{
		Clusters: clusters,
	}

	servicesObj, diags := types.ObjectValueFrom(t.Context(), projectipaddresses.ServicesObjectType.AttrTypes, servicesValue)
	if diags.HasError() {
		t.Fatalf("unexpected errors found when creating services object: %s", diags.Errors()[0].Summary())
	}

	return servicesObj
}
