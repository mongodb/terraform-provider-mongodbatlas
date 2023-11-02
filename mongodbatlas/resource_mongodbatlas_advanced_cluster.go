package mongodbatlas

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
	"golang.org/x/exp/slices"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mwielbut/pointy"
	"github.com/spf13/cast"
)

type acCtxKey string

const (
	errorClusterAdvancedCreate             = "error creating MongoDB ClusterAdvanced: %s"
	errorClusterAdvancedRead               = "error reading MongoDB ClusterAdvanced (%s): %s"
	errorClusterAdvancedDelete             = "error deleting MongoDB ClusterAdvanced (%s): %s"
	errorClusterAdvancedUpdate             = "error updating MongoDB ClusterAdvanced (%s): %s"
	errorClusterAdvancedSetting            = "error setting `%s` for MongoDB ClusterAdvanced (%s): %s"
	errorAdvancedClusterAdvancedConfUpdate = "error updating Advanced Configuration Option form MongoDB Cluster (%s): %s"
	errorAdvancedClusterAdvancedConfRead   = "error reading Advanced Configuration Option form MongoDB Cluster (%s): %s"
	errorAdvancedClusterListStatus         = "error awaiting MongoDB ClusterAdvanced List IDLE: %s"
)

var upgradeRequestCtxKey acCtxKey = "upgradeRequest"

func resourceMongoDBAtlasAdvancedCluster() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceMongoDBAtlasAdvancedClusterCreate,
		ReadWithoutTimeout:   resourceMongoDBAtlasAdvancedClusterRead,
		UpdateWithoutTimeout: resourceMongoDBAtlasAdvancedClusterUpdateOrUpgrade,
		DeleteWithoutTimeout: resourceMongoDBAtlasAdvancedClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasAdvancedClusterImportState,
		},
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceMongoDBAtlasAdvancedClusterResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceMongoDBAtlasAdvancedClusterStateUpgradeV0,
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
			"backup_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"retain_backups_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Flag that indicates whether to retain backup snapshots for the deleted dedicated cluster",
			},
			"bi_connector_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
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
				Required: true,
			},
			"connection_strings": clusterConnectionStringsSchema(),
			"create_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
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
			"labels": {
				Type:       schema.TypeSet,
				Optional:   true,
				Set:        HashFunctionForKeyValuePair,
				Deprecated: fmt.Sprintf(DeprecationByDateWithReplacement, "September 2024", "tags"),
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
			"tags": &tagsSchema,
			"mongo_db_major_version": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				StateFunc: formatMongoDBMajorVersion,
			},
			"mongo_db_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"paused": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"pit_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"replication_specs": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"num_shards": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: validation.IntBetween(1, 50),
						},
						"region_configs": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"analytics_specs": advancedClusterRegionConfigsSpecsSchema(),
									"auto_scaling": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"disk_gb_enabled": {
													Type:     schema.TypeBool,
													Optional: true,
													Computed: true,
												},
												"compute_enabled": {
													Type:     schema.TypeBool,
													Optional: true,
													Computed: true,
												},
												"compute_scale_down_enabled": {
													Type:     schema.TypeBool,
													Optional: true,
													Computed: true,
												},
												"compute_min_instance_size": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"compute_max_instance_size": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
											},
										},
									},
									"analytics_auto_scaling": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"disk_gb_enabled": {
													Type:     schema.TypeBool,
													Optional: true,
													Computed: true,
												},
												"compute_enabled": {
													Type:     schema.TypeBool,
													Optional: true,
													Computed: true,
												},
												"compute_scale_down_enabled": {
													Type:     schema.TypeBool,
													Optional: true,
													Computed: true,
												},
												"compute_min_instance_size": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"compute_max_instance_size": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
											},
										},
									},
									"backing_provider_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"electable_specs": advancedClusterRegionConfigsSpecsSchema(),
									"priority": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"provider_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"read_only_specs": advancedClusterRegionConfigsSpecsSchema(),
									"region_name": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"container_id": {
							Type: schema.TypeMap,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed: true,
						},
						"zone_name": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "ZoneName managed by Terraform",
						},
					},
				},
				// Set: replicationSpecsHashSet,
			},
			"root_cert_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"state_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"termination_protection_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"version_release_system": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"LTS", "CONTINUOUS"}, false),
			},
			"advanced_configuration": clusterAdvancedConfigurationSchema(),
			"accept_data_risks_and_force_replica_set_reconfig": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Submit this field alongside your topology reconfiguration to request a new regional outage resistant topology",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Hour),
			Update: schema.DefaultTimeout(3 * time.Hour),
			Delete: schema.DefaultTimeout(3 * time.Hour),
		},
	}
}

var tagsSchema = schema.Schema{
	Type:     schema.TypeSet,
	Optional: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	},
}

func HashFunctionForKeyValuePair(v any) int {
	var buf bytes.Buffer
	m := v.(map[string]any)
	buf.WriteString(m["key"].(string))
	buf.WriteString(m["value"].(string))
	return HashCodeString(buf.String())
}

func advancedClusterRegionConfigsSpecsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"disk_iops": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"ebs_volume_type": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"instance_size": {
					Type:     schema.TypeString,
					Required: true,
				},
				"node_count": {
					Type:     schema.TypeInt,
					Optional: true,
				},
			},
		},
	}
}

