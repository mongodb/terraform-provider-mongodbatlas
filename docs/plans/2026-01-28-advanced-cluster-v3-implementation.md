# Advanced Cluster v3.0.0 Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implement breaking changes for advanced_cluster resource v3.0.0 with simplified replication_specs handling and new effective_replication_specs attribute.

**Architecture:** Remove Optional+Computed complexity from replication_specs, making all children Optional only. Add UseEffectiveFieldsReplicationSpecs(true) to all API calls. Add effective_replication_specs as Computed attribute to data sources. Remove plan modifier entirely.

**Tech Stack:** Go, Terraform Plugin Framework, MongoDB Atlas SDK v20250312013

---

## Task 1: Update Schema - Remove use_effective_fields and Change Version

**Files:**
- Modify: `internal/service/advancedcluster/schema.go`

**Step 1: Update schema version from 2 to 3**

In `resourceSchema()` function (line 49), change:
```go
// Before
return schema.Schema{
    Version: 2,

// After
return schema.Schema{
    Version: 3,
```

**Step 2: Remove use_effective_fields from resource schema**

Delete lines 341-347:
```go
// DELETE THIS BLOCK
"use_effective_fields": schema.BoolAttribute{
    Optional: true,
    Validators: []validator.Bool{
        UseEffectiveFieldsValidator{},
    },
    MarkdownDescription: descUseEffectiveFields,
},
```

**Step 3: Remove use_effective_fields from plural data source schema**

In `pluralDataSourceSchema()` function (lines 359-370), remove the `OverridenRootFields` block:
```go
// Before
func pluralDataSourceSchema(ctx context.Context) dsschema.Schema {
    return conversion.PluralDataSourceSchemaFromResource(resourceSchema(ctx), &conversion.PluralDataSourceSchemaRequest{
        RequiredFields:  []string{"project_id"},
        OverridenFields: dataSourceOverridenFields(),
        OverridenRootFields: map[string]dsschema.Attribute{
            "use_effective_fields": dsschema.BoolAttribute{
                Optional:            true,
                MarkdownDescription: descUseEffectiveFields,
            },
        },
    })
}

// After
func pluralDataSourceSchema(ctx context.Context) dsschema.Schema {
    return conversion.PluralDataSourceSchemaFromResource(resourceSchema(ctx), &conversion.PluralDataSourceSchemaRequest{
        RequiredFields:  []string{"project_id"},
        OverridenFields: dataSourceOverridenFields(),
    })
}
```

**Step 4: Remove use_effective_fields from dataSourceOverridenFields**

In `dataSourceOverridenFields()` function (lines 372-383), remove the `use_effective_fields` entry:
```go
// DELETE THIS BLOCK from the map
"use_effective_fields": dsschema.BoolAttribute{
    Optional:            true,
    MarkdownDescription: descUseEffectiveFields,
},
```

**Step 5: Remove descUseEffectiveFields constant**

Delete line 24:
```go
// DELETE
descUseEffectiveFields = "Controls how hardware specification fields are returned..."
```

**Step 6: Remove UseEffectiveFields from TFModel**

Delete line 702:
```go
// DELETE from TFModel struct
UseEffectiveFields types.Bool `tfsdk:"use_effective_fields"`
```

**Step 7: Remove UseEffectiveFields from TFModelDS**

Delete line 734:
```go
// DELETE from TFModelDS struct
UseEffectiveFields types.Bool `tfsdk:"use_effective_fields"`
```

**Step 8: Remove UseEffectiveFields from TFModelPluralDS**

Delete line 740:
```go
// DELETE from TFModelPluralDS struct
UseEffectiveFields types.Bool `tfsdk:"use_effective_fields"`
```

**Step 9: Build to verify**

Run: `make build && golangci-lint run`
Expected: Build may fail due to references to UseEffectiveFields - this is expected, will fix in next tasks.

**Step 10: Commit**

```bash
git add internal/service/advancedcluster/schema.go
git commit -m "feat!: Remove use_effective_fields from advanced_cluster schema (v3)

BREAKING CHANGE: use_effective_fields attribute removed from resource and data sources.
Schema version updated from 2 to 3.
"
```

---

## Task 2: Remove effective_*_specs from Data Source Schema

**Files:**
- Modify: `internal/service/advancedcluster/schema.go`

**Step 1: Remove effective specs from replicationSpecsSchemaDS**

In `replicationSpecsSchemaDS()` function (lines 385-443), remove the three effective_*_specs attributes from region_configs:
```go
// DELETE these three lines from the region_configs attributes map (around line 412-414)
"effective_analytics_specs": specsSchemaDS(),
"effective_electable_specs": specsSchemaDS(),
"effective_read_only_specs": specsSchemaDS(),
```

**Step 2: Remove effective specs from TFRegionConfigsDSModel**

Delete lines 850-852 from `TFRegionConfigsDSModel` struct:
```go
// DELETE from TFRegionConfigsDSModel struct
EffectiveAnalyticsSpecs types.Object `tfsdk:"effective_analytics_specs"`
EffectiveElectableSpecs types.Object `tfsdk:"effective_electable_specs"`
EffectiveReadOnlySpecs  types.Object `tfsdk:"effective_read_only_specs"`
```

