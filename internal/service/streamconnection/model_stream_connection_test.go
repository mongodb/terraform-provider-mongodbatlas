package streamconnection_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streamconnection"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312013/admin"
)

const (
	connectionName = "Connection"
	clusterName    = "Cluster0"
	dummyProjectID = "111111111111111111111111"
	instanceName   = "InstanceName"
	authMechanism  = "PLAIN"
	authUsername   = "user1"
	clientID       = "auth0Client"
	clientSecret   = "secret"
	// #nosec G101
	tokenEndpointURL          = "https://your-domain.com/"
	scope                     = "read:messages write:messages"
	saslOauthbearerExtentions = "logicalCluster=cluster-kmo17m,identityPoolId=pool-l7Arl"
	method                    = "OIDC"
	securityProtocol          = "SASL_SSL"
	bootstrapServers          = "localhost:9092,another.host:9092"
	dbRole                    = "customRole"
	dbRoleType                = "CUSTOM"
	sampleConnectionName      = "sample_stream_solar"
	networkingType            = "PUBLIC"
	awslambdaConnectionName   = "aws_lambda_connection"
	sampleRoleArn             = "rn:aws:iam::123456789123:role/sample"
	httpsURL                  = "https://example.com"
)

var (
	configMap = map[string]string{
		"auto.offset.reset": "earliest",
	}
	headersMap = map[string]string{
		"header1": "value1",
	}
)

type sdkToTFModelTestCase struct {
	SDKResp                          *admin.StreamsConnection
	providedProjID                   string
	providedInstanceName             string
	providedAuthConfig               *types.Object
	providedSchemaRegistryAuthConfig *types.Object
	expectedTFModel                  *streamconnection.TFStreamConnectionModel
	name                             string
}