func resourceMongoDBAtlasAdvancedClusterCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	if v, ok := d.GetOk("accept_data_risks_and_force_replica_set_reconfig"); ok {
		if v.(string) != "" {
			return diag.FromErr(fmt.Errorf("accept_data_risks_and_force_replica_set_reconfig can not be set in creation, only in update"))
		}
	}

	conn := meta.(*MongoDBClient).Atlas

	projectID := d.Get("project_id").(string)

	request := &matlas.AdvancedCluster{
		Name:             d.Get("name").(string),
		ClusterType:      cast.ToString(d.Get("cluster_type")),
		ReplicationSpecs: expandAdvancedReplicationSpecs(d.Get("replication_specs").([]any)),
	}

	if v, ok := d.GetOk("backup_enabled"); ok {
		request.BackupEnabled = pointy.Bool(v.(bool))
	}
	if _, ok := d.GetOk("bi_connector_config"); ok {
		biConnector, err := expandBiConnectorConfig(d)
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorClusterAdvancedCreate, err))
		}
		request.BiConnector = biConnector
	}
	if v, ok := d.GetOk("disk_size_gb"); ok {
		request.DiskSizeGB = pointy.Float64(v.(float64))
	}
	if v, ok := d.GetOk("encryption_at_rest_provider"); ok {
		request.EncryptionAtRestProvider = v.(string)
	}

	if _, ok := d.GetOk("labels"); ok && containsLabelOrKey(expandLabelSliceFromSetSchema(d), defaultLabel) {
		return diag.FromErr(fmt.Errorf("you should not set `Infrastructure Tool` label, it is used for internal purposes"))
	}
	request.Labels = append(expandLabelSliceFromSetSchema(d), defaultLabel)

	if _, ok := d.GetOk("tags"); ok {
		request.Tags = expandTagSliceFromSetSchema(d)
	}
	if v, ok := d.GetOk("mongo_db_major_version"); ok {
		request.MongoDBMajorVersion = formatMongoDBMajorVersion(v.(string))
	}
	if v, ok := d.GetOk("pit_enabled"); ok {
		request.PitEnabled = pointy.Bool(v.(bool))
	}
	if v, ok := d.GetOk("root_cert_type"); ok {
		request.RootCertType = v.(string)
	}
	if v, ok := d.GetOk("termination_protection_enabled"); ok {
		request.TerminationProtectionEnabled = pointy.Bool(v.(bool))
	}
	if v, ok := d.GetOk("version_release_system"); ok {
		request.VersionReleaseSystem = v.(string)
	}

	// We need to validate the oplog_size_mb attr of the advanced configuration option to show the error
	// before that the cluster is created
	if oplogSizeMB, ok := d.GetOkExists("advanced_configuration.0.oplog_size_mb"); ok {
		if cast.ToInt64(oplogSizeMB) <= 0 {
			return diag.FromErr(fmt.Errorf("`advanced_configuration.oplog_size_mb` cannot be <= 0"))
		}
	}

	cluster, _, err := conn.AdvancedClusters.Create(ctx, projectID, request)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterAdvancedCreate, err))
	}

	timeout := d.Timeout(schema.TimeoutCreate)
	stateConf := &retry.StateChangeConf{
		Pending:    []string{"CREATING", "UPDATING", "REPAIRING", "REPEATING", "PENDING"},
		Target:     []string{"IDLE"},
		Refresh:    resourceClusterAdvancedRefreshFunc(ctx, d.Get("name").(string), projectID, conn),
		Timeout:    timeout,
		MinTimeout: 1 * time.Minute,
		Delay:      3 * time.Minute,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterAdvancedCreate, err))
	}

	/*
		So far, the cluster has created correctly, so we need to set up
		the advanced configuration option to attach it
	*/
	ac, ok := d.GetOk("advanced_configuration")
	if aclist, ok1 := ac.([]any); ok1 && len(aclist) > 0 {
		advancedConfReq := expandProcessArgs(d, aclist[0].(map[string]any))

		if ok {
			_, _, err := conn.Clusters.UpdateProcessArgs(ctx, projectID, cluster.Name, advancedConfReq)
			if err != nil {
				return diag.FromErr(fmt.Errorf(errorAdvancedClusterAdvancedConfUpdate, cluster.Name, err))
			}
		}
	}

	// To pause a cluster
	if v := d.Get("paused").(bool); v {
		request = &matlas.AdvancedCluster{
			Paused: pointy.Bool(v),
		}

		_, _, err = updateAdvancedCluster(ctx, conn, request, projectID, d.Get("name").(string), timeout)
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorClusterAdvancedUpdate, d.Get("name").(string), err))
		}
	}

	d.SetId(encodeStateID(map[string]string{
		"cluster_id":   cluster.ID,
		"project_id":   projectID,
		"cluster_name": cluster.Name,
	}))

	return resourceMongoDBAtlasAdvancedClusterRead(ctx, d, meta)
}

