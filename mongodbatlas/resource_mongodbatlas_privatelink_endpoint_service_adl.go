package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorADLServiceEndpointAdd = "error adding MongoDB ADL PrivateLink Endpoint Connection(%s): %s"
)

func resourceMongoDBAtlasPrivateLinkEndpointServiceADL() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBAtlasPrivateLinkEndpointServiceADLCreate,
		ReadContext:   resourceMongoDBAtlasPrivateLinkEndpointServiceADLRead,
		UpdateContext: resourceMongoDBAtlasPrivateLinkEndpointServiceADLUpdate,
		DeleteContext: resourceMongoDBAtlasPrivateLinkEndpointServiceADLDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasPrivateLinkEndpointServiceADLImportState,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"endpoint_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"provider_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"AWS"}, false),
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"DATA_LAKE"}, false),
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceMongoDBAtlasPrivateLinkEndpointServiceADLUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	endpointID := ids["endpoint_id"]

	privateLink, _, err := conn.DataLakes.GetPrivateLinkEndpoint(ctx, projectID, endpointID)
	if err != nil {
		return diag.Errorf("error getting ADL PrivateLink Endpoint Information: %s", err)
	}

	if d.HasChange("comment") {
		privateLink.Comment = d.Get("comment").(string)
	}

	_, _, err = conn.DataLakes.CreatePrivateLinkEndpoint(context.Background(), projectID, privateLink)
	if err != nil {
		return diag.Errorf("error updating ADL PrivateLink endpoint (%s): %s", privateLink.EndpointID, err)
	}

	return resourceMongoDBAtlasPrivateLinkEndpointServiceADLRead(ctx, d, meta)
}

func resourceMongoDBAtlasPrivateLinkEndpointServiceADLCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)

	privateLinkRequest := &matlas.PrivateLinkEndpointDataLake{
		EndpointID: d.Get("endpoint_id").(string),
		Type:       d.Get("type").(string),
		Provider:   d.Get("provider_name").(string),
		Comment:    d.Get("comment").(string),
	}

	_, _, err := conn.DataLakes.CreatePrivateLinkEndpoint(ctx, projectID, privateLinkRequest)
	if err != nil {
		return diag.Errorf(errorADLServiceEndpointAdd, privateLinkRequest.EndpointID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":  projectID,
		"endpoint_id": privateLinkRequest.EndpointID,
	}))

	return resourceMongoDBAtlasPrivateLinkEndpointServiceADLRead(ctx, d, meta)
}

func resourceMongoDBAtlasPrivateLinkEndpointServiceADLRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	endpointID := ids["endpoint_id"]

	privateLinkResponse, _, err := conn.DataLakes.GetPrivateLinkEndpoint(ctx, projectID, endpointID)
	if err != nil {
		// case 404
		// deleted in the backend case
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
			return nil
		}

		return diag.Errorf("error getting adl private link endpoint  information: %s", err)
	}

	if err := d.Set("endpoint_id", privateLinkResponse.EndpointID); err != nil {
		return diag.Errorf("error setting `endpoint_id` for endpoint_id (%s): %s", d.Id(), err)
	}

	if err := d.Set("type", privateLinkResponse.Type); err != nil {
		return diag.Errorf("error setting `type` for endpoint_id (%s): %s", d.Id(), err)
	}

	if err := d.Set("comment", privateLinkResponse.Comment); err != nil {
		return diag.Errorf("error setting `comment` for endpoint_id (%s): %s", d.Id(), err)
	}

	if err := d.Set("provider_name", privateLinkResponse.Provider); err != nil {
		return diag.Errorf("error setting `provider_name` for endpoint_id (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceMongoDBAtlasPrivateLinkEndpointServiceADLDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	endpointID := ids["endpoint_id"]

	_, err := conn.DataLakes.DeletePrivateLinkEndpoint(ctx, projectID, endpointID)
	if err != nil {
		return diag.Errorf("error deleting adl private link endpoint(%s): %s", endpointID, err)
	}

	return nil
}

func resourceMongoDBAtlasPrivateLinkEndpointServiceADLImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*MongoDBClient).Atlas

	parts := strings.SplitN(d.Id(), "--", 2)
	if len(parts) != 2 {
		return nil, errors.New("import format error: to import a search index, use the format {project_id}--{endpoint_id}")
	}

	projectID := parts[0]
	endpointID := parts[1]

	_, _, err := conn.DataLakes.GetPrivateLinkEndpoint(ctx, projectID, endpointID)
	if err != nil {
		return nil, fmt.Errorf("couldn't adl private link endpoint (%s) in projectID (%s) , error: %s", endpointID, projectID, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		log.Printf("[WARN] Error setting project_id for (%s): %s", projectID, err)
	}

	if err := d.Set("endpoint_id", endpointID); err != nil {
		log.Printf("[WARN] Error setting index_id for (%s): %s", endpointID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":  projectID,
		"endpoint_id": endpointID,
	}))

	return []*schema.ResourceData{d}, nil
}
