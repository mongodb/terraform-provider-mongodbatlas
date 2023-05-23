package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasDataLakePipeline() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasDataLakePipelineRead,
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
			"snapshots":           dataSourceSchemaDataLakePipelineSnapshots(),
			"ingestion_schedules": dataSourceSchemaDataLakePipelineIngestionSchedules(),
		},
	}
}

func dataSourceSchemaDataLakePipelineIngestionSchedules() *schema.Schema {
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

func dataSourceSchemaDataLakePipelineSnapshots() *schema.Schema {
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

func dataSourceMongoDBAtlasDataLakePipelineRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)
	name := d.Get("name").(string)

	dataLakePipeline, _, err := conn.DataLakePipeline.Get(ctx, projectID, name)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineRead, name, err))
	}

	snapshots, _, err := conn.DataLakePipeline.ListSnapshots(ctx, projectID, name, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineRead, name, err))
	}

	ingestionSchedules, _, err := conn.DataLakePipeline.ListIngestionSchedules(ctx, projectID, name)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineRead, name, err))
	}

	return setDataLakeResourceData(d, dataLakePipeline, snapshots, ingestionSchedules)
}

func setDataLakeResourceData(
	d *schema.ResourceData,
	pipeline *matlas.DataLakePipeline,
	snapshots *matlas.DataLakePipelineSnapshotsResponse,
	ingestionSchedules []*matlas.DataLakePipelineIngestionSchedule) diag.Diagnostics {
	if err := d.Set("id", pipeline.ID); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakeSetting, "id", pipeline.Name, err))
	}

	if err := d.Set("state", pipeline.State); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineSetting, "state", pipeline.Name, err))
	}

	if err := d.Set("created_date", pipeline.CreatedDate); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineSetting, "created_date", pipeline.Name, err))
	}

	if err := d.Set("last_updated_date", pipeline.LastUpdatedDate); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineSetting, "last_updated_date", pipeline.Name, err))
	}

	if err := d.Set("sink", flattenDataLakePipelineSink(pipeline.Sink)); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineSetting, "sink", pipeline.Name, err))
	}

	if err := d.Set("source", flattenDataLakePipelineSource(pipeline.Source)); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineSetting, "source", pipeline.Name, err))
	}

	if err := d.Set("transformations", flattenDataLakePipelineTransformations(pipeline.Transformations)); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineSetting, "transformations", pipeline.Name, err))
	}

	if err := d.Set("snapshots", flattenDataLakePipelineSnapshots(snapshots.Results)); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineSetting, "snapshots", pipeline.Name, err))
	}

	if err := d.Set("ingestion_schedules", flattenDataLakePipelineIngestionSchedules(ingestionSchedules)); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineSetting, "ingestion_schedules", pipeline.Name, err))
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": pipeline.GroupID,
		"name":       pipeline.Name,
	}))

	return nil
}
