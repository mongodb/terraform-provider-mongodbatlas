package customdnsconfigurationclusteraws

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"

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
	conn := meta.(*config.MongoDBClient).Atlas
	orgID := d.Get("project_id").(string)

	// Creating(Updating) the Custom DNS Configuration for Atlas Clusters on AWS
	_, _, err := conn.CustomAWSDNS.Update(ctx, orgID,
		&matlas.AWSCustomDNSSetting{
			Enabled: d.Get("enabled").(bool),
		})
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorCreate, err))
	}

	d.SetId(orgID)

	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).Atlas

	dnsResp, resp, err := conn.CustomAWSDNS.Get(context.Background(), d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorRead, err))
	}

	if err = d.Set("enabled", dnsResp.Enabled); err != nil {
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
