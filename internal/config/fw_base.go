package config

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

const (
	errorConfigureSummary = "Unexpected Resource Configure Type"
	errorConfigure        = "expected *MongoDBClient, got: %T. Please report this issue to the provider developers"
)

// RSCommon is used as an embedded struct for all framework resources. Implements the following plugin-framework defined functions:
// - Metadata
// - Configure
// Client is left empty and populated by the framework when envoking Configure method.
// ResourceName must be defined when creating an instance of a resource.
type RSCommon struct {
	Client       *MongoDBClient
	ResourceName string
}

func (r *RSCommon) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, r.ResourceName)
}

func (r *RSCommon) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, err := configureClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError(errorConfigureSummary, err.Error())
		return
	}
	r.Client = client
}

// DSCommon is used as an embedded struct for all framework data sources. Implements the following plugin-framework defined functions:
// - Metadata
// - Configure
// Client is left empty and populated by the framework when envoking Configure method.
// DataSourceName must be defined when creating an instance of a data source.
type DSCommon struct {
	Client         *MongoDBClient
	DataSourceName string
}

func (d *DSCommon) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, d.DataSourceName)
}

func (d *DSCommon) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, err := configureClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError(errorConfigureSummary, err.Error())
		return
	}
	d.Client = client
}

func configureClient(providerData any) (*MongoDBClient, error) {
	if providerData == nil {
		return nil, nil
	}

	if client, ok := providerData.(*MongoDBClient); ok {
		return client, nil
	}

	return nil, fmt.Errorf(errorConfigure, providerData)
}
