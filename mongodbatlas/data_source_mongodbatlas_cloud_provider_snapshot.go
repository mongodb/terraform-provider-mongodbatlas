package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasCloudProviderSnapshot() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasCloudProviderSnapshotRead,
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
		},
	}
}

func dataSourceMongoDBAtlasCloudProviderSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	requestParameters := &matlas.SnapshotReqPathParameters{
		SnapshotID:  d.Get("snapshot_id").(string),
		GroupID:     d.Get("project_id").(string),
		ClusterName: d.Get("cluster_name").(string),
	}

	snapshotRes, _, err := conn.CloudProviderSnapshots.GetOneCloudProviderSnapshot(context.Background(), requestParameters)
	if err != nil {
		return fmt.Errorf("error getting cloudProviderSnapshot Information: %s", err)
	}

	if err = d.Set("created_at", snapshotRes.CreatedAt); err != nil {
		return fmt.Errorf("error setting `created_at` for cloudProviderSnapshot (%s): %s", d.Id(), err)
	}

	if err = d.Set("description", snapshotRes.Description); err != nil {
		return fmt.Errorf("error setting `description` for cloudProviderSnapshot (%s): %s", d.Id(), err)
	}

	if err = d.Set("expires_at", snapshotRes.ExpiresAt); err != nil {
		return fmt.Errorf("error setting `expires_at` for cloudProviderSnapshot (%s): %s", d.Id(), err)
	}

	if err = d.Set("master_key_uuid", snapshotRes.MasterKeyUUID); err != nil {
		return fmt.Errorf("error setting `master_key_uuid` for cloudProviderSnapshot (%s): %s", d.Id(), err)
	}

	if err = d.Set("mongod_version", snapshotRes.MongodVersion); err != nil {
		return fmt.Errorf("error setting `mongod_version` for cloudProviderSnapshot (%s): %s", d.Id(), err)
	}

	if err = d.Set("snapshot_type", snapshotRes.SnapshotType); err != nil {
		return fmt.Errorf("error setting `snapshot_type` for cloudProviderSnapshot (%s): %s", d.Id(), err)
	}

	if err = d.Set("status", snapshotRes.Status); err != nil {
		return fmt.Errorf("error setting `status` for cloudProviderSnapshot (%s): %s", d.Id(), err)
	}

	if err = d.Set("storage_size_bytes", snapshotRes.StorageSizeBytes); err != nil {
		return fmt.Errorf("error setting `storage_size_bytes` for cloudProviderSnapshot (%s): %s", d.Id(), err)
	}

	if err = d.Set("type", snapshotRes.Type); err != nil {
		return fmt.Errorf("error setting `type` for cloudProviderSnapshot (%s): %s", d.Id(), err)
	}

	d.SetId(snapshotRes.ID)

	return nil
}
