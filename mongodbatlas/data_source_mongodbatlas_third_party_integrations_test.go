package mongodbatlas

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"go.mongodb.org/atlas/mongodbatlas"
)

func TestAccdataSourceMongoDBAtlasThirdPartyIntegrations_basic(t *testing.T) {
	var (
		projectID = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		intgTypes = []string{"NEW_RELIC"}
		hclConfig = make([]*thirdPartyConfig, 0, len(intgTypes))
		dsName    = "data.mongodbatlas_third_party_integrations.test"
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
					return *out
				}(),
			},
		)
	}

	intgResourcesHCL := testAccMongoDBAtlasThirdPartyIntegrationsDataSourceConfig(hclConfig, projectID)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: intgResourcesHCL,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mongodbatlas_third_party_integration.test_NEW_RELIC", "id"),
				),
			}, {

				Config: testAccMongoDBAtlasThirdPartyIntegrationsDataSourceConfigWithDS(intgResourcesHCL, projectID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsName, "project_id", projectID),
					resource.TestCheckResourceAttrSet(dsName, "project_id"),
					resource.TestCheckResourceAttr(dsName, "results.#", "1"),
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
