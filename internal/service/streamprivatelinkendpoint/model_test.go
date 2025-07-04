package streamprivatelinkendpoint_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streamprivatelinkendpoint"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20250312005/admin"
)

type sdkToTFModelTestCase struct {
	SDKResp         *admin.StreamsPrivateLinkConnection
	expectedTFModel *streamprivatelinkendpoint.TFModel
	projectID       string
}

var (
	projectID             = "projectID"
	id                    = "id"
	arn                   = "arn"
	dnsDomain             = "dnsDomain"
	dnsSubDomain          = "dnsSubDomain"
	interfaceEndpointID   = "interfaceEndpointId"
	interfaceEndpointName = "interfaceEndpointName"
	provider              = "AWS"
	providerAccountID     = "providerAccountId"
	region                = "us-east-1"
	serviceEndpointID     = "serviceEndpointId"
	state                 = "DONE"
	vendor                = "CONFLUENT"
)

func TestStreamPrivatelinkEndpointSDKToTFModel(t *testing.T) {
	testCases := map[string]sdkToTFModelTestCase{
		"Complete SDK response": {
			SDKResp: &admin.StreamsPrivateLinkConnection{
				Id:                    &id,
				Arn:                   &arn,
				DnsDomain:             &dnsDomain,
				DnsSubDomain:          conversion.Pointer([]string{dnsSubDomain, dnsSubDomain}),
				InterfaceEndpointId:   &interfaceEndpointID,
				InterfaceEndpointName: &interfaceEndpointName,
				Provider:              provider,
				ProviderAccountId:     &providerAccountID,
				Region:                &region,
				ServiceEndpointId:     &serviceEndpointID,
				State:                 &state,
				Vendor:                &vendor,
			},
			projectID: projectID,
			expectedTFModel: &streamprivatelinkendpoint.TFModel{
				Id:        types.StringValue(id),
				Arn:       types.StringValue(arn),
				DnsDomain: types.StringValue(dnsDomain),
				DnsSubDomain: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue(dnsSubDomain),
					types.StringValue(dnsSubDomain),
				}),
				ProjectId:             types.StringValue(projectID),
				InterfaceEndpointId:   types.StringValue(interfaceEndpointID),
				InterfaceEndpointName: types.StringValue(interfaceEndpointName),
				Provider:              types.StringValue(provider),
				ProviderAccountId:     types.StringValue(providerAccountID),
				Region:                types.StringValue(region),
				ServiceEndpointId:     types.StringValue(serviceEndpointID),
				State:                 types.StringValue(state),
				Vendor:                types.StringValue(vendor),
			},
		},
		"SDK response without dns subdomains": {
			SDKResp: &admin.StreamsPrivateLinkConnection{
				Id:                  &id,
				Arn:                 &arn,
				DnsDomain:           &dnsDomain,
				InterfaceEndpointId: &interfaceEndpointID,
				Provider:            provider,
				Region:              &region,
				ServiceEndpointId:   &serviceEndpointID,
				State:               &state,
				Vendor:              &vendor,
			},
			projectID: projectID,
			expectedTFModel: &streamprivatelinkendpoint.TFModel{
				Id:                  types.StringValue(id),
				Arn:                 types.StringValue(arn),
				DnsDomain:           types.StringValue(dnsDomain),
				DnsSubDomain:        types.ListNull(types.StringType),
				ProjectId:           types.StringValue(projectID),
				InterfaceEndpointId: types.StringValue(interfaceEndpointID),
				Provider:            types.StringValue(provider),
				Region:              types.StringValue(region),
				ServiceEndpointId:   types.StringValue(serviceEndpointID),
				State:               types.StringValue(state),
				Vendor:              types.StringValue(vendor),
			},
		},
		"SDK response without arn": {
			SDKResp: &admin.StreamsPrivateLinkConnection{
				Id:                  &id,
				DnsDomain:           &dnsDomain,
				DnsSubDomain:        conversion.Pointer([]string{dnsSubDomain, dnsSubDomain}),
				InterfaceEndpointId: &interfaceEndpointID,
				Provider:            provider,
				Region:              &region,
				ServiceEndpointId:   &serviceEndpointID,
				State:               &state,
				Vendor:              &vendor,
			},
			projectID: projectID,
			expectedTFModel: &streamprivatelinkendpoint.TFModel{
				Id:        types.StringValue(id),
				DnsDomain: types.StringValue(dnsDomain),
				DnsSubDomain: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue(dnsSubDomain),
					types.StringValue(dnsSubDomain),
				}),
				ProjectId:           types.StringValue(projectID),
				InterfaceEndpointId: types.StringValue(interfaceEndpointID),
				Provider:            types.StringValue(provider),
				Region:              types.StringValue(region),
				ServiceEndpointId:   types.StringValue(serviceEndpointID),
				State:               types.StringValue(state),
				Vendor:              types.StringValue(vendor),
			},
		},
		"Empty SDK response": {
			SDKResp: &admin.StreamsPrivateLinkConnection{},
			expectedTFModel: &streamprivatelinkendpoint.TFModel{
				ProjectId:    types.StringValue(""),
				Provider:     types.StringValue(""),
				DnsSubDomain: types.ListNull(types.StringType),
			},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			resultModel, diags := streamprivatelinkendpoint.NewTFModel(t.Context(), tc.projectID, tc.SDKResp)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedTFModel, resultModel, "created terraform model did not match expected output")
		})
	}
}

