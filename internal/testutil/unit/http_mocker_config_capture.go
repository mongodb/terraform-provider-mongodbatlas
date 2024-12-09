package unit

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func defaultIsDiff(rt *RoundTrip) bool {
	return rt.Request.Method != "GET" && !strings.HasSuffix(rt.Request.Path, ":validate")
}

func NewCaptureMockConfigClientModifier(t *testing.T, expectedStepCount int) *CaptureMockConfigClientModifier {
	t.Helper()
	return &CaptureMockConfigClientModifier{
		t:                 t,
		expectedStepCount: expectedStepCount,
		isDiff:            defaultIsDiff,
		capturedData:      NewMockHTTPData(expectedStepCount),
	}
}

type CaptureMockConfigClientModifier struct {
	oldTransport      http.RoundTripper
	t                 *testing.T
	isDiff            func(*RoundTrip) bool
	capturedData      MockHTTPData
	expectedStepCount int
	responseIndex     int
	stepNumber        int
}

func (c *CaptureMockConfigClientModifier) IncreaseStepNumber() {
	c.stepNumber++
}

func (c *CaptureMockConfigClientModifier) ModifyHTTPClient(httpClient *http.Client) error {
	if !IsCapture() {
		return fmt.Errorf("cannot use capture modifier without %s='yes|true|1'", EnvNameHTTPMockerCapture)
	}
	c.oldTransport = httpClient.Transport
	httpClient.Transport = c
	return nil
}

func (c *CaptureMockConfigClientModifier) RoundTrip(req *http.Request) (*http.Response, error) {
	// Capture request body to avoid it being consumed
	originalBody, normalizedBody, err := extractAndNormalizePayload(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body = io.NopCloser(strings.NewReader(originalBody))

	resp, err := c.oldTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	c.responseIndex++
	specPaths := apiSpecPaths[req.Method]
	rt, parseError := parseRoundTrip(req, resp, c.responseIndex, c.stepNumber, &specPaths, normalizedBody)
	if parseError != nil {
		c.t.Logf("error parsing round trip: %s", parseError)
		return resp, err
	}
	addError := c.capturedData.AddRoundtrip(c.t, rt, c.isDiff(rt))
	if addError != nil {
		c.t.Logf("error adding round trip: %s", addError)
	}
	return resp, err
}

func (c *CaptureMockConfigClientModifier) ConfigYaml() (string, error) {
	initialYaml := strings.Builder{}
	e := yaml.NewEncoder(&initialYaml)
	e.SetIndent(1)
	err := e.Encode(c.capturedData)
	return initialYaml.String(), err
}

func (c *CaptureMockConfigClientModifier) WriteCapturedData(filePath string) error {
	if c.stepNumber != c.expectedStepCount {
		filePath = FailedFilename(filePath)
		c.t.Logf("expected %d steps, but got %d, skipping config dump", c.expectedStepCount, c.stepNumber)
	}
	configYaml, err := c.ConfigYaml()
	if err != nil {
		return err
	}
	// will override content if file exists
	err = os.WriteFile(filePath, []byte(configYaml), 0o600)
	if err != nil {
		return err
	}
	return nil
}

func FailedFilename(filePath string) string {
	dirName := path.Dir(filePath)
	formattedTime := time.Now().Format("2006-01-02-15-04")
	stem, _ := strings.CutSuffix(path.Base(filePath), configFileExtension)
	return path.Join(dirName, fmt.Sprintf("%s_failed_%s", stem, formattedTime)) + configFileExtension
}

func parseRoundTrip(req *http.Request, resp *http.Response, responseIndex, stepNumber int, specPaths *[]APISpecPath, requestPayload string) (*RoundTrip, error) {
	version := ExtractVersionRequestResponse(req.Header.Get("Accept"), resp.Header.Get("Content-Type"))
	if version == "" {
		return nil, fmt.Errorf("could not find version in request or response headers for responseIndex %d", responseIndex)
	}
	normalizedPath, err := FindNormalizedPath(req.URL.Path, specPaths)
	if err != nil {
		return nil, err
	}
	originalResponsePayload, responsePayload, err := extractAndNormalizePayload(resp.Body)
	if err != nil {
		return nil, err
	}
	// Write back response body to support reading it again
	resp.Body = io.NopCloser(strings.NewReader(originalResponsePayload))
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
