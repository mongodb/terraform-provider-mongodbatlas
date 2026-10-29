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

// TestConfigAzureStorageResources pins the scope of the resource group's tag handling: an automated
// leaked-resource detector can tag the resource group while a long-running test uses it, so the group ignores
// tag changes, while the other resources keep the tags the config sets under Terraform's control.
func TestConfigAzureStorageResources(t *testing.T) {
	config := acc.ConfigAzureStorageResources("blob", "test-rg", "testsa", "test-container", "test-principal")

	file, diags := hclwrite.ParseConfig([]byte(config), "main.tf", hcl.InitialPos)
	require.False(t, diags.HasErrors(), "generated config must be valid HCL: %s", diags)

	resourceGroup := findResource(t, file, "azurerm_resource_group", "blob_rg")
	assert.Equal(t, []string{"tags"}, lifecycleIgnoreChanges(t, resourceGroup), "the resource group must ignore externally added tags")

	for _, resource := range []struct{ resourceType, name string }{
		{"azurerm_storage_account", "blob_storage"},
		{"azurerm_storage_container", "blob_container"},
		{"azurerm_role_assignment", "blob_contributor"},
	} {
		block := findResource(t, file, resource.resourceType, resource.name)
		assert.Empty(t, lifecycleIgnoreChanges(t, block), "%s must not ignore tag changes", resource.resourceType)
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
