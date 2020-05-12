package mongodbatlas

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/mwielbut/pointy"
	"github.com/spf13/cast"

	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

const (
	errorClusterCreate      = "error creating MongoDB Cluster: %s"
	errorClusterRead        = "error reading MongoDB Cluster (%s): %s"
	errorClusterDelete      = "error deleting MongoDB Cluster (%s): %s"
	errorClusterUpdate      = "error updating MongoDB Cluster (%s): %s"
	errorClusterSetting     = "error setting `%s` for MongoDB Cluster (%s): %s"
	errorAdvancedConfUpdate = "error updating Advanced Configuration Option form MongoDB Cluster (%s): %s"
	errorAdvancedConfRead   = "error reading Advanced Configuration Option form MongoDB Cluster (%s): %s"
)

var defaultLabel = matlas.Label{Key: "Infrastructure Tool", Value: "MongoDB Atlas Terraform Provider"}

func resourceMongoDBAtlasCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoDBAtlasClusterCreate,
		Read:   resourceMongoDBAtlasClusterRead,
		Update: resourceMongoDBAtlasClusterUpdate,
		Delete: resourceMongoDBAtlasClusterDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMongoDBAtlasClusterImportState,
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
				Optional: true,
				Computed: true,
			},
			"backup_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"bi_connector": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeString,
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
			"connection_strings": {
				Type:     schema.TypeList,
				MinItems: 1,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"standard": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"standard_srv": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"aws_private_link": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"aws_private_link_srv": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"private": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_srv": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mongo_db_major_version": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				StateFunc: formatMongoDBMajorVersion,
			},
			"num_shards": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"provider_backup_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"provider_instance_size_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"provider_name": {
				Type:     schema.TypeString,
				ForceNew: true,
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
				Type:     schema.TypeBool,
				Optional: true,
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
			"replication_factor": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"replication_specs": {
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
						"num_shards": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"regions_config": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region_name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
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
			"advanced_configuration": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"fail_index_key_too_long": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"javascript_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"minimum_enabled_tls_protocol": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"no_table_scan": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"oplog_size_mb": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"sample_size_bi_connector": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"sample_refresh_interval_bi_connector": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"labels": {
				Type:     schema.TypeSet,
				Optional: true,
				Set: func(v interface{}) int {
					var buf bytes.Buffer
					m := v.(map[string]interface{})
					buf.WriteString(m["key"].(string))
					buf.WriteString(m["value"].(string))
					return hashcode.String(buf.String())
				},
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"snapshot_backup_policy": computedCloudProviderSnapshotBackupPolicySchema(),
			"container_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceMongoDBAtlasClusterCreate(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)
	providerName := d.Get("provider_name").(string)

	autoScaling := matlas.AutoScaling{
		DiskGBEnabled: pointy.Bool(true),
	}

	if diskGBEnabled, ok := d.GetOkExists("auto_scaling_disk_gb_enabled"); ok {
		autoScaling = matlas.AutoScaling{
			DiskGBEnabled: pointy.Bool(diskGBEnabled.(bool)),
		}
	}

	//validate cluster_type conditional
	if _, ok := d.GetOk("replication_specs"); ok {
		if _, ok1 := d.GetOk("cluster_type"); !ok1 {
			return fmt.Errorf("`cluster_type` should be set when `replication_specs` is set")
		}

		if _, ok1 := d.GetOk("num_shards"); !ok1 {
			return fmt.Errorf("`num_shards` should be set when `replication_specs` is set")
		}
	}

	if providerName != "AWS" {
		if _, ok := d.GetOk("provider_disk_iops"); ok {
			return fmt.Errorf("`provider_disk_iops` shouldn't be set when provider name is `GCP` or `AZURE`")
		}
		if _, ok := d.GetOk("provider_encrypt_ebs_volume"); ok {
			return fmt.Errorf("`provider_encrypt_ebs_volume` shouldn't be set when provider name is `GCP` or `AZURE`")
		}
		if _, ok := d.GetOk("provider_volume_type"); ok {
			return fmt.Errorf("`provider_volume_type` shouldn't be set when provider name is `GCP` or `AZURE`")
		}
	}

	if providerName != "AZURE" {
		if _, ok := d.GetOk("provider_disk_type_name"); ok {
			return fmt.Errorf("`provider_disk_type_name` shouldn't be set when provider name is `GCP` or `AWS`")
		}
	}

	if providerName == "AZURE" {
		if _, ok := d.GetOk("disk_size_gb"); ok {
			return fmt.Errorf("`disk_size_gb` cannot be used with Azure clusters")
		}
	}

	if providerName == "TENANT" {
		if diskGBEnabled := d.Get("auto_scaling_disk_gb_enabled"); diskGBEnabled.(bool) {
			return fmt.Errorf("`auto_scaling_disk_gb_enabled` cannot be true when provider name is TENANT")
		}
		autoScaling = matlas.AutoScaling{
			DiskGBEnabled: pointy.Bool(false),
		}
	}

	// We need to validate the oplog_size_mb attr of the advanced configuration option to show the error
	// before that the cluster is created
	if oplogSizeMB, ok := d.GetOk("advanced_configuration.oplog_size_mb"); ok {
		if cast.ToInt64(oplogSizeMB) <= 0 {
			return fmt.Errorf("`advanced_configuration.oplog_size_mb` cannot be <= 0")
		}
	}

	biConnector, err := expandBiConnector(d)
	if err != nil {
		return fmt.Errorf(errorClusterCreate, err)
	}

	providerSettings := expandProviderSetting(d)

	replicationSpecs, err := expandReplicationSpecs(d)
	if err != nil {
		return fmt.Errorf(errorClusterCreate, err)
	}

	clusterRequest := &matlas.Cluster{
		Name:                     d.Get("name").(string),
		EncryptionAtRestProvider: d.Get("encryption_at_rest_provider").(string),
		ClusterType:              cast.ToString(d.Get("cluster_type")),
		BackupEnabled:            pointy.Bool(d.Get("backup_enabled").(bool)),
		ProviderBackupEnabled:    pointy.Bool(d.Get("provider_backup_enabled").(bool)),
		PitEnabled:               pointy.Bool(d.Get("pit_enabled").(bool)),
		AutoScaling:              autoScaling,
		BiConnector:              biConnector,
		ProviderSettings:         &providerSettings,
		ReplicationSpecs:         replicationSpecs,
	}

	if containsLabelOrKey(expandLabelSliceFromSetSchema(d), defaultLabel) {
		return fmt.Errorf("you should not set `Infrastructure Tool` label, it is used for internal purposes.")
	}

	clusterRequest.Labels = append(expandLabelSliceFromSetSchema(d), defaultLabel)

	if v, ok := d.GetOk("disk_size_gb"); ok {
		clusterRequest.DiskSizeGB = pointy.Float64(v.(float64))
	}

	if v, ok := d.GetOk("mongo_db_major_version"); ok {
		clusterRequest.MongoDBMajorVersion = formatMongoDBMajorVersion(v.(string))
	}

	if r, ok := d.GetOk("replication_factor"); ok {
		clusterRequest.ReplicationFactor = pointy.Int64(cast.ToInt64(r))
	}

	if n, ok := d.GetOk("num_shards"); ok {
		clusterRequest.NumShards = pointy.Int64(cast.ToInt64(n))
	}

	cluster, _, err := conn.Clusters.Create(context.Background(), projectID, clusterRequest)
	if err != nil {
		return fmt.Errorf(errorClusterCreate, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATING", "UPDATING", "REPAIRING", "REPEATING"},
		Target:     []string{"IDLE"},
		Refresh:    resourceClusterRefreshFunc(d.Get("name").(string), projectID, conn),
		Timeout:    3 * time.Hour,
		MinTimeout: 30 * time.Second,
		Delay:      1 * time.Minute,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(errorClusterCreate, err)
	}

	/*
		So far, the cluster has created correctly, so we need to set up
		the advanced configuration option to attach it
	*/
	ac, ok := d.GetOk("advanced_configuration")
	advancedConfReq := expandProcessArgs(ac.(map[string]interface{}))
	if ok {
		_, _, err := conn.Clusters.UpdateProcessArgs(context.Background(), projectID, cluster.Name, advancedConfReq)
		if err != nil {
			return fmt.Errorf(errorAdvancedConfUpdate, cluster.Name, err)
		}
	}

	d.SetId(encodeStateID(map[string]string{
		"cluster_id":    cluster.ID,
		"project_id":    projectID,
		"cluster_name":  cluster.Name,
		"provider_name": providerName,
	}))

	return resourceMongoDBAtlasClusterRead(d, meta)
}

func resourceMongoDBAtlasClusterRead(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]
	providerName := ids["provider_name"]

	cluster, resp, err := conn.Clusters.Get(context.Background(), projectID, clusterName)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf(errorClusterRead, clusterName, err)
	}

	log.Printf("[DEBUG] GET Cluster %+v", cluster)

	if err := d.Set("cluster_id", cluster.ID); err != nil {
		return fmt.Errorf(errorClusterSetting, "cluster_id", clusterName, err)
	}
	if err := d.Set("auto_scaling_disk_gb_enabled", cluster.AutoScaling.DiskGBEnabled); err != nil {
		return fmt.Errorf(errorClusterSetting, "auto_scaling_disk_gb_enabled", clusterName, err)
	}
	if err := d.Set("backup_enabled", cluster.BackupEnabled); err != nil {
		return fmt.Errorf(errorClusterSetting, "backup_enabled", clusterName, err)
	}
	if err := d.Set("provider_backup_enabled", cluster.ProviderBackupEnabled); err != nil {
		return fmt.Errorf(errorClusterSetting, "provider_backup_enabled", clusterName, err)
	}
	if err := d.Set("cluster_type", cluster.ClusterType); err != nil {
		return fmt.Errorf(errorClusterSetting, "cluster_type", clusterName, err)
	}
	if err := d.Set("connection_strings", flattenConnectionStrings(cluster.ConnectionStrings)); err != nil {
		return fmt.Errorf(errorClusterSetting, "connection_strings", clusterName, err)
	}
	if err := d.Set("disk_size_gb", cluster.DiskSizeGB); err != nil {
		return fmt.Errorf(errorClusterSetting, "disk_size_gb", clusterName, err)
	}
	if err := d.Set("encryption_at_rest_provider", cluster.EncryptionAtRestProvider); err != nil {
		return fmt.Errorf(errorClusterSetting, "encryption_at_rest_provider", clusterName, err)
	}
	if err := d.Set("mongo_db_major_version", cluster.MongoDBMajorVersion); err != nil {
		return fmt.Errorf(errorClusterSetting, "mongo_db_major_version", clusterName, err)
	}

	//Avoid Global Cluster issues. (NumShards is not present in Global Clusters)
	if cluster.NumShards != nil {
		if err := d.Set("num_shards", cluster.NumShards); err != nil {
			return fmt.Errorf(errorClusterSetting, "num_shards", clusterName, err)
		}
	}

	if err := d.Set("mongo_db_version", cluster.MongoDBVersion); err != nil {
		return fmt.Errorf(errorClusterSetting, "mongo_db_version", clusterName, err)
	}
	if err := d.Set("mongo_uri", cluster.MongoURI); err != nil {
		return fmt.Errorf(errorClusterSetting, "mongo_uri", clusterName, err)
	}
	if err := d.Set("mongo_uri_updated", cluster.MongoURIUpdated); err != nil {
		return fmt.Errorf(errorClusterSetting, "mongo_uri_updated", clusterName, err)
	}
	if err := d.Set("mongo_uri_with_options", cluster.MongoURIWithOptions); err != nil {
		return fmt.Errorf(errorClusterSetting, "mongo_uri_with_options", clusterName, err)
	}

	if err := d.Set("pit_enabled", cluster.PitEnabled); err != nil {
		return fmt.Errorf(errorClusterSetting, "pit_enabled", clusterName, err)
	}

	if err := d.Set("paused", cluster.Paused); err != nil {
		return fmt.Errorf(errorClusterSetting, "paused", clusterName, err)
	}
	if err := d.Set("srv_address", cluster.SrvAddress); err != nil {
		return fmt.Errorf(errorClusterSetting, "srv_address", clusterName, err)
	}
	if err := d.Set("state_name", cluster.StateName); err != nil {
		return fmt.Errorf(errorClusterSetting, "state_name", clusterName, err)
	}
	if err := d.Set("bi_connector", flattenBiConnector(cluster.BiConnector)); err != nil {
		return fmt.Errorf(errorClusterSetting, "bi_connector", clusterName, err)
	}
	if cluster.ProviderSettings != nil {
		flattenProviderSettings(d, cluster.ProviderSettings, clusterName)
	}
	if err := d.Set("replication_specs", flattenReplicationSpecs(cluster.ReplicationSpecs)); err != nil {
		return fmt.Errorf(errorClusterSetting, "replication_specs", clusterName, err)
	}
	if err := d.Set("replication_factor", cluster.ReplicationFactor); err != nil {
		return fmt.Errorf(errorClusterSetting, "replication_factor", clusterName, err)
	}

	if err := d.Set("labels", flattenLabels(removeLabel(cluster.Labels, defaultLabel))); err != nil {
		return fmt.Errorf(errorClusterSetting, "labels", clusterName, err)
	}

	containers, _, err := conn.Containers.List(context.Background(), projectID,
		&matlas.ContainersListOptions{ProviderName: providerName})
	if err != nil {
		return fmt.Errorf(errorClusterRead, clusterName, err)
	}

	if err := d.Set("container_id", getContainerID(containers, cluster)); err != nil {
		return fmt.Errorf(errorClusterSetting, "container_id", clusterName, err)
	}

	/*
		Get the advaced configuration options and set up to the terraform state
	*/
	processArgs, _, err := conn.Clusters.GetProcessArgs(context.Background(), projectID, clusterName)
	if err != nil {
		return fmt.Errorf(errorAdvancedConfRead, clusterName, err)
	}

	if err := d.Set("advanced_configuration", flattenProcessArgs(processArgs)); err != nil {
		return fmt.Errorf(errorClusterSetting, "advanced_configuration", clusterName, err)
	}

	// Get the snapshot policy and set the data
	snapshotBackupPolicy, err := flattenCloudProviderSnapshotBackupPolicy(d, conn, projectID, clusterName)
	if err != nil {
		return err
	}
	if err := d.Set("snapshot_backup_policy", snapshotBackupPolicy); err != nil {
		return err
	}

	return nil
}

func resourceMongoDBAtlasClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	cluster := new(matlas.Cluster)

	if d.HasChange("bi_connector") {
		cluster.BiConnector, _ = expandBiConnector(d)
	}

	providerSettings := matlas.ProviderSettings{}

	// If at least one of the provider settings argument has changed, expand all provider settings
	if d.HasChange("provider_disk_iops") || d.HasChange("provider_encrypt_ebs_volume") ||
		d.HasChange("backing_provider_name") || d.HasChange("provider_disk_type_name") ||
		d.HasChange("provider_instance_size_name") || d.HasChange("provider_instance_size_name") ||
		d.HasChange("provider_instance_size_name") || d.HasChange("provider_name") ||
		d.HasChange("provider_region_name") || d.HasChange("provider_volume_type") {
		providerSettings = expandProviderSetting(d)
	}

	//Check if Provider setting was changed.
	if !reflect.DeepEqual(providerSettings, matlas.ProviderSettings{}) {
		cluster.ProviderSettings = &providerSettings
	}

	if d.HasChange("replication_specs") {
		replicationSpecs, err := expandReplicationSpecs(d)
		if err != nil {
			return fmt.Errorf(errorClusterUpdate, clusterName, err)
		}
		cluster.ReplicationSpecs = replicationSpecs
	}

	if d.HasChange("auto_scaling_disk_gb_enabled") {
		cluster.AutoScaling.DiskGBEnabled = pointy.Bool(d.Get("auto_scaling_disk_gb_enabled").(bool))
	}
	if d.HasChange("encryption_at_rest_provider") {
		cluster.EncryptionAtRestProvider = d.Get("encryption_at_rest_provider").(string)
	}
	if d.HasChange("mongo_db_major_version") {
		cluster.MongoDBMajorVersion = formatMongoDBMajorVersion(d.Get("mongo_db_major_version"))
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
	if d.HasChange("provider_backup_enabled") {
		cluster.ProviderBackupEnabled = pointy.Bool(d.Get("provider_backup_enabled").(bool))
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
	if d.HasChange("labels") {
		if containsLabelOrKey(expandLabelSliceFromSetSchema(d), defaultLabel) {
			return fmt.Errorf("you should not set `Infrastructure Tool` label, it is used for internal purposes.")
		}

		cluster.Labels = append(expandLabelSliceFromSetSchema(d), defaultLabel)
	}

	// Has changes
	if !reflect.DeepEqual(cluster, matlas.Cluster{}) {
		_, _, err := conn.Clusters.Update(context.Background(), projectID, clusterName, cluster)
		if err != nil {
			return fmt.Errorf(errorClusterUpdate, clusterName, err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATING", "UPDATING", "REPAIRING"},
		Target:     []string{"IDLE"},
		Refresh:    resourceClusterRefreshFunc(clusterName, projectID, conn),
		Timeout:    3 * time.Hour,
		MinTimeout: 30 * time.Second,
		Delay:      1 * time.Minute,
	}

	// Wait, catching any errors
	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(errorClusterCreate, err)
	}

	/*
		Check if advaced configuration option has a changes to update it
	*/
	if d.HasChange("advanced_configuration") {
		advancedConfReq := expandProcessArgs(d.Get("advanced_configuration").(map[string]interface{}))

		if !reflect.DeepEqual(advancedConfReq, matlas.ProcessArgs{}) {
			_, _, err := conn.Clusters.UpdateProcessArgs(context.Background(), projectID, clusterName, advancedConfReq)
			if err != nil {
				return fmt.Errorf(errorAdvancedConfUpdate, clusterName, err)
			}
		}
	}

	return resourceMongoDBAtlasClusterRead(d, meta)
}

func resourceMongoDBAtlasClusterDelete(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	_, err := conn.Clusters.Delete(context.Background(), projectID, clusterName)

	if err != nil {
		return fmt.Errorf(errorClusterDelete, clusterName, err)
	}

	log.Println("[INFO] Waiting for MongoDB Cluster to be destroyed")

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"IDLE", "CREATING", "UPDATING", "REPAIRING", "DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    resourceClusterRefreshFunc(clusterName, projectID, conn),
		Timeout:    3 * time.Hour,
		MinTimeout: 30 * time.Second,
		Delay:      1 * time.Minute, // Wait 30 secs before starting
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(errorClusterDelete, clusterName, err)
	}
	return nil
}

func resourceMongoDBAtlasClusterImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*matlas.Client)

	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 2 {
		return nil, errors.New("import format error: to import a cluster, use the format {project_id}-{name}")
	}

	projectID := parts[0]
	name := parts[1]

	u, _, err := conn.Clusters.Get(context.Background(), projectID, name)
	if err != nil {
		return nil, fmt.Errorf("couldn't import cluster %s in project %s, error: %s", name, projectID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"cluster_id":   u.ID,
		"project_id":   projectID,
		"cluster_name": u.Name,
	}))

	if err := d.Set("project_id", u.GroupID); err != nil {
		log.Printf(errorClusterSetting, "project_id", u.ID, err)
	}
	if err := d.Set("name", u.Name); err != nil {
		log.Printf(errorClusterSetting, "name", u.ID, err)
	}

	return []*schema.ResourceData{d}, nil
}

