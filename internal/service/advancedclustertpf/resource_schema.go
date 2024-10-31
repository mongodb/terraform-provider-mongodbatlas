package advancedclustertpf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Version: 1, // TODO: as in current resource
		Attributes: map[string]schema.Attribute{
			"accept_data_risks_and_force_replica_set_reconfig": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "If reconfiguration is necessary to regain a primary due to a regional outage, submit this field alongside your topology reconfiguration to request a new regional outage resistant topology. Forced reconfigurations during an outage of the majority of electable nodes carry a risk of data loss if replicated writes (even majority committed writes) have not been replicated to the new primary node. MongoDB Atlas docs contain more information. To proceed with an operation which carries that risk, set **acceptDataRisksAndForceReplicaSetReconfig** to the current date.",
			},
			"backup_enabled": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Flag that indicates whether the cluster can perform backups. If set to `true`, the cluster can perform backups. You must set this value to `true` for NVMe clusters. Backup uses [Cloud Backups](https://docs.atlas.mongodb.com/backup/cloud-backup/overview/) for dedicated clusters and [Shared Cluster Backups](https://docs.atlas.mongodb.com/backup/shared-tier/overview/) for tenant clusters. If set to `false`, the cluster doesn't use backups.",
			},
			"bi_connector": schema.SingleNestedAttribute{
				// TODO: MaxItems: 1
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Settings needed to configure the MongoDB Connector for Business Intelligence for this cluster.",
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Computed:            true,
						Optional:            true,
						MarkdownDescription: "Flag that indicates whether MongoDB Connector for Business Intelligence is enabled on the specified cluster.",
					},
					"read_preference": schema.StringAttribute{
						Computed:            true,
						Optional:            true,
						MarkdownDescription: "Data source node designated for the MongoDB Connector for Business Intelligence on MongoDB Cloud. The MongoDB Connector for Business Intelligence on MongoDB Cloud reads data from the primary, secondary, or analytics node based on your read preferences. Defaults to `ANALYTICS` node, or `SECONDARY` if there are no `ANALYTICS` nodes.",
					},
				},
			},
			"cluster_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Configuration of nodes that comprise the cluster.",
			},
			"config_server_management_mode": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Config Server Management Mode for creating or updating a sharded cluster.\n\nWhen configured as ATLAS_MANAGED, atlas may automatically switch the cluster's config server type for optimal performance and savings.\n\nWhen configured as FIXED_TO_DEDICATED, the cluster will always use a dedicated config server.",
			},
			"config_server_type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Describes a sharded cluster's config server type.",
			},
			"connection_strings": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Collection of Uniform Resource Locators that point to the MongoDB database.",
				Attributes: map[string]schema.Attribute{
					"aws_private_link": schema.MapAttribute{ // TODO: not in current resource, decide if keep
						Computed:            true,
						MarkdownDescription: "Private endpoint-aware connection strings that use AWS-hosted clusters with Amazon Web Services (AWS) PrivateLink. Each key identifies an Amazon Web Services (AWS) interface endpoint. Each value identifies the related `mongodb://` connection string that you use to connect to MongoDB Cloud through the interface endpoint that the key names.",
						ElementType:         types.StringType,
					},
					"aws_private_link_srv": schema.MapAttribute{ // TODO: not in current resource, decide if keep
						Computed:            true,
						MarkdownDescription: "Private endpoint-aware connection strings that use AWS-hosted clusters with Amazon Web Services (AWS) PrivateLink. Each key identifies an Amazon Web Services (AWS) interface endpoint. Each value identifies the related `mongodb://` connection string that you use to connect to Atlas through the interface endpoint that the key names. If the cluster uses an optimized connection string, `awsPrivateLinkSrv` contains the optimized connection string. If the cluster has the non-optimized (legacy) connection string, `awsPrivateLinkSrv` contains the non-optimized connection string even if an optimized connection string is also present.",
						ElementType:         types.StringType,
					},
					"private": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Network peering connection strings for each interface Virtual Private Cloud (VPC) endpoint that you configured to connect to this cluster. This connection string uses the `mongodb+srv://` protocol. The resource returns this parameter once someone creates a network peering connection to this cluster. This protocol tells the application to look up the host seed list in the Domain Name System (DNS). This list synchronizes with the nodes in a cluster. If the connection string uses this Uniform Resource Identifier (URI) format, you don't need to append the seed list or change the URI if the nodes change. Use this URI format if your driver supports it. If it doesn't, use connectionStrings.private. For Amazon Web Services (AWS) clusters, this resource returns this parameter only if you enable custom DNS.",
					},
					"private_endpoint": schema.ListNestedAttribute{
						Computed:            true,
						MarkdownDescription: "List of private endpoint-aware connection strings that you can use to connect to this cluster through a private endpoint. This parameter returns only if you deployed a private endpoint to all regions to which you deployed this clusters' nodes.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"connection_string": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Private endpoint-aware connection string that uses the `mongodb://` protocol to connect to MongoDB Cloud through a private endpoint.",
								},
								"endpoints": schema.ListNestedAttribute{
									Computed:            true,
									MarkdownDescription: "List that contains the private endpoints through which you connect to MongoDB Cloud when you use **connectionStrings.privateEndpoint[n].connectionString** or **connectionStrings.privateEndpoint[n].srvConnectionString**.",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"endpoint_id": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "Unique string that the cloud provider uses to identify the private endpoint.",
											},
											"provider_name": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "Cloud provider in which MongoDB Cloud deploys the private endpoint.",
											},
											"region": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "Region where the private endpoint is deployed.",
											},
										},
									},
								},
								"srv_connection_string": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Private endpoint-aware connection string that uses the `mongodb+srv://` protocol to connect to MongoDB Cloud through a private endpoint. The `mongodb+srv` protocol tells the driver to look up the seed list of hosts in the Domain Name System (DNS). This list synchronizes with the nodes in a cluster. If the connection string uses this Uniform Resource Identifier (URI) format, you don't need to append the seed list or change the Uniform Resource Identifier (URI) if the nodes change. Use this Uniform Resource Identifier (URI) format if your application supports it. If it doesn't, use connectionStrings.privateEndpoint[n].connectionString.",
								},
								"srv_shard_optimized_connection_string": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Private endpoint-aware connection string optimized for sharded clusters that uses the `mongodb+srv://` protocol to connect to MongoDB Cloud through a private endpoint. If the connection string uses this Uniform Resource Identifier (URI) format, you don't need to change the Uniform Resource Identifier (URI) if the nodes change. Use this Uniform Resource Identifier (URI) format if your application and Atlas cluster supports it. If it doesn't, use and consult the documentation for connectionStrings.privateEndpoint[n].srvConnectionString.",
								},
								"type": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "MongoDB process type to which your application connects. Use `MONGOD` for replica sets and `MONGOS` for sharded clusters.",
								},
							},
						},
					},
					"private_srv": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Network peering connection strings for each interface Virtual Private Cloud (VPC) endpoint that you configured to connect to this cluster. This connection string uses the `mongodb+srv://` protocol. The resource returns this parameter when someone creates a network peering connection to this cluster. This protocol tells the application to look up the host seed list in the Domain Name System (DNS). This list synchronizes with the nodes in a cluster. If the connection string uses this Uniform Resource Identifier (URI) format, you don't need to append the seed list or change the Uniform Resource Identifier (URI) if the nodes change. Use this Uniform Resource Identifier (URI) format if your driver supports it. If it doesn't, use `connectionStrings.private`. For Amazon Web Services (AWS) clusters, this parameter returns only if you [enable custom DNS](https://docs.atlas.mongodb.com/reference/api/aws-custom-dns-update/).",
					},
					"standard": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Public connection string that you can use to connect to this cluster. This connection string uses the `mongodb://` protocol.",
					},
					"standard_srv": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Public connection string that you can use to connect to this cluster. This connection string uses the `mongodb+srv://` protocol.",
					},
				},
			},
			"create_date": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Date and time when MongoDB Cloud created this cluster. This parameter expresses its value in ISO 8601 format in UTC.",
			},
			"disk_warming_mode": schema.StringAttribute{ // TODO: not in current resource, decide if keep
				Computed:            true,
				MarkdownDescription: "Disk warming mode selection.",
			},
			"encryption_at_rest_provider": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Cloud service provider that manages your customer keys to provide an additional layer of encryption at rest for the cluster. To enable customer key management for encryption at rest, the cluster **replicationSpecs[n].regionConfigs[m].{type}Specs.instanceSize** setting must be `M10` or higher and `\"backupEnabled\" : false` or omitted entirely.",
			},
			"feature_compatibility_version": schema.StringAttribute{ // TODO: not in current resource, decide if keep
				Computed:            true,
				MarkdownDescription: "Feature compatibility version of the cluster.",
			},
			"feature_compatibility_version_expiration_date": schema.StringAttribute{ // TODO: not in current resource, decide if keep
				Computed:            true,
				MarkdownDescription: "Feature compatibility version expiration date.",
			},
			"global_cluster_self_managed_sharding": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Set this field to configure the Sharding Management Mode when creating a new Global Cluster.\n\nWhen set to false, the management mode is set to Atlas-Managed Sharding. This mode fully manages the sharding of your Global Cluster and is built to provide a seamless deployment experience.\n\nWhen set to true, the management mode is set to Self-Managed Sharding. This mode leaves the management of shards in your hands and is built to provide an advanced and flexible deployment experience.\n\nThis setting cannot be changed once the cluster is deployed.",
			},
			"project_id": schema.StringAttribute{ // TODO: fail if trying to update
				Required:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.",
			},
			"cluster_id": schema.StringAttribute{ // TODO: was generated as id
				Computed:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the cluster.",
			},
			"labels": schema.ListNestedAttribute{ // TODO: database_user is using SetNestedBlock, probably better to align
				Computed:            true, // TODO: not sure if it should be computed
				Optional:            true,
				MarkdownDescription: "Collection of key-value pairs between 1 to 255 characters in length that tag and categorize the cluster. The MongoDB Cloud console doesn't display your labels.\n\nCluster labels are deprecated and will be removed in a future release. We strongly recommend that you use [resource tags](https://dochub.mongodb.org/core/add-cluster-tag-atlas) instead.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Computed:            true,
							Optional:            true,
							MarkdownDescription: "Key applied to tag and categorize this component.",
						},
						"value": schema.StringAttribute{
							Computed:            true,
							Optional:            true,
							MarkdownDescription: "Value set to the Key applied to tag and categorize this component.",
						},
					},
				},
			},
			"mongo_db_employee_access_grant": schema.SingleNestedAttribute{ // TODO: not in current resource, already in mongodbemployeeaccessgrant, will probably delete, was generated as mongo_dbemployee_access_grant
				Computed:            true,
				MarkdownDescription: "MongoDB employee granted access level and expiration for a cluster.",
				Attributes: map[string]schema.Attribute{
					"expiration_time": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Expiration date for the employee access grant.",
					},
					"grant_type": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Level of access to grant to MongoDB Employees.",
					},
				},
			},
			"mongo_db_major_version": schema.StringAttribute{ // TODO: was generated as mongo_dbmajor_version
				// TODO: watch out new error, Error code: "ATLAS_CLUSTER_VERSION_DEPRECATED") Detail: MongoDB version is deprecated in Atlas. Reason: Bad Request. Params: [], BadRequestDetail:
				Computed: true,
				Optional: true,
				// TODO: StateFunc: FormatMongoDBMajorVersion,
				MarkdownDescription: "MongoDB major version of the cluster.\n\nOn creation: Choose from the available versions of MongoDB, or leave unspecified for the current recommended default in the MongoDB Cloud platform. The recommended version is a recent Long Term Support version. The default is not guaranteed to be the most recently released version throughout the entire release cycle. For versions available in a specific project, see the linked documentation or use the API endpoint for [project LTS versions endpoint](#tag/Projects/operation/getProjectLTSVersions).\n\n On update: Increase version only by 1 major version at a time. If the cluster is pinned to a MongoDB feature compatibility version exactly one major version below the current MongoDB version, the MongoDB version can be downgraded to the previous major version.",
			},
			"mongo_db_version": schema.StringAttribute{ // TODO: was generated as mongo_dbversion
				Computed:            true,
				MarkdownDescription: "Version of MongoDB that the cluster runs.",
			},
			"name": schema.StringAttribute{ // TODO: fail if trying to update
				Required:            true,
				MarkdownDescription: "Human-readable label that identifies this cluster.",
			},
			"paused": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Flag that indicates whether the cluster is paused.",
			},
			"pit_enabled": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Flag that indicates whether the cluster uses continuous cloud backups.",
			},
			"redact_client_log_data": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Enable or disable log redaction.\n\nThis setting configures the ``mongod`` or ``mongos`` to redact any document field contents from a message accompanying a given log event before logging. This prevents the program from writing potentially sensitive data stored on the database to the diagnostic log. Metadata such as error or operation codes, line numbers, and source file names are still visible in the logs.\n\nUse ``redactClientLogData`` in conjunction with Encryption at Rest and TLS/SSL (Transport Encryption) to assist compliance with regulatory requirements.\n\n*Note*: changing this setting on a cluster will trigger a rolling restart as soon as the cluster is updated.",
			},
			"replica_set_scaling_strategy": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Set this field to configure the replica set scaling mode for your cluster.\n\nBy default, Atlas scales under WORKLOAD_TYPE. This mode allows Atlas to scale your analytics nodes in parallel to your operational nodes.\n\nWhen configured as SEQUENTIAL, Atlas scales all nodes sequentially. This mode is intended for steady-state workloads and applications performing latency-sensitive secondary reads.\n\nWhen configured as NODE_TYPE, Atlas scales your electable nodes in parallel with your read-only and analytics nodes. This mode is intended for large, dynamic workloads requiring frequent and timely cluster tier scaling. This is the fastest scaling strategy, but it might impact latency of workloads when performing extensive secondary reads.",
			},
			"replication_specs": schema.ListNestedAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "List of settings that configure your cluster regions. This array has one object per shard representing node configurations in each shard. For replica sets there is only one object representing node configurations.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							// TODO: Deprecated: DeprecationMsgOldSchema,
							Computed:            true,
							MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the replication object for a shard in a Cluster. If you include existing shard replication configurations in the request, you must specify this parameter. If you add a new shard to an existing Cluster, you may specify this parameter. The request deletes any existing shards  in the Cluster that you exclude from the request. This corresponds to Shard ID displayed in the UI.",
						},
						"container_id": schema.MapAttribute{ // TODO: added as in current resource
							ElementType:         types.StringType,
							Optional:            true,
							MarkdownDescription: "container_id", // TODO: add description
						},
						"external_id": schema.StringAttribute{ // TODO: added as in current resource
							Computed:            true,
							MarkdownDescription: "external_id", // TODO: add description
						},
						"num_shards": schema.Int64Attribute{ // TODO: added as in current resource
							// TODO: Deprecated: DeprecationMsgOldSchema,
							// TODO: not sure if add valitadation here: ValidateFunc: validation.IntBetween(1, 50),
							Computed:            true,
							Optional:            true,
							Default:             int64default.StaticInt64(1),
							MarkdownDescription: "num_shards", // TODO: add description
						},
						"region_configs": schema.ListNestedAttribute{
							Computed:            true,
							Optional:            true,
							MarkdownDescription: "Hardware specifications for nodes set for a given region. Each **regionConfigs** object describes the region's priority in elections and the number and type of MongoDB nodes that MongoDB Cloud deploys to the region. Each **regionConfigs** object must have either an **analyticsSpecs** object, **electableSpecs** object, or **readOnlySpecs** object. Tenant clusters only require **electableSpecs. Dedicated** clusters can specify any of these specifications, but must have at least one **electableSpecs** object within a **replicationSpec**.\n\n**Example:**\n\nIf you set `\"replicationSpecs[n].regionConfigs[m].analyticsSpecs.instanceSize\" : \"M30\"`, set `\"replicationSpecs[n].regionConfigs[m].electableSpecs.instanceSize\" : `\"M30\"` if you have electable nodes and `\"replicationSpecs[n].regionConfigs[m].readOnlySpecs.instanceSize\" : `\"M30\"` if you have read-only nodes.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"analytics_auto_scaling": AutoScalingSchema(),
									"analytics_specs":        SpecsSchema("Hardware specifications for read-only nodes in the region. Read-only nodes can never become the primary member, but can enable local reads. If you don't specify this parameter, no read-only nodes are deployed to the region."),
									"auto_scaling":           AutoScalingSchema(),
									"backing_provider_name": schema.StringAttribute{
										Optional:            true,
										MarkdownDescription: "Cloud service provider on which MongoDB Cloud provisioned the multi-tenant cluster. The resource returns this parameter when **providerName** is `TENANT` and **electableSpecs.instanceSize** is `M0`, `M2` or `M5`.",
									},
									"electable_specs": SpecsSchema("Hardware specifications for all electable nodes deployed in the region. Electable nodes can become the primary and can enable local reads. If you don't specify this option, MongoDB Cloud deploys no electable nodes to the region."),
									"priority": schema.Int64Attribute{
										Required:            true,
										MarkdownDescription: "Precedence is given to this region when a primary election occurs. If your **regionConfigs** has only **readOnlySpecs**, **analyticsSpecs**, or both, set this value to `0`. If you have multiple **regionConfigs** objects (your cluster is multi-region or multi-cloud), they must have priorities in descending order. The highest priority is `7`.\n\n**Example:** If you have three regions, their priorities would be `7`, `6`, and `5` respectively. If you added two more regions for supporting electable nodes, the priorities of those regions would be `4` and `3` respectively.",
									},
									"provider_name": schema.StringAttribute{
										// TODO: probably leave validation just in the server, ValidateDiagFunc: validate.StringIsUppercase(),
										Required:            true,
										MarkdownDescription: "Cloud service provider on which MongoDB Cloud provisions the hosts. Set dedicated clusters to `AWS`, `GCP`, `AZURE` or `TENANT`.",
									},
									"read_only_specs": SpecsSchema("Hardware specifications for read-only nodes in the region. Read-only nodes can never become the primary member, but can enable local reads. If you don't specify this parameter, no read-only nodes are deployed to the region."),
									"region_name": schema.StringAttribute{
										// TODO: probably leave validation just in the server, ValidateDiagFunc: validate.StringIsUppercase(),
										Required:            true,
										MarkdownDescription: "Physical location of your MongoDB cluster nodes. The region you choose can affect network latency for clients accessing your databases. The region name is only returned in the response for single-region clusters. When MongoDB Cloud deploys a dedicated cluster, it checks if a VPC or VPC connection exists for that provider and region. If not, MongoDB Cloud creates them as part of the deployment. It assigns the VPC a Classless Inter-Domain Routing (CIDR) block. To limit a new VPC peering connection to one Classless Inter-Domain Routing (CIDR) block and region, create the connection first. Deploy the cluster after the connection starts. GCP Clusters and Multi-region clusters require one VPC peering connection for each region. MongoDB nodes can use only the peering connection that resides in the same region as the nodes to communicate with the peered VPC.",
									},
								},
							},
						},
						"zone_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the zone in a Global Cluster. This value can be used to configure Global Cluster backup policies.",
						},
						"zone_name": schema.StringAttribute{
							Computed:            true,
							Optional:            true,
							Default:             stringdefault.StaticString("ZoneName managed by Terraform"), // TODO: as in current resource
							MarkdownDescription: "Human-readable label that describes the zone this shard belongs to in a Global Cluster. Provide this value only if \"clusterType\" : \"GEOSHARDED\" but not \"selfManagedSharding\" : true.",
						},
					},
				},
			},
			"root_cert_type": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Root Certificate Authority that MongoDB Cloud cluster uses. MongoDB Cloud supports Internet Security Research Group.",
			},
			"state_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Human-readable label that indicates the current operating condition of this cluster.",
			},
			"tags": schema.MapAttribute{ // TODO: was ListNestedAttribute, changed to align with flex cluster, might be breaking change
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Map that contains key-value pairs between 1 to 255 characters in length for tagging and categorizing the instance.",
			},
			"termination_protection_enabled": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Flag that indicates whether termination protection is enabled on the cluster. If set to `true`, MongoDB Cloud won't delete the cluster. If set to `false`, MongoDB Cloud will delete the cluster.",
			},
			"version_release_system": schema.StringAttribute{
				// TODO: probably leave validation just in the server, ValidateFunc: validation.StringInSlice([]string{"LTS", "CONTINUOUS"}, false),
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Method by which the cluster maintains the MongoDB versions. If value is `CONTINUOUS`, you must not specify **mongoDBMajorVersion**.",
			},
			"retain_backups_enabled": schema.BoolAttribute{ // TODO: not exposed in API, used in Delete operation
				Optional:            true,
				MarkdownDescription: "Flag that indicates whether to retain backup snapshots for the deleted dedicated cluster.",
			},
			"disk_size_gb": schema.Float64Attribute{ // TODO: not exposed in latest API, deprecated in root in current resource
				// Deprecated: DeprecationMsgOldSchema,
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Storage capacity of instance data volumes expressed in gigabytes. Increase this number to add capacity.\n\n This value must be equal for all shards and node types.\n\n This value is not configurable on M0/M2/M5 clusters.\n\n MongoDB Cloud requires this parameter if you set **replicationSpecs**.\n\n If you specify a disk size below the minimum (10 GB), this parameter defaults to the minimum disk size value. \n\n Storage charge calculations depend on whether you choose the default value or a custom value.\n\n The maximum value for disk storage cannot exceed 50 times the maximum RAM for the selected cluster. If you require more storage space, consider upgrading your cluster to a higher tier.",
			},
			"advanced_configuration": AdvancedConfigurationSchema(ctx),
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
			}),
		},
	}
}

func AutoScalingSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		// TODO: MaxItems: 1
		Computed:            true,
		Optional:            true,
		MarkdownDescription: "Options that determine how this cluster handles resource scaling.",
		Attributes: map[string]schema.Attribute{
			"compute_enabled": schema.BoolAttribute{ // TODO: was nested in compute
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Flag that indicates whether someone enabled instance size auto-scaling.\n\n- Set to `true` to enable instance size auto-scaling. If enabled, you must specify a value for **replicationSpecs[n].regionConfigs[m].autoScaling.compute.maxInstanceSize**.\n- Set to `false` to disable instance size automatic scaling.",
			},
			"compute_max_instance_size": schema.StringAttribute{ // TODO: was nested in compute
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Minimum instance size to which your cluster can automatically scale. MongoDB Cloud requires this parameter if `\"replicationSpecs[n].regionConfigs[m].autoScaling.compute.scaleDownEnabled\" : true`.",
			},
			"compute_min_instance_size": schema.StringAttribute{ // TODO: was nested in compute
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Minimum instance size to which your cluster can automatically scale. MongoDB Cloud requires this parameter if `\"replicationSpecs[n].regionConfigs[m].autoScaling.compute.scaleDownEnabled\" : true`.",
			},
			"compute_scale_down_enabled": schema.BoolAttribute{ // TODO: was nested in compute
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Flag that indicates whether the instance size may scale down. MongoDB Cloud requires this parameter if `\"replicationSpecs[n].regionConfigs[m].autoScaling.compute.enabled\" : true`. If you enable this option, specify a value for **replicationSpecs[n].regionConfigs[m].autoScaling.compute.minInstanceSize**.",
			},
			"disk_gb_enabled": schema.BoolAttribute{ // TODO: was nested in disk_gb
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Flag that indicates whether this cluster enables disk auto-scaling. The maximum memory allowed for the selected cluster tier and the oplog size can limit storage auto-scaling.",
			},
		},
	}
}

