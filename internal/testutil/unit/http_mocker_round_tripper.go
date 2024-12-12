package unit

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func MockRoundTripper(t *testing.T, config *MockHTTPDataConfig, data *MockHTTPData) (http.RoundTripper, *mockRoundTripper) {
	t.Helper()
	myTransport := httpmock.NewMockTransport()
	var mockTransport http.RoundTripper = myTransport
	tracker := newMockRoundTripper(t, data)
	if config != nil {
		tracker.allowMissingRequests = config.AllowMissingRequests
	}
	for _, method := range []string{"GET", "POST", "PUT", "DELETE", "PATCH"} {
		myTransport.RegisterRegexpResponder(method, regexp.MustCompile(".*"), tracker.receiveRequest(method))
	}
	return mockTransport, tracker
}
func parseTestDataConfigYAML(filePath string) (*MockHTTPData, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var testData MockHTTPData
	err = yaml.Unmarshal(data, &testData)
	if err != nil {
		return nil, err
	}
	return &testData, nil
}

func newMockRoundTripper(t *testing.T, data *MockHTTPData) *mockRoundTripper {
	t.Helper()
	return &mockRoundTripper{
		t:                t,
		g:                goldie.New(t, goldie.WithTestNameForDir(true), goldie.WithNameSuffix(".json")),
		data:             data,
		usedVars:         map[string]string{},
		logRequests:      os.Getenv("TF_LOG") == "DEBUG",
		currentStepIndex: -1, // increased on the start of the test
	}
}

type mockRoundTripper struct {
	t    *testing.T
	g    *goldie.Goldie
	data *MockHTTPData

	usedVars             map[string]string
	usedResponses        map[string]int
	foundsDiffs          map[int]string
	currentStepIndex     int
	diffResponseIndex    int
	allowMissingRequests bool
	logRequests          bool
}

func (r *mockRoundTripper) IncreaseStepNumberAndInit() {
	r.currentStepIndex++
	err := r.initStep()
	require.NoError(r.t, err)
}

func (r *mockRoundTripper) allowReUse(req *RequestInfo) bool {
	isGet := req.Method == "GET"
	customReReadOk := req.Method == "POST" && strings.HasSuffix(req.Path, ":validate")
	return isGet || customReReadOk
}

func (r *mockRoundTripper) requestFilename(requestID string, index int) string {
	return strings.ReplaceAll(fmt.Sprintf("%02d_%02d_%s", r.currentStepIndex+1, index+1, requestID), "/", "_")
}

func (r *mockRoundTripper) manualFilenameIfExist(requestID string, index int) string {
	defaultFilestem := strings.ReplaceAll(fmt.Sprintf("%02d_%02d_%s", r.currentStepIndex+1, index+1, requestID), "/", "_")
	manualFilestem := defaultFilestem + "_manual"
	if _, err := os.Stat("testdata" + "/" + r.t.Name() + "/" + manualFilestem + ".json"); err == nil {
		return manualFilestem
	}
	return defaultFilestem
}

func (r *mockRoundTripper) initStep() error {
	r.usedResponses = map[string]int{}
	r.foundsDiffs = map[int]string{}
	step := r.currentStep()
	if step == nil {
		return nil
	}
	for index, req := range step.DiffRequests {
		err := r.g.Update(r.t, r.requestFilename(req.idShort(), index), []byte(req.Text))
		if err != nil {
			return err
		}
	}
	r.nextDiffResponseIndex()
	return nil
}

func (r *mockRoundTripper) nextDiffResponseIndex() {
	step := r.currentStep()
	if step == nil {
		r.t.Fatal("no more steps, in testCase")
	}
	for index, req := range step.DiffRequests {
		if _, ok := r.foundsDiffs[index]; !ok {
			r.diffResponseIndex = req.Responses[0].ResponseIndex
			return
		}
	}
	// no more diffs in current step, any response index will do, assuming never more than 100k responses
	r.diffResponseIndex = 99999
}

func (r *mockRoundTripper) currentStep() *stepRequests {
	if r.currentStepIndex >= len(r.data.Steps) {
		return nil
	}
	return &r.data.Steps[r.currentStepIndex]
}