**Step 3: Remove effective specs from regionConfigsDSObjType**

Delete lines 864-866 from `regionConfigsDSObjType`:
```go
// DELETE from regionConfigsDSObjType map
"effective_analytics_specs": specsObjType,
"effective_electable_specs": specsObjType,
"effective_read_only_specs": specsObjType,
```

**Step 4: Commit**

```bash
git add internal/service/advancedcluster/schema.go
git commit -m "refactor: Remove effective_*_specs from data source schema

Part of v3.0.0 changes - these will be replaced by effective_replication_specs.
"
```

---

## Task 3: Add effective_replication_specs to Data Source Schema

**Files:**
- Modify: `internal/service/advancedcluster/schema.go`

**Step 1: Add effective_replication_specs to dataSourceOverridenFields**

In `dataSourceOverridenFields()` function, add after the `replication_specs` entry:
```go
func dataSourceOverridenFields() map[string]dsschema.Attribute {
    return map[string]dsschema.Attribute{
        "accept_data_risks_and_force_replica_set_reconfig": nil,
        "delete_on_create_timeout":                         nil,
        "retain_backups_enabled":                           nil,
        "replication_specs":                                replicationSpecsSchemaDS(),
        "effective_replication_specs":                      effectiveReplicationSpecsSchemaDS(),
    }
}
```

**Step 2: Create effectiveReplicationSpecsSchemaDS function**

Add this new function after `replicationSpecsSchemaDS()`:
```go
func effectiveReplicationSpecsSchemaDS() dsschema.ListNestedAttribute {
    return dsschema.ListNestedAttribute{
        Computed:            true,
        MarkdownDescription: "Effective replication specifications representing the actual running configuration as computed by Atlas. This may differ from replication_specs when auto-scaling adjusts instance sizes or other values.",
        NestedObject: dsschema.NestedAttributeObject{
            Attributes: map[string]dsschema.Attribute{
                "container_id": dsschema.MapAttribute{
                    ElementType:         types.StringType,
                    Computed:            true,
                    MarkdownDescription: descContainerID,
                },
                "external_id": dsschema.StringAttribute{
                    Computed:            true,
                    MarkdownDescription: descExternalID,
                },
                "region_configs": dsschema.ListNestedAttribute{
                    Computed:            true,
                    MarkdownDescription: descRegionConfigs,
                    NestedObject: dsschema.NestedAttributeObject{
                        Attributes: map[string]dsschema.Attribute{
                            "analytics_auto_scaling": autoScalingSchemaDS(),
                            "analytics_specs":        specsSchemaDS(),
                            "auto_scaling":           autoScalingSchemaDS(),
                            "backing_provider_name": dsschema.StringAttribute{
                                Computed:            true,
                                MarkdownDescription: descBackingProviderNameTenant,
                            },
                            "electable_specs": specsSchemaDS(),
                            "priority": dsschema.Int64Attribute{
                                Computed:            true,
                                MarkdownDescription: descPriority,
                            },
                            "provider_name": dsschema.StringAttribute{
                                Computed:            true,
                                MarkdownDescription: descProviderName,
                            },
                            "read_only_specs": specsSchemaDS(),
                            "region_name": dsschema.StringAttribute{
                                Computed:            true,
                                MarkdownDescription: descRegionName,
                            },
                        },
                    },
                },
                "zone_id": dsschema.StringAttribute{
                    Computed:            true,
                    MarkdownDescription: descZoneID,
                },
                "zone_name": dsschema.StringAttribute{
                    Computed:            true,
                    MarkdownDescription: descZoneName,
                },
            },
        },
    }
}
```

**Step 3: Add EffectiveReplicationSpecs to TFModelDS**

Add to `TFModelDS` struct:
```go
type TFModelDS struct {
    // ... existing fields ...
    EffectiveReplicationSpecs types.List `tfsdk:"effective_replication_specs"`
}
```

**Step 4: Add EffectiveReplicationSpecs to TFModelPluralDS results**

The TFModelPluralDS uses TFModelDS for results, so this is automatically included.

**Step 5: Build to verify**

Run: `make build && golangci-lint run`
Expected: May still fail due to other references - continue with next tasks.

**Step 6: Commit**

```bash
git add internal/service/advancedcluster/schema.go
git commit -m "feat: Add effective_replication_specs to data source schema

New computed attribute shows actual running configuration from Atlas.
"
```

---

## Task 4: Change replication_specs Children to Optional Only

**Files:**
- Modify: `internal/service/advancedcluster/schema.go`

**Step 1: Update autoScalingSchema to Optional only**

In `autoScalingSchema()` function (lines 445-477), remove `Computed: true` from all attributes:
```go
func autoScalingSchema() schema.SingleNestedAttribute {
    return schema.SingleNestedAttribute{
        Optional:            true,  // Remove Computed: true
        MarkdownDescription: descAutoScaling,
        Attributes: map[string]schema.Attribute{
            "compute_enabled": schema.BoolAttribute{
                Optional:            true,  // Remove Computed: true
                MarkdownDescription: descComputeEnabled,
            },
            "compute_max_instance_size": schema.StringAttribute{
                Optional:            true,  // Remove Computed: true
                MarkdownDescription: descComputeMinMaxInstanceSize,
            },
            "compute_min_instance_size": schema.StringAttribute{
                Optional:            true,  // Remove Computed: true
                MarkdownDescription: descComputeMinMaxInstanceSize,
            },
            "compute_scale_down_enabled": schema.BoolAttribute{
                Optional:            true,  // Remove Computed: true
                MarkdownDescription: descComputeScaleDownEnabled,
            },
            "disk_gb_enabled": schema.BoolAttribute{
                Optional:            true,  // Remove Computed: true
                MarkdownDescription: descDiskGBEnabled,
            },
        },
    }
}
```

