package project

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"go.mongodb.org/atlas-sdk/v20250312002/admin"
)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"org_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cluster_count": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"created": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_owner_id": schema.StringAttribute{
				Optional: true,
			},
			"with_default_alerts_settings": schema.BoolAttribute{
				// Default values also must be Computed otherwise Terraform throws error:
				// Schema Using Attribute Default For Non-Computed Attribute
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"is_collect_database_specifics_statistics_enabled": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_data_explorer_enabled": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_extended_storage_sizes_enabled": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_performance_advisor_enabled": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_realtime_performance_panel_enabled": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_schema_advisor_enabled": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_slow_operation_thresholding_enabled": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"region_usage_restrictions": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"ip_addresses": schema.SingleNestedAttribute{
				Computed:           true,
				DeprecationMessage: fmt.Sprintf(constant.DeprecationParamByVersionWithReplacement, "1.21.0", "mongodbatlas_project_ip_addresses data source"),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"services": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"clusters": schema.ListNestedAttribute{
								Computed: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"cluster_name": schema.StringAttribute{
											Computed: true,
										},
										"inbound": schema.ListAttribute{
											ElementType: types.StringType,
											Computed:    true,
										},
										"outbound": schema.ListAttribute{
											ElementType: types.StringType,
											Computed:    true,
										},
									},
								},
							},
						},
					},
				},
			},
			"tags": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"teams": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"team_id": schema.StringAttribute{
							Required: true,
						},
						"role_names": schema.SetAttribute{
							Required:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
			"limits": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,
						},
						"value": schema.Int64Attribute{
							Required: true,
						},
						"current_usage": schema.Int64Attribute{
							Computed: true,
						},
						"default_limit": schema.Int64Attribute{
							Computed: true,
						},
						"maximum_limit": schema.Int64Attribute{
							Computed: true,
						},
					},
				},
				// https://discuss.hashicorp.com/t/computed-attributes-and-plan-modifiers/45830/12
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

type TFProjectRSModel struct {
	Limits                                      types.Set    `tfsdk:"limits"`
	Teams                                       types.Set    `tfsdk:"teams"`
	Tags                                        types.Map    `tfsdk:"tags"`
	IPAddresses                                 types.Object `tfsdk:"ip_addresses"`
	RegionUsageRestrictions                     types.String `tfsdk:"region_usage_restrictions"`
	Name                                        types.String `tfsdk:"name"`
	OrgID                                       types.String `tfsdk:"org_id"`
	Created                                     types.String `tfsdk:"created"`
	ProjectOwnerID                              types.String `tfsdk:"project_owner_id"`
	ID                                          types.String `tfsdk:"id"`
	ClusterCount                                types.Int64  `tfsdk:"cluster_count"`
	IsDataExplorerEnabled                       types.Bool   `tfsdk:"is_data_explorer_enabled"`
	IsPerformanceAdvisorEnabled                 types.Bool   `tfsdk:"is_performance_advisor_enabled"`
	IsRealtimePerformancePanelEnabled           types.Bool   `tfsdk:"is_realtime_performance_panel_enabled"`
	IsSchemaAdvisorEnabled                      types.Bool   `tfsdk:"is_schema_advisor_enabled"`
	IsExtendedStorageSizesEnabled               types.Bool   `tfsdk:"is_extended_storage_sizes_enabled"`
	IsCollectDatabaseSpecificsStatisticsEnabled types.Bool   `tfsdk:"is_collect_database_specifics_statistics_enabled"`
	WithDefaultAlertsSettings                   types.Bool   `tfsdk:"with_default_alerts_settings"`
	IsSlowOperationThresholdingEnabled          types.Bool   `tfsdk:"is_slow_operation_thresholding_enabled"`
}

type TFTeamModel struct {
	TeamID    types.String `tfsdk:"team_id"`
	RoleNames types.Set    `tfsdk:"role_names"`
}

type TFLimitModel struct {
	Name         types.String `tfsdk:"name"`
	Value        types.Int64  `tfsdk:"value"`
	CurrentUsage types.Int64  `tfsdk:"current_usage"`
	DefaultLimit types.Int64  `tfsdk:"default_limit"`
	MaximumLimit types.Int64  `tfsdk:"maximum_limit"`
}

type TFIPAddressesModel struct {
	Services TFServicesModel `tfsdk:"services"`
}

type TFServicesModel struct {
	Clusters []TFClusterIPsModel `tfsdk:"clusters"`
}

type TFClusterIPsModel struct {
	ClusterName types.String `tfsdk:"cluster_name"`
	Inbound     types.List   `tfsdk:"inbound"`
	Outbound    types.List   `tfsdk:"outbound"`
}

var IPAddressesObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"services": ServicesObjectType,
}}

var ServicesObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"clusters": types.ListType{ElemType: ClusterIPsObjectType},
}}

var ClusterIPsObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"cluster_name": types.StringType,
	"inbound":      types.ListType{ElemType: types.StringType},
	"outbound":     types.ListType{ElemType: types.StringType},
}}

var TfTeamObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"team_id":    types.StringType,
	"role_names": types.SetType{ElemType: types.StringType},
}}
var TfLimitObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"name":          types.StringType,
	"value":         types.Int64Type,
	"current_usage": types.Int64Type,
	"default_limit": types.Int64Type,
	"maximum_limit": types.Int64Type,
}}

// Resources that need to be cleaned up before a project can be deleted
type AtlasProjectDependants struct {
	AdvancedClusters *admin.PaginatedClusterDescription20240805
}
