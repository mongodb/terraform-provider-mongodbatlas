package mongodbatlas

import (
	"context"
	"fmt"

	atlasSDK "go.mongodb.org/atlas-sdk/v20231001001/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// This datasource does not have a resource: we tested it manually
func dataSourceMongoDBAtlasSharedTierSnapshots() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasSharedTierSnapshotsRead,
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
						"snapshot_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mongo_db_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expiration": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"finish_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scheduled_time": {
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

func dataSourceMongoDBAtlasSharedTierSnapshotsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).AtlasV2

	projectID := d.Get("project_id").(string)
	clusterName := d.Get("cluster_name").(string)
	snapshots, _, err := conn.SharedTierSnapshotsApi.ListSharedClusterBackups(ctx, projectID, clusterName).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting shard-tier snapshots for cluster '%s': %w", clusterName, err))
	}

	if err := d.Set("results", flattenSharedTierSnapshots(snapshots.Results)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `results`: %w", err))
	}

	if err := d.Set("total_count", snapshots.TotalCount); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `total_count`: %w", err))
	}

	d.SetId(id.UniqueId())
	return nil
}

func flattenSharedTierSnapshots(sharedTierSnapshots []atlasSDK.BackupTenantSnapshot) []map[string]interface{} {
	if len(sharedTierSnapshots) == 0 {
		return nil
	}

	results := make([]map[string]interface{}, len(sharedTierSnapshots))
	for k, sharedTierSnapshot := range sharedTierSnapshots {
		results[k] = map[string]interface{}{
			"snapshot_id":      sharedTierSnapshot.Id,
			"mongo_db_version": sharedTierSnapshot.MongoDBVersion,
			"status":           sharedTierSnapshot.Status,
			"start_time":       sharedTierSnapshot.GetStartTime().String(),
			"finish_time":      sharedTierSnapshot.GetFinishTime().String(),
			"scheduled_time":   sharedTierSnapshot.GetScheduledTime().String(),
			"expiration":       sharedTierSnapshot.GetExpiration().String(),
		}
	}

	return results
}
