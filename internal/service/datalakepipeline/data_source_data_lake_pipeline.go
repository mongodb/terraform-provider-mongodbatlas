package datalakepipeline

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext:        dataSourceRead,
		DeprecationMessage: "Data Lake is deprecated. As of September 2024, Data Lake is deprecated and will reach end-of-life. To learn more, see https://dochub.mongodb.org/core/data-lake-deprecation",
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
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
			"sink": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"provider": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
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
								},
							},
						},
					},
				},
			},
			"source": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"collection_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"database_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"transformations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"snapshots":           dataSourceSchemaSnapshots(),
			"ingestion_schedules": dataSourceSchemaIngestionSchedules(),
		},
	}
}

func dataSourceSchemaIngestionSchedules() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"frequency_type": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"retention_unit": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"retention_value": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"frequency_interval": {
					Type:     schema.TypeInt,
					Computed: true,
				},
			},
		},
	}
}

func dataSourceSchemaSnapshots() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"provider": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"created_at": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"expires_at": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"frequency_yype": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"master_key": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"mongod_version": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"replica_set_name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"type": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"status": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"size": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"copy_region": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"policies": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)
	name := d.Get("name").(string)

	pipeline, _, err := connV2.DataLakePipelinesApi.GetPipeline(ctx, projectID, name).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineRead, name, err))
	}

	snapshots, _, err := connV2.DataLakePipelinesApi.GetAvailablePipelineSnapshots(ctx, projectID, name).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineRead, name, err))
	}

	ingestionSchedules, _, err := connV2.DataLakePipelinesApi.GetAvailablePipelineSchedules(ctx, projectID, name).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineRead, name, err))
	}

	pipelineName := pipeline.GetName()

	if err := d.Set("id", pipeline.GetId()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorDataLakeSetting, "id", pipelineName, err))
	}

	if err := d.Set("state", pipeline.GetState()); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineSetting, "state", pipelineName, err))
	}

	if err := d.Set("created_date", conversion.TimePtrToStringPtr(pipeline.CreatedDate)); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineSetting, "created_date", pipelineName, err))
	}

	if err := d.Set("last_updated_date", conversion.TimePtrToStringPtr(pipeline.LastUpdatedDate)); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineSetting, "last_updated_date", pipelineName, err))
	}

	if err := d.Set("sink", flattenSink(pipeline.Sink)); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineSetting, "sink", pipelineName, err))
	}

	if err := d.Set("source", flattenSource(pipeline.Source)); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineSetting, "source", pipelineName, err))
	}

	if err := d.Set("transformations", flattenTransformations(pipeline.GetTransformations())); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineSetting, "transformations", pipelineName, err))
	}

	if err := d.Set("snapshots", flattenSnapshots(snapshots.GetResults())); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineSetting, "snapshots", pipelineName, err))
	}

	if err := d.Set("ingestion_schedules", flattenIngestionSchedules(ingestionSchedules)); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineSetting, "ingestion_schedules", pipelineName, err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id": pipeline.GetGroupId(),
		"name":       pipelineName,
	}))

	return nil
}
