package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
		"collection_type": {
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
		"schedule": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"type": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"end_hour": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"end_minute": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"start_hour": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"start_minute": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"day_of_month": {
						Type:     schema.TypeInt,
						Optional: true,
					},
					"day_of_week": {
						Type:     schema.TypeInt,
						Optional: true,
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
	connV2 := meta.(*MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)
	clusterName := d.Get("cluster_name").(string)
	archiveID := d.Get("archive_id").(string)

	archive, _, err := connV2.OnlineArchiveApi.GetOnlineArchive(ctx, projectID, archiveID, clusterName).Execute()

	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading Online Archive datasource with id %s: %s", archiveID, err.Error()))
	}

	onlineArchiveMap := fromOnlineArchiveToMap(archive)

	for key, val := range onlineArchiveMap {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(fmt.Errorf("error reading Online Archive datasource with id %s: %s", archiveID, err.Error()))
		}
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":   projectID,
		"cluster_name": clusterName,
		"archive_id":   archiveID,
	}))

	return nil
}

func dataSourceMongoDBAtlasOnlineArchivesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	connV2 := meta.(*MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)
	clusterName := d.Get("cluster_name").(string)

	archives, _, err := connV2.OnlineArchiveApi.ListOnlineArchives(ctx, projectID, clusterName).Execute()

	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting Online Archives list for project(%s) in cluster (%s): (%s)", projectID, clusterName, err.Error()))
	}

	results := make([]map[string]interface{}, 0, len(archives.Results))

	for i := range archives.Results {
		archiveData := fromOnlineArchiveToMap(&archives.Results[i])
		archiveData["project_id"] = projectID
		results = append(results, archiveData)
	}

	if err = d.Set("results", results); err != nil {
		return diag.FromErr(fmt.Errorf("error getting Online Archives list for project(%s) in cluster (%s): (%s)", projectID, clusterName, err.Error()))
	}

	if err = d.Set("total_count", archives.TotalCount); err != nil {
		return diag.FromErr(fmt.Errorf("error getting Online Archives list for project(%s) in cluster (%s): (%s)", projectID, clusterName, err.Error()))
	}

	d.SetId(id.UniqueId())

	return nil
}
