package conversion_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/stretchr/testify/assert"
)

func TestDataSourceSchemaFromResource(t *testing.T) {
	s := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"requiredAttrString": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "desc requiredAttrString",
			},
			"requiredAttrInt64": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "desc requiredAttrInt64",
			},
			"attrString": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "desc attrString",
			},
			"attrInt64": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "desc attrInt64",
			},
			"attrFloat64": schema.Float64Attribute{
				Computed:            true,
				MarkdownDescription: "desc attrFloat64",
			},
			"attrBool": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "desc attrBool",
			},
			"attrDelete": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "desc attrDelete",
			},
			"attrSensitive": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "desc attrSensitive",
			},
			"mapAttr": schema.MapAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "desc mapAttr",
			},
			"listAttr": schema.ListAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "desc listAttr",
			},
			"setAttr": schema.SetAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "desc setAttr",
			},
			"singleNestedAttribute": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "desc singleNestedAttribute",
				Attributes: map[string]schema.Attribute{
					"singleNestedAttributeAttr": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "desc singleNestedAttributeAttr",
					},
				},
			},
			"listNestedAttribute": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "desc listNestedAttribute",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"nestedAttr": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "desc nested nestList",
						},
						"requiredAttrString": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "desc required not matched nested",
						},
					},
				},
			},
			"setNestedAttribute": schema.SetNestedAttribute{
				Computed:            true,
				MarkdownDescription: "desc setNestedAttribute",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"nestedAttr": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "desc nested nestSet",
						},
					},
				},
			},
			"timeouts": timeouts.Attributes(t.Context(), timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
			}),
		},
	}

	expected := dsschema.Schema{
		Attributes: map[string]dsschema.Attribute{
			"requiredAttrString": dsschema.StringAttribute{
				Required:            true,
				MarkdownDescription: "desc requiredAttrString",
				Description:         "desc requiredAttrString",
			},
			"requiredAttrInt64": dsschema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "desc requiredAttrInt64",
				Description:         "desc requiredAttrInt64",
			},
			"attrString": dsschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "desc attrString",
				Description:         "desc attrString",
			},
			"attrInt64": dsschema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "desc attrInt64",
				Description:         "desc attrInt64",
			},
			"attrFloat64": dsschema.Float64Attribute{
				Computed:            true,
				MarkdownDescription: "desc attrFloat64",
				Description:         "desc attrFloat64",
			},
			"attrBool": dsschema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "desc attrBool",
				Description:         "desc attrBool",
			},
			"attrSensitive": dsschema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "desc attrSensitive",
				Description:         "desc attrSensitive",
			},
			"mapAttr": dsschema.MapAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "desc mapAttr",
				Description:         "desc mapAttr",
			},
			"listAttr": dsschema.ListAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "desc listAttr",
				Description:         "desc listAttr",
			},
			"setAttr": dsschema.SetAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "desc setAttr",
				Description:         "desc setAttr",
			},
			"singleNestedAttribute": dsschema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "desc singleNestedAttribute",
				Description:         "desc singleNestedAttribute",
				Attributes: map[string]dsschema.Attribute{
					"singleNestedAttributeAttr": dsschema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "desc singleNestedAttributeAttr",
						Description:         "desc singleNestedAttributeAttr",
					},
				},
			},
			"listNestedAttribute": dsschema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "desc listNestedAttribute",
				Description:         "desc listNestedAttribute",
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"nestedAttr": dsschema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "desc nested nestList",
							Description:         "desc nested nestList",
						},
						"requiredAttrString": dsschema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "desc required not matched nested",
							Description:         "desc required not matched nested",
						},
					},
				},
			},
			"setNestedAttribute": dsschema.SetNestedAttribute{
				Computed:            true,
				MarkdownDescription: "desc setNestedAttribute",
				Description:         "desc setNestedAttribute",
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"nestedAttr": dsschema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "desc nested nestSet",
							Description:         "desc nested nestSet",
						},
					},
				},
			},
			"overridenString": dsschema.StringAttribute{
				Computed:            true,
				Description:         "desc overridenString",
				MarkdownDescription: "desc overridenString",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("otherAttr")),
				},
			},
		},
	}

	ds := conversion.DataSourceSchemaFromResource(s, &conversion.DataSourceSchemaRequest{
		RequiredFields: []string{"requiredAttrString", "requiredAttrInt64"},
		OverridenFields: map[string]dsschema.Attribute{
			"overridenString": dsschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "desc overridenString",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("otherAttr")),
				},
			},
			"attrDelete": nil,
		},
	})
	assert.Equal(t, expected, ds)
}

