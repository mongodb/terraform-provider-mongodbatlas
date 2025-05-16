package globalclusterconfig

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"go.mongodb.org/atlas-sdk/v20250312002/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	errorGlobalClusterCreate = "error creating MongoDB Global Cluster Configuration: %s"
	errorGlobalClusterRead   = "error reading MongoDB Global Cluster Configuration (%s): %s"
	errorGlobalClusterDelete = "error deleting MongoDB Global Cluster Configuration (%s): %s"
	errorGlobalClusterUpdate = "error updating MongoDB Global Cluster Configuration (%s): %s"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate,
		ReadContext:   resourceRead,
		UpdateContext: resourceUpdate,
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
				Deprecated: fmt.Sprintf(constant.DeprecationParamByVersionWithReplacement, "1.23.0", "custom_zone_mapping_zone_id"),
				Type:       schema.TypeMap,
				Computed:   true,
			},
			"custom_zone_mapping_zone_id": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)
	clusterName := d.Get("cluster_name").(string)

	if v, ok := d.GetOk("managed_namespaces"); ok {
		for _, m := range v.(*schema.Set).List() {
			mn := m.(map[string]any)

			req := &admin.ManagedNamespaces{
				Collection:     mn["collection"].(string),
				Db:             mn["db"].(string),
				CustomShardKey: mn["custom_shard_key"].(string),
			}

			if isCustomShardKeyHashed, okCustomShard := mn["is_custom_shard_key_hashed"]; okCustomShard {
				req.IsCustomShardKeyHashed = conversion.Pointer[bool](isCustomShardKeyHashed.(bool))
			}

			if isShardKeyUnique, okShard := mn["is_shard_key_unique"]; okShard {
				req.IsShardKeyUnique = conversion.Pointer[bool](isShardKeyUnique.(bool))
			}

			err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
				_, _, err := connV2.GlobalClustersApi.CreateManagedNamespace(ctx, projectID, clusterName, req).Execute()
				if err != nil {
					if admin.IsErrorCode(err, "DUPLICATE_MANAGED_NAMESPACE") {
						if err := removeManagedNamespaces(ctx, connV2, v.(*schema.Set).List(), projectID, clusterName); err != nil {
							return retry.NonRetryableError(fmt.Errorf(errorGlobalClusterCreate, err))
						}
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(fmt.Errorf(errorGlobalClusterCreate, err))
				}
				return nil
			})
			if err != nil {
				return diag.FromErr(fmt.Errorf(errorGlobalClusterCreate, err))
			}
		}
	}

	if v, ok := d.GetOk("custom_zone_mappings"); ok {
		_, _, err := connV2.GlobalClustersApi.CreateCustomZoneMapping(ctx, projectID, clusterName, &admin.CustomZoneMappings{
			CustomZoneMappings: newCustomZoneMappings(v.(*schema.Set).List()),
		}).Execute()

		if err != nil {
			if v2, ok2 := d.GetOk("managed_namespaces"); ok2 {
				if err := removeManagedNamespaces(ctx, connV2, v2.(*schema.Set).List(), projectID, clusterName); err != nil {
					return diag.FromErr(fmt.Errorf(errorGlobalClusterCreate, err))
				}
			}
			return diag.FromErr(fmt.Errorf(errorGlobalClusterCreate, err))
		}
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":   projectID,
		"cluster_name": clusterName,
	}))

	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	notFound, err := readGlobalClusterConfig(ctx, meta, projectID, clusterName, d)
	if err != nil {
		return diag.FromErr(err)
	}
	if notFound {
		d.SetId("")
	}
	return nil
}

