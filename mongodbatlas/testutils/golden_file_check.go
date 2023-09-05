package testutils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var defaultIgnoredFields = []string{".primary.id", ".primary.attributes.id"}

// CheckStateMatchesGoldenFile is a TestCheckFunc (used in terraform-plugin-testing) that will compare the current terraform state againt a golden file.
// If a diff is detected, the test will fail with a detail message of the diff and a command that can be used to update the golden file.
// ingoredFields parameter can be used if certain fields within the terraform state need to be ignored (attributes that have values unique to each execution).
func CheckStateMatchesGoldenFile(resourceName, goldenFilePath string, ignoredFields []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		curStateStr, err := json.MarshalIndent(rs, "", "  ")
		if err != nil {
			return err
		}
		goldenStr, err := ioutil.ReadFile(goldenFilePath)
		if err != nil {
			return err
		}

		// Unmarshal the JSON content into maps
		var curStateMap, goldenMap map[string]interface{}
		if err := json.Unmarshal(curStateStr, &curStateMap); err != nil {
			return fmt.Errorf("failed to unmarshal current state: %v", err)
		}
		if err := json.Unmarshal(goldenStr, &goldenMap); err != nil {
			cpCmd, err := generateCpCommandWithCurState(curStateStr, goldenFilePath)
			if err != nil {
				return err
			}
			return fmt.Errorf("golden file has malformed json. To update the golden file with current state: %s", cpCmd)
		}

		optIgnoredFields := ignoreFieldsOpt(ignoredFields)

		if !cmp.Equal(curStateMap, goldenMap, optIgnoredFields) {
			diffMsg := cmp.Diff(curStateMap, goldenMap, optIgnoredFields)
			cpCmd, err := generateCpCommandWithCurState(curStateStr, goldenFilePath)
			if err != nil {
				return err
			}
			return fmt.Errorf("serialized state doesn't match the golden file.\nDifferences: %s.\nTo update the golden file, run:\n%s", diffMsg, cpCmd)
		}

		return nil
	}
}

func ignoreFieldsOpt(ignoredFields []string) cmp.Option {
	ignored := defaultIgnoredFields
	ignored = append(ignored, ignoredFields...)
	ignoreFunc := func(path cmp.Path) bool {
		completePath := ""
		for _, p := range path {
			if t, ok := p.(cmp.MapIndex); ok {
				completePath += "." + t.Key().String()
			}
		}
		for _, ignore := range ignored {
			if completePath == ignore { // TODO: this can be adjusted to support regex for ingored fields, useful for arrays
				return true
			}
		}
		return false
	}
	return cmp.FilterPath(ignoreFunc, cmp.Ignore())
}

func generateCpCommandWithCurState(state []byte, goldenFilePath string) (string, error) {
	// Writing the current state to a tmp file
	tmpDir := os.TempDir()
	tmpFilePath := filepath.Join(tmpDir, "current-state.json")
	err := ioutil.WriteFile(tmpFilePath, state, 0600) // TODO revise permission
	if err != nil {
		return "", fmt.Errorf("failed to write current state to tmp file: %v", err)
	}

	return fmt.Sprintf("cp %s ./mongodbatlas/%s", tmpFilePath, goldenFilePath), nil
}