func TestStreamConnectionSDKToTFModel(t *testing.T) {
	var authConfigWithPasswordDefined = tfAuthenticationObject(t, authMechanism, authUsername, "raw password")
	var authConfigWithOAuth = tfAuthenticationObjectForOAuth(t, authMechanism, clientID, clientSecret, tokenEndpointURL, scope, saslOauthbearerExtentions, method)

	testCases := []sdkToTFModelTestCase{
		{
			name: "Cluster connection type SDK response",
			SDKResp: &admin.StreamsConnection{
				Name:        admin.PtrString(connectionName),
				Type:        admin.PtrString("Cluster"),
				ClusterName: admin.PtrString(clusterName),
				DbRoleToExecute: &admin.DBRoleToExecute{
					Role: admin.PtrString(dbRole),
					Type: admin.PtrString(dbRoleType),
				},
			},
			providedProjID:       dummyProjectID,
			providedInstanceName: instanceName,
			providedAuthConfig:   nil,
			expectedTFModel: &streamconnection.TFStreamConnectionModel{
				ProjectID:                    types.StringValue(dummyProjectID),
				WorkspaceName:                types.StringValue(instanceName),
				ConnectionName:               types.StringValue(connectionName),
				Type:                         types.StringValue("Cluster"),
				ClusterName:                  types.StringValue(clusterName),
				Authentication:               types.ObjectNull(streamconnection.ConnectionAuthenticationObjectType.AttrTypes),
				Config:                       types.MapNull(types.StringType),
				Security:                     types.ObjectNull(streamconnection.ConnectionSecurityObjectType.AttrTypes),
				DBRoleToExecute:              tfDBRoleToExecuteObject(t, dbRole, dbRoleType),
				Networking:                   types.ObjectNull(streamconnection.NetworkingObjectType.AttrTypes),
				AWS:                          types.ObjectNull(streamconnection.AWSObjectType.AttrTypes),
				Headers:                      types.MapNull(types.StringType),
				SchemaRegistryURLs:           types.ListNull(types.StringType),
				SchemaRegistryAuthentication: types.ObjectNull(streamconnection.SchemaRegistryAuthenticationObjectType.AttrTypes),
			},
		},
		{
			name: "Cluster cross project connection type SDK response",
			SDKResp: &admin.StreamsConnection{
				Name:           admin.PtrString(connectionName),
				Type:           admin.PtrString("Cluster"),
				ClusterName:    admin.PtrString(clusterName),
				ClusterGroupId: admin.PtrString("foo"),
				DbRoleToExecute: &admin.DBRoleToExecute{
					Role: admin.PtrString(dbRole),
					Type: admin.PtrString(dbRoleType),
				},
			},
			providedProjID:       dummyProjectID,
			providedInstanceName: instanceName,
			providedAuthConfig:   nil,
			expectedTFModel: &streamconnection.TFStreamConnectionModel{
				ProjectID:                    types.StringValue(dummyProjectID),
				WorkspaceName:                types.StringValue(instanceName),
				ConnectionName:               types.StringValue(connectionName),
				Type:                         types.StringValue("Cluster"),
				ClusterName:                  types.StringValue(clusterName),
				ClusterProjectID:             types.StringValue("foo"),
				Authentication:               types.ObjectNull(streamconnection.ConnectionAuthenticationObjectType.AttrTypes),
				Config:                       types.MapNull(types.StringType),
				Security:                     types.ObjectNull(streamconnection.ConnectionSecurityObjectType.AttrTypes),
				DBRoleToExecute:              tfDBRoleToExecuteObject(t, dbRole, dbRoleType),
				Networking:                   types.ObjectNull(streamconnection.NetworkingObjectType.AttrTypes),
				AWS:                          types.ObjectNull(streamconnection.AWSObjectType.AttrTypes),
				Headers:                      types.MapNull(types.StringType),
				SchemaRegistryURLs:           types.ListNull(types.StringType),
				SchemaRegistryAuthentication: types.ObjectNull(streamconnection.SchemaRegistryAuthenticationObjectType.AttrTypes),
			},
		},
		{
			name: "Kafka connection type SDK response",
			SDKResp: &admin.StreamsConnection{
				Name: admin.PtrString(connectionName),
				Type: admin.PtrString("Kafka"),
				Authentication: &admin.StreamsKafkaAuthentication{
					Mechanism: admin.PtrString(authMechanism),
					Username:  admin.PtrString(authUsername),
				},
				BootstrapServers: admin.PtrString(bootstrapServers),
				Config:           &configMap,
				Security: &admin.StreamsKafkaSecurity{
					Protocol:                admin.PtrString(securityProtocol),
					BrokerPublicCertificate: admin.PtrString(DummyCACert),
				},
			},
			providedProjID:       dummyProjectID,
			providedInstanceName: instanceName,
			providedAuthConfig:   &authConfigWithPasswordDefined,
			expectedTFModel: &streamconnection.TFStreamConnectionModel{
				ProjectID:                    types.StringValue(dummyProjectID),
				WorkspaceName:                types.StringValue(instanceName),
				ConnectionName:               types.StringValue(connectionName),
				Type:                         types.StringValue("Kafka"),
				Authentication:               tfAuthenticationObject(t, authMechanism, authUsername, "raw password"), // password value is obtained from config, not api resp.
				BootstrapServers:             types.StringValue(bootstrapServers),
				Config:                       tfConfigMap(t, configMap),
				Security:                     tfSecurityObject(t, DummyCACert, securityProtocol),
				DBRoleToExecute:              types.ObjectNull(streamconnection.DBRoleToExecuteObjectType.AttrTypes),
				Networking:                   types.ObjectNull(streamconnection.NetworkingObjectType.AttrTypes),
				AWS:                          types.ObjectNull(streamconnection.AWSObjectType.AttrTypes),
				Headers:                      types.MapNull(types.StringType),
				SchemaRegistryURLs:           types.ListNull(types.StringType),
				SchemaRegistryAuthentication: types.ObjectNull(streamconnection.SchemaRegistryAuthenticationObjectType.AttrTypes),
			},
		},
		{
			name: "Kafka connection type SDK response for OAuthBearer authentication",
			SDKResp: &admin.StreamsConnection{
				Name: admin.PtrString(connectionName),
				Type: admin.PtrString("Kafka"),
				Authentication: &admin.StreamsKafkaAuthentication{
					Mechanism:                 admin.PtrString(authMechanism),
					Method:                    admin.PtrString(method),
					ClientId:                  admin.PtrString(clientID),
					TokenEndpointUrl:          admin.PtrString(tokenEndpointURL),
					Scope:                     admin.PtrString(scope),
					SaslOauthbearerExtensions: admin.PtrString(saslOauthbearerExtentions),
				},
				BootstrapServers: admin.PtrString(bootstrapServers),
				Config:           &configMap,
				Security: &admin.StreamsKafkaSecurity{
					Protocol:                admin.PtrString(securityProtocol),
					BrokerPublicCertificate: admin.PtrString(DummyCACert),
				},
			},
			providedProjID:       dummyProjectID,
			providedInstanceName: instanceName,
			providedAuthConfig:   &authConfigWithOAuth,
			expectedTFModel: &streamconnection.TFStreamConnectionModel{
				ProjectID:                    types.StringValue(dummyProjectID),
				WorkspaceName:                types.StringValue(instanceName),
				ConnectionName:               types.StringValue(connectionName),
				Type:                         types.StringValue("Kafka"),
				Authentication:               tfAuthenticationObjectForOAuth(t, authMechanism, clientID, clientSecret, tokenEndpointURL, scope, saslOauthbearerExtentions, method), // password value is obtained from config, not api resp.
				BootstrapServers:             types.StringValue(bootstrapServers),
				Config:                       tfConfigMap(t, configMap),
				Security:                     tfSecurityObject(t, DummyCACert, securityProtocol),
				DBRoleToExecute:              types.ObjectNull(streamconnection.DBRoleToExecuteObjectType.AttrTypes),
				Networking:                   types.ObjectNull(streamconnection.NetworkingObjectType.AttrTypes),
				AWS:                          types.ObjectNull(streamconnection.AWSObjectType.AttrTypes),
				Headers:                      types.MapNull(types.StringType),
				SchemaRegistryURLs:           types.ListNull(types.StringType),
				SchemaRegistryAuthentication: types.ObjectNull(streamconnection.SchemaRegistryAuthenticationObjectType.AttrTypes),
			},
		},
		{
			name: "Kafka connection type SDK response with no optional values provided",
			SDKResp: &admin.StreamsConnection{
				Name: admin.PtrString(connectionName),
				Type: admin.PtrString("Kafka"),
			},
			providedProjID:       dummyProjectID,
			providedInstanceName: instanceName,
			providedAuthConfig:   nil,
			expectedTFModel: &streamconnection.TFStreamConnectionModel{
				ProjectID:                    types.StringValue(dummyProjectID),
				WorkspaceName:                types.StringValue(instanceName),
				ConnectionName:               types.StringValue(connectionName),
				Type:                         types.StringValue("Kafka"),
				Authentication:               types.ObjectNull(streamconnection.ConnectionAuthenticationObjectType.AttrTypes),
				Config:                       types.MapNull(types.StringType),
				Security:                     types.ObjectNull(streamconnection.ConnectionSecurityObjectType.AttrTypes),
				DBRoleToExecute:              types.ObjectNull(streamconnection.DBRoleToExecuteObjectType.AttrTypes),
				Networking:                   types.ObjectNull(streamconnection.NetworkingObjectType.AttrTypes),
				AWS:                          types.ObjectNull(streamconnection.AWSObjectType.AttrTypes),
				Headers:                      types.MapNull(types.StringType),
				SchemaRegistryURLs:           types.ListNull(types.StringType),
				SchemaRegistryAuthentication: types.ObjectNull(streamconnection.SchemaRegistryAuthenticationObjectType.AttrTypes),
			},
		},
		{
			name: "Kafka connection type with config that does not have authentication value (case of imports and data sources)",
			SDKResp: &admin.StreamsConnection{
				Name: admin.PtrString(connectionName),
				Type: admin.PtrString("Kafka"),
				Authentication: &admin.StreamsKafkaAuthentication{
					Mechanism: admin.PtrString(authMechanism),
					Username:  admin.PtrString(authUsername),
				},
				BootstrapServers: admin.PtrString(bootstrapServers),
				Config:           &configMap,
				Security: &admin.StreamsKafkaSecurity{
					Protocol:                admin.PtrString(securityProtocol),
					BrokerPublicCertificate: admin.PtrString(DummyCACert),
				},
			},
			providedProjID:       dummyProjectID,
			providedInstanceName: instanceName,
			providedAuthConfig:   nil,
			expectedTFModel: &streamconnection.TFStreamConnectionModel{
				ProjectID:                    types.StringValue(dummyProjectID),
				WorkspaceName:                types.StringValue(instanceName),
				ConnectionName:               types.StringValue(connectionName),
				Type:                         types.StringValue("Kafka"),
				Authentication:               tfAuthenticationObjectWithNoPassword(t, authMechanism, authUsername),
				BootstrapServers:             types.StringValue(bootstrapServers),
				Config:                       tfConfigMap(t, configMap),
				Security:                     tfSecurityObject(t, DummyCACert, securityProtocol),
				DBRoleToExecute:              types.ObjectNull(streamconnection.DBRoleToExecuteObjectType.AttrTypes),
				Networking:                   types.ObjectNull(streamconnection.NetworkingObjectType.AttrTypes),
				AWS:                          types.ObjectNull(streamconnection.AWSObjectType.AttrTypes),
				Headers:                      types.MapNull(types.StringType),
				SchemaRegistryURLs:           types.ListNull(types.StringType),
				SchemaRegistryAuthentication: types.ObjectNull(streamconnection.SchemaRegistryAuthenticationObjectType.AttrTypes),
			},
		},
		{
			name: "Sample connection type sample_stream_solar sample",
			SDKResp: &admin.StreamsConnection{
				Name: admin.PtrString(sampleConnectionName),
				Type: admin.PtrString("Sample"),
			},
			providedProjID:       dummyProjectID,
			providedInstanceName: instanceName,
			expectedTFModel: &streamconnection.TFStreamConnectionModel{
				ProjectID:                    types.StringValue(dummyProjectID),
				WorkspaceName:                types.StringValue(instanceName),
				ConnectionName:               types.StringValue(sampleConnectionName),
				Type:                         types.StringValue("Sample"),
				Authentication:               types.ObjectNull(streamconnection.ConnectionAuthenticationObjectType.AttrTypes),
				Config:                       types.MapNull(types.StringType),
				Security:                     types.ObjectNull(streamconnection.ConnectionSecurityObjectType.AttrTypes),
				DBRoleToExecute:              types.ObjectNull(streamconnection.DBRoleToExecuteObjectType.AttrTypes),
				Networking:                   types.ObjectNull(streamconnection.NetworkingObjectType.AttrTypes),
				AWS:                          types.ObjectNull(streamconnection.AWSObjectType.AttrTypes),
				Headers:                      types.MapNull(types.StringType),
				SchemaRegistryURLs:           types.ListNull(types.StringType),
				SchemaRegistryAuthentication: types.ObjectNull(streamconnection.SchemaRegistryAuthenticationObjectType.AttrTypes),
			},
		},
		{
			name: "AWSLambda connection type with roleArn",
			SDKResp: &admin.StreamsConnection{
				Name: admin.PtrString(awslambdaConnectionName),
				Type: admin.PtrString("AWSLambda"),
				Aws:  &admin.StreamsAWSConnectionConfig{RoleArn: admin.PtrString(sampleRoleArn)},
			},
			providedProjID:       dummyProjectID,
			providedInstanceName: instanceName,
			expectedTFModel: &streamconnection.TFStreamConnectionModel{
				ProjectID:                    types.StringValue(dummyProjectID),
				WorkspaceName:                types.StringValue(instanceName),
				ConnectionName:               types.StringValue(awslambdaConnectionName),
				Type:                         types.StringValue("AWSLambda"),
				Authentication:               types.ObjectNull(streamconnection.ConnectionAuthenticationObjectType.AttrTypes),
				Config:                       types.MapNull(types.StringType),
				Security:                     types.ObjectNull(streamconnection.ConnectionSecurityObjectType.AttrTypes),
				DBRoleToExecute:              types.ObjectNull(streamconnection.DBRoleToExecuteObjectType.AttrTypes),
				Networking:                   types.ObjectNull(streamconnection.NetworkingObjectType.AttrTypes),
				AWS:                          tfAWSLambdaConfigObject(t, sampleRoleArn),
				Headers:                      types.MapNull(types.StringType),
				SchemaRegistryURLs:           types.ListNull(types.StringType),
				SchemaRegistryAuthentication: types.ObjectNull(streamconnection.SchemaRegistryAuthenticationObjectType.AttrTypes),
			},
		},
		{
			name: "SchemaRegistry connection type",
			SDKResp: &admin.StreamsConnection{
				Name:     admin.PtrString(connectionName),
				Type:     admin.PtrString("SchemaRegistry"),
				Provider: admin.PtrString("CONFLUENT"),
				SchemaRegistryUrls: &[]string{
					"https://schemaregistry1.com",
					"https://schemaregistry2.com",
				},
				SchemaRegistryAuthentication: &admin.SchemaRegistryAuthentication{
					Type:     "USER_INFO",
					Username: admin.PtrString("schemaUser"),
					Password: admin.PtrString("schemaPass"),
				},
			},
			providedProjID:       dummyProjectID,
			providedInstanceName: instanceName,
			expectedTFModel: &streamconnection.TFStreamConnectionModel{
				ProjectID:              types.StringValue(dummyProjectID),
				WorkspaceName:          types.StringValue(instanceName),
				ConnectionName:         types.StringValue(connectionName),
				Type:                   types.StringValue("SchemaRegistry"),
				Authentication:         types.ObjectNull(streamconnection.ConnectionAuthenticationObjectType.AttrTypes),
				Config:                 types.MapNull(types.StringType),
				Security:               types.ObjectNull(streamconnection.ConnectionSecurityObjectType.AttrTypes),
				DBRoleToExecute:        types.ObjectNull(streamconnection.DBRoleToExecuteObjectType.AttrTypes),
				Networking:             types.ObjectNull(streamconnection.NetworkingObjectType.AttrTypes),
				AWS:                    types.ObjectNull(streamconnection.AWSObjectType.AttrTypes),
				Headers:                types.MapNull(types.StringType),
				SchemaRegistryProvider: types.StringValue("CONFLUENT"),
				SchemaRegistryURLs: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("https://schemaregistry1.com"),
					types.StringValue("https://schemaregistry2.com"),
				}),
				SchemaRegistryAuthentication: tfSchemaRegistryAuthObjectNoPassword(t, "USER_INFO", "schemaUser"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel, diags := streamconnection.NewTFStreamConnection(t.Context(), tc.providedProjID, "", tc.providedInstanceName, tc.providedAuthConfig, tc.providedSchemaRegistryAuthConfig, tc.SDKResp, nil)
			if diags.HasError() {
				t.Fatalf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			tc.expectedTFModel.ID = resultModel.ID // id is auto-generated, have no way of defining within expected model
			if !assert.Equal(t, tc.expectedTFModel, resultModel) {
				t.Fatalf("created terraform model did not match expected output")
			}
		})
	}
}

