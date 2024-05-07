package cloudbackupsnapshotexportbucket

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"export_bucket_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
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
	}
}

func datasourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2

	projectID := d.Get("project_id").(string)
	bucketID := d.Get("id").(string)

	bucket, _, err := conn.CloudBackupsApi.GetExportBucket(ctx, projectID, bucketID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting CloudProviderSnapshotExportBuckets Information: %s", err))
	}

	if err = d.Set("export_bucket_id", bucket.GetId()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `export_bucket_id` for CloudProviderSnapshotExportBuckets (%s): %s", d.Id(), err))
	}

	if err = d.Set("bucket_name", bucket.GetBucketName()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `bucket_name` for CloudProviderSnapshotExportBuckets (%s): %s", d.Id(), err))
	}

	if err = d.Set("cloud_provider", bucket.GetCloudProvider()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `cloud_provider` for CloudProviderSnapshotExportBuckets (%s): %s", d.Id(), err))
	}

	if err = d.Set("iam_role_id", bucket.GetIamRoleId()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `iam_role_id` for CloudProviderSnapshotExportBuckets (%s): %s", d.Id(), err))
	}

	d.SetId(bucket.GetId())

	return nil
}
