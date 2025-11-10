package config

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func NewAnalyticsResourceSDKv2(d *schema.Resource, name string) *schema.Resource {
	analyticsResource := &AnalyticsResourceSDKv2{
		resource: d,
		name:     name,
	}
	/*
		We are not initializing deprecated fields, for example Update to avoid the message:
			resource mongodbatlas_cloud_backup_snapshot: All fields are ForceNew or Computed w/out Optional, Update is superfluous

		Ensure no deprecated fields are used by running `staticcheck ./internal/service/... | grep -v 'd.GetOkExists'` and looking for (SA1019)
			GetOkExists we are using in many places; therefore, we use -v (invert match) to filter out lines with different deprecations
			Example line:
				internal/service/cluster/model_cluster.go:306:14: d.GetOkExists is deprecated: usage is discouraged due to undefined behaviors and may be removed in a future version of the SDK  (SA1019)
	*/
	resource := &schema.Resource{
		CustomizeDiff:                     d.CustomizeDiff,
		DeprecationMessage:                d.DeprecationMessage,
		Description:                       d.Description,
		EnableLegacyTypeSystemApplyErrors: d.EnableLegacyTypeSystemApplyErrors,
		EnableLegacyTypeSystemPlanErrors:  d.EnableLegacyTypeSystemPlanErrors,
		Identity:                          d.Identity,
		ResourceBehavior:                  d.ResourceBehavior,
		Schema:                            d.Schema,
		SchemaFunc:                        d.SchemaFunc,
		SchemaVersion:                     d.SchemaVersion,
		StateUpgraders:                    d.StateUpgraders,
		Timeouts:                          d.Timeouts,
		UpdateWithoutTimeout:              d.UpdateWithoutTimeout,
		UseJSONNumber:                     d.UseJSONNumber,
		ValidateRawResourceConfigFuncs:    d.ValidateRawResourceConfigFuncs,
	}
	importer := d.Importer
	if importer != nil {
		resource.Importer = &schema.ResourceImporter{
			StateContext: analyticsResource.resourceImport,
		}
	}
	// CreateContext or CreateWithoutTimeout, cannot use both
	if d.CreateContext != nil {
		resource.CreateContext = analyticsResource.CreateContext
	}
	if d.CreateWithoutTimeout != nil {
		resource.CreateWithoutTimeout = analyticsResource.CreateWithoutTimeout
	}
	// ReadContext or ReadWithoutTimeout, cannot use both
	if d.ReadContext != nil {
		resource.ReadContext = analyticsResource.ReadContext
	}
	if d.ReadWithoutTimeout != nil {
		resource.ReadWithoutTimeout = analyticsResource.ReadWithoutTimeout
	}
	// UpdateContext is not set on all resources
	if d.UpdateContext != nil {
		resource.UpdateContext = analyticsResource.UpdateContext
	}
	if d.UpdateWithoutTimeout != nil {
		resource.UpdateWithoutTimeout = analyticsResource.UpdateWithoutTimeout
	}
	// DeleteContext or DeleteWithoutTimeout, cannot use both
	if d.DeleteContext != nil {
		resource.DeleteContext = analyticsResource.DeleteContext
	}
	if d.DeleteWithoutTimeout != nil {
		resource.DeleteWithoutTimeout = analyticsResource.DeleteWithoutTimeout
	}
	return resource
}

type ProviderMetaSDKv2 struct {
	UserAgentExtra map[string]string `cty:"user_agent_extra"`
	ModuleName     *string           `cty:"module_name"`
	ModuleVersion  *string           `cty:"module_version"`
}

type AnalyticsResourceSDKv2 struct {
	resource *schema.Resource
	name     string
}

func (a *AnalyticsResourceSDKv2) CreateContext(ctx context.Context, r *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta, err := parseProviderMeta(r)
	if err != nil {
		return a.resource.CreateContext(ctx, r, m)
	}
	ctx = a.updateContextWithProviderMeta(ctx, meta, UserAgentOperationValueCreate)
	return a.resource.CreateContext(ctx, r, m)
}