// TestNewTFStreamConnectionCustomTimeoutsOverrideDefault verifies that when a user specifies
// custom timeouts in their Terraform configuration, those values override the default timeout.
// This ensures users can configure longer timeouts for slow-provisioning connections or
// shorter timeouts to fail fast.
func TestNewTFStreamConnectionCustomTimeoutsOverrideDefault(t *testing.T) {
	defaultTimeout := 20 * time.Minute
	customCreateTimeout := 30 * time.Minute
	customUpdateTimeout := 45 * time.Minute

	// User specifies custom timeouts in their Terraform config
	userConfiguredTimeouts := timeouts.Value{
		Object: types.ObjectValueMust(
			map[string]attr.Type{
				"create": types.StringType,
				"read":   types.StringType,
				"update": types.StringType,
				"delete": types.StringType,
			},
			map[string]attr.Value{
				"create": types.StringValue("30m"),
				"read":   types.StringNull(),
				"update": types.StringValue("45m"),
				"delete": types.StringNull(),
			},
		),
	}

	apiResp := &admin.StreamsConnection{
		Name: admin.PtrString("TestConnection"),
		Type: admin.PtrString("Cluster"),
	}

	resultModel, diags := streamconnection.NewTFStreamConnection(
		t.Context(),
		dummyProjectID,
		"",           // instanceName (deprecated)
		instanceName, // workspaceName
		nil,          // currAuthConfig
		nil,          // currSchemaRegistryAuthConfig
		apiResp,
		&userConfiguredTimeouts, // planTimeouts - user configured custom timeouts
	)

	require.False(t, diags.HasError(), "unexpected errors: %v", diags)
	require.NotNil(t, resultModel)

	// Verify user-configured timeouts are preserved in the model
	assert.Equal(t, userConfiguredTimeouts, resultModel.Timeouts)

	// Verify custom timeouts override the default (20m) when extracted
	createTimeout, localDiags := resultModel.Timeouts.Create(t.Context(), defaultTimeout)
	require.False(t, localDiags.HasError())
	assert.Equal(t, customCreateTimeout, createTimeout, "user-configured create timeout (30m) should override default (20m)")

	updateTimeout, localDiags := resultModel.Timeouts.Update(t.Context(), defaultTimeout)
	require.False(t, localDiags.HasError())
	assert.Equal(t, customUpdateTimeout, updateTimeout, "user-configured update timeout (45m) should override default (20m)")
}

