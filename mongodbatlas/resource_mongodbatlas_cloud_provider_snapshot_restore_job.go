package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/spf13/cast"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func resourceMongoDBAtlasCloudProviderSnapshotRestoreJob() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBAtlasCloudProviderSnapshotRestoreJobCreate,
		ReadContext:   resourceMongoDBAtlasCloudProviderSnapshotRestoreJobRead,
		DeleteContext: resourceMongoDBAtlasCloudProviderSnapshotRestoreJobDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasCloudProviderSnapshotRestoreJobImportState,
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
			"snapshot_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"delivery_type": {
				Type:          schema.TypeMap,
				Optional:      true,
				ForceNew:      true,
				Deprecated:    "use delivery_type_config instead",
				ConflictsWith: []string{"delivery_type_config"},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(map[string]interface{})

					_, automated := v["automated"]
					_, download := v["download"]
					_, pointInTime := v["point_in_time"]

					if (v["automated"] == "true" && v["download"] == "true" && v["point_in_time"] == "true") ||
						(v["automated"] == "false" && v["download"] == "false" && v["point_in_time"] == "false") ||
						(!automated && !download && !pointInTime) {
						errs = append(errs, fmt.Errorf("%q you can only submit one type of restore job: automated, download or point_in_time", key))
					}
					if v["automated"] == "true" && (v["download"] == "false" || v["download"] == "" || !download) {
						if targetClusterName, ok := v["target_cluster_name"]; !ok || targetClusterName == "" {
							errs = append(errs, fmt.Errorf("%q target_cluster_name must be set", key))
						}
						if targetGroupID, ok := v["target_project_id"]; !ok || targetGroupID == "" {
							errs = append(errs, fmt.Errorf("%q target_project_id must be set", key))
						}
					}
					if v["download"] == "true" && (v["automated"] == "false" || v["automated"] == "" || !automated) &&
						(v["point_in_time"] == "false" || v["point_in_time"] == "" || !pointInTime) {
						if targetClusterName, ok := v["target_cluster_name"]; ok || targetClusterName == "" {
							errs = append(errs, fmt.Errorf("%q it's not necessary implement target_cluster_name when you are using download delivery type", key))
						}
						if targetGroupID, ok := v["target_project_id"]; ok || targetGroupID == "" {
							errs = append(errs, fmt.Errorf("%q it's not necessary implement target_project_id when you are using download delivery type", key))
						}
					}
					if v["point_in_time"] == "true" && (v["download"] == "false" || v["download"] == "" || !download) &&
						(v["automated"] == "false" || v["automated"] == "" || !automated) {
						_, oplogTS := v["oplog_ts"]
						_, pointTimeUTC := v["point_in_time_utc_seconds"]
						_, oplogInc := v["oplog_inc"]
						if targetClusterName, ok := v["target_cluster_name"]; !ok || targetClusterName == "" {
							errs = append(errs, fmt.Errorf("%q target_cluster_name must be set", key))
						}
						if targetGroupID, ok := v["target_project_id"]; !ok || targetGroupID == "" {
							errs = append(errs, fmt.Errorf("%q target_project_id must be set", key))
						}
						if !pointTimeUTC && !oplogTS && !oplogInc {
							errs = append(errs, fmt.Errorf("%q point_in_time_utc_seconds or oplog_ts and oplog_inc must be set", key))
						}
						if (oplogTS && !oplogInc) || (!oplogTS && oplogInc) {
							errs = append(errs, fmt.Errorf("%q if oplog_ts or oplog_inc is provided, oplog_inc and oplog_ts must be set", key))
						}
						if pointTimeUTC && (oplogTS || oplogInc) {
							errs = append(errs, fmt.Errorf("%q you can't use both point_in_time_utc_seconds and oplog_ts or oplog_inc", key))
						}
					}
					return
				},
			},
			"delivery_type_config": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"delivery_type"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"download": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"automated": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"point_in_time": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"target_cluster_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"target_project_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"oplog_ts": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"point_in_time_utc_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"oplog_inc": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"delivery_url": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cancelled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
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
			"timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_restore_job_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		DeprecationMessage: "this resource is deprecated, please transition as soon as possible to mongodbatlas_cloud_backup_snapshot_restore_job",
	}
}

