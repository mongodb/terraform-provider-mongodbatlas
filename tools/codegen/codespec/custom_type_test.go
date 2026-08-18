package codespec_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
)

func TestCustomCollectionPackages(t *testing.T) {
	tests := map[string]struct {
		customType *codespec.CustomType
		packages   []codespec.CustomTypePackage
	}{
		"List of JSON": {
			customType: codespec.NewCustomListType(codespec.CustomTypeJSON),
			packages:   []codespec.CustomTypePackage{codespec.CustomTypesPkg, codespec.JSONTypesPkg},
		},
		"Map of JSON": {
			customType: codespec.NewCustomMapType(codespec.CustomTypeJSON),
			packages:   []codespec.CustomTypePackage{codespec.CustomTypesPkg, codespec.JSONTypesPkg},
		},
		"Set of JSON": {
			customType: codespec.NewCustomSetType(codespec.CustomTypeJSON),
			packages:   []codespec.CustomTypePackage{codespec.CustomTypesPkg, codespec.JSONTypesPkg},
		},
		"List of strings": {
			customType: codespec.NewCustomListType(codespec.String),
			packages:   []codespec.CustomTypePackage{codespec.CustomTypesPkg},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.packages, test.customType.Packages)
		})
	}
}

func TestCustomTypeUnmarshalYAML(t *testing.T) {
	tests := map[string]struct {
		yaml     string
		packages []codespec.CustomTypePackage
	}{
		"Legacy package": {
			yaml: `package: customtypes
model: customtypes.ListValue[types.String]
schema: customtypes.NewListType[types.String](ctx)`,
			packages: []codespec.CustomTypePackage{codespec.CustomTypesPkg},
		},
		"Multiple packages": {
			yaml: `packages:
  - customtypes
  - jsontypes
model: customtypes.ListValue[jsontypes.Normalized]
schema: customtypes.NewListType[jsontypes.Normalized](ctx)`,
			packages: []codespec.CustomTypePackage{codespec.CustomTypesPkg, codespec.JSONTypesPkg},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var customType codespec.CustomType
			require.NoError(t, yaml.Unmarshal([]byte(test.yaml), &customType))
			assert.Equal(t, test.packages, customType.Packages)
		})
	}
}
