package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// This datasource does not have a resource: we tested it manually
func dataSourceMongoDBAtlasSharedTierSnapshot() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasSharedTierSnapshotRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
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
	}
}

func dataSourceMongoDBAtlasSharedTierSnapshotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).AtlasV2

	projectID := d.Get("project_id").(string)
	clusterName := d.Get("cluster_name").(string)
	snapshotID := d.Get("snapshot_id").(string)
	snapshot, resp, err := conn.SharedTierSnapshotsApi.GetSharedClusterBackup(ctx, projectID, clusterName, snapshotID).Execute()
	defer resp.Body.Close()

	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting shard-tier snapshot '%s': %w", snapshotID, err))
	}

	if err = d.Set("status", snapshot.Status); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `status` for shard-tier snapshot '%s': %w", snapshotID, err))
	}

	if err = d.Set("mongo_db_version", snapshot.MongoDBVersion); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `mongo_db_version` for shard-tier snapshot '%s': %w", snapshotID, err))
	}

	if err = d.Set("start_time", snapshot.StartTime.String()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `start_time` for shard-tier snapshot '%s': %w", snapshotID, err))
	}

	if err = d.Set("expiration", snapshot.Expiration.String()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `expiration` for shard-tier snapshot '%s': %w", snapshotID, err))
	}

	if err = d.Set("finish_time", snapshot.FinishTime.String()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `finish_time` for shard-tier snapshot '%s': %w", snapshotID, err))
	}

	if err = d.Set("scheduled_time", snapshot.ScheduledTime.String()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `scheduled_time` for shard-tier snapshot '%s': %w", snapshotID, err))
	}

	d.SetId(*snapshot.Id)
	return nil
}
