package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasCloudProviderSnapshotRestoreJobs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasCloudProviderSnapshotRestoreJobsRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
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

func dataSourceMongoDBAtlasCloudProviderSnapshotRestoreJobsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	requestParameters := &matlas.SnapshotReqPathParameters{
		GroupID:     d.Get("project_id").(string),
		ClusterName: d.Get("cluster_name").(string),
	}

	cloudProviderSnapshotRestoreJobs, _, err := conn.CloudProviderSnapshotRestoreJobs.List(context.Background(), requestParameters)
	if err != nil {
		return fmt.Errorf("error getting cloudProviderSnapshotRestoreJobs information: %s", err)
	}
	if err := d.Set("results", flattenCloudProviderSnapshotRestoreJobs(cloudProviderSnapshotRestoreJobs.Results)); err != nil {
		return fmt.Errorf("error setting `results`: %s", err)
	}
	if err := d.Set("total_count", cloudProviderSnapshotRestoreJobs.TotalCount); err != nil {
		return fmt.Errorf("error setting `total_count`: %s", err)
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
				"id":                  cloudProviderSnapshotRestoreJob.ID,
				"cancelled":           cloudProviderSnapshotRestoreJob.Cancelled,
				"created_at":          cloudProviderSnapshotRestoreJob.CreatedAt,
				"delivery_type":       cloudProviderSnapshotRestoreJob.DeliveryType,
				"delivery_url":        cloudProviderSnapshotRestoreJob.DeliveryURL,
				"expired":             cloudProviderSnapshotRestoreJob.Expired,
				"expires_at":          cloudProviderSnapshotRestoreJob.ExpiresAt,
				"finished_at":         cloudProviderSnapshotRestoreJob.FinishedAt,
				"snapshot_id":         cloudProviderSnapshotRestoreJob.SnapshotID,
				"target_project_id":   cloudProviderSnapshotRestoreJob.TargetGroupID,
				"target_cluster_name": cloudProviderSnapshotRestoreJob.TargetClusterName,
				"timestamp":           cloudProviderSnapshotRestoreJob.Timestamp,
			}
		}
	}
	return results
}
