package unit

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func NewMockRoundTripper(t *testing.T, config *MockHTTPDataConfig, data *MockHTTPData) (http.RoundTripper, *MockRoundTripper) {
	t.Helper()
	myTransport := httpmock.NewMockTransport()
	var mockTransport http.RoundTripper = myTransport
	tracker := newMockRoundTripper(t, data)
	if config != nil {
		tracker.allowMissingRequests = config.AllowMissingRequests
		tracker.allowOutOfOrder = config.AllowOutOfOrder
	}
	for _, method := range []string{"GET", "POST", "PUT", "DELETE", "PATCH"} {
		myTransport.RegisterRegexpResponder(method, regexp.MustCompile(".*"), tracker.receiveRequest(method))
	}
	return mockTransport, tracker
}
func ParseTestDataConfigYAML(filePath string) (*MockHTTPData, error) {
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

func newMockRoundTripper(t *testing.T, data *MockHTTPData) *MockRoundTripper {
	t.Helper()
	return &MockRoundTripper{
		t:                t,
		g:                goldie.New(t, goldie.WithTestNameForDir(true), goldie.WithNameSuffix(".json")),
		data:             data,
		logRequests:      IsTfLogDebug(),
		currentStepIndex: -1, // increased on the start of the test
	}
}

type MockRoundTripper struct {
	t                    *testing.T
	g                    *goldie.Goldie
	data                 *MockHTTPData
	usedResponses        map[string]int
	foundsDiffs          map[int]string
	currentStepIndex     int
	diffResponseIndex    int
	reReadCounter        int
	mu                   sync.Mutex // as requests are in parallel, there is a chance of concurrent modification while reading/updating variables
	allowMissingRequests bool
	allowOutOfOrder      bool
	logRequests          bool
}

func (r *MockRoundTripper) IncreaseStepNumberAndInit() {
	r.currentStepIndex++
	err := r.initStep()
	require.NoError(r.t, err)
}

func (r *MockRoundTripper) canReturnResponse(responseIndex int) bool {
	isAfter := responseIndex > r.diffResponseIndex
	if r.allowOutOfOrder && isAfter {
		r.t.Logf("allowwingOutOfOrder: response_index=%d is after nextDiffResponse=%d", responseIndex, r.diffResponseIndex)
	}
	return r.allowOutOfOrder || !isAfter
}

func (r *MockRoundTripper) allowReUse(req *RequestInfo) bool {
	isGet := req.Method == "GET"
	customReReadOk := req.Method == "POST" && strings.HasSuffix(req.Path, ":validate")
	return isGet || customReReadOk
}

func (r *MockRoundTripper) requestFilename(requestID string, index int) string {
	return strings.ReplaceAll(fmt.Sprintf("%02d_%02d_%s", r.currentStepIndex+1, index+1, requestID), "/", "_")
}

func (r *MockRoundTripper) manualFilenameIfExist(requestID string, index int) string {
	defaultFilestem := strings.ReplaceAll(fmt.Sprintf("%02d_%02d_%s", r.currentStepIndex+1, index+1, requestID), "/", "_")
	manualFilestem := defaultFilestem + "_manual"
	if _, err := os.Stat("testdata" + "/" + r.t.Name() + "/" + manualFilestem + ".json"); err == nil {
		return manualFilestem
	}
	return defaultFilestem
}

func (r *MockRoundTripper) initStep() error {
	r.usedResponses = map[string]int{}
	r.foundsDiffs = map[int]string{}
	r.reReadCounter = 0
	step := r.currentStep()
	if step == nil {
		return nil
	}
	for index, req := range step.DiffRequests {
		err := r.g.Update(r.t, r.requestFilename(req.IdShort(), index), []byte(req.Text))
		if err != nil {
			return err
		}
	}
	r.nextDiffResponseIndex()
	return nil
}

func (r *MockRoundTripper) nextDiffResponseIndex() {
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

func (r *MockRoundTripper) currentStep() *stepRequests {
	if r.currentStepIndex >= len(r.data.Steps) {
		return nil
	}
	return &r.data.Steps[r.currentStepIndex]
}

func (r *MockRoundTripper) CheckStepRequests(_ *terraform.State) error {
	missingRequests := []string{}
	step := r.currentStep()
	for _, req := range step.RequestResponses {
		missingRequestsCount := len(req.Responses) - r.usedResponses[req.id()]
		if missingRequestsCount > 0 {
			missingIndexes := []string{}
			for i := range missingRequestsCount {
				missingResponse := (len(req.Responses) - missingRequestsCount) + i
				missingIndexes = append(missingIndexes, fmt.Sprintf("%d", req.Responses[missingResponse].ResponseIndex))
			}
			missingIndexesStr := strings.Join(missingIndexes, ", ")
			missingRequests = append(missingRequests, fmt.Sprintf("missing %d requests of %s (%s)", missingRequestsCount, req.IdShort(), missingIndexesStr))
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
			missingDiffs = append(missingDiffs, fmt.Sprintf("missing diff request %s", req.IdShort()))
		}
	}
	assert.Empty(r.t, missingDiffs)
	for index, payload := range r.foundsDiffs {
		diff := step.DiffRequests[index]
		filename := r.manualFilenameIfExist(diff.IdShort(), index)
		r.t.Logf("checking diff %s", filename)
		payloadWithVars := useVars(r.data.Variables, payload)
		r.g.Assert(r.t, filename, []byte(payloadWithVars))
		if IsDataUpdate() {
			r.t.Logf("updating diff %s", filename)
			UpdateMockDataDiffRequest(r.t, r.currentStepIndex, index, payloadWithVars)
		}
	}
	return nil
}

func (r *MockRoundTripper) receiveRequest(method string) func(req *http.Request) (*http.Response, error) {
	return func(req *http.Request) (*http.Response, error) {
		r.mu.Lock()
		defer r.mu.Unlock()
		acceptHeader := req.Header.Get("Accept")
		version, err := ExtractVersion(acceptHeader)
		if err != nil {
			return nil, err
		}
		_, payload, err := extractAndNormalizePayload(req.Body)
		if r.logRequests {
			r.t.Logf("received request\n %s %s?%s %s\n%s\n", method, req.URL.Path, req.URL.RawQuery, version, payload)
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
func (r *MockRoundTripper) matchRequest(method, version, payload string, reqURL *url.URL) (response string, statusCode int, err error) {
	step := r.currentStep()
	if step == nil {
		return "", 0, fmt.Errorf("no more steps in mock data")
	}
	isDiff := false
	for index, request := range step.DiffRequests {
		if !request.Match(r.t, method, version, reqURL, r.data) {
			continue
		}
		if _, ok := r.foundsDiffs[index]; ok {
			continue
		}
		r.foundsDiffs[index] = payload
		r.nextDiffResponseIndex()
		isDiff = true
		break
	}
	nextDiffResponse := r.diffResponseIndex

	for _, request := range step.RequestResponses {
		if !request.Match(r.t, method, version, reqURL, r.data) {
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
		// cannot return a response that is sent after a diff response, unless it is a diff or we ignore order with allowOutOfOrder
		if !isDiff && !r.canReturnResponse(response.ResponseIndex) {
			prevIndex := nextIndex - 1
			if prevIndex >= 0 && r.allowReUse(&request) {
				r.reReadCounter++
				if r.reReadCounter > 20 {
					return "", 0, fmt.Errorf("stuck in a loop trying to re-read the same request: %s %s %s", method, version, reqURL.Path)
				}
				response = request.Responses[prevIndex]
				r.t.Logf("re-reading %s request with response_index=%d as diff hasn't been returned yet (%d)", request.Method, response.ResponseIndex, nextDiffResponse)
				return replaceVars(response.Text, r.data.Variables), response.Status, nil
			}
			continue
		}
		r.usedResponses[requestID]++
		return replaceVars(response.Text, r.data.Variables), response.Status, nil
	}
	return "", 0, fmt.Errorf("no matching request found %s %s\n%s\nnextDiffResponse=%d", method, version, reqURL.Path, nextDiffResponse)
}
