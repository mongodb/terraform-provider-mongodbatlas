package cloudbackupsnapshot

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/mwielbut/pointy"
	"go.mongodb.org/atlas-sdk/v20231115006/admin"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate,
		ReadContext:   resourceRead,
		DeleteContext: resourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceImport,
		},
		Schema: map[string]*schema.Schema{
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
			"description": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"retention_in_days": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"created_at": {
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
			"snapshot_id": {
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

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).Atlas
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	groupID := d.Get("project_id").(string)
	clusterName := d.Get("cluster_name").(string)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"CREATING", "UPDATING", "REPAIRING", "REPEATING"},
		Target:     []string{"IDLE"},
		Refresh:    advancedcluster.ResourceClusterRefreshFunc(ctx, d.Get("cluster_name").(string), d.Get("project_id").(string), advancedcluster.ServiceFromClient(conn)),
		Timeout:    15 * time.Minute,
		MinTimeout: 30 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	params := &admin.DiskBackupOnDemandSnapshotRequest{
		Description:     conversion.StringPtr(d.Get("description").(string)),
		RetentionInDays: pointy.Int(d.Get("retention_in_days").(int)),
	}
	snapshot, _, err := connV2.CloudBackupsApi.TakeSnapshot(ctx, groupID, clusterName, params).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error taking a snapshot: %s", err))
	}

	requestParams := &admin.GetReplicaSetBackupApiParams{
		GroupId:     groupID,
		ClusterName: clusterName,
		SnapshotId:  snapshot.GetId(),
	}

	stateConf = &retry.StateChangeConf{
		Pending:    []string{"queued", "inProgress"},
		Target:     []string{"completed", "failed"},
		Refresh:    resourceRefreshFunc(ctx, requestParams, connV2),
		Timeout:    1 * time.Hour,
		MinTimeout: 60 * time.Second,
		Delay:      1 * time.Minute,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":   groupID,
		"cluster_name": clusterName,
		"snapshot_id":  snapshot.GetId(),
	}))

	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	groupID := ids["project_id"]
	clusterName := ids["cluster_name"]
	snapshotID := ids["snapshot_id"]

	snapshot, resp, err := connV2.CloudBackupsApi.GetReplicaSetBackup(ctx, groupID, clusterName, snapshotID).Execute()
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error getting snapshot Information: %s", err))
	}

	if err = d.Set("snapshot_id", snapshot.GetId()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `snapshot_id` for snapshot (%s): %s", snapshotID, err))
	}

	if err = d.Set("created_at", conversion.TimePtrToStringPtr(snapshot.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `created_at` for snapshot (%s): %s", snapshotID, err))
	}

	if err = d.Set("expires_at", conversion.TimePtrToStringPtr(snapshot.ExpiresAt)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `expires_at` for snapshot (%s): %s", snapshotID, err))
	}

	if err = d.Set("master_key_uuid", snapshot.GetMasterKeyUUID()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `master_key_uuid` for snapshot (%s): %s", snapshotID, err))
	}

	if err = d.Set("mongod_version", snapshot.GetMongodVersion()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `mongod_version` for snapshot (%s): %s", snapshotID, err))
	}

	if err = d.Set("snapshot_type", snapshot.GetSnapshotType()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `snapshot_type` for snapshot (%s): %s", snapshotID, err))
	}

	if err = d.Set("status", snapshot.GetStatus()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `status` for snapshot (%s): %s", snapshotID, err))
	}

	if err = d.Set("storage_size_bytes", snapshot.GetStorageSizeBytes()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `storage_size_bytes` for snapshot (%s): %s", snapshotID, err))
	}

	if err = d.Set("type", snapshot.GetType()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `type` for snapshot (%s): %s", snapshotID, err))
	}

	if err = d.Set("cloud_provider", snapshot.GetCloudProvider()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `cloud_provider` for snapshot (%s): %s", snapshotID, err))
	}

	if err = d.Set("replica_set_name", snapshot.GetReplicaSetName()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `replica_set_name` for snapshot (%s): %s", snapshotID, err))
	}

	sharded, _, _ := connV2.CloudBackupsApi.GetShardedClusterBackup(ctx, groupID, clusterName, snapshotID).Execute()
	if sharded != nil {
		if err = d.Set("members", flattenCloudMembers(sharded.GetMembers())); err != nil {
			return diag.FromErr(fmt.Errorf("error setting `members` for snapshot (%s): %s", snapshotID, err))
		}

		if err = d.Set("snapshot_ids", sharded.GetSnapshotIds()); err != nil {
			return diag.FromErr(fmt.Errorf("error setting `snapshot_ids` for snapshot (%s): %s", snapshotID, err))
		}
	}
	return nil
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	groupID := ids["project_id"]
	clusterName := ids["cluster_name"]
	snapshotID := ids["snapshot_id"]
	_, _, err := connV2.CloudBackupsApi.DeleteReplicaSetBackup(ctx, groupID, clusterName, snapshotID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting a snapshot (%s): %s", snapshotID, err))
	}
	return nil
}

func resourceImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	params, err := SplitSnapshotImportID(d.Id())
	if err != nil {
		return nil, err
	}

	snapshot, _, err := connV2.CloudBackupsApi.GetReplicaSetBackupWithParams(ctx, params).Execute()
	if err != nil {
		return nil, fmt.Errorf("couldn't import snapshot %s in project %s, error: %s", params.ClusterName, params.GroupId, err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":   params.GroupId,
		"cluster_name": params.ClusterName,
		"snapshot_id":  params.SnapshotId,
	}))

	if err := d.Set("project_id", params.GroupId); err != nil {
		log.Printf("[WARN] Error setting project_id for (%s): %s", params.SnapshotId, err)
	}

	if err := d.Set("cluster_name", params.ClusterName); err != nil {
		log.Printf("[WARN] Error setting cluster_name for (%s): %s", params.SnapshotId, err)
	}

	if err := d.Set("description", snapshot.GetDescription()); err != nil {
		log.Printf("[WARN] Error setting description for (%s): %s", params.SnapshotId, err)
	}

	return []*schema.ResourceData{d}, nil
}

func resourceRefreshFunc(ctx context.Context, requestParams *admin.GetReplicaSetBackupApiParams, connV2 *admin.APIClient) retry.StateRefreshFunc {
	return func() (any, string, error) {
		snapshot, resp, err := connV2.CloudBackupsApi.GetReplicaSetBackupWithParams(ctx, requestParams).Execute()
		if err != nil {
			return nil, "failed", err
		}
		if resp.StatusCode == http.StatusNotFound {
			return "", "DELETED", nil
		}
		status := snapshot.GetStatus()
		if status == "failed" {
			return nil, status, fmt.Errorf("error creating MongoDB snapshot(%s) status was: %s", requestParams.SnapshotId, status)
		}
		return snapshot, status, nil
	}
}

func flattenCloudMembers(members []admin.DiskBackupShardedClusterSnapshotMember) []map[string]any {
	if len(members) == 0 {
		return nil
	}
	ret := make([]map[string]any, 0)
	for _, member := range members {
		ret = append(ret, map[string]any{
			"id":               member.GetId(),
			"cloud_provider":   member.GetCloudProvider(),
			"replica_set_name": member.GetReplicaSetName(),
		})
	}
	return ret
}
