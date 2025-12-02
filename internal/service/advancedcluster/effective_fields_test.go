package advancedcluster_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccAdvancedCluster_effectiveBasic(t *testing.T) {
	var (
		set = baseEffectiveReq(t).withInstanceSize("M10").withFlag()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: set.config(),
				Check:  set.check(),
			},
			// Ignore replication_specs differences as import doesn't use flag so non-effective specs not in the config are set in the state.
			acc.TestStepImportCluster(resourceName, "use_effective_fields", "replication_specs"),
		},
	})
}

func TestAccAdvancedCluster_effectiveUnsetToSet(t *testing.T) {
	var (
		set   = baseEffectiveReq(t).withInstanceSize("M10").withFlag()
		unset = set.withoutFlag()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: unset.config(),
				Check:  unset.check(),
			},
			{
				Config: set.config(),
				Check:  set.check(),
			},
		},
	})
}

func TestAccAdvancedCluster_effectiveSetToUnset(t *testing.T) {
	var (
		set   = baseEffectiveReq(t).withInstanceSize("M10").withFlag()
		unset = set.withoutFlag()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: set.config(),
				Check:  set.check(),
			},
			{
				Config: unset.config(),
				Check:  unset.check(),
			},
		},
	})
}

func TestAccAdvancedCluster_effectiveTenantFlex(t *testing.T) {
	var (
		projectID, clusterName = acc.ProjectIDExecutionWithCluster(t, 0)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config:      configEffectiveTenantFlex(projectID, clusterName, "TENANT", true),
				ExpectError: regexp.MustCompile("use_effective_fields cannot be set for Flex or Tenant clusters"),
			},
			{
				Config:      configEffectiveTenantFlex(projectID, clusterName, "FLEX", true),
				ExpectError: regexp.MustCompile("use_effective_fields cannot be set for Flex or Tenant clusters"),
			},
			{
				Config:      configEffectiveTenantFlex(projectID, clusterName, "TENANT", false),
				ExpectError: regexp.MustCompile("attribute electableSpecs was not specified"), // Try to create cluster when flag is not set.
			},
		},
	})
}

func TestAccAdvancedCluster_effectiveWithOtherChanges(t *testing.T) {
	var (
		unset      = baseEffectiveReq(t).withInstanceSize("M10").withoutFlag()
		setUpdated = unset.withInstanceSize("M20").withFlag()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: unset.config(),
				Check:  unset.check(),
			},
			{
				Config:      setUpdated.config(),
				ExpectError: regexp.MustCompile("Cannot change use_effective_fields with other cluster changes"),
			},
		},
	})
}

func TestAccAdvancedCluster_effectiveComputeAutoScalingInstanceSize(t *testing.T) {
	var (
		initial = baseEffectiveReq(t).withFlag().withComputeMaxInstanceSize("M40").withInstanceSize("M10")
		updated = initial.withInstanceSize("M20").withEffectiveValues(initial)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: initial.config(),
				Check:  initial.check(),
			},
			{
				Config: updated.config(),
				Check:  updated.check(), // Config values echoed in state, but effective specs show actual running values.
			},
		},
	})
}

func TestAccAdvancedCluster_effectiveComputeAutoScalingAll(t *testing.T) {
	var (
		initial = baseEffectiveReq(t).withFlag().withComputeMaxInstanceSize("M40").withInstanceSize("M10").withDiskSizeGB(10).withDiskIOPS(3000)
		updated = initial.withInstanceSize("M20").withDiskSizeGB(15).withDiskIOPS(3010).withEffectiveValues(initial)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: initial.config(),
				Check:  initial.check(),
			},
			{
				Config: updated.config(),
				Check:  updated.check(), // Config values echoed in state, but effective specs show actual running values.
			},
		},
	})
}

func TestAccAdvancedCluster_effectiveDiskAutoScalingAll(t *testing.T) {
	var (
		initial = baseEffectiveReq(t).withFlag().withDiskAutoScaling().withInstanceSize("M10").withDiskSizeGB(10).withDiskIOPS(3000)
		updated = initial.withInstanceSize("M20").withDiskSizeGB(15).withDiskIOPS(3010).withEffectiveValues(initial)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: initial.config(),
				Check:  initial.check(),
			},
			{
				Config: updated.config(),
				Check:  updated.check(), // Config values echoed in state, but effective specs show actual running values.
			},
		},
	})
}

func TestAccAdvancedCluster_effectiveDiskFieldsWithoutAutoScaling(t *testing.T) {
	var (
		initial = baseEffectiveReq(t).withFlag().withInstanceSize("M10").withDiskSizeGB(10).withDiskIOPS(3000)
		updated = initial.withInstanceSize("M20").withDiskSizeGB(15).withDiskIOPS(3010)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: initial.config(),
				Check:  initial.check(),
			},
			{
				Config: updated.config(),
				Check:  updated.check(), // Without auto-scaling, disk fields should update normally.
			},
		},
	})
}

