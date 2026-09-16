package metricintegration_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/hcl"
)

const (
	resourceName         = "mongodbatlas_metric_integration.test"
	dataSourceName       = "data.mongodbatlas_metric_integration.test"
	pluralDataSourceName = "data.mongodbatlas_metric_integrations.test"
	datasourcesConfig    = `
		data "mongodbatlas_metric_integration" "test" {
			project_id            = mongodbatlas_metric_integration.test.project_id
			metric_integration_id = mongodbatlas_metric_integration.test.metric_integration_id
		}

		data "mongodbatlas_metric_integrations" "test" {
			project_id = mongodbatlas_metric_integration.test.project_id
			depends_on = [mongodbatlas_metric_integration.test]
		}
	`
	// Dummy endpoints for OAuth integrations.
	oauthTokenEndpoint = "https://192.0.2.2/oauth2/token" //nolint:gosec // Test data
	oauthClientID      = "atlas-otel-test"                //nolint:gosec // Test data
	oauthEndpoint      = "https://192.0.2.1/v1/metrics"
	// Pre-rendered token_request_params HCL fragments, passed to the config builders.
	trParamsClientSecret  = `token_request_params = { resource = "atlas:otel:test" }`
	trParamsPrivateKeyJWT = `token_request_params = { resource = "JPMC:URI:OTel" }`
)

// TestAccMetricIntegration_basic covers the base HEADER auth path. Serial because the project
// allows at most 2 metric integrations, and the OAuth tests share the same project.
func TestAccMetricIntegration_basic(t *testing.T) {
	var (
		projectID       = acc.ProjectIDExecution(t)
		integrationType = "OTEL"
		providerType    = "CUSTOM"
		aggregation     = "DELTA"
		endpoint        = os.Getenv("MONGODB_ATLAS_METRIC_INTEGRATION_ENDPOINT")
		headerValue     = os.Getenv("MONGODB_ATLAS_METRIC_INTEGRATION_API_KEY")
		metricSelection = []string{"ATLAS_STREAM_PROCESSING"}
		extraHeader     = true
		withDS          = true
	)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t); preCheckMetricIntegration(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, integrationType, providerType, aggregation, endpoint, headerValue, metricSelection, !extraHeader, withDS),
				Check:  checkBasic(integrationType, providerType, aggregation, endpoint, metricSelection, !extraHeader, withDS),
			},
			{
				Config: configBasic(projectID, integrationType, providerType, aggregation, endpoint, headerValue, metricSelection, extraHeader, !withDS),
				Check:  checkBasic(integrationType, providerType, aggregation, endpoint, metricSelection, extraHeader, !withDS),
			},
			{
				Config:                               configBasic(projectID, integrationType, providerType, aggregation, endpoint, headerValue, metricSelection, extraHeader, false),
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    importStateIDFunc(resourceName),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "metric_integration_id",
				ImportStateVerifyIgnore:              []string{"headers"}, // headers is write-only and not returned on GET (only headers_redacted is)
			},
		},
	})
}

func preCheckMetricIntegration(tb testing.TB) {
	tb.Helper()
	if os.Getenv("MONGODB_ATLAS_METRIC_INTEGRATION_ENDPOINT") == "" || os.Getenv("MONGODB_ATLAS_METRIC_INTEGRATION_API_KEY") == "" {
		tb.Fatal("`MONGODB_ATLAS_METRIC_INTEGRATION_ENDPOINT` and `MONGODB_ATLAS_METRIC_INTEGRATION_API_KEY` must be set for acceptance testing")
	}
}

func configBasic(projectID, integrationType, providerType, aggregation, endpoint, headerValue string, metricSelection []string, extraHeader, withDS bool) string {
	selectionHCL := hcl.StringSliceToHCL(metricSelection)
	extraHeaderHCL := ""
	if extraHeader {
		extraHeaderHCL = `,
				{
					name  = "x-custom-header"
					value = "custom-value"
				}`
	}
	dsConfig := ""
	if withDS {
		dsConfig = datasourcesConfig
	}
	return fmt.Sprintf(`
		resource "mongodbatlas_metric_integration" "test" {
			project_id              = %[1]q
			integration_type        = %[2]q
			provider_type           = %[3]q
			auth_type               = "HEADER"
			aggregation_temporality = %[4]q
			endpoint                = %[5]q
			metric_selection        = %[6]s

			headers = [
				{
					name  = "dd-api-key"
					value = %[7]q
				}%[8]s
			]
		}

		%[9]s
	`, projectID, integrationType, providerType, aggregation, endpoint, selectionHCL, headerValue, extraHeaderHCL, dsConfig)
}

