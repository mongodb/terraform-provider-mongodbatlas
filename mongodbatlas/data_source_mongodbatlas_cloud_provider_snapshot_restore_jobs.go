package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasCloudProviderSnapshotRestoreJobs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasCloudProviderSnapshotRestoreJobsRead,
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
				},
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceMongoDBAtlasCloudProviderSnapshotRestoreJobsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	requestParameters := &matlas.SnapshotReqPathParameters{
		GroupID:     d.Get("project_id").(string),
		ClusterName: d.Get("cluster_name").(string),
	}
	options := &matlas.ListOptions{
		PageNum:      d.Get("page_num").(int),
		ItemsPerPage: d.Get("items_per_page").(int),
	}

	cloudProviderSnapshotRestoreJobs, _, err := conn.CloudProviderSnapshotRestoreJobs.List(ctx, requestParameters, options)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting cloudProviderSnapshotRestoreJobs information: %s", err))
	}

	if err := d.Set("results", flattenCloudProviderSnapshotRestoreJobs(cloudProviderSnapshotRestoreJobs.Results)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `results`: %s", err))
	}

	if err := d.Set("total_count", cloudProviderSnapshotRestoreJobs.TotalCount); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `total_count`: %s", err))
	}

	d.SetId(resource.UniqueId())

	return nil
}

func flattenCloudProviderSnapshotRestoreJobs(cloudProviderSnapshotRestoreJobs []*matlas.CloudProviderSnapshotRestoreJob) []map[string]interface{} {
	var results []map[string]interface{}

	if len(cloudProviderSnapshotRestoreJobs) > 0 {
		results = make([]map[string]interface{}, len(cloudProviderSnapshotRestoreJobs))

		for k, cloudProviderSnapshotRestoreJob := range cloudProviderSnapshotRestoreJobs {
			results[k] = map[string]interface{}{
				"id":                        cloudProviderSnapshotRestoreJob.ID,
				"cancelled":                 cloudProviderSnapshotRestoreJob.Cancelled,
				"created_at":                cloudProviderSnapshotRestoreJob.CreatedAt,
				"delivery_type":             cloudProviderSnapshotRestoreJob.DeliveryType,
				"delivery_url":              cloudProviderSnapshotRestoreJob.DeliveryURL,
				"expired":                   cloudProviderSnapshotRestoreJob.Expired,
				"expires_at":                cloudProviderSnapshotRestoreJob.ExpiresAt,
				"finished_at":               cloudProviderSnapshotRestoreJob.FinishedAt,
				"snapshot_id":               cloudProviderSnapshotRestoreJob.SnapshotID,
				"target_project_id":         cloudProviderSnapshotRestoreJob.TargetGroupID,
				"target_cluster_name":       cloudProviderSnapshotRestoreJob.TargetClusterName,
				"timestamp":                 cloudProviderSnapshotRestoreJob.Timestamp,
				"oplog_ts":                  cloudProviderSnapshotRestoreJob.OplogTs,
				"point_in_time_utc_seconds": cloudProviderSnapshotRestoreJob.PointInTimeUTCSeconds,
				"oplog_inc":                 cloudProviderSnapshotRestoreJob.OplogInc,
			}
		}
	}

	return results
}
