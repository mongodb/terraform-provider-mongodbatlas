package unit

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"
)

const resourceID = "mongodbatlas_example.this"

func MockPlanChecksAndRun(t *testing.T, mockConfig MockHTTPDataConfig, importInput, importConfig, importResourceName string, testStep *resource.TestStep) {
	t.Helper()
	exampleResourceConfig := createConfig(importInput)
	testCase := &resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: exampleResourceConfig,
				Check:  resource.TestCheckResourceAttr("mongodbatlas_example.this", "import_id", importInput),
			},
			{
				Config:             exampleResourceConfig + importConfig,
				ResourceName:       importResourceName,
				ImportStateIdFunc:  importStateImportID(),
				ImportState:        true,
				ImportStatePersist: true, // save the state to use it in the next plan
			},
			*testStep,
		},
	}
	err := enableReplayForTestCase(
		t,
		&mockConfig,
		testCase,
	)
	require.NoError(t, err)
	resource.ParallelTest(t, *testCase)
}

func createConfig(importID string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_example" "this" {
			import_id = %[1]q
		}
		`, importID)
}
func importStateImportID() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceID]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return rs.Primary.Attributes["import_id"], nil
	}
}
