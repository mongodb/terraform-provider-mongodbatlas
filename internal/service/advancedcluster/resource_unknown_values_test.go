package advancedcluster_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
)

const unknownValuesProjectID = "111111111111111111111111"

// TestAccMockableAdvancedCluster_unknownTagsAndLabels covers updates whose tags or labels come from an
// expression that is unresolved at plan time, such as a module output or merge(). The whole map is unknown
// then, and converting it into the Atlas request failed with a Value Conversion Error before this fix,
// during plan and before any Atlas call.
func TestAccMockableAdvancedCluster_unknownTagsAndLabels(t *testing.T) {
	testCases := map[string]string{
		"unknown tags map": `
resource "terraform_data" "values" {
  input = { owner = "team-atlas" }
}
` + unknownValuesClusterConfig(`tags = terraform_data.values.output`),
		"unknown labels map": `
resource "terraform_data" "values" {
  input = { owner = "team-atlas" }
}
` + unknownValuesClusterConfig(`labels = terraform_data.values.output`),
		// A known map holding an unknown value is a different shape that already worked, keep it covered.
		"unknown tag value": `
resource "terraform_data" "owner" {
  input = "team-atlas"
}
` + unknownValuesClusterConfig(`tags = { owner = terraform_data.owner.output }`),
	}
	for name, updatedConfig := range testCases {
		t.Run(name, func(t *testing.T) {
			for _, envVar := range []string{"ASSUME_ROLE_ARN", "TF_VAR_ASSUME_ROLE_ARN", "MONGODB_ATLAS_CLIENT_ID", "MONGODB_ATLAS_CLIENT_SECRET"} {
				t.Setenv(envVar, "")
			}
			mock := &unknownValuesHTTPMock{}
			resource.UnitTest(t, resource.TestCase{
				PreCheck:                 func() { require.NoError(t, unit.MockConfigAdvancedCluster.RunBeforeEach()) },
				ProtoV6ProviderFactories: unit.TestAccProviderV6FactoriesWithMock(t, mock),
				Steps: []resource.TestStep{
					// The cluster must already exist: ModifyPlan returns early while the state is null.
					{Config: unknownValuesProviderConfig + unknownValuesClusterConfig("")},
					{
						Config: unknownValuesProviderConfig + updatedConfig,
						Check: resource.TestCheckResourceAttr(
							"mongodbatlas_advanced_cluster.test", "replication_specs.0.region_configs.0.electable_specs.node_count", "3"),
					},
				},
			})
		})
	}
}

const unknownValuesProviderConfig = `
provider "mongodbatlas" {
  public_key  = "test-public-key"
  private_key = "test-private-key"
  base_url    = "https://atlas.invalid/"
}
`

func unknownValuesClusterConfig(extraAttribute string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_advanced_cluster" "test" {
  project_id   = %[1]q
  name         = "example"
  cluster_type = "REPLICASET"
  %[2]s
  replication_specs = [{
    region_configs = [{
      provider_name   = "AWS"
      region_name     = "US_EAST_1"
      priority        = 7
      electable_specs = { instance_size = "M10", node_count = 3 }
    }]
  }]
}
`, unknownValuesProjectID, extraAttribute)
}

// unknownValuesHTTPMock serves a minimal cluster lifecycle and echoes back the tags and labels it receives.
type unknownValuesHTTPMock struct {
	tags    []any
	labels  []any
	mu      sync.Mutex
	deleted bool
}

func (m *unknownValuesHTTPMock) ModifyHTTPClient(client *http.Client) error {
	client.Transport = m
	return nil
}

func (m *unknownValuesHTTPMock) ResetHTTPClient(*http.Client) {}

func (m *unknownValuesHTTPMock) RoundTrip(req *http.Request) (*http.Response, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var (
		groupPath   = "/api/atlas/v2/groups/" + unknownValuesProjectID
		clustersURL = groupPath + "/clusters"
		clusterURL  = clustersURL + "/example"
		status      = http.StatusOK
		body        string
	)
	switch {
	case req.Method == http.MethodPost && req.URL.Path == clustersURL,
		req.Method == http.MethodPatch && req.URL.Path == clusterURL:
		var payload map[string]any
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			return nil, err
		}
		m.tags, _ = payload["tags"].([]any)
		m.labels, _ = payload["labels"].([]any)
		body = m.clusterJSON()
	case req.Method == http.MethodGet && req.URL.Path == clusterURL:
		body = m.clusterJSON()
		if m.deleted {
			status, body = http.StatusNotFound, `{"errorCode":"CLUSTER_NOT_FOUND","error":404}`
		}
	case req.Method == http.MethodGet && req.URL.Path == groupPath+"/containers":
		body = `{"results":[{"id":"222222222222222222222222","providerName":"AWS","regionName":"US_EAST_1"}],"totalCount":1}`
	case req.Method == http.MethodGet && req.URL.Path == clusterURL+"/processArgs":
		body = "{}"
	case req.Method == http.MethodDelete && req.URL.Path == clusterURL:
		m.deleted, body = true, "{}"
	default:
		return nil, fmt.Errorf("unexpected mocked Atlas request: %s %s", req.Method, req.URL.Path)
	}
	return &http.Response{
		StatusCode: status, Header: http.Header{"Content-Type": {"application/json"}},
		Body: io.NopCloser(strings.NewReader(body)), Request: req,
	}, nil
}

func (m *unknownValuesHTTPMock) clusterJSON() string {
	cluster := map[string]any{
		"id": "333333333333333333333333", "groupId": unknownValuesProjectID, "name": "example",
		"clusterType": "REPLICASET", "stateName": "IDLE", "mongoDBMajorVersion": "8.0",
		"mongoDBVersion": "8.0.5", "versionReleaseSystem": "LTS", "rootCertType": "ISRGROOTX1",
		"encryptionAtRestProvider": "NONE", "replicaSetScalingStrategy": "SEQUENTIAL",
		"tags": m.tags, "labels": m.labels,
		"replicationSpecs": []any{map[string]any{
			"id": "444444444444444444444444", "zoneId": "555555555555555555555555", "zoneName": "Zone 1",
			"regionConfigs": []any{map[string]any{
				"providerName": "AWS", "regionName": "US_EAST_1", "priority": 7,
				"electableSpecs": map[string]any{
					"instanceSize": "M10", "nodeCount": 3, "diskSizeGB": 10, "diskIOPS": 3000, "ebsVolumeType": "STANDARD",
				},
			}},
		}},
	}
	encoded, err := json.Marshal(cluster)
	if err != nil {
		return "{}"
	}
	return string(encoded)
}
