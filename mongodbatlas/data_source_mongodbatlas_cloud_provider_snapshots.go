package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasCloudProviderSnapshots() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasCloudProviderSnapshotsRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"page_num": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"items_per_page": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expires_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_key_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mongod_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snapshot_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_size_bytes": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
		DeprecationMessage: "This data source is deprecated. Please transition to mongodbatlas_cloud_backup_snapshots as soon as possible",
	}
}

func dataSourceMongoDBAtlasCloudProviderSnapshotsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas

	requestParameters := &matlas.SnapshotReqPathParameters{
		GroupID:     d.Get("project_id").(string),
		ClusterName: d.Get("cluster_name").(string),
	}
	options := &matlas.ListOptions{
		PageNum:      d.Get("page_num").(int),
		ItemsPerPage: d.Get("items_per_page").(int),
	}

	cloudProviderSnapshots, _, err := conn.CloudProviderSnapshots.GetAllCloudProviderSnapshots(ctx, requestParameters, options)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting cloudProviderSnapshots information: %s", err))
	}

	if err := d.Set("results", flattenCloudProviderSnapshots(cloudProviderSnapshots.Results)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `results`: %s", err))
	}

	if err := d.Set("total_count", cloudProviderSnapshots.TotalCount); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `total_count`: %s", err))
	}

	d.SetId(resource.UniqueId())

	return nil
}

func flattenCloudProviderSnapshots(cloudProviderSnapshots []*matlas.CloudProviderSnapshot) []map[string]interface{} {
	var results []map[string]interface{}

	if len(cloudProviderSnapshots) > 0 {
		results = make([]map[string]interface{}, len(cloudProviderSnapshots))

		for k, cloudProviderSnapshot := range cloudProviderSnapshots {
			results[k] = map[string]interface{}{
				"id":                 cloudProviderSnapshot.ID,
				"created_at":         cloudProviderSnapshot.CreatedAt,
				"description":        cloudProviderSnapshot.Description,
				"expires_at":         cloudProviderSnapshot.ExpiresAt,
				"master_key_uuid":    cloudProviderSnapshot.MasterKeyUUID,
				"mongod_version":     cloudProviderSnapshot.MongodVersion,
				"snapshot_type":      cloudProviderSnapshot.SnapshotType,
				"status":             cloudProviderSnapshot.Status,
				"storage_size_bytes": cloudProviderSnapshot.StorageSizeBytes,
				"type":               cloudProviderSnapshot.Type,
			}
		}
	}

	return results
}
