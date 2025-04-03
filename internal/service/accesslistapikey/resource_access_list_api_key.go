package accesslistapikey

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"go.mongodb.org/atlas-sdk/v20250312002/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate,
		ReadContext:   resourceRead,
		UpdateContext: resourceUpdate,
		DeleteContext: resourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceImport,
		},
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"api_key_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cidr_block": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"ip_address"},
				ValidateFunc: func(i any, k string) (s []string, es []error) {
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
				ConflictsWith: []string{"cidr_block"},
				ValidateFunc:  validation.IsIPAddress,
			},
		},
	}
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	orgID := d.Get("org_id").(string)
	apiKeyID := d.Get("api_key_id").(string)
	IPAddress := d.Get("ip_address").(string)
	CIDRBlock := d.Get("cidr_block").(string)

	var entry string

	switch {
	case CIDRBlock != "":
		parts := strings.SplitN(CIDRBlock, "/", 2)
		if parts[1] == "32" {
			entry = parts[0]
		} else {
			entry = CIDRBlock
		}
	case IPAddress != "":
		entry = IPAddress
	default:
		entry = IPAddress
	}

	accessList := &[]admin.UserAccessListRequest{
		{
			CidrBlock: conversion.StringPtr(CIDRBlock),
			IpAddress: conversion.StringPtr(IPAddress),
		},
	}

	_, resp, err := connV2.ProgrammaticAPIKeysApi.CreateApiKeyAccessList(ctx, orgID, apiKeyID, accessList).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error create API key: %s", err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"org_id":     orgID,
		"api_key_id": apiKeyID,
		"entry":      entry,
	}))

	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	orgID := ids["org_id"]
	apiKeyID := ids["api_key_id"]
	ipAddress := ids["entry"]

	apiKey, resp, err := connV2.ProgrammaticAPIKeysApi.GetApiKeyAccessList(ctx, orgID, ipAddress, apiKeyID).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) || validate.StatusBadRequest(resp) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error getting api key information: %s", err))
	}

	if err := d.Set("api_key_id", apiKeyID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `api_key_id`: %s", err))
	}

	if err := d.Set("ip_address", apiKey.IpAddress); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `ip_address`: %s", err))
	}

	if err := d.Set("cidr_block", apiKey.CidrBlock); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `cidr_block`: %s", err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"org_id":     orgID,
		"api_key_id": apiKeyID,
		"entry":      ipAddress,
	}))

	return nil
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return resourceRead(ctx, d, meta)
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	orgID := ids["org_id"]
	apiKeyID := ids["api_key_id"]
	ipAddress := ids["entry"]

	_, _, err := connV2.ProgrammaticAPIKeysApi.DeleteApiKeyAccessListEntry(ctx, orgID, apiKeyID, ipAddress).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting API Key: %s", err))
	}
	return nil
}

func resourceImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	parts := strings.SplitN(d.Id(), "-", 3)
	if len(parts) != 3 {
		return nil, errors.New("import format error: to import a api key use the format {org_id}-{api_key_id}-{ip_address}")
	}

	orgID := parts[0]
	apiKeyID := parts[1]
	ipAddress := parts[2]

	r, _, err := connV2.ProgrammaticAPIKeysApi.GetApiKeyAccessList(ctx, orgID, ipAddress, apiKeyID).Execute()
	if err != nil {
		return nil, fmt.Errorf("couldn't import api key %s in project %s, error: %s", orgID, apiKeyID, err)
	}

	if err := d.Set("org_id", orgID); err != nil {
		return nil, fmt.Errorf("error setting `org_id`: %s", err)
	}

	if err := d.Set("ip_address", r.IpAddress); err != nil {
		return nil, fmt.Errorf("error setting `ip_address`: %s", err)
	}

	if err := d.Set("cidr_block", r.CidrBlock); err != nil {
		return nil, fmt.Errorf("error setting `cidr_block`: %s", err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"org_id":     orgID,
		"api_key_id": apiKeyID,
		"entry":      ipAddress,
	}))

	return []*schema.ResourceData{d}, nil
}
