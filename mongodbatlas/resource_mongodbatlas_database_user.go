package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"log"

	"strings"

	"github.com/hashicorp/terraform/helper/schema"

	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func resourceMongoDBAtlasDatabaseUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoDBAtlasDatabaseUserCreate,
		Read:   resourceMongoDBAtlasDatabaseUserRead,
		Update: resourceMongoDBAtlasDatabaseUserUpdate,
		Delete: resourceMongoDBAtlasDatabaseUserDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMongoDBAtlasDatabaseUserImportState,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"x509_type"},
			},
			"x509_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "NONE",
			},
			"roles": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"collection_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"database_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"labels": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceMongoDBAtlasDatabaseUserRead(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	username := ids["username"]
	authDatabaseName := ids["auth_database_name"]

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
	if err := d.Set("roles", flattenRoles(dbUser.Roles)); err != nil {
		return fmt.Errorf("error setting `roles` for database user (%s): %s", d.Id(), err)
	}
	if err := d.Set("labels", flattenLabels(dbUser.Labels)); err != nil {
		return fmt.Errorf("error setting `labels` for database user (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceMongoDBAtlasDatabaseUserCreate(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)

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

	dbUserReq := &matlas.DatabaseUser{
		Roles:        expandRoles(d),
		GroupID:      projectID,
		Username:     d.Get("username").(string),
		Password:     d.Get("password").(string),
		X509Type:     d.Get("x509_type").(string),
		DatabaseName: authDatabaseName,
		Labels:       expandLabelSliceFromSetSchema(d),
	}

	dbUserRes, _, err := conn.DatabaseUsers.Create(context.Background(), projectID, dbUserReq)
	if err != nil {
		return fmt.Errorf("error creating database user: %s", err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":         projectID,
		"username":           dbUserRes.Username,
		"auth_database_name": authDatabaseName,
	}))

	return resourceMongoDBAtlasDatabaseUserRead(d, meta)
}

func resourceMongoDBAtlasDatabaseUserUpdate(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	username := ids["username"]
	authDatabaseName := ids["auth_database_name"]

	dbUser, _, err := conn.DatabaseUsers.Get(context.Background(), authDatabaseName, projectID, username)
	if err != nil {
		return fmt.Errorf("error getting database user information: %s", err)
	}

	if d.HasChange("password") {
		dbUser.Password = d.Get("password").(string)
	}

	if d.HasChange("roles") {
		dbUser.Roles = expandRoles(d)
	}

	if d.HasChange("labels") {
		dbUser.Labels = expandLabelSliceFromSetSchema(d)
	}

	_, _, err = conn.DatabaseUsers.Update(context.Background(), projectID, username, dbUser)
	if err != nil {
		return fmt.Errorf("error updating database user(%s): %s", username, err)
	}

	return resourceMongoDBAtlasDatabaseUserRead(d, meta)
}

func resourceMongoDBAtlasDatabaseUserDelete(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	username := ids["username"]
	authDatabaseName := ids["auth_database_name"]

	_, err := conn.DatabaseUsers.Delete(context.Background(), authDatabaseName, projectID, username)
	if err != nil {
		return fmt.Errorf("error deleting database user (%s): %s", username, err)
	}
	return nil
}

func resourceMongoDBAtlasDatabaseUserImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*matlas.Client)

	parts := strings.SplitN(d.Id(), "-", 3)
	if len(parts) != 3 {
		return nil, errors.New("import format error: to import a database user, use the format {project_id}-{username}-{auth_database_name}")
	}

	projectID := parts[0]
	username := parts[1]
	authDatabaseName := parts[2]

	u, _, err := conn.DatabaseUsers.Get(context.Background(), authDatabaseName, projectID, username)
	if err != nil {
		return nil, fmt.Errorf("couldn't import user(%s) in project(%s), error: %s", username, projectID, err)
	}

	if err := d.Set("project_id", u.GroupID); err != nil {
		log.Printf("[WARN] Error setting project_id for (%s): %s", d.Id(), err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":         projectID,
		"username":           username,
		"auth_database_name": authDatabaseName,
	}))

	return []*schema.ResourceData{d}, nil
}

func expandRoles(d *schema.ResourceData) []matlas.Role {
	var roles []matlas.Role
	if v, ok := d.GetOk("roles"); ok {
		if rs := v.([]interface{}); len(rs) > 0 {
			roles = make([]matlas.Role, len(rs))
			for k, r := range rs {
				roleMap := r.(map[string]interface{})
				roles[k] = matlas.Role{
					RoleName:       roleMap["role_name"].(string),
					DatabaseName:   roleMap["database_name"].(string),
					CollectionName: roleMap["collection_name"].(string),
				}
			}
		}
	}
	return roles
}

func flattenRoles(roles []matlas.Role) []map[string]interface{} {
	roleList := make([]map[string]interface{}, 0)
	for _, v := range roles {
		roleList = append(roleList, map[string]interface{}{
			"role_name":       v.RoleName,
			"database_name":   v.DatabaseName,
			"collection_name": v.CollectionName,
		})
	}
	return roleList
}
