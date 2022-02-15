package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func resourceMongoDBAtlasCloudBackupSnapshotExportJob() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBAtlasCloudBackupSnapshotExportJobCreate,
		ReadContext:   resourceMongoDBAtlasCloudBackupSnapshotExportJobRead,
		DeleteContext: schema.NoopContext,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasCloudBackupSnapshotExportJobImportState,
		},
		Schema: returnCloudBackupSnapshotExportJobSchema(),
	}
}

func returnCloudBackupSnapshotExportJobSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"export_job_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
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
		"snapshot_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"export_bucket_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"custom_data": {
			Type:     schema.TypeSet,
			Required: true,
			ForceNew: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:     schema.TypeString,
						Required: true,
						ForceNew: true,
					},
					"value": {
						Type:     schema.TypeString,
						Required: true,
						ForceNew: true,
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
	}
}

func resourceMongoDBAtlasCloudBackupSnapshotExportJobRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]
	exportJobID := ids["export_job_id"]

	exportJob, _, err := conn.CloudProviderSnapshotExportJobs.Get(ctx, projectID, clusterName, exportJobID)
	if err != nil {
		// case 404
		// deleted in the backend case
		reset := strings.Contains(err.Error(), "404") && !d.IsNewResource()

		if reset {
			d.SetId("")
			return nil
		}

		return diag.Errorf("error getting snapshot export job information: %s", err)
	}

	if err := d.Set("export_job_id", exportJob.ID); err != nil {
		return diag.Errorf("error setting `export_job_id` for snapshot export job (%s): %s", d.Id(), err)
	}

	if err := d.Set("snapshot_id", exportJob.SnapshotID); err != nil {
		return diag.Errorf("error setting `snapshot_id` for snapshot export job (%s): %s", d.Id(), err)
	}

	if err := d.Set("custom_data", flattenExportJobsCustomData(exportJob.CustomData)); err != nil {
		return diag.Errorf("error setting `custom_data` for snapshot export job (%s): %s", d.Id(), err)
	}

	if err := d.Set("components", flattenExportJobsComponents(exportJob.Components)); err != nil {
		return diag.Errorf("error setting `components` for snapshot export job (%s): %s", d.Id(), err)
	}

	if err := d.Set("created_at", exportJob.CreatedAt); err != nil {
		return diag.Errorf("error setting `created_at` for snapshot export job (%s): %s", d.Id(), err)
	}

	if err := d.Set("err_msg", exportJob.ErrMsg); err != nil {
		return diag.Errorf("error setting `created_at` for snapshot export job (%s): %s", d.Id(), err)
	}

	if err := d.Set("export_bucket_id", exportJob.ExportBucketID); err != nil {
		return diag.Errorf("error setting `created_at` for snapshot export job (%s): %s", d.Id(), err)
	}

	if exportJob.ExportStatus != nil {
		if err := d.Set("export_status_exported_collections", exportJob.ExportStatus.ExportedCollections); err != nil {
			return diag.Errorf("error setting `export_status_exported_collections` for snapshot export job (%s): %s", d.Id(), err)
		}

		if err := d.Set("export_status_total_collections", exportJob.ExportStatus.TotalCollections); err != nil {
			return diag.Errorf("error setting `export_status_total_collections` for snapshot export job (%s): %s", d.Id(), err)
		}
	}

	if err := d.Set("finished_at", exportJob.FinishedAt); err != nil {
		return diag.Errorf("error setting `finished_at` for snapshot export job (%s): %s", d.Id(), err)
	}

	if err := d.Set("prefix", exportJob.Prefix); err != nil {
		return diag.Errorf("error setting `prefix` for snapshot export job (%s): %s", d.Id(), err)
	}

	if err := d.Set("state", exportJob.State); err != nil {
		return diag.Errorf("error setting `prefix` for snapshot export job (%s): %s", d.Id(), err)
	}

	return nil
}

func flattenExportJobsComponents(components []*matlas.CloudProviderSnapshotExportJobComponent) []map[string]interface{} {
	if len(components) == 0 {
		return nil
	}

	customData := make([]map[string]interface{}, 0)

	for i := range components {
		customData = append(customData, map[string]interface{}{
			"export_id":        components[i].ExportID,
			"replica_set_name": components[i].ReplicaSetName,
		})
	}

	return customData
}

func flattenExportJobsCustomData(data []*matlas.CloudProviderSnapshotExportJobCustomData) []map[string]interface{} {
	if len(data) == 0 {
		return nil
	}

	customData := make([]map[string]interface{}, 0)

	for i := range data {
		customData = append(customData, map[string]interface{}{
			"key":   data[i].Key,
			"value": data[i].Value,
		})
	}

	return customData
}

func resourceMongoDBAtlasCloudBackupSnapshotExportJobCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)
	clusterName := d.Get("cluster_name").(string)

	request := &matlas.CloudProviderSnapshotExportJob{
		SnapshotID:     d.Get("snapshot_id").(string),
		ExportBucketID: d.Get("export_bucket_id").(string),
		CustomData:     expandExportJobCustomData(d),
	}

	jobResponse, _, err := conn.CloudProviderSnapshotExportJobs.Create(ctx, projectID, clusterName, request)
	if err != nil {
		return diag.Errorf("error creating snapshot export job: %s", err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":    projectID,
		"cluster_name":  clusterName,
		"export_job_id": jobResponse.ID,
	}))

	return resourceMongoDBAtlasCloudBackupSnapshotExportJobRead(ctx, d, meta)
}

func expandExportJobCustomData(d *schema.ResourceData) []*matlas.CloudProviderSnapshotExportJobCustomData {
	customData := d.Get("custom_data").(*schema.Set)
	res := make([]*matlas.CloudProviderSnapshotExportJobCustomData, customData.Len())

	for i, val := range customData.List() {
		v := val.(map[string]interface{})
		res[i] = &matlas.CloudProviderSnapshotExportJobCustomData{
			Key:   v["key"].(string),
			Value: v["value"].(string),
		}
	}

	return res
}

func resourceMongoDBAtlasCloudBackupSnapshotExportJobImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*MongoDBClient).Atlas

	parts := strings.SplitN(d.Id(), "--", 3)
	if len(parts) != 3 {
		return nil, errors.New("import format error: to import snapshot export job, use the format {project_id}--{cluster_name}--{export_job_id}")
	}

	projectID := parts[0]
	clusterName := parts[1]
	exportJobID := parts[2]

	_, _, err := conn.CloudProviderSnapshotExportJobs.Get(ctx, projectID, clusterName, exportJobID)
	if err != nil {
		return nil, fmt.Errorf("couldn't import snapshot export job %s in project %s and cluster %s, error: %s", exportJobID, projectID, clusterName, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":    projectID,
		"cluster_name":  clusterName,
		"export_job_id": exportJobID,
	}))

	return []*schema.ResourceData{d}, nil
}