func expandBiConnector(d *schema.ResourceData) (matlas.BiConnector, error) {
	var biConnector matlas.BiConnector

	if v, ok := d.GetOk("bi_connector"); ok {
		biConnMap := v.(map[string]interface{})

		enabled := cast.ToBool(biConnMap["enabled"])

		biConnector = matlas.BiConnector{
			Enabled:        &enabled,
			ReadPreference: cast.ToString(biConnMap["read_preference"]),
		}
	}
	return biConnector, nil
}

func flattenBiConnector(biConnector matlas.BiConnector) map[string]interface{} {
	biConnectorMap := make(map[string]interface{})

	if biConnector.Enabled != nil {
		biConnectorMap["enabled"] = strconv.FormatBool(*biConnector.Enabled)
	}

	if biConnector.ReadPreference != "" {
		biConnectorMap["read_preference"] = biConnector.ReadPreference
	}

	return biConnectorMap
}

func expandProviderSetting(d *schema.ResourceData) matlas.ProviderSettings {
	providerSettings := matlas.ProviderSettings{}

	if d.Get("provider_name") == "AWS" {

		// Check if the Provider Disk IOS sets in the Terraform configuration.
		// If it didn't, the MongoDB Atlas server would set it to the default for the amount of storage.
		if v, ok := d.GetOk("provider_disk_iops"); ok {
			providerSettings.DiskIOPS = pointy.Int64(cast.ToInt64(v))
		}

		providerSettings.EncryptEBSVolume = pointy.Bool(true)
		if encryptEBSVolume, ok := d.GetOkExists("provider_encrypt_ebs_volume"); ok {
			providerSettings.EncryptEBSVolume = pointy.Bool(cast.ToBool(encryptEBSVolume))
		}
	}

	region, _ := valRegion(d.Get("provider_region_name"))

	providerSettings.BackingProviderName = cast.ToString(d.Get("backing_provider_name"))
	providerSettings.InstanceSizeName = cast.ToString(d.Get("provider_instance_size_name"))
	providerSettings.ProviderName = cast.ToString(d.Get("provider_name"))
	providerSettings.RegionName = region
	providerSettings.VolumeType = cast.ToString(d.Get("provider_volume_type"))
	providerSettings.DiskTypeName = cast.ToString(d.Get("provider_disk_type_name"))

	return providerSettings
}

