package advancedclustertpf

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"accept_data_risks_and_force_replica_set_reconfig": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "If reconfiguration is necessary to regain a primary due to a regional outage, submit this field alongside your topology reconfiguration to request a new regional outage resistant topology. Forced reconfigurations during an outage of the majority of electable nodes carry a risk of data loss if replicated writes (even majority committed writes) have not been replicated to the new primary node. MongoDB Atlas docs contain more information. To proceed with an operation which carries that risk, set **acceptDataRisksAndForceReplicaSetReconfig** to the current date.",
				MarkdownDescription: "If reconfiguration is necessary to regain a primary due to a regional outage, submit this field alongside your topology reconfiguration to request a new regional outage resistant topology. Forced reconfigurations during an outage of the majority of electable nodes carry a risk of data loss if replicated writes (even majority committed writes) have not been replicated to the new primary node. MongoDB Atlas docs contain more information. To proceed with an operation which carries that risk, set **acceptDataRisksAndForceReplicaSetReconfig** to the current date.",
			},
			"backup_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Flag that indicates whether the cluster can perform backups. If set to `true`, the cluster can perform backups. You must set this value to `true` for NVMe clusters. Backup uses [Cloud Backups](https://docs.atlas.mongodb.com/backup/cloud-backup/overview/) for dedicated clusters and [Shared Cluster Backups](https://docs.atlas.mongodb.com/backup/shared-tier/overview/) for tenant clusters. If set to `false`, the cluster doesn't use backups.",
				MarkdownDescription: "Flag that indicates whether the cluster can perform backups. If set to `true`, the cluster can perform backups. You must set this value to `true` for NVMe clusters. Backup uses [Cloud Backups](https://docs.atlas.mongodb.com/backup/cloud-backup/overview/) for dedicated clusters and [Shared Cluster Backups](https://docs.atlas.mongodb.com/backup/shared-tier/overview/) for tenant clusters. If set to `false`, the cluster doesn't use backups.",
				Default:             booldefault.StaticBool(false),
			},
			"bi_connector": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Description:         "Flag that indicates whether MongoDB Connector for Business Intelligence is enabled on the specified cluster.",
						MarkdownDescription: "Flag that indicates whether MongoDB Connector for Business Intelligence is enabled on the specified cluster.",
					},
					"read_preference": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Description:         "Data source node designated for the MongoDB Connector for Business Intelligence on MongoDB Cloud. The MongoDB Connector for Business Intelligence on MongoDB Cloud reads data from the primary, secondary, or analytics node based on your read preferences. Defaults to `ANALYTICS` node, or `SECONDARY` if there are no `ANALYTICS` nodes.",
						MarkdownDescription: "Data source node designated for the MongoDB Connector for Business Intelligence on MongoDB Cloud. The MongoDB Connector for Business Intelligence on MongoDB Cloud reads data from the primary, secondary, or analytics node based on your read preferences. Defaults to `ANALYTICS` node, or `SECONDARY` if there are no `ANALYTICS` nodes.",
					},
				},
				CustomType: BiConnectorType{
					ObjectType: types.ObjectType{
						AttrTypes: BiConnectorValue{}.AttributeTypes(ctx),
					},
				},
				Optional:            true,
				Computed:            true,
				Description:         "Settings needed to configure the MongoDB Connector for Business Intelligence for this cluster.",
				MarkdownDescription: "Settings needed to configure the MongoDB Connector for Business Intelligence for this cluster.",
			},
			"cluster_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Human-readable label that identifies this cluster.",
				MarkdownDescription: "Human-readable label that identifies this cluster.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
					stringvalidator.RegexMatches(regexp.MustCompile("^([a-zA-Z0-9][a-zA-Z0-9-]*)?[a-zA-Z0-9]+$"), ""),
				},
			},
			"cluster_type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Configuration of nodes that comprise the cluster.",
				MarkdownDescription: "Configuration of nodes that comprise the cluster.",
			},
			"config_server_management_mode": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Config Server Management Mode for creating or updating a sharded cluster.\n\nWhen configured as ATLAS_MANAGED, atlas may automatically switch the cluster's config server type for optimal performance and savings.\n\nWhen configured as FIXED_TO_DEDICATED, the cluster will always use a dedicated config server.",
				MarkdownDescription: "Config Server Management Mode for creating or updating a sharded cluster.\n\nWhen configured as ATLAS_MANAGED, atlas may automatically switch the cluster's config server type for optimal performance and savings.\n\nWhen configured as FIXED_TO_DEDICATED, the cluster will always use a dedicated config server.",
				Default:             stringdefault.StaticString("ATLAS_MANAGED"),
			},
			"config_server_type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Describes a sharded cluster's config server type.",
				MarkdownDescription: "Describes a sharded cluster's config server type.",
			},
			"connection_strings": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"aws_private_link": schema.MapAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						Computed:            true,
						Description:         "Private endpoint-aware connection strings that use AWS-hosted clusters with Amazon Web Services (AWS) PrivateLink. Each key identifies an Amazon Web Services (AWS) interface endpoint. Each value identifies the related `mongodb://` connection string that you use to connect to MongoDB Cloud through the interface endpoint that the key names.",
						MarkdownDescription: "Private endpoint-aware connection strings that use AWS-hosted clusters with Amazon Web Services (AWS) PrivateLink. Each key identifies an Amazon Web Services (AWS) interface endpoint. Each value identifies the related `mongodb://` connection string that you use to connect to MongoDB Cloud through the interface endpoint that the key names.",
					},
					"aws_private_link_srv": schema.MapAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						Computed:            true,
						Description:         "Private endpoint-aware connection strings that use AWS-hosted clusters with Amazon Web Services (AWS) PrivateLink. Each key identifies an Amazon Web Services (AWS) interface endpoint. Each value identifies the related `mongodb://` connection string that you use to connect to Atlas through the interface endpoint that the key names.",
						MarkdownDescription: "Private endpoint-aware connection strings that use AWS-hosted clusters with Amazon Web Services (AWS) PrivateLink. Each key identifies an Amazon Web Services (AWS) interface endpoint. Each value identifies the related `mongodb://` connection string that you use to connect to Atlas through the interface endpoint that the key names.",
					},
					"private": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Description:         "Network peering connection strings for each interface Virtual Private Cloud (VPC) endpoint that you configured to connect to this cluster. This connection string uses the `mongodb+srv://` protocol. The resource returns this parameter once someone creates a network peering connection to this cluster. This protocol tells the application to look up the host seed list in the Domain Name System (DNS). This list synchronizes with the nodes in a cluster. If the connection string uses this Uniform Resource Identifier (URI) format, you don't need to append the seed list or change the URI if the nodes change. Use this URI format if your driver supports it. If it doesn't, use connectionStrings.private. For Amazon Web Services (AWS) clusters, this resource returns this parameter only if you enable custom DNS.",
						MarkdownDescription: "Network peering connection strings for each interface Virtual Private Cloud (VPC) endpoint that you configured to connect to this cluster. This connection string uses the `mongodb+srv://` protocol. The resource returns this parameter once someone creates a network peering connection to this cluster. This protocol tells the application to look up the host seed list in the Domain Name System (DNS). This list synchronizes with the nodes in a cluster. If the connection string uses this Uniform Resource Identifier (URI) format, you don't need to append the seed list or change the URI if the nodes change. Use this URI format if your driver supports it. If it doesn't, use connectionStrings.private. For Amazon Web Services (AWS) clusters, this resource returns this parameter only if you enable custom DNS.",
					},
					"private_endpoint": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"connection_string": schema.StringAttribute{
									Optional:            true,
									Computed:            true,
									Description:         "Private endpoint-aware connection string that uses the `mongodb://` protocol to connect to MongoDB Cloud through a private endpoint.",
									MarkdownDescription: "Private endpoint-aware connection string that uses the `mongodb://` protocol to connect to MongoDB Cloud through a private endpoint.",
								},
								"endpoints": schema.ListNestedAttribute{
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"endpoint_id": schema.StringAttribute{
												Optional:            true,
												Computed:            true,
												Description:         "Unique string that the cloud provider uses to identify the private endpoint.",
												MarkdownDescription: "Unique string that the cloud provider uses to identify the private endpoint.",
											},
											"provider_name": schema.StringAttribute{
												Optional:            true,
												Computed:            true,
												Description:         "Cloud provider in which MongoDB Cloud deploys the private endpoint.",
												MarkdownDescription: "Cloud provider in which MongoDB Cloud deploys the private endpoint.",
											},
											"region": schema.StringAttribute{
												Optional:            true,
												Computed:            true,
												Description:         "Region where the private endpoint is deployed.",
												MarkdownDescription: "Region where the private endpoint is deployed.",
											},
										},
										CustomType: EndpointsType{
											ObjectType: types.ObjectType{
												AttrTypes: EndpointsValue{}.AttributeTypes(ctx),
											},
										},
									},
									Optional:            true,
									Computed:            true,
									Description:         "List that contains the private endpoints through which you connect to MongoDB Cloud when you use **connectionStrings.privateEndpoint[n].connectionString** or **connectionStrings.privateEndpoint[n].srvConnectionString**.",
									MarkdownDescription: "List that contains the private endpoints through which you connect to MongoDB Cloud when you use **connectionStrings.privateEndpoint[n].connectionString** or **connectionStrings.privateEndpoint[n].srvConnectionString**.",
								},
								"srv_connection_string": schema.StringAttribute{
									Optional:            true,
									Computed:            true,
									Description:         "Private endpoint-aware connection string that uses the `mongodb+srv://` protocol to connect to MongoDB Cloud through a private endpoint. The `mongodb+srv` protocol tells the driver to look up the seed list of hosts in the Domain Name System (DNS). This list synchronizes with the nodes in a cluster. If the connection string uses this Uniform Resource Identifier (URI) format, you don't need to append the seed list or change the Uniform Resource Identifier (URI) if the nodes change. Use this Uniform Resource Identifier (URI) format if your application supports it. If it doesn't, use connectionStrings.privateEndpoint[n].connectionString.",
									MarkdownDescription: "Private endpoint-aware connection string that uses the `mongodb+srv://` protocol to connect to MongoDB Cloud through a private endpoint. The `mongodb+srv` protocol tells the driver to look up the seed list of hosts in the Domain Name System (DNS). This list synchronizes with the nodes in a cluster. If the connection string uses this Uniform Resource Identifier (URI) format, you don't need to append the seed list or change the Uniform Resource Identifier (URI) if the nodes change. Use this Uniform Resource Identifier (URI) format if your application supports it. If it doesn't, use connectionStrings.privateEndpoint[n].connectionString.",
								},
								"srv_shard_optimized_connection_string": schema.StringAttribute{
									Optional:            true,
									Computed:            true,
									Description:         "Private endpoint-aware connection string optimized for sharded clusters that uses the `mongodb+srv://` protocol to connect to MongoDB Cloud through a private endpoint. If the connection string uses this Uniform Resource Identifier (URI) format, you don't need to change the Uniform Resource Identifier (URI) if the nodes change. Use this Uniform Resource Identifier (URI) format if your application and Atlas cluster supports it. If it doesn't, use and consult the documentation for connectionStrings.privateEndpoint[n].srvConnectionString.",
									MarkdownDescription: "Private endpoint-aware connection string optimized for sharded clusters that uses the `mongodb+srv://` protocol to connect to MongoDB Cloud through a private endpoint. If the connection string uses this Uniform Resource Identifier (URI) format, you don't need to change the Uniform Resource Identifier (URI) if the nodes change. Use this Uniform Resource Identifier (URI) format if your application and Atlas cluster supports it. If it doesn't, use and consult the documentation for connectionStrings.privateEndpoint[n].srvConnectionString.",
								},
								"type": schema.StringAttribute{
									Optional:            true,
									Computed:            true,
									Description:         "MongoDB process type to which your application connects. Use `MONGOD` for replica sets and `MONGOS` for sharded clusters.",
									MarkdownDescription: "MongoDB process type to which your application connects. Use `MONGOD` for replica sets and `MONGOS` for sharded clusters.",
								},
							},
							CustomType: PrivateEndpointType{
								ObjectType: types.ObjectType{
									AttrTypes: PrivateEndpointValue{}.AttributeTypes(ctx),
								},
							},
						},
						Optional:            true,
						Computed:            true,
						Description:         "List of private endpoint-aware connection strings that you can use to connect to this cluster through a private endpoint. This parameter returns only if you deployed a private endpoint to all regions to which you deployed this clusters' nodes.",
						MarkdownDescription: "List of private endpoint-aware connection strings that you can use to connect to this cluster through a private endpoint. This parameter returns only if you deployed a private endpoint to all regions to which you deployed this clusters' nodes.",
					},
					"private_srv": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Description:         "Network peering connection strings for each interface Virtual Private Cloud (VPC) endpoint that you configured to connect to this cluster. This connection string uses the `mongodb+srv://` protocol. The resource returns this parameter when someone creates a network peering connection to this cluster. This protocol tells the application to look up the host seed list in the Domain Name System (DNS). This list synchronizes with the nodes in a cluster. If the connection string uses this Uniform Resource Identifier (URI) format, you don't need to append the seed list or change the Uniform Resource Identifier (URI) if the nodes change. Use this Uniform Resource Identifier (URI) format if your driver supports it. If it doesn't, use `connectionStrings.private`. For Amazon Web Services (AWS) clusters, this parameter returns only if you [enable custom DNS](https://docs.atlas.mongodb.com/reference/api/aws-custom-dns-update/).",
						MarkdownDescription: "Network peering connection strings for each interface Virtual Private Cloud (VPC) endpoint that you configured to connect to this cluster. This connection string uses the `mongodb+srv://` protocol. The resource returns this parameter when someone creates a network peering connection to this cluster. This protocol tells the application to look up the host seed list in the Domain Name System (DNS). This list synchronizes with the nodes in a cluster. If the connection string uses this Uniform Resource Identifier (URI) format, you don't need to append the seed list or change the Uniform Resource Identifier (URI) if the nodes change. Use this Uniform Resource Identifier (URI) format if your driver supports it. If it doesn't, use `connectionStrings.private`. For Amazon Web Services (AWS) clusters, this parameter returns only if you [enable custom DNS](https://docs.atlas.mongodb.com/reference/api/aws-custom-dns-update/).",
					},
					"standard": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Description:         "Public connection string that you can use to connect to this cluster. This connection string uses the `mongodb://` protocol.",
						MarkdownDescription: "Public connection string that you can use to connect to this cluster. This connection string uses the `mongodb://` protocol.",
					},
					"standard_srv": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Description:         "Public connection string that you can use to connect to this cluster. This connection string uses the `mongodb+srv://` protocol.",
						MarkdownDescription: "Public connection string that you can use to connect to this cluster. This connection string uses the `mongodb+srv://` protocol.",
					},
				},
				CustomType: ConnectionStringsType{
					ObjectType: types.ObjectType{
						AttrTypes: ConnectionStringsValue{}.AttributeTypes(ctx),
					},
				},
				Optional:            true,
				Computed:            true,
				Description:         "Collection of Uniform Resource Locators that point to the MongoDB database.",
				MarkdownDescription: "Collection of Uniform Resource Locators that point to the MongoDB database.",
			},
			"create_date": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Date and time when MongoDB Cloud created this cluster. This parameter expresses its value in ISO 8601 format in UTC.",
				MarkdownDescription: "Date and time when MongoDB Cloud created this cluster. This parameter expresses its value in ISO 8601 format in UTC.",
			},
			"disk_warming_mode": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Disk warming mode selection.",
				MarkdownDescription: "Disk warming mode selection.",
				Default:             stringdefault.StaticString("FULLY_WARMED"),
			},
			"encryption_at_rest_provider": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Cloud service provider that manages your customer keys to provide an additional layer of encryption at rest for the cluster. To enable customer key management for encryption at rest, the cluster **replicationSpecs[n].regionConfigs[m].{type}Specs.instanceSize** setting must be `M10` or higher and `\"backupEnabled\" : false` or omitted entirely.",
				MarkdownDescription: "Cloud service provider that manages your customer keys to provide an additional layer of encryption at rest for the cluster. To enable customer key management for encryption at rest, the cluster **replicationSpecs[n].regionConfigs[m].{type}Specs.instanceSize** setting must be `M10` or higher and `\"backupEnabled\" : false` or omitted entirely.",
			},
			"feature_compatibility_version": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Feature compatibility version of the cluster.",
				MarkdownDescription: "Feature compatibility version of the cluster.",
			},
			"feature_compatibility_version_expiration_date": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Feature compatibility version expiration date.",
				MarkdownDescription: "Feature compatibility version expiration date.",
			},
			"global_cluster_self_managed_sharding": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Set this field to configure the Sharding Management Mode when creating a new Global Cluster.\n\nWhen set to false, the management mode is set to Atlas-Managed Sharding. This mode fully manages the sharding of your Global Cluster and is built to provide a seamless deployment experience.\n\nWhen set to true, the management mode is set to Self-Managed Sharding. This mode leaves the management of shards in your hands and is built to provide an advanced and flexible deployment experience.\n\nThis setting cannot be changed once the cluster is deployed.",
				MarkdownDescription: "Set this field to configure the Sharding Management Mode when creating a new Global Cluster.\n\nWhen set to false, the management mode is set to Atlas-Managed Sharding. This mode fully manages the sharding of your Global Cluster and is built to provide a seamless deployment experience.\n\nWhen set to true, the management mode is set to Self-Managed Sharding. This mode leaves the management of shards in your hands and is built to provide an advanced and flexible deployment experience.\n\nThis setting cannot be changed once the cluster is deployed.",
			},
			"group_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Unique 24-hexadecimal character string that identifies the project.",
				MarkdownDescription: "Unique 24-hexadecimal character string that identifies the project.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(24, 24),
					stringvalidator.RegexMatches(regexp.MustCompile("^([a-f0-9]{24})$"), ""),
				},
			},
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Unique 24-hexadecimal digit string that identifies the cluster.",
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the cluster.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(24, 24),
					stringvalidator.RegexMatches(regexp.MustCompile("^([a-f0-9]{24})$"), ""),
				},
			},
			"labels": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "Key applied to tag and categorize this component.",
							MarkdownDescription: "Key applied to tag and categorize this component.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 255),
							},
						},
						"value": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "Value set to the Key applied to tag and categorize this component.",
							MarkdownDescription: "Value set to the Key applied to tag and categorize this component.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 255),
							},
						},
					},
					CustomType: LabelsType{
						ObjectType: types.ObjectType{
							AttrTypes: LabelsValue{}.AttributeTypes(ctx),
						},
					},
				},
				Optional:            true,
				Computed:            true,
				Description:         "Collection of key-value pairs between 1 to 255 characters in length that tag and categorize the cluster. The MongoDB Cloud console doesn't display your labels.\n\nCluster labels are deprecated and will be removed in a future release. We strongly recommend that you use [resource tags](https://dochub.mongodb.org/core/add-cluster-tag-atlas) instead.",
				MarkdownDescription: "Collection of key-value pairs between 1 to 255 characters in length that tag and categorize the cluster. The MongoDB Cloud console doesn't display your labels.\n\nCluster labels are deprecated and will be removed in a future release. We strongly recommend that you use [resource tags](https://dochub.mongodb.org/core/add-cluster-tag-atlas) instead.",
				DeprecationMessage:  "This attribute is deprecated.",
			},
			"links": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"href": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
							MarkdownDescription: "Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
						},
						"rel": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "Uniform Resource Locator (URL) that defines the semantic relationship between this resource and another API resource. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
							MarkdownDescription: "Uniform Resource Locator (URL) that defines the semantic relationship between this resource and another API resource. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
						},
					},
					CustomType: LinksType{
						ObjectType: types.ObjectType{
							AttrTypes: LinksValue{}.AttributeTypes(ctx),
						},
					},
				},
				Optional:            true,
				Computed:            true,
				Description:         "List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.",
				MarkdownDescription: "List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.",
			},
			"mongo_dbemployee_access_grant": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"expiration_time": schema.StringAttribute{
						Required:            true,
						Description:         "Expiration date for the employee access grant.",
						MarkdownDescription: "Expiration date for the employee access grant.",
					},
					"grant_type": schema.StringAttribute{
						Required:            true,
						Description:         "Level of access to grant to MongoDB Employees.",
						MarkdownDescription: "Level of access to grant to MongoDB Employees.",
					},
					"links": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"href": schema.StringAttribute{
									Optional:            true,
									Computed:            true,
									Description:         "Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
									MarkdownDescription: "Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
								},
								"rel": schema.StringAttribute{
									Optional:            true,
									Computed:            true,
									Description:         "Uniform Resource Locator (URL) that defines the semantic relationship between this resource and another API resource. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
									MarkdownDescription: "Uniform Resource Locator (URL) that defines the semantic relationship between this resource and another API resource. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
								},
							},
							CustomType: LinksType{
								ObjectType: types.ObjectType{
									AttrTypes: LinksValue{}.AttributeTypes(ctx),
								},
							},
						},
						Optional:            true,
						Computed:            true,
						Description:         "List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.",
						MarkdownDescription: "List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.",
					},
				},
				CustomType: MongoDbemployeeAccessGrantType{
					ObjectType: types.ObjectType{
						AttrTypes: MongoDbemployeeAccessGrantValue{}.AttributeTypes(ctx),
					},
				},
				Optional:            true,
				Computed:            true,
				Description:         "MongoDB employee granted access level and expiration for a cluster.",
				MarkdownDescription: "MongoDB employee granted access level and expiration for a cluster.",
			},
			"mongo_dbmajor_version": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "MongoDB major version of the cluster.\n\nOn creation: Choose from the available versions of MongoDB, or leave unspecified for the current recommended default in the MongoDB Cloud platform. The recommended version is a recent Long Term Support version. The default is not guaranteed to be the most recently released version throughout the entire release cycle. For versions available in a specific project, see the linked documentation or use the API endpoint for [project LTS versions endpoint](#tag/Projects/operation/getProjectLTSVersions).\n\n On update: Increase version only by 1 major version at a time. If the cluster is pinned to a MongoDB feature compatibility version exactly one major version below the current MongoDB version, the MongoDB version can be downgraded to the previous major version.",
				MarkdownDescription: "MongoDB major version of the cluster.\n\nOn creation: Choose from the available versions of MongoDB, or leave unspecified for the current recommended default in the MongoDB Cloud platform. The recommended version is a recent Long Term Support version. The default is not guaranteed to be the most recently released version throughout the entire release cycle. For versions available in a specific project, see the linked documentation or use the API endpoint for [project LTS versions endpoint](#tag/Projects/operation/getProjectLTSVersions).\n\n On update: Increase version only by 1 major version at a time. If the cluster is pinned to a MongoDB feature compatibility version exactly one major version below the current MongoDB version, the MongoDB version can be downgraded to the previous major version.",
			},
			"mongo_dbversion": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Version of MongoDB that the cluster runs.",
				MarkdownDescription: "Version of MongoDB that the cluster runs.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile("([\\d]+\\.[\\d]+\\.[\\d]+)"), ""),
				},
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Human-readable label that identifies the cluster.",
				MarkdownDescription: "Human-readable label that identifies the cluster.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
					stringvalidator.RegexMatches(regexp.MustCompile("^([a-zA-Z0-9][a-zA-Z0-9-]*)?[a-zA-Z0-9]+$"), ""),
				},
			},
			"paused": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Flag that indicates whether the cluster is paused.",
				MarkdownDescription: "Flag that indicates whether the cluster is paused.",
			},
			"pit_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Flag that indicates whether the cluster uses continuous cloud backups.",
				MarkdownDescription: "Flag that indicates whether the cluster uses continuous cloud backups.",
			},
			"redact_client_log_data": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Enable or disable log redaction.\n\nThis setting configures the ``mongod`` or ``mongos`` to redact any document field contents from a message accompanying a given log event before logging. This prevents the program from writing potentially sensitive data stored on the database to the diagnostic log. Metadata such as error or operation codes, line numbers, and source file names are still visible in the logs.\n\nUse ``redactClientLogData`` in conjunction with Encryption at Rest and TLS/SSL (Transport Encryption) to assist compliance with regulatory requirements.\n\n*Note*: changing this setting on a cluster will trigger a rolling restart as soon as the cluster is updated.",
				MarkdownDescription: "Enable or disable log redaction.\n\nThis setting configures the ``mongod`` or ``mongos`` to redact any document field contents from a message accompanying a given log event before logging. This prevents the program from writing potentially sensitive data stored on the database to the diagnostic log. Metadata such as error or operation codes, line numbers, and source file names are still visible in the logs.\n\nUse ``redactClientLogData`` in conjunction with Encryption at Rest and TLS/SSL (Transport Encryption) to assist compliance with regulatory requirements.\n\n*Note*: changing this setting on a cluster will trigger a rolling restart as soon as the cluster is updated.",
			},
			"replica_set_scaling_strategy": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Set this field to configure the replica set scaling mode for your cluster.\n\nBy default, Atlas scales under WORKLOAD_TYPE. This mode allows Atlas to scale your analytics nodes in parallel to your operational nodes.\n\nWhen configured as SEQUENTIAL, Atlas scales all nodes sequentially. This mode is intended for steady-state workloads and applications performing latency-sensitive secondary reads.\n\nWhen configured as NODE_TYPE, Atlas scales your electable nodes in parallel with your read-only and analytics nodes. This mode is intended for large, dynamic workloads requiring frequent and timely cluster tier scaling. This is the fastest scaling strategy, but it might impact latency of workloads when performing extensive secondary reads.",
				MarkdownDescription: "Set this field to configure the replica set scaling mode for your cluster.\n\nBy default, Atlas scales under WORKLOAD_TYPE. This mode allows Atlas to scale your analytics nodes in parallel to your operational nodes.\n\nWhen configured as SEQUENTIAL, Atlas scales all nodes sequentially. This mode is intended for steady-state workloads and applications performing latency-sensitive secondary reads.\n\nWhen configured as NODE_TYPE, Atlas scales your electable nodes in parallel with your read-only and analytics nodes. This mode is intended for large, dynamic workloads requiring frequent and timely cluster tier scaling. This is the fastest scaling strategy, but it might impact latency of workloads when performing extensive secondary reads.",
				Default:             stringdefault.StaticString("WORKLOAD_TYPE"),
			},
			"replication_specs": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "Unique 24-hexadecimal digit string that identifies the replication object for a shard in a Cluster. If you include existing shard replication configurations in the request, you must specify this parameter. If you add a new shard to an existing Cluster, you may specify this parameter. The request deletes any existing shards  in the Cluster that you exclude from the request. This corresponds to Shard ID displayed in the UI.",
							MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the replication object for a shard in a Cluster. If you include existing shard replication configurations in the request, you must specify this parameter. If you add a new shard to an existing Cluster, you may specify this parameter. The request deletes any existing shards  in the Cluster that you exclude from the request. This corresponds to Shard ID displayed in the UI.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(24, 24),
								stringvalidator.RegexMatches(regexp.MustCompile("^([a-f0-9]{24})$"), ""),
							},
						},
						"region_configs": schema.ListNestedAttribute{
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"analytics_auto_scaling": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"compute": schema.SingleNestedAttribute{
												Attributes: map[string]schema.Attribute{
													"enabled": schema.BoolAttribute{
														Optional:            true,
														Computed:            true,
														Description:         "Flag that indicates whether someone enabled instance size auto-scaling.\n\n- Set to `true` to enable instance size auto-scaling. If enabled, you must specify a value for **replicationSpecs[n].regionConfigs[m].autoScaling.compute.maxInstanceSize**.\n- Set to `false` to disable instance size automatic scaling.",
														MarkdownDescription: "Flag that indicates whether someone enabled instance size auto-scaling.\n\n- Set to `true` to enable instance size auto-scaling. If enabled, you must specify a value for **replicationSpecs[n].regionConfigs[m].autoScaling.compute.maxInstanceSize**.\n- Set to `false` to disable instance size automatic scaling.",
													},
													"max_instance_size": schema.StringAttribute{
														Optional:            true,
														Computed:            true,
														Description:         "Minimum instance size to which your cluster can automatically scale. MongoDB Cloud requires this parameter if `\"replicationSpecs[n].regionConfigs[m].autoScaling.compute.scaleDownEnabled\" : true`.",
														MarkdownDescription: "Minimum instance size to which your cluster can automatically scale. MongoDB Cloud requires this parameter if `\"replicationSpecs[n].regionConfigs[m].autoScaling.compute.scaleDownEnabled\" : true`.",
													},
													"min_instance_size": schema.StringAttribute{
														Optional:            true,
														Computed:            true,
														Description:         "Minimum instance size to which your cluster can automatically scale. MongoDB Cloud requires this parameter if `\"replicationSpecs[n].regionConfigs[m].autoScaling.compute.scaleDownEnabled\" : true`.",
														MarkdownDescription: "Minimum instance size to which your cluster can automatically scale. MongoDB Cloud requires this parameter if `\"replicationSpecs[n].regionConfigs[m].autoScaling.compute.scaleDownEnabled\" : true`.",
													},
													"scale_down_enabled": schema.BoolAttribute{
														Optional:            true,
														Computed:            true,
														Description:         "Flag that indicates whether the instance size may scale down. MongoDB Cloud requires this parameter if `\"replicationSpecs[n].regionConfigs[m].autoScaling.compute.enabled\" : true`. If you enable this option, specify a value for **replicationSpecs[n].regionConfigs[m].autoScaling.compute.minInstanceSize**.",
														MarkdownDescription: "Flag that indicates whether the instance size may scale down. MongoDB Cloud requires this parameter if `\"replicationSpecs[n].regionConfigs[m].autoScaling.compute.enabled\" : true`. If you enable this option, specify a value for **replicationSpecs[n].regionConfigs[m].autoScaling.compute.minInstanceSize**.",
													},
												},
												CustomType: ComputeType{
													ObjectType: types.ObjectType{
														AttrTypes: ComputeValue{}.AttributeTypes(ctx),
													},
												},
												Optional:            true,
												Computed:            true,
												Description:         "Options that determine how this cluster handles CPU scaling.",
												MarkdownDescription: "Options that determine how this cluster handles CPU scaling.",
											},
											"disk_gb": schema.SingleNestedAttribute{
												Attributes: map[string]schema.Attribute{
													"enabled": schema.BoolAttribute{
														Optional:            true,
														Computed:            true,
														Description:         "Flag that indicates whether this cluster enables disk auto-scaling. The maximum memory allowed for the selected cluster tier and the oplog size can limit storage auto-scaling.",
														MarkdownDescription: "Flag that indicates whether this cluster enables disk auto-scaling. The maximum memory allowed for the selected cluster tier and the oplog size can limit storage auto-scaling.",
													},
												},
												CustomType: DiskGbType{
													ObjectType: types.ObjectType{
														AttrTypes: DiskGbValue{}.AttributeTypes(ctx),
													},
												},
												Optional:            true,
												Computed:            true,
												Description:         "Setting that enables disk auto-scaling.",
												MarkdownDescription: "Setting that enables disk auto-scaling.",
											},
										},
										CustomType: AnalyticsAutoScalingType{
											ObjectType: types.ObjectType{
												AttrTypes: AnalyticsAutoScalingValue{}.AttributeTypes(ctx),
											},
										},
										Optional:            true,
										Computed:            true,
										Description:         "Options that determine how this cluster handles resource scaling.",
										MarkdownDescription: "Options that determine how this cluster handles resource scaling.",
									},
									"analytics_specs": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"disk_iops": schema.Int64Attribute{
												Optional:            true,
												Computed:            true,
												Description:         "Target throughput desired for storage attached to your Azure-provisioned cluster. Change this parameter if you:\n\n- set `\"replicationSpecs[n].regionConfigs[m].providerName\" : \"Azure\"`.\n- set `\"replicationSpecs[n].regionConfigs[m].electableSpecs.instanceSize\" : \"M40\"` or greater not including `Mxx_NVME` tiers.\n\nThe maximum input/output operations per second (IOPS) depend on the selected **.instanceSize** and **.diskSizeGB**.\nThis parameter defaults to the cluster tier's standard IOPS value.\nChanging this value impacts cluster cost.",
												MarkdownDescription: "Target throughput desired for storage attached to your Azure-provisioned cluster. Change this parameter if you:\n\n- set `\"replicationSpecs[n].regionConfigs[m].providerName\" : \"Azure\"`.\n- set `\"replicationSpecs[n].regionConfigs[m].electableSpecs.instanceSize\" : \"M40\"` or greater not including `Mxx_NVME` tiers.\n\nThe maximum input/output operations per second (IOPS) depend on the selected **.instanceSize** and **.diskSizeGB**.\nThis parameter defaults to the cluster tier's standard IOPS value.\nChanging this value impacts cluster cost.",
											},
											"disk_size_gb": schema.Float64Attribute{
												Optional:            true,
												Computed:            true,
												Description:         "Storage capacity of instance data volumes expressed in gigabytes. Increase this number to add capacity.\n\n This value must be equal for all shards and node types.\n\n This value is not configurable on M0/M2/M5 clusters.\n\n MongoDB Cloud requires this parameter if you set **replicationSpecs**.\n\n If you specify a disk size below the minimum (10 GB), this parameter defaults to the minimum disk size value. \n\n Storage charge calculations depend on whether you choose the default value or a custom value.\n\n The maximum value for disk storage cannot exceed 50 times the maximum RAM for the selected cluster. If you require more storage space, consider upgrading your cluster to a higher tier.",
												MarkdownDescription: "Storage capacity of instance data volumes expressed in gigabytes. Increase this number to add capacity.\n\n This value must be equal for all shards and node types.\n\n This value is not configurable on M0/M2/M5 clusters.\n\n MongoDB Cloud requires this parameter if you set **replicationSpecs**.\n\n If you specify a disk size below the minimum (10 GB), this parameter defaults to the minimum disk size value. \n\n Storage charge calculations depend on whether you choose the default value or a custom value.\n\n The maximum value for disk storage cannot exceed 50 times the maximum RAM for the selected cluster. If you require more storage space, consider upgrading your cluster to a higher tier.",
											},
											"ebs_volume_type": schema.StringAttribute{
												Optional:            true,
												Computed:            true,
												Description:         "Type of storage you want to attach to your AWS-provisioned cluster.\n\n- `STANDARD` volume types can't exceed the default input/output operations per second (IOPS) rate for the selected volume size. \n\n- `PROVISIONED` volume types must fall within the allowable IOPS range for the selected volume size. You must set this value to (`PROVISIONED`) for NVMe clusters.",
												MarkdownDescription: "Type of storage you want to attach to your AWS-provisioned cluster.\n\n- `STANDARD` volume types can't exceed the default input/output operations per second (IOPS) rate for the selected volume size. \n\n- `PROVISIONED` volume types must fall within the allowable IOPS range for the selected volume size. You must set this value to (`PROVISIONED`) for NVMe clusters.",
												Default:             stringdefault.StaticString("STANDARD"),
											},
											"instance_size": schema.StringAttribute{
												Optional:            true,
												Computed:            true,
												Description:         "Hardware specification for the instance sizes in this region in this shard. Each instance size has a default storage and memory capacity. Electable nodes and read-only nodes (known as \"base nodes\") within a single shard must use the same instance size. Analytics nodes can scale independently from base nodes within a shard. Both base nodes and analytics nodes can scale independently from their equivalents in other shards.",
												MarkdownDescription: "Hardware specification for the instance sizes in this region in this shard. Each instance size has a default storage and memory capacity. Electable nodes and read-only nodes (known as \"base nodes\") within a single shard must use the same instance size. Analytics nodes can scale independently from base nodes within a shard. Both base nodes and analytics nodes can scale independently from their equivalents in other shards.",
											},
											"node_count": schema.Int64Attribute{
												Optional:            true,
												Computed:            true,
												Description:         "Number of nodes of the given type for MongoDB Cloud to deploy to the region.",
												MarkdownDescription: "Number of nodes of the given type for MongoDB Cloud to deploy to the region.",
											},
										},
										CustomType: AnalyticsSpecsType{
											ObjectType: types.ObjectType{
												AttrTypes: AnalyticsSpecsValue{}.AttributeTypes(ctx),
											},
										},
										Optional:            true,
										Computed:            true,
										Description:         "Hardware specifications for read-only nodes in the region. Read-only nodes can never become the primary member, but can enable local reads. If you don't specify this parameter, no read-only nodes are deployed to the region.",
										MarkdownDescription: "Hardware specifications for read-only nodes in the region. Read-only nodes can never become the primary member, but can enable local reads. If you don't specify this parameter, no read-only nodes are deployed to the region.",
									},
									"auto_scaling": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"compute": schema.SingleNestedAttribute{
												Attributes: map[string]schema.Attribute{
													"enabled": schema.BoolAttribute{
														Optional:            true,
														Computed:            true,
														Description:         "Flag that indicates whether someone enabled instance size auto-scaling.\n\n- Set to `true` to enable instance size auto-scaling. If enabled, you must specify a value for **replicationSpecs[n].regionConfigs[m].autoScaling.compute.maxInstanceSize**.\n- Set to `false` to disable instance size automatic scaling.",
														MarkdownDescription: "Flag that indicates whether someone enabled instance size auto-scaling.\n\n- Set to `true` to enable instance size auto-scaling. If enabled, you must specify a value for **replicationSpecs[n].regionConfigs[m].autoScaling.compute.maxInstanceSize**.\n- Set to `false` to disable instance size automatic scaling.",
													},
													"max_instance_size": schema.StringAttribute{
														Optional:            true,
														Computed:            true,
														Description:         "Minimum instance size to which your cluster can automatically scale. MongoDB Cloud requires this parameter if `\"replicationSpecs[n].regionConfigs[m].autoScaling.compute.scaleDownEnabled\" : true`.",
														MarkdownDescription: "Minimum instance size to which your cluster can automatically scale. MongoDB Cloud requires this parameter if `\"replicationSpecs[n].regionConfigs[m].autoScaling.compute.scaleDownEnabled\" : true`.",
													},
													"min_instance_size": schema.StringAttribute{
														Optional:            true,
														Computed:            true,
														Description:         "Minimum instance size to which your cluster can automatically scale. MongoDB Cloud requires this parameter if `\"replicationSpecs[n].regionConfigs[m].autoScaling.compute.scaleDownEnabled\" : true`.",
														MarkdownDescription: "Minimum instance size to which your cluster can automatically scale. MongoDB Cloud requires this parameter if `\"replicationSpecs[n].regionConfigs[m].autoScaling.compute.scaleDownEnabled\" : true`.",
													},
													"scale_down_enabled": schema.BoolAttribute{
														Optional:            true,
														Computed:            true,
														Description:         "Flag that indicates whether the instance size may scale down. MongoDB Cloud requires this parameter if `\"replicationSpecs[n].regionConfigs[m].autoScaling.compute.enabled\" : true`. If you enable this option, specify a value for **replicationSpecs[n].regionConfigs[m].autoScaling.compute.minInstanceSize**.",
														MarkdownDescription: "Flag that indicates whether the instance size may scale down. MongoDB Cloud requires this parameter if `\"replicationSpecs[n].regionConfigs[m].autoScaling.compute.enabled\" : true`. If you enable this option, specify a value for **replicationSpecs[n].regionConfigs[m].autoScaling.compute.minInstanceSize**.",
													},
												},
												CustomType: ComputeType{
													ObjectType: types.ObjectType{
														AttrTypes: ComputeValue{}.AttributeTypes(ctx),
													},
												},
												Optional:            true,
												Computed:            true,
												Description:         "Options that determine how this cluster handles CPU scaling.",
												MarkdownDescription: "Options that determine how this cluster handles CPU scaling.",
											},
											"disk_gb": schema.SingleNestedAttribute{
												Attributes: map[string]schema.Attribute{
													"enabled": schema.BoolAttribute{
														Optional:            true,
														Computed:            true,
														Description:         "Flag that indicates whether this cluster enables disk auto-scaling. The maximum memory allowed for the selected cluster tier and the oplog size can limit storage auto-scaling.",
														MarkdownDescription: "Flag that indicates whether this cluster enables disk auto-scaling. The maximum memory allowed for the selected cluster tier and the oplog size can limit storage auto-scaling.",
													},
												},
												CustomType: DiskGbType{
													ObjectType: types.ObjectType{
														AttrTypes: DiskGbValue{}.AttributeTypes(ctx),
													},
												},
												Optional:            true,
												Computed:            true,
												Description:         "Setting that enables disk auto-scaling.",
												MarkdownDescription: "Setting that enables disk auto-scaling.",
											},
										},
										CustomType: AutoScalingType{
											ObjectType: types.ObjectType{
												AttrTypes: AutoScalingValue{}.AttributeTypes(ctx),
											},
										},
										Optional:            true,
										Computed:            true,
										Description:         "Options that determine how this cluster handles resource scaling.",
										MarkdownDescription: "Options that determine how this cluster handles resource scaling.",
									},
									"backing_provider_name": schema.StringAttribute{
										Optional:            true,
										Computed:            true,
										Description:         "Cloud service provider on which MongoDB Cloud provisioned the multi-tenant cluster. The resource returns this parameter when **providerName** is `TENANT` and **electableSpecs.instanceSize** is `M0`, `M2` or `M5`.",
										MarkdownDescription: "Cloud service provider on which MongoDB Cloud provisioned the multi-tenant cluster. The resource returns this parameter when **providerName** is `TENANT` and **electableSpecs.instanceSize** is `M0`, `M2` or `M5`.",
									},
									"electable_specs": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"disk_iops": schema.Int64Attribute{
												Optional:            true,
												Computed:            true,
												Description:         "Target throughput desired for storage attached to your Azure-provisioned cluster. Change this parameter if you:\n\n- set `\"replicationSpecs[n].regionConfigs[m].providerName\" : \"Azure\"`.\n- set `\"replicationSpecs[n].regionConfigs[m].electableSpecs.instanceSize\" : \"M40\"` or greater not including `Mxx_NVME` tiers.\n\nThe maximum input/output operations per second (IOPS) depend on the selected **.instanceSize** and **.diskSizeGB**.\nThis parameter defaults to the cluster tier's standard IOPS value.\nChanging this value impacts cluster cost.",
												MarkdownDescription: "Target throughput desired for storage attached to your Azure-provisioned cluster. Change this parameter if you:\n\n- set `\"replicationSpecs[n].regionConfigs[m].providerName\" : \"Azure\"`.\n- set `\"replicationSpecs[n].regionConfigs[m].electableSpecs.instanceSize\" : \"M40\"` or greater not including `Mxx_NVME` tiers.\n\nThe maximum input/output operations per second (IOPS) depend on the selected **.instanceSize** and **.diskSizeGB**.\nThis parameter defaults to the cluster tier's standard IOPS value.\nChanging this value impacts cluster cost.",
											},
											"disk_size_gb": schema.Float64Attribute{
												Optional:            true,
												Computed:            true,
												Description:         "Storage capacity of instance data volumes expressed in gigabytes. Increase this number to add capacity.\n\n This value must be equal for all shards and node types.\n\n This value is not configurable on M0/M2/M5 clusters.\n\n MongoDB Cloud requires this parameter if you set **replicationSpecs**.\n\n If you specify a disk size below the minimum (10 GB), this parameter defaults to the minimum disk size value. \n\n Storage charge calculations depend on whether you choose the default value or a custom value.\n\n The maximum value for disk storage cannot exceed 50 times the maximum RAM for the selected cluster. If you require more storage space, consider upgrading your cluster to a higher tier.",
												MarkdownDescription: "Storage capacity of instance data volumes expressed in gigabytes. Increase this number to add capacity.\n\n This value must be equal for all shards and node types.\n\n This value is not configurable on M0/M2/M5 clusters.\n\n MongoDB Cloud requires this parameter if you set **replicationSpecs**.\n\n If you specify a disk size below the minimum (10 GB), this parameter defaults to the minimum disk size value. \n\n Storage charge calculations depend on whether you choose the default value or a custom value.\n\n The maximum value for disk storage cannot exceed 50 times the maximum RAM for the selected cluster. If you require more storage space, consider upgrading your cluster to a higher tier.",
											},
											"ebs_volume_type": schema.StringAttribute{
												Optional:            true,
												Computed:            true,
												Description:         "Type of storage you want to attach to your AWS-provisioned cluster.\n\n- `STANDARD` volume types can't exceed the default input/output operations per second (IOPS) rate for the selected volume size. \n\n- `PROVISIONED` volume types must fall within the allowable IOPS range for the selected volume size. You must set this value to (`PROVISIONED`) for NVMe clusters.",
												MarkdownDescription: "Type of storage you want to attach to your AWS-provisioned cluster.\n\n- `STANDARD` volume types can't exceed the default input/output operations per second (IOPS) rate for the selected volume size. \n\n- `PROVISIONED` volume types must fall within the allowable IOPS range for the selected volume size. You must set this value to (`PROVISIONED`) for NVMe clusters.",
												Default:             stringdefault.StaticString("STANDARD"),
											},
											"instance_size": schema.StringAttribute{
												Optional:            true,
												Computed:            true,
												Description:         "Hardware specification for the instances in this M0/M2/M5 tier cluster.",
												MarkdownDescription: "Hardware specification for the instances in this M0/M2/M5 tier cluster.",
											},
											"node_count": schema.Int64Attribute{
												Optional:            true,
												Computed:            true,
												Description:         "Number of nodes of the given type for MongoDB Cloud to deploy to the region.",
												MarkdownDescription: "Number of nodes of the given type for MongoDB Cloud to deploy to the region.",
											},
										},
										CustomType: ElectableSpecsType{
											ObjectType: types.ObjectType{
												AttrTypes: ElectableSpecsValue{}.AttributeTypes(ctx),
											},
										},
										Optional:            true,
										Computed:            true,
										Description:         "Hardware specifications for all electable nodes deployed in the region. Electable nodes can become the primary and can enable local reads. If you don't specify this option, MongoDB Cloud deploys no electable nodes to the region.",
										MarkdownDescription: "Hardware specifications for all electable nodes deployed in the region. Electable nodes can become the primary and can enable local reads. If you don't specify this option, MongoDB Cloud deploys no electable nodes to the region.",
									},
									"priority": schema.Int64Attribute{
										Optional:            true,
										Computed:            true,
										Description:         "Precedence is given to this region when a primary election occurs. If your **regionConfigs** has only **readOnlySpecs**, **analyticsSpecs**, or both, set this value to `0`. If you have multiple **regionConfigs** objects (your cluster is multi-region or multi-cloud), they must have priorities in descending order. The highest priority is `7`.\n\n**Example:** If you have three regions, their priorities would be `7`, `6`, and `5` respectively. If you added two more regions for supporting electable nodes, the priorities of those regions would be `4` and `3` respectively.",
										MarkdownDescription: "Precedence is given to this region when a primary election occurs. If your **regionConfigs** has only **readOnlySpecs**, **analyticsSpecs**, or both, set this value to `0`. If you have multiple **regionConfigs** objects (your cluster is multi-region or multi-cloud), they must have priorities in descending order. The highest priority is `7`.\n\n**Example:** If you have three regions, their priorities would be `7`, `6`, and `5` respectively. If you added two more regions for supporting electable nodes, the priorities of those regions would be `4` and `3` respectively.",
										Validators: []validator.Int64{
											int64validator.Between(0, 7),
										},
									},
									"provider_name": schema.StringAttribute{
										Optional:            true,
										Computed:            true,
										Description:         "Cloud service provider on which MongoDB Cloud provisions the hosts. Set dedicated clusters to `AWS`, `GCP`, `AZURE` or `TENANT`.",
										MarkdownDescription: "Cloud service provider on which MongoDB Cloud provisions the hosts. Set dedicated clusters to `AWS`, `GCP`, `AZURE` or `TENANT`.",
									},
									"read_only_specs": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"disk_iops": schema.Int64Attribute{
												Optional:            true,
												Computed:            true,
												Description:         "Target throughput desired for storage attached to your Azure-provisioned cluster. Change this parameter if you:\n\n- set `\"replicationSpecs[n].regionConfigs[m].providerName\" : \"Azure\"`.\n- set `\"replicationSpecs[n].regionConfigs[m].electableSpecs.instanceSize\" : \"M40\"` or greater not including `Mxx_NVME` tiers.\n\nThe maximum input/output operations per second (IOPS) depend on the selected **.instanceSize** and **.diskSizeGB**.\nThis parameter defaults to the cluster tier's standard IOPS value.\nChanging this value impacts cluster cost.",
												MarkdownDescription: "Target throughput desired for storage attached to your Azure-provisioned cluster. Change this parameter if you:\n\n- set `\"replicationSpecs[n].regionConfigs[m].providerName\" : \"Azure\"`.\n- set `\"replicationSpecs[n].regionConfigs[m].electableSpecs.instanceSize\" : \"M40\"` or greater not including `Mxx_NVME` tiers.\n\nThe maximum input/output operations per second (IOPS) depend on the selected **.instanceSize** and **.diskSizeGB**.\nThis parameter defaults to the cluster tier's standard IOPS value.\nChanging this value impacts cluster cost.",
											},
											"disk_size_gb": schema.Float64Attribute{
												Optional:            true,
												Computed:            true,
												Description:         "Storage capacity of instance data volumes expressed in gigabytes. Increase this number to add capacity.\n\n This value must be equal for all shards and node types.\n\n This value is not configurable on M0/M2/M5 clusters.\n\n MongoDB Cloud requires this parameter if you set **replicationSpecs**.\n\n If you specify a disk size below the minimum (10 GB), this parameter defaults to the minimum disk size value. \n\n Storage charge calculations depend on whether you choose the default value or a custom value.\n\n The maximum value for disk storage cannot exceed 50 times the maximum RAM for the selected cluster. If you require more storage space, consider upgrading your cluster to a higher tier.",
												MarkdownDescription: "Storage capacity of instance data volumes expressed in gigabytes. Increase this number to add capacity.\n\n This value must be equal for all shards and node types.\n\n This value is not configurable on M0/M2/M5 clusters.\n\n MongoDB Cloud requires this parameter if you set **replicationSpecs**.\n\n If you specify a disk size below the minimum (10 GB), this parameter defaults to the minimum disk size value. \n\n Storage charge calculations depend on whether you choose the default value or a custom value.\n\n The maximum value for disk storage cannot exceed 50 times the maximum RAM for the selected cluster. If you require more storage space, consider upgrading your cluster to a higher tier.",
											},
											"ebs_volume_type": schema.StringAttribute{
												Optional:            true,
												Computed:            true,
												Description:         "Type of storage you want to attach to your AWS-provisioned cluster.\n\n- `STANDARD` volume types can't exceed the default input/output operations per second (IOPS) rate for the selected volume size. \n\n- `PROVISIONED` volume types must fall within the allowable IOPS range for the selected volume size. You must set this value to (`PROVISIONED`) for NVMe clusters.",
												MarkdownDescription: "Type of storage you want to attach to your AWS-provisioned cluster.\n\n- `STANDARD` volume types can't exceed the default input/output operations per second (IOPS) rate for the selected volume size. \n\n- `PROVISIONED` volume types must fall within the allowable IOPS range for the selected volume size. You must set this value to (`PROVISIONED`) for NVMe clusters.",
												Default:             stringdefault.StaticString("STANDARD"),
											},
											"instance_size": schema.StringAttribute{
												Optional:            true,
												Computed:            true,
												Description:         "Hardware specification for the instance sizes in this region in this shard. Each instance size has a default storage and memory capacity. Electable nodes and read-only nodes (known as \"base nodes\") within a single shard must use the same instance size. Analytics nodes can scale independently from base nodes within a shard. Both base nodes and analytics nodes can scale independently from their equivalents in other shards.",
												MarkdownDescription: "Hardware specification for the instance sizes in this region in this shard. Each instance size has a default storage and memory capacity. Electable nodes and read-only nodes (known as \"base nodes\") within a single shard must use the same instance size. Analytics nodes can scale independently from base nodes within a shard. Both base nodes and analytics nodes can scale independently from their equivalents in other shards.",
											},
											"node_count": schema.Int64Attribute{
												Optional:            true,
												Computed:            true,
												Description:         "Number of nodes of the given type for MongoDB Cloud to deploy to the region.",
												MarkdownDescription: "Number of nodes of the given type for MongoDB Cloud to deploy to the region.",
											},
										},
										CustomType: ReadOnlySpecsType{
											ObjectType: types.ObjectType{
												AttrTypes: ReadOnlySpecsValue{}.AttributeTypes(ctx),
											},
										},
										Optional:            true,
										Computed:            true,
										Description:         "Hardware specifications for read-only nodes in the region. Read-only nodes can never become the primary member, but can enable local reads. If you don't specify this parameter, no read-only nodes are deployed to the region.",
										MarkdownDescription: "Hardware specifications for read-only nodes in the region. Read-only nodes can never become the primary member, but can enable local reads. If you don't specify this parameter, no read-only nodes are deployed to the region.",
									},
									"region_name": schema.StringAttribute{
										Optional:            true,
										Computed:            true,
										Description:         "Physical location of your MongoDB cluster nodes. The region you choose can affect network latency for clients accessing your databases. The region name is only returned in the response for single-region clusters. When MongoDB Cloud deploys a dedicated cluster, it checks if a VPC or VPC connection exists for that provider and region. If not, MongoDB Cloud creates them as part of the deployment. It assigns the VPC a Classless Inter-Domain Routing (CIDR) block. To limit a new VPC peering connection to one Classless Inter-Domain Routing (CIDR) block and region, create the connection first. Deploy the cluster after the connection starts. GCP Clusters and Multi-region clusters require one VPC peering connection for each region. MongoDB nodes can use only the peering connection that resides in the same region as the nodes to communicate with the peered VPC.",
										MarkdownDescription: "Physical location of your MongoDB cluster nodes. The region you choose can affect network latency for clients accessing your databases. The region name is only returned in the response for single-region clusters. When MongoDB Cloud deploys a dedicated cluster, it checks if a VPC or VPC connection exists for that provider and region. If not, MongoDB Cloud creates them as part of the deployment. It assigns the VPC a Classless Inter-Domain Routing (CIDR) block. To limit a new VPC peering connection to one Classless Inter-Domain Routing (CIDR) block and region, create the connection first. Deploy the cluster after the connection starts. GCP Clusters and Multi-region clusters require one VPC peering connection for each region. MongoDB nodes can use only the peering connection that resides in the same region as the nodes to communicate with the peered VPC.",
									},
								},
								CustomType: RegionConfigsType{
									ObjectType: types.ObjectType{
										AttrTypes: RegionConfigsValue{}.AttributeTypes(ctx),
									},
								},
							},
							Optional:            true,
							Computed:            true,
							Description:         "Hardware specifications for nodes set for a given region. Each **regionConfigs** object describes the region's priority in elections and the number and type of MongoDB nodes that MongoDB Cloud deploys to the region. Each **regionConfigs** object must have either an **analyticsSpecs** object, **electableSpecs** object, or **readOnlySpecs** object. Tenant clusters only require **electableSpecs. Dedicated** clusters can specify any of these specifications, but must have at least one **electableSpecs** object within a **replicationSpec**.\n\n**Example:**\n\nIf you set `\"replicationSpecs[n].regionConfigs[m].analyticsSpecs.instanceSize\" : \"M30\"`, set `\"replicationSpecs[n].regionConfigs[m].electableSpecs.instanceSize\" : `\"M30\"` if you have electable nodes and `\"replicationSpecs[n].regionConfigs[m].readOnlySpecs.instanceSize\" : `\"M30\"` if you have read-only nodes.",
							MarkdownDescription: "Hardware specifications for nodes set for a given region. Each **regionConfigs** object describes the region's priority in elections and the number and type of MongoDB nodes that MongoDB Cloud deploys to the region. Each **regionConfigs** object must have either an **analyticsSpecs** object, **electableSpecs** object, or **readOnlySpecs** object. Tenant clusters only require **electableSpecs. Dedicated** clusters can specify any of these specifications, but must have at least one **electableSpecs** object within a **replicationSpec**.\n\n**Example:**\n\nIf you set `\"replicationSpecs[n].regionConfigs[m].analyticsSpecs.instanceSize\" : \"M30\"`, set `\"replicationSpecs[n].regionConfigs[m].electableSpecs.instanceSize\" : `\"M30\"` if you have electable nodes and `\"replicationSpecs[n].regionConfigs[m].readOnlySpecs.instanceSize\" : `\"M30\"` if you have read-only nodes.",
						},
						"zone_id": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "Unique 24-hexadecimal digit string that identifies the zone in a Global Cluster. This value can be used to configure Global Cluster backup policies.",
							MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the zone in a Global Cluster. This value can be used to configure Global Cluster backup policies.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(24, 24),
								stringvalidator.RegexMatches(regexp.MustCompile("^([a-f0-9]{24})$"), ""),
							},
						},
						"zone_name": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "Human-readable label that describes the zone this shard belongs to in a Global Cluster. Provide this value only if \"clusterType\" : \"GEOSHARDED\" but not \"selfManagedSharding\" : true.",
							MarkdownDescription: "Human-readable label that describes the zone this shard belongs to in a Global Cluster. Provide this value only if \"clusterType\" : \"GEOSHARDED\" but not \"selfManagedSharding\" : true.",
						},
					},
					CustomType: ReplicationSpecsType{
						ObjectType: types.ObjectType{
							AttrTypes: ReplicationSpecsValue{}.AttributeTypes(ctx),
						},
					},
				},
				Optional:            true,
				Computed:            true,
				Description:         "List of settings that configure your cluster regions. This array has one object per shard representing node configurations in each shard. For replica sets there is only one object representing node configurations.",
				MarkdownDescription: "List of settings that configure your cluster regions. This array has one object per shard representing node configurations in each shard. For replica sets there is only one object representing node configurations.",
			},
			"root_cert_type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Root Certificate Authority that MongoDB Cloud cluster uses. MongoDB Cloud supports Internet Security Research Group.",
				MarkdownDescription: "Root Certificate Authority that MongoDB Cloud cluster uses. MongoDB Cloud supports Internet Security Research Group.",
				Default:             stringdefault.StaticString("ISRGROOTX1"),
			},
			"state_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Human-readable label that indicates the current operating condition of this cluster.",
				MarkdownDescription: "Human-readable label that indicates the current operating condition of this cluster.",
			},
			"tags": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Required:            true,
							Description:         "Constant that defines the set of the tag. For example, `environment` in the `environment : production` tag.",
							MarkdownDescription: "Constant that defines the set of the tag. For example, `environment` in the `environment : production` tag.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 255),
							},
						},
						"value": schema.StringAttribute{
							Required:            true,
							Description:         "Variable that belongs to the set of the tag. For example, `production` in the `environment : production` tag.",
							MarkdownDescription: "Variable that belongs to the set of the tag. For example, `production` in the `environment : production` tag.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 255),
							},
						},
					},
					CustomType: TagsType{
						ObjectType: types.ObjectType{
							AttrTypes: TagsValue{}.AttributeTypes(ctx),
						},
					},
				},
				Optional:            true,
				Computed:            true,
				Description:         "List that contains key-value pairs between 1 to 255 characters in length for tagging and categorizing the cluster.",
				MarkdownDescription: "List that contains key-value pairs between 1 to 255 characters in length for tagging and categorizing the cluster.",
			},
			"termination_protection_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Flag that indicates whether termination protection is enabled on the cluster. If set to `true`, MongoDB Cloud won't delete the cluster. If set to `false`, MongoDB Cloud will delete the cluster.",
				MarkdownDescription: "Flag that indicates whether termination protection is enabled on the cluster. If set to `true`, MongoDB Cloud won't delete the cluster. If set to `false`, MongoDB Cloud will delete the cluster.",
				Default:             booldefault.StaticBool(false),
			},
			"version_release_system": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Method by which the cluster maintains the MongoDB versions. If value is `CONTINUOUS`, you must not specify **mongoDBMajorVersion**.",
				MarkdownDescription: "Method by which the cluster maintains the MongoDB versions. If value is `CONTINUOUS`, you must not specify **mongoDBMajorVersion**.",
				Default:             stringdefault.StaticString("LTS"),
			},
		},
	}
}