**Step 2: Update specsSchema to Optional only**

In `specsSchema()` function (lines 509-544), remove `Computed: true` from all attributes:
```go
func specsSchema() schema.SingleNestedAttribute {
    return schema.SingleNestedAttribute{
        Optional:            true,  // Remove Computed: true
        MarkdownDescription: descSpecs,
        Attributes: map[string]schema.Attribute{
            "disk_iops": schema.Int64Attribute{
                Optional:            true,  // Remove Computed: true
                MarkdownDescription: descDiskIops,
            },
            "disk_size_gb": schema.Float64Attribute{
                Optional:            true,  // Remove Computed: true
                MarkdownDescription: descDiskSizeGb,
            },
            "ebs_volume_type": schema.StringAttribute{
                Optional:            true,  // Remove Computed: true
                MarkdownDescription: descEbsVolumeType,
            },
            "instance_size": schema.StringAttribute{
                Optional: true,  // Remove Computed: true
                PlanModifiers: []planmodifier.String{
                    customplanmodifier.InstanceSizeStringAttributePlanModifier(),
                },
                MarkdownDescription: descInstanceSize,
            },
            "node_count": schema.Int64Attribute{
                Optional:            true,  // Remove Computed: true
                MarkdownDescription: descNodeCount,
            },
        },
    }
}
```

**Step 3: Update zone_name in replication_specs to Optional only**

In `resourceSchema()`, update zone_name (around line 281):
```go
"zone_name": schema.StringAttribute{
    Optional:            true,  // Remove Computed: true
    MarkdownDescription: descZoneName,
},
```

**Step 4: Commit**

```bash
git add internal/service/advancedcluster/schema.go
git commit -m "feat!: Change replication_specs children to Optional only

BREAKING CHANGE: All replication_specs nested attributes are now Optional only,
not Optional+Computed. This means user config is the source of truth.
"
```

---

## Task 5: Delete Plan Modifier

**Files:**
- Delete: `internal/service/advancedcluster/plan_modifier.go`
- Modify: `internal/service/advancedcluster/resource.go`

**Step 1: Remove ModifyPlan method from resource**

In `resource.go`, remove the `resource.ResourceWithModifyPlan` interface (line 25):
```go
// DELETE this line
var _ resource.ResourceWithModifyPlan = &rs{}
```

**Step 2: Remove the ModifyPlan method**

Delete the entire `ModifyPlan` method (lines 78-105):
```go
// DELETE THIS ENTIRE METHOD
func (r *rs) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
    // ... all content ...
}
```

**Step 3: Delete plan_modifier.go file**

```bash
rm internal/service/advancedcluster/plan_modifier.go
```

**Step 4: Delete plan_modifier_test.go file**

```bash
rm internal/service/advancedcluster/plan_modifier_test.go
```

**Step 5: Build to verify**

Run: `make build && golangci-lint run`
Expected: May fail due to unused imports - fix them.

**Step 6: Commit**

```bash
git add -A
git commit -m "refactor: Remove plan modifier for advanced_cluster

No longer needed since replication_specs children are Optional only.
Terraform handles plan/state comparison naturally without custom logic.
"
```

---

## Task 6: Remove UseEffectiveFields from Resource and Common Functions

**Files:**
- Modify: `internal/service/advancedcluster/resource.go`
- Modify: `internal/service/advancedcluster/common_admin_sdk.go`

**Step 1: Remove UseEffectiveFields from ClusterWaitParams**

The `ClusterWaitParams` struct is likely in a common file. Find and remove `UseEffectiveFields` field.

In `resource.go`, find the `ClusterWaitParams` struct or where it's used and remove the field.

**Step 2: Update resolveClusterWaitParams**

In `resource.go`, update `resolveClusterWaitParams()` (lines 468-482) to remove UseEffectiveFields:
```go
func resolveClusterWaitParams(ctx context.Context, model *TFModel, diags *diag.Diagnostics, operation string) *ClusterWaitParams {
    projectID := model.ProjectID.ValueString()
    clusterName := model.Name.ValueString()
    operationTimeout := cleanup.ResolveTimeout(ctx, &model.Timeouts, operation, diags)
    if diags.HasError() {
        return nil
    }
    return &ClusterWaitParams{
        ProjectID:   projectID,
        ClusterName: clusterName,
        Timeout:     operationTimeout,
        IsDelete:    operation == operationDelete,
        // Remove UseEffectiveFields line
    }
}
```

**Step 3: Update GetClusterDetails signature**

