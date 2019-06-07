package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/resource"

	"github.com/hashicorp/terraform/helper/schema"

	matlas "github.com/mongodb-partners/go-client-mongodb-atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasDatabaseUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasDatabaseUsersRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
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
				},
			},
		},
	}
}

func dataSourceMongoDBAtlasDatabaseUsersRead(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)

	groupID := d.Get("group_id").(string)

	dbUsers, _, err := conn.DatabaseUsers.List(context.Background(), groupID, nil)

	if err != nil {
		return fmt.Errorf("error getting database users information: %s", err)
	}

	if err := d.Set("results", flattenDbUsers(dbUsers)); err != nil {
		return fmt.Errorf("error setting `result` for database users: %s", err)
	}

	d.SetId(resource.UniqueId())

	return nil
}

func flattenDbUsers(dbUsers []matlas.DatabaseUser) []map[string]interface{} {
	var dbUsersMap []map[string]interface{}

	if len(dbUsers) > 0 {
		dbUsersMap = make([]map[string]interface{}, len(dbUsers))

		for k, dbUser := range dbUsers {
			dbUsersMap[k] = map[string]interface{}{
				"roles":         flattenRoles(dbUser.Roles),
				"username":      dbUser.Username,
				"group_id":      dbUser.GroupID,
				"database_name": dbUser.DatabaseName,
			}
		}
	}
	return dbUsersMap
}
