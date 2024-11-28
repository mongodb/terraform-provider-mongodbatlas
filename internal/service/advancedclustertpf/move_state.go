package advancedclustertpf

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20241113001/admin"
)

const (
	MoveModeEnvVarName   = "MONGODB_ATLAS_TEST_MOVE_MODE"
	MoveModeValPreferred = "preferred"
	MoveModeValRawState  = "rawstate"
	MoveModeValJSON      = "json"
)

// TODO: We temporarily use mongodbatlas_database_user instead of mongodbatlas_cluster to set up the initial environment
func (r *rs) MoveState(context.Context) []resource.StateMover {
	return []resource.StateMover{
		{
			SourceSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
					},
					"project_id": schema.StringAttribute{
						Required: true,
					},
					"auth_database_name": schema.StringAttribute{
						Required: true,
					},
					"username": schema.StringAttribute{
						Required: true,
					},
					"password": schema.StringAttribute{
						Optional:  true,
						Sensitive: true,
					},
					"x509_type": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"oidc_auth_type": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"ldap_auth_type": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"aws_iam_type": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
				},
				Blocks: map[string]schema.Block{
					"roles": schema.SetNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"collection_name": schema.StringAttribute{
									Optional: true,
								},
								"database_name": schema.StringAttribute{
									Required: true,
								},
								"role_name": schema.StringAttribute{
									Required: true,
								},
							},
						},
					},
					"labels": schema.SetNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"key": schema.StringAttribute{
									Optional: true,
									Computed: true,
								},
								"value": schema.StringAttribute{
									Optional: true,
									Computed: true,
								},
							},
						},
					},
					"scopes": schema.SetNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Optional: true,
								},
								"type": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},
				},
			},
			StateMover: stateMoverTemporaryPreferred,
		},
		{
			StateMover: stateMoverTemporaryRawState,
		},
		{
			StateMover: stateMoverTemporaryJSON,
		},
	}
}

func stateMoverTemporaryPreferred(ctx context.Context, req resource.MoveStateRequest, resp *resource.MoveStateResponse) {
	if !isSource(req, "database_user", MoveModeValPreferred) {
		return
	}
	type model struct {
		ID               types.String `tfsdk:"id"`
		ProjectID        types.String `tfsdk:"project_id"`
		AuthDatabaseName types.String `tfsdk:"auth_database_name"`
		Username         types.String `tfsdk:"username"`
		Password         types.String `tfsdk:"password"`
		X509Type         types.String `tfsdk:"x509_type"`
		OIDCAuthType     types.String `tfsdk:"oidc_auth_type"`
		LDAPAuthType     types.String `tfsdk:"ldap_auth_type"`
		AWSIAMType       types.String `tfsdk:"aws_iam_type"`
		Roles            types.Set    `tfsdk:"roles"`
		Labels           types.Set    `tfsdk:"labels"`
		Scopes           types.Set    `tfsdk:"scopes"`
	}
	var state model
	resp.Diagnostics.Append(req.SourceState.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	setMoveState(ctx, state.ProjectID.ValueString(), state.Username.ValueString(), resp)
}

func stateMoverTemporaryRawState(ctx context.Context, req resource.MoveStateRequest, resp *resource.MoveStateResponse) {
	if !isSource(req, "database_user", MoveModeValRawState) {
		return
	}
	// TODO: not need to define the full model if using IgnoreUndefinedAttributes, as in JSON case
	rawStateValue, err := req.SourceRawState.UnmarshalWithOpts(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"id":                 tftypes.String,
			"project_id":         tftypes.String,
			"auth_database_name": tftypes.String,
			"username":           tftypes.String,
			"password":           tftypes.String,
			"x509_type":          tftypes.String,
			"oidc_auth_type":     tftypes.String,
			"ldap_auth_type":     tftypes.String,
			"aws_iam_type":       tftypes.String,
			"roles": tftypes.Set{ElementType: tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"collection_name": tftypes.String,
					"database_name":   tftypes.String,
					"role_name":       tftypes.String,
				},
			}},
			"labels": tftypes.Set{ElementType: tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"key":   tftypes.String,
					"value": tftypes.String,
				},
			}},
			"scopes": tftypes.Set{ElementType: tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"name": tftypes.String,
					"type": tftypes.String,
				},
			}},
		},
	}, tfprotov6.UnmarshalOpts{ValueFromJSONOpts: tftypes.ValueFromJSONOpts{IgnoreUndefinedAttributes: true}})
	if err != nil {
		resp.Diagnostics.AddError("Unable to Unmarshal Source State", err.Error())
		return
	}
	var rawState map[string]tftypes.Value
	if err := rawStateValue.As(&rawState); err != nil {
		resp.Diagnostics.AddError("Unable to Convert Source State", err.Error())
		return
	}
	var projectID *string // TODO: take username as the cluster name
	if err := rawState["project_id"].As(&projectID); err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("project_id"), "Unable to Convert Source State", err.Error())
		return
	}
	var clusterName *string // TODO: take username as the cluster name
	if err := rawState["username"].As(&clusterName); err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("username"), "Unable to Convert Source State", err.Error())
		return
	}
	setMoveState(ctx, *projectID, *clusterName, resp)
}

func stateMoverTemporaryJSON(ctx context.Context, req resource.MoveStateRequest, resp *resource.MoveStateResponse) {
	if !isSource(req, "database_user", MoveModeValJSON) {
		return
	}
	type model struct {
		ProjectID   string `json:"project_id"`
		ClusterName string `json:"username"` // TODO: take username as the cluster name
	}
	var state model
	if err := json.Unmarshal(req.SourceRawState.JSON, &state); err != nil {
		resp.Diagnostics.AddError("Unable to Unmarshal Source State", err.Error())
		return
	}
	setMoveState(ctx, state.ProjectID, state.ClusterName, resp)
}

func isSource(req resource.MoveStateRequest, resourceName, moveMode string) bool {
	return os.Getenv(MoveModeEnvVarName) == moveMode &&
		req.SourceTypeName == "mongodbatlas_"+resourceName &&
		req.SourceSchemaVersion == 0 &&
		strings.HasSuffix(req.SourceProviderAddress, "/mongodbatlas")
}

func setMoveState(ctx context.Context, projectID, clusterName string, resp *resource.MoveStateResponse) {
	// TODO: timeout should be read from source if provided
	timeout := timeouts.Value{
		Object: types.ObjectValueMust(
			map[string]attr.Type{
				"create": types.StringType,
				"update": types.StringType,
				"delete": types.StringType,
			},
			map[string]attr.Value{
				"create": types.StringValue("30m"),
				"update": types.StringValue("30m"),
				"delete": types.StringValue("30m"),
			}),
	}
	// TODO: we need to have a good state (all attributes known or null) but not need to be the final ones as Read is called after
	model := NewTFModel(ctx, &admin.ClusterDescription20240805{
		GroupId: conversion.StringPtr(projectID),
		Name:    conversion.StringPtr(clusterName),
	}, timeout, &resp.Diagnostics, nil)
	if resp.Diagnostics.HasError() {
		return
	}
	AddAdvancedConfig(ctx, model, nil, nil, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.TargetState.Set(ctx, model)...)
}