func SpecsSchema(markdownDescription string) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		// TODO: MaxItems: 1
		Computed:            true,
		Optional:            true,
		MarkdownDescription: markdownDescription,
		Attributes: map[string]schema.Attribute{
			"disk_iops": schema.Int64Attribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Target throughput desired for storage attached to your Azure-provisioned cluster. Change this parameter if you:\n\n- set `\"replicationSpecs[n].regionConfigs[m].providerName\" : \"Azure\"`.\n- set `\"replicationSpecs[n].regionConfigs[m].electableSpecs.instanceSize\" : \"M40\"` or greater not including `Mxx_NVME` tiers.\n\nThe maximum input/output operations per second (IOPS) depend on the selected **.instanceSize** and **.diskSizeGB**.\nThis parameter defaults to the cluster tier's standard IOPS value.\nChanging this value impacts cluster cost.",
			},
			"disk_size_gb": schema.Float64Attribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Storage capacity of instance data volumes expressed in gigabytes. Increase this number to add capacity.\n\n This value must be equal for all shards and node types.\n\n This value is not configurable on M0/M2/M5 clusters.\n\n MongoDB Cloud requires this parameter if you set **replicationSpecs**.\n\n If you specify a disk size below the minimum (10 GB), this parameter defaults to the minimum disk size value. \n\n Storage charge calculations depend on whether you choose the default value or a custom value.\n\n The maximum value for disk storage cannot exceed 50 times the maximum RAM for the selected cluster. If you require more storage space, consider upgrading your cluster to a higher tier.",
			},
			"ebs_volume_type": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Type of storage you want to attach to your AWS-provisioned cluster.\n\n- `STANDARD` volume types can't exceed the default input/output operations per second (IOPS) rate for the selected volume size. \n\n- `PROVISIONED` volume types must fall within the allowable IOPS range for the selected volume size. You must set this value to (`PROVISIONED`) for NVMe clusters.",
			},
			"instance_size": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Hardware specification for the instance sizes in this region in this shard. Each instance size has a default storage and memory capacity. Electable nodes and read-only nodes (known as \"base nodes\") within a single shard must use the same instance size. Analytics nodes can scale independently from base nodes within a shard. Both base nodes and analytics nodes can scale independently from their equivalents in other shards.",
			},
			"node_count": schema.Int64Attribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Number of nodes of the given type for MongoDB Cloud to deploy to the region.",
			},
		},
	}
}

