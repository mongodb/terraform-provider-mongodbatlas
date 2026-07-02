package organization2_test

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/organization2"
)

func TestRotationDue(t *testing.T) {
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)

	t.Run("neither due", func(t *testing.T) {
		nextRenewal := now.Add(time.Hour)
		expiresAt := now.Add(2 * time.Hour)
		assert.False(t, organization2.RotationDue(now, nextRenewal, expiresAt))
	})

	t.Run("renewal due", func(t *testing.T) {
		nextRenewal := now.Add(-time.Second)
		expiresAt := now.Add(time.Hour)
		assert.True(t, organization2.RotationDue(now, nextRenewal, expiresAt))
	})

	t.Run("expiry only due", func(t *testing.T) {
		nextRenewal := now.Add(time.Hour)
		expiresAt := now.Add(-time.Second)
		assert.True(t, organization2.RotationDue(now, nextRenewal, expiresAt))
	})
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
	nextRenewal := now.Add(time.Hour)
	expiresAt := now.Add(2 * time.Hour)
	rotation := rotationValues{
		interval:        "1h",
		secretVersion:   int64(1),
		nextRenewal:     nextRenewal.Format(time.RFC3339),
		expiresAt:       expiresAt.Format(time.RFC3339),
		currentSecretID: "sid-current",
		oldSecretID:     nil,
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
	nextRenewal := now.Add(-time.Second)
	expiresAt := now.Add(time.Minute)
	rotation := rotationValues{
		interval:        "1h",
		secretVersion:   int64(1),
		nextRenewal:     nextRenewal.Format(time.RFC3339),
		expiresAt:       expiresAt.Format(time.RFC3339),
		currentSecretID: "sid-current",
		oldSecretID:     nil,
	}
	plan, state := buildPlanAndState(t, &rotation, &rotation)
	req := resource.ModifyPlanRequest{Plan: plan, State: state}
	resp := &resource.ModifyPlanResponse{Plan: plan}

	modifyPlan(t, req, resp)
	require.False(t, resp.Diagnostics.HasError())

	var planModel organization2.TFModel
	resp.Diagnostics.Append(resp.Plan.Get(t.Context(), &planModel)...)
	require.False(t, resp.Diagnostics.HasError())

	rotationModel, diags := organization2.RotationFromObject(t.Context(), planModel.ClientSecretRotation)
	require.False(t, diags.HasError())
	assert.Equal(t, int64(2), rotationModel.SecretVersion.ValueInt64())
	assert.Equal(t, "sid-current", rotationModel.OldSecretID.ValueString())
	assert.True(t, rotationModel.CurrentSecretID.IsUnknown())
	assert.True(t, rotationModel.NextRenewal.IsUnknown())
	assert.True(t, rotationModel.ExpiresAt.IsUnknown())
	assert.True(t, planModel.ClientSecret.IsUnknown())
}

func TestModifyPlan_forceSecretVersion(t *testing.T) {
	now := time.Now().UTC()
	nextRenewal := now.Add(240 * time.Hour)
	expiresAt := now.Add(480 * time.Hour)
	stateRotation := rotationValues{
		interval:        "240h",
		secretVersion:   int64(1),
		nextRenewal:     nextRenewal.Format(time.RFC3339),
		expiresAt:       expiresAt.Format(time.RFC3339),
		currentSecretID: "sid-before-force",
		oldSecretID:     nil,
	}
	planRotation := stateRotation
	planRotation.secretVersion = int64(2)
	plan, state := buildPlanAndState(t, &planRotation, &stateRotation)
	req := resource.ModifyPlanRequest{Plan: plan, State: state}
	resp := &resource.ModifyPlanResponse{Plan: plan}

	modifyPlan(t, req, resp)
	require.False(t, resp.Diagnostics.HasError())

	var planModel organization2.TFModel
	resp.Diagnostics.Append(resp.Plan.Get(t.Context(), &planModel)...)
	require.False(t, resp.Diagnostics.HasError())

	rotationModel, diags := organization2.RotationFromObject(t.Context(), planModel.ClientSecretRotation)
	require.False(t, diags.HasError())
	assert.Equal(t, int64(2), rotationModel.SecretVersion.ValueInt64())
	assert.Equal(t, "sid-before-force", rotationModel.OldSecretID.ValueString())
	assert.True(t, rotationModel.CurrentSecretID.IsUnknown())
	assert.True(t, planModel.ClientSecret.IsUnknown())
}

type rotationValues struct {
	oldSecretID     *string
	interval        string
	nextRenewal     string
	expiresAt       string
	currentSecretID string
	secretVersion   int64
}

func buildPlanAndState(t *testing.T, planRotation, stateRotation *rotationValues) (tfsdk.Plan, tfsdk.State) {
	t.Helper()
	ctx := t.Context()
	schema := organization2.ResourceSchema(ctx)

	buildValue := func(rotation *rotationValues) tftypes.Value {
		attrs := map[string]tftypes.Value{
			"name":          tftypes.NewValue(tftypes.String, "test-org"),
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
		"interval":          tftypes.String,
		"secret_version":    tftypes.Number,
		"next_renewal":      tftypes.String,
		"expires_at":        tftypes.String,
		"current_secret_id": tftypes.String,
		"old_secret_id":     tftypes.String,
	}}
}

func rotationTerraformMap(rotation *rotationValues) map[string]tftypes.Value {
	oldSecretID := any(nil)
	if rotation.oldSecretID != nil {
		oldSecretID = *rotation.oldSecretID
	}
	return map[string]tftypes.Value{
		"interval":          tftypes.NewValue(tftypes.String, rotation.interval),
		"secret_version":    tftypes.NewValue(tftypes.Number, float64(rotation.secretVersion)),
		"next_renewal":      tftypes.NewValue(tftypes.String, rotation.nextRenewal),
		"expires_at":        tftypes.NewValue(tftypes.String, rotation.expiresAt),
		"current_secret_id": tftypes.NewValue(tftypes.String, rotation.currentSecretID),
		"old_secret_id":     tftypes.NewValue(tftypes.String, oldSecretID),
	}
}

func modifyPlan(t *testing.T, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	t.Helper()
	rs := organization2.Resource()
	rsMp, ok := rs.(resource.ResourceWithModifyPlan)
	require.True(t, ok)
	rsMp.ModifyPlan(t.Context(), req, resp)
}

func assertPlanRotationUnchanged(t *testing.T, plan tfsdk.Plan) {
	t.Helper()
	var planModel organization2.TFModel
	diags := plan.Get(t.Context(), &planModel)
	require.False(t, diags.HasError())
	require.False(t, planModel.ClientSecret.IsUnknown())
}
