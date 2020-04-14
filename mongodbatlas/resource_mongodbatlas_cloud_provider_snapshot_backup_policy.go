package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/spf13/cast"

	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

const (
	errorSnapshotBackupPolicyUpdate  = "error updating a Cloud Provider Snapshot Backup Policy: %s"
	errorSnapshotBackupPolicyRead    = "error getting a Cloud Provider Snapshot Backup Policy for the cluster(%s): %s"
	errorSnapshotBackupPolicySetting = "error setting `%s` for Cloud Provider Snapshot Backup Policy(%s): %s"
)

func resourceMongoDBAtlasCloudProviderSnapshotBackupPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoDBAtlasCloudProviderSnapshotBackupPolicyCreate,
		Update: resourceMongoDBAtlasCloudProviderSnapshotBackupPolicyUpdate,
		Read:   resourceMongoDBAtlasCloudProviderSnapshotBackupPolicyRead,
		Delete: resourceMongoDBAtlasCloudProviderSnapshotBackupPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMongoDBAtlasCloudProviderSnapshotBackupPolicyImportState,
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
			"cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
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
			"next_snapshot": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policies": {
				Type:     schema.TypeList,
				Required: true,
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
		},
	}
}

func resourceMongoDBAtlasCloudProviderSnapshotBackupPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)
	clusterName := d.Get("cluster_name").(string)

	// there is not an entry point to create a snapshot backup policy until it will use the update entry point
	if err := snapshotScheduleUpdate(d, conn, projectID, clusterName); err != nil {
		return err
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":   projectID,
		"cluster_name": clusterName,
	}))

	return resourceMongoDBAtlasCloudProviderSnapshotBackupPolicyRead(d, meta)
}

func resourceMongoDBAtlasCloudProviderSnapshotBackupPolicyRead(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)

	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	backupPolicy, _, err := conn.CloudProviderSnapshotBackupPolicies.Get(context.Background(), projectID, clusterName)
	if err != nil {
		return fmt.Errorf(errorSnapshotBackupPolicyRead, clusterName, err)
	}

	if err := d.Set("cluster_id", backupPolicy.ClusterID); err != nil {
		return fmt.Errorf(errorSnapshotBackupPolicySetting, "cluster_id", clusterName, err)
	}
	if err := d.Set("reference_hour_of_day", backupPolicy.ReferenceHourOfDay); err != nil {
		return fmt.Errorf(errorSnapshotBackupPolicySetting, "reference_hour_of_day", clusterName, err)
	}
	if err := d.Set("reference_minute_of_hour", backupPolicy.ReferenceMinuteOfHour); err != nil {
		return fmt.Errorf(errorSnapshotBackupPolicySetting, "reference_minute_of_hour", clusterName, err)
	}
	if err := d.Set("restore_window_days", backupPolicy.RestoreWindowDays); err != nil {
		return fmt.Errorf(errorSnapshotBackupPolicySetting, "restore_window_days", clusterName, err)
	}
	if err := d.Set("update_snapshots", backupPolicy.UpdateSnapshots); err != nil {
		return fmt.Errorf(errorSnapshotBackupPolicySetting, "update_snapshots", clusterName, err)
	}
	if err := d.Set("next_snapshot", backupPolicy.NextSnapshot); err != nil {
		return fmt.Errorf(errorSnapshotBackupPolicySetting, "next_snapshot", clusterName, err)
	}
	if err := d.Set("policies", flattenPolicies(backupPolicy.Policies)); err != nil {
		return fmt.Errorf(errorSnapshotBackupPolicySetting, "policies", clusterName, err)
	}

	return nil
}

func resourceMongoDBAtlasCloudProviderSnapshotBackupPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	if err := snapshotScheduleUpdate(d, conn, projectID, clusterName); err != nil {
		return err
	}

	return resourceMongoDBAtlasCloudProviderSnapshotBackupPolicyRead(d, meta)
}

func resourceMongoDBAtlasCloudProviderSnapshotBackupPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	// There is no resource to delete a backup policy, it can only be updated.
	return nil
}

func resourceMongoDBAtlasCloudProviderSnapshotBackupPolicyImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*matlas.Client)

	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 2 {
		return nil, errors.New("import format error: to import a Cloud Provider Snapshot Backup Policy use the format {project_id}-{cluster_name}")
	}

	projectID := parts[0]
	clusterName := parts[1]

	_, _, err := conn.CloudProviderSnapshotBackupPolicies.Get(context.Background(), projectID, clusterName)
	if err != nil {
		return nil, fmt.Errorf(errorSnapshotBackupPolicyRead, clusterName, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		return nil, fmt.Errorf(errorSnapshotBackupPolicySetting, "project_id", clusterName, err)
	}
	if err := d.Set("cluster_name", clusterName); err != nil {
		return nil, fmt.Errorf(errorSnapshotBackupPolicySetting, "cluster_name", clusterName, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":   projectID,
		"cluster_name": clusterName,
	}))

	return []*schema.ResourceData{d}, nil
}