// TODO: generated from processArgs API endpoint
func AdvancedConfigurationSchema(ctx context.Context) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Computed: true,
		Optional: true,
		// TODO: MaxItems: 1,
		MarkdownDescription: "advanced_configuration", // TODO: add description
		Attributes: map[string]schema.Attribute{
			"change_stream_options_pre_and_post_images_expire_after_seconds": schema.Int64Attribute{
				Computed:            true,
				Optional:            true,
				Default:             int64default.StaticInt64(-1), // TODO: think if default in the server only
				MarkdownDescription: "The minimum pre- and post-image retention time in seconds.",
			},
			"chunk_migration_concurrency": schema.Int64Attribute{ // TODO: not in current resource, decide if keep
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Number of threads on the source shard and the receiving shard for chunk migration. The number of threads should not exceed the half the total number of CPU cores in the sharded cluster.",
			},
			"default_max_time_ms": schema.Int64Attribute{ // TODO: not in current resource, decide if keep
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Default time limit in milliseconds for individual read operations to complete.",
			},
			"default_write_concern": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Default level of acknowledgment requested from MongoDB for write operations when none is specified by the driver.",
			},
			"javascript_enabled": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Flag that indicates whether the cluster allows execution of operations that perform server-side executions of JavaScript. When using 8.0+, we recommend disabling server-side JavaScript and using operators of aggregation pipeline as more performant alternative.",
			},
			"minimum_enabled_tls_protocol": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Minimum Transport Layer Security (TLS) version that the cluster accepts for incoming connections. Clusters using TLS 1.0 or 1.1 should consider setting TLS 1.2 as the minimum TLS protocol version.",
			},
			"no_table_scan": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Flag that indicates whether the cluster disables executing any query that requires a collection scan to return results.",
			},
			"oplog_min_retention_hours": schema.Float64Attribute{
				Optional:            true,
				MarkdownDescription: "Minimum retention window for cluster's oplog expressed in hours. A value of null indicates that the cluster uses the default minimum oplog window that MongoDB Cloud calculates.",
			},
			"oplog_size_mb": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Storage limit of cluster's oplog expressed in megabytes. A value of null indicates that the cluster uses the default oplog size that MongoDB Cloud calculates.",
			},
			"query_stats_log_verbosity": schema.Int64Attribute{ // TODO: not in current resource, decide if keep
				Optional:            true,
				MarkdownDescription: "May be set to 1 (disabled) or 3 (enabled). When set to 3, Atlas will include redacted and anonymized $queryStats output in MongoDB logs. $queryStats output does not contain literals or field values. Enabling this setting might impact the performance of your cluster.",
			},
			"sample_refresh_interval_bi_connector": schema.Int64Attribute{ // TODO was sample_refresh_interval_biconnector
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Interval in seconds at which the mongosqld process re-samples data to create its relational schema.",
			},
			"sample_size_bi_connector": schema.Int64Attribute{ // TODO was sample_size_biconnector
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Number of documents per database to sample when gathering schema information.",
			},
			"transaction_lifetime_limit_seconds": schema.Int64Attribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Lifetime, in seconds, of multi-document transactions. Atlas considers the transactions that exceed this limit as expired and so aborts them through a periodic cleanup process.",
			},
			"default_read_concern": schema.StringAttribute{ // TODO: not exposed in latest API, deprecated in current resource
				// TODO: Deprecated: DeprecationMsgOldSchema,
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "default_read_concern", // TODO: add description
			},
			"fail_index_key_too_long": schema.BoolAttribute{ // TODO: not exposed in latest API, deprecated in current resource
				// TODO: Deprecated: DeprecationMsgOldSchema,
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "fail_index_key_too_long", // TODO: add description
			},
		},
	}
}

