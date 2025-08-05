package streamworkspace_test

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	streamworkspace "github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streamworkspace"
	"go.mongodb.org/atlas-sdk/v20250312005/admin"
)

const (
	dummyProjectID         = "111111111111111111111111"
	dummyStreamWorkspaceID = "222222222222222222222222"
	cloudProvider          = "AWS"
	region                 = "VIRGINIA_USA"
	workspaceName          = "WorkspaceName"
	tier                   = "SP30"
)

var hostnames = &[]string{"atlas-stream.virginia-usa.a.query.mongodb-dev.net"}

type sdkToTFModelTestCase struct {
	SDKResp         *admin.StreamsTenant
	expectedTFModel *streamworkspace.TFStreamWorkspaceModel
	name            string
}

func TestStreamWorkspaceSDKToTFModel(t *testing.T) {
	testCases := []sdkToTFModelTestCase{
		{
			name: "Complete SDK response",
			SDKResp: &admin.StreamsTenant{
				Id: admin.PtrString(dummyStreamWorkspaceID),
				DataProcessRegion: &admin.StreamsDataProcessRegion{
					CloudProvider: cloudProvider,
					Region:        region,
				},
				StreamConfig: &admin.StreamConfig{
					Tier: admin.PtrString(tier),
				},
				GroupId:   admin.PtrString(dummyProjectID),
				Hostnames: hostnames,
				Name:      admin.PtrString(workspaceName),
			},
			expectedTFModel: &streamworkspace.TFStreamWorkspaceModel{
				ID:                types.StringValue(dummyStreamWorkspaceID),
				DataProcessRegion: tfRegionObject(t, cloudProvider, region),
				ProjectID:         types.StringValue(dummyProjectID),
				Hostnames:         tfHostnamesList(t, hostnames),
				WorkspaceName:     types.StringValue(workspaceName),
				StreamConfig:      tfStreamConfigObject(t, tier),
			},
		},
		{
			name: "Empty hostnames, streamConfig and dataProcessRegion in response", // should never happen, but verifying it is handled gracefully
			SDKResp: &admin.StreamsTenant{
				Id:      admin.PtrString(dummyStreamWorkspaceID),
				GroupId: admin.PtrString(dummyProjectID),
				Name:    admin.PtrString(workspaceName),
			},
			expectedTFModel: &streamworkspace.TFStreamWorkspaceModel{
				ID:                types.StringValue(dummyStreamWorkspaceID),
				DataProcessRegion: types.ObjectNull(streamworkspace.ProcessRegionObjectType.AttrTypes),
				ProjectID:         types.StringValue(dummyProjectID),
				Hostnames:         types.ListNull(types.StringType),
				WorkspaceName:     types.StringValue(workspaceName),
				StreamConfig:      types.ObjectNull(streamworkspace.StreamConfigObjectType.AttrTypes),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel, diags := streamworkspace.NewTFStreamWorkspace(t.Context(), tc.SDKResp)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			if !reflect.DeepEqual(resultModel, tc.expectedTFModel) {
				t.Errorf("created terraform model did not match expected output")
			}
		})
	}
}

type paginatedWorkspacesSDKToTFModelTestCase struct {
	SDKResp         *admin.PaginatedApiStreamsTenant
	providedConfig  *streamworkspace.TFStreamWorkspacesModel
	expectedTFModel *streamworkspace.TFStreamWorkspacesModel
	name            string
}