func TestDataSourceSchemaFromResource_blocksToAttrs(t *testing.T) {
	s := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"requiredAttrString": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "desc requiredAttrString",
			},
			"attrString": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "desc attrString",
			},
		},
		Blocks: map[string]schema.Block{
			"setNestedBlock": schema.SetNestedBlock{
				MarkdownDescription: "desc setNestedBlock",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"setNestedBlockAttr": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "desc setNestedBlockAttr",
						},
					},
					Blocks: map[string]schema.Block{
						"bb 1": schema.SingleNestedBlock{
							MarkdownDescription: "desc bb 1",
							Attributes: map[string]schema.Attribute{
								"bb attr 1": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "desc bb attr 1",
								},
							},
						},
					},
				},
			},
			"listNestedBlock": schema.ListNestedBlock{
				MarkdownDescription: "desc listNestedBlock",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"listNestedBlockAttr": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "desc listNestedBlockAttr",
						},
					},
					Blocks: map[string]schema.Block{
						"bb 2": schema.ListNestedBlock{
							MarkdownDescription: "desc bb 2",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"bb attr 2": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "desc bb attr 2",
									},
								},
							},
						},
					},
				},
			},
			"singleNestedBlock": schema.SingleNestedBlock{
				MarkdownDescription: "desc singleNestedBlock",
				Attributes: map[string]schema.Attribute{
					"nestattr": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "desc nestattr",
					},
				},
				Blocks: map[string]schema.Block{
					"bb 3": schema.ListNestedBlock{
						MarkdownDescription: "desc bb 3",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"bb attr 3": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "desc bb attr 3",
								},
							},
						},
					},
				},
			},
		},
	}

	expected := dsschema.Schema{
		Attributes: map[string]dsschema.Attribute{
			"requiredAttrString": dsschema.StringAttribute{
				Required:            true,
				MarkdownDescription: "desc requiredAttrString",
				Description:         "desc requiredAttrString",
			},
			"attrString": dsschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "desc attrString",
				Description:         "desc attrString",
			},
			"setNestedBlock": dsschema.SetNestedAttribute{
				Computed:            true,
				Description:         "desc setNestedBlock",
				MarkdownDescription: "desc setNestedBlock",
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"setNestedBlockAttr": dsschema.StringAttribute{
							Computed:            true,
							Description:         "desc setNestedBlockAttr",
							MarkdownDescription: "desc setNestedBlockAttr",
						},
						"bb 1": dsschema.SingleNestedAttribute{
							Computed:            true,
							Description:         "desc bb 1",
							MarkdownDescription: "desc bb 1",
							Attributes: map[string]dsschema.Attribute{
								"bb attr 1": dsschema.StringAttribute{
									Computed:            true,
									Description:         "desc bb attr 1",
									MarkdownDescription: "desc bb attr 1",
								},
							},
						},
					},
				},
			},
			"listNestedBlock": dsschema.ListNestedAttribute{
				Computed:            true,
				Description:         "desc listNestedBlock",
				MarkdownDescription: "desc listNestedBlock",
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"listNestedBlockAttr": dsschema.StringAttribute{
							Computed:            true,
							Description:         "desc listNestedBlockAttr",
							MarkdownDescription: "desc listNestedBlockAttr",
						},
						"bb 2": dsschema.ListNestedAttribute{
							Computed:            true,
							Description:         "desc bb 2",
							MarkdownDescription: "desc bb 2",
							NestedObject: dsschema.NestedAttributeObject{
								Attributes: map[string]dsschema.Attribute{
									"bb attr 2": dsschema.StringAttribute{
										Computed:            true,
										Description:         "desc bb attr 2",
										MarkdownDescription: "desc bb attr 2",
									},
								},
							},
						},
					},
				},
			},
			"singleNestedBlock": dsschema.SingleNestedAttribute{
				Computed:            true,
				Description:         "desc singleNestedBlock",
				MarkdownDescription: "desc singleNestedBlock",
				Attributes: map[string]dsschema.Attribute{
					"nestattr": dsschema.StringAttribute{
						Computed:            true,
						Description:         "desc nestattr",
						MarkdownDescription: "desc nestattr",
					},
					"bb 3": dsschema.ListNestedAttribute{
						Computed:            true,
						Description:         "desc bb 3",
						MarkdownDescription: "desc bb 3",
						NestedObject: dsschema.NestedAttributeObject{
							Attributes: map[string]dsschema.Attribute{
								"bb attr 3": dsschema.StringAttribute{
									Computed:            true,
									Description:         "desc bb attr 3",
									MarkdownDescription: "desc bb attr 3",
								},
							},
						},
					},
				},
			},
		},
	}

	ds := conversion.DataSourceSchemaFromResource(s, &conversion.DataSourceSchemaRequest{
		RequiredFields: []string{"requiredAttrString"},
	})
	assert.Equal(t, expected, ds)
}

