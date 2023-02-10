package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type permCtxKey string

const (
	errorPrivateEndpointRegionalModeRead    = "error reading MongoDB Group `%s Private Endpoints Regional Mode: %s"
	errorPrivateEndpointRegionalModeSetting = "error setting `%s` on MongoDB Group `%s` Private Endpoints Regional Mode: %s"
	errorPrivateEndpointRegionalModeUpdate  = "error updating MongoDB Group `%s` Private Endpoints Regional Mode: %s"
)

var regionalModeTimeoutCtxKey permCtxKey = "regionalModeTimeout"

func resourceMongoDBAtlasPrivateEndpointRegionalMode() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBAtlasPrivateEndpointRegionalModeCreate,
		ReadContext:   resourceMongoDBAtlasPrivateEndpointRegionalModeRead,
		UpdateContext: resourceMongoDBAtlasPrivateEndpointRegionalModeUpdate,
		DeleteContext: resourceMongoDBAtlasPrivateEndpointRegionalModeDelete,
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
				Optional: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Hour),
			Update: schema.DefaultTimeout(3 * time.Hour),
			Delete: schema.DefaultTimeout(3 * time.Hour),
		},
	}
}

func resourceMongoDBAtlasPrivateEndpointRegionalModeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(d.Get("project_id").(string))
	err := resourceMongoDBAtlasPrivateEndpointRegionalModeUpdate(context.WithValue(ctx, regionalModeTimeoutCtxKey, schema.TimeoutCreate), d, meta)

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

		return diag.Errorf(errorPrivateEndpointRegionalModeRead, projectID, err)
	}

	if err := d.Set("enabled", setting.Enabled); err != nil {
		return diag.Errorf(errorPrivateEndpointRegionalModeSetting, "enabled", projectID, err)
	}

	return nil
}

func resourceMongoDBAtlasPrivateEndpointRegionalModeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	projectID := d.Id()
	enabled := d.Get("enabled").(bool)
	timeoutKey := ctx.Value(regionalModeTimeoutCtxKey)

	if timeoutKey == nil {
		timeoutKey = schema.TimeoutUpdate
	}

	_, resp, err := conn.PrivateEndpoints.UpdateRegionalizedPrivateEndpointSetting(ctx, projectID, enabled)
	if err != nil {
		if resp != nil && resp.Response.StatusCode == 404 {
			return nil
		}

		return diag.Errorf(errorPrivateEndpointRegionalModeUpdate, projectID, err)
	}

	log.Println("[INFO] Waiting for MongoDB Clusters' Private Endpoints to be updated")

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"REPEATING", "PENDING"},
		Target:     []string{"IDLE", "DELETED"},
		Refresh:    resourceClusterListAdvancedRefreshFunc(ctx, projectID, conn),
		Timeout:    d.Timeout(timeoutKey.(string)),
		MinTimeout: 5 * time.Second,
		Delay:      3 * time.Second,
	}
	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(errorPrivateEndpointRegionalModeUpdate, projectID, err)
	}

	return nil
}

func resourceMongoDBAtlasPrivateEndpointRegionalModeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if err := d.Set("enabled", false); err == nil {
		resourceMongoDBAtlasPrivateEndpointRegionalModeUpdate(context.WithValue(ctx, regionalModeTimeoutCtxKey, schema.TimeoutDelete), d, meta)
	} else {
		log.Printf(errorPrivateEndpointRegionalModeSetting, "enabled", d.Id(), err)
	}

	d.SetId("")

	return nil
}

func resourceMongoDBAtlasPrivateEndpointRegionalModeImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Id()

	setting, _, err := conn.PrivateEndpoints.GetRegionalizedPrivateEndpointSetting(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("couldn't import Private Endpoint Regional Mode for project %s error: %s", projectID, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		log.Printf(errorPrivateEndpointRegionalModeSetting, "project_id", projectID, err)
	}

	if err := d.Set("enabled", setting.Enabled); err != nil {
		log.Printf(errorPrivateEndpointRegionalModeSetting, "enabled", projectID, err)
	}

	d.SetId(projectID)

	return []*schema.ResourceData{d}, nil
}
