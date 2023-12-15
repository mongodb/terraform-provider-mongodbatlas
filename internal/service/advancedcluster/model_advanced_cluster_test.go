package advancedcluster_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/stretchr/testify/assert"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

var (
	standard     = "standard"
	standardSrv  = "standardSrv"
	private      = "private"
	privateSrv   = "privateSrv"
	region       = "region"
	providerName = "providerName"
	endpointID   = "endpointID"
	endpoint     = matlas.Endpoint{
		Region:       region,
		ProviderName: providerName,
		EndpointID:   endpointID,
	}
	privateEndpoint = matlas.PrivateEndpoint{
		ConnectionString:                  connectionString,
		SRVConnectionString:               srvConnectionString,
		SRVShardOptimizedConnectionString: srvShardOptimizedConnectionString,
		Type:                              endpointType,
		Endpoints:                         []matlas.Endpoint{endpoint},
	}
	connectionString                  = "connectionString"
	srvConnectionString               = "srvConnectionString"
	srvShardOptimizedConnectionString = "SrvShardOptimizedConnectionString"
	endpointType                      = "EndpointType"
	tfEndpointModel                   = []advancedcluster.TfEndpointModel{
		{
			Region:       types.StringValue(region),
			ProviderName: types.StringValue(providerName),
			EndpointID:   types.StringValue(endpointID),
		},
	}
	endpointList, _        = types.ListValueFrom(context.Background(), advancedcluster.TfEndpointType, tfEndpointModel)
	tfPrivateEndpointModel = []advancedcluster.TfPrivateEndpointModel{
		{
			ConnectionString:                  types.StringValue(connectionString),
			SrvConnectionString:               types.StringValue(srvConnectionString),
			SrvShardOptimizedConnectionString: types.StringValue(srvShardOptimizedConnectionString),
			EndpointType:                      types.StringValue(endpointType),
			Endpoints:                         endpointList,
		},
	}
	privateEndpointList, _ = types.ListValueFrom(context.Background(), advancedcluster.TfPrivateEndpointType, tfPrivateEndpointModel)
	diskIOPS               = int64(10)
	ebsVolumeType          = "volumeType"
	instanceSize           = "instanceType"
	nodeCount              = int(1)
)

func TestNewTFEndpointModel(t *testing.T) {
	expectedEmpty, _ := types.ListValueFrom(context.Background(), advancedcluster.TfEndpointType, make([]advancedcluster.TfEndpointModel, 0))
	expectedWithValues := endpointList
	testCases := []struct {
		name           string
		expectedResult types.List
		endpoints      []matlas.Endpoint
	}{
		{
			name:           "Endpoint is empty",
			endpoints:      []matlas.Endpoint{},
			expectedResult: expectedEmpty,
		},
		{
			name: "Endpoint has values ",
			endpoints: []matlas.Endpoint{
				endpoint,
			},
			expectedResult: expectedWithValues,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := advancedcluster.NewTFEndpointModel(context.Background(), tc.endpoints)

			assert.Equal(t, tc.expectedResult, resultModel)
		})
	}
}

func TestNewTfPrivateEndpointModel(t *testing.T) {
	expectedEmpty, _ := types.ListValueFrom(context.Background(), advancedcluster.TfPrivateEndpointType, make([]advancedcluster.TfPrivateEndpointModel, 0))

	expectedWithValues := privateEndpointList
	testCases := []struct {
		name            string
		expectedResult  types.List
		privateEndpoint []matlas.PrivateEndpoint
	}{
		{
			name:            "PrivateEndpoint is empty",
			privateEndpoint: []matlas.PrivateEndpoint{},
			expectedResult:  expectedEmpty,
		},
		{
			name: "PrivateEndpoint has values ",
			privateEndpoint: []matlas.PrivateEndpoint{
				{
					ConnectionString:                  connectionString,
					SRVConnectionString:               srvConnectionString,
					SRVShardOptimizedConnectionString: srvShardOptimizedConnectionString,
					Type:                              endpointType,
					Endpoints:                         []matlas.Endpoint{endpoint},
				},
			},
			expectedResult: expectedWithValues,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := advancedcluster.NewTfPrivateEndpointModel(context.Background(), tc.privateEndpoint)

			assert.Equal(t, tc.expectedResult, resultModel)
		})
	}
}