type tfToSDKModelTestCase struct {
	tfModel        *streamprivatelinkendpoint.TFModel
	expectedSDKReq *admin.StreamsPrivateLinkConnection
}

func TestStreamPrivatelinkEndpointTFModelToSDK(t *testing.T) {
	testCases := map[string]tfToSDKModelTestCase{
		"Complete TF state": {
			tfModel: &streamprivatelinkendpoint.TFModel{
				Id:        types.StringValue(id),
				Arn:       types.StringValue(arn),
				DnsDomain: types.StringValue(dnsDomain),
				DnsSubDomain: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue(dnsSubDomain),
					types.StringValue(dnsSubDomain),
				}),
				ProjectId:             types.StringValue(projectID),
				InterfaceEndpointId:   types.StringValue(interfaceEndpointID),
				InterfaceEndpointName: types.StringValue(interfaceEndpointName),
				Provider:              types.StringValue(provider),
				ProviderAccountId:     types.StringValue(providerAccountID),
				Region:                types.StringValue(region),
				ServiceEndpointId:     types.StringValue(serviceEndpointID),
				State:                 types.StringValue(state),
				Vendor:                types.StringValue(vendor),
			},
			expectedSDKReq: &admin.StreamsPrivateLinkConnection{
				Arn:               &arn,
				DnsDomain:         &dnsDomain,
				DnsSubDomain:      conversion.Pointer([]string{dnsSubDomain, dnsSubDomain}),
				Provider:          provider,
				Region:            &region,
				ServiceEndpointId: &serviceEndpointID,
				State:             &state,
				Vendor:            &vendor,
			},
		},
		"TF state without dns subdomains": {
			tfModel: &streamprivatelinkendpoint.TFModel{
				Id:                  types.StringValue(id),
				Arn:                 types.StringValue(arn),
				DnsDomain:           types.StringValue(dnsDomain),
				ProjectId:           types.StringValue(projectID),
				InterfaceEndpointId: types.StringValue(interfaceEndpointID),
				Provider:            types.StringValue(provider),
				Region:              types.StringValue(region),
				ServiceEndpointId:   types.StringValue(serviceEndpointID),
				State:               types.StringValue(state),
				Vendor:              types.StringValue(vendor),
			},
			expectedSDKReq: &admin.StreamsPrivateLinkConnection{
				Arn:               &arn,
				DnsDomain:         &dnsDomain,
				Provider:          provider,
				Region:            &region,
				ServiceEndpointId: &serviceEndpointID,
				State:             &state,
				Vendor:            &vendor,
			},
		},
		"TF state without arn": {
			tfModel: &streamprivatelinkendpoint.TFModel{
				Id:        types.StringValue(id),
				DnsDomain: types.StringValue(dnsDomain),
				DnsSubDomain: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue(dnsSubDomain),
					types.StringValue(dnsSubDomain),
				}),
				ProjectId:             types.StringValue(projectID),
				InterfaceEndpointId:   types.StringValue(interfaceEndpointID),
				InterfaceEndpointName: types.StringValue(interfaceEndpointName),
				Provider:              types.StringValue(provider),
				ProviderAccountId:     types.StringValue(providerAccountID),
				Region:                types.StringValue(region),
				ServiceEndpointId:     types.StringValue(serviceEndpointID),
				State:                 types.StringValue(state),
				Vendor:                types.StringValue(vendor),
			},
			expectedSDKReq: &admin.StreamsPrivateLinkConnection{
				DnsDomain:         &dnsDomain,
				DnsSubDomain:      conversion.Pointer([]string{dnsSubDomain, dnsSubDomain}),
				Provider:          provider,
				Region:            &region,
				ServiceEndpointId: &serviceEndpointID,
				State:             &state,
				Vendor:            &vendor,
			},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			apiReqResult, diags := streamprivatelinkendpoint.NewAtlasReq(t.Context(), tc.tfModel)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedSDKReq, apiReqResult, "created sdk model did not match expected output")
		})
	}
}
