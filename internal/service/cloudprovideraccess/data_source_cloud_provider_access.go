package cloudprovideraccess

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext:        dataSourceMongoDBAtlasCloudProviderAccessRead,
		DeprecationMessage: fmt.Sprintf(constant.DeprecationResourceByDateWithReplacement, "v1.14.0", "mongodbatlas_cloud_provider_access_setup"),
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"aws_iam_roles": {
				Type:     schema.TypeList,
				Elem:     dataSourceMongoDBAtlasCloudProviderAccess(),
				Computed: true,
			},
		},
	}
}

func dataSourceMongoDBAtlasCloudProviderAccess() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"atlas_aws_account_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"atlas_assumed_role_external_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"authorized_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"iam_assumed_role_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provider_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"feature_usages": {
				Type:     schema.TypeList,
				Elem:     featureUsagesSchema(),
				Computed: true,
			},
		},
	}
}

func featureUsagesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"feature_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"feature_id": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceMongoDBAtlasCloudProviderAccessRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)

	roles, _, err := conn.CloudProviderAccess.ListRoles(ctx, projectID)

	if err != nil {
		return diag.FromErr(fmt.Errorf(ErrorCloudProviderGetRead, err))
	}

	if err = d.Set("aws_iam_roles", flatCloudProviderAccessRolesAWS(roles)); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorCloudProviderGetRead, err))
	}

	d.SetId(id.UniqueId())

	return nil
}

func flatCloudProviderAccessRolesAWS(roles *matlas.CloudProviderAccessRoles) (list []map[string]any) {
	list = make([]map[string]any, 0, len(roles.AWSIAMRoles))
	for i := range roles.AWSIAMRoles {
		role := &(roles.AWSIAMRoles[i])
		list = append(list, roleToSchemaAWS(role))
	}

	return list
}

func roleToSchemaAWS(role *matlas.CloudProviderAccessRole) map[string]any {
	out := map[string]any{
		"atlas_aws_account_arn":          role.AtlasAWSAccountARN,
		"atlas_assumed_role_external_id": role.AtlasAssumedRoleExternalID,
		"authorized_date":                role.AuthorizedDate,
		"created_date":                   role.CreatedDate,
		"iam_assumed_role_arn":           role.IAMAssumedRoleARN,
		"provider_name":                  role.ProviderName,
		"role_id":                        role.RoleID,
	}

	features := make([]map[string]any, 0, len(role.FeatureUsages))

	for _, featureUsage := range role.FeatureUsages {
		features = append(features, featureToSchema(featureUsage))
	}

	out["feature_usages"] = features

	return out
}

func featureToSchema(feature *matlas.FeatureUsage) map[string]any {
	return map[string]any{
		"feature_type": feature.FeatureType,
		"feature_id":   feature.FeatureID,
	}
}

func featureUsagesSchemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"feature_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"feature_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
