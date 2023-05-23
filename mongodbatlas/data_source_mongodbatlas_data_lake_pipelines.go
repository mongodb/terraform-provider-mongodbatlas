package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const errorDataLakePipelineList = "error creating MongoDB Atlas DataLake Pipelines: %s"

func dataSourceMongoDBAtlasDataLakePipelines() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasDataLakePipelinesRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					ReadContext: dataSourceMongoDBAtlasDataLakePipelineRead,
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
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
					},
				},
			},
		},
	}
}

func dataSourceMongoDBAtlasDataLakePipelinesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)

	dataLakePipelines, _, err := conn.DataLakePipeline.List(ctx, projectID)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineList, err))
	}

	if err := d.Set("results", flattenDataLakePipelines(dataLakePipelines)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `result` for DataLake Pipelines: %s", err))
	}

	d.SetId(resource.UniqueId())

	return nil
}

func flattenDataLakePipelines(peers []*matlas.DataLakePipeline) []map[string]interface{} {
	if len(peers) == 0 {
		return nil
	}

	pipelines := make([]map[string]interface{}, len(peers))
	for i := range peers {
		pipelines[i] = map[string]interface{}{
			"project_id":        peers[i].GroupID,
			"name":              peers[i].Name,
			"id":                peers[i].ID,
			"created_date":      peers[i].CreatedDate,
			"last_updated_date": peers[i].LastUpdatedDate,
			"state":             peers[i].State,
			"sink":              flattenDataLakePipelineSink(peers[i].Sink),
			"source":            flattenDataLakePipelineSource(peers[i].Source),
			"transformations":   flattenDataLakePipelineTransformations(peers[i].Transformations),
		}
	}

	return pipelines
}