func TestStreamWorkspacesSDKToTFModel(t *testing.T) {
	testCases := []paginatedWorkspacesSDKToTFModelTestCase{
		{
			name: "Complete SDK response with configured page options",
			SDKResp: &admin.PaginatedApiStreamsTenant{
				Results: &[]admin.StreamsTenant{
					{
						Id: admin.PtrString(dummyStreamWorkspaceID),
						DataProcessRegion: &admin.StreamsDataProcessRegion{
							CloudProvider: cloudProvider,
							Region:        region,
						},
						GroupId:   admin.PtrString(dummyProjectID),
						Hostnames: hostnames,
						Name:      admin.PtrString(workspaceName),
						StreamConfig: &admin.StreamConfig{
							Tier: admin.PtrString(tier),
						},
					},
				},
				TotalCount: admin.PtrInt(1),
			},
			providedConfig: &streamworkspace.TFStreamWorkspacesModel{
				ProjectID:    types.StringValue(dummyProjectID),
				PageNum:      types.Int64Value(1),
				ItemsPerPage: types.Int64Value(2),
			},
			expectedTFModel: &streamworkspace.TFStreamWorkspacesModel{
				ProjectID:    types.StringValue(dummyProjectID),
				PageNum:      types.Int64Value(1),
				ItemsPerPage: types.Int64Value(2),
				TotalCount:   types.Int64Value(1),
				Results: []streamworkspace.TFStreamWorkspaceModel{
					{
						ID:                types.StringValue(dummyStreamWorkspaceID),
						DataProcessRegion: tfRegionObject(t, cloudProvider, region),
						ProjectID:         types.StringValue(dummyProjectID),
						Hostnames:         tfHostnamesList(t, hostnames),
						WorkspaceName:     types.StringValue(workspaceName),
						StreamConfig:      tfStreamConfigObject(t, tier),
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
			providedConfig: &streamworkspace.TFStreamWorkspacesModel{
				ProjectID: types.StringValue(dummyProjectID),
			},
			expectedTFModel: &streamworkspace.TFStreamWorkspacesModel{
				ProjectID:    types.StringValue(dummyProjectID),
				PageNum:      types.Int64Null(),
				ItemsPerPage: types.Int64Null(),
				TotalCount:   types.Int64Value(0),
				Results:      []streamworkspace.TFStreamWorkspaceModel{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel, diags := streamworkspace.NewTFStreamWorkspaces(t.Context(), tc.providedConfig, tc.SDKResp)
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
	tfModel        *streamworkspace.TFStreamWorkspaceModel
	expectedSDKReq *admin.StreamsTenant
	name           string
}

func TestStreamWorkspaceTFToSDKCreateModel(t *testing.T) {
	testCases := []tfToSDKCreateModelTestCase{
		{
			name: "Complete TF state",
			tfModel: &streamworkspace.TFStreamWorkspaceModel{
				DataProcessRegion: tfRegionObject(t, cloudProvider, region),
				ProjectID:         types.StringValue(dummyProjectID),
				WorkspaceName:     types.StringValue(workspaceName),
				StreamConfig:      tfStreamConfigObject(t, tier),
			},
			expectedSDKReq: &admin.StreamsTenant{
				DataProcessRegion: &admin.StreamsDataProcessRegion{
					CloudProvider: cloudProvider,
					Region:        region,
				},
				GroupId: admin.PtrString(dummyProjectID),
				Name:    admin.PtrString(workspaceName),
				StreamConfig: &admin.StreamConfig{
					Tier: admin.PtrString(tier),
				},
			},
		},
		{
			name: "TF State without StreamConfig",
			tfModel: &streamworkspace.TFStreamWorkspaceModel{
				DataProcessRegion: tfRegionObject(t, cloudProvider, region),
				ProjectID:         types.StringValue(dummyProjectID),
				WorkspaceName:     types.StringValue(workspaceName),
			},
			expectedSDKReq: &admin.StreamsTenant{
				DataProcessRegion: &admin.StreamsDataProcessRegion{
					CloudProvider: cloudProvider,
					Region:        region,
				},
				GroupId: admin.PtrString(dummyProjectID),
				Name:    admin.PtrString(workspaceName),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			apiReqResult, diags := streamworkspace.NewStreamWorkspaceCreateReq(t.Context(), tc.tfModel)
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
	tfModel        *streamworkspace.TFStreamWorkspaceModel
	expectedSDKReq *admin.StreamsDataProcessRegion
	name           string
}

func TestStreamWorkspaceTFToSDKUpdateModel(t *testing.T) {
	testCases := []tfToSDKUpdateModelTestCase{
		{
			name: "Complete TF state",
			tfModel: &streamworkspace.TFStreamWorkspaceModel{
				ID:                types.StringValue(dummyStreamWorkspaceID),
				DataProcessRegion: tfRegionObject(t, cloudProvider, region),
				ProjectID:         types.StringValue(dummyProjectID),
				Hostnames:         tfHostnamesList(t, hostnames),
				WorkspaceName:     types.StringValue(workspaceName),
			},
			expectedSDKReq: &admin.StreamsDataProcessRegion{
				CloudProvider: cloudProvider,
				Region:        region,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			apiReqResult, diags := streamworkspace.NewStreamWorkspaceUpdateReq(t.Context(), tc.tfModel)
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
	dataProcessRegion, diags := types.ObjectValueFrom(t.Context(), streamworkspace.ProcessRegionObjectType.AttrTypes, streamworkspace.TFWorkspaceProcessRegionSpecModel{
		CloudProvider: types.StringValue(cloudProvider),
		Region:        types.StringValue(region),
	})
	if diags.HasError() {
		t.Errorf("failed to create terraform data process region model: %s", diags.Errors()[0].Summary())
	}
	return dataProcessRegion
}

func tfStreamConfigObject(t *testing.T, tier string) types.Object {
	t.Helper()
	streamConfig, diags := types.ObjectValueFrom(t.Context(), streamworkspace.StreamConfigObjectType.AttrTypes, streamworkspace.TFWorkspaceStreamConfigSpecModel{
		Tier: types.StringValue(tier),
	})
	if diags.HasError() {
		t.Errorf("failed to create terraform data process region model: %s", diags.Errors()[0].Summary())
	}
	return streamConfig
}

func tfHostnamesList(t *testing.T, hostnames *[]string) types.List {
	t.Helper()
	resultList, diags := types.ListValueFrom(t.Context(), types.StringType, hostnames)
	if diags.HasError() {
		t.Errorf("failed to create terraform hostnames list: %s", diags.Errors()[0].Summary())
	}
	return resultList
}
