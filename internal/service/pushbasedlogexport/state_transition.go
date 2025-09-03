package pushbasedlogexport

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20250312007/admin"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
)

const (
	ActiveState                   = "ACTIVE"
	UnconfiguredState             = "UNCONFIGURED"
	InitiatingState               = "INITIATING"
	BucketVerifiedState           = "BUCKET_VERIFIED"
	BucketVerificationFailedState = "BUCKET_VERIFICATION_FAILED"
	AssumeRoleFailedState         = "ASSUME_ROLE_FAILED"
)

var failureStates = []string{BucketVerificationFailedState, AssumeRoleFailedState}

func WaitStateTransition(ctx context.Context, projectID string, client admin.PushBasedLogExportApi,
	timeConfig retrystrategy.TimeConfig) (*admin.PushBasedLogExportProject, error) {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{InitiatingState, BucketVerifiedState},
		Target:     []string{ActiveState},
		Refresh:    refreshFunc(ctx, projectID, client),
		Timeout:    timeConfig.Timeout,
		MinTimeout: timeConfig.MinTimeout,
		Delay:      timeConfig.Delay,
	}

	result, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, err
	}
	if logConfig, ok := result.(*admin.PushBasedLogExportProject); ok && logConfig != nil {
		return logConfig, nil
	}
	return nil, errors.New("did not obtain valid result when waiting for push-based log export configuration state transition")
}

func WaitResourceDelete(ctx context.Context, projectID string, client admin.PushBasedLogExportApi, timeConfig retrystrategy.TimeConfig) error {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{ActiveState, InitiatingState, BucketVerifiedState},
		Target:     []string{UnconfiguredState},
		Refresh:    refreshFunc(ctx, projectID, client),
		Timeout:    timeConfig.Timeout,
		MinTimeout: timeConfig.MinTimeout,
		Delay:      timeConfig.Delay,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func refreshFunc(ctx context.Context, projectID string, client admin.PushBasedLogExportApi) retry.StateRefreshFunc {
	return func() (any, string, error) {
		logConfig, resp, err := client.GetLogExport(ctx, projectID).Execute()
		if err != nil && logConfig == nil && resp == nil {
			return nil, "", err
		}
		if err != nil {
			if validate.StatusNotFound(resp) {
				return "", retrystrategy.RetryStrategyDeletedState, nil
			}
			return nil, "", err
		}

		if conversion.IsStringPresent(logConfig.State) {
			tflog.Debug(ctx, fmt.Sprintf("push-based log export configuration status: %s", *logConfig.State))
			return logConfig, *logConfig.State, nil
		}
		return logConfig, "", nil
	}
}
