package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const errorDataLakePipelineRunList = "error reading MongoDB Atlas DataLake Runs (%s): %s"

func dataSourceMongoDBAtlasDataLakePipelineRuns() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasDataLakeRunsRead,
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

func dataSourceMongoDBAtlasDataLakeRunsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)
	name := d.Get("pipeline_name").(string)

	dataLakeRuns, _, err := conn.DataLakePipeline.ListRuns(ctx, projectID, name)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineRunList, projectID, err))
	}

	if err := d.Set("results", flattenDataLakePipelineRunResult(dataLakeRuns.Results)); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakeSetting, "results", projectID, err))
	}

	d.SetId(id.UniqueId())

	return nil
}

func flattenDataLakePipelineRunResult(datalakePipelineRuns []*matlas.DataLakePipelineRun) []map[string]interface{} {
	var results []map[string]interface{}

	if len(datalakePipelineRuns) == 0 {
		return results
	}

	results = make([]map[string]interface{}, len(datalakePipelineRuns))

	for k, run := range datalakePipelineRuns {
		results[k] = map[string]interface{}{
			"id":                    run.ID,
			"created_date":          run.CreatedDate,
			"last_updated_date":     run.LastUpdatedDate,
			"state":                 run.State,
			"pipeline_id":           run.PipelineID,
			"snapshot_id":           run.SnapshotID,
			"backup_frequency_type": run.BackupFrequencyType,
			"stats":                 flattenDataLakePipelineRunStats(run.Stats),
		}
	}

	return results
}
