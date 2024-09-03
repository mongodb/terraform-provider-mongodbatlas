package cloudbackupsnapshotexportbucket

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20240805003/admin"
)

func PluralDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
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
						"role_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
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

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2

	projectID := d.Get("project_id").(string)
	itemsPerPage := d.Get("items_per_page").(int)
	pageNum := d.Get("page_num").(int)

	buckets, _, err := conn.CloudBackupsApi.ListExportBuckets(ctx, projectID).ItemsPerPage(itemsPerPage).PageNum(pageNum).Execute()
	if err != nil {
		return diag.Errorf("error getting CloudProviderSnapshotExportBuckets information: %s", err)
	}

	if err := d.Set("results", flattenBuckets(buckets.GetResults())); err != nil {
		return diag.Errorf("error setting `results`: %s", err)
	}

	if err := d.Set("total_count", buckets.GetTotalCount()); err != nil {
		return diag.Errorf("error setting `total_count`: %s", err)
	}

	d.SetId(id.UniqueId())

	return nil
}

func flattenBuckets(buckets []admin.DiskBackupSnapshotExportBucket) []map[string]any {
	var results []map[string]any

	if len(buckets) == 0 {
		return results
	}

	results = make([]map[string]any, len(buckets))

	for k, bucket := range buckets {
		results[k] = map[string]any{
			"export_bucket_id": bucket.GetId(),
			"bucket_name":      bucket.GetBucketName(),
			"cloud_provider":   bucket.GetCloudProvider(),
			"iam_role_id":      bucket.GetIamRoleId(),
			"role_id":          bucket.GetRoleId(),
			"service_url":      bucket.GetServiceUrl(),
			"tenant_id":        bucket.GetTenantId(),
		}
	}

	return results
}
