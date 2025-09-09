package cloudbackupsnapshot

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
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

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	groupID := d.Get("project_id").(string)
	clusterName := d.Get("cluster_name").(string)
	snapshotID := d.Get("snapshot_id").(string)

	snapshot, _, err := connV2.CloudBackupsApi.GetClusterBackupSnapshot(ctx, groupID, clusterName, snapshotID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting cloudProviderSnapshot Information: %s", err))
	}

	if err = d.Set("created_at", conversion.TimePtrToStringPtr(snapshot.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `created_at` for cloudProviderSnapshot (%s): %s", d.Id(), err))
	}

	if err = d.Set("description", snapshot.GetDescription()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `description` for cloudProviderSnapshot (%s): %s", d.Id(), err))
	}

	if err = d.Set("expires_at", conversion.TimePtrToStringPtr(snapshot.ExpiresAt)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `expires_at` for cloudProviderSnapshot (%s): %s", d.Id(), err))
	}

	if err = d.Set("master_key_uuid", snapshot.GetMasterKeyUUID()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `master_key_uuid` for cloudProviderSnapshot (%s): %s", d.Id(), err))
	}

	if err = d.Set("mongod_version", snapshot.GetMongodVersion()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `mongod_version` for cloudProviderSnapshot (%s): %s", d.Id(), err))
	}

	if err = d.Set("snapshot_type", snapshot.GetSnapshotType()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `snapshot_type` for cloudProviderSnapshot (%s): %s", d.Id(), err))
	}

	if err = d.Set("status", snapshot.GetStatus()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `status` for cloudProviderSnapshot (%s): %s", d.Id(), err))
	}

	if err = d.Set("storage_size_bytes", snapshot.GetStorageSizeBytes()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `storage_size_bytes` for cloudProviderSnapshot (%s): %s", d.Id(), err))
	}

	if err = d.Set("type", snapshot.GetType()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `type` for cloudProviderSnapshot (%s): %s", d.Id(), err))
	}

	if err = d.Set("cloud_provider", snapshot.GetCloudProvider()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `cloud_provider` for snapshot (%s): %s", d.Id(), err))
	}

	if err = d.Set("replica_set_name", snapshot.GetReplicaSetName()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `replica_set_name` for snapshot (%s): %s", d.Id(), err))
	}

	sharded, _, _ := connV2.CloudBackupsApi.GetBackupShardedCluster(ctx, groupID, clusterName, snapshotID).Execute()
	if sharded != nil {
		if err = d.Set("members", flattenCloudMembers(sharded.GetMembers())); err != nil {
			return diag.FromErr(fmt.Errorf("error setting `members` for snapshot (%s): %s", d.Id(), err))
		}
		if err = d.Set("snapshot_ids", sharded.GetSnapshotIds()); err != nil {
			return diag.FromErr(fmt.Errorf("error setting `snapshot_ids` for snapshot (%s): %s", d.Id(), err))
		}
	}
	d.SetId(snapshot.GetId())
	return nil
}