func checkBasic(integrationType, providerType, aggregation, endpoint string, metricSelection []string, extraHeader, withDS bool) resource.TestCheckFunc {
	headerCount := "1"
	if extraHeader {
		headerCount = "2"
	}
	setChecks := []string{"project_id", "metric_integration_id"}
	mapChecks := map[string]string{
		"integration_type":        integrationType,
		"provider_type":           providerType,
		"auth_type":               "HEADER",
		"aggregation_temporality": aggregation,
		"endpoint":                endpoint,
		"metric_selection.#":      strconv.Itoa(len(metricSelection)),
		"headers_redacted.#":      headerCount,
	}
	// headers is a write-only request field, so it exists only on the resource, not the data sources.
	checks := []resource.TestCheckFunc{resource.TestCheckResourceAttr(resourceName, "headers.#", headerCount)}
	var dsName *string
	if withDS {
		dsName = new(dataSourceName)
		checks = append(checks, resource.TestCheckResourceAttrWith(pluralDataSourceName, "results.#", acc.IntGreatThan(0)))
	}
	checks = append(checks, acc.CheckRSAndDS(resourceName, dsName, nil, setChecks, mapChecks, checkExists(resourceName)))
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		projectID := rs.Primary.Attributes["project_id"]
		integrationID := rs.Primary.Attributes["metric_integration_id"]
		if projectID == "" || integrationID == "" {
			return fmt.Errorf("checkExists, attributes not found for: %s", resourceName)
		}

		resp, err := acc.GetMetricIntegration(context.Background(), projectID, integrationID)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			return nil
		}
		return fmt.Errorf("metric integration for project_id %s with id %s does not exist, status %d", projectID, integrationID, resp.StatusCode)
	}
}

func checkDestroy(state *terraform.State) error {
	for name, rs := range state.RootModule().Resources {
		if name != resourceName {
			continue
		}
		projectID := rs.Primary.Attributes["project_id"]
		integrationID := rs.Primary.Attributes["metric_integration_id"]
		if projectID == "" || integrationID == "" {
			return fmt.Errorf("checkDestroy, attributes not found for: %s", resourceName)
		}
		resp, err := acc.GetMetricIntegration(context.Background(), projectID, integrationID)
		if err != nil {
			return err
		}
		resp.Body.Close()
		if resp.StatusCode == http.StatusNotFound {
			return nil
		}
		if resp.StatusCode == http.StatusOK {
			return fmt.Errorf("metric integration for project_id %s with id %s still exists", projectID, integrationID)
		}
		return fmt.Errorf("checkDestroy, unexpected status %d for project_id %s with id %s", resp.StatusCode, projectID, integrationID)
	}
	return nil
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		projectID := rs.Primary.Attributes["project_id"]
		integrationID := rs.Primary.Attributes["metric_integration_id"]
		if projectID == "" || integrationID == "" {
			return "", fmt.Errorf("import, attributes not found for: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", projectID, integrationID), nil
	}
}

// configOauth renders an OAUTH2 metric integration resource. The optional oauth block
// attributes are appended after client_id, each on its own line at the block indentation.
func configOauth(projectID, endpoint, tokenEndpoint, clientID, clientAuthMethod string, withDS bool, oauthAttrs ...string) string {
	dsConfig := ""
	if withDS {
		dsConfig = datasourcesConfig
	}
	var attrsBlock strings.Builder
	for _, a := range oauthAttrs {
		if a != "" {
			attrsBlock.WriteString("\n\t\t\t\t" + a)
		}
	}
	return fmt.Sprintf(`
		resource "mongodbatlas_metric_integration" "test" {
			project_id              = %[1]q
			integration_type        = "OTEL"
			provider_type           = "CUSTOM"
			auth_type               = "OAUTH2"
			aggregation_temporality = "DELTA"
			endpoint                = %[2]q
			metric_selection        = ["ATLAS_STREAM_PROCESSING"]

			oauth = {
				client_auth_method = %[3]q
				token_endpoint     = %[4]q
				client_id          = %[5]q%[6]s
			}
		}

		%[7]s
	`, projectID, endpoint, clientAuthMethod, tokenEndpoint, clientID, attrsBlock.String(), dsConfig)
}

