package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mwielbut/pointy"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorGlobalClusterCreate = "error creating MongoDB Global Cluster Configuration: %s"
	errorGlobalClusterRead   = "error reading MongoDB Global Cluster Configuration (%s): %s"
	errorGlobalClusterDelete = "error deleting MongoDB Global Cluster Configuration (%s): %s"
	errorGlobalClusterUpdate = "error updating MongoDB Global Cluster Configuration (%s): %s"
)

func resourceMongoDBAtlasGlobalCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBAtlasGlobalClusterCreate,
		ReadContext:   resourceMongoDBAtlasGlobalClusterRead,
		UpdateContext: resourceMongoDBAtlasGlobalClusterUpdate,
		DeleteContext: resourceMongoDBAtlasGlobalClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasGlobalClusterImportState,
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
						"is_custom_shard_key_hashed": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"is_shard_key_unique": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
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
			"custom_zone_mapping": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func resourceMongoDBAtlasGlobalClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
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

			if isCustomShardKeyHashed, okCustomShard := mn["is_custom_shard_key_hashed"]; okCustomShard {
				addManagedNamespace.IsCustomShardKeyHashed = pointy.Bool(isCustomShardKeyHashed.(bool))
			}

			if isShardKeyUnique, okShard := mn["is_shard_key_unique"]; okShard {
				addManagedNamespace.IsShardKeyUnique = pointy.Bool(isShardKeyUnique.(bool))
			}

			_, _, err := conn.GlobalClusters.AddManagedNamespace(ctx, projectID, clusterName, addManagedNamespace)

			if err != nil {
				return diag.FromErr(fmt.Errorf(errorGlobalClusterCreate, err))
			}
		}
	}

	if v, ok := d.GetOk("custom_zone_mappings"); ok {
		var customZoneMappings []matlas.CustomZoneMapping

		for _, czms := range v.(*schema.Set).List() {
			cz := czms.(map[string]interface{})

			log.Printf("[DEBUG] custom zone mappings %+v", cz)

			customZoneMapping := matlas.CustomZoneMapping{
				Location: cz["location"].(string),
				Zone:     cz["zone"].(string),
			}

			customZoneMappings = append(customZoneMappings, customZoneMapping)
		}

		if len(customZoneMappings) > 0 {
			_, _, err := conn.GlobalClusters.AddCustomZoneMappings(ctx, projectID, clusterName, &matlas.CustomZoneMappingsRequest{
				CustomZoneMappings: customZoneMappings,
			})

			if err != nil {
				return diag.FromErr(fmt.Errorf(errorGlobalClusterCreate, err))
			}
		}
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":   projectID,
		"cluster_name": clusterName,
	}))

	return resourceMongoDBAtlasGlobalClusterRead(ctx, d, meta)
}

func resourceMongoDBAtlasGlobalClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	globalCluster, resp, err := conn.GlobalClusters.Get(ctx, projectID, clusterName)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorGlobalClusterRead, clusterName, err))
	}

	if err := d.Set("managed_namespaces", flattenManagedNamespaces(globalCluster.ManagedNamespaces)); err != nil {
		return diag.FromErr(fmt.Errorf(errorGlobalClusterRead, clusterName, err))
	}

	if err := d.Set("custom_zone_mapping", globalCluster.CustomZoneMapping); err != nil {
		return diag.FromErr(fmt.Errorf(errorGlobalClusterRead, clusterName, err))
	}

	return nil
}

func resourceMongoDBAtlasGlobalClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	if d.HasChange("managed_namespaces") {
		old, newMN := d.GetChange("managed_namespaces")
		oldSet := old.(*schema.Set)
		newSet := newMN.(*schema.Set)

		remove := oldSet.Difference(newSet).List()
		add := newSet.Difference(oldSet).List()

		if len(remove) > 0 {
			if err := removeManagedNamespaces(ctx, conn, remove, projectID, clusterName); err != nil {
				return diag.FromErr(fmt.Errorf(errorGlobalClusterUpdate, clusterName, err))
			}
		}

		if len(add) > 0 {
			if err := addManagedNamespaces(ctx, conn, add, projectID, clusterName); err != nil {
				return diag.FromErr(fmt.Errorf(errorGlobalClusterUpdate, clusterName, err))
			}
		}
	}

	return resourceMongoDBAtlasGlobalClusterRead(ctx, d, meta)
}

func resourceMongoDBAtlasGlobalClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	if v, ok := d.GetOk("managed_namespaces"); ok {
		if err := removeManagedNamespaces(ctx, conn, v.(*schema.Set).List(), projectID, clusterName); err != nil {
			return diag.FromErr(fmt.Errorf(errorGlobalClusterDelete, clusterName, err))
		}
	}

	if v, ok := d.GetOk("custom_zone_mappings"); ok {
		if v.(*schema.Set).Len() > 0 {
			if _, _, err := conn.GlobalClusters.DeleteCustomZoneMappings(ctx, projectID, clusterName); err != nil {
				return diag.FromErr(fmt.Errorf(errorGlobalClusterDelete, clusterName, err))
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
				"db":                         managedNamespace.Db,
				"collection":                 managedNamespace.Collection,
				"custom_shard_key":           managedNamespace.CustomShardKey,
				"is_custom_shard_key_hashed": *managedNamespace.IsCustomShardKeyHashed,
				"is_shard_key_unique":        *managedNamespace.IsShardKeyUnique,
			}
		}
	}

	return results
}

func resourceMongoDBAtlasGlobalClusterImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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

	if err := d.Set("project_id", projectID); err != nil {
		return []*schema.ResourceData{d}, err
	}

	if err := d.Set("cluster_name", name); err != nil {
		return []*schema.ResourceData{d}, err
	}

	return []*schema.ResourceData{d}, nil
}

func removeManagedNamespaces(ctx context.Context, conn *matlas.Client, remove []interface{}, projectID, clusterName string) error {
	for _, m := range remove {
		mn := m.(map[string]interface{})
		addManagedNamespace := &matlas.ManagedNamespace{
			Collection:     mn["collection"].(string),
			Db:             mn["db"].(string),
			CustomShardKey: mn["custom_shard_key"].(string),
		}

		if isCustomShardKeyHashed, okCustomShard := mn["is_custom_shard_key_hashed"]; okCustomShard {
			addManagedNamespace.IsCustomShardKeyHashed = pointy.Bool(isCustomShardKeyHashed.(bool))
		}

		if isShardKeyUnique, okShard := mn["is_shard_key_unique"]; okShard {
			addManagedNamespace.IsShardKeyUnique = pointy.Bool(isShardKeyUnique.(bool))
		}
		_, _, err := conn.GlobalClusters.DeleteManagedNamespace(ctx, projectID, clusterName, addManagedNamespace)

		if err != nil {
			return err
		}
	}

	return nil
}

func addManagedNamespaces(ctx context.Context, conn *matlas.Client, add []interface{}, projectID, clusterName string) error {
	for _, m := range add {
		mn := m.(map[string]interface{})

		addManagedNamespace := &matlas.ManagedNamespace{
			Collection:     mn["collection"].(string),
			Db:             mn["db"].(string),
			CustomShardKey: mn["custom_shard_key"].(string),
		}

		if isCustomShardKeyHashed, okCustomShard := mn["is_custom_shard_key_hashed"]; okCustomShard {
			addManagedNamespace.IsCustomShardKeyHashed = pointy.Bool(isCustomShardKeyHashed.(bool))
		}

		if isShardKeyUnique, okShard := mn["is_shard_key_unique"]; okShard {
			addManagedNamespace.IsShardKeyUnique = pointy.Bool(isShardKeyUnique.(bool))
		}
		_, _, err := conn.GlobalClusters.AddManagedNamespace(ctx, projectID, clusterName, addManagedNamespace)

		if err != nil {
			return err
		}
	}

	return nil
}
