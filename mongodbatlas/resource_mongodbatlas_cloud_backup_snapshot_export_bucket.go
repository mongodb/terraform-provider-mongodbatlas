package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func resourceMongoDBAtlasCloudBackupSnapshotExportBucket() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBAtlasCloudBackupSnapshotExportBucketCreate,
		ReadContext:   resourceMongoDBAtlasCloudBackupSnapshotExportBucketRead,
		DeleteContext: resourceMongoDBAtlasCloudBackupSnapshotExportBucketDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasCloudBackupSnapshotExportBucketImportState,
		},
		Schema: returnCloudBackupSnapshotExportBucketSchema(),
	}
}

func returnCloudBackupSnapshotExportBucketSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"export_bucket_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"project_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"bucket_name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"cloud_provider": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"iam_role_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
	}
}

func resourceMongoDBAtlasCloudBackupSnapshotExportBucketCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)

	cloudProvider := d.Get("cloud_provider").(string)
	if cloudProvider != "AWS" {
		return diag.Errorf("atlas only supports AWS")
	}

	request := &matlas.CloudProviderSnapshotExportBucket{
		IAMRoleID:     d.Get("iam_role_id").(string),
		BucketName:    d.Get("bucket_name").(string),
		CloudProvider: cloudProvider,
	}

	bucketResponse, _, err := conn.CloudProviderSnapshotExportBuckets.Create(ctx, projectID, request)
	if err != nil {
		return diag.Errorf("error creating snapshot export bucket: %s", err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"id":         bucketResponse.ID,
	}))

	return resourceMongoDBAtlasCloudBackupSnapshotExportBucketRead(ctx, d, meta)
}

func resourceMongoDBAtlasCloudBackupSnapshotExportBucketRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	exportJobID := ids["id"]

	exportBackup, _, err := conn.CloudProviderSnapshotExportBuckets.Get(ctx, projectID, exportJobID)
	if err != nil {
		// case 404
		// deleted in the backend case
		reset := strings.Contains(err.Error(), "404") && !d.IsNewResource()

		if reset {
			d.SetId("")
			return nil
		}

		return diag.Errorf("error getting snapshot export backup information: %s", err)
	}

	if err := d.Set("export_bucket_id", exportBackup.ID); err != nil {
		return diag.Errorf("error setting `export_bucket_id` for snapshot export bucket (%s): %s", d.Id(), err)
	}

	if err := d.Set("bucket_name", exportBackup.BucketName); err != nil {
		return diag.Errorf("error setting `bucket_name` for snapshot export bucket (%s): %s", d.Id(), err)
	}

	if err := d.Set("cloud_provider", exportBackup.CloudProvider); err != nil {
		return diag.Errorf("error setting `bucket_name` for snapshot export bucket (%s): %s", d.Id(), err)
	}

	if err := d.Set("iam_role_id", exportBackup.IAMRoleID); err != nil {
		return diag.Errorf("error setting `iam_role_id` for snapshot export bucket (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceMongoDBAtlasCloudBackupSnapshotExportBucketDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	exportJobID := ids["id"]

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING", "REPEATING"},
		Target:     []string{"DELETED"},
		Refresh:    resourceCloudBackupSnapshotExportBucketRefreshFunc(ctx, conn, projectID, exportJobID),
		Timeout:    1 * time.Hour,
		MinTimeout: 5 * time.Second,
		Delay:      3 * time.Second,
	}
	// Wait, catching any errors
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error deleting snapshot export bucket %s %s", projectID, err)
	}

	_, err = conn.CloudProviderSnapshotExportBuckets.Delete(ctx, projectID, exportJobID)

	if err != nil {
		return diag.Errorf("error deleting snapshot export bucket (%s): %s", exportJobID, err)
	}

	return nil
}

func resourceMongoDBAtlasCloudBackupSnapshotExportBucketImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*MongoDBClient).Atlas

	projectID, id, err := splitCloudBackupSnapshotExportBucketImportID(d.Id())
	if err != nil {
		return nil, err
	}

	_, _, err = conn.CloudProviderSnapshotExportBuckets.Get(ctx, *projectID, *id)
	if err != nil {
		return nil, fmt.Errorf("couldn't import snapshot export bucket %s in project %s, error: %s", *id, *projectID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": *projectID,
		"id":         *id,
	}))

	return []*schema.ResourceData{d}, nil
}

func splitCloudBackupSnapshotExportBucketImportID(id string) (projectID, exportJobID *string, err error) {
	var re = regexp.MustCompile(`(?s)^([0-9a-fA-F]{24})-(.*)$`)
	parts := re.FindStringSubmatch(id)

	if len(parts) != 3 {
		err = errors.New("import format error: to import a serverless instance, use the format {project_id}-{id}")
		return
	}

	projectID = &parts[1]
	exportJobID = &parts[2]

	return
}

func resourceCloudBackupSnapshotExportBucketRefreshFunc(ctx context.Context, client *matlas.Client, projectID, exportBucketID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		clusters, resp, err := client.Clusters.List(ctx, projectID, nil)
		if err != nil {
			// For our purposes, no clusters is equivalent to all changes having been APPLIED
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				return "", "APPLIED", nil
			}
			return nil, "REPEATING", err
		}

		for i := range clusters {
			backupPolicy, _, err := client.CloudProviderSnapshotBackupPolicies.Get(context.Background(), projectID, clusters[i].Name)
			if err != nil {
				continue
			}
			// find cluster that has export id attached to its config
			if backupPolicy.Export != nil {
				if backupPolicy.Export.ExportBucketID == exportBucketID {
					if clusters[i].StateName == "IDLE" {
						return clusters, "PENDING", nil
					}
					if clusters[i].StateName == "UPDATING" {
						return clusters, "PENDING", nil
					}

					s, resp, err := client.Clusters.Status(ctx, projectID, clusters[i].Name)

					if err != nil && strings.Contains(err.Error(), "reset by peer") {
						return nil, "REPEATING", nil
					}

					if err != nil {
						if resp != nil && resp.StatusCode == http.StatusNotFound {
							return "", "DELETED", nil
						}

						if resp.StatusCode == 404 {
							// The cluster no longer exists, consider this equivalent to status APPLIED
							continue
						}
						if resp.StatusCode == 503 {
							return "", "PENDING", nil
						}
						return nil, "REPEATING", err
					}

					if s.ChangeStatus == matlas.ChangeStatusPending {
						return clusters, "PENDING", nil
					}
				}
			}
		}

		// If all clusters were properly read, and none are PENDING, all changes have been APPLIED.
		return clusters, "DELETED", nil
	}
}
