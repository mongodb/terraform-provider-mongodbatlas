package organization3_test

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/organization3"
)

func TestRenewalDue(t *testing.T) {
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)
	expiresAt := now.Add(720 * time.Hour)

	t.Run("before renew window", func(t *testing.T) {
		assert.False(t, organization3.RenewalDue(now, expiresAt, 360))
	})

	t.Run("inside renew window", func(t *testing.T) {
		renewAt := expiresAt.Add(-360 * time.Hour)
		assert.True(t, organization3.RenewalDue(renewAt, expiresAt, 360))
	})

	t.Run("at expiry", func(t *testing.T) {
		assert.True(t, organization3.RenewalDue(expiresAt, expiresAt, 360))
	})

	t.Run("after expiry", func(t *testing.T) {
		assert.True(t, organization3.RenewalDue(expiresAt.Add(time.Hour), expiresAt, 360))
	})
}

func TestEffectiveRotateBeforeExpiryHours(t *testing.T) {
	assert.Equal(t, int64(360), organization3.EffectiveRotateBeforeExpiryHoursForTest(720, types.Int64Null()))
	assert.Equal(t, int64(100), organization3.EffectiveRotateBeforeExpiryHoursForTest(720, types.Int64Value(100)))
}

func TestShouldDeleteOldSecret(t *testing.T) {
	empty := &organization3.TFSecretMetadataModel{SecretID: types.StringNull()}
	withID := &organization3.TFSecretMetadataModel{SecretID: types.StringValue("sid-old")}

	assert.False(t, organization3.ShouldDeleteOldSecretForTest(1, empty))
	assert.False(t, organization3.ShouldDeleteOldSecretForTest(2, empty))
	assert.False(t, organization3.ShouldDeleteOldSecretForTest(1, withID))
	assert.True(t, organization3.ShouldDeleteOldSecretForTest(2, withID))
	assert.True(t, organization3.ShouldDeleteOldSecretForTest(3, withID))
}

func TestModifyPlan_noRotationBlock(t *testing.T) {
	plan, state := buildPlanAndState(t, nil, nil)
	req := resource.ModifyPlanRequest{Plan: plan, State: state}
	resp := &resource.ModifyPlanResponse{Plan: plan}

	modifyPlan(t, req, resp)
	assert.False(t, resp.Diagnostics.HasError())
	assertPlanRotationUnchanged(t, resp.Plan)
}

func TestModifyPlan_blockNotDue(t *testing.T) {
	now := time.Now().UTC()
	expiresAt := now.Add(720 * time.Hour)
	rotation := rotationValues{
		expiresAfterHours: 720,
		rotateBefore:      360,
		secretVersion:     1,
		currentSecret: &secretMetadataValues{
			secretID:  "sid-current",
			createdAt: now.Format(time.RFC3339),
			expiresAt: expiresAt.Format(time.RFC3339),
		},
	}
	plan, state := buildPlanAndState(t, &rotation, &rotation)
	req := resource.ModifyPlanRequest{Plan: plan, State: state}
	resp := &resource.ModifyPlanResponse{Plan: plan}

	modifyPlan(t, req, resp)
	assert.False(t, resp.Diagnostics.HasError())
	assertPlanRotationUnchanged(t, resp.Plan)
}

func TestModifyPlan_blockDue(t *testing.T) {
	now := time.Now().UTC()
	expiresAt := now.Add(24 * time.Hour)
	rotation := rotationValues{
		expiresAfterHours: 720,
		rotateBefore:      87600,
		secretVersion:     1,
		currentSecret: &secretMetadataValues{
			secretID:  "sid-current",
			createdAt: now.Add(-48 * time.Hour).Format(time.RFC3339),
			expiresAt: expiresAt.Format(time.RFC3339),
		},
	}
	plan, state := buildPlanAndState(t, &rotation, &rotation)
	req := resource.ModifyPlanRequest{Plan: plan, State: state}
	resp := &resource.ModifyPlanResponse{Plan: plan}

	modifyPlan(t, req, resp)
	require.False(t, resp.Diagnostics.HasError())

	var planModel organization3.TFModel
	resp.Diagnostics.Append(resp.Plan.Get(t.Context(), &planModel)...)
	require.False(t, resp.Diagnostics.HasError())

	rotationModel, diags := organization3.RotationFromObject(t.Context(), planModel.ClientSecretRotation)
	require.False(t, diags.HasError())
	assert.Equal(t, int64(2), rotationModel.SecretVersion.ValueInt64())
	oldMeta, oldDiags := organization3.SecretMetadataFromObject(t.Context(), rotationModel.OldSecret)
	require.False(t, oldDiags.HasError())
	assert.Equal(t, "sid-current", oldMeta.SecretID.ValueString())
	assert.True(t, rotationModel.CurrentSecret.IsUnknown())
	assert.True(t, planModel.ClientSecret.IsUnknown())
}