func resourceMongoDBAtlasCloudProviderSnapshotRestoreJobCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas

	requestParameters := &matlas.SnapshotReqPathParameters{
		GroupID:     d.Get("project_id").(string),
		ClusterName: d.Get("cluster_name").(string),
	}

	err := validateDeliveryType(d.Get("delivery_type_config").([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	snapshotReq := buildRequestSnapshotReq(d)

	if _, ok := d.GetOk("delivery_type"); ok {
		deliveryType := "automated"
		if aut, _ := d.Get("delivery_type.download").(string); aut != "true" {
			deliveryType = "download"
		}

		if aut, _ := d.Get("delivery_type.point_in_time").(string); aut == "true" {
			deliveryType = "pointInTime"
		}

		snapshotReq = &matlas.CloudProviderSnapshotRestoreJob{
			SnapshotID:            getEncodedID(d.Get("snapshot_id").(string), "snapshot_id"),
			DeliveryType:          deliveryType,
			TargetClusterName:     d.Get("delivery_type.target_cluster_name").(string),
			TargetGroupID:         d.Get("delivery_type.target_project_id").(string),
			OplogTs:               cast.ToInt64(d.Get("delivery_type.oplog_ts")),
			OplogInc:              cast.ToInt64(d.Get("delivery_type.oplog_inc")),
			PointInTimeUTCSeconds: cast.ToInt64(d.Get("delivery_type.point_in_time_utc_seconds")),
		}
	}

	cloudProviderSnapshotRestoreJob, _, err := conn.CloudProviderSnapshotRestoreJobs.Create(ctx, requestParameters, snapshotReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error restore a snapshot: %s", err))
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":              d.Get("project_id").(string),
		"cluster_name":            d.Get("cluster_name").(string),
		"snapshot_restore_job_id": cloudProviderSnapshotRestoreJob.ID,
	}))

	return resourceMongoDBAtlasCloudProviderSnapshotRestoreJobRead(ctx, d, meta)
}

func resourceMongoDBAtlasCloudProviderSnapshotRestoreJobRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())

	requestParameters := &matlas.SnapshotReqPathParameters{
		JobID:       ids["snapshot_restore_job_id"],
		GroupID:     ids["project_id"],
		ClusterName: ids["cluster_name"],
	}

	snapshotReq, resp, err := conn.CloudProviderSnapshotRestoreJobs.Get(context.Background(), requestParameters)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf("error getting cloudProviderSnapshotRestoreJob Information: %s", err))
	}

	if err = d.Set("delivery_url", snapshotReq.DeliveryURL); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `delivery_url` for cloudProviderSnapshotRestoreJob (%s): %s", ids["snapshot_restore_job_id"], err))
	}

	if err = d.Set("cancelled", snapshotReq.Cancelled); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `cancelled` for cloudProviderSnapshotRestoreJob (%s): %s", ids["snapshot_restore_job_id"], err))
	}

	if err = d.Set("created_at", snapshotReq.CreatedAt); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `created_at` for cloudProviderSnapshotRestoreJob (%s): %s", ids["snapshot_restore_job_id"], err))
	}

	if err = d.Set("expired", snapshotReq.Expired); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `expired` for cloudProviderSnapshotRestoreJob (%s): %s", ids["snapshot_restore_job_id"], err))
	}

	if err = d.Set("expires_at", snapshotReq.ExpiresAt); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `expires_at` for cloudProviderSnapshotRestoreJob (%s): %s", ids["snapshot_restore_job_id"], err))
	}

	if err = d.Set("finished_at", snapshotReq.FinishedAt); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `Finished_at` for cloudProviderSnapshotRestoreJob (%s): %s", ids["snapshot_restore_job_id"], err))
	}

	if err = d.Set("timestamp", snapshotReq.Timestamp); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `timestamp` for cloudProviderSnapshotRestoreJob (%s): %s", ids["snapshot_restore_job_id"], err))
	}

	if err = d.Set("snapshot_restore_job_id", snapshotReq.ID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `snapshot_restore_job_id` for cloudProviderSnapshotRestoreJob (%s): %s", ids["snapshot_restore_job_id"], err))
	}

	return nil
}

