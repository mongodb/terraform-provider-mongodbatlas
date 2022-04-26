package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorPrivateEndpointRegionalModeRead    = "error reading MongoDB Group `%s Private Endpoints Regional Mode: %s"
	errorPrivateEndpointRegionalModeSetting = "error setting `%s` on MongoDB Group `%s` Private Endpoints Regional Mode: %s"
)

func resourceMongoDBAtlasPrivateEndpointRegionalMode() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBAtlasPrivateEndpointRegionalModeCreate,
		ReadContext:   resourceMongoDBAtlasPrivateEndpointRegionalModeRead,
		UpdateContext: resourceMongoDBAtlasPrivateEndpointRegionalModeUpdate,
		DeleteContext: schema.NoopContext,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasPrivateEndpointRegionalModeImportState,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceMongoDBAtlasPrivateEndpointRegionalModeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(d.Get("project_id").(string))
	err := resourceMongoDBAtlasPrivateEndpointRegionalModeUpdate(ctx, d, meta)

	if err != nil {
		return err
	}

	return resourceMongoDBAtlasPrivateEndpointRegionalModeRead(ctx, d, meta)
}

func resourceMongoDBAtlasPrivateEndpointRegionalModeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	projectID := d.Id()

	setting, resp, err := conn.PrivateEndpoints.GetRegionalizedPrivateEndpointSetting(context.Background(), projectID)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorPrivateEndpointRegionalModeRead, projectID, err))
	}

	if err := d.Set("enabled", setting.Enabled); err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateEndpointRegionalModeSetting, "enabled", projectID, err))
	}

	return nil
}

func resourceMongoDBAtlasPrivateEndpointRegionalModeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	projectID := d.Id()
	enabled := d.Get("enabled").(bool)

	_, resp, err := conn.PrivateEndpoints.UpdateRegionalizedPrivateEndpointSetting(ctx, projectID, enabled)
	if err != nil {
		if resp != nil && resp.Response.StatusCode == 404 {
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsDelete, projectID, err))
	}

	log.Println("[INFO] Waiting for MongoDB Private Endpoints Connection to be destroyed")

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING", "REPEATING"},
		Target:     []string{"APPLIED"},
		Refresh:    resourcePrivateEndpointRegionalModeRefreshFunc(ctx, conn, projectID),
		Timeout:    1 * time.Hour,
		MinTimeout: 5 * time.Second,
		Delay:      3 * time.Second,
	}
	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsDelete, projectID, err))
	}

	return nil
}

func resourceMongoDBAtlasPrivateEndpointRegionalModeImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Id()

	setting, _, err := conn.PrivateEndpoints.GetRegionalizedPrivateEndpointSetting(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("couldn't import regional mode for project %s error: %s", projectID, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		log.Printf(errorPrivateLinkEndpointsSetting, "project_id", projectID, err)
	}

	if err := d.Set("enabled", setting.Enabled); err != nil {
		log.Printf(errorPrivateLinkEndpointsSetting, "enabled", projectID, err)
	}

	d.SetId(projectID)

	return []*schema.ResourceData{d}, nil
}

func resourcePrivateEndpointRegionalModeRefreshFunc(ctx context.Context, client *matlas.Client, projectID string) resource.StateRefreshFunc {
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
			s, resp, err := client.Clusters.Status(ctx, projectID, clusters[i].Name)

			if err != nil && strings.Contains(err.Error(), "reset by peer") {
				return nil, "REPEATING", nil
			}

			if err != nil {
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

		// If all clusters were properly read, and none are PENDING, all changes have been APPLIED.
		return clusters, "APPLIED", nil
	}
}