func resourceMongoDBAtlasAdvancedClusterRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	cluster, resp, err := conn.AdvancedClusters.Get(ctx, projectID, clusterName)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorClusterAdvancedRead, clusterName, err))
	}

	log.Printf("[DEBUG] GET ClusterAdvanced %+v", cluster)

	if err := d.Set("cluster_id", cluster.ID); err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterAdvancedSetting, "cluster_id", clusterName, err))
	}

	if err := d.Set("backup_enabled", cluster.BackupEnabled); err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterAdvancedSetting, "backup_enabled", clusterName, err))
	}

	if err := d.Set("bi_connector_config", flattenBiConnectorConfig(cluster.BiConnector)); err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterAdvancedSetting, "bi_connector_config", clusterName, err))
	}

	if err := d.Set("cluster_type", cluster.ClusterType); err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterAdvancedSetting, "cluster_type", clusterName, err))
	}

	if err := d.Set("connection_strings", flattenConnectionStrings(cluster.ConnectionStrings)); err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterAdvancedSetting, "connection_strings", clusterName, err))
	}

	if err := d.Set("create_date", cluster.CreateDate); err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterAdvancedSetting, "create_date", clusterName, err))
	}

	if err := d.Set("disk_size_gb", cluster.DiskSizeGB); err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterAdvancedSetting, "disk_size_gb", clusterName, err))
	}

	if err := d.Set("encryption_at_rest_provider", cluster.EncryptionAtRestProvider); err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterAdvancedSetting, "encryption_at_rest_provider", clusterName, err))
	}

	if err := d.Set("labels", flattenLabels(removeLabel(cluster.Labels, defaultLabel))); err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterAdvancedSetting, "labels", clusterName, err))
	}

	if err := d.Set("tags", flattenTags(&cluster.Tags)); err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterAdvancedSetting, "tags", clusterName, err))
	}

	if err := d.Set("mongo_db_major_version", cluster.MongoDBMajorVersion); err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterAdvancedSetting, "mongo_db_major_version", clusterName, err))
	}

	if err := d.Set("mongo_db_version", cluster.MongoDBVersion); err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterAdvancedSetting, "mongo_db_version", clusterName, err))
	}

	if err := d.Set("name", cluster.Name); err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterAdvancedSetting, "name", clusterName, err))
	}

	if err := d.Set("paused", cluster.Paused); err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterAdvancedSetting, "paused", clusterName, err))
	}

	if err := d.Set("pit_enabled", cluster.PitEnabled); err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterAdvancedSetting, "pit_enabled", clusterName, err))
	}

	replicationSpecs, err := flattenAdvancedReplicationSpecs(ctx, cluster.ReplicationSpecs, d.Get("replication_specs").([]any), d, conn)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterAdvancedSetting, "replication_specs", clusterName, err))
	}

	if err := d.Set("replication_specs", replicationSpecs); err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterAdvancedSetting, "replication_specs", clusterName, err))
	}

	if err := d.Set("root_cert_type", cluster.RootCertType); err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterAdvancedSetting, "state_name", clusterName, err))
	}

	if err := d.Set("state_name", cluster.StateName); err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterAdvancedSetting, "state_name", clusterName, err))
	}

	if err := d.Set("termination_protection_enabled", cluster.TerminationProtectionEnabled); err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterAdvancedSetting, "termination_protection_enabled", clusterName, err))
	}

	if err := d.Set("version_release_system", cluster.VersionReleaseSystem); err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterAdvancedSetting, "version_release_system", clusterName, err))
	}

	if err := d.Set("accept_data_risks_and_force_replica_set_reconfig", cluster.AcceptDataRisksAndForceReplicaSetReconfig); err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterAdvancedSetting, "accept_data_risks_and_force_replica_set_reconfig", clusterName, err))
	}

	/*
		Get the advaced configuration options and set up to the terraform state
	*/
	processArgs, _, err := conn.Clusters.GetProcessArgs(ctx, projectID, clusterName)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorAdvancedClusterAdvancedConfRead, clusterName, err))
	}

	if err := d.Set("advanced_configuration", flattenProcessArgs(processArgs)); err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterAdvancedSetting, "advanced_configuration", clusterName, err))
	}

	return nil
}

func resourceMongoDBAtlasAdvancedClusterUpdateOrUpgrade(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	if upgradeRequest := getUpgradeRequest(d); upgradeRequest != nil {
		upgradeCtx := context.WithValue(ctx, upgradeRequestCtxKey, upgradeRequest)
		return resourceMongoDBAtlasAdvancedClusterUpgrade(upgradeCtx, d, meta)
	}

	return resourceMongoDBAtlasAdvancedClusterUpdate(ctx, d, meta)
}

func resourceMongoDBAtlasAdvancedClusterUpgrade(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	upgradeRequest := ctx.Value(upgradeRequestCtxKey).(*matlas.Cluster)

	if upgradeRequest == nil {
		return diag.FromErr(fmt.Errorf("upgrade called without %s in ctx", string(upgradeRequestCtxKey)))
	}

	upgradeResponse, _, err := upgradeCluster(ctx, conn, upgradeRequest, projectID, clusterName, d.Timeout(schema.TimeoutUpdate))

	if err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterAdvancedUpdate, clusterName, err))
	}

	d.SetId(encodeStateID(map[string]string{
		"cluster_id":   upgradeResponse.ID,
		"project_id":   projectID,
		"cluster_name": clusterName,
	}))

	return resourceMongoDBAtlasAdvancedClusterRead(ctx, d, meta)
}

