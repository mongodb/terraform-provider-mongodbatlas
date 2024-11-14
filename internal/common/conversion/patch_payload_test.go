package conversion_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20241023002/admin"
)

func TestJsonPatchReplicationSpecs(t *testing.T) {
	var (
		idGlobal                    = "id_root"
		idReplicationSpec1          = "id_replicationSpec1"
		replicationSpec1ZoneNameOld = "replicationSpec1_zoneName_old"
		replicationSpec1ZoneNameNew = "replicationSpec1_zoneName_new"
		replicationSpec1ZoneID      = "replicationSpec1_zoneId"
		replicationSpec2ZoneName    = "replicationSpec2_zoneName"
		rootName                    = "my-cluster"
		rootNameUpdated             = "my-cluster-updated"
		state                       = admin.ClusterDescription20240805{
			Id:   &idGlobal,
			Name: &rootName,
			ReplicationSpecs: &[]admin.ReplicationSpec20240805{
				{
					Id:       &idReplicationSpec1,
					ZoneId:   &replicationSpec1ZoneID,
					ZoneName: &replicationSpec1ZoneNameOld,
				},
			},
		}
		planOptionalUpdated = admin.ClusterDescription20240805{
			Name: &rootName,
			ReplicationSpecs: &[]admin.ReplicationSpec20240805{
				{
					ZoneName: &replicationSpec1ZoneNameNew,
				},
			},
		}
		planNewListEntry = admin.ClusterDescription20240805{
			ReplicationSpecs: &[]admin.ReplicationSpec20240805{
				{
					ZoneName: &replicationSpec1ZoneNameOld,
				},
				{
					ZoneName: &replicationSpec2ZoneName,
				},
			},
		}
		planNameDifferentAndEnableBackup = admin.ClusterDescription20240805{
			Name:          &rootNameUpdated,
			BackupEnabled: conversion.Pointer(true),
		}
		planNoChanges = admin.ClusterDescription20240805{
			ReplicationSpecs: &[]admin.ReplicationSpec20240805{
				{
					ZoneName: &replicationSpec1ZoneNameOld,
				},
			},
		}
		testCases = map[string]struct {
			state         *admin.ClusterDescription20240805
			plan          *admin.ClusterDescription20240805
			patchExpected *admin.ClusterDescription20240805
			noChanges     bool
		}{
			"ComputedValues from the state are added to plan and unchanged attributes are not included": {
				state: &state,
				plan:  &planOptionalUpdated,
				patchExpected: &admin.ClusterDescription20240805{
					ReplicationSpecs: &[]admin.ReplicationSpec20240805{
						{
							Id:       &idReplicationSpec1,
							ZoneId:   &replicationSpec1ZoneID,
							ZoneName: &replicationSpec1ZoneNameNew,
						},
					},
				},
			},
			"New list entry added should be included": {
				state: &state,
				plan:  &planNewListEntry,
				patchExpected: &admin.ClusterDescription20240805{
					ReplicationSpecs: &[]admin.ReplicationSpec20240805{
						{
							Id:       &idReplicationSpec1,
							ZoneId:   &replicationSpec1ZoneID,
							ZoneName: &replicationSpec1ZoneNameOld,
						},
						{
							ZoneName: &replicationSpec2ZoneName,
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
				noChanges:     true,
				patchExpected: &admin.ClusterDescription20240805{},
			},
		}
	)
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			patchReq := &admin.ClusterDescription20240805{}
			noChanges, err := conversion.PatchPayloadNoChanges(tc.state, tc.plan, patchReq)
			require.NoError(t, err)
			assert.Equal(t, tc.noChanges, noChanges)
			assert.Equal(t, tc.patchExpected, patchReq)
		})
	}
}
