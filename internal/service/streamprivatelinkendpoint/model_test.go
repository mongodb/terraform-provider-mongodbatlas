package streamprivatelinkendpoint_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streamprivatelinkendpoint"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20250312007/admin"
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
	providerAccountID     = "providerAccountId"
	region                = "us-east-1"
	serviceEndpointID     = "serviceEndpointId"
	serviceEndpointIDS3   = "com.amazonaws.us-east-1.s3"
	state                 = "DONE"
	vendorConfluent       = "CONFLUENT"
	vendorS3              = "S3"
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
				Provider:              constant.AWS,
				ProviderAccountId:     &providerAccountID,
				Region:                &region,
				ServiceEndpointId:     &serviceEndpointID,
				State:                 &state,
				Vendor:                &vendorConfluent,
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
				Provider:              types.StringValue(constant.AWS),
				ProviderAccountId:     types.StringValue(providerAccountID),
				Region:                types.StringValue(region),
				ServiceEndpointId:     types.StringValue(serviceEndpointID),
				State:                 types.StringValue(state),
				Vendor:                types.StringValue(vendorConfluent),
			},
		},
		"SDK response without dns subdomains": {
			SDKResp: &admin.StreamsPrivateLinkConnection{
				Id:                  &id,
				Arn:                 &arn,
				DnsDomain:           &dnsDomain,
				InterfaceEndpointId: &interfaceEndpointID,
				Provider:            constant.AWS,
				Region:              &region,
				ServiceEndpointId:   &serviceEndpointID,
				State:               &state,
				Vendor:              &vendorConfluent,
			},
			projectID: projectID,
			expectedTFModel: &streamprivatelinkendpoint.TFModel{
				Id:                  types.StringValue(id),
				Arn:                 types.StringValue(arn),
				DnsDomain:           types.StringValue(dnsDomain),
				DnsSubDomain:        types.ListNull(types.StringType),
				ProjectId:           types.StringValue(projectID),
				InterfaceEndpointId: types.StringValue(interfaceEndpointID),
				Provider:            types.StringValue(constant.AWS),
				Region:              types.StringValue(region),
				ServiceEndpointId:   types.StringValue(serviceEndpointID),
				State:               types.StringValue(state),
				Vendor:              types.StringValue(vendorConfluent),
			},
		},
		"SDK response without arn": {
			SDKResp: &admin.StreamsPrivateLinkConnection{
				Id:                  &id,
				DnsDomain:           &dnsDomain,
				DnsSubDomain:        conversion.Pointer([]string{dnsSubDomain, dnsSubDomain}),
				InterfaceEndpointId: &interfaceEndpointID,
				Provider:            constant.AWS,
				Region:              &region,
				ServiceEndpointId:   &serviceEndpointID,
				State:               &state,
				Vendor:              &vendorConfluent,
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
				Provider:            types.StringValue(constant.AWS),
				Region:              types.StringValue(region),
				ServiceEndpointId:   types.StringValue(serviceEndpointID),
				State:               types.StringValue(state),
				Vendor:              types.StringValue(vendorConfluent),
			},
		},
		"SDK response with vendor S3": {
			SDKResp: &admin.StreamsPrivateLinkConnection{
				Id:       &id,
				Provider: constant.AWS,
				Vendor:   &vendorS3,
				Region:   &region,
			},
			projectID: projectID,
			expectedTFModel: &streamprivatelinkendpoint.TFModel{
				Id:           types.StringValue(id),
				Provider:     types.StringValue(constant.AWS),
				Vendor:       types.StringValue(vendorS3),
				Region:       types.StringValue(region),
				ProjectId:    types.StringValue(projectID),
				DnsSubDomain: types.ListNull(types.StringType),
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
				Provider:              types.StringValue(constant.AWS),
				ProviderAccountId:     types.StringValue(providerAccountID),
				Region:                types.StringValue(region),
				ServiceEndpointId:     types.StringValue(serviceEndpointID),
				State:                 types.StringValue(state),
				Vendor:                types.StringValue(vendorConfluent),
			},
			expectedSDKReq: &admin.StreamsPrivateLinkConnection{
				Arn:               &arn,
				DnsDomain:         &dnsDomain,
				DnsSubDomain:      conversion.Pointer([]string{dnsSubDomain, dnsSubDomain}),
				Provider:          constant.AWS,
				Region:            &region,
				ServiceEndpointId: &serviceEndpointID,
				State:             &state,
				Vendor:            &vendorConfluent,
			},
		},
		"TF state without dns subdomains": {
			tfModel: &streamprivatelinkendpoint.TFModel{
				Id:                  types.StringValue(id),
				Arn:                 types.StringValue(arn),
				DnsDomain:           types.StringValue(dnsDomain),
				ProjectId:           types.StringValue(projectID),
				InterfaceEndpointId: types.StringValue(interfaceEndpointID),
				Provider:            types.StringValue(constant.AWS),
				Region:              types.StringValue(region),
				ServiceEndpointId:   types.StringValue(serviceEndpointID),
				State:               types.StringValue(state),
				Vendor:              types.StringValue(vendorConfluent),
			},
			expectedSDKReq: &admin.StreamsPrivateLinkConnection{
				Arn:               &arn,
				DnsDomain:         &dnsDomain,
				Provider:          constant.AWS,
				Region:            &region,
				ServiceEndpointId: &serviceEndpointID,
				State:             &state,
				Vendor:            &vendorConfluent,
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
				Provider:              types.StringValue(constant.AWS),
				ProviderAccountId:     types.StringValue(providerAccountID),
				Region:                types.StringValue(region),
				ServiceEndpointId:     types.StringValue(serviceEndpointID),
				State:                 types.StringValue(state),
				Vendor:                types.StringValue(vendorConfluent),
			},
			expectedSDKReq: &admin.StreamsPrivateLinkConnection{
				DnsDomain:         &dnsDomain,
				DnsSubDomain:      conversion.Pointer([]string{dnsSubDomain, dnsSubDomain}),
				Provider:          constant.AWS,
				Region:            &region,
				ServiceEndpointId: &serviceEndpointID,
				State:             &state,
				Vendor:            &vendorConfluent,
			},
		},
		"TF state with s3 vendor": {
			tfModel: &streamprivatelinkendpoint.TFModel{
				Id:                types.StringValue(id),
				Provider:          types.StringValue(constant.AWS),
				Vendor:            types.StringValue(vendorS3),
				Region:            types.StringValue(region),
				ServiceEndpointId: types.StringValue(serviceEndpointIDS3),
			},
			expectedSDKReq: &admin.StreamsPrivateLinkConnection{
				Provider:          constant.AWS,
				Vendor:            &vendorS3,
				Region:            &region,
				ServiceEndpointId: &serviceEndpointIDS3,
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
