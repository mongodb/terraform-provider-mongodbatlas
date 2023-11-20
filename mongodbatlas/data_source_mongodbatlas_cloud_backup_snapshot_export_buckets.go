package mongodbatlas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/mongodbatlas/client"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func datasourceMongoDBAtlasCloudBackupSnapshotExportBuckets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasCloudBackupSnapshotsExportBucketsRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
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
						"export_bucket_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bucket_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloud_provider": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"iam_role_id": {
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
	}
}

func dataSourceMongoDBAtlasCloudBackupSnapshotsExportBucketsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*client.MongoDBClient).Atlas

	projectID := d.Get("project_id").(string)

	options := &matlas.ListOptions{
		PageNum:      d.Get("page_num").(int),
		ItemsPerPage: d.Get("items_per_page").(int),
	}

	buckets, _, err := conn.CloudProviderSnapshotExportBuckets.List(ctx, projectID, options)
	if err != nil {
		return diag.Errorf("error getting CloudProviderSnapshotExportBuckets information: %s", err)
	}

	if err := d.Set("results", flattenCloudBackupSnapshotExportBuckets(buckets.Results)); err != nil {
		return diag.Errorf("error setting `results`: %s", err)
	}

	if err := d.Set("total_count", buckets.TotalCount); err != nil {
		return diag.Errorf("error setting `total_count`: %s", err)
	}

	d.SetId(id.UniqueId())

	return nil
}

func flattenCloudBackupSnapshotExportBuckets(buckets []*matlas.CloudProviderSnapshotExportBucket) []map[string]any {
	var results []map[string]any

	if len(buckets) == 0 {
		return results
	}

	results = make([]map[string]any, len(buckets))

	for k, bucket := range buckets {
		results[k] = map[string]any{
			"export_bucket_id": bucket.ID,
			"bucket_name":      bucket.BucketName,
			"cloud_provider":   bucket.CloudProvider,
			"iam_role_id":      bucket.IAMRoleID,
		}
	}

	return results
}