type paginatedConnectionsSDKToTFModelTestCase struct {
	SDKResp         *admin.PaginatedApiStreamsConnection
	providedConfig  *streamconnection.TFStreamConnectionsDSModel
	expectedTFModel *streamconnection.TFStreamConnectionsDSModel
	name            string
}

func TestStreamConnectionsSDKToTFModel(t *testing.T) {
	testCases := []paginatedConnectionsSDKToTFModelTestCase{
		{
			name: "Complete SDK response with configured page options",
			SDKResp: &admin.PaginatedApiStreamsConnection{
				Results: &[]admin.StreamsConnection{
					{
						Name: admin.PtrString(connectionName),
						Type: admin.PtrString("Kafka"),
						Authentication: &admin.StreamsKafkaAuthentication{
							Mechanism: admin.PtrString(authMechanism),
							Username:  admin.PtrString(authUsername),
						},
						BootstrapServers: admin.PtrString(bootstrapServers),
						Config:           &configMap,
						Security: &admin.StreamsKafkaSecurity{
							Protocol:                admin.PtrString(securityProtocol),
							BrokerPublicCertificate: admin.PtrString(DummyCACert),
						},
						Networking: &admin.StreamsKafkaNetworking{
							Access: &admin.StreamsKafkaNetworkingAccess{
								Type: admin.PtrString(networkingType),
							},
						},
					},
					{
						Name:        admin.PtrString(connectionName),
						Type:        admin.PtrString("Cluster"),
						ClusterName: admin.PtrString(clusterName),
						DbRoleToExecute: &admin.DBRoleToExecute{
							Role: admin.PtrString(dbRole),
							Type: admin.PtrString(dbRoleType),
						},
					},
					{
						Name: admin.PtrString(sampleConnectionName),
						Type: admin.PtrString("Sample"),
					},
					{
						Name: admin.PtrString(awslambdaConnectionName),
						Type: admin.PtrString("AWSLambda"),
						Aws: &admin.StreamsAWSConnectionConfig{
							RoleArn: admin.PtrString(sampleRoleArn),
						},
					},
					{
						Name:    admin.PtrString(connectionName),
						Type:    admin.PtrString("Https"),
						Url:     admin.PtrString(httpsURL),
						Headers: &headersMap,
					},
					{
						Name:     admin.PtrString(connectionName),
						Type:     admin.PtrString("SchemaRegistry"),
						Provider: admin.PtrString("CONFLUENT"),
						SchemaRegistryUrls: &[]string{
							"https://schemaregistry1.com",
							"https://schemaregistry2.com",
						},
						SchemaRegistryAuthentication: &admin.SchemaRegistryAuthentication{
							Type:     "USER_INFO",
							Username: admin.PtrString("schemaUser"),
							Password: admin.PtrString("schemaPass"),
						},
					},
				},
				TotalCount: admin.PtrInt(6),
			},
			providedConfig: &streamconnection.TFStreamConnectionsDSModel{
				ProjectID:    types.StringValue(dummyProjectID),
				InstanceName: types.StringValue(instanceName),
				PageNum:      types.Int64Value(1),
				ItemsPerPage: types.Int64Value(3),
			},
			expectedTFModel: &streamconnection.TFStreamConnectionsDSModel{
				ProjectID:    types.StringValue(dummyProjectID),
				InstanceName: types.StringValue(instanceName),
				PageNum:      types.Int64Value(1),
				ItemsPerPage: types.Int64Value(3),
				TotalCount:   types.Int64Value(6),
				Results: []streamconnection.TFStreamConnectionModel{
					{
						ID:                           types.StringValue(fmt.Sprintf("%s-%s-%s", instanceName, dummyProjectID, connectionName)),
						ProjectID:                    types.StringValue(dummyProjectID),
						InstanceName:                 types.StringValue(instanceName),
						ConnectionName:               types.StringValue(connectionName),
						Type:                         types.StringValue("Kafka"),
						Authentication:               tfAuthenticationObjectWithNoPassword(t, authMechanism, authUsername),
						BootstrapServers:             types.StringValue(bootstrapServers),
						Config:                       tfConfigMap(t, configMap),
						Security:                     tfSecurityObject(t, DummyCACert, securityProtocol),
						DBRoleToExecute:              types.ObjectNull(streamconnection.DBRoleToExecuteObjectType.AttrTypes),
						Networking:                   tfNetworkingObject(t, networkingType, nil),
						AWS:                          types.ObjectNull(streamconnection.AWSObjectType.AttrTypes),
						Headers:                      types.MapNull(types.StringType),
						SchemaRegistryURLs:           types.ListNull(types.StringType),
						SchemaRegistryAuthentication: types.ObjectNull(streamconnection.SchemaRegistryAuthenticationObjectType.AttrTypes),
					},
					{
						ID:                           types.StringValue(fmt.Sprintf("%s-%s-%s", instanceName, dummyProjectID, connectionName)),
						ProjectID:                    types.StringValue(dummyProjectID),
						InstanceName:                 types.StringValue(instanceName),
						ConnectionName:               types.StringValue(connectionName),
						Type:                         types.StringValue("Cluster"),
						ClusterName:                  types.StringValue(clusterName),
						Authentication:               types.ObjectNull(streamconnection.ConnectionAuthenticationObjectType.AttrTypes),
						Config:                       types.MapNull(types.StringType),
						Security:                     types.ObjectNull(streamconnection.ConnectionSecurityObjectType.AttrTypes),
						DBRoleToExecute:              tfDBRoleToExecuteObject(t, dbRole, dbRoleType),
						Networking:                   types.ObjectNull(streamconnection.NetworkingObjectType.AttrTypes),
						AWS:                          types.ObjectNull(streamconnection.AWSObjectType.AttrTypes),
						Headers:                      types.MapNull(types.StringType),
						SchemaRegistryURLs:           types.ListNull(types.StringType),
						SchemaRegistryAuthentication: types.ObjectNull(streamconnection.SchemaRegistryAuthenticationObjectType.AttrTypes),
					},
					{
						ID:                           types.StringValue(fmt.Sprintf("%s-%s-%s", instanceName, dummyProjectID, sampleConnectionName)),
						ProjectID:                    types.StringValue(dummyProjectID),
						InstanceName:                 types.StringValue(instanceName),
						ConnectionName:               types.StringValue(sampleConnectionName),
						Type:                         types.StringValue("Sample"),
						ClusterName:                  types.StringNull(),
						Authentication:               types.ObjectNull(streamconnection.ConnectionAuthenticationObjectType.AttrTypes),
						Config:                       types.MapNull(types.StringType),
						Security:                     types.ObjectNull(streamconnection.ConnectionSecurityObjectType.AttrTypes),
						DBRoleToExecute:              types.ObjectNull(streamconnection.DBRoleToExecuteObjectType.AttrTypes),
						Networking:                   types.ObjectNull(streamconnection.NetworkingObjectType.AttrTypes),
						AWS:                          types.ObjectNull(streamconnection.AWSObjectType.AttrTypes),
						Headers:                      types.MapNull(types.StringType),
						SchemaRegistryURLs:           types.ListNull(types.StringType),
						SchemaRegistryAuthentication: types.ObjectNull(streamconnection.SchemaRegistryAuthenticationObjectType.AttrTypes),
					},
					{
						ID:                           types.StringValue(fmt.Sprintf("%s-%s-%s", instanceName, dummyProjectID, awslambdaConnectionName)),
						ProjectID:                    types.StringValue(dummyProjectID),
						InstanceName:                 types.StringValue(instanceName),
						ConnectionName:               types.StringValue(awslambdaConnectionName),
						Type:                         types.StringValue("AWSLambda"),
						ClusterName:                  types.StringNull(),
						Authentication:               types.ObjectNull(streamconnection.ConnectionAuthenticationObjectType.AttrTypes),
						Config:                       types.MapNull(types.StringType),
						Security:                     types.ObjectNull(streamconnection.ConnectionSecurityObjectType.AttrTypes),
						DBRoleToExecute:              types.ObjectNull(streamconnection.DBRoleToExecuteObjectType.AttrTypes),
						Networking:                   types.ObjectNull(streamconnection.NetworkingObjectType.AttrTypes),
						AWS:                          tfAWSLambdaConfigObject(t, sampleRoleArn),
						Headers:                      types.MapNull(types.StringType),
						SchemaRegistryURLs:           types.ListNull(types.StringType),
						SchemaRegistryAuthentication: types.ObjectNull(streamconnection.SchemaRegistryAuthenticationObjectType.AttrTypes),
					},
					{
						ID:                           types.StringValue(fmt.Sprintf("%s-%s-%s", instanceName, dummyProjectID, connectionName)),
						ProjectID:                    types.StringValue(dummyProjectID),
						InstanceName:                 types.StringValue(instanceName),
						ConnectionName:               types.StringValue(connectionName),
						Type:                         types.StringValue("Https"),
						ClusterName:                  types.StringNull(),
						Authentication:               types.ObjectNull(streamconnection.ConnectionAuthenticationObjectType.AttrTypes),
						Config:                       types.MapNull(types.StringType),
						Security:                     types.ObjectNull(streamconnection.ConnectionSecurityObjectType.AttrTypes),
						DBRoleToExecute:              types.ObjectNull(streamconnection.DBRoleToExecuteObjectType.AttrTypes),
						Networking:                   types.ObjectNull(streamconnection.NetworkingObjectType.AttrTypes),
						AWS:                          types.ObjectNull(streamconnection.AWSObjectType.AttrTypes),
						Headers:                      tfConfigMap(t, headersMap),
						URL:                          types.StringValue(httpsURL),
						SchemaRegistryURLs:           types.ListNull(types.StringType),
						SchemaRegistryAuthentication: types.ObjectNull(streamconnection.SchemaRegistryAuthenticationObjectType.AttrTypes),
					},
					{
						ID:                     types.StringValue(fmt.Sprintf("%s-%s-%s", instanceName, dummyProjectID, connectionName)),
						ProjectID:              types.StringValue(dummyProjectID),
						InstanceName:           types.StringValue(instanceName),
						ConnectionName:         types.StringValue(connectionName),
						Type:                   types.StringValue("SchemaRegistry"),
						ClusterName:            types.StringNull(),
						Authentication:         types.ObjectNull(streamconnection.ConnectionAuthenticationObjectType.AttrTypes),
						Config:                 types.MapNull(types.StringType),
						Security:               types.ObjectNull(streamconnection.ConnectionSecurityObjectType.AttrTypes),
						DBRoleToExecute:        types.ObjectNull(streamconnection.DBRoleToExecuteObjectType.AttrTypes),
						Networking:             types.ObjectNull(streamconnection.NetworkingObjectType.AttrTypes),
						AWS:                    types.ObjectNull(streamconnection.AWSObjectType.AttrTypes),
						Headers:                types.MapNull(types.StringType),
						SchemaRegistryProvider: types.StringValue("CONFLUENT"),
						SchemaRegistryURLs: types.ListValueMust(types.StringType, []attr.Value{
							types.StringValue("https://schemaregistry1.com"),
							types.StringValue("https://schemaregistry2.com"),
						}),
						SchemaRegistryAuthentication: tfSchemaRegistryAuthObjectNoPassword(t, "USER_INFO", "schemaUser"),
					},
				},
			},
		},
		{
			name: "Without defining page options",
			SDKResp: &admin.PaginatedApiStreamsConnection{
				Results:    &[]admin.StreamsConnection{},
				TotalCount: admin.PtrInt(0),
			},
			providedConfig: &streamconnection.TFStreamConnectionsDSModel{
				ProjectID:    types.StringValue(dummyProjectID),
				InstanceName: types.StringValue(instanceName),
			},
			expectedTFModel: &streamconnection.TFStreamConnectionsDSModel{
				ProjectID:    types.StringValue(dummyProjectID),
				InstanceName: types.StringValue(instanceName),
				PageNum:      types.Int64Null(),
				ItemsPerPage: types.Int64Null(),
				TotalCount:   types.Int64Value(0),
				Results:      []streamconnection.TFStreamConnectionModel{},
			},
		},
		{
			name: "With workspace name and no page options",
			SDKResp: &admin.PaginatedApiStreamsConnection{
				Results:    &[]admin.StreamsConnection{},
				TotalCount: admin.PtrInt(0),
			},
			providedConfig: &streamconnection.TFStreamConnectionsDSModel{
				ProjectID:     types.StringValue(dummyProjectID),
				WorkspaceName: types.StringValue(instanceName),
			},
			expectedTFModel: &streamconnection.TFStreamConnectionsDSModel{
				ProjectID:     types.StringValue(dummyProjectID),
				WorkspaceName: types.StringValue(instanceName),
				PageNum:       types.Int64Null(),
				ItemsPerPage:  types.Int64Null(),
				TotalCount:    types.Int64Value(0),
				Results:       []streamconnection.TFStreamConnectionModel{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel, diags := streamconnection.NewTFStreamConnections(t.Context(), tc.providedConfig, tc.SDKResp)
			if diags.HasError() {
				t.Fatalf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			tc.expectedTFModel.ID = resultModel.ID // id is auto-generated, have no way of defining within expected model
			if !assert.Equal(t, tc.expectedTFModel, resultModel) {
				t.Fatalf("created terraform model did not match expected output")
			}
		})
	}
}

type tfToSDKCreateModelTestCase struct {
	tfModel        *streamconnection.TFStreamConnectionModel
	expectedSDKReq *admin.StreamsConnection
	name           string
}

func TestStreamInstanceTFToSDKCreateModel(t *testing.T) {
	testCases := []tfToSDKCreateModelTestCase{
		{
			name: "Cluster type complete TF state",
			tfModel: &streamconnection.TFStreamConnectionModel{
				ProjectID:       types.StringValue(dummyProjectID),
				InstanceName:    types.StringValue(instanceName),
				ConnectionName:  types.StringValue(connectionName),
				Type:            types.StringValue("Cluster"),
				ClusterName:     types.StringValue(clusterName),
				DBRoleToExecute: tfDBRoleToExecuteObject(t, dbRole, dbRoleType),
			},
			expectedSDKReq: &admin.StreamsConnection{
				Name:        admin.PtrString(connectionName),
				Type:        admin.PtrString("Cluster"),
				ClusterName: admin.PtrString(clusterName),
				DbRoleToExecute: &admin.DBRoleToExecute{
					Role: admin.PtrString(dbRole),
					Type: admin.PtrString(dbRoleType),
				},
			},
		},
		{
			name: "Kafka type complete TF state",
			tfModel: &streamconnection.TFStreamConnectionModel{
				ProjectID:        types.StringValue(dummyProjectID),
				InstanceName:     types.StringValue(instanceName),
				ConnectionName:   types.StringValue(connectionName),
				Type:             types.StringValue("Kafka"),
				Authentication:   tfAuthenticationObject(t, authMechanism, authUsername, "raw password"),
				BootstrapServers: types.StringValue(bootstrapServers),
				Config:           tfConfigMap(t, configMap),
				Security:         tfSecurityObject(t, DummyCACert, securityProtocol),
			},
			expectedSDKReq: &admin.StreamsConnection{
				Name: admin.PtrString(connectionName),
				Type: admin.PtrString("Kafka"),
				Authentication: &admin.StreamsKafkaAuthentication{
					Mechanism: admin.PtrString(authMechanism),
					Username:  admin.PtrString(authUsername),
					Password:  admin.PtrString("raw password"),
				},
				BootstrapServers: admin.PtrString(bootstrapServers),
				Config:           &configMap,
				Security: &admin.StreamsKafkaSecurity{
					Protocol:                admin.PtrString(securityProtocol),
					BrokerPublicCertificate: admin.PtrString(DummyCACert),
				},
			},
		},
		{
			name: "Kafka type TF state with no optional attributes",
			tfModel: &streamconnection.TFStreamConnectionModel{
				ProjectID:      types.StringValue(dummyProjectID),
				InstanceName:   types.StringValue(instanceName),
				ConnectionName: types.StringValue(connectionName),
				Type:           types.StringValue("Kafka"),
			},
			expectedSDKReq: &admin.StreamsConnection{
				Name: admin.PtrString(connectionName),
				Type: admin.PtrString("Kafka"),
			},
		},
		{
			name: "Sample type TF state",
			tfModel: &streamconnection.TFStreamConnectionModel{
				ProjectID:      types.StringValue(dummyProjectID),
				InstanceName:   types.StringValue(instanceName),
				ConnectionName: types.StringValue(connectionName),
				Type:           types.StringValue("Sample"),
			},
			expectedSDKReq: &admin.StreamsConnection{
				Name: admin.PtrString(connectionName),
				Type: admin.PtrString("Sample"),
			},
		},
		{
			name: "AWSLambda type TF state",
			tfModel: &streamconnection.TFStreamConnectionModel{
				ProjectID:      types.StringValue(dummyProjectID),
				InstanceName:   types.StringValue(instanceName),
				ConnectionName: types.StringValue(awslambdaConnectionName),
				Type:           types.StringValue("AWSLambda"),
				AWS:            tfAWSLambdaConfigObject(t, sampleRoleArn),
			},
			expectedSDKReq: &admin.StreamsConnection{
				Name: admin.PtrString(awslambdaConnectionName),
				Type: admin.PtrString("AWSLambda"),
				Aws: &admin.StreamsAWSConnectionConfig{
					RoleArn: admin.PtrString(sampleRoleArn),
				},
			},
		},
		{
			name: "Https type TF state",
			tfModel: &streamconnection.TFStreamConnectionModel{
				ProjectID:      types.StringValue(dummyProjectID),
				InstanceName:   types.StringValue(instanceName),
				ConnectionName: types.StringValue(connectionName),
				Type:           types.StringValue("Https"),
				URL:            types.StringValue(httpsURL),
				Headers:        tfConfigMap(t, headersMap),
			},
			expectedSDKReq: &admin.StreamsConnection{
				Name:    admin.PtrString(connectionName),
				Type:    admin.PtrString("Https"),
				Url:     admin.PtrString(httpsURL),
				Headers: &headersMap,
			},
		},
		{
			name: "SchemaRegistry type USER_INFO TF state",
			tfModel: &streamconnection.TFStreamConnectionModel{
				ProjectID:              types.StringValue(dummyProjectID),
				InstanceName:           types.StringValue(instanceName),
				ConnectionName:         types.StringValue(connectionName),
				Type:                   types.StringValue("SchemaRegistry"),
				SchemaRegistryProvider: types.StringValue("CONFLUENT"),
				SchemaRegistryURLs: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("https://schemaregistry1.com"),
					types.StringValue("https://schemaregistry2.com"),
				}),
				SchemaRegistryAuthentication: tfSchemaRegistryAuthObject(t, "USER_INFO", "schemaUser", "schemaPass"),
			},
			expectedSDKReq: &admin.StreamsConnection{
				Name:     admin.PtrString(connectionName),
				Type:     admin.PtrString("SchemaRegistry"),
				Provider: admin.PtrString("CONFLUENT"),
				SchemaRegistryUrls: &[]string{
					"https://schemaregistry1.com",
					"https://schemaregistry2.com",
				},
				SchemaRegistryAuthentication: &admin.SchemaRegistryAuthentication{
					Type:     "USER_INFO",
					Username: admin.PtrString("schemaUser"),
					Password: admin.PtrString("schemaPass"),
				},
			},
		},
		{
			name: "SchemaRegistry type SASL_INHERIT TF state",
			tfModel: &streamconnection.TFStreamConnectionModel{
				ProjectID:              types.StringValue(dummyProjectID),
				InstanceName:           types.StringValue(instanceName),
				ConnectionName:         types.StringValue(connectionName),
				Type:                   types.StringValue("SchemaRegistry"),
				SchemaRegistryProvider: types.StringValue("CONFLUENT"),
				SchemaRegistryURLs: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("https://schemaregistry1.com"),
					types.StringValue("https://schemaregistry2.com"),
				}),
				SchemaRegistryAuthentication: tfSchemaRegistryAuthObject(t, "SASL_INHERIT", "", ""),
			},
			expectedSDKReq: &admin.StreamsConnection{
				Name:     admin.PtrString(connectionName),
				Type:     admin.PtrString("SchemaRegistry"),
				Provider: admin.PtrString("CONFLUENT"),
				SchemaRegistryUrls: &[]string{
					"https://schemaregistry1.com",
					"https://schemaregistry2.com",
				},
				SchemaRegistryAuthentication: &admin.SchemaRegistryAuthentication{
					Type: "SASL_INHERIT",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			apiReqResult, diags := streamconnection.NewStreamConnectionReq(t.Context(), tc.tfModel)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			if !assert.Equal(t, tc.expectedSDKReq, apiReqResult) {
				t.Errorf("created sdk model did not match expected output")
			}
		})
	}
}