In `common_admin_sdk.go`, update `GetClusterDetails()` function signature (line 127):
```go
// Before
func GetClusterDetails(ctx context.Context, diags *diag.Diagnostics, projectID, clusterName string, client *config.MongoDBClient, fcvPresentInState, useEffectiveFields bool) (...)

// After
func GetClusterDetails(ctx context.Context, diags *diag.Diagnostics, projectID, clusterName string, client *config.MongoDBClient, fcvPresentInState bool) (...)
```

**Step 4: Update GetClusterDetails implementation**

Replace `UseEffectiveInstanceFields(useEffectiveFields)` with `UseEffectiveFieldsReplicationSpecs(true)`:
```go
func GetClusterDetails(ctx context.Context, diags *diag.Diagnostics, projectID, clusterName string, client *config.MongoDBClient, fcvPresentInState bool) (clusterDesc *admin.ClusterDescription20240805, flexClusterResp *admin.FlexClusterDescription20241113) {
    isFlex := false
    clusterDesc, resp, err := client.AtlasV2.ClustersApi.GetCluster(ctx, projectID, clusterName).UseEffectiveFieldsReplicationSpecs(true).Execute()
    // ... rest of function ...
}
```

**Step 5: Update all GetClusterDetails call sites**

Update calls in:
- `resource.go` Read method (line 198)
- `resource.go` Update method (line 277)
- `data_source.go` (line 50)
- `resource_test.go` (line 1145)

Example update:
```go
// Before
cluster, flexCluster := GetClusterDetails(ctx, diags, projectID, clusterName, r.Client, !state.PinnedFCV.IsNull(), state.UseEffectiveFields.ValueBool())

// After
cluster, flexCluster := GetClusterDetails(ctx, diags, projectID, clusterName, r.Client, !state.PinnedFCV.IsNull())
```

**Step 6: Commit**

```bash
git add -A
git commit -m "refactor: Remove UseEffectiveFields parameter, always use UseEffectiveFieldsReplicationSpecs(true)

The new API flag is always enabled internally.
"
```

---

## Task 7: Update API Calls to Use UseEffectiveFieldsReplicationSpecs

**Files:**
- Modify: `internal/service/advancedcluster/resource.go`
- Modify: `internal/service/advancedcluster/plural_data_source.go`

**Step 1: Update createCluster function**

In `resource.go`, update `createCluster()` (around line 444):
```go
// Before
_, _, err := client.AtlasV2.ClustersApi.CreateCluster(ctx, waitParams.ProjectID, req).UseEffectiveInstanceFields(waitParams.UseEffectiveFields).Execute()

// After
_, _, err := client.AtlasV2.ClustersApi.CreateCluster(ctx, waitParams.ProjectID, req).UseEffectiveFieldsReplicationSpecs(true).Execute()
```

**Step 2: Update updateCluster function**

In `resource.go`, update `updateCluster()` (around line 460):
```go
// Before
_, _, err := client.AtlasV2.ClustersApi.UpdateCluster(ctx, waitParams.ProjectID, waitParams.ClusterName, req).UseEffectiveInstanceFields(waitParams.UseEffectiveFields).Execute()

// After
_, _, err := client.AtlasV2.ClustersApi.UpdateCluster(ctx, waitParams.ProjectID, waitParams.ClusterName, req).UseEffectiveFieldsReplicationSpecs(true).Execute()
```

**Step 3: Update plural data source ListClusters**

In `plural_data_source.go`, update `getBasicClusters()` (around line 78-83):
```go
// Before
params := admin.ListClustersApiParams{
    GroupId:                    projectID,
    UseEffectiveInstanceFields: conversion.Pointer(useEffectiveFields.ValueBool()),
}

// After
params := admin.ListClustersApiParams{
    GroupId:                            projectID,
    UseEffectiveFieldsReplicationSpecs: conversion.Pointer(true),
}
```

**Step 4: Remove useEffectiveFields parameter from plural data source methods**

Update `getBasicClusters()` and `getFlexClusters()` signatures to remove `useEffectiveFields types.Bool` parameter.

**Step 5: Update readClusters method**

In `plural_data_source.go`, update calls to `getBasicClusters` and `getFlexClusters`:
```go
func (d *pluralDS) readClusters(ctx context.Context, diags *diag.Diagnostics, pluralModel *TFModelPluralDS) (*TFModelPluralDS, *diag.Diagnostics) {
    projectID := pluralModel.ProjectID.ValueString()
    outs := &TFModelPluralDS{
        ProjectID: pluralModel.ProjectID,
    }
    basicClusters := d.getBasicClusters(ctx, diags, projectID)
    // ... rest ...
    flexClusters := d.getFlexClusters(ctx, diags, projectID)
    // ... rest ...
}
```

**Step 6: Remove UseEffectiveFields assignment in appendClusterModelIfValid**

In `plural_data_source.go`, remove line 131:
```go
// DELETE
modelOutDS.UseEffectiveFields = useEffectiveFields
```

**Step 7: Remove UseEffectiveFields assignment in Read method**

In `plural_data_source.go`, remove line 50:
```go
// DELETE
model.UseEffectiveFields = state.UseEffectiveFields
```

**Step 8: Commit**

```bash
git add -A
git commit -m "feat: Use UseEffectiveFieldsReplicationSpecs(true) for all API calls

Always enable the new effective fields flag internally.
"
```

---

## Task 8: Update Model Conversion for effective_replication_specs

