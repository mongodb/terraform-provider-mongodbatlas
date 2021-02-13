package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasOnlineArchive() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceMongoDBAtlasOnlineArchiveRead,
		Schema: schemaOnlineArchive(),
	}
}

func schemaOnlineArchive() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// argument values
		"project_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"cluster_name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"atlas_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"coll_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"db_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"criteria": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"type": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"date_field": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"date_format": {
						Type:     schema.TypeString,
						Computed: true, // api will set the default
					},
					"expire_after_days": {
						Type:     schema.TypeFloat,
						Computed: true,
					},
					"query": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"partition_fields": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"field_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"order": {
						Type:     schema.TypeFloat,
						Computed: true,
					},
					"field_type": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"paused": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"state": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func dataSourceMongoDBAtlasOnlineArchiveRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)
	clusterName := d.Get("cluster_name").(string)
	atlasID := d.Get("atlas_id").(string)

	outOnlineArchive, _, err := conn.OnlineArchives.Get(context.Background(), projectID, clusterName, atlasID)

	if err != nil {
		return fmt.Errorf("error reading Online Archive datasource with id %s: %s", atlasID, err.Error())
	}

	if err = syncSchema(d, outOnlineArchive); err != nil {
		return fmt.Errorf("error reading Online Archive datasource with id %s: %s", atlasID, err.Error())
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":   projectID,
		"cluster_name": outOnlineArchive.ClusterName,
		"atlas_id":     outOnlineArchive.ID,
	}))

	return nil
}