func TestAccAdvancedCluster_effectiveBothAutoScalingEnabled(t *testing.T) {
	var (
		initial = baseEffectiveReq(t).withFlag().withComputeMaxInstanceSize("M40").withDiskAutoScaling().withInstanceSize("M10").withDiskSizeGB(10).withDiskIOPS(3000)
		updated = initial.withInstanceSize("M20").withDiskSizeGB(15).withDiskIOPS(3010).withEffectiveValues(initial)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: initial.config(),
				Check:  initial.check(),
			},
			{
				Config: updated.config(),
				Check:  updated.check(), // Config values echoed in state, but effective specs show actual running values.
			},
		},
	})
}

func TestAccAdvancedCluster_effectiveToggleAutoScaling(t *testing.T) {
	var (
		withoutAutoScaling     = baseEffectiveReq(t).withFlag().withInstanceSize("M10").withDiskSizeGB(10).withDiskIOPS(3000)
		withAutoScaling        = withoutAutoScaling.withComputeMaxInstanceSize("M40").withEffectiveValues(withoutAutoScaling)
		backWithoutAutoScaling = withAutoScaling.withComputeMaxInstanceSize("").withInstanceSize("M20").withDiskSizeGB(15).withDiskIOPS(3010)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: withoutAutoScaling.config(),
				Check:  withoutAutoScaling.check(),
			},
			{
				Config: withAutoScaling.config(),
				Check:  withAutoScaling.check(),
			},
			{
				Config: backWithoutAutoScaling.config(),
				Check:  backWithoutAutoScaling.check(),
			},
		},
	})
}

type effectiveReq struct {
	projectID              string
	clusterName            string
	instanceSize           string
	computeMaxInstanceSize string
	effectiveInstanceSize  string
	nodeCountElectable     int
	diskIOPS               int
	diskSizeGB             int
	effectiveDiskIOPS      int
	effectiveDiskSizeGB    int
	useEffectiveFields     bool
	diskAutoScaling        bool
}

func baseEffectiveReq(t *testing.T) effectiveReq {
	t.Helper()
	var (
		nodeCountElectable     = 3
		projectID, clusterName = acc.ProjectIDExecutionWithCluster(t, nodeCountElectable)
	)
	return effectiveReq{
		projectID:          projectID,
		clusterName:        clusterName,
		nodeCountElectable: nodeCountElectable,
	}
}

func (req effectiveReq) withFlag() effectiveReq {
	req.useEffectiveFields = true
	return req
}

func (req effectiveReq) withoutFlag() effectiveReq {
	req.useEffectiveFields = false
	return req
}

func (req effectiveReq) withInstanceSize(instanceSize string) effectiveReq {
	req.instanceSize = instanceSize
	return req
}

func (req effectiveReq) withDiskSizeGB(diskSizeGB int) effectiveReq {
	req.diskSizeGB = diskSizeGB
	return req
}

func (req effectiveReq) withDiskIOPS(diskIOPS int) effectiveReq {
	req.diskIOPS = diskIOPS
	return req
}

func (req effectiveReq) withComputeMaxInstanceSize(computeMaxInstanceSize string) effectiveReq {
	req.computeMaxInstanceSize = computeMaxInstanceSize
	return req
}

func (req effectiveReq) withDiskAutoScaling() effectiveReq {
	req.diskAutoScaling = true
	return req
}

func (req effectiveReq) withEffectiveValues(effectiveReq effectiveReq) effectiveReq {
	req.effectiveInstanceSize = effectiveReq.instanceSize
	req.effectiveDiskSizeGB = effectiveReq.diskSizeGB
	req.effectiveDiskIOPS = effectiveReq.diskIOPS
	return req
}

