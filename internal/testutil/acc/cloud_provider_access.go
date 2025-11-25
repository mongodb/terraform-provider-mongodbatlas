package acc

import "fmt"

func ConfigSetupAzure(projectID, atlasAzureAppID, servicePrincipalID, tenantID string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_cloud_provider_access_setup" "test" {
		project_id = %[1]q
		provider_name = "AZURE"
		azure_config {
			atlas_azure_app_id = %[2]q
			service_principal_id = %[3]q
			tenant_id = %[4]q
		}
	 }
	`, projectID, atlasAzureAppID, servicePrincipalID, tenantID)
}
