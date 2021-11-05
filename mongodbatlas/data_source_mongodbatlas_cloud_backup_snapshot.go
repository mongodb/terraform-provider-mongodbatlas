package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasCloudBackupSnapshot() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasCloudBackupSnapshotRead,
		Schema: map[string]*schema.Schema{
			"snapshot_id": {
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
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expires_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"master_key_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mongod_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_size_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloud_provider": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"members": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloud_provider": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
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
			"replica_set_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceMongoDBAtlasCloudBackupSnapshotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	requestParameters := &matlas.SnapshotReqPathParameters{
		SnapshotID:  d.Get("snapshot_id").(string),
		GroupID:     d.Get("project_id").(string),
		ClusterName: d.Get("cluster_name").(string),
	}

	snapshot, _, err := conn.CloudProviderSnapshots.GetOneCloudProviderSnapshot(ctx, requestParameters)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting cloudProviderSnapshot Information: %s", err))
	}

	if err = d.Set("created_at", snapshot.CreatedAt); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `created_at` for cloudProviderSnapshot (%s): %s", d.Id(), err))
	}

	if err = d.Set("description", snapshot.Description); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `description` for cloudProviderSnapshot (%s): %s", d.Id(), err))
	}

	if err = d.Set("expires_at", snapshot.ExpiresAt); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `expires_at` for cloudProviderSnapshot (%s): %s", d.Id(), err))
	}

	if err = d.Set("master_key_uuid", snapshot.MasterKeyUUID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `master_key_uuid` for cloudProviderSnapshot (%s): %s", d.Id(), err))
	}

	if err = d.Set("mongod_version", snapshot.MongodVersion); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `mongod_version` for cloudProviderSnapshot (%s): %s", d.Id(), err))
	}

	if err = d.Set("snapshot_type", snapshot.SnapshotType); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `snapshot_type` for cloudProviderSnapshot (%s): %s", d.Id(), err))
	}

	if err = d.Set("status", snapshot.Status); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `status` for cloudProviderSnapshot (%s): %s", d.Id(), err))
	}

	if err = d.Set("storage_size_bytes", snapshot.StorageSizeBytes); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `storage_size_bytes` for cloudProviderSnapshot (%s): %s", d.Id(), err))
	}

	if err = d.Set("type", snapshot.Type); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `type` for cloudProviderSnapshot (%s): %s", d.Id(), err))
	}

	if err = d.Set("cloud_provider", snapshot.CloudProvider); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `cloud_provider` for snapshot (%s): %s", d.Id(), err))
	}

	if err = d.Set("members", flattenCloudMembers(snapshot.Members)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `members` for snapshot (%s): %s", d.Id(), err))
	}

	if err = d.Set("replica_set_name", snapshot.ReplicaSetName); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `replica_set_name` for snapshot (%s): %s", d.Id(), err))
	}

	if err = d.Set("snapshot_ids", snapshot.SnapshotsIds); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `snapshot_ids` for snapshot (%s): %s", d.Id(), err))
	}

	d.SetId(snapshot.ID)

	return nil
}