**Files:**
- Modify: `internal/service/advancedcluster/model_ClusterDescription20240805.go`

**Step 1: Remove effective specs from newRegionConfigsDSObjType**

In `newRegionConfigsDSObjType()` function (lines 255-272), remove the effective specs assignments:
```go
// DELETE these three lines
dsModel.EffectiveAnalyticsSpecs = newSpecsObjType(ctx, item.EffectiveAnalyticsSpecs, diags)
dsModel.EffectiveElectableSpecs = newSpecsObjType(ctx, item.EffectiveElectableSpecs, diags)
dsModel.EffectiveReadOnlySpecs = newSpecsObjType(ctx, item.EffectiveReadOnlySpecs, diags)
```

**Step 2: Update newTFModelDS to populate EffectiveReplicationSpecs**

In `newTFModelDS()` function (lines 59-67):
```go
func newTFModelDS(ctx context.Context, input *admin.ClusterDescription20240805, diags *diag.Diagnostics, containerIDs map[string]string) *TFModelDS {
    resourceModel := newTFModel(ctx, input, diags, containerIDs)
    if diags.HasError() {
        return nil
    }
    dsModel := conversion.CopyModel[TFModelDS](resourceModel)
    dsModel.ReplicationSpecs = newReplicationSpecsDSObjType(ctx, input.ReplicationSpecs, diags, containerIDs)
    dsModel.EffectiveReplicationSpecs = newEffectiveReplicationSpecsObjType(ctx, input.EffectiveReplicationSpecs, diags, containerIDs)
    return dsModel
}
```

**Step 3: Create newEffectiveReplicationSpecsObjType function**

Add this new function:
```go
func newEffectiveReplicationSpecsObjType(ctx context.Context, input *[]admin.ReplicationSpec20240805, diags *diag.Diagnostics, containerIDs map[string]string) types.List {
    if input == nil {
        return types.ListNull(replicationSpecsDSObjType)
    }
    // Reuse the same conversion logic as replication_specs since structure is identical
    tfModels := convertReplicationSpecs(ctx, input, diags, containerIDs, newRegionConfigsDSObjType)
    if diags.HasError() {
        return types.ListNull(replicationSpecsDSObjType)
    }
    listType, diagsLocal := types.ListValueFrom(ctx, replicationSpecsDSObjType, *tfModels)
    diags.Append(diagsLocal...)
    return listType
}
```

**Step 4: Verify API response includes EffectiveReplicationSpecs**

Check that the Atlas SDK's `ClusterDescription20240805` struct has an `EffectiveReplicationSpecs` field. If not, this will need to come from the API response when `UseEffectiveFieldsReplicationSpecs(true)` is set.

**Step 5: Commit**

```bash
git add -A
git commit -m "feat: Add effective_replication_specs conversion for data sources

Populates the new computed attribute from API response.
"
```

---

## Task 9: Update State Upgrade for v2 to v3

**Files:**
- Modify: `internal/service/advancedcluster/move_upgrade_state.go`

**Step 1: Add v2 to v3 state upgrader to UpgradeState**

Update `UpgradeState()` method (lines 30-34):
```go
func (r *rs) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
    return map[int64]resource.StateUpgrader{
        1: {StateUpgrader: stateUpgraderFromV1},
        2: {StateUpgrader: stateUpgraderFromV2},
    }
}
```

**Step 2: Create stateUpgraderFromV2 function**

Add this new function:
```go
func stateUpgraderFromV2(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
    // V2 to V3 upgrade is straightforward - just copy state
    // The main change is removing use_effective_fields which will be handled by schema
    setStateResponse(ctx, &resp.Diagnostics, req.RawState, &resp.State, true)
}
```

**Step 3: Commit**

```bash
git add internal/service/advancedcluster/move_upgrade_state.go
git commit -m "feat: Add state upgrade path from v2 to v3

Simple upgrade that copies existing state fields to new schema.
"
```

---

## Task 10: Update Data Source Implementation

**Files:**
- Modify: `internal/service/advancedcluster/data_source.go`

**Step 1: Update readCluster method**

Remove the `useEffectiveFields` parameter from GetClusterDetails call and remove the UseEffectiveFields assignment:
```go
func (d *ds) readCluster(ctx context.Context, diags *diag.Diagnostics, modelDS *TFModelDS) *TFModelDS {
    clusterName := modelDS.Name.ValueString()
    projectID := modelDS.ProjectID.ValueString()
    clusterResp, flexClusterResp := GetClusterDetails(ctx, diags, projectID, clusterName, d.Client, false)
    if diags.HasError() {
        return nil
    }
    if flexClusterResp == nil && clusterResp == nil {
        return nil
    }
    var result *TFModelDS
    if flexClusterResp != nil {
        result = convertFlexClusterToDS(ctx, diags, flexClusterResp)
    } else {
        result = convertBasicClusterToDS(ctx, diags, d.Client, clusterResp)
    }
    // Remove: result.UseEffectiveFields = modelDS.UseEffectiveFields
    return result
}
```

**Step 2: Commit**

```bash
git add internal/service/advancedcluster/data_source.go
git commit -m "refactor: Remove use_effective_fields from data source

Always uses new effective fields behavior internally.
"
```

---

