package autogen_test

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleImport(t *testing.T) {
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
			expectedError: conversion.StringPtr(fmt.Sprintf(autogen.ExpectedErrorMsg, "project_id/cluster_name/region")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			attrValues, err := autogen.ProcessImportID(tc.importID, tc.idAttributes)

			if tc.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, *tc.expectedError, err.Error())
				return
			}
			assert.Equal(t, tc.expectedAttrs, attrValues)
		})
	}
}

func TestHandleImportWithCustomHook(t *testing.T) {
	hook := &testPreImportHook{
		preImportFunc: func(id string) (string, error) {
			if strings.Contains(id, "/") {
				return id, nil
			}

			re := regexp.MustCompile(`^(.*)-([0-9a-fA-F]{24})-(.*)$`)
			matches := re.FindStringSubmatch(id)
			if len(matches) != 4 || matches[1] == "" || matches[3] == "" {
				return "", fmt.Errorf("use one of the formats: {project_id}/{workspace_name}/{connection_name} or {workspace_name}-{project_id}-{connection_name}")
			}
			return fmt.Sprintf("%s/%s/%s", matches[2], matches[1], matches[3]), nil
		},
	}

	normalizedID, err := hook.PreImport("myWorkspace-507f1f77bcf86cd799439011-myConnection")
	require.NoError(t, err)
	assert.Equal(t, "507f1f77bcf86cd799439011/myWorkspace/myConnection", normalizedID)

	req := resource.ImportStateRequest{ID: "bad-format"}
	resp := &resource.ImportStateResponse{}
	autogen.HandleImport(context.Background(), []string{"project_id", "workspace_name", "connection_name"}, req, resp, hook)
	require.True(t, resp.Diagnostics.HasError())

	_, err = hook.PreImport("bad-format")
	require.Error(t, err)
}

type testPreImportHook struct {
	preImportFunc func(string) (string, error)
}

func (h *testPreImportHook) PreImport(id string) (string, error) {
	return h.preImportFunc(id)
}
