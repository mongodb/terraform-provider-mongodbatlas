package streamconnection_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streamconnection"
	"github.com/stretchr/testify/require"
)

func TestResourceSchemaAuthenticationAWSRoleARNOptional(t *testing.T) {
	resourceSchema := streamconnection.ResourceSchema(t.Context())
	authentication, ok := resourceSchema.Attributes["authentication"].(schema.SingleNestedAttribute)
	require.True(t, ok)

	aws, ok := authentication.Attributes["aws"].(schema.SingleNestedAttribute)
	require.True(t, ok)

	roleARN, ok := aws.Attributes["role_arn"].(schema.StringAttribute)
	require.True(t, ok)
	require.True(t, roleARN.Optional, "authentication.aws.role_arn must not be required when authentication.aws is omitted")
	require.False(t, roleARN.Required)
}