func TestPluralDataSourceSchemaFromResource(t *testing.T) {
	s := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"requiredAttrString": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "desc requiredAttrString",
			},
			"attrString": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "desc attrString",
			},
			"attrDelete": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "desc attrDelete",
			},
			"overridenString": dsschema.StringAttribute{
				Computed:            true,
				Description:         "desc overridenString",
				MarkdownDescription: "desc overridenString",
			},
		},
		Blocks: map[string]schema.Block{
			"nested": schema.ListNestedBlock{
				MarkdownDescription: "desc nested",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"nested attr": schema.StringAttribute{
							MarkdownDescription: "desc nested attr",
						},
					},
				},
			},
		},
	}

	expected := dsschema.Schema{
		Attributes: map[string]dsschema.Attribute{
			"requiredAttrString": dsschema.StringAttribute{
				Required:            true,
				MarkdownDescription: "desc requiredAttrString",
				Description:         "desc requiredAttrString",
			},
			"overridenRootStringOptional": dsschema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "desc overridenRootStringOptional",
				Description:         "desc overridenRootStringOptional",
			},
			"overridenRootStringRequired": dsschema.StringAttribute{
				Required:            true,
				MarkdownDescription: "desc overridenRootStringRequired",
				Description:         "desc overridenRootStringRequired",
			},
			"results": dsschema.ListNestedAttribute{
				Computed: true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"requiredAttrString": dsschema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "desc requiredAttrString",
							Description:         "desc requiredAttrString",
						},
						"attrString": dsschema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "desc attrString",
							Description:         "desc attrString",
						},
						"overridenString": dsschema.StringAttribute{
							Computed:            true,
							Description:         "desc overridenString",
							MarkdownDescription: "desc overridenString",
							Validators: []validator.String{
								stringvalidator.ConflictsWith(path.MatchRoot("otherAttr")),
							},
						},
						"nested": dsschema.ListNestedAttribute{
							Computed:            true,
							Description:         "desc nested",
							MarkdownDescription: "desc nested",
							NestedObject: dsschema.NestedAttributeObject{
								Attributes: map[string]dsschema.Attribute{
									"nested attr": dsschema.StringAttribute{
										Computed:            true,
										Description:         "desc nested attr",
										MarkdownDescription: "desc nested attr",
									},
								},
							},
						},
					},
				},
				Description:         "List of documents that MongoDB Cloud returns for this request.",
				MarkdownDescription: "List of documents that MongoDB Cloud returns for this request.",
			},
		},
	}

	ds := conversion.PluralDataSourceSchemaFromResource(s, &conversion.PluralDataSourceSchemaRequest{
		RequiredFields: []string{"requiredAttrString"},
		OverridenFields: map[string]dsschema.Attribute{
			"overridenString": dsschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "desc overridenString",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("otherAttr")),
				},
			},
			"attrDelete": nil,
		},
		OverridenRootFields: map[string]dsschema.Attribute{
			"overridenRootStringOptional": dsschema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "desc overridenRootStringOptional",
			},
			"overridenRootStringRequired": dsschema.StringAttribute{
				Required:            true,
				MarkdownDescription: "desc overridenRootStringRequired",
			},
		},
	})
	assert.Equal(t, expected, ds)
}

