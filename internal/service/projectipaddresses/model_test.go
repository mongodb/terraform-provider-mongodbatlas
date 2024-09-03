package projectipaddresses_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/projectipaddresses"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20240805003/admin"
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
			expectedTFModel: &projectipaddresses.TFProjectIpAddressesModel{
				ProjectId: types.StringValue(dummyProjectID),
				Services: createExpectedServices(t, []projectipaddresses.ClustersValue{
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
					},
					{
						ClusterName: types.StringValue("cluster2"),
						Inbound: types.ListValueMust(types.StringType, []attr.Value{
							types.StringValue("192.168.2.1"),
						}),
						Outbound: types.ListValueMust(types.StringType, []attr.Value{
							types.StringValue("10.0.1.1"),
						}),
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
				Services:  createExpectedServices(t, []projectipaddresses.ClustersValue{}),
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

func createExpectedServices(t *testing.T, clusters []projectipaddresses.ClustersValue) types.Object {
	t.Helper()
	servicesValue := projectipaddresses.ServicesValue{
		Clusters: clusters,
	}

	servicesObj, diags := types.ObjectValueFrom(context.Background(), projectipaddresses.ServicesObjectType.AttrTypes, servicesValue)
	if diags.HasError() {
		t.Fatalf("unexpected errors found when creating services object: %s", diags.Errors()[0].Summary())
	}

	return servicesObj
}
