package mongodbatlas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceMongoDBAtlasCloudBackupSnapshotExportJob() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasCloudBackupSnapshotsExportJobRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"export_job_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
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
			"err_msg": {
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
	}
}

func dataSourceMongoDBAtlasCloudBackupSnapshotsExportJobRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]
	exportJobID := ids["export_job_id"]

	exportJob, _, err := conn.CloudProviderSnapshotExportJobs.Get(ctx, projectID, clusterName, exportJobID)
	if err != nil {
		return diag.Errorf("error getting snapshot export job information: %s", err)
	}

	if err := d.Set("export_job_id", exportJob.ID); err != nil {
		return diag.Errorf("error setting `export_job_id` for snapshot export job (%s): %s", d.Id(), err)
	}

	if err := d.Set("snapshot_id", exportJob.SnapshotID); err != nil {
		return diag.Errorf("error setting `snapshot_id` for snapshot export job (%s): %s", d.Id(), err)
	}

	if err := d.Set("custom_data", flattenExportJobsCustomData(exportJob.CustomData)); err != nil {
		return diag.Errorf("error setting `custom_data` for snapshot export job (%s): %s", d.Id(), err)
	}

	if err := d.Set("components", flattenExportJobsComponents(exportJob.Components)); err != nil {
		return diag.Errorf("error setting `components` for snapshot export job (%s): %s", d.Id(), err)
	}

	if err := d.Set("created_at", exportJob.CreatedAt); err != nil {
		return diag.Errorf("error setting `created_at` for snapshot export job (%s): %s", d.Id(), err)
	}

	if err := d.Set("err_msg", exportJob.ErrMsg); err != nil {
		return diag.Errorf("error setting `created_at` for snapshot export job (%s): %s", d.Id(), err)
	}

	if err := d.Set("export_bucket_id", exportJob.ExportBucketID); err != nil {
		return diag.Errorf("error setting `created_at` for snapshot export job (%s): %s", d.Id(), err)
	}

	if exportJob.ExportStatus != nil {
		if err := d.Set("export_status_exported_collections", exportJob.ExportStatus.ExportedCollections); err != nil {
			return diag.Errorf("error setting `export_status_exported_collections` for snapshot export job (%s): %s", d.Id(), err)
		}

		if err := d.Set("export_status_total_collections", exportJob.ExportStatus.TotalCollections); err != nil {
			return diag.Errorf("error setting `export_status_total_collections` for snapshot export job (%s): %s", d.Id(), err)
		}
	}

	if err := d.Set("finished_at", exportJob.FinishedAt); err != nil {
		return diag.Errorf("error setting `finished_at` for snapshot export job (%s): %s", d.Id(), err)
	}

	if err := d.Set("prefix", exportJob.Prefix); err != nil {
		return diag.Errorf("error setting `prefix` for snapshot export job (%s): %s", d.Id(), err)
	}

	if err := d.Set("state", exportJob.State); err != nil {
		return diag.Errorf("error setting `prefix` for snapshot export job (%s): %s", d.Id(), err)
	}

	d.SetId(exportJob.ID)

	return nil
}