func TestNewTfConnectionStringsModel(t *testing.T) {
	expectedWithValues := privateEndpointList
	testCases := []struct {
		name           string
		connString     *matlas.ConnectionStrings
		expectedResult []*advancedcluster.TfConnectionStringModel
	}{
		{
			name:           "ConnectionString is nil",
			connString:     nil,
			expectedResult: []*advancedcluster.TfConnectionStringModel{},
		},
		{
			name: "ConnectionString not nil",
			connString: &matlas.ConnectionStrings{
				Standard:        standard,
				StandardSrv:     standardSrv,
				Private:         private,
				PrivateSrv:      privateSrv,
				PrivateEndpoint: []matlas.PrivateEndpoint{privateEndpoint},
			},
			expectedResult: []*advancedcluster.TfConnectionStringModel{
				{
					Standard:        types.StringValue(standard),
					StandardSrv:     types.StringValue(standardSrv),
					Private:         types.StringValue(private),
					PrivateSrv:      types.StringValue(privateSrv),
					PrivateEndpoint: expectedWithValues,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := advancedcluster.NewTfConnectionStringsModel(context.Background(), tc.connString)

			assert.Equal(t, tc.expectedResult, resultModel)
		})
	}
}

func TestNewTfRegionsConfigSpecsModel(t *testing.T) {
	testCases := []struct {
		name           string
		specs          *matlas.Specs
		expectedResult []*advancedcluster.TfRegionsConfigSpecsModel
	}{
		{
			name:           "Specs is nil",
			specs:          nil,
			expectedResult: []*advancedcluster.TfRegionsConfigSpecsModel{},
		},
		{
			name: "Specs not nil",
			specs: &matlas.Specs{
				DiskIOPS:      &diskIOPS,
				EbsVolumeType: ebsVolumeType,
				InstanceSize:  instanceSize,
				NodeCount:     &nodeCount,
			},
			expectedResult: []*advancedcluster.TfRegionsConfigSpecsModel{
				{
					DiskIOPS:      types.Int64PointerValue(&diskIOPS),
					EBSVolumeType: types.StringValue(ebsVolumeType),
					InstanceSize:  types.StringValue(instanceSize),
					NodeCount:     types.Int64Value(int64(nodeCount)),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := advancedcluster.NewTfRegionsConfigSpecsModel(tc.specs)

			assert.Equal(t, tc.expectedResult, resultModel)
		})
	}
}

var (
	enabled          = true
	scaleDownEnabled = false
	minInstanceSize  = "minInstanceSize"
	maxInstanceSize  = "maxInstanceSize"
	readPreference   = "readPreference"
	key              = "key"
	value            = "value"
)

func TestNewTfRegionsConfigAutoScalingSpecsModel(t *testing.T) {
	testCases := []struct {
		name                string
		advancedAutoScaling *matlas.AdvancedAutoScaling
		expectedResult      []*advancedcluster.TfRegionsConfigAutoScalingSpecsModel
	}{
		{
			name:                "AdvancedAutoScaling is nil",
			advancedAutoScaling: nil,
			expectedResult:      []*advancedcluster.TfRegionsConfigAutoScalingSpecsModel{},
		},
		{
			name: "AdvancedAutoScaling is not nil, Cumpute is nil",
			advancedAutoScaling: &matlas.AdvancedAutoScaling{
				Compute: nil,
			},
			expectedResult: []*advancedcluster.TfRegionsConfigAutoScalingSpecsModel{},
		},
		{
			name: "AdvancedAutoScaling not nil, Compute not nil",
			advancedAutoScaling: &matlas.AdvancedAutoScaling{
				DiskGB: &matlas.DiskGB{
					Enabled: &enabled,
				},
				Compute: &matlas.Compute{
					Enabled:          &enabled,
					ScaleDownEnabled: &scaleDownEnabled,
					MinInstanceSize:  minInstanceSize,
					MaxInstanceSize:  maxInstanceSize,
				},
			},
			expectedResult: []*advancedcluster.TfRegionsConfigAutoScalingSpecsModel{
				{
					DiskGBEnabled:           types.BoolValue(enabled),
					ComputeEnabled:          types.BoolValue(enabled),
					ComputeScaleDownEnabled: types.BoolValue(scaleDownEnabled),
					ComputeMinInstanceSize:  types.StringValue(minInstanceSize),
					ComputeMaxInstanceSize:  types.StringValue(maxInstanceSize),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := advancedcluster.NewTfRegionsConfigAutoScalingSpecsModel(tc.advancedAutoScaling)

			assert.Equal(t, tc.expectedResult, resultModel)
		})
	}
}

func TestNewTfBiConnectorConfigModel(t *testing.T) {
	testCases := []struct {
		name           string
		biConnector    *matlas.BiConnector
		expectedResult []*advancedcluster.TfBiConnectorConfigModel
	}{
		{
			name:           "BiConnector is nil",
			biConnector:    nil,
			expectedResult: []*advancedcluster.TfBiConnectorConfigModel{},
		},
		{
			name: "BiConnector is not nil",
			biConnector: &matlas.BiConnector{
				Enabled:        &enabled,
				ReadPreference: readPreference,
			},
			expectedResult: []*advancedcluster.TfBiConnectorConfigModel{
				{
					Enabled:        types.BoolValue(enabled),
					ReadPreference: types.StringValue(readPreference),
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := advancedcluster.NewTfBiConnectorConfigModel(tc.biConnector)

			assert.Equal(t, tc.expectedResult, resultModel)
		})
	}
}

func TestNewTfTagsModel(t *testing.T) {
	testCases := []struct {
		name           string
		tags           *[]*matlas.Tag
		expectedResult []*advancedcluster.TfTagModel
	}{
		{
			name:           "Tags is empty",
			tags:           nil,
			expectedResult: []*advancedcluster.TfTagModel{},
		},
		{
			name:           "Tags is empty",
			tags:           &[]*matlas.Tag{},
			expectedResult: []*advancedcluster.TfTagModel{},
		},
		{
			name: "Tags is not nil",
			tags: &[]*matlas.Tag{
				{
					Key:   key,
					Value: value,
				},
			},
			expectedResult: []*advancedcluster.TfTagModel{
				{
					Key:   types.StringValue(key),
					Value: types.StringValue(value),
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := advancedcluster.NewTfTagsModel(tc.tags)

			assert.Equal(t, tc.expectedResult, resultModel)
		})
	}
}

func TestNewTfLabelsModel(t *testing.T) {
	testCases := []struct {
		name           string
		labels         []matlas.Label
		expectedResult []advancedcluster.TfLabelModel
	}{
		{
			name:           "Labels is nil",
			labels:         nil,
			expectedResult: []advancedcluster.TfLabelModel{},
		},
		{
			name: "Labels is not nil",
			labels: []matlas.Label{
				{
					Key:   key,
					Value: value,
				},
			},
			expectedResult: []advancedcluster.TfLabelModel{
				{
					Key:   types.StringValue(key),
					Value: types.StringValue(value),
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := advancedcluster.NewTfLabelsModel(tc.labels)

			assert.Equal(t, tc.expectedResult, resultModel)
		})
	}
}
