package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
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
				ValidateFunc: func(i interface{}, k string) (s []string, es []error) {
					v, ok := i.(string)
					if !ok {
						es = append(es, fmt.Errorf("expected type of %s to be string", k))
						return
					}

					_, ipnet, err := net.ParseCIDR(v)
					if err != nil {
						es = append(es, fmt.Errorf(
							"expected %s to contain a valid CIDR, got: %s with err: %s", k, v, err))
						return
					}

					if ipnet == nil || v != ipnet.String() {
						es = append(es, fmt.Errorf(
							"expected %s to contain a valid network CIDR, expected %s, got %s",
							k, ipnet, v))
						return
					}
					return
				},
			},
			"ip_address": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.SingleIP(),
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

	for _, entry := range resp {
		if (req.CIDRBlock != "" && entry.CIDRBlock == req.CIDRBlock) ||
			(req.IPAddress != "" && entry.IPAddress == req.IPAddress) {
			d.SetId(encodeStateID(map[string]string{
				"project_id": projectID,
				"cidr_block": entry.CIDRBlock,
			}))
			return resourceMongoDBAtlasProjectIPWhitelistRead(d, meta)
		}
	}
	return fmt.Errorf("MongoDB Project IP Whitelist with CIDR block: %s and IP Address: %s could not be found in the response from MongoDB Atlas", req.CIDRBlock, req.IPAddress)
}

func resourceMongoDBAtlasProjectIPWhitelistRead(d *schema.ResourceData, meta interface{}) error {
	//Get the client connection.
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	whitelistEntry := ids["cidr_block"]

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

func resourceMongoDBAtlasProjectIPWhitelistUpdate(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	whitelistEntry := ids["cidr_block"]

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
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	whitelistEntry := ids["cidr_block"]

	_, err := conn.ProjectIPWhitelist.Delete(context.Background(), projectID, whitelistEntry)
	if err != nil {
		return fmt.Errorf("error deleting project IP whitelist: %s", err)
	}
	return nil
}

func resourceMongoDBAtlasProjectIPWhitelistImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*matlas.Client)

	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 2 {
		return nil, errors.New("To import an ip whitelist, use the format {project_id}-{cidr block}")
	}
	projectID := parts[0]
	whitelistEntry := parts[1]

	log.Printf("[DEBUG] whitelist entry: %s", whitelistEntry)

	ipEntry, _, err := conn.ProjectIPWhitelist.Get(context.Background(), projectID, whitelistEntry)
	if err != nil {
		return nil, fmt.Errorf("Couldn't import ip whitelist %s in project_id %s, error: %s", whitelistEntry, projectID, err.Error())
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": ipEntry.GroupID,
		"cidr_block": ipEntry.CIDRBlock,
	}))

	if err := d.Set("project_id", ipEntry.GroupID); err != nil {
		log.Printf("[WARN] Error setting project_id for (%s): %s", ipEntry.CIDRBlock, err)
		return []*schema.ResourceData{d}, err
	}
	return []*schema.ResourceData{d}, nil
}
