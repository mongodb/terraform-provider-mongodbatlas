package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const errorDataLakePipelineRunRead = "error reading MongoDB Atlas DataLake Run (%s): %s"

func dataSourceMongoDBAtlasDataLakePipelineRun() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasDataLakeRunRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pipeline_run_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_updated_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dataset_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"phase": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pipeline_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_frequency_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"stats": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 0,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bytes_exported": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"num_docs": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceMongoDBAtlasDataLakeRunRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)
	name := d.Get("name").(string)
	pipelineRunID := d.Get("pipeline_run_id").(string)

	dataLakeRun, resp, err := conn.DataLakePipeline.GetRun(ctx, projectID, name, pipelineRunID)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorDataLakePipelineRunRead, name, err))
	}

	if err := d.Set("id", dataLakeRun.ID); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakeSetting, "hostnames", name, err))
	}

	if err := d.Set("project_id", dataLakeRun.GroupID); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakeSetting, "state", name, err))
	}

	if err := d.Set("created_date", dataLakeRun.CreatedDate); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakeSetting, "storage_databases", name, err))
	}

	if err := d.Set("last_updated_date", dataLakeRun.LastUpdatedDate); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakeSetting, "storage_databases", name, err))
	}

	if err := d.Set("state", dataLakeRun.State); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakeSetting, "storage_databases", name, err))
	}

	if err := d.Set("phase", dataLakeRun.Phase); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakeSetting, "storage_databases", name, err))
	}

	if err := d.Set("pipeline_id", dataLakeRun.PipelineID); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakeSetting, "storage_stores", name, err))
	}

	if err := d.Set("dataset_name", dataLakeRun.DatasetName); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakeSetting, "storage_stores", name, err))
	}

	if err := d.Set("snapshot_id", dataLakeRun.SnapshotID); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakeSetting, "storage_stores", name, err))
	}

	if err := d.Set("backup_frequency_type", dataLakeRun.BackupFrequencyType); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakeSetting, "storage_stores", name, err))
	}

	if err := d.Set("stats", flattenDataLakePipelineRunStats(dataLakeRun.Stats)); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakeSetting, "storage_stores", name, err))
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":      projectID,
		"name":            name,
		"pipeline_run_id": pipelineRunID,
	}))

	return nil
}

func flattenDataLakePipelineRunStats(datalakeRunStats *matlas.DataLakePipelineRunStats) []map[string]interface{} {
	if datalakeRunStats == nil {
		return nil
	}

	maps := make([]map[string]interface{}, 1)
	maps[0] = map[string]interface{}{
		"bytes_exported": datalakeRunStats.BytesExported,
		"num_docs":       datalakeRunStats.NumDocs,
	}
	return maps
}
