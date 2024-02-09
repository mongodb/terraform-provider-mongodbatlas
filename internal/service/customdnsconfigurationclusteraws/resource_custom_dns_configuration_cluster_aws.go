package customdnsconfigurationclusteraws

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"

	"go.mongodb.org/atlas-sdk/v20231115006/admin"
	matlas "go.mongodb.org/atlas/mongodbatlas"
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
	_, _, err := connV2.AWSClustersDNSApi.ToggleAWSCustomDNS(ctx, projectID, params).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorCreate, err))
	}
	d.SetId(projectID)
	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	dnsResp, resp, err := connV2.AWSClustersDNSApi.GetAWSCustomDNS(context.Background(), d.Id()).Execute()
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf(errorRead, err))
	}
	if err = d.Set("enabled", dnsResp.GetEnabled()); err != nil {
		return diag.FromErr(fmt.Errorf(errorSetting, "enabled", d.Id(), err))
	}
	if err = d.Set("project_id", d.Id()); err != nil {
		return diag.FromErr(fmt.Errorf(errorSetting, "project_id", d.Id(), err))
	}
	return nil
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).Atlas

	if d.HasChange("enabled") {
		_, _, err := conn.CustomAWSDNS.Update(ctx, d.Id(), &matlas.AWSCustomDNSSetting{
			Enabled: d.Get("enabled").(bool),
		})
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorUpdate, err))
		}
	}

	return resourceRead(ctx, d, meta)
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).Atlas

	_, _, err := conn.CustomAWSDNS.Update(ctx, d.Id(), &matlas.AWSCustomDNSSetting{
		Enabled: false,
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorDelete, d.Id(), err))
	}

	d.SetId("")

	return nil
}
