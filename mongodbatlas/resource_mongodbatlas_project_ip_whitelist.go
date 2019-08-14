package mongodbatlas

import (
	"context"
	"fmt"
	"net"

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
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"entry": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr_block": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
							// ConflictsWith: []string{"ip_address"},
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
				},
			},
			"whitelist": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr_block": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"comment": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceMongoDBAtlasProjectIPWhitelistCreate(d *schema.ResourceData, meta interface{}) error {
	//Get the client connection.
	conn := meta.(*matlas.Client)

	projectID := d.Get("project_id").(string)

	req := expandProjectIPWhitelist(d)
	_, _, err := conn.ProjectIPWhitelist.Create(context.Background(), projectID, req)

	if err != nil {
		return fmt.Errorf("error creating project IP whitelist: %s", err)
	}

	d.SetId(projectID)
	return resourceMongoDBAtlasProjectIPWhitelistRead(d, meta)
}

func resourceMongoDBAtlasProjectIPWhitelistRead(d *schema.ResourceData, meta interface{}) error {
	//Get the client connection.
	conn := meta.(*matlas.Client)
	projectID := d.Id()

	var options *matlas.ListOptions
	resp, _, err := conn.ProjectIPWhitelist.List(context.Background(), projectID, options)
	if err != nil {
		return fmt.Errorf(errorGetInfo, err)
	}

	if err := d.Set("whitelist", flattenProjectIPWhitelist(resp)); err != nil {
		return fmt.Errorf(errorGetInfo, err)
	}
	return nil
}

func resourceMongoDBAtlasProjectIPWhitelistUpdate(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)

	if d.HasChange("whitelist") {
		whitelist := expandProjectIPWhitelist(d)

		_, _, err := conn.ProjectIPWhitelist.Update(context.Background(), d.Id(), "", whitelist)
		if err != nil {
			return fmt.Errorf("error updating project ip whitelist (%s): %s", d.Id(), err)
		}
	}

	return resourceMongoDBAtlasProjectIPWhitelistRead(d, meta)
}

func resourceMongoDBAtlasProjectIPWhitelistDelete(d *schema.ResourceData, meta interface{}) error {
	//Get the client connection.
	conn := meta.(*matlas.Client)
	projectID := d.Id()

	var options *matlas.ListOptions
	whitelist, _, err := conn.ProjectIPWhitelist.List(context.Background(), projectID, options)
	if err != nil {
		return fmt.Errorf(errorGetInfo, err)
	}

	deleteEntryWhiteList := func(f func(string)) {
		for _, entry := range whitelist {
			if entry.CIDRBlock != "" {
				f(entry.CIDRBlock)
			} else if entry.IPAddress != "" {
				f(entry.IPAddress)
			}
		}
	}

	var er error
	deleteEntryWhiteList(func(entry string) {
		_, err := conn.ProjectIPWhitelist.Delete(context.Background(), projectID, entry)
		if err != nil {
			er = fmt.Errorf("error deleting project IP whitelist: %s", err)
		}
	})
	if er != nil {
		return er
	}
	return nil
}

func flattenProjectIPWhitelist(whitelists []matlas.ProjectIPWhitelist) []map[string]interface{} {
	var results []map[string]interface{}

	if len(whitelists) > 0 {
		results = make([]map[string]interface{}, len(whitelists))

		for k, whitelist := range whitelists {
			results[k] = map[string]interface{}{
				"project_id": whitelist.GroupID,
				"cidr_block": whitelist.CIDRBlock,
				"ip_address": whitelist.IPAddress,
				"comment":    whitelist.Comment,
			}
		}
	}
	return results
}

func expandProjectIPWhitelist(d *schema.ResourceData) []*matlas.ProjectIPWhitelist {
	var whitelist []*matlas.ProjectIPWhitelist
	if v, ok := d.GetOk("entry"); ok {
		if rs := v.(*schema.Set).List(); len(rs) > 0 {
			whitelist = make([]*matlas.ProjectIPWhitelist, len(rs))
			for k, r := range rs {
				roleMap := r.(map[string]interface{})
				whitelist[k] = &matlas.ProjectIPWhitelist{
					CIDRBlock: roleMap["cidr_block"].(string),
					IPAddress: roleMap["ip_address"].(string),
					Comment:   roleMap["comment"].(string),
				}
			}
		}
	}
	return whitelist
}