type TFModel struct {
	Labels                                    types.List                      `tfsdk:"labels"`
	Tags                                      types.List                      `tfsdk:"tags"`
	ReplicationSpecs                          types.List                      `tfsdk:"replication_specs"`
	Links                                     types.List                      `tfsdk:"links"`
	CreateDate                                types.String                    `tfsdk:"create_date"`
	ClusterName                               types.String                    `tfsdk:"cluster_name"`
	ConfigServerType                          types.String                    `tfsdk:"config_server_type"`
	VersionReleaseSystem                      types.String                    `tfsdk:"version_release_system"`
	AcceptDataRisksAndForceReplicaSetReconfig types.String                    `tfsdk:"accept_data_risks_and_force_replica_set_reconfig"`
	DiskWarmingMode                           types.String                    `tfsdk:"disk_warming_mode"`
	EncryptionAtRestProvider                  types.String                    `tfsdk:"encryption_at_rest_provider"`
	FeatureCompatibilityVersion               types.String                    `tfsdk:"feature_compatibility_version"`
	FeatureCompatibilityVersionExpirationDate types.String                    `tfsdk:"feature_compatibility_version_expiration_date"`
	StateName                                 types.String                    `tfsdk:"state_name"`
	GroupId                                   types.String                    `tfsdk:"group_id"`
	Id                                        types.String                    `tfsdk:"id"`
	ClusterType                               types.String                    `tfsdk:"cluster_type"`
	ConfigServerManagementMode                types.String                    `tfsdk:"config_server_management_mode"`
	RootCertType                              types.String                    `tfsdk:"root_cert_type"`
	MongoDbmajorVersion                       types.String                    `tfsdk:"mongo_dbmajor_version"`
	MongoDbversion                            types.String                    `tfsdk:"mongo_dbversion"`
	Name                                      types.String                    `tfsdk:"name"`
	ReplicaSetScalingStrategy                 types.String                    `tfsdk:"replica_set_scaling_strategy"`
	ConnectionStrings                         ConnectionStringsValue          `tfsdk:"connection_strings"`
	MongoDbemployeeAccessGrant                MongoDbemployeeAccessGrantValue `tfsdk:"mongo_dbemployee_access_grant"`
	BiConnector                               BiConnectorValue                `tfsdk:"bi_connector"`
	PitEnabled                                types.Bool                      `tfsdk:"pit_enabled"`
	RedactClientLogData                       types.Bool                      `tfsdk:"redact_client_log_data"`
	Paused                                    types.Bool                      `tfsdk:"paused"`
	GlobalClusterSelfManagedSharding          types.Bool                      `tfsdk:"global_cluster_self_managed_sharding"`
	BackupEnabled                             types.Bool                      `tfsdk:"backup_enabled"`
	TerminationProtectionEnabled              types.Bool                      `tfsdk:"termination_protection_enabled"`
}

