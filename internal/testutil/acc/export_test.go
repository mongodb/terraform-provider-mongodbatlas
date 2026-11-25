package acc

import "time"

// Test helpers exported only for package acc_test (see atlas_test.go and cloud_backup_collection_restore_fixture_test.go).

type (
	ClusterCreateFuncs      = clusterCreateFuncs
	SampleDatasetLoadTiming = sampleDatasetLoadTiming
	SeededClusterFuncs      = seededClusterFuncs
)

func NewClusterCreateFuncs(create func(region string) error, wait, deleteFn func() error) ClusterCreateFuncs {
	return clusterCreateFuncs{create: create, wait: wait, delete: deleteFn}
}

func NewSampleDatasetLoadTiming(totalTimeout, jobTimeout, pollInterval, initialDelay, retryDelay time.Duration, maxAttempts int) SampleDatasetLoadTiming {
	return sampleDatasetLoadTiming{
		totalTimeout: totalTimeout,
		jobTimeout:   jobTimeout,
		pollInterval: pollInterval,
		initialDelay: initialDelay,
		retryDelay:   retryDelay,
		maxAttempts:  maxAttempts,
	}
}

func NewSeededClusterFuncs(create func() (string, error), populate, deleteFn func(clusterName string) error) SeededClusterFuncs {
	return seededClusterFuncs{
		create:   create,
		populate: populate,
		delete:   deleteFn,
	}
}

var (
	CreateClusterWithRegionFallback   = createClusterWithRegionFallback
	IsFailedProvisioning              = isFailedProvisioning
	DeleteProjectUntil                = deleteProjectUntil
	DeleteLocalProjectSharedResources = deleteLocalProjectSharedResources
	DeleteExecutionProjects           = deleteExecutionProjects
	PopulateWithSampleDataWithTiming  = populateWithSampleData
	IsTaintedSampleDatasetClusterErr  = isTaintedSampleDatasetClusterError
	CreateSeededClusterWithRetry      = createSeededClusterWithRetry
	RequireSnapshotSeedCollections    = requireSnapshotSeedCollections
)
