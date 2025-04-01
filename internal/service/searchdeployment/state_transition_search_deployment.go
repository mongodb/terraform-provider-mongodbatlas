package searchdeployment

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"go.mongodb.org/atlas-sdk/v20250312001/admin"
)

const SearchDeploymentDoesNotExistsError = "ATLAS_SEARCH_DEPLOYMENT_DOES_NOT_EXIST"
var pendingStates = []string{retrystrategy.RetryStrategyUpdatingState, retrystrategy.RetryStrategyPausedState}

func WaitSearchNodeStateTransition(ctx context.Context, projectID, clusterName string, client admin.AtlasSearchApi,
	timeConfig retrystrategy.TimeConfig, extraTargetStates ...string) (*admin.ApiSearchDeploymentResponse, error) {
	targetStates := []string{retrystrategy.RetryStrategyIdleState}
	targetStates = append(targetStates, extraTargetStates...)
	stateConf := &retry.StateChangeConf{
		Pending:    pendingStates,
		Target:     targetStates,
		Refresh:    searchDeploymentRefreshFunc(ctx, projectID, clusterName, client),
		Timeout:    timeConfig.Timeout,
		MinTimeout: timeConfig.MinTimeout,
		Delay:      timeConfig.Delay,
	}

	result, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, err
	}
	if deploymentResp, ok := result.(*admin.ApiSearchDeploymentResponse); ok && deploymentResp != nil {
		return deploymentResp, nil
	}
	return nil, errors.New("did not obtain valid result when waiting for search deployment state transition")
}

func WaitSearchNodeDelete(ctx context.Context, projectID, clusterName string, client admin.AtlasSearchApi, timeConfig retrystrategy.TimeConfig) error {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{retrystrategy.RetryStrategyIdleState, retrystrategy.RetryStrategyUpdatingState, retrystrategy.RetryStrategyPausedState},
		Target:     []string{retrystrategy.RetryStrategyDeletedState},
		Refresh:    searchDeploymentRefreshFunc(ctx, projectID, clusterName, client),
		Timeout:    timeConfig.Timeout,
		MinTimeout: timeConfig.MinTimeout,
		Delay:      timeConfig.Delay,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func searchDeploymentRefreshFunc(ctx context.Context, projectID, clusterName string, client admin.AtlasSearchApi) retry.StateRefreshFunc {
	return func() (any, string, error) {
		deploymentResp, resp, err := client.GetAtlasSearchDeployment(ctx, projectID, clusterName).Execute()
		if err != nil && deploymentResp == nil && resp == nil {
			return nil, "", err
		}
		if err != nil {
			if validate.StatusNotFound(resp) && strings.Contains(err.Error(), SearchDeploymentDoesNotExistsError) {
				return "", retrystrategy.RetryStrategyDeletedState, nil
			}
			if validate.StatusServiceUnavailable(resp) {
				return "", retrystrategy.RetryStrategyUpdatingState, nil
			}
			return nil, "", err
		}

		if IsNotFoundDeploymentResponse(deploymentResp) {
			return "", retrystrategy.RetryStrategyDeletedState, nil
		}

		if conversion.IsStringPresent(deploymentResp.StateName) {
			tflog.Debug(ctx, fmt.Sprintf("search deployment status: %s", *deploymentResp.StateName))
			return deploymentResp, *deploymentResp.StateName, nil
		}
		return deploymentResp, "", nil
	}
}
