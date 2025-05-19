package flexsnapshot_test

import (
	"testing"
	"time"

	"go.mongodb.org/atlas-sdk/v20250312003/admin"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/flexsnapshot"
)

type sdkToTFModelTestCase struct {
	SDKResp         *admin.FlexBackupSnapshot20241113
	expectedTFModel *flexsnapshot.TFModel
	ProjectID       string
	Name            string
}

func TestFlexSnapshotSDKToTFModel(t *testing.T) {
	var (
		projectID      = "projectID"
		name           = "name"
		id             = "id"
		MongoDBVersion = "MongoDBVersion"
		status         = "status"
		now            = time.Now()
	)

	testCases := map[string]sdkToTFModelTestCase{
		"Complete SDK response": {
			ProjectID: projectID,
			Name:      name,
			SDKResp: &admin.FlexBackupSnapshot20241113{
				Expiration:     &now,
				FinishTime:     &now,
				Id:             &id,
				MongoDBVersion: &MongoDBVersion,
				ScheduledTime:  &now,
				StartTime:      &now,
				Status:         &status,
			},
			expectedTFModel: &flexsnapshot.TFModel{
				ProjectId:      types.StringValue(projectID),
				Name:           types.StringValue(name),
				Expiration:     types.StringPointerValue(conversion.TimePtrToStringPtr(&now)),
				FinishTime:     types.StringPointerValue(conversion.TimePtrToStringPtr(&now)),
				SnapshotId:     types.StringValue(id),
				MongoDBVersion: types.StringValue(MongoDBVersion),
				ScheduledTime:  types.StringPointerValue(conversion.TimePtrToStringPtr(&now)),
				StartTime:      types.StringPointerValue(conversion.TimePtrToStringPtr(&now)),
				Status:         types.StringValue(status),
			},
		},
		"Empty SDK response": {
			ProjectID: projectID,
			Name:      name,
			SDKResp:   &admin.FlexBackupSnapshot20241113{},
			expectedTFModel: &flexsnapshot.TFModel{
				ProjectId: types.StringValue(projectID),
				Name:      types.StringValue(name),
			},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			resultModel := flexsnapshot.NewTFModel(tc.ProjectID, tc.Name, tc.SDKResp)
			assert.Equal(t, tc.expectedTFModel, resultModel, "created terraform model did not match expected output")
		})
	}
}
