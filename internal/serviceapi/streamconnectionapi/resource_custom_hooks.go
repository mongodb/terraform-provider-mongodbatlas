package streamconnectionapi

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
)

var _ autogen.PostReadAPICallHook = (*rs)(nil)
var _ autogen.ResourceSchemaHook = (*rs)(nil)

type TFExpandedModel struct {
	Id types.String `tfsdk:"id" autogen:"omitjson"`
}

func (r *rs) ResourceSchema(ctx context.Context, baseSchema schema.Schema) schema.Schema {
	baseSchema.Attributes["id"] = schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
	return baseSchema
}

func (r *rs) PostReadAPICall(req autogen.HandleReadReq, result autogen.APICallResult) autogen.APICallResult {
	if result.Err != nil {
		return result
	}

	model, ok := req.State.(*TFModel)
	if !ok || model.WorkspaceName.IsNull() || model.ProjectId.IsNull() || model.ConnectionName.IsNull() {
		return result
	}

	craftedID := fmt.Sprintf("%s-%s-%s", model.WorkspaceName.ValueString(), model.ProjectId.ValueString(), model.ConnectionName.ValueString())

	var obj map[string]any
	if err := json.Unmarshal(result.Body, &obj); err != nil {
		return autogen.APICallResult{Body: nil, Err: err, Resp: result.Resp}
	}

	obj["id"] = craftedID

	body, err := json.Marshal(obj)
	if err != nil {
		return autogen.APICallResult{Body: nil, Err: err, Resp: result.Resp}
	}

	return autogen.APICallResult{
		Body: body,
		Err:  nil,
		Resp: result.Resp,
	}
}

func (r *rs) PreImport(id string) (string, error) {
	if strings.Contains(id, "/") {
		return id, nil
	}

	normalizedID, err := parseLegacyImportID(id)
	if err == nil {
		return normalizedID, nil
	}

	return "", fmt.Errorf("use one of the formats: {project_id}/{workspace_name}/{connection_name} or {workspace_name}-{project_id}-{connection_name}")
}

func parseLegacyImportID(id string) (string, error) {
	re := regexp.MustCompile(`^(.*)-([0-9a-fA-F]{24})-(.*)$`)
	m := re.FindStringSubmatch(id)
	if len(m) != 4 || m[1] == "" || m[3] == "" {
		return "", fmt.Errorf("invalid legacy import format")
	}

	workspaceName := m[1]
	projectID := m[2]
	connectionName := m[3]
	// Normalize to default format: {project_id}/{workspace_name}/{connection_name}
	return fmt.Sprintf("%s/%s/%s", projectID, workspaceName, connectionName), nil
}
