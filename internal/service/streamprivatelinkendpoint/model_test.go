package streamprivatelinkendpoint_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streamprivatelinkendpoint"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20250312008/admin"
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
				ServiceAttachmentUris: types.ListNull(types.StringType),
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
				Id:                    types.StringValue(id),
				Arn:                   types.StringValue(arn),
				DnsDomain:             types.StringValue(dnsDomain),
				DnsSubDomain:          types.ListNull(types.StringType),
				ProjectId:             types.StringValue(projectID),
				InterfaceEndpointId:   types.StringValue(interfaceEndpointID),
				Provider:              types.StringValue(constant.AWS),
				Region:                types.StringValue(region),
				ServiceEndpointId:     types.StringValue(serviceEndpointID),
				ServiceAttachmentUris: types.ListNull(types.StringType),
				State:                 types.StringValue(state),
				Vendor:                types.StringValue(vendorConfluent),
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
				ProjectId:             types.StringValue(projectID),
				InterfaceEndpointId:   types.StringValue(interfaceEndpointID),
				Provider:              types.StringValue(constant.AWS),
				Region:                types.StringValue(region),
				ServiceEndpointId:     types.StringValue(serviceEndpointID),
				ServiceAttachmentUris: types.ListNull(types.StringType),
				State:                 types.StringValue(state),
				Vendor:                types.StringValue(vendorConfluent),
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
				Id:                    types.StringValue(id),
				Provider:              types.StringValue(constant.AWS),
				Vendor:                types.StringValue(vendorS3),
				Region:                types.StringValue(region),
				ProjectId:             types.StringValue(projectID),
				DnsSubDomain:          types.ListNull(types.StringType),
				ServiceAttachmentUris: types.ListNull(types.StringType),
			},
		},
		"SDK response with GCP Confluent": {
			SDKResp: &admin.StreamsPrivateLinkConnection{
				Id:                       &id,
				DnsDomain:                &dnsDomain,
				DnsSubDomain:             conversion.Pointer([]string{dnsSubDomain, dnsSubDomain}),
				GcpServiceAttachmentUris: conversion.Pointer([]string{"projects/test-project/regions/us-west1/serviceAttachments/test-attachment-1", "projects/test-project/regions/us-west1/serviceAttachments/test-attachment-2"}),
				Provider:                 constant.GCP,
				Region:                   &region,
				State:                    &state,
				Vendor:                   &vendorConfluent,
			},
			projectID: projectID,
			expectedTFModel: &streamprivatelinkendpoint.TFModel{
				Id:        types.StringValue(id),
				DnsDomain: types.StringValue(dnsDomain),
				DnsSubDomain: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue(dnsSubDomain),
					types.StringValue(dnsSubDomain),
				}),
				ServiceAttachmentUris: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("projects/test-project/regions/us-west1/serviceAttachments/test-attachment-1"),
					types.StringValue("projects/test-project/regions/us-west1/serviceAttachments/test-attachment-2"),
				}),
				ProjectId: types.StringValue(projectID),
				Provider:  types.StringValue(constant.GCP),
				Region:    types.StringValue(region),
				State:     types.StringValue(state),
				Vendor:    types.StringValue(vendorConfluent),
			},
		},
		"Empty SDK response": {
			SDKResp: &admin.StreamsPrivateLinkConnection{},
			expectedTFModel: &streamprivatelinkendpoint.TFModel{
				ProjectId:             types.StringValue(""),
				Provider:              types.StringValue(""),
				DnsSubDomain:          types.ListNull(types.StringType),
				ServiceAttachmentUris: types.ListNull(types.StringType),
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
				ServiceAttachmentUris: types.ListNull(types.StringType),
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
				Id:                    types.StringValue(id),
				Arn:                   types.StringValue(arn),
				DnsDomain:             types.StringValue(dnsDomain),
				ProjectId:             types.StringValue(projectID),
				InterfaceEndpointId:   types.StringValue(interfaceEndpointID),
				Provider:              types.StringValue(constant.AWS),
				Region:                types.StringValue(region),
				ServiceEndpointId:     types.StringValue(serviceEndpointID),
				ServiceAttachmentUris: types.ListNull(types.StringType),
				State:                 types.StringValue(state),
				Vendor:                types.StringValue(vendorConfluent),
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
				ServiceAttachmentUris: types.ListNull(types.StringType),
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
				Id:                    types.StringValue(id),
				Provider:              types.StringValue(constant.AWS),
				Vendor:                types.StringValue(vendorS3),
				Region:                types.StringValue(region),
				ServiceEndpointId:     types.StringValue(serviceEndpointIDS3),
				ServiceAttachmentUris: types.ListNull(types.StringType),
			},
			expectedSDKReq: &admin.StreamsPrivateLinkConnection{
				Provider:          constant.AWS,
				Vendor:            &vendorS3,
				Region:            &region,
				ServiceEndpointId: &serviceEndpointIDS3,
			},
		},
		"TF state with GCP Confluent": {
			tfModel: &streamprivatelinkendpoint.TFModel{
				Id:        types.StringValue(id),
				DnsDomain: types.StringValue(dnsDomain),
				DnsSubDomain: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue(dnsSubDomain),
					types.StringValue(dnsSubDomain),
				}),
				ServiceAttachmentUris: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("projects/test-project/regions/us-west1/serviceAttachments/test-attachment-1"),
					types.StringValue("projects/test-project/regions/us-west1/serviceAttachments/test-attachment-2"),
				}),
				Provider: types.StringValue(constant.GCP),
				Region:   types.StringValue(region),
				State:    types.StringValue(state),
				Vendor:   types.StringValue(vendorConfluent),
			},
			expectedSDKReq: &admin.StreamsPrivateLinkConnection{
				DnsDomain:                &dnsDomain,
				DnsSubDomain:             conversion.Pointer([]string{dnsSubDomain, dnsSubDomain}),
				GcpServiceAttachmentUris: conversion.Pointer([]string{"projects/test-project/regions/us-west1/serviceAttachments/test-attachment-1", "projects/test-project/regions/us-west1/serviceAttachments/test-attachment-2"}),
				Provider:                 constant.GCP,
				Region:                   &region,
				State:                    &state,
				Vendor:                   &vendorConfluent,
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

func TestStreamPrivatelinkEndpointValidation(t *testing.T) {
	testCases := map[string]struct {
		tfModel     *streamprivatelinkendpoint.TFModel
		expectError bool
		errorCount  int
	}{
		"GCP Confluent with all required fields": {
			tfModel: &streamprivatelinkendpoint.TFModel{
				Provider:              types.StringValue(constant.GCP),
				Vendor:                types.StringValue(streamprivatelinkendpoint.VendorConfluent),
				ServiceAttachmentUris: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("projects/test-project/regions/us-west1/serviceAttachments/test-attachment")}),
				ServiceEndpointId:     types.StringNull(),
				DnsDomain:             types.StringValue("example.com"),
				Region:                types.StringValue("us-west1"),
			},
			expectError: false,
			errorCount:  0,
		},
		"GCP Confluent missing service_attachment_uris": {
			tfModel: &streamprivatelinkendpoint.TFModel{
				Provider:          types.StringValue(constant.GCP),
				Vendor:            types.StringValue(streamprivatelinkendpoint.VendorConfluent),
				ServiceEndpointId: types.StringNull(),
				DnsDomain:         types.StringValue("example.com"),
				Region:            types.StringValue("us-west1"),
			},
			expectError: true,
			errorCount:  1,
		},
		"GCP Confluent missing dns_domain": {
			tfModel: &streamprivatelinkendpoint.TFModel{
				Provider:              types.StringValue(constant.GCP),
				Vendor:                types.StringValue(streamprivatelinkendpoint.VendorConfluent),
				ServiceAttachmentUris: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("projects/test-project/regions/us-west1/serviceAttachments/test-attachment")}),
				ServiceEndpointId:     types.StringNull(),
				DnsDomain:             types.StringNull(),
				Region:                types.StringValue("us-west1"),
			},
			expectError: true,
			errorCount:  1,
		},
		"GCP Confluent missing region": {
			tfModel: &streamprivatelinkendpoint.TFModel{
				Provider:              types.StringValue(constant.GCP),
				Vendor:                types.StringValue(streamprivatelinkendpoint.VendorConfluent),
				ServiceAttachmentUris: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("projects/test-project/regions/us-west1/serviceAttachments/test-attachment")}),
				ServiceEndpointId:     types.StringNull(),
				DnsDomain:             types.StringValue("example.com"),
				Region:                types.StringNull(),
			},
			expectError: true,
			errorCount:  1,
		},
		"AWS Confluent": {
			tfModel: &streamprivatelinkendpoint.TFModel{
				Provider:          types.StringValue(constant.AWS),
				Vendor:            types.StringValue(streamprivatelinkendpoint.VendorConfluent),
				ServiceEndpointId: types.StringValue("vpce-12345"),
				DnsDomain:         types.StringValue("example.com"),
				Region:            types.StringValue("us-west-2"),
			},
			expectError: false,
			errorCount:  0,
		},
		"Both service_endpoint_id and service_attachment_uris provided": {
			tfModel: &streamprivatelinkendpoint.TFModel{
				Provider:              types.StringValue(constant.GCP),
				Vendor:                types.StringValue(streamprivatelinkendpoint.VendorConfluent),
				ServiceAttachmentUris: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("projects/test-project/regions/us-west1/serviceAttachments/test-attachment")}),
				ServiceEndpointId:     types.StringValue("vpce-12345"),
				DnsDomain:             types.StringValue("example.com"),
				Region:                types.StringValue("us-west1"),
			},
			expectError: true,
			errorCount:  1,
		},
		"Neither service_endpoint_id nor service_attachment_uris provided": {
			tfModel: &streamprivatelinkendpoint.TFModel{
				Provider:          types.StringValue(constant.GCP),
				Vendor:            types.StringValue(streamprivatelinkendpoint.VendorConfluent),
				ServiceEndpointId: types.StringNull(),
				DnsDomain:         types.StringValue("example.com"),
				Region:            types.StringValue("us-west1"),
			},
			expectError: true,
			errorCount:  1,
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			_, diags := streamprivatelinkendpoint.NewAtlasReq(t.Context(), tc.tfModel)

			if tc.expectError {
				if !diags.HasError() {
					t.Errorf("Expected validation errors but got none")
				}
				if len(diags.Errors()) != tc.errorCount {
					t.Errorf("Expected %d validation errors but got %d", tc.errorCount, len(diags.Errors()))
				}
			} else if diags.HasError() {
				t.Errorf("Expected no validation errors but got: %v", diags.Errors())
			}
		})
	}
}
