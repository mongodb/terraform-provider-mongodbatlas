package clusteroutagesimulation_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312027/admin"
	"go.mongodb.org/atlas-sdk/v20250312027/mockadmin"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/clusteroutagesimulation"
)

const (
	testProjectID   = "projectID"
	testClusterName = "clusterName"
	testMinTimeout  = 10 * time.Millisecond
	testWaitTimeout = 5 * time.Second
)

func simulationInState(state string) *admin.ClusterOutageSimulation {
	return &admin.ClusterOutageSimulation{State: conversion.StringPtr(state)}
}

// TestWaitForDeletableStatePollsUntilDeletable verifies that the wait polls the simulation
// every min timeout and returns as soon as it reaches a deletable state.
func TestWaitForDeletableStatePollsUntilDeletable(t *testing.T) {
	m := mockadmin.NewClusterOutageSimulationAPI(t)
	m.EXPECT().GetOutageSimulation(mock.Anything, testProjectID, testClusterName).Return(admin.GetOutageSimulationApiRequest{ApiService: m})
	m.EXPECT().GetOutageSimulationExecute(mock.Anything).Return(simulationInState("STARTING"), nil, nil).Twice()
	m.EXPECT().GetOutageSimulationExecute(mock.Anything).Return(simulationInState("SIMULATING"), nil, nil).Once()

	result, err := clusteroutagesimulation.WaitForDeletableStateWithMinTimeout(t.Context(), m, testProjectID, testClusterName, testWaitTimeout, testMinTimeout)

	require.NoError(t, err)
	assert.Equal(t, "SIMULATING", result.GetState())
}

// TestWaitForDeletableStateSimulationAlreadyGone verifies that a 404 is reported as a deleted
// simulation instead of failing or panicking.
func TestWaitForDeletableStateSimulationAlreadyGone(t *testing.T) {
	m := mockadmin.NewClusterOutageSimulationAPI(t)
	m.EXPECT().GetOutageSimulation(mock.Anything, testProjectID, testClusterName).Return(admin.GetOutageSimulationApiRequest{ApiService: m}).Once()
	m.EXPECT().GetOutageSimulationExecute(mock.Anything).Return(nil, &http.Response{StatusCode: http.StatusNotFound}, errors.New("not found")).Once()

	result, err := clusteroutagesimulation.WaitForDeletableStateWithMinTimeout(t.Context(), m, testProjectID, testClusterName, testWaitTimeout, testMinTimeout)

	require.NoError(t, err)
	assert.Equal(t, "DELETED", result.GetState())
}

// TestEndOutageSimulationAndWaitPollsUntilDeleted verifies that ending the simulation polls it
// every min timeout and returns as soon as it is deleted.
func TestEndOutageSimulationAndWaitPollsUntilDeleted(t *testing.T) {
	m := mockadmin.NewClusterOutageSimulationAPI(t)
	m.EXPECT().EndOutageSimulation(mock.Anything, testProjectID, testClusterName).Return(admin.EndOutageSimulationApiRequest{ApiService: m}).Once()
	m.EXPECT().EndOutageSimulationExecute(mock.Anything).Return(simulationInState("RECOVERY_REQUESTED"), nil, nil).Once()
	m.EXPECT().GetOutageSimulation(mock.Anything, testProjectID, testClusterName).Return(admin.GetOutageSimulationApiRequest{ApiService: m})
	m.EXPECT().GetOutageSimulationExecute(mock.Anything).Return(simulationInState("RECOVERING"), nil, nil).Once()
	m.EXPECT().GetOutageSimulationExecute(mock.Anything).Return(simulationInState("DELETED"), nil, nil).Once()

	err := clusteroutagesimulation.EndOutageSimulationAndWaitWithMinTimeout(t.Context(), m, testProjectID, testClusterName, testWaitTimeout, testMinTimeout)

	require.NoError(t, err)
}
