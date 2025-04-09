package autogeneration_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogeneration"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenericImportOperation(t *testing.T) {
	testCases := []struct {
		expectedError *string
		expectedAttrs map[string]string
		name          string
		importID      string
		idAttributes  []string
	}{
		{
			name:         "Single attribute ID",
			importID:     "5c9d0a239ccf643e6a35ddasdf",
			idAttributes: []string{"project_id"},
			expectedAttrs: map[string]string{
				"project_id": "5c9d0a239ccf643e6a35ddasdf",
			},
		},
		{
			name:         "Multiple attribute ID",
			importID:     "5c9d0a239ccf643e6a35ddasdf/myCluster/us-east-1",
			idAttributes: []string{"project_id", "cluster_name", "region"},
			expectedAttrs: map[string]string{
				"project_id":   "5c9d0a239ccf643e6a35ddasdf",
				"cluster_name": "myCluster",
				"region":       "us-east-1",
			},
		},
		{
			name:          "Error: Wrong number of attributes",
			importID:      "5c9d0a239ccf643e6a35ddasdf/myCluster",
			idAttributes:  []string{"project_id", "cluster_name", "region"},
			expectedError: conversion.StringPtr("Expected format: project_id/cluster_name/region"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			attrValues, err := autogeneration.ProcessImportID(tc.importID, tc.idAttributes)

			if tc.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, *tc.expectedError, err.Error())
				return
			}
			assert.Equal(t, tc.expectedAttrs, attrValues)
		})
	}
}
