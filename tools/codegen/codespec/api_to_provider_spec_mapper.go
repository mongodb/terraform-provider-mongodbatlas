//nolint:gocritic
package codespec

import (
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/openapi"
)

// using blank identifiers for now, will be removed in follow-up PRs once logic for conversion is added
func ToProviderSpecModel(atlasAdminAPISpecFilePath, configPath string, resourceName *string) *CodeSpecification {
	_, err := openapi.ParseAtlasAdminAPI(atlasAdminAPISpecFilePath)
	if err != nil {
		panic(err)
	}

	genConfig, _ := config.ParseGenConfigYAML(configPath)

	// var resourceSpec config.Resource
	if resourceName != nil {
		_ = genConfig.Resources[*resourceName]
	}

	// TODO: remove after ToProviderSpecModel() implemented
	return TestExampleCodeSpecification()
}

func TestExampleCodeSpecification() *CodeSpecification {
	testFieldDesc := "Test field description"
	return &CodeSpecification{
		Resources: Resource{
			Schema: &Schema{
				Attributes: Attributes{
					Attribute{
						Name:        "project_id",
						IsRequired:  conversion.Pointer(true),
						String:      &StringAttribute{},
						Description: conversion.StringPtr("Overridden project_id description"),
					},
					Attribute{
						Name:        "bucket_name",
						IsRequired:  conversion.Pointer(true),
						String:      &StringAttribute{},
						Description: &testFieldDesc,
					},
					Attribute{
						Name:        "iam_role_id",
						IsRequired:  conversion.Pointer(true),
						String:      &StringAttribute{},
						Description: &testFieldDesc,
					},
					Attribute{
						Name:        "state",
						IsComputed:  conversion.Pointer(true),
						String:      &StringAttribute{},
						Description: &testFieldDesc,
					},
					Attribute{
						Name:        "prefix_path",
						String:      &StringAttribute{},
						IsComputed:  conversion.Pointer(true),
						IsOptional:  conversion.Pointer(true),
						Description: &testFieldDesc,
					},
				},
			},
			Name: "test_resource",
		},
	}
}