var _ basetypes.ObjectTypable = BiConnectorType{}

type BiConnectorType struct {
	basetypes.ObjectType
}

func (t BiConnectorType) Equal(o attr.Type) bool {
	other, ok := o.(BiConnectorType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t BiConnectorType) String() string {
	return "BiConnectorType"
}

func (t BiConnectorType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	enabledAttribute, ok := attributes["enabled"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`enabled is missing from object`)

		return nil, diags
	}

	enabledVal, ok := enabledAttribute.(basetypes.BoolValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`enabled expected to be basetypes.BoolValue, was: %T`, enabledAttribute))
	}

	readPreferenceAttribute, ok := attributes["read_preference"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`read_preference is missing from object`)

		return nil, diags
	}

	readPreferenceVal, ok := readPreferenceAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`read_preference expected to be basetypes.StringValue, was: %T`, readPreferenceAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return BiConnectorValue{
		Enabled:        enabledVal,
		ReadPreference: readPreferenceVal,
		state:          attr.ValueStateKnown,
	}, diags
}

func NewBiConnectorValueNull() BiConnectorValue {
	return BiConnectorValue{
		state: attr.ValueStateNull,
	}
}

func NewBiConnectorValueUnknown() BiConnectorValue {
	return BiConnectorValue{
		state: attr.ValueStateUnknown,
	}
}

func NewBiConnectorValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (BiConnectorValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing BiConnectorValue Attribute Value",
				"While creating a BiConnectorValue value, a missing attribute value was detected. "+
					"A BiConnectorValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("BiConnectorValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid BiConnectorValue Attribute Type",
				"While creating a BiConnectorValue value, an invalid attribute value was detected. "+
					"A BiConnectorValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("BiConnectorValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("BiConnectorValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra BiConnectorValue Attribute Value",
				"While creating a BiConnectorValue value, an extra attribute value was detected. "+
					"A BiConnectorValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra BiConnectorValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewBiConnectorValueUnknown(), diags
	}

	enabledAttribute, ok := attributes["enabled"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`enabled is missing from object`)

		return NewBiConnectorValueUnknown(), diags
	}

	enabledVal, ok := enabledAttribute.(basetypes.BoolValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`enabled expected to be basetypes.BoolValue, was: %T`, enabledAttribute))
	}

	readPreferenceAttribute, ok := attributes["read_preference"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`read_preference is missing from object`)

		return NewBiConnectorValueUnknown(), diags
	}

	readPreferenceVal, ok := readPreferenceAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`read_preference expected to be basetypes.StringValue, was: %T`, readPreferenceAttribute))
	}

	if diags.HasError() {
		return NewBiConnectorValueUnknown(), diags
	}

	return BiConnectorValue{
		Enabled:        enabledVal,
		ReadPreference: readPreferenceVal,
		state:          attr.ValueStateKnown,
	}, diags
}

func NewBiConnectorValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) BiConnectorValue {
	object, diags := NewBiConnectorValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewBiConnectorValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t BiConnectorType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewBiConnectorValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewBiConnectorValueUnknown(), nil
	}

	if in.IsNull() {
		return NewBiConnectorValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewBiConnectorValueMust(BiConnectorValue{}.AttributeTypes(ctx), attributes), nil
}

func (t BiConnectorType) ValueType(ctx context.Context) attr.Value {
	return BiConnectorValue{}
}

var _ basetypes.ObjectValuable = BiConnectorValue{}

type BiConnectorValue struct {
	ReadPreference basetypes.StringValue `tfsdk:"read_preference"`
	Enabled        basetypes.BoolValue   `tfsdk:"enabled"`
	state          attr.ValueState
}

func (v BiConnectorValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 2)

	var val tftypes.Value
	var err error

	attrTypes["enabled"] = basetypes.BoolType{}.TerraformType(ctx)
	attrTypes["read_preference"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 2)

		val, err = v.Enabled.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["enabled"] = val

		val, err = v.ReadPreference.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["read_preference"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v BiConnectorValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v BiConnectorValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v BiConnectorValue) String() string {
	return "BiConnectorValue"
}

func (v BiConnectorValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := map[string]attr.Type{
		"enabled":         basetypes.BoolType{},
		"read_preference": basetypes.StringType{},
	}

	if v.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	}

	if v.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	}

	objVal, diags := types.ObjectValue(
		attributeTypes,
		map[string]attr.Value{
			"enabled":         v.Enabled,
			"read_preference": v.ReadPreference,
		})

	return objVal, diags
}

func (v BiConnectorValue) Equal(o attr.Value) bool {
	other, ok := o.(BiConnectorValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Enabled.Equal(other.Enabled) {
		return false
	}

	if !v.ReadPreference.Equal(other.ReadPreference) {
		return false
	}

	return true
}

func (v BiConnectorValue) Type(ctx context.Context) attr.Type {
	return BiConnectorType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v BiConnectorValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"enabled":         basetypes.BoolType{},
		"read_preference": basetypes.StringType{},
	}
}

var _ basetypes.ObjectTypable = ConnectionStringsType{}

type ConnectionStringsType struct {
	basetypes.ObjectType
}

func (t ConnectionStringsType) Equal(o attr.Type) bool {
	other, ok := o.(ConnectionStringsType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t ConnectionStringsType) String() string {
	return "ConnectionStringsType"
}

func (t ConnectionStringsType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	awsPrivateLinkAttribute, ok := attributes["aws_private_link"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`aws_private_link is missing from object`)

		return nil, diags
	}

	awsPrivateLinkVal, ok := awsPrivateLinkAttribute.(basetypes.MapValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`aws_private_link expected to be basetypes.MapValue, was: %T`, awsPrivateLinkAttribute))
	}

	awsPrivateLinkSrvAttribute, ok := attributes["aws_private_link_srv"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`aws_private_link_srv is missing from object`)

		return nil, diags
	}

	awsPrivateLinkSrvVal, ok := awsPrivateLinkSrvAttribute.(basetypes.MapValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`aws_private_link_srv expected to be basetypes.MapValue, was: %T`, awsPrivateLinkSrvAttribute))
	}

	privateAttribute, ok := attributes["private"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`private is missing from object`)

		return nil, diags
	}

	privateVal, ok := privateAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`private expected to be basetypes.StringValue, was: %T`, privateAttribute))
	}

	privateEndpointAttribute, ok := attributes["private_endpoint"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`private_endpoint is missing from object`)

		return nil, diags
	}

	privateEndpointVal, ok := privateEndpointAttribute.(basetypes.ListValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`private_endpoint expected to be basetypes.ListValue, was: %T`, privateEndpointAttribute))
	}

	privateSrvAttribute, ok := attributes["private_srv"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`private_srv is missing from object`)

		return nil, diags
	}

	privateSrvVal, ok := privateSrvAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`private_srv expected to be basetypes.StringValue, was: %T`, privateSrvAttribute))
	}

	standardAttribute, ok := attributes["standard"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`standard is missing from object`)

		return nil, diags
	}

	standardVal, ok := standardAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`standard expected to be basetypes.StringValue, was: %T`, standardAttribute))
	}

	standardSrvAttribute, ok := attributes["standard_srv"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`standard_srv is missing from object`)

		return nil, diags
	}

	standardSrvVal, ok := standardSrvAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`standard_srv expected to be basetypes.StringValue, was: %T`, standardSrvAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return ConnectionStringsValue{
		AwsPrivateLink:    awsPrivateLinkVal,
		AwsPrivateLinkSrv: awsPrivateLinkSrvVal,
		Private:           privateVal,
		PrivateEndpoint:   privateEndpointVal,
		PrivateSrv:        privateSrvVal,
		Standard:          standardVal,
		StandardSrv:       standardSrvVal,
		state:             attr.ValueStateKnown,
	}, diags
}

func NewConnectionStringsValueNull() ConnectionStringsValue {
	return ConnectionStringsValue{
		state: attr.ValueStateNull,
	}
}

func NewConnectionStringsValueUnknown() ConnectionStringsValue {
	return ConnectionStringsValue{
		state: attr.ValueStateUnknown,
	}
}

func NewConnectionStringsValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (ConnectionStringsValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing ConnectionStringsValue Attribute Value",
				"While creating a ConnectionStringsValue value, a missing attribute value was detected. "+
					"A ConnectionStringsValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("ConnectionStringsValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid ConnectionStringsValue Attribute Type",
				"While creating a ConnectionStringsValue value, an invalid attribute value was detected. "+
					"A ConnectionStringsValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("ConnectionStringsValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("ConnectionStringsValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra ConnectionStringsValue Attribute Value",
				"While creating a ConnectionStringsValue value, an extra attribute value was detected. "+
					"A ConnectionStringsValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra ConnectionStringsValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewConnectionStringsValueUnknown(), diags
	}

	awsPrivateLinkAttribute, ok := attributes["aws_private_link"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`aws_private_link is missing from object`)

		return NewConnectionStringsValueUnknown(), diags
	}

	awsPrivateLinkVal, ok := awsPrivateLinkAttribute.(basetypes.MapValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`aws_private_link expected to be basetypes.MapValue, was: %T`, awsPrivateLinkAttribute))
	}

	awsPrivateLinkSrvAttribute, ok := attributes["aws_private_link_srv"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`aws_private_link_srv is missing from object`)

		return NewConnectionStringsValueUnknown(), diags
	}

	awsPrivateLinkSrvVal, ok := awsPrivateLinkSrvAttribute.(basetypes.MapValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`aws_private_link_srv expected to be basetypes.MapValue, was: %T`, awsPrivateLinkSrvAttribute))
	}

	privateAttribute, ok := attributes["private"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`private is missing from object`)

		return NewConnectionStringsValueUnknown(), diags
	}

	privateVal, ok := privateAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`private expected to be basetypes.StringValue, was: %T`, privateAttribute))
	}

	privateEndpointAttribute, ok := attributes["private_endpoint"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`private_endpoint is missing from object`)

		return NewConnectionStringsValueUnknown(), diags
	}

	privateEndpointVal, ok := privateEndpointAttribute.(basetypes.ListValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`private_endpoint expected to be basetypes.ListValue, was: %T`, privateEndpointAttribute))
	}

	privateSrvAttribute, ok := attributes["private_srv"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`private_srv is missing from object`)

		return NewConnectionStringsValueUnknown(), diags
	}

	privateSrvVal, ok := privateSrvAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`private_srv expected to be basetypes.StringValue, was: %T`, privateSrvAttribute))
	}

	standardAttribute, ok := attributes["standard"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`standard is missing from object`)

		return NewConnectionStringsValueUnknown(), diags
	}

	standardVal, ok := standardAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`standard expected to be basetypes.StringValue, was: %T`, standardAttribute))
	}

	standardSrvAttribute, ok := attributes["standard_srv"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`standard_srv is missing from object`)

		return NewConnectionStringsValueUnknown(), diags
	}

	standardSrvVal, ok := standardSrvAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`standard_srv expected to be basetypes.StringValue, was: %T`, standardSrvAttribute))
	}

	if diags.HasError() {
		return NewConnectionStringsValueUnknown(), diags
	}

	return ConnectionStringsValue{
		AwsPrivateLink:    awsPrivateLinkVal,
		AwsPrivateLinkSrv: awsPrivateLinkSrvVal,
		Private:           privateVal,
		PrivateEndpoint:   privateEndpointVal,
		PrivateSrv:        privateSrvVal,
		Standard:          standardVal,
		StandardSrv:       standardSrvVal,
		state:             attr.ValueStateKnown,
	}, diags
}

func NewConnectionStringsValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) ConnectionStringsValue {
	object, diags := NewConnectionStringsValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewConnectionStringsValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t ConnectionStringsType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewConnectionStringsValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewConnectionStringsValueUnknown(), nil
	}

	if in.IsNull() {
		return NewConnectionStringsValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewConnectionStringsValueMust(ConnectionStringsValue{}.AttributeTypes(ctx), attributes), nil
}

func (t ConnectionStringsType) ValueType(ctx context.Context) attr.Value {
	return ConnectionStringsValue{}
}

var _ basetypes.ObjectValuable = ConnectionStringsValue{}

type ConnectionStringsValue struct {
	AwsPrivateLink    basetypes.MapValue    `tfsdk:"aws_private_link"`
	AwsPrivateLinkSrv basetypes.MapValue    `tfsdk:"aws_private_link_srv"`
	Private           basetypes.StringValue `tfsdk:"private"`
	PrivateEndpoint   basetypes.ListValue   `tfsdk:"private_endpoint"`
	PrivateSrv        basetypes.StringValue `tfsdk:"private_srv"`
	Standard          basetypes.StringValue `tfsdk:"standard"`
	StandardSrv       basetypes.StringValue `tfsdk:"standard_srv"`
	state             attr.ValueState
}

func (v ConnectionStringsValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 7)

	var val tftypes.Value
	var err error

	attrTypes["aws_private_link"] = basetypes.MapType{
		ElemType: types.StringType,
	}.TerraformType(ctx)
	attrTypes["aws_private_link_srv"] = basetypes.MapType{
		ElemType: types.StringType,
	}.TerraformType(ctx)
	attrTypes["private"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["private_endpoint"] = basetypes.ListType{
		ElemType: PrivateEndpointValue{}.Type(ctx),
	}.TerraformType(ctx)
	attrTypes["private_srv"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["standard"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["standard_srv"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 7)

		val, err = v.AwsPrivateLink.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["aws_private_link"] = val

		val, err = v.AwsPrivateLinkSrv.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["aws_private_link_srv"] = val

		val, err = v.Private.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["private"] = val

		val, err = v.PrivateEndpoint.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["private_endpoint"] = val

		val, err = v.PrivateSrv.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["private_srv"] = val

		val, err = v.Standard.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["standard"] = val

		val, err = v.StandardSrv.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["standard_srv"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v ConnectionStringsValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v ConnectionStringsValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v ConnectionStringsValue) String() string {
	return "ConnectionStringsValue"
}

func (v ConnectionStringsValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	privateEndpoint := types.ListValueMust(
		PrivateEndpointType{
			basetypes.ObjectType{
				AttrTypes: PrivateEndpointValue{}.AttributeTypes(ctx),
			},
		},
		v.PrivateEndpoint.Elements(),
	)

	if v.PrivateEndpoint.IsNull() {
		privateEndpoint = types.ListNull(
			PrivateEndpointType{
				basetypes.ObjectType{
					AttrTypes: PrivateEndpointValue{}.AttributeTypes(ctx),
				},
			},
		)
	}

	if v.PrivateEndpoint.IsUnknown() {
		privateEndpoint = types.ListUnknown(
			PrivateEndpointType{
				basetypes.ObjectType{
					AttrTypes: PrivateEndpointValue{}.AttributeTypes(ctx),
				},
			},
		)
	}

	var awsPrivateLinkVal basetypes.MapValue
	switch {
	case v.AwsPrivateLink.IsUnknown():
		awsPrivateLinkVal = types.MapUnknown(types.StringType)
	case v.AwsPrivateLink.IsNull():
		awsPrivateLinkVal = types.MapNull(types.StringType)
	default:
		var d diag.Diagnostics
		awsPrivateLinkVal, d = types.MapValue(types.StringType, v.AwsPrivateLink.Elements())
		diags.Append(d...)
	}

	if diags.HasError() {
		return types.ObjectUnknown(map[string]attr.Type{
			"aws_private_link": basetypes.MapType{
				ElemType: types.StringType,
			},
			"aws_private_link_srv": basetypes.MapType{
				ElemType: types.StringType,
			},
			"private": basetypes.StringType{},
			"private_endpoint": basetypes.ListType{
				ElemType: PrivateEndpointValue{}.Type(ctx),
			},
			"private_srv":  basetypes.StringType{},
			"standard":     basetypes.StringType{},
			"standard_srv": basetypes.StringType{},
		}), diags
	}

	var awsPrivateLinkSrvVal basetypes.MapValue
	switch {
	case v.AwsPrivateLinkSrv.IsUnknown():
		awsPrivateLinkSrvVal = types.MapUnknown(types.StringType)
	case v.AwsPrivateLinkSrv.IsNull():
		awsPrivateLinkSrvVal = types.MapNull(types.StringType)
	default:
		var d diag.Diagnostics
		awsPrivateLinkSrvVal, d = types.MapValue(types.StringType, v.AwsPrivateLinkSrv.Elements())
		diags.Append(d...)
	}

	if diags.HasError() {
		return types.ObjectUnknown(map[string]attr.Type{
			"aws_private_link": basetypes.MapType{
				ElemType: types.StringType,
			},
			"aws_private_link_srv": basetypes.MapType{
				ElemType: types.StringType,
			},
			"private": basetypes.StringType{},
			"private_endpoint": basetypes.ListType{
				ElemType: PrivateEndpointValue{}.Type(ctx),
			},
			"private_srv":  basetypes.StringType{},
			"standard":     basetypes.StringType{},
			"standard_srv": basetypes.StringType{},
		}), diags
	}

	attributeTypes := map[string]attr.Type{
		"aws_private_link": basetypes.MapType{
			ElemType: types.StringType,
		},
		"aws_private_link_srv": basetypes.MapType{
			ElemType: types.StringType,
		},
		"private": basetypes.StringType{},
		"private_endpoint": basetypes.ListType{
			ElemType: PrivateEndpointValue{}.Type(ctx),
		},
		"private_srv":  basetypes.StringType{},
		"standard":     basetypes.StringType{},
		"standard_srv": basetypes.StringType{},
	}

	if v.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	}

	if v.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	}

	objVal, diags := types.ObjectValue(
		attributeTypes,
		map[string]attr.Value{
			"aws_private_link":     awsPrivateLinkVal,
			"aws_private_link_srv": awsPrivateLinkSrvVal,
			"private":              v.Private,
			"private_endpoint":     privateEndpoint,
			"private_srv":          v.PrivateSrv,
			"standard":             v.Standard,
			"standard_srv":         v.StandardSrv,
		})

	return objVal, diags
}

func (v ConnectionStringsValue) Equal(o attr.Value) bool {
	other, ok := o.(ConnectionStringsValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.AwsPrivateLink.Equal(other.AwsPrivateLink) {
		return false
	}

	if !v.AwsPrivateLinkSrv.Equal(other.AwsPrivateLinkSrv) {
		return false
	}

	if !v.Private.Equal(other.Private) {
		return false
	}

	if !v.PrivateEndpoint.Equal(other.PrivateEndpoint) {
		return false
	}

	if !v.PrivateSrv.Equal(other.PrivateSrv) {
		return false
	}

	if !v.Standard.Equal(other.Standard) {
		return false
	}

	if !v.StandardSrv.Equal(other.StandardSrv) {
		return false
	}

	return true
}

func (v ConnectionStringsValue) Type(ctx context.Context) attr.Type {
	return ConnectionStringsType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v ConnectionStringsValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"aws_private_link": basetypes.MapType{
			ElemType: types.StringType,
		},
		"aws_private_link_srv": basetypes.MapType{
			ElemType: types.StringType,
		},
		"private": basetypes.StringType{},
		"private_endpoint": basetypes.ListType{
			ElemType: PrivateEndpointValue{}.Type(ctx),
		},
		"private_srv":  basetypes.StringType{},
		"standard":     basetypes.StringType{},
		"standard_srv": basetypes.StringType{},
	}
}

var _ basetypes.ObjectTypable = PrivateEndpointType{}

type PrivateEndpointType struct {
	basetypes.ObjectType
}

func (t PrivateEndpointType) Equal(o attr.Type) bool {
	other, ok := o.(PrivateEndpointType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t PrivateEndpointType) String() string {
	return "PrivateEndpointType"
}

func (t PrivateEndpointType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	connectionStringAttribute, ok := attributes["connection_string"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`connection_string is missing from object`)

		return nil, diags
	}

	connectionStringVal, ok := connectionStringAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`connection_string expected to be basetypes.StringValue, was: %T`, connectionStringAttribute))
	}

	endpointsAttribute, ok := attributes["endpoints"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`endpoints is missing from object`)

		return nil, diags
	}

	endpointsVal, ok := endpointsAttribute.(basetypes.ListValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`endpoints expected to be basetypes.ListValue, was: %T`, endpointsAttribute))
	}

	srvConnectionStringAttribute, ok := attributes["srv_connection_string"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`srv_connection_string is missing from object`)

		return nil, diags
	}

	srvConnectionStringVal, ok := srvConnectionStringAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`srv_connection_string expected to be basetypes.StringValue, was: %T`, srvConnectionStringAttribute))
	}

	srvShardOptimizedConnectionStringAttribute, ok := attributes["srv_shard_optimized_connection_string"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`srv_shard_optimized_connection_string is missing from object`)

		return nil, diags
	}

	srvShardOptimizedConnectionStringVal, ok := srvShardOptimizedConnectionStringAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`srv_shard_optimized_connection_string expected to be basetypes.StringValue, was: %T`, srvShardOptimizedConnectionStringAttribute))
	}

	typeAttribute, ok := attributes["type"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`type is missing from object`)

		return nil, diags
	}

	typeVal, ok := typeAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`type expected to be basetypes.StringValue, was: %T`, typeAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return PrivateEndpointValue{
		ConnectionString:                  connectionStringVal,
		Endpoints:                         endpointsVal,
		SrvConnectionString:               srvConnectionStringVal,
		SrvShardOptimizedConnectionString: srvShardOptimizedConnectionStringVal,
		PrivateEndpointType:               typeVal,
		state:                             attr.ValueStateKnown,
	}, diags
}

func NewPrivateEndpointValueNull() PrivateEndpointValue {
	return PrivateEndpointValue{
		state: attr.ValueStateNull,
	}
}

func NewPrivateEndpointValueUnknown() PrivateEndpointValue {
	return PrivateEndpointValue{
		state: attr.ValueStateUnknown,
	}
}

func NewPrivateEndpointValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (PrivateEndpointValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing PrivateEndpointValue Attribute Value",
				"While creating a PrivateEndpointValue value, a missing attribute value was detected. "+
					"A PrivateEndpointValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("PrivateEndpointValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid PrivateEndpointValue Attribute Type",
				"While creating a PrivateEndpointValue value, an invalid attribute value was detected. "+
					"A PrivateEndpointValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("PrivateEndpointValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("PrivateEndpointValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra PrivateEndpointValue Attribute Value",
				"While creating a PrivateEndpointValue value, an extra attribute value was detected. "+
					"A PrivateEndpointValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra PrivateEndpointValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewPrivateEndpointValueUnknown(), diags
	}

	connectionStringAttribute, ok := attributes["connection_string"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`connection_string is missing from object`)

		return NewPrivateEndpointValueUnknown(), diags
	}

	connectionStringVal, ok := connectionStringAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`connection_string expected to be basetypes.StringValue, was: %T`, connectionStringAttribute))
	}

	endpointsAttribute, ok := attributes["endpoints"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`endpoints is missing from object`)

		return NewPrivateEndpointValueUnknown(), diags
	}

	endpointsVal, ok := endpointsAttribute.(basetypes.ListValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`endpoints expected to be basetypes.ListValue, was: %T`, endpointsAttribute))
	}

	srvConnectionStringAttribute, ok := attributes["srv_connection_string"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`srv_connection_string is missing from object`)

		return NewPrivateEndpointValueUnknown(), diags
	}

	srvConnectionStringVal, ok := srvConnectionStringAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`srv_connection_string expected to be basetypes.StringValue, was: %T`, srvConnectionStringAttribute))
	}

	srvShardOptimizedConnectionStringAttribute, ok := attributes["srv_shard_optimized_connection_string"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`srv_shard_optimized_connection_string is missing from object`)

		return NewPrivateEndpointValueUnknown(), diags
	}

	srvShardOptimizedConnectionStringVal, ok := srvShardOptimizedConnectionStringAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`srv_shard_optimized_connection_string expected to be basetypes.StringValue, was: %T`, srvShardOptimizedConnectionStringAttribute))
	}

	typeAttribute, ok := attributes["type"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`type is missing from object`)

		return NewPrivateEndpointValueUnknown(), diags
	}

	typeVal, ok := typeAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`type expected to be basetypes.StringValue, was: %T`, typeAttribute))
	}

	if diags.HasError() {
		return NewPrivateEndpointValueUnknown(), diags
	}

	return PrivateEndpointValue{
		ConnectionString:                  connectionStringVal,
		Endpoints:                         endpointsVal,
		SrvConnectionString:               srvConnectionStringVal,
		SrvShardOptimizedConnectionString: srvShardOptimizedConnectionStringVal,
		PrivateEndpointType:               typeVal,
		state:                             attr.ValueStateKnown,
	}, diags
}

func NewPrivateEndpointValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) PrivateEndpointValue {
	object, diags := NewPrivateEndpointValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewPrivateEndpointValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t PrivateEndpointType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewPrivateEndpointValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewPrivateEndpointValueUnknown(), nil
	}

	if in.IsNull() {
		return NewPrivateEndpointValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewPrivateEndpointValueMust(PrivateEndpointValue{}.AttributeTypes(ctx), attributes), nil
}

func (t PrivateEndpointType) ValueType(ctx context.Context) attr.Value {
	return PrivateEndpointValue{}
}

var _ basetypes.ObjectValuable = PrivateEndpointValue{}

type PrivateEndpointValue struct {
	ConnectionString                  basetypes.StringValue `tfsdk:"connection_string"`
	Endpoints                         basetypes.ListValue   `tfsdk:"endpoints"`
	SrvConnectionString               basetypes.StringValue `tfsdk:"srv_connection_string"`
	SrvShardOptimizedConnectionString basetypes.StringValue `tfsdk:"srv_shard_optimized_connection_string"`
	PrivateEndpointType               basetypes.StringValue `tfsdk:"type"`
	state                             attr.ValueState
}

func (v PrivateEndpointValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 5)

	var val tftypes.Value
	var err error

	attrTypes["connection_string"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["endpoints"] = basetypes.ListType{
		ElemType: EndpointsValue{}.Type(ctx),
	}.TerraformType(ctx)
	attrTypes["srv_connection_string"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["srv_shard_optimized_connection_string"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["type"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 5)

		val, err = v.ConnectionString.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["connection_string"] = val

		val, err = v.Endpoints.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["endpoints"] = val

		val, err = v.SrvConnectionString.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["srv_connection_string"] = val

		val, err = v.SrvShardOptimizedConnectionString.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["srv_shard_optimized_connection_string"] = val

		val, err = v.PrivateEndpointType.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["type"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v PrivateEndpointValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v PrivateEndpointValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v PrivateEndpointValue) String() string {
	return "PrivateEndpointValue"
}

func (v PrivateEndpointValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	endpoints := types.ListValueMust(
		EndpointsType{
			basetypes.ObjectType{
				AttrTypes: EndpointsValue{}.AttributeTypes(ctx),
			},
		},
		v.Endpoints.Elements(),
	)

	if v.Endpoints.IsNull() {
		endpoints = types.ListNull(
			EndpointsType{
				basetypes.ObjectType{
					AttrTypes: EndpointsValue{}.AttributeTypes(ctx),
				},
			},
		)
	}

	if v.Endpoints.IsUnknown() {
		endpoints = types.ListUnknown(
			EndpointsType{
				basetypes.ObjectType{
					AttrTypes: EndpointsValue{}.AttributeTypes(ctx),
				},
			},
		)
	}

	attributeTypes := map[string]attr.Type{
		"connection_string": basetypes.StringType{},
		"endpoints": basetypes.ListType{
			ElemType: EndpointsValue{}.Type(ctx),
		},
		"srv_connection_string":                 basetypes.StringType{},
		"srv_shard_optimized_connection_string": basetypes.StringType{},
		"type":                                  basetypes.StringType{},
	}

	if v.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	}

	if v.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	}

	objVal, diags := types.ObjectValue(
		attributeTypes,
		map[string]attr.Value{
			"connection_string":                     v.ConnectionString,
			"endpoints":                             endpoints,
			"srv_connection_string":                 v.SrvConnectionString,
			"srv_shard_optimized_connection_string": v.SrvShardOptimizedConnectionString,
			"type":                                  v.PrivateEndpointType,
		})

	return objVal, diags
}

func (v PrivateEndpointValue) Equal(o attr.Value) bool {
	other, ok := o.(PrivateEndpointValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.ConnectionString.Equal(other.ConnectionString) {
		return false
	}

	if !v.Endpoints.Equal(other.Endpoints) {
		return false
	}

	if !v.SrvConnectionString.Equal(other.SrvConnectionString) {
		return false
	}

	if !v.SrvShardOptimizedConnectionString.Equal(other.SrvShardOptimizedConnectionString) {
		return false
	}

	if !v.PrivateEndpointType.Equal(other.PrivateEndpointType) {
		return false
	}

	return true
}

func (v PrivateEndpointValue) Type(ctx context.Context) attr.Type {
	return PrivateEndpointType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v PrivateEndpointValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"connection_string": basetypes.StringType{},
		"endpoints": basetypes.ListType{
			ElemType: EndpointsValue{}.Type(ctx),
		},
		"srv_connection_string":                 basetypes.StringType{},
		"srv_shard_optimized_connection_string": basetypes.StringType{},
		"type":                                  basetypes.StringType{},
	}
}

var _ basetypes.ObjectTypable = EndpointsType{}

type EndpointsType struct {
	basetypes.ObjectType
}

func (t EndpointsType) Equal(o attr.Type) bool {
	other, ok := o.(EndpointsType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t EndpointsType) String() string {
	return "EndpointsType"
}

func (t EndpointsType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	endpointIdAttribute, ok := attributes["endpoint_id"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`endpoint_id is missing from object`)

		return nil, diags
	}

	endpointIdVal, ok := endpointIdAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`endpoint_id expected to be basetypes.StringValue, was: %T`, endpointIdAttribute))
	}

	providerNameAttribute, ok := attributes["provider_name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`provider_name is missing from object`)

		return nil, diags
	}

	providerNameVal, ok := providerNameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`provider_name expected to be basetypes.StringValue, was: %T`, providerNameAttribute))
	}

	regionAttribute, ok := attributes["region"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`region is missing from object`)

		return nil, diags
	}

	regionVal, ok := regionAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`region expected to be basetypes.StringValue, was: %T`, regionAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return EndpointsValue{
		EndpointId:   endpointIdVal,
		ProviderName: providerNameVal,
		Region:       regionVal,
		state:        attr.ValueStateKnown,
	}, diags
}

func NewEndpointsValueNull() EndpointsValue {
	return EndpointsValue{
		state: attr.ValueStateNull,
	}
}

func NewEndpointsValueUnknown() EndpointsValue {
	return EndpointsValue{
		state: attr.ValueStateUnknown,
	}
}

func NewEndpointsValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (EndpointsValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing EndpointsValue Attribute Value",
				"While creating a EndpointsValue value, a missing attribute value was detected. "+
					"A EndpointsValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("EndpointsValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid EndpointsValue Attribute Type",
				"While creating a EndpointsValue value, an invalid attribute value was detected. "+
					"A EndpointsValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("EndpointsValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("EndpointsValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra EndpointsValue Attribute Value",
				"While creating a EndpointsValue value, an extra attribute value was detected. "+
					"A EndpointsValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra EndpointsValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewEndpointsValueUnknown(), diags
	}

	endpointIdAttribute, ok := attributes["endpoint_id"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`endpoint_id is missing from object`)

		return NewEndpointsValueUnknown(), diags
	}

	endpointIdVal, ok := endpointIdAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`endpoint_id expected to be basetypes.StringValue, was: %T`, endpointIdAttribute))
	}

	providerNameAttribute, ok := attributes["provider_name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`provider_name is missing from object`)

		return NewEndpointsValueUnknown(), diags
	}

	providerNameVal, ok := providerNameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`provider_name expected to be basetypes.StringValue, was: %T`, providerNameAttribute))
	}

	regionAttribute, ok := attributes["region"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`region is missing from object`)

		return NewEndpointsValueUnknown(), diags
	}

	regionVal, ok := regionAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`region expected to be basetypes.StringValue, was: %T`, regionAttribute))
	}

	if diags.HasError() {
		return NewEndpointsValueUnknown(), diags
	}

	return EndpointsValue{
		EndpointId:   endpointIdVal,
		ProviderName: providerNameVal,
		Region:       regionVal,
		state:        attr.ValueStateKnown,
	}, diags
}

func NewEndpointsValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) EndpointsValue {
	object, diags := NewEndpointsValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewEndpointsValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t EndpointsType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewEndpointsValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewEndpointsValueUnknown(), nil
	}

	if in.IsNull() {
		return NewEndpointsValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewEndpointsValueMust(EndpointsValue{}.AttributeTypes(ctx), attributes), nil
}

func (t EndpointsType) ValueType(ctx context.Context) attr.Value {
	return EndpointsValue{}
}

var _ basetypes.ObjectValuable = EndpointsValue{}

type EndpointsValue struct {
	EndpointId   basetypes.StringValue `tfsdk:"endpoint_id"`
	ProviderName basetypes.StringValue `tfsdk:"provider_name"`
	Region       basetypes.StringValue `tfsdk:"region"`
	state        attr.ValueState
}

func (v EndpointsValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 3)

	var val tftypes.Value
	var err error

	attrTypes["endpoint_id"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["provider_name"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["region"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 3)

		val, err = v.EndpointId.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["endpoint_id"] = val

		val, err = v.ProviderName.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["provider_name"] = val

		val, err = v.Region.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["region"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v EndpointsValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v EndpointsValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v EndpointsValue) String() string {
	return "EndpointsValue"
}

func (v EndpointsValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := map[string]attr.Type{
		"endpoint_id":   basetypes.StringType{},
		"provider_name": basetypes.StringType{},
		"region":        basetypes.StringType{},
	}

	if v.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	}

	if v.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	}

	objVal, diags := types.ObjectValue(
		attributeTypes,
		map[string]attr.Value{
			"endpoint_id":   v.EndpointId,
			"provider_name": v.ProviderName,
			"region":        v.Region,
		})

	return objVal, diags
}

func (v EndpointsValue) Equal(o attr.Value) bool {
	other, ok := o.(EndpointsValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.EndpointId.Equal(other.EndpointId) {
		return false
	}

	if !v.ProviderName.Equal(other.ProviderName) {
		return false
	}

	if !v.Region.Equal(other.Region) {
		return false
	}

	return true
}

func (v EndpointsValue) Type(ctx context.Context) attr.Type {
	return EndpointsType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v EndpointsValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"endpoint_id":   basetypes.StringType{},
		"provider_name": basetypes.StringType{},
		"region":        basetypes.StringType{},
	}
}

var _ basetypes.ObjectTypable = LabelsType{}

type LabelsType struct {
	basetypes.ObjectType
}

func (t LabelsType) Equal(o attr.Type) bool {
	other, ok := o.(LabelsType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t LabelsType) String() string {
	return "LabelsType"
}

func (t LabelsType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	keyAttribute, ok := attributes["key"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`key is missing from object`)

		return nil, diags
	}

	keyVal, ok := keyAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`key expected to be basetypes.StringValue, was: %T`, keyAttribute))
	}

	valueAttribute, ok := attributes["value"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`value is missing from object`)

		return nil, diags
	}

	valueVal, ok := valueAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`value expected to be basetypes.StringValue, was: %T`, valueAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return LabelsValue{
		Key:   keyVal,
		Value: valueVal,
		state: attr.ValueStateKnown,
	}, diags
}

