package codespec_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/stretchr/testify/assert"
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
