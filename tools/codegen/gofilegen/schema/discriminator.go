package schema

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
)

const (
	importCustomValidator  = "github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen/customvalidator"
	customValidatorPkgName = "customvalidator"
)

// discriminatorValidatorProperty generates a single ValidateDiscriminator call for a discriminator string attribute.
// It returns just the validator entry; the wrapping Validators: []validator.String{...} block is handled by commonProperties.
func discriminatorValidatorProperty(disc *codespec.Discriminator) CodeStatement {
	var mappingEntries []string
	keys := sortedDiscriminatorKeys(disc.Mapping)
	for _, key := range keys {
		variant := disc.Mapping[key]
		entry := variantDefinitionCode(key, variant)
		mappingEntries = append(mappingEntries, entry)
	}

	mappingCode := strings.Join(mappingEntries, "\n")
	code := fmt.Sprintf(`%s.ValidateDiscriminator(%s.DiscriminatorDefinition{
			Mapping: map[string]%s.VariantDefinition{
				%s
			},
		})`, customValidatorPkgName, customValidatorPkgName, customValidatorPkgName, mappingCode)

	return CodeStatement{
		Code:    code,
		Imports: []string{importCustomValidator},
	}
}

func variantDefinitionCode(key string, variant codespec.DiscriminatorType) string {
	allowed := sortedTFSchemaNames(variant.Allowed)
	required := sortedTFSchemaNames(variant.Required)

	if len(allowed) == 0 && len(required) == 0 {
		return fmt.Sprintf(`%q: {},`, key)
	}

	var fields []string
	if len(allowed) > 0 {
		fields = append(fields, fmt.Sprintf("Allowed: []string{%s}", quotedStringList(allowed)))
	}
	if len(required) > 0 {
		fields = append(fields, fmt.Sprintf("Required: []string{%s}", quotedStringList(required)))
	}

	return fmt.Sprintf("%q: {\n%s,\n},", key, strings.Join(fields, ",\n"))
}

func sortedDiscriminatorKeys(mapping map[string]codespec.DiscriminatorType) []string {
	keys := make([]string, 0, len(mapping))
	for k := range mapping {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortedTFSchemaNames(names []codespec.DiscriminatorAttrName) []string {
	result := make([]string, len(names))
	for i, n := range names {
		result[i] = n.TFSchemaName
	}
	sort.Strings(result)
	return result
}

func quotedStringList(names []string) string {
	quoted := make([]string, len(names))
	for i, n := range names {
		quoted[i] = fmt.Sprintf("%q", n)
	}
	return strings.Join(quoted, ", ")
}
