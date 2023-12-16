package cluster

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/mwielbut/pointy"
	"github.com/spf13/cast"
)

const (
	errorClusterCreate            = "error creating MongoDB Cluster: %s"
	errorClusterRead              = "error reading MongoDB Cluster (%s): %s"
	errorClusterDelete            = "error deleting MongoDB Cluster (%s): %s"
	errorClusterUpdate            = "error updating MongoDB Cluster (%s): %s"
	errorAdvancedConfUpdate       = "error updating Advanced Configuration Option form MongoDB Cluster (%s): %s"
	ErrorSnapshotBackupPolicyRead = "error getting a Cloud Provider Snapshot Backup Policy for the cluster(%s): %s"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceMongoDBAtlasClusterCreate,
		ReadWithoutTimeout:   resourceMongoDBAtlasClusterRead,
		UpdateWithoutTimeout: resourceMongoDBAtlasClusterUpdate,
		DeleteWithoutTimeout: resourceMongoDBAtlasClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasClusterImportState,
		},
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    ResourceClusterResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceMongoDBAtlasClusterStateUpgradeV0,
				Version: 0,
			},
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_scaling_disk_gb_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"auto_scaling_compute_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"auto_scaling_compute_scale_down_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"backup_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Clusters running MongoDB FCV 4.2 or later and any new Atlas clusters of any type do not support this parameter",
			},
			"retain_backups_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Flag that indicates whether to retain backup snapshots for the deleted dedicated cluster",
			},
			"bi_connector_config": {
				Type:       schema.TypeList,
				Optional:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Computed:   true,
				MaxItems:   1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"read_preference": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"cluster_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"connection_strings": advancedcluster.ClusterConnectionStringsSchema(),
			"disk_size_gb": {
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
			},
			"encryption_at_rest_provider": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"mongo_db_major_version": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				StateFunc: advancedcluster.FormatMongoDBMajorVersion,
			},
			"num_shards": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"cloud_backup": {
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"backup_enabled"},
			},
			"provider_instance_size_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"provider_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pit_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"backing_provider_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"provider_disk_iops": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"provider_disk_type_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"provider_encrypt_ebs_volume": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "All EBS volumes are encrypted by default, the option to disable encryption has been removed",
				Computed:   true,
			},
			"provider_encrypt_ebs_volume_flag": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"provider_region_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"provider_volume_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"provider_auto_scaling_compute_min_instance_size": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: isEqualProviderAutoScalingMinInstanceSize,
			},
			"provider_auto_scaling_compute_max_instance_size": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: isEqualProviderAutoScalingMaxInstanceSize,
			},
			"replication_factor": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"replication_specs": {
				Type:       schema.TypeSet,
				Optional:   true,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"num_shards": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"regions_config": {
							Type:       schema.TypeSet,
							Optional:   true,
							Computed:   true,
							ConfigMode: schema.SchemaConfigModeAttr,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"electable_nodes": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"priority": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"read_only_nodes": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  0,
									},
									"analytics_nodes": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  0,
									},
								},
							},
						},
						"zone_name": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "ZoneName managed by Terraform",
						},
					},
				},
				Set: func(v any) int {
					var buf bytes.Buffer
					m := v.(map[string]any)
					buf.WriteString(fmt.Sprintf("%d", m["num_shards"].(int)))
					buf.WriteString(m["zone_name"].(string))
					buf.WriteString(fmt.Sprintf("%+v", m["regions_config"].(*schema.Set)))
					return advancedcluster.HashCodeString(buf.String())
				},
			},
			"mongo_db_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mongo_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mongo_uri_updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mongo_uri_with_options": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"paused": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"srv_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"advanced_configuration": advancedcluster.ClusterAdvancedConfigurationSchema(),
			"labels": {
				Type:       schema.TypeSet,
				Optional:   true,
				Set:        advancedcluster.HashFunctionForKeyValuePair,
				Deprecated: fmt.Sprintf(constant.DeprecationParamByDateWithReplacement, "September 2024", "tags"),
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"tags":                   &advancedcluster.RSTagsSchema,
			"snapshot_backup_policy": computedCloudProviderSnapshotBackupPolicySchema(),
			"termination_protection_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"container_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version_release_system": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"LTS", "CONTINUOUS"}, false),
			},
			"accept_data_risks_and_force_replica_set_reconfig": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Submit this field alongside your topology reconfiguration to request a new regional outage resistant topology",
			},
		},
		CustomizeDiff: resourceClusterCustomizeDiff,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Hour),
			Update: schema.DefaultTimeout(3 * time.Hour),
			Delete: schema.DefaultTimeout(3 * time.Hour),
		},
	}
}

func resourceMongoDBAtlasClusterCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	if v, ok := d.GetOk("accept_data_risks_and_force_replica_set_reconfig"); ok {
		if v.(string) != "" {
			return diag.FromErr(fmt.Errorf("accept_data_risks_and_force_replica_set_reconfig can not be set in creation, only in update"))
		}
	}

	conn := meta.(*config.MongoDBClient).Atlas

	projectID := d.Get("project_id").(string)
	providerName := d.Get("provider_name").(string)

	computeEnabled := d.Get("auto_scaling_compute_enabled").(bool)
	scaleDownEnabled := d.Get("auto_scaling_compute_scale_down_enabled").(bool)
	minInstanceSize := d.Get("provider_auto_scaling_compute_min_instance_size").(string)
	maxInstanceSize := d.Get("provider_auto_scaling_compute_max_instance_size").(string)

	if scaleDownEnabled && !computeEnabled {
		return diag.FromErr(fmt.Errorf("`auto_scaling_compute_scale_down_enabled` must be set when `auto_scaling_compute_enabled` is set"))
	}

	if computeEnabled && maxInstanceSize == "" {
		return diag.FromErr(fmt.Errorf("`provider_auto_scaling_compute_max_instance_size` must be set when `auto_scaling_compute_enabled` is set"))
	}

	if scaleDownEnabled && minInstanceSize == "" {
		return diag.FromErr(fmt.Errorf("`provider_auto_scaling_compute_min_instance_size` must be set when `auto_scaling_compute_scale_down_enabled` is set"))
	}

	autoScaling := &matlas.AutoScaling{
		DiskGBEnabled: pointy.Bool(d.Get("auto_scaling_disk_gb_enabled").(bool)),
		Compute: &matlas.Compute{
			Enabled:          &computeEnabled,
			ScaleDownEnabled: &scaleDownEnabled,
		},
	}

	// validate cluster_type conditional
	if _, ok := d.GetOk("replication_specs"); ok {
		if _, ok1 := d.GetOk("cluster_type"); !ok1 {
			return diag.FromErr(fmt.Errorf("`cluster_type` should be set when `replication_specs` is set"))
		}
	}

	if providerName != "AWS" {
		if _, ok := d.GetOk("provider_disk_iops"); ok {
			return diag.Errorf("`provider_disk_iops` shouldn't be set when provider name is `GCP` or `AZURE`")
		}

		if _, ok := d.GetOk("provider_volume_type"); ok {
			return diag.FromErr(fmt.Errorf("`provider_volume_type` shouldn't be set when provider name is `GCP` or `AZURE`"))
		}
	}

	if providerName != "AZURE" {
		if _, ok := d.GetOk("provider_disk_type_name"); ok {
			return diag.FromErr(fmt.Errorf("`provider_disk_type_name` shouldn't be set when provider name is `GCP` or `AWS`"))
		}
	}

	if providerName == "AZURE" {
		if _, ok := d.GetOk("disk_size_gb"); ok {
			return diag.FromErr(fmt.Errorf("`disk_size_gb` cannot be used with Azure clusters"))
		}
	}

	tenantDisksize := pointy.Float64(0)
	if providerName == "TENANT" {
		autoScaling = nil

		if instanceSizeName, ok := d.GetOk("provider_instance_size_name"); ok {
			if instanceSizeName == "M2" {
				if diskSizeGB, ok := d.GetOk("disk_size_gb"); ok {
					if cast.ToFloat64(diskSizeGB) != 2 {
						return diag.FromErr(fmt.Errorf("`disk_size_gb` must be 2 for M2 shared tier"))
					}
				}
			}
			if instanceSizeName == "M5" {
				if diskSizeGB, ok := d.GetOk("disk_size_gb"); ok {
					if cast.ToFloat64(diskSizeGB) != 5 {
						return diag.FromErr(fmt.Errorf("`disk_size_gb` must be 5 for M5 shared tier"))
					}
				}
			}
		}
	}

	// We need to validate the oplog_size_mb attr of the advanced configuration option to show the error
	// before that the cluster is created
	if oplogSizeMB, ok := d.GetOkExists("advanced_configuration.0.oplog_size_mb"); ok {
		if cast.ToInt64(oplogSizeMB) <= 0 {
			return diag.FromErr(fmt.Errorf("`advanced_configuration.oplog_size_mb` cannot be <= 0"))
		}
	}

	providerSettings, err := expandProviderSetting(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterCreate, err))
	}

	replicationSpecs, err := expandReplicationSpecs(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterCreate, err))
	}

	clusterRequest := &matlas.Cluster{
		Name:                     d.Get("name").(string),
		EncryptionAtRestProvider: d.Get("encryption_at_rest_provider").(string),
		ClusterType:              cast.ToString(d.Get("cluster_type")),
		BackupEnabled:            pointy.Bool(d.Get("backup_enabled").(bool)),
		PitEnabled:               pointy.Bool(d.Get("pit_enabled").(bool)),
		AutoScaling:              autoScaling,
		ProviderSettings:         providerSettings,
		ReplicationSpecs:         replicationSpecs,
	}
	if v, ok := d.GetOk("cloud_backup"); ok {
		clusterRequest.ProviderBackupEnabled = pointy.Bool(v.(bool))
	}

	if _, ok := d.GetOk("bi_connector_config"); ok {
		biConnector, err := advancedcluster.ExpandBiConnectorConfig(d)
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorClusterCreate, err))
		}
		clusterRequest.BiConnector = biConnector
	}

	if advancedcluster.ContainsLabelOrKey(advancedcluster.ExpandLabelSliceFromSetSchema(d), advancedcluster.DefaultLabel) {
		return diag.FromErr(fmt.Errorf("you should not set `Infrastructure Tool` label, it is used for internal purposes"))
	}

	clusterRequest.Labels = append(advancedcluster.ExpandLabelSliceFromSetSchema(d), advancedcluster.DefaultLabel)

	if _, ok := d.GetOk("tags"); ok {
		tagsSlice := advancedcluster.ExpandTagSliceFromSetSchema(d)
		clusterRequest.Tags = &tagsSlice
	}

	if v, ok := d.GetOk("disk_size_gb"); ok {
		clusterRequest.DiskSizeGB = pointy.Float64(v.(float64))
	}
	if cast.ToFloat64(tenantDisksize) != 0 {
		clusterRequest.DiskSizeGB = tenantDisksize
	}
	if v, ok := d.GetOk("mongo_db_major_version"); ok {
		clusterRequest.MongoDBMajorVersion = advancedcluster.FormatMongoDBMajorVersion(v.(string))
	}

	if r, ok := d.GetOk("replication_factor"); ok {
		clusterRequest.ReplicationFactor = pointy.Int64(cast.ToInt64(r))
	}

	if n, ok := d.GetOk("num_shards"); ok {
		clusterRequest.NumShards = pointy.Int64(cast.ToInt64(n))
	}

	if v, ok := d.GetOk("termination_protection_enabled"); ok {
		clusterRequest.TerminationProtectionEnabled = pointy.Bool(v.(bool))
	}

	if v, ok := d.GetOk("version_release_system"); ok {
		clusterRequest.VersionReleaseSystem = v.(string)
	}

	cluster, _, err := conn.Clusters.Create(ctx, projectID, clusterRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterCreate, err))
	}

	timeout := d.Timeout(schema.TimeoutCreate)
	stateConf := &retry.StateChangeConf{
		Pending:    []string{"CREATING", "UPDATING", "REPAIRING", "REPEATING", "PENDING"},
		Target:     []string{"IDLE"},
		Refresh:    advancedcluster.ResourceClusterRefreshFunc(ctx, d.Get("name").(string), projectID, conn),
		Timeout:    timeout,
		MinTimeout: 1 * time.Minute,
		Delay:      3 * time.Minute,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterCreate, err))
	}

	/*
		So far, the cluster has created correctly, so we need to set up
		the advanced configuration option to attach it
	*/
	ac, ok := d.GetOk("advanced_configuration")
	if aclist, ok1 := ac.([]any); ok1 && len(aclist) > 0 {
		advancedConfReq := advancedcluster.ExpandProcessArgs(d, aclist[0].(map[string]any))

		if ok {
			_, _, err := conn.Clusters.UpdateProcessArgs(ctx, projectID, cluster.Name, advancedConfReq)
			if err != nil {
				return diag.FromErr(fmt.Errorf(errorAdvancedConfUpdate, cluster.Name, err))
			}
		}
	}

	// To pause a cluster
	if v := d.Get("paused").(bool); v {
		clusterRequest = &matlas.Cluster{
			Paused: pointy.Bool(v),
		}

		_, _, err = updateCluster(ctx, conn, clusterRequest, projectID, d.Get("name").(string), timeout)
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorClusterUpdate, d.Get("name").(string), err))
		}
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"cluster_id":    cluster.ID,
		"project_id":    projectID,
		"cluster_name":  cluster.Name,
		"provider_name": providerName,
	}))

	return resourceMongoDBAtlasClusterRead(ctx, d, meta)
}

func resourceMongoDBAtlasClusterRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*config.MongoDBClient).Atlas
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]
	providerName := ids["provider_name"]

	cluster, resp, err := conn.Clusters.Get(ctx, projectID, clusterName)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorClusterRead, clusterName, err))
	}

	log.Printf("[DEBUG] GET Cluster %+v", cluster)

	if err := d.Set("cluster_id", cluster.ID); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "cluster_id", clusterName, err))
	}

	if err := d.Set("auto_scaling_compute_enabled", cluster.AutoScaling.Compute.Enabled); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "auto_scaling_compute_enabled", clusterName, err))
	}

	if err := d.Set("auto_scaling_compute_scale_down_enabled", cluster.AutoScaling.Compute.ScaleDownEnabled); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "auto_scaling_compute_scale_down_enabled", clusterName, err))
	}

	if err := d.Set("provider_auto_scaling_compute_min_instance_size", cluster.ProviderSettings.AutoScaling.Compute.MinInstanceSize); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "provider_auto_scaling_compute_min_instance_size", clusterName, err))
	}

	if err := d.Set("provider_auto_scaling_compute_max_instance_size", cluster.ProviderSettings.AutoScaling.Compute.MaxInstanceSize); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "provider_auto_scaling_compute_max_instance_size", clusterName, err))
	}

	if err := d.Set("backup_enabled", cluster.BackupEnabled); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "backup_enabled", clusterName, err))
	}

	if _, ok := d.GetOk("cloud_backup"); ok {
		if err := d.Set("cloud_backup", cluster.ProviderBackupEnabled); err != nil {
			return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "cloud_backup", clusterName, err))
		}
	}

	if err := d.Set("cluster_type", cluster.ClusterType); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "cluster_type", clusterName, err))
	}

	if err := d.Set("connection_strings", advancedcluster.FlattenConnectionStrings(cluster.ConnectionStrings)); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "connection_strings", clusterName, err))
	}

	if err := d.Set("disk_size_gb", cluster.DiskSizeGB); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "disk_size_gb", clusterName, err))
	}

	if err := d.Set("encryption_at_rest_provider", cluster.EncryptionAtRestProvider); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "encryption_at_rest_provider", clusterName, err))
	}

	if err := d.Set("mongo_db_major_version", cluster.MongoDBMajorVersion); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "mongo_db_major_version", clusterName, err))
	}

	// Avoid Global Cluster issues. (NumShards is not present in Global Clusters)
	if cluster.NumShards != nil {
		if err := d.Set("num_shards", cluster.NumShards); err != nil {
			return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "num_shards", clusterName, err))
		}
	}

	if err := d.Set("mongo_db_version", cluster.MongoDBVersion); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "mongo_db_version", clusterName, err))
	}

	if err := d.Set("mongo_uri", cluster.MongoURI); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "mongo_uri", clusterName, err))
	}

	if err := d.Set("mongo_uri_updated", cluster.MongoURIUpdated); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "mongo_uri_updated", clusterName, err))
	}

	if err := d.Set("mongo_uri_with_options", cluster.MongoURIWithOptions); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "mongo_uri_with_options", clusterName, err))
	}

	if err := d.Set("pit_enabled", cluster.PitEnabled); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "pit_enabled", clusterName, err))
	}

	if err := d.Set("paused", cluster.Paused); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "paused", clusterName, err))
	}

	if err := d.Set("srv_address", cluster.SrvAddress); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "srv_address", clusterName, err))
	}

	if err := d.Set("state_name", cluster.StateName); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "state_name", clusterName, err))
	}

	if err := d.Set("termination_protection_enabled", cluster.TerminationProtectionEnabled); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "termination_protection_enabled", clusterName, err))
	}

	if err := d.Set("bi_connector_config", advancedcluster.FlattenBiConnectorConfig(cluster.BiConnector)); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "bi_connector_config", clusterName, err))
	}

	if cluster.ProviderSettings != nil {
		flattenProviderSettings(d, cluster.ProviderSettings, clusterName)
	}

	if err := d.Set("replication_specs", flattenReplicationSpecs(cluster.ReplicationSpecs)); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "replication_specs", clusterName, err))
	}

	if err := d.Set("replication_factor", cluster.ReplicationFactor); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "replication_factor", clusterName, err))
	}

	if err := d.Set("labels", advancedcluster.FlattenLabels(advancedcluster.RemoveLabel(cluster.Labels, advancedcluster.DefaultLabel))); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "labels", clusterName, err))
	}

	if err := d.Set("tags", advancedcluster.FlattenTags(cluster.Tags)); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "tags", clusterName, err))
	}

	if err := d.Set("version_release_system", cluster.VersionReleaseSystem); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "version_release_system", clusterName, err))
	}

	if err := d.Set("accept_data_risks_and_force_replica_set_reconfig", cluster.AcceptDataRisksAndForceReplicaSetReconfig); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "accept_data_risks_and_force_replica_set_reconfig", clusterName, err))
	}

	if providerName != "TENANT" {
		containers, _, err := conn.Containers.List(ctx, projectID,
			&matlas.ContainersListOptions{ProviderName: providerName})
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorClusterRead, clusterName, err))
		}

		if err := d.Set("container_id", getContainerID(containers, cluster)); err != nil {
			return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "container_id", clusterName, err))
		}

		if err := d.Set("auto_scaling_disk_gb_enabled", cluster.AutoScaling.DiskGBEnabled); err != nil {
			return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "auto_scaling_disk_gb_enabled", clusterName, err))
		}
	}

	/*
		Get the advaced configuration options and set up to the terraform state
	*/
	processArgs, _, err := conn.Clusters.GetProcessArgs(ctx, projectID, clusterName)
	if err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorAdvancedConfRead, clusterName, err))
	}

	if err := d.Set("advanced_configuration", advancedcluster.FlattenProcessArgs(processArgs)); err != nil {
		return diag.FromErr(fmt.Errorf(advancedcluster.ErrorClusterSetting, "advanced_configuration", clusterName, err))
	}

	// Get the snapshot policy and set the data
	snapshotBackupPolicy, err := flattenCloudProviderSnapshotBackupPolicy(ctx, d, conn, projectID, clusterName)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("snapshot_backup_policy", snapshotBackupPolicy); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceMongoDBAtlasClusterUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*config.MongoDBClient).Atlas
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	cluster := new(matlas.Cluster)
	clusterChangeDetect := new(matlas.Cluster)
	clusterChangeDetect.AutoScaling = &matlas.AutoScaling{Compute: &matlas.Compute{}}

	if d.HasChange("name") {
		cluster.Name, _ = d.Get("name").(string)
	}

	if d.HasChange("bi_connector_config") {
		cluster.BiConnector, _ = advancedcluster.ExpandBiConnectorConfig(d)
	}

	// If at least one of the provider settings argument has changed, expand all provider settings
	if d.HasChange("provider_disk_iops") ||
		d.HasChange("backing_provider_name") ||
		d.HasChange("provider_disk_type_name") ||
		d.HasChange("provider_instance_size_name") ||
		d.HasChange("provider_name") ||
		d.HasChange("provider_region_name") ||
		d.HasChange("provider_volume_type") ||
		d.HasChange("provider_auto_scaling_compute_min_instance_size") ||
		d.HasChange("provider_auto_scaling_compute_max_instance_size") {
		var err error
		cluster.ProviderSettings, err = expandProviderSetting(d)
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorClusterUpdate, clusterName, err))
		}
	}

	if d.HasChange("replication_specs") {
		replicationSpecs, err := expandReplicationSpecs(d)
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorClusterUpdate, clusterName, err))
		}

		cluster.ReplicationSpecs = replicationSpecs
	}

	cluster.AutoScaling = &matlas.AutoScaling{Compute: &matlas.Compute{}}

	if d.HasChange("auto_scaling_disk_gb_enabled") {
		cluster.AutoScaling.DiskGBEnabled = pointy.Bool(d.Get("auto_scaling_disk_gb_enabled").(bool))
	}

	if d.HasChange("auto_scaling_compute_enabled") {
		cluster.AutoScaling.Compute.Enabled = pointy.Bool(d.Get("auto_scaling_compute_enabled").(bool))
	}

	if d.HasChange("auto_scaling_compute_scale_down_enabled") {
		cluster.AutoScaling.Compute.ScaleDownEnabled = pointy.Bool(d.Get("auto_scaling_compute_scale_down_enabled").(bool))
	}

	if d.HasChange("encryption_at_rest_provider") {
		cluster.EncryptionAtRestProvider = d.Get("encryption_at_rest_provider").(string)
	}

	if d.HasChange("mongo_db_major_version") {
		cluster.MongoDBMajorVersion = advancedcluster.FormatMongoDBMajorVersion(d.Get("mongo_db_major_version"))
	}

	if d.HasChange("cluster_type") {
		cluster.ClusterType = d.Get("cluster_type").(string)
	}

	if d.HasChange("backup_enabled") {
		cluster.BackupEnabled = pointy.Bool(d.Get("backup_enabled").(bool))
	}

	if d.HasChange("disk_size_gb") {
		cluster.DiskSizeGB = pointy.Float64(d.Get("disk_size_gb").(float64))
	}

	if d.HasChange("cloud_backup") {
		cluster.ProviderBackupEnabled = pointy.Bool(d.Get("cloud_backup").(bool))
	}

	if d.HasChange("pit_enabled") {
		cluster.PitEnabled = pointy.Bool(d.Get("pit_enabled").(bool))
	}

	if d.HasChange("replication_factor") {
		cluster.ReplicationFactor = pointy.Int64(cast.ToInt64(d.Get("replication_factor")))
	}

	if d.HasChange("num_shards") {
		cluster.NumShards = pointy.Int64(cast.ToInt64(d.Get("num_shards")))
	}

	if d.HasChange("version_release_system") {
		cluster.VersionReleaseSystem = d.Get("version_release_system").(string)
	}

	if d.HasChange("accept_data_risks_and_force_replica_set_reconfig") {
		cluster.AcceptDataRisksAndForceReplicaSetReconfig = d.Get("accept_data_risks_and_force_replica_set_reconfig").(string)
	}

	if d.HasChange("termination_protection_enabled") {
		cluster.TerminationProtectionEnabled = pointy.Bool(d.Get("termination_protection_enabled").(bool))
	}

	if d.HasChange("labels") {
		if advancedcluster.ContainsLabelOrKey(advancedcluster.ExpandLabelSliceFromSetSchema(d), advancedcluster.DefaultLabel) {
			return diag.FromErr(fmt.Errorf("you should not set `Infrastructure Tool` label, it is used for internal purposes"))
		}

		cluster.Labels = append(advancedcluster.ExpandLabelSliceFromSetSchema(d), advancedcluster.DefaultLabel)
	}

	if d.HasChange("tags") {
		tagsSlice := advancedcluster.ExpandTagSliceFromSetSchema(d)
		cluster.Tags = &tagsSlice
	}

	// when Provider instance type changes this argument must be passed explicitly in patch request
	if d.HasChange("provider_instance_size_name") {
		if _, ok := d.GetOk("cloud_backup"); ok {
			cluster.ProviderBackupEnabled = pointy.Bool(d.Get("cloud_backup").(bool))
		}
	}

	if d.HasChange("paused") && !d.Get("paused").(bool) {
		cluster.Paused = pointy.Bool(d.Get("paused").(bool))
	}

	timeout := d.Timeout(schema.TimeoutUpdate)

	/*
		Check if advaced configuration option has a changes to update it
	*/
	if d.HasChange("advanced_configuration") {
		ac := d.Get("advanced_configuration")
		if aclist, ok1 := ac.([]any); ok1 && len(aclist) > 0 {
			advancedConfReq := advancedcluster.ExpandProcessArgs(d, aclist[0].(map[string]any))
			if !reflect.DeepEqual(advancedConfReq, matlas.ProcessArgs{}) {
				argResp, _, err := conn.Clusters.UpdateProcessArgs(ctx, projectID, clusterName, advancedConfReq)
				if err != nil {
					return diag.FromErr(fmt.Errorf(errorAdvancedConfUpdate, clusterName+argResp.DefaultReadConcern, err))
				}
			}
		}
	}

	if isUpgradeRequired(d) {
		updatedCluster, _, err := advancedcluster.UpgradeCluster(ctx, conn, cluster, projectID, clusterName, timeout)

		if err != nil {
			return diag.FromErr(fmt.Errorf(errorClusterUpdate, clusterName, err))
		}

		d.SetId(conversion.EncodeStateID(map[string]string{
			"cluster_id":    updatedCluster.ID,
			"project_id":    projectID,
			"cluster_name":  updatedCluster.Name,
			"provider_name": updatedCluster.ProviderSettings.ProviderName,
		}))
	} else if !reflect.DeepEqual(cluster, clusterChangeDetect) {
		err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
			_, _, err := updateCluster(ctx, conn, cluster, projectID, clusterName, timeout)

			if didErrOnPausedCluster(err) {
				clusterRequest := &matlas.Cluster{
					Paused: pointy.Bool(false),
				}

				_, _, err = updateCluster(ctx, conn, clusterRequest, projectID, clusterName, timeout)
			}

			if err != nil {
				return retry.NonRetryableError(fmt.Errorf(errorClusterUpdate, clusterName, err))
			}

			return nil
		})

		if err != nil {
			return diag.FromErr(fmt.Errorf(errorClusterUpdate, clusterName, err))
		}
	}

	if d.Get("paused").(bool) && !advancedcluster.IsSharedTier(d.Get("provider_instance_size_name").(string)) {
		clusterRequest := &matlas.Cluster{
			Paused: pointy.Bool(true),
		}

		_, _, err := updateCluster(ctx, conn, clusterRequest, projectID, clusterName, timeout)
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorClusterUpdate, clusterName, err))
		}
	}

	return resourceMongoDBAtlasClusterRead(ctx, d, meta)
}

