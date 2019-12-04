package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"

	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

const (
	errorCreateAtlasUser = "error creating MongoDB Atlas User: %s"
	errorReadAtlasUser   = "error reading MongoDB Atlas User (%s): %s"
)

func resourceMongoDBAtlasAtlasUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoDBAtlasUserCreate,
		Read:   resourceMongoDBAtlasUserRead,
		Delete: resourceMongoDBAtlasUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				ForceNew:  true,
			},
			"country": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"email_address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"mobile_number": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"first_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"roles": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"org_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"role_name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"team_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceMongoDBAtlasUserCreate(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)

	userRequest := &matlas.AtlasUser{
		EmailAddress: d.Get("email_address").(string),
		FirstName:    d.Get("first_name").(string),
		LastName:     d.Get("last_name").(string),
		Username:     d.Get("username").(string),
		MobileNumber: d.Get("mobile_number").(string),
		Country:      d.Get("country").(string),
		Password:     d.Get("password").(string),
		Roles:        expandAtlasRoles(d.Get("roles").([]interface{})),
	}

	user, _, err := conn.AtlasUsers.Create(context.Background(), userRequest)
	if err != nil {
		return fmt.Errorf(errorCreateAtlasUser, err)
	}

	d.SetId(user.ID)

	return resourceMongoDBAtlasUserRead(d, meta)
}

func resourceMongoDBAtlasUserRead(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)

	user, resp, err := conn.AtlasUsers.Get(context.Background(), d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf(errorReadAtlasUser, d.Id(), err)
	}

	if err := d.Set("user_id", user.ID); err != nil {
		return fmt.Errorf(errorReadAtlasUser, d.Id(), err)
	}

	if err := d.Set("team_ids", user.TeamIds); err != nil {
		return fmt.Errorf(errorReadAtlasUser, d.Id(), err)
	}

	return nil
}

func resourceMongoDBAtlasUserDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func expandAtlasRoles(atlasRoles []interface{}) (roles []matlas.AtlasRole) {
	for _, v := range atlasRoles {
		roleMap := v.(map[string]interface{})

		role := matlas.AtlasRole{
			RoleName: roleMap["role_name"].(string),
			GroupID:  roleMap["project_id"].(string),
			OrgID:    roleMap["org_id"].(string),
		}

		roles = append(roles, role)
	}
	return roles
}
