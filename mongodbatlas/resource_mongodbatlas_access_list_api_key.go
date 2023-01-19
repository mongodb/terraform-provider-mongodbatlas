package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func resourceMongoDBAtlasAccessListAPIKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBAtlasAccessListAPIKeyCreate,
		ReadContext:   resourceMongoDBAtlasAccessListAPIKeyRead,
		UpdateContext: resourceMongoDBAtlasAccessListAPIKeyUpdate,
		DeleteContext: resourceMongoDBAtlasAccessListAPIKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasAccessListAPIKeyImportState,
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
				ConflictsWith: []string{"cidr_block"},
				ValidateFunc:  validation.IsIPAddress,
			},
		},
	}
}

func resourceMongoDBAtlasAccessListAPIKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
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

	createReq := matlas.AccessListAPIKeysReq{}
	createReq.CidrBlock = CIDRBlock
	createReq.IPAddress = IPAddress

	createRequest := []*matlas.AccessListAPIKeysReq{}
	createRequest = append(createRequest, &createReq)

	_, resp, err := conn.AccessListAPIKeys.Create(ctx, orgID, apiKeyID, createRequest)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf("error create API key: %s", err))
	}

	d.SetId(encodeStateID(map[string]string{
		"org_id":     orgID,
		"api_key_id": apiKeyID,
		"entry":      entry,
	}))

	return resourceMongoDBAtlasAccessListAPIKeyRead(ctx, d, meta)
}

func resourceMongoDBAtlasAccessListAPIKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	orgID := ids["org_id"]
	apiKeyID := ids["api_key_id"]

	apiKey, _, err := conn.AccessListAPIKeys.Get(ctx, orgID, apiKeyID, strings.ReplaceAll(ids["entry"], "/", "%2F"))
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting api key information: %s", err))
	}

	if err := d.Set("api_key_id", apiKeyID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `api_key_id`: %s", err))
	}

	if err := d.Set("ip_address", apiKey.IPAddress); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `ip_address`: %s", err))
	}

	if err := d.Set("cidr_block", apiKey.CidrBlock); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `cidr_block`: %s", err))
	}

	d.SetId(encodeStateID(map[string]string{
		"org_id":     orgID,
		"api_key_id": apiKeyID,
		"entry":      ids["entry"],
	}))

	return nil
}

func resourceMongoDBAtlasAccessListAPIKeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceMongoDBAtlasAccessListAPIKeyRead(ctx, d, meta)
}

func resourceMongoDBAtlasAccessListAPIKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	orgID := ids["org_id"]
	apiKeyID := ids["api_key_id"]

	_, err := conn.AccessListAPIKeys.Delete(ctx, orgID, apiKeyID, strings.ReplaceAll(ids["entry"], "/", "%2F"))
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting API Key: %s", err))
	}
	return nil
}

func resourceMongoDBAtlasAccessListAPIKeyImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*MongoDBClient).Atlas

	parts := strings.SplitN(d.Id(), "-", 3)
	if len(parts) != 3 {
		return nil, errors.New("import format error: to import a api key use the format {org_id}-{api_key_id}-{ip_address}")
	}

	orgID := parts[0]
	apiKeyID := parts[1]
	entry := parts[2]

	r, _, err := conn.AccessListAPIKeys.Get(ctx, orgID, apiKeyID, strings.ReplaceAll(entry, "/", "%2F"))
	if err != nil {
		return nil, fmt.Errorf("couldn't import api key %s in project %s, error: %s", orgID, apiKeyID, err)
	}

	if err := d.Set("org_id", orgID); err != nil {
		return nil, fmt.Errorf("error setting `org_id`: %s", err)
	}

	if err := d.Set("ip_address", r.IPAddress); err != nil {
		return nil, fmt.Errorf("error setting `ip_address`: %s", err)
	}

	if err := d.Set("cidr_block", r.CidrBlock); err != nil {
		return nil, fmt.Errorf("error setting `cidr_block`: %s", err)
	}

	d.SetId(encodeStateID(map[string]string{
		"org_id":     orgID,
		"api_key_id": apiKeyID,
		"entry":      entry,
	}))

	return []*schema.ResourceData{d}, nil
}

func flattenAccessListAPIKeys(ctx context.Context, conn *matlas.Client, orgID string, accessListAPIKeys []*matlas.AccessListAPIKey) []map[string]interface{} {
	var results []map[string]interface{}

	if len(accessListAPIKeys) > 0 {
		results = make([]map[string]interface{}, len(accessListAPIKeys))
		for k, accessListAPIKey := range accessListAPIKeys {
			results[k] = map[string]interface{}{
				"ip_address":        accessListAPIKey.IPAddress,
				"cidr_block":        accessListAPIKey.CidrBlock,
				"created":           accessListAPIKey.Created,
				"access_count":      accessListAPIKey.Count,
				"last_used":         accessListAPIKey.LastUsed,
				"last_used_address": accessListAPIKey.LastUsedAddress,
			}
		}
	}
	return results
}
