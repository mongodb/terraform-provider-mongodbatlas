package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mwielbut/pointy"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorPrivateIPModeCreate = "error setting MongoDB Only Private IP Mode for Peering Connections: %s"
	errorPrivateIPModeRead   = "error reading MongoDB Only Private IP Mode for Peering Connections (%s): %s"
	errorPrivateIPModeDelete = "error deleting MongoDB Only Private IP Mode for Peering Connections (%s): %s"
)

func resourceMongoDBAtlasPrivateIPMode() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBAtlasPrivateIPModeCreate,
		ReadContext:   resourceMongoDBAtlasPrivateIPModeRead,
		UpdateContext: resourceMongoDBAtlasPrivateIPModeCreate,
		DeleteContext: resourceMongoDBAtlasPrivateIPModeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasPrivateIPModeImportState,
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

func resourceMongoDBAtlasPrivateIPModeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)

	// Get the required ones
	privateIPModeRequest := &matlas.PrivateIPMode{
		Enabled: pointy.Bool(d.Get("enabled").(bool)),
	}

	_, _, err := conn.PrivateIPMode.Update(ctx, projectID, privateIPModeRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateIPModeCreate, err))
	}

	d.SetId(projectID)

	return resourceMongoDBAtlasPrivateIPModeRead(ctx, d, meta)
}

func resourceMongoDBAtlasPrivateIPModeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas

	projectID := d.Id()

	privateIPMode, resp, err := conn.PrivateIPMode.Get(ctx, projectID)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorPrivateIPModeRead, projectID, err))
	}

	if err := d.Set("enabled", privateIPMode.Enabled); err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateIPModeRead, projectID, err))
	}

	return nil
}

func resourceMongoDBAtlasPrivateIPModeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Id()

	// Get the required ones
	privateIPModeRequest := &matlas.PrivateIPMode{
		Enabled: pointy.Bool(false),
	}

	_, _, err := conn.PrivateIPMode.Update(ctx, projectID, privateIPModeRequest)

	if err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateIPModeDelete, projectID, err))
	}

	return nil
}

func resourceMongoDBAtlasPrivateIPModeImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := d.Set("project_id", d.Id()); err != nil {
		log.Printf("[WARN] Error setting project_id for private IP Mode: %s", err)
	}

	return []*schema.ResourceData{d}, nil
}
