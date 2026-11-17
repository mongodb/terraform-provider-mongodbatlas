package thirdpartyintegration_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/thirdpartyintegration"
)

func TestPluralDataSourceCoreConfigSchema_ResultsSensitive(t *testing.T) {
	// results is computed-only, so its inner credential sensitivity does not reach Terraform.
	results := thirdpartyintegration.PluralDataSource().CoreConfigSchema().Attributes["results"]
	require.True(t, results.Sensitive, "results must be sensitive so Terraform redacts the integration credentials")
}
