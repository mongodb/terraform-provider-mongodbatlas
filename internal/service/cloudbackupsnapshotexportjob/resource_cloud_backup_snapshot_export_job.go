package cloudbackupsnapshotexportjob

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312004/admin"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate,
		ReadContext:   resourceRead,
		DeleteContext: resourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceImportState,
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
			Optional: true,
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

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	exportJob, err := readExportJob(ctx, meta, d)
	if err != nil {
		reset := strings.Contains(err.Error(), "404") && !d.IsNewResource()

		if reset {
			d.SetId("")
			return nil
		}

		return diag.Errorf("error getting snapshot export job information: %s", err)
	}
	return setExportJobFields(d, exportJob)
}

func readExportJob(ctx context.Context, meta any, d *schema.ResourceData) (*admin.DiskBackupExportJob, error) {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID, clusterName, exportID := getRequiredFields(d)
	if d.Id() != "" && (projectID == "" || clusterName == "" || exportID == "") {
		ids := conversion.DecodeStateID(d.Id())
		projectID = ids["project_id"]
		clusterName = ids["cluster_name"]
		exportID = ids["export_job_id"]
	}
	exportJob, _, err := connV2.CloudBackupsApi.GetBackupExportJob(ctx, projectID, clusterName, exportID).Execute()
	if err == nil {
		d.SetId(conversion.EncodeStateID(map[string]string{
			"project_id":    projectID,
			"cluster_name":  clusterName,
			"export_job_id": exportJob.GetId(),
		}))
	}
	return exportJob, err
}

func getRequiredFields(d *schema.ResourceData) (projectID, clusterName, exportID string) {
	projectID = d.Get("project_id").(string)
	clusterName = d.Get("cluster_name").(string)
	exportID = d.Get("export_job_id").(string)
	return projectID, clusterName, exportID
}

func setExportJobFields(d *schema.ResourceData, exportJob *admin.DiskBackupExportJob) diag.Diagnostics {
	if err := d.Set("export_job_id", exportJob.GetId()); err != nil {
		return diag.Errorf("error setting `export_job_id` for snapshot export job (%s): %s", d.Id(), err)
	}

	if err := d.Set("snapshot_id", exportJob.GetSnapshotId()); err != nil {
		return diag.Errorf("error setting `snapshot_id` for snapshot export job (%s): %s", d.Id(), err)
	}

	if err := d.Set("custom_data", flattenExportJobsCustomData(exportJob.GetCustomData())); err != nil {
		return diag.Errorf("error setting `custom_data` for snapshot export job (%s): %s", d.Id(), err)
	}

	if err := d.Set("components", flattenExportJobsComponents(exportJob.GetComponents())); err != nil {
		return diag.Errorf("error setting `components` for snapshot export job (%s): %s", d.Id(), err)
	}

	if err := d.Set("created_at", conversion.TimePtrToStringPtr(exportJob.CreatedAt)); err != nil {
		return diag.Errorf("error setting `created_at` for snapshot export job (%s): %s", d.Id(), err)
	}

	if err := d.Set("export_bucket_id", exportJob.GetExportBucketId()); err != nil {
		return diag.Errorf("error setting `created_at` for snapshot export job (%s): %s", d.Id(), err)
	}

	if exportJob.ExportStatus != nil {
		if err := d.Set("export_status_exported_collections", exportJob.ExportStatus.GetExportedCollections()); err != nil {
			return diag.Errorf("error setting `export_status_exported_collections` for snapshot export job (%s): %s", d.Id(), err)
		}

		if err := d.Set("export_status_total_collections", exportJob.ExportStatus.GetTotalCollections()); err != nil {
			return diag.Errorf("error setting `export_status_total_collections` for snapshot export job (%s): %s", d.Id(), err)
		}
	}

	if err := d.Set("finished_at", conversion.TimePtrToStringPtr(exportJob.FinishedAt)); err != nil {
		return diag.Errorf("error setting `finished_at` for snapshot export job (%s): %s", d.Id(), err)
	}

	if err := d.Set("prefix", exportJob.GetPrefix()); err != nil {
		return diag.Errorf("error setting `prefix` for snapshot export job (%s): %s", d.Id(), err)
	}

	if err := d.Set("state", exportJob.GetState()); err != nil {
		return diag.Errorf("error setting `prefix` for snapshot export job (%s): %s", d.Id(), err)
	}

	return nil
}

func flattenExportJobsComponents(components []admin.DiskBackupExportMember) []map[string]any {
	if len(components) == 0 {
		return nil
	}

	customData := make([]map[string]any, 0)

	for i := range components {
		customData = append(customData, map[string]any{
			"export_id":        (components)[i].GetExportId(),
			"replica_set_name": (components)[i].GetReplicaSetName(),
		})
	}

	return customData
}

func flattenExportJobsCustomData(data []admin.BackupLabel) []map[string]any {
	if len(data) == 0 {
		return nil
	}

	customData := make([]map[string]any, 0)

	for i := range data {
		customData = append(customData, map[string]any{
			"key":   data[i].GetKey(),
			"value": data[i].GetValue(),
		})
	}

	return customData
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)
	clusterName := d.Get("cluster_name").(string)

	request := &admin.DiskBackupExportJobRequest{
		SnapshotId:     d.Get("snapshot_id").(string),
		ExportBucketId: d.Get("export_bucket_id").(string),
		CustomData:     expandExportJobCustomData(d),
	}

	jobResponse, _, err := connV2.CloudBackupsApi.CreateBackupExportJob(ctx, projectID, clusterName, request).Execute()
	if err != nil {
		return diag.Errorf("error creating snapshot export job: %s", err)
	}

	if err := d.Set("export_job_id", jobResponse.Id); err != nil {
		return diag.Errorf("error setting `export_job_id` for snapshot export job (%s): %s", *jobResponse.Id, err)
	}
	return resourceRead(ctx, d, meta)
}

func expandExportJobCustomData(d *schema.ResourceData) *[]admin.BackupLabel {
	customData := d.Get("custom_data").(*schema.Set)
	res := make([]admin.BackupLabel, customData.Len())

	for i, val := range customData.List() {
		v := val.(map[string]any)
		res[i] = admin.BackupLabel{
			Key:   conversion.Pointer(v["key"].(string)),
			Value: conversion.Pointer(v["value"].(string)),
		}
	}

	return &res
}

func resourceImportState(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	parts := strings.SplitN(d.Id(), "--", 3)
	if len(parts) != 3 {
		return nil, errors.New("import format error: to import snapshot export job, use the format {project_id}--{cluster_name}--{export_job_id}")
	}

	projectID := parts[0]
	clusterName := parts[1]
	exportID := parts[2]

	_, _, err := connV2.CloudBackupsApi.GetBackupExportJob(ctx, projectID, clusterName, exportID).Execute()
	if err != nil {
		return nil, fmt.Errorf("couldn't import snapshot export job %s in project %s and cluster %s, error: %s", exportID, projectID, clusterName, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		return nil, fmt.Errorf("error setting `project_id` for snapshot export job (%s): %s", d.Id(), err)
	}
	if err := d.Set("cluster_name", clusterName); err != nil {
		return nil, fmt.Errorf("error setting `cluster_name` for snapshot export job (%s): %s", d.Id(), err)
	}
	if err := d.Set("export_job_id", exportID); err != nil {
		return nil, fmt.Errorf("error setting `export_job_id` for snapshot export job (%s): %s", d.Id(), err)
	}
	return []*schema.ResourceData{d}, nil
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	d.SetId("")
	return nil
}
