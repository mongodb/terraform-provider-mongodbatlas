package unit

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func replaceVars(text string, vars map[string]string) string {
	for key, value := range vars {
		text = strings.ReplaceAll(text, fmt.Sprintf("{%s}", key), value)
	}
	return text
}

type statusText struct {
	Text   string `yaml:"text"`
	Status int    `yaml:"status"`
}

type RequestInfo struct {
	Version   string       `yaml:"version"`
	Method    string       `yaml:"method"`
	Path      string       `yaml:"path"`
	Text      string       `yaml:"text"`
	Responses []statusText `yaml:"responses"`
}

func (i *RequestInfo) id() string {
	return fmt.Sprintf("%s_%s_%s", i.Method, i.Path, i.Version)
}

func (i *RequestInfo) Match(method, urlPath, version string, vars map[string]string) bool {
	if i.Method != method {
		return false
	}
	selfPath := replaceVars(i.Path, vars)
	return selfPath == urlPath && i.Version == version
}

type stepRequests struct {
	DiffRequests     []RequestInfo `yaml:"diff_requests"`
	RequestResponses []RequestInfo `yaml:"request_responses"`
}

type mockHTTPData struct {
	Variables map[string]string `yaml:"variables"`
	Steps     []stepRequests    `yaml:"steps"`
	StepCount int               `yaml:"step_count"`
}

type MockHTTPDataConfig struct {
	AllowMissingRequests bool
	AllowReReadGet       bool
}

func parseTestDataConfigYAML(filePath string) (*mockHTTPData, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var testData mockHTTPData
	err = yaml.Unmarshal(data, &testData)
	if err != nil {
		return nil, err
	}
	return &testData, nil
}

func MockRoundTripper(t *testing.T, vars map[string]string, config *MockHTTPDataConfig) (http.RoundTripper, resource.TestCheckFunc) {
	t.Helper()
	testDir := "testdata"
	httpDataPath := path.Join(testDir, t.Name()+".yaml")
	data, err := parseTestDataConfigYAML(httpDataPath)
	require.NoError(t, err)
	myTransport := httpmock.NewMockTransport()
	var mockTransport http.RoundTripper = myTransport
	g := goldie.New(t, goldie.WithTestNameForDir(true), goldie.WithNameSuffix(".json"))
	tracker := requestTracker{data: data, g: g, vars: vars, t: t}
	if config != nil {
		tracker.allowMissingRequests = config.AllowMissingRequests
		tracker.allowReReadGet = config.AllowReReadGet
	}
	err = tracker.initStep()
	require.NoError(t, err)
	for _, method := range []string{"GET", "POST", "PUT", "DELETE", "PATCH"} {
		myTransport.RegisterRegexpResponder(method, regexp.MustCompile(".*"), tracker.receiveRequest(method))
	}
	return mockTransport, tracker.checkStepRequests
}

var _versionDatePattern = regexp.MustCompile(`(\d{4}-\d{2}-\d{2})`)

func ExtractVersion(contentType string) (string, error) {
	match := _versionDatePattern.FindStringSubmatch(contentType)
	if len(match) > 1 {
		return match[1], nil
	}
	return "", fmt.Errorf("could not extract version from %s header", contentType)
}

type requestTracker struct {
	t                    *testing.T
	g                    *goldie.Goldie
	data                 *mockHTTPData
	vars                 map[string]string
	usedResponses        map[string]int
	foundsDiffs          map[string]string
	currentStepIndex     int
	allowMissingRequests bool
	allowReReadGet       bool
}

func (r *requestTracker) requestFilename(requestID string) string {
	return strings.ReplaceAll(fmt.Sprintf("%02d_%s", r.currentStepIndex+1, requestID), "/", "_")
}

