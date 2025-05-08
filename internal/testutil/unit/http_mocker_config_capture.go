package unit

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func configureIsDiff(config *MockHTTPDataConfig) func(*RoundTrip) bool {
	return func(rt *RoundTrip) bool {
		if rt.Request.Method == "GET" {
			return false
		}
		if config == nil {
			return true
		}
		if config.IsDiffSkipSuffixes != nil {
			for _, suffix := range config.IsDiffSkipSuffixes {
				if strings.HasSuffix(rt.Request.Path, suffix) {
					return false
				}
			}
		}
		if config.IsDiffMustSubstrings != nil {
			for _, substring := range config.IsDiffMustSubstrings {
				if !strings.Contains(rt.Request.Path, substring) {
					return false
				}
			}
		}
		return true
	}
}

func configureQueryVars(config *MockHTTPDataConfig) []string {
	if config == nil {
		return nil
	}
	vars := config.QueryVars
	sort.Strings(vars)
	return vars
}

func NewCaptureMockConfigClientModifier(t *testing.T, config *MockHTTPDataConfig, data *MockHTTPData) *CaptureMockConfigClientModifier {
	t.Helper()
	return &CaptureMockConfigClientModifier{
		t:            t,
		isDiff:       configureIsDiff(config),
		queryVars:    configureQueryVars(config),
		capturedData: data,
	}
}

type CaptureMockConfigClientModifier struct {
	oldTransport  http.RoundTripper
	t             *testing.T
	isDiff        func(*RoundTrip) bool
	capturedData  *MockHTTPData
	queryVars     []string
	responseIndex int
	stepNumber    int
	mu            sync.Mutex // as requests are in parallel, there is a chance of concurrent modification while storing round trip variables
}

func (c *CaptureMockConfigClientModifier) IncreaseStepNumber() {
	c.stepNumber++
}

func (c *CaptureMockConfigClientModifier) ModifyHTTPClient(httpClient *http.Client) error {
	c.oldTransport = httpClient.Transport
	httpClient.Transport = c
	return nil
}

func (c *CaptureMockConfigClientModifier) ResetHTTPClient(httpClient *http.Client) {
	if c.oldTransport != nil {
		httpClient.Transport = c.oldTransport
	}
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
	rt, parseError := parseRoundTrip(req, resp, c.responseIndex, c.stepNumber, &specPaths, normalizedBody, c.queryVars)
	if parseError != nil {
		c.t.Logf("error parsing round trip: %s", parseError)
		return resp, err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	addError := c.capturedData.AddRoundtrip(c.t, rt, c.isDiff(rt))
	if addError != nil {
		c.t.Logf("error adding round trip: %s", addError)
	}
	return resp, err
}

func (c *CaptureMockConfigClientModifier) NormalizeCapturedData() {
	c.capturedData.Normalize()
}

func (c *CaptureMockConfigClientModifier) ConfigYaml() (string, error) {
	capturedData := c.capturedData
	return ConfigYaml(capturedData)
}

func ConfigYaml(capturedData *MockHTTPData) (string, error) {
	initialYaml := strings.Builder{}
	e := yaml.NewEncoder(&initialYaml)
	e.SetIndent(1)
	err := e.Encode(capturedData)
	return initialYaml.String(), err
}

func (c *CaptureMockConfigClientModifier) WriteCapturedData(filePath string) error {
	configYaml, err := c.ConfigYaml()
	if err != nil {
		return err
	}
	return WriteConfigYaml(filePath, configYaml)
}

func WriteConfigYaml(filePath, configYaml string) error {
	dirPath := path.Dir(filePath)
	if !FileExist(dirPath) {
		err := os.Mkdir(dirPath, 0o755)
		if err != nil {
			return err
		}
	}
	// will override content if file exists
	err := os.WriteFile(filePath, []byte(configYaml), 0o600)
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

func parseRoundTrip(req *http.Request, resp *http.Response, responseIndex, stepNumber int, specPaths *[]APISpecPath, requestPayload string, queryVars []string) (*RoundTrip, error) {
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
		QueryString: relevantQuery(queryVars, req.URL.Query()),
		Variables:   normalizedPath.Variables(req.URL.Path),
		StepNumber:  stepNumber,
		Request: RequestInfo{
			Version: version,
			Path:    removeQueryParamsAndTrim(req.URL.Path),
			Method:  req.Method,
			Text:    requestPayload,
		},
		Response: StatusText{
			Text:          responsePayload,
			Status:        resp.StatusCode,
			ResponseIndex: responseIndex,
		},
	}, nil
}
