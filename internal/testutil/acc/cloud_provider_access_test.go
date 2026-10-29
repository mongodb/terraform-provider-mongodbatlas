package acc_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

// TestAzureResourceGroupIgnoresTags pins the scope of the Azure resource group tag handling: an automated
// leaked-resource detector can tag a resource group while a long-running test uses it, so each resource group
// these fixtures create ignores tag changes, while the other resources keep their tags under Terraform's control.
func TestAzureResourceGroupIgnoresTags(t *testing.T) {
	testCases := map[string]struct {
		config         string
		resourceGroup  string
		storageAccount string
	}{
		"storage resources": {
			config:         acc.ConfigAzureStorageResources("blob", "test-rg", "testsa", "test-container", "test-principal"),
			resourceGroup:  "blob_rg",
			storageAccount: "blob_storage",
		},
		"privatelink endpoint": {
			config:         acc.GetCompleteAzureBlobStorageConfig("test-project", "test-cluster", "test-subscription", "test-client", "test-secret", "test-tenant", "test-rg", "testsa"),
			resourceGroup:  "this",
			storageAccount: "this",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			file, diags := hclwrite.ParseConfig([]byte(tc.config), "main.tf", hcl.InitialPos)
			require.False(t, diags.HasErrors(), "generated config must be valid HCL: %s", diags)

			resourceGroup := findResource(t, file, "azurerm_resource_group", tc.resourceGroup)
			assert.Equal(t, []string{"tags"}, lifecycleIgnoreChanges(t, resourceGroup), "the resource group must ignore externally added tags")

			storageAccount := findResource(t, file, "azurerm_storage_account", tc.storageAccount)
			assert.Empty(t, lifecycleIgnoreChanges(t, storageAccount), "the storage account must not ignore tag changes")
		})
	}
}

func findResource(t *testing.T, file *hclwrite.File, resourceType, name string) *hclwrite.Block {
	t.Helper()
	for _, block := range file.Body().Blocks() {
		if block.Type() == "resource" && len(block.Labels()) == 2 && block.Labels()[0] == resourceType && block.Labels()[1] == name {
			return block
		}
	}
	require.Failf(t, "resource block not found", "resource %s.%s not found in generated config", resourceType, name)
	return nil
}

func lifecycleIgnoreChanges(t *testing.T, resource *hclwrite.Block) []string {
	t.Helper()
	for _, block := range resource.Body().Blocks() {
		if block.Type() != "lifecycle" {
			continue
		}
		attribute := block.Body().GetAttribute("ignore_changes")
		if attribute == nil {
			return nil
		}
		expression := strings.TrimSpace(string(attribute.Expr().BuildTokens(nil).Bytes()))
		expression = strings.TrimPrefix(strings.TrimSuffix(expression, "]"), "[")
		keys := []string{}
		for key := range strings.SplitSeq(expression, ",") {
			if key = strings.TrimSpace(key); key != "" {
				keys = append(keys, key)
			}
		}
		return keys
	}
	return nil
}
