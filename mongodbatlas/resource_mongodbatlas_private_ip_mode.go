package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

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
		Create: resourceMongoDBAtlasPrivateIPModeCreate,
		Read:   resourceMongoDBAtlasPrivateIPModeRead,
		Update: resourceMongoDBAtlasPrivateIPModeCreate,
		Delete: resourceMongoDBAtlasPrivateIPModeDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMongoDBAtlasPrivateIPModeImportState,
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

func resourceMongoDBAtlasPrivateIPModeCreate(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)

	// Get the required ones
	privateIPModeRequest := &matlas.PrivateIPMode{
		Enabled: pointy.Bool(d.Get("enabled").(bool)),
	}

	_, _, err := conn.PrivateIPMode.Update(context.Background(), projectID, privateIPModeRequest)
	if err != nil {
		return fmt.Errorf(errorPrivateIPModeCreate, err)
	}

	d.SetId(projectID)

	return resourceMongoDBAtlasPrivateIPModeRead(d, meta)
}

func resourceMongoDBAtlasPrivateIPModeRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*matlas.Client)

	projectID := d.Id()

	privateIPMode, resp, err := conn.PrivateIPMode.Get(context.Background(), projectID)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil
		}

		return fmt.Errorf(errorPrivateIPModeRead, projectID, err)
	}

	if err := d.Set("enabled", privateIPMode.Enabled); err != nil {
		return fmt.Errorf(errorPrivateIPModeRead, projectID, err)
	}

	return nil
}

func resourceMongoDBAtlasPrivateIPModeDelete(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*matlas.Client)
	projectID := d.Id()

	// Get the required ones
	privateIPModeRequest := &matlas.PrivateIPMode{
		Enabled: pointy.Bool(false),
	}

	_, _, err := conn.PrivateIPMode.Update(context.Background(), projectID, privateIPModeRequest)

	if err != nil {
		return fmt.Errorf(errorPrivateIPModeDelete, projectID, err)
	}

	return nil
}

func resourceMongoDBAtlasPrivateIPModeImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := d.Set("project_id", d.Id()); err != nil {
		log.Printf("[WARN] Error setting project_id for private IP Mode: %s", err)
	}

	return []*schema.ResourceData{d}, nil
}
