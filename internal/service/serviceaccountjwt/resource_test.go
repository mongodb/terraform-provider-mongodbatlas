package serviceaccountjwt_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var versionChecks = []tfversion.TerraformVersionCheck{
	tfversion.SkipBelow(tfversion.Version1_10_0),
}

func TestAccServiceAccountJWT_providerCredentials(t *testing.T) {
	acc.SkipIfNotSA(t)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
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
	acc.SkipIfNotSA(t)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
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
	acc.SkipIfNotSA(t)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
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
				ExpectError: regexp.MustCompile(config.ErrPartialCreds),
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
