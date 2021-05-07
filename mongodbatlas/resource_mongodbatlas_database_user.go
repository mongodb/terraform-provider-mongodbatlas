package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
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
				ConflictsWith: []string{"x509_type", "ldap_auth_type", "aws_iam_type"},
			},
			"x509_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "NONE",
				ValidateFunc: validation.StringInSlice([]string{"NONE", "MANAGED", "CUSTOMER"}, false),
			},
			"ldap_auth_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "NONE",
				ValidateFunc: validation.StringInSlice([]string{"NONE", "USER", "GROUP"}, false),
			},
			"aws_iam_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "NONE",
				ValidateFunc: validation.StringInSlice([]string{"NONE", "USER", "ROLE"}, false),
			},
			"roles": {
				Type:     schema.TypeSet,
				Required: true,
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
			"scopes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceMongoDBAtlasDatabaseUserRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	username := ids["username"]
	authDatabaseName := ids["auth_database_name"]

	if authDatabaseName == "" {
		if dbName, ok := d.GetOk("database_name"); ok {
			authDatabaseName = dbName.(string)
		} else {
			authDatabaseName = d.Get("auth_database_name").(string)
		}
	}
	usernameEscaped := url.PathEscape(username)

	dbUser, _, err := conn.DatabaseUsers.Get(context.Background(), authDatabaseName, projectID, usernameEscaped)
	if err != nil {
		// case 404
		// deleted in the backend case
		reset := strings.Contains(err.Error(), "404") && !d.IsNewResource()

		if reset {
			d.SetId("")
			return nil
		}

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

func resourceMongoDBAtlasDatabaseUserCreate(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
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
		AWSIAMType:   d.Get("aws_iam_type").(string),
		LDAPAuthType: d.Get("ldap_auth_type").(string),
		DatabaseName: authDatabaseName,
		Labels:       expandLabelSliceFromSetSchema(d),
		Scopes:       expandScopes(d),
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
	// Get client connection.
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	username := ids["username"]
	authDatabaseName := ids["auth_database_name"]

	usernameEscaped := url.PathEscape(username)

	dbUser, _, err := conn.DatabaseUsers.Get(context.Background(), authDatabaseName, projectID, usernameEscaped)
	if err != nil {
		return fmt.Errorf("error getting database user information to update it: %s", err)
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

	if d.HasChange("scopes") {
		dbUser.Scopes = expandScopes(d)
	}

	_, _, err = conn.DatabaseUsers.Update(context.Background(), projectID, usernameEscaped, dbUser)
	if err != nil {
		return fmt.Errorf("error updating database user(%s): %s", username, err)
	}

	return resourceMongoDBAtlasDatabaseUserRead(d, meta)
}

func resourceMongoDBAtlasDatabaseUserDelete(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	username := ids["username"]
	authDatabaseName := ids["auth_database_name"]

	usernameEscaped := url.PathEscape(username)

	_, err := conn.DatabaseUsers.Delete(context.Background(), authDatabaseName, projectID, usernameEscaped)
	if err != nil {
		return fmt.Errorf("error deleting database user (%s): %s", username, err)
	}

	return nil
}

func resourceMongoDBAtlasDatabaseUserImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*matlas.Client)

	projectID, username, authDatabaseName, err := splitDatabaseUserImportID(d.Id())
	if err != nil {
		return nil, err
	}

	usernameEscaped := url.PathEscape(*username)

	u, _, err := conn.DatabaseUsers.Get(context.Background(), *authDatabaseName, *projectID, usernameEscaped)
	if err != nil {
		return nil, fmt.Errorf("couldn't import user(%s) in project(%s), error: %s", *username, *projectID, err)
	}

	if err := d.Set("project_id", u.GroupID); err != nil {
		return nil, fmt.Errorf("error setting `project_id` for database user (%s): %s", d.Id(), err)
	}

	if err := d.Set("auth_database_name", u.DatabaseName); err != nil {
		return nil, fmt.Errorf("error setting `auth_database_name` for database user (%s): %s", d.Id(), err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":         *projectID,
		"username":           u.Username,
		"auth_database_name": u.DatabaseName,
	}))

	return []*schema.ResourceData{d}, nil
}

func splitDatabaseUserImportID(id string) (projectID, username, authDatabaseName *string, err error) {
	var re = regexp.MustCompile(`(?s)^([0-9a-fA-F]{24})-(.*)-([$a-z]{1,15})$`)
	parts := re.FindStringSubmatch(id)

	if len(parts) != 4 {
		err = errors.New("import format error: to import a Database User, use the format {project_id}-{username}-{auth_database_name}")
		return
	}

	projectID = &parts[1]
	username = &parts[2]
	authDatabaseName = &parts[3]

	return
}

func expandRoles(d *schema.ResourceData) []matlas.Role {
	var roles []matlas.Role

	if v, ok := d.GetOk("roles"); ok {
		if rs := v.(*schema.Set); rs.Len() > 0 {
			roles = make([]matlas.Role, rs.Len())

			for k, r := range rs.List() {
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

func flattenRoles(roles []matlas.Role) []interface{} {
	roleList := make([]interface{}, 0)
	for _, v := range roles {
		roleList = append(roleList, map[string]interface{}{
			"role_name":       v.RoleName,
			"database_name":   v.DatabaseName,
			"collection_name": v.CollectionName,
		})
	}

	return roleList
}

func flattenScopes(l []matlas.Scope) []map[string]interface{} {
	scopes := make([]map[string]interface{}, len(l))
	for i, v := range l {
		scopes[i] = map[string]interface{}{
			"name": v.Name,
			"type": v.Type,
		}
	}

	return scopes
}

func expandScopes(d *schema.ResourceData) []matlas.Scope {
	list := d.Get("scopes").(*schema.Set)
	res := []matlas.Scope{}
	for _, val := range list.List() {
		v := val.(map[string]interface{})
		scope := matlas.Scope{
			Type: v["type"].(string),
			Name: v["name"].(string),
		}
		res = append(res, scope)
	}

	return res
}
