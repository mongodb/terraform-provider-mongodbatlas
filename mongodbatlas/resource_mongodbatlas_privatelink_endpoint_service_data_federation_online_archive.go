package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorPrivateEndpointServiceDataFederationOnlineArchiveCreate = "error creating a Private Endpoing for projectId %s: %s"
	errorPrivateEndpointServiceDataFederationOnlineArchiveDelete = "error deleting Private Endpoing %s for projectId %s: %s"
	errorPrivateEndpointServiceDataFederationOnlineArchiveRead   = "error reading Private Endpoing %s for projectId %s: %s"
	errorPrivateEndpointServiceDataFederationOnlineArchiveImport = "error importing Private Endpoing %s for projectId %s: %w"
	endpointType                                                 = "DATA_LAKE"
)

func resourceMongoDBAtlasPrivatelinkEndpointServiceDataFederationOnlineArchive() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBAtlasPrivatelinkEndpointServiceDataFederationOnlineArchiveCreate,
		ReadContext:   resourceMongoDBAtlasPrivatelinkEndpointServiceDataFederationOnlineArchiveRead,
		DeleteContext: resourceMongoDBAtlasPrivatelinkEndpointServiceDataFederationOnlineArchiveDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasPrivatelinkEndpointServiceDataFederationOnlineArchiveImportState,
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
				ForceNew: true,
			},
			"provider_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceMongoDBAtlasPrivatelinkEndpointServiceDataFederationOnlineArchiveCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)
	endpointID := d.Get("endpoint_id").(string)

	_, _, err := conn.DataLakes.CreatePrivateLinkEndpoint(ctx, projectID, newPrivateLinkEndpointDataLake(d))
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateEndpointServiceDataFederationOnlineArchiveCreate, projectID, err))
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":  projectID,
		"endpoint_id": endpointID,
	}))

	return resourceMongoDBAtlasPrivatelinkEndpointServiceDataFederationOnlineArchiveRead(ctx, d, meta)
}

func resourceMongoDBAtlasPrivatelinkEndpointServiceDataFederationOnlineArchiveRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	endopointID := ids["endpoint_id"]

	privateEndpoint, resp, err := conn.DataLakes.GetPrivateLinkEndpoint(context.Background(), projectID, endopointID)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.Errorf(errorPrivateEndpointServiceDataFederationOnlineArchiveRead, endopointID, projectID, err)
	}

	if err := d.Set("comment", privateEndpoint.Comment); err != nil {
		return diag.Errorf(errorPrivateEndpointServiceDataFederationOnlineArchiveRead, endopointID, projectID, err)
	}

	if err := d.Set("provider_name", privateEndpoint.Provider); err != nil {
		return diag.Errorf(errorPrivateEndpointServiceDataFederationOnlineArchiveRead, endopointID, projectID, err)
	}

	if err := d.Set("type", privateEndpoint.Type); err != nil {
		return diag.Errorf(errorPrivateEndpointServiceDataFederationOnlineArchiveRead, endopointID, projectID, err)
	}

	return nil
}

func resourceMongoDBAtlasPrivatelinkEndpointServiceDataFederationOnlineArchiveDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	endpointID := ids["endpoint_id"]

	_, err := conn.DataLakes.DeletePrivateLinkEndpoint(ctx, projectID, endpointID)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateEndpointServiceDataFederationOnlineArchiveDelete, endpointID, projectID, err))
	}

	d.SetId("")

	return nil
}

func resourceMongoDBAtlasPrivatelinkEndpointServiceDataFederationOnlineArchiveImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*MongoDBClient).Atlas
	projectID, endpointID, err := splitAtlasPrivatelinkEndpointServiceDataFederationOnlineArchive(d.Id())
	if err != nil {
		return nil, err
	}

	privateEndpoint, _, err := conn.DataLakes.GetPrivateLinkEndpoint(ctx, projectID, endpointID)
	if err != nil {
		return nil, fmt.Errorf(errorPrivateEndpointServiceDataFederationOnlineArchiveImport, endpointID, projectID, err)
	}

	if err := d.Set("comment", privateEndpoint.Comment); err != nil {
		return nil, fmt.Errorf(errorPrivateEndpointServiceDataFederationOnlineArchiveImport, endpointID, projectID, err)
	}

	if err := d.Set("provider_name", privateEndpoint.Provider); err != nil {
		return nil, fmt.Errorf(errorPrivateEndpointServiceDataFederationOnlineArchiveImport, endpointID, projectID, err)
	}

	if err := d.Set("type", privateEndpoint.Type); err != nil {
		return nil, fmt.Errorf(errorPrivateEndpointServiceDataFederationOnlineArchiveImport, endpointID, projectID, err)
	}

	if err := d.Set("endpoint_id", privateEndpoint.EndpointID); err != nil {
		return nil, fmt.Errorf(errorPrivateEndpointServiceDataFederationOnlineArchiveImport, endpointID, projectID, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		return nil, fmt.Errorf(errorPrivateEndpointServiceDataFederationOnlineArchiveImport, endpointID, projectID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":  projectID,
		"endpoint_id": endpointID,
	}))

	return []*schema.ResourceData{d}, nil
}

func splitAtlasPrivatelinkEndpointServiceDataFederationOnlineArchive(id string) (projectID, endpointID string, err error) {
	var parts = strings.Split(id, "--")

	if len(parts) != 2 {
		err = errors.New("import format error: to import a Data Lake, use the format {project_id}--{name}")
		return
	}

	projectID = parts[0]
	endpointID = parts[1]

	return
}

func newPrivateLinkEndpointDataLake(d *schema.ResourceData) *matlas.PrivateLinkEndpointDataLake {
	out := matlas.PrivateLinkEndpointDataLake{
		EndpointID: d.Get("endpoint_id").(string),
		Type:       endpointType,
	}

	if v, ok := d.GetOk("comment"); ok {
		out.Comment = v.(string)
	}

	if v, ok := d.GetOk("provider_name"); ok && v != "" {
		out.Provider = strings.ToUpper(v.(string))
	}

	return &out
}
