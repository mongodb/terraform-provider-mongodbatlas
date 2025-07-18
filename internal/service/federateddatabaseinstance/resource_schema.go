package federateddatabaseinstance

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func cloudProviderConfig(isDataSource bool) *schema.Schema {
	var computed, optional, required bool
	var maxItems int
	if isDataSource {
		computed = true
		maxItems = 0
	} else {
		required = true
		optional = true
		maxItems = 1
	}

	return &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: maxItems,
		Computed: true,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"aws": {
					Type:     schema.TypeList,
					MaxItems: maxItems,
					Optional: true,
					Computed: computed,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"role_id": {
								Type:     schema.TypeString,
								Required: required,
								Computed: computed,
							},
							"test_s3_bucket": {
								Type:     schema.TypeString,
								Required: required,
								Optional: isDataSource,
							},
							"iam_assumed_role_arn": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"iam_user_arn": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"external_id": {
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},
				"azure": {
					Type:     schema.TypeList,
					MaxItems: maxItems,
					Optional: optional,
					Computed: computed,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"role_id": {
								Type:     schema.TypeString,
								Required: required,
								Computed: computed,
							},
							"atlas_app_id": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"service_principal_id": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"tenant_id": {
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}