func didErrOnPausedCluster(err error) bool {
	if err == nil {
		return false
	}

	var target *matlas.ErrorResponse

	return errors.As(err, &target) && target.ErrorCode == "CANNOT_UPDATE_PAUSED_CLUSTER"
}

func resourceMongoDBAtlasClusterDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*config.MongoDBClient).Atlas
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	var options *matlas.DeleteAdvanceClusterOptions
	if v, ok := d.GetOkExists("retain_backups_enabled"); ok {
		options = &matlas.DeleteAdvanceClusterOptions{
			RetainBackups: pointy.Bool(v.(bool)),
		}
	}

	_, err := conn.Clusters.Delete(ctx, projectID, clusterName, options)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterDelete, clusterName, err))
	}

	log.Println("[INFO] Waiting for MongoDB Cluster to be destroyed")

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"IDLE", "CREATING", "UPDATING", "REPAIRING", "DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    advancedcluster.ResourceClusterRefreshFunc(ctx, clusterName, projectID, conn),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		MinTimeout: 30 * time.Second,
		Delay:      1 * time.Minute, // Wait 30 secs before starting
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterDelete, clusterName, err))
	}

	return nil
}

func resourceMongoDBAtlasClusterImportState(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	conn := meta.(*config.MongoDBClient).Atlas

	projectID, name, err := splitSClusterImportID(d.Id())
	if err != nil {
		return nil, err
	}

	u, _, err := conn.Clusters.Get(ctx, *projectID, *name)
	if err != nil {
		return nil, fmt.Errorf("couldn't import cluster %s in project %s, error: %s", *name, *projectID, err)
	}

	if err := d.Set("project_id", u.GroupID); err != nil {
		log.Printf(advancedcluster.ErrorClusterSetting, "project_id", u.ID, err)
	}

	if err := d.Set("name", u.Name); err != nil {
		log.Printf(advancedcluster.ErrorClusterSetting, "name", u.ID, err)
	}

	if err := d.Set("cloud_backup", u.ProviderBackupEnabled); err != nil {
		return nil, fmt.Errorf("couldn't import cluster backup configuration %s in project %s, error: %s", *name, *projectID, err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"cluster_id":    u.ID,
		"project_id":    *projectID,
		"cluster_name":  u.Name,
		"provider_name": u.ProviderSettings.ProviderName,
	}))

	return []*schema.ResourceData{d}, nil
}

