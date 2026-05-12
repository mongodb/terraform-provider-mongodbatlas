package apiresource

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TFModel is the in-memory representation of mongodbatlas_api_resource.
type TFModel struct {
	IDAttribute        types.List    `tfsdk:"id_attribute"`
	CreateOnlyBodyKeys types.Set     `tfsdk:"create_only_body_keys"`
	ID                 types.String  `tfsdk:"id"`
	Path               types.String  `tfsdk:"path"`
	CreateMethod       types.String  `tfsdk:"create_method"`
	UpdateMethod       types.String  `tfsdk:"update_method"`
	VersionHeader      types.String  `tfsdk:"version_header"`
	Body               types.Dynamic `tfsdk:"body"`
	SensitiveBody      types.Dynamic `tfsdk:"sensitive_body"`
	Output             types.Dynamic `tfsdk:"output"`
	Preview            types.Bool    `tfsdk:"preview"`
}

// TFModelDS is the data source equivalent: read-only, no body / sensitive_body.
type TFModelDS struct {
	ID            types.String  `tfsdk:"id"`
	Path          types.String  `tfsdk:"path"`
	VersionHeader types.String  `tfsdk:"version_header"`
	Output        types.Dynamic `tfsdk:"output"`
	Preview       types.Bool    `tfsdk:"preview"`
}

const (
	previewVersionHeader = "application/vnd.atlas.preview+json"
	defaultCreateMethod  = "POST"
	defaultReadMethod    = "GET"
	defaultUpdateMethod  = "PATCH"
	defaultDeleteMethod  = "DELETE"
)

// TFModelUpdate is the in-memory representation of mongodbatlas_api_update.
// Subset of TFModel: no id_attribute (path IS the read URL), no create_method
// (we never create), no create_only_body_keys.
type TFModelUpdate struct {
	ID            types.String  `tfsdk:"id"`
	Path          types.String  `tfsdk:"path"`
	UpdateMethod  types.String  `tfsdk:"update_method"`
	VersionHeader types.String  `tfsdk:"version_header"`
	Body          types.Dynamic `tfsdk:"body"`
	SensitiveBody types.Dynamic `tfsdk:"sensitive_body"`
	Output        types.Dynamic `tfsdk:"output"`
	Preview       types.Bool    `tfsdk:"preview"`
}

// todayVersionHeader returns the Atlas version media type for today's date in
// UTC. Atlas snaps the date down to the latest published API version on or
// before it, so this effectively means "use the latest version Atlas has".
// Captured at resource Create time and persisted in state so subsequent
// operations on the same resource use the same version for the lifetime of
// the resource.
func todayVersionHeader() string {
	return "application/vnd.atlas." + time.Now().UTC().Format("2006-01-02") + "+json"
}
