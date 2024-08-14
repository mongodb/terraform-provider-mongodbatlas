package encryptionatrestprivateendpoint

import (
	"context"

	"go.mongodb.org/atlas-sdk/v20240805001/admin"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

func PluralDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cloud_provider": schema.StringAttribute{
				Required:            true,
				Description:         "Human-readable label that identifies the cloud provider for the private endpoints to return.",
				MarkdownDescription: "Human-readable label that identifies the cloud provider for the private endpoints to return.",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				Description:         "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.",
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.",
			},
			"results": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"cloud_provider": schema.StringAttribute{
							Computed:            true,
							Description:         "Human-readable label that identifies the cloud provider for the Encryption At Rest private endpoint.",
							MarkdownDescription: "Human-readable label that identifies the cloud provider for the Encryption At Rest private endpoint.",
						},
						"error_message": schema.StringAttribute{
							Computed:            true,
							Description:         "Error message for failures associated with the Encryption At Rest private endpoint.",
							MarkdownDescription: "Error message for failures associated with the Encryption At Rest private endpoint.",
						},
						"id": schema.StringAttribute{
							Computed:            true,
							Description:         "Unique 24-hexadecimal digit string that identifies the Private Endpoint Service.",
							MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the Private Endpoint Service.",
						},
						"private_endpoint_connection_name": schema.StringAttribute{
							Computed:            true,
							Description:         "Connection name of the Azure Private Endpoint.",
							MarkdownDescription: "Connection name of the Azure Private Endpoint.",
						},
						"region_name": schema.StringAttribute{
							Computed:            true,
							Description:         "Cloud provider region in which the Encryption At Rest private endpoint is located.",
							MarkdownDescription: "Cloud provider region in which the Encryption At Rest private endpoint is located.",
						},
						"status": schema.StringAttribute{
							Computed:            true,
							Description:         "State of the Encryption At Rest private endpoint.",
							MarkdownDescription: "State of the Encryption At Rest private endpoint.",
						},
					},
				},
				Computed:            true,
				Description:         "List of returned documents that MongoDB Cloud providers when completing this request.",
				MarkdownDescription: "List of returned documents that MongoDB Cloud providers when completing this request.",
			},
			"total_count": schema.Int64Attribute{
				Computed:            true,
				Description:         "Total number of documents available. MongoDB Cloud omits this value if `includeCount` is set to `false`.",
				MarkdownDescription: "Total number of documents available. MongoDB Cloud omits this value if `includeCount` is set to `false`.",
			},
		},
	}
}

type TFEncryptionAtRestPrivateEndpointsDSModel struct {
	CloudProvider types.String                  `tfsdk:"cloud_provider"`
	ProjectID     types.String                  `tfsdk:"project_id"`
	Results       []TFEarPrivateEndpointDSModel `tfsdk:"results"`
	TotalCount    types.Int64                   `tfsdk:"total_count"`
}

type TFEarPrivateEndpointDSModel struct {
	CloudProvider                 types.String `tfsdk:"cloud_provider"`
	ErrorMessage                  types.String `tfsdk:"error_message"`
	ID                            types.String `tfsdk:"id"`
	PrivateEndpointConnectionName types.String `tfsdk:"private_endpoint_connection_name"`
	RegionName                    types.String `tfsdk:"region_name"`
	Status                        types.String `tfsdk:"status"`
}

func NewTFEarPrivateEndpoints(projectID, cloudProvider string, totalCount *int, results []admin.EARPrivateEndpoint) *TFEncryptionAtRestPrivateEndpointsDSModel {
	return &TFEncryptionAtRestPrivateEndpointsDSModel{
		ProjectID:     types.StringValue(projectID),
		TotalCount:    types.Int64PointerValue(conversion.IntPtrToInt64Ptr(totalCount)),
		CloudProvider: types.StringValue(cloudProvider),
		Results:       NewTFEarPrivateEndpointsDS(results),
	}
}

func NewTFEarPrivateEndpointsDS(endpoints []admin.EARPrivateEndpoint) []TFEarPrivateEndpointDSModel {
	results := make([]TFEarPrivateEndpointDSModel, len(endpoints))

	for i, v := range endpoints {
		results[i] = TFEarPrivateEndpointDSModel{
			CloudProvider:                 conversion.StringNullIfEmpty(v.GetCloudProvider()),
			ErrorMessage:                  conversion.StringNullIfEmpty(v.GetErrorMessage()),
			ID:                            conversion.StringNullIfEmpty(v.GetId()),
			RegionName:                    conversion.StringNullIfEmpty(v.GetRegionName()),
			Status:                        conversion.StringNullIfEmpty(v.GetStatus()),
			PrivateEndpointConnectionName: conversion.StringNullIfEmpty(v.GetPrivateEndpointConnectionName()),
		}
	}
	return results
}