func flattenProviderSettings(d *schema.ResourceData, settings *matlas.ProviderSettings, clusterName string) {

	if err := d.Set("backing_provider_name", settings.BackingProviderName); err != nil {
		log.Printf(errorClusterSetting, "backing_provider_name", clusterName, err)
	}

	if settings.DiskIOPS != nil && *settings.DiskIOPS != 0 {
		if err := d.Set("provider_disk_iops", *settings.DiskIOPS); err != nil {
			log.Printf(errorClusterSetting, "provider_disk_iops", clusterName, err)
		}
	}

	if err := d.Set("provider_disk_type_name", settings.DiskTypeName); err != nil {
		log.Printf(errorClusterSetting, "provider_disk_type_name", clusterName, err)
	}

	if err := d.Set("provider_encrypt_ebs_volume", settings.EncryptEBSVolume); err != nil {
		log.Printf(errorClusterSetting, "provider_encrypt_ebs_volume", clusterName, err)
	}

	if err := d.Set("provider_instance_size_name", settings.InstanceSizeName); err != nil {
		log.Printf(errorClusterSetting, "provider_instance_size_name", clusterName, err)
	}

	if err := d.Set("provider_name", settings.ProviderName); err != nil {
		log.Printf(errorClusterSetting, "provider_name", clusterName, err)
	}

	if err := d.Set("provider_region_name", settings.RegionName); err != nil {
		log.Printf(errorClusterSetting, "provider_region_name", clusterName, err)
	}

	if err := d.Set("provider_volume_type", settings.VolumeType); err != nil {
		log.Printf(errorClusterSetting, "provider_volume_type", clusterName, err)
	}
}

