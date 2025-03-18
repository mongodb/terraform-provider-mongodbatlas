package unit

import (
	"net/http"
	"os"
	"path"
	"strconv"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/hcl"
	"github.com/stretchr/testify/require"
)

const (
	EnvNameHTTPMockerCapture    = "HTTP_MOCKER_CAPTURE"
	EnvNameHTTPMockerReplay     = "HTTP_MOCKER_REPLAY"
	EnvNameHTTPMockerDataUpdate = "HTTP_MOCKER_DATA_UPDATE"
	configFileExtension         = ".yaml"
)

type MockHTTPDataConfig struct {
	SideEffect           func() error
	IsDiffSkipSuffixes   []string
	IsDiffMustSubstrings []string
	QueryVars            []string
	AllowMissingRequests bool
	AllowOutOfOrder      bool
	RequestHandler       ManualRequestHandler
}

func (c MockHTTPDataConfig) WithAllowOutOfOrder() MockHTTPDataConfig { //nolint: gocritic // Want each test run to have its own config (hugeParam: c is heavy (112 bytes); consider passing it by pointer)
	c.AllowOutOfOrder = true
	return c
}

func IsCapture() bool {
	val, _ := strconv.ParseBool(os.Getenv(EnvNameHTTPMockerCapture))
	return val
}

func IsReplay() bool {
	val, _ := strconv.ParseBool(os.Getenv(EnvNameHTTPMockerReplay))
	return val
}

func IsDataUpdate() bool {
	val, _ := strconv.ParseBool(os.Getenv(EnvNameHTTPMockerDataUpdate))
	return val
}

func CaptureOrMockTestCaseAndRun(t *testing.T, config MockHTTPDataConfig, testCase *resource.TestCase) { //nolint: gocritic // Want each test run to have its own config (hugeParam: config is heavy (112 bytes); consider passing it by pointer)
	t.Helper()
	var err error
	noneSet := !IsCapture() && !IsReplay()
	bothSet := IsCapture() && IsReplay()
	switch {
	case bothSet:
		t.Fatalf("Both %s and %s are set, only one of them should be set", EnvNameHTTPMockerCapture, EnvNameHTTPMockerReplay)
	case noneSet:
		t.Logf("Neither %s nor %s is set, running test case without modifications", EnvNameHTTPMockerCapture, EnvNameHTTPMockerReplay)
	case IsReplay():
		err = enableReplayForTestCase(t, &config, testCase)
	case IsCapture():
		err = enableCaptureForTestCase(t, &config, testCase)
	}
	require.NoError(t, err)
	resource.ParallelTest(t, *testCase)
}

func IsTfLogDebug() bool {
	return os.Getenv("TF_LOG") == "DEBUG"
}

type mockClientModifier struct {
	config           *MockHTTPDataConfig
	mockRoundTripper http.RoundTripper
	oldRoundTripper  http.RoundTripper
}

func (c *mockClientModifier) ModifyHTTPClient(httpClient *http.Client) error {
	c.oldRoundTripper = httpClient.Transport
	httpClient.Transport = c.mockRoundTripper
	return nil
}

func (c *mockClientModifier) ResetHTTPClient(httpClient *http.Client) {
	if c.oldRoundTripper != nil {
		httpClient.Transport = c.oldRoundTripper
	}
}

func MockConfigFilePath(t *testing.T) string {
	t.Helper()
	testDir := "testdata"
	return path.Join(testDir, t.Name()+configFileExtension)
}

func ReadMockData(t *testing.T, tfConfigs []string) *MockHTTPData {
	t.Helper()
	httpDataPath := MockConfigFilePath(t)
	data, err := ParseTestDataConfigYAML(httpDataPath)
	require.NoError(t, err)
	oldVariables := data.Variables
	data.Variables = map[string]string{}
	data.useTFConfigs(t, tfConfigs)
	newVariables := data.Variables
	for key, value := range oldVariables {
		if _, ok := newVariables[key]; !ok {
			t.Logf("Variable %s not found from TF Config, will use variable from the mock data with value %s", key, value)
			data.Variables[key] = value
		}
	}
	for key, value := range newVariables {
		if _, ok := oldVariables[key]; !ok {
			t.Logf("Variable %s=%s not found in Mock Data, has the TF Config updated?", key, value)
		}
	}
	return data
}

func UpdateMockDataDiffRequest(t *testing.T, stepIndex, diffRequestIndex int, newText string) {
	t.Helper()
	httpDataPath := MockConfigFilePath(t)
	data, err := ParseTestDataConfigYAML(httpDataPath)
	require.NoError(t, err)
	data.Steps[stepIndex].DiffRequests[diffRequestIndex].Text = newText
	configYaml, err := ConfigYaml(data)
	require.NoError(t, err)
	err = WriteConfigYaml(httpDataPath, configYaml)
	require.NoError(t, err)
}

