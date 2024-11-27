package conversion_test

import (
	"context"
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
			"nestSingle": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "desc nestSingle",
				Attributes: map[string]schema.Attribute{
					"nestedSingleAttr": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "desc nestedSingleAttr",
					},
				},
			},
			"nestList": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "desc nestList",
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
			"nestSet": schema.SetNestedAttribute{
				Computed:            true,
				MarkdownDescription: "desc nestSet",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"nestedAttr": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "desc nested nestSet",
						},
					},
				},
			},
			"timeouts": timeouts.Attributes(context.Background(), timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
			}),
		},
		Blocks: map[string]schema.Block{
			"nestBlock": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"nestBlockAttr": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "desc nestBlockAttr",
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
			"nestSingle": dsschema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "desc nestSingle",
				Description:         "desc nestSingle",
				Attributes: map[string]dsschema.Attribute{
					"nestedSingleAttr": dsschema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "desc nestedSingleAttr",
						Description:         "desc nestedSingleAttr",
					},
				},
			},
			"nestList": dsschema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "desc nestList",
				Description:         "desc nestList",
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
			"nestSet": dsschema.SetNestedAttribute{
				Computed:            true,
				MarkdownDescription: "desc nestSet",
				Description:         "desc nestSet",
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
		Blocks: map[string]dsschema.Block{
			"nestBlock": dsschema.SetNestedBlock{
				NestedObject: dsschema.NestedBlockObject{
					Attributes: map[string]dsschema.Attribute{
						"nestBlockAttr": dsschema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "desc nestBlockAttr",
							Description:         "desc nestBlockAttr",
						},
					},
				},
			},
		},
	}

	requiredFields := []string{"requiredAttrString", "requiredAttrInt64"}
	overridenFields := map[string]dsschema.Attribute{
		"overridenString": dsschema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "desc overridenString",
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRoot("otherAttr")),
			},
		},
		"attrDelete": nil,
	}
	ds := conversion.DataSourceSchemaFromResource(s, requiredFields, overridenFields)
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