func expandReplicationSpecs(d *schema.ResourceData) ([]matlas.ReplicationSpec, error) {
	rSpecs := make([]matlas.ReplicationSpec, 0)

	if v, ok := d.GetOk("replication_specs"); ok {
		for _, s := range v.([]interface{}) {
			spec := s.(map[string]interface{})

			regionsConfig, err := expandRegionsConfig(spec["regions_config"].(*schema.Set).List())
			if err != nil {
				return rSpecs, err
			}

			rSpec := matlas.ReplicationSpec{
				ID:            cast.ToString(spec["id"]),
				NumShards:     pointy.Int64(cast.ToInt64(spec["num_shards"])),
				ZoneName:      cast.ToString(spec["zone_name"]),
				RegionsConfig: regionsConfig,
			}
			rSpecs = append(rSpecs, rSpec)
		}
	}

	return rSpecs, nil
}

func flattenReplicationSpecs(rSpecs []matlas.ReplicationSpec) []map[string]interface{} {
	specs := make([]map[string]interface{}, 0)
	for _, rSpec := range rSpecs {
		spec := map[string]interface{}{
			"id":             rSpec.ID,
			"num_shards":     rSpec.NumShards,
			"zone_name":      cast.ToString(rSpec.ZoneName),
			"regions_config": flattenRegionsConfig(rSpec.RegionsConfig),
		}
		specs = append(specs, spec)
	}
	return specs
}

