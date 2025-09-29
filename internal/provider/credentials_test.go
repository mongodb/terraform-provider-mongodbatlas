package provider_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_deriveSTSRegionFromEndpoint(t *testing.T) {
	testCases := map[string]struct {
		input    string
		expected string
	}{
		"empty endpoint": {
			input:    "",
			expected: "",
		},
		"global endpoint": {
			input:    "https://sts.amazonaws.com",
			expected: provider.DefaultRegionSTS,
		},
		"regional": {
			input:    "https://sts.us-east-1.amazonaws.com/",
			expected: "us-east-1",
		},
		"regional eu-north-1": {
			input:    "https://sts.eu-north-1.amazonaws.com/",
			expected: "eu-north-1",
		},
		"malformed url": {
			input:    "://not-a-url",
			expected: provider.DefaultRegionSTS,
		},
		"unexpected host shape": {
			input:    "https://sts.something-weird",
			expected: provider.DefaultRegionSTS,
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			t.Parallel()
			got := provider.DeriveSTSRegionFromEndpoint(tc.input)
			if got != tc.expected {
				t.Fatalf("deriveSTSRegionFromEndpoint(%q) = %q; want %q", tc.input, got, tc.expected)
			}
		})
	}
}

func Test_resolveSTSEndpoint(t *testing.T) {
	testCases := map[string]struct {
		stsEndpoint   string
		secretsRegion string
		expectedURL   string
		expectedSign  string
	}{
		"explicit regional endpoint": {
			stsEndpoint:   "https://sts.eu-north-1.amazonaws.com/",
			secretsRegion: "us-east-1",
			expectedURL:   "https://sts.eu-north-1.amazonaws.com/",
			expectedSign:  "eu-north-1",
		},
		"global endpoint - us-east-1 signing": {
			stsEndpoint:   "https://sts.amazonaws.com",
			secretsRegion: "eu-west-1",
			expectedURL:   "https://sts.amazonaws.com",
			expectedSign:  provider.DefaultRegionSTS,
		},
		"no endpoint - uses secrets region": {
			stsEndpoint:   "",
			secretsRegion: "us-west-2",
			expectedURL:   "https://sts.us-west-2.amazonaws.com/",
			expectedSign:  "us-west-2",
		},
		"no endpoint and empty region": {
			stsEndpoint:   "",
			secretsRegion: "",
			expectedURL:   "https://sts.us-east-1.amazonaws.com/",
			expectedSign:  provider.DefaultRegionSTS,
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			ep, err := provider.ResolveSTSEndpoint(tc.stsEndpoint, tc.secretsRegion)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedURL, ep.URL)
			assert.Equal(t, tc.expectedSign, ep.SigningRegion)
		})
	}
}