type TFModel struct {
	DiskSizeGB                                types.Float64  `tfsdk:"disk_size_gb"`
	Labels                                    types.List     `tfsdk:"labels"`
	ReplicationSpecs                          types.List     `tfsdk:"replication_specs"`
	Tags                                      types.Map      `tfsdk:"tags"`
	DiskWarmingMode                           types.String   `tfsdk:"disk_warming_mode"`
	StateName                                 types.String   `tfsdk:"state_name"`
	ConnectionStrings                         types.Object   `tfsdk:"connection_strings"`
	CreateDate                                types.String   `tfsdk:"create_date"`
	AcceptDataRisksAndForceReplicaSetReconfig types.String   `tfsdk:"accept_data_risks_and_force_replica_set_reconfig"`
	EncryptionAtRestProvider                  types.String   `tfsdk:"encryption_at_rest_provider"`
	FeatureCompatibilityVersion               types.String   `tfsdk:"feature_compatibility_version"`
	FeatureCompatibilityVersionExpirationDate types.String   `tfsdk:"feature_compatibility_version_expiration_date"`
	Timeouts                                  timeouts.Value `tfsdk:"timeouts"`
	ProjectID                                 types.String   `tfsdk:"project_id"`
	ClusterID                                 types.String   `tfsdk:"cluster_id"`
	ConfigServerManagementMode                types.String   `tfsdk:"config_server_management_mode"`
	MongoDBEmployeeAccessGrant                types.Object   `tfsdk:"mongo_db_employee_access_grant"`
	MongoDBMajorVersion                       types.String   `tfsdk:"mongo_db_major_version"`
	MongoDBVersion                            types.String   `tfsdk:"mongo_db_version"`
	Name                                      types.String   `tfsdk:"name"`
	VersionReleaseSystem                      types.String   `tfsdk:"version_release_system"`
	BiConnector                               types.Object   `tfsdk:"bi_connector"`
	ConfigServerType                          types.String   `tfsdk:"config_server_type"`
	ReplicaSetScalingStrategy                 types.String   `tfsdk:"replica_set_scaling_strategy"`
	ClusterType                               types.String   `tfsdk:"cluster_type"`
	RootCertType                              types.String   `tfsdk:"root_cert_type"`
	RedactClientLogData                       types.Bool     `tfsdk:"redact_client_log_data"`
	PitEnabled                                types.Bool     `tfsdk:"pit_enabled"`
	TerminationProtectionEnabled              types.Bool     `tfsdk:"termination_protection_enabled"`
	Paused                                    types.Bool     `tfsdk:"paused"`
	RetainBackupsEnabled                      types.Bool     `tfsdk:"retain_backups_enabled"`
	BackupEnabled                             types.Bool     `tfsdk:"backup_enabled"`
	GlobalClusterSelfManagedSharding          types.Bool     `tfsdk:"global_cluster_self_managed_sharding"`
	AdvancedConfiguration                     types.Object   `tfsdk:"advanced_configuration"`
}

