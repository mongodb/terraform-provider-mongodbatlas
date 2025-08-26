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
	resource := &schema.Resource{
		Schema:                            d.Schema,
		SchemaFunc:                        d.SchemaFunc,
		SchemaVersion:                     d.SchemaVersion,
		Identity:                          d.Identity,
		MigrateState:                      d.MigrateState,
		StateUpgraders:                    d.StateUpgraders,
		Create:                            d.Create,
		Read:                              d.Read,
		Update:                            d.Update,
		Delete:                            d.Delete,
		Exists:                            d.Exists,
		CreateWithoutTimeout:              d.CreateWithoutTimeout,
		ReadWithoutTimeout:                d.ReadWithoutTimeout,
		UpdateWithoutTimeout:              d.UpdateWithoutTimeout,
		DeleteWithoutTimeout:              d.DeleteWithoutTimeout,
		CustomizeDiff:                     d.CustomizeDiff,
		DeprecationMessage:                d.DeprecationMessage,
		Timeouts:                          d.Timeouts,
		Description:                       d.Description,
		UseJSONNumber:                     d.UseJSONNumber,
		EnableLegacyTypeSystemApplyErrors: d.EnableLegacyTypeSystemApplyErrors,
		EnableLegacyTypeSystemPlanErrors:  d.EnableLegacyTypeSystemPlanErrors,
		ResourceBehavior:                  d.ResourceBehavior,
		ValidateRawResourceConfigFuncs:    d.ValidateRawResourceConfigFuncs,
		CreateContext:                     analyticsResource.CreateContext,
		ReadContext:                       analyticsResource.ReadContext,
		UpdateContext:                     analyticsResource.UpdateContext,
		DeleteContext:                     analyticsResource.DeleteContext,
	}
	importer := d.Importer
	if importer != nil {
		resource.Importer = &schema.ResourceImporter{
			State:        importer.State,
			StateContext: analyticsResource.resourceImport,
		}
	}
	return resource
}

type AnalyticsResourceSDKv2 struct {
	resource *schema.Resource
	name     string
}

func (a *AnalyticsResourceSDKv2) CreateContext(ctx context.Context, r *schema.ResourceData, m interface{}) diag.Diagnostics {
	// TODO: Add analytics
	return a.resource.CreateContext(ctx, r, m)
}

type ProviderMetaSDKv2 struct {
	UserAgentExtra map[string]string `cty:"user_agent_extra"`
	ModuleName     *string           `cty:"module_name"`
	ModuleVersion  *string           `cty:"module_version"`
}

// See Resource documentation.
func (a *AnalyticsResourceSDKv2) ReadContext(ctx context.Context, r *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := ProviderMetaSDKv2{}
	err := r.GetProviderMeta(&meta)
	if err != nil {
		log.Printf("[WARN] failed to decode provider meta: %s, meta: %v", err, meta)
		return a.resource.ReadContext(ctx, r, m)
	}
	moduleName := ""
	if meta.ModuleName != nil {
		moduleName = *meta.ModuleName
	}
	moduleVersion := ""
	if meta.ModuleVersion != nil {
		moduleVersion = *meta.ModuleVersion
	}

	uaExtra := UserAgentExtra{
		Name:          a.name,
		Operation:     UserAgentOperationValueRead,
		Extras:        meta.UserAgentExtra,
		ModuleName:    moduleName,
		ModuleVersion: moduleVersion,
	}
	ctx = AddUserAgentExtra(ctx, uaExtra)
	return a.resource.ReadContext(ctx, r, m)
}

// See Resource documentation.
func (a *AnalyticsResourceSDKv2) UpdateContext(ctx context.Context, r *schema.ResourceData, m interface{}) diag.Diagnostics {
	return a.resource.UpdateContext(ctx, r, m)
}

// See Resource documentation.
func (a *AnalyticsResourceSDKv2) DeleteContext(ctx context.Context, r *schema.ResourceData, m interface{}) diag.Diagnostics {
	return a.resource.DeleteContext(ctx, r, m)
}

func (a *AnalyticsResourceSDKv2) resourceImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	return a.resource.Importer.StateContext(ctx, d, meta)
}
