package acc

import (
	"fmt"
	"os"
	"testing"

	"go.mongodb.org/atlas-sdk/v20250312007/admin"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

// ClusterRequest contains configuration for a cluster where all fields are optional and AddDefaults is used for required fields.
// Used together with GetClusterInfo which will set ProjectID if it is unset.
type ClusterRequest struct {
	Tags                   map[string]string
	ProjectID              string
	ResourceSuffix         string
	AdvancedConfiguration  map[string]any
	ResourceDependencyName string
	ClusterName            string
	MongoDBMajorVersion    string
	ReplicationSpecs       []ReplicationSpecRequest
	DiskSizeGb             int
	CloudBackup            bool
	Geosharded             bool
	RetainBackupsEnabled   bool
	PitEnabled             bool
}

// AddDefaults ensures the required fields are populated to generate a resource.
func (r *ClusterRequest) AddDefaults() {
	if r.ResourceSuffix == "" {
		r.ResourceSuffix = defaultClusterResourceSuffix
	}
	if len(r.ReplicationSpecs) == 0 {
		r.ReplicationSpecs = []ReplicationSpecRequest{{}}
	}
	if r.ClusterName == "" {
		r.ClusterName = RandomClusterName()
	}
}

func (r *ClusterRequest) ClusterType() string {
	if r.Geosharded {
		return "GEOSHARDED"
	}
	return "REPLICASET"
}

type ClusterInfo struct {
	ProjectID        string
	Name             string
	ResourceName     string
	TerraformNameRef string
	TerraformStr     string
}

const defaultClusterResourceSuffix = "cluster_info"

// GetClusterInfo is used to obtain a project and cluster configuration resource.
// When `MONGODB_ATLAS_CLUSTER_NAME` and `MONGODB_ATLAS_PROJECT_ID` are defined, a data source is created instead. This is useful for local execution but not intended for CI executions.
// Clusters will be created in project ProjectIDExecution or in req.ProjectID which can be both a direct id, e.g., `664610ec80cc36255e634074` or a config reference `mongodbatlas_project.test.id`.
func GetClusterInfo(tb testing.TB, req *ClusterRequest) ClusterInfo {
	tb.Helper()
	if req == nil {
		req = new(ClusterRequest)
	}
	hclCreator := ClusterResourceHcl
	if req.ProjectID == "" {
		if ExistingClusterUsed() {
			projectID, clusterName := existingProjectIDClusterName()
			req.ProjectID = projectID
			req.ClusterName = clusterName
			hclCreator = ClusterDatasourceHcl
		} else {
			req.ProjectID = ProjectIDExecution(tb)
		}
	}
	clusterTerraformStr, clusterName, clusterResourceName, err := hclCreator(req)
	if err != nil {
		tb.Error(err)
	}
	return ClusterInfo{
		ProjectID:        req.ProjectID,
		Name:             clusterName,
		TerraformNameRef: fmt.Sprintf("%s.name", clusterResourceName),
		ResourceName:     clusterResourceName,
		TerraformStr:     clusterTerraformStr,
	}
}

func ExistingClusterUsed() bool {
	projectID, clusterName := existingProjectIDClusterName()
	return clusterName != "" && projectID != ""
}

func existingProjectIDClusterName() (projectID, clusterName string) {
	return os.Getenv("MONGODB_ATLAS_PROJECT_ID"), os.Getenv("MONGODB_ATLAS_CLUSTER_NAME")
}

func existingStreamInstanceUsed() bool {
	return existingStreamInstanceName() != "" && projectIDLocal() != ""
}

func existingStreamInstanceName() string {
	return os.Getenv("MONGODB_ATLAS_STREAM_INSTANCE_NAME")
}

// ReplicationSpecRequest can be used to customize the ReplicationSpecs of a Cluster.
// No fields are required.
// Use `ExtraRegionConfigs` to specify multiple region configs.
type ReplicationSpecRequest struct {
	ZoneName                 string
	Region                   string
	InstanceSize             string
	ProviderName             string
	EbsVolumeType            string
	ExtraRegionConfigs       []ReplicationSpecRequest
	NodeCount                int
	NodeCountReadOnly        int
	Priority                 int
	AutoScalingDiskGbEnabled bool
}

func (r *ReplicationSpecRequest) AddDefaults() {
	if r.Priority == 0 {
		r.Priority = 7
	}
	if r.NodeCount == 0 {
		r.NodeCount = 3
	}
	if r.ZoneName == "" {
		r.ZoneName = "Zone 1"
	}
	if r.Region == "" {
		r.Region = "US_WEST_2"
	}
	if r.InstanceSize == "" {
		r.InstanceSize = "M10"
	}
	if r.ProviderName == "" {
		r.ProviderName = constant.AWS
	}
}

func (r *ReplicationSpecRequest) AllRegionConfigs() []admin.CloudRegionConfig20240805 {
	config := cloudRegionConfig(*r)
	configs := []admin.CloudRegionConfig20240805{config}
	for i := range r.ExtraRegionConfigs {
		extra := r.ExtraRegionConfigs[i]
		configs = append(configs, cloudRegionConfig(extra))
	}
	return configs
}

func replicationSpec(req *ReplicationSpecRequest) admin.ReplicationSpec20240805 {
	if req == nil {
		req = new(ReplicationSpecRequest)
	}
	req.AddDefaults()
	regionConfigs := req.AllRegionConfigs()
	return admin.ReplicationSpec20240805{
		ZoneName:      &req.ZoneName,
		RegionConfigs: &regionConfigs,
	}
}

func cloudRegionConfig(req ReplicationSpecRequest) admin.CloudRegionConfig20240805 {
	req.AddDefaults()
	var readOnly admin.DedicatedHardwareSpec20240805
	if req.NodeCountReadOnly != 0 {
		readOnly = admin.DedicatedHardwareSpec20240805{
			NodeCount:    &req.NodeCountReadOnly,
			InstanceSize: &req.InstanceSize,
		}
	}
	return admin.CloudRegionConfig20240805{
		RegionName:   &req.Region,
		Priority:     &req.Priority,
		ProviderName: &req.ProviderName,
		ElectableSpecs: &admin.HardwareSpec20240805{
			InstanceSize:  &req.InstanceSize,
			NodeCount:     &req.NodeCount,
			EbsVolumeType: conversion.StringPtr(req.EbsVolumeType),
		},
		ReadOnlySpecs: &readOnly,
		AutoScaling: &admin.AdvancedAutoScalingSettings{
			DiskGB: &admin.DiskGBAutoScaling{Enabled: &req.AutoScalingDiskGbEnabled},
		},
	}
}
