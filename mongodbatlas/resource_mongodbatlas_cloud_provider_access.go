package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorCloudProviderAccessCreate = "error creating cloud provider access %s"
)

func resourceMongoDBAtlasCloudProviderAccess() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoDBAtlasCloudProviderAccessCreate,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"provider_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"AWS"}, false),
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

func resourceMongoDBAtlasCloudProviderAccessCreate(d *schema.ResourceData, meta interface{}) error {
	projectID := d.Get("project_id").(string)

	conn := meta.(*matlas.Client)

	requestParameters := &matlas.CloudProviderAccessRoleRequest{
		ProviderName: d.Get("provider_name").(string),
	}

	role, _, err := conn.CloudProviderAccess.CreateRole(context.Background(), projectID, requestParameters)

	if err != nil {
		return fmt.Errorf(errorCloudProviderAccessCreate, err)
	}

	schema := roleToSchema(*role)

	for key, val := range schema {
		if err := d.Set(key, val); err != nil {
			return fmt.Errorf(errorCloudProviderAccessCreate, err)
		}
	}

	d.SetId(encodeStateID(map[string]string{
		"id":         role.RoleID,
		"project_id": projectID,
	}))

	return nil
}
