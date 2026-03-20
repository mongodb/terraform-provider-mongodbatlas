package serviceaccountjwt_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/serviceaccountjwt"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var versionChecks = []tfversion.TerraformVersionCheck{
	tfversion.SkipBelow(tfversion.Version1_10_0),
}

func newES(erd *config.EphemeralResourceData) *serviceaccountjwt.ES {
	return &serviceaccountjwt.ES{
		ESCommon: config.ESCommon{
			EphemeralResourceData: erd,
			ResourceName:          serviceaccountjwt.ResourceTypeName,
		},
	}
}

func TestResolveCredentials_Order(t *testing.T) {
	r := newES(&config.EphemeralResourceData{
		ClientID:         "provider-id",
		ClientSecret:     "provider-secret",
		BaseURL:          "https://provider.example.com/",
		TerraformVersion: "1.10.0",
	})

	t.Run("resource attributes first", func(t *testing.T) {
		model := serviceaccountjwt.TFModel{
			ClientID:     types.StringValue("resource-id"),
			ClientSecret: types.StringValue("resource-secret"),
		}
		id, secret, baseURL, diags := r.ResolveCredentials(&model)
		require.False(t, diags.HasError())
		require.Equal(t, "resource-id", id)
		require.Equal(t, "resource-secret", secret)
		require.Equal(t, "https://provider.example.com/", baseURL)
	})

	t.Run("provider fallback when no resource attributes", func(t *testing.T) {
		model := serviceaccountjwt.TFModel{}
		id, secret, baseURL, diags := r.ResolveCredentials(&model)
		require.False(t, diags.HasError())
		require.Equal(t, "provider-id", id)
		require.Equal(t, "provider-secret", secret)
		require.Equal(t, "https://provider.example.com/", baseURL)
	})
}

func TestResolveCredentials_NonSAProviderAuth(t *testing.T) {
	t.Run("nil provider data", func(t *testing.T) {
		r := newES(nil)
		model := serviceaccountjwt.TFModel{}
		clientID, clientSecret, baseURL, diags := r.ResolveCredentials(&model)
		require.Empty(t, clientID)
		require.Empty(t, clientSecret)
		require.Empty(t, baseURL)
		require.True(t, diags.HasError())
		require.Contains(t, diags.Errors()[0].Detail(), "Service Account credentials")
	})

	t.Run("provider configured with PAK (no SA credentials)", func(t *testing.T) {
		r := newES(&config.EphemeralResourceData{
			BaseURL:          "https://cloud.mongodb.com/",
			TerraformVersion: "1.11.0",
		})
		model := serviceaccountjwt.TFModel{}
		clientID, clientSecret, baseURL, diags := r.ResolveCredentials(&model)
		require.Empty(t, clientID)
		require.Empty(t, clientSecret)
		require.Empty(t, baseURL)
		require.True(t, diags.HasError())
		require.Contains(t, diags.Errors()[0].Detail(), "different authentication method")
	})
}

func TestResolveCredentials_DoesNotMixSources(t *testing.T) {
	t.Run("partial resource credentials do not fallback", func(t *testing.T) {
		r := newES(&config.EphemeralResourceData{
			ClientID:     "provider-id",
			ClientSecret: "provider-secret",
			BaseURL:      "https://provider.example.com/",
		})
		model := serviceaccountjwt.TFModel{ClientID: types.StringValue("resource-id")}
		clientID, clientSecret, baseURL, diags := r.ResolveCredentials(&model)
		require.Empty(t, clientID)
		require.Empty(t, clientSecret)
		require.Empty(t, baseURL)
		require.True(t, diags.HasError())
	})

	t.Run("partial provider credentials do not fallback", func(t *testing.T) {
		r := newES(&config.EphemeralResourceData{
			ClientID: "provider-id",
			BaseURL:  "https://provider.example.com/",
		})
		model := serviceaccountjwt.TFModel{}
		clientID, clientSecret, baseURL, diags := r.ResolveCredentials(&model)
		require.Empty(t, clientID)
		require.Empty(t, clientSecret)
		require.Empty(t, baseURL)
		require.True(t, diags.HasError())
	})
}

