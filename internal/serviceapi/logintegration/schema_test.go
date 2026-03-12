package logintegration_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccLogIntegration_discriminatorMissingRequired(t *testing.T) {
	testCases := []struct {
		name          string
		typeName      string
		extraAttrs    string
		expectedError string
	}{
		{
			name:          "datadog missing api_key",
			typeName:      "DATADOG_LOG_EXPORT",
			extraAttrs:    `region = "US1"`,
			expectedError: `Attribute "api_key" must be set when type is "DATADOG_LOG_EXPORT"`,
		},
		{
			name:     "splunk missing hec_token",
			typeName: "SPLUNK_LOG_EXPORT",
			extraAttrs: `
				hec_url = "https://splunk.example.com:8088"
			`,
			expectedError: `Attribute "hec_token" must be set when type is "SPLUNK_LOG_EXPORT"`,
		},
		{
			name:     "s3 missing bucket_name",
			typeName: "S3_LOG_EXPORT",
			extraAttrs: `
				iam_role_id = "111111111111111111111111"
				prefix_path = "prefix-path"
			`,
			expectedError: `Attribute "bucket_name" must be set when type is "S3_LOG_EXPORT"`,
		},
		{
			name:     "azure missing prefix_path",
			typeName: "AZURE_LOG_EXPORT",
			extraAttrs: `
				role_id                = "111111111111111111111111"
				storage_account_name   = "storageaccountname"
				storage_container_name = "storage-container-name"
			`,
			expectedError: `Attribute "prefix_path" must be set when type is "AZURE_LOG_EXPORT"`,
		},
		{
			name:     "gcs missing bucket_name",
			typeName: "GCS_LOG_EXPORT",
			extraAttrs: `
				prefix_path = "prefix-path"
				role_id     = "111111111111111111111111"
			`,
			expectedError: `Attribute "bucket_name" must be set when type is "GCS_LOG_EXPORT"`,
		},
		{
			name:     "otel missing otel_endpoint",
			typeName: "OTEL_LOG_EXPORT",
			extraAttrs: `
				otel_supplied_headers = [{
					name  = "Authorization"
					value = "Bearer token"
				}]
			`,
			expectedError: `Attribute "otel_endpoint" must be set when type is "OTEL_LOG_EXPORT"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Steps: []resource.TestStep{
					{
						Config:      config(tc.typeName, tc.extraAttrs),
						ExpectError: regexp.MustCompile(tc.expectedError),
					},
				},
			})
		})
	}
}

func TestAccLogIntegration_discriminatorDisallowedAttrs(t *testing.T) {
	testCases := []struct {
		name          string
		typeName      string
		extraAttrs    string
		expectedError string
	}{
		{
			name:     "datadog disallows bucket_name",
			typeName: "DATADOG_LOG_EXPORT",
			extraAttrs: `
				api_key     = "test-dd-api-key"
				region      = "US1"
				bucket_name = "test-bucket"
			`,
			expectedError: `"bucket_name" is not allowed when type is "DATADOG_LOG_EXPORT"`,
		},
		{
			name:     "s3 disallows hec_url",
			typeName: "S3_LOG_EXPORT",
			extraAttrs: `
				bucket_name = "test-bucket"
				iam_role_id = "111111111111111111111111"
				prefix_path = "prefix-path"
				hec_url     = "https://splunk.example.com:8088"
			`,
			expectedError: `"hec_url" is not allowed when type is "S3_LOG_EXPORT"`,
		},
		{
			name:     "splunk disallows otel_endpoint",
			typeName: "SPLUNK_LOG_EXPORT",
			extraAttrs: `
				hec_token     = "test-hec-token"
				hec_url       = "https://splunk.example.com:8088"
				otel_endpoint = "https://otel.example.com:4317/v1/logs"
			`,
			expectedError: `"otel_endpoint" is not allowed when type is "SPLUNK_LOG_EXPORT"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Steps: []resource.TestStep{
					{
						Config:      config(tc.typeName, tc.extraAttrs),
						ExpectError: regexp.MustCompile(tc.expectedError),
					},
				},
			})
		})
	}
}

func TestAccLogIntegration_discriminatorValidConfig(t *testing.T) {
	testCases := []struct {
		name   string
		config string
	}{
		{
			name: "valid datadog config",
			config: config("DATADOG_LOG_EXPORT", `
				api_key = "test-dd-api-key"
				region  = "US1"
			`),
		},
		{
			name: "valid splunk config",
			config: config("SPLUNK_LOG_EXPORT", `
				hec_token = "test-hec-token"
				hec_url   = "https://splunk.example.com:8088"
			`),
		},
		{
			name: "valid otel config",
			config: config("OTEL_LOG_EXPORT", `
				otel_endpoint = "https://otel.example.com:4317/v1/logs"
				otel_supplied_headers = [{
					name  = "Authorization"
					value = "Bearer token"
				}]
			`),
		},
		{
			name: "valid datadog config with unknown api_key at plan time",
			config: configWithPrefix(`
				resource "terraform_data" "unknown" {}
			`, "DATADOG_LOG_EXPORT", `
				api_key = terraform_data.unknown.id
				region  = "US1"
			`),
		},
		{
			name: "unknown integration type skips discriminator validation",
			config: config("FUTURE_UNKNOWN_LOG_EXPORT", `
				api_key = "test-dd-api-key"
			`),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Steps: []resource.TestStep{
					{
						Config:             tc.config,
						PlanOnly:           true,
						ExpectNonEmptyPlan: true,
					},
				},
			})
		})
	}
}

func config(typeName, extraAttrs string) string {
	return configWithPrefix("", typeName, extraAttrs)
}

func configWithPrefix(prefix, typeName, extraAttrs string) string {
	return fmt.Sprintf(`
		%[1]s
		resource "mongodbatlas_log_integration" "test" {
			project_id = "111111111111111111111111"
			type       = %[2]q
			log_types  = ["MONGOD"]
			%[3]s
		}
	`, prefix, typeName, extraAttrs)
}
