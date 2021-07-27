package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasCloudProviderSnapshotRestoreJob() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasCloudProviderSnapshotRestoreJobRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cancelled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"delivery_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"delivery_url": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"expired": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"expires_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"finished_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_cluster_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"oplog_ts": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"point_in_time_utc_seconds": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"oplog_inc": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceMongoDBAtlasCloudProviderSnapshotRestoreJobRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	requestParameters := &matlas.SnapshotReqPathParameters{
		JobID:       getEncodedID(d.Get("job_id").(string), "snapshot_restore_job_id"),
		GroupID:     d.Get("project_id").(string),
		ClusterName: d.Get("cluster_name").(string),
	}

	snapshotRes, _, err := conn.CloudProviderSnapshotRestoreJobs.Get(ctx, requestParameters)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting cloudProviderSnapshotRestoreJob Information: %s", err))
	}

	if err = d.Set("cancelled", snapshotRes.Cancelled); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `cancelled` for cloudProviderSnapshotRestoreJob (%s): %s", d.Id(), err))
	}

	if err = d.Set("created_at", snapshotRes.CreatedAt); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `created_at` for cloudProviderSnapshotRestoreJob (%s): %s", d.Id(), err))
	}

	if err = d.Set("delivery_type", snapshotRes.DeliveryType); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `delivery_type` for cloudProviderSnapshotRestoreJob (%s): %s", d.Id(), err))
	}

	if err = d.Set("delivery_url", snapshotRes.DeliveryURL); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `delivery_url` for cloudProviderSnapshotRestoreJob (%s): %s", d.Id(), err))
	}

	if err = d.Set("expired", snapshotRes.Expired); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `expired` for cloudProviderSnapshotRestoreJob (%s): %s", d.Id(), err))
	}

	if err = d.Set("expires_at", snapshotRes.ExpiresAt); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `expires_at` for cloudProviderSnapshotRestoreJob (%s): %s", d.Id(), err))
	}

	if err = d.Set("finished_at", snapshotRes.FinishedAt); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `finished_at` for cloudProviderSnapshotRestoreJob (%s): %s", d.Id(), err))
	}

	if err = d.Set("snapshot_id", snapshotRes.SnapshotID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `snapshotId` for cloudProviderSnapshotRestoreJob (%s): %s", d.Id(), err))
	}

	if err = d.Set("target_project_id", snapshotRes.TargetGroupID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `targetGroupId` for cloudProviderSnapshotRestoreJob (%s): %s", d.Id(), err))
	}

	if err = d.Set("target_cluster_name", snapshotRes.TargetClusterName); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `targetClusterName` for cloudProviderSnapshotRestoreJob (%s): %s", d.Id(), err))
	}

	if err = d.Set("timestamp", snapshotRes.Timestamp); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `timestamp` for cloudProviderSnapshotRestoreJob (%s): %s", d.Id(), err))
	}

	if err = d.Set("oplog_ts", snapshotRes.OplogTs); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `oplog_ts` for cloudProviderSnapshotRestoreJob (%s): %s", d.Id(), err))
	}

	if err = d.Set("point_in_time_utc_seconds", snapshotRes.PointInTimeUTCSeconds); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `point_in_time_utc_seconds` for cloudProviderSnapshotRestoreJob (%s): %s", d.Id(), err))
	}

	if err = d.Set("oplog_inc", snapshotRes.OplogInc); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `oplog_inc` for cloudProviderSnapshotRestoreJob (%s): %s", d.Id(), err))
	}

	d.SetId(snapshotRes.ID)

	return nil
}
