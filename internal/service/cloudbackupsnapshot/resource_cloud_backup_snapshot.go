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
	"go.mongodb.org/atlas-sdk/v20231115006/admin"
	matlas "go.mongodb.org/atlas/mongodbatlas"
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
	// Get client connection.
	conn := meta.(*config.MongoDBClient).Atlas

	requestParameters := &matlas.SnapshotReqPathParameters{
		GroupID:     d.Get("project_id").(string),
		ClusterName: d.Get("cluster_name").(string),
	}

	snapshotReq := &matlas.CloudProviderSnapshot{
		Description:     d.Get("description").(string),
		RetentionInDays: d.Get("retention_in_days").(int),
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"CREATING", "UPDATING", "REPAIRING", "REPEATING"},
		Target:     []string{"IDLE"},
		Refresh:    advancedcluster.ResourceClusterRefreshFunc(ctx, d.Get("cluster_name").(string), d.Get("project_id").(string), advancedcluster.ServiceFromClient(conn)),
		Timeout:    15 * time.Minute,
		MinTimeout: 30 * time.Second,
	}

	// Wait, catching any errors
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	snapshot, _, err := conn.CloudProviderSnapshots.Create(ctx, requestParameters, snapshotReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error taking a snapshot: %s", err))
	}

	requestParameters.SnapshotID = snapshot.ID

	stateConf = &retry.StateChangeConf{
		Pending:    []string{"queued", "inProgress"},
		Target:     []string{"completed", "failed"},
		Refresh:    resourceRefreshFunc(ctx, requestParameters, conn),
		Timeout:    1 * time.Hour,
		MinTimeout: 60 * time.Second,
		Delay:      1 * time.Minute,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":   d.Get("project_id").(string),
		"cluster_name": d.Get("cluster_name").(string),
		"snapshot_id":  snapshot.ID,
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
	// Get client connection.
	conn := meta.(*config.MongoDBClient).Atlas
	ids := conversion.DecodeStateID(d.Id())

	requestParameters := &matlas.SnapshotReqPathParameters{
		SnapshotID:  ids["snapshot_id"],
		GroupID:     ids["project_id"],
		ClusterName: ids["cluster_name"],
	}

	_, err := conn.CloudProviderSnapshots.Delete(ctx, requestParameters)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting a snapshot (%s): %s", ids["snapshot_id"], err))
	}

	return nil
}

func resourceImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	conn := meta.(*config.MongoDBClient).Atlas

	requestParameters, err := SplitSnapshotImportID(d.Id())
	if err != nil {
		return nil, err
	}

	u, _, err := conn.CloudProviderSnapshots.GetOneCloudProviderSnapshot(ctx, requestParameters)
	if err != nil {
		return nil, fmt.Errorf("couldn't import snapshot %s in project %s, error: %s", requestParameters.ClusterName, requestParameters.GroupID, err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":   requestParameters.GroupID,
		"cluster_name": requestParameters.ClusterName,
		"snapshot_id":  requestParameters.SnapshotID,
	}))

	if err := d.Set("project_id", requestParameters.GroupID); err != nil {
		log.Printf("[WARN] Error setting project_id for (%s): %s", requestParameters.SnapshotID, err)
	}

	if err := d.Set("cluster_name", requestParameters.ClusterName); err != nil {
		log.Printf("[WARN] Error setting cluster_name for (%s): %s", requestParameters.SnapshotID, err)
	}

	if err := d.Set("description", u.Description); err != nil {
		log.Printf("[WARN] Error setting description for (%s): %s", requestParameters.SnapshotID, err)
	}

	return []*schema.ResourceData{d}, nil
}

func resourceRefreshFunc(ctx context.Context, requestParameters *matlas.SnapshotReqPathParameters, client *matlas.Client) retry.StateRefreshFunc {
	return func() (any, string, error) {
		c, resp, err := client.CloudProviderSnapshots.GetOneCloudProviderSnapshot(ctx, requestParameters)

		switch {
		case err != nil:
			return nil, "failed", err
		case resp.StatusCode == http.StatusNotFound:
			return "", "DELETED", nil
		case c.Status == "failed":
			return nil, c.Status, fmt.Errorf("error creating MongoDB snapshot(%s) status was: %s", requestParameters.SnapshotID, c.Status)
		}

		if c.Status != "" {
			log.Printf("[DEBUG] status for MongoDB snapshot: %s: %s", requestParameters.SnapshotID, c.Status)
		}

		return c, c.Status, nil
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
