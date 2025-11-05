package privateendpointregionalmode

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"go.mongodb.org/atlas-sdk/v20250312009/admin"
)

type permCtxKey string

const (
	errorPrivateEndpointRegionalModeRead    = "error reading MongoDB Group `%s Private Endpoints Regional Mode: %s"
	errorPrivateEndpointRegionalModeSetting = "error setting `%s` on MongoDB Group `%s` Private Endpoints Regional Mode: %s"
	errorPrivateEndpointRegionalModeUpdate  = "error updating MongoDB Group `%s` Private Endpoints Regional Mode: %s"
)

var regionalModeTimeoutCtxKey permCtxKey = "regionalModeTimeout"

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate,
		ReadContext:   resourceRead,
		UpdateContext: resourceUpdate,
		DeleteContext: resourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceImportState,
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
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Hour),
			Update: schema.DefaultTimeout(3 * time.Hour),
			Delete: schema.DefaultTimeout(3 * time.Hour),
		},
	}
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	d.SetId(d.Get("project_id").(string))
	err := resourceUpdate(context.WithValue(ctx, regionalModeTimeoutCtxKey, schema.TimeoutCreate), d, meta)

	if err != nil {
		return err
	}

	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2

	projectID := d.Id()

	setting, resp, err := conn.PrivateEndpointServicesApi.GetRegionalEndpointMode(ctx, projectID).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) {
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

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2

	projectID := d.Id()
	enabled := d.Get("enabled").(bool)
	timeoutKey := ctx.Value(regionalModeTimeoutCtxKey)

	if timeoutKey == nil {
		timeoutKey = schema.TimeoutUpdate
	}
	settingParam := admin.ProjectSettingItem{
		Enabled: enabled,
	}
	_, resp, err := conn.PrivateEndpointServicesApi.ToggleRegionalEndpointMode(ctx, projectID, &settingParam).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) {
			return nil
		}

		return diag.Errorf(errorPrivateEndpointRegionalModeUpdate, projectID, err)
	}

	log.Println("[INFO] Waiting for MongoDB Clusters' Private Endpoints to be updated")

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"REPEATING", "PENDING"},
		Target:     []string{"IDLE", "DELETED"},
		Refresh:    advancedcluster.ResourceClusterListAdvancedRefreshFunc(ctx, projectID, conn.ClustersApi),
		Timeout:    d.Timeout(timeoutKey.(string)),
		MinTimeout: 15 * time.Second,
		Delay:      30 * time.Second, // give time for cluster connection strings to be updated
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(errorPrivateEndpointRegionalModeUpdate, projectID, err)
	}

	return nil
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	if err := d.Set("enabled", false); err == nil {
		resourceUpdate(context.WithValue(ctx, regionalModeTimeoutCtxKey, schema.TimeoutDelete), d, meta)
	} else {
		log.Printf(errorPrivateEndpointRegionalModeSetting, "enabled", d.Id(), err)
	}

	d.SetId("")

	return nil
}

func resourceImportState(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	conn := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Id()

	setting, _, err := conn.PrivateEndpointServicesApi.GetRegionalEndpointMode(ctx, projectID).Execute()
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
