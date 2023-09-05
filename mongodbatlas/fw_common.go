package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	conf "github.com/mongodb/terraform-provider-mongodbatlas/config"
)

// RSCommon is used as an embedded struct for all framework resources. Implements the following plugin-framework defined functions:
// - Metadata
// - Configure
// client is left empty and populated by the framework when envoking Configure method.
// resourceName must be defined when creating an instance of a resource.
type RSCommon struct {
	config       *conf.Config
	resourceName string
}

func (r *RSCommon) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, r.resourceName)
}

func (r *RSCommon) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	c, err := configureClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError(errorConfigureSummary, err.Error())
		return
	}
	r.config = c
}

// DSCommon is used as an embedded struct for all framework data sources. Implements the following plugin-framework defined functions:
// - Metadata
// - Configure
// client is left empty and populated by the framework when envoking Configure method.
// dataSourceName must be defined when creating an instance of a data source.
type DSCommon struct {
	config         *conf.Config
	dataSourceName string
}

func (d *DSCommon) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, d.dataSourceName)
}

func (d *DSCommon) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, err := configureClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError(errorConfigureSummary, err.Error())
		return
	}
	d.config = client
}

func configureClient(providerData any) (*conf.Config, error) {
	if providerData == nil {
		return nil, nil
	}

	if c, ok := providerData.(*conf.Config); ok {
		return c, nil
	}

	return nil, fmt.Errorf(errorConfigure, providerData)
}