## Task 11: Delete UseEffectiveFieldsValidator

**Files:**
- Search for and delete `UseEffectiveFieldsValidator`

**Step 1: Find the validator**

```bash
grep -r "UseEffectiveFieldsValidator" internal/
```

**Step 2: Delete the validator struct and method**

The validator is likely in schema.go or a separate file. Delete the entire struct and its `ValidateBool` method.

**Step 3: Remove import if no longer needed**

Clean up any unused imports.

**Step 4: Commit**

```bash
git add -A
git commit -m "refactor: Remove UseEffectiveFieldsValidator

No longer needed since use_effective_fields is removed.
"
```

---

## Task 12: Update Tests

**Files:**
- Modify: `internal/service/advancedcluster/effective_fields_test.go`
- Possibly delete and recreate with new tests

**Step 1: Rename test file**

```bash
mv internal/service/advancedcluster/effective_fields_test.go internal/service/advancedcluster/effective_replication_specs_test.go
```

**Step 2: Update test content**

Create new test content that tests:
1. Basic cluster creation with effective_replication_specs populated in data source
2. Auto-scaling scenarios where effective values differ from configured
3. Data source shows both replication_specs and effective_replication_specs

Example test structure:
```go
package advancedcluster_test

import (
    "testing"
    "github.com/hashicorp/terraform-plugin-testing/helper/resource"
    "github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccAdvancedCluster_effectiveReplicationSpecs(t *testing.T) {
    // Test that data source populates effective_replication_specs
}

func TestAccAdvancedCluster_effectiveReplicationSpecsWithAutoScaling(t *testing.T) {
    // Test auto-scaling scenario where effective values differ
}
```

**Step 3: Remove old tests that tested use_effective_fields toggle**

Delete tests like:
- `TestAccAdvancedCluster_effectiveUnsetToSet`
- `TestAccAdvancedCluster_effectiveSetToUnset`
- `TestAccAdvancedCluster_effectiveTenantFlex`
- `TestAccAdvancedCluster_effectiveToggleFlagWithRemovedSpecs`

**Step 4: Update remaining tests to not use use_effective_fields**

Remove all references to `use_effective_fields` in test configs.

**Step 5: Commit**

```bash
git add -A
git commit -m "test: Update tests for effective_replication_specs

Remove old use_effective_fields toggle tests.
Add new tests for effective_replication_specs in data sources.
"
```

---

## Task 13: Delete Old Examples and Create New Example

**Files:**
- Delete: `examples/mongodbatlas_advanced_cluster/effective_fields/` (entire directory)
- Create: `examples/mongodbatlas_advanced_cluster/effective_replication_specs/`

**Step 1: Delete old example**

```bash
rm -rf examples/mongodbatlas_advanced_cluster/effective_fields/
```

**Step 2: Create new example directory**

```bash
mkdir -p examples/mongodbatlas_advanced_cluster/effective_replication_specs/
```

**Step 3: Create README.md**

Create `examples/mongodbatlas_advanced_cluster/effective_replication_specs/README.md`:
```markdown
# Effective Replication Specs Example

This example demonstrates how to use `effective_replication_specs` in Terraform Provider 3.x to see the actual running configuration of your MongoDB Atlas cluster.

## Breaking Change Warning

**Important:** In Terraform Provider 3.0.0, the behavior of `replication_specs` has changed:

- All `replication_specs` children are now **Optional only** (not Computed)
- If you previously **removed** `read_only_specs` or `analytics_specs` from your configuration in Provider 2.x, the plan modifier preserved them - they continued running even though they were no longer in your config
- In 3.0.0, Terraform will now detect this mismatch and **plan to delete** those nodes

**Before upgrading:** Review your actual cluster state in Atlas. If your cluster has read-only or analytics nodes that you previously removed from your Terraform config:
1. Add them back to your configuration before upgrading to preserve them, or
2. Accept that Terraform will remove them after upgrading

## Usage

The `effective_replication_specs` attribute is available on data sources and shows the actual running configuration as computed by Atlas. This is useful when using auto-scaling, as the actual instance sizes may differ from your configured values.

```terraform
data "mongodbatlas_advanced_cluster" "example" {
  project_id = mongodbatlas_advanced_cluster.example.project_id
  name       = mongodbatlas_advanced_cluster.example.name
}

# Configured values (what you specified)
output "configured_instance_size" {
  value = data.mongodbatlas_advanced_cluster.example.replication_specs[0].region_configs[0].electable_specs.instance_size
}

# Actual running values (what Atlas computed)
output "actual_instance_size" {
  value = data.mongodbatlas_advanced_cluster.example.effective_replication_specs[0].region_configs[0].electable_specs.instance_size
}
```
```

**Step 4: Create main.tf**

