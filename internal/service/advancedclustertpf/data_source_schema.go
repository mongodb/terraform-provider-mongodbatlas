package advancedclustertpf

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"accept_data_risks_and_force_replica_set_reconfig": schema.StringAttribute{
				Computed:            true,
				Description:         "If reconfiguration is necessary to regain a primary due to a regional outage, submit this field alongside your topology reconfiguration to request a new regional outage resistant topology. Forced reconfigurations during an outage of the majority of electable nodes carry a risk of data loss if replicated writes (even majority committed writes) have not been replicated to the new primary node. MongoDB Atlas docs contain more information. To proceed with an operation which carries that risk, set **acceptDataRisksAndForceReplicaSetReconfig** to the current date.",
				MarkdownDescription: "If reconfiguration is necessary to regain a primary due to a regional outage, submit this field alongside your topology reconfiguration to request a new regional outage resistant topology. Forced reconfigurations during an outage of the majority of electable nodes carry a risk of data loss if replicated writes (even majority committed writes) have not been replicated to the new primary node. MongoDB Atlas docs contain more information. To proceed with an operation which carries that risk, set **acceptDataRisksAndForceReplicaSetReconfig** to the current date.",
			},
			"backup_enabled": schema.BoolAttribute{
				Computed:            true,
				Description:         "Flag that indicates whether the cluster can perform backups. If set to `true`, the cluster can perform backups. You must set this value to `true` for NVMe clusters. Backup uses [Cloud Backups](https://docs.atlas.mongodb.com/backup/cloud-backup/overview/) for dedicated clusters and [Shared Cluster Backups](https://docs.atlas.mongodb.com/backup/shared-tier/overview/) for tenant clusters. If set to `false`, the cluster doesn't use backups.",
				MarkdownDescription: "Flag that indicates whether the cluster can perform backups. If set to `true`, the cluster can perform backups. You must set this value to `true` for NVMe clusters. Backup uses [Cloud Backups](https://docs.atlas.mongodb.com/backup/cloud-backup/overview/) for dedicated clusters and [Shared Cluster Backups](https://docs.atlas.mongodb.com/backup/shared-tier/overview/) for tenant clusters. If set to `false`, the cluster doesn't use backups.",
			},
			"bi_connector": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Computed:            true,
						Description:         "Flag that indicates whether MongoDB Connector for Business Intelligence is enabled on the specified cluster.",
						MarkdownDescription: "Flag that indicates whether MongoDB Connector for Business Intelligence is enabled on the specified cluster.",
					},
					"read_preference": schema.StringAttribute{
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
				Computed:            true,
				Description:         "Settings needed to configure the MongoDB Connector for Business Intelligence for this cluster.",
				MarkdownDescription: "Settings needed to configure the MongoDB Connector for Business Intelligence for this cluster.",
			},
			"cluster_name": schema.StringAttribute{
				Required:            true,
				Description:         "Human-readable label that identifies this cluster.",
				MarkdownDescription: "Human-readable label that identifies this cluster.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
					stringvalidator.RegexMatches(regexp.MustCompile("^([a-zA-Z0-9][a-zA-Z0-9-]*)?[a-zA-Z0-9]+$"), ""),
				},
			},
			"cluster_type": schema.StringAttribute{
				Computed:            true,
				Description:         "Configuration of nodes that comprise the cluster.",
				MarkdownDescription: "Configuration of nodes that comprise the cluster.",
			},
			"config_server_management_mode": schema.StringAttribute{
				Computed:            true,
				Description:         "Config Server Management Mode for creating or updating a sharded cluster.\n\nWhen configured as ATLAS_MANAGED, atlas may automatically switch the cluster's config server type for optimal performance and savings.\n\nWhen configured as FIXED_TO_DEDICATED, the cluster will always use a dedicated config server.",
				MarkdownDescription: "Config Server Management Mode for creating or updating a sharded cluster.\n\nWhen configured as ATLAS_MANAGED, atlas may automatically switch the cluster's config server type for optimal performance and savings.\n\nWhen configured as FIXED_TO_DEDICATED, the cluster will always use a dedicated config server.",
			},
			"config_server_type": schema.StringAttribute{
				Computed:            true,
				Description:         "Describes a sharded cluster's config server type.",
				MarkdownDescription: "Describes a sharded cluster's config server type.",
			},
			"connection_strings": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"aws_private_link": schema.MapAttribute{
						ElementType:         types.StringType,
						Computed:            true,
						Description:         "Private endpoint-aware connection strings that use AWS-hosted clusters with Amazon Web Services (AWS) PrivateLink. Each key identifies an Amazon Web Services (AWS) interface endpoint. Each value identifies the related `mongodb://` connection string that you use to connect to MongoDB Cloud through the interface endpoint that the key names.",
						MarkdownDescription: "Private endpoint-aware connection strings that use AWS-hosted clusters with Amazon Web Services (AWS) PrivateLink. Each key identifies an Amazon Web Services (AWS) interface endpoint. Each value identifies the related `mongodb://` connection string that you use to connect to MongoDB Cloud through the interface endpoint that the key names.",
					},
					"aws_private_link_srv": schema.MapAttribute{
						ElementType:         types.StringType,
						Computed:            true,
						Description:         "Private endpoint-aware connection strings that use AWS-hosted clusters with Amazon Web Services (AWS) PrivateLink. Each key identifies an Amazon Web Services (AWS) interface endpoint. Each value identifies the related `mongodb://` connection string that you use to connect to Atlas through the interface endpoint that the key names.",
						MarkdownDescription: "Private endpoint-aware connection strings that use AWS-hosted clusters with Amazon Web Services (AWS) PrivateLink. Each key identifies an Amazon Web Services (AWS) interface endpoint. Each value identifies the related `mongodb://` connection string that you use to connect to Atlas through the interface endpoint that the key names.",
					},
					"private": schema.StringAttribute{
						Computed:            true,
						Description:         "Network peering connection strings for each interface Virtual Private Cloud (VPC) endpoint that you configured to connect to this cluster. This connection string uses the `mongodb+srv://` protocol. The resource returns this parameter once someone creates a network peering connection to this cluster. This protocol tells the application to look up the host seed list in the Domain Name System (DNS). This list synchronizes with the nodes in a cluster. If the connection string uses this Uniform Resource Identifier (URI) format, you don't need to append the seed list or change the URI if the nodes change. Use this URI format if your driver supports it. If it doesn't, use connectionStrings.private. For Amazon Web Services (AWS) clusters, this resource returns this parameter only if you enable custom DNS.",
						MarkdownDescription: "Network peering connection strings for each interface Virtual Private Cloud (VPC) endpoint that you configured to connect to this cluster. This connection string uses the `mongodb+srv://` protocol. The resource returns this parameter once someone creates a network peering connection to this cluster. This protocol tells the application to look up the host seed list in the Domain Name System (DNS). This list synchronizes with the nodes in a cluster. If the connection string uses this Uniform Resource Identifier (URI) format, you don't need to append the seed list or change the URI if the nodes change. Use this URI format if your driver supports it. If it doesn't, use connectionStrings.private. For Amazon Web Services (AWS) clusters, this resource returns this parameter only if you enable custom DNS.",
					},
					"private_endpoint": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"connection_string": schema.StringAttribute{
									Computed:            true,
									Description:         "Private endpoint-aware connection string that uses the `mongodb://` protocol to connect to MongoDB Cloud through a private endpoint.",
									MarkdownDescription: "Private endpoint-aware connection string that uses the `mongodb://` protocol to connect to MongoDB Cloud through a private endpoint.",
								},
								"endpoints": schema.ListNestedAttribute{
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"endpoint_id": schema.StringAttribute{
												Computed:            true,
												Description:         "Unique string that the cloud provider uses to identify the private endpoint.",
												MarkdownDescription: "Unique string that the cloud provider uses to identify the private endpoint.",
											},
											"provider_name": schema.StringAttribute{
												Computed:            true,
												Description:         "Cloud provider in which MongoDB Cloud deploys the private endpoint.",
												MarkdownDescription: "Cloud provider in which MongoDB Cloud deploys the private endpoint.",
											},
											"region": schema.StringAttribute{
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
									Computed:            true,
									Description:         "List that contains the private endpoints through which you connect to MongoDB Cloud when you use **connectionStrings.privateEndpoint[n].connectionString** or **connectionStrings.privateEndpoint[n].srvConnectionString**.",
									MarkdownDescription: "List that contains the private endpoints through which you connect to MongoDB Cloud when you use **connectionStrings.privateEndpoint[n].connectionString** or **connectionStrings.privateEndpoint[n].srvConnectionString**.",
								},
								"srv_connection_string": schema.StringAttribute{
									Computed:            true,
									Description:         "Private endpoint-aware connection string that uses the `mongodb+srv://` protocol to connect to MongoDB Cloud through a private endpoint. The `mongodb+srv` protocol tells the driver to look up the seed list of hosts in the Domain Name System (DNS). This list synchronizes with the nodes in a cluster. If the connection string uses this Uniform Resource Identifier (URI) format, you don't need to append the seed list or change the Uniform Resource Identifier (URI) if the nodes change. Use this Uniform Resource Identifier (URI) format if your application supports it. If it doesn't, use connectionStrings.privateEndpoint[n].connectionString.",
									MarkdownDescription: "Private endpoint-aware connection string that uses the `mongodb+srv://` protocol to connect to MongoDB Cloud through a private endpoint. The `mongodb+srv` protocol tells the driver to look up the seed list of hosts in the Domain Name System (DNS). This list synchronizes with the nodes in a cluster. If the connection string uses this Uniform Resource Identifier (URI) format, you don't need to append the seed list or change the Uniform Resource Identifier (URI) if the nodes change. Use this Uniform Resource Identifier (URI) format if your application supports it. If it doesn't, use connectionStrings.privateEndpoint[n].connectionString.",
								},
								"srv_shard_optimized_connection_string": schema.StringAttribute{
									Computed:            true,
									Description:         "Private endpoint-aware connection string optimized for sharded clusters that uses the `mongodb+srv://` protocol to connect to MongoDB Cloud through a private endpoint. If the connection string uses this Uniform Resource Identifier (URI) format, you don't need to change the Uniform Resource Identifier (URI) if the nodes change. Use this Uniform Resource Identifier (URI) format if your application and Atlas cluster supports it. If it doesn't, use and consult the documentation for connectionStrings.privateEndpoint[n].srvConnectionString.",
									MarkdownDescription: "Private endpoint-aware connection string optimized for sharded clusters that uses the `mongodb+srv://` protocol to connect to MongoDB Cloud through a private endpoint. If the connection string uses this Uniform Resource Identifier (URI) format, you don't need to change the Uniform Resource Identifier (URI) if the nodes change. Use this Uniform Resource Identifier (URI) format if your application and Atlas cluster supports it. If it doesn't, use and consult the documentation for connectionStrings.privateEndpoint[n].srvConnectionString.",
								},
								"type": schema.StringAttribute{
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
						Computed:            true,
						Description:         "List of private endpoint-aware connection strings that you can use to connect to this cluster through a private endpoint. This parameter returns only if you deployed a private endpoint to all regions to which you deployed this clusters' nodes.",
						MarkdownDescription: "List of private endpoint-aware connection strings that you can use to connect to this cluster through a private endpoint. This parameter returns only if you deployed a private endpoint to all regions to which you deployed this clusters' nodes.",
					},
					"private_srv": schema.StringAttribute{
						Computed:            true,
						Description:         "Network peering connection strings for each interface Virtual Private Cloud (VPC) endpoint that you configured to connect to this cluster. This connection string uses the `mongodb+srv://` protocol. The resource returns this parameter when someone creates a network peering connection to this cluster. This protocol tells the application to look up the host seed list in the Domain Name System (DNS). This list synchronizes with the nodes in a cluster. If the connection string uses this Uniform Resource Identifier (URI) format, you don't need to append the seed list or change the Uniform Resource Identifier (URI) if the nodes change. Use this Uniform Resource Identifier (URI) format if your driver supports it. If it doesn't, use `connectionStrings.private`. For Amazon Web Services (AWS) clusters, this parameter returns only if you [enable custom DNS](https://docs.atlas.mongodb.com/reference/api/aws-custom-dns-update/).",
						MarkdownDescription: "Network peering connection strings for each interface Virtual Private Cloud (VPC) endpoint that you configured to connect to this cluster. This connection string uses the `mongodb+srv://` protocol. The resource returns this parameter when someone creates a network peering connection to this cluster. This protocol tells the application to look up the host seed list in the Domain Name System (DNS). This list synchronizes with the nodes in a cluster. If the connection string uses this Uniform Resource Identifier (URI) format, you don't need to append the seed list or change the Uniform Resource Identifier (URI) if the nodes change. Use this Uniform Resource Identifier (URI) format if your driver supports it. If it doesn't, use `connectionStrings.private`. For Amazon Web Services (AWS) clusters, this parameter returns only if you [enable custom DNS](https://docs.atlas.mongodb.com/reference/api/aws-custom-dns-update/).",
					},
					"standard": schema.StringAttribute{
						Computed:            true,
						Description:         "Public connection string that you can use to connect to this cluster. This connection string uses the `mongodb://` protocol.",
						MarkdownDescription: "Public connection string that you can use to connect to this cluster. This connection string uses the `mongodb://` protocol.",
					},
					"standard_srv": schema.StringAttribute{
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
				Computed:            true,
				Description:         "Collection of Uniform Resource Locators that point to the MongoDB database.",
				MarkdownDescription: "Collection of Uniform Resource Locators that point to the MongoDB database.",
			},
			"create_date": schema.StringAttribute{
				Computed:            true,
				Description:         "Date and time when MongoDB Cloud created this cluster. This parameter expresses its value in ISO 8601 format in UTC.",
				MarkdownDescription: "Date and time when MongoDB Cloud created this cluster. This parameter expresses its value in ISO 8601 format in UTC.",
			},
			"disk_warming_mode": schema.StringAttribute{
				Computed:            true,
				Description:         "Disk warming mode selection.",
				MarkdownDescription: "Disk warming mode selection.",
			},
			"encryption_at_rest_provider": schema.StringAttribute{
				Computed:            true,
				Description:         "Cloud service provider that manages your customer keys to provide an additional layer of encryption at rest for the cluster. To enable customer key management for encryption at rest, the cluster **replicationSpecs[n].regionConfigs[m].{type}Specs.instanceSize** setting must be `M10` or higher and `\"backupEnabled\" : false` or omitted entirely.",
				MarkdownDescription: "Cloud service provider that manages your customer keys to provide an additional layer of encryption at rest for the cluster. To enable customer key management for encryption at rest, the cluster **replicationSpecs[n].regionConfigs[m].{type}Specs.instanceSize** setting must be `M10` or higher and `\"backupEnabled\" : false` or omitted entirely.",
			},
			"feature_compatibility_version": schema.StringAttribute{
				Computed:            true,
				Description:         "Feature compatibility version of the cluster.",
				MarkdownDescription: "Feature compatibility version of the cluster.",
			},
			"feature_compatibility_version_expiration_date": schema.StringAttribute{
				Computed:            true,
				Description:         "Feature compatibility version expiration date.",
				MarkdownDescription: "Feature compatibility version expiration date.",
			},
			"global_cluster_self_managed_sharding": schema.BoolAttribute{
				Computed:            true,
				Description:         "Set this field to configure the Sharding Management Mode when creating a new Global Cluster.\n\nWhen set to false, the management mode is set to Atlas-Managed Sharding. This mode fully manages the sharding of your Global Cluster and is built to provide a seamless deployment experience.\n\nWhen set to true, the management mode is set to Self-Managed Sharding. This mode leaves the management of shards in your hands and is built to provide an advanced and flexible deployment experience.\n\nThis setting cannot be changed once the cluster is deployed.",
				MarkdownDescription: "Set this field to configure the Sharding Management Mode when creating a new Global Cluster.\n\nWhen set to false, the management mode is set to Atlas-Managed Sharding. This mode fully manages the sharding of your Global Cluster and is built to provide a seamless deployment experience.\n\nWhen set to true, the management mode is set to Self-Managed Sharding. This mode leaves the management of shards in your hands and is built to provide an advanced and flexible deployment experience.\n\nThis setting cannot be changed once the cluster is deployed.",
			},
			"group_id": schema.StringAttribute{
				Required:            true,
				Description:         "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.",
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(24, 24),
					stringvalidator.RegexMatches(regexp.MustCompile("^([a-f0-9]{24})$"), ""),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "Unique 24-hexadecimal digit string that identifies the cluster.",
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the cluster.",
			},
			"labels": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Computed:            true,
							Description:         "Key applied to tag and categorize this component.",
							MarkdownDescription: "Key applied to tag and categorize this component.",
						},
						"value": schema.StringAttribute{
							Computed:            true,
							Description:         "Value set to the Key applied to tag and categorize this component.",
							MarkdownDescription: "Value set to the Key applied to tag and categorize this component.",
						},
					},
					CustomType: LabelsType{
						ObjectType: types.ObjectType{
							AttrTypes: LabelsValue{}.AttributeTypes(ctx),
						},
					},
				},
				Computed:            true,
				Description:         "Collection of key-value pairs between 1 to 255 characters in length that tag and categorize the cluster. The MongoDB Cloud console doesn't display your labels.\n\nCluster labels are deprecated and will be removed in a future release. We strongly recommend that you use [resource tags](https://dochub.mongodb.org/core/add-cluster-tag-atlas) instead.",
				MarkdownDescription: "Collection of key-value pairs between 1 to 255 characters in length that tag and categorize the cluster. The MongoDB Cloud console doesn't display your labels.\n\nCluster labels are deprecated and will be removed in a future release. We strongly recommend that you use [resource tags](https://dochub.mongodb.org/core/add-cluster-tag-atlas) instead.",
				DeprecationMessage:  "This attribute is deprecated.",
			},
			"links": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"href": schema.StringAttribute{
							Computed:            true,
							Description:         "Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
							MarkdownDescription: "Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
						},
						"rel": schema.StringAttribute{
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
				Computed:            true,
				Description:         "List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.",
				MarkdownDescription: "List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.",
			},
			"mongo_dbemployee_access_grant": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"expiration_time": schema.StringAttribute{
						Computed:            true,
						Description:         "Expiration date for the employee access grant.",
						MarkdownDescription: "Expiration date for the employee access grant.",
					},
					"grant_type": schema.StringAttribute{
						Computed:            true,
						Description:         "Level of access to grant to MongoDB Employees.",
						MarkdownDescription: "Level of access to grant to MongoDB Employees.",
					},
					"links": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"href": schema.StringAttribute{
									Computed:            true,
									Description:         "Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
									MarkdownDescription: "Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
								},
								"rel": schema.StringAttribute{
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
				Computed:            true,
				Description:         "MongoDB employee granted access level and expiration for a cluster.",
				MarkdownDescription: "MongoDB employee granted access level and expiration for a cluster.",
			},
			"mongo_dbmajor_version": schema.StringAttribute{
				Computed:            true,
				Description:         "MongoDB major version of the cluster.\n\nOn creation: Choose from the available versions of MongoDB, or leave unspecified for the current recommended default in the MongoDB Cloud platform. The recommended version is a recent Long Term Support version. The default is not guaranteed to be the most recently released version throughout the entire release cycle. For versions available in a specific project, see the linked documentation or use the API endpoint for [project LTS versions endpoint](#tag/Projects/operation/getProjectLTSVersions).\n\n On update: Increase version only by 1 major version at a time. If the cluster is pinned to a MongoDB feature compatibility version exactly one major version below the current MongoDB version, the MongoDB version can be downgraded to the previous major version.",
				MarkdownDescription: "MongoDB major version of the cluster.\n\nOn creation: Choose from the available versions of MongoDB, or leave unspecified for the current recommended default in the MongoDB Cloud platform. The recommended version is a recent Long Term Support version. The default is not guaranteed to be the most recently released version throughout the entire release cycle. For versions available in a specific project, see the linked documentation or use the API endpoint for [project LTS versions endpoint](#tag/Projects/operation/getProjectLTSVersions).\n\n On update: Increase version only by 1 major version at a time. If the cluster is pinned to a MongoDB feature compatibility version exactly one major version below the current MongoDB version, the MongoDB version can be downgraded to the previous major version.",
			},
			"mongo_dbversion": schema.StringAttribute{
				Computed:            true,
				Description:         "Version of MongoDB that the cluster runs.",
				MarkdownDescription: "Version of MongoDB that the cluster runs.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				Description:         "Human-readable label that identifies the cluster.",
				MarkdownDescription: "Human-readable label that identifies the cluster.",
			},
			"paused": schema.BoolAttribute{
				Computed:            true,
				Description:         "Flag that indicates whether the cluster is paused.",
				MarkdownDescription: "Flag that indicates whether the cluster is paused.",
			},
			"pit_enabled": schema.BoolAttribute{
				Computed:            true,
				Description:         "Flag that indicates whether the cluster uses continuous cloud backups.",
				MarkdownDescription: "Flag that indicates whether the cluster uses continuous cloud backups.",
			},
			"redact_client_log_data": schema.BoolAttribute{
				Computed:            true,
				Description:         "Enable or disable log redaction.\n\nThis setting configures the ``mongod`` or ``mongos`` to redact any document field contents from a message accompanying a given log event before logging. This prevents the program from writing potentially sensitive data stored on the database to the diagnostic log. Metadata such as error or operation codes, line numbers, and source file names are still visible in the logs.\n\nUse ``redactClientLogData`` in conjunction with Encryption at Rest and TLS/SSL (Transport Encryption) to assist compliance with regulatory requirements.\n\n*Note*: changing this setting on a cluster will trigger a rolling restart as soon as the cluster is updated.",
				MarkdownDescription: "Enable or disable log redaction.\n\nThis setting configures the ``mongod`` or ``mongos`` to redact any document field contents from a message accompanying a given log event before logging. This prevents the program from writing potentially sensitive data stored on the database to the diagnostic log. Metadata such as error or operation codes, line numbers, and source file names are still visible in the logs.\n\nUse ``redactClientLogData`` in conjunction with Encryption at Rest and TLS/SSL (Transport Encryption) to assist compliance with regulatory requirements.\n\n*Note*: changing this setting on a cluster will trigger a rolling restart as soon as the cluster is updated.",
			},
			"replica_set_scaling_strategy": schema.StringAttribute{
				Computed:            true,
				Description:         "Set this field to configure the replica set scaling mode for your cluster.\n\nBy default, Atlas scales under WORKLOAD_TYPE. This mode allows Atlas to scale your analytics nodes in parallel to your operational nodes.\n\nWhen configured as SEQUENTIAL, Atlas scales all nodes sequentially. This mode is intended for steady-state workloads and applications performing latency-sensitive secondary reads.\n\nWhen configured as NODE_TYPE, Atlas scales your electable nodes in parallel with your read-only and analytics nodes. This mode is intended for large, dynamic workloads requiring frequent and timely cluster tier scaling. This is the fastest scaling strategy, but it might impact latency of workloads when performing extensive secondary reads.",
				MarkdownDescription: "Set this field to configure the replica set scaling mode for your cluster.\n\nBy default, Atlas scales under WORKLOAD_TYPE. This mode allows Atlas to scale your analytics nodes in parallel to your operational nodes.\n\nWhen configured as SEQUENTIAL, Atlas scales all nodes sequentially. This mode is intended for steady-state workloads and applications performing latency-sensitive secondary reads.\n\nWhen configured as NODE_TYPE, Atlas scales your electable nodes in parallel with your read-only and analytics nodes. This mode is intended for large, dynamic workloads requiring frequent and timely cluster tier scaling. This is the fastest scaling strategy, but it might impact latency of workloads when performing extensive secondary reads.",
			},
			"replication_specs": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							Description:         "Unique 24-hexadecimal digit string that identifies the replication object for a shard in a Cluster. If you include existing shard replication configurations in the request, you must specify this parameter. If you add a new shard to an existing Cluster, you may specify this parameter. The request deletes any existing shards  in the Cluster that you exclude from the request. This corresponds to Shard ID displayed in the UI.",
							MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the replication object for a shard in a Cluster. If you include existing shard replication configurations in the request, you must specify this parameter. If you add a new shard to an existing Cluster, you may specify this parameter. The request deletes any existing shards  in the Cluster that you exclude from the request. This corresponds to Shard ID displayed in the UI.",
						},
						"region_configs": schema.ListNestedAttribute{
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"analytics_auto_scaling": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"compute": schema.SingleNestedAttribute{
												Attributes: map[string]schema.Attribute{
													"enabled": schema.BoolAttribute{
														Computed:            true,
														Description:         "Flag that indicates whether someone enabled instance size auto-scaling.\n\n- Set to `true` to enable instance size auto-scaling. If enabled, you must specify a value for **replicationSpecs[n].regionConfigs[m].autoScaling.compute.maxInstanceSize**.\n- Set to `false` to disable instance size automatic scaling.",
														MarkdownDescription: "Flag that indicates whether someone enabled instance size auto-scaling.\n\n- Set to `true` to enable instance size auto-scaling. If enabled, you must specify a value for **replicationSpecs[n].regionConfigs[m].autoScaling.compute.maxInstanceSize**.\n- Set to `false` to disable instance size automatic scaling.",
													},
													"max_instance_size": schema.StringAttribute{
														Computed:            true,
														Description:         "Minimum instance size to which your cluster can automatically scale. MongoDB Cloud requires this parameter if `\"replicationSpecs[n].regionConfigs[m].autoScaling.compute.scaleDownEnabled\" : true`.",
														MarkdownDescription: "Minimum instance size to which your cluster can automatically scale. MongoDB Cloud requires this parameter if `\"replicationSpecs[n].regionConfigs[m].autoScaling.compute.scaleDownEnabled\" : true`.",
													},
													"min_instance_size": schema.StringAttribute{
														Computed:            true,
														Description:         "Minimum instance size to which your cluster can automatically scale. MongoDB Cloud requires this parameter if `\"replicationSpecs[n].regionConfigs[m].autoScaling.compute.scaleDownEnabled\" : true`.",
														MarkdownDescription: "Minimum instance size to which your cluster can automatically scale. MongoDB Cloud requires this parameter if `\"replicationSpecs[n].regionConfigs[m].autoScaling.compute.scaleDownEnabled\" : true`.",
													},
													"scale_down_enabled": schema.BoolAttribute{
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
												Computed:            true,
												Description:         "Options that determine how this cluster handles CPU scaling.",
												MarkdownDescription: "Options that determine how this cluster handles CPU scaling.",
											},
											"disk_gb": schema.SingleNestedAttribute{
												Attributes: map[string]schema.Attribute{
													"enabled": schema.BoolAttribute{
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
										Computed:            true,
										Description:         "Options that determine how this cluster handles resource scaling.",
										MarkdownDescription: "Options that determine how this cluster handles resource scaling.",
									},
									"analytics_specs": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"disk_iops": schema.Int64Attribute{
												Computed:            true,
												Description:         "Target throughput desired for storage attached to your Azure-provisioned cluster. Change this parameter if you:\n\n- set `\"replicationSpecs[n].regionConfigs[m].providerName\" : \"Azure\"`.\n- set `\"replicationSpecs[n].regionConfigs[m].electableSpecs.instanceSize\" : \"M40\"` or greater not including `Mxx_NVME` tiers.\n\nThe maximum input/output operations per second (IOPS) depend on the selected **.instanceSize** and **.diskSizeGB**.\nThis parameter defaults to the cluster tier's standard IOPS value.\nChanging this value impacts cluster cost.",
												MarkdownDescription: "Target throughput desired for storage attached to your Azure-provisioned cluster. Change this parameter if you:\n\n- set `\"replicationSpecs[n].regionConfigs[m].providerName\" : \"Azure\"`.\n- set `\"replicationSpecs[n].regionConfigs[m].electableSpecs.instanceSize\" : \"M40\"` or greater not including `Mxx_NVME` tiers.\n\nThe maximum input/output operations per second (IOPS) depend on the selected **.instanceSize** and **.diskSizeGB**.\nThis parameter defaults to the cluster tier's standard IOPS value.\nChanging this value impacts cluster cost.",
											},
											"disk_size_gb": schema.Float64Attribute{
												Computed:            true,
												Description:         "Storage capacity of instance data volumes expressed in gigabytes. Increase this number to add capacity.\n\n This value must be equal for all shards and node types.\n\n This value is not configurable on M0/M2/M5 clusters.\n\n MongoDB Cloud requires this parameter if you set **replicationSpecs**.\n\n If you specify a disk size below the minimum (10 GB), this parameter defaults to the minimum disk size value. \n\n Storage charge calculations depend on whether you choose the default value or a custom value.\n\n The maximum value for disk storage cannot exceed 50 times the maximum RAM for the selected cluster. If you require more storage space, consider upgrading your cluster to a higher tier.",
												MarkdownDescription: "Storage capacity of instance data volumes expressed in gigabytes. Increase this number to add capacity.\n\n This value must be equal for all shards and node types.\n\n This value is not configurable on M0/M2/M5 clusters.\n\n MongoDB Cloud requires this parameter if you set **replicationSpecs**.\n\n If you specify a disk size below the minimum (10 GB), this parameter defaults to the minimum disk size value. \n\n Storage charge calculations depend on whether you choose the default value or a custom value.\n\n The maximum value for disk storage cannot exceed 50 times the maximum RAM for the selected cluster. If you require more storage space, consider upgrading your cluster to a higher tier.",
											},
											"ebs_volume_type": schema.StringAttribute{
												Computed:            true,
												Description:         "Type of storage you want to attach to your AWS-provisioned cluster.\n\n- `STANDARD` volume types can't exceed the default input/output operations per second (IOPS) rate for the selected volume size. \n\n- `PROVISIONED` volume types must fall within the allowable IOPS range for the selected volume size. You must set this value to (`PROVISIONED`) for NVMe clusters.",
												MarkdownDescription: "Type of storage you want to attach to your AWS-provisioned cluster.\n\n- `STANDARD` volume types can't exceed the default input/output operations per second (IOPS) rate for the selected volume size. \n\n- `PROVISIONED` volume types must fall within the allowable IOPS range for the selected volume size. You must set this value to (`PROVISIONED`) for NVMe clusters.",
											},
											"instance_size": schema.StringAttribute{
												Computed:            true,
												Description:         "Hardware specification for the instance sizes in this region in this shard. Each instance size has a default storage and memory capacity. Electable nodes and read-only nodes (known as \"base nodes\") within a single shard must use the same instance size. Analytics nodes can scale independently from base nodes within a shard. Both base nodes and analytics nodes can scale independently from their equivalents in other shards.",
												MarkdownDescription: "Hardware specification for the instance sizes in this region in this shard. Each instance size has a default storage and memory capacity. Electable nodes and read-only nodes (known as \"base nodes\") within a single shard must use the same instance size. Analytics nodes can scale independently from base nodes within a shard. Both base nodes and analytics nodes can scale independently from their equivalents in other shards.",
											},
											"node_count": schema.Int64Attribute{
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
										Computed:            true,
										Description:         "Hardware specifications for read-only nodes in the region. Read-only nodes can never become the primary member, but can enable local reads. If you don't specify this parameter, no read-only nodes are deployed to the region.",
										MarkdownDescription: "Hardware specifications for read-only nodes in the region. Read-only nodes can never become the primary member, but can enable local reads. If you don't specify this parameter, no read-only nodes are deployed to the region.",
									},
									"auto_scaling": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"compute": schema.SingleNestedAttribute{
												Attributes: map[string]schema.Attribute{
													"enabled": schema.BoolAttribute{
														Computed:            true,
														Description:         "Flag that indicates whether someone enabled instance size auto-scaling.\n\n- Set to `true` to enable instance size auto-scaling. If enabled, you must specify a value for **replicationSpecs[n].regionConfigs[m].autoScaling.compute.maxInstanceSize**.\n- Set to `false` to disable instance size automatic scaling.",
														MarkdownDescription: "Flag that indicates whether someone enabled instance size auto-scaling.\n\n- Set to `true` to enable instance size auto-scaling. If enabled, you must specify a value for **replicationSpecs[n].regionConfigs[m].autoScaling.compute.maxInstanceSize**.\n- Set to `false` to disable instance size automatic scaling.",
													},
													"max_instance_size": schema.StringAttribute{
														Computed:            true,
														Description:         "Minimum instance size to which your cluster can automatically scale. MongoDB Cloud requires this parameter if `\"replicationSpecs[n].regionConfigs[m].autoScaling.compute.scaleDownEnabled\" : true`.",
														MarkdownDescription: "Minimum instance size to which your cluster can automatically scale. MongoDB Cloud requires this parameter if `\"replicationSpecs[n].regionConfigs[m].autoScaling.compute.scaleDownEnabled\" : true`.",
													},
													"min_instance_size": schema.StringAttribute{
														Computed:            true,
														Description:         "Minimum instance size to which your cluster can automatically scale. MongoDB Cloud requires this parameter if `\"replicationSpecs[n].regionConfigs[m].autoScaling.compute.scaleDownEnabled\" : true`.",
														MarkdownDescription: "Minimum instance size to which your cluster can automatically scale. MongoDB Cloud requires this parameter if `\"replicationSpecs[n].regionConfigs[m].autoScaling.compute.scaleDownEnabled\" : true`.",
													},
													"scale_down_enabled": schema.BoolAttribute{
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
												Computed:            true,
												Description:         "Options that determine how this cluster handles CPU scaling.",
												MarkdownDescription: "Options that determine how this cluster handles CPU scaling.",
											},
											"disk_gb": schema.SingleNestedAttribute{
												Attributes: map[string]schema.Attribute{
													"enabled": schema.BoolAttribute{
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
										Computed:            true,
										Description:         "Options that determine how this cluster handles resource scaling.",
										MarkdownDescription: "Options that determine how this cluster handles resource scaling.",
									},
									"backing_provider_name": schema.StringAttribute{
										Computed:            true,
										Description:         "Cloud service provider on which MongoDB Cloud provisioned the multi-tenant cluster. The resource returns this parameter when **providerName** is `TENANT` and **electableSpecs.instanceSize** is `M0`, `M2` or `M5`.",
										MarkdownDescription: "Cloud service provider on which MongoDB Cloud provisioned the multi-tenant cluster. The resource returns this parameter when **providerName** is `TENANT` and **electableSpecs.instanceSize** is `M0`, `M2` or `M5`.",
									},
									"electable_specs": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"disk_iops": schema.Int64Attribute{
												Computed:            true,
												Description:         "Target throughput desired for storage attached to your Azure-provisioned cluster. Change this parameter if you:\n\n- set `\"replicationSpecs[n].regionConfigs[m].providerName\" : \"Azure\"`.\n- set `\"replicationSpecs[n].regionConfigs[m].electableSpecs.instanceSize\" : \"M40\"` or greater not including `Mxx_NVME` tiers.\n\nThe maximum input/output operations per second (IOPS) depend on the selected **.instanceSize** and **.diskSizeGB**.\nThis parameter defaults to the cluster tier's standard IOPS value.\nChanging this value impacts cluster cost.",
												MarkdownDescription: "Target throughput desired for storage attached to your Azure-provisioned cluster. Change this parameter if you:\n\n- set `\"replicationSpecs[n].regionConfigs[m].providerName\" : \"Azure\"`.\n- set `\"replicationSpecs[n].regionConfigs[m].electableSpecs.instanceSize\" : \"M40\"` or greater not including `Mxx_NVME` tiers.\n\nThe maximum input/output operations per second (IOPS) depend on the selected **.instanceSize** and **.diskSizeGB**.\nThis parameter defaults to the cluster tier's standard IOPS value.\nChanging this value impacts cluster cost.",
											},
											"disk_size_gb": schema.Float64Attribute{
												Computed:            true,
												Description:         "Storage capacity of instance data volumes expressed in gigabytes. Increase this number to add capacity.\n\n This value must be equal for all shards and node types.\n\n This value is not configurable on M0/M2/M5 clusters.\n\n MongoDB Cloud requires this parameter if you set **replicationSpecs**.\n\n If you specify a disk size below the minimum (10 GB), this parameter defaults to the minimum disk size value. \n\n Storage charge calculations depend on whether you choose the default value or a custom value.\n\n The maximum value for disk storage cannot exceed 50 times the maximum RAM for the selected cluster. If you require more storage space, consider upgrading your cluster to a higher tier.",
												MarkdownDescription: "Storage capacity of instance data volumes expressed in gigabytes. Increase this number to add capacity.\n\n This value must be equal for all shards and node types.\n\n This value is not configurable on M0/M2/M5 clusters.\n\n MongoDB Cloud requires this parameter if you set **replicationSpecs**.\n\n If you specify a disk size below the minimum (10 GB), this parameter defaults to the minimum disk size value. \n\n Storage charge calculations depend on whether you choose the default value or a custom value.\n\n The maximum value for disk storage cannot exceed 50 times the maximum RAM for the selected cluster. If you require more storage space, consider upgrading your cluster to a higher tier.",
											},
											"ebs_volume_type": schema.StringAttribute{
												Computed:            true,
												Description:         "Type of storage you want to attach to your AWS-provisioned cluster.\n\n- `STANDARD` volume types can't exceed the default input/output operations per second (IOPS) rate for the selected volume size. \n\n- `PROVISIONED` volume types must fall within the allowable IOPS range for the selected volume size. You must set this value to (`PROVISIONED`) for NVMe clusters.",
												MarkdownDescription: "Type of storage you want to attach to your AWS-provisioned cluster.\n\n- `STANDARD` volume types can't exceed the default input/output operations per second (IOPS) rate for the selected volume size. \n\n- `PROVISIONED` volume types must fall within the allowable IOPS range for the selected volume size. You must set this value to (`PROVISIONED`) for NVMe clusters.",
											},
											"instance_size": schema.StringAttribute{
												Computed:            true,
												Description:         "Hardware specification for the instances in this M0/M2/M5 tier cluster.",
												MarkdownDescription: "Hardware specification for the instances in this M0/M2/M5 tier cluster.",
											},
											"node_count": schema.Int64Attribute{
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
										Computed:            true,
										Description:         "Hardware specifications for all electable nodes deployed in the region. Electable nodes can become the primary and can enable local reads. If you don't specify this option, MongoDB Cloud deploys no electable nodes to the region.",
										MarkdownDescription: "Hardware specifications for all electable nodes deployed in the region. Electable nodes can become the primary and can enable local reads. If you don't specify this option, MongoDB Cloud deploys no electable nodes to the region.",
									},
									"priority": schema.Int64Attribute{
										Computed:            true,
										Description:         "Precedence is given to this region when a primary election occurs. If your **regionConfigs** has only **readOnlySpecs**, **analyticsSpecs**, or both, set this value to `0`. If you have multiple **regionConfigs** objects (your cluster is multi-region or multi-cloud), they must have priorities in descending order. The highest priority is `7`.\n\n**Example:** If you have three regions, their priorities would be `7`, `6`, and `5` respectively. If you added two more regions for supporting electable nodes, the priorities of those regions would be `4` and `3` respectively.",
										MarkdownDescription: "Precedence is given to this region when a primary election occurs. If your **regionConfigs** has only **readOnlySpecs**, **analyticsSpecs**, or both, set this value to `0`. If you have multiple **regionConfigs** objects (your cluster is multi-region or multi-cloud), they must have priorities in descending order. The highest priority is `7`.\n\n**Example:** If you have three regions, their priorities would be `7`, `6`, and `5` respectively. If you added two more regions for supporting electable nodes, the priorities of those regions would be `4` and `3` respectively.",
									},
									"provider_name": schema.StringAttribute{
										Computed:            true,
										Description:         "Cloud service provider on which MongoDB Cloud provisions the hosts. Set dedicated clusters to `AWS`, `GCP`, `AZURE` or `TENANT`.",
										MarkdownDescription: "Cloud service provider on which MongoDB Cloud provisions the hosts. Set dedicated clusters to `AWS`, `GCP`, `AZURE` or `TENANT`.",
									},
									"read_only_specs": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"disk_iops": schema.Int64Attribute{
												Computed:            true,
												Description:         "Target throughput desired for storage attached to your Azure-provisioned cluster. Change this parameter if you:\n\n- set `\"replicationSpecs[n].regionConfigs[m].providerName\" : \"Azure\"`.\n- set `\"replicationSpecs[n].regionConfigs[m].electableSpecs.instanceSize\" : \"M40\"` or greater not including `Mxx_NVME` tiers.\n\nThe maximum input/output operations per second (IOPS) depend on the selected **.instanceSize** and **.diskSizeGB**.\nThis parameter defaults to the cluster tier's standard IOPS value.\nChanging this value impacts cluster cost.",
												MarkdownDescription: "Target throughput desired for storage attached to your Azure-provisioned cluster. Change this parameter if you:\n\n- set `\"replicationSpecs[n].regionConfigs[m].providerName\" : \"Azure\"`.\n- set `\"replicationSpecs[n].regionConfigs[m].electableSpecs.instanceSize\" : \"M40\"` or greater not including `Mxx_NVME` tiers.\n\nThe maximum input/output operations per second (IOPS) depend on the selected **.instanceSize** and **.diskSizeGB**.\nThis parameter defaults to the cluster tier's standard IOPS value.\nChanging this value impacts cluster cost.",
											},
											"disk_size_gb": schema.Float64Attribute{
												Computed:            true,
												Description:         "Storage capacity of instance data volumes expressed in gigabytes. Increase this number to add capacity.\n\n This value must be equal for all shards and node types.\n\n This value is not configurable on M0/M2/M5 clusters.\n\n MongoDB Cloud requires this parameter if you set **replicationSpecs**.\n\n If you specify a disk size below the minimum (10 GB), this parameter defaults to the minimum disk size value. \n\n Storage charge calculations depend on whether you choose the default value or a custom value.\n\n The maximum value for disk storage cannot exceed 50 times the maximum RAM for the selected cluster. If you require more storage space, consider upgrading your cluster to a higher tier.",
												MarkdownDescription: "Storage capacity of instance data volumes expressed in gigabytes. Increase this number to add capacity.\n\n This value must be equal for all shards and node types.\n\n This value is not configurable on M0/M2/M5 clusters.\n\n MongoDB Cloud requires this parameter if you set **replicationSpecs**.\n\n If you specify a disk size below the minimum (10 GB), this parameter defaults to the minimum disk size value. \n\n Storage charge calculations depend on whether you choose the default value or a custom value.\n\n The maximum value for disk storage cannot exceed 50 times the maximum RAM for the selected cluster. If you require more storage space, consider upgrading your cluster to a higher tier.",
											},
											"ebs_volume_type": schema.StringAttribute{
												Computed:            true,
												Description:         "Type of storage you want to attach to your AWS-provisioned cluster.\n\n- `STANDARD` volume types can't exceed the default input/output operations per second (IOPS) rate for the selected volume size. \n\n- `PROVISIONED` volume types must fall within the allowable IOPS range for the selected volume size. You must set this value to (`PROVISIONED`) for NVMe clusters.",
												MarkdownDescription: "Type of storage you want to attach to your AWS-provisioned cluster.\n\n- `STANDARD` volume types can't exceed the default input/output operations per second (IOPS) rate for the selected volume size. \n\n- `PROVISIONED` volume types must fall within the allowable IOPS range for the selected volume size. You must set this value to (`PROVISIONED`) for NVMe clusters.",
											},
											"instance_size": schema.StringAttribute{
												Computed:            true,
												Description:         "Hardware specification for the instance sizes in this region in this shard. Each instance size has a default storage and memory capacity. Electable nodes and read-only nodes (known as \"base nodes\") within a single shard must use the same instance size. Analytics nodes can scale independently from base nodes within a shard. Both base nodes and analytics nodes can scale independently from their equivalents in other shards.",
												MarkdownDescription: "Hardware specification for the instance sizes in this region in this shard. Each instance size has a default storage and memory capacity. Electable nodes and read-only nodes (known as \"base nodes\") within a single shard must use the same instance size. Analytics nodes can scale independently from base nodes within a shard. Both base nodes and analytics nodes can scale independently from their equivalents in other shards.",
											},
											"node_count": schema.Int64Attribute{
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
										Computed:            true,
										Description:         "Hardware specifications for read-only nodes in the region. Read-only nodes can never become the primary member, but can enable local reads. If you don't specify this parameter, no read-only nodes are deployed to the region.",
										MarkdownDescription: "Hardware specifications for read-only nodes in the region. Read-only nodes can never become the primary member, but can enable local reads. If you don't specify this parameter, no read-only nodes are deployed to the region.",
									},
									"region_name": schema.StringAttribute{
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
							Computed:            true,
							Description:         "Hardware specifications for nodes set for a given region. Each **regionConfigs** object describes the region's priority in elections and the number and type of MongoDB nodes that MongoDB Cloud deploys to the region. Each **regionConfigs** object must have either an **analyticsSpecs** object, **electableSpecs** object, or **readOnlySpecs** object. Tenant clusters only require **electableSpecs. Dedicated** clusters can specify any of these specifications, but must have at least one **electableSpecs** object within a **replicationSpec**.\n\n**Example:**\n\nIf you set `\"replicationSpecs[n].regionConfigs[m].analyticsSpecs.instanceSize\" : \"M30\"`, set `\"replicationSpecs[n].regionConfigs[m].electableSpecs.instanceSize\" : `\"M30\"` if you have electable nodes and `\"replicationSpecs[n].regionConfigs[m].readOnlySpecs.instanceSize\" : `\"M30\"` if you have read-only nodes.",
							MarkdownDescription: "Hardware specifications for nodes set for a given region. Each **regionConfigs** object describes the region's priority in elections and the number and type of MongoDB nodes that MongoDB Cloud deploys to the region. Each **regionConfigs** object must have either an **analyticsSpecs** object, **electableSpecs** object, or **readOnlySpecs** object. Tenant clusters only require **electableSpecs. Dedicated** clusters can specify any of these specifications, but must have at least one **electableSpecs** object within a **replicationSpec**.\n\n**Example:**\n\nIf you set `\"replicationSpecs[n].regionConfigs[m].analyticsSpecs.instanceSize\" : \"M30\"`, set `\"replicationSpecs[n].regionConfigs[m].electableSpecs.instanceSize\" : `\"M30\"` if you have electable nodes and `\"replicationSpecs[n].regionConfigs[m].readOnlySpecs.instanceSize\" : `\"M30\"` if you have read-only nodes.",
						},
						"zone_id": schema.StringAttribute{
							Computed:            true,
							Description:         "Unique 24-hexadecimal digit string that identifies the zone in a Global Cluster. This value can be used to configure Global Cluster backup policies.",
							MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the zone in a Global Cluster. This value can be used to configure Global Cluster backup policies.",
						},
						"zone_name": schema.StringAttribute{
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
				Computed:            true,
				Description:         "List of settings that configure your cluster regions. This array has one object per shard representing node configurations in each shard. For replica sets there is only one object representing node configurations.",
				MarkdownDescription: "List of settings that configure your cluster regions. This array has one object per shard representing node configurations in each shard. For replica sets there is only one object representing node configurations.",
			},
			"root_cert_type": schema.StringAttribute{
				Computed:            true,
				Description:         "Root Certificate Authority that MongoDB Cloud cluster uses. MongoDB Cloud supports Internet Security Research Group.",
				MarkdownDescription: "Root Certificate Authority that MongoDB Cloud cluster uses. MongoDB Cloud supports Internet Security Research Group.",
			},
			"state_name": schema.StringAttribute{
				Computed:            true,
				Description:         "Human-readable label that indicates the current operating condition of this cluster.",
				MarkdownDescription: "Human-readable label that indicates the current operating condition of this cluster.",
			},
			"tags": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Computed:            true,
							Description:         "Constant that defines the set of the tag. For example, `environment` in the `environment : production` tag.",
							MarkdownDescription: "Constant that defines the set of the tag. For example, `environment` in the `environment : production` tag.",
						},
						"value": schema.StringAttribute{
							Computed:            true,
							Description:         "Variable that belongs to the set of the tag. For example, `production` in the `environment : production` tag.",
							MarkdownDescription: "Variable that belongs to the set of the tag. For example, `production` in the `environment : production` tag.",
						},
					},
					CustomType: TagsType{
						ObjectType: types.ObjectType{
							AttrTypes: TagsValue{}.AttributeTypes(ctx),
						},
					},
				},
				Computed:            true,
				Description:         "List that contains key-value pairs between 1 to 255 characters in length for tagging and categorizing the cluster.",
				MarkdownDescription: "List that contains key-value pairs between 1 to 255 characters in length for tagging and categorizing the cluster.",
			},
			"termination_protection_enabled": schema.BoolAttribute{
				Computed:            true,
				Description:         "Flag that indicates whether termination protection is enabled on the cluster. If set to `true`, MongoDB Cloud won't delete the cluster. If set to `false`, MongoDB Cloud will delete the cluster.",
				MarkdownDescription: "Flag that indicates whether termination protection is enabled on the cluster. If set to `true`, MongoDB Cloud won't delete the cluster. If set to `false`, MongoDB Cloud will delete the cluster.",
			},
			"version_release_system": schema.StringAttribute{
				Computed:            true,
				Description:         "Method by which the cluster maintains the MongoDB versions. If value is `CONTINUOUS`, you must not specify **mongoDBMajorVersion**.",
				MarkdownDescription: "Method by which the cluster maintains the MongoDB versions. If value is `CONTINUOUS`, you must not specify **mongoDBMajorVersion**.",
			},
		},
	}
}

type ModelDS struct {
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
