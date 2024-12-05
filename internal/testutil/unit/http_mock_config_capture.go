package unit

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"testing"

	"gopkg.in/yaml.v3"
)

func defaultIsDiff(rt *RoundTrip) bool {
	return rt.Request.Method != "GET"
}

func newCaptureMockConfigClientModifier(t *testing.T, expectedStepCount int) *captureMockConfigClientModifier {
	t.Helper()
	// TODO: Support reading these from a file
	apiSpecPaths := map[string][]APISpecPath{}
	return &captureMockConfigClientModifier{
		t:                 t,
		apiSpecPaths:      apiSpecPaths,
		expectedStepCount: expectedStepCount,
		isDiff:            defaultIsDiff,
		capturedData:      newMockHTTPData(expectedStepCount),
	}
}

type captureMockConfigClientModifier struct {
	oldTransport      http.RoundTripper
	t                 *testing.T
	isDiff            func(*RoundTrip) bool
	apiSpecPaths      map[string][]APISpecPath
	capturedData      mockHTTPData
	expectedStepCount int
	responseIndex     int
	stepNumber        int
}

func (c *captureMockConfigClientModifier) IncreaseStepNumber() {
	c.stepNumber++
}

func (c *captureMockConfigClientModifier) ModifyHTTPClient(httpClient *http.Client) error {
	if !IsCapture() {
		return fmt.Errorf("cannot use capture modifier without %s='yes|true|1'", envNameHTTPMockerCapture)
	}
	c.oldTransport = httpClient.Transport
	httpClient.Transport = c
	return nil
}

func (c *captureMockConfigClientModifier) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := c.oldTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	c.responseIndex++
	specPaths := c.apiSpecPaths[req.Method]
	rt, parseError := parseRoundTrip(c.t, req, resp, c.responseIndex, c.stepNumber, &specPaths)
	if parseError != nil {
		c.t.Logf("error parsing round trip: %s", parseError)
		return resp, err
	}
	addError := c.capturedData.AddRT(c.t, rt, c.isDiff(rt))
	if addError != nil {
		c.t.Logf("error adding round trip: %s", addError)
	}
	return resp, err
}

func (c *captureMockConfigClientModifier) WriteCapturedData(path string) error {
	if c.stepNumber != c.expectedStepCount {
		c.t.Logf("expected %d steps, but got %d, skipping config dump", c.expectedStepCount, c.stepNumber)
	}
	// will override content if file exists
	configYaml, err := yaml.Marshal(c.capturedData)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, configYaml, 0o600)
	if err != nil {
		return err
	}
	return nil
}

var datePattern = regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)

func extractVersion(headerValueRequest, headerValueResponse string) string {
	found := datePattern.FindString(headerValueRequest)
	if found != "" {
		return found
	}
	return datePattern.FindString(headerValueResponse)
}

func parseRoundTrip(t *testing.T, req *http.Request, resp *http.Response, responseIndex, stepNumber int, specPaths *[]APISpecPath) (*RoundTrip, error) {
	t.Helper()
	version := extractVersion(req.Header.Get("Accept"), resp.Header.Get("Content-Type"))
	if version == "" {
		t.Logf("could not find version in request or response headers for responseIndex %d", responseIndex)
	}
	normalizedPath, err := findNormalizedPath(req.URL.Path, specPaths)
	if err != nil {
		return nil, err
	}
	requestPayload, err := extractPayload(req.Body)
	if err != nil {
		return nil, err
	}
	responsePayload, err := extractPayload(resp.Body)
	if err != nil {
		return nil, err
	}
	return &RoundTrip{
		Request:    parseRequestInfo(req, version, requestPayload),
		Response:   parseStatusText(resp, responsePayload, responseIndex),
		Variables:  normalizedPath.Variables(req.URL.Path),
		StepNumber: stepNumber,
	}, nil
}

func parseRequestInfo(req *http.Request, version, payload string) RequestInfo {
	return RequestInfo{
		Version: version,
		Path:    req.URL.Path,
		Method:  req.Method,
		Text:    payload,
	}
}

func parseStatusText(resp *http.Response, payload string, responseIndex int) statusText {
	return statusText{
		Text:          payload,
		Status:        resp.StatusCode,
		ResponseIndex: responseIndex,
	}
}
