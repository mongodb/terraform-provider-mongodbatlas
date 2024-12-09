package unit

import (
	"errors"
	"net/http"
	"os"
	"path"
	"slices"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"
)

const (
	EnvNameHTTPMockerCapture = "HTTP_MOCKER_CAPTURE"
	configFileExtension      = ".yaml"
)

type MockHTTPDataConfig struct {
	SideEffect           func() error
	AllowMissingRequests bool
	AllowReReadGet       bool
}

func IsCapture() bool {
	return slices.Contains([]string{"yes", "1", "true"}, strings.ToLower(os.Getenv(EnvNameHTTPMockerCapture)))
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
