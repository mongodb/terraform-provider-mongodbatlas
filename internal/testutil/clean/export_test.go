package clean

// Test helpers exported only for package clean_test (see project_clean_test.go).

var (
	RemoveClusters           = removeClusters
	RemoveFlexClusters       = removeFlexClusters
	RemoveFederatedDatabases = removeFederatedDatabases
)
