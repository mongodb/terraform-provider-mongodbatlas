package organization2

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ resource.ResourceWithConfigure = &organization2RS{}
var _ resource.ResourceWithImportState = &organization2RS{}
var _ resource.ResourceWithModifyPlan = &organization2RS{}

const (
	resourceName     = "organization2"
	fullResourceName = "mongodbatlas_" + resourceName
)

func Resource() resource.Resource {
	return &organization2RS{
		RSCommon: config.RSCommon{
			ResourceName: resourceName,
		},
	}
}

type organization2RS struct {
	config.RSCommon
}

func (r *organization2RS) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
	conversion.UpdateSchemaDescription(&resp.Schema)
}

func (r *organization2RS) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	var plan, state TFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ClientSecretRotation.IsNull() {
		return
	}

	planRotation, diags := RotationFromObject(ctx, plan.ClientSecretRotation)
	resp.Diagnostics.Append(diags...)
	stateRotation, stateDiags := RotationFromObject(ctx, state.ClientSecretRotation)
	resp.Diagnostics.Append(stateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	stateVersion := int64(0)
	if !stateRotation.SecretVersion.IsNull() {
		stateVersion = stateRotation.SecretVersion.ValueInt64()
	}

	targetVersion, shouldRotate := rotationTargetVersion(&planRotation, &stateRotation, stateVersion)
	if !shouldRotate {
		return
	}

	planRotation.SecretVersion = types.Int64Value(targetVersion)
	if !stateRotation.CurrentSecretID.IsNull() && stateRotation.CurrentSecretID.ValueString() != "" {
		planRotation.OldSecretID = stateRotation.CurrentSecretID
	} else {
		planRotation.OldSecretID = types.StringNull()
	}
	planRotation.CurrentSecretID = types.StringUnknown()
	planRotation.NextRenewal = types.StringUnknown()
	planRotation.ExpiresAt = types.StringUnknown()
	plan.ClientSecret = types.StringUnknown()

	rotationObject, objectDiags := rotationToObject(ctx, &planRotation)
	resp.Diagnostics.Append(objectDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ClientSecretRotation = rotationObject
	resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
}

func rotationTargetVersion(planRotation, stateRotation *TFClientSecretRotationModel, stateVersion int64) (int64, bool) {
	if !planRotation.SecretVersion.IsUnknown() && !planRotation.SecretVersion.IsNull() {
		configVersion := planRotation.SecretVersion.ValueInt64()
		if configVersion > stateVersion {
			return configVersion, true
		}
	}

	if stateRotation.NextRenewal.IsNull() || stateRotation.ExpiresAt.IsNull() {
		return 0, false
	}

	nextRenewal, err := time.Parse(time.RFC3339, stateRotation.NextRenewal.ValueString())
	if err != nil {
		return 0, false
	}
	expiresAt, err := time.Parse(time.RFC3339, stateRotation.ExpiresAt.ValueString())
	if err != nil {
		return 0, false
	}
	if !RotationDue(time.Now(), nextRenewal, expiresAt) {
		return 0, false
	}
	return stateVersion + 1, true
}

func (r *organization2RS) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	now := time.Now()
	name := plan.Name.ValueString()
	clientID := mintClientID(name)
	clientSecret, currentSecretID := mintCredentials(name, 1)

	state := &orgState{
		name:            name,
		orgID:           mintOrgID(name),
		clientID:        clientID,
		clientSecret:    clientSecret,
		currentSecretID: currentSecretID,
		secretCreatedAt: now,
	}

	if !plan.ClientSecretRotation.IsNull() {
		rotation, diags := RotationFromObject(ctx, plan.ClientSecretRotation)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		intervalDuration, err := parseInterval(rotation.Interval.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(errorCreate, fmt.Sprintf("invalid interval: %s", err))
			return
		}

		nextRenewal, expiresAt := computeRenewalTimes(now, intervalDuration)
		state.hasRotationBlock = true
		state.interval = rotation.Interval.ValueString()
		state.secretVersion = 1
		state.nextRenewal = nextRenewal
		state.expiresAt = expiresAt
	}

	if err := putStoreEntry(name, state); err != nil {
		resp.Diagnostics.AddError(errorCreate, fmt.Sprintf("failed to persist mock store: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, stateToTFModel(ctx, state))...)
}

func (r *organization2RS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entry, ok := getStoreEntry(state.Name.ValueString())
	if !ok {
		resp.State.RemoveResource(ctx)
		return
	}

	model := stateToTFModel(ctx, entry)
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *organization2RS) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.ValueString()

	err := updateStoreEntry(name, func(entry *orgState) error {
		entry.name = name

		if plan.ClientSecretRotation.IsNull() {
			entry.hasRotationBlock = false
			entry.interval = ""
			return nil
		}

		planRotation, diags := RotationFromObject(ctx, plan.ClientSecretRotation)
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return fmt.Errorf("invalid rotation object")
		}

		if !entry.hasRotationBlock {
			intervalDuration, parseErr := parseInterval(planRotation.Interval.ValueString())
			if parseErr != nil {
				return parseErr
			}
			now := time.Now()
			nextRenewal, expiresAt := computeRenewalTimes(now, intervalDuration)
			entry.hasRotationBlock = true
			entry.interval = planRotation.Interval.ValueString()
			entry.secretVersion = 1
			entry.secretCreatedAt = now
			entry.nextRenewal = nextRenewal
			entry.expiresAt = expiresAt
			entry.clientID = mintClientID(name)
			entry.clientSecret, entry.currentSecretID = mintCredentials(name, 1)
			entry.oldSecretID = ""
			return nil
		}

		entry.interval = planRotation.Interval.ValueString()

		planVersion := entry.secretVersion
		if !planRotation.SecretVersion.IsNull() && !planRotation.SecretVersion.IsUnknown() {
			planVersion = planRotation.SecretVersion.ValueInt64()
		}

		if planVersion > entry.secretVersion {
			entry.oldSecretID = entry.currentSecretID
			entry.clientSecret, entry.currentSecretID = mintCredentials(name, planVersion)
			entry.secretVersion = planVersion
			now := time.Now()
			entry.secretCreatedAt = now
			intervalDuration, parseErr := parseInterval(entry.interval)
			if parseErr != nil {
				return parseErr
			}
			entry.nextRenewal, entry.expiresAt = computeRenewalTimes(now, intervalDuration)
		}
		return nil
	})
	if err != nil {
		if resp.Diagnostics.HasError() {
			return
		}
		resp.Diagnostics.AddError(errorUpdate, err.Error())
		return
	}

	entry, ok := getStoreEntry(name)
	if !ok {
		resp.Diagnostics.AddError(errorRead, fmt.Sprintf("organization2 %q not found", name))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, stateToTFModel(ctx, entry))...)
}

