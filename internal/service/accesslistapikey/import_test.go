package accesslistapikey_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/accesslistapikey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSplitAccessListAPIKeyImportID(t *testing.T) {
	tests := map[string]struct {
		importID    string
		orgID       string
		apiKeyID    string
		entry       string
		legacy      bool
		expectError bool
	}{
		"slash format with CIDR entry": {
			importID: "org/key/10.0.0.0/8",
			orgID:    "org",
			apiKeyID: "key",
			entry:    "10.0.0.0/8",
		},
		"legacy format with CIDR entry": {
			importID: "org-key-10.0.0.0/8",
			orgID:    "org",
			apiKeyID: "key",
			entry:    "10.0.0.0/8",
			legacy:   true,
		},
		"invalid format": {
			importID:    "org-key",
			expectError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			orgID, apiKeyID, entry, legacy, err := accesslistapikey.SplitAccessListAPIKeyImportID(tc.importID)
			if tc.expectError {
				require.Error(t, err)
				assert.False(t, legacy)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.orgID, orgID)
			assert.Equal(t, tc.apiKeyID, apiKeyID)
			assert.Equal(t, tc.entry, entry)
			assert.Equal(t, tc.legacy, legacy)
		})
	}
}
