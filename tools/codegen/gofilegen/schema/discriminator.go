package schema

import (
	"fmt"
	"maps"
	"slices"
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
	code := fmt.Sprintf(`%[1]s.ValidateDiscriminator(%[1]s.DiscriminatorDefinition{
			Mapping: map[string]%[1]s.VariantDefinition{
				%[2]s
			},
		})`, customValidatorPkgName, mappingCode)

	return CodeStatement{
		Code:    code,
		Imports: []string{importCustomValidator},
	}
}

func variantDefinitionCode(key string, variant codespec.DiscriminatorType) string {
	if len(variant.Allowed) == 0 && len(variant.Required) == 0 {
		return fmt.Sprintf(`%q: {},`, key)
	}

	allowed := tfSchemaNames(variant.Allowed)
	required := tfSchemaNames(variant.Required)
	var fields []string
	if len(allowed) > 0 {
		fields = append(fields, fmt.Sprintf("Allowed: []string{%s}", `"`+strings.Join(allowed, `", "`)+`"`))
	}
	if len(required) > 0 {
		fields = append(fields, fmt.Sprintf("Required: []string{%s}", `"`+strings.Join(required, `", "`)+`"`))
	}

	return fmt.Sprintf("%q: {\n%s,\n},", key, strings.Join(fields, ",\n"))
}

func sortedDiscriminatorKeys(mapping map[string]codespec.DiscriminatorType) []string {
	return slices.Sorted(maps.Keys(mapping))
}

func tfSchemaNames(names []codespec.DiscriminatorAttrName) []string {
	result := make([]string, len(names))
	for i, n := range names {
		result[i] = n.TFSchemaName
	}
	return result
}
