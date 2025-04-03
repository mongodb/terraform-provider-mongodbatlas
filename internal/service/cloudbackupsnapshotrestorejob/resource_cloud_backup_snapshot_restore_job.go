package cloudbackupsnapshotrestorejob

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312002/admin"
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
			"snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"delivery_type_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				ForceNew: true,
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
			"expired": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"expires_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"failed": {
				Type:     schema.TypeBool,
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
	}
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2

	projectID := d.Get("project_id").(string)
	clusterName := d.Get("cluster_name").(string)
	err := validateDeliveryType(d.Get("delivery_type_config").([]any))
	if err != nil {
		return diag.FromErr(err)
	}

	snapshotReq := buildRequestSnapshotReq(d)

	cloudProviderSnapshotRestoreJob, _, err := conn.CloudBackupsApi.CreateBackupRestoreJob(ctx, projectID, clusterName, snapshotReq).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error restore a snapshot: %s", err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":              d.Get("project_id").(string),
		"cluster_name":            d.Get("cluster_name").(string),
		"snapshot_restore_job_id": cloudProviderSnapshotRestoreJob.GetId(),
	}))

	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())

	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]
	restoreID := ids["snapshot_restore_job_id"]
	snapshotReq, resp, err := conn.CloudBackupsApi.GetBackupRestoreJob(ctx, projectID, clusterName, restoreID).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf("error getting cloudProviderSnapshotRestoreJob Information: %s", err))
	}

	if err = d.Set("snapshot_restore_job_id", snapshotReq.GetId()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `snapshot_restore_job_id` for cloudProviderSnapshotRestoreJob (%s): %s", restoreID, err))
	}

	return setCommonFields(d, snapshotReq, restoreID)
}

func setCommonFields(d *schema.ResourceData, snapshotReq *admin.DiskBackupSnapshotRestoreJob, restoreID string) diag.Diagnostics {
	var err error

	if err = d.Set("delivery_url", snapshotReq.GetDeliveryUrl()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `delivery_url` for cloudProviderSnapshotRestoreJob (%s): %s", restoreID, err))
	}
	if err = d.Set("cancelled", snapshotReq.GetCancelled()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `cancelled` for cloudProviderSnapshotRestoreJob (%s): %s", restoreID, err))
	}
	if err = d.Set("expired", snapshotReq.GetExpired()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `expired` for cloudProviderSnapshotRestoreJob (%s): %s", restoreID, err))
	}

	if err = d.Set("expires_at", conversion.TimePtrToStringPtr(snapshotReq.ExpiresAt)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `expires_at` for cloudProviderSnapshotRestoreJob (%s): %s", restoreID, err))
	}

	if err = d.Set("failed", snapshotReq.GetFailed()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `failed` for cloudProviderSnapshotRestoreJob (%s): %s", restoreID, err))
	}

	if err = d.Set("finished_at", conversion.TimePtrToStringPtr(snapshotReq.FinishedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `Finished_at` for cloudProviderSnapshotRestoreJob (%s): %s", restoreID, err))
	}

	if err = d.Set("timestamp", conversion.TimePtrToStringPtr(snapshotReq.Timestamp)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `timestamp` for cloudProviderSnapshotRestoreJob (%s): %s", restoreID, err))
	}

	return nil
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())

	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]
	restoreID := ids["snapshot_restore_job_id"]

	shouldDelete := true

	// Validate because automated restore can not be cancelled
	if aut, _ := d.Get("delivery_type_config.0.automated").(bool); aut {
		log.Print("Automated restore cannot be cancelled")
		shouldDelete = false
	}

	if aut, ok := d.Get("delivery_type_config.0.point_in_time").(bool); ok && aut {
		log.Print("Point in time restore cannot be cancelled")
		shouldDelete = false
	}

	if shouldDelete {
		_, _, err := conn.CloudBackupsApi.CancelBackupRestoreJob(ctx, projectID, clusterName, restoreID).Execute()
		if err != nil {
			return diag.FromErr(fmt.Errorf("error deleting a cloudProviderSnapshotRestoreJob (%s): %s", ids["snapshot_restore_job_id"], err))
		}
	}

	return nil
}

func resourceImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	conn := meta.(*config.MongoDBClient).AtlasV2

	projectID, clusterName, restoreID, err := splitSnapshotRestoreJobImportID(d.Id())
	if err != nil {
		return nil, err
	}
	u, _, err := conn.CloudBackupsApi.GetBackupRestoreJob(ctx, projectID, clusterName, restoreID).Execute()
	if err != nil {
		return nil, fmt.Errorf("couldn't import cloudProviderSnapshotRestoreJob %s in project %s, error: %s", clusterName, projectID, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		log.Printf("[WARN] Error setting project_id for (%s): %s", d.Id(), err)
	}

	if err := d.Set("cluster_name", clusterName); err != nil {
		log.Printf("[WARN] Error setting cluster_name for (%s): %s", d.Id(), err)
	}

	if err := d.Set("snapshot_id", u.GetSnapshotId()); err != nil {
		log.Printf("[WARN] Error setting snapshot_id for (%s): %s", d.Id(), err)
	}

	deliveryType := make(map[string]any)
	deliveryTypeConfig := make(map[string]any)

	if u.DeliveryType == "automated" {
		deliveryType["automated"] = "true"
		deliveryType["target_cluster_name"] = u.GetTargetClusterName()
		deliveryType["target_project_id"] = u.GetTargetGroupId()
		deliveryTypeConfig["automated"] = true
		deliveryTypeConfig["target_cluster_name"] = u.GetTargetClusterName()
		deliveryTypeConfig["target_project_id"] = u.GetTargetGroupId()
	}

	if err := d.Set("delivery_type_config", []any{deliveryTypeConfig}); err != nil {
		log.Printf("[WARN] Error setting delivery_type for (%s): %s", d.Id(), err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":              projectID,
		"cluster_name":            clusterName,
		"snapshot_restore_job_id": restoreID}))

	return []*schema.ResourceData{d}, nil
}