func resourceMongoDBAtlasAdvancedClusterUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	cluster := new(matlas.AdvancedCluster)
	clusterChangeDetect := new(matlas.AdvancedCluster)

	if d.HasChange("backup_enabled") {
		cluster.BackupEnabled = pointy.Bool(d.Get("backup_enabled").(bool))
	}

	if d.HasChange("bi_connector_config") {
		cluster.BiConnector, _ = expandBiConnectorConfig(d)
	}

	if d.HasChange("cluster_type") {
		cluster.ClusterType = d.Get("cluster_type").(string)
	}

	if d.HasChange("disk_size_gb") {
		cluster.DiskSizeGB = pointy.Float64(d.Get("disk_size_gb").(float64))
	}

	if d.HasChange("encryption_at_rest_provider") {
		cluster.EncryptionAtRestProvider = d.Get("encryption_at_rest_provider").(string)
	}

	if d.HasChange("labels") {
		if containsLabelOrKey(expandLabelSliceFromSetSchema(d), defaultLabel) {
			return diag.FromErr(fmt.Errorf("you should not set `Infrastructure Tool` label, it is used for internal purposes"))
		}

		cluster.Labels = append(expandLabelSliceFromSetSchema(d), defaultLabel)
	}

	if d.HasChange("tags") {
		cluster.Tags = expandTagSliceFromSetSchema(d)
	}

	if d.HasChange("mongo_db_major_version") {
		cluster.MongoDBMajorVersion = formatMongoDBMajorVersion(d.Get("mongo_db_major_version"))
	}

	if d.HasChange("pit_enabled") {
		cluster.PitEnabled = pointy.Bool(d.Get("pit_enabled").(bool))
	}

	if d.HasChange("replication_specs") {
		cluster.ReplicationSpecs = expandAdvancedReplicationSpecs(d.Get("replication_specs").([]any))
	}

	if d.HasChange("root_cert_type") {
		cluster.RootCertType = d.Get("root_cert_type").(string)
	}

	if d.HasChange("termination_protection_enabled") {
		cluster.TerminationProtectionEnabled = pointy.Bool(d.Get("termination_protection_enabled").(bool))
	}

	if d.HasChange("version_release_system") {
		cluster.VersionReleaseSystem = d.Get("version_release_system").(string)
	}

	if d.HasChange("accept_data_risks_and_force_replica_set_reconfig") {
		cluster.AcceptDataRisksAndForceReplicaSetReconfig = d.Get("accept_data_risks_and_force_replica_set_reconfig").(string)
	}

	if d.HasChange("paused") && !d.Get("paused").(bool) {
		cluster.Paused = pointy.Bool(d.Get("paused").(bool))
	}

	timeout := d.Timeout(schema.TimeoutUpdate)

	if d.HasChange("advanced_configuration") {
		ac := d.Get("advanced_configuration")
		if aclist, ok := ac.([]any); ok && len(aclist) > 0 {
			advancedConfReq := expandProcessArgs(d, aclist[0].(map[string]any))
			if !reflect.DeepEqual(advancedConfReq, matlas.ProcessArgs{}) {
				_, _, err := conn.Clusters.UpdateProcessArgs(ctx, projectID, clusterName, advancedConfReq)
				if err != nil {
					return diag.FromErr(fmt.Errorf(errorAdvancedClusterAdvancedConfUpdate, clusterName, err))
				}
			}
		}
	}

	// Has changes
	if !reflect.DeepEqual(cluster, clusterChangeDetect) {
		err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
			_, _, err := updateAdvancedCluster(ctx, conn, cluster, projectID, clusterName, timeout)
			if err != nil {
				var target *matlas.ErrorResponse
				if errors.As(err, &target) && target.ErrorCode == "CANNOT_UPDATE_PAUSED_CLUSTER" {
					clusterRequest := &matlas.AdvancedCluster{
						Paused: pointy.Bool(false),
					}
					_, _, err := updateAdvancedCluster(ctx, conn, clusterRequest, projectID, clusterName, timeout)
					if err != nil {
						return retry.NonRetryableError(fmt.Errorf(errorClusterAdvancedUpdate, clusterName, err))
					}
				}
				if errors.As(err, &target) && target.HTTPCode == 400 {
					return retry.NonRetryableError(fmt.Errorf(errorClusterAdvancedUpdate, clusterName, err))
				}
			}
			return nil
		})
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorClusterAdvancedUpdate, clusterName, err))
		}
	}

	if d.Get("paused").(bool) {
		clusterRequest := &matlas.AdvancedCluster{
			Paused: pointy.Bool(true),
		}

		_, _, err := updateAdvancedCluster(ctx, conn, clusterRequest, projectID, clusterName, timeout)
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorClusterAdvancedUpdate, clusterName, err))
		}
	}

	return resourceMongoDBAtlasAdvancedClusterRead(ctx, d, meta)
}

func resourceMongoDBAtlasAdvancedClusterDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	var options *matlas.DeleteAdvanceClusterOptions
	if v, ok := d.GetOkExists("retain_backups_enabled"); ok {
		options = &matlas.DeleteAdvanceClusterOptions{
			RetainBackups: pointy.Bool(v.(bool)),
		}
	}

	_, err := conn.AdvancedClusters.Delete(ctx, projectID, clusterName, options)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterAdvancedDelete, clusterName, err))
	}

	log.Println("[INFO] Waiting for MongoDB ClusterAdvanced to be destroyed")

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"IDLE", "CREATING", "UPDATING", "REPAIRING", "DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    resourceClusterAdvancedRefreshFunc(ctx, clusterName, projectID, conn),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		MinTimeout: 30 * time.Second,
		Delay:      1 * time.Minute, // Wait 30 secs before starting
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterAdvancedDelete, clusterName, err))
	}

	return nil
}

func resourceMongoDBAtlasAdvancedClusterImportState(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	conn := meta.(*MongoDBClient).Atlas

	projectID, name, err := splitSClusterAdvancedImportID(d.Id())
	if err != nil {
		return nil, err
	}

	u, _, err := conn.AdvancedClusters.Get(ctx, *projectID, *name)
	if err != nil {
		return nil, fmt.Errorf("couldn't import cluster %s in project %s, error: %s", *name, *projectID, err)
	}

	if err := d.Set("project_id", u.GroupID); err != nil {
		log.Printf(errorClusterAdvancedSetting, "project_id", u.ID, err)
	}

	if err := d.Set("name", u.Name); err != nil {
		log.Printf(errorClusterAdvancedSetting, "name", u.ID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"cluster_id":   u.ID,
		"project_id":   *projectID,
		"cluster_name": u.Name,
	}))

	return []*schema.ResourceData{d}, nil
}

