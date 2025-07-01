package streamconnection_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streamconnection"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20250312005/admin"
)

const (
	connectionName            = "Connection"
	typeValue                 = ""
	clusterName               = "Cluster0"
	dummyProjectID            = "111111111111111111111111"
	instanceName              = "InstanceName"
	authMechanism             = "PLAIN"
	authUsername              = "user1"
	securityProtocol          = "SASL_SSL"
	bootstrapServers          = "localhost:9092,another.host:9092"
	dbRole                    = "customRole"
	dbRoleType                = "CUSTOM"
	sampleConnectionName      = "sample_stream_solar"
	networkingType            = "PUBLIC"
	privatelinkNetworkingType = "PRIVATE_LINK"
	awslambdaConnectionName   = "aws_lambda_connection"
	awss3ConnectionName       = "aws_s3_connection"
	sampleRoleArn             = "rn:aws:iam::123456789123:role/sample"
	sampleTestBucket          = "sample_test_bucket"
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
	SDKResp              *admin.StreamsConnection
	providedProjID       string
	providedInstanceName string
	providedAuthConfig   *types.Object
	expectedTFModel      *streamconnection.TFStreamConnectionModel
	name                 string
}

func TestStreamConnectionSDKToTFModel(t *testing.T) {
	var authConfigWithPasswordDefined = tfAuthenticationObject(t, authMechanism, authUsername, "raw password")

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
				ProjectID:       types.StringValue(dummyProjectID),
				InstanceName:    types.StringValue(instanceName),
				ConnectionName:  types.StringValue(connectionName),
				Type:            types.StringValue("Cluster"),
				ClusterName:     types.StringValue(clusterName),
				Authentication:  types.ObjectNull(streamconnection.ConnectionAuthenticationObjectType.AttrTypes),
				Config:          types.MapNull(types.StringType),
				Security:        types.ObjectNull(streamconnection.ConnectionSecurityObjectType.AttrTypes),
				DBRoleToExecute: tfDBRoleToExecuteObject(t, dbRole, dbRoleType),
				Networking:      types.ObjectNull(streamconnection.NetworkingObjectType.AttrTypes),
				AWS:             types.ObjectNull(streamconnection.AWSObjectType.AttrTypes),
				Headers:         types.MapNull(types.StringType),
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
				ProjectID:        types.StringValue(dummyProjectID),
				InstanceName:     types.StringValue(instanceName),
				ConnectionName:   types.StringValue(connectionName),
				Type:             types.StringValue("Cluster"),
				ClusterName:      types.StringValue(clusterName),
				ClusterProjectID: types.StringValue("foo"),
				Authentication:   types.ObjectNull(streamconnection.ConnectionAuthenticationObjectType.AttrTypes),
				Config:           types.MapNull(types.StringType),
				Security:         types.ObjectNull(streamconnection.ConnectionSecurityObjectType.AttrTypes),
				DBRoleToExecute:  tfDBRoleToExecuteObject(t, dbRole, dbRoleType),
				Networking:       types.ObjectNull(streamconnection.NetworkingObjectType.AttrTypes),
				AWS:              types.ObjectNull(streamconnection.AWSObjectType.AttrTypes),
				Headers:          types.MapNull(types.StringType),
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
				ProjectID:        types.StringValue(dummyProjectID),
				InstanceName:     types.StringValue(instanceName),
				ConnectionName:   types.StringValue(connectionName),
				Type:             types.StringValue("Kafka"),
				Authentication:   tfAuthenticationObject(t, authMechanism, authUsername, "raw password"), // password value is obtained from config, not api resp.
				BootstrapServers: types.StringValue(bootstrapServers),
				Config:           tfConfigMap(t, configMap),
				Security:         tfSecurityObject(t, DummyCACert, securityProtocol),
				DBRoleToExecute:  types.ObjectNull(streamconnection.DBRoleToExecuteObjectType.AttrTypes),
				Networking:       types.ObjectNull(streamconnection.NetworkingObjectType.AttrTypes),
				AWS:              types.ObjectNull(streamconnection.AWSObjectType.AttrTypes),
				Headers:          types.MapNull(types.StringType),
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
				ProjectID:       types.StringValue(dummyProjectID),
				InstanceName:    types.StringValue(instanceName),
				ConnectionName:  types.StringValue(connectionName),
				Type:            types.StringValue("Kafka"),
				Authentication:  types.ObjectNull(streamconnection.ConnectionAuthenticationObjectType.AttrTypes),
				Config:          types.MapNull(types.StringType),
				Security:        types.ObjectNull(streamconnection.ConnectionSecurityObjectType.AttrTypes),
				DBRoleToExecute: types.ObjectNull(streamconnection.DBRoleToExecuteObjectType.AttrTypes),
				Networking:      types.ObjectNull(streamconnection.NetworkingObjectType.AttrTypes),
				AWS:             types.ObjectNull(streamconnection.AWSObjectType.AttrTypes),
				Headers:         types.MapNull(types.StringType),
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
				ProjectID:        types.StringValue(dummyProjectID),
				InstanceName:     types.StringValue(instanceName),
				ConnectionName:   types.StringValue(connectionName),
				Type:             types.StringValue("Kafka"),
				Authentication:   tfAuthenticationObjectWithNoPassword(t, authMechanism, authUsername),
				BootstrapServers: types.StringValue(bootstrapServers),
				Config:           tfConfigMap(t, configMap),
				Security:         tfSecurityObject(t, DummyCACert, securityProtocol),
				DBRoleToExecute:  types.ObjectNull(streamconnection.DBRoleToExecuteObjectType.AttrTypes),
				Networking:       types.ObjectNull(streamconnection.NetworkingObjectType.AttrTypes),
				AWS:              types.ObjectNull(streamconnection.AWSObjectType.AttrTypes),
				Headers:          types.MapNull(types.StringType),
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
				ProjectID:       types.StringValue(dummyProjectID),
				InstanceName:    types.StringValue(instanceName),
				ConnectionName:  types.StringValue(sampleConnectionName),
				Type:            types.StringValue("Sample"),
				Authentication:  types.ObjectNull(streamconnection.ConnectionAuthenticationObjectType.AttrTypes),
				Config:          types.MapNull(types.StringType),
				Security:        types.ObjectNull(streamconnection.ConnectionSecurityObjectType.AttrTypes),
				DBRoleToExecute: types.ObjectNull(streamconnection.DBRoleToExecuteObjectType.AttrTypes),
				Networking:      types.ObjectNull(streamconnection.NetworkingObjectType.AttrTypes),
				AWS:             types.ObjectNull(streamconnection.AWSObjectType.AttrTypes),
				Headers:         types.MapNull(types.StringType),
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
				ProjectID:       types.StringValue(dummyProjectID),
				InstanceName:    types.StringValue(instanceName),
				ConnectionName:  types.StringValue(awslambdaConnectionName),
				Type:            types.StringValue("AWSLambda"),
				Authentication:  types.ObjectNull(streamconnection.ConnectionAuthenticationObjectType.AttrTypes),
				Config:          types.MapNull(types.StringType),
				Security:        types.ObjectNull(streamconnection.ConnectionSecurityObjectType.AttrTypes),
				DBRoleToExecute: types.ObjectNull(streamconnection.DBRoleToExecuteObjectType.AttrTypes),
				Networking:      types.ObjectNull(streamconnection.NetworkingObjectType.AttrTypes),
				AWS:             tfAWSLambdaConfigObject(t, sampleRoleArn),
				Headers:         types.MapNull(types.StringType),
			},
		},
		{
			name: "AWS S3 connection type with roleArn",
			SDKResp: &admin.StreamsConnection{
				Name: admin.PtrString(awss3ConnectionName),
				Type: admin.PtrString("S3"),
				Aws:  &admin.StreamsAWSConnectionConfig{RoleArn: admin.PtrString(sampleRoleArn), TestBucket: admin.PtrString(sampleTestBucket)},
			},
			providedProjID:       dummyProjectID,
			providedInstanceName: instanceName,
			expectedTFModel: &streamconnection.TFStreamConnectionModel{
				ProjectID:       types.StringValue(dummyProjectID),
				InstanceName:    types.StringValue(instanceName),
				ConnectionName:  types.StringValue(awss3ConnectionName),
				Type:            types.StringValue("S3"),
				Authentication:  types.ObjectNull(streamconnection.ConnectionAuthenticationObjectType.AttrTypes),
				Config:          types.MapNull(types.StringType),
				Security:        types.ObjectNull(streamconnection.ConnectionSecurityObjectType.AttrTypes),
				DBRoleToExecute: types.ObjectNull(streamconnection.DBRoleToExecuteObjectType.AttrTypes),
				Networking:      types.ObjectNull(streamconnection.NetworkingObjectType.AttrTypes),
				AWS:             tfAWSS3ConfigObject(t, sampleRoleArn, sampleTestBucket),
				Headers:         types.MapNull(types.StringType),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel, diags := streamconnection.NewTFStreamConnection(t.Context(), tc.providedProjID, tc.providedInstanceName, tc.providedAuthConfig, tc.SDKResp)
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
						Name: admin.PtrString(awss3ConnectionName),
						Type: admin.PtrString("S3"),
						Aws: &admin.StreamsAWSConnectionConfig{
							RoleArn:    admin.PtrString(sampleRoleArn),
							TestBucket: admin.PtrString(sampleTestBucket),
						},
					},
					{
						Name:    admin.PtrString(connectionName),
						Type:    admin.PtrString("Https"),
						Url:     admin.PtrString(httpsURL),
						Headers: &headersMap,
					},
				},
				TotalCount: admin.PtrInt(5),
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
				TotalCount:   types.Int64Value(5),
				Results: []streamconnection.TFStreamConnectionModel{
					{
						ID:               types.StringValue(fmt.Sprintf("%s-%s-%s", instanceName, dummyProjectID, connectionName)),
						ProjectID:        types.StringValue(dummyProjectID),
						InstanceName:     types.StringValue(instanceName),
						ConnectionName:   types.StringValue(connectionName),
						Type:             types.StringValue("Kafka"),
						Authentication:   tfAuthenticationObjectWithNoPassword(t, authMechanism, authUsername),
						BootstrapServers: types.StringValue(bootstrapServers),
						Config:           tfConfigMap(t, configMap),
						Security:         tfSecurityObject(t, DummyCACert, securityProtocol),
						DBRoleToExecute:  types.ObjectNull(streamconnection.DBRoleToExecuteObjectType.AttrTypes),
						Networking:       tfNetworkingObject(t, networkingType, nil),
						AWS:              types.ObjectNull(streamconnection.AWSObjectType.AttrTypes),
						Headers:          types.MapNull(types.StringType),
					},
					{
						ID:              types.StringValue(fmt.Sprintf("%s-%s-%s", instanceName, dummyProjectID, connectionName)),
						ProjectID:       types.StringValue(dummyProjectID),
						InstanceName:    types.StringValue(instanceName),
						ConnectionName:  types.StringValue(connectionName),
						Type:            types.StringValue("Cluster"),
						ClusterName:     types.StringValue(clusterName),
						Authentication:  types.ObjectNull(streamconnection.ConnectionAuthenticationObjectType.AttrTypes),
						Config:          types.MapNull(types.StringType),
						Security:        types.ObjectNull(streamconnection.ConnectionSecurityObjectType.AttrTypes),
						DBRoleToExecute: tfDBRoleToExecuteObject(t, dbRole, dbRoleType),
						Networking:      types.ObjectNull(streamconnection.NetworkingObjectType.AttrTypes),
						AWS:             types.ObjectNull(streamconnection.AWSObjectType.AttrTypes),
						Headers:         types.MapNull(types.StringType),
					},
					{
						ID:              types.StringValue(fmt.Sprintf("%s-%s-%s", instanceName, dummyProjectID, sampleConnectionName)),
						ProjectID:       types.StringValue(dummyProjectID),
						InstanceName:    types.StringValue(instanceName),
						ConnectionName:  types.StringValue(sampleConnectionName),
						Type:            types.StringValue("Sample"),
						ClusterName:     types.StringNull(),
						Authentication:  types.ObjectNull(streamconnection.ConnectionAuthenticationObjectType.AttrTypes),
						Config:          types.MapNull(types.StringType),
						Security:        types.ObjectNull(streamconnection.ConnectionSecurityObjectType.AttrTypes),
						DBRoleToExecute: types.ObjectNull(streamconnection.DBRoleToExecuteObjectType.AttrTypes),
						Networking:      types.ObjectNull(streamconnection.NetworkingObjectType.AttrTypes),
						AWS:             types.ObjectNull(streamconnection.AWSObjectType.AttrTypes),
						Headers:         types.MapNull(types.StringType),
					},
					{
						ID:              types.StringValue(fmt.Sprintf("%s-%s-%s", instanceName, dummyProjectID, awslambdaConnectionName)),
						ProjectID:       types.StringValue(dummyProjectID),
						InstanceName:    types.StringValue(instanceName),
						ConnectionName:  types.StringValue(awslambdaConnectionName),
						Type:            types.StringValue("AWSLambda"),
						ClusterName:     types.StringNull(),
						Authentication:  types.ObjectNull(streamconnection.ConnectionAuthenticationObjectType.AttrTypes),
						Config:          types.MapNull(types.StringType),
						Security:        types.ObjectNull(streamconnection.ConnectionSecurityObjectType.AttrTypes),
						DBRoleToExecute: types.ObjectNull(streamconnection.DBRoleToExecuteObjectType.AttrTypes),
						Networking:      types.ObjectNull(streamconnection.NetworkingObjectType.AttrTypes),
						AWS:             tfAWSLambdaConfigObject(t, sampleRoleArn),
						Headers:         types.MapNull(types.StringType),
					},
					{
						ID:              types.StringValue(fmt.Sprintf("%s-%s-%s", instanceName, dummyProjectID, awss3ConnectionName)),
						ProjectID:       types.StringValue(dummyProjectID),
						InstanceName:    types.StringValue(instanceName),
						ConnectionName:  types.StringValue(awss3ConnectionName),
						Type:            types.StringValue("S3"),
						ClusterName:     types.StringNull(),
						Authentication:  types.ObjectNull(streamconnection.ConnectionAuthenticationObjectType.AttrTypes),
						Config:          types.MapNull(types.StringType),
						Security:        types.ObjectNull(streamconnection.ConnectionSecurityObjectType.AttrTypes),
						DBRoleToExecute: types.ObjectNull(streamconnection.DBRoleToExecuteObjectType.AttrTypes),
						Networking:      types.ObjectNull(streamconnection.NetworkingObjectType.AttrTypes),
						AWS:             tfAWSS3ConfigObject(t, sampleRoleArn, sampleTestBucket),
						Headers:         types.MapNull(types.StringType),
					},
					{
						ID:              types.StringValue(fmt.Sprintf("%s-%s-%s", instanceName, dummyProjectID, connectionName)),
						ProjectID:       types.StringValue(dummyProjectID),
						InstanceName:    types.StringValue(instanceName),
						ConnectionName:  types.StringValue(connectionName),
						Type:            types.StringValue("Https"),
						ClusterName:     types.StringNull(),
						Authentication:  types.ObjectNull(streamconnection.ConnectionAuthenticationObjectType.AttrTypes),
						Config:          types.MapNull(types.StringType),
						Security:        types.ObjectNull(streamconnection.ConnectionSecurityObjectType.AttrTypes),
						DBRoleToExecute: types.ObjectNull(streamconnection.DBRoleToExecuteObjectType.AttrTypes),
						Networking:      types.ObjectNull(streamconnection.NetworkingObjectType.AttrTypes),
						AWS:             types.ObjectNull(streamconnection.AWSObjectType.AttrTypes),
						Headers:         tfConfigMap(t, headersMap),
						URL:             types.StringValue(httpsURL),
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
			name: "AWS S3 type TF state",
			tfModel: &streamconnection.TFStreamConnectionModel{
				ProjectID:      types.StringValue(dummyProjectID),
				InstanceName:   types.StringValue(instanceName),
				ConnectionName: types.StringValue(awslambdaConnectionName),
				Type:           types.StringValue("S3"),
				AWS:            tfAWSS3ConfigObject(t, sampleRoleArn, sampleTestBucket),
			},
			expectedSDKReq: &admin.StreamsConnection{
				Name: admin.PtrString(awslambdaConnectionName),
				Type: admin.PtrString("S3"),
				Aws: &admin.StreamsAWSConnectionConfig{
					RoleArn:    admin.PtrString(sampleRoleArn),
					TestBucket: admin.PtrString(sampleTestBucket),
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

func tfAWSS3ConfigObject(t *testing.T, roleArn, testBucket string) types.Object {
	t.Helper()
	aws, diags := types.ObjectValueFrom(t.Context(), streamconnection.AWSObjectType.AttrTypes, streamconnection.TFAWSModel{
		RoleArn:    types.StringValue(roleArn),
		TestBucket: types.StringValue(testBucket),
	})
	if diags.HasError() {
		t.Errorf("failed to create terraform data model: %s", diags.Errors()[0].Summary())
	}
	return aws
}