func splitSClusterImportID(id string) (projectID, clusterName *string, err error) {
	var re = regexp.MustCompile(`(?s)^([0-9a-fA-F]{24})-(.*)$`)
	parts := re.FindStringSubmatch(id)

	if len(parts) != 3 {
		err = errors.New("import format error: to import a cluster, use the format {project_id}-{name}")
		return
	}

	projectID = &parts[1]
	clusterName = &parts[2]

	return
}

func getInstanceSizeToInt(instanceSize string) int {
	regex := regexp.MustCompile(`\d+`)
	num := regex.FindString(instanceSize)

	return cast.ToInt(num) // if the string is empty it always return 0
}

func expandProviderSetting(d *schema.ResourceData) (*matlas.ProviderSettings, error) {
	var (
		region, _          = conversion.ValRegion(d.Get("provider_region_name"))
		minInstanceSize    = getInstanceSizeToInt(d.Get("provider_auto_scaling_compute_min_instance_size").(string))
		maxInstanceSize    = getInstanceSizeToInt(d.Get("provider_auto_scaling_compute_max_instance_size").(string))
		instanceSize       = getInstanceSizeToInt(d.Get("provider_instance_size_name").(string))
		compute            *matlas.Compute
		autoScalingEnabled = d.Get("auto_scaling_compute_enabled").(bool)
		providerName       = cast.ToString(d.Get("provider_name"))
	)

	if minInstanceSize != 0 && autoScalingEnabled {
		if instanceSize < minInstanceSize {
			return nil, fmt.Errorf("`provider_auto_scaling_compute_min_instance_size` must be lower than `provider_instance_size_name`")
		}

		compute = &matlas.Compute{
			MinInstanceSize: d.Get("provider_auto_scaling_compute_min_instance_size").(string),
		}
	}

	if maxInstanceSize != 0 && autoScalingEnabled {
		if instanceSize > maxInstanceSize {
			return nil, fmt.Errorf("`provider_auto_scaling_compute_max_instance_size` must be higher than `provider_instance_size_name`")
		}

		if compute == nil {
			compute = &matlas.Compute{}
		}
		compute.MaxInstanceSize = d.Get("provider_auto_scaling_compute_max_instance_size").(string)
	}

	providerSettings := &matlas.ProviderSettings{
		InstanceSizeName: cast.ToString(d.Get("provider_instance_size_name")),
		ProviderName:     providerName,
		RegionName:       region,
		VolumeType:       cast.ToString(d.Get("provider_volume_type")),
	}

	if d.HasChange("provider_disk_type_name") {
		_, newdiskTypeName := d.GetChange("provider_disk_type_name")
		diskTypeName := cast.ToString(newdiskTypeName)
		if diskTypeName != "" { // ensure disk type is not included in request if attribute is removed, prevents errors in NVME intances
			providerSettings.DiskTypeName = diskTypeName
		}
	}

	if providerName == "TENANT" {
		providerSettings.BackingProviderName = cast.ToString(d.Get("backing_provider_name"))
	}

	if autoScalingEnabled {
		providerSettings.AutoScaling = &matlas.AutoScaling{Compute: compute}
	}

	if d.Get("provider_name") == "AWS" {
		// Check if the Provider Disk IOS sets in the Terraform configuration and if the instance size name is not NVME.
		// If it didn't, the MongoDB Atlas server would set it to the default for the amount of storage.
		if v, ok := d.GetOk("provider_disk_iops"); ok && !strings.Contains(providerSettings.InstanceSizeName, "NVME") {
			providerSettings.DiskIOPS = pointy.Int64(cast.ToInt64(v))
		}

		providerSettings.EncryptEBSVolume = pointy.Bool(true)
	}

	return providerSettings, nil
}