func (a *AnalyticsResourceSDKv2) CreateWithoutTimeout(ctx context.Context, r *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta, err := parseProviderMeta(r)
	if err != nil {
		return a.resource.CreateWithoutTimeout(ctx, r, m)
	}
	ctx = a.updateContextWithProviderMeta(ctx, meta, UserAgentOperationValueCreate)
	return a.resource.CreateWithoutTimeout(ctx, r, m)
}

func (a *AnalyticsResourceSDKv2) ReadWithoutTimeout(ctx context.Context, r *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta, err := parseProviderMeta(r)
	if err != nil {
		return a.resource.ReadWithoutTimeout(ctx, r, m)
	}
	ctx = a.updateContextWithProviderMeta(ctx, meta, UserAgentOperationValueRead)
	return a.resource.ReadWithoutTimeout(ctx, r, m)
}

func (a *AnalyticsResourceSDKv2) ReadContext(ctx context.Context, r *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta, err := parseProviderMeta(r)
	if err != nil {
		return a.resource.ReadContext(ctx, r, m)
	}
	ctx = a.updateContextWithProviderMeta(ctx, meta, UserAgentOperationValueRead)
	return a.resource.ReadContext(ctx, r, m)
}

func (a *AnalyticsResourceSDKv2) UpdateContext(ctx context.Context, r *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta, err := parseProviderMeta(r)
	if err != nil {
		return a.resource.UpdateContext(ctx, r, m)
	}
	ctx = a.updateContextWithProviderMeta(ctx, meta, UserAgentOperationValueUpdate)
	return a.resource.UpdateContext(ctx, r, m)
}
func (a *AnalyticsResourceSDKv2) UpdateWithoutTimeout(ctx context.Context, r *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta, err := parseProviderMeta(r)
	if err != nil {
		return a.resource.UpdateWithoutTimeout(ctx, r, m)
	}
	ctx = a.updateContextWithProviderMeta(ctx, meta, UserAgentOperationValueUpdate)
	return a.resource.UpdateWithoutTimeout(ctx, r, m)
}

func (a *AnalyticsResourceSDKv2) DeleteContext(ctx context.Context, r *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta, err := parseProviderMeta(r)
	if err != nil {
		return a.resource.DeleteContext(ctx, r, m)
	}
	ctx = a.updateContextWithProviderMeta(ctx, meta, UserAgentOperationValueDelete)
	return a.resource.DeleteContext(ctx, r, m)
}

func (a *AnalyticsResourceSDKv2) DeleteWithoutTimeout(ctx context.Context, r *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta, err := parseProviderMeta(r)
	if err != nil {
		return a.resource.DeleteWithoutTimeout(ctx, r, m)
	}
	ctx = a.updateContextWithProviderMeta(ctx, meta, UserAgentOperationValueDelete)
	return a.resource.DeleteWithoutTimeout(ctx, r, m)
}

func (a *AnalyticsResourceSDKv2) resourceImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	// Import doesn't have providerMeta
	ctx = AddUserAgentExtra(ctx, UserAgentExtra{
		Name:      a.name,
		Operation: UserAgentOperationValueImport,
	})
	return a.resource.Importer.StateContext(ctx, d, meta)
}

func (a *AnalyticsResourceSDKv2) updateContextWithProviderMeta(ctx context.Context, meta ProviderMetaSDKv2, operationName string) context.Context {
	moduleName := ""
	if meta.ModuleName != nil {
		moduleName = *meta.ModuleName
	}
	moduleVersion := ""
	if meta.ModuleVersion != nil {
		moduleVersion = *meta.ModuleVersion
	}

	uaExtra := UserAgentExtra{
		Name:          userAgentNameValue(a.name),
		Operation:     operationName,
		Extras:        meta.UserAgentExtra,
		ModuleName:    moduleName,
		ModuleVersion: moduleVersion,
	}
	ctx = AddUserAgentExtra(ctx, uaExtra)
	return ctx
}

func parseProviderMeta(r *schema.ResourceData) (ProviderMetaSDKv2, error) {
	meta := ProviderMetaSDKv2{}
	err := r.GetProviderMeta(&meta)
	if err != nil {
		log.Printf("[WARN] failed to decode provider meta: %s, meta: %v", err, meta)
	}
	return meta, err
}