func TestPluralDataSourceSchemaFromResource_legacyFields(t *testing.T) {
	s := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"requiredAttrString": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "desc requiredAttrString",
			},
			"attrString": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "desc attrString",
			},
		},
	}

	expected := dsschema.Schema{
		Attributes: map[string]dsschema.Attribute{
			"requiredAttrString": dsschema.StringAttribute{
				Required:            true,
				MarkdownDescription: "desc requiredAttrString",
				Description:         "desc requiredAttrString",
			},
			"id": dsschema.StringAttribute{
				Computed: true,
			},
			"page_num": dsschema.Int64Attribute{
				Optional: true,
			},
			"items_per_page": dsschema.Int64Attribute{
				Optional: true,
			},
			"total_count": dsschema.Int64Attribute{
				Computed: true,
			},
			"results": dsschema.ListNestedAttribute{
				Computed: true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"requiredAttrString": dsschema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "desc requiredAttrString",
							Description:         "desc requiredAttrString",
						},
						"attrString": dsschema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "desc attrString",
							Description:         "desc attrString",
						},
					},
				},
				Description:         "List of documents that MongoDB Cloud returns for this request.",
				MarkdownDescription: "List of documents that MongoDB Cloud returns for this request.",
			},
		},
	}
	ds := conversion.PluralDataSourceSchemaFromResource(s, &conversion.PluralDataSourceSchemaRequest{
		RequiredFields:  []string{"requiredAttrString"},
		HasLegacyFields: true,
	})
	assert.Equal(t, expected, ds)
}

func TestPluralDataSourceSchemaFromResource_overrideResultsDoc(t *testing.T) {
	doc := "results doc"

	s := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"attrString": dsschema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "desc attrString",
			},
		},
	}

	expected := dsschema.Schema{
		Attributes: map[string]dsschema.Attribute{
			"results": dsschema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: doc,
				Description:         doc,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"attrString": dsschema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "desc attrString",
							Description:         "desc attrString",
						},
					},
				},
			},
		},
	}
	ds := conversion.PluralDataSourceSchemaFromResource(s, &conversion.PluralDataSourceSchemaRequest{
		OverrideResultsDoc: doc,
	})
	assert.Equal(t, expected, ds)
}

func TestUpdateSchemaDescription(t *testing.T) {
	s := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "mddesc project_id",
			},
			"nested": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "mddesc nested",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"attr1": schema.StringAttribute{
							Description: "desc attr1",
							Computed:    true,
						},
						"attr2": schema.StringAttribute{
							MarkdownDescription: "mddesc attr2",
							Computed:            true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"list": schema.ListNestedBlock{
				MarkdownDescription: "mddesc list",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"attr3": schema.BoolAttribute{
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "mddesc attr3",
						},
					},
				},
			},
			"set": schema.SetNestedBlock{
				MarkdownDescription: "mddesc set",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"attr4": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "mddesc attr4",
						},
						"attr5": schema.Int64Attribute{
							Required:            true,
							MarkdownDescription: "mddesc attr5",
						},
					},
				},
			},
		},
	}

	expected := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				Description:         "mddesc project_id",
				MarkdownDescription: "mddesc project_id",
			},
			"nested": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "mddesc nested",
				MarkdownDescription: "mddesc nested",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"attr1": schema.StringAttribute{
							Description:         "desc attr1",
							MarkdownDescription: "desc attr1",
							Computed:            true,
						},
						"attr2": schema.StringAttribute{
							Description:         "mddesc attr2",
							MarkdownDescription: "mddesc attr2",
							Computed:            true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"list": schema.ListNestedBlock{
				Description:         "mddesc list",
				MarkdownDescription: "mddesc list",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"attr3": schema.BoolAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "mddesc attr3",
							MarkdownDescription: "mddesc attr3",
						},
					},
				},
			},
			"set": schema.SetNestedBlock{
				Description:         "mddesc set",
				MarkdownDescription: "mddesc set",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"attr4": schema.StringAttribute{
							Optional:            true,
							Description:         "mddesc attr4",
							MarkdownDescription: "mddesc attr4",
						},
						"attr5": schema.Int64Attribute{
							Required:            true,
							Description:         "mddesc attr5",
							MarkdownDescription: "mddesc attr5",
						},
					},
				},
			},
		},
	}

	conversion.UpdateSchemaDescription(&s)
	assert.Equal(t, expected, s)
}

func TestUpdateAttrPanic(t *testing.T) {
	testCases := map[string]any{
		"not ptr, please fix caller":    "no ptr",
		"not struct, please fix caller": conversion.Pointer("no struct"),
		"invalid desc fields, please fix caller": conversion.Pointer(struct {
			Description         int
			MarkdownDescription int
		}{}),
		"both descriptions exist, please fix caller: description": conversion.Pointer(struct {
			Description         string
			MarkdownDescription string
		}{
			Description:         "description",
			MarkdownDescription: "markdown description",
		}),
		"not map, please fix caller: Attributes": conversion.Pointer(struct {
			Attributes string
		}{}),
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.PanicsWithValue(t, name, func() {
				conversion.UpdateAttr(tc)
			})
		})
	}
}