func configOauthClientSecret(projectID, endpoint, tokenEndpoint, clientID, clientSecret string, scopes []string, tokenRequestParams string, withDS bool) string {
	return configOauth(projectID, endpoint, tokenEndpoint, clientID, "CLIENT_SECRET", withDS,
		fmt.Sprintf("client_secret = %q", clientSecret),
		"scopes = "+hcl.StringSliceToHCL(scopes),
		tokenRequestParams,
	)
}

func configOauthPrivateKeyJWT(projectID, endpoint, tokenEndpoint, clientID string, scopes []string, tokenRequestParams string, withDS bool) string {
	oauthAttrs := []string{}
	if len(scopes) > 0 {
		oauthAttrs = append(oauthAttrs, "scopes = "+hcl.StringSliceToHCL(scopes))
	}
	oauthAttrs = append(oauthAttrs, tokenRequestParams)
	return configOauth(projectID, endpoint, tokenEndpoint, clientID, "PRIVATE_KEY_JWT", withDS, oauthAttrs...)
}

// TestAccMetricIntegration_oauthClientSecret covers the CLIENT_SECRET OAuth path
func TestAccMetricIntegration_oauthClientSecret(t *testing.T) {
	// TODO: remove once the OAuth metric integration fields are available in prod and before merging to master
	acc.SkipTestForCI(t)
	projectID := acc.ProjectIDExecution(t)
	var (
		secret1 = "client-secret-initial"
		secret2 = "client-secret-rotated"
		scopes1 = []string{"metrics.write"}
		scopes2 = []string{"metrics.write", "monitoring.read"}
		dsName  = new(dataSourceName)
	)

	// Test is serial because the project allows at most 2 metric integrations and tests share it.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configOauthClientSecret(projectID, oauthEndpoint, oauthTokenEndpoint, oauthClientID, secret1, scopes1, "", true),
				Check:  checkOauthClientSecret(secret1, scopes1, 0, dsName),
			},
			{
				Config: configOauthClientSecret(projectID, oauthEndpoint, oauthTokenEndpoint, oauthClientID, secret2, scopes2, trParamsClientSecret, false),
				Check:  checkOauthClientSecret(secret2, scopes2, 1, nil),
			},
			{
				Config:                               configOauthClientSecret(projectID, oauthEndpoint, oauthTokenEndpoint, oauthClientID, secret2, scopes2, trParamsClientSecret, false),
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    importStateIDFunc(resourceName),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "metric_integration_id",
				// client_secret is write-only and not returned on GET, so it cannot match on import.
				ImportStateVerifyIgnore: []string{"oauth.client_secret"},
			},
		},
	})
}

// TestAccMetricIntegration_oauthPrivateKeyJWT covers the PRIVATE_KEY_JWT OAuth path
func TestAccMetricIntegration_oauthPrivateKeyJWT(t *testing.T) {
	// TODO: remove once the OAuth metric integration fields are available in prod and before merging to master
	acc.SkipTestForCI(t)
	projectID := acc.ProjectIDExecution(t)
	var (
		scopes1 = []string{}
		scopes2 = []string{"metrics.write"}
		dsName  = new(dataSourceName)
	)

	// Test is serial because the project allows at most 2 metric integrations and tests share it.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configOauthPrivateKeyJWT(projectID, oauthEndpoint, oauthTokenEndpoint, oauthClientID, scopes1, trParamsPrivateKeyJWT, true),
				Check:  checkOauthPrivateKeyJWT(scopes1, 1, dsName),
			},
			{
				Config: configOauthPrivateKeyJWT(projectID, oauthEndpoint, oauthTokenEndpoint, oauthClientID, scopes2, trParamsPrivateKeyJWT, false),
				Check:  checkOauthPrivateKeyJWT(scopes2, 1, nil),
			},
			{
				Config:                               configOauthPrivateKeyJWT(projectID, oauthEndpoint, oauthTokenEndpoint, oauthClientID, scopes2, trParamsPrivateKeyJWT, false),
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    importStateIDFunc(resourceName),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "metric_integration_id",
			},
		},
	})
}

