package config

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func NewAnalyticsResource(d *schema.Resource) *schema.Resource {
	analyticsResource := AnalyticsResourceSDKv2{
		resource: d,
	}
	return &schema.Resource{
		// Original field mapping
		Schema:     d.Schema,
		SchemaFunc: d.SchemaFunc,
		// Overriding the CRUDI methods
		CreateContext: analyticsResource.CreateContext,
		ReadContext:   analyticsResource.ReadContext,
		UpdateContext: analyticsResource.UpdateContext,
		DeleteContext: analyticsResource.DeleteContext,
	}
}

type AnalyticsResourceSDKv2 struct {
	resource *schema.Resource
}

func (a *AnalyticsResourceSDKv2) CreateContext(ctx context.Context, r *schema.ResourceData, m interface{}) diag.Diagnostics {
	// TODO: Add analytics
	return a.resource.CreateContext(ctx, r, m)
}

// See Resource documentation.
func (a *AnalyticsResourceSDKv2) ReadContext(ctx context.Context, r *schema.ResourceData, m interface{}) diag.Diagnostics {
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
