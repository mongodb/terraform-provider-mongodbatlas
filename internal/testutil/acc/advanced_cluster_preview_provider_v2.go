package acc

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/hcl"
	"github.com/zclconf/go-cty/cty"

	"github.com/stretchr/testify/assert"
)

// IsTestSDKv2ToTPF returns if we want to run migration tests from SDKv2 to TPF.
func IsTestSDKv2ToTPF() bool {
	env, _ := strconv.ParseBool(os.Getenv("MONGODB_ATLAS_TEST_SDKV2_TO_TPF"))
	return env
}

func CheckRSAndDSPreviewProviderV2(usePreviewProvider bool, resourceName string, dataSourceName, pluralDataSourceName *string, attrsSet []string, attrsMap map[string]string, extra ...resource.TestCheckFunc) resource.TestCheckFunc {
	modifiedSet := ConvertToPreviewProviderV2AttrsSet(usePreviewProvider, attrsSet)
	modifiedMap := ConvertToPreviewProviderV2AttrsMap(usePreviewProvider, attrsMap)
	return CheckRSAndDS(resourceName, dataSourceName, pluralDataSourceName, modifiedSet, modifiedMap, extra...)
}

func TestCheckResourceAttrPreviewProviderV2(usePreviewProvider bool, name, key, value string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttr(name, AttrNameToPreviewProviderV2(usePreviewProvider, key), value)
}

func TestCheckResourceAttrSetPreviewProviderV2(usePreviewProvider bool, name, key string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(name, AttrNameToPreviewProviderV2(usePreviewProvider, key))
}

func TestCheckResourceAttrWithPreviewProviderV2(usePreviewProvider bool, name, key string, checkValueFunc resource.CheckResourceAttrWithFunc) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrWith(name, AttrNameToPreviewProviderV2(usePreviewProvider, key), checkValueFunc)
}

func TestCheckTypeSetElemNestedAttrsPreviewProviderV2(usePreviewProvider bool, name, key string, values map[string]string) resource.TestCheckFunc {
	return resource.TestCheckTypeSetElemNestedAttrs(name, AttrNameToPreviewProviderV2(usePreviewProvider, key), values)
}

func AddAttrChecksPreviewProviderV2(usePreviewProvider bool, name string, checks []resource.TestCheckFunc, mapChecks map[string]string) []resource.TestCheckFunc {
	return AddAttrChecks(name, checks, ConvertToPreviewProviderV2AttrsMap(usePreviewProvider, mapChecks))
}

func AddAttrSetChecksPreviewProviderV2(usePreviewProvider bool, name string, checks []resource.TestCheckFunc, attrNames ...string) []resource.TestCheckFunc {
	return AddAttrSetChecks(name, checks, ConvertToPreviewProviderV2AttrsSet(usePreviewProvider, attrNames)...)
}

func AddAttrChecksPrefixPreviewProviderV2(usePreviewProvider bool, name string, checks []resource.TestCheckFunc, mapChecks map[string]string, prefix string, skipNames ...string) []resource.TestCheckFunc {
	return AddAttrChecksPrefix(name, checks, ConvertToPreviewProviderV2AttrsMap(usePreviewProvider, mapChecks), prefix, skipNames...)
}

func ConvertToPreviewProviderV2AttrsMap(usePreviewProvider bool, attrsMap map[string]string) map[string]string {
	if skipPreviewProviderV2Work(usePreviewProvider) {
		return attrsMap
	}
	ret := make(map[string]string, len(attrsMap))
	for name, value := range attrsMap {
		ret[AttrNameToPreviewProviderV2(usePreviewProvider, name)] = value
	}
	return ret
}

func ConvertToPreviewProviderV2AttrsSet(usePreviewProvider bool, attrsSet []string) []string {
	if skipPreviewProviderV2Work(usePreviewProvider) {
		return attrsSet
	}
	ret := make([]string, 0, len(attrsSet))
	for _, name := range attrsSet {
		ret = append(ret, AttrNameToPreviewProviderV2(usePreviewProvider, name))
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
	"pinned_fcv",
	"timeouts",
	"connection_strings",
	"tags",
}

func AttrNameToPreviewProviderV2(usePreviewProvider bool, name string) string {
	if skipPreviewProviderV2Work(usePreviewProvider) {
		return name
	}
	for _, singleAttrName := range tpfSingleNestedAttrs {
		name = strings.ReplaceAll(name, singleAttrName+".0", singleAttrName)
	}
	return name
}

func ConvertAdvancedClusterToPreviewProviderV2(t *testing.T, usePreviewProvider bool, def string) string {
	t.Helper()
	if skipPreviewProviderV2Work(usePreviewProvider) {
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
		convertAttrs(t, "replication_specs", writeBody, true, getReplicationSpecs)
		convertAttrs(t, "advanced_configuration", writeBody, false, hcl.GetAttrVal)
		convertAttrs(t, "bi_connector_config", writeBody, false, hcl.GetAttrVal)
		convertAttrs(t, "pinned_fcv", writeBody, false, hcl.GetAttrVal)
		convertAttrs(t, "timeouts", writeBody, false, hcl.GetAttrVal)
		convertKeyValueAttrs(t, "labels", writeBody)
		convertKeyValueAttrs(t, "tags", writeBody)
	}
	result := string(parse.Bytes())
	result = AttrNameToPreviewProviderV2(usePreviewProvider, result) // useful for lifecycle ingore definitions
	return result
}

func skipPreviewProviderV2Work(usePreviewProvider bool) bool {
	return !config.PreviewProviderV2AdvancedCluster() || !usePreviewProvider
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
		writeBody.RemoveBlock(match) // RemoveBlock doesn't remove newline just after the block so an extra line is added
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

func convertKeyValueAttrs(t *testing.T, name string, writeBody *hclwrite.Body) {
	t.Helper()
	vals := make(map[string]cty.Value)
	for {
		match := writeBody.FirstMatchingBlock(name, nil)
		if match == nil {
			break
		}
		attrs := hcl.GetAttrVal(t, hcl.GetBlockBody(t, match))
		key := attrs.GetAttr("key")
		value := attrs.GetAttr("value")
		vals[key.AsString()] = value
		writeBody.RemoveBlock(match) // RemoveBlock doesn't remove newline just after the block so an extra line is added
	}
	if len(vals) > 0 {
		writeBody.SetAttributeValue(name, cty.ObjectVal(vals))
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
	attributeValues := map[string]cty.Value{
		name: cty.TupleVal(vals),
	}
	hcl.AddAttributes(t, body, attributeValues)
	return cty.ObjectVal(attributeValues)
}