func NewLabelsValueNull() LabelsValue {
	return LabelsValue{
		state: attr.ValueStateNull,
	}
}

func NewLabelsValueUnknown() LabelsValue {
	return LabelsValue{
		state: attr.ValueStateUnknown,
	}
}

func NewLabelsValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (LabelsValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing LabelsValue Attribute Value",
				"While creating a LabelsValue value, a missing attribute value was detected. "+
					"A LabelsValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("LabelsValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid LabelsValue Attribute Type",
				"While creating a LabelsValue value, an invalid attribute value was detected. "+
					"A LabelsValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("LabelsValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("LabelsValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra LabelsValue Attribute Value",
				"While creating a LabelsValue value, an extra attribute value was detected. "+
					"A LabelsValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra LabelsValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewLabelsValueUnknown(), diags
	}

	keyAttribute, ok := attributes["key"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`key is missing from object`)

		return NewLabelsValueUnknown(), diags
	}

	keyVal, ok := keyAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`key expected to be basetypes.StringValue, was: %T`, keyAttribute))
	}

	valueAttribute, ok := attributes["value"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`value is missing from object`)

		return NewLabelsValueUnknown(), diags
	}

	valueVal, ok := valueAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`value expected to be basetypes.StringValue, was: %T`, valueAttribute))
	}

	if diags.HasError() {
		return NewLabelsValueUnknown(), diags
	}

	return LabelsValue{
		Key:   keyVal,
		Value: valueVal,
		state: attr.ValueStateKnown,
	}, diags
}

func NewLabelsValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) LabelsValue {
	object, diags := NewLabelsValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewLabelsValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t LabelsType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewLabelsValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewLabelsValueUnknown(), nil
	}

	if in.IsNull() {
		return NewLabelsValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewLabelsValueMust(LabelsValue{}.AttributeTypes(ctx), attributes), nil
}

func (t LabelsType) ValueType(ctx context.Context) attr.Value {
	return LabelsValue{}
}

var _ basetypes.ObjectValuable = LabelsValue{}

type LabelsValue struct {
	Key   basetypes.StringValue `tfsdk:"key"`
	Value basetypes.StringValue `tfsdk:"value"`
	state attr.ValueState
}

func (v LabelsValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 2)

	var val tftypes.Value
	var err error

	attrTypes["key"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["value"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 2)

		val, err = v.Key.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["key"] = val

		val, err = v.Value.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["value"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v LabelsValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v LabelsValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v LabelsValue) String() string {
	return "LabelsValue"
}

func (v LabelsValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := map[string]attr.Type{
		"key":   basetypes.StringType{},
		"value": basetypes.StringType{},
	}

	if v.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	}

	if v.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	}

	objVal, diags := types.ObjectValue(
		attributeTypes,
		map[string]attr.Value{
			"key":   v.Key,
			"value": v.Value,
		})

	return objVal, diags
}

func (v LabelsValue) Equal(o attr.Value) bool {
	other, ok := o.(LabelsValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Key.Equal(other.Key) {
		return false
	}

	if !v.Value.Equal(other.Value) {
		return false
	}

	return true
}

func (v LabelsValue) Type(ctx context.Context) attr.Type {
	return LabelsType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v LabelsValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"key":   basetypes.StringType{},
		"value": basetypes.StringType{},
	}
}

var _ basetypes.ObjectTypable = LinksType{}

type LinksType struct {
	basetypes.ObjectType
}

func (t LinksType) Equal(o attr.Type) bool {
	other, ok := o.(LinksType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t LinksType) String() string {
	return "LinksType"
}

func (t LinksType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	hrefAttribute, ok := attributes["href"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`href is missing from object`)

		return nil, diags
	}

	hrefVal, ok := hrefAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`href expected to be basetypes.StringValue, was: %T`, hrefAttribute))
	}

	relAttribute, ok := attributes["rel"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`rel is missing from object`)

		return nil, diags
	}

	relVal, ok := relAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`rel expected to be basetypes.StringValue, was: %T`, relAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return LinksValue{
		Href:  hrefVal,
		Rel:   relVal,
		state: attr.ValueStateKnown,
	}, diags
}

func NewLinksValueNull() LinksValue {
	return LinksValue{
		state: attr.ValueStateNull,
	}
}

func NewLinksValueUnknown() LinksValue {
	return LinksValue{
		state: attr.ValueStateUnknown,
	}
}

func NewLinksValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (LinksValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing LinksValue Attribute Value",
				"While creating a LinksValue value, a missing attribute value was detected. "+
					"A LinksValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("LinksValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid LinksValue Attribute Type",
				"While creating a LinksValue value, an invalid attribute value was detected. "+
					"A LinksValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("LinksValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("LinksValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra LinksValue Attribute Value",
				"While creating a LinksValue value, an extra attribute value was detected. "+
					"A LinksValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra LinksValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewLinksValueUnknown(), diags
	}

	hrefAttribute, ok := attributes["href"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`href is missing from object`)

		return NewLinksValueUnknown(), diags
	}

	hrefVal, ok := hrefAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`href expected to be basetypes.StringValue, was: %T`, hrefAttribute))
	}

	relAttribute, ok := attributes["rel"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`rel is missing from object`)

		return NewLinksValueUnknown(), diags
	}

	relVal, ok := relAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`rel expected to be basetypes.StringValue, was: %T`, relAttribute))
	}

	if diags.HasError() {
		return NewLinksValueUnknown(), diags
	}

	return LinksValue{
		Href:  hrefVal,
		Rel:   relVal,
		state: attr.ValueStateKnown,
	}, diags
}

func NewLinksValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) LinksValue {
	object, diags := NewLinksValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewLinksValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t LinksType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewLinksValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewLinksValueUnknown(), nil
	}

	if in.IsNull() {
		return NewLinksValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewLinksValueMust(LinksValue{}.AttributeTypes(ctx), attributes), nil
}

func (t LinksType) ValueType(ctx context.Context) attr.Value {
	return LinksValue{}
}

var _ basetypes.ObjectValuable = LinksValue{}

type LinksValue struct {
	Href  basetypes.StringValue `tfsdk:"href"`
	Rel   basetypes.StringValue `tfsdk:"rel"`
	state attr.ValueState
}

func (v LinksValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 2)

	var val tftypes.Value
	var err error

	attrTypes["href"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["rel"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 2)

		val, err = v.Href.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["href"] = val

		val, err = v.Rel.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["rel"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v LinksValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v LinksValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v LinksValue) String() string {
	return "LinksValue"
}

func (v LinksValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := map[string]attr.Type{
		"href": basetypes.StringType{},
		"rel":  basetypes.StringType{},
	}

	if v.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	}

	if v.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	}

	objVal, diags := types.ObjectValue(
		attributeTypes,
		map[string]attr.Value{
			"href": v.Href,
			"rel":  v.Rel,
		})

	return objVal, diags
}

func (v LinksValue) Equal(o attr.Value) bool {
	other, ok := o.(LinksValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Href.Equal(other.Href) {
		return false
	}

	if !v.Rel.Equal(other.Rel) {
		return false
	}

	return true
}

func (v LinksValue) Type(ctx context.Context) attr.Type {
	return LinksType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v LinksValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"href": basetypes.StringType{},
		"rel":  basetypes.StringType{},
	}
}

var _ basetypes.ObjectTypable = MongoDbemployeeAccessGrantType{}

type MongoDbemployeeAccessGrantType struct {
	basetypes.ObjectType
}

func (t MongoDbemployeeAccessGrantType) Equal(o attr.Type) bool {
	other, ok := o.(MongoDbemployeeAccessGrantType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t MongoDbemployeeAccessGrantType) String() string {
	return "MongoDbemployeeAccessGrantType"
}

func (t MongoDbemployeeAccessGrantType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	expirationTimeAttribute, ok := attributes["expiration_time"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`expiration_time is missing from object`)

		return nil, diags
	}

	expirationTimeVal, ok := expirationTimeAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`expiration_time expected to be basetypes.StringValue, was: %T`, expirationTimeAttribute))
	}

	grantTypeAttribute, ok := attributes["grant_type"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`grant_type is missing from object`)

		return nil, diags
	}

	grantTypeVal, ok := grantTypeAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`grant_type expected to be basetypes.StringValue, was: %T`, grantTypeAttribute))
	}

	linksAttribute, ok := attributes["links"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`links is missing from object`)

		return nil, diags
	}

	linksVal, ok := linksAttribute.(basetypes.ListValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`links expected to be basetypes.ListValue, was: %T`, linksAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return MongoDbemployeeAccessGrantValue{
		ExpirationTime: expirationTimeVal,
		GrantType:      grantTypeVal,
		Links:          linksVal,
		state:          attr.ValueStateKnown,
	}, diags
}

func NewMongoDbemployeeAccessGrantValueNull() MongoDbemployeeAccessGrantValue {
	return MongoDbemployeeAccessGrantValue{
		state: attr.ValueStateNull,
	}
}

func NewMongoDbemployeeAccessGrantValueUnknown() MongoDbemployeeAccessGrantValue {
	return MongoDbemployeeAccessGrantValue{
		state: attr.ValueStateUnknown,
	}
}

func NewMongoDbemployeeAccessGrantValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (MongoDbemployeeAccessGrantValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing MongoDbemployeeAccessGrantValue Attribute Value",
				"While creating a MongoDbemployeeAccessGrantValue value, a missing attribute value was detected. "+
					"A MongoDbemployeeAccessGrantValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("MongoDbemployeeAccessGrantValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid MongoDbemployeeAccessGrantValue Attribute Type",
				"While creating a MongoDbemployeeAccessGrantValue value, an invalid attribute value was detected. "+
					"A MongoDbemployeeAccessGrantValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("MongoDbemployeeAccessGrantValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("MongoDbemployeeAccessGrantValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra MongoDbemployeeAccessGrantValue Attribute Value",
				"While creating a MongoDbemployeeAccessGrantValue value, an extra attribute value was detected. "+
					"A MongoDbemployeeAccessGrantValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra MongoDbemployeeAccessGrantValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewMongoDbemployeeAccessGrantValueUnknown(), diags
	}

	expirationTimeAttribute, ok := attributes["expiration_time"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`expiration_time is missing from object`)

		return NewMongoDbemployeeAccessGrantValueUnknown(), diags
	}

	expirationTimeVal, ok := expirationTimeAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`expiration_time expected to be basetypes.StringValue, was: %T`, expirationTimeAttribute))
	}

	grantTypeAttribute, ok := attributes["grant_type"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`grant_type is missing from object`)

		return NewMongoDbemployeeAccessGrantValueUnknown(), diags
	}

	grantTypeVal, ok := grantTypeAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`grant_type expected to be basetypes.StringValue, was: %T`, grantTypeAttribute))
	}

	linksAttribute, ok := attributes["links"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`links is missing from object`)

		return NewMongoDbemployeeAccessGrantValueUnknown(), diags
	}

	linksVal, ok := linksAttribute.(basetypes.ListValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`links expected to be basetypes.ListValue, was: %T`, linksAttribute))
	}

	if diags.HasError() {
		return NewMongoDbemployeeAccessGrantValueUnknown(), diags
	}

	return MongoDbemployeeAccessGrantValue{
		ExpirationTime: expirationTimeVal,
		GrantType:      grantTypeVal,
		Links:          linksVal,
		state:          attr.ValueStateKnown,
	}, diags
}

func NewMongoDbemployeeAccessGrantValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) MongoDbemployeeAccessGrantValue {
	object, diags := NewMongoDbemployeeAccessGrantValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewMongoDbemployeeAccessGrantValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t MongoDbemployeeAccessGrantType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewMongoDbemployeeAccessGrantValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewMongoDbemployeeAccessGrantValueUnknown(), nil
	}

	if in.IsNull() {
		return NewMongoDbemployeeAccessGrantValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewMongoDbemployeeAccessGrantValueMust(MongoDbemployeeAccessGrantValue{}.AttributeTypes(ctx), attributes), nil
}

func (t MongoDbemployeeAccessGrantType) ValueType(ctx context.Context) attr.Value {
	return MongoDbemployeeAccessGrantValue{}
}

var _ basetypes.ObjectValuable = MongoDbemployeeAccessGrantValue{}

type MongoDbemployeeAccessGrantValue struct {
	ExpirationTime basetypes.StringValue `tfsdk:"expiration_time"`
	GrantType      basetypes.StringValue `tfsdk:"grant_type"`
	Links          basetypes.ListValue   `tfsdk:"links"`
	state          attr.ValueState
}

func (v MongoDbemployeeAccessGrantValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 3)

	var val tftypes.Value
	var err error

	attrTypes["expiration_time"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["grant_type"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["links"] = basetypes.ListType{
		ElemType: LinksValue{}.Type(ctx),
	}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 3)

		val, err = v.ExpirationTime.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["expiration_time"] = val

		val, err = v.GrantType.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["grant_type"] = val

		val, err = v.Links.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["links"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v MongoDbemployeeAccessGrantValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v MongoDbemployeeAccessGrantValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v MongoDbemployeeAccessGrantValue) String() string {
	return "MongoDbemployeeAccessGrantValue"
}

func (v MongoDbemployeeAccessGrantValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	links := types.ListValueMust(
		LinksType{
			basetypes.ObjectType{
				AttrTypes: LinksValue{}.AttributeTypes(ctx),
			},
		},
		v.Links.Elements(),
	)

	if v.Links.IsNull() {
		links = types.ListNull(
			LinksType{
				basetypes.ObjectType{
					AttrTypes: LinksValue{}.AttributeTypes(ctx),
				},
			},
		)
	}

	if v.Links.IsUnknown() {
		links = types.ListUnknown(
			LinksType{
				basetypes.ObjectType{
					AttrTypes: LinksValue{}.AttributeTypes(ctx),
				},
			},
		)
	}

	attributeTypes := map[string]attr.Type{
		"expiration_time": basetypes.StringType{},
		"grant_type":      basetypes.StringType{},
		"links": basetypes.ListType{
			ElemType: LinksValue{}.Type(ctx),
		},
	}

	if v.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	}

	if v.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	}

	objVal, diags := types.ObjectValue(
		attributeTypes,
		map[string]attr.Value{
			"expiration_time": v.ExpirationTime,
			"grant_type":      v.GrantType,
			"links":           links,
		})

	return objVal, diags
}

func (v MongoDbemployeeAccessGrantValue) Equal(o attr.Value) bool {
	other, ok := o.(MongoDbemployeeAccessGrantValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.ExpirationTime.Equal(other.ExpirationTime) {
		return false
	}

	if !v.GrantType.Equal(other.GrantType) {
		return false
	}

	if !v.Links.Equal(other.Links) {
		return false
	}

	return true
}

func (v MongoDbemployeeAccessGrantValue) Type(ctx context.Context) attr.Type {
	return MongoDbemployeeAccessGrantType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v MongoDbemployeeAccessGrantValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"expiration_time": basetypes.StringType{},
		"grant_type":      basetypes.StringType{},
		"links": basetypes.ListType{
			ElemType: LinksValue{}.Type(ctx),
		},
	}
}

var _ basetypes.ObjectTypable = ReplicationSpecsType{}

type ReplicationSpecsType struct {
	basetypes.ObjectType
}

func (t ReplicationSpecsType) Equal(o attr.Type) bool {
	other, ok := o.(ReplicationSpecsType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t ReplicationSpecsType) String() string {
	return "ReplicationSpecsType"
}

func (t ReplicationSpecsType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	idAttribute, ok := attributes["id"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`id is missing from object`)

		return nil, diags
	}

	idVal, ok := idAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`id expected to be basetypes.StringValue, was: %T`, idAttribute))
	}

	regionConfigsAttribute, ok := attributes["region_configs"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`region_configs is missing from object`)

		return nil, diags
	}

	regionConfigsVal, ok := regionConfigsAttribute.(basetypes.ListValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`region_configs expected to be basetypes.ListValue, was: %T`, regionConfigsAttribute))
	}

	zoneIdAttribute, ok := attributes["zone_id"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`zone_id is missing from object`)

		return nil, diags
	}

	zoneIdVal, ok := zoneIdAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`zone_id expected to be basetypes.StringValue, was: %T`, zoneIdAttribute))
	}

	zoneNameAttribute, ok := attributes["zone_name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`zone_name is missing from object`)

		return nil, diags
	}

	zoneNameVal, ok := zoneNameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`zone_name expected to be basetypes.StringValue, was: %T`, zoneNameAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return ReplicationSpecsValue{
		Id:            idVal,
		RegionConfigs: regionConfigsVal,
		ZoneId:        zoneIdVal,
		ZoneName:      zoneNameVal,
		state:         attr.ValueStateKnown,
	}, diags
}

func NewReplicationSpecsValueNull() ReplicationSpecsValue {
	return ReplicationSpecsValue{
		state: attr.ValueStateNull,
	}
}

func NewReplicationSpecsValueUnknown() ReplicationSpecsValue {
	return ReplicationSpecsValue{
		state: attr.ValueStateUnknown,
	}
}

func NewReplicationSpecsValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (ReplicationSpecsValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing ReplicationSpecsValue Attribute Value",
				"While creating a ReplicationSpecsValue value, a missing attribute value was detected. "+
					"A ReplicationSpecsValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("ReplicationSpecsValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid ReplicationSpecsValue Attribute Type",
				"While creating a ReplicationSpecsValue value, an invalid attribute value was detected. "+
					"A ReplicationSpecsValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("ReplicationSpecsValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("ReplicationSpecsValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra ReplicationSpecsValue Attribute Value",
				"While creating a ReplicationSpecsValue value, an extra attribute value was detected. "+
					"A ReplicationSpecsValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra ReplicationSpecsValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewReplicationSpecsValueUnknown(), diags
	}

	idAttribute, ok := attributes["id"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`id is missing from object`)

		return NewReplicationSpecsValueUnknown(), diags
	}

	idVal, ok := idAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`id expected to be basetypes.StringValue, was: %T`, idAttribute))
	}

	regionConfigsAttribute, ok := attributes["region_configs"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`region_configs is missing from object`)

		return NewReplicationSpecsValueUnknown(), diags
	}

	regionConfigsVal, ok := regionConfigsAttribute.(basetypes.ListValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`region_configs expected to be basetypes.ListValue, was: %T`, regionConfigsAttribute))
	}

	zoneIdAttribute, ok := attributes["zone_id"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`zone_id is missing from object`)

		return NewReplicationSpecsValueUnknown(), diags
	}

	zoneIdVal, ok := zoneIdAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`zone_id expected to be basetypes.StringValue, was: %T`, zoneIdAttribute))
	}

	zoneNameAttribute, ok := attributes["zone_name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`zone_name is missing from object`)

		return NewReplicationSpecsValueUnknown(), diags
	}

	zoneNameVal, ok := zoneNameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`zone_name expected to be basetypes.StringValue, was: %T`, zoneNameAttribute))
	}

	if diags.HasError() {
		return NewReplicationSpecsValueUnknown(), diags
	}

	return ReplicationSpecsValue{
		Id:            idVal,
		RegionConfigs: regionConfigsVal,
		ZoneId:        zoneIdVal,
		ZoneName:      zoneNameVal,
		state:         attr.ValueStateKnown,
	}, diags
}

func NewReplicationSpecsValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) ReplicationSpecsValue {
	object, diags := NewReplicationSpecsValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewReplicationSpecsValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t ReplicationSpecsType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewReplicationSpecsValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewReplicationSpecsValueUnknown(), nil
	}

	if in.IsNull() {
		return NewReplicationSpecsValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewReplicationSpecsValueMust(ReplicationSpecsValue{}.AttributeTypes(ctx), attributes), nil
}

func (t ReplicationSpecsType) ValueType(ctx context.Context) attr.Value {
	return ReplicationSpecsValue{}
}

var _ basetypes.ObjectValuable = ReplicationSpecsValue{}

type ReplicationSpecsValue struct {
	Id            basetypes.StringValue `tfsdk:"id"`
	RegionConfigs basetypes.ListValue   `tfsdk:"region_configs"`
	ZoneId        basetypes.StringValue `tfsdk:"zone_id"`
	ZoneName      basetypes.StringValue `tfsdk:"zone_name"`
	state         attr.ValueState
}

func (v ReplicationSpecsValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 4)

	var val tftypes.Value
	var err error

	attrTypes["id"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["region_configs"] = basetypes.ListType{
		ElemType: RegionConfigsValue{}.Type(ctx),
	}.TerraformType(ctx)
	attrTypes["zone_id"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["zone_name"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 4)

		val, err = v.Id.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["id"] = val

		val, err = v.RegionConfigs.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["region_configs"] = val

		val, err = v.ZoneId.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["zone_id"] = val

		val, err = v.ZoneName.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["zone_name"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v ReplicationSpecsValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v ReplicationSpecsValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v ReplicationSpecsValue) String() string {
	return "ReplicationSpecsValue"
}

func (v ReplicationSpecsValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	regionConfigs := types.ListValueMust(
		RegionConfigsType{
			basetypes.ObjectType{
				AttrTypes: RegionConfigsValue{}.AttributeTypes(ctx),
			},
		},
		v.RegionConfigs.Elements(),
	)

	if v.RegionConfigs.IsNull() {
		regionConfigs = types.ListNull(
			RegionConfigsType{
				basetypes.ObjectType{
					AttrTypes: RegionConfigsValue{}.AttributeTypes(ctx),
				},
			},
		)
	}

	if v.RegionConfigs.IsUnknown() {
		regionConfigs = types.ListUnknown(
			RegionConfigsType{
				basetypes.ObjectType{
					AttrTypes: RegionConfigsValue{}.AttributeTypes(ctx),
				},
			},
		)
	}

	attributeTypes := map[string]attr.Type{
		"id": basetypes.StringType{},
		"region_configs": basetypes.ListType{
			ElemType: RegionConfigsValue{}.Type(ctx),
		},
		"zone_id":   basetypes.StringType{},
		"zone_name": basetypes.StringType{},
	}

	if v.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	}

	if v.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	}

	objVal, diags := types.ObjectValue(
		attributeTypes,
		map[string]attr.Value{
			"id":             v.Id,
			"region_configs": regionConfigs,
			"zone_id":        v.ZoneId,
			"zone_name":      v.ZoneName,
		})

	return objVal, diags
}

func (v ReplicationSpecsValue) Equal(o attr.Value) bool {
	other, ok := o.(ReplicationSpecsValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Id.Equal(other.Id) {
		return false
	}

	if !v.RegionConfigs.Equal(other.RegionConfigs) {
		return false
	}

	if !v.ZoneId.Equal(other.ZoneId) {
		return false
	}

	if !v.ZoneName.Equal(other.ZoneName) {
		return false
	}

	return true
}

func (v ReplicationSpecsValue) Type(ctx context.Context) attr.Type {
	return ReplicationSpecsType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v ReplicationSpecsValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"id": basetypes.StringType{},
		"region_configs": basetypes.ListType{
			ElemType: RegionConfigsValue{}.Type(ctx),
		},
		"zone_id":   basetypes.StringType{},
		"zone_name": basetypes.StringType{},
	}
}

var _ basetypes.ObjectTypable = RegionConfigsType{}

type RegionConfigsType struct {
	basetypes.ObjectType
}