func resourceMongoDBAtlasCloudProviderSnapshotRestoreJobDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())

	requestParameters := &matlas.SnapshotReqPathParameters{
		JobID:       ids["snapshot_restore_job_id"],
		GroupID:     ids["project_id"],
		ClusterName: ids["cluster_name"],
	}

	shouldDelete := true

	// Validate because atomated restore can not be cancelled
	if aut, _ := d.Get("delivery_type.automated").(string); aut == "true" {
		log.Print("Automated restore cannot be cancelled")
		shouldDelete = false
	}

	if aut, _ := d.Get("delivery_type_config.0.automated").(bool); aut {
		log.Print("Automated restore cannot be cancelled")
		shouldDelete = false
	}

	if shouldDelete {
		_, err := conn.CloudProviderSnapshotRestoreJobs.Delete(context.Background(), requestParameters)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error deleting a cloudProviderSnapshotRestoreJob (%s): %s", ids["snapshot_restore_job_id"], err))
		}
	}

	return nil
}

func resourceMongoDBAtlasCloudProviderSnapshotRestoreJobImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*MongoDBClient).Atlas

	projectID, clusterName, snapshotJobID, err := splitSnapshotRestoreJobImportID(d.Id())
	if err != nil {
		return nil, err
	}

	requestParameters := &matlas.SnapshotReqPathParameters{
		GroupID:     *projectID,
		ClusterName: *clusterName,
		JobID:       *snapshotJobID,
	}

	u, _, err := conn.CloudProviderSnapshotRestoreJobs.Get(ctx, requestParameters)
	if err != nil {
		return nil, fmt.Errorf("couldn't import cloudProviderSnapshotRestoreJob %s in project %s, error: %s", requestParameters.ClusterName, requestParameters.GroupID, err)
	}

	if err := d.Set("project_id", requestParameters.GroupID); err != nil {
		log.Printf("[WARN] Error setting project_id for (%s): %s", d.Id(), err)
	}

	if err := d.Set("cluster_name", requestParameters.ClusterName); err != nil {
		log.Printf("[WARN] Error setting cluster_name for (%s): %s", d.Id(), err)
	}

	if err := d.Set("snapshot_id", u.SnapshotID); err != nil {
		log.Printf("[WARN] Error setting snapshot_id for (%s): %s", d.Id(), err)
	}

	deliveryType := make(map[string]interface{})
	deliveryTypeConfig := make(map[string]interface{})

	if u.DeliveryType == "automated" {
		deliveryType["automated"] = "true"
		deliveryType["target_cluster_name"] = u.TargetClusterName
		deliveryType["target_project_id"] = u.TargetGroupID
		// For delivery_type_config
		deliveryTypeConfig["automated"] = true
		deliveryTypeConfig["target_cluster_name"] = u.TargetClusterName
		deliveryTypeConfig["target_project_id"] = u.TargetGroupID
	}

	if _, ok := d.GetOk("delivery_type"); ok {
		if err := d.Set("delivery_type", deliveryType); err != nil {
			log.Printf("[WARN] Error setting delivery_type for (%s): %s", d.Id(), err)
		}
	}

	if err := d.Set("delivery_type_config", []interface{}{deliveryTypeConfig}); err != nil {
		log.Printf("[WARN] Error setting delivery_type for (%s): %s", d.Id(), err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":              *projectID,
		"cluster_name":            *clusterName,
		"snapshot_restore_job_id": *snapshotJobID,
	}))

	return []*schema.ResourceData{d}, nil
}