func flattenProviderSettings(d *schema.ResourceData, settings *matlas.ProviderSettings, clusterName string) {
	if settings.ProviderName == "TENANT" {
		if err := d.Set("backing_provider_name", settings.BackingProviderName); err != nil {
			log.Printf(advancedcluster.ErrorClusterSetting, "backing_provider_name", clusterName, err)
		}
	}

	if settings.DiskIOPS != nil && *settings.DiskIOPS != 0 {
		if err := d.Set("provider_disk_iops", *settings.DiskIOPS); err != nil {
			log.Printf(advancedcluster.ErrorClusterSetting, "provider_disk_iops", clusterName, err)
		}
	}

	if err := d.Set("provider_disk_type_name", settings.DiskTypeName); err != nil {
		log.Printf(advancedcluster.ErrorClusterSetting, "provider_disk_type_name", clusterName, err)
	}

	if settings.EncryptEBSVolume != nil {
		if err := d.Set("provider_encrypt_ebs_volume_flag", *settings.EncryptEBSVolume); err != nil {
			log.Printf(advancedcluster.ErrorClusterSetting, "provider_encrypt_ebs_volume_flag", clusterName, err)
		}
	}

	if err := d.Set("provider_instance_size_name", settings.InstanceSizeName); err != nil {
		log.Printf(advancedcluster.ErrorClusterSetting, "provider_instance_size_name", clusterName, err)
	}

	if err := d.Set("provider_name", settings.ProviderName); err != nil {
		log.Printf(advancedcluster.ErrorClusterSetting, "provider_name", clusterName, err)
	}

	if err := d.Set("provider_region_name", settings.RegionName); err != nil {
		log.Printf(advancedcluster.ErrorClusterSetting, "provider_region_name", clusterName, err)
	}

	if err := d.Set("provider_volume_type", settings.VolumeType); err != nil {
		log.Printf(advancedcluster.ErrorClusterSetting, "provider_volume_type", clusterName, err)
	}
}

