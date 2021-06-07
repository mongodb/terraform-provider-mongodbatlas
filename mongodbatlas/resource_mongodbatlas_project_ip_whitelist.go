package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorWhitelistCreate = "error creating Project IP Whitelist information: %s"
	errorWhitelistRead   = "error getting Project IP Whitelist information: %s"
	// errorWhitelistUpdate  = "error updating Project IP Whitelist information: %s"
	errorWhitelistDelete  = "error deleting Project IP Whitelist information: %s"
	errorWhitelistSetting = "error setting `%s` for Project IP Whitelist (%s): %s"
)

func resourceMongoDBAtlasProjectIPWhitelist() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoDBAtlasProjectIPWhitelistCreate,
		Read:   resourceMongoDBAtlasProjectIPWhitelistRead,
		Delete: resourceMongoDBAtlasProjectIPWhitelistDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMongoDBAtlasIPWhitelistImportState,
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
				ConflictsWith: []string{"aws_security_group", "ip_address"},
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
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"aws_security_group", "cidr_block"},
				ValidateFunc:  validation.IsIPAddress,
			},
			// You must configure VPC peering for your project before you can whitelist an AWS security group.
			"aws_security_group": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"ip_address", "cidr_block"},
			},
			"comment": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Read:   schema.DefaultTimeout(45 * time.Minute),
			Delete: schema.DefaultTimeout(45 * time.Minute),
		},
	}
}

func resourceMongoDBAtlasProjectIPWhitelistCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)
	cidrBlock := d.Get("cidr_block").(string)
	ipAddress := d.Get("ip_address").(string)
	awsSecurityGroup := d.Get("aws_security_group").(string)

	if cidrBlock == "" && ipAddress == "" && awsSecurityGroup == "" {
		return errors.New("cidr_block, ip_address or aws_security_group needs to contain a value")
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"created", "failed"},
		Refresh: func() (interface{}, string, error) {
			whitelist, _, err := conn.ProjectIPWhitelist.Create(context.Background(), projectID, []*matlas.ProjectIPWhitelist{
				{
					AwsSecurityGroup: awsSecurityGroup,
					CIDRBlock:        cidrBlock,
					IPAddress:        ipAddress,
					Comment:          d.Get("comment").(string),
				},
			})
			if err != nil {
				if strings.Contains(fmt.Sprint(err), "Unexpected error") ||
					strings.Contains(fmt.Sprint(err), "UNEXPECTED_ERROR") ||
					strings.Contains(fmt.Sprint(err), "500") {
					return nil, "pending", nil
				}
				return nil, "failed", fmt.Errorf(errorWhitelistCreate, err)
			}

			if len(whitelist) > 0 {
				whiteListEntry := ipAddress
				if cidrBlock != "" {
					whiteListEntry = cidrBlock
				}

				for _, entry := range whitelist {
					if entry.IPAddress == whiteListEntry || entry.CIDRBlock == whiteListEntry {
						return whitelist, "created", nil
					}
				}
				return nil, "pending", nil
			}

			return whitelist, "created", nil
		},
		Timeout:    45 * time.Minute,
		Delay:      4 * time.Second,
		MinTimeout: 2 * time.Second,
	}

	// Wait, catching any errors
	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(errorWhitelistCreate, err)
	}

	var entry string

	switch {
	case cidrBlock != "":
		entry = cidrBlock
	case ipAddress != "":
		entry = ipAddress
	default:
		entry = awsSecurityGroup
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

	return resource.Retry(2*time.Minute, func() *resource.RetryError {
		whitelist, _, err := conn.ProjectIPWhitelist.Get(context.Background(), ids["project_id"], ids["entry"])
		if err != nil {
			switch {
			case strings.Contains(fmt.Sprint(err), "500"):
				return resource.RetryableError(err)
			case strings.Contains(fmt.Sprint(err), "404"):
				if !d.IsNewResource() {
					d.SetId("")
					return nil
				}
				return resource.RetryableError(err)
			default:
				return resource.NonRetryableError(fmt.Errorf(errorWhitelistRead, err))
			}
		}

		if whitelist != nil {
			if err := d.Set("aws_security_group", whitelist.AwsSecurityGroup); err != nil {
				return resource.NonRetryableError(fmt.Errorf(errorWhitelistSetting, "aws_security_group", ids["project_id"], err))
			}

			if err := d.Set("cidr_block", whitelist.CIDRBlock); err != nil {
				return resource.NonRetryableError(fmt.Errorf(errorWhitelistSetting, "cidr_block", ids["project_id"], err))
			}

			if err := d.Set("ip_address", whitelist.IPAddress); err != nil {
				return resource.NonRetryableError(fmt.Errorf(errorWhitelistSetting, "ip_address", ids["project_id"], err))
			}

			if err := d.Set("comment", whitelist.Comment); err != nil {
				return resource.NonRetryableError(fmt.Errorf(errorWhitelistSetting, "comment", ids["project_id"], err))
			}
		}

		return nil
	})
}

func resourceMongoDBAtlasProjectIPWhitelistDelete(d *schema.ResourceData, meta interface{}) error {
	// Get the client connection.
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := conn.ProjectIPWhitelist.Delete(context.Background(), ids["project_id"], ids["entry"])
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "500") ||
				strings.Contains(fmt.Sprint(err), "Unexpected error") ||
				strings.Contains(fmt.Sprint(err), "UNEXPECTED_ERROR") {
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(fmt.Errorf(errorWhitelistDelete, err))
		}

		entry, _, err := conn.ProjectIPWhitelist.Get(context.Background(), ids["project_id"], ids["entry"])
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "404") ||
				strings.Contains(fmt.Sprint(err), "ATLAS_WHITELIST_NOT_FOUND") {
				return nil
			}

			return resource.RetryableError(err)
		}

		if entry != nil {
			return resource.RetryableError(fmt.Errorf(errorWhitelistDelete, "Whitelist still exists"))
		}

		return nil
	})
}

func resourceMongoDBAtlasIPWhitelistImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*matlas.Client)

	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 2 {
		return nil, errors.New("import format error: to import a peer, use the format {project_id}-{whitelist_entry}")
	}

	projectID := parts[0]
	entry := parts[1]

	_, _, err := conn.ProjectIPWhitelist.Get(context.Background(), projectID, entry)
	if err != nil {
		return nil, fmt.Errorf("couldn't import entry whitelist %s in project %s, error: %s", entry, projectID, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		log.Printf("[WARN] Error setting project_id for (%s): %s", projectID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"entry":      entry,
	}))

	return []*schema.ResourceData{d}, nil
}
