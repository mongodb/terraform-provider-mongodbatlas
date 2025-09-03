package flexrestorejob_test

import (
	"testing"
	"time"

	"go.mongodb.org/atlas-sdk/v20250312007/admin"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/flexrestorejob"
)

type sdkToTFModelTestCase struct {
	SDKResp         *admin.FlexBackupRestoreJob20241113
	expectedTFModel *flexrestorejob.TFModel
}

func TestFlexRestoreJobSDKToTFModel(t *testing.T) {
	var (
		projectID            = "projectID"
		instanceName         = "instanceName"
		deliveryType         = "deliveryType"
		id                   = "id"
		snapshotID           = "snapshotID"
		snapshotURL          = "snapshotURL"
		status               = "status"
		targetDeploymentName = "targetDeploymentName"
		targetProjectID      = "targetProjectID"
		now                  = time.Now()
	)

	testCases := map[string]sdkToTFModelTestCase{
		"Complete SDK response": {
			SDKResp: &admin.FlexBackupRestoreJob20241113{
				ProjectId:                &projectID,
				InstanceName:             &instanceName,
				DeliveryType:             &deliveryType,
				ExpirationDate:           &now,
				Id:                       &id,
				RestoreFinishedDate:      &now,
				RestoreScheduledDate:     &now,
				SnapshotFinishedDate:     &now,
				SnapshotId:               &snapshotID,
				SnapshotUrl:              &snapshotURL,
				Status:                   &status,
				TargetDeploymentItemName: &targetDeploymentName,
				TargetProjectId:          &targetProjectID,
			},
			expectedTFModel: &flexrestorejob.TFModel{
				ProjectID:                types.StringValue(projectID),
				Name:                     types.StringValue(instanceName),
				DeliveryType:             types.StringValue(deliveryType),
				ExpirationDate:           types.StringPointerValue(conversion.TimePtrToStringPtr(&now)),
				RestoreJobID:             types.StringValue(id),
				RestoreFinishedDate:      types.StringPointerValue(conversion.TimePtrToStringPtr(&now)),
				RestoreScheduledDate:     types.StringPointerValue(conversion.TimePtrToStringPtr(&now)),
				SnapshotFinishedDate:     types.StringPointerValue(conversion.TimePtrToStringPtr(&now)),
				SnapshotID:               types.StringValue(snapshotID),
				SnapshotUrl:              types.StringValue(snapshotURL),
				Status:                   types.StringValue(status),
				TargetDeploymentItemName: types.StringValue(targetDeploymentName),
				TargetProjectID:          types.StringValue(targetProjectID),
			},
		},
		"Empty SDK response": {
			SDKResp:         &admin.FlexBackupRestoreJob20241113{},
			expectedTFModel: &flexrestorejob.TFModel{},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			resultModel := flexrestorejob.NewTFModel(tc.SDKResp)
			assert.Equal(t, tc.expectedTFModel, resultModel, "created terraform model did not match expected output")
		})
	}
}