func (t RegionConfigsType) Equal(o attr.Type) bool {
	other, ok := o.(RegionConfigsType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t RegionConfigsType) String() string {
	return "RegionConfigsType"
}

func (t RegionConfigsType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	analyticsAutoScalingAttribute, ok := attributes["analytics_auto_scaling"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`analytics_auto_scaling is missing from object`)

		return nil, diags
	}

	analyticsAutoScalingVal, ok := analyticsAutoScalingAttribute.(basetypes.ObjectValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`analytics_auto_scaling expected to be basetypes.ObjectValue, was: %T`, analyticsAutoScalingAttribute))
	}

	analyticsSpecsAttribute, ok := attributes["analytics_specs"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`analytics_specs is missing from object`)

		return nil, diags
	}

	analyticsSpecsVal, ok := analyticsSpecsAttribute.(basetypes.ObjectValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`analytics_specs expected to be basetypes.ObjectValue, was: %T`, analyticsSpecsAttribute))
	}

	autoScalingAttribute, ok := attributes["auto_scaling"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`auto_scaling is missing from object`)

		return nil, diags
	}

	autoScalingVal, ok := autoScalingAttribute.(basetypes.ObjectValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`auto_scaling expected to be basetypes.ObjectValue, was: %T`, autoScalingAttribute))
	}

	backingProviderNameAttribute, ok := attributes["backing_provider_name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`backing_provider_name is missing from object`)

		return nil, diags
	}

	backingProviderNameVal, ok := backingProviderNameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`backing_provider_name expected to be basetypes.StringValue, was: %T`, backingProviderNameAttribute))
	}

	electableSpecsAttribute, ok := attributes["electable_specs"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`electable_specs is missing from object`)

		return nil, diags
	}

	electableSpecsVal, ok := electableSpecsAttribute.(basetypes.ObjectValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`electable_specs expected to be basetypes.ObjectValue, was: %T`, electableSpecsAttribute))
	}

	priorityAttribute, ok := attributes["priority"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`priority is missing from object`)

		return nil, diags
	}

	priorityVal, ok := priorityAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`priority expected to be basetypes.Int64Value, was: %T`, priorityAttribute))
	}

	providerNameAttribute, ok := attributes["provider_name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`provider_name is missing from object`)

		return nil, diags
	}

	providerNameVal, ok := providerNameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`provider_name expected to be basetypes.StringValue, was: %T`, providerNameAttribute))
	}

	readOnlySpecsAttribute, ok := attributes["read_only_specs"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`read_only_specs is missing from object`)

		return nil, diags
	}

	readOnlySpecsVal, ok := readOnlySpecsAttribute.(basetypes.ObjectValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`read_only_specs expected to be basetypes.ObjectValue, was: %T`, readOnlySpecsAttribute))
	}

	regionNameAttribute, ok := attributes["region_name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`region_name is missing from object`)

		return nil, diags
	}

	regionNameVal, ok := regionNameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`region_name expected to be basetypes.StringValue, was: %T`, regionNameAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return RegionConfigsValue{
		AnalyticsAutoScaling: analyticsAutoScalingVal,
		AnalyticsSpecs:       analyticsSpecsVal,
		AutoScaling:          autoScalingVal,
		BackingProviderName:  backingProviderNameVal,
		ElectableSpecs:       electableSpecsVal,
		Priority:             priorityVal,
		ProviderName:         providerNameVal,
		ReadOnlySpecs:        readOnlySpecsVal,
		RegionName:           regionNameVal,
		state:                attr.ValueStateKnown,
	}, diags
}

func NewRegionConfigsValueNull() RegionConfigsValue {
	return RegionConfigsValue{
		state: attr.ValueStateNull,
	}
}

func NewRegionConfigsValueUnknown() RegionConfigsValue {
	return RegionConfigsValue{
		state: attr.ValueStateUnknown,
	}
}

func NewRegionConfigsValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (RegionConfigsValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing RegionConfigsValue Attribute Value",
				"While creating a RegionConfigsValue value, a missing attribute value was detected. "+
					"A RegionConfigsValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("RegionConfigsValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid RegionConfigsValue Attribute Type",
				"While creating a RegionConfigsValue value, an invalid attribute value was detected. "+
					"A RegionConfigsValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("RegionConfigsValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("RegionConfigsValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra RegionConfigsValue Attribute Value",
				"While creating a RegionConfigsValue value, an extra attribute value was detected. "+
					"A RegionConfigsValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra RegionConfigsValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewRegionConfigsValueUnknown(), diags
	}

	analyticsAutoScalingAttribute, ok := attributes["analytics_auto_scaling"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`analytics_auto_scaling is missing from object`)

		return NewRegionConfigsValueUnknown(), diags
	}

	analyticsAutoScalingVal, ok := analyticsAutoScalingAttribute.(basetypes.ObjectValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`analytics_auto_scaling expected to be basetypes.ObjectValue, was: %T`, analyticsAutoScalingAttribute))
	}

	analyticsSpecsAttribute, ok := attributes["analytics_specs"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`analytics_specs is missing from object`)

		return NewRegionConfigsValueUnknown(), diags
	}

	analyticsSpecsVal, ok := analyticsSpecsAttribute.(basetypes.ObjectValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`analytics_specs expected to be basetypes.ObjectValue, was: %T`, analyticsSpecsAttribute))
	}

	autoScalingAttribute, ok := attributes["auto_scaling"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`auto_scaling is missing from object`)

		return NewRegionConfigsValueUnknown(), diags
	}

	autoScalingVal, ok := autoScalingAttribute.(basetypes.ObjectValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`auto_scaling expected to be basetypes.ObjectValue, was: %T`, autoScalingAttribute))
	}

	backingProviderNameAttribute, ok := attributes["backing_provider_name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`backing_provider_name is missing from object`)

		return NewRegionConfigsValueUnknown(), diags
	}

	backingProviderNameVal, ok := backingProviderNameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`backing_provider_name expected to be basetypes.StringValue, was: %T`, backingProviderNameAttribute))
	}

	electableSpecsAttribute, ok := attributes["electable_specs"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`electable_specs is missing from object`)

		return NewRegionConfigsValueUnknown(), diags
	}

	electableSpecsVal, ok := electableSpecsAttribute.(basetypes.ObjectValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`electable_specs expected to be basetypes.ObjectValue, was: %T`, electableSpecsAttribute))
	}

	priorityAttribute, ok := attributes["priority"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`priority is missing from object`)

		return NewRegionConfigsValueUnknown(), diags
	}

	priorityVal, ok := priorityAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`priority expected to be basetypes.Int64Value, was: %T`, priorityAttribute))
	}

	providerNameAttribute, ok := attributes["provider_name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`provider_name is missing from object`)

		return NewRegionConfigsValueUnknown(), diags
	}

	providerNameVal, ok := providerNameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`provider_name expected to be basetypes.StringValue, was: %T`, providerNameAttribute))
	}

	readOnlySpecsAttribute, ok := attributes["read_only_specs"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`read_only_specs is missing from object`)

		return NewRegionConfigsValueUnknown(), diags
	}

	readOnlySpecsVal, ok := readOnlySpecsAttribute.(basetypes.ObjectValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`read_only_specs expected to be basetypes.ObjectValue, was: %T`, readOnlySpecsAttribute))
	}

	regionNameAttribute, ok := attributes["region_name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`region_name is missing from object`)

		return NewRegionConfigsValueUnknown(), diags
	}

	regionNameVal, ok := regionNameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`region_name expected to be basetypes.StringValue, was: %T`, regionNameAttribute))
	}

	if diags.HasError() {
		return NewRegionConfigsValueUnknown(), diags
	}

	return RegionConfigsValue{
		AnalyticsAutoScaling: analyticsAutoScalingVal,
		AnalyticsSpecs:       analyticsSpecsVal,
		AutoScaling:          autoScalingVal,
		BackingProviderName:  backingProviderNameVal,
		ElectableSpecs:       electableSpecsVal,
		Priority:             priorityVal,
		ProviderName:         providerNameVal,
		ReadOnlySpecs:        readOnlySpecsVal,
		RegionName:           regionNameVal,
		state:                attr.ValueStateKnown,
	}, diags
}

func NewRegionConfigsValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) RegionConfigsValue {
	object, diags := NewRegionConfigsValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewRegionConfigsValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t RegionConfigsType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewRegionConfigsValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewRegionConfigsValueUnknown(), nil
	}

	if in.IsNull() {
		return NewRegionConfigsValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewRegionConfigsValueMust(RegionConfigsValue{}.AttributeTypes(ctx), attributes), nil
}

func (t RegionConfigsType) ValueType(ctx context.Context) attr.Value {
	return RegionConfigsValue{}
}

var _ basetypes.ObjectValuable = RegionConfigsValue{}

type RegionConfigsValue struct {
	AnalyticsAutoScaling basetypes.ObjectValue `tfsdk:"analytics_auto_scaling"`
	AnalyticsSpecs       basetypes.ObjectValue `tfsdk:"analytics_specs"`
	AutoScaling          basetypes.ObjectValue `tfsdk:"auto_scaling"`
	BackingProviderName  basetypes.StringValue `tfsdk:"backing_provider_name"`
	ElectableSpecs       basetypes.ObjectValue `tfsdk:"electable_specs"`
	ProviderName         basetypes.StringValue `tfsdk:"provider_name"`
	ReadOnlySpecs        basetypes.ObjectValue `tfsdk:"read_only_specs"`
	RegionName           basetypes.StringValue `tfsdk:"region_name"`
	Priority             basetypes.Int64Value  `tfsdk:"priority"`
	state                attr.ValueState
}

func (v RegionConfigsValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 9)

	var val tftypes.Value
	var err error

	attrTypes["analytics_auto_scaling"] = basetypes.ObjectType{
		AttrTypes: AnalyticsAutoScalingValue{}.AttributeTypes(ctx),
	}.TerraformType(ctx)
	attrTypes["analytics_specs"] = basetypes.ObjectType{
		AttrTypes: AnalyticsSpecsValue{}.AttributeTypes(ctx),
	}.TerraformType(ctx)
	attrTypes["auto_scaling"] = basetypes.ObjectType{
		AttrTypes: AutoScalingValue{}.AttributeTypes(ctx),
	}.TerraformType(ctx)
	attrTypes["backing_provider_name"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["electable_specs"] = basetypes.ObjectType{
		AttrTypes: ElectableSpecsValue{}.AttributeTypes(ctx),
	}.TerraformType(ctx)
	attrTypes["priority"] = basetypes.Int64Type{}.TerraformType(ctx)
	attrTypes["provider_name"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["read_only_specs"] = basetypes.ObjectType{
		AttrTypes: ReadOnlySpecsValue{}.AttributeTypes(ctx),
	}.TerraformType(ctx)
	attrTypes["region_name"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 9)

		val, err = v.AnalyticsAutoScaling.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["analytics_auto_scaling"] = val

		val, err = v.AnalyticsSpecs.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["analytics_specs"] = val

		val, err = v.AutoScaling.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["auto_scaling"] = val

		val, err = v.BackingProviderName.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["backing_provider_name"] = val

		val, err = v.ElectableSpecs.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["electable_specs"] = val

		val, err = v.Priority.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["priority"] = val

		val, err = v.ProviderName.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["provider_name"] = val

		val, err = v.ReadOnlySpecs.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["read_only_specs"] = val

		val, err = v.RegionName.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["region_name"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v RegionConfigsValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v RegionConfigsValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v RegionConfigsValue) String() string {
	return "RegionConfigsValue"
}

func (v RegionConfigsValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	var analyticsAutoScaling basetypes.ObjectValue

	if v.AnalyticsAutoScaling.IsNull() {
		analyticsAutoScaling = types.ObjectNull(
			AnalyticsAutoScalingValue{}.AttributeTypes(ctx),
		)
	}

	if v.AnalyticsAutoScaling.IsUnknown() {
		analyticsAutoScaling = types.ObjectUnknown(
			AnalyticsAutoScalingValue{}.AttributeTypes(ctx),
		)
	}

	if !v.AnalyticsAutoScaling.IsNull() && !v.AnalyticsAutoScaling.IsUnknown() {
		analyticsAutoScaling = types.ObjectValueMust(
			AnalyticsAutoScalingValue{}.AttributeTypes(ctx),
			v.AnalyticsAutoScaling.Attributes(),
		)
	}

	var analyticsSpecs basetypes.ObjectValue

	if v.AnalyticsSpecs.IsNull() {
		analyticsSpecs = types.ObjectNull(
			AnalyticsSpecsValue{}.AttributeTypes(ctx),
		)
	}

	if v.AnalyticsSpecs.IsUnknown() {
		analyticsSpecs = types.ObjectUnknown(
			AnalyticsSpecsValue{}.AttributeTypes(ctx),
		)
	}

	if !v.AnalyticsSpecs.IsNull() && !v.AnalyticsSpecs.IsUnknown() {
		analyticsSpecs = types.ObjectValueMust(
			AnalyticsSpecsValue{}.AttributeTypes(ctx),
			v.AnalyticsSpecs.Attributes(),
		)
	}

	var autoScaling basetypes.ObjectValue

	if v.AutoScaling.IsNull() {
		autoScaling = types.ObjectNull(
			AutoScalingValue{}.AttributeTypes(ctx),
		)
	}

	if v.AutoScaling.IsUnknown() {
		autoScaling = types.ObjectUnknown(
			AutoScalingValue{}.AttributeTypes(ctx),
		)
	}

	if !v.AutoScaling.IsNull() && !v.AutoScaling.IsUnknown() {
		autoScaling = types.ObjectValueMust(
			AutoScalingValue{}.AttributeTypes(ctx),
			v.AutoScaling.Attributes(),
		)
	}

	var electableSpecs basetypes.ObjectValue

	if v.ElectableSpecs.IsNull() {
		electableSpecs = types.ObjectNull(
			ElectableSpecsValue{}.AttributeTypes(ctx),
		)
	}

	if v.ElectableSpecs.IsUnknown() {
		electableSpecs = types.ObjectUnknown(
			ElectableSpecsValue{}.AttributeTypes(ctx),
		)
	}

	if !v.ElectableSpecs.IsNull() && !v.ElectableSpecs.IsUnknown() {
		electableSpecs = types.ObjectValueMust(
			ElectableSpecsValue{}.AttributeTypes(ctx),
			v.ElectableSpecs.Attributes(),
		)
	}

	var readOnlySpecs basetypes.ObjectValue

	if v.ReadOnlySpecs.IsNull() {
		readOnlySpecs = types.ObjectNull(
			ReadOnlySpecsValue{}.AttributeTypes(ctx),
		)
	}

	if v.ReadOnlySpecs.IsUnknown() {
		readOnlySpecs = types.ObjectUnknown(
			ReadOnlySpecsValue{}.AttributeTypes(ctx),
		)
	}

	if !v.ReadOnlySpecs.IsNull() && !v.ReadOnlySpecs.IsUnknown() {
		readOnlySpecs = types.ObjectValueMust(
			ReadOnlySpecsValue{}.AttributeTypes(ctx),
			v.ReadOnlySpecs.Attributes(),
		)
	}

	attributeTypes := map[string]attr.Type{
		"analytics_auto_scaling": basetypes.ObjectType{
			AttrTypes: AnalyticsAutoScalingValue{}.AttributeTypes(ctx),
		},
		"analytics_specs": basetypes.ObjectType{
			AttrTypes: AnalyticsSpecsValue{}.AttributeTypes(ctx),
		},
		"auto_scaling": basetypes.ObjectType{
			AttrTypes: AutoScalingValue{}.AttributeTypes(ctx),
		},
		"backing_provider_name": basetypes.StringType{},
		"electable_specs": basetypes.ObjectType{
			AttrTypes: ElectableSpecsValue{}.AttributeTypes(ctx),
		},
		"priority":      basetypes.Int64Type{},
		"provider_name": basetypes.StringType{},
		"read_only_specs": basetypes.ObjectType{
			AttrTypes: ReadOnlySpecsValue{}.AttributeTypes(ctx),
		},
		"region_name": basetypes.StringType{},
	}

	if v.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	}

	if v.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	}

	objVal, diags := types.ObjectValue(
		attributeTypes,
		map[string]attr.Value{
			"analytics_auto_scaling": analyticsAutoScaling,
			"analytics_specs":        analyticsSpecs,
			"auto_scaling":           autoScaling,
			"backing_provider_name":  v.BackingProviderName,
			"electable_specs":        electableSpecs,
			"priority":               v.Priority,
			"provider_name":          v.ProviderName,
			"read_only_specs":        readOnlySpecs,
			"region_name":            v.RegionName,
		})

	return objVal, diags
}

func (v RegionConfigsValue) Equal(o attr.Value) bool {
	other, ok := o.(RegionConfigsValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.AnalyticsAutoScaling.Equal(other.AnalyticsAutoScaling) {
		return false
	}

	if !v.AnalyticsSpecs.Equal(other.AnalyticsSpecs) {
		return false
	}

	if !v.AutoScaling.Equal(other.AutoScaling) {
		return false
	}

	if !v.BackingProviderName.Equal(other.BackingProviderName) {
		return false
	}

	if !v.ElectableSpecs.Equal(other.ElectableSpecs) {
		return false
	}

	if !v.Priority.Equal(other.Priority) {
		return false
	}

	if !v.ProviderName.Equal(other.ProviderName) {
		return false
	}

	if !v.ReadOnlySpecs.Equal(other.ReadOnlySpecs) {
		return false
	}

	if !v.RegionName.Equal(other.RegionName) {
		return false
	}

	return true
}

func (v RegionConfigsValue) Type(ctx context.Context) attr.Type {
	return RegionConfigsType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v RegionConfigsValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"analytics_auto_scaling": basetypes.ObjectType{
			AttrTypes: AnalyticsAutoScalingValue{}.AttributeTypes(ctx),
		},
		"analytics_specs": basetypes.ObjectType{
			AttrTypes: AnalyticsSpecsValue{}.AttributeTypes(ctx),
		},
		"auto_scaling": basetypes.ObjectType{
			AttrTypes: AutoScalingValue{}.AttributeTypes(ctx),
		},
		"backing_provider_name": basetypes.StringType{},
		"electable_specs": basetypes.ObjectType{
			AttrTypes: ElectableSpecsValue{}.AttributeTypes(ctx),
		},
		"priority":      basetypes.Int64Type{},
		"provider_name": basetypes.StringType{},
		"read_only_specs": basetypes.ObjectType{
			AttrTypes: ReadOnlySpecsValue{}.AttributeTypes(ctx),
		},
		"region_name": basetypes.StringType{},
	}
}

var _ basetypes.ObjectTypable = AnalyticsAutoScalingType{}

type AnalyticsAutoScalingType struct {
	basetypes.ObjectType
}

func (t AnalyticsAutoScalingType) Equal(o attr.Type) bool {
	other, ok := o.(AnalyticsAutoScalingType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t AnalyticsAutoScalingType) String() string {
	return "AnalyticsAutoScalingType"
}

func (t AnalyticsAutoScalingType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	computeAttribute, ok := attributes["compute"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`compute is missing from object`)

		return nil, diags
	}

	computeVal, ok := computeAttribute.(basetypes.ObjectValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`compute expected to be basetypes.ObjectValue, was: %T`, computeAttribute))
	}

	diskGbAttribute, ok := attributes["disk_gb"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`disk_gb is missing from object`)

		return nil, diags
	}

	diskGbVal, ok := diskGbAttribute.(basetypes.ObjectValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`disk_gb expected to be basetypes.ObjectValue, was: %T`, diskGbAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return AnalyticsAutoScalingValue{
		Compute: computeVal,
		DiskGb:  diskGbVal,
		state:   attr.ValueStateKnown,
	}, diags
}

func NewAnalyticsAutoScalingValueNull() AnalyticsAutoScalingValue {
	return AnalyticsAutoScalingValue{
		state: attr.ValueStateNull,
	}
}

func NewAnalyticsAutoScalingValueUnknown() AnalyticsAutoScalingValue {
	return AnalyticsAutoScalingValue{
		state: attr.ValueStateUnknown,
	}
}

func NewAnalyticsAutoScalingValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (AnalyticsAutoScalingValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing AnalyticsAutoScalingValue Attribute Value",
				"While creating a AnalyticsAutoScalingValue value, a missing attribute value was detected. "+
					"A AnalyticsAutoScalingValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("AnalyticsAutoScalingValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid AnalyticsAutoScalingValue Attribute Type",
				"While creating a AnalyticsAutoScalingValue value, an invalid attribute value was detected. "+
					"A AnalyticsAutoScalingValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("AnalyticsAutoScalingValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("AnalyticsAutoScalingValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra AnalyticsAutoScalingValue Attribute Value",
				"While creating a AnalyticsAutoScalingValue value, an extra attribute value was detected. "+
					"A AnalyticsAutoScalingValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra AnalyticsAutoScalingValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewAnalyticsAutoScalingValueUnknown(), diags
	}

	computeAttribute, ok := attributes["compute"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`compute is missing from object`)

		return NewAnalyticsAutoScalingValueUnknown(), diags
	}

	computeVal, ok := computeAttribute.(basetypes.ObjectValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`compute expected to be basetypes.ObjectValue, was: %T`, computeAttribute))
	}

	diskGbAttribute, ok := attributes["disk_gb"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`disk_gb is missing from object`)

		return NewAnalyticsAutoScalingValueUnknown(), diags
	}

	diskGbVal, ok := diskGbAttribute.(basetypes.ObjectValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`disk_gb expected to be basetypes.ObjectValue, was: %T`, diskGbAttribute))
	}

	if diags.HasError() {
		return NewAnalyticsAutoScalingValueUnknown(), diags
	}

	return AnalyticsAutoScalingValue{
		Compute: computeVal,
		DiskGb:  diskGbVal,
		state:   attr.ValueStateKnown,
	}, diags
}

func NewAnalyticsAutoScalingValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) AnalyticsAutoScalingValue {
	object, diags := NewAnalyticsAutoScalingValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewAnalyticsAutoScalingValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t AnalyticsAutoScalingType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewAnalyticsAutoScalingValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewAnalyticsAutoScalingValueUnknown(), nil
	}

	if in.IsNull() {
		return NewAnalyticsAutoScalingValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewAnalyticsAutoScalingValueMust(AnalyticsAutoScalingValue{}.AttributeTypes(ctx), attributes), nil
}

func (t AnalyticsAutoScalingType) ValueType(ctx context.Context) attr.Value {
	return AnalyticsAutoScalingValue{}
}

var _ basetypes.ObjectValuable = AnalyticsAutoScalingValue{}

type AnalyticsAutoScalingValue struct {
	Compute basetypes.ObjectValue `tfsdk:"compute"`
	DiskGb  basetypes.ObjectValue `tfsdk:"disk_gb"`
	state   attr.ValueState
}

func (v AnalyticsAutoScalingValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 2)

	var val tftypes.Value
	var err error

	attrTypes["compute"] = basetypes.ObjectType{
		AttrTypes: ComputeValue{}.AttributeTypes(ctx),
	}.TerraformType(ctx)
	attrTypes["disk_gb"] = basetypes.ObjectType{
		AttrTypes: DiskGbValue{}.AttributeTypes(ctx),
	}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 2)

		val, err = v.Compute.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["compute"] = val

		val, err = v.DiskGb.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["disk_gb"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v AnalyticsAutoScalingValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v AnalyticsAutoScalingValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v AnalyticsAutoScalingValue) String() string {
	return "AnalyticsAutoScalingValue"
}

func (v AnalyticsAutoScalingValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	var compute basetypes.ObjectValue

	if v.Compute.IsNull() {
		compute = types.ObjectNull(
			ComputeValue{}.AttributeTypes(ctx),
		)
	}

	if v.Compute.IsUnknown() {
		compute = types.ObjectUnknown(
			ComputeValue{}.AttributeTypes(ctx),
		)
	}

	if !v.Compute.IsNull() && !v.Compute.IsUnknown() {
		compute = types.ObjectValueMust(
			ComputeValue{}.AttributeTypes(ctx),
			v.Compute.Attributes(),
		)
	}

	var diskGb basetypes.ObjectValue

	if v.DiskGb.IsNull() {
		diskGb = types.ObjectNull(
			DiskGbValue{}.AttributeTypes(ctx),
		)
	}

	if v.DiskGb.IsUnknown() {
		diskGb = types.ObjectUnknown(
			DiskGbValue{}.AttributeTypes(ctx),
		)
	}

	if !v.DiskGb.IsNull() && !v.DiskGb.IsUnknown() {
		diskGb = types.ObjectValueMust(
			DiskGbValue{}.AttributeTypes(ctx),
			v.DiskGb.Attributes(),
		)
	}

	attributeTypes := map[string]attr.Type{
		"compute": basetypes.ObjectType{
			AttrTypes: ComputeValue{}.AttributeTypes(ctx),
		},
		"disk_gb": basetypes.ObjectType{
			AttrTypes: DiskGbValue{}.AttributeTypes(ctx),
		},
	}

	if v.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	}

	if v.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	}

	objVal, diags := types.ObjectValue(
		attributeTypes,
		map[string]attr.Value{
			"compute": compute,
			"disk_gb": diskGb,
		})

	return objVal, diags
}

