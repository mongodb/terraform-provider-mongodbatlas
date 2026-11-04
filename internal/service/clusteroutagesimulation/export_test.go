package clusteroutagesimulation

// Test helpers exported only for package clusteroutagesimulation_test (see state_transition_test.go).

var (
	WaitForDeletableState      = waitForDeletableState
	EndOutageSimulationAndWait = endOutageSimulationAndWait
)