func tfAuthenticationObject(t *testing.T, mechanism, username, password string) types.Object {
	t.Helper()
	auth, diags := types.ObjectValueFrom(t.Context(), streamconnection.ConnectionAuthenticationObjectType.AttrTypes, streamconnection.TFConnectionAuthenticationModel{
		Mechanism: types.StringValue(mechanism),
		Username:  types.StringValue(username),
		Password:  types.StringValue(password),
	})
	if diags.HasError() {
		t.Errorf("failed to create terraform data model: %s", diags.Errors()[0].Summary())
	}
	return auth
}

func tfAuthenticationObjectForOAuth(t *testing.T, mechanism, clientID, clientSecret, tokenEndpointURL, scope, saslOauthbearerExtensions, method string) types.Object {
	t.Helper()
	auth, diags := types.ObjectValueFrom(t.Context(), streamconnection.ConnectionAuthenticationObjectType.AttrTypes, streamconnection.TFConnectionAuthenticationModel{
		Mechanism:                 types.StringValue(mechanism),
		Method:                    types.StringValue(method),
		ClientID:                  types.StringValue(clientID),
		ClientSecret:              types.StringValue(clientSecret),
		TokenEndpointURL:          types.StringValue(tokenEndpointURL),
		Scope:                     types.StringValue(scope),
		SaslOauthbearerExtensions: types.StringValue(saslOauthbearerExtensions),
	})
	if diags.HasError() {
		t.Errorf("failed to create terraform data model: %s", diags.Errors()[0].Summary())
	}
	return auth
}

