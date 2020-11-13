package mongodbatlas

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceMongoDBAtlasCloudProviderAccess() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
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
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
