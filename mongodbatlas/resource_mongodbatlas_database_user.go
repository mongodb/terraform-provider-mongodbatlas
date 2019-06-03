package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	matlas "github.com/mongodb-partners/go-client-mongodbatlas/mongodbatlas"
)

func resourceMongoDBAtlasDatabaseUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoDBAtlasDatabaseUserCreate,
		Read:   resourceMongoDBAtlasDatabaseUserRead,
		Update: resourceMongoDBAtlasDatabaseUserUpdate,
		Delete: resourceMongoDBAtlasDatabaseUserDelete,
		Schema: map[string]*schema.Schema{
			"group_id": {
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
			"delete_after_date": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.ValidateRFC3339TimeString,
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
	return nil
}

func resourceMongoDBAtlasDatabaseUserCreate(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	groupID := d.Get("group_id").(string)

	dbUserReq := &matlas.DatabaseUser{
		Roles:        expandRoles(d),
		GroupID:      groupID,
		Username:     d.Get("username").(string),
		DatabaseName: d.Get("database_name").(string),
	}

	if v, ok := d.GetOk("password"); ok {
		dbUserReq.Password = v.(string)
	}

	if v, ok := d.GetOk("delete_after_date"); ok {
		dbUserReq.DeleteAfterDate = v.(string)
	}

	dbUserRes, _, err := conn.DatabaseUsers.Create(context.Background(), groupID, dbUserReq)

	if err != nil {
		return fmt.Errorf("error creating database user: %s", err)
	}

	d.SetId(dbUserRes.Username)

	return resourceMongoDBAtlasDatabaseUserRead(d, meta)
}

func resourceMongoDBAtlasDatabaseUserUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceMongoDBAtlasDatabaseUserDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
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
