package cloudbackupsnapshotrestorejob

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"snapshot_restore_job_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
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
				ForceNew: true,
				Optional: true,
				// When deprecating, change snapshot_restore_job_id to Required: true and implementation below
				Deprecated: fmt.Sprintf(constant.DeprecationParamByVersion, "1.18.0") + " Use snapshot_restore_job_id instead.",
			},
			"cancelled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"created_at": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: fmt.Sprintf(constant.DeprecationParamByVersion, "1.18.0"),
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

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2

	var restoreID string
	var err error
	restoreIDRaw, restoreIDInField := d.GetOk("snapshot_restore_job_id")
	if restoreIDInField {
		restoreID = restoreIDRaw.(string)
	} else {
		restoreID = conversion.GetEncodedID(d.Get("job_id").(string), "snapshot_restore_job_id")
		if err = d.Set("snapshot_restore_job_id", restoreID); err != nil {
			return diag.FromErr(fmt.Errorf("error setting `snapshot_restore_job_id` for cloudProviderSnapshotRestoreJob (%s): %s", d.Id(), err))
		}
	}
	projectID := d.Get("project_id").(string)
	clusterName := d.Get("cluster_name").(string)

	snapshotRes, _, err := conn.CloudBackupsApi.GetBackupRestoreJob(ctx, projectID, clusterName, restoreID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting cloudProviderSnapshotRestoreJob Information: %s", err))
	}

	if err = d.Set("delivery_type", snapshotRes.GetDeliveryType()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `delivery_type` for cloudProviderSnapshotRestoreJob (%s): %s", d.Id(), err))
	}

	if err = d.Set("snapshot_id", snapshotRes.GetSnapshotId()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `snapshotId` for cloudProviderSnapshotRestoreJob (%s): %s", d.Id(), err))
	}

	if err = d.Set("target_project_id", snapshotRes.GetTargetGroupId()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `targetGroupId` for cloudProviderSnapshotRestoreJob (%s): %s", d.Id(), err))
	}

	if err = d.Set("target_cluster_name", snapshotRes.GetTargetClusterName()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `targetClusterName` for cloudProviderSnapshotRestoreJob (%s): %s", d.Id(), err))
	}

	if err = d.Set("oplog_ts", snapshotRes.GetOplogTs()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `oplog_ts` for cloudProviderSnapshotRestoreJob (%s): %s", d.Id(), err))
	}

	if err = d.Set("point_in_time_utc_seconds", snapshotRes.GetPointInTimeUTCSeconds()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `point_in_time_utc_seconds` for cloudProviderSnapshotRestoreJob (%s): %s", d.Id(), err))
	}

	if err = d.Set("oplog_inc", snapshotRes.GetOplogInc()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `oplog_inc` for cloudProviderSnapshotRestoreJob (%s): %s", d.Id(), err))
	}

	d.SetId(snapshotRes.GetId())

	return nil
}
