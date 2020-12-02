package mongodbatlas

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorCloudProviderAccessCreate   = "error creating cloud provider access %s"
	errorCloudProviderAccessUpdate   = "error updating cloud provider access %s"
	errorCloudProviderAccessDelete   = "error deleting cloud provider access %s"
	errorCloudProviderAccessImporter = "error importing cloud provider access %s"
)

func resourceMongoDBAtlasCloudProviderAccess() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoDBAtlasCloudProviderAccessCreate,
		Read:   resourceMongoDBAtlasCloudProviderAccessRead,
		Update: resourceMongoDBAtlasCloudProviderAccessUpdate,
		Delete: resourceMongoDBAtlasCloudProviderAccessDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMongoDBAtlasCloudProviderAccessImportState,
		},
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
				Optional: true,
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

	roleSchema := roleToSchema(role)

	d.SetId(encodeStateID(map[string]string{
		"id":            role.RoleID,
		"project_id":    projectID,
		"provider_name": role.ProviderName,
	}))

	for key, val := range roleSchema {
		if err := d.Set(key, val); err != nil {
			return fmt.Errorf(errorCloudProviderAccessCreate, err)
		}
	}

	return nil
}

func resourceMongoDBAtlasCloudProviderAccessRead(d *schema.ResourceData, meta interface{}) error {
	// sadly there is no just get API
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]

	roles, _, err := conn.CloudProviderAccess.ListRoles(context.Background(), projectID)

	if err != nil {
		return fmt.Errorf(errorGetRead, err)
	}

	var targetRole matlas.AWSIAMRole

	// searching in roles
	for i := range roles.AWSIAMRoles {
		role := &(roles.AWSIAMRoles[i])

		if role.RoleID == ids["id"] && role.ProviderName == ids["provider_name"] {
			targetRole = *role
		}
	}

	// Not Found
	if targetRole.RoleID == "" && !d.IsNewResource() {
		d.SetId("")
		return nil
	}

	roleSchema := roleToSchema(&targetRole)

	for key, val := range roleSchema {
		if err := d.Set(key, val); err != nil {
			return fmt.Errorf(errorGetRead, err)
		}
	}

	return nil
}

func resourceMongoDBAtlasCloudProviderAccessUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())

	projectID := ids["project_id"]
	roleID := ids["id"]

	if d.HasChanges("provider_name", "iam_assumed_role_arn") {
		req := &matlas.CloudProviderAuthorizationRequest{
			ProviderName:      d.Get("provider_name").(string),
			IAMAssumedRoleARN: d.Get("iam_assumed_role_arn").(string),
		}

		role, _, err := conn.CloudProviderAccess.AuthorizeRole(context.Background(), projectID, roleID, req)
		if err != nil {
			return fmt.Errorf(errorCloudProviderAccessUpdate, err)
		}

		roleSchema := roleToSchema(role)

		for key, val := range roleSchema {
			if err := d.Set(key, val); err != nil {
				return fmt.Errorf(errorGetRead, err)
			}
		}
	}

	return nil
}

func resourceMongoDBAtlasCloudProviderAccessDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())

	projectID := ids["project_id"]
	roleID := ids["id"]
	providerName := ids["provider_name"]

	req := &matlas.CloudProviderDeauthorizationRequest{
		ProviderName: providerName,
		RoleID:       roleID,
		GroupID:      projectID,
	}

	_, err := conn.CloudProviderAccess.DeauthorizeRole(context.Background(), req)

	if err != nil {
		return fmt.Errorf(errorCloudProviderAccessDelete, err)
	}

	return nil
}

func resourceMongoDBAtlasCloudProviderAccessImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	projectID, providerName, roleID, err := splitCloudProviderAccessID(d.Id())

	if err != nil {
		return nil, fmt.Errorf(errorCloudProviderAccessImporter, err)
	}

	// searching id in internal format
	d.SetId(encodeStateID(map[string]string{
		"id":            roleID,
		"project_id":    projectID,
		"provider_name": providerName,
	}))

	err = resourceMongoDBAtlasCloudProviderAccessRead(d, meta)

	if err != nil {
		return nil, fmt.Errorf(errorCloudProviderAccessImporter, err)
	}

	// case of not found
	if d.Id() == "" {
		return nil, fmt.Errorf(errorCloudProviderAccessImporter, " Resource not found at the cloud please check your id")
	}

	// params syncing
	if err = d.Set("project_id", projectID); err != nil {
		return nil, fmt.Errorf(errorCloudProviderAccessImporter, err)
	}

	return []*schema.ResourceData{d}, nil
}

// format  {project_id}-{provider-name}-{role-id}
func splitCloudProviderAccessID(id string) (projectID, providerName, roleID string, err error) {
	var re = regexp.MustCompile(`(?s)^([0-9a-fA-F]{24})-(.*)-(.*)$`)
	parts := re.FindStringSubmatch(id)

	if len(parts) != 4 {
		err = fmt.Errorf(errorCloudProviderAccessImporter, "format please use {project_id}-{provider-name}-{role-id}")
		return
	}

	projectID, providerName, roleID = parts[1], parts[2], parts[3]

	return
}