func tfAuthenticationObjectWithNoPassword(t *testing.T, mechanism, username string) types.Object {
	t.Helper()
	auth, diags := types.ObjectValueFrom(t.Context(), streamconnection.ConnectionAuthenticationObjectType.AttrTypes, streamconnection.TFConnectionAuthenticationModel{
		Mechanism: types.StringValue(mechanism),
		Username:  types.StringValue(username),
	})
	if diags.HasError() {
		t.Errorf("failed to create terraform data model: %s", diags.Errors()[0].Summary())
	}
	return auth
}

func tfSecurityObject(t *testing.T, brokerPublicCertificate, protocol string) types.Object {
	t.Helper()
	auth, diags := types.ObjectValueFrom(t.Context(), streamconnection.ConnectionSecurityObjectType.AttrTypes, streamconnection.TFConnectionSecurityModel{
		BrokerPublicCertificate: types.StringValue(brokerPublicCertificate),
		Protocol:                types.StringValue(protocol),
	})
	if diags.HasError() {
		t.Errorf("failed to create terraform data model: %s", diags.Errors()[0].Summary())
	}
	return auth
}

func tfConfigMap(t *testing.T, config map[string]string) types.Map {
	t.Helper()
	mapValue, diags := types.MapValueFrom(t.Context(), types.StringType, config)
	if diags.HasError() {
		t.Errorf("failed to create terraform data model: %s", diags.Errors()[0].Summary())
	}
	return mapValue
}

