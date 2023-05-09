package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/openlyinc/pointy"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func resourceMongoDBAtlasOrganization() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBAtlasOrganizationCreate,
		ReadContext:   resourceMongoDBAtlasOrganizationRead,
		UpdateContext: resourceMongoDBAtlasOrganizationUpdate,
		DeleteContext: resourceMongoDBAtlasOrganizationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasOrganizationImportState,
		},
		Schema: map[string]*schema.Schema{
			"org_owner_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"org_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"private_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"role_names": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceMongoDBAtlasOrganizationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	createRequest := new(matlas.CreateOrganizationRequest)
	apiKey := &matlas.APIKeyInput{}

	createRequest.Name = d.Get("name").(string)
	createRequest.OrgOwnerID = pointy.String(d.Get("org_owner_id").(string))
	apiKey.Desc = d.Get("description").(string)
	apiKey.Roles = expandStringList(d.Get("role_names").(*schema.Set).List())

	organization, resp, err := conn.Organizations.Create(ctx, &matlas.CreateOrganizationRequest{
		Name:       createRequest.Name,
		OrgOwnerID: createRequest.OrgOwnerID,
		APIKey:     apiKey,
	})
	//organization, resp, err := conn.Organizations.Create(ctx, createRequest)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf("error create Organization: %s", err))
	}

	if err := d.Set("private_key", organization.APIKey.PrivateKey); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `public_key`: %s", err))
	}

	if err := d.Set("public_key", organization.APIKey.PublicKey); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `public_key`: %s", err))
	}

	d.SetId(encodeStateID(map[string]string{
		"org_id": organization.Organization.ID,
	}))

	return resourceMongoDBAtlasOrganizationRead(ctx, d, meta)
}

func resourceMongoDBAtlasOrganizationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	config := Config{
		PublicKey:  d.Get("public_key").(string),
		PrivateKey: d.Get("private_key").(string),
	}

	clients, _ := config.NewClient(ctx)
	conn := clients.(*MongoDBClient).Atlas

	ids := decodeStateID(d.Id())
	orgID := ids["org_id"]

	resp, err := conn.Organizations.Delete(ctx, orgID)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusBadRequest {
			log.Printf("warning Organization deleted will recreate: %s \n", err.Error())
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error deleting organization information: %s", err))
	}
	d.SetId("")

	return nil
}

func resourceMongoDBAtlasOrganizationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceMongoDBAtlasOrganizationRead(ctx, d, meta)
}

func resourceMongoDBAtlasOrganizationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	config := Config{
		PublicKey:  d.Get("public_key").(string),
		PrivateKey: d.Get("private_key").(string),
	}

	clients, _ := config.NewClient(ctx)
	conn := clients.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	orgID := ids["org_id"]

	_, err := conn.Organizations.Delete(ctx, orgID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error Organization: %s", err))
	}
	return nil
}

func resourceMongoDBAtlasOrganizationImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*MongoDBClient).Atlas
	orgID := d.Id()

	r, _, err := conn.Organizations.Get(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("couldn't import organization %s , error: %s", orgID, err)
	}

	if err := d.Set("name", r.Name); err != nil {
		return nil, fmt.Errorf("error setting `name`: %s", err)
	}

	if err := d.Set("org_id", r.ID); err != nil {
		return nil, fmt.Errorf("error setting `org_id`: %s", err)
	}

	d.SetId(encodeStateID(map[string]string{
		"org_id": orgID,
	}))

	return []*schema.ResourceData{d}, nil
}
