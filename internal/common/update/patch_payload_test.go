package update_test

import (
	"fmt"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/update"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312004/admin"
)

func TestPatchReplicationSpecs(t *testing.T) {
	var (
		rp1         = replicationSpec{placeholderIndex: 1}.toAdmin()
		rp2         = replicationSpec{placeholderIndex: 2}.toAdmin()
		rp3         = replicationSpec{placeholderIndex: 3}.toAdmin()
		rp1ZoneName = rp1.GetZoneName()
		rp1ID       = rp1.GetId()
		rp1ZoneID   = rp1.GetZoneId()
		idGlobal    = "id_root"

		clusterName           = "my-cluster"
		rootNameUpdated       = "my-cluster-updated"
		stateReplicationSpecs = []admin.ReplicationSpec20240805{
			rp1,
		}
		state = admin.ClusterDescription20240805{
			Id:               &idGlobal,
			Name:             &clusterName,
			ReplicationSpecs: &stateReplicationSpecs,
		}
		stateWithReplicationSpecs = func(specs []replicationSpec, id, name string) *admin.ClusterDescription20240805 {
			newSpecs := make([]admin.ReplicationSpec20240805, len(specs))
			for i := range specs {
				newSpecs[i] = specs[i].toAdmin()
			}
			cd := admin.ClusterDescription20240805{
				ReplicationSpecs: &newSpecs,
			}
			if id != "" {
				cd.Id = &id
			}
			if name != "" {
				cd.Name = &name
			}
			return &cd
		}
		planNameDifferentAndEnableBackup = admin.ClusterDescription20240805{
			Name:          &rootNameUpdated,
			BackupEnabled: conversion.Pointer(true),
		}
		planNoChanges = admin.ClusterDescription20240805{
			ReplicationSpecs: &[]admin.ReplicationSpec20240805{
				{
					ZoneName: &rp1ZoneName,
				},
			},
		}
		testCases = map[string]struct {
			state         *admin.ClusterDescription20240805
			plan          *admin.ClusterDescription20240805
			patchExpected *admin.ClusterDescription20240805
			options       []update.PatchOptions
		}{
			"ComputedValues from the state are added to nested attribute plan and unchanged attributes are not included": {
				state: &state,
				plan:  stateWithReplicationSpecs([]replicationSpec{{zoneName: "newName"}}, "", ""),
				patchExpected: &admin.ClusterDescription20240805{
					ReplicationSpecs: &[]admin.ReplicationSpec20240805{
						{
							Id:       &rp1ID,
							ZoneId:   &rp1ZoneID,
							ZoneName: conversion.Pointer("newName"),
						},
					},
				},
			},
			"New list entry added should be included": {
				state: &state,
				plan:  stateWithReplicationSpecs([]replicationSpec{{placeholderIndex: 1}, {zoneName: "zone2"}}, "", ""),
				patchExpected: &admin.ClusterDescription20240805{
					ReplicationSpecs: &[]admin.ReplicationSpec20240805{
						rp1,
						{
							ZoneName: conversion.Pointer("zone2"),
						},
					},
				},
			},
			"Removed list entry should be detected": {
				state: stateWithReplicationSpecs([]replicationSpec{{placeholderIndex: 1}, {placeholderIndex: 2}}, "", ""),
				plan:  stateWithReplicationSpecs([]replicationSpec{{placeholderIndex: 1}}, "", ""),
				patchExpected: &admin.ClusterDescription20240805{
					ReplicationSpecs: &[]admin.ReplicationSpec20240805{
						rp1,
					},
				},
			},
			"Added list entry in the middle should be detected": {
				state: stateWithReplicationSpecs([]replicationSpec{{placeholderIndex: 1}, {placeholderIndex: 2}}, "", ""),
				plan:  stateWithReplicationSpecs([]replicationSpec{{placeholderIndex: 1}, {placeholderIndex: 3}, {placeholderIndex: 2}}, "", ""),
				patchExpected: &admin.ClusterDescription20240805{
					ReplicationSpecs: &[]admin.ReplicationSpec20240805{
						rp1,
						rp3,
						rp2,
					},
				},
			},
			"Removed list entry in the middle should be detected": {
				state: stateWithReplicationSpecs([]replicationSpec{{placeholderIndex: 1}, {placeholderIndex: 2}, {placeholderIndex: 3}}, "", ""),
				plan:  stateWithReplicationSpecs([]replicationSpec{{placeholderIndex: 1}, {placeholderIndex: 3}}, "", ""),
				patchExpected: &admin.ClusterDescription20240805{
					ReplicationSpecs: &[]admin.ReplicationSpec20240805{
						rp1,
						rp3,
					},
				},
			},
			"Region Config changes are included in patch": {
				state: &admin.ClusterDescription20240805{
					ReplicationSpecs: &[]admin.ReplicationSpec20240805{
						{
							Id: &rp1ID,
							RegionConfigs: &[]admin.CloudRegionConfig20240805{
								{
									Priority: conversion.Pointer(1),
								},
							},
						},
					},
				},
				plan: &admin.ClusterDescription20240805{
					ReplicationSpecs: &[]admin.ReplicationSpec20240805{
						{
							Id: &rp1ID,
							RegionConfigs: &[]admin.CloudRegionConfig20240805{
								{
									Priority: conversion.Pointer(1),
								},
								{
									Priority: conversion.Pointer(2),
								},
							},
						},
					},
				},
				patchExpected: &admin.ClusterDescription20240805{
					ReplicationSpecs: &[]admin.ReplicationSpec20240805{
						{
							Id: &rp1ID,
							RegionConfigs: &[]admin.CloudRegionConfig20240805{
								{
									Priority: conversion.Pointer(1),
								},
								{
									Priority: conversion.Pointer(2),
								},
							},
						},
					},
				},
			},
			"Name change and backup enabled added": {
				state: &state,
				plan:  &planNameDifferentAndEnableBackup,
				patchExpected: &admin.ClusterDescription20240805{
					Name:          &rootNameUpdated,
					BackupEnabled: conversion.Pointer(true),
				},
			},
			"No Changes when only computed attributes are not in plan": {
				state:         &state,
				plan:          &planNoChanges,
				patchExpected: nil,
			},
			"Forced changes when forceUpdateAttr set": {
				state: &state,
				plan:  &planNoChanges,
				patchExpected: &admin.ClusterDescription20240805{
					ReplicationSpecs: &stateReplicationSpecs,
				},
				options: []update.PatchOptions{
					{ForceUpdateAttr: []string{"replicationSpecs"}},
				},
			},
			"Force changes when forceUpdateAttr set and state==plan": {
				state: &state,
				plan:  &state,
				patchExpected: &admin.ClusterDescription20240805{
					ReplicationSpecs: &stateReplicationSpecs,
				},
				options: []update.PatchOptions{
					{ForceUpdateAttr: []string{"replicationSpecs"}},
				},
			},
			"Empty array should return no changes": {
				state: &admin.ClusterDescription20240805{
					Labels: &[]admin.ComponentLabel{},
				},
				plan: &admin.ClusterDescription20240805{
					Labels: &[]admin.ComponentLabel{},
				},
				patchExpected: nil,
			},
			"diskSizeGb ignored in state": {
				state:         clusterDescriptionDiskSizeNodeCount(50.0, 3, conversion.Pointer(50.0), 0, conversion.Pointer(3500)),
				plan:          clusterDescriptionDiskSizeNodeCount(55.0, 3, nil, 0, nil),
				patchExpected: clusterDescriptionDiskSizeNodeCount(55.0, 3, nil, 0, conversion.Pointer(3500)),
				options: []update.PatchOptions{
					{
						IgnoreInStateSuffix: []string{"diskSizeGB"},
					},
				},
			},
			"regionConfigs ignored in state but diskIOPS included": {
				state:         clusterDescriptionDiskSizeNodeCount(50.0, 3, conversion.Pointer(50.0), 0, conversion.Pointer(3500)),
				plan:          clusterDescriptionDiskSizeNodeCount(55.0, 3, nil, 0, nil),
				patchExpected: clusterDescriptionDiskSizeNodeCount(55.0, 3, nil, 0, conversion.Pointer(3500)),
				options: []update.PatchOptions{
					{
						IgnoreInStatePrefix:  []string{"regionConfigs"},
						IncludeInStateSuffix: []string{"diskIOPS"},
					},
				},
			},
		}
	)
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			patchReq, err := update.PatchPayload(tc.state, tc.plan, tc.options...)
			require.NoError(t, err)
			assert.Equal(t, tc.patchExpected, patchReq)
		})
	}
}