func TestAccServiceAccountJWT_providerCredentials(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckSA(t) },
		TerraformVersionChecks:   versionChecks,
		ProtoV6ProviderFactories: acc.ProtoV6FactoriesWithEcho(),
		Steps: []resource.TestStep{
			{
				Config: configProviderCredentials(false),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.jwt", tfjsonpath.New("data").AtMapKey("token_type"), knownvalue.StringExact("Bearer")),
					statecheck.ExpectKnownValue("echo.jwt", tfjsonpath.New("data").AtMapKey("expires_in"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.jwt", tfjsonpath.New("data").AtMapKey("access_token"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccServiceAccountJWT_revokeOnClosure(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckSA(t) },
		TerraformVersionChecks:   versionChecks,
		ProtoV6ProviderFactories: acc.ProtoV6FactoriesWithEcho(),
		Steps: []resource.TestStep{
			{
				Config: configProviderCredentials(true),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.jwt", tfjsonpath.New("data").AtMapKey("token_type"), knownvalue.StringExact("Bearer")),
					statecheck.ExpectKnownValue("echo.jwt", tfjsonpath.New("data").AtMapKey("access_token"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccServiceAccountJWT_explicitCredentials(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckSA(t) },
		TerraformVersionChecks:   versionChecks,
		ProtoV6ProviderFactories: acc.ProtoV6FactoriesWithEcho(),
		Steps: []resource.TestStep{
			{
				Config: configExplicitCredentials(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.jwt", tfjsonpath.New("data").AtMapKey("token_type"), knownvalue.StringExact("Bearer")),
					statecheck.ExpectKnownValue("echo.jwt", tfjsonpath.New("data").AtMapKey("access_token"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccServiceAccountJWT_invalidCredentials(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		TerraformVersionChecks:   versionChecks,
		ProtoV6ProviderFactories: acc.ProtoV6FactoriesWithEcho(),
		Steps: []resource.TestStep{
			{
				Config:      configInvalidCredentials(),
				ExpectError: regexp.MustCompile(`Error generating Service Account JWT`),
			},
		},
	})
}

func TestAccServiceAccountJWT_partialResourceCredentials(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		TerraformVersionChecks:   versionChecks,
		ProtoV6ProviderFactories: acc.ProtoV6FactoriesWithEcho(),
		Steps: []resource.TestStep{
			{
				Config:      configPartialCredentials(),
				ExpectError: regexp.MustCompile(`(?s)both client_id and\s+client_secret must be provided`),
			},
		},
	})
}

func configProviderCredentials(revokeOnClosure bool) string {
	return fmt.Sprintf(`
ephemeral "mongodbatlas_service_account_jwt" "test" {
  revoke_on_closure = %t
}

provider "echo" {
  data = ephemeral.mongodbatlas_service_account_jwt.test
}

resource "echo" "jwt" {}
`, revokeOnClosure)
}

func configExplicitCredentials() string {
	return fmt.Sprintf(`
ephemeral "mongodbatlas_service_account_jwt" "test" {
  client_id     = %q
  client_secret = %q
}

provider "echo" {
  data = ephemeral.mongodbatlas_service_account_jwt.test
}

resource "echo" "jwt" {}
`, os.Getenv("MONGODB_ATLAS_CLIENT_ID"), os.Getenv("MONGODB_ATLAS_CLIENT_SECRET"))
}

func configInvalidCredentials() string {
	return `
ephemeral "mongodbatlas_service_account_jwt" "test" {
  client_id     = "mdb_sa_id_000000000000000000000000"
  client_secret = "not-a-real-secret"
}

provider "echo" {
  data = ephemeral.mongodbatlas_service_account_jwt.test
}

resource "echo" "jwt" {}
`
}

func configPartialCredentials() string {
	return `
ephemeral "mongodbatlas_service_account_jwt" "test" {
  client_id = "mdb_sa_id_000000000000000000000000"
}

provider "echo" {
  data = ephemeral.mongodbatlas_service_account_jwt.test
}

resource "echo" "jwt" {}
`
}
