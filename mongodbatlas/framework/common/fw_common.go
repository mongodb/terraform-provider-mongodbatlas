package common

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/mongodbatlas/client"
)

const (
	errorConfigureSummary = "Unexpected Resource Configure Type"
	errorConfigure        = "expected *MongoDBClient, got: %T. Please report this issue to the provider developers"
)

// RSCommon is used as an embedded struct for all framework resources. Implements the following plugin-framework defined functions:
// - Metadata
// - Configure
// client is left empty and populated by the framework when envoking Configure method.
// resourceName must be defined when creating an instance of a resource.
type RSCommon struct {
	Client       *client.MongoDBClient
	ResourceName string
}

func (r *RSCommon) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, r.ResourceName)
}

func (r *RSCommon) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	c, err := configureClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError(errorConfigureSummary, err.Error())
		return
	}
	r.Client = c
}

// DSCommon is used as an embedded struct for all framework data sources. Implements the following plugin-framework defined functions:
// - Metadata
// - Configure
// client is left empty and populated by the framework when envoking Configure method.
// dataSourceName must be defined when creating an instance of a data source.
type DSCommon struct {
	Client         *client.MongoDBClient
	DataSourceName string
}

func (d *DSCommon) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, d.DataSourceName)
}

func (d *DSCommon) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	c, err := configureClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError(errorConfigureSummary, err.Error())
		return
	}
	d.Client = c
}

func configureClient(providerData any) (*client.MongoDBClient, error) {
	if providerData == nil {
		return nil, nil
	}

	if c, ok := providerData.(*client.MongoDBClient); ok {
		return c, nil
	}

	return nil, fmt.Errorf(errorConfigure, providerData)
}