func (r *organization2RS) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := deleteStoreEntry(state.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError(errorUpdate, fmt.Sprintf("failed to persist mock store: %s", err))
	}
}

func (r *organization2RS) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), req.ID)...)
}

const (
	errorCreate = "error creating resource " + fullResourceName
	errorRead   = "error reading resource " + fullResourceName
	errorUpdate = "error updating resource " + fullResourceName
)

var rotationObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"interval":          types.StringType,
	"secret_version":    types.Int64Type,
	"next_renewal":      types.StringType,
	"expires_at":        types.StringType,
	"current_secret_id": types.StringType,
	"old_secret_id":     types.StringType,
}}

func RotationFromObject(ctx context.Context, object types.Object) (TFClientSecretRotationModel, diag.Diagnostics) {
	if object.IsNull() {
		return TFClientSecretRotationModel{}, nil
	}
	var rotation TFClientSecretRotationModel
	diags := object.As(ctx, &rotation, basetypes.ObjectAsOptions{})
	return rotation, diags
}

func rotationToObject(ctx context.Context, rotation *TFClientSecretRotationModel) (types.Object, diag.Diagnostics) {
	return types.ObjectValueFrom(ctx, rotationObjectType.AttrTypes, rotation)
}

func stateToTFModel(ctx context.Context, state *orgState) TFModel {
	model := TFModel{
		Name:         types.StringValue(state.name),
		OrgID:        types.StringValue(state.orgID),
		ClientID:     types.StringValue(state.clientID),
		ClientSecret: types.StringValue(state.clientSecret),
	}

	if !state.hasRotationBlock {
		model.ClientSecretRotation = types.ObjectNull(rotationObjectType.AttrTypes)
		return model
	}

	oldSecretID := types.StringNull()
	if state.oldSecretID != "" {
		oldSecretID = types.StringValue(state.oldSecretID)
	}

	rotation := TFClientSecretRotationModel{
		Interval:        types.StringValue(state.interval),
		SecretVersion:   types.Int64Value(state.secretVersion),
		NextRenewal:     types.StringValue(formatRFC3339(state.nextRenewal)),
		ExpiresAt:       types.StringValue(formatRFC3339(state.expiresAt)),
		CurrentSecretID: types.StringValue(state.currentSecretID),
		OldSecretID:     oldSecretID,
	}
	object, _ := rotationToObject(ctx, &rotation)
	model.ClientSecretRotation = object
	return model
}

func mintOrgID(name string) string {
	return "org-" + randomHex(12) + "-" + name
}

func mintCredentials(name string, version int64) (clientSecret, secretID string) {
	suffix := randomHex(6)
	clientSecret = fmt.Sprintf("secret-%s-v%d-%s", name, version, suffix)
	secretID = fmt.Sprintf("sid-%s-v%d-%s", name, version, suffix)
	return clientSecret, secretID
}

func mintClientID(name string) string {
	return fmt.Sprintf("client-%s-%s", name, randomHex(6))
}

func randomHex(byteCount int) string {
	buf := make([]byte, byteCount)
	if _, err := rand.Read(buf); err != nil {
		panic(fmt.Sprintf("failed to read random bytes: %s", err))
	}
	return hex.EncodeToString(buf)
}

// ResetStoreForTest clears the mock store. Tests only.
func ResetStoreForTest() {
	resetStoreLocked()
}

// HasStoreEntry reports whether name exists in the mock store. Tests only.
func HasStoreEntry(name string) bool {
	_, ok := getStoreEntry(name)
	return ok
}