func (v AnalyticsAutoScalingValue) Equal(o attr.Value) bool {
	other, ok := o.(AnalyticsAutoScalingValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Compute.Equal(other.Compute) {
		return false
	}

	if !v.DiskGb.Equal(other.DiskGb) {
		return false
	}

	return true
}

func (v AnalyticsAutoScalingValue) Type(ctx context.Context) attr.Type {
	return AnalyticsAutoScalingType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v AnalyticsAutoScalingValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"compute": basetypes.ObjectType{
			AttrTypes: ComputeValue{}.AttributeTypes(ctx),
		},
		"disk_gb": basetypes.ObjectType{
			AttrTypes: DiskGbValue{}.AttributeTypes(ctx),
		},
	}
}

var _ basetypes.ObjectTypable = ComputeType{}

type ComputeType struct {
	basetypes.ObjectType
}

func (t ComputeType) Equal(o attr.Type) bool {
	other, ok := o.(ComputeType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t ComputeType) String() string {
	return "ComputeType"
}

func (t ComputeType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	enabledAttribute, ok := attributes["enabled"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`enabled is missing from object`)

		return nil, diags
	}

	enabledVal, ok := enabledAttribute.(basetypes.BoolValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`enabled expected to be basetypes.BoolValue, was: %T`, enabledAttribute))
	}

	maxInstanceSizeAttribute, ok := attributes["max_instance_size"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`max_instance_size is missing from object`)

		return nil, diags
	}

	maxInstanceSizeVal, ok := maxInstanceSizeAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`max_instance_size expected to be basetypes.StringValue, was: %T`, maxInstanceSizeAttribute))
	}

	minInstanceSizeAttribute, ok := attributes["min_instance_size"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`min_instance_size is missing from object`)

		return nil, diags
	}

	minInstanceSizeVal, ok := minInstanceSizeAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`min_instance_size expected to be basetypes.StringValue, was: %T`, minInstanceSizeAttribute))
	}

	scaleDownEnabledAttribute, ok := attributes["scale_down_enabled"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`scale_down_enabled is missing from object`)

		return nil, diags
	}

	scaleDownEnabledVal, ok := scaleDownEnabledAttribute.(basetypes.BoolValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`scale_down_enabled expected to be basetypes.BoolValue, was: %T`, scaleDownEnabledAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return ComputeValue{
		Enabled:          enabledVal,
		MaxInstanceSize:  maxInstanceSizeVal,
		MinInstanceSize:  minInstanceSizeVal,
		ScaleDownEnabled: scaleDownEnabledVal,
		state:            attr.ValueStateKnown,
	}, diags
}

func NewComputeValueNull() ComputeValue {
	return ComputeValue{
		state: attr.ValueStateNull,
	}
}

func NewComputeValueUnknown() ComputeValue {
	return ComputeValue{
		state: attr.ValueStateUnknown,
	}
}

func NewComputeValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (ComputeValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing ComputeValue Attribute Value",
				"While creating a ComputeValue value, a missing attribute value was detected. "+
					"A ComputeValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("ComputeValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid ComputeValue Attribute Type",
				"While creating a ComputeValue value, an invalid attribute value was detected. "+
					"A ComputeValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("ComputeValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("ComputeValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra ComputeValue Attribute Value",
				"While creating a ComputeValue value, an extra attribute value was detected. "+
					"A ComputeValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra ComputeValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewComputeValueUnknown(), diags
	}

	enabledAttribute, ok := attributes["enabled"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`enabled is missing from object`)

		return NewComputeValueUnknown(), diags
	}

	enabledVal, ok := enabledAttribute.(basetypes.BoolValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`enabled expected to be basetypes.BoolValue, was: %T`, enabledAttribute))
	}

	maxInstanceSizeAttribute, ok := attributes["max_instance_size"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`max_instance_size is missing from object`)

		return NewComputeValueUnknown(), diags
	}

	maxInstanceSizeVal, ok := maxInstanceSizeAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`max_instance_size expected to be basetypes.StringValue, was: %T`, maxInstanceSizeAttribute))
	}

	minInstanceSizeAttribute, ok := attributes["min_instance_size"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`min_instance_size is missing from object`)

		return NewComputeValueUnknown(), diags
	}

	minInstanceSizeVal, ok := minInstanceSizeAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`min_instance_size expected to be basetypes.StringValue, was: %T`, minInstanceSizeAttribute))
	}

	scaleDownEnabledAttribute, ok := attributes["scale_down_enabled"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`scale_down_enabled is missing from object`)

		return NewComputeValueUnknown(), diags
	}

	scaleDownEnabledVal, ok := scaleDownEnabledAttribute.(basetypes.BoolValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`scale_down_enabled expected to be basetypes.BoolValue, was: %T`, scaleDownEnabledAttribute))
	}

	if diags.HasError() {
		return NewComputeValueUnknown(), diags
	}

	return ComputeValue{
		Enabled:          enabledVal,
		MaxInstanceSize:  maxInstanceSizeVal,
		MinInstanceSize:  minInstanceSizeVal,
		ScaleDownEnabled: scaleDownEnabledVal,
		state:            attr.ValueStateKnown,
	}, diags
}

func NewComputeValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) ComputeValue {
	object, diags := NewComputeValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewComputeValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t ComputeType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewComputeValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewComputeValueUnknown(), nil
	}

	if in.IsNull() {
		return NewComputeValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewComputeValueMust(ComputeValue{}.AttributeTypes(ctx), attributes), nil
}

func (t ComputeType) ValueType(ctx context.Context) attr.Value {
	return ComputeValue{}
}

var _ basetypes.ObjectValuable = ComputeValue{}

type ComputeValue struct {
	MaxInstanceSize  basetypes.StringValue `tfsdk:"max_instance_size"`
	MinInstanceSize  basetypes.StringValue `tfsdk:"min_instance_size"`
	Enabled          basetypes.BoolValue   `tfsdk:"enabled"`
	ScaleDownEnabled basetypes.BoolValue   `tfsdk:"scale_down_enabled"`
	state            attr.ValueState
}

func (v ComputeValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 4)

	var val tftypes.Value
	var err error

	attrTypes["enabled"] = basetypes.BoolType{}.TerraformType(ctx)
	attrTypes["max_instance_size"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["min_instance_size"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["scale_down_enabled"] = basetypes.BoolType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 4)

		val, err = v.Enabled.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["enabled"] = val

		val, err = v.MaxInstanceSize.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["max_instance_size"] = val

		val, err = v.MinInstanceSize.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["min_instance_size"] = val

		val, err = v.ScaleDownEnabled.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["scale_down_enabled"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v ComputeValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v ComputeValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v ComputeValue) String() string {
	return "ComputeValue"
}

func (v ComputeValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := map[string]attr.Type{
		"enabled":            basetypes.BoolType{},
		"max_instance_size":  basetypes.StringType{},
		"min_instance_size":  basetypes.StringType{},
		"scale_down_enabled": basetypes.BoolType{},
	}

	if v.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	}

	if v.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	}

	objVal, diags := types.ObjectValue(
		attributeTypes,
		map[string]attr.Value{
			"enabled":            v.Enabled,
			"max_instance_size":  v.MaxInstanceSize,
			"min_instance_size":  v.MinInstanceSize,
			"scale_down_enabled": v.ScaleDownEnabled,
		})

	return objVal, diags
}

func (v ComputeValue) Equal(o attr.Value) bool {
	other, ok := o.(ComputeValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Enabled.Equal(other.Enabled) {
		return false
	}

	if !v.MaxInstanceSize.Equal(other.MaxInstanceSize) {
		return false
	}

	if !v.MinInstanceSize.Equal(other.MinInstanceSize) {
		return false
	}

	if !v.ScaleDownEnabled.Equal(other.ScaleDownEnabled) {
		return false
	}

	return true
}

func (v ComputeValue) Type(ctx context.Context) attr.Type {
	return ComputeType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v ComputeValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"enabled":            basetypes.BoolType{},
		"max_instance_size":  basetypes.StringType{},
		"min_instance_size":  basetypes.StringType{},
		"scale_down_enabled": basetypes.BoolType{},
	}
}

var _ basetypes.ObjectTypable = DiskGbType{}

type DiskGbType struct {
	basetypes.ObjectType
}

func (t DiskGbType) Equal(o attr.Type) bool {
	other, ok := o.(DiskGbType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t DiskGbType) String() string {
	return "DiskGbType"
}

func (t DiskGbType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	enabledAttribute, ok := attributes["enabled"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`enabled is missing from object`)

		return nil, diags
	}

	enabledVal, ok := enabledAttribute.(basetypes.BoolValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`enabled expected to be basetypes.BoolValue, was: %T`, enabledAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return DiskGbValue{
		Enabled: enabledVal,
		state:   attr.ValueStateKnown,
	}, diags
}

func NewDiskGbValueNull() DiskGbValue {
	return DiskGbValue{
		state: attr.ValueStateNull,
	}
}

func NewDiskGbValueUnknown() DiskGbValue {
	return DiskGbValue{
		state: attr.ValueStateUnknown,
	}
}

func NewDiskGbValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (DiskGbValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing DiskGbValue Attribute Value",
				"While creating a DiskGbValue value, a missing attribute value was detected. "+
					"A DiskGbValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("DiskGbValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid DiskGbValue Attribute Type",
				"While creating a DiskGbValue value, an invalid attribute value was detected. "+
					"A DiskGbValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("DiskGbValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("DiskGbValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra DiskGbValue Attribute Value",
				"While creating a DiskGbValue value, an extra attribute value was detected. "+
					"A DiskGbValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra DiskGbValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewDiskGbValueUnknown(), diags
	}

	enabledAttribute, ok := attributes["enabled"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`enabled is missing from object`)

		return NewDiskGbValueUnknown(), diags
	}

	enabledVal, ok := enabledAttribute.(basetypes.BoolValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`enabled expected to be basetypes.BoolValue, was: %T`, enabledAttribute))
	}

	if diags.HasError() {
		return NewDiskGbValueUnknown(), diags
	}

	return DiskGbValue{
		Enabled: enabledVal,
		state:   attr.ValueStateKnown,
	}, diags
}

func NewDiskGbValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) DiskGbValue {
	object, diags := NewDiskGbValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewDiskGbValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t DiskGbType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewDiskGbValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewDiskGbValueUnknown(), nil
	}

	if in.IsNull() {
		return NewDiskGbValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewDiskGbValueMust(DiskGbValue{}.AttributeTypes(ctx), attributes), nil
}

func (t DiskGbType) ValueType(ctx context.Context) attr.Value {
	return DiskGbValue{}
}

var _ basetypes.ObjectValuable = DiskGbValue{}

type DiskGbValue struct {
	Enabled basetypes.BoolValue `tfsdk:"enabled"`
	state   attr.ValueState
}

func (v DiskGbValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 1)

	var val tftypes.Value
	var err error

	attrTypes["enabled"] = basetypes.BoolType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 1)

		val, err = v.Enabled.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["enabled"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v DiskGbValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v DiskGbValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v DiskGbValue) String() string {
	return "DiskGbValue"
}

func (v DiskGbValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := map[string]attr.Type{
		"enabled": basetypes.BoolType{},
	}

	if v.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	}

	if v.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	}

	objVal, diags := types.ObjectValue(
		attributeTypes,
		map[string]attr.Value{
			"enabled": v.Enabled,
		})

	return objVal, diags
}

func (v DiskGbValue) Equal(o attr.Value) bool {
	other, ok := o.(DiskGbValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Enabled.Equal(other.Enabled) {
		return false
	}

	return true
}

func (v DiskGbValue) Type(ctx context.Context) attr.Type {
	return DiskGbType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v DiskGbValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"enabled": basetypes.BoolType{},
	}
}

var _ basetypes.ObjectTypable = AnalyticsSpecsType{}

type AnalyticsSpecsType struct {
	basetypes.ObjectType
}

func (t AnalyticsSpecsType) Equal(o attr.Type) bool {
	other, ok := o.(AnalyticsSpecsType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t AnalyticsSpecsType) String() string {
	return "AnalyticsSpecsType"
}

func (t AnalyticsSpecsType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	diskIopsAttribute, ok := attributes["disk_iops"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`disk_iops is missing from object`)

		return nil, diags
	}

	diskIopsVal, ok := diskIopsAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`disk_iops expected to be basetypes.Int64Value, was: %T`, diskIopsAttribute))
	}

	diskSizeGbAttribute, ok := attributes["disk_size_gb"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`disk_size_gb is missing from object`)

		return nil, diags
	}

	diskSizeGbVal, ok := diskSizeGbAttribute.(basetypes.Float64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`disk_size_gb expected to be basetypes.Float64Value, was: %T`, diskSizeGbAttribute))
	}

	ebsVolumeTypeAttribute, ok := attributes["ebs_volume_type"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`ebs_volume_type is missing from object`)

		return nil, diags
	}

	ebsVolumeTypeVal, ok := ebsVolumeTypeAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`ebs_volume_type expected to be basetypes.StringValue, was: %T`, ebsVolumeTypeAttribute))
	}

	instanceSizeAttribute, ok := attributes["instance_size"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`instance_size is missing from object`)

		return nil, diags
	}

	instanceSizeVal, ok := instanceSizeAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`instance_size expected to be basetypes.StringValue, was: %T`, instanceSizeAttribute))
	}

	nodeCountAttribute, ok := attributes["node_count"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`node_count is missing from object`)

		return nil, diags
	}

	nodeCountVal, ok := nodeCountAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`node_count expected to be basetypes.Int64Value, was: %T`, nodeCountAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return AnalyticsSpecsValue{
		DiskIops:      diskIopsVal,
		DiskSizeGb:    diskSizeGbVal,
		EbsVolumeType: ebsVolumeTypeVal,
		InstanceSize:  instanceSizeVal,
		NodeCount:     nodeCountVal,
		state:         attr.ValueStateKnown,
	}, diags
}

func NewAnalyticsSpecsValueNull() AnalyticsSpecsValue {
	return AnalyticsSpecsValue{
		state: attr.ValueStateNull,
	}
}

func NewAnalyticsSpecsValueUnknown() AnalyticsSpecsValue {
	return AnalyticsSpecsValue{
		state: attr.ValueStateUnknown,
	}
}

func NewAnalyticsSpecsValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (AnalyticsSpecsValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing AnalyticsSpecsValue Attribute Value",
				"While creating a AnalyticsSpecsValue value, a missing attribute value was detected. "+
					"A AnalyticsSpecsValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("AnalyticsSpecsValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid AnalyticsSpecsValue Attribute Type",
				"While creating a AnalyticsSpecsValue value, an invalid attribute value was detected. "+
					"A AnalyticsSpecsValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("AnalyticsSpecsValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("AnalyticsSpecsValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra AnalyticsSpecsValue Attribute Value",
				"While creating a AnalyticsSpecsValue value, an extra attribute value was detected. "+
					"A AnalyticsSpecsValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra AnalyticsSpecsValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewAnalyticsSpecsValueUnknown(), diags
	}

	diskIopsAttribute, ok := attributes["disk_iops"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`disk_iops is missing from object`)

		return NewAnalyticsSpecsValueUnknown(), diags
	}

	diskIopsVal, ok := diskIopsAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`disk_iops expected to be basetypes.Int64Value, was: %T`, diskIopsAttribute))
	}

	diskSizeGbAttribute, ok := attributes["disk_size_gb"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`disk_size_gb is missing from object`)

		return NewAnalyticsSpecsValueUnknown(), diags
	}

	diskSizeGbVal, ok := diskSizeGbAttribute.(basetypes.Float64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`disk_size_gb expected to be basetypes.Float64Value, was: %T`, diskSizeGbAttribute))
	}

	ebsVolumeTypeAttribute, ok := attributes["ebs_volume_type"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`ebs_volume_type is missing from object`)

		return NewAnalyticsSpecsValueUnknown(), diags
	}

	ebsVolumeTypeVal, ok := ebsVolumeTypeAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`ebs_volume_type expected to be basetypes.StringValue, was: %T`, ebsVolumeTypeAttribute))
	}

	instanceSizeAttribute, ok := attributes["instance_size"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`instance_size is missing from object`)

		return NewAnalyticsSpecsValueUnknown(), diags
	}

	instanceSizeVal, ok := instanceSizeAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`instance_size expected to be basetypes.StringValue, was: %T`, instanceSizeAttribute))
	}

	nodeCountAttribute, ok := attributes["node_count"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`node_count is missing from object`)

		return NewAnalyticsSpecsValueUnknown(), diags
	}

	nodeCountVal, ok := nodeCountAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`node_count expected to be basetypes.Int64Value, was: %T`, nodeCountAttribute))
	}

	if diags.HasError() {
		return NewAnalyticsSpecsValueUnknown(), diags
	}

	return AnalyticsSpecsValue{
		DiskIops:      diskIopsVal,
		DiskSizeGb:    diskSizeGbVal,
		EbsVolumeType: ebsVolumeTypeVal,
		InstanceSize:  instanceSizeVal,
		NodeCount:     nodeCountVal,
		state:         attr.ValueStateKnown,
	}, diags
}

func NewAnalyticsSpecsValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) AnalyticsSpecsValue {
	object, diags := NewAnalyticsSpecsValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewAnalyticsSpecsValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t AnalyticsSpecsType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewAnalyticsSpecsValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewAnalyticsSpecsValueUnknown(), nil
	}

	if in.IsNull() {
		return NewAnalyticsSpecsValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewAnalyticsSpecsValueMust(AnalyticsSpecsValue{}.AttributeTypes(ctx), attributes), nil
}

func (t AnalyticsSpecsType) ValueType(ctx context.Context) attr.Value {
	return AnalyticsSpecsValue{}
}

var _ basetypes.ObjectValuable = AnalyticsSpecsValue{}

type AnalyticsSpecsValue struct {
	DiskSizeGb    basetypes.Float64Value `tfsdk:"disk_size_gb"`
	EbsVolumeType basetypes.StringValue  `tfsdk:"ebs_volume_type"`
	InstanceSize  basetypes.StringValue  `tfsdk:"instance_size"`
	DiskIops      basetypes.Int64Value   `tfsdk:"disk_iops"`
	NodeCount     basetypes.Int64Value   `tfsdk:"node_count"`
	state         attr.ValueState
}

func (v AnalyticsSpecsValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 5)

	var val tftypes.Value
	var err error

	attrTypes["disk_iops"] = basetypes.Int64Type{}.TerraformType(ctx)
	attrTypes["disk_size_gb"] = basetypes.Float64Type{}.TerraformType(ctx)
	attrTypes["ebs_volume_type"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["instance_size"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["node_count"] = basetypes.Int64Type{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 5)

		val, err = v.DiskIops.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["disk_iops"] = val

		val, err = v.DiskSizeGb.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["disk_size_gb"] = val

		val, err = v.EbsVolumeType.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["ebs_volume_type"] = val

		val, err = v.InstanceSize.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["instance_size"] = val

		val, err = v.NodeCount.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["node_count"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v AnalyticsSpecsValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v AnalyticsSpecsValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v AnalyticsSpecsValue) String() string {
	return "AnalyticsSpecsValue"
}

func (v AnalyticsSpecsValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := map[string]attr.Type{
		"disk_iops":       basetypes.Int64Type{},
		"disk_size_gb":    basetypes.Float64Type{},
		"ebs_volume_type": basetypes.StringType{},
		"instance_size":   basetypes.StringType{},
		"node_count":      basetypes.Int64Type{},
	}

	if v.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	}

	if v.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	}

	objVal, diags := types.ObjectValue(
		attributeTypes,
		map[string]attr.Value{
			"disk_iops":       v.DiskIops,
			"disk_size_gb":    v.DiskSizeGb,
			"ebs_volume_type": v.EbsVolumeType,
			"instance_size":   v.InstanceSize,
			"node_count":      v.NodeCount,
		})

	return objVal, diags
}

func (v AnalyticsSpecsValue) Equal(o attr.Value) bool {
	other, ok := o.(AnalyticsSpecsValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.DiskIops.Equal(other.DiskIops) {
		return false
	}

	if !v.DiskSizeGb.Equal(other.DiskSizeGb) {
		return false
	}

	if !v.EbsVolumeType.Equal(other.EbsVolumeType) {
		return false
	}

	if !v.InstanceSize.Equal(other.InstanceSize) {
		return false
	}

	if !v.NodeCount.Equal(other.NodeCount) {
		return false
	}

	return true
}

func (v AnalyticsSpecsValue) Type(ctx context.Context) attr.Type {
	return AnalyticsSpecsType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v AnalyticsSpecsValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"disk_iops":       basetypes.Int64Type{},
		"disk_size_gb":    basetypes.Float64Type{},
		"ebs_volume_type": basetypes.StringType{},
		"instance_size":   basetypes.StringType{},
		"node_count":      basetypes.Int64Type{},
	}
}

var _ basetypes.ObjectTypable = AutoScalingType{}

type AutoScalingType struct {
	basetypes.ObjectType
}

func (t AutoScalingType) Equal(o attr.Type) bool {
	other, ok := o.(AutoScalingType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t AutoScalingType) String() string {
	return "AutoScalingType"
}

func (t AutoScalingType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	computeAttribute, ok := attributes["compute"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`compute is missing from object`)

		return nil, diags
	}

	computeVal, ok := computeAttribute.(basetypes.ObjectValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`compute expected to be basetypes.ObjectValue, was: %T`, computeAttribute))
	}

	diskGbAttribute, ok := attributes["disk_gb"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`disk_gb is missing from object`)

		return nil, diags
	}

	diskGbVal, ok := diskGbAttribute.(basetypes.ObjectValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`disk_gb expected to be basetypes.ObjectValue, was: %T`, diskGbAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return AutoScalingValue{
		Compute: computeVal,
		DiskGb:  diskGbVal,
		state:   attr.ValueStateKnown,
	}, diags
}

func NewAutoScalingValueNull() AutoScalingValue {
	return AutoScalingValue{
		state: attr.ValueStateNull,
	}
}

func NewAutoScalingValueUnknown() AutoScalingValue {
	return AutoScalingValue{
		state: attr.ValueStateUnknown,
	}
}

func NewAutoScalingValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (AutoScalingValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing AutoScalingValue Attribute Value",
				"While creating a AutoScalingValue value, a missing attribute value was detected. "+
					"A AutoScalingValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("AutoScalingValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid AutoScalingValue Attribute Type",
				"While creating a AutoScalingValue value, an invalid attribute value was detected. "+
					"A AutoScalingValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("AutoScalingValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("AutoScalingValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra AutoScalingValue Attribute Value",
				"While creating a AutoScalingValue value, an extra attribute value was detected. "+
					"A AutoScalingValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra AutoScalingValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewAutoScalingValueUnknown(), diags
	}

	computeAttribute, ok := attributes["compute"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`compute is missing from object`)

		return NewAutoScalingValueUnknown(), diags
	}

	computeVal, ok := computeAttribute.(basetypes.ObjectValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`compute expected to be basetypes.ObjectValue, was: %T`, computeAttribute))
	}

	diskGbAttribute, ok := attributes["disk_gb"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`disk_gb is missing from object`)

		return NewAutoScalingValueUnknown(), diags
	}

	diskGbVal, ok := diskGbAttribute.(basetypes.ObjectValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`disk_gb expected to be basetypes.ObjectValue, was: %T`, diskGbAttribute))
	}

	if diags.HasError() {
		return NewAutoScalingValueUnknown(), diags
	}

	return AutoScalingValue{
		Compute: computeVal,
		DiskGb:  diskGbVal,
		state:   attr.ValueStateKnown,
	}, diags
}

func NewAutoScalingValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) AutoScalingValue {
	object, diags := NewAutoScalingValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewAutoScalingValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t AutoScalingType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewAutoScalingValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewAutoScalingValueUnknown(), nil
	}

	if in.IsNull() {
		return NewAutoScalingValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewAutoScalingValueMust(AutoScalingValue{}.AttributeTypes(ctx), attributes), nil
}

func (t AutoScalingType) ValueType(ctx context.Context) attr.Value {
	return AutoScalingValue{}
}

var _ basetypes.ObjectValuable = AutoScalingValue{}

type AutoScalingValue struct {
	Compute basetypes.ObjectValue `tfsdk:"compute"`
	DiskGb  basetypes.ObjectValue `tfsdk:"disk_gb"`
	state   attr.ValueState
}

func (v AutoScalingValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 2)

	var val tftypes.Value
	var err error

	attrTypes["compute"] = basetypes.ObjectType{
		AttrTypes: ComputeValue{}.AttributeTypes(ctx),
	}.TerraformType(ctx)
	attrTypes["disk_gb"] = basetypes.ObjectType{
		AttrTypes: DiskGbValue{}.AttributeTypes(ctx),
	}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 2)

		val, err = v.Compute.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["compute"] = val

		val, err = v.DiskGb.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["disk_gb"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v AutoScalingValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v AutoScalingValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v AutoScalingValue) String() string {
	return "AutoScalingValue"
}

