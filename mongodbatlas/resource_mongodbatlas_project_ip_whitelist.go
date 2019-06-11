package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	matlas "github.com/mongodb-partners/go-client-mongodb-atlas/mongodbatlas"
)

const (
	errorGetInfo = "error getting project IP whitelist information: %s"
)

func resourceMongoDBAtlasProjectIPWhitelist() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoDBAtlasProjectIPWhitelistCreate,
		Read:   resourceMongoDBAtlasProjectIPWhitelistRead,
		Update: resourceMongoDBAtlasProjectIPWhitelistUpdate,
		Delete: resourceMongoDBAtlasProjectIPWhitelistDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMongoDBAtlasProjectIPWhitelistImportState,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cidr_block": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"ip_address"},
			},
			"ip_address": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"comment": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}

func resourceMongoDBAtlasProjectIPWhitelistRead(d *schema.ResourceData, meta interface{}) error {
	//Get the client connection.
	conn := meta.(*matlas.Client)

	projectID := d.Get("project_id").(string)
	whitelistEntry := d.Id()

	resp, _, err := conn.ProjectIPWhitelist.Get(context.Background(), projectID, whitelistEntry)
	if err != nil {
		return fmt.Errorf(errorGetInfo, err)
	}

	if err := d.Set("cidr_block", resp.CIDRBlock); err != nil {
		return fmt.Errorf(errorGetInfo, err)
	}

	if err := d.Set("ip_address", resp.IPAddress); err != nil {
		return fmt.Errorf(errorGetInfo, err)
	}

	if err := d.Set("comment", resp.Comment); err != nil {
		return fmt.Errorf(errorGetInfo, err)
	}

	return nil
}

func resourceMongoDBAtlasProjectIPWhitelistCreate(d *schema.ResourceData, meta interface{}) error {
	//Get the client connection.
	conn := meta.(*matlas.Client)

	projectID := d.Get("project_id").(string)

	req := &matlas.ProjectIPWhitelist{}

	if v, ok := d.GetOk("cidr_block"); ok {
		req.CIDRBlock = v.(string)
	}

	if v, ok := d.GetOk("ip_address"); ok {
		req.IPAddress = v.(string)
	}

	if v, ok := d.GetOk("comment"); ok {
		req.Comment = v.(string)
	}

	resp, _, err := conn.ProjectIPWhitelist.Create(context.Background(), projectID, []*matlas.ProjectIPWhitelist{req})

	if err != nil {
		return fmt.Errorf("error creating project IP whitelist: %s", err)
	}

	//Get the project ip whitelist created.
	projectIPWhitelist := resp[0]

	d.SetId(projectIPWhitelist.CIDRBlock)

	if projectIPWhitelist.CIDRBlock == "" {
		d.SetId(projectIPWhitelist.IPAddress)
	}

	return resourceMongoDBAtlasProjectIPWhitelistRead(d, meta)
}

func resourceMongoDBAtlasProjectIPWhitelistUpdate(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)
	whitelistEntry := d.Id()

	projectIPWhitelist, _, err := conn.ProjectIPWhitelist.Get(context.Background(), projectID, whitelistEntry)

	if err != nil {
		return fmt.Errorf(errorGetInfo, err)
	}

	if d.HasChange("comment") {
		projectIPWhitelist.Comment = d.Get("comment").(string)
	}

	req := []*matlas.ProjectIPWhitelist{projectIPWhitelist}

	_, _, err = conn.ProjectIPWhitelist.Update(context.Background(), projectID, whitelistEntry, req)

	if err != nil {
		return fmt.Errorf("error updating project ip whitelist (%s): %s", whitelistEntry, err)
	}
	return resourceMongoDBAtlasProjectIPWhitelistRead(d, meta)
}

func resourceMongoDBAtlasProjectIPWhitelistDelete(d *schema.ResourceData, meta interface{}) error {
	//Get the client connection.
	conn := meta.(*matlas.Client)

	projectID := d.Get("project_id").(string)
	whitelistEntry := d.Id()

	_, err := conn.ProjectIPWhitelist.Delete(context.Background(), projectID, whitelistEntry)
	if err != nil {
		return fmt.Errorf("error deleting project IP whitelist: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceMongoDBAtlasProjectIPWhitelistImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
