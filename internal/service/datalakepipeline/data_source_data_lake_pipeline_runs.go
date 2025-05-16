package datalakepipeline

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312003/admin"
)

const errorDataLakePipelineRunList = "error reading MongoDB Atlas DataLake Runs (%s): %s"

func PluralDataSourceRun() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "Data Lake is deprecated. As of September 2024, Data Lake is deprecated and will reach end-of-life. To learn more, see https://dochub.mongodb.org/core/data-lake-deprecation",
		ReadContext:        dataSourcePluralRunRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pipeline_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
				},
			},
		},
	}
}

func dataSourcePluralRunRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)
	name := d.Get("pipeline_name").(string)
	runs, _, err := connV2.DataLakePipelinesApi.ListPipelineRuns(ctx, projectID, name).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineRunList, projectID, err))
	}
	if err := d.Set("results", flattenRunResults(runs.GetResults())); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorDataLakeSetting, "results", projectID, err))
	}
	d.SetId(id.UniqueId())
	return nil
}

func flattenRunResults(datalakePipelineRuns []admin.IngestionPipelineRun) []map[string]any {
	if len(datalakePipelineRuns) == 0 {
		return nil
	}
	results := make([]map[string]any, len(datalakePipelineRuns))

	for k, run := range datalakePipelineRuns {
		results[k] = map[string]any{
			"id":                    run.GetId(),
			"created_date":          conversion.TimePtrToStringPtr(run.CreatedDate),
			"last_updated_date":     conversion.TimePtrToStringPtr(run.LastUpdatedDate),
			"state":                 run.GetState(),
			"pipeline_id":           run.GetPipelineId(),
			"snapshot_id":           run.GetSnapshotId(),
			"backup_frequency_type": run.GetBackupFrequencyType(),
			"stats":                 flattenRunStats(run.Stats),
		}
	}
	return results
}
