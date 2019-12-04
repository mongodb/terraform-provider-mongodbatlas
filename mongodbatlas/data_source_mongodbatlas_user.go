package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"

	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasAtlasUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasUserRead,
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"country": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"email_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mobile_number": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"first_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"roles": {
				Type:     schema.TypeList,
				Computed: true,
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
		},
	}
}

func dataSourceMongoDBAtlasUserRead(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)

	userID := d.Get("user_id").(string)

	user, resp, err := conn.AtlasUsers.Get(context.Background(), userID)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf(errorReadAtlasUser, d.Id(), err)
	}

	if err := d.Set("username", user.Username); err != nil {
		return fmt.Errorf(errorReadAtlasUser, d.Id(), err)
	}

	if err := d.Set("country", user.Country); err != nil {
		return fmt.Errorf(errorReadAtlasUser, d.Id(), err)
	}

	if err := d.Set("email_address", user.EmailAddress); err != nil {
		return fmt.Errorf(errorReadAtlasUser, d.Id(), err)
	}

	if err := d.Set("mobile_number", user.MobileNumber); err != nil {
		return fmt.Errorf(errorReadAtlasUser, d.Id(), err)
	}

	if err := d.Set("first_name", user.FirstName); err != nil {
		return fmt.Errorf(errorReadAtlasUser, d.Id(), err)
	}

	if err := d.Set("last_name", user.LastName); err != nil {
		return fmt.Errorf(errorReadAtlasUser, d.Id(), err)
	}

	if err := d.Set("roles", flattenAtlasRoles(user.Roles)); err != nil {
		return fmt.Errorf(errorReadAtlasUser, d.Id(), err)
	}

	if err := d.Set("team_ids", user.TeamIds); err != nil {
		return fmt.Errorf(errorReadAtlasUser, d.Id(), err)
	}

	d.SetId(user.ID)

	return nil
}

func flattenAtlasRoles(roles []matlas.AtlasRole) []map[string]interface{} {
	var atlasRoles []map[string]interface{}

	for _, role := range roles {
		atlasRole := map[string]interface{}{
			"project_id": role.GroupID,
			"org_id":     role.OrgID,
			"role_name":  role.RoleName,
		}
		atlasRoles = append(atlasRoles, atlasRole)
	}
	return atlasRoles
}
