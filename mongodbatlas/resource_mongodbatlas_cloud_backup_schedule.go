package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/spf13/cast"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/mwielbut/pointy"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorSnapshotBackupScheduleUpdate  = "error updating a Cloud Backup Schedule: %s"
	errorSnapshotBackupScheduleRead    = "error getting a Cloud Backup Schedule for the cluster(%s): %s"
	errorSnapshotBackupScheduleSetting = "error setting `%s` for Cloud Backup Schedule(%s): %s"
)

// https://docs.atlas.mongodb.com/reference/api/cloud-backup/schedule/modify-one-schedule/
// same as resourceMongoDBAtlasCloudProviderSnapshotBackupPolicy
func resourceMongoDBAtlasCloudBackupSchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoDBAtlasCloudBackupScheduleCreate,
		Read:   resourceMongoDBAtlasCloudBackupScheduleRead,
		Update: resourceMongoDBAtlasCloudBackupScheduleUpdate, // To review
		Delete: resourceMongoDBAtlasCloudBackupScheduleDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMongoDBAtlasCloudBackupScheduleImportState,
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
			},
			"policies": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"policy_item": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"frequency_interval": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"frequency_type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"hourly", "daily", "weekly", "monthly"}, false),
									},
									"retention_unit": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"days", "weeks", "months"}, false),
									},
									"retention_value": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			// Optionals
			"reference_hour_of_day": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 0 || v > 23 {
						errs = append(errs, fmt.Errorf("%q value should be between 0 and 23, got: %d", key, v))
					}
					return
				},
			},
			"reference_minute_of_hour": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 0 || v > 59 {
						errs = append(errs, fmt.Errorf("%q value should be between 0 and 59, got: %d", key, v))
					}
					return
				},
			},
			"restore_window_days": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"update_snapshots": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			// Only computed
			"cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"next_snapshot": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceMongoDBAtlasCloudBackupScheduleCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)
	clusterName := d.Get("cluster_name").(string)

	performUpdate := false

	// Reading default configuration
	backupPolicy, _, err := conn.CloudProviderSnapshotBackupPolicies.Get(context.Background(), projectID, clusterName)

	if err != nil {
		return err
	}

	var config []matlas.Policy

	// if no policies getting default policies
	_, ok := d.GetOk("policies")

	if !ok {
		config = backupPolicy.Policies
	} else {
		config = expandPolicies(d)
	}

	// compare as set
	if needsUpdate(config, backupPolicy.Policies) {
		performUpdate = true
	}

	req, performUpdate := buildRequestCloudBackupSchedule(d, backupPolicy, config, performUpdate)

	if performUpdate {
		_, _, err := conn.CloudProviderSnapshotBackupPolicies.Update(context.Background(), projectID, clusterName, req)
		if err != nil {
			return fmt.Errorf(errorSnapshotBackupScheduleUpdate, err)
		}
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":   projectID,
		"cluster_name": clusterName,
	}))

	return resourceMongoDBAtlasCloudBackupScheduleRead(d, meta)
}

func resourceMongoDBAtlasCloudBackupScheduleRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*matlas.Client)

	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	backupPolicy, _, err := conn.CloudProviderSnapshotBackupPolicies.Get(context.Background(), projectID, clusterName)
	if err != nil {
		return fmt.Errorf(errorSnapshotBackupScheduleRead, clusterName, err)
	}

	if err := d.Set("cluster_id", backupPolicy.ClusterID); err != nil {
		return fmt.Errorf(errorSnapshotBackupScheduleSetting, "cluster_id", clusterName, err)
	}

	if err := d.Set("reference_hour_of_day", backupPolicy.ReferenceHourOfDay); err != nil {
		return fmt.Errorf(errorSnapshotBackupScheduleSetting, "reference_hour_of_day", clusterName, err)
	}

	if err := d.Set("reference_minute_of_hour", backupPolicy.ReferenceMinuteOfHour); err != nil {
		return fmt.Errorf(errorSnapshotBackupScheduleSetting, "reference_minute_of_hour", clusterName, err)
	}

	if err := d.Set("restore_window_days", backupPolicy.RestoreWindowDays); err != nil {
		return fmt.Errorf(errorSnapshotBackupScheduleSetting, "restore_window_days", clusterName, err)
	}

	if err := d.Set("update_snapshots", backupPolicy.UpdateSnapshots); err != nil {
		return fmt.Errorf(errorSnapshotBackupScheduleSetting, "update_snapshots", clusterName, err)
	}

	if err := d.Set("next_snapshot", backupPolicy.NextSnapshot); err != nil {
		return fmt.Errorf(errorSnapshotBackupScheduleSetting, "next_snapshot", clusterName, err)
	}

	if err := d.Set("policies", flattenPolicies(backupPolicy.Policies)); err != nil {
		return fmt.Errorf(errorSnapshotBackupScheduleSetting, "policies", clusterName, err)
	}

	return nil
}

