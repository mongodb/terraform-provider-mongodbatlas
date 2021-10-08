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
		DeprecationMessage: "This data source is deprecated. Please transition to mongodbatlas_cloud_backup_snapshot_restore_jobs as soon as possible",
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
