package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const alertConfigurationsDataSourceName = "alert_configurations"

var _ datasource.DataSource = &AlertConfigurationDS{}
var _ datasource.DataSourceWithConfigure = &AlertConfigurationDS{}

type tfAlertConfigurationsDSModel struct {
	ID          types.String                  `tfsdk:"id"`
	ProjectID   types.String                  `tfsdk:"project_id"`
	OutputType  []string                      `tfsdk:"output_type"`
	ListOptions []tfListOptionsModel          `tfsdk:"list_options"`
	Results     []tfAlertConfigurationDSModel `tfsdk:"results"`
	TotalCount  types.Int64                   `tfsdk:"total_count"`
}

type tfListOptionsModel struct {
	PageNum      types.Int64 `tfsdk:"page_num"`
	ItemsPerPage types.Int64 `tfsdk:"items_per_page"`
	IncludeCount types.Bool  `tfsdk:"include_count"`
}

func NewAlertConfigurationsDS() datasource.DataSource {
	return &AlertConfigurationsDS{
		DSCommon: DSCommon{
			dataSourceName: alertConfigurationsDataSourceName,
		},
	}
}

type AlertConfigurationsDS struct {
	DSCommon
}

func (d *AlertConfigurationsDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"project_id": schema.StringAttribute{
				Required: true,
			},
			"total_count": schema.Int64Attribute{
				Computed: true,
			},
			"output_type": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(stringvalidator.OneOf("resource_hcl", "resource_import")),
				},
			},
			"results": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: copyAndAdd(alertConfigDSSchemaAttributes,
						"output",
						schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"type": schema.StringAttribute{
										Computed: true,
									},
									"label": schema.StringAttribute{
										Computed: true,
									},
									"value": schema.StringAttribute{
										Computed: true,
									},
								},
							},
						}),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"list_options": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"page_num": schema.Int64Attribute{
							Optional: true,
						},
						"items_per_page": schema.Int64Attribute{
							Optional: true,
						},
						"include_count": schema.BoolAttribute{
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func copyAndAdd(m map[string]schema.Attribute, k string, v schema.Attribute) map[string]schema.Attribute {
	newMap := make(map[string]schema.Attribute, len(m)+1)

	for key, value := range m {
		newMap[key] = value
	}

	newMap[k] = v
	return newMap
}

func (d *AlertConfigurationsDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var alertConfigurationsConfig tfAlertConfigurationsDSModel
	conn := d.client.Atlas

	resp.Diagnostics.Append(req.Config.Get(ctx, &alertConfigurationsConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := alertConfigurationsConfig.ProjectID.ValueString()

	alertConfigurationsConfig.ListOptions = setDefaultValuesInListOptions(alertConfigurationsConfig.ListOptions)

	alerts, _, err := conn.AlertConfigurations.List(ctx, projectID, newListOptions(alertConfigurationsConfig.ListOptions))
	if err != nil {
		resp.Diagnostics.AddError(errorReadAlertConf, err.Error())
		return
	}

	alertConfigurationsConfig.Results = newTFAlertConfigurationDSModelList(alerts, projectID, alertConfigurationsConfig.OutputType)
	alertConfigurationsConfig.ID = types.StringValue(encodeStateID(map[string]string{
		"project_id": projectID,
	}))

	resp.Diagnostics.Append(resp.State.Set(ctx, &alertConfigurationsConfig)...)
}

func newTFAlertConfigurationDSModelList(alerts []matlas.AlertConfiguration, projectID string, definedOutputs []string) []tfAlertConfigurationDSModel {
	outputConfigurations := make([]tfAlertConfigurationOutputModel, len(definedOutputs))
	for i, output := range definedOutputs {
		outputConfigurations[i] = tfAlertConfigurationOutputModel{
			Type: types.StringValue(output),
		}
	}

	results := make([]tfAlertConfigurationDSModel, len(alerts))

	for i := 0; i < len(alerts); i++ {
		alert := alerts[i]
		label := fmt.Sprintf("%s_%d", alert.EventTypeName, i)
		resultAlertConfigModel := newTFAlertConfigurationDSModel(&alerts[i], projectID)
		computedOutputs := computeAlertConfigurationOutput(&alert, outputConfigurations, label)
		resultAlertConfigModel.Output = computedOutputs
		results[i] = resultAlertConfigModel
	}

	return results
}

const listOptionDefaultPageNum = 0
const listOptionDefaultItemsPerPage = 100
const listOptionDefaultIncludeCount = false

func setDefaultValuesInListOptions(listOptionsArr []tfListOptionsModel) []tfListOptionsModel {
	var result = make([]tfListOptionsModel, len(listOptionsArr))
	for i, v := range listOptionsArr {
		result[i] = tfListOptionsModel{
			PageNum:      types.Int64Value(listOptionDefaultPageNum),
			ItemsPerPage: types.Int64Value(listOptionDefaultItemsPerPage),
			IncludeCount: types.BoolValue(listOptionDefaultIncludeCount),
		}
		if !v.PageNum.IsNull() {
			result[i].PageNum = v.PageNum
		}
		if !v.ItemsPerPage.IsNull() {
			result[i].ItemsPerPage = v.ItemsPerPage
		}
		if !v.IncludeCount.IsNull() {
			result[i].IncludeCount = v.IncludeCount
		}
	}
	return result
}

func newListOptions(listOptionsArr []tfListOptionsModel) *matlas.ListOptions {
	result := &matlas.ListOptions{
		PageNum:      listOptionDefaultPageNum,
		ItemsPerPage: listOptionDefaultItemsPerPage,
		IncludeCount: listOptionDefaultIncludeCount,
	}

	if len(listOptionsArr) > 0 {
		listOption := listOptionsArr[0]
		if !listOption.PageNum.IsNull() {
			result.PageNum = int(listOption.PageNum.ValueInt64())
		}
		if !listOption.ItemsPerPage.IsNull() {
			result.ItemsPerPage = int(listOption.ItemsPerPage.ValueInt64())
		}
		if !listOption.IncludeCount.IsNull() {
			result.IncludeCount = listOption.IncludeCount.ValueBool()
		}
	}
	return result
}