func isUpgradeRequired(d *schema.ResourceData) bool {
	currentSize, updatedSize := d.GetChange("provider_instance_size_name")

	return currentSize != updatedSize && advancedcluster.IsSharedTier(currentSize.(string))
}

func expandReplicationSpecs(d *schema.ResourceData) ([]matlas.ReplicationSpec, error) {
	rSpecs := make([]matlas.ReplicationSpec, 0)

	vRSpecs, okRSpecs := d.GetOk("replication_specs")
	vPRName, okPRName := d.GetOk("provider_region_name")

	if okRSpecs {
		for _, s := range vRSpecs.(*schema.Set).List() {
			spec := s.(map[string]any)

			replaceRegion := ""
			originalRegion := ""
			id := ""

			if okPRName && d.Get("provider_name").(string) == "GCP" && cast.ToString(d.Get("cluster_type")) == "REPLICASET" {
				if d.HasChange("provider_region_name") {
					replaceRegion = vPRName.(string)
					original, _ := d.GetChange("provider_region_name")
					originalRegion = original.(string)
				}
			}

			if d.HasChange("replication_specs") {
				// Get original and new object
				var oldSpecs map[string]any
				original, _ := d.GetChange("replication_specs")
				for _, s := range original.(*schema.Set).List() {
					oldSpecs = s.(map[string]any)
					if spec["zone_name"].(string) == cast.ToString(oldSpecs["zone_name"]) {
						id = oldSpecs["id"].(string)
						break
					}
				}
				// If there was an item before and after then use the same id assuming it's the same replication spec
				if id == "" && oldSpecs != nil && len(vRSpecs.(*schema.Set).List()) == 1 && len(original.(*schema.Set).List()) == 1 {
					id = oldSpecs["id"].(string)
				}
			}

			regionsConfig, err := expandRegionsConfig(spec["regions_config"].(*schema.Set).List(), originalRegion, replaceRegion)
			if err != nil {
				return rSpecs, err
			}

			rSpec := matlas.ReplicationSpec{
				ID:            id,
				NumShards:     pointy.Int64(cast.ToInt64(spec["num_shards"])),
				ZoneName:      cast.ToString(spec["zone_name"]),
				RegionsConfig: regionsConfig,
			}
			rSpecs = append(rSpecs, rSpec)
		}
	}

	return rSpecs, nil
}

func flattenReplicationSpecs(rSpecs []matlas.ReplicationSpec) []map[string]any {
	specs := make([]map[string]any, 0)

	for _, rSpec := range rSpecs {
		spec := map[string]any{
			"id":             rSpec.ID,
			"num_shards":     rSpec.NumShards,
			"zone_name":      cast.ToString(rSpec.ZoneName),
			"regions_config": flattenRegionsConfig(rSpec.RegionsConfig),
		}
		specs = append(specs, spec)
	}

	return specs
}

func expandRegionsConfig(regions []any, originalRegion, replaceRegion string) (map[string]matlas.RegionsConfig, error) {
	regionsConfig := make(map[string]matlas.RegionsConfig)

	for _, r := range regions {
		region := r.(map[string]any)

		r, err := conversion.ValRegion(region["region_name"])
		if err != nil {
			return regionsConfig, err
		}

		if replaceRegion != "" && r == originalRegion {
			r, err = conversion.ValRegion(replaceRegion)
		}
		if err != nil {
			return regionsConfig, err
		}

		regionsConfig[r] = matlas.RegionsConfig{
			AnalyticsNodes: pointy.Int64(cast.ToInt64(region["analytics_nodes"])),
			ElectableNodes: pointy.Int64(cast.ToInt64(region["electable_nodes"])),
			Priority:       pointy.Int64(cast.ToInt64(region["priority"])),
			ReadOnlyNodes:  pointy.Int64(cast.ToInt64(region["read_only_nodes"])),
		}
	}

	return regionsConfig, nil
}

