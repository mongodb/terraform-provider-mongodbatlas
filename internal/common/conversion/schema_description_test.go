package conversion_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedclustertpf"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/mongodbemployeeaccessgrant"
	"github.com/stretchr/testify/assert"
)

func TestDataSourceSchemaFromResourceTemporary(t *testing.T) {
	s := mongodbemployeeaccessgrant.ResourceSchema(context.Background())
	ds := conversion.DataSourceSchemaFromResource(s)
	fmt.Println(ds)
}

func TestDataSourceSchemasTemporary(t *testing.T) {
	ds := advancedclustertpf.DataSource()
	schemaRequest := datasource.SchemaRequest{}
	schemaResponse := &datasource.SchemaResponse{}
	ds.Schema(context.Background(), schemaRequest, schemaResponse)
}

func TestDataSourceSchemaFromResource_basic(t *testing.T) {
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
			"mapAttr": schema.MapAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "desc mapAttr",
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
			"mapAttr": dsschema.MapAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "desc mapAttr",
				Description:         "desc mapAttr",
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
		},
	}
	ds := conversion.DataSourceSchemaFromResource(s, "requiredAttrString", "requiredAttrInt64")
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
