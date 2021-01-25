package mongodbatlas

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceMongoDBAtlasOnlineArchive() *schema.Resource {
	return &schema.Resource{
		Schema: getMongoDBAtlasOnlineArchiveSchema(),
	}
}

func getMongoDBAtlasOnlineArchiveSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"project_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"cluster_name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"coll_name": {},
		"db_name":   {},
		"criteria":  {},
	}

}