func splitSClusterAdvancedImportID(id string) (projectID, clusterName *string, err error) {
	var re = regexp.MustCompile(`(?s)^([0-9a-fA-F]{24})-(.*)$`)
	parts := re.FindStringSubmatch(id)

	if len(parts) != 3 {
		err = errors.New("import format error: to import a advanced cluster, use the format {project_id}-{name}")
		return
	}

	projectID = &parts[1]
	clusterName = &parts[2]

	return
}

func expandAdvancedReplicationSpec(tfMap map[string]any) *matlas.AdvancedReplicationSpec {
	if tfMap == nil {
		return nil
	}

	apiObject := &matlas.AdvancedReplicationSpec{
		NumShards:     tfMap["num_shards"].(int),
		ZoneName:      tfMap["zone_name"].(string),
		RegionConfigs: expandRegionConfigs(tfMap["region_configs"].([]any)),
	}

	return apiObject
}

func expandAdvancedReplicationSpecs(tfList []any) []*matlas.AdvancedReplicationSpec {
	if len(tfList) == 0 {
		return nil
	}

	var apiObjects []*matlas.AdvancedReplicationSpec

	for _, tfMapRaw := range tfList {
		tfMap, ok := tfMapRaw.(map[string]any)

		if !ok {
			continue
		}

		apiObject := expandAdvancedReplicationSpec(tfMap)

		apiObjects = append(apiObjects, apiObject)
	}

	return apiObjects
}

func expandRegionConfig(tfMap map[string]any) *matlas.AdvancedRegionConfig {
	if tfMap == nil {
		return nil
	}

	providerName := tfMap["provider_name"].(string)
	apiObject := &matlas.AdvancedRegionConfig{
		Priority:     pointy.Int(cast.ToInt(tfMap["priority"])),
		ProviderName: providerName,
		RegionName:   tfMap["region_name"].(string),
	}

	if v, ok := tfMap["analytics_specs"]; ok && len(v.([]any)) > 0 {
		apiObject.AnalyticsSpecs = expandRegionConfigSpec(v.([]any), providerName)
	}
	if v, ok := tfMap["electable_specs"]; ok && len(v.([]any)) > 0 {
		apiObject.ElectableSpecs = expandRegionConfigSpec(v.([]any), providerName)
	}
	if v, ok := tfMap["read_only_specs"]; ok && len(v.([]any)) > 0 {
		apiObject.ReadOnlySpecs = expandRegionConfigSpec(v.([]any), providerName)
	}
	if v, ok := tfMap["auto_scaling"]; ok && len(v.([]any)) > 0 {
		apiObject.AutoScaling = expandRegionConfigAutoScaling(v.([]any))
	}
	if v, ok := tfMap["analytics_auto_scaling"]; ok && len(v.([]any)) > 0 {
		apiObject.AnalyticsAutoScaling = expandRegionConfigAutoScaling(v.([]any))
	}
	if v, ok := tfMap["backing_provider_name"]; ok {
		apiObject.BackingProviderName = v.(string)
	}

	return apiObject
}

func expandRegionConfigs(tfList []any) []*matlas.AdvancedRegionConfig {
	if len(tfList) == 0 {
		return nil
	}

	var apiObjects []*matlas.AdvancedRegionConfig

	for _, tfMapRaw := range tfList {
		tfMap, ok := tfMapRaw.(map[string]any)

		if !ok {
			continue
		}

		apiObject := expandRegionConfig(tfMap)

		apiObjects = append(apiObjects, apiObject)
	}

	return apiObjects
}

func expandRegionConfigSpec(tfList []any, providerName string) *matlas.Specs {
	if tfList == nil && len(tfList) > 0 {
		return nil
	}

	tfMap, _ := tfList[0].(map[string]any)

	apiObject := &matlas.Specs{}

	if providerName == "AWS" {
		if v, ok := tfMap["disk_iops"]; ok && v.(int) > 0 {
			apiObject.DiskIOPS = pointy.Int64(cast.ToInt64(v.(int)))
		}
		if v, ok := tfMap["ebs_volume_type"]; ok {
			apiObject.EbsVolumeType = v.(string)
		}
	}
	if v, ok := tfMap["instance_size"]; ok {
		apiObject.InstanceSize = v.(string)
	}
	if v, ok := tfMap["node_count"]; ok {
		apiObject.NodeCount = pointy.Int(v.(int))
	}

	return apiObject
}

func expandRegionConfigAutoScaling(tfList []any) *matlas.AdvancedAutoScaling {
	if tfList == nil && len(tfList) > 0 {
		return nil
	}

	tfMap, _ := tfList[0].(map[string]any)

	advancedAutoScaling := &matlas.AdvancedAutoScaling{}
	diskGB := &matlas.DiskGB{}
	compute := &matlas.Compute{}

	if v, ok := tfMap["disk_gb_enabled"]; ok {
		diskGB.Enabled = pointy.Bool(v.(bool))
	}
	if v, ok := tfMap["compute_enabled"]; ok {
		compute.Enabled = pointy.Bool(v.(bool))
	}
	if v, ok := tfMap["compute_scale_down_enabled"]; ok {
		compute.ScaleDownEnabled = pointy.Bool(v.(bool))
	}
	if v, ok := tfMap["compute_min_instance_size"]; ok {
		value := compute.ScaleDownEnabled
		if *value {
			compute.MinInstanceSize = v.(string)
		}
	}
	if v, ok := tfMap["compute_max_instance_size"]; ok {
		value := compute.Enabled
		if *value {
			compute.MaxInstanceSize = v.(string)
		}
	}

	advancedAutoScaling.DiskGB = diskGB
	advancedAutoScaling.Compute = compute

	return advancedAutoScaling
}

