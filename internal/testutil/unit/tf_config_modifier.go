package unit

import (
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/hcl"
)

type TFConfigReplacementType int

const (
	TFConfigReplacementString TFConfigReplacementType = iota
)

// Current assumption, variable name must match API Spec Path Param name
var variableAttributes = map[string]func(string, string) string{
	"name": func(resourceName string, attrName string) string {
		return shortName(resourceName) + "Name"
	},
	"org_id": func(resourceName string, attrName string) string {
		return "orgId"
	},
	"project_id": func(resourceName string, attrName string) string {
		return "groupId"
	},
}

func ExtractConfigVariables(t *testing.T, config string) map[string]string {
	t.Helper()
	if config == "" {
		return nil
	}
	vars := map[string]string{}
	parse := hcl.GetDefParser(t, config)
	for _, resource := range parse.Body().Blocks() {
		if resource.Type() != "resource" {
			continue
		}
		for name, attr := range resource.Body().Attributes() {
			varNameFunc, ok := variableAttributes[name]
			if !ok {
				continue
			}
			varName := varNameFunc(resource.Labels()[0], name)
			varValue := extractStringValue(attr.BuildTokens(nil))
			if varValue != "" {
				vars[varName] = varValue
			}
		}
	}
	return vars
}

func fullResourceName(resourceName string) string {
	if strings.HasPrefix(resourceName, "mongodbatlas_") {
		return resourceName
	}
	return "mongodbatlas_" + resourceName
}

func shortName(resourceName string) string {
	parts := strings.Split(resourceName, "_")
	return parts[len(parts)-1]
}

func extractStringValue(tokens hclwrite.Tokens) string {
	var str string
	for _, token := range tokens {
		if token.Type == hclsyntax.TokenQuotedLit {
			str = string(token.Bytes)
			break
		}
	}
	return str
}