Create `examples/mongodbatlas_advanced_cluster/effective_replication_specs/main.tf`:
```terraform
resource "mongodbatlas_advanced_cluster" "example" {
  project_id   = var.project_id
  name         = var.cluster_name
  cluster_type = "REPLICASET"

  replication_specs = [
    {
      region_configs = [
        {
          electable_specs = {
            instance_size = "M10"
            node_count    = 3
          }
          auto_scaling = {
            compute_enabled            = true
            compute_scale_down_enabled = true
            compute_min_instance_size  = "M10"
            compute_max_instance_size  = "M30"
          }
          provider_name = "AWS"
          priority      = 7
          region_name   = "US_EAST_1"
        }
      ]
    }
  ]
}

data "mongodbatlas_advanced_cluster" "example" {
  project_id = mongodbatlas_advanced_cluster.example.project_id
  name       = mongodbatlas_advanced_cluster.example.name
  depends_on = [mongodbatlas_advanced_cluster.example]
}

# Configured values (what you specified in Terraform)
output "configured_instance_size" {
  description = "The instance size you configured in Terraform"
  value       = data.mongodbatlas_advanced_cluster.example.replication_specs[0].region_configs[0].electable_specs.instance_size
}

# Actual running values (what Atlas computed after auto-scaling)
output "actual_instance_size" {
  description = "The actual instance size running in Atlas (may differ due to auto-scaling)"
  value       = data.mongodbatlas_advanced_cluster.example.effective_replication_specs[0].region_configs[0].electable_specs.instance_size
}
```

**Step 5: Create variables.tf**

Create `examples/mongodbatlas_advanced_cluster/effective_replication_specs/variables.tf`:
```terraform
variable "project_id" {
  description = "MongoDB Atlas Project ID"
  type        = string
}

variable "cluster_name" {
  description = "Name of the cluster"
  type        = string
  default     = "effective-specs-example"
}
```

**Step 6: Create versions.tf**

Create `examples/mongodbatlas_advanced_cluster/effective_replication_specs/versions.tf`:
```terraform
terraform {
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = ">= 3.0.0"
    }
  }
  required_version = ">= 1.0"
}
```

**Step 7: Commit**

```bash
git add -A
git commit -m "doc: Add effective_replication_specs example, remove old effective_fields example

New example demonstrates v3.0.0 effective_replication_specs usage.
Includes breaking change warning about spec removal behavior.
"
```

---

## Task 14: Create Migration Guide

**Files:**
- Create: `docs/guides/3.0.0-advanced-cluster-migration.md`

**Step 1: Create migration guide**

Create `docs/guides/3.0.0-advanced-cluster-migration.md`:
```markdown
---
page_title: "Migrating advanced_cluster to Provider 3.0.0"
subcategory: "Migration Guides"
---

# Migrating advanced_cluster to Terraform Provider 3.0.0

This guide helps you migrate your `mongodbatlas_advanced_cluster` resource configuration from Provider 2.x to 3.0.0.

## Breaking Changes Summary

1. **`use_effective_fields` attribute removed** - This attribute no longer exists. The new behavior is always enabled.

2. **`effective_electable_specs`, `effective_read_only_specs`, `effective_analytics_specs` removed** - These are replaced by `effective_replication_specs` on data sources.

3. **`replication_specs` children are now Optional only** - Previously Optional+Computed, now just Optional. This is the most impactful change.

## Critical: Previously Removed Specs Will Now Be Deleted

**This is the most important change to understand.**

In Provider 2.x, if you **removed** `read_only_specs` or `analytics_specs` from your Terraform configuration, the plan modifier preserved them - they continued running in Atlas even though they weren't in your config.

In Provider 3.0.0, since all `replication_specs` children are Optional only (not Computed), Terraform will now detect this mismatch and **plan to delete** those nodes that exist in your cluster but are missing from your configuration.

### Before Upgrading

1. **Check your Atlas clusters** - Log into Atlas and verify what nodes are actually running
2. **Compare with your Terraform config** - Look for any `read_only_specs` or `analytics_specs` that exist in Atlas but not in your config
3. **Update your config** - Add any missing specs back to your configuration if you want to keep them

### Example

If your cluster in Atlas has:
- 3 electable nodes
- 2 read-only nodes
- 1 analytics node

But your Terraform config only has:
```terraform
replication_specs = [
  {
    region_configs = [
      {
        electable_specs = {
          instance_size = "M10"
          node_count    = 3
        }
        # read_only_specs - NOT DEFINED but running in Atlas!
        # analytics_specs - NOT DEFINED but running in Atlas!
        provider_name = "AWS"
        priority      = 7
        region_name   = "US_EAST_1"
      }
    ]
  }
]
```

After upgrading to 3.0.0, Terraform will plan to **delete** the read-only and analytics nodes.

**To preserve them, update your config before upgrading:**
```terraform
replication_specs = [
  {
    region_configs = [
      {
        electable_specs = {
          instance_size = "M10"
          node_count    = 3
        }
        read_only_specs = {
          instance_size = "M10"
          node_count    = 2
        }
        analytics_specs = {
          instance_size = "M10"
          node_count    = 1
        }
        provider_name = "AWS"
        priority      = 7
        region_name   = "US_EAST_1"
      }
    ]
  }
]
```

## Migration Steps

### Step 1: Remove use_effective_fields

If you have `use_effective_fields = true` in your resource:

```terraform
# Before (2.x)
resource "mongodbatlas_advanced_cluster" "example" {
  use_effective_fields = true  # REMOVE THIS LINE
  # ...
}