func TestPatchAdvancedConfig(t *testing.T) {
	var (
		state = admin.ClusterDescriptionProcessArgs20240805{
			JavascriptEnabled: conversion.Pointer(true),
		}
		testCases = map[string]struct {
			state         *admin.ClusterDescriptionProcessArgs20240805
			plan          *admin.ClusterDescriptionProcessArgs20240805
			patchExpected *admin.ClusterDescriptionProcessArgs20240805
			options       []update.PatchOptions
		}{
			"JavascriptEnabled is set to false": {
				state: &state,
				plan: &admin.ClusterDescriptionProcessArgs20240805{
					JavascriptEnabled: conversion.Pointer(false),
				},
				patchExpected: &admin.ClusterDescriptionProcessArgs20240805{
					JavascriptEnabled: conversion.Pointer(false),
				},
			},
			"JavascriptEnabled is set to null leads to no changes": {
				state:         &state,
				plan:          &admin.ClusterDescriptionProcessArgs20240805{},
				patchExpected: nil,
			},
			"JavascriptEnabled state equals plan leads to no changes": {
				state:         &state,
				plan:          &state,
				patchExpected: nil,
			},
			"Adding NoTableScan changes the plan payload and but doesn't include old value of JavascriptEnabled": {
				state: &state,
				plan: &admin.ClusterDescriptionProcessArgs20240805{
					NoTableScan: conversion.Pointer(true),
				},
				patchExpected: &admin.ClusterDescriptionProcessArgs20240805{
					NoTableScan: conversion.Pointer(true),
				},
			},
			"Nil plan should return no changes": {
				state:         &state,
				plan:          nil,
				patchExpected: nil,
			},
		}
	)
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			patchReq, err := update.PatchPayload(tc.state, tc.plan, tc.options...)
			require.NoError(t, err)
			assert.Equal(t, tc.patchExpected, patchReq)
		})
	}
}