type TFBiConnectorModel struct {
	ReadPreference types.String `tfsdk:"read_preference"`
	Enabled        types.Bool   `tfsdk:"enabled"`
}

var BiConnectorObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"enabled":         types.BoolType,
	"read_preference": types.StringType,
}}

type TFConnectionStringsModel struct {
	AwsPrivateLink    types.Map    `tfsdk:"aws_private_link"`
	AwsPrivateLinkSrv types.Map    `tfsdk:"aws_private_link_srv"`
	Private           types.String `tfsdk:"private"`
	PrivateEndpoint   types.List   `tfsdk:"private_endpoint"`
	PrivateSrv        types.String `tfsdk:"private_srv"`
	Standard          types.String `tfsdk:"standard"`
	StandardSrv       types.String `tfsdk:"standard_srv"`
}

var ConnectionStringsObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"aws_private_link":     types.MapType{ElemType: types.StringType},
	"aws_private_link_srv": types.MapType{ElemType: types.StringType},
	"private":              types.StringType,
	"private_endpoint":     types.ListType{ElemType: types.StringType},
	"private_srv":          types.StringType,
	"standard":             types.StringType,
	"standard_srv":         types.StringType,
}}

type TFPrivateEndpointModel struct {
	ConnectionString                  types.String `tfsdk:"connection_string"`
	Endpoints                         types.List   `tfsdk:"endpoints"`
	SrvConnectionString               types.String `tfsdk:"srv_connection_string"`
	SrvShardOptimizedConnectionString types.String `tfsdk:"srv_shard_optimized_connection_string"`
	Type                              types.String `tfsdk:"type"`
}

var PrivateEndpointObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"connection_string":                     types.StringType,
	"endpoints":                             types.ListType{ElemType: types.StringType},
	"srv_connection_string":                 types.StringType,
	"srv_shard_optimized_connection_string": types.StringType,
	"type":                                  types.StringType,
}}

type TFEndpointsModel struct {
	EndpointId   types.String `tfsdk:"endpoint_id"`
	ProviderName types.String `tfsdk:"provider_name"`
	Region       types.String `tfsdk:"region"`
}

var EndpointsObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"endpoint_id":   types.StringType,
	"provider_name": types.StringType,
	"region":        types.StringType,
}}

type TFLabelsModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

var LabelsObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"key":   types.StringType,
	"value": types.StringType,
}}

type TFMongoDbemployeeAccessGrantModel struct {
	ExpirationTime types.String `tfsdk:"expiration_time"`
	GrantType      types.String `tfsdk:"grant_type"`
}

var MongoDbemployeeAccessGrantObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"expiration_time": types.StringType,
	"grant_type":      types.StringType,
}}

type TFReplicationSpecsModel struct {
	Id            types.String `tfsdk:"id"`
	RegionConfigs types.List   `tfsdk:"region_configs"`
	ZoneId        types.String `tfsdk:"zone_id"`
	ZoneName      types.String `tfsdk:"zone_name"`
}

var ReplicationSpecsObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"id":             types.StringType,
	"region_configs": types.ListType{ElemType: types.StringType},
	"zone_id":        types.StringType,
	"zone_name":      types.StringType,
}}

type TFRegionConfigsModel struct {
	AnalyticsAutoScaling types.Object `tfsdk:"analytics_auto_scaling"`
	AnalyticsSpecs       types.Object `tfsdk:"analytics_specs"`
	AutoScaling          types.Object `tfsdk:"auto_scaling"`
	BackingProviderName  types.String `tfsdk:"backing_provider_name"`
	ElectableSpecs       types.Object `tfsdk:"electable_specs"`
	ProviderName         types.String `tfsdk:"provider_name"`
	ReadOnlySpecs        types.Object `tfsdk:"read_only_specs"`
	RegionName           types.String `tfsdk:"region_name"`
	Priority             types.Int64  `tfsdk:"priority"`
}

var RegionConfigsObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"analytics_auto_scaling": AutoScalingObjType,
	"analytics_specs":        SpecsObjType,
	"auto_scaling":           AutoScalingObjType,
	"backing_provider_name":  types.StringType,
	"electable_specs":        SpecsObjType,
	"priority":               types.Int64Type,
	"provider_name":          types.StringType,
	"read_only_specs":        SpecsObjType,
	"region_name":            types.StringType,
}}

type TFAutoScalingModel struct {
	ComputeMaxInstanceSize  types.String `tfsdk:"max_instance_size"`
	ComputeMinInstanceSize  types.String `tfsdk:"min_instance_size"`
	ComputeEnabled          types.Bool   `tfsdk:"enabled"`
	ComputeScaleDownEnabled types.Bool   `tfsdk:"scale_down_enabled"`
	DiskGBEnabled           types.Bool   `tfsdk:"disk_gb_enabled"`
}

var AutoScalingObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"compute_enabled":            types.BoolType,
	"compute_max_instance_size":  types.StringType,
	"compute_min_instance_size":  types.StringType,
	"compute_scale_down_enabled": types.BoolType,
	"disk_gb_enabled":            types.BoolType,
}}

type TFSpecsModel struct {
	DiskSizeGb    types.Float64 `tfsdk:"disk_size_gb"`
	EbsVolumeType types.String  `tfsdk:"ebs_volume_type"`
	InstanceSize  types.String  `tfsdk:"instance_size"`
	DiskIops      types.Int64   `tfsdk:"disk_iops"`
	NodeCount     types.Int64   `tfsdk:"node_count"`
}

var SpecsObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"disk_iops":       types.Int64Type,
	"disk_size_gb":    types.Float64Type,
	"ebs_volume_type": types.StringType,
	"instance_size":   types.StringType,
	"node_count":      types.Int64Type,
}}

type TFTagsModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

var TagsObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"key":   types.StringType,
	"value": types.StringType,
}}

type TFAdvancedConfigurationModel struct {
	OplogMinRetentionHours                                types.Float64 `tfsdk:"oplog_min_retention_hours"`
	ProjectId                                             types.String  `tfsdk:"project_id"`
	ClusterName                                           types.String  `tfsdk:"cluster_name"`
	MinimumEnabledTlsProtocol                             types.String  `tfsdk:"minimum_enabled_tls_protocol"`
	DefaultWriteConcern                                   types.String  `tfsdk:"default_write_concern"`
	DefaultMaxTimeMs                                      types.Int64   `tfsdk:"default_max_time_ms"`
	ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds types.Int64   `tfsdk:"change_stream_options_pre_and_post_images_expire_after_seconds"`
	ChunkMigrationConcurrency                             types.Int64   `tfsdk:"chunk_migration_concurrency"`
	OplogSizeMb                                           types.Int64   `tfsdk:"oplog_size_mb"`
	QueryStatsLogVerbosity                                types.Int64   `tfsdk:"query_stats_log_verbosity"`
	SampleRefreshIntervalBiconnector                      types.Int64   `tfsdk:"sample_refresh_interval_biconnector"`
	SampleSizeBiconnector                                 types.Int64   `tfsdk:"sample_size_biconnector"`
	TransactionLifetimeLimitSeconds                       types.Int64   `tfsdk:"transaction_lifetime_limit_seconds"`
	JavascriptEnabled                                     types.Bool    `tfsdk:"javascript_enabled"`
	NoTableScan                                           types.Bool    `tfsdk:"no_table_scan"`
}

var TFAdvancedConfigurationObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
	// TODO: to be implemented
}}
