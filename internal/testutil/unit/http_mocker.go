package unit

import (
	"errors"
	"net/http"
	"os"
	"path"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/stretchr/testify/require"
)

const (
	EnvNameHTTPMockerCapture = "HTTP_MOCKER_CAPTURE"
	EnvNameHTTPMockerReplay  = "HTTP_MOCKER_REPLAY"
	configFileExtension      = ".yaml"
)

type MockHTTPDataConfig struct {
	SideEffect           func() error
	IsDiffSkipSuffixes   []string
	IsDiffMustSubstrings []string
	QueryVars            []string
	AllowMissingRequests bool
}

func IsCapture() bool {
	return slices.Contains([]string{"yes", "1", "true"}, strings.ToLower(os.Getenv(EnvNameHTTPMockerCapture)))
}

func IsReplay() bool {
	return slices.Contains([]string{"yes", "1", "true"}, strings.ToLower(os.Getenv(EnvNameHTTPMockerReplay)))
}

func CaptureOrMockTestCaseAndRun(t *testing.T, config *MockHTTPDataConfig, testCase *resource.TestCase) {
	t.Helper()
	var err error
	noneSet := !IsCapture() && !IsReplay()
	switch {
	case noneSet:
		t.Logf("Neither %s nor %s is set, defaulting to capture mode", EnvNameHTTPMockerCapture, EnvNameHTTPMockerReplay)
		err = enableCaptureForTestCase(t, config, testCase)
	case IsCapture():
		err = enableCaptureForTestCase(t, config, testCase)
	case IsReplay():
		err = enableMockingForTestCase(t, config, testCase)
	}
	require.NoError(t, err)
	resource.ParallelTest(t, *testCase)
}

type mockClientModifier struct {
	config           *MockHTTPDataConfig
	mockRoundTripper http.RoundTripper
	oldRoundTripper  http.RoundTripper
}

func (c *mockClientModifier) ModifyHTTPClient(httpClient *http.Client) error {
	if IsCapture() {
		return errors.New("cannot capture requests when using MockTestCaseAndRun")
	}
	c.oldRoundTripper = httpClient.Transport
	httpClient.Transport = c.mockRoundTripper
	return nil
}

func (c *mockClientModifier) ResetHTTPClient(httpClient *http.Client) {
	if c.oldRoundTripper != nil {
		httpClient.Transport = c.oldRoundTripper
	}
}

func enableMockingForTestCase(t *testing.T, config *MockHTTPDataConfig, testCase *resource.TestCase) error {
	t.Helper()
	roundTripper, nextStep, checkFunc := MockRoundTripper(t, config)
	httpClientModifier := mockClientModifier{config: config, mockRoundTripper: roundTripper}
	testCase.ProtoV6ProviderFactories = TestAccProviderV6FactoriesWithMock(t, &httpClientModifier)
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
	clientModifier := NewCaptureMockConfigClientModifier(t, stepCount, config)
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

func MockConfigFilePath(t *testing.T) string {
	t.Helper()
	testDir := "testdata"
	return path.Join(testDir, t.Name()+configFileExtension)
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
