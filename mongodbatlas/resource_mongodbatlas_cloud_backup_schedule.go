package mongodbatlas

import (
	"context"
	"fmt"
	"reflect"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/mwielbut/pointy"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

// https://docs.atlas.mongodb.com/reference/api/cloud-backup/schedule/modify-one-schedule/
// sasme as resourceMongoDBAtlasCloudProviderSnapshotBackupPolicy
func resourceMongoDBAtlasCloudBackupSchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoDBAtlasCloudBackupScheduleCreate,
		Read:   resourceMongoDBAtlasCloudProviderSnapshotBackupPolicyRead,
		Update: resourceMongoDBAtlasCloudProviderSnapshotBackupPolicyUpdate, // To review
		Delete: resourceMongoDBAtlasCloudProviderSnapshotBackupPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMongoDBAtlasCloudProviderSnapshotBackupPolicyImportState,
		},
		// delete is pending to check
		// Delete: resourceMongoDBAtlasCloudProviderSnapshotBackupPolicyDelete,

		Schema: map[string]*schema.Schema{
			// Required
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
				// remember we can import this , and of course it can
				// return computed items as well
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
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Required: true,
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
		// if there was no policies in config set what we get from response
		config = backupPolicy.Policies
	} else {
		config = expandPolicies(d)
	}

	// compare as set
	if needsUpdate(config, backupPolicy.Policies) {
		performUpdate = true
	}

	// tenative request
	req := &matlas.CloudProviderSnapshotBackupPolicy{
		Policies: config,
	}

	// Refactor all of this into a function this is ugly
	hourDay, ok := d.GetOk("reference_hour_of_day")

	if ok {
		value := pointy.Int64(hourDay.(int64))
		if compareInt64(value, backupPolicy.ReferenceHourOfDay) {
			performUpdate = true
			req.ReferenceHourOfDay = value
		}
	}

	minHour, ok := d.GetOk("reference_minute_of_hour")

	if ok {
		value := pointy.Int64(minHour.(int64))
		if compareInt64(value, backupPolicy.ReferenceMinuteOfHour) {
			performUpdate = true
			req.ReferenceMinuteOfHour = value
		}
	}

	winDays, ok := d.GetOk("restore_window_days")

	if ok {
		value := pointy.Int64(winDays.(int64))
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

	if performUpdate {
		_, _, err := conn.CloudProviderSnapshotBackupPolicies.Update(context.Background(), projectID, clusterName, req)
		if err != nil {
			return fmt.Errorf(errorSnapshotBackupPolicyUpdate, err)
		}
	}

	// otherwise set read config
	return resourceMongoDBAtlasCloudProviderSnapshotBackupPolicyRead(d, meta)
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
