package serviceaccountaccesslistentry

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
)

const (
	cidrBlockDesc = "Range of IP addresses in CIDR notation to be added to the access list. You can set a value for this parameter or **ip_address**, but not for both."
	ipAddressDesc = "IP address to be added to the access list. You can set a value for this parameter or **cidr_block**, but not for both."
)

func ResourceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"org_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the organization that contains your projects.",
			},
			"client_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The Client ID of the Service Account.",
			},
			"cidr_block": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: cidrBlockDesc,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validate.ValidCIDR(),
					stringvalidator.ConflictsWith(path.MatchRoot("ip_address")),
				},
			},
			"ip_address": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: ipAddressDesc,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validate.ValidIP(),
					stringvalidator.ConflictsWith(path.MatchRoot("cidr_block")),
				},
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Date the entry was added to the access list. This attribute expresses its value in the ISO 8601 timestamp format in UTC.",
			},
			"last_used_address": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Network address that issued the most recent request to the API.",
			},
			"last_used_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Date when the API received the most recent request that originated from this network address.",
			},
			"request_count": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The number of requests that has originated from this network address.",
			},
		},
	}
}

type TFServiceAccountAccessListEntryModel struct {
	OrgID           types.String `tfsdk:"org_id"`
	ClientID        types.String `tfsdk:"client_id"`
	IPAddress       types.String `tfsdk:"ip_address"`
	CIDRBlock       types.String `tfsdk:"cidr_block"`
	CreatedAt       types.String `tfsdk:"created_at"`
	LastUsedAddress types.String `tfsdk:"last_used_address"`
	LastUsedAt      types.String `tfsdk:"last_used_at"`
	RequestCount    types.Int64  `tfsdk:"request_count"`
}

type TFServiceAccountAccessListEntriesPluralDSModel struct {
	OrgID    types.String                            `tfsdk:"org_id"`
	ClientID types.String                            `tfsdk:"client_id"`
	Results  []*TFServiceAccountAccessListEntryModel `tfsdk:"results"`
}