func snapshotScheduleUpdate(d *schema.ResourceData, conn *matlas.Client, projectID, clusterName string) error {
	req := &matlas.CloudProviderSnapshotBackupPolicy{
		ReferenceHourOfDay:    d.Get("reference_hour_of_day").(int),
		ReferenceMinuteOfHour: d.Get("reference_minute_of_hour").(int),
		RestoreWindowDays:     d.Get("restore_window_days").(int),
		UpdateSnapshots:       cast.ToBool(d.Get("update_snapshots").(bool)),
		Policies:              expandPolicies(d),
	}

	_, _, err := conn.CloudProviderSnapshotBackupPolicies.Update(context.Background(), projectID, clusterName, req)
	if err != nil {
		return fmt.Errorf(errorSnapshotBackupPolicyUpdate, err)
	}

	return nil
}

func expandPolicies(d *schema.ResourceData) []matlas.Policy {
	policies := make([]matlas.Policy, len(d.Get("policies").([]interface{})))

	for k, v := range d.Get("policies").([]interface{}) {
		policy := v.(map[string]interface{})
		policies[k] = matlas.Policy{
			ID:          policy["id"].(string),
			PolicyItems: expandPolicyItems(policy["policy_item"].([]interface{})),
		}
	}
	return policies
}

func flattenPolicies(policies []matlas.Policy) []map[string]interface{} {
	actionList := make([]map[string]interface{}, 0)
	for _, v := range policies {
		actionList = append(actionList, map[string]interface{}{
			"id":          v.ID,
			"policy_item": flattenPolicyItems(v.PolicyItems),
		})
	}
	return actionList
}

func expandPolicyItems(p []interface{}) []matlas.PolicyItem {
	policyItems := make([]matlas.PolicyItem, len(p))

	for k, v := range p {
		item := v.(map[string]interface{})
		policyItems[k] = matlas.PolicyItem{
			ID:                item["id"].(string),
			FrequencyInterval: item["frequency_interval"].(int),
			FrequencyType:     item["frequency_type"].(string),
			RetentionUnit:     item["retention_unit"].(string),
			RetentionValue:    item["retention_value"].(int),
		}
	}
	return policyItems
}

func flattenPolicyItems(items []matlas.PolicyItem) []map[string]interface{} {
	policyItems := make([]map[string]interface{}, 0)
	for _, v := range items {
		policyItems = append(policyItems, map[string]interface{}{
			"id":                 v.ID,
			"frequency_interval": v.FrequencyInterval,
			"frequency_type":     v.FrequencyType,
			"retention_unit":     v.RetentionUnit,
			"retention_value":    v.RetentionValue,
		})
	}
	return policyItems
}

func flattenCloudProviderSnapshotBackupPolicy(d *schema.ResourceData, conn *matlas.Client, projectID, clusterName string) ([]map[string]interface{}, error) {
	backupPolicy, _, err := conn.CloudProviderSnapshotBackupPolicies.Get(context.Background(), projectID, clusterName)
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "BACKUP_CONFIG_NOT_FOUND") || strings.Contains(fmt.Sprint(err), "400") {
			return []map[string]interface{}{}, nil
		}
		return []map[string]interface{}{}, fmt.Errorf(errorSnapshotBackupPolicyRead, clusterName, err)
	}

	return []map[string]interface{}{
		{
			"cluster_id":               backupPolicy.ClusterID,
			"cluster_name":             backupPolicy.ClusterName,
			"next_snapshot":            backupPolicy.NextSnapshot,
			"reference_hour_of_day":    backupPolicy.ReferenceHourOfDay,
			"reference_minute_of_hour": backupPolicy.ReferenceMinuteOfHour,
			"restore_window_days":      backupPolicy.RestoreWindowDays,
			"update_snapshots":         cast.ToBool(backupPolicy.UpdateSnapshots),
			"policies":                 flattenPolicies(backupPolicy.Policies),
		},
	}, nil
}

func computedCloudProviderSnapshotBackupPolicySchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"cluster_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"cluster_name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"next_snapshot": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"reference_hour_of_day": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"reference_minute_of_hour": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"restore_window_days": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"update_snapshots": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"policies": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"id": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"policy_item": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"id": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"frequency_interval": {
											Type:     schema.TypeInt,
											Computed: true,
										},
										"frequency_type": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"retention_unit": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"retention_value": {
											Type:     schema.TypeInt,
											Computed: true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
