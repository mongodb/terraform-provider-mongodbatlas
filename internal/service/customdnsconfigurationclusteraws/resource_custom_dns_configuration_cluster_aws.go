package customdnsconfigurationclusteraws

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"

	"go.mongodb.org/atlas-sdk/v20250312006/admin"
)

const (
	errorCreate  = "error creating custom dns configuration cluster aws information: %s"
	errorRead    = "error getting custom dns configuration cluster aws information: %s"
	errorUpdate  = "error updating custom dns configuration cluster aws information: %s"
	errorDelete  = "error deleting custom dns configuration cluster aws (%s): %s"
	errorSetting = "error setting `%s` for custom dns configuration cluster aws (%s): %s"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate,
		ReadContext:   resourceRead,
		UpdateContext: resourceUpdate,
		DeleteContext: resourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)
	params := &admin.AWSCustomDNSEnabled{
		Enabled: d.Get("enabled").(bool),
	}
	_, _, err := connV2.AWSClustersDNSApi.ToggleAwsCustomDns(ctx, projectID, params).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorCreate, err))
	}
	d.SetId(projectID)
	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Id()
	dnsResp, resp, err := connV2.AWSClustersDNSApi.GetAwsCustomDns(context.Background(), projectID).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf(errorRead, err))
	}
	if err = d.Set("enabled", dnsResp.GetEnabled()); err != nil {
		return diag.FromErr(fmt.Errorf(errorSetting, "enabled", projectID, err))
	}
	if err = d.Set("project_id", d.Id()); err != nil {
		return diag.FromErr(fmt.Errorf(errorSetting, "project_id", projectID, err))
	}
	return nil
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Id()

	if d.HasChange("enabled") {
		params := &admin.AWSCustomDNSEnabled{
			Enabled: d.Get("enabled").(bool),
		}
		_, _, err := connV2.AWSClustersDNSApi.ToggleAwsCustomDns(ctx, projectID, params).Execute()
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorUpdate, err))
		}
	}

	return resourceRead(ctx, d, meta)
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Id()
	params := &admin.AWSCustomDNSEnabled{
		Enabled: false,
	}
	_, _, err := connV2.AWSClustersDNSApi.ToggleAwsCustomDns(ctx, projectID, params).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorDelete, projectID, err))
	}
	d.SetId("")
	return nil
}
