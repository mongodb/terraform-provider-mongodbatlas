package mongodbatlas

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"

	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
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
				Type:     schema.TypeString,
				Required: true,
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
		},
	}
}

func dataSourceMongoDBAtlasDatabaseUserRead(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)
	username := d.Get("username").(string)
	databaseName := d.Get("database_name").(string)

	dbUser, _, err := conn.DatabaseUsers.Get(context.Background(), databaseName, projectID, username)
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
	log.Printf("LOG___ dbUser.Labels): %#+v\n", flattenLabels(dbUser.Labels))
	if err := d.Set("labels", flattenLabels(dbUser.Labels)); err != nil {
		return fmt.Errorf("error setting `labels` for database user (%s): %s", d.Id(), err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":    projectID,
		"username":      username,
		"database_name": databaseName,
	}))

	return nil
}
