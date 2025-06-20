// Code generated by terraform-provider-mongodbatlas using `make generate-resource`. DO NOT EDIT.

package streaminstanceapi

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique 24-hexadecimal character string that identifies the project.",
			},
			"cloud_provider": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Label that identifies the cloud service provider where MongoDB Cloud performs stream processing. Currently, this parameter only supports AWS and AZURE.",
			},
			"connections": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of connections configured in the stream instance.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"authentication": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "User credentials required to connect to a Kafka Cluster. Includes the authentication type, as well as the parameters for that authentication mode.",
							Attributes: map[string]schema.Attribute{
								"links": schema.ListNestedAttribute{
									Computed:            true,
									MarkdownDescription: "List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"href": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
											},
											"rel": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "Uniform Resource Locator (URL) that defines the semantic relationship between this resource and another API resource. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
											},
										},
									},
								},
								"mechanism": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Style of authentication. Can be one of PLAIN, SCRAM-256, or SCRAM-512.",
								},
								"password": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Password of the account to connect to the Kafka cluster.",
									Sensitive:           true,
								},
								"ssl_certificate": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "SSL certificate for client authentication to Kafka.",
								},
								"ssl_key": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "SSL key for client authentication to Kafka.",
								},
								"ssl_key_password": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Password for the SSL key, if it is password protected.",
								},
								"username": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Username of the account to connect to the Kafka cluster.",
								},
							},
						},
						"aws": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "AWS configurations for AWS-based connection types.",
							Attributes: map[string]schema.Attribute{
								"links": schema.ListNestedAttribute{
									Computed:            true,
									MarkdownDescription: "List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"href": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
											},
											"rel": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "Uniform Resource Locator (URL) that defines the semantic relationship between this resource and another API resource. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
											},
										},
									},
								},
								"role_arn": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Amazon Resource Name (ARN) that identifies the Amazon Web Services (AWS) Identity and Access Management (IAM) role that MongoDB Cloud assumes when it accesses resources in your AWS account.",
								},
								"test_bucket": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The name of an S3 bucket used to check authorization of the passed-in IAM role ARN.",
								},
							},
						},
						"bootstrap_servers": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Comma separated list of server addresses.",
						},
						"cluster_group_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The id of the group that the cluster belongs to.",
						},
						"cluster_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Name of the cluster configured for this connection.",
						},
						"config": schema.MapAttribute{
							Computed:            true,
							MarkdownDescription: "A map of Kafka key-value pairs for optional configuration. This is a flat object, and keys can have '.' characters.",
							ElementType:         types.StringType,
						},
						"db_role_to_execute": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "The name of a Built in or Custom DB Role to connect to an Atlas Cluster.",
							Attributes: map[string]schema.Attribute{
								"links": schema.ListNestedAttribute{
									Computed:            true,
									MarkdownDescription: "List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"href": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
											},
											"rel": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "Uniform Resource Locator (URL) that defines the semantic relationship between this resource and another API resource. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
											},
										},
									},
								},
								"role": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The name of the role to use. Can be a built in role or a custom role.",
								},
								"type": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Type of the DB role. Can be either BuiltIn or Custom.",
								},
							},
						},
						"headers": schema.MapAttribute{
							Computed:            true,
							MarkdownDescription: "A map of key-value pairs that will be passed as headers for the request.",
							ElementType:         types.StringType,
						},
						"links": schema.ListNestedAttribute{
							Computed:            true,
							MarkdownDescription: "List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"href": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
									},
									"rel": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "Uniform Resource Locator (URL) that defines the semantic relationship between this resource and another API resource. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
									},
								},
							},
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Human-readable label that identifies the stream connection. In the case of the Sample type, this is the name of the sample source.",
						},
						"networking": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Networking Access Type can either be 'PUBLIC' (default) or VPC. VPC type is in public preview, please file a support ticket to enable VPC Network Access.",
							Attributes: map[string]schema.Attribute{
								"access": schema.SingleNestedAttribute{
									Computed:            true,
									MarkdownDescription: "Information about the networking access.",
									Attributes: map[string]schema.Attribute{
										"connection_id": schema.StringAttribute{
											Computed:            true,
											MarkdownDescription: "Reserved. Will be used by PRIVATE_LINK connection type.",
										},
										"links": schema.ListNestedAttribute{
											Computed:            true,
											MarkdownDescription: "List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.",
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"href": schema.StringAttribute{
														Computed:            true,
														MarkdownDescription: "Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
													},
													"rel": schema.StringAttribute{
														Computed:            true,
														MarkdownDescription: "Uniform Resource Locator (URL) that defines the semantic relationship between this resource and another API resource. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
													},
												},
											},
										},
										"name": schema.StringAttribute{
											Computed:            true,
											MarkdownDescription: "Reserved. Will be used by PRIVATE_LINK connection type.",
										},
										"tgw_id": schema.StringAttribute{
											Computed:            true,
											MarkdownDescription: "Reserved. Will be used by TRANSIT_GATEWAY connection type.",
										},
										"tgw_route_id": schema.StringAttribute{
											Computed:            true,
											MarkdownDescription: "Reserved. Will be used by TRANSIT_GATEWAY connection type.",
										},
										"type": schema.StringAttribute{
											Computed:            true,
											MarkdownDescription: "Selected networking type. Either PUBLIC, VPC, PRIVATE_LINK, or TRANSIT_GATEWAY. Defaults to PUBLIC. For VPC, ensure that VPC peering exists and connectivity has been established between Atlas VPC and the VPC where Kafka cluster is hosted for the connection to function properly. TRANSIT_GATEWAY support is coming soon.",
										},
										"vpc_cidr": schema.StringAttribute{
											Computed:            true,
											MarkdownDescription: "Reserved. Will be used by TRANSIT_GATEWAY connection type.",
										},
									},
								},
								"links": schema.ListNestedAttribute{
									Computed:            true,
									MarkdownDescription: "List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"href": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
											},
											"rel": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "Uniform Resource Locator (URL) that defines the semantic relationship between this resource and another API resource. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
											},
										},
									},
								},
							},
						},
						"security": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Properties for the secure transport connection to Kafka. For SSL, this can include the trusted certificate to use.",
							Attributes: map[string]schema.Attribute{
								"broker_public_certificate": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "A trusted, public x509 certificate for connecting to Kafka over SSL.",
								},
								"links": schema.ListNestedAttribute{
									Computed:            true,
									MarkdownDescription: "List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"href": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
											},
											"rel": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "Uniform Resource Locator (URL) that defines the semantic relationship between this resource and another API resource. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
											},
										},
									},
								},
								"protocol": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Describes the transport type. Can be either SASL_PLAINTEXT, SASL_SSL, or SSL.",
								},
							},
						},
						"type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Type of the connection.",
						},
						"url": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The url to be used for the request.",
						},
					},
				},
			},
			"data_process_region": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Information about the cloud provider region in which MongoDB Cloud processes the stream.",
				Attributes: map[string]schema.Attribute{
					"cloud_provider": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Label that identifies the cloud service provider where MongoDB Cloud performs stream processing. Currently, this parameter only supports AWS and AZURE.",
					},
					"links": schema.ListNestedAttribute{
						Computed:            true,
						MarkdownDescription: "List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"href": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
								},
								"rel": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Uniform Resource Locator (URL) that defines the semantic relationship between this resource and another API resource. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
								},
							},
						},
					},
					"region": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Name of the cloud provider region hosting Atlas Stream Processing.",
					},
				},
			},
			"group_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.",
			},
			"hostnames": schema.ListAttribute{
				Computed:            true,
				MarkdownDescription: "List that contains the hostnames assigned to the stream instance.",
				ElementType:         types.StringType,
			},
			"links": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"href": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
						},
						"rel": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Uniform Resource Locator (URL) that defines the semantic relationship between this resource and another API resource. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
						},
					},
				},
			},
			"name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Human-readable label that identifies the stream instance.",
			},
			"region": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the cloud provider region hosting Atlas Stream Processing.",
			},
			"sample_connections": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Sample connections to add to SPI.",
				Attributes: map[string]schema.Attribute{
					"links": schema.ListNestedAttribute{
						Computed:            true,
						MarkdownDescription: "List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"href": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
								},
								"rel": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Uniform Resource Locator (URL) that defines the semantic relationship between this resource and another API resource. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
								},
							},
						},
					},
					"solar": schema.BoolAttribute{
						Computed:            true,
						Optional:            true,
						MarkdownDescription: "Flag that indicates whether to add a 'sample_stream_solar' connection.",
					},
				},
			},
			"stream_config": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Configuration options for an Atlas Stream Processing Instance.",
				Attributes: map[string]schema.Attribute{
					"links": schema.ListNestedAttribute{
						Computed:            true,
						MarkdownDescription: "List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"href": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
								},
								"rel": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Uniform Resource Locator (URL) that defines the semantic relationship between this resource and another API resource. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
								},
							},
						},
					},
					"tier": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Selected tier for the Stream Instance. Configures Memory / VCPU allowances.",
					},
				},
			},
		},
	}
}

