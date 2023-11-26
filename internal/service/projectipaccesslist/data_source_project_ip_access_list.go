package projectipaccesslist

import (
	"bytes"
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	projectIPAccessList = "project_ip_access_list"
)

type projectIPAccessListDS struct {
	config.DSCommon
}

func NewProjectIPAccessListDS() datasource.DataSource {
	return &projectIPAccessListDS{
		DSCommon: config.DSCommon{
			DataSourceName: projectIPAccessList,
		},
	}
}

var _ datasource.DataSource = &projectIPAccessListDS{}
var _ datasource.DataSourceWithConfigure = &projectIPAccessListDS{}

type tfProjectIPAccessListDSModel struct {
	ID               types.String `tfsdk:"id"`
	ProjectID        types.String `tfsdk:"project_id"`
	CIDRBlock        types.String `tfsdk:"cidr_block"`
	IPAddress        types.String `tfsdk:"ip_address"`
	AWSSecurityGroup types.String `tfsdk:"aws_security_group"`
	Comment          types.String `tfsdk:"comment"`
}

func (d *projectIPAccessListDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"project_id": schema.StringAttribute{
				Required: true,
			},
			"cidr_block": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					validate.ValidCIDR(),
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRelative().AtParent().AtName("aws_security_group"),
						path.MatchRelative().AtParent().AtName("ip_address"),
					}...),
				},
			},
			"ip_address": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					validate.ValidIP(),
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRelative().AtParent().AtName("aws_security_group"),
						path.MatchRelative().AtParent().AtName("cidr_block"),
					}...),
				},
			},
			"aws_security_group": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRelative().AtParent().AtName("ip_address"),
						path.MatchRelative().AtParent().AtName("cidr_block"),
					}...),
				},
			},
			"comment": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *projectIPAccessListDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var databaseDSUserConfig *tfProjectIPAccessListDSModel
	var err error
	resp.Diagnostics.Append(req.Config.Get(ctx, &databaseDSUserConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if databaseDSUserConfig.CIDRBlock.IsNull() && databaseDSUserConfig.IPAddress.IsNull() && databaseDSUserConfig.AWSSecurityGroup.IsNull() {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("validation error", "One of cidr_block, ip_address or aws_security_group needs to contain a value"))
		return
	}

	var entry bytes.Buffer
	entry.WriteString(databaseDSUserConfig.CIDRBlock.ValueString())
	if !databaseDSUserConfig.IPAddress.IsNull() {
		entry.WriteString(databaseDSUserConfig.IPAddress.ValueString())
	} else if !databaseDSUserConfig.AWSSecurityGroup.IsNull() {
		entry.WriteString(databaseDSUserConfig.AWSSecurityGroup.ValueString())
	}

	conn := d.Client.Atlas
	accessList, _, err := conn.ProjectIPAccessList.Get(ctx, databaseDSUserConfig.ProjectID.ValueString(), entry.String())
	if err != nil {
		resp.Diagnostics.AddError("error getting access list entry", err.Error())
		return
	}

	accessListEntry, diagnostic := newTFProjectIPAccessListDSModel(ctx, accessList)
	resp.Diagnostics.Append(diagnostic...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &accessListEntry)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func newTFProjectIPAccessListDSModel(ctx context.Context, accessList *matlas.ProjectIPAccessList) (*tfProjectIPAccessListDSModel, diag.Diagnostics) {
	databaseUserModel := &tfProjectIPAccessListDSModel{
		ProjectID:        types.StringValue(accessList.GroupID),
		Comment:          types.StringValue(accessList.Comment),
		CIDRBlock:        types.StringValue(accessList.CIDRBlock),
		IPAddress:        types.StringValue(accessList.IPAddress),
		AWSSecurityGroup: types.StringValue(accessList.AwsSecurityGroup),
	}

	entry := accessList.CIDRBlock
	if accessList.IPAddress != "" {
		entry = accessList.IPAddress
	} else if accessList.AwsSecurityGroup != "" {
		entry = accessList.AwsSecurityGroup
	}

	id := conversion.EncodeStateID(map[string]string{
		"entry":      entry,
		"project_id": accessList.GroupID,
	})

	databaseUserModel.ID = types.StringValue(id)
	return databaseUserModel, nil
}
