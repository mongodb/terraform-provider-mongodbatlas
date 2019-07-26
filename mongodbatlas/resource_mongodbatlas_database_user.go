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
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
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
		},
	}
}

func resourceMongoDBAtlasDatabaseUserRead(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	username := ids["username"]

	dbUser, _, err := conn.DatabaseUsers.Get(context.Background(), projectID, username)

	if err != nil {
		return fmt.Errorf("error getting database user information: %s", err)
	}
	if err := d.Set("username", dbUser.Username); err != nil {
		return fmt.Errorf("error setting `username` for database user (%s): %s", d.Id(), err)
	}
	if err := d.Set("database_name", dbUser.DatabaseName); err != nil {
		return fmt.Errorf("error setting `database_name` for database user (%s): %s", d.Id(), err)
	}
	if err := d.Set("roles", flattenRoles(dbUser.Roles)); err != nil {
		return fmt.Errorf("error setting `roles` for database user (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceMongoDBAtlasDatabaseUserCreate(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)

	dbUserReq := &matlas.DatabaseUser{
		Roles:        expandRoles(d),
		GroupID:      projectID,
		Username:     d.Get("username").(string),
		DatabaseName: d.Get("database_name").(string),
	}

	if v, ok := d.GetOk("password"); ok {
		dbUserReq.Password = v.(string)
	}

	dbUserRes, _, err := conn.DatabaseUsers.Create(context.Background(), projectID, dbUserReq)

	if err != nil {
		return fmt.Errorf("error creating database user: %s", err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"username":   dbUserRes.Username,
	}))

	return resourceMongoDBAtlasDatabaseUserRead(d, meta)
}

func resourceMongoDBAtlasDatabaseUserUpdate(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	username := ids["username"]

	dbUser, _, err := conn.DatabaseUsers.Get(context.Background(), projectID, username)

	if err != nil {
		return fmt.Errorf("error getting database user information: %s", err)
	}

	if d.HasChange("password") {
		dbUser.Password = d.Get("password").(string)
	}

	if d.HasChange("roles") {
		dbUser.Roles = expandRoles(d)
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

	_, err := conn.DatabaseUsers.Delete(context.Background(), projectID, username)

	if err != nil {
		return fmt.Errorf("error deleting database user (%s): %s", username, err)
	}
	return nil
}

func resourceMongoDBAtlasDatabaseUserImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*matlas.Client)

	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 2 {
		return nil, errors.New("import format error: to import a database user, use the format {project_id}-{username}")
	}

	projectID := parts[0]
	username := parts[1]

	u, _, err := conn.DatabaseUsers.Get(context.Background(), projectID, username)
	if err != nil {
		return nil, fmt.Errorf("couldn't import user %s in project %s, error: %s", username, projectID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"username":   u.Username,
	}))

	if err := d.Set("project_id", u.GroupID); err != nil {
		log.Printf("[WARN] Error setting project_id for (%s): %s", d.Id(), err)
	}

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
