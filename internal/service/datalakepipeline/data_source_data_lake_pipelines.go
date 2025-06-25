package datalakepipeline

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312004/admin"
)

const errorDataLakePipelineList = "error creating MongoDB Atlas DataLake Pipelines: %s"

func PluralDataSource() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "Data Lake is deprecated. As of September 2024, Data Lake is deprecated and will reach end-of-life. To learn more, see https://dochub.mongodb.org/core/data-lake-deprecation",
		ReadContext:        dataSourcePluralRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					ReadContext: dataSourceRead,
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

func dataSourcePluralRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)

	pipelines, _, err := connV2.DataLakePipelinesApi.ListPipelines(ctx, projectID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineList, err))
	}

	if err := d.Set("results", flattenDataLakePipelines(pipelines)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `result` for DataLake Pipelines: %s", err))
	}

	d.SetId(id.UniqueId())
	return nil
}

func flattenDataLakePipelines(peers []admin.DataLakeIngestionPipeline) []map[string]any {
	pipelines := make([]map[string]any, len(peers))
	for i := range peers {
		pipelines[i] = map[string]any{
			"project_id":        peers[i].GetGroupId(),
			"name":              peers[i].GetName(),
			"id":                peers[i].GetId(),
			"created_date":      conversion.TimePtrToStringPtr(peers[i].CreatedDate),
			"last_updated_date": conversion.TimePtrToStringPtr(peers[i].LastUpdatedDate),
			"state":             peers[i].GetState(),
			"sink":              flattenSink(peers[i].Sink),
			"source":            flattenSource(peers[i].Source),
			"transformations":   flattenTransformations(peers[i].GetTransformations()),
		}
	}
	return pipelines
}
