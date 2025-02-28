package cloudbackupsnapshotexportbucket

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/atlas-sdk/v20250219001/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate,
		ReadContext:   resourceRead,
		DeleteContext: resourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceImportState,
		},
		Schema: Schema(),
	}
}

func Schema() map[string]*schema.Schema {
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
			Optional: true,
			ForceNew: true,
		},
		"role_id": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"service_url": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"tenant_id": {
			Type:     schema.TypeString,
			Optional: true, // attribute is only used as a computed, this is called out in docs and configuration of optional argument can be eventually removed implying a breaking change. To be removed in https://jira.mongodb.org/browse/CLOUDP-293142
			Computed: true,
		},
	}
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)

	cloudProvider := d.Get("cloud_provider").(string)
	request := &admin.DiskBackupSnapshotExportBucketRequest{
		IamRoleId:     conversion.StringPtr(d.Get("iam_role_id").(string)),
		BucketName:    d.Get("bucket_name").(string),
		RoleId:        conversion.StringPtr(d.Get("role_id").(string)),
		ServiceUrl:    conversion.StringPtr(d.Get("service_url").(string)),
		TenantId:      conversion.StringPtr(d.Get("tenant_id").(string)),
		CloudProvider: cloudProvider,
	}
	bucketResponse, _, err := conn.CloudBackupsApi.CreateExportBucket(ctx, projectID, request).Execute()
	if err != nil {
		return diag.Errorf("error creating snapshot export bucket: %s", err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id": projectID,
		"id":         bucketResponse.GetId(),
	}))

	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	bucketID := ids["id"]

	exportBackup, _, err := conn.CloudBackupsApi.GetExportBucket(ctx, projectID, bucketID).Execute()
	if err != nil {
		reset := strings.Contains(err.Error(), "404") && !d.IsNewResource()

		if reset {
			d.SetId("")
			return nil
		}

		return diag.Errorf("error getting snapshot export backup information: %s", err)
	}

	if err := d.Set("export_bucket_id", exportBackup.GetId()); err != nil {
		return diag.Errorf("error setting `export_bucket_id` for snapshot export bucket (%s): %s", d.Id(), err)
	}

	if err := d.Set("bucket_name", exportBackup.GetBucketName()); err != nil {
		return diag.Errorf("error setting `bucket_name` for snapshot export bucket (%s): %s", d.Id(), err)
	}

	if err := d.Set("cloud_provider", exportBackup.GetCloudProvider()); err != nil {
		return diag.Errorf("error setting `bucket_name` for snapshot export bucket (%s): %s", d.Id(), err)
	}

	if err := d.Set("iam_role_id", exportBackup.IamRoleId); err != nil {
		return diag.Errorf("error setting `iam_role_id` for snapshot export bucket (%s): %s", d.Id(), err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		return diag.Errorf("error setting `project_id` for snapshot export bucket (%s): %s", d.Id(), err)
	}

	if err := d.Set("service_url", exportBackup.ServiceUrl); err != nil {
		return diag.Errorf("error setting `service_url` for snapshot export bucket (%s): %s", d.Id(), err)
	}

	if err := d.Set("role_id", exportBackup.RoleId); err != nil {
		return diag.Errorf("error setting `role_id` for snapshot export bucket (%s): %s", d.Id(), err)
	}

	if err := d.Set("tenant_id", exportBackup.TenantId); err != nil {
		return diag.Errorf("error setting `tenant_id` for snapshot export bucket (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	bucketID := ids["id"]

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"PENDING", "REPEATING"},
		Target:     []string{"DELETED"},
		Refresh:    resourceRefresh(ctx, conn, projectID, bucketID),
		Timeout:    1 * time.Hour,
		MinTimeout: 5 * time.Second,
		Delay:      3 * time.Second,
	}
	_, _, err := conn.CloudBackupsApi.DeleteExportBucket(ctx, projectID, bucketID).Execute()

	if err != nil {
		return diag.Errorf("error deleting snapshot export bucket (%s): %s", bucketID, err)
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error deleting snapshot export bucket %s %s", projectID, err)
	}

	return nil
}

func resourceImportState(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	conn := meta.(*config.MongoDBClient).AtlasV2

	projectID, id, err := splitImportID(d.Id())
	if err != nil {
		return nil, err
	}

	_, _, err = conn.CloudBackupsApi.GetExportBucket(ctx, *projectID, *id).Execute()
	if err != nil {
		return nil, fmt.Errorf("couldn't import snapshot export bucket %s in project %s, error: %s", *id, *projectID, err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id": *projectID,
		"id":         *id,
	}))

	return []*schema.ResourceData{d}, nil
}

func splitImportID(id string) (projectID, bucketID *string, err error) {
	var re = regexp.MustCompile(`(?s)^([0-9a-fA-F]{24})-(.*)$`)
	parts := re.FindStringSubmatch(id)

	if len(parts) != 3 {
		err = errors.New("import format error: to import a serverless instance, use the format {project_id}-{id}")
		return
	}

	projectID = &parts[1]
	bucketID = &parts[2]

	return
}

func resourceRefresh(ctx context.Context, client *admin.APIClient, projectID, exportBucketID string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		clustersPaginated, resp, err := client.ClustersApi.ListClusters(ctx, projectID).Execute()
		if err != nil {
			// For our purposes, no clusters is equivalent to all changes having been APPLIED
			if validate.StatusNotFound(resp) {
				return "", "APPLIED", nil
			}
			return nil, "REPEATING", err
		}
		clusters := clustersPaginated.GetResults()

		for i := range clusters {
			backupPolicy, _, err := client.CloudBackupsApi.GetBackupSchedule(context.Background(), projectID, clusters[i].GetName()).Execute()
			if err != nil {
				continue
			}
			// find cluster that has export id attached to its config
			if backupPolicy.Export != nil {
				if backupPolicy.Export.GetExportBucketId() == exportBucketID {
					if clusters[i].GetStateName() == "IDLE" {
						return clusters, "PENDING", nil
					}
					if clusters[i].GetStateName() == "UPDATING" {
						return clusters, "PENDING", nil
					}

					s, resp, err := client.ClustersApi.GetClusterStatus(ctx, projectID, clusters[i].GetName()).Execute()

					if err != nil && strings.Contains(err.Error(), "reset by peer") {
						return nil, "REPEATING", nil
					}

					if err != nil {
						if validate.StatusNotFound(resp) {
							return "", "DELETED", nil
						}

						if validate.StatusNotFound(resp) {
							// The cluster no longer exists, consider this equivalent to status APPLIED
							continue
						}
						if validate.StatusServiceUnavailable(resp) {
							return "", "PENDING", nil
						}
						return nil, "REPEATING", err
					}

					if s.GetChangeStatus() == "PENDING" {
						return clusters, "PENDING", nil
					}
				}
			}
		}

		// If all clusters were properly read, and none are PENDING, all changes have been APPLIED.
		return clusters, "DELETED", nil
	}
}
