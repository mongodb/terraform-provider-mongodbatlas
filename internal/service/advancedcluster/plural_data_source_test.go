package advancedcluster_test

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemoveClusterNotFoundErrors(t *testing.T) {
	tests := map[string]struct {
		input    diag.Diagnostics
		expected diag.Diagnostics
	}{
		"empty diagnostics": {
			input:    diag.Diagnostics{},
			expected: diag.Diagnostics{},
		},
		"single cluster not found": {
			input: diag.Diagnostics{
				diag.NewErrorDiagnostic("Cluster Not Found", "CLUSTER_NOT_FOUND"),
			},
			expected: diag.Diagnostics{},
		},
		"only cluster not found errors": {
			input: diag.Diagnostics{
				diag.NewErrorDiagnostic("Error", "CLUSTER_NOT_FOUND"),
				diag.NewErrorDiagnostic("Error", "CLUSTER_NOT_FOUND"),
			},
			expected: diag.Diagnostics{},
		},
		"mixed errors with cluster not found": {
			input: diag.Diagnostics{
				diag.NewErrorDiagnostic("Other Error", "Some other error"),
				diag.NewErrorDiagnostic("Cluster Not Found", "CLUSTER_NOT_FOUND"),
			},
			expected: diag.Diagnostics{
				diag.NewErrorDiagnostic("Other Error", "Some other error"),
			},
		},
		"other errors only": {
			input: diag.Diagnostics{
				diag.NewErrorDiagnostic("Error 1", "Some error"),
				diag.NewErrorDiagnostic("Error 2", "Another error"),
			},
			expected: diag.Diagnostics{
				diag.NewErrorDiagnostic("Error 1", "Some error"),
				diag.NewErrorDiagnostic("Error 2", "Another error"),
			},
		},
		"warnings with cluster not found error": {
			input: diag.Diagnostics{
				diag.NewWarningDiagnostic("Warning", "Some warning"),
				diag.NewErrorDiagnostic("Cluster Not Found", "CLUSTER_NOT_FOUND"),
			},
			expected: diag.Diagnostics{
				diag.NewWarningDiagnostic("Warning", "Some warning"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			advancedcluster.RemoveClusterNotFoundErrors(&tc.input)
			assert.Equal(t, tc.expected, tc.input)
		})
	}
}

// TestPluralDataSourceRetriesTransientClusterNotFound covers a cluster list that fails
// with CLUSTER_NOT_FOUND before succeeding, which is what Atlas can return while a
// cluster is deleted concurrently. The retry must let the data source return the
// project's clusters instead of an empty result set.
func TestPluralDataSourceRetriesTransientClusterNotFound(t *testing.T) {
	storageConfigCredentials(t)
	mock := &pluralDataSourceHTTPMock{listFailures: 2}
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { require.NoError(t, unit.MockConfigAdvancedCluster.RunBeforeEach()) },
		ProtoV6ProviderFactories: unit.TestAccProviderV6FactoriesWithMock(t, mock),
		Steps: []resource.TestStep{
			{
				Config: pluralDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourcePluralName, "results.#", "1"),
					resource.TestCheckResourceAttr(dataSourcePluralName, "results.0.name", "example"),
				),
			},
		},
	})
	require.Greater(t, mock.listCalls, mock.listFailures, "the cluster list should be retried after the transient error")
}

// TestPluralDataSourceEmptyProjectReturnsEmptyList covers a project without clusters:
// `results` must be an empty list instead of null.
func TestPluralDataSourceEmptyProjectReturnsEmptyList(t *testing.T) {
	storageConfigCredentials(t)
	mock := &pluralDataSourceHTTPMock{emptyList: true}
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { require.NoError(t, unit.MockConfigAdvancedCluster.RunBeforeEach()) },
		ProtoV6ProviderFactories: unit.TestAccProviderV6FactoriesWithMock(t, mock),
		Steps: []resource.TestStep{
			{
				Config: pluralDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(dataSourcePluralName, tfjsonpath.New("results"), knownvalue.ListExact([]knownvalue.Check{})),
				},
			},
		},
	})
}

const pluralDataSourceConfig = `
provider "mongodbatlas" {
  public_key  = "test-public-key"
  private_key = "test-private-key"
  base_url    = "https://atlas.invalid/"
}

data "mongodbatlas_advanced_clusters" "test" {
  project_id = "111111111111111111111111"
}
`

const pluralDataSourceClusterResponse = `{
  "id": "333333333333333333333333", "groupId": "111111111111111111111111", "name": "example",
  "clusterType": "REPLICASET", "databaseEdition": "INFINITE", "stateName": "IDLE",
  "mongoDBMajorVersion": "8.0", "mongoDBVersion": "8.0.5", "versionReleaseSystem": "LTS",
  "rootCertType": "ISRGROOTX1", "encryptionAtRestProvider": "NONE", "replicaSetScalingStrategy": "SEQUENTIAL",
  "replicationSpecs": [{
    "id": "444444444444444444444444", "zoneId": "555555555555555555555555", "zoneName": "Zone 1",
    "regionConfigs": [{
      "providerName": "AWS", "regionName": "US_EAST_1", "priority": 7,
      "electableSpecs": {"instanceSize": "M10", "nodeCount": 2, "diskSizeGB": 10, "diskIOPS": 3000, "ebsVolumeType": "STANDARD"},
      "autoScaling": {"compute": {"enabled": true, "maxInstanceSize": "M20", "minInstanceSize": "M10", "scaleDownEnabled": false}, "diskGB": {"enabled": false}}
    }]
  }]
}`

// pluralDataSourceHTTPMock serves a cluster list that returns CLUSTER_NOT_FOUND for the
// first listFailures requests before succeeding, and the supporting reads the plural
// data source performs. Set emptyList to serve an empty cluster list.
type pluralDataSourceHTTPMock struct {
	mu           sync.Mutex
	listFailures int
	listCalls    int
	emptyList    bool
}

func (m *pluralDataSourceHTTPMock) ModifyHTTPClient(client *http.Client) error {
	client.Transport = m
	return nil
}

func (m *pluralDataSourceHTTPMock) ResetHTTPClient(*http.Client) {}

func (m *pluralDataSourceHTTPMock) RoundTrip(req *http.Request) (*http.Response, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	const groupPath = "/api/atlas/v2/groups/111111111111111111111111"
	status := http.StatusOK
	body := "{}"
	switch {
	case req.Method == http.MethodGet && req.URL.Path == groupPath+"/clusters":
		m.listCalls++
		switch {
		case m.listCalls <= m.listFailures:
			status, body = http.StatusNotFound, `{"errorCode":"CLUSTER_NOT_FOUND","error":404}`
		case m.emptyList:
			body = `{"results":[],"totalCount":0}`
		default:
			body = `{"results":[` + pluralDataSourceClusterResponse + `],"totalCount":1}`
		}
	case req.Method == http.MethodGet && req.URL.Path == groupPath+"/containers":
		body = `{"results":[{"id":"222222222222222222222222","providerName":"AWS","regionName":"US_EAST_1"}],"totalCount":1}`
	case req.Method == http.MethodGet && req.URL.Path == groupPath+"/clusters/example/processArgs":
	case req.Method == http.MethodGet && req.URL.Path == groupPath+"/flexClusters":
		body = `{"results":[],"totalCount":0}`
	default:
		return nil, fmt.Errorf("unexpected mocked Atlas request: %s %s", req.Method, req.URL.Path)
	}
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": {"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
		Request:    req,
	}, nil
}