func flattenAdvancedReplicationSpec(ctx context.Context, apiObject *matlas.AdvancedReplicationSpec, tfMapObject map[string]any,
	d *schema.ResourceData, conn *matlas.Client) (map[string]any, error) {
	if apiObject == nil {
		return nil, nil
	}

	tfMap := map[string]any{}
	tfMap["num_shards"] = apiObject.NumShards
	tfMap["id"] = apiObject.ID
	if tfMapObject != nil {
		object, containerIds, err := flattenAdvancedReplicationSpecRegionConfigs(ctx, apiObject.RegionConfigs, tfMapObject["region_configs"].([]any), d, conn)
		if err != nil {
			return nil, err
		}
		tfMap["region_configs"] = object
		tfMap["container_id"] = containerIds
	} else {
		object, containerIds, err := flattenAdvancedReplicationSpecRegionConfigs(ctx, apiObject.RegionConfigs, nil, d, conn)
		if err != nil {
			return nil, err
		}
		tfMap["region_configs"] = object
		tfMap["container_id"] = containerIds
	}
	tfMap["zone_name"] = apiObject.ZoneName

	return tfMap, nil
}

func doesAdvancedReplicationSpecMatchAPI(tfObject map[string]any, apiObject *matlas.AdvancedReplicationSpec) bool {
	return tfObject["id"] == apiObject.ID || (tfObject["id"] == nil && tfObject["zone_name"] == apiObject.ZoneName)
}

func flattenAdvancedReplicationSpecs(ctx context.Context, rawAPIObjects []*matlas.AdvancedReplicationSpec, tfMapObjects []any,
	d *schema.ResourceData, conn *matlas.Client) ([]map[string]any, error) {
	var apiObjects []*matlas.AdvancedReplicationSpec

	for _, advancedReplicationSpec := range rawAPIObjects {
		if advancedReplicationSpec != nil {
			apiObjects = append(apiObjects, advancedReplicationSpec)
		}
	}

	if len(apiObjects) == 0 {
		return nil, nil
	}

	tfList := make([]map[string]any, len(apiObjects))
	wasAPIObjectUsed := make([]bool, len(apiObjects))

	for i := 0; i < len(tfList); i++ {
		var tfMapObject map[string]any

		if len(tfMapObjects) > i {
			tfMapObject = tfMapObjects[i].(map[string]any)
		}

		for j := 0; j < len(apiObjects); j++ {
			if wasAPIObjectUsed[j] {
				continue
			}

			if !doesAdvancedReplicationSpecMatchAPI(tfMapObject, apiObjects[j]) {
				continue
			}

			advancedReplicationSpec, err := flattenAdvancedReplicationSpec(ctx, apiObjects[j], tfMapObject, d, conn)

			if err != nil {
				return nil, err
			}

			tfList[i] = advancedReplicationSpec
			wasAPIObjectUsed[j] = true
			break
		}
	}

	for i, tfo := range tfList {
		var tfMapObject map[string]any

		if tfo != nil {
			continue
		}

		if len(tfMapObjects) > i {
			tfMapObject = tfMapObjects[i].(map[string]any)
		}

		j := slices.IndexFunc(wasAPIObjectUsed, func(isUsed bool) bool { return !isUsed })
		advancedReplicationSpec, err := flattenAdvancedReplicationSpec(ctx, apiObjects[j], tfMapObject, d, conn)

		if err != nil {
			return nil, err
		}

		tfList[i] = advancedReplicationSpec
		wasAPIObjectUsed[j] = true
	}

	return tfList, nil
}

func flattenAdvancedReplicationSpecRegionConfig(apiObject *matlas.AdvancedRegionConfig, tfMapObject map[string]any) map[string]any {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]any{}
	if tfMapObject != nil {
		if v, ok := tfMapObject["analytics_specs"]; ok && len(v.([]any)) > 0 {
			tfMap["analytics_specs"] = flattenAdvancedReplicationSpecRegionConfigSpec(apiObject.AnalyticsSpecs, apiObject.ProviderName, tfMapObject["analytics_specs"].([]any))
		}
		if v, ok := tfMapObject["electable_specs"]; ok && len(v.([]any)) > 0 {
			tfMap["electable_specs"] = flattenAdvancedReplicationSpecRegionConfigSpec(apiObject.ElectableSpecs, apiObject.ProviderName, tfMapObject["electable_specs"].([]any))
		}
		if v, ok := tfMapObject["read_only_specs"]; ok && len(v.([]any)) > 0 {
			tfMap["read_only_specs"] = flattenAdvancedReplicationSpecRegionConfigSpec(apiObject.ReadOnlySpecs, apiObject.ProviderName, tfMapObject["read_only_specs"].([]any))
		}
		if v, ok := tfMapObject["auto_scaling"]; ok && len(v.([]any)) > 0 {
			tfMap["auto_scaling"] = flattenAdvancedReplicationSpecAutoScaling(apiObject.AutoScaling)
		}
		if v, ok := tfMapObject["analytics_auto_scaling"]; ok && len(v.([]any)) > 0 {
			tfMap["analytics_auto_scaling"] = flattenAdvancedReplicationSpecAutoScaling(apiObject.AnalyticsAutoScaling)
		}
	} else {
		tfMap["analytics_specs"] = flattenAdvancedReplicationSpecRegionConfigSpec(apiObject.AnalyticsSpecs, apiObject.ProviderName, nil)
		tfMap["electable_specs"] = flattenAdvancedReplicationSpecRegionConfigSpec(apiObject.ElectableSpecs, apiObject.ProviderName, nil)
		tfMap["read_only_specs"] = flattenAdvancedReplicationSpecRegionConfigSpec(apiObject.ReadOnlySpecs, apiObject.ProviderName, nil)
		tfMap["auto_scaling"] = flattenAdvancedReplicationSpecAutoScaling(apiObject.AutoScaling)
		tfMap["analytics_auto_scaling"] = flattenAdvancedReplicationSpecAutoScaling(apiObject.AnalyticsAutoScaling)
	}

	tfMap["region_name"] = apiObject.RegionName
	tfMap["provider_name"] = apiObject.ProviderName
	tfMap["backing_provider_name"] = apiObject.BackingProviderName
	tfMap["priority"] = apiObject.Priority

	return tfMap
}

