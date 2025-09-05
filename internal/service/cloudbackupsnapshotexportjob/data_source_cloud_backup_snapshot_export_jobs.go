package cloudbackupsnapshotexportjob

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312007/admin"
)

func PluralDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"page_num": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"items_per_page": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"export_job_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"custom_data": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								}},
						},
						"components": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"export_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"replica_set_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"export_bucket_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"export_status_exported_collections": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"export_status_total_collections": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"finished_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"prefix": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	projectID := d.Get("project_id").(string)
	clusterName := d.Get("cluster_name").(string)
	request := admin.ListBackupExportsApiParams{
		GroupId:      projectID,
		ClusterName:  clusterName,
		PageNum:      conversion.IntPtr(d.Get("page_num").(int)),
		ItemsPerPage: conversion.IntPtr(d.Get("items_per_page").(int)),
	}

	jobs, _, err := connV2.CloudBackupsApi.ListBackupExportsWithParams(ctx, &request).Execute()
	if err != nil {
		return diag.Errorf("error getting CloudProviderSnapshotExportJobs information: %s", err)
	}

	if err := d.Set("results", flattenCloudBackupSnapshotExportJobs(jobs.GetResults())); err != nil {
		return diag.Errorf("error setting `results`: %s", err)
	}

	if err := d.Set("total_count", jobs.GetTotalCount()); err != nil {
		return diag.Errorf("error setting `total_count`: %s", err)
	}

	d.SetId(id.UniqueId())

	return nil
}

func flattenCloudBackupSnapshotExportJobs(jobs []admin.DiskBackupExportJob) []map[string]any {
	var results []map[string]any

	if len(jobs) == 0 {
		return results
	}

	results = make([]map[string]any, len(jobs))

	for k, job := range jobs {
		results[k] = map[string]any{
			"export_job_id":                      job.GetId(),
			"created_at":                         conversion.TimePtrToStringPtr(job.CreatedAt),
			"components":                         flattenExportJobsComponents(job.GetComponents()),
			"custom_data":                        flattenExportJobsCustomData(job.GetCustomData()),
			"export_bucket_id":                   job.GetExportBucketId(),
			"export_status_exported_collections": job.ExportStatus.GetExportedCollections(),
			"export_status_total_collections":    job.ExportStatus.GetTotalCollections(),
			"finished_at":                        conversion.TimePtrToStringPtr(job.FinishedAt),
			"prefix":                             job.GetPrefix(),
			"snapshot_id":                        job.GetSnapshotId(),
			"state":                              job.GetState(),
		}
	}

	return results
}