// TestAccMetricIntegration_oauthPrivateKeyJWTRejectsClientSecret verifies the API rejects a
// client_secret set on a PRIVATE_KEY_JWT integration. The validation is server-side, so this
// step performs a real apply that is expected to fail with the API's 400.
func TestAccMetricIntegration_oauthPrivateKeyJWTRejectsClientSecret(t *testing.T) {
	// TODO: remove once the OAuth metric integration fields are available in prod and before merging to master
	acc.SkipTestForCI(t)
	projectID := acc.ProjectIDExecution(t)

	// Test is serial because the project allows at most 2 metric integrations and tests share it.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configOauthPrivateKeyJWTWithClientSecret(projectID),
				ExpectError: regexp.MustCompile(
					`oauth.clientSecret must not be set when clientAuthMethod is 'PRIVATE_KEY_JWT'`,
				),
			},
		},
	})
}

func configOauthPrivateKeyJWTWithClientSecret(projectID string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_metric_integration" "test" {
			project_id              = %[1]q
			integration_type        = "OTEL"
			provider_type           = "CUSTOM"
			auth_type               = "OAUTH2"
			aggregation_temporality = "DELTA"
			endpoint                = %[2]q
			metric_selection        = ["ATLAS_STREAM_PROCESSING"]

			oauth = {
				client_auth_method = "PRIVATE_KEY_JWT"
				token_endpoint     = %[3]q
				client_id          = %[4]q
				client_secret      = "should-be-rejected"
			}
		}
	`, projectID, oauthEndpoint, oauthTokenEndpoint, oauthClientID)
}

func checkOauthClientSecret(clientSecret string, scopes []string, tokenRequestParamsCount int, dsName *string) resource.TestCheckFunc {
	mapChecks := map[string]string{
		"auth_type":                    "OAUTH2",
		"oauth.client_auth_method":     "CLIENT_SECRET",
		"oauth.client_id":              oauthClientID,
		"oauth.token_endpoint":         oauthTokenEndpoint,
		"oauth.scopes.#":               strconv.Itoa(len(scopes)),
		"oauth.token_request_params.%": strconv.Itoa(tokenRequestParamsCount),
		"headers_redacted.#":           "0",
	}
	setChecks := []string{"project_id", "metric_integration_id"}
	// client_secret is write-only and only present on the resource (never on data sources).
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(resourceName, "oauth.client_secret", clientSecret),
	}
	if dsName != nil {
		checks = append(checks, resource.TestCheckResourceAttrWith(pluralDataSourceName, "results.#", acc.IntGreatThan(0)))
	}
	checks = append(checks, acc.CheckRSAndDS(resourceName, dsName, nil, setChecks, mapChecks, checkExists(resourceName)))
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func checkOauthPrivateKeyJWT(scopes []string, tokenRequestParamsCount int, dsName *string) resource.TestCheckFunc {
	mapChecks := map[string]string{
		"auth_type":                        "OAUTH2",
		"oauth.client_auth_method":         "PRIVATE_KEY_JWT",
		"oauth.client_id":                  oauthClientID,
		"oauth.token_endpoint":             oauthTokenEndpoint,
		"oauth.scopes.#":                   strconv.Itoa(len(scopes)),
		"oauth.token_request_params.%":     strconv.Itoa(tokenRequestParamsCount),
		"oauth.signing_key_info.algorithm": "RS256",
		"headers_redacted.#":               "0",
	}
	// signing_key_info is read-only and present on both resource and data sources for PRIVATE_KEY_JWT.
	setChecks := []string{
		"project_id", "metric_integration_id",
		"oauth.signing_key_info.kid",
		"oauth.signing_key_info.jwks_uri",
		"oauth.signing_key_info.created_at",
	}
	checks := []resource.TestCheckFunc{
		// No client_secret is sent for PRIVATE_KEY_JWT, so it must not appear in state.
		resource.TestCheckNoResourceAttr(resourceName, "oauth.client_secret"),
	}
	if dsName != nil {
		checks = append(checks, resource.TestCheckResourceAttrWith(pluralDataSourceName, "results.#", acc.IntGreatThan(0)))
	}
	checks = append(checks, acc.CheckRSAndDS(resourceName, dsName, nil, setChecks, mapChecks, checkExists(resourceName)))
	return resource.ComposeAggregateTestCheckFunc(checks...)
}