func enableReplayForTestCase(t *testing.T, config *MockHTTPDataConfig, testCase *resource.TestCase) error {
	t.Helper()
	tfConfigs := extractAndNormalizeConfig(t, testCase)
	data := ReadMockData(t, tfConfigs)
	roundTripper, mockRoundTripper := NewMockRoundTripper(t, config, data)
	httpClientModifier := mockClientModifier{config: config, mockRoundTripper: roundTripper}
	testCase.ProtoV6ProviderFactories = TestAccProviderV6FactoriesWithMock(t, &httpClientModifier)
	testCase.PreCheck = func() {
		if config.SideEffect != nil {
			require.NoError(t, config.SideEffect())
		}
	}
	require.Equal(t, len(testCase.Steps), len(data.Steps), "Number of steps in test case and mock data should match")
	checkFunc := mockRoundTripper.CheckStepRequests
	for i := range testCase.Steps {
		step := &testCase.Steps[i]
		oldSkip := step.SkipFunc
		step.SkipFunc = func() (bool, error) {
			mockRoundTripper.IncreaseStepNumberAndInit()
			logConfig(t, tfConfigs, i)
			var shouldSkip bool
			var err error
			if oldSkip != nil {
				shouldSkip, err = oldSkip()
			}
			return shouldSkip, err
		}
		if i == len(testCase.Steps)-1 {
			// Last check done in checkDestroy to support checking DELETE calls
			step.Check = wrapClientDuringCheck(step.Check, &httpClientModifier)
		} else {
			step.Check = wrapClientDuringCheck(step.Check, &httpClientModifier, checkFunc)
		}
	}
	testCase.CheckDestroy = wrapClientDuringCheck(testCase.CheckDestroy, &httpClientModifier, checkFunc)
	return nil
}

func enableCaptureForTestCase(t *testing.T, config *MockHTTPDataConfig, testCase *resource.TestCase) error {
	t.Helper()
	stepCount := len(testCase.Steps)
	tfConfigs := extractAndNormalizeConfig(t, testCase)
	capturedData := NewMockHTTPData(t, stepCount, tfConfigs)
	clientModifier := NewCaptureMockConfigClientModifier(t, config, capturedData)
	testCase.ProtoV6ProviderFactories = TestAccProviderV6FactoriesWithMock(t, clientModifier)
	for i := range stepCount {
		step := &testCase.Steps[i]
		oldSkip := step.SkipFunc
		step.SkipFunc = func() (bool, error) {
			clientModifier.IncreaseStepNumber()
			logConfig(t, tfConfigs, i)
			var shouldSkip bool
			var err error
			if oldSkip != nil {
				shouldSkip, err = oldSkip()
			}
			return shouldSkip, err
		}
		step.Check = wrapClientDuringCheck(step.Check, clientModifier)
	}

	writeCapturedData := func() {
		clientModifier.NormalizeCapturedData()
		filePath := MockConfigFilePath(t)
		if t.Failed() {
			filePath = FailedFilename(filePath)
		}
		err := clientModifier.WriteCapturedData(filePath)
		require.NoError(t, err)
	}
	t.Cleanup(writeCapturedData)
	testCase.CheckDestroy = wrapClientDuringCheck(testCase.CheckDestroy, clientModifier)
	return nil
}

func logConfig(t *testing.T, tfConfigs []string, i int) {
	t.Helper()
	if IsTfLogDebug() && tfConfigs[i] != "" {
		t.Logf("Step %d:\n%s\n", i+1, tfConfigs[i])
	}
}

func extractAndNormalizeConfig(t *testing.T, testCase *resource.TestCase) []string {
	t.Helper()
	stepCount := len(testCase.Steps)
	tfConfigs := make([]string, stepCount)
	for i := range testCase.Steps {
		tfConfigs[i] = hcl.PrettyHCL(t, testCase.Steps[i].Config)
	}
	return tfConfigs
}

var accClientLock = &sync.Mutex{}

func wrapClientDuringCheck(oldCheck resource.TestCheckFunc, clientModifier HTTPClientModifier, extraChecks ...resource.TestCheckFunc) resource.TestCheckFunc {
	if oldCheck == nil && len(extraChecks) == 0 {
		return nil
	}
	return func(s *terraform.State) error {
		accClientLock.Lock()
		accClient := acc.ConnV2().GetConfig().HTTPClient
		modifyErr := clientModifier.ModifyHTTPClient(accClient)
		defer func() {
			clientModifier.ResetHTTPClient(accClient)
			accClientLock.Unlock()
		}()
		if modifyErr != nil {
			return modifyErr
		}
		if oldCheck != nil {
			if err := oldCheck(s); err != nil {
				return err
			}
			for _, check := range extraChecks {
				if err := check(s); err != nil {
					return err
				}
			}
		}
		return nil
	}
}
