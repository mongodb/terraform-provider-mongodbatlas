package customplanmodifier_test

import (
	"context"
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/customplanmodifier"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedclustertpf"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
	"github.com/stretchr/testify/assert"
)

var _ resource.ResourceWithConfigure = &rs{}
var _ resource.ResourceWithImportState = &rs{}
var _ resource.ResourceWithModifyPlan = &rs{}

type BaseResourcePlanModify interface {
	resource.Resource
	resource.ResourceWithConfigure
	resource.ResourceWithImportState
	resource.ResourceWithModifyPlan
}

type planModifyRunData struct {
	keepUnknownCalls []string
	attributeChanges customplanmodifier.AttributeChanges
}

type replaceUnknownResourceInfo struct {
	anyMap map[string]any
}

type replaceUnknownTestCall customplanmodifier.UnknownReplacementCall[replaceUnknownResourceInfo]

func WrappedResource(base BaseResourcePlanModify, info *replaceUnknownResourceInfo, runData *planModifyRunData, attributeReplaceUnknowns map[string]replaceUnknownTestCall) func() resource.Resource {
	return func() resource.Resource {
		return &rs{base: base, info: info, runData: runData, attributeReplaceUnknowns: attributeReplaceUnknowns}
	}
}

type rs struct {
	base                     BaseResourcePlanModify
	runData                  *planModifyRunData
	info                     *replaceUnknownResourceInfo
	attributeReplaceUnknowns map[string]replaceUnknownTestCall
}

func (r *rs) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.base.Metadata(ctx, req, resp)
}

func (r *rs) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.base.Configure(ctx, req, resp)
}

// ModifyPlan is the only method overridden in this test.
func (r *rs) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	schema := advancedclustertpf.ResourceSchema(ctx)
	if req.Plan.Raw.IsFullyKnown() {
		return
	}
	unknownReplacements := customplanmodifier.NewUnknownReplacements(ctx, &req.State, &resp.Plan, &resp.Diagnostics, schema, *r.info)
	for attrName, replacer := range r.attributeReplaceUnknowns {
		modifiedReplacer := func(ctx context.Context, stateValue customplanmodifier.ParsedAttrValue, req *customplanmodifier.UnknownReplacementRequest[replaceUnknownResourceInfo]) attr.Value {
			r.runData.keepUnknownCalls = append(r.runData.keepUnknownCalls, req.Path.String())
			return replacer(ctx, stateValue, req)
		}
		unknownReplacements.AddReplacement(attrName, modifiedReplacer)
	}
	unknownReplacements.ApplyReplacements(ctx, &resp.Diagnostics)
	r.runData.attributeChanges = unknownReplacements.Differ.AttributeChanges
}

func (r *rs) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.base.Schema(ctx, req, resp)
}

func (r *rs) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.base.Create(ctx, req, resp)
}

func (r *rs) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	r.base.Read(ctx, req, resp)
}

func (r *rs) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.base.Update(ctx, req, resp)
}

func (r *rs) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	r.base.Delete(ctx, req, resp)
}

func (r *rs) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	r.base.ImportState(ctx, req, resp)
}

func configureResources(info *replaceUnknownResourceInfo, runData *planModifyRunData, attributeReplaceUnknowns map[string]replaceUnknownTestCall) []func() resource.Resource {
	return []func() resource.Resource{
		WrappedResource(advancedclustertpf.Resource().(BaseResourcePlanModify), info, runData, attributeReplaceUnknowns),
	}
}

type unknownReplacementTestCase struct {
	attributeReplaceUnknowns map[string]replaceUnknownTestCall
	info                     replaceUnknownResourceInfo
	ImportName               string
	ConfigFilename           string
	CheckUnknowns            []tfjsonpath.Path
	CheckKnownValues         []tfjsonpath.Path
	ExtraChecks              []func(string) plancheck.PlanCheck
	expectedAttributeChanges customplanmodifier.AttributeChanges
	expectedKeepUnknownCalls []string
}

