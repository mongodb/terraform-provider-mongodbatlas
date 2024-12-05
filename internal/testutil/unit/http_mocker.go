package unit

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

const envNameHTTPMockerCapture = "HTTP_MOCKER_CAPTURE"

func replaceVars(text string, vars map[string]string) string {
	for key, value := range vars {
		text = strings.ReplaceAll(text, fmt.Sprintf("{%s}", key), value)
	}
	return text
}

type MockHTTPDataConfig struct {
	SideEffect           func() error
	AllowMissingRequests bool
	AllowReReadGet       bool
}

func IsCapture() bool {
	return slices.Contains([]string{"yes", "1", "true"}, strings.ToLower(os.Getenv(envNameHTTPMockerCapture)))
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

type mockClientModifier struct {
	config           *MockHTTPDataConfig
	mockRoundTripper http.RoundTripper
}

func (c *mockClientModifier) ModifyHTTPClient(httpClient *http.Client) error {
	if IsCapture() {
		return errors.New("cannot capture requests when using MockTestCaseAndRun")
	}
	httpClient.Transport = c.mockRoundTripper
	return nil
}

func MockTestCaseAndRun(t *testing.T, vars map[string]string, config *MockHTTPDataConfig, testCase *resource.TestCase) {
	t.Helper()
	if IsCapture() {
		stepCount := len(testCase.Steps)
		clientModifier := newCaptureMockConfigClientModifier(t, stepCount)
		testCase.ProtoV6ProviderFactories = TestAccProviderV6FactoriesWithMock(t, clientModifier)
		for i := range stepCount {
			step := &testCase.Steps[i]
			oldSkip := step.SkipFunc
			step.SkipFunc = func() (bool, error) {
				clientModifier.IncreaseStepNumber()
				var shouldSkip bool
				var err error
				if oldSkip != nil {
					shouldSkip, err = oldSkip()
				}
				return shouldSkip, err
			}
		}
		oldCheckDestroy := testCase.CheckDestroy
		newCheckDestroy := func(s *terraform.State) error {
			if oldCheckDestroy != nil {
				if err := oldCheckDestroy(s); err != nil {
					return err
				}
			}
			return clientModifier.WriteCapturedData(MockConfigFilePath(t))
		}
		testCase.CheckDestroy = newCheckDestroy
		resource.ParallelTest(t, *testCase)
		return
	}
	roundTripper, checkFunc := MockRoundTripper(t, vars, config)
	testCase.ProtoV6ProviderFactories = TestAccProviderV6FactoriesWithMock(t, &mockClientModifier{config: config, mockRoundTripper: roundTripper})
	testCase.PreCheck = func() {
		if config.SideEffect != nil {
			require.NoError(t, config.SideEffect())
		}
	}
	stepCount := len(testCase.Steps)
	for i := range stepCount - 1 {
		step := &testCase.Steps[i]
		if oldCheck := step.Check; oldCheck != nil {
			step.Check = resource.ComposeAggregateTestCheckFunc(oldCheck, checkFunc)
		}
	}
	// Using CheckDestroy for the final step assertions to allow mocked responses in cleanup
	oldCheckDestroy := testCase.CheckDestroy
	newCheckDestroy := func(s *terraform.State) error {
		if oldCheckDestroy != nil {
			if err := oldCheckDestroy(s); err != nil {
				return err
			}
		}
		return checkFunc(s)
	}
	testCase.CheckDestroy = newCheckDestroy
	resource.ParallelTest(t, *testCase)
}

func MockConfigFilePath(t *testing.T) string {
	t.Helper()
	testDir := "testdata"
	return path.Join(testDir, t.Name()+".yaml")
}

func MockRoundTripper(t *testing.T, vars map[string]string, config *MockHTTPDataConfig) (http.RoundTripper, resource.TestCheckFunc) {
	t.Helper()
	httpDataPath := MockConfigFilePath(t)
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

var versionDatePattern = regexp.MustCompile(`(\d{4}-\d{2}-\d{2})`)

func ExtractVersion(contentType string) (string, error) {
	match := versionDatePattern.FindStringSubmatch(contentType)
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
	foundsDiffs          map[int]string
	currentStepIndex     int
	diffResponseIndex    int
	allowMissingRequests bool
	allowReReadGet       bool
}

func (r *requestTracker) allowReUse(method string) bool {
	return r.allowReReadGet && method == "GET"
}

func (r *requestTracker) requestFilename(requestID string, index int) string {
	return strings.ReplaceAll(fmt.Sprintf("%02d_%02d_%s", r.currentStepIndex+1, index+1, requestID), "/", "_")
}

func (r *requestTracker) manualFilenameIfExist(requestID string, index int) string {
	defaultFilestem := strings.ReplaceAll(fmt.Sprintf("%02d_%02d_%s", r.currentStepIndex+1, index+1, requestID), "/", "_")
	manualFilestem := defaultFilestem + "_manual"
	if _, err := os.Stat("testdata" + "/" + r.t.Name() + "/" + manualFilestem + ".json"); err == nil {
		return manualFilestem
	}
	return defaultFilestem
}

func (r *requestTracker) initStep() error {
	require.Len(r.t, r.data.Steps, r.data.StepCount, "step count didn't match steps")
	usedKeys := strings.Join(acc.SortStringMapKeys(r.vars), ", ")
	expectedKeys := strings.Join(acc.SortStringMapKeys(r.data.Variables), ", ")
	require.Equal(r.t, expectedKeys, usedKeys, "mock variables didn't match mock data variables")
	r.usedResponses = map[string]int{}
	r.foundsDiffs = map[int]string{}
	step := r.currentStep()
	if step == nil {
		return nil
	}
	for index, req := range step.DiffRequests {
		err := r.g.Update(r.t, r.requestFilename(req.idShort(), index), []byte(replaceVars(req.Text, r.vars)))
		if err != nil {
			return err
		}
	}
	r.nextDiffResponseIndex()
	return nil
}

func (r *requestTracker) nextDiffResponseIndex() {
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

func (r *requestTracker) currentStep() *stepRequests {
	if r.currentStepIndex >= r.data.StepCount {
		return nil
	}
	return &r.data.Steps[r.currentStepIndex]
}

func (r *requestTracker) checkStepRequests(_ *terraform.State) error {
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
		r.g.Assert(r.t, filename, []byte(payload))
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
		payload, err := extractPayload(req.Body)
		if err != nil {
			return nil, err
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

func extractPayload(body io.Reader) (string, error) {
	var payload string
	if body != nil {
		payloadBytes, err := io.ReadAll(body)
		if err != nil {
			return "", err
		}
		payload = string(payloadBytes)
	}
	payload, err := normalizePayload(payload)
	if err != nil {
		return "", err
	}
	return payload, nil
}

func normalizePayload(payload string) (string, error) {
	if payload == "" {
		return "", nil
	}
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
	if step == nil {
		return "", 0, fmt.Errorf("no more steps in mock data")
	}
	for index, request := range step.DiffRequests {
		if !request.Match(method, urlPath, version, r.vars) {
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
		response := request.Responses[nextIndex]
		// cannot return a response that is sent after a diff response
		if response.ResponseIndex > nextDiffResponse {
			prevIndex := nextIndex - 1
			if prevIndex >= 0 && r.allowReUse(method) {
				response = request.Responses[prevIndex]
				r.t.Logf("re-reading GET request with response_index=%d as diff hasn't been returned yet (%d)", response.ResponseIndex, nextDiffResponse)
				return replaceVars(response.Text, r.vars), response.Status, nil
			}
			continue
		}
		r.usedResponses[requestID]++
		return replaceVars(response.Text, r.vars), response.Status, nil
	}
	return "", 0, fmt.Errorf("no matching request found %s %s %s", method, urlPath, version)
}
