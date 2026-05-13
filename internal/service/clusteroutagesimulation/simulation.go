package clusteroutagesimulation

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/cleanup"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"go.mongodb.org/atlas-sdk/v20250312018/admin"
)

const defaultActionTimeout = 25 * time.Minute

func ParseActionTimeout(val types.String) (time.Duration, error) {
	if val.IsNull() || val.IsUnknown() || val.ValueString() == "" {
		return defaultActionTimeout, nil
	}
	return time.ParseDuration(val.ValueString())
}

// SimulateOutage starts a cluster outage simulation and waits for it to reach SIMULATING state.
// deleteOnTimeout controls whether to end the simulation if the wait times out.
func SimulateOutage(ctx context.Context, api admin.ClusterOutageSimulationApi, projectID, clusterName string, filters []admin.AtlasClusterOutageSimulationOutageFilter, deleteOnTimeout bool, tc retrystrategy.TimeConfig) error {
	_, _, err := api.StartOutageSimulation(ctx, projectID, clusterName, &admin.ClusterOutageSimulation{
		OutageFilters: &filters,
	}).Execute()
	if err != nil {
		return fmt.Errorf(errorClusterOutageSimulationCreate, projectID, clusterName, err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"START_REQUESTED", "STARTING"},
		Target:     []string{"SIMULATING"},
		Refresh:    outageAPIRefreshFunc(ctx, clusterName, projectID, api),
		Timeout:    tc.Timeout,
		MinTimeout: tc.MinTimeout,
		Delay:      tc.Delay,
	}

	_, errWait := stateConf.WaitForStateContext(ctx)
	return cleanup.HandleCreateTimeout(deleteOnTimeout, errWait, func(ctxCleanup context.Context) error {
		return cleanupOutageSimulation(ctxCleanup, api, projectID, clusterName, tc.Timeout)
	})
}

// StopSimulation ends an active cluster outage simulation and waits for DELETED state.
func StopSimulation(ctx context.Context, api admin.ClusterOutageSimulationApi, projectID, clusterName string, tc retrystrategy.TimeConfig) error {
	_, _, err := api.EndOutageSimulation(ctx, projectID, clusterName).Execute()
	if err != nil {
		return fmt.Errorf(errorClusterOutageSimulationDelete, projectID, clusterName, err)
	}

	log.Println("[INFO] Waiting for MongoDB Cluster Outage Simulation to end")

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"RECOVERY_REQUESTED", "RECOVERING", "COMPLETE"},
		Target:     []string{"DELETED"},
		Refresh:    outageAPIRefreshFunc(ctx, clusterName, projectID, api),
		Timeout:    tc.Timeout,
		MinTimeout: tc.MinTimeout,
		Delay:      tc.Delay,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf(errorClusterOutageSimulationDelete, projectID, clusterName, err)
	}
	return nil
}

func outageAPIRefreshFunc(ctx context.Context, clusterName, projectID string, api admin.ClusterOutageSimulationApi) retry.StateRefreshFunc {
	return func() (any, string, error) {
		outageSimulation, resp, err := api.GetOutageSimulation(ctx, projectID, clusterName).Execute()
		if err != nil {
			if validate.StatusNotFound(resp) {
				return "", "DELETED", nil
			}
			return nil, "", err
		}
		state := outageSimulation.GetState()
		if outageSimulation.State != nil {
			log.Printf("[DEBUG] status for MongoDB cluster outage simulation: %s: %s", clusterName, state)
		}
		return outageSimulation, state, nil
	}
}

func cleanupOutageSimulation(ctx context.Context, api admin.ClusterOutageSimulationApi, projectID, clusterName string, waitTimeout time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{"START_REQUESTED", "STARTING"},
		Target:     []string{"SIMULATING", "FAILED", "DELETED"},
		Refresh:    outageAPIRefreshFunc(ctx, clusterName, projectID, api),
		Timeout:    waitTimeout,
		MinTimeout: timeout,
		Delay:      timeout,
	}

	result, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil // don't fail cleanup if we can't reach a deletable state
	}

	simulation, ok := result.(*admin.ClusterOutageSimulation)
	if !ok || simulation == nil || simulation.GetState() != "SIMULATING" {
		return nil
	}

	return StopSimulation(ctx, api, projectID, clusterName, retrystrategy.TimeConfig{
		Timeout:    waitTimeout,
		MinTimeout: timeout,
		Delay:      timeout,
	})
}
