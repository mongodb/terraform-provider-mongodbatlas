package encryptionatrestprivateendpoint_test

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
// 	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
// )

// func TestAccEncryptionAtRestPrivateEndpoint_basic(t *testing.T) {
// 	acc.SkipTestForCI(t) // needs Azure configuration

// 	resource.ParallelTest(t, resource.TestCase{
// 		PreCheck:                 func() { acc.PreCheckBasic(t) },
// 		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
// 		//		CheckDestroy:             checkDestroyEncryptionatrestprivateendpoint,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: encryptionatrestprivateendpointConfig(),
// 				Check:  encryptionatrestprivateendpointAttributeChecks(),
// 			},
// 			{
// 				Config:            encryptionatrestprivateendpointConfig(),
// 				ResourceName:      resourceName,
// 				ImportStateIdFunc: checkEncryptionatrestprivateendpointImportStateIDFunc,
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 		},
// 	},
// 	)
// }

// func configAzureBasic(projectID string, azure *admin.AzureKeyVault, region string) string {
// 	var encryptionAtRestConfig string

// 	return fmt.Sprintf(`
// 		%s

// 		resource "mongodbatlas_encryption_at_rest_private_endpoint" "test" {
// 		    project_id = mongodbatlas_encryption_at_rest.test.project_id
// 		    cloud_provider = "AZURE"
// 		    region_name =
// 		}
// 	`, encryptionAtRestConfig, region)
// }