func TestModifyPlan_forceSecretVersion(t *testing.T) {
	now := time.Now().UTC()
	expiresAt := now.Add(720 * time.Hour)
	stateRotation := rotationValues{
		expiresAfterHours: 720,
		rotateBefore:      360,
		secretVersion:     1,
		currentSecret: &secretMetadataValues{
			secretID:  "sid-before-force",
			createdAt: now.Format(time.RFC3339),
			expiresAt: expiresAt.Format(time.RFC3339),
		},
	}
	planRotation := stateRotation
	planRotation.secretVersion = 2
	plan, state := buildPlanAndState(t, &planRotation, &stateRotation)
	req := resource.ModifyPlanRequest{Plan: plan, State: state}
	resp := &resource.ModifyPlanResponse{Plan: plan}

	modifyPlan(t, req, resp)
	require.False(t, resp.Diagnostics.HasError())

	var planModel organization3.TFModel
	resp.Diagnostics.Append(resp.Plan.Get(t.Context(), &planModel)...)
	require.False(t, resp.Diagnostics.HasError())

	rotationModel, diags := organization3.RotationFromObject(t.Context(), planModel.ClientSecretRotation)
	require.False(t, diags.HasError())
	assert.Equal(t, int64(2), rotationModel.SecretVersion.ValueInt64())
	assert.True(t, planModel.ClientSecret.IsUnknown())
}

type secretMetadataValues struct {
	lastUsedAt *string
	secretID   string
	createdAt  string
	expiresAt  string
}

type rotationValues struct {
	oldSecret         *secretMetadataValues
	currentSecret     *secretMetadataValues
	expiresAfterHours int64
	rotateBefore      int64
	secretVersion     int64
}

func buildPlanAndState(t *testing.T, planRotation, stateRotation *rotationValues) (tfsdk.Plan, tfsdk.State) {
	t.Helper()
	ctx := t.Context()
	schema := organization3.ResourceSchema(ctx)

	buildValue := func(rotation *rotationValues) tftypes.Value {
		attrs := map[string]tftypes.Value{
			"name":          tftypes.NewValue(tftypes.String, "test-org"),
			"org_owner_id":  tftypes.NewValue(tftypes.String, "owner-id"),
			"org_id":        tftypes.NewValue(tftypes.String, "org-test"),
			"client_id":     tftypes.NewValue(tftypes.String, "client-test"),
			"client_secret": tftypes.NewValue(tftypes.String, "secret-test"),
		}
		if rotation == nil {
			attrs["client_secret_rotation"] = tftypes.NewValue(rotationObjectTerraformType(), nil)
		} else {
			attrs["client_secret_rotation"] = tftypes.NewValue(rotationObjectTerraformType(), rotationTerraformMap(rotation))
		}
		return tftypes.NewValue(schema.Type().TerraformType(ctx), attrs)
	}

	planRaw := buildValue(planRotation)
	stateRaw := buildValue(stateRotation)
	return tfsdk.Plan{Schema: schema, Raw: planRaw}, tfsdk.State{Schema: schema, Raw: stateRaw}
}

func rotationObjectTerraformType() tftypes.Type {
	return tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"expires_after_hours":        tftypes.Number,
		"rotate_before_expiry_hours": tftypes.Number,
		"secret_version":             tftypes.Number,
		"current_secret":             secretMetadataTerraformType(),
		"old_secret":                 secretMetadataTerraformType(),
	}}
}

func secretMetadataTerraformType() tftypes.Type {
	return tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"secret_id":    tftypes.String,
		"created_at":   tftypes.String,
		"expires_at":   tftypes.String,
		"last_used_at": tftypes.String,
	}}
}

func secretMetadataTerraformMap(meta secretMetadataValues) map[string]tftypes.Value {
	lastUsed := any(nil)
	if meta.lastUsedAt != nil {
		lastUsed = *meta.lastUsedAt
	}
	return map[string]tftypes.Value{
		"secret_id":    tftypes.NewValue(tftypes.String, meta.secretID),
		"created_at":   tftypes.NewValue(tftypes.String, meta.createdAt),
		"expires_at":   tftypes.NewValue(tftypes.String, meta.expiresAt),
		"last_used_at": tftypes.NewValue(tftypes.String, lastUsed),
	}
}

func rotationTerraformMap(rotation *rotationValues) map[string]tftypes.Value {
	oldSecret := any(nil)
	if rotation.oldSecret != nil {
		oldSecret = secretMetadataTerraformMap(*rotation.oldSecret)
	}
	return map[string]tftypes.Value{
		"expires_after_hours":        tftypes.NewValue(tftypes.Number, float64(rotation.expiresAfterHours)),
		"rotate_before_expiry_hours": tftypes.NewValue(tftypes.Number, float64(rotation.rotateBefore)),
		"secret_version":             tftypes.NewValue(tftypes.Number, float64(rotation.secretVersion)),
		"current_secret":             tftypes.NewValue(secretMetadataTerraformType(), secretMetadataTerraformMap(*rotation.currentSecret)),
		"old_secret":                 tftypes.NewValue(secretMetadataTerraformType(), oldSecret),
	}
}

func modifyPlan(t *testing.T, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	t.Helper()
	rs := organization3.Resource()
	rsMp, ok := rs.(resource.ResourceWithModifyPlan)
	require.True(t, ok)
	rsMp.ModifyPlan(t.Context(), req, resp)
}

func assertPlanRotationUnchanged(t *testing.T, plan tfsdk.Plan) {
	t.Helper()
	var planModel organization3.TFModel
	diags := plan.Get(t.Context(), &planModel)
	require.False(t, diags.HasError())
	require.False(t, planModel.ClientSecret.IsUnknown())
}
