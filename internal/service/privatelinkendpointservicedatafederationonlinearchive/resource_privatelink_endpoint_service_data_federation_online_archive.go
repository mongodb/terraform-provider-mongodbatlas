package privatelinkendpointservicedatafederationonlinearchive

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/atlas-sdk/v20250312007/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	errorPrivateEndpointServiceDataFederationOnlineArchiveCreate = "error creating a Private Endpoint for projectId %s: %s"
	errorPrivateEndpointServiceDataFederationOnlineArchiveDelete = "error deleting Private Endpoint %s for projectId %s: %s"
	errorPrivateEndpointServiceDataFederationOnlineArchiveRead   = "error reading Private Endpoint %s for projectId %s: %s"
	errorPrivateEndpointServiceDataFederationOnlineArchiveUpdate = "error updating a Private Endpoint for projectId %s: %s"
	errorPrivateEndpointServiceDataFederationOnlineArchiveImport = "error importing Private Endpoint %s for projectId %s: %w"
	endpointType                                                 = "DATA_LAKE"
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
			},
			"region": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ValidateDiagFunc: validate.StringIsUppercase(),
			},
			"customer_endpoint_dns_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		// Timeouts are not being used but kept to avoid breaking changes.
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Hour),
			Delete: schema.DefaultTimeout(2 * time.Hour),
		},
	}
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)
	endpointID := d.Get("endpoint_id").(string)

	_, _, err := connV2.DataFederationApi.CreateDataFederationPrivateEndpoint(ctx, projectID, newPrivateNetworkEndpointIDEntry(d)).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateEndpointServiceDataFederationOnlineArchiveCreate, projectID, err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":  projectID,
		"endpoint_id": endpointID,
	}))

	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	endpointID := ids["endpoint_id"]

	privateEndpoint, resp, err := connV2.DataFederationApi.GetDataFederationPrivateEndpoint(ctx, projectID, endpointID).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) {
			d.SetId("")
			return nil
		}
		return diag.Errorf(errorPrivateEndpointServiceDataFederationOnlineArchiveRead, endpointID, projectID, err)
	}

	err = populateResourceData(d, privateEndpoint, projectID, endpointID)
	if err != nil {
		return diag.Errorf(errorPrivateEndpointServiceDataFederationOnlineArchiveRead, endpointID, projectID, err)
	}
	return nil
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	endpointID := ids["endpoint_id"]
	privateEndpoint, resp, err := connV2.DataFederationApi.GetDataFederationPrivateEndpoint(ctx, projectID, endpointID).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) {
			d.SetId("")
			return nil
		}
		return diag.Errorf(errorPrivateEndpointServiceDataFederationOnlineArchiveRead, endpointID, projectID, err)
	}

	if d.HasChange("comment") {
		privateEndpoint.Comment = conversion.StringPtr(d.Get("comment").(string))
		_, _, err = connV2.DataFederationApi.CreateDataFederationPrivateEndpoint(ctx, projectID, privateEndpoint).Execute()
		if err != nil {
			return diag.Errorf(errorPrivateEndpointServiceDataFederationOnlineArchiveCreate, endpointID, projectID)
		}
	}
	return resourceRead(ctx, d, meta)
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	endpointID := ids["endpoint_id"]

	_, err := connV2.DataFederationApi.DeleteDataFederationPrivateEndpoint(ctx, projectID, endpointID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateEndpointServiceDataFederationOnlineArchiveDelete, endpointID, projectID, err))
	}

	d.SetId("")

	return nil
}

func resourceImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID, endpointID, err := splitAtlasPrivatelinkEndpointServiceDataFederationOnlineArchive(d.Id())
	if err != nil {
		return nil, err
	}

	privateEndpoint, _, err := connV2.DataFederationApi.GetDataFederationPrivateEndpoint(ctx, projectID, endpointID).Execute()
	if err != nil {
		return nil, fmt.Errorf(errorPrivateEndpointServiceDataFederationOnlineArchiveImport, endpointID, projectID, err)
	}
	err = populateResourceData(d, privateEndpoint, projectID, endpointID)
	if err != nil {
		return nil, fmt.Errorf(errorPrivateEndpointServiceDataFederationOnlineArchiveImport, endpointID, projectID, err)
	}

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

func newPrivateNetworkEndpointIDEntry(d *schema.ResourceData) *admin.PrivateNetworkEndpointIdEntry {
	endpointType := endpointType
	out := admin.PrivateNetworkEndpointIdEntry{
		EndpointId: d.Get("endpoint_id").(string),
		Type:       &endpointType,
	}

	if v, ok := d.GetOk("comment"); ok {
		comment := v.(string)
		out.Comment = &comment
	}

	if v, ok := d.GetOk("provider_name"); ok && v != "" {
		providerName := strings.ToUpper(v.(string))
		out.Provider = &providerName
	}

	if v, ok := d.GetOk("region"); ok && v != "" {
		region := v.(string)
		out.Region = &region
	}

	if v, ok := d.GetOk("customer_endpoint_dns_name"); ok && v != "" {
		customerEndpointDNSName := v.(string)
		out.CustomerEndpointDNSName = &customerEndpointDNSName
	}

	return &out
}

func populateResourceData(d *schema.ResourceData, privateEndpoint *admin.PrivateNetworkEndpointIdEntry, projectID, endpointID string) error {
	if err := d.Set("comment", privateEndpoint.GetComment()); err != nil {
		return err
	}

	if err := d.Set("provider_name", privateEndpoint.GetProvider()); err != nil {
		return err
	}

	if err := d.Set("type", privateEndpoint.GetType()); err != nil {
		return err
	}

	if err := d.Set("region", privateEndpoint.GetRegion()); err != nil {
		return err
	}

	if err := d.Set("endpoint_id", privateEndpoint.GetEndpointId()); err != nil {
		return err
	}

	if err := d.Set("project_id", projectID); err != nil {
		return err
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":  projectID,
		"endpoint_id": endpointID,
	}))

	return d.Set("customer_endpoint_dns_name", privateEndpoint.GetCustomerEndpointDNSName())
}
