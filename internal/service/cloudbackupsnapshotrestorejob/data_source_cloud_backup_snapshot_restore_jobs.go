package cloudbackupsnapshotrestorejob

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312001/admin"
)

func PluralDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePluralRead,
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
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cancelled": {
							Type:     schema.TypeBool,
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
						"failed": {
							Type:     schema.TypeBool,
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
				},
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourcePluralRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2

	projectID := d.Get("project_id").(string)
	clusterName := d.Get("cluster_name").(string)
	pageNum := d.Get("page_num").(int)
	itermsPerPage := d.Get("items_per_page").(int)

	cloudProviderSnapshotRestoreJobs, _, err := conn.CloudBackupsApi.ListBackupRestoreJobs(ctx, projectID, clusterName).PageNum(pageNum).ItemsPerPage(itermsPerPage).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting cloudProviderSnapshotRestoreJobs information: %s", err))
	}

	if err := d.Set("results", flattenCloudProviderSnapshotRestoreJobs(cloudProviderSnapshotRestoreJobs.GetResults())); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `results`: %s", err))
	}

	if err := d.Set("total_count", cloudProviderSnapshotRestoreJobs.GetTotalCount()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `total_count`: %s", err))
	}

	d.SetId(id.UniqueId())

	return nil
}

func flattenCloudProviderSnapshotRestoreJobs(cloudProviderSnapshotRestoreJobs []admin.DiskBackupSnapshotRestoreJob) []map[string]any {
	var results []map[string]any

	if len(cloudProviderSnapshotRestoreJobs) > 0 {
		results = make([]map[string]any, len(cloudProviderSnapshotRestoreJobs))

		for k := range cloudProviderSnapshotRestoreJobs {
			cloudProviderSnapshotRestoreJob := cloudProviderSnapshotRestoreJobs[k]
			results[k] = map[string]any{
				"id":                        cloudProviderSnapshotRestoreJob.GetId(),
				"cancelled":                 cloudProviderSnapshotRestoreJob.GetCancelled(),
				"delivery_type":             cloudProviderSnapshotRestoreJob.GetDeliveryType(),
				"delivery_url":              cloudProviderSnapshotRestoreJob.GetDeliveryUrl(),
				"expired":                   cloudProviderSnapshotRestoreJob.GetExpired(),
				"expires_at":                conversion.TimePtrToStringPtr(cloudProviderSnapshotRestoreJob.ExpiresAt),
				"failed":                    cloudProviderSnapshotRestoreJob.GetFailed(),
				"finished_at":               conversion.TimePtrToStringPtr(cloudProviderSnapshotRestoreJob.FinishedAt),
				"snapshot_id":               cloudProviderSnapshotRestoreJob.GetSnapshotId(),
				"target_project_id":         cloudProviderSnapshotRestoreJob.GetTargetGroupId(),
				"target_cluster_name":       cloudProviderSnapshotRestoreJob.GetTargetClusterName(),
				"timestamp":                 conversion.TimePtrToStringPtr(cloudProviderSnapshotRestoreJob.Timestamp),
				"oplog_ts":                  cloudProviderSnapshotRestoreJob.GetOplogTs(),
				"point_in_time_utc_seconds": cloudProviderSnapshotRestoreJob.GetPointInTimeUTCSeconds(),
				"oplog_inc":                 cloudProviderSnapshotRestoreJob.GetOplogInc(),
			}
		}
	}

	return results
}