# After (3.0.0)
resource "mongodbatlas_advanced_cluster" "example" {
  # use_effective_fields is removed - new behavior is automatic
  # ...
}
```

### Step 2: Update Data Source References

If you were using `effective_electable_specs`, `effective_read_only_specs`, or `effective_analytics_specs`:

```terraform
# Before (2.x)
output "actual_instance_size" {
  value = data.mongodbatlas_advanced_cluster.example.replication_specs[0].region_configs[0].effective_electable_specs.instance_size
}

# After (3.0.0)
output "actual_instance_size" {
  value = data.mongodbatlas_advanced_cluster.example.effective_replication_specs[0].region_configs[0].electable_specs.instance_size
}
```

### Step 3: Ensure All Desired Specs Are in Config

Review your configuration and add any `read_only_specs` or `analytics_specs` that you want to keep.

### Step 4: Run terraform plan

After updating your configuration, run `terraform plan` to verify the changes look correct before applying.

## New effective_replication_specs Attribute

In 3.0.0, data sources have a new `effective_replication_specs` attribute that shows the actual running configuration:

```terraform
data "mongodbatlas_advanced_cluster" "example" {
  project_id = var.project_id
  name       = "my-cluster"
}

# What you configured
output "configured" {
  value = data.mongodbatlas_advanced_cluster.example.replication_specs[0].region_configs[0].electable_specs.instance_size
}

# What's actually running (may differ due to auto-scaling)
output "actual" {
  value = data.mongodbatlas_advanced_cluster.example.effective_replication_specs[0].region_configs[0].electable_specs.instance_size
}
```
```

**Step 2: Commit**

```bash
git add docs/guides/3.0.0-advanced-cluster-migration.md
git commit -m "doc: Add migration guide for advanced_cluster v3.0.0

Covers breaking changes and migration steps from 2.x to 3.0.0.
"
```

---

## Task 15: Update Resource Documentation

**Files:**
- Modify: `docs/resources/advanced_cluster.md`

**Step 1: Remove use_effective_fields references**

Search and remove/update:
- Remove the "IMPORTANT" note about `use_effective_fields` (around line 15)
- Remove the auto-scaling with effective fields example (lines 60-107) or update it
- Remove `use_effective_fields` from argument reference
- Update any examples that use `use_effective_fields`

**Step 2: Add note about Optional-only replication_specs**

Add a note explaining that all `replication_specs` children must be explicitly defined:

```markdown
~> **IMPORTANT:** All `replication_specs` nested attributes (`electable_specs`, `read_only_specs`, `analytics_specs`, `auto_scaling`, etc.) must be explicitly defined in your configuration. If an attribute is not defined, it will not be created or will be removed from your cluster.
```

**Step 3: Update examples**

Update examples to remove `use_effective_fields`:
- Remove the flag from all resource examples
- Update data source examples to use `effective_replication_specs`

**Step 4: Commit**

```bash
git add docs/resources/advanced_cluster.md
git commit -m "doc: Update advanced_cluster resource docs for v3.0.0

Remove use_effective_fields references.
Add note about Optional-only replication_specs behavior.
"
```

---

## Task 16: Update Data Source Documentation

**Files:**
- Modify: `docs/data-sources/advanced_cluster.md`
- Modify: `docs/data-sources/advanced_clusters.md`

**Step 1: Remove use_effective_fields from both files**

Remove from argument reference.

**Step 2: Remove effective_*_specs from attribute reference**

Remove:
- `effective_electable_specs`
- `effective_read_only_specs`
- `effective_analytics_specs`

**Step 3: Add effective_replication_specs to attribute reference**

Add documentation for the new attribute with its structure.

**Step 4: Update examples**

Update any examples to use the new `effective_replication_specs` attribute.

**Step 5: Commit**

```bash
git add docs/data-sources/advanced_cluster.md docs/data-sources/advanced_clusters.md
git commit -m "doc: Update advanced_cluster data source docs for v3.0.0

Remove use_effective_fields and effective_*_specs.
Add effective_replication_specs documentation.
"
```

---

## Task 17: Final Build and Lint Verification

**Step 1: Build the provider**

Run: `make build`
Expected: SUCCESS

**Step 2: Run linter**

Run: `golangci-lint run`
Expected: SUCCESS (or only pre-existing warnings)

**Step 3: Fix any remaining issues**

Address any compilation errors or lint warnings.

**Step 4: Run unit tests**

Run: `make test`
Expected: Tests pass (some may need updates)

**Step 5: Final commit**

```bash
git add -A
git commit -m "chore: Fix remaining build and lint issues for v3.0.0"
```

---

## Task 18: Summary Commit

**Step 1: Create summary commit or tag**

```bash
git tag -a v3.0.0-poc -m "PoC for advanced_cluster v3.0.0 breaking changes"
```

---

## Verification Checklist

- [ ] `make build` succeeds
- [ ] `golangci-lint run` passes
- [ ] Schema version is 3
- [ ] `use_effective_fields` removed from resource and data sources
- [ ] `effective_*_specs` removed from data source region_configs
- [ ] `effective_replication_specs` added to data sources
- [ ] All `replication_specs` children are Optional only
- [ ] Plan modifier deleted
- [ ] `UseEffectiveFieldsReplicationSpecs(true)` used in all API calls
- [ ] State upgrade v2→v3 added
- [ ] Tests updated
- [ ] Examples updated
- [ ] Documentation updated
- [ ] Migration guide created
