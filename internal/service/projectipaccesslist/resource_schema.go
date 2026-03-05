package projectipaccesslist

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
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
	cidrBlockDesc        = "Range of IP addresses in CIDR notation to be added to the access list. Mutually exclusive with `ip_address` and `aws_security_group`."
	ipAddressDesc        = "Single IP address to be added to the access list. Mutually exclusive with `cidr_block` and `aws_security_group`."
	awsSecurityGroupDesc = "Unique identifier of the AWS security group to add to the access list. Mutually exclusive with `cidr_block` and `ip_address`."
)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Provides an IP Access List entry resource. The access list grants access from IPs, CIDRs or AWS Security Groups (if VPC Peering is enabled) to clusters within the Project.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier used for terraform for internal management.",
			},
			"project_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project.",
			},
			"cidr_block": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validate.ValidCIDR(),
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRelative().AtParent().AtName("aws_security_group"),
						path.MatchRelative().AtParent().AtName("ip_address"),
					}...),
				},
				MarkdownDescription: cidrBlockDesc,
			},
			"ip_address": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validate.ValidIP(),
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRelative().AtParent().AtName("aws_security_group"),
						path.MatchRelative().AtParent().AtName("cidr_block"),
					}...),
				},
				MarkdownDescription: ipAddressDesc,
			},
			"aws_security_group": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRelative().AtParent().AtName("ip_address"),
						path.MatchRelative().AtParent().AtName("cidr_block"),
					}...),
				},
				MarkdownDescription: awsSecurityGroupDesc,
			},
			"comment": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "Remark that explains the purpose or scope of this IP access list entry.",
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Read:   true,
				Delete: true,
			}),
		},
	}
}

type TfProjectIPAccessListModel struct {
	ID               types.String   `tfsdk:"id"`
	ProjectID        types.String   `tfsdk:"project_id"`
	CIDRBlock        types.String   `tfsdk:"cidr_block"`
	IPAddress        types.String   `tfsdk:"ip_address"`
	AWSSecurityGroup types.String   `tfsdk:"aws_security_group"`
	Comment          types.String   `tfsdk:"comment"`
	Timeouts         timeouts.Value `tfsdk:"timeouts"`
}
