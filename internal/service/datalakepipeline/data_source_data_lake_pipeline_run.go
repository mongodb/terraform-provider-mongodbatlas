package datalakepipeline

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312006/admin"
)

const errorDataLakePipelineRunRead = "error reading MongoDB Atlas DataLake Run (%s): %s"

func DataSourceRun() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "Data Lake is deprecated. As of September 2024, Data Lake is deprecated and will reach end-of-life. To learn more, see https://dochub.mongodb.org/core/data-lake-deprecation",
		ReadContext:        dataSourceRunRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pipeline_name": {
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

func dataSourceRunRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)
	name := d.Get("pipeline_name").(string)
	pipelineRunID := d.Get("pipeline_run_id").(string)

	run, resp, err := connV2.DataLakePipelinesApi.GetPipelineRun(ctx, projectID, name, pipelineRunID).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorDataLakePipelineRunRead, name, err))
	}

	if err := d.Set("id", run.GetId()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorDataLakeSetting, "hostnames", name, err))
	}

	if err := d.Set("project_id", run.GetGroupId()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorDataLakeSetting, "state", name, err))
	}

	if err := d.Set("created_date", conversion.TimePtrToStringPtr(run.CreatedDate)); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorDataLakeSetting, "storage_databases", name, err))
	}

	if err := d.Set("last_updated_date", conversion.TimePtrToStringPtr(run.LastUpdatedDate)); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorDataLakeSetting, "storage_databases", name, err))
	}

	if err := d.Set("state", run.GetState()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorDataLakeSetting, "storage_databases", name, err))
	}

	if err := d.Set("phase", run.GetPhase()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorDataLakeSetting, "storage_databases", name, err))
	}

	if err := d.Set("pipeline_id", run.GetPipelineId()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorDataLakeSetting, "storage_stores", name, err))
	}

	if err := d.Set("dataset_name", run.GetDatasetName()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorDataLakeSetting, "storage_stores", name, err))
	}

	if err := d.Set("snapshot_id", run.GetSnapshotId()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorDataLakeSetting, "storage_stores", name, err))
	}

	if err := d.Set("backup_frequency_type", run.GetBackupFrequencyType()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorDataLakeSetting, "storage_stores", name, err))
	}

	if err := d.Set("stats", flattenRunStats(run.Stats)); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorDataLakeSetting, "storage_stores", name, err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":      projectID,
		"pipeline_name":   name,
		"pipeline_run_id": pipelineRunID,
	}))

	return nil
}

func flattenRunStats(stats *admin.PipelineRunStats) []map[string]any {
	if stats == nil {
		return nil
	}
	maps := make([]map[string]any, 1)
	maps[0] = map[string]any{
		"bytes_exported": stats.GetBytesExported(),
		"num_docs":       stats.GetNumDocs(),
	}
	return maps
}
