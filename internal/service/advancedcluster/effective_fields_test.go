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
		flag   = baseEffectiveReq(t).withFlag()
		noFlag = flag.withoutFlag()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: flag.config(),
				Check:  flag.check(),
			},
			{
				Config: noFlag.config(),
				Check:  noFlag.check(),
			},
			{
				Config: flag.config(),
				Check:  flag.check(),
			},
			acc.TestStepImportCluster(flag.clusterName),
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

type effectiveReq struct {
	projectID          string
	clusterName        string
	instanceSize       string
	nodeCountElectable int
	useEffectiveFields bool
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
		instanceSize:       "M10",
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

func (req effectiveReq) config() string {
	useEffectiveFieldsConfig := ""
	if req.useEffectiveFields {
		useEffectiveFieldsConfig = "use_effective_fields = true"
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
							}
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
	`, req.projectID, req.clusterName, req.nodeCountElectable, req.instanceSize, useEffectiveFieldsConfig)
}

func (req effectiveReq) check() resource.TestCheckFunc {
	attrsMap := map[string]string{
		"replication_specs.0.region_configs.0.electable_specs.instance_size": req.instanceSize,
		"replication_specs.0.region_configs.0.electable_specs.node_count":    fmt.Sprintf("%d", req.nodeCountElectable),
	}
	extraChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrWith(dataSourcePluralName, "results.#", acc.IntGreatThan(0)),
	}
	if req.useEffectiveFields {
		attrsMap["use_effective_fields"] = "true"
		extraChecks = append(extraChecks,
			resource.TestCheckResourceAttr(dataSourcePluralName, "use_effective_fields", "true"),
		)
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
