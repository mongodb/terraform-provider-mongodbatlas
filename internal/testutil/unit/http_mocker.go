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
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

const (
	EnvNameHTTPMockerCapture = "HTTP_MOCKER_CAPTURE"
	configFileExtension      = ".yaml"
)

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
	return slices.Contains([]string{"yes", "1", "true"}, strings.ToLower(os.Getenv(EnvNameHTTPMockerCapture)))
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
	var err error
	if IsCapture() {
		err = enableCaptureForTestCase(t, testCase)
	} else {
		err = enableMockingForTestCase(t, vars, config, testCase)
	}
	require.NoError(t, err)
	resource.ParallelTest(t, *testCase)
}

func enableMockingForTestCase(t *testing.T, vars map[string]string, config *MockHTTPDataConfig, testCase *resource.TestCase) error {
	t.Helper()
	roundTripper, nextStep, checkFunc := MockRoundTripper(t, vars, config)
	testCase.ProtoV6ProviderFactories = TestAccProviderV6FactoriesWithMock(t, &mockClientModifier{config: config, mockRoundTripper: roundTripper})
	testCase.PreCheck = func() {
		if config.SideEffect != nil {
			require.NoError(t, config.SideEffect())
		}
	}
	for i := range testCase.Steps {
		step := &testCase.Steps[i]
		oldSkip := step.SkipFunc
		step.SkipFunc = func() (bool, error) {
			nextStep()
			var shouldSkip bool
			var err error
			if oldSkip != nil {
				shouldSkip, err = oldSkip()
			}
			return shouldSkip, err
		}
		if oldCheck := step.Check; oldCheck != nil {
			newCheck := func(s *terraform.State) error {
				if err := oldCheck(s); err != nil {
					errText := err.Error()
					// TODO: Support also mocking the acc.Connv2
					if !strings.Contains(errText, "not found") {
						return err
					}
				}
				return checkFunc(s)
			}
			step.Check = newCheck
		}
	}
	return nil
}

func enableCaptureForTestCase(t *testing.T, testCase *resource.TestCase) error {
	t.Helper()
	stepCount := len(testCase.Steps)
	clientModifier := NewCaptureMockConfigClientModifier(t, stepCount)
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
	return nil
}

func MockConfigFilePath(t *testing.T) string {
	t.Helper()
	testDir := "testdata"
	return path.Join(testDir, t.Name()+configFileExtension)
}

var versionDatePattern = regexp.MustCompile(`(\d{4}-\d{2}-\d{2})`)

func ExtractVersion(contentType string) (string, error) {
	match := versionDatePattern.FindStringSubmatch(contentType)
	if len(match) > 1 {
		return match[1], nil
	}
	return "", fmt.Errorf("could not extract version from %s header", contentType)
}

func ExtractVersionRequestResponse(headerValueRequest, headerValueResponse string) string {
	found := versionDatePattern.FindString(headerValueRequest)
	if found != "" {
		return found
	}
	return versionDatePattern.FindString(headerValueResponse)
}

func extractAndNormalizePayload(body io.Reader) (originalPayload, normalizedPayload string, err error) {
	if body != nil {
		payloadBytes, err := io.ReadAll(body)
		if err != nil {
			return "", "", err
		}
		originalPayload = string(payloadBytes)
	}
	normalizedPayload, err = normalizePayload(originalPayload)
	if err != nil {
		return "", "", err
	}
	return originalPayload, normalizedPayload, nil
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
