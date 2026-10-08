package advancedcluster_test

import (
	"fmt"
	"maps"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSpecChangeWarning_specTypes(t *testing.T) {
	changes := map[string]any{"instance_size": "M20", "disk_size_gb": float64(30), "disk_iops": int64(4000)}
	for _, tc := range []struct{ specName, scaling, fields string }{
		{"electable_specs", "compute_enabled", "instance_size, disk_size_gb, disk_iops"},
		{"electable_specs", "disk_gb_enabled", "instance_size, disk_size_gb, disk_iops"},
		{"read_only_specs", "compute_enabled", "instance_size, disk_size_gb, disk_iops"},
		{"read_only_specs", "disk_gb_enabled", "instance_size, disk_size_gb, disk_iops"},
		{"analytics_specs", "compute_enabled", "instance_size"},
		{"analytics_specs", "disk_gb_enabled", ""},
	} {
		t.Run(tc.specName+"/"+tc.scaling, func(t *testing.T) {
			prior := specWarningModel(specWarningRegion(tc.specName, nil, tc.scaling))
			planned := specWarningModel(specWarningRegion(tc.specName, changes, tc.scaling))
			assertSpecWarning(t, runSpecWarningPlan(t, prior, planned, planned), tc.fields, 0, tc.specName)
		})
	}
}

func TestSpecChangeWarning_valueChanges(t *testing.T) {
	for name, tc := range map[string]struct {
		before, after map[string]any
		fields        string
	}{
		"unchanged":           {nil, nil, ""},
		"instance size":       {nil, map[string]any{"instance_size": "M20"}, "instance_size"},
		"disk size":           {nil, map[string]any{"disk_size_gb": float64(30)}, "disk_size_gb"},
		"disk IOPS":           {nil, map[string]any{"disk_iops": int64(4000)}, "disk_iops"},
		"AUTO":                {nil, map[string]any{"instance_size": "AUTO"}, ""},
		"unrelated field":     {nil, map[string]any{"node_count": int64(5)}, ""},
		"unknown prior value": {map[string]any{"instance_size": tftypes.UnknownValue}, map[string]any{"instance_size": "M20"}, ""},
		"AUTO with disk changes": {nil, map[string]any{
			"instance_size": "AUTO", "disk_size_gb": float64(30), "disk_iops": int64(4000),
		}, "disk_size_gb, disk_iops"},
		"previously omitted disk fields": {
			map[string]any{"disk_size_gb": nil, "disk_iops": nil},
			map[string]any{"disk_size_gb": float64(30), "disk_iops": int64(4000)},
			"disk_size_gb, disk_iops"},
		"removed disk fields": {
			map[string]any{"disk_size_gb": float64(30), "disk_iops": int64(4000)},
			map[string]any{"disk_size_gb": nil, "disk_iops": nil},
			""},
	} {
		t.Run(name, func(t *testing.T) {
			prior := specWarningModel(specWarningRegion("electable_specs", tc.before, "compute_enabled"))
			planned := specWarningModel(specWarningRegion("electable_specs", tc.after, "compute_enabled"))
			assertSpecWarning(t, runSpecWarningPlan(t, prior, planned, planned), tc.fields, 0, "electable_specs")
		})
	}
}

func TestSpecChangeWarning_conditions(t *testing.T) {
	for name, change := range map[string]func(prior, planned, config map[string]any){
		"effective fields disabled": func(_, planned, _ map[string]any) { planned["use_effective_fields"] = false },
		"effective fields unknown":  func(_, planned, _ map[string]any) { planned["use_effective_fields"] = tftypes.UnknownValue },
		"effective fields omitted": func(_, planned, config map[string]any) {
			delete(planned, "use_effective_fields")
			delete(config, "use_effective_fields")
		},
		"omitted specs": func(_, _, config map[string]any) { warningRegion(config)["electable_specs"] = nil },
		"omitted field": func(_, _, config map[string]any) {
			delete(warningRegion(config)["electable_specs"].(map[string]any), "instance_size")
		},
		"unknown field": func(_, _, config map[string]any) {
			warningRegion(config)["electable_specs"].(map[string]any)["instance_size"] = tftypes.UnknownValue
		},
		"unknown specs": func(_, _, config map[string]any) { warningRegion(config)["electable_specs"] = tftypes.UnknownValue },
		"new node type": func(prior, _, _ map[string]any) { warningRegion(prior)["electable_specs"] = nil },
		"unrelated scaling": func(prior, planned, config map[string]any) {
			for _, model := range []map[string]any{prior, planned, config} {
				region := warningRegion(model)
				region["analytics_auto_scaling"] = region["auto_scaling"]
				region["auto_scaling"] = nil
			}
		},
	} {
		t.Run(name, func(t *testing.T) {
			prior := specWarningModel(specWarningRegion("electable_specs", nil, "compute_enabled"))
			planned := specWarningModel(specWarningRegion("electable_specs", map[string]any{"instance_size": "M20"}, "compute_enabled"))
			config := specWarningModel(specWarningRegion("electable_specs", map[string]any{"instance_size": "M20"}, "compute_enabled"))
			change(prior, planned, config)
			require.Empty(t, runSpecWarningPlan(t, prior, planned, config))
		})
	}
	t.Run("create", func(t *testing.T) {
		model := specWarningModel(specWarningRegion("electable_specs", nil, "compute_enabled"))
		require.Empty(t, runSpecWarningPlan(t, nil, model, model))
	})
	t.Run("destroy", func(t *testing.T) {
		model := specWarningModel(specWarningRegion("electable_specs", nil, "compute_enabled"))
		require.Empty(t, runSpecWarningPlan(t, model, nil, model))
	})
}

func TestSpecChangeWarning_autoScalingTransitions(t *testing.T) {
	for _, specName := range []string{"electable_specs", "analytics_specs"} {
		for name, tc := range map[string]struct {
			before, after any
			fields        string
		}{
			"disabled":        {map[string]any{"compute_enabled": false, "disk_gb_enabled": false}, map[string]any{"compute_enabled": false, "disk_gb_enabled": false}, ""},
			"enable compute":  {map[string]any{"compute_enabled": false, "disk_gb_enabled": false}, map[string]any{"compute_enabled": true, "disk_gb_enabled": false}, ""},
			"disable compute": {map[string]any{"compute_enabled": true, "disk_gb_enabled": false}, map[string]any{"compute_enabled": false, "disk_gb_enabled": false}, ""},
			"enable disk while compute stays enabled":  {map[string]any{"compute_enabled": true, "disk_gb_enabled": false}, map[string]any{"compute_enabled": true, "disk_gb_enabled": true}, ""},
			"disable disk while compute stays enabled": {map[string]any{"compute_enabled": true, "disk_gb_enabled": true}, map[string]any{"compute_enabled": true, "disk_gb_enabled": false}, ""},
			"unknown flag":     {map[string]any{"compute_enabled": true, "disk_gb_enabled": false}, map[string]any{"compute_enabled": true, "disk_gb_enabled": tftypes.UnknownValue}, "instance_size"},
			"unknown settings": {map[string]any{"compute_enabled": true, "disk_gb_enabled": false}, tftypes.UnknownValue, ""},
		} {
			t.Run(specName+"/"+name, func(t *testing.T) {
				prior := specWarningRegion(specName, nil, "compute_enabled")
				planned := specWarningRegion(specName, map[string]any{"instance_size": "M20"}, "compute_enabled")
				key := "auto_scaling"
				if specName == "analytics_specs" {
					key = "analytics_auto_scaling"
				}
				prior[key], planned[key] = tc.before, tc.after
				assertSpecWarning(t, runSpecWarningPlan(t, specWarningModel(prior), specWarningModel(planned), specWarningModel(planned)), tc.fields, 0, specName)
			})
		}
	}
}

func TestSpecChangeWarning_topology(t *testing.T) {
	// cluster_type is intentionally irrelevant: the warning keys only on the number of replication_specs entries.
	for _, clusterType := range []string{"REPLICASET", "SHARDED", "GEOSHARDED"} {
		for _, counts := range [][2]int{{1, 1}, {1, 2}, {2, 1}, {2, 2}} {
			t.Run(fmt.Sprintf("%s/%d-to-%d", clusterType, counts[0], counts[1]), func(t *testing.T) {
				model := func(count int, size string) map[string]any {
					region := specWarningRegion("electable_specs", map[string]any{"instance_size": size}, "compute_enabled")
					result := specWarningModel(region)
					result["cluster_type"] = clusterType
					specs := make([]any, count)
					for i := range count {
						specs[i] = map[string]any{"region_configs": []any{region}}
					}
					result["replication_specs"] = specs
					return result
				}
				diagnostics := runSpecWarningPlan(t, model(counts[0], "M10"), model(counts[1], "M20"), model(counts[1], "M20"))
				fields := ""
				if counts == [2]int{1, 1} {
					fields = "instance_size"
				}
				assertSpecWarning(t, diagnostics, fields, 0, "electable_specs")
			})
		}
	}
	t.Run("region count change", func(t *testing.T) {
		// A region is added while the existing region's instance_size changes: index matching is not reliable, so no warning.
		prior := specWarningModel(specWarningRegion("electable_specs", nil, "compute_enabled"))
		changed := specWarningRegion("electable_specs", map[string]any{"instance_size": "M20"}, "compute_enabled")
		added := specWarningRegion("electable_specs", map[string]any{"instance_size": "M20"}, "compute_enabled")
		added["region_name"] = "US_WEST_2"
		require.Empty(t, runSpecWarningPlan(t, prior, specWarningModel(changed, added), specWarningModel(changed, added)))
	})
}

func TestSpecChangeWarning_nodeTypeIndependence(t *testing.T) {
	// Electable auto-scaling must not warn for analytics changes: analytics uses analytics_auto_scaling only.
	region := func(analyticsChanges map[string]any) map[string]any {
		result := specWarningRegion("analytics_specs", analyticsChanges, "compute_enabled")
		result["analytics_auto_scaling"] = map[string]any{"compute_enabled": false, "disk_gb_enabled": false}
		result["electable_specs"] = map[string]any{"instance_size": "M10", "disk_size_gb": float64(20), "disk_iops": int64(3000), "node_count": int64(3)}
		result["auto_scaling"] = map[string]any{"compute_enabled": true, "disk_gb_enabled": true}
		return result
	}
	prior := specWarningModel(region(nil))
	for name, changes := range map[string]map[string]any{
		"instance size": {"instance_size": "M20"},
		"disk size":     {"disk_size_gb": float64(30)},
		"disk IOPS":     {"disk_iops": int64(4000)},
	} {
		t.Run(name, func(t *testing.T) {
			planned := specWarningModel(region(changes))
			require.Empty(t, runSpecWarningPlan(t, prior, planned, planned))
		})
	}
}

func TestSpecChangeWarning_computedPlan(t *testing.T) {
	for name, tc := range map[string]struct {
		change func(prior, planned, config map[string]any)
		fields string
	}{
		"unrelated computed value": {func(_, _, _ map[string]any) {}, "instance_size"},
		"auto-scaling retained from state": {func(_, planned, config map[string]any) {
			warningRegion(planned)["auto_scaling"] = tftypes.UnknownValue
			warningRegion(config)["auto_scaling"] = nil
		}, "instance_size"},
		"storage config on auto-scaling": {func(prior, planned, config map[string]any) {
			for _, model := range []map[string]any{prior, planned, config} {
				warningRegion(model)["auto_scaling"].(map[string]any)["storage_config"] = map[string]any{"shard_size_limit_gb": int64(1024)}
			}
		}, "instance_size"},
		"unknown requested value": {func(_, planned, config map[string]any) {
			for _, model := range []map[string]any{planned, config} {
				warningRegion(model)["electable_specs"].(map[string]any)["instance_size"] = tftypes.UnknownValue
			}
		}, ""},
		"unknown replication specs": {func(_, planned, config map[string]any) {
			planned["replication_specs"], config["replication_specs"] = tftypes.UnknownValue, tftypes.UnknownValue
		}, ""},
		"unknown configured replication specs": {func(_, _, config map[string]any) {
			config["replication_specs"] = tftypes.UnknownValue
		}, ""},
	} {
		t.Run(name, func(t *testing.T) {
			prior := specWarningModel(specWarningRegion("electable_specs", nil, "compute_enabled"))
			planned := specWarningModel(specWarningRegion("electable_specs", map[string]any{"instance_size": "M20"}, "compute_enabled"))
			config := specWarningModel(specWarningRegion("electable_specs", map[string]any{"instance_size": "M20"}, "compute_enabled"))
			prior["cluster_id"] = "333333333333333333333333"
			planned["cluster_id"] = tftypes.UnknownValue
			tc.change(prior, planned, config)
			assertSpecWarning(t, runSpecWarningPlan(t, prior, planned, config), tc.fields, 0, "electable_specs")
		})
	}
}

func TestSpecChangeWarning_combinesRegionsAndNodeTypes(t *testing.T) {
	region := func(name string) map[string]any {
		result := specWarningRegion("electable_specs", nil, "compute_enabled")
		result["region_name"] = name
		result["read_only_specs"] = maps.Clone(result["electable_specs"].(map[string]any))
		result["analytics_specs"] = maps.Clone(result["electable_specs"].(map[string]any))
		result["analytics_auto_scaling"] = maps.Clone(result["auto_scaling"].(map[string]any))
		return result
	}
	prior := specWarningModel(region("US_EAST_1"), region("US_WEST_2"))
	first, second := region("US_EAST_1"), region("US_WEST_2")
	first["electable_specs"].(map[string]any)["instance_size"] = "M20"
	first["electable_specs"].(map[string]any)["disk_size_gb"] = float64(30)
	first["read_only_specs"].(map[string]any)["disk_iops"] = int64(4000)
	first["analytics_specs"].(map[string]any)["instance_size"] = "M20"
	first["analytics_specs"].(map[string]any)["disk_size_gb"] = float64(30) // Analytics disk changes must not enter the warning.
	second["electable_specs"].(map[string]any)["instance_size"] = "AUTO"
	second["electable_specs"].(map[string]any)["disk_size_gb"] = float64(30)
	second["analytics_specs"].(map[string]any)["instance_size"] = "M20"
	second["analytics_auto_scaling"].(map[string]any)["compute_enabled"] = false // Toggling scaling must not enter the warning.
	planned := specWarningModel(first, second)

	assertCombinedSpecWarning(t, runSpecWarningPlan(t, prior, planned, planned), []string{
		"replication_specs[0].region_configs[0].electable_specs.instance_size",
		"replication_specs[0].region_configs[0].electable_specs.disk_size_gb",
		"replication_specs[0].region_configs[0].read_only_specs.disk_iops",
		"replication_specs[0].region_configs[0].analytics_specs.instance_size",
		"replication_specs[0].region_configs[1].electable_specs.disk_size_gb",
	})
}

func TestSpecChangeWarning_providerProtocol(t *testing.T) {
	ctx := t.Context()
	_, typ := clusterSchema(ctx, t)
	prior := specWarningModel(specWarningRegion("electable_specs", nil, "compute_enabled"))
	prior["cluster_id"] = "333333333333333333333333"
	changes := map[string]any{"instance_size": "M20", "disk_size_gb": float64(30), "disk_iops": int64(4000)}
	for name, tc := range map[string]struct {
		value    any
		identity string
	}{
		"update":             {},
		"changed name":       {identity: "name", value: "replacement"},
		"changed project":    {identity: "project_id", value: "444444444444444444444444"},
		"unresolved name":    {identity: "name", value: tftypes.UnknownValue},
		"unresolved project": {identity: "project_id", value: tftypes.UnknownValue},
	} {
		t.Run(name, func(t *testing.T) {
			config := specWarningModel(specWarningRegion("electable_specs", changes, "compute_enabled"))
			if tc.identity != "" {
				config[tc.identity] = tc.value
			}
			proposed := maps.Clone(config)
			proposed["cluster_id"] = tftypes.UnknownValue

			// Exercise framework plan modifiers and resource ModifyPlan through the protocol, without configuring an Atlas client.
			server, err := acc.TestAccProviderV6Factories["mongodbatlas"]()
			require.NoError(t, err)
			response, err := server.PlanResourceChange(ctx, &tfprotov6.PlanResourceChangeRequest{
				TypeName:         "mongodbatlas_advanced_cluster",
				PriorState:       clusterDynamic(t, typ, prior),
				ProposedNewState: clusterDynamic(t, typ, proposed),
				Config:           clusterDynamic(t, typ, config),
			})
			require.NoError(t, err)
			if tc.identity != "" {
				require.Empty(t, response.Diagnostics, "replacement clusters apply the requested specs")
				require.NotNil(t, response.PlannedState)
				require.Len(t, response.RequiresReplace, 1)
				assert.True(t, tftypes.NewAttributePath().WithAttributeName(tc.identity).Equal(response.RequiresReplace[0]))
				return
			}
			require.Empty(t, response.RequiresReplace)
			require.NotNil(t, response.PlannedState)
			require.Len(t, response.Diagnostics, 1)
			warning := response.Diagnostics[0]
			assert.Equal(t, tfprotov6.DiagnosticSeverityWarning, warning.Severity)
			assert.Equal(t, "Spec changes are ignored while use_effective_fields and auto-scaling remain enabled", warning.Summary)
			assert.True(t, tftypes.NewAttributePath().WithAttributeName("replication_specs").Equal(warning.Attribute))
			assert.Contains(t, warning.Detail, "\n\n- replication_specs[0].region_configs[0].electable_specs.instance_size\n"+
				"- replication_specs[0].region_configs[0].electable_specs.disk_size_gb\n"+
				"- replication_specs[0].region_configs[0].electable_specs.disk_iops\n\n")

			// Warning preserves requested changes while existing plan logic restores computed state values.
			planned, err := response.PlannedState.Unmarshal(typ)
			require.NoError(t, err)
			specsPath := tftypes.NewAttributePath().WithAttributeName("replication_specs").WithElementKeyInt(0).
				WithAttributeName("region_configs").WithElementKeyInt(0).WithAttributeName("electable_specs")
			for field, expected := range changes {
				value, _, err := tftypes.WalkAttributePath(planned, specsPath.WithAttributeName(field))
				require.NoError(t, err)
				actual := value.(tftypes.Value)
				assert.True(t, actual.Equal(tftypes.NewValue(actual.Type(), expected)), "%s must retain its configured value", field)
			}
			clusterID, _, err := tftypes.WalkAttributePath(planned, tftypes.NewAttributePath().WithAttributeName("cluster_id"))
			require.NoError(t, err)
			assert.Equal(t, tftypes.NewValue(tftypes.String, prior["cluster_id"]), clusterID)
		})
	}
}

func specWarningRegion(specName string, changes map[string]any, scaling string) map[string]any {
	scalingName := "auto_scaling"
	if specName == "analytics_specs" {
		scalingName = "analytics_auto_scaling"
	}
	specs := map[string]any{"instance_size": "M10", "disk_size_gb": float64(20), "disk_iops": int64(3000), "node_count": int64(3)}
	maps.Copy(specs, changes)
	return map[string]any{
		"provider_name": "AWS", "region_name": "US_EAST_1", "priority": int64(7),
		specName: specs, scalingName: map[string]any{"compute_enabled": scaling == "compute_enabled", "disk_gb_enabled": scaling == "disk_gb_enabled"},
	}
}

func specWarningModel(regions ...any) map[string]any {
	return map[string]any{
		"name": "example", "project_id": "111111111111111111111111", "cluster_type": "REPLICASET",
		"use_effective_fields": true,
		"replication_specs":    []any{map[string]any{"region_configs": regions}},
	}
}

func warningRegion(model map[string]any) map[string]any {
	return model["replication_specs"].([]any)[0].(map[string]any)["region_configs"].([]any)[0].(map[string]any)
}

func runSpecWarningPlan(t *testing.T, prior, planned, config map[string]any) diag.Diagnostics {
	t.Helper()
	ctx := t.Context()
	schema, typ := clusterSchema(ctx, t)
	raw := func(model map[string]any) tftypes.Value {
		if model == nil {
			return tftypes.NewValue(typ, nil)
		}
		return planTestValue(typ, model)
	}
	request := resource.ModifyPlanRequest{
		State:  tfsdk.State{Schema: schema, Raw: raw(prior)},
		Plan:   tfsdk.Plan{Schema: schema, Raw: raw(planned)},
		Config: tfsdk.Config{Schema: schema, Raw: raw(config)},
	}
	response := resource.ModifyPlanResponse{Plan: request.Plan}
	advancedcluster.Resource().(resource.ResourceWithModifyPlan).ModifyPlan(ctx, request, &response)
	require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
	if request.Plan.Raw.IsFullyKnown() {
		assert.Equal(t, request.Plan.Raw, response.Plan.Raw, "emitting a warning must not change a fully known plan")
	}
	return response.Diagnostics
}

func assertSpecWarning(t *testing.T, diagnostics diag.Diagnostics, fields string, regionIndex int, specName string) {
	t.Helper()
	if fields == "" {
		require.Empty(t, diagnostics)
		return
	}
	var expectedPaths []string
	for field := range strings.SplitSeq(fields, ", ") {
		expectedPaths = append(expectedPaths, fmt.Sprintf("replication_specs[0].region_configs[%d].%s.%s", regionIndex, specName, field))
	}
	assertCombinedSpecWarning(t, diagnostics, expectedPaths)
}

func assertCombinedSpecWarning(t *testing.T, diagnostics diag.Diagnostics, expectedPaths []string) {
	t.Helper()
	require.Len(t, diagnostics, 1)
	warning := diagnostics[0]
	assert.Equal(t, diag.SeverityWarning, warning.Severity())
	assert.Equal(t, "Spec changes are ignored while use_effective_fields and auto-scaling remain enabled", warning.Summary())
	assert.Contains(t, warning.Detail(), "Atlas ignores changes to the following attributes, although Terraform stores their new values in state:")
	assert.Contains(t, warning.Detail(), "\n\n- "+strings.Join(expectedPaths, "\n- ")+"\n\n")
	advice := "disable auto-scaling and apply the desired values, then re-enable auto-scaling in a separate apply"
	assert.Equal(t, 1, strings.Count(warning.Detail(), advice), "the advice must appear only once")
	assert.Contains(t, warning.Detail(), "#manually-updating-specs-with-use_effective_fields")
	attributeWarning, ok := warning.(diag.DiagnosticWithPath)
	require.True(t, ok)
	assert.True(t, path.Root("replication_specs").Equal(attributeWarning.Path()), "the combined warning must identify replication_specs")
}
