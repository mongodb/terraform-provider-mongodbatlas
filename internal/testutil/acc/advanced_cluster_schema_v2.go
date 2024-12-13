package acc

import (
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/hcl"
	"github.com/zclconf/go-cty/cty"

	"github.com/stretchr/testify/assert"
)

func ConvertToTPFAttrsMap(attrsMap map[string]string) map[string]string {
	if !config.AdvancedClusterV2Schema() {
		return attrsMap
	}
	ret := make(map[string]string, len(attrsMap))
	for name, value := range attrsMap {
		ret[attrNameToSchemaV2(name)] = value
	}
	return ret
}

func ConvertToTPFAttrsSet(attrsSet []string) []string {
	if !config.AdvancedClusterV2Schema() {
		return attrsSet
	}
	ret := make([]string, 0, len(attrsSet))
	for _, name := range attrsSet {
		ret = append(ret, attrNameToSchemaV2(name))
	}
	return ret
}

var tpfSingleNestedAttrs = []string{
	"analytics_specs",
	"electable_specs",
	"read_only_specs",
	"auto_scaling", // includes analytics_auto_scaling
	"advanced_configuration",
	"bi_connector_config",
}

func attrNameToSchemaV2(name string) string {
	for _, singleAttrName := range tpfSingleNestedAttrs {
		name = strings.ReplaceAll(name, singleAttrName+".0", singleAttrName)
	}
	return name
}

func ConvertAdvancedClusterToTPF(t *testing.T, def string) string {
	t.Helper()
	if !config.AdvancedClusterV2Schema() {
		return def
	}
	parse := hcl.GetDefParser(t, def)
	for _, resource := range parse.Body().Blocks() {
		isResource := resource.Type() == "resource"
		resourceName := resource.Labels()[0]
		if !isResource || resourceName != "mongodbatlas_advanced_cluster" {
			continue
		}
		writeBody := resource.Body()
		convertAttrs(t, "labels", writeBody, true, hcl.GetAttrVal)
		convertAttrs(t, "tags", writeBody, true, hcl.GetAttrVal)
		convertAttrs(t, "replication_specs", writeBody, true, getReplicationSpecs)
		convertAttrs(t, "advanced_configuration", writeBody, false, hcl.GetAttrVal)
		convertAttrs(t, "bi_connector_config", writeBody, false, hcl.GetAttrVal)
	}
	content := parse.Bytes()
	return string(content)
}

func ConvertAdvancedClusterToTPFIfEnabled(t *testing.T, enabled bool, def string) string {
	t.Helper()
	if enabled {
		return ConvertAdvancedClusterToTPF(t, def)
	}
	return def
}

func AssertEqualHCL(t *testing.T, expected, actual string, msgAndArgs ...interface{}) {
	t.Helper()
	assert.Equal(t, hcl.CanonicalHCL(t, expected), hcl.CanonicalHCL(t, actual), msgAndArgs...)
}

func convertAttrs(t *testing.T, name string, writeBody *hclwrite.Body, isList bool, getOneAttr func(*testing.T, *hclsyntax.Body) cty.Value) {
	t.Helper()
	var vals []cty.Value
	for {
		match := writeBody.FirstMatchingBlock(name, nil)
		if match == nil {
			break
		}
		vals = append(vals, getOneAttr(t, hcl.GetBlockBody(t, match)))
		writeBody.RemoveBlock(match) // TODO: RemoveBlock doesn't remove newline just after the block so an extra line is added
	}
	if len(vals) == 0 {
		return
	}
	if isList {
		writeBody.SetAttributeValue(name, cty.TupleVal(vals))
	} else {
		assert.Len(t, vals, 1, "can be only one of %s", name)
		writeBody.SetAttributeValue(name, vals[0])
	}
}

func getReplicationSpecs(t *testing.T, body *hclsyntax.Body) cty.Value {
	t.Helper()
	const name = "region_configs"
	var vals []cty.Value
	for _, block := range body.Blocks {
		assert.Equal(t, name, block.Type, "unexpected block type: %s", block.Type)
		vals = append(vals, hcl.GetAttrVal(t, block.Body))
	}
	return cty.ObjectVal(map[string]cty.Value{
		name: cty.TupleVal(vals),
	})
}