func expandRegionsConfig(regions []interface{}) (map[string]matlas.RegionsConfig, error) {
	regionsConfig := make(map[string]matlas.RegionsConfig)
	for _, r := range regions {
		region := r.(map[string]interface{})

		r, err := valRegion(region["region_name"])
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

func flattenRegionsConfig(regionsConfig map[string]matlas.RegionsConfig) []map[string]interface{} {
	regions := make([]map[string]interface{}, 0)

	for regionName, regionConfig := range regionsConfig {
		region := map[string]interface{}{
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

func expandProcessArgs(p map[string]interface{}) *matlas.ProcessArgs {
	res := &matlas.ProcessArgs{
		FailIndexKeyTooLong:              pointy.Bool(cast.ToBool(p["fail_index_key_too_long"])),
		JavascriptEnabled:                pointy.Bool(cast.ToBool(p["javascript_enabled"])),
		MinimumEnabledTLSProtocol:        cast.ToString(p["minimum_enabled_tls_protocol"]),
		NoTableScan:                      pointy.Bool(cast.ToBool(p["no_table_scan"])),
		SampleSizeBIConnector:            pointy.Int64(cast.ToInt64(p["sample_size_bi_connector"])),
		SampleRefreshIntervalBIConnector: pointy.Int64(cast.ToInt64(p["sample_refresh_interval_bi_connector"])),
	}
	if sizeMB := cast.ToInt64(p["oplog_size_mb"]); sizeMB != 0 {
		res.OplogSizeMB = pointy.Int64(cast.ToInt64(p["oplog_size_mb"]))
	} else {
		log.Printf(errorClusterSetting, `oplog_size_mb`, "", cast.ToString(sizeMB))
	}
	return res
}

func flattenProcessArgs(p *matlas.ProcessArgs) map[string]interface{} {
	return map[string]interface{}{
		"fail_index_key_too_long":              cast.ToString(*p.FailIndexKeyTooLong),
		"javascript_enabled":                   cast.ToString(*p.JavascriptEnabled),
		"minimum_enabled_tls_protocol":         cast.ToString(p.MinimumEnabledTLSProtocol),
		"no_table_scan":                        cast.ToString(*p.NoTableScan),
		"oplog_size_mb":                        cast.ToString(p.OplogSizeMB),
		"sample_size_bi_connector":             cast.ToString(p.SampleSizeBIConnector),
		"sample_refresh_interval_bi_connector": cast.ToString(p.SampleRefreshIntervalBIConnector),
	}
}

func resourceClusterRefreshFunc(name, projectID string, client *matlas.Client) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		c, resp, err := client.Clusters.Get(context.Background(), projectID, name)

		if err != nil && strings.Contains(err.Error(), "reset by peer") {
			return nil, "REPEATING", nil
		}

		if err != nil && c == nil && resp == nil {
			log.Printf(errorClusterRead, name, err)
			return nil, "", err
		} else if err != nil {
			if resp.StatusCode == 404 {
				return 42, "DELETED", nil
			}
			log.Printf(errorClusterRead, name, err)
			return nil, "", err
		}

		if c.StateName != "" {
			log.Printf("[DEBUG] status for MongoDB cluster: %s: %s", name, c.StateName)
		}

		return c, c.StateName, nil
	}
}

func formatMongoDBMajorVersion(val interface{}) string {
	if strings.Contains(val.(string), ".") {
		return val.(string)
	}
	return fmt.Sprintf("%.1f", cast.ToFloat32(val))
}

func flattenConnectionStrings(connectionStrings *matlas.ConnectionStrings) []map[string]interface{} {
	connections := make([]map[string]interface{}, 0)

	connections = append(connections, map[string]interface{}{
		"standard":             connectionStrings.Standard,
		"standard_srv":         connectionStrings.StandardSrv,
		"aws_private_link":     connectionStrings.AwsPrivateLink,
		"aws_private_link_srv": connectionStrings.AwsPrivateLinkSrv,
		"private":              connectionStrings.Private,
		"private_srv":          connectionStrings.PrivateSrv,
	})
	return connections
}

func getContainerID(containers []matlas.Container, cluster *matlas.Cluster) string {
	if len(containers) != 0 {
		for _, container := range containers {
			if cluster.ProviderSettings.ProviderName == "GCP" {
				return container.ID
			}
			if container.ProviderName == cluster.ProviderSettings.ProviderName &&
				container.Region == cluster.ProviderSettings.RegionName || // For Azure
				container.RegionName == cluster.ProviderSettings.RegionName { // For AWS
				return container.ID
			}
		}
	}
	return ""
}