func TestReplaceUnknownLogicByWrappingAdvancedClusterTPF(t *testing.T) {
	instanceSizeChanged := customplanmodifier.AttributeChanges{
		"replication_specs",
		"replication_specs[0]",
		"replication_specs[0].region_configs",
		"replication_specs[0].region_configs[0]",
		"replication_specs[0].region_configs[0].electable_specs",
		"replication_specs[0].region_configs[0].electable_specs.instance_size",
		"timeouts",
		"timeouts.create",
	}
	repSpec0 := tfjsonpath.New("replication_specs").AtSliceIndex(0)
	repSpec1 := tfjsonpath.New("replication_specs").AtSliceIndex(1)
	regionConfigPath := repSpec0.AtMapKey("region_configs").AtSliceIndex(0)
	nodeCountChanged := customplanmodifier.AttributeChanges{
		"replication_specs",
		"replication_specs[0]",
		"replication_specs[0].region_configs",
		"replication_specs[0].region_configs[0]",
		"replication_specs[0].region_configs[0].electable_specs",
		"replication_specs[0].region_configs[0].electable_specs.node_count",
		"timeouts",
		"timeouts.create",
	}
	alwaysUnknown := func(ctx context.Context, stateValue customplanmodifier.ParsedAttrValue, req *customplanmodifier.UnknownReplacementRequest[replaceUnknownResourceInfo]) attr.Value {
		return req.Unknown
	}
	for name, tc := range map[string]unknownReplacementTestCase{
		"mongo db major version changed should show in attribute changes and mongo_db_version replace unknown should be called": {
			ImportName:     unit.ImportNameClusterReplicasetOneRegion,
			ConfigFilename: "main_mongo_db_major_version_changed.tf",
			CheckUnknowns: []tfjsonpath.Path{
				tfjsonpath.New("mongo_db_version"),
			},
			attributeReplaceUnknowns: map[string]replaceUnknownTestCall{
				"mongo_db_version": alwaysUnknown,
			},
			expectedAttributeChanges: customplanmodifier.AttributeChanges{"mongo_db_major_version", "timeouts", "timeouts.create"},
			expectedKeepUnknownCalls: []string{"mongo_db_version"},
		},
		"instance_size changed should show changes in parent attributes too": {
			ImportName:               unit.ImportNameClusterReplicasetOneRegion,
			ConfigFilename:           "main_instance_size_changed.tf",
			expectedAttributeChanges: instanceSizeChanged,
		},
		"auto scaling removed should show changes and call replace unknown": {
			ImportName:     unit.ImportNameClusterReplicasetOneRegion,
			ConfigFilename: "main_auto_scaling_removed_node_count_changed.tf",
			CheckUnknowns: []tfjsonpath.Path{
				regionConfigPath.AtMapKey("auto_scaling"),
			},
			attributeReplaceUnknowns: map[string]replaceUnknownTestCall{
				"auto_scaling": alwaysUnknown,
			},
			expectedKeepUnknownCalls: []string{"replication_specs[0].region_configs[0].auto_scaling"},
			expectedAttributeChanges: nodeCountChanged,
		},
		"auto scaling removed but state value returned should update plan": {
			ImportName:     unit.ImportNameClusterReplicasetOneRegion,
			ConfigFilename: "main_auto_scaling_removed_node_count_changed.tf",
			attributeReplaceUnknowns: map[string]replaceUnknownTestCall{
				"auto_scaling": func(ctx context.Context, stateValue customplanmodifier.ParsedAttrValue, req *customplanmodifier.UnknownReplacementRequest[replaceUnknownResourceInfo]) attr.Value {
					return stateValue.AsObject()
				},
			},
			CheckKnownValues: []tfjsonpath.Path{
				regionConfigPath.AtMapKey("auto_scaling"),
			},
			expectedKeepUnknownCalls: []string{"replication_specs[0].region_configs[0].auto_scaling"},
			expectedAttributeChanges: nodeCountChanged,
		},
		"use resource info in value replacement for read_only_specs": {
			ImportName:     unit.ImportNameClusterReplicasetOneRegion,
			ConfigFilename: "main_instance_size_changed.tf",
			attributeReplaceUnknowns: map[string]replaceUnknownTestCall{
				"read_only_specs": func(ctx context.Context, stateValue customplanmodifier.ParsedAttrValue, req *customplanmodifier.UnknownReplacementRequest[replaceUnknownResourceInfo]) attr.Value {
					infoValue, found := req.Info.anyMap["node_count"]
					if !found {
						return req.Unknown
					}
					newValue := customplanmodifier.ReadStateStructValue[advancedclustertpf.TFSpecsModel](ctx, req.Differ, req.Path)
					newValue.NodeCount = types.Int64Value(infoValue.(int64))
					newValue.InstanceSize = types.StringUnknown()
					return conversion.AsObjectValue(ctx, newValue, stateValue.AsObject().AttributeTypes(ctx))
				},
			},
			info: replaceUnknownResourceInfo{
				anyMap: map[string]any{
					"node_count": int64(99),
				},
			},
			ExtraChecks: []func(string) plancheck.PlanCheck{
				func(resourceName string) plancheck.PlanCheck {
					return plancheck.ExpectKnownValue(resourceName, regionConfigPath.AtMapKey("read_only_specs").AtMapKey("node_count"), knownvalue.Int64Exact(99))
				},
			},
			CheckUnknowns: []tfjsonpath.Path{
				regionConfigPath.AtMapKey("read_only_specs").AtMapKey("instance_size"),
			},
			expectedKeepUnknownCalls: []string{"replication_specs[0].region_configs[0].read_only_specs"},
			expectedAttributeChanges: instanceSizeChanged,
		},
		"remove a replication_spec should not call replace unknown": {
			ImportName:     unit.ImportNameClusterTwoRepSpecsWithAutoScalingAndSpecs,
			ConfigFilename: "main_removed_replication_spec.tf",
			attributeReplaceUnknowns: map[string]replaceUnknownTestCall{
				"analytics_auto_scaling": alwaysUnknown,
			},
			expectedAttributeChanges: []string{
				"replication_specs",
				"replication_specs[-1]",
			},
		},
		"add a region config should not call replace unknown in the new region config": {
			ImportName:     unit.ImportNameClusterReplicasetOneRegion,
			ConfigFilename: "main_add_region_config.tf",
			attributeReplaceUnknowns: map[string]replaceUnknownTestCall{
				"analytics_auto_scaling": alwaysUnknown,
			},
			expectedAttributeChanges: []string{
				"replication_specs",
				"replication_specs[0]",
				"replication_specs[0].region_configs",
				"replication_specs[0].region_configs[+1]",
				"replication_specs[0].region_configs[1]",
				"timeouts",
				"timeouts.create",
			},
			expectedKeepUnknownCalls: []string{"replication_specs[0].region_configs[0].analytics_auto_scaling"},
		},
		"add a replication spec should not call replace unknown in the new region config": {
			ImportName:     unit.ImportNameClusterReplicasetOneRegion,
			ConfigFilename: "main_add_replication_spec.tf",
			attributeReplaceUnknowns: map[string]replaceUnknownTestCall{
				"analytics_auto_scaling": alwaysUnknown,
				"id":                     alwaysUnknown,
			},
			expectedAttributeChanges: []string{
				"replication_specs",
				"replication_specs[+1]",
				"replication_specs[1]",
				"timeouts",
				"timeouts.create",
			},
			expectedKeepUnknownCalls: []string{
				"replication_specs[0].id",
				"replication_specs[0].region_configs[0].analytics_auto_scaling",
			},
			CheckUnknowns: []tfjsonpath.Path{
				repSpec1.AtMapKey("id"),
				repSpec1.AtMapKey("region_configs").AtSliceIndex(0).AtMapKey("analytics_auto_scaling"),
			},
		},
		"add tags should not call replace unknown": {
			ImportName:     unit.ImportNameClusterReplicasetOneRegion,
			ConfigFilename: "main_with_tags.tf",
			attributeReplaceUnknowns: map[string]replaceUnknownTestCall{
				"analytics_auto_scaling": alwaysUnknown,
				"id":                     alwaysUnknown,
			},
			expectedAttributeChanges: []string{
				"tags",
				"tags[\"id\"]",
				"timeouts",
				"timeouts.create",
			},
			expectedKeepUnknownCalls: []string{
				"replication_specs[0].id",
				"replication_specs[0].region_configs[0].analytics_auto_scaling",
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			runData := planModifyRunData{}
			mockConfig := unit.MockConfigAdvancedClusterTPF.WithResources(configureResources(&tc.info, &runData, tc.attributeReplaceUnknowns))
			baseConfig := unit.NewMockPlanChecksConfig(t, &mockConfig, tc.ImportName)
			baseConfig.TestdataPrefix = unit.PackagePath("advancedclustertpf")
			checks := make([]plancheck.PlanCheck, 0, len(tc.CheckUnknowns)+len(tc.CheckKnownValues)+len(tc.ExtraChecks))
			for _, checkUnknown := range tc.CheckUnknowns {
				checks = append(checks, plancheck.ExpectUnknownValue(baseConfig.ResourceName, checkUnknown))
			}
			for _, checkKnown := range tc.CheckKnownValues {
				checks = append(checks, plancheck.ExpectKnownValue(baseConfig.ResourceName, checkKnown, knownvalue.NotNull()))
			}
			for _, extraCheck := range tc.ExtraChecks {
				checks = append(checks, extraCheck(baseConfig.ResourceName))
			}
			unit.MockPlanChecksAndRun(t, baseConfig.WithPlanCheckTest(unit.PlanCheckTest{
				ConfigFilename: tc.ConfigFilename,
				Checks:         checks,
			}))
			assert.Equal(t, tc.expectedAttributeChanges, runData.attributeChanges)
			slices.Sort(runData.keepUnknownCalls)
			assert.Equal(t, tc.expectedKeepUnknownCalls, runData.keepUnknownCalls)
		})
	}
}

func TestUnknownReplacements_AddReplacementSameNameShouldPanic(t *testing.T) {
	unknownReplacements := customplanmodifier.UnknownReplacements[replaceUnknownResourceInfo]{
		Differ:       nil,
		Replacements: map[string]customplanmodifier.UnknownReplacementCall[replaceUnknownResourceInfo]{},
		Info:         replaceUnknownResourceInfo{},
	}
	unknownReplacements.AddReplacement("name", nil)
	assert.Panics(t, func() {
		unknownReplacements.AddReplacement("name", nil)
	})
}
