package streaminstance_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streaminstance"
	"go.mongodb.org/atlas-sdk/v20231115007/admin"
)

const (
	dummyProjectID        = "111111111111111111111111"
	dummyStreamInstanceID = "222222222222222222222222"
	cloudProvider         = "AWS"
	region                = "VIRGINIA_USA"
	instanceName          = "InstanceName"
)

var hostnames = &[]string{"atlas-stream.virginia-usa.a.query.mongodb-dev.net"}

type sdkToTFModelTestCase struct {
	SDKResp         *admin.StreamsTenant
	expectedTFModel *streaminstance.TFStreamInstanceModel
	name            string
}

func TestStreamInstanceSDKToTFModel(t *testing.T) {
	testCases := []sdkToTFModelTestCase{
		{
			name: "Complete SDK response",
			SDKResp: &admin.StreamsTenant{
				Id: admin.PtrString(dummyStreamInstanceID),
				DataProcessRegion: &admin.StreamsDataProcessRegion{
					CloudProvider: cloudProvider,
					Region:        region,
				},
				GroupId:   admin.PtrString(dummyProjectID),
				Hostnames: hostnames,
				Name:      admin.PtrString(instanceName),
			},
			expectedTFModel: &streaminstance.TFStreamInstanceModel{
				ID:                types.StringValue(dummyStreamInstanceID),
				DataProcessRegion: tfRegionObject(t, cloudProvider, region),
				ProjectID:         types.StringValue(dummyProjectID),
				Hostnames:         tfHostnamesList(t, hostnames),
				InstanceName:      types.StringValue(instanceName),
			},
		},
		{
			name: "Empty hostnames and dataProcessRegion in response", // should never happen, but verifying it is handled gracefully
			SDKResp: &admin.StreamsTenant{
				Id:      admin.PtrString(dummyStreamInstanceID),
				GroupId: admin.PtrString(dummyProjectID),
				Name:    admin.PtrString(instanceName),
			},
			expectedTFModel: &streaminstance.TFStreamInstanceModel{
				ID:                types.StringValue(dummyStreamInstanceID),
				DataProcessRegion: types.ObjectNull(streaminstance.ProcessRegionObjectType.AttrTypes),
				ProjectID:         types.StringValue(dummyProjectID),
				Hostnames:         types.ListNull(types.StringType),
				InstanceName:      types.StringValue(instanceName),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel, diags := streaminstance.NewTFStreamInstance(context.Background(), tc.SDKResp)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			if !reflect.DeepEqual(resultModel, tc.expectedTFModel) {
				t.Errorf("created terraform model did not match expected output")
			}
		})
	}
}

type paginatedInstancesSDKToTFModelTestCase struct {
	SDKResp         *admin.PaginatedApiStreamsTenant
	providedConfig  *streaminstance.TFStreamInstancesModel
	expectedTFModel *streaminstance.TFStreamInstancesModel
	name            string
}