func flattenAdvancedReplicationSpecRegionConfigs(ctx context.Context, apiObjects []*matlas.AdvancedRegionConfig, tfMapObjects []any,
	d *schema.ResourceData, conn *matlas.Client) (tfResult []map[string]any, containersIDs map[string]string, err error) {
	if len(apiObjects) == 0 {
		return nil, nil, nil
	}

	var tfList []map[string]any
	containerIds := make(map[string]string)

	for i, apiObject := range apiObjects {
		if apiObject == nil {
			continue
		}

		if len(tfMapObjects) > i {
			tfMapObject := tfMapObjects[i].(map[string]any)
			tfList = append(tfList, flattenAdvancedReplicationSpecRegionConfig(apiObject, tfMapObject))
		} else {
			tfList = append(tfList, flattenAdvancedReplicationSpecRegionConfig(apiObject, nil))
		}

		if apiObject.ProviderName != "TENANT" {
			containers, _, err := conn.Containers.List(ctx, d.Get("project_id").(string),
				&matlas.ContainersListOptions{ProviderName: apiObject.ProviderName})
			if err != nil {
				return nil, nil, err
			}
			if result := getAdvancedClusterContainerID(containers, apiObject); result != "" {
				// Will print as "providerName:regionName" = "containerId" in terraform show
				containerIds[fmt.Sprintf("%s:%s", apiObject.ProviderName, apiObject.RegionName)] = result
			}
		}
	}

	return tfList, containerIds, nil
}

func flattenAdvancedReplicationSpecRegionConfigSpec(apiObject *matlas.Specs, providerName string, tfMapObjects []any) []map[string]any {
	if apiObject == nil {
		return nil
	}
	var tfList []map[string]any

	tfMap := map[string]any{}

	if len(tfMapObjects) > 0 {
		tfMapObject := tfMapObjects[0].(map[string]any)

		if providerName == "AWS" {
			if cast.ToInt64(apiObject.DiskIOPS) > 0 {
				if v, ok := tfMapObject["disk_iops"]; ok && v.(int) > 0 {
					tfMap["disk_iops"] = apiObject.DiskIOPS
				}
			}
			if v, ok := tfMapObject["ebs_volume_type"]; ok && v.(string) != "" {
				tfMap["ebs_volume_type"] = apiObject.EbsVolumeType
			}
		}
		if _, ok := tfMapObject["node_count"]; ok {
			tfMap["node_count"] = apiObject.NodeCount
		}
		if v, ok := tfMapObject["instance_size"]; ok && v.(string) != "" {
			tfMap["instance_size"] = apiObject.InstanceSize
			tfList = append(tfList, tfMap)
		}
	} else {
		tfMap["disk_iops"] = apiObject.DiskIOPS
		tfMap["ebs_volume_type"] = apiObject.EbsVolumeType
		tfMap["node_count"] = apiObject.NodeCount
		tfMap["instance_size"] = apiObject.InstanceSize
		tfList = append(tfList, tfMap)
	}

	return tfList
}

func flattenAdvancedReplicationSpecAutoScaling(apiObject *matlas.AdvancedAutoScaling) []map[string]any {
	if apiObject == nil {
		return nil
	}

	var tfList []map[string]any

	tfMap := map[string]any{}
	if apiObject.DiskGB != nil {
		tfMap["disk_gb_enabled"] = apiObject.DiskGB.Enabled
	}
	if apiObject.Compute != nil {
		tfMap["compute_enabled"] = apiObject.Compute.Enabled
		tfMap["compute_scale_down_enabled"] = apiObject.Compute.ScaleDownEnabled
		tfMap["compute_min_instance_size"] = apiObject.Compute.MinInstanceSize
		tfMap["compute_max_instance_size"] = apiObject.Compute.MaxInstanceSize
	}

	tfList = append(tfList, tfMap)

	return tfList
}

func resourceClusterAdvancedRefreshFunc(ctx context.Context, name, projectID string, client *matlas.Client) retry.StateRefreshFunc {
	return func() (any, string, error) {
		c, resp, err := client.AdvancedClusters.Get(ctx, projectID, name)

		if err != nil && strings.Contains(err.Error(), "reset by peer") {
			return nil, "REPEATING", nil
		}

		if err != nil && c == nil && resp == nil {
			return nil, "", err
		}

		if err != nil {
			if resp.StatusCode == 404 {
				return "", "DELETED", nil
			}
			if resp.StatusCode == 503 {
				return "", "PENDING", nil
			}
			return nil, "", err
		}

		if c.StateName != "" {
			log.Printf("[DEBUG] status for MongoDB cluster: %s: %s", name, c.StateName)
		}

		return c, c.StateName, nil
	}
}

