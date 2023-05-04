package mongodbatlas

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.mongodb.org/atlas/mongodbatlas"
)

func TestAccConfigDSThirdPartyIntegrations_basic(t *testing.T) {
	SkipTest(t)
	var (
		projectID       = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		intgTypes       = []string{"NEW_RELIC", "OPS_GENIE", "DATADOG", "VICTOR_OPS", "WEBHOOK", "PROMETHEUS"}
		hclConfig       = make([]*thirdPartyConfig, 0, len(intgTypes))
		dsName          = "data.mongodbatlas_third_party_integrations.test"
		integrationType = ""
	)

	for _, intg := range intgTypes {
		hclConfig = append(
			hclConfig,
			&thirdPartyConfig{
				Name:      "test_" + intg,
				ProjectID: projectID,
				Integration: func() mongodbatlas.ThirdPartyIntegration {
					out := testAccCreateThirdPartyIntegrationConfig()
					out.Type = intg
					integrationType = "test_" + intg
					out.APIKey = testOpsGenieOrDatadog(intg)
					return *out
				}(),
			},
		)
	}

	intgResourcesHCL := testAccMongoDBAtlasThirdPartyIntegrationsDataSourceConfig(hclConfig, projectID)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: intgResourcesHCL,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(fmt.Sprintf("mongodbatlas_third_party_integration.%s", integrationType), "id"),
				),
			}, {

				Config: testAccMongoDBAtlasThirdPartyIntegrationsDataSourceConfigWithDS(intgResourcesHCL, projectID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsName, "project_id", projectID),
					resource.TestCheckResourceAttrSet(dsName, "project_id"),
					resource.TestCheckResourceAttrSet(dsName, "results.#"),
					resource.TestCheckResourceAttrSet(dsName, "results.0.type"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasThirdPartyIntegrationsDataSourceConfig(configs []*thirdPartyConfig, projectID string) string {
	var factory strings.Builder

	for _, cfg := range configs {
		hclEntry := testAccMongoDBAtlasThirdPartyIntegrationResourceConfig(cfg)
		factory.WriteString(hclEntry + "\n")
	}

	return factory.String()
}

func testAccMongoDBAtlasThirdPartyIntegrationsDataSourceConfigWithDS(resources, projectID string) string {
	return fmt.Sprintf(`
	%s
	data "mongodbatlas_third_party_integrations" "test" {
		project_id = "%s"
	}
`, resources, projectID)
}

func testOpsGenieOrDatadog(intTtype string) string {
	if intTtype == "DATADOG" {
		return os.Getenv("DD_API_KEY")
	}
	if intTtype == "OPS_GENIE" {
		return os.Getenv("OPS_GENIE_API_KEY")
	}
	if intTtype == "VICTOR_OPS" {
		return os.Getenv("VICTOR_OPS_API_KEY")
	}
	return ""
}