func TestIsEmpty(t *testing.T) {
	assert.True(t, update.IsZeroValues(&admin.ClusterDescription20240805{}))
	var myVar admin.ClusterDescription20240805
	assert.True(t, update.IsZeroValues(&myVar))
	assert.False(t, update.IsZeroValues(&admin.ClusterDescription20240805{Name: conversion.Pointer("my-cluster")}))
}

type replicationSpec struct {
	id               string
	zoneName         string
	zoneID           string
	placeholderIndex int
}

func (r replicationSpec) toAdmin() admin.ReplicationSpec20240805 {
	var (
		placeholderID       = "replicationSpec%d_id"
		placeholderZoneID   = "replicationSpec%d_zoneId"
		placeholderZoneName = "replicationSpec%d_zoneName"
	)
	index := r.placeholderIndex
	if index != 0 {
		if r.id == "" {
			r.id = fmt.Sprintf(placeholderID, index)
		}
		if r.zoneName == "" {
			r.zoneName = fmt.Sprintf(placeholderZoneName, index)
		}
		if r.zoneID == "" {
			r.zoneID = fmt.Sprintf(placeholderZoneID, index)
		}
	}
	spec := admin.ReplicationSpec20240805{}
	if r.id != "" {
		spec.SetId(r.id)
	}
	if r.zoneID != "" {
		spec.SetZoneId(r.zoneID)
	}
	if r.zoneName != "" {
		spec.SetZoneName(r.zoneName)
	}
	return spec
}

func clusterDescriptionDiskSizeNodeCount(diskSizeGBElectable float64, nodeCountElectable int, diskSizeGBReadOnly *float64, nodeCountReadOnly int, diskIopsState *int) *admin.ClusterDescription20240805 {
	return &admin.ClusterDescription20240805{
		ReplicationSpecs: &[]admin.ReplicationSpec20240805{
			{
				RegionConfigs: &[]admin.CloudRegionConfig20240805{
					{
						ElectableSpecs: &admin.HardwareSpec20240805{
							NodeCount:  &nodeCountElectable,
							DiskSizeGB: &diskSizeGBElectable,
							DiskIOPS:   diskIopsState,
						},
						ReadOnlySpecs: &admin.DedicatedHardwareSpec20240805{
							NodeCount:  &nodeCountReadOnly,
							DiskSizeGB: diskSizeGBReadOnly,
							DiskIOPS:   diskIopsState,
						},
					},
				},
			},
		},
	}
}
