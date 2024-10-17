package streamprocessor

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"go.mongodb.org/atlas-sdk/v20240805005/admin"
)

const (
	InitiatingState = "INIT"
	CreatingState   = "CREATING"
	CreatedState    = "CREATED"
	StartedState    = "STARTED"
	StoppedState    = "STOPPED"
	DroppedState    = "DROPPED"
	FailedState     = "FAILED"
)

func WaitStateTransition(ctx context.Context, requestParams *admin.GetStreamProcessorApiParams, client admin.StreamsApi, pendingStates, desiredStates []string) (*admin.StreamsProcessorWithStats, error) {
	stateConf := &retry.StateChangeConf{
		Pending:    pendingStates,
		Target:     desiredStates,
		Refresh:    refreshFunc(ctx, requestParams, client),
		Timeout:    5 * time.Minute, // big pipelines can take a while to stop due to checkpointing. We prefer the API to raise the error (~ 3min) than having to expose custom timeouts.
		MinTimeout: 3 * time.Second,
		Delay:      0,
	}

	streamProcessorResp, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, err
	}

	if streamProcessor, ok := streamProcessorResp.(*admin.StreamsProcessorWithStats); ok && streamProcessor != nil {
		return streamProcessor, nil
	}

	return nil, errors.New("did not obtain valid result when waiting for stream processor state transition")
}

func refreshFunc(ctx context.Context, requestParams *admin.GetStreamProcessorApiParams, client admin.StreamsApi) retry.StateRefreshFunc {
	return func() (any, string, error) {
		streamProcessor, resp, err := client.GetStreamProcessorWithParams(ctx, requestParams).Execute()
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return "", DroppedState, err
			}
			return nil, FailedState, err
		}
		state := streamProcessor.GetState()
		if state == FailedState {
			return nil, state, fmt.Errorf("error creating MongoDB Stream Processor(%s) status was: %s", requestParams.ProcessorName, state)
		}
		return streamProcessor, state, nil
	}
}
