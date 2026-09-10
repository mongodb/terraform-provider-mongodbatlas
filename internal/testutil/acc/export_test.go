package acc

// Test helpers exported only for package acc_test (see atlas_test.go).

type ClusterCreateFuncs = clusterCreateFuncs

func NewClusterCreateFuncs(create func(region string) error, wait, deleteFn func() error) ClusterCreateFuncs {
	return clusterCreateFuncs{create: create, wait: wait, delete: deleteFn}
}

var (
	CreateClusterWithRegionFallback = createClusterWithRegionFallback
	IsFailedProvisioning            = isFailedProvisioning
)
