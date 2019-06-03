package mongodbatlas

import (
	"github.com/hashicorp/terraform/helper/schema"
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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
	return nil
}

func resourceMongoDBAtlasDatabaseUserUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceMongoDBAtlasDatabaseUserDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
