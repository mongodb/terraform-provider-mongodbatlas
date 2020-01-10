package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"

	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasClusterRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"auto_scaling_disk_gb_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"backup_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"bi_connector": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeString, //Convert to Bool
							Computed: true,
						},
						"read_preference": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"cluster_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_size_gb": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"encryption_at_rest_provider": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mongo_db_major_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"num_shards": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"provider_backup_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"provider_instance_size_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provider_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backing_provider_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provider_disk_iops": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"provider_disk_type_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provider_encrypt_ebs_volume": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"provider_region_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provider_volume_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"replication_factor": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"replication_specs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"num_shards": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"regions_config": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"electable_nodes": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"priority": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"read_only_nodes": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"analytics_nodes": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"zone_name": {
							Type:     schema.TypeString,
							Computed: true,
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
			"labels": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"plugin": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceMongoDBAtlasClusterRead(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)
	name := d.Get("name").(string)

	cluster, resp, err := conn.Clusters.Get(context.Background(), projectID, name)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf(errorRead, name, err)
	}

	if err := d.Set("auto_scaling_disk_gb_enabled", cluster.AutoScaling.DiskGBEnabled); err != nil {
		return fmt.Errorf(errorRead, name, err)
	}

	if err := d.Set("backup_enabled", cluster.BackupEnabled); err != nil {
		return fmt.Errorf(errorRead, name, err)
	}

	if err := d.Set("provider_backup_enabled", cluster.ProviderBackupEnabled); err != nil {
		return fmt.Errorf(errorRead, name, err)
	}

	if err := d.Set("cluster_type", cluster.ClusterType); err != nil {
		return fmt.Errorf(errorRead, name, err)
	}

	if err := d.Set("disk_size_gb", cluster.DiskSizeGB); err != nil {
		return fmt.Errorf(errorRead, name, err)
	}

	if err := d.Set("encryption_at_rest_provider", cluster.EncryptionAtRestProvider); err != nil {
		return fmt.Errorf(errorRead, name, err)
	}

	if err := d.Set("mongo_db_major_version", cluster.MongoDBMajorVersion); err != nil {
		return fmt.Errorf(errorRead, name, err)
	}

	//Avoid Global Cluster issues. (NumShards is not present in Global Clusters)
	if cluster.NumShards != nil {
		if err := d.Set("num_shards", cluster.NumShards); err != nil {
			return fmt.Errorf(errorRead, name, err)
		}
	}

	if err := d.Set("mongo_db_version", cluster.MongoDBVersion); err != nil {
		return fmt.Errorf(errorRead, name, err)
	}

	if err := d.Set("mongo_uri", cluster.MongoURI); err != nil {
		return fmt.Errorf(errorRead, name, err)
	}

	if err := d.Set("mongo_uri_updated", cluster.MongoURIUpdated); err != nil {
		return fmt.Errorf(errorRead, name, err)
	}

	if err := d.Set("mongo_uri_with_options", cluster.MongoURIWithOptions); err != nil {
		return fmt.Errorf(errorRead, name, err)
	}

	if err := d.Set("paused", cluster.Paused); err != nil {
		return fmt.Errorf(errorRead, name, err)
	}

	if err := d.Set("srv_address", cluster.SrvAddress); err != nil {
		return fmt.Errorf(errorRead, name, err)
	}

	if err := d.Set("state_name", cluster.StateName); err != nil {
		return fmt.Errorf(errorRead, name, err)
	}

	if err := d.Set("bi_connector", flattenBiConnector(cluster.BiConnector)); err != nil {
		return fmt.Errorf(errorRead, name, err)
	}

	if cluster.ProviderSettings != nil {
		flattenProviderSettings(d, *cluster.ProviderSettings)
	}

	if err := d.Set("replication_specs", flattenReplicationSpecs(cluster.ReplicationSpecs)); err != nil {
		return fmt.Errorf(errorRead, name, err)
	}

	if err := d.Set("replication_factor", cluster.ReplicationFactor); err != nil {
		return fmt.Errorf(errorRead, name, err)
	}

	if err := d.Set("labels", flattenLabels(cluster.Labels)); err != nil {
		return fmt.Errorf("error setting `labels` for database user (%s): %s", d.Id(), err)
	}

	if err := d.Set("plugin", map[string]interface{}{
		"name":    "Terraform MongoDB Atlas Provider",
		"version": getPluginVersion(),
	}); err != nil {
		return fmt.Errorf("error setting `plugin` for database user (%s): %s", d.Id(), err)
	}

	d.SetId(cluster.ID)

	return nil
}
