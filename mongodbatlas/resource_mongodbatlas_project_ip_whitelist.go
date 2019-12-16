package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

const (
	errorWhitelistCreate  = "error creating Project IP Whitelist information: %s"
	errorWhitelistRead    = "error getting Project IP Whitelist information: %s"
	errorWhitelistUpdate  = "error updating Project IP Whitelist information: %s"
	errorWhitelistDelete  = "error deleting Project IP Whitelist information: %s"
	errorWhitelistSetting = "error setting `%s` for Project IP Whitelist (%s): %s"
)

func resourceMongoDBAtlasProjectIPWhitelist() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoDBAtlasProjectIPWhitelistCreate,
		Update: resourceMongoDBAtlasProjectIPWhitelistUpdate,
		Read:   resourceMongoDBAtlasProjectIPWhitelistRead,
		Delete: resourceMongoDBAtlasProjectIPWhitelistDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cidr_block": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: func(i interface{}, k string) (s []string, es []error) {
					v, ok := i.(string)
					if !ok {
						es = append(es, fmt.Errorf("expected type of %s to be string", k))
						return
					}

					_, ipnet, err := net.ParseCIDR(v)
					if err != nil {
						es = append(es, fmt.Errorf("expected %s to contain a valid CIDR, got: %s with err: %s", k, v, err))
						return
					}

					if ipnet == nil || v != ipnet.String() {
						es = append(es, fmt.Errorf("expected %s to contain a valid network CIDR, expected %s, got %s", k, ipnet, v))
						return
					}
					return
				},
			},
			"ip_address": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
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
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)

	if d.Get("cidr_block") == "" && d.Get("ip_address") == "" {
		return errors.New("cidr_block or ip_address needs to be setted")
	}

	_, _, err := conn.ProjectIPWhitelist.Create(context.Background(), projectID, []*matlas.ProjectIPWhitelist{
		{
			CIDRBlock: d.Get("cidr_block").(string),
			IPAddress: d.Get("ip_address").(string),
			Comment:   d.Get("comment").(string),
		},
	})
	if err != nil {
		return fmt.Errorf(errorWhitelistCreate, err)
	}

	entry := d.Get("ip_address").(string)
	if entry == "" {
		entry = d.Get("cidr_block").(string)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"entry":      entry,
	}))

	return resourceMongoDBAtlasProjectIPWhitelistRead(d, meta)
}

func resourceMongoDBAtlasProjectIPWhitelistRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())

	whitelist, _, err := conn.ProjectIPWhitelist.Get(context.Background(), ids["project_id"], ids["entry"])
	if err != nil {
		return fmt.Errorf(errorWhitelistRead, err)
	}

	if err := d.Set("ip_address", whitelist.IPAddress); err != nil {
		return fmt.Errorf(errorWhitelistSetting, "ip_address", ids["project_id"], err)
	}
	if err := d.Set("cidr_block", whitelist.CIDRBlock); err != nil {
		return fmt.Errorf(errorWhitelistSetting, "cidr_block", ids["project_id"], err)
	}
	if err := d.Set("comment", whitelist.Comment); err != nil {
		return fmt.Errorf(errorWhitelistSetting, "comment", ids["project_id"], err)
	}

	return nil
}

func resourceMongoDBAtlasProjectIPWhitelistUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())

	req := &matlas.ProjectIPWhitelist{}

	if d.HasChange("comment") {
		req.Comment = d.Get("comment").(string)
	}

	_, _, err := conn.ProjectIPWhitelist.Update(context.Background(), ids["project_id"], []*matlas.ProjectIPWhitelist{req})
	if err != nil {
		return fmt.Errorf(errorWhitelistUpdate, err)
	}

	return resourceMongoDBAtlasProjectIPWhitelistRead(d, meta)
}

func resourceMongoDBAtlasProjectIPWhitelistDelete(d *schema.ResourceData, meta interface{}) error {
	//Get the client connection.
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())

	_, err := conn.ProjectIPWhitelist.Delete(context.Background(), ids["project_id"], ids["entry"])
	if err != nil {
		return fmt.Errorf(errorWhitelistDelete, err)
	}
	return nil
}