func flattenRegionsConfig(regionsConfig map[string]matlas.RegionsConfig) []map[string]any {
	regions := make([]map[string]any, 0)

	for regionName, regionConfig := range regionsConfig {
		region := map[string]any{
			"region_name":     regionName,
			"priority":        regionConfig.Priority,
			"analytics_nodes": regionConfig.AnalyticsNodes,
			"electable_nodes": regionConfig.ElectableNodes,
			"read_only_nodes": regionConfig.ReadOnlyNodes,
		}
		regions = append(regions, region)
	}

	return regions
}

func resourceClusterCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, meta any) error {
	var err error
	currentProvider, updatedProvider := d.GetChange("provider_name")

	willProviderChange := currentProvider != updatedProvider
	willLeaveTenant := willProviderChange && currentProvider == "TENANT"

	if willLeaveTenant {
		err = d.SetNewComputed("backing_provider_name")
	} else if willProviderChange {
		err = d.ForceNew("provider_name")
	}

	return err
}

func getContainerID(containers []matlas.Container, cluster *matlas.Cluster) string {
	if len(containers) != 0 {
		for i := range containers {
			if cluster.ProviderSettings.ProviderName == "GCP" {
				return containers[i].ID
			}

			if containers[i].ProviderName == cluster.ProviderSettings.ProviderName &&
				containers[i].Region == cluster.ProviderSettings.RegionName || // For Azure
				containers[i].RegionName == cluster.ProviderSettings.RegionName { // For AWS
				return containers[i].ID
			}
		}
	}

	return ""
}

func isEqualProviderAutoScalingMinInstanceSize(k, old, newStr string, d *schema.ResourceData) bool {
	canScaleDown, scaleDownOK := d.GetOk("auto_scaling_compute_scale_down_enabled")
	canScaleUp, scaleUpOk := d.GetOk("auto_scaling_compute_enabled")

	if !scaleDownOK || !scaleUpOk {
		return true // if the return is true, it means that both values are the same and there's nothing to do
	}

	if canScaleUp.(bool) && canScaleDown.(bool) {
		if old != newStr {
			return false
		}
	}
	return true
}

func isEqualProviderAutoScalingMaxInstanceSize(k, old, newStr string, d *schema.ResourceData) bool {
	canScaleUp, _ := d.GetOk("auto_scaling_compute_enabled")
	if canScaleUp != nil && canScaleUp.(bool) {
		if old != newStr {
			return false
		}
	}
	return true
}

func updateCluster(ctx context.Context, conn *matlas.Client, request *matlas.Cluster, projectID, name string, timeout time.Duration) (*matlas.Cluster, *matlas.Response, error) {
	cluster, resp, err := conn.Clusters.Update(ctx, projectID, name, request)
	if err != nil {
		return nil, nil, err
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"CREATING", "UPDATING", "REPAIRING"},
		Target:     []string{"IDLE"},
		Refresh:    advancedcluster.ResourceClusterRefreshFunc(ctx, name, projectID, conn),
		Timeout:    timeout,
		MinTimeout: 30 * time.Second,
		Delay:      1 * time.Minute,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, nil, err
	}

	return cluster, resp, nil
}

func computedCloudProviderSnapshotBackupPolicySchema() *schema.Schema {
	return &schema.Schema{
		Type:       schema.TypeList,
		Computed:   true,
		ConfigMode: schema.SchemaConfigModeAttr,
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
					Type:       schema.TypeList,
					Computed:   true,
					ConfigMode: schema.SchemaConfigModeAttr,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"id": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"policy_item": {
								Type:       schema.TypeList,
								Computed:   true,
								ConfigMode: schema.SchemaConfigModeAttr,
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

func flattenCloudProviderSnapshotBackupPolicy(ctx context.Context, d *schema.ResourceData, conn *matlas.Client, projectID, clusterName string) ([]map[string]any, error) {
	backupPolicy, res, err := conn.CloudProviderSnapshotBackupPolicies.Get(ctx, projectID, clusterName)
	if err != nil {
		if res.StatusCode == http.StatusNotFound ||
			strings.Contains(err.Error(), "BACKUP_CONFIG_NOT_FOUND") ||
			strings.Contains(err.Error(), "Not Found") ||
			strings.Contains(err.Error(), "404") {
			return []map[string]any{}, nil
		}

		return []map[string]any{}, fmt.Errorf(ErrorSnapshotBackupPolicyRead, clusterName, err)
	}

	return []map[string]any{
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

func flattenPolicies(policies []matlas.Policy) []map[string]any {
	actionList := make([]map[string]any, 0)
	for _, v := range policies {
		actionList = append(actionList, map[string]any{
			"id":          v.ID,
			"policy_item": flattenPolicyItems(v.PolicyItems),
		})
	}

	return actionList
}

func flattenPolicyItems(items []matlas.PolicyItem) []map[string]any {
	policyItems := make([]map[string]any, 0)
	for _, v := range items {
		policyItems = append(policyItems, map[string]any{
			"id":                 v.ID,
			"frequency_interval": v.FrequencyInterval,
			"frequency_type":     v.FrequencyType,
			"retention_unit":     v.RetentionUnit,
			"retention_value":    v.RetentionValue,
		})
	}

	return policyItems
}
