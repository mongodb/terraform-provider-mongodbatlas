// To test the internals of unknownReplacements we create a new resource that wraps advanced_cluster TPF but replace the modify plan call to store the attribute changes and the calls to keepUnknown.
package customplanmodifier_test

import (
	"context"
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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

type replaceUnknownResourceInfo struct{} // used to store specific info about the resource, for example upgrade request or sharding schema upgrade

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
	if req.Plan.Raw.IsFullyKnown() {
		return
	}
	schema := advancedclustertpf.ResourceSchema(ctx)
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
	ImportName               string
	ConfigFilename           string
	expectedAttributeChanges customplanmodifier.AttributeChanges
	expectedKeepUnknownCalls []string
}

func alwaysUnknown(ctx context.Context, stateValue customplanmodifier.ParsedAttrValue, req *customplanmodifier.UnknownReplacementRequest[replaceUnknownResourceInfo]) attr.Value {
	return req.Unknown
}

func alwaysState(ctx context.Context, stateValue customplanmodifier.ParsedAttrValue, req *customplanmodifier.UnknownReplacementRequest[replaceUnknownResourceInfo]) attr.Value {
	return stateValue.AsObject()
}

func TestReplaceUnknownLogicByWrappingAdvancedClusterTPF(t *testing.T) {
	var (
		attributeReplaceUnknowns = map[string]replaceUnknownTestCall{
			"auto_scaling":           alwaysState,
			"analytics_auto_scaling": alwaysUnknown,
			"id":                     alwaysUnknown,
		}
		defaultReplaceUnknownCalls = []string{
			"replication_specs[0].id",
			"replication_specs[0].region_configs[0].analytics_auto_scaling",
		}
		defaultAttributeChanges = []string{
			"timeouts",
			"timeouts.create",
		}
	)
	for name, tc := range map[string]unknownReplacementTestCase{
		"no config changes should show the default changes and unknown calls": {
			ImportName:               unit.ImportNameClusterReplicasetOneRegion,
			ConfigFilename:           "main.tf",
			expectedKeepUnknownCalls: defaultReplaceUnknownCalls,
			expectedAttributeChanges: defaultAttributeChanges,
		},
		"root level change should show in attribute changes": {
			ImportName:               unit.ImportNameClusterReplicasetOneRegion,
			ConfigFilename:           "main_mongo_db_major_version_changed.tf",
			expectedAttributeChanges: slices.Concat([]string{"mongo_db_major_version"}, defaultAttributeChanges),
			expectedKeepUnknownCalls: defaultReplaceUnknownCalls,
		},
		"nested change should show changes in parent attributes too": {
			ImportName:     unit.ImportNameClusterReplicasetOneRegion,
			ConfigFilename: "main_instance_size_changed.tf",
			expectedAttributeChanges: slices.Concat([]string{
				"replication_specs",
				"replication_specs[0]",
				"replication_specs[0].region_configs",
				"replication_specs[0].region_configs[0]",
				"replication_specs[0].region_configs[0].electable_specs",
				"replication_specs[0].region_configs[0].electable_specs.instance_size",
			}, defaultAttributeChanges),
			expectedKeepUnknownCalls: defaultReplaceUnknownCalls,
		},
		"auto scaling removed should call replace unknown": {
			ImportName:               unit.ImportNameClusterReplicasetOneRegion,
			ConfigFilename:           "main_auto_scaling_removed_node_count_changed.tf",
			expectedKeepUnknownCalls: slices.Concat(defaultReplaceUnknownCalls, []string{"replication_specs[0].region_configs[0].auto_scaling"}),
			expectedAttributeChanges: slices.Concat([]string{
				"replication_specs",
				"replication_specs[0]",
				"replication_specs[0].region_configs",
				"replication_specs[0].region_configs[0]",
				"replication_specs[0].region_configs[0].electable_specs",
				"replication_specs[0].region_configs[0].electable_specs.node_count",
			}, defaultAttributeChanges),
		},
		"add a region config should not call replace unknown in the new region config": {
			ImportName:     unit.ImportNameClusterReplicasetOneRegion,
			ConfigFilename: "main_add_region_config.tf",
			attributeReplaceUnknowns: map[string]replaceUnknownTestCall{
				"analytics_auto_scaling": alwaysUnknown,
			},
			expectedAttributeChanges: slices.Concat([]string{
				"replication_specs",
				"replication_specs[0]",
				"replication_specs[0].region_configs",
				"replication_specs[0].region_configs[+1]",
				"replication_specs[0].region_configs[1]",
			}, defaultAttributeChanges),
			expectedKeepUnknownCalls: defaultReplaceUnknownCalls,
		},
		"add a replication spec should not call replace unknown in the new replication spec": {
			ImportName:     unit.ImportNameClusterReplicasetOneRegion,
			ConfigFilename: "main_add_replication_spec.tf",
			expectedAttributeChanges: slices.Concat([]string{
				"replication_specs",
				"replication_specs[+1]",
				"replication_specs[1]",
			}, defaultAttributeChanges),
			expectedKeepUnknownCalls: defaultReplaceUnknownCalls,
		},
		"add mapAttribute (tags) should show in attributeChanges but not with '+'": {
			ImportName:     unit.ImportNameClusterReplicasetOneRegion,
			ConfigFilename: "main_with_tags.tf",
			expectedAttributeChanges: slices.Concat([]string{
				"tags",
				"tags[\"id\"]",
			}, defaultAttributeChanges),
			expectedKeepUnknownCalls: defaultReplaceUnknownCalls,
		},
		"add setAttribute (custom_openssl_cipher_config_tls12) should show in attributeChanges but not with '+'": {
			ImportName:     unit.ImportNameClusterReplicasetOneRegion,
			ConfigFilename: "main_tls_cipher_config_mode_with_custom_openssl_cipher_config_tls12.tf",
			expectedAttributeChanges: slices.Concat([]string{
				"advanced_configuration",
				"advanced_configuration.custom_openssl_cipher_config_tls12",
				"advanced_configuration.custom_openssl_cipher_config_tls12[Value(\"ECDHE-RSA-AES256-GCM-SHA384\")]",
				"advanced_configuration.tls_cipher_config_mode",
			}, defaultAttributeChanges),
			expectedKeepUnknownCalls: defaultReplaceUnknownCalls,
		},
		// Different ImportName to test multiple replication specs
		"remove a replication_spec should not call replace unknown in removed spec": {
			ImportName:     unit.ImportNameClusterTwoRepSpecsWithAutoScalingAndSpecs,
			ConfigFilename: "main_removed_replication_spec.tf",
			expectedAttributeChanges: []string{
				"replication_specs",
				"replication_specs[-1]",
			},
			expectedKeepUnknownCalls: []string{
				"replication_specs[0].id",
			},
		},
		"update replication spec1 should not show changes to replication spec0": {
			ImportName:     unit.ImportNameClusterTwoRepSpecsWithAutoScalingAndSpecs,
			ConfigFilename: "main_removed_blocks_from_config_and_instance_change.tf",
			expectedAttributeChanges: []string{
				"replication_specs",
				"replication_specs[1]",
				"replication_specs[1].region_configs",
				"replication_specs[1].region_configs[0]",
				"replication_specs[1].region_configs[0].electable_specs",
				"replication_specs[1].region_configs[0].electable_specs.instance_size",
			},
			expectedKeepUnknownCalls: []string{
				"replication_specs[0].id",
				"replication_specs[0].region_configs[0].analytics_auto_scaling",
				"replication_specs[0].region_configs[0].auto_scaling",
				"replication_specs[1].id",
				"replication_specs[1].region_configs[0].analytics_auto_scaling",
				"replication_specs[1].region_configs[0].auto_scaling",
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			runData := planModifyRunData{}
			mockConfig := unit.MockConfigAdvancedClusterTPF.WithResources(configureResources(&replaceUnknownResourceInfo{}, &runData, attributeReplaceUnknowns))
			baseConfig := unit.NewMockPlanChecksConfig(t, &mockConfig, tc.ImportName)
			baseConfig.TestdataPrefix = unit.PackagePath("advancedclustertpf")
			unit.MockPlanChecksAndRun(t, baseConfig.WithPlanCheckTest(unit.PlanCheckTest{ConfigFilename: tc.ConfigFilename}))
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
