package flexsnapshot_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/flexsnapshot"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20241113004/admin"
)

type sdkToTFModelTestCase struct {
	ProjectID       string
	Name            string
	SDKResp         *admin.FlexBackupSnapshot20241113
	expectedTFModel *flexsnapshot.TFModel
}

func TestFlexSnapshotSDKToTFModel(t *testing.T) {
	testCases := map[string]sdkToTFModelTestCase{ // TODO: consider adding test cases to contemplate all possible API responses
		"Complete SDK response": {
			ProjectID:       "projectID",
			Name:            "name",
			SDKResp:         &admin.FlexBackupSnapshot20241113{},
			expectedTFModel: &flexsnapshot.TFModel{},
		},
		"Empty SDK response": {
			ProjectID:       "",
			Name:            "",
			SDKResp:         &admin.FlexBackupSnapshot20241113{},
			expectedTFModel: &flexsnapshot.TFModel{},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			resultModel := flexsnapshot.NewTFModel(tc.SDKResp)
			assert.Equal(t, tc.expectedTFModel, resultModel, "created terraform model did not match expected output")
		})
	}
}