func tfDBRoleToExecuteObject(t *testing.T, role, roleType string) types.Object {
	t.Helper()
	auth, diags := types.ObjectValueFrom(t.Context(), streamconnection.DBRoleToExecuteObjectType.AttrTypes, streamconnection.TFDbRoleToExecuteModel{
		Role: types.StringValue(role),
		Type: types.StringValue(roleType),
	})
	if diags.HasError() {
		t.Errorf("failed to create terraform data model: %s", diags.Errors()[0].Summary())
	}
	return auth
}

func tfNetworkingObject(t *testing.T, networkingType string, connectionID *string) types.Object {
	t.Helper()
	networkingAccessModel, diags := types.ObjectValueFrom(t.Context(), streamconnection.NetworkingAccessObjectType.AttrTypes, streamconnection.TFNetworkingAccessModel{
		Type:         types.StringValue(networkingType),
		ConnectionID: types.StringPointerValue(connectionID),
	})
	if diags.HasError() {
		t.Errorf("failed to create terraform data model: %s", diags.Errors()[0].Summary())
	}
	networking, diags := types.ObjectValueFrom(t.Context(), streamconnection.NetworkingObjectType.AttrTypes, streamconnection.TFNetworkingModel{
		Access: networkingAccessModel,
	})
	if diags.HasError() {
		t.Errorf("failed to create terraform data model: %s", diags.Errors()[0].Summary())
	}
	return networking
}