func resourceMongoDBAtlasCloudBackupScheduleUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	if err := snapshotScheduleUpdate(d, conn, projectID, clusterName); err != nil {
		return err
	}

	if restoreWindowDays, ok := d.GetOk("restore_window_days"); ok {
		if cast.ToInt64(restoreWindowDays) <= 0 {
			return fmt.Errorf("`restore_window_days` cannot be <= 0")
		}
	}

	req := &matlas.CloudProviderSnapshotBackupPolicy{
		ReferenceHourOfDay:    pointy.Int64(cast.ToInt64(d.Get("reference_hour_of_day"))),
		ReferenceMinuteOfHour: pointy.Int64(cast.ToInt64(d.Get("reference_minute_of_hour"))),
		UpdateSnapshots:       pointy.Bool(cast.ToBool(d.Get("update_snapshots").(bool))),
	}

	if rwd, ok := d.GetOk("restore_window_days"); ok {
		req.RestoreWindowDays = pointy.Int64(cast.ToInt64(rwd))
	}

	if d.HasChange("policies") {
		req.Policies = expandPolicies(d)
	}

	_, _, err := conn.CloudProviderSnapshotBackupPolicies.Update(context.Background(), projectID, clusterName, req)
	if err != nil {
		return fmt.Errorf(errorSnapshotBackupScheduleUpdate, err)
	}

	return resourceMongoDBAtlasCloudBackupScheduleRead(d, meta)
}

func resourceMongoDBAtlasCloudBackupScheduleDelete(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	_, _, err := conn.CloudProviderSnapshotBackupPolicies.Delete(context.Background(), projectID, clusterName)
	if err != nil {
		return fmt.Errorf("error deleting MongoDB Cloud Backup Schedule (%s): %s", clusterName, err)
	}

	return nil
}

func needsUpdate(a, b []matlas.Policy) bool {
	if len(a) != len(b) {
		return true
	}

	// deeper diff sort everything XD
	sort.Slice(a, func(i, j int) bool {
		// sort each item
		return a[i].ID <= a[j].ID
	})

	sort.Slice(b, func(i, j int) bool {
		return b[i].ID <= b[j].ID
	})

	// sort subitems
	sortItems := func(array []matlas.Policy) {
		for _, item := range array {
			sort.Slice(item.PolicyItems, func(i, j int) bool {
				return item.PolicyItems[i].ID <= item.PolicyItems[j].ID
			})
		}
	}

	sortItems(a)
	sortItems(b)

	return reflect.DeepEqual(a, b)
}

func compareInt64(a, b *int64) bool {
	if b == nil {
		return true
	} else if *a != *b {
		return true
	}
	return false
}

func buildRequestCloudBackupSchedule(d *schema.ResourceData, backupPolicy *matlas.CloudProviderSnapshotBackupPolicy,
	policies []matlas.Policy, performUpdate bool) (*matlas.CloudProviderSnapshotBackupPolicy, bool) {
	// tenative request
	req := &matlas.CloudProviderSnapshotBackupPolicy{
		Policies: policies,
	}

	hourDay, ok := d.GetOk("reference_hour_of_day")
	if ok {
		value := pointy.Int64(cast.ToInt64(hourDay))
		if compareInt64(value, backupPolicy.ReferenceHourOfDay) {
			performUpdate = true
			req.ReferenceHourOfDay = value
		}
	}
	minHour, ok := d.GetOk("reference_minute_of_hour")
	if ok {
		value := pointy.Int64(cast.ToInt64(minHour))
		if compareInt64(value, backupPolicy.ReferenceMinuteOfHour) {
			performUpdate = true
			req.ReferenceMinuteOfHour = value
		}
	}
	winDays, ok := d.GetOk("restore_window_days")
	if ok {
		value := pointy.Int64(cast.ToInt64(winDays))
		if compareInt64(value, backupPolicy.RestoreWindowDays) {
			performUpdate = true
			req.RestoreWindowDays = value
		}
	}
	updateSnap, ok := d.GetOk("update_snapshots")
	if ok {
		value := pointy.Bool(updateSnap.(bool))
		if backupPolicy.UpdateSnapshots != nil {
			// just when true sending back
			if *value {
				performUpdate = true
				req.UpdateSnapshots = value
			}
		} else if *backupPolicy.UpdateSnapshots != *value {
			performUpdate = true
			req.UpdateSnapshots = value
		}
	}

	return req, performUpdate
}

func resourceMongoDBAtlasCloudBackupScheduleImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*matlas.Client)

	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 2 {
		return nil, errors.New("import format error: to import a Cloud Backup Schedule use the format {project_id}-{cluster_name}")
	}

	projectID := parts[0]
	clusterName := parts[1]

	_, _, err := conn.CloudProviderSnapshotBackupPolicies.Get(context.Background(), projectID, clusterName)
	if err != nil {
		return nil, fmt.Errorf(errorSnapshotBackupScheduleRead, clusterName, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		return nil, fmt.Errorf(errorSnapshotBackupScheduleSetting, "project_id", clusterName, err)
	}

	if err := d.Set("cluster_name", clusterName); err != nil {
		return nil, fmt.Errorf(errorSnapshotBackupScheduleSetting, "cluster_name", clusterName, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":   projectID,
		"cluster_name": clusterName,
	}))

	return []*schema.ResourceData{d}, nil
}
