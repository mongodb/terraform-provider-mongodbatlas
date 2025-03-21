package pushbasedlogexport_test

import (
	"testing"
	"time"

	"go.mongodb.org/atlas-sdk/v20250312001/admin"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/pushbasedlogexport"
)

var (
	testBucketName  = "test-bucket-name"
	testIAMRoleID   = "661fe3ad234b02027ddee196"
	testPrefixPath  = "prefix/path"
	prefixPathEmpty = ""
	testProjectID   = "661fe3ad234b02027dabcabc"
)

type sdkToTFModelTestCase struct {
	apiResp         *admin.PushBasedLogExportProject
	timeout         *timeouts.Value
	expectedTFModel *pushbasedlogexport.TFPushBasedLogExportRSModel
	name            string
	projectID       string
}

func TestNewTFPushBasedLogExport(t *testing.T) {
	currentTime := time.Now()

	testCases := []sdkToTFModelTestCase{
		{
			name:      "Complete API response",
			projectID: testProjectID,
			apiResp: &admin.PushBasedLogExportProject{
				BucketName: admin.PtrString(testBucketName),
				CreateDate: admin.PtrTime(currentTime),
				IamRoleId:  admin.PtrString(testIAMRoleID),
				PrefixPath: admin.PtrString(testPrefixPath),
				State:      admin.PtrString(activeState),
			},
			expectedTFModel: &pushbasedlogexport.TFPushBasedLogExportRSModel{
				ProjectID:  types.StringValue(testProjectID),
				BucketName: types.StringValue(testBucketName),
				IamRoleID:  types.StringValue(testIAMRoleID),
				PrefixPath: types.StringValue(testPrefixPath),
				State:      types.StringValue(activeState),
				CreateDate: types.StringPointerValue(conversion.TimePtrToStringPtr(&currentTime)),
			},
		},
		{
			name:      "Complete API response with empty prefix path",
			projectID: testProjectID,
			apiResp: &admin.PushBasedLogExportProject{
				BucketName: admin.PtrString(testBucketName),
				CreateDate: admin.PtrTime(currentTime),
				IamRoleId:  admin.PtrString(testIAMRoleID),
				PrefixPath: admin.PtrString(prefixPathEmpty),
				State:      admin.PtrString(activeState),
			},
			expectedTFModel: &pushbasedlogexport.TFPushBasedLogExportRSModel{
				ProjectID:  types.StringValue(testProjectID),
				BucketName: types.StringValue(testBucketName),
				IamRoleID:  types.StringValue(testIAMRoleID),
				PrefixPath: types.StringValue(prefixPathEmpty),
				State:      types.StringValue(activeState),
				CreateDate: types.StringPointerValue(conversion.TimePtrToStringPtr(&currentTime)),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel, _ := pushbasedlogexport.NewTFPushBasedLogExport(t.Context(), tc.projectID, tc.apiResp, tc.timeout)
			if !assert.Equal(t, tc.expectedTFModel, resultModel) {
				t.Errorf("result model does not match expected output: expected %+v, got %+v", tc.expectedTFModel, resultModel)
			}
		})
	}
}

type pushBasedLogExportReqTestCase struct {
	input             *pushbasedlogexport.TFPushBasedLogExportRSModel
	expectedCreateReq *admin.CreatePushBasedLogExportProjectRequest
	expectedUpdateReq *admin.PushBasedLogExportProject
	name              string
}

func TestNewPushBasedLogExportReq(t *testing.T) {
	testCases := []pushBasedLogExportReqTestCase{
		{
			name: "Valid TF state",
			input: &pushbasedlogexport.TFPushBasedLogExportRSModel{
				BucketName: types.StringValue(testBucketName),
				IamRoleID:  types.StringValue(testIAMRoleID),
				PrefixPath: types.StringValue(testPrefixPath),
			},
			expectedCreateReq: &admin.CreatePushBasedLogExportProjectRequest{
				BucketName: testBucketName,
				IamRoleId:  testIAMRoleID,
				PrefixPath: testPrefixPath,
			},
			expectedUpdateReq: &admin.PushBasedLogExportProject{
				BucketName: admin.PtrString(testBucketName),
				IamRoleId:  admin.PtrString(testIAMRoleID),
				PrefixPath: admin.PtrString(testPrefixPath),
			},
		},
		{
			name: "Valid TF state with empty prefix path",
			input: &pushbasedlogexport.TFPushBasedLogExportRSModel{
				BucketName: types.StringValue(testBucketName),
				IamRoleID:  types.StringValue(testIAMRoleID),
				PrefixPath: types.StringValue(prefixPathEmpty),
			},
			expectedCreateReq: &admin.CreatePushBasedLogExportProjectRequest{
				BucketName: testBucketName,
				IamRoleId:  testIAMRoleID,
				PrefixPath: prefixPathEmpty,
			},
			expectedUpdateReq: &admin.PushBasedLogExportProject{
				BucketName: admin.PtrString(testBucketName),
				IamRoleId:  admin.PtrString(testIAMRoleID),
				PrefixPath: admin.PtrString(prefixPathEmpty),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name+" Create", func(t *testing.T) {
			createReq := pushbasedlogexport.NewPushBasedLogExportCreateReq(tc.input)
			if !assert.Equal(t, tc.expectedCreateReq, createReq) {
				t.Errorf("Create request does not match expected output: expected %+v, got %+v", tc.expectedCreateReq, createReq)
			}
		})
		t.Run(tc.name+" Update", func(t *testing.T) {
			updateReq := pushbasedlogexport.NewPushBasedLogExportUpdateReq(tc.input)
			if !assert.Equal(t, tc.expectedUpdateReq, updateReq) {
				t.Errorf("Update request does not match expected output: expected %+v, got %+v", tc.expectedUpdateReq, updateReq)
			}
		})
	}
}
