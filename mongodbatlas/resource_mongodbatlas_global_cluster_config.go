package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"strings"

	"github.com/hashicorp/terraform/helper/schema"

	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

const (
	errorGlobalClusterCreate = "error creating MongoDB Global Cluster Configuration: %s"
	errorGlobalClusterRead   = "error reading MongoDB Global Cluster Configuration (%s): %s"
	errorGlobalClusterDelete = "error deleting MongoDB Global Cluster Configuration (%s): %s"
	errorGlobalClusterUpdate = "error updating MongoDB Global Cluster Configuration (%s): %s"
)

func resourceMongoDBAtlasGlobalCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoDBAtlasGlobalClusterCreate,
		Read:   resourceMongoDBAtlasGlobalClusterRead,
		Update: resourceMongoDBAtlasGlobalClusterUpdate,
		Delete: resourceMongoDBAtlasGlobalClusterDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMongoDBAtlasGlobalClusterImportState,
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
			"managed_namespaces": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"collection": {
							Type:     schema.TypeString,
							Required: true,
						},
						"custom_shard_key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"db": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"custom_zone_mappings": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"location": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"zone": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			//Computed:
			"custom_zone_mapping": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func resourceMongoDBAtlasGlobalClusterCreate(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)
	clusterName := d.Get("cluster_name").(string)

	if v, ok := d.GetOk("managed_namespaces"); ok {
		for _, m := range v.(*schema.Set).List() {
			mn := m.(map[string]interface{})

			log.Printf("[DEBUG] managed namespaces %+v", mn)

			addManagedNamespace := &matlas.ManagedNamespace{
				Collection:     mn["collection"].(string),
				Db:             mn["db"].(string),
				CustomShardKey: mn["custom_shard_key"].(string),
			}
			_, _, err := conn.GlobalClusters.AddManagedNamespace(context.Background(), projectID, clusterName, addManagedNamespace)

			if err != nil {
				return fmt.Errorf(errorGlobalClusterCreate, err)
			}
		}
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":   projectID,
		"cluster_name": clusterName,
	}))

	return resourceMongoDBAtlasGlobalClusterRead(d, meta)
}
func resourceMongoDBAtlasGlobalClusterRead(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	globalCluster, resp, err := conn.GlobalClusters.Get(context.Background(), projectID, clusterName)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {

			return nil
		}
		return fmt.Errorf(errorGlobalClusterRead, clusterName, err)
	}

	if err := d.Set("managed_namespaces", flattenManagedNamespaces(globalCluster.ManagedNamespaces)); err != nil {
		return fmt.Errorf(errorGlobalClusterRead, clusterName, err)
	}

	if err := d.Set("custom_zone_mapping", globalCluster.CustomZoneMapping); err != nil {
		return fmt.Errorf(errorGlobalClusterRead, clusterName, err)
	}

	return nil
}

func resourceMongoDBAtlasGlobalClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	if d.HasChange("managed_namespaces") {
		old, new := d.GetChange("managed_namespaces")
		oldSet := old.(*schema.Set)
		newSet := new.(*schema.Set)

		remove := oldSet.Difference(newSet).List()
		add := newSet.Difference(oldSet).List()

		if len(add) > 0 {
			if err := addManagedNamespaces(conn, add, projectID, clusterName); err != nil {
				return fmt.Errorf(errorGlobalClusterUpdate, clusterName, err)
			}

		}

		if len(remove) > 0 {
			if err := removeManagedNamespaces(conn, remove, projectID, clusterName); err != nil {
				return fmt.Errorf(errorGlobalClusterUpdate, clusterName, err)
			}
		}
	}

	return resourceMongoDBAtlasGlobalClusterRead(d, meta)
}

func resourceMongoDBAtlasGlobalClusterDelete(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	if v, ok := d.GetOk("managed_namespaces"); ok {
		if err := removeManagedNamespaces(conn, v.(*schema.Set).List(), projectID, clusterName); err != nil {
			return fmt.Errorf(errorGlobalClusterDelete, clusterName, err)
		}
	}

	if v, ok := d.GetOk("custom_zone_mappings"); ok {
		if v.(*schema.Set).Len() > 0 {
			if _, _, err := conn.GlobalClusters.DeleteCustomZoneMappings(context.Background(), projectID, clusterName); err != nil {
				return fmt.Errorf(errorGlobalClusterDelete, clusterName, err)
			}
		}
	}

	return nil
}

func flattenManagedNamespaces(managedNamespaces []matlas.ManagedNamespace) []map[string]interface{} {
	var results []map[string]interface{}

	if len(managedNamespaces) > 0 {
		results = make([]map[string]interface{}, len(managedNamespaces))

		for k, managedNamespace := range managedNamespaces {
			results[k] = map[string]interface{}{
				"db":               managedNamespace.Db,
				"collection":       managedNamespace.Collection,
				"custom_shard_key": managedNamespace.CustomShardKey,
			}
		}
	}
	return results
}

func resourceMongoDBAtlasGlobalClusterImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	//conn := meta.(*matlas.Client)

	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 2 {
		return nil, errors.New("import format error: to import a global cluster, use the format {project_id}-{cluster-name}")
	}

	projectID := parts[0]
	name := parts[1]

	d.SetId(encodeStateID(map[string]string{
		"project_id":   projectID,
		"cluster_name": name,
	}))

	return []*schema.ResourceData{d}, nil
}

func removeManagedNamespaces(conn *matlas.Client, remove []interface{}, projectID, clusterName string) error {
	for _, m := range remove {
		mn := m.(map[string]interface{})
		addManagedNamespace := &matlas.ManagedNamespace{
			Collection:     mn["collection"].(string),
			Db:             mn["db"].(string),
			CustomShardKey: mn["custom_shard_key"].(string),
		}
		_, _, err := conn.GlobalClusters.DeleteManagedNamespace(context.Background(), projectID, clusterName, addManagedNamespace)

		if err != nil {
			return err
		}
	}
	return nil
}

func addManagedNamespaces(conn *matlas.Client, add []interface{}, projectID, clusterName string) error {
	for _, m := range add {
		mn := m.(map[string]interface{})

		addManagedNamespace := &matlas.ManagedNamespace{
			Collection:     mn["collection"].(string),
			Db:             mn["db"].(string),
			CustomShardKey: mn["custom_shard_key"].(string),
		}
		_, _, err := conn.GlobalClusters.AddManagedNamespace(context.Background(), projectID, clusterName, addManagedNamespace)

		if err != nil {
			return err
		}
	}

	return nil
}
