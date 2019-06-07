package mongodbatlas

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"

	matlas "github.com/mongodb-partners/go-client-mongodb-atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasDatabaseUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasDatabaseUserRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},

			"database_name": {
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
		},
	}
}

func dataSourceMongoDBAtlasDatabaseUserRead(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	groupID := d.Get("group_id").(string)
	username := d.Get("username").(string)

	dbUser, _, err := conn.DatabaseUsers.Get(context.Background(), groupID, username)

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

	d.SetId(dbUser.Username)

	return nil
}
