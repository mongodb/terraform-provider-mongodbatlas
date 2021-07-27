package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasOnlineArchives() *schema.Resource {
	singleElement := schemaOnlineArchive()
	// overwritten to make them read only
	singleElement["project_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	singleElement["cluster_name"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	singleElement["archive_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasOnlineArchivesRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"cluster_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: singleElement,
				},
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceMongoDBAtlasOnlineArchive() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasOnlineArchiveRead,
		Schema:      schemaOnlineArchive(),
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
		"archive_id": {
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
			Type:     schema.TypeList,
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
						Type:     schema.TypeInt,
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
						Type:     schema.TypeInt,
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

func dataSourceMongoDBAtlasOnlineArchiveRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)
	clusterName := d.Get("cluster_name").(string)
	atlasID := d.Get("archive_id").(string)

	archive, _, err := conn.OnlineArchives.Get(ctx, projectID, clusterName, atlasID)

	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading Online Archive datasource with id %s: %s", atlasID, err.Error()))
	}

	onlineArchiveMap := fromOnlineArchiveToMap(archive)

	for key, val := range onlineArchiveMap {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(fmt.Errorf("error reading Online Archive datasource with id %s: %s", atlasID, err.Error()))
		}
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":   projectID,
		"cluster_name": archive.ClusterName,
		"archive_id":   archive.ID,
	}))

	return nil
}

func dataSourceMongoDBAtlasOnlineArchivesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)
	clusterName := d.Get("cluster_name").(string)
	archives, _, err := conn.OnlineArchives.List(ctx, projectID, clusterName, &matlas.ListOptions{})

	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting Online Archives list for project(%s) in cluster (%s): (%s)", projectID, clusterName, err.Error()))
	}

	results := make([]map[string]interface{}, 0, len(archives.Results))

	for _, archive := range archives.Results {
		archiveData := fromOnlineArchiveToMap(archive)
		archiveData["project_id"] = projectID
		results = append(results, archiveData)
	}

	if err = d.Set("results", results); err != nil {
		return diag.FromErr(fmt.Errorf("error getting Online Archives list for project(%s) in cluster (%s): (%s)", projectID, clusterName, err.Error()))
	}

	if err = d.Set("total_count", archives.TotalCount); err != nil {
		return diag.FromErr(fmt.Errorf("error getting Online Archives list for project(%s) in cluster (%s): (%s)", projectID, clusterName, err.Error()))
	}

	d.SetId(resource.UniqueId())

	return nil
}
