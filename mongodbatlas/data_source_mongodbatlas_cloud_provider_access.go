package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorGetRead = "error reading cloud provider access %s"
)

func dataSourceMongoDBAtlasCloudProviderAccessList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasCloudProviderAccessRead,
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

func dataSourceMongoDBAtlasCloudProviderAccessRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)

	roles, _, err := conn.CloudProviderAccess.ListRoles(context.Background(), projectID)

	if err != nil {
		return fmt.Errorf(errorGetRead, err)
	}

	if err = d.Set("aws_iam_roles", flatCloudProviderAccessRoles(roles)); err != nil {
		return fmt.Errorf(errorGetRead, err)
	}

	d.SetId(resource.UniqueId())

	return nil
}

func flatCloudProviderAccessRoles(roles *matlas.CloudProviderAccessRoles) (list []map[string]interface{}) {
	list = make([]map[string]interface{}, 0, len(roles.AWSIAMRoles))

	for i := range roles.AWSIAMRoles {
		role := &(roles.AWSIAMRoles[i])
		list = append(list, roleToSchema(role))
	}

	return list
}

func roleToSchema(role *matlas.AWSIAMRole) map[string]interface{} {
	out := map[string]interface{}{
		"atlas_aws_account_arn":          role.AtlasAWSAccountARN,
		"atlas_assumed_role_external_id": role.AtlasAssumedRoleExternalID,
		"authorized_date":                role.AuthorizedDate,
		"created_date":                   role.CreatedDate,
		"iam_assumed_role_arn":           role.IAMAssumedRoleARN,
		"provider_name":                  role.ProviderName,
		"role_id":                        role.RoleID,
	}

	features := make([]map[string]interface{}, 0, len(role.FeatureUsages))

	for _, featureUsage := range role.FeatureUsages {
		features = append(features, featureToSchema(featureUsage))
	}

	out["feature_usages"] = features

	return out
}

func featureToSchema(feature *matlas.FeatureUsage) map[string]interface{} {
	return map[string]interface{}{
		"feature_type": feature.FeatureType,
		"feature_id":   feature.FeatureID,
	}
}