func (req effectiveReq) config() string {
	var (
		extraRoot         = ""
		extraRegionConfig = ""
		extraSpecs        = ""
	)
	if req.useEffectiveFields {
		extraRoot += "use_effective_fields = true\n"
	}
	if req.computeMaxInstanceSize != "" || req.diskAutoScaling {
		extraRegionConfig += "auto_scaling = {\n"
		if req.computeMaxInstanceSize != "" {
			extraRegionConfig += "\t\t\tcompute_enabled = true\n"
			extraRegionConfig += fmt.Sprintf("\t\t\tcompute_max_instance_size = %q\n", req.computeMaxInstanceSize)
		}
		if req.diskAutoScaling {
			extraRegionConfig += "\t\t\tdisk_gb_enabled = true\n"
		}
		extraRegionConfig += "\t\t}\n"
	}
	if req.diskIOPS != 0 {
		extraSpecs += fmt.Sprintf("disk_iops = %d\n", req.diskIOPS)
	}
	if req.diskSizeGB != 0 {
		extraSpecs += fmt.Sprintf("disk_size_gb = %d\n", req.diskSizeGB)
	}
	return fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id   = %[1]q
			name         = %[2]q
			cluster_type = "REPLICASET"
			replication_specs = [
				{
					region_configs = [
						{
							priority      = 7
							provider_name = "AWS"
							region_name   = "US_EAST_1"
							electable_specs = {
								node_count    = %[3]d
								instance_size = %[4]q
								%[7]s
							}
							%[6]s
						}
					]
				}
			]
			%[5]s
		}

		data "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			name       = mongodbatlas_advanced_cluster.test.name
			%[5]s
			depends_on = [mongodbatlas_advanced_cluster.test]
		}

		data "mongodbatlas_advanced_clusters" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			%[5]s
			depends_on = [mongodbatlas_advanced_cluster.test]
		}
	`, req.projectID, req.clusterName, req.nodeCountElectable, req.instanceSize, extraRoot, extraRegionConfig, extraSpecs)
}

func (req effectiveReq) check() resource.TestCheckFunc {
	const (
		specsPath           = "replication_specs.0.region_configs.0.electable_specs."
		effectivePath       = "replication_specs.0.region_configs.0.effective_electable_specs."
		effectivePathPlural = "results.0.replication_specs.0.region_configs.0.effective_electable_specs."
		autoScalingPath     = "replication_specs.0.region_configs.0.auto_scaling."
	)
	attrsMap := map[string]string{
		specsPath + "instance_size": req.instanceSize,
		specsPath + "node_count":    fmt.Sprintf("%d", req.nodeCountElectable),
	}
	extraChecks := []resource.TestCheckFunc{
		// Effective fields in singular data source.
		resource.TestCheckResourceAttrSet(dataSourceName, effectivePath+"node_count"),
		resource.TestCheckResourceAttrSet(dataSourceName, effectivePath+"ebs_volume_type"),

		// Effective fields in plural data source - verify they are populated.
		resource.TestCheckResourceAttrWith(dataSourcePluralName, "results.#", acc.IntGreatThan(0)),
		resource.TestCheckResourceAttrSet(dataSourcePluralName, effectivePathPlural+"instance_size"),
		resource.TestCheckResourceAttrSet(dataSourcePluralName, effectivePathPlural+"node_count"),
		resource.TestCheckResourceAttrSet(dataSourcePluralName, effectivePathPlural+"disk_size_gb"),
		resource.TestCheckResourceAttrSet(dataSourcePluralName, effectivePathPlural+"disk_iops"),
		resource.TestCheckResourceAttrSet(dataSourcePluralName, effectivePathPlural+"ebs_volume_type"),
	}
	// Check effective values if specified, otherwise just verify they're set
	if req.effectiveInstanceSize != "" {
		extraChecks = append(extraChecks, resource.TestCheckResourceAttr(dataSourceName, effectivePath+"instance_size", req.effectiveInstanceSize))
	} else {
		extraChecks = append(extraChecks, resource.TestCheckResourceAttrSet(dataSourceName, effectivePath+"instance_size"))
	}
	if req.effectiveDiskSizeGB != 0 {
		extraChecks = append(extraChecks, resource.TestCheckResourceAttr(dataSourceName, effectivePath+"disk_size_gb", fmt.Sprintf("%d", req.effectiveDiskSizeGB)))
	} else {
		extraChecks = append(extraChecks, resource.TestCheckResourceAttrSet(dataSourceName, effectivePath+"disk_size_gb"))
	}
	if req.effectiveDiskIOPS != 0 {
		extraChecks = append(extraChecks, resource.TestCheckResourceAttr(dataSourceName, effectivePath+"disk_iops", fmt.Sprintf("%d", req.effectiveDiskIOPS)))
	} else {
		extraChecks = append(extraChecks, resource.TestCheckResourceAttrSet(dataSourceName, effectivePath+"disk_iops"))
	}
	if req.diskSizeGB != 0 {
		attrsMap[specsPath+"disk_size_gb"] = fmt.Sprintf("%d", req.diskSizeGB)
	}
	if req.diskIOPS != 0 {
		attrsMap[specsPath+"disk_iops"] = fmt.Sprintf("%d", req.diskIOPS)
	}
	if req.computeMaxInstanceSize != "" {
		attrsMap[autoScalingPath+"compute_enabled"] = "true"
		attrsMap[autoScalingPath+"compute_max_instance_size"] = req.computeMaxInstanceSize
	}
	if req.diskAutoScaling {
		attrsMap[autoScalingPath+"disk_gb_enabled"] = "true"
	}
	if req.useEffectiveFields {
		attrsMap["use_effective_fields"] = "true"
		extraChecks = append(extraChecks, resource.TestCheckResourceAttr(dataSourcePluralName, "use_effective_fields", "true"))
	}
	return checkAggr(nil, attrsMap, extraChecks...)
}

// configEffectiveTenantFlex creates a recognizable but incomplete tenant or flex cluster config
// as we're only checking use_effective_fields and the cluster is not actually created.
func configEffectiveTenantFlex(projectID, clusterName, providerName string, useEffectiveFields bool) string {
	var extraConfig string
	if useEffectiveFields {
		extraConfig = "use_effective_fields = true"
	}
	return fmt.Sprintf(`
	resource "mongodbatlas_advanced_cluster" "test" {
		project_id   = %[1]q
		name         = %[2]q
		cluster_type = "REPLICASET"
		replication_specs = [
			{
				region_configs = [
					{
						provider_name = %[3]q
						region_name   = "US_EAST_1"
						priority      = 7
						backing_provider_name = "AWS"
					}
				]
			}
		]
		%[4]s
	}
`, projectID, clusterName, providerName, extraConfig)
}