func resourceClusterListAdvancedRefreshFunc(ctx context.Context, projectID string, client *matlas.Client) retry.StateRefreshFunc {
	return func() (any, string, error) {
		clusters, resp, err := client.AdvancedClusters.List(ctx, projectID, nil)

		if err != nil && strings.Contains(err.Error(), "reset by peer") {
			return nil, "REPEATING", nil
		}

		if err != nil && clusters == nil && resp == nil {
			return nil, "", err
		}

		if err != nil {
			if resp.StatusCode == 404 {
				return "", "DELETED", nil
			}
			if resp.StatusCode == 503 {
				return "", "PENDING", nil
			}
			return nil, "", err
		}

		for i := range clusters.Results {
			if clusters.Results[i].StateName != "IDLE" {
				return clusters.Results[i], "PENDING", nil
			}
		}

		return clusters, "IDLE", nil
	}
}

func replicationSpecsHashSet(v any) int {
	var buf bytes.Buffer
	m := v.(map[string]any)
	buf.WriteString(fmt.Sprintf("%d", m["num_shards"].(int)))
	buf.WriteString(fmt.Sprintf("%+v", m["region_configs"].(*schema.Set)))
	buf.WriteString(m["zone_name"].(string))
	return schema.HashString(buf.String())
}

func getUpgradeRequest(d *schema.ResourceData) *matlas.Cluster {
	if !d.HasChange("replication_specs") {
		return nil
	}

	cs, us := d.GetChange("replication_specs")
	currentSpecs := expandAdvancedReplicationSpecs(cs.([]any))
	updatedSpecs := expandAdvancedReplicationSpecs(us.([]any))

	if len(currentSpecs) != 1 || len(updatedSpecs) != 1 || len(currentSpecs[0].RegionConfigs) != 1 || len(updatedSpecs[0].RegionConfigs) != 1 {
		return nil
	}

	currentRegion := currentSpecs[0].RegionConfigs[0]
	updatedRegion := updatedSpecs[0].RegionConfigs[0]
	currentSize := currentRegion.ElectableSpecs.InstanceSize

	if currentRegion.ElectableSpecs.InstanceSize == updatedRegion.ElectableSpecs.InstanceSize || !isSharedTier(currentSize) {
		return nil
	}

	return &matlas.Cluster{
		ProviderSettings: &matlas.ProviderSettings{
			ProviderName:     updatedRegion.ProviderName,
			InstanceSizeName: updatedRegion.ElectableSpecs.InstanceSize,
			RegionName:       updatedRegion.RegionName,
		},
	}
}

func updateAdvancedCluster(
	ctx context.Context,
	conn *matlas.Client,
	request *matlas.AdvancedCluster,
	projectID, name string,
	timeout time.Duration,
) (*matlas.AdvancedCluster, *matlas.Response, error) {
	cluster, resp, err := conn.AdvancedClusters.Update(ctx, projectID, name, request)
	if err != nil {
		return nil, nil, err
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"CREATING", "UPDATING", "REPAIRING"},
		Target:     []string{"IDLE"},
		Refresh:    resourceClusterAdvancedRefreshFunc(ctx, name, projectID, conn),
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

func getAdvancedClusterContainerID(containers []matlas.Container, cluster *matlas.AdvancedRegionConfig) string {
	if len(containers) != 0 {
		for i := range containers {
			if cluster.ProviderName == "GCP" {
				return containers[i].ID
			}

			if containers[i].ProviderName == cluster.ProviderName &&
				containers[i].Region == cluster.RegionName || // For Azure
				containers[i].RegionName == cluster.RegionName { // For AWS
				return containers[i].ID
			}
		}
	}

	return ""
}

func flattenLabels(l []matlas.Label) []map[string]any {
	labels := make([]map[string]any, len(l))
	for i, v := range l {
		labels[i] = map[string]any{
			"key":   v.Key,
			"value": v.Value,
		}
	}

	return labels
}

func flattenTags(l *[]*matlas.Tag) []map[string]any {
	if l == nil {
		return []map[string]any{}
	}
	tags := make([]map[string]any, len(*l))
	for i, v := range *l {
		tags[i] = map[string]any{
			"key":   v.Key,
			"value": v.Value,
		}
	}
	return tags
}

func expandLabelSliceFromSetSchema(d *schema.ResourceData) []matlas.Label {
	list := d.Get("labels").(*schema.Set)
	res := make([]matlas.Label, list.Len())

	for i, val := range list.List() {
		v := val.(map[string]any)
		res[i] = matlas.Label{
			Key:   v["key"].(string),
			Value: v["value"].(string),
		}
	}

	return res
}

func expandTagSliceFromSetSchema(d *schema.ResourceData) []*matlas.Tag {
	list := d.Get("tags").(*schema.Set)
	res := make([]*matlas.Tag, list.Len())
	for i, val := range list.List() {
		v := val.(map[string]any)
		res[i] = &matlas.Tag{
			Key:   v["key"].(string),
			Value: v["value"].(string),
		}
	}
	return res
}

func containsLabelOrKey(list []matlas.Label, item matlas.Label) bool {
	for _, v := range list {
		if reflect.DeepEqual(v, item) || v.Key == item.Key {
			return true
		}
	}

	return false
}

func removeLabel(list []matlas.Label, item matlas.Label) []matlas.Label {
	var pos int

	for _, v := range list {
		if reflect.DeepEqual(v, item) {
			list = append(list[:pos], list[pos+1:]...)

			if pos > 0 {
				pos--
			}

			continue
		}
		pos++
	}

	return list
}