func splitSnapshotRestoreJobImportID(id string) (projectID, clusterName, snapshotJobID string, err error) {
	var re = regexp.MustCompile(`(?s)^([0-9a-fA-F]{24})-(.*)-([0-9a-fA-F]{24})$`)
	parts := re.FindStringSubmatch(id)

	if len(parts) != 4 {
		err = errors.New("import format error: to import a cloudProviderSnapshotRestoreJob, use the format {project_id}-{cluster_name}-{snapshot_restore_job_id}")
		return
	}

	projectID = parts[1]
	clusterName = parts[2]
	snapshotJobID = parts[3]

	return
}

func validateDeliveryType(dt []any) error {
	if len(dt) == 0 {
		return nil
	}

	v := dt[0].(map[string]any)
	key := "delivery_type_config"

	a, aOk := v["automated"]
	automated := aOk && a != nil && a.(bool)
	d, dOk := v["download"]
	download := dOk && d != nil && d.(bool)
	p, pOk := v["point_in_time"]
	pointInTime := pOk && p != nil && p.(bool)

	hasDeliveryType := automated || download || pointInTime

	if !hasDeliveryType ||
		(automated && download) ||
		(automated && pointInTime) ||
		(download && pointInTime) {
		return fmt.Errorf("%q you must submit exactly one type of restore job: automated, download or point_in_time", key)
	}

	if automated || pointInTime {
		if targetClusterName, ok := v["target_cluster_name"]; !ok || targetClusterName == "" {
			return fmt.Errorf("%q target_cluster_name must be set", key)
		}

		if targetProjectID, ok := v["target_project_id"]; !ok || targetProjectID == "" {
			return fmt.Errorf("%q target_project_id must be set", key)
		}
	} else {
		if targetClusterName, ok := v["target_cluster_name"]; ok && targetClusterName.(string) != "" {
			return fmt.Errorf("%q it's not necessary implement target_cluster_name when you are using download delivery type", key)
		}

		if targetProjectID, ok := v["target_project_id"]; ok && targetProjectID.(string) != "" {
			return fmt.Errorf("%q it's not necessary implement target_project_id when you are using download delivery type", key)
		}
	}

	if automated || download {
		return nil
	}

	pointTimeUTC, pointTimeUTCOk := v["point_in_time_utc_seconds"]
	isPITSet := pointTimeUTCOk && pointTimeUTC != nil && (pointTimeUTC.(int) > 0)
	oplogTS, oplogTSOk := v["oplog_ts"]
	isOpTSSet := oplogTSOk && oplogTS != nil && (oplogTS.(int) > 0)
	oplogInc, oplogIncOk := v["oplog_inc"]
	isOpIncSet := oplogIncOk && oplogInc != nil && (oplogInc.(int) > 0)

	if !isPITSet && (!isOpTSSet || !isOpIncSet) {
		return fmt.Errorf("%q point_in_time_utc_seconds or oplog_ts and oplog_inc must be set", key)
	}
	if isPITSet && (isOpTSSet || isOpIncSet) {
		return fmt.Errorf("%q you can't use both point_in_time_utc_seconds and oplog_ts or oplog_inc", key)
	}

	return nil
}

func buildRequestSnapshotReq(d *schema.ResourceData) *admin.DiskBackupSnapshotRestoreJob {
	if _, ok := d.GetOk("delivery_type_config"); ok {
		deliveryList := d.Get("delivery_type_config").([]any)

		delivery := deliveryList[0].(map[string]any)

		deliveryType := "automated"
		if aut, _ := delivery["download"].(bool); aut {
			deliveryType = "download"
		}

		if aut, _ := delivery["point_in_time"].(bool); aut {
			deliveryType = "pointInTime"
		}
		snapshotID := conversion.GetEncodedID(d.Get("snapshot_id").(string), "snapshot_id")
		return &admin.DiskBackupSnapshotRestoreJob{
			SnapshotId:            conversion.StringPtr(snapshotID),
			DeliveryType:          deliveryType,
			TargetClusterName:     conversion.StringPtr(delivery["target_cluster_name"].(string)),
			TargetGroupId:         conversion.StringPtr(delivery["target_project_id"].(string)),
			OplogTs:               conversion.IntPtr(delivery["oplog_ts"].(int)),
			OplogInc:              conversion.IntPtr(delivery["oplog_inc"].(int)),
			PointInTimeUTCSeconds: conversion.IntPtr(delivery["point_in_time_utc_seconds"].(int)),
		}
	}

	return &admin.DiskBackupSnapshotRestoreJob{}
}
