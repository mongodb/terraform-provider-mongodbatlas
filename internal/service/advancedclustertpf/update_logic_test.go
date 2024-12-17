package advancedclustertpf_test

import (
	"context"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedclustertpf"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20241113003/admin"
)

func TestAlignStateReplicationSpecs(t *testing.T) {
	var (
		zone1Region1Spec = admin.ReplicationSpec20240805{
			ZoneName: conversion.Pointer("zone1"),
			RegionConfigs: &[]admin.CloudRegionConfig20240805{
				{RegionName: conversion.Pointer("region1")},
			},
		}
		zone1Region2Spec = admin.ReplicationSpec20240805{
			ZoneName: conversion.Pointer("zone1"),
			RegionConfigs: &[]admin.CloudRegionConfig20240805{
				{RegionName: conversion.Pointer("region2")},
			},
		}
		zone1Region1And2Spec = admin.ReplicationSpec20240805{
			ZoneName: conversion.Pointer("zone1"),
			RegionConfigs: &[]admin.CloudRegionConfig20240805{
				{RegionName: conversion.Pointer("region1")},
				{RegionName: conversion.Pointer("region2")},
			},
		}
		zone2Region2Spec = admin.ReplicationSpec20240805{
			ZoneName: conversion.Pointer("zone2"),
			RegionConfigs: &[]admin.CloudRegionConfig20240805{
				{RegionName: conversion.Pointer("region2")},
			},
		}
		emptyReplicationSpec = admin.ReplicationSpec20240805{}
	)

	tests := map[string]struct {
		state    []admin.ReplicationSpec20240805
		plan     []admin.ReplicationSpec20240805
		expected []admin.ReplicationSpec20240805
	}{
		"No replication specs in state or plan": {
			state:    []admin.ReplicationSpec20240805{},
			plan:     []admin.ReplicationSpec20240805{},
			expected: []admin.ReplicationSpec20240805{},
		},
		"Matching replication specs": {
			state:    []admin.ReplicationSpec20240805{zone1Region1Spec},
			plan:     []admin.ReplicationSpec20240805{zone1Region1Spec},
			expected: []admin.ReplicationSpec20240805{zone1Region1Spec},
		},
		"Different replication specs, should lead to an empty replication spec": {
			state:    []admin.ReplicationSpec20240805{zone1Region1Spec},
			plan:     []admin.ReplicationSpec20240805{zone2Region2Spec},
			expected: []admin.ReplicationSpec20240805{emptyReplicationSpec},
		},
		"Only region update should re-use replication spec": {
			state:    []admin.ReplicationSpec20240805{zone1Region1Spec},
			plan:     []admin.ReplicationSpec20240805{zone1Region2Spec},
			expected: []admin.ReplicationSpec20240805{zone1Region1Spec},
		},
		"State has more replication specs than plan": {
			state:    []admin.ReplicationSpec20240805{zone1Region1Spec, zone2Region2Spec},
			plan:     []admin.ReplicationSpec20240805{zone1Region1Spec},
			expected: []admin.ReplicationSpec20240805{zone1Region1Spec},
		},
		"Plan has more replication specs than state, should add an empty replication spec": {
			state:    []admin.ReplicationSpec20240805{zone1Region1Spec},
			plan:     []admin.ReplicationSpec20240805{zone1Region1Spec, zone2Region2Spec},
			expected: []admin.ReplicationSpec20240805{zone1Region1Spec, emptyReplicationSpec},
		},
		"Adding a replication spec in the middle should match the two old ones correctly": {
			state:    []admin.ReplicationSpec20240805{zone1Region1Spec, zone2Region2Spec},
			plan:     []admin.ReplicationSpec20240805{zone1Region1Spec, zone1Region2Spec, zone2Region2Spec},
			expected: []admin.ReplicationSpec20240805{zone1Region1Spec, emptyReplicationSpec, zone2Region2Spec},
		},
		"Adding a region config should match the existing one": {
			state:    []admin.ReplicationSpec20240805{zone1Region1Spec},
			plan:     []admin.ReplicationSpec20240805{zone1Region1And2Spec},
			expected: []admin.ReplicationSpec20240805{zone1Region1Spec},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			state := &admin.ClusterDescription20240805{ReplicationSpecs: &tc.state}
			plan := &admin.ClusterDescription20240805{ReplicationSpecs: &tc.plan}
			advancedclustertpf.AlignStateReplicationSpecs(context.Background(), state, plan)
			assert.Equal(t, &tc.expected, state.ReplicationSpecs)
		})
	}
}
