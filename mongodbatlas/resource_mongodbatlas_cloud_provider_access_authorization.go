package mongodbatlas

import (
	"context"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func resourceMongoDBAtlasCloudProviderAccessAuthorization() *schema.Resource {
	return &schema.Resource{
		Read: resourceMongoDBAtlasCloudProviderAccessAuthorizationRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"aws": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"iam_assumed_role_arn": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"feature_usages": {
				Type:     schema.TypeList,
				Elem:     featureUsagesSchema(),
				Computed: true,
			},
		},
	}
}

func resourceMongoDBAtlasCloudProviderAccessAuthorizationRead(d *schema.ResourceData, meta interface{}) error {
	// sadly there is no just get API
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())

	roleID := ids["id"] // atlas ID
	projectID := ids["project_id"]

	roles, _, err := conn.CloudProviderAccess.ListRoles(context.Background(), projectID)

	if err != nil {
		return fmt.Errorf(errorGetRead, err)
	}

	// for future implementations is aws?
	var targetRole matlas.AWSIAMRole

	sort.Slice(roles.AWSIAMRoles,
		func(i, j int) bool { return roles.AWSIAMRoles[i].RoleID < roles.AWSIAMRoles[j].RoleID })

	index := sort.Search(len(roles.AWSIAMRoles), func(i int) bool { return roles.AWSIAMRoles[i].RoleID >= roleID })

	if index < len(roles.AWSIAMRoles) && roles.AWSIAMRoles[index].RoleID == roleID {
		targetRole = roles.AWSIAMRoles[index]
		roleSchema := roleToSchemaAuthorization(&targetRole)

		for key, val := range roleSchema {
			if err := d.Set(key, val); err != nil {
				return fmt.Errorf(errorGetRead, err)
			}
		}
	}

	// Not Found when more providers added this must be an interface
	if targetRole.RoleID == "" && !d.IsNewResource() {
		d.SetId("")
		return nil
	}

	return nil
}

func roleToSchemaAuthorization(role *matlas.AWSIAMRole) map[string]interface{} {
	out := map[string]interface{}{
		"role_id": role.RoleID,
		"aws": map[string]interface{}{
			"iam_assumed_role_arn": role.IAMAssumedRoleARN,
		},
		"authorized_date": role.AuthorizedDate,
	}

	// features
	features := make([]map[string]interface{}, 0, len(role.FeatureUsages))

	for _, featureUsage := range role.FeatureUsages {
		features = append(features, featureToSchema(featureUsage))
	}

	out["feature_usages"] = features

	return out
}
