package mongodbatlas

import (
	"context"
	"fmt"

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
				Type:     schema.TypeString,
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

	if err = flatCloudProviderAccessRoles(roles, d); err != nil {
		return fmt.Errorf(errorGetRead, err)
	}

	return nil
}

func flatCloudProviderAccessRoles(roles *matlas.CloudProviderAccessRoles, d *schema.ResourceData) error {
	return nil
}