func TestStreamInstancesSDKToTFModel(t *testing.T) {
	testCases := []paginatedInstancesSDKToTFModelTestCase{
		{
			name: "Complete SDK response with configured page options",
			SDKResp: &admin.PaginatedApiStreamsTenant{
				Results: &[]admin.StreamsTenant{
					{
						Id: admin.PtrString(dummyStreamInstanceID),
						DataProcessRegion: &admin.StreamsDataProcessRegion{
							CloudProvider: cloudProvider,
							Region:        region,
						},
						GroupId:   admin.PtrString(dummyProjectID),
						Hostnames: hostnames,
						Name:      admin.PtrString(instanceName),
					},
				},
				TotalCount: admin.PtrInt(1),
			},
			providedConfig: &streaminstance.TFStreamInstancesModel{
				ProjectID:    types.StringValue(dummyProjectID),
				PageNum:      types.Int64Value(1),
				ItemsPerPage: types.Int64Value(2),
			},
			expectedTFModel: &streaminstance.TFStreamInstancesModel{
				ProjectID:    types.StringValue(dummyProjectID),
				PageNum:      types.Int64Value(1),
				ItemsPerPage: types.Int64Value(2),
				TotalCount:   types.Int64Value(1),
				Results: []streaminstance.TFStreamInstanceModel{
					{
						ID:                types.StringValue(dummyStreamInstanceID),
						DataProcessRegion: tfRegionObject(t, cloudProvider, region),
						ProjectID:         types.StringValue(dummyProjectID),
						Hostnames:         tfHostnamesList(t, hostnames),
						InstanceName:      types.StringValue(instanceName),
					},
				},
			},
		},
		{
			name: "Without defining page options",
			SDKResp: &admin.PaginatedApiStreamsTenant{
				Results:    &[]admin.StreamsTenant{},
				TotalCount: admin.PtrInt(0),
			},
			providedConfig: &streaminstance.TFStreamInstancesModel{
				ProjectID: types.StringValue(dummyProjectID),
			},
			expectedTFModel: &streaminstance.TFStreamInstancesModel{
				ProjectID:    types.StringValue(dummyProjectID),
				PageNum:      types.Int64Null(),
				ItemsPerPage: types.Int64Null(),
				TotalCount:   types.Int64Value(0),
				Results:      []streaminstance.TFStreamInstanceModel{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel, diags := streaminstance.NewTFStreamInstances(context.Background(), tc.providedConfig, tc.SDKResp)
			tc.expectedTFModel.ID = resultModel.ID // id is auto-generated, have no way of defining within expected model
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			if !reflect.DeepEqual(resultModel, tc.expectedTFModel) {
				t.Errorf("created terraform model did not match expected output")
			}
		})
	}
}

type tfToSDKCreateModelTestCase struct {
	tfModel        *streaminstance.TFStreamInstanceModel
	expectedSDKReq *admin.StreamsTenant
	name           string
}

func TestStreamInstanceTFToSDKCreateModel(t *testing.T) {
	testCases := []tfToSDKCreateModelTestCase{
		{
			name: "Complete TF state",
			tfModel: &streaminstance.TFStreamInstanceModel{
				DataProcessRegion: tfRegionObject(t, cloudProvider, region),
				ProjectID:         types.StringValue(dummyProjectID),
				InstanceName:      types.StringValue(instanceName),
			},
			expectedSDKReq: &admin.StreamsTenant{
				DataProcessRegion: &admin.StreamsDataProcessRegion{
					CloudProvider: cloudProvider,
					Region:        region,
				},
				GroupId: admin.PtrString(dummyProjectID),
				Name:    admin.PtrString(instanceName),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			apiReqResult, diags := streaminstance.NewStreamInstanceCreateReq(context.Background(), tc.tfModel)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			if !reflect.DeepEqual(apiReqResult, tc.expectedSDKReq) {
				t.Errorf("created sdk model did not match expected output")
			}
		})
	}
}

type tfToSDKUpdateModelTestCase struct {
	tfModel        *streaminstance.TFStreamInstanceModel
	expectedSDKReq *admin.StreamsDataProcessRegion
	name           string
}

func TestStreamInstanceTFToSDKUpdateModel(t *testing.T) {
	testCases := []tfToSDKUpdateModelTestCase{
		{
			name: "Complete TF state",
			tfModel: &streaminstance.TFStreamInstanceModel{
				ID:                types.StringValue(dummyStreamInstanceID),
				DataProcessRegion: tfRegionObject(t, cloudProvider, region),
				ProjectID:         types.StringValue(dummyProjectID),
				Hostnames:         tfHostnamesList(t, hostnames),
				InstanceName:      types.StringValue(instanceName),
			},
			expectedSDKReq: &admin.StreamsDataProcessRegion{
				CloudProvider: cloudProvider,
				Region:        region,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			apiReqResult, diags := streaminstance.NewStreamInstanceUpdateReq(context.Background(), tc.tfModel)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			if !reflect.DeepEqual(apiReqResult, tc.expectedSDKReq) {
				t.Errorf("created sdk model did not match expected output")
			}
		})
	}
}

func tfRegionObject(t *testing.T, cloudProvider, region string) types.Object {
	t.Helper()
	dataProcessRegion, diags := types.ObjectValueFrom(context.Background(), streaminstance.ProcessRegionObjectType.AttrTypes, streaminstance.TFInstanceProcessRegionSpecModel{
		CloudProvider: types.StringValue(cloudProvider),
		Region:        types.StringValue(region),
	})
	if diags.HasError() {
		t.Errorf("failed to create terraform data process region model: %s", diags.Errors()[0].Summary())
	}
	return dataProcessRegion
}

func tfHostnamesList(t *testing.T, hostnames *[]string) types.List {
	t.Helper()
	resultList, diags := types.ListValueFrom(context.Background(), types.StringType, hostnames)
	if diags.HasError() {
		t.Errorf("failed to create terraform hostnames list: %s", diags.Errors()[0].Summary())
	}
	return resultList
}