func readGlobalClusterConfig(ctx context.Context, meta any, projectID, clusterName string, d *schema.ResourceData) (notFound bool, err error) {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	connV220240530 := meta.(*config.MongoDBClient).AtlasV220240530
	resp, httpResp, err := connV2.GlobalClustersApi.GetManagedNamespace(ctx, projectID, clusterName).Execute()
	if err != nil {
		if validate.StatusNotFound(httpResp) {
			return true, nil
		}
		return false, fmt.Errorf(errorGlobalClusterRead, clusterName, err)
	}
	if err := d.Set("managed_namespaces", flattenManagedNamespaces(resp.GetManagedNamespaces())); err != nil {
		return false, fmt.Errorf(errorGlobalClusterRead, clusterName, err)
	}
	if err := d.Set("custom_zone_mapping_zone_id", resp.GetCustomZoneMapping()); err != nil {
		return false, fmt.Errorf(errorGlobalClusterRead, clusterName, err)
	}

	oldResp, httpResp, err := connV220240530.GlobalClustersApi.GetManagedNamespace(ctx, projectID, clusterName).Execute()
	if err != nil {
		if validate.StatusNotFound(httpResp) {
			return true, nil
		}
		if validate.ErrorClusterIsAsymmetrics(err) {
			// Avoid non-empty plan by setting an empty custom_zone_mapping.
			if err := d.Set("custom_zone_mapping", map[string]string{}); err != nil {
				return false, fmt.Errorf(errorGlobalClusterRead, clusterName, err)
			}
			return false, nil
		}
		return false, fmt.Errorf(errorGlobalClusterRead, clusterName, err)
	}
	if err := d.Set("custom_zone_mapping", oldResp.GetCustomZoneMapping()); err != nil {
		return false, fmt.Errorf(errorGlobalClusterRead, clusterName, err)
	}
	return false, nil
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	if d.HasChange("managed_namespaces") {
		oldMN, newMN := d.GetChange("managed_namespaces")
		oldList := oldMN.(*schema.Set).List()
		newList := newMN.(*schema.Set).List()
		if err := updateManagedNamespaces(ctx, connV2, projectID, clusterName, oldList, newList); err != nil {
			return diag.FromErr(fmt.Errorf(errorGlobalClusterUpdate, clusterName, err))
		}
	}

	if d.HasChange("custom_zone_mappings") {
		oldZN, newZN := d.GetChange("custom_zone_mappings")
		oldSet := oldZN.(*schema.Set)
		newSet := newZN.(*schema.Set)
		if err := updateCustomZoneMappings(ctx, connV2, projectID, clusterName, oldSet, newSet); err != nil {
			return diag.FromErr(fmt.Errorf(errorGlobalClusterUpdate, clusterName, err))
		}
	}
	return resourceRead(ctx, d, meta)
}

// convertInterfaceSlice is a helper function that converts []map[string]any into []any
func convertInterfaceSlice(input []map[string]any) []any {
	var out []any
	for _, v := range input {
		out = append(out, v)
	}
	return out
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	if v, ok := d.GetOk("managed_namespaces"); ok {
		if err := removeManagedNamespaces(ctx, connV2, v.(*schema.Set).List(), projectID, clusterName); err != nil {
			return diag.FromErr(fmt.Errorf(errorGlobalClusterDelete, clusterName, err))
		}
	}

	if v, ok := d.GetOk("custom_zone_mappings"); ok {
		if v.(*schema.Set).Len() > 0 {
			if _, _, err := connV2.GlobalClustersApi.DeleteAllCustomZoneMappings(ctx, projectID, clusterName).Execute(); err != nil {
				return diag.FromErr(fmt.Errorf(errorGlobalClusterDelete, clusterName, err))
			}
		}
	}

	return nil
}

func flattenManagedNamespaces(managedNamespaces []admin.ManagedNamespaces) []map[string]any {
	var results []map[string]any

	if len(managedNamespaces) > 0 {
		results = make([]map[string]any, len(managedNamespaces))

		for k, managedNamespace := range managedNamespaces {
			results[k] = map[string]any{
				"db":                         managedNamespace.GetDb(),
				"collection":                 managedNamespace.GetCollection(),
				"custom_shard_key":           managedNamespace.GetCustomShardKey(),
				"is_custom_shard_key_hashed": managedNamespace.GetIsCustomShardKeyHashed(),
				"is_shard_key_unique":        managedNamespace.GetIsShardKeyUnique(),
			}
		}
	}
	return results
}

func resourceImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 2 {
		return nil, errors.New("import format error: to import a global cluster, use the format {project_id}-{cluster-name}")
	}

	projectID := parts[0]
	name := parts[1]

	d.SetId(conversion.EncodeStateID(map[string]string{
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

func removeManagedNamespaces(ctx context.Context, connV2 *admin.APIClient, remove []any, projectID, clusterName string) error {
	for _, m := range remove {
		mn := m.(map[string]any)
		managedNamespace := &admin.DeleteManagedNamespaceApiParams{
			Collection:  conversion.StringPtr(mn["collection"].(string)),
			Db:          conversion.StringPtr(mn["db"].(string)),
			ClusterName: clusterName,
			GroupId:     projectID,
		}

		_, _, err := connV2.GlobalClustersApi.DeleteManagedNamespaceWithParams(ctx, managedNamespace).Execute()

		if err != nil {
			return err
		}
	}
	return nil
}

func newCustomZoneMapping(tfMap map[string]any) *admin.ZoneMapping {
	if tfMap == nil {
		return nil
	}

	apiObject := &admin.ZoneMapping{
		Location: tfMap["location"].(string),
		Zone:     tfMap["zone"].(string),
	}

	return apiObject
}

func newCustomZoneMappings(tfList []any) *[]admin.ZoneMapping {
	if len(tfList) == 0 {
		return nil
	}

	apiObjects := make([]admin.ZoneMapping, len(tfList))
	if len(tfList) > 0 {
		for i, tfMapRaw := range tfList {
			if tfMap, ok := tfMapRaw.(map[string]any); ok {
				apiObject := newCustomZoneMapping(tfMap)
				if apiObject == nil {
					continue
				}
				apiObjects[i] = *apiObject
			}
		}
	}

	return &apiObjects
}

func addManagedNamespaces(ctx context.Context, connV2 *admin.APIClient, add []any, projectID, clusterName string) error {
	for _, m := range add {
		mn := m.(map[string]any)

		addManagedNamespace := &admin.ManagedNamespaces{
			Collection:     mn["collection"].(string),
			Db:             mn["db"].(string),
			CustomShardKey: mn["custom_shard_key"].(string),
		}
		if isCustomShardKeyHashed, okCustomShard := mn["is_custom_shard_key_hashed"]; okCustomShard {
			addManagedNamespace.IsCustomShardKeyHashed = conversion.Pointer[bool](isCustomShardKeyHashed.(bool))
		}
		if isShardKeyUnique, okShard := mn["is_shard_key_unique"]; okShard {
			addManagedNamespace.IsShardKeyUnique = conversion.Pointer[bool](isShardKeyUnique.(bool))
		}
		_, _, err := connV2.GlobalClustersApi.CreateManagedNamespace(ctx, projectID, clusterName, addManagedNamespace).Execute()
		if err != nil {
			return err
		}
	}
	return nil
}

// buildManagedNamespacesMap converts a list of managed_namespace entries into a map keyed by "collection:db"
func buildManagedNamespacesMap(list []any) map[string]map[string]any {
	namespacesMap := make(map[string]map[string]any)
	for _, item := range list {
		m := item.(map[string]any)
		key := fmt.Sprintf("%s:%s", m["collection"].(string), m["db"].(string))
		namespacesMap[key] = m
	}
	return namespacesMap
}

// diffManagedNamespaces calculates the difference between old and new managed_namespaces.
// Returns slices of namespaces to add and remove; errors out on modifications.
func diffManagedNamespaces(oldList, newList []any) (toAdd, toRemove []map[string]any, err error) {
	oldMap := buildManagedNamespacesMap(oldList)
	newMap := buildManagedNamespacesMap(newList)
	for key, oldEntry := range oldMap {
		if newEntry, exists := newMap[key]; exists {
			// Modification is not allowed.
			if !reflect.DeepEqual(oldEntry, newEntry) {
				return nil, nil, fmt.Errorf("managed namespace for collection '%s' in db '%s' cannot be modified", oldEntry["collection"], oldEntry["db"])
			}
		} else {
			toRemove = append(toRemove, oldEntry)
		}
	}
	for key, newEntry := range newMap {
		if _, exists := oldMap[key]; !exists {
			toAdd = append(toAdd, newEntry)
		}
	}
	return toAdd, toRemove, nil
}

// updateManagedNamespaces encapsulates diffing and applying removals/additions.
func updateManagedNamespaces(ctx context.Context, connV2 *admin.APIClient, projectID, clusterName string, oldList, newList []any) error {
	toAdd, toRemove, err := diffManagedNamespaces(oldList, newList)
	if err != nil {
		return err
	}
	if len(toRemove) > 0 {
		if err := removeManagedNamespaces(ctx, connV2, convertInterfaceSlice(toRemove), projectID, clusterName); err != nil {
			return err
		}
	}
	if len(toAdd) > 0 {
		if err := addManagedNamespaces(ctx, connV2, convertInterfaceSlice(toAdd), projectID, clusterName); err != nil {
			return err
		}
	}
	return nil
}

// updateCustomZoneMappings encapsulates diffing and applying changes for custom_zone_mappings.
func updateCustomZoneMappings(ctx context.Context, connV2 *admin.APIClient, projectID, clusterName string, oldSet, newSet *schema.Set) error {
	removed := oldSet.Difference(newSet).List()
	added := newSet.Difference(oldSet).List()

	if len(removed) > 0 {
		// Allow deletion only if all mappings are removed.
		if newSet.Len() != 0 {
			return fmt.Errorf("partial deletion of custom_zone_mappings is not allowed; remove either all mappings or none")
		}
		if _, _, err := connV2.GlobalClustersApi.DeleteAllCustomZoneMappings(ctx, projectID, clusterName).Execute(); err != nil {
			return err
		}
	}
	if len(added) > 0 {
		if _, _, err := connV2.GlobalClustersApi.CreateCustomZoneMapping(ctx, projectID, clusterName, &admin.CustomZoneMappings{
			CustomZoneMappings: newCustomZoneMappings(added),
		}).Execute(); err != nil {
			return err
		}
	}
	return nil
}
