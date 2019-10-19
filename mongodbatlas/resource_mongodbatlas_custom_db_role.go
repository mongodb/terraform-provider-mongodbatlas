package mongodbatlas

import "github.com/hashicorp/terraform/helper/schema"

func resourceMongoDBAtlasCustomDBRole() *schema.Resource {
	return &schema.Resource{
		Create:   nil,
		Read:     nil,
		Update:   nil,
		Delete:   nil,
		Importer: nil,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"actions": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:     schema.TypeString,
							Required: true,
						},
						"resources": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"collection_name": {
										Type:          schema.TypeString,
										Optional:      true,
										ConflictsWith: []string{"cluster"},
									},
									"database_name": {
										Type:          schema.TypeString,
										Optional:      true,
										ConflictsWith: []string{"cluster"},
									},
									"cluster": {
										Type:          schema.TypeBool,
										Optional:      true,
										ConflictsWith: []string{"database_name", "collection_name"},
									},
								},
							},
						},
					},
				},
			},
			"inherited_roles": {
				Type:     schema.TypeList,
				Required: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"role_name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}
