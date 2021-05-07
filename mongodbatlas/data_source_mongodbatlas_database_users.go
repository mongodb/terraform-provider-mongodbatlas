package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasDatabaseUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasDatabaseUsersRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auth_database_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"x509_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"aws_iam_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ldap_auth_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"roles": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"role_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"collection_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"database_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"labels": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"scopes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceMongoDBAtlasDatabaseUsersRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*matlas.Client)

	projectID := d.Get("project_id").(string)

	dbUsers, _, err := conn.DatabaseUsers.List(context.Background(), projectID, nil)
	if err != nil {
		return fmt.Errorf("error getting database users information: %s", err)
	}

	if err := d.Set("results", flattenDBUsers(dbUsers)); err != nil {
		return fmt.Errorf("error setting `result` for database users: %s", err)
	}

	d.SetId(resource.UniqueId())

	return nil
}

func flattenDBUsers(dbUsers []matlas.DatabaseUser) []map[string]interface{} {
	var dbUsersMap []map[string]interface{}

	if len(dbUsers) > 0 {
		dbUsersMap = make([]map[string]interface{}, len(dbUsers))

		for i := range dbUsers {
			dbUsersMap[i] = map[string]interface{}{
				"roles":              flattenRoles(dbUsers[i].Roles),
				"username":           dbUsers[i].Username,
				"project_id":         dbUsers[i].GroupID,
				"auth_database_name": dbUsers[i].DatabaseName,
				"x509_type":          dbUsers[i].X509Type,
				"aws_iam_type":       dbUsers[i].AWSIAMType,
				"ldap_auth_type":     dbUsers[i].LDAPAuthType,
				"labels":             flattenLabels(dbUsers[i].Labels),
				"scopes":             flattenScopes(dbUsers[i].Scopes),
			}
		}
	}

	return dbUsersMap
}