func (v AutoScalingValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	var compute basetypes.ObjectValue

	if v.Compute.IsNull() {
		compute = types.ObjectNull(
			ComputeValue{}.AttributeTypes(ctx),
		)
	}

	if v.Compute.IsUnknown() {
		compute = types.ObjectUnknown(
			ComputeValue{}.AttributeTypes(ctx),
		)
	}

	if !v.Compute.IsNull() && !v.Compute.IsUnknown() {
		compute = types.ObjectValueMust(
			ComputeValue{}.AttributeTypes(ctx),
			v.Compute.Attributes(),
		)
	}

	var diskGb basetypes.ObjectValue

	if v.DiskGb.IsNull() {
		diskGb = types.ObjectNull(
			DiskGbValue{}.AttributeTypes(ctx),
		)
	}

	if v.DiskGb.IsUnknown() {
		diskGb = types.ObjectUnknown(
			DiskGbValue{}.AttributeTypes(ctx),
		)
	}

	if !v.DiskGb.IsNull() && !v.DiskGb.IsUnknown() {
		diskGb = types.ObjectValueMust(
			DiskGbValue{}.AttributeTypes(ctx),
			v.DiskGb.Attributes(),
		)
	}

	attributeTypes := map[string]attr.Type{
		"compute": basetypes.ObjectType{
			AttrTypes: ComputeValue{}.AttributeTypes(ctx),
		},
		"disk_gb": basetypes.ObjectType{
			AttrTypes: DiskGbValue{}.AttributeTypes(ctx),
		},
	}

	if v.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	}

	if v.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	}

	objVal, diags := types.ObjectValue(
		attributeTypes,
		map[string]attr.Value{
			"compute": compute,
			"disk_gb": diskGb,
		})

	return objVal, diags
}

func (v AutoScalingValue) Equal(o attr.Value) bool {
	other, ok := o.(AutoScalingValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Compute.Equal(other.Compute) {
		return false
	}

	if !v.DiskGb.Equal(other.DiskGb) {
		return false
	}

	return true
}

func (v AutoScalingValue) Type(ctx context.Context) attr.Type {
	return AutoScalingType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v AutoScalingValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"compute": basetypes.ObjectType{
			AttrTypes: ComputeValue{}.AttributeTypes(ctx),
		},
		"disk_gb": basetypes.ObjectType{
			AttrTypes: DiskGbValue{}.AttributeTypes(ctx),
		},
	}
}

var _ basetypes.ObjectTypable = ElectableSpecsType{}

type ElectableSpecsType struct {
	basetypes.ObjectType
}

func (t ElectableSpecsType) Equal(o attr.Type) bool {
	other, ok := o.(ElectableSpecsType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t ElectableSpecsType) String() string {
	return "ElectableSpecsType"
}

func (t ElectableSpecsType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	diskIopsAttribute, ok := attributes["disk_iops"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`disk_iops is missing from object`)

		return nil, diags
	}

	diskIopsVal, ok := diskIopsAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`disk_iops expected to be basetypes.Int64Value, was: %T`, diskIopsAttribute))
	}

	diskSizeGbAttribute, ok := attributes["disk_size_gb"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`disk_size_gb is missing from object`)

		return nil, diags
	}

	diskSizeGbVal, ok := diskSizeGbAttribute.(basetypes.Float64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`disk_size_gb expected to be basetypes.Float64Value, was: %T`, diskSizeGbAttribute))
	}

	ebsVolumeTypeAttribute, ok := attributes["ebs_volume_type"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`ebs_volume_type is missing from object`)

		return nil, diags
	}

	ebsVolumeTypeVal, ok := ebsVolumeTypeAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`ebs_volume_type expected to be basetypes.StringValue, was: %T`, ebsVolumeTypeAttribute))
	}

	instanceSizeAttribute, ok := attributes["instance_size"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`instance_size is missing from object`)

		return nil, diags
	}

	instanceSizeVal, ok := instanceSizeAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`instance_size expected to be basetypes.StringValue, was: %T`, instanceSizeAttribute))
	}

	nodeCountAttribute, ok := attributes["node_count"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`node_count is missing from object`)

		return nil, diags
	}

	nodeCountVal, ok := nodeCountAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`node_count expected to be basetypes.Int64Value, was: %T`, nodeCountAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return ElectableSpecsValue{
		DiskIops:      diskIopsVal,
		DiskSizeGb:    diskSizeGbVal,
		EbsVolumeType: ebsVolumeTypeVal,
		InstanceSize:  instanceSizeVal,
		NodeCount:     nodeCountVal,
		state:         attr.ValueStateKnown,
	}, diags
}

func NewElectableSpecsValueNull() ElectableSpecsValue {
	return ElectableSpecsValue{
		state: attr.ValueStateNull,
	}
}

func NewElectableSpecsValueUnknown() ElectableSpecsValue {
	return ElectableSpecsValue{
		state: attr.ValueStateUnknown,
	}
}

func NewElectableSpecsValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (ElectableSpecsValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing ElectableSpecsValue Attribute Value",
				"While creating a ElectableSpecsValue value, a missing attribute value was detected. "+
					"A ElectableSpecsValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("ElectableSpecsValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid ElectableSpecsValue Attribute Type",
				"While creating a ElectableSpecsValue value, an invalid attribute value was detected. "+
					"A ElectableSpecsValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("ElectableSpecsValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("ElectableSpecsValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra ElectableSpecsValue Attribute Value",
				"While creating a ElectableSpecsValue value, an extra attribute value was detected. "+
					"A ElectableSpecsValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra ElectableSpecsValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewElectableSpecsValueUnknown(), diags
	}

	diskIopsAttribute, ok := attributes["disk_iops"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`disk_iops is missing from object`)

		return NewElectableSpecsValueUnknown(), diags
	}

	diskIopsVal, ok := diskIopsAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`disk_iops expected to be basetypes.Int64Value, was: %T`, diskIopsAttribute))
	}

	diskSizeGbAttribute, ok := attributes["disk_size_gb"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`disk_size_gb is missing from object`)

		return NewElectableSpecsValueUnknown(), diags
	}

	diskSizeGbVal, ok := diskSizeGbAttribute.(basetypes.Float64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`disk_size_gb expected to be basetypes.Float64Value, was: %T`, diskSizeGbAttribute))
	}

	ebsVolumeTypeAttribute, ok := attributes["ebs_volume_type"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`ebs_volume_type is missing from object`)

		return NewElectableSpecsValueUnknown(), diags
	}

	ebsVolumeTypeVal, ok := ebsVolumeTypeAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`ebs_volume_type expected to be basetypes.StringValue, was: %T`, ebsVolumeTypeAttribute))
	}

	instanceSizeAttribute, ok := attributes["instance_size"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`instance_size is missing from object`)

		return NewElectableSpecsValueUnknown(), diags
	}

	instanceSizeVal, ok := instanceSizeAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`instance_size expected to be basetypes.StringValue, was: %T`, instanceSizeAttribute))
	}

	nodeCountAttribute, ok := attributes["node_count"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`node_count is missing from object`)

		return NewElectableSpecsValueUnknown(), diags
	}

	nodeCountVal, ok := nodeCountAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`node_count expected to be basetypes.Int64Value, was: %T`, nodeCountAttribute))
	}

	if diags.HasError() {
		return NewElectableSpecsValueUnknown(), diags
	}

	return ElectableSpecsValue{
		DiskIops:      diskIopsVal,
		DiskSizeGb:    diskSizeGbVal,
		EbsVolumeType: ebsVolumeTypeVal,
		InstanceSize:  instanceSizeVal,
		NodeCount:     nodeCountVal,
		state:         attr.ValueStateKnown,
	}, diags
}

func NewElectableSpecsValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) ElectableSpecsValue {
	object, diags := NewElectableSpecsValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewElectableSpecsValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t ElectableSpecsType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewElectableSpecsValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewElectableSpecsValueUnknown(), nil
	}

	if in.IsNull() {
		return NewElectableSpecsValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewElectableSpecsValueMust(ElectableSpecsValue{}.AttributeTypes(ctx), attributes), nil
}

func (t ElectableSpecsType) ValueType(ctx context.Context) attr.Value {
	return ElectableSpecsValue{}
}

var _ basetypes.ObjectValuable = ElectableSpecsValue{}

type ElectableSpecsValue struct {
	DiskSizeGb    basetypes.Float64Value `tfsdk:"disk_size_gb"`
	EbsVolumeType basetypes.StringValue  `tfsdk:"ebs_volume_type"`
	InstanceSize  basetypes.StringValue  `tfsdk:"instance_size"`
	DiskIops      basetypes.Int64Value   `tfsdk:"disk_iops"`
	NodeCount     basetypes.Int64Value   `tfsdk:"node_count"`
	state         attr.ValueState
}

func (v ElectableSpecsValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 5)

	var val tftypes.Value
	var err error

	attrTypes["disk_iops"] = basetypes.Int64Type{}.TerraformType(ctx)
	attrTypes["disk_size_gb"] = basetypes.Float64Type{}.TerraformType(ctx)
	attrTypes["ebs_volume_type"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["instance_size"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["node_count"] = basetypes.Int64Type{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 5)

		val, err = v.DiskIops.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["disk_iops"] = val

		val, err = v.DiskSizeGb.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["disk_size_gb"] = val

		val, err = v.EbsVolumeType.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["ebs_volume_type"] = val

		val, err = v.InstanceSize.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["instance_size"] = val

		val, err = v.NodeCount.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["node_count"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v ElectableSpecsValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v ElectableSpecsValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v ElectableSpecsValue) String() string {
	return "ElectableSpecsValue"
}

func (v ElectableSpecsValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := map[string]attr.Type{
		"disk_iops":       basetypes.Int64Type{},
		"disk_size_gb":    basetypes.Float64Type{},
		"ebs_volume_type": basetypes.StringType{},
		"instance_size":   basetypes.StringType{},
		"node_count":      basetypes.Int64Type{},
	}

	if v.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	}

	if v.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	}

	objVal, diags := types.ObjectValue(
		attributeTypes,
		map[string]attr.Value{
			"disk_iops":       v.DiskIops,
			"disk_size_gb":    v.DiskSizeGb,
			"ebs_volume_type": v.EbsVolumeType,
			"instance_size":   v.InstanceSize,
			"node_count":      v.NodeCount,
		})

	return objVal, diags
}

func (v ElectableSpecsValue) Equal(o attr.Value) bool {
	other, ok := o.(ElectableSpecsValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.DiskIops.Equal(other.DiskIops) {
		return false
	}

	if !v.DiskSizeGb.Equal(other.DiskSizeGb) {
		return false
	}

	if !v.EbsVolumeType.Equal(other.EbsVolumeType) {
		return false
	}

	if !v.InstanceSize.Equal(other.InstanceSize) {
		return false
	}

	if !v.NodeCount.Equal(other.NodeCount) {
		return false
	}

	return true
}

func (v ElectableSpecsValue) Type(ctx context.Context) attr.Type {
	return ElectableSpecsType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v ElectableSpecsValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"disk_iops":       basetypes.Int64Type{},
		"disk_size_gb":    basetypes.Float64Type{},
		"ebs_volume_type": basetypes.StringType{},
		"instance_size":   basetypes.StringType{},
		"node_count":      basetypes.Int64Type{},
	}
}

var _ basetypes.ObjectTypable = ReadOnlySpecsType{}

type ReadOnlySpecsType struct {
	basetypes.ObjectType
}

func (t ReadOnlySpecsType) Equal(o attr.Type) bool {
	other, ok := o.(ReadOnlySpecsType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t ReadOnlySpecsType) String() string {
	return "ReadOnlySpecsType"
}

func (t ReadOnlySpecsType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	diskIopsAttribute, ok := attributes["disk_iops"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`disk_iops is missing from object`)

		return nil, diags
	}

	diskIopsVal, ok := diskIopsAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`disk_iops expected to be basetypes.Int64Value, was: %T`, diskIopsAttribute))
	}

	diskSizeGbAttribute, ok := attributes["disk_size_gb"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`disk_size_gb is missing from object`)

		return nil, diags
	}

	diskSizeGbVal, ok := diskSizeGbAttribute.(basetypes.Float64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`disk_size_gb expected to be basetypes.Float64Value, was: %T`, diskSizeGbAttribute))
	}

	ebsVolumeTypeAttribute, ok := attributes["ebs_volume_type"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`ebs_volume_type is missing from object`)

		return nil, diags
	}

	ebsVolumeTypeVal, ok := ebsVolumeTypeAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`ebs_volume_type expected to be basetypes.StringValue, was: %T`, ebsVolumeTypeAttribute))
	}

	instanceSizeAttribute, ok := attributes["instance_size"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`instance_size is missing from object`)

		return nil, diags
	}

	instanceSizeVal, ok := instanceSizeAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`instance_size expected to be basetypes.StringValue, was: %T`, instanceSizeAttribute))
	}

	nodeCountAttribute, ok := attributes["node_count"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`node_count is missing from object`)

		return nil, diags
	}

	nodeCountVal, ok := nodeCountAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`node_count expected to be basetypes.Int64Value, was: %T`, nodeCountAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return ReadOnlySpecsValue{
		DiskIops:      diskIopsVal,
		DiskSizeGb:    diskSizeGbVal,
		EbsVolumeType: ebsVolumeTypeVal,
		InstanceSize:  instanceSizeVal,
		NodeCount:     nodeCountVal,
		state:         attr.ValueStateKnown,
	}, diags
}

func NewReadOnlySpecsValueNull() ReadOnlySpecsValue {
	return ReadOnlySpecsValue{
		state: attr.ValueStateNull,
	}
}

func NewReadOnlySpecsValueUnknown() ReadOnlySpecsValue {
	return ReadOnlySpecsValue{
		state: attr.ValueStateUnknown,
	}
}

func NewReadOnlySpecsValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (ReadOnlySpecsValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing ReadOnlySpecsValue Attribute Value",
				"While creating a ReadOnlySpecsValue value, a missing attribute value was detected. "+
					"A ReadOnlySpecsValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("ReadOnlySpecsValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid ReadOnlySpecsValue Attribute Type",
				"While creating a ReadOnlySpecsValue value, an invalid attribute value was detected. "+
					"A ReadOnlySpecsValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("ReadOnlySpecsValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("ReadOnlySpecsValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra ReadOnlySpecsValue Attribute Value",
				"While creating a ReadOnlySpecsValue value, an extra attribute value was detected. "+
					"A ReadOnlySpecsValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra ReadOnlySpecsValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewReadOnlySpecsValueUnknown(), diags
	}

	diskIopsAttribute, ok := attributes["disk_iops"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`disk_iops is missing from object`)

		return NewReadOnlySpecsValueUnknown(), diags
	}

	diskIopsVal, ok := diskIopsAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`disk_iops expected to be basetypes.Int64Value, was: %T`, diskIopsAttribute))
	}

	diskSizeGbAttribute, ok := attributes["disk_size_gb"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`disk_size_gb is missing from object`)

		return NewReadOnlySpecsValueUnknown(), diags
	}

	diskSizeGbVal, ok := diskSizeGbAttribute.(basetypes.Float64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`disk_size_gb expected to be basetypes.Float64Value, was: %T`, diskSizeGbAttribute))
	}

	ebsVolumeTypeAttribute, ok := attributes["ebs_volume_type"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`ebs_volume_type is missing from object`)

		return NewReadOnlySpecsValueUnknown(), diags
	}

	ebsVolumeTypeVal, ok := ebsVolumeTypeAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`ebs_volume_type expected to be basetypes.StringValue, was: %T`, ebsVolumeTypeAttribute))
	}

	instanceSizeAttribute, ok := attributes["instance_size"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`instance_size is missing from object`)

		return NewReadOnlySpecsValueUnknown(), diags
	}

	instanceSizeVal, ok := instanceSizeAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`instance_size expected to be basetypes.StringValue, was: %T`, instanceSizeAttribute))
	}

	nodeCountAttribute, ok := attributes["node_count"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`node_count is missing from object`)

		return NewReadOnlySpecsValueUnknown(), diags
	}

	nodeCountVal, ok := nodeCountAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`node_count expected to be basetypes.Int64Value, was: %T`, nodeCountAttribute))
	}

	if diags.HasError() {
		return NewReadOnlySpecsValueUnknown(), diags
	}

	return ReadOnlySpecsValue{
		DiskIops:      diskIopsVal,
		DiskSizeGb:    diskSizeGbVal,
		EbsVolumeType: ebsVolumeTypeVal,
		InstanceSize:  instanceSizeVal,
		NodeCount:     nodeCountVal,
		state:         attr.ValueStateKnown,
	}, diags
}

func NewReadOnlySpecsValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) ReadOnlySpecsValue {
	object, diags := NewReadOnlySpecsValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewReadOnlySpecsValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t ReadOnlySpecsType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewReadOnlySpecsValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewReadOnlySpecsValueUnknown(), nil
	}

	if in.IsNull() {
		return NewReadOnlySpecsValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewReadOnlySpecsValueMust(ReadOnlySpecsValue{}.AttributeTypes(ctx), attributes), nil
}

func (t ReadOnlySpecsType) ValueType(ctx context.Context) attr.Value {
	return ReadOnlySpecsValue{}
}

var _ basetypes.ObjectValuable = ReadOnlySpecsValue{}

type ReadOnlySpecsValue struct {
	DiskSizeGb    basetypes.Float64Value `tfsdk:"disk_size_gb"`
	EbsVolumeType basetypes.StringValue  `tfsdk:"ebs_volume_type"`
	InstanceSize  basetypes.StringValue  `tfsdk:"instance_size"`
	DiskIops      basetypes.Int64Value   `tfsdk:"disk_iops"`
	NodeCount     basetypes.Int64Value   `tfsdk:"node_count"`
	state         attr.ValueState
}

func (v ReadOnlySpecsValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 5)

	var val tftypes.Value
	var err error

	attrTypes["disk_iops"] = basetypes.Int64Type{}.TerraformType(ctx)
	attrTypes["disk_size_gb"] = basetypes.Float64Type{}.TerraformType(ctx)
	attrTypes["ebs_volume_type"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["instance_size"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["node_count"] = basetypes.Int64Type{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 5)

		val, err = v.DiskIops.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["disk_iops"] = val

		val, err = v.DiskSizeGb.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["disk_size_gb"] = val

		val, err = v.EbsVolumeType.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["ebs_volume_type"] = val

		val, err = v.InstanceSize.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["instance_size"] = val

		val, err = v.NodeCount.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["node_count"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v ReadOnlySpecsValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v ReadOnlySpecsValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v ReadOnlySpecsValue) String() string {
	return "ReadOnlySpecsValue"
}

func (v ReadOnlySpecsValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := map[string]attr.Type{
		"disk_iops":       basetypes.Int64Type{},
		"disk_size_gb":    basetypes.Float64Type{},
		"ebs_volume_type": basetypes.StringType{},
		"instance_size":   basetypes.StringType{},
		"node_count":      basetypes.Int64Type{},
	}

	if v.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	}

	if v.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	}

	objVal, diags := types.ObjectValue(
		attributeTypes,
		map[string]attr.Value{
			"disk_iops":       v.DiskIops,
			"disk_size_gb":    v.DiskSizeGb,
			"ebs_volume_type": v.EbsVolumeType,
			"instance_size":   v.InstanceSize,
			"node_count":      v.NodeCount,
		})

	return objVal, diags
}

func (v ReadOnlySpecsValue) Equal(o attr.Value) bool {
	other, ok := o.(ReadOnlySpecsValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.DiskIops.Equal(other.DiskIops) {
		return false
	}

	if !v.DiskSizeGb.Equal(other.DiskSizeGb) {
		return false
	}

	if !v.EbsVolumeType.Equal(other.EbsVolumeType) {
		return false
	}

	if !v.InstanceSize.Equal(other.InstanceSize) {
		return false
	}

	if !v.NodeCount.Equal(other.NodeCount) {
		return false
	}

	return true
}

func (v ReadOnlySpecsValue) Type(ctx context.Context) attr.Type {
	return ReadOnlySpecsType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v ReadOnlySpecsValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"disk_iops":       basetypes.Int64Type{},
		"disk_size_gb":    basetypes.Float64Type{},
		"ebs_volume_type": basetypes.StringType{},
		"instance_size":   basetypes.StringType{},
		"node_count":      basetypes.Int64Type{},
	}
}

var _ basetypes.ObjectTypable = TagsType{}

type TagsType struct {
	basetypes.ObjectType
}

func (t TagsType) Equal(o attr.Type) bool {
	other, ok := o.(TagsType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t TagsType) String() string {
	return "TagsType"
}

func (t TagsType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	keyAttribute, ok := attributes["key"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`key is missing from object`)

		return nil, diags
	}

	keyVal, ok := keyAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`key expected to be basetypes.StringValue, was: %T`, keyAttribute))
	}

	valueAttribute, ok := attributes["value"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`value is missing from object`)

		return nil, diags
	}

	valueVal, ok := valueAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`value expected to be basetypes.StringValue, was: %T`, valueAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return TagsValue{
		Key:   keyVal,
		Value: valueVal,
		state: attr.ValueStateKnown,
	}, diags
}

func NewTagsValueNull() TagsValue {
	return TagsValue{
		state: attr.ValueStateNull,
	}
}

func NewTagsValueUnknown() TagsValue {
	return TagsValue{
		state: attr.ValueStateUnknown,
	}
}

func NewTagsValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (TagsValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing TagsValue Attribute Value",
				"While creating a TagsValue value, a missing attribute value was detected. "+
					"A TagsValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("TagsValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid TagsValue Attribute Type",
				"While creating a TagsValue value, an invalid attribute value was detected. "+
					"A TagsValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("TagsValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("TagsValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra TagsValue Attribute Value",
				"While creating a TagsValue value, an extra attribute value was detected. "+
					"A TagsValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra TagsValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewTagsValueUnknown(), diags
	}

	keyAttribute, ok := attributes["key"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`key is missing from object`)

		return NewTagsValueUnknown(), diags
	}

	keyVal, ok := keyAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`key expected to be basetypes.StringValue, was: %T`, keyAttribute))
	}

	valueAttribute, ok := attributes["value"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`value is missing from object`)

		return NewTagsValueUnknown(), diags
	}

	valueVal, ok := valueAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`value expected to be basetypes.StringValue, was: %T`, valueAttribute))
	}

	if diags.HasError() {
		return NewTagsValueUnknown(), diags
	}

	return TagsValue{
		Key:   keyVal,
		Value: valueVal,
		state: attr.ValueStateKnown,
	}, diags
}

func NewTagsValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) TagsValue {
	object, diags := NewTagsValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewTagsValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t TagsType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewTagsValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewTagsValueUnknown(), nil
	}

	if in.IsNull() {
		return NewTagsValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewTagsValueMust(TagsValue{}.AttributeTypes(ctx), attributes), nil
}

func (t TagsType) ValueType(ctx context.Context) attr.Value {
	return TagsValue{}
}

var _ basetypes.ObjectValuable = TagsValue{}

type TagsValue struct {
	Key   basetypes.StringValue `tfsdk:"key"`
	Value basetypes.StringValue `tfsdk:"value"`
	state attr.ValueState
}

func (v TagsValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 2)

	var val tftypes.Value
	var err error

	attrTypes["key"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["value"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 2)

		val, err = v.Key.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["key"] = val

		val, err = v.Value.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["value"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v TagsValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v TagsValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v TagsValue) String() string {
	return "TagsValue"
}

func (v TagsValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := map[string]attr.Type{
		"key":   basetypes.StringType{},
		"value": basetypes.StringType{},
	}

	if v.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	}

	if v.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	}

	objVal, diags := types.ObjectValue(
		attributeTypes,
		map[string]attr.Value{
			"key":   v.Key,
			"value": v.Value,
		})

	return objVal, diags
}

func (v TagsValue) Equal(o attr.Value) bool {
	other, ok := o.(TagsValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Key.Equal(other.Key) {
		return false
	}

	if !v.Value.Equal(other.Value) {
		return false
	}

	return true
}

func (v TagsValue) Type(ctx context.Context) attr.Type {
	return TagsType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v TagsValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"key":   basetypes.StringType{},
		"value": basetypes.StringType{},
	}
}
