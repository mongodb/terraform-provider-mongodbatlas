package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// This datasource does not have a resource: we tested it manually
func dataSourceMongoDBAtlasCloudSharedTierRestoreJob() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasCloudSharedTierRestoreJobsRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_deployment_item_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_finished_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"restore_scheduled_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"restore_finished_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"delivery_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expiration_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceMongoDBAtlasCloudSharedTierRestoreJobsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).AtlasV2

	jobID := d.Get("job_id").(string)
	projectID := d.Get("project_id").(string)
	clusterName := d.Get("cluster_name").(string)

	job, _, err := conn.SharedTierRestoreJobsApi.GetSharedClusterBackupRestoreJob(ctx, clusterName, projectID, jobID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting shared tier restore job '%s': %w", jobID, err))
	}

	if err = d.Set("status", job.Status); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `status` for shared tier restore job '%s': %w", jobID, err))
	}

	if err = d.Set("target_project_id", job.TargetProjectId); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `target_project_id` for shared tier restore job '%s': %w", jobID, err))
	}

	if err = d.Set("target_deployment_item_name", job.TargetDeploymentItemName); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `target_deployment_item_name` for shared tier restore job '%s': %w", jobID, err))
	}

	if err = d.Set("snapshot_url", job.SnapshotUrl); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `snapshot_url` for shared tier restore job '%s': %w", jobID, err))
	}

	if err = d.Set("snapshot_id", job.SnapshotId); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `snapshot_id` for shared tier restore job '%s': %w", jobID, err))
	}

	if err = d.Set("snapshot_finished_date", job.GetSnapshotFinishedDate().String()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `snapshot_finished_date` for shared tier restore job '%s': %w", jobID, err))
	}

	if err = d.Set("restore_scheduled_date", job.GetRestoreScheduledDate().String()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `restore_scheduled_date` for shared tier restore job '%s': %w", jobID, err))
	}

	if err = d.Set("restore_finished_date", job.GetRestoreFinishedDate().String()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `restore_finished_date` for shared tier restore job '%s': %w", jobID, err))
	}

	if err = d.Set("delivery_type", job.DeliveryType); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `delivery_type` for shared tier restore job '%s': %w", jobID, err))
	}

	if err = d.Set("expiration_date", job.GetExpirationDate().String()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `expiration_date` for shared tier restore job '%s': %w", jobID, err))
	}

	d.SetId(*job.Id)

	return nil
}
