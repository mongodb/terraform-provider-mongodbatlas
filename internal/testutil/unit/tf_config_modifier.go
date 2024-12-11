package unit

import (
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/hcl"
	"github.com/stretchr/testify/require"
)

type TFConfigReplacementType int

const (
	TFConfigReplacementString TFConfigReplacementType = iota
)

type TFConfigReplacement struct {
	ResourceName  string
	AttributeName string
	Type          TFConfigReplacementType
}

func ApplyConfigModifiers(t *testing.T, oldConfig, newConfig string, modifiers []TFConfigReplacement) string {
	t.Helper()
	if oldConfig == "" || newConfig == "" {
		return ""
	}
	for _, modifier := range modifiers {
		switch modifier.Type {
		case TFConfigReplacementString:
			newConfig = stringModifier(t, oldConfig, newConfig, modifier)
		default:
			t.Fatalf("unsupported config modifier type: %d", modifier.Type)
		}
	}
	return newConfig
}

func fullResourceName(resourceName string) string {
	if strings.HasPrefix(resourceName, "mongodbatlas_") {
		return resourceName
	}
	return "mongodbatlas_" + resourceName
}

func stringModifier(t *testing.T, oldConfig, newConfig string, modifier TFConfigReplacement) string {
	t.Helper()
	resourceName := fullResourceName(modifier.ResourceName)
	oldAttribute := findAttribute(t, oldConfig, resourceName, modifier.AttributeName)
	require.NotNil(t, oldAttribute, "attribute %s not found in old config for resource %s\n%s", modifier.AttributeName, resourceName, oldConfig)
	newAttribute := findAttribute(t, newConfig, resourceName, modifier.AttributeName)
	require.NotNil(t, newAttribute, "attribute %s not found in new config for resource %s\n%s", modifier.AttributeName, resourceName, newConfig)
	oldStatement := string(oldAttribute.BuildTokens(nil).Bytes())
	newStatement := string(newAttribute.BuildTokens(nil).Bytes())
	newConfig = strings.Replace(newConfig, newStatement, oldStatement, 1)
	return newConfig
}

func findAttribute(t *testing.T, config, resourceName, attributeName string) *hclwrite.Attribute {
	t.Helper()
	parse := hcl.GetDefParser(t, config)
	for _, resource := range parse.Body().Blocks() {
		isResource := resource.Type() == "resource"
		iResourceName := resource.Labels()[0]
		if !isResource || iResourceName != resourceName {
			continue
		}
		writeBody := resource.Body()
		return writeBody.GetAttribute(attributeName)
	}
	return nil
}