type TFModel struct {
	_Id               types.String `tfsdk:"_id" autogen:"omitjson"`
	CloudProvider     types.String `tfsdk:"cloud_provider"`
	Connections       types.List   `tfsdk:"connections" autogen:"omitjson"`
	DataProcessRegion types.Object `tfsdk:"data_process_region" autogen:"omitjsonupdate"`
	GroupId           types.String `tfsdk:"group_id" autogen:"omitjson"`
	Hostnames         types.List   `tfsdk:"hostnames" autogen:"omitjson"`
	Links             types.List   `tfsdk:"links" autogen:"omitjson"`
	Name              types.String `tfsdk:"name" autogen:"omitjsonupdate"`
	Region            types.String `tfsdk:"region"`
	SampleConnections types.Object `tfsdk:"sample_connections" autogen:"omitjsonupdate"`
	StreamConfig      types.Object `tfsdk:"stream_config" autogen:"omitjsonupdate"`
}
type TFConnectionsModel struct {
	Authentication   types.Object `tfsdk:"authentication" autogen:"omitjson"`
	Aws              types.Object `tfsdk:"aws" autogen:"omitjson"`
	BootstrapServers types.String `tfsdk:"bootstrap_servers" autogen:"omitjson"`
	ClusterGroupId   types.String `tfsdk:"cluster_group_id" autogen:"omitjson"`
	ClusterName      types.String `tfsdk:"cluster_name" autogen:"omitjson"`
	Config           types.Map    `tfsdk:"config" autogen:"omitjson"`
	DbRoleToExecute  types.Object `tfsdk:"db_role_to_execute" autogen:"omitjson"`
	Headers          types.Map    `tfsdk:"headers" autogen:"omitjson"`
	Links            types.List   `tfsdk:"links" autogen:"omitjson"`
	Name             types.String `tfsdk:"name" autogen:"omitjson"`
	Networking       types.Object `tfsdk:"networking" autogen:"omitjson"`
	Security         types.Object `tfsdk:"security" autogen:"omitjson"`
	Type             types.String `tfsdk:"type" autogen:"omitjson"`
	Url              types.String `tfsdk:"url" autogen:"omitjson"`
}
type TFConnectionsAuthenticationModel struct {
	Links          types.List   `tfsdk:"links" autogen:"omitjson"`
	Mechanism      types.String `tfsdk:"mechanism" autogen:"omitjson"`
	Password       types.String `tfsdk:"password" autogen:"omitjson"`
	SslCertificate types.String `tfsdk:"ssl_certificate" autogen:"omitjson"`
	SslKey         types.String `tfsdk:"ssl_key" autogen:"omitjson"`
	SslKeyPassword types.String `tfsdk:"ssl_key_password" autogen:"omitjson"`
	Username       types.String `tfsdk:"username" autogen:"omitjson"`
}
type TFConnectionsAuthenticationLinksModel struct {
	Href types.String `tfsdk:"href" autogen:"omitjson"`
	Rel  types.String `tfsdk:"rel" autogen:"omitjson"`
}
type TFConnectionsAwsModel struct {
	Links      types.List   `tfsdk:"links" autogen:"omitjson"`
	RoleArn    types.String `tfsdk:"role_arn" autogen:"omitjson"`
	TestBucket types.String `tfsdk:"test_bucket" autogen:"omitjson"`
}
type TFConnectionsAwsLinksModel struct {
	Href types.String `tfsdk:"href" autogen:"omitjson"`
	Rel  types.String `tfsdk:"rel" autogen:"omitjson"`
}
type TFConnectionsDbRoleToExecuteModel struct {
	Links types.List   `tfsdk:"links" autogen:"omitjson"`
	Role  types.String `tfsdk:"role" autogen:"omitjson"`
	Type  types.String `tfsdk:"type" autogen:"omitjson"`
}
type TFConnectionsDbRoleToExecuteLinksModel struct {
	Href types.String `tfsdk:"href" autogen:"omitjson"`
	Rel  types.String `tfsdk:"rel" autogen:"omitjson"`
}
type TFConnectionsLinksModel struct {
	Href types.String `tfsdk:"href" autogen:"omitjson"`
	Rel  types.String `tfsdk:"rel" autogen:"omitjson"`
}
type TFConnectionsNetworkingModel struct {
	Access types.Object `tfsdk:"access" autogen:"omitjson"`
	Links  types.List   `tfsdk:"links" autogen:"omitjson"`
}
type TFConnectionsNetworkingAccessModel struct {
	ConnectionId types.String `tfsdk:"connection_id" autogen:"omitjson"`
	Links        types.List   `tfsdk:"links" autogen:"omitjson"`
	Name         types.String `tfsdk:"name" autogen:"omitjson"`
	TgwId        types.String `tfsdk:"tgw_id" autogen:"omitjson"`
	TgwRouteId   types.String `tfsdk:"tgw_route_id" autogen:"omitjson"`
	Type         types.String `tfsdk:"type" autogen:"omitjson"`
	VpcCidr      types.String `tfsdk:"vpc_cidr" autogen:"omitjson"`
}
type TFConnectionsNetworkingAccessLinksModel struct {
	Href types.String `tfsdk:"href" autogen:"omitjson"`
	Rel  types.String `tfsdk:"rel" autogen:"omitjson"`
}
type TFConnectionsNetworkingLinksModel struct {
	Href types.String `tfsdk:"href" autogen:"omitjson"`
	Rel  types.String `tfsdk:"rel" autogen:"omitjson"`
}
type TFConnectionsSecurityModel struct {
	BrokerPublicCertificate types.String `tfsdk:"broker_public_certificate" autogen:"omitjson"`
	Links                   types.List   `tfsdk:"links" autogen:"omitjson"`
	Protocol                types.String `tfsdk:"protocol" autogen:"omitjson"`
}
type TFConnectionsSecurityLinksModel struct {
	Href types.String `tfsdk:"href" autogen:"omitjson"`
	Rel  types.String `tfsdk:"rel" autogen:"omitjson"`
}
type TFDataProcessRegionModel struct {
	CloudProvider types.String `tfsdk:"cloud_provider"`
	Links         types.List   `tfsdk:"links" autogen:"omitjson"`
	Region        types.String `tfsdk:"region"`
}
type TFDataProcessRegionLinksModel struct {
	Href types.String `tfsdk:"href" autogen:"omitjson"`
	Rel  types.String `tfsdk:"rel" autogen:"omitjson"`
}
type TFLinksModel struct {
	Href types.String `tfsdk:"href" autogen:"omitjson"`
	Rel  types.String `tfsdk:"rel" autogen:"omitjson"`
}
type TFSampleConnectionsModel struct {
	Links types.List `tfsdk:"links" autogen:"omitjson"`
	Solar types.Bool `tfsdk:"solar"`
}
type TFSampleConnectionsLinksModel struct {
	Href types.String `tfsdk:"href" autogen:"omitjson"`
	Rel  types.String `tfsdk:"rel" autogen:"omitjson"`
}
type TFStreamConfigModel struct {
	Links types.List   `tfsdk:"links" autogen:"omitjson"`
	Tier  types.String `tfsdk:"tier"`
}
type TFStreamConfigLinksModel struct {
	Href types.String `tfsdk:"href" autogen:"omitjson"`
	Rel  types.String `tfsdk:"rel" autogen:"omitjson"`
}
