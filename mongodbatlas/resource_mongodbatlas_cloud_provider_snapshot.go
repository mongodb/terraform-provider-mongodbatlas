package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func resourceMongoDBAtlasCloudProviderSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoDBAtlasCloudProviderSnapshotCreate,
		Read:   resourceMongoDBAtlasCloudProviderSnapshotRead,
		Delete: resourceMongoDBAtlasCloudProviderSnapshotDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMongoDBAtlasCloudProviderSnapshotImportState,
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
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
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
		},
	}
}

func resourceMongoDBAtlasCloudProviderSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())

	requestParameters := &matlas.SnapshotReqPathParameters{
		SnapshotID:  ids["snapshot_id"],
		GroupID:     ids["project_id"],
		ClusterName: ids["cluster_name"],
	}

	snapshotReq, _, err := conn.CloudProviderSnapshots.GetOneCloudProviderSnapshot(context.Background(), requestParameters)
	if err != nil {
		return fmt.Errorf("error getting snapshot Information: %s", err)
	}

	if err = d.Set("snapshot_id", snapshotReq.ID); err != nil {
		return fmt.Errorf("error setting `snapshot_id` for snapshot (%s): %s", ids["snapshot_id"], err)
	}

	if err = d.Set("created_at", snapshotReq.CreatedAt); err != nil {
		return fmt.Errorf("error setting `created_at` for snapshot (%s): %s", ids["snapshot_id"], err)
	}

	if err = d.Set("expires_at", snapshotReq.ExpiresAt); err != nil {
		return fmt.Errorf("error setting `expires_at` for snapshot (%s): %s", ids["snapshot_id"], err)
	}

	if err = d.Set("master_key_uuid", snapshotReq.MasterKeyUUID); err != nil {
		return fmt.Errorf("error setting `master_key_uuid` for snapshot (%s): %s", ids["snapshot_id"], err)
	}

	if err = d.Set("mongod_version", snapshotReq.MongodVersion); err != nil {
		return fmt.Errorf("error setting `mongod_version` for snapshot (%s): %s", ids["snapshot_id"], err)
	}

	if err = d.Set("snapshot_type", snapshotReq.SnapshotType); err != nil {
		return fmt.Errorf("error setting `snapshot_type` for snapshot (%s): %s", ids["snapshot_id"], err)
	}

	if err = d.Set("status", snapshotReq.Status); err != nil {
		return fmt.Errorf("error setting `status` for snapshot (%s): %s", ids["snapshot_id"], err)
	}

	if err = d.Set("storage_size_bytes", snapshotReq.StorageSizeBytes); err != nil {
		return fmt.Errorf("error setting `storage_size_bytes` for snapshot (%s): %s", ids["snapshot_id"], err)
	}

	if err = d.Set("type", snapshotReq.Type); err != nil {
		return fmt.Errorf("error setting `type` for snapshot (%s): %s", ids["snapshot_id"], err)
	}

	return nil
}

func resourceMongoDBAtlasCloudProviderSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*matlas.Client)

	requestParameters := &matlas.SnapshotReqPathParameters{
		GroupID:     d.Get("project_id").(string),
		ClusterName: d.Get("cluster_name").(string),
	}

	snapshotReq := &matlas.CloudProviderSnapshot{
		Description:     d.Get("description").(string),
		RetentionInDays: d.Get("retention_in_days").(int),
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATING", "UPDATING", "REPAIRING", "REPEATING"},
		Target:     []string{"IDLE"},
		Refresh:    resourceClusterRefreshFunc(d.Get("cluster_name").(string), d.Get("project_id").(string), conn),
		Timeout:    10 * time.Minute,
		MinTimeout: 10 * time.Second,
		Delay:      3 * time.Minute,
	}

	// Wait, catching any errors
	_, err := stateConf.WaitForState()
	if err != nil {
		return err
	}

	snapshot, _, err := conn.CloudProviderSnapshots.Create(context.Background(), requestParameters, snapshotReq)
	if err != nil {
		return fmt.Errorf("error taking a snapshot: %s", err)
	}

	requestParameters.SnapshotID = snapshot.ID

	stateConf = &resource.StateChangeConf{
		Pending:    []string{"queued", "inProgress"},
		Target:     []string{"completed", "failed"},
		Refresh:    resourceCloudProviderSnapshotRefreshFunc(requestParameters, conn),
		Timeout:    1 * time.Hour,
		MinTimeout: 60 * time.Second,
		Delay:      1 * time.Minute,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForState()
	if err != nil {
		return err
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":   d.Get("project_id").(string),
		"cluster_name": d.Get("cluster_name").(string),
		"snapshot_id":  snapshot.ID,
	}))

	return resourceMongoDBAtlasCloudProviderSnapshotRead(d, meta)
}

func resourceMongoDBAtlasCloudProviderSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())

	requestParameters := &matlas.SnapshotReqPathParameters{
		SnapshotID:  ids["snapshot_id"],
		GroupID:     ids["project_id"],
		ClusterName: ids["cluster_name"],
	}

	_, err := conn.CloudProviderSnapshots.Delete(context.Background(), requestParameters)
	if err != nil {
		return fmt.Errorf("error deleting a snapshot (%s): %s", ids["snapshot_id"], err)
	}

	return nil
}

func resourceCloudProviderSnapshotRefreshFunc(requestParameters *matlas.SnapshotReqPathParameters, client *matlas.Client) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		c, resp, err := client.CloudProviderSnapshots.GetOneCloudProviderSnapshot(context.Background(), requestParameters)

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

func resourceMongoDBAtlasCloudProviderSnapshotImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*matlas.Client)

	requestParameters, err := splitSnapshotImportID(d.Id())
	if err != nil {
		return nil, err
	}

	u, _, err := conn.CloudProviderSnapshots.GetOneCloudProviderSnapshot(context.Background(), requestParameters)
	if err != nil {
		return nil, fmt.Errorf("couldn't import snapshot %s in project %s, error: %s", requestParameters.ClusterName, requestParameters.GroupID, err)
	}

	d.SetId(encodeStateID(map[string]string{
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

func splitSnapshotImportID(id string) (*matlas.SnapshotReqPathParameters, error) {
	var re = regexp.MustCompile(`(?s)^([0-9a-fA-F]{24})-(.*)-([0-9a-fA-F]{24})$`)

	parts := re.FindStringSubmatch(id)

	if len(parts) != 4 {
		return nil, errors.New("import format error: to import a snapshot, use the format {project_id}-{cluster_name}-{snapshot_id}")
	}

	return &matlas.SnapshotReqPathParameters{
		GroupID:     parts[1],
		ClusterName: parts[2],
		SnapshotID:  parts[3],
	}, nil
}
