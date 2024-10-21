package flexcluster

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"go.mongodb.org/atlas-sdk/v20240805004/admin"
)

const (
	IdleState      = "IDLE"
	CreatingState  = "CREATING"
	UpdatingState  = "UPDATING"
	DeletingState  = "DELETING"
	RepairingState = "REPAIRING"
)

func WaitStateTransition(ctx context.Context, requestParams *admin.GetFlexClusterApiParams, client admin.FlexClustersApi, pendingStates, desiredStates []string) (*admin.FlexClusterDescription20250101, error) {
	stateConf := &retry.StateChangeConf{
		Pending:    pendingStates,
		Target:     desiredStates,
		Refresh:    refreshFunc(ctx, requestParams, client),
		Timeout:    5 * time.Minute,
		MinTimeout: 3 * time.Second,
		Delay:      0,
	}

	flexClusterResp, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, err
	}

	if flexCluster, ok := flexClusterResp.(*admin.FlexClusterDescription20250101); ok && flexCluster != nil {
		return flexCluster, nil
	}

	return nil, errors.New("did not obtain valid result when waiting for flex cluster state transition")
}

func refreshFunc(ctx context.Context, requestParams *admin.GetFlexClusterApiParams, client admin.FlexClustersApi) retry.StateRefreshFunc {
	return func() (any, string, error) {
		flexCluster, _, err := client.GetFlexClusterWithParams(ctx, requestParams).Execute()
		if err != nil {
			return nil, "", err
		}
		state := flexCluster.GetStateName()
		return flexCluster, state, nil
	}
}