func (r *mockRoundTripper) CheckStepRequests(_ *terraform.State) error {
	missingRequests := []string{}
	step := r.currentStep()
	for _, req := range step.RequestResponses {
		missingRequestsCount := len(req.Responses) - r.usedResponses[req.id()]
		if missingRequestsCount > 0 {
			missingIndexes := []string{}
			for i := 0; i < missingRequestsCount; i++ {
				missingResponse := (len(req.Responses) - missingRequestsCount) + i
				missingIndexes = append(missingIndexes, fmt.Sprintf("%d", req.Responses[missingResponse].ResponseIndex))
			}
			missingIndexesStr := strings.Join(missingIndexes, ", ")
			missingRequests = append(missingRequests, fmt.Sprintf("missing %d requests of %s (%s)", missingRequestsCount, req.idShort(), missingIndexesStr))
		}
	}
	if r.allowMissingRequests {
		if len(missingRequests) > 0 {
			r.t.Logf("missing requests:\n%s", strings.Join(missingRequests, "\n"))
		}
	} else {
		assert.Empty(r.t, missingRequests)
	}
	missingDiffs := []string{}
	for i, req := range step.DiffRequests {
		if _, ok := r.foundsDiffs[i]; !ok {
			missingDiffs = append(missingDiffs, fmt.Sprintf("missing diff request %s", req.idShort()))
		}
	}
	assert.Empty(r.t, missingDiffs)
	for index, payload := range r.foundsDiffs {
		diff := step.DiffRequests[index]
		filename := r.manualFilenameIfExist(diff.idShort(), index)
		r.t.Logf("checking diff %s", filename)
		payloadWithVars := useVars(r.usedVars, payload)
		r.g.Assert(r.t, filename, []byte(payloadWithVars))
	}
	return nil
}

func (r *mockRoundTripper) receiveRequest(method string) func(req *http.Request) (*http.Response, error) {
	return func(req *http.Request) (*http.Response, error) {
		acceptHeader := req.Header.Get("Accept")
		version, err := ExtractVersion(acceptHeader)
		if err != nil {
			return nil, err
		}
		_, payload, err := extractAndNormalizePayload(req.Body)
		if r.logRequests {
			r.t.Logf("received request\n %s %s %s\n%s\n", method, req.URL.Path, version, payload)
		}
		if err != nil {
			return nil, err
		}
		text, status, err := r.matchRequest(method, version, payload, req.URL)
		if err != nil {
			return nil, err
		}
		if r.logRequests {
			r.t.Logf("responding with\n%d\n%s\n", status, text)
		}
		response := httpmock.NewStringResponse(status, text)
		response.Header.Set("Content-Type", fmt.Sprintf("application/vnd.atlas.%s+json;charset=utf-8", version))
		return response, nil
	}
}
func (r *mockRoundTripper) matchRequest(method, version, payload string, reqURL *url.URL) (response string, statusCode int, err error) {
	step := r.currentStep()
	if step == nil {
		return "", 0, fmt.Errorf("no more steps in mock data")
	}
	for index, request := range step.DiffRequests {
		if !request.Match(r.t, method, version, reqURL, r.usedVars) {
			continue
		}
		if _, ok := r.foundsDiffs[index]; ok {
			continue
		}
		r.foundsDiffs[index] = payload
		r.nextDiffResponseIndex()
		break
	}
	nextDiffResponse := r.diffResponseIndex

	for _, request := range step.RequestResponses {
		if !request.Match(r.t, method, version, reqURL, r.usedVars) {
			continue
		}
		requestID := request.id()
		nextIndex := r.usedResponses[requestID]
		if nextIndex >= len(request.Responses) {
			if r.allowReUse(&request) {
				nextIndex = len(request.Responses) - 1
			} else {
				continue
			}
		}
		response := request.Responses[nextIndex]
		// cannot return a response that is sent after a diff response
		if response.ResponseIndex > nextDiffResponse {
			prevIndex := nextIndex - 1
			if prevIndex >= 0 && r.allowReUse(&request) {
				response = request.Responses[prevIndex]
				r.t.Logf("re-reading %s request with response_index=%d as diff hasn't been returned yet (%d)", request.Method, response.ResponseIndex, nextDiffResponse)
				return replaceVars(response.Text, r.usedVars), response.Status, nil
			}
			continue
		}
		r.usedResponses[requestID]++
		return replaceVars(response.Text, r.usedVars), response.Status, nil
	}
	return "", 0, fmt.Errorf("no matching request found %s %s?%s %s", method, reqURL.Path, reqURL.RawQuery, version)
}
