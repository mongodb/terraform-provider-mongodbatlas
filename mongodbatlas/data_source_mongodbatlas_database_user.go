package mongodbatlas

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasDatabaseUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasDatabaseUserRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"database_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"auth_database_name"},
				Deprecated:    "use auth_database_name instead",
			},
			"auth_database_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"database_name"},
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
				Type:     schema.TypeList,
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
	}
}

func dataSourceMongoDBAtlasDatabaseUserRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)
	username := d.Get("username").(string)

	dbName, dbNameOk := d.GetOk("database_name")
	authDBName, authDBNameOk := d.GetOk("auth_database_name")

	if !dbNameOk && !authDBNameOk {
		return errors.New("one of database_name or auth_database_name must be configured")
	}

	var authDatabaseName string
	if dbNameOk {
		authDatabaseName = dbName.(string)
	} else {
		authDatabaseName = authDBName.(string)
	}

	dbUser, _, err := conn.DatabaseUsers.Get(context.Background(), authDatabaseName, projectID, username)
	if err != nil {
		return fmt.Errorf("error getting database user information: %s", err)
	}

	if err := d.Set("username", dbUser.Username); err != nil {
		return fmt.Errorf("error setting `username` for database user (%s): %s", d.Id(), err)
	}

	if _, ok := d.GetOk("auth_database_name"); ok {
		if err := d.Set("auth_database_name", dbUser.DatabaseName); err != nil {
			return fmt.Errorf("error setting `auth_database_name` for database user (%s): %s", d.Id(), err)
		}
	} else {
		if err := d.Set("database_name", dbUser.DatabaseName); err != nil {
			return fmt.Errorf("error setting `database_name` for database user (%s): %s", d.Id(), err)
		}
	}

	if err := d.Set("x509_type", dbUser.X509Type); err != nil {
		return fmt.Errorf("error setting `x509_type` for database user (%s): %s", d.Id(), err)
	}

	if err := d.Set("aws_iam_type", dbUser.AWSIAMType); err != nil {
		return fmt.Errorf("error setting `aws_iam_type` for database user (%s): %s", d.Id(), err)
	}

	if err := d.Set("ldap_auth_type", dbUser.LDAPAuthType); err != nil {
		return fmt.Errorf("error setting `ldap_auth_type` for database user (%s): %s", d.Id(), err)
	}

	if err := d.Set("roles", flattenRoles(dbUser.Roles)); err != nil {
		return fmt.Errorf("error setting `roles` for database user (%s): %s", d.Id(), err)
	}

	if err := d.Set("labels", flattenLabels(dbUser.Labels)); err != nil {
		return fmt.Errorf("error setting `labels` for database user (%s): %s", d.Id(), err)
	}

	if err := d.Set("scopes", flattenScopes(dbUser.Scopes)); err != nil {
		return fmt.Errorf("error setting `scopes` for database user (%s): %s", d.Id(), err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":         projectID,
		"username":           username,
		"auth_database_name": authDatabaseName,
	}))

	return nil
}
