package cloudbackupsnapshotexportjob

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasCloudBackupSnapshotsExportJobRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: fmt.Sprintf(constant.DeprecationParamByVersion, "1.18.0") + " Will not be an input parameter, only computed.",
			},
			"export_job_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_data": {
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
					}},
			},
			"components": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"export_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"replica_set_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"err_msg": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: fmt.Sprintf(constant.DeprecationParamByVersion, "1.20.0"),
			},
			"export_bucket_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"export_status_exported_collections": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"export_status_total_collections": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"finished_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"prefix": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceMongoDBAtlasCloudBackupSnapshotsExportJobRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	exportJob, err := readExportJob(ctx, meta, d)
	if err != nil {
		return diag.Errorf("error reading snapshot export job information: %s", err)
	}
	return setExportJobFields(d, exportJob)
}
