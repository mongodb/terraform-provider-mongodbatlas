package mongodbatlas

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasProjectIPWhitelist() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasProjectIPWhitelistRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cidr_block": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
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
				ConflictsWith: []string{"aws_security_group", "cidr_block"},
				ValidateFunc:  validation.IsIPAddress,
			},
			"aws_security_group": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"ip_address", "cidr_block"},
			},
			"comment": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceMongoDBAtlasProjectIPWhitelistRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)
	cidrBlock := d.Get("cidr_block").(string)
	ipAddress := d.Get("ip_address").(string)
	awsSecurityGroup := d.Get("aws_security_group").(string)

	if cidrBlock == "" && ipAddress == "" && awsSecurityGroup == "" {
		return errors.New("cidr_block, ip_address or aws_security_group needs to contain a value")
	}
	var entry bytes.Buffer

	entry.WriteString(cidrBlock)
	entry.WriteString(ipAddress)
	entry.WriteString(awsSecurityGroup)

	whitelist, _, err := conn.ProjectIPWhitelist.Get(context.Background(), projectID, entry.String())
	if err != nil {
		return fmt.Errorf("error getting whitelist information: %s", err)
	}

	if err := d.Set("cidr_block", whitelist.CIDRBlock); err != nil {
		return fmt.Errorf("error setting `cidr_block` for Project whitelist: %s", err)
	}
	if err := d.Set("ip_address", whitelist.IPAddress); err != nil {
		return fmt.Errorf("error setting `ip_address` for Project whitelist: %s", err)
	}
	if err := d.Set("aws_security_group", whitelist.AwsSecurityGroup); err != nil {
		return fmt.Errorf("error setting `aws_security_group` for Project whitelist: %s", err)
	}
	if err := d.Set("comment", whitelist.Comment); err != nil {
		return fmt.Errorf("error setting `comment` for Project whitelist: %s", err)
	}

	d.SetId(resource.UniqueId())

	return nil
}