func tfAWSLambdaConfigObject(t *testing.T, roleArn string) types.Object {
	t.Helper()
	aws, diags := types.ObjectValueFrom(t.Context(), streamconnection.AWSObjectType.AttrTypes, streamconnection.TFAWSModel{
		RoleArn: types.StringValue(roleArn),
	})
	if diags.HasError() {
		t.Errorf("failed to create terraform data model: %s", diags.Errors()[0].Summary())
	}
	return aws
}

func tfSchemaRegistryAuthObject(t *testing.T, authType, username, password string) types.Object {
	t.Helper()
	tfAuth := streamconnection.TFSchemaRegistryAuthenticationModel{Type: types.StringValue(authType)}
	if authType == "USER_INFO" {
		tfAuth.Username = types.StringValue(username)
		tfAuth.Password = types.StringValue(password)
	}

	auth, diags := types.ObjectValueFrom(t.Context(), streamconnection.SchemaRegistryAuthenticationObjectType.AttrTypes, tfAuth)
	if diags.HasError() {
		t.Errorf("failed to create terraform data model: %s", diags.Errors()[0].Summary())
	}
	return auth
}

func tfSchemaRegistryAuthObjectNoPassword(t *testing.T, authType, username string) types.Object {
	t.Helper()
	tfAuth := streamconnection.TFSchemaRegistryAuthenticationModel{Type: types.StringValue(authType)}
	if authType == "USER_INFO" {
		tfAuth.Username = types.StringValue(username)
		tfAuth.Password = types.StringNull()
	}

	auth, diags := types.ObjectValueFrom(t.Context(), streamconnection.SchemaRegistryAuthenticationObjectType.AttrTypes, tfAuth)
	if diags.HasError() {
		t.Errorf("failed to create terraform data model: %s", diags.Errors()[0].Summary())
	}
	return auth
}