func (r *requestTracker) initStep() error {
	r.usedResponses = map[string]int{}
	r.foundsDiffs = map[string]string{}
	step := r.currentStep()
	if step == nil {
		return nil
	}
	for _, req := range step.DiffRequests {
		err := r.g.Update(r.t, r.requestFilename(req.id()), []byte(req.Text))
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *requestTracker) currentStep() *stepRequests {
	if r.currentStepIndex >= r.data.StepCount {
		return nil
	}
	return &r.data.Steps[r.currentStepIndex]
}

func (r *requestTracker) checkStepRequests(_ *terraform.State) error {
	missingRequests := []string{}
	for _, req := range r.currentStep().RequestResponses {
		missingRequestsCount := len(req.Responses) - r.usedResponses[req.id()]
		if missingRequestsCount > 0 {
			missingRequests = append(missingRequests, fmt.Sprintf("missing %d requests of %s", missingRequestsCount, req.id()))
		}
	}
	if r.allowMissingRequests {
		if len(missingRequests) > 0 {
			r.t.Logf("missing requests: %s", strings.Join(missingRequests, ", "))
		}
	} else {
		assert.Empty(r.t, missingRequests)
	}
	missingDiffs := []string{}
	for _, req := range r.currentStep().DiffRequests {
		if _, ok := r.foundsDiffs[req.id()]; !ok {
			missingDiffs = append(missingDiffs, fmt.Sprintf("missing diff request %s", req.id()))
		}
	}
	assert.Empty(r.t, missingDiffs)
	for id, payload := range r.foundsDiffs {
		r.g.Assert(r.t, r.requestFilename(id), []byte(payload))
	}
	r.currentStepIndex++
	return r.initStep()
}

func (r *requestTracker) receiveRequest(method string) func(req *http.Request) (*http.Response, error) {
	return func(req *http.Request) (*http.Response, error) {
		acceptHeader := req.Header.Get("Accept")
		version, err := ExtractVersion(acceptHeader)
		if err != nil {
			return nil, err
		}
		var payload string
		if req.Body != nil {
			payloadBytes, err := io.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}
			payload = string(payloadBytes)
		}
		text, status, err := r.matchRequest(method, req.URL.Path, version, payload)
		if err != nil {
			return nil, err
		}
		response := httpmock.NewStringResponse(status, text)
		response.Header.Set("Content-Type", fmt.Sprintf("application/vnd.atlas.%s+json;charset=utf-8", version))
		return response, nil
	}
}

func normalizePayload(payload string) (string, error) {
	var tempHolder any
	err := json.Unmarshal([]byte(payload), &tempHolder)
	if err != nil {
		return "", err
	}
	sb := strings.Builder{}
	encoder := json.NewEncoder(&sb)
	encoder.SetIndent("", " ")
	err = encoder.Encode(tempHolder)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(sb.String(), "\n"), nil
}

func (r *requestTracker) matchRequest(method, urlPath, version, payload string) (response string, statusCode int, err error) {
	step := r.currentStep()
	for _, request := range step.DiffRequests {
		if !request.Match(method, urlPath, version, r.vars) {
			continue
		}
		requestID := request.id()
		if _, ok := r.foundsDiffs[requestID]; ok {
			continue
		}
		normalizedPayload, err := normalizePayload(payload)
		if err != nil {
			return "", 0, err
		}
		r.foundsDiffs[requestID] = normalizedPayload
	}

	for _, request := range step.RequestResponses {
		if !request.Match(method, urlPath, version, r.vars) {
			continue
		}
		requestID := request.id()
		nextIndex := r.usedResponses[requestID]
		if nextIndex >= len(request.Responses) {
			if r.allowReReadGet && method == "GET" {
				nextIndex = len(request.Responses) - 1
			} else {
				continue
			}
		}
		r.usedResponses[requestID]++
		response := request.Responses[nextIndex]
		return replaceVars(response.Text, r.vars), response.Status, nil
	}
	return "", 0, fmt.Errorf("no matching request found %s %s %s", method, urlPath, version)
}
