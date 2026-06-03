package organization3

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/atlas-sdk/v20250312020/admin"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ resource.ResourceWithConfigure = &organization3RS{}
var _ resource.ResourceWithModifyPlan = &organization3RS{}

const (
	resourceName     = "organization3"
	fullResourceName = "mongodbatlas_" + resourceName
)

func Resource() resource.Resource {
	return &organization3RS{
		RSCommon: config.RSCommon{
			ResourceName: resourceName,
		},
	}
}

type organization3RS struct {
	config.RSCommon
}

func (r *organization3RS) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
	conversion.UpdateSchemaDescription(&resp.Schema)
}

func (r *organization3RS) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
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

	targetVersion, shouldRotate := rotationTargetVersion(ctx, &planRotation, &stateRotation, stateVersion, time.Now())
	if !shouldRotate {
		return
	}

	planRotation.SecretVersion = types.Int64Value(targetVersion)
	promotedOld, promoteDiags := secretMetadataForRotationPromotion(ctx, &planRotation, &stateRotation)
	resp.Diagnostics.Append(promoteDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	oldSecretObject, oldDiags := secretMetadataToObject(ctx, &promotedOld)
	resp.Diagnostics.Append(oldDiags...)
	planRotation.OldSecret = oldSecretObject
	planRotation.CurrentSecret = types.ObjectUnknown(secretMetadataObjectType.AttrTypes)
	plan.ClientSecret = types.StringUnknown()

	rotationObject, objectDiags := rotationToObject(ctx, &planRotation)
	resp.Diagnostics.Append(objectDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ClientSecretRotation = rotationObject
	resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
}

func resolveRotationPlanVersion(
	ctx context.Context,
	plan *TFModel,
	planRotation, stateRotation *TFClientSecretRotationModel,
	stateVersion int64,
) int64 {
	planVersion := stateVersion
	if !planRotation.SecretVersion.IsNull() && !planRotation.SecretVersion.IsUnknown() {
		planVersion = planRotation.SecretVersion.ValueInt64()
	}
	if planVersion > stateVersion {
		return planVersion
	}
	targetVersion, shouldRotate := rotationTargetVersion(ctx, planRotation, stateRotation, stateVersion, time.Now())
	if shouldRotate {
		return targetVersion
	}
	if planRotation.CurrentSecret.IsUnknown() && plan.ClientSecret.IsUnknown() {
		return stateVersion + 1
	}
	return stateVersion
}

func fillRotationSecretsFromStateIfUnknown(planRotation, stateRotation *TFClientSecretRotationModel) {
	if planRotation.CurrentSecret.IsUnknown() && !stateRotation.CurrentSecret.IsUnknown() && !stateRotation.CurrentSecret.IsNull() {
		planRotation.CurrentSecret = stateRotation.CurrentSecret
	}
	if planRotation.OldSecret.IsUnknown() && !stateRotation.OldSecret.IsUnknown() && !stateRotation.OldSecret.IsNull() {
		planRotation.OldSecret = stateRotation.OldSecret
	}
}

func (r *organization3RS) finalizeRotationState(
	ctx context.Context,
	conn *admin.APIClient,
	state *TFModel,
	rotation *TFClientSecretRotationModel,
) diag.Diagnostics {
	var diags diag.Diagnostics
	sa, err := getServiceAccount(ctx, conn, state.OrgID.ValueString(), state.ClientID.ValueString())
	if err != nil {
		diags.AddError(errorUpdate, err.Error())
		return diags
	}
	if sa == nil {
		diags.AddError(errorUpdate, "service account not found during update")
		return diags
	}
	updated, refreshDiags := refreshRotationSecrets(ctx, rotation, sa.GetSecrets())
	diags.Append(refreshDiags...)
	*rotation = updated
	return diags
}

func (r *organization3RS) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	expiresAfterHours := defaultExpiresAfterHours
	var planRotation TFClientSecretRotationModel
	if !plan.ClientSecretRotation.IsNull() {
		var diags diag.Diagnostics
		planRotation, diags = RotationFromObject(ctx, plan.ClientSecretRotation)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		expiresAfterHours = planRotation.ExpiresAfterHours.ValueInt64()
	}

	conn := providerAtlasV2(r.Client)
	createResp, initialSecret, err := createOrganization(
		ctx,
		conn,
		plan.Name.ValueString(),
		plan.OrgOwnerID.ValueString(),
		int(expiresAfterHours),
	)
	if err != nil {
		resp.Diagnostics.AddError(errorCreate, err.Error())
		return
	}

	org := createResp.GetOrganization()
	sa, _ := createResp.GetServiceAccountOk()
	plaintext := ""
	if initialSecret.HasSecret() {
		plaintext = initialSecret.GetSecret()
	}

	state := plan
	state.OrgID = types.StringValue(org.GetId())
	state.ClientID = types.StringValue(sa.GetClientId())
	state.ClientSecret = types.StringValue(plaintext)

	if !plan.ClientSecretRotation.IsNull() {
		rotationWithDefaults(&planRotation, expiresAfterHours)
		currentMeta := secretMetadataFromAPI(initialSecret)
		currentObject, currentDiags := secretMetadataToObject(ctx, &currentMeta)
		resp.Diagnostics.Append(currentDiags...)
		oldObject, oldDiags := secretMetadataNullObject(ctx)
		resp.Diagnostics.Append(oldDiags...)
		planRotation.SecretVersion = types.Int64Value(1)
		planRotation.CurrentSecret = currentObject
		planRotation.OldSecret = oldObject
		rotationObject, rotDiags := rotationToObject(ctx, &planRotation)
		resp.Diagnostics.Append(rotDiags...)
		state.ClientSecretRotation = rotationObject
	} else {
		state.ClientSecretRotation = types.ObjectNull(rotationObjectType.AttrTypes)
	}

	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(validateRotationStateForWrite(ctx, state.ClientSecretRotation)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *organization3RS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgID := state.OrgID.ValueString()
	conn := r.atlasV2(ctx, &state)

	org, err := getOrganization(ctx, conn, orgID)
	if err != nil {
		resp.Diagnostics.AddError(errorRead, err.Error())
		return
	}
	if org == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(org.GetName())

	if state.ClientSecretRotation.IsNull() {
		resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
		return
	}

	stateRotation, diags := RotationFromObject(ctx, state.ClientSecretRotation)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sa, err := getServiceAccount(ctx, conn, orgID, state.ClientID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(errorRead, err.Error())
		return
	}
	if sa == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	secrets := sa.GetSecrets()

	stateRotation, readDiags := refreshRotationSecrets(ctx, &stateRotation, secrets)
	resp.Diagnostics.Append(readDiags...)

	rotationObject, rotDiags := rotationToObject(ctx, &stateRotation)
	resp.Diagnostics.Append(rotDiags...)
	state.ClientSecretRotation = rotationObject

	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func refreshRotationSecrets(ctx context.Context, rotation *TFClientSecretRotationModel, apiSecrets []admin.ServiceAccountSecret) (TFClientSecretRotationModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	rotationWithDefaults(rotation, rotation.ExpiresAfterHours.ValueInt64())

	if rotation.CurrentSecret.IsUnknown() {
		return *rotation, diags
	}
	currentMeta, currentDiags := SecretMetadataFromObject(ctx, rotation.CurrentSecret)
	diags.Append(currentDiags...)
	if !currentMeta.SecretID.IsNull() && currentMeta.SecretID.ValueString() != "" {
		if apiSecret, ok := findSecretByID(apiSecrets, currentMeta.SecretID.ValueString()); ok {
			merged := mergeSecretMetadataFromAPI(&currentMeta, apiSecret)
			currentObject, objDiags := secretMetadataToObject(ctx, &merged)
			diags.Append(objDiags...)
			rotation.CurrentSecret = currentObject
		}
	}

	if rotation.OldSecret.IsNull() || rotation.OldSecret.IsUnknown() {
		return *rotation, diags
	}
	oldMeta, oldDiags := SecretMetadataFromObject(ctx, rotation.OldSecret)
	diags.Append(oldDiags...)
	if oldMeta.SecretID.IsNull() || oldMeta.SecretID.ValueString() == "" {
		nullOld, nullDiags := secretMetadataNullObject(ctx)
		diags.Append(nullDiags...)
		rotation.OldSecret = nullOld
		return *rotation, diags
	}
	if apiSecret, ok := findSecretByID(apiSecrets, oldMeta.SecretID.ValueString()); ok {
		merged := mergeSecretMetadataFromAPI(&oldMeta, apiSecret)
		oldObject, objDiags := secretMetadataToObject(ctx, &merged)
		diags.Append(objDiags...)
		rotation.OldSecret = oldObject
	} else {
		nullOld, nullDiags := secretMetadataNullObject(ctx)
		diags.Append(nullDiags...)
		rotation.OldSecret = nullOld
	}
	return *rotation, diags
}

func (r *organization3RS) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state TFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	conn := r.atlasV2(ctx, &state)

	if !plan.Name.Equal(state.Name) {
		if err := updateOrganizationName(ctx, conn, state.OrgID.ValueString(), plan.Name.ValueString()); err != nil {
			resp.Diagnostics.AddError(errorUpdate, err.Error())
			return
		}
	}

	if plan.ClientSecretRotation.IsNull() {
		state.Name = plan.Name
		state.ClientSecretRotation = types.ObjectNull(rotationObjectType.AttrTypes)
		resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
		return
	}

	planRotation, diags := RotationFromObject(ctx, plan.ClientSecretRotation)
	resp.Diagnostics.Append(diags...)
	stateRotation, stateDiags := RotationFromObject(ctx, state.ClientSecretRotation)
	resp.Diagnostics.Append(stateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	fillRotationSecretsFromStateIfUnknown(&planRotation, &stateRotation)

	stateVersion := int64(0)
	if !stateRotation.SecretVersion.IsNull() && !stateRotation.SecretVersion.IsUnknown() {
		stateVersion = stateRotation.SecretVersion.ValueInt64()
	}

	planVersion := resolveRotationPlanVersion(ctx, &plan, &planRotation, &stateRotation, stateVersion)
	planRotation.SecretVersion = types.Int64Value(planVersion)

	if planVersion > stateVersion {
		if err := r.applyRotation(
			ctx,
			conn,
			&plan,
			&state,
			&planRotation,
			&stateRotation,
			stateVersion,
			planVersion,
		); err != nil {
			resp.Diagnostics.AddError(errorUpdate, err.Error())
			return
		}
	}

	fillRotationSecretsFromStateIfUnknown(&planRotation, &stateRotation)
	resp.Diagnostics.Append(r.finalizeRotationState(ctx, conn, &state, &planRotation)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Name = plan.Name
	rotationObject, rotDiags := rotationToObject(ctx, &planRotation)
	resp.Diagnostics.Append(rotDiags...)
	state.ClientSecretRotation = rotationObject
	state.ClientSecret = plan.ClientSecret

	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(validateRotationStateForWrite(ctx, state.ClientSecretRotation)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func deleteOldSecretBeforeRotation(
	ctx context.Context,
	conn *admin.APIClient,
	orgID, clientID string,
	stateRotation *TFClientSecretRotationModel,
	stateVersion int64,
) error {
	if stateVersion < 2 {
		return nil
	}
	stateOld, diags := SecretMetadataFromObject(ctx, stateRotation.OldSecret)
	if diags.HasError() {
		return fmt.Errorf("invalid old_secret in state")
	}
	if !shouldDeleteOldSecret(stateVersion, &stateOld) {
		return fmt.Errorf(
			"cannot rotate at secret_version %d: old_secret.secret_id must be set in state so the overlap secret managed by this resource can be deleted before creating the next secret; run terraform refresh if Atlas already has two secrets",
			stateVersion,
		)
	}
	if err := deleteServiceAccountSecret(ctx, conn, orgID, clientID, stateOld.SecretID.ValueString()); err != nil {
		return fmt.Errorf("delete old secret before rotation: %w", err)
	}
	return nil
}

func (r *organization3RS) applyRotation(
	ctx context.Context,
	conn *admin.APIClient,
	plan, state *TFModel,
	planRotation, stateRotation *TFClientSecretRotationModel,
	stateVersion, planVersion int64,
) error {
	orgID := state.OrgID.ValueString()
	clientID := state.ClientID.ValueString()
	if err := deleteOldSecretBeforeRotation(ctx, conn, orgID, clientID, stateRotation, stateVersion); err != nil {
		return err
	}

	expiresAfter := planRotation.ExpiresAfterHours.ValueInt64()
	newSecret, err := createServiceAccountSecret(ctx, conn, orgID, clientID, int(expiresAfter))
	if err != nil {
		return fmt.Errorf("create rotated secret: %w", err)
	}

	plaintext := ""
	if newSecret.HasSecret() {
		plaintext = newSecret.GetSecret()
	}

	oldObject := stateRotation.CurrentSecret
	currentMeta := secretMetadataFromAPI(newSecret)
	currentObject, currentDiags := secretMetadataToObject(ctx, &currentMeta)
	if currentDiags.HasError() {
		return fmt.Errorf("build current_secret")
	}

	planRotation.OldSecret = oldObject
	planRotation.CurrentSecret = currentObject
	planRotation.SecretVersion = types.Int64Value(planVersion)
	rotationWithDefaults(planRotation, expiresAfter)
	plan.ClientSecret = types.StringValue(plaintext)
	state.ClientSecret = plan.ClientSecret
	return nil
}

func (r *organization3RS) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	conn := providerAtlasV2(r.Client)
	if err := deleteOrganization(ctx, conn, state.OrgID.ValueString()); err != nil {
		resp.Diagnostics.AddError(errorUpdate, fmt.Sprintf("delete organization: %s", err))
	}
}

const (
	errorCreate = "error creating resource " + fullResourceName
	errorRead   = "error reading resource " + fullResourceName
	errorUpdate = "error updating resource " + fullResourceName
)
