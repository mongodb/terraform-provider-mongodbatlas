package mongodbatlas

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceMongoDBAtlasProjectIPAccessList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasProjectIPAccessListRead,
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

func dataSourceMongoDBAtlasProjectIPAccessListRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)
	cidrBlock := d.Get("cidr_block").(string)
	ipAddress := d.Get("ip_address").(string)
	awsSecurityGroup := d.Get("aws_security_group").(string)

	if cidrBlock == "" && ipAddress == "" && awsSecurityGroup == "" {
		return diag.FromErr(errors.New("cidr_block, ip_address or aws_security_group needs to contain a value"))
	}
	var entry bytes.Buffer

	entry.WriteString(cidrBlock)
	entry.WriteString(ipAddress)
	entry.WriteString(awsSecurityGroup)

	accessList, _, err := conn.ProjectIPAccessList.Get(ctx, projectID, entry.String())
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting access list information: %s", err))
	}

	if err := d.Set("cidr_block", accessList.CIDRBlock); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `cidr_block` for the project access list: %s", err))
	}
	if err := d.Set("ip_address", accessList.IPAddress); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `ip_address` for the project access list: %s", err))
	}
	if err := d.Set("aws_security_group", accessList.AwsSecurityGroup); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `aws_security_group` for the project access list: %s", err))
	}
	if err := d.Set("comment", accessList.Comment); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `comment` for the project access list: %s", err))
	}

	d.SetId(resource.UniqueId())

	return nil
}
