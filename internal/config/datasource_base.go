package config

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ImplementedDataSource interface {
	datasource.DataSourceWithConfigure
	GetName() string
	SetClient(*MongoDBClient)
}

func AnalyticsDataSourceFunc(iDataSource datasource.DataSource) func() datasource.DataSource {
	commonDataSource, ok := iDataSource.(ImplementedDataSource)
	if !ok {
		panic(fmt.Sprintf("data source %T didn't comply with the ImplementedDataSource interface", iDataSource))
	}
	return func() datasource.DataSource {
		return analyticsDataSource(commonDataSource)
	}
}

// DSCommon is used as an embedded struct for all framework data sources. Implements the following plugin-framework defined functions:
// - Metadata
// - Configure
// Client is left empty and populated by the framework when envoking Configure method.
// DataSourceName must be defined when creating an instance of a data source.
//
// When used as a wrapper (ImplementedDataSource is set), it intercepts Read to add analytics tracking.
// When embedded in a data source struct, the data source's own Read method is used.
type DSCommon struct {
	ImplementedDataSource // Set when used as a wrapper, nil when embedded
	Client                *MongoDBClient
	DataSourceName        string
}

func (d *DSCommon) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, d.DataSourceName)
}

func (d *DSCommon) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	if d.ImplementedDataSource != nil {
		// When used as a wrapper, delegate to the wrapped data source
		d.ImplementedDataSource.Schema(ctx, req, resp)
	}
	// When embedded, the data source's own Schema method is used
}

func (d *DSCommon) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, err := configureClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError(errorConfigureSummary, err.Error())
		return
	}
	d.Client = client
	// If used as a wrapper, set the client on the wrapped data source
	if d.ImplementedDataSource != nil {
		d.ImplementedDataSource.SetClient(client)
	}
}

// Read intercepts the Read operation when DSCommon is used as a wrapper to add analytics tracking.
// When DSCommon is embedded, this method is not used (the data source's own Read method is called).
func (d *DSCommon) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.ImplementedDataSource == nil {
		// This shouldn't happen, but if DSCommon is embedded, the data source's Read is used instead
		return
	}
	extra := asUserAgentExtraFromProviderMeta(ctx, d.DataSourceName, UserAgentOperationValueRead, true, req.ProviderMeta)
	ctx = AddUserAgentExtra(ctx, extra)
	d.ImplementedDataSource.Read(ctx, req, resp)
}

func (d *DSCommon) GetName() string {
	return d.DataSourceName
}

func (d *DSCommon) SetClient(client *MongoDBClient) {
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

// analyticsDataSource wraps an ImplementedDataSource with DSCommon to add analytics tracking.
// We cannot return iDataSource directly because we need to intercept the Read operation
// to inject provider_meta information into the context before calling the actual data source method.
func analyticsDataSource(iDataSource ImplementedDataSource) datasource.DataSource {
	return &DSCommon{
		DataSourceName:        iDataSource.GetName(),
		ImplementedDataSource: iDataSource,
	}
}

// asUserAgentExtraFromProviderMeta extracts UserAgentExtra from provider_meta.
// This is a shared function used by both resources and data sources.
func asUserAgentExtraFromProviderMeta(ctx context.Context, name, reqOperation string, isDataSource bool, reqProviderMeta tfsdk.Config) UserAgentExtra {
	var meta ProviderMeta
	var nameValue string
	if isDataSource {
		nameValue = userAgentNameValueDataSource(name)
	} else {
		nameValue = userAgentNameValue(name)
	}
	uaExtra := UserAgentExtra{
		Name:      nameValue,
		Operation: reqOperation,
	}
	if reqProviderMeta.Raw.IsNull() {
		return uaExtra
	}
	diags := reqProviderMeta.Get(ctx, &meta)
	if diags.HasError() {
		return uaExtra
	}

	extrasLen := len(meta.UserAgentExtra.Elements())
	userExtras := make(map[string]types.String, extrasLen)
	diags.Append(meta.UserAgentExtra.ElementsAs(ctx, &userExtras, false)...)
	if diags.HasError() {
		return uaExtra
	}
	userExtrasString := make(map[string]string, extrasLen)
	for k, v := range userExtras {
		userExtrasString[k] = v.ValueString()
	}
	return uaExtra.Combine(UserAgentExtra{
		Extras:        userExtrasString,
		ModuleName:    meta.ModuleName.ValueString(),
		ModuleVersion: meta.ModuleVersion.ValueString(),
	})
}
