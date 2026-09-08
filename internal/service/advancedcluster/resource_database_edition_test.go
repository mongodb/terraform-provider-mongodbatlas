package advancedcluster_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccAdvancedCluster_freeDatabaseEditionUpgrade(t *testing.T) {
	for _, tc := range []struct {
		freeEdition      *string
		dedicatedEdition *string
		name             string
		nodeCount        int
	}{
		{name: "omitted", nodeCount: 3},
		{name: "omitted_to_infinite", dedicatedEdition: new("INFINITE"), nodeCount: 2},
		{name: "core", freeEdition: new("CORE"), dedicatedEdition: new("CORE"), nodeCount: 3},
		{name: "core_to_infinite", freeEdition: new("CORE"), dedicatedEdition: new("INFINITE"), nodeCount: 2},
		{name: "infinite", freeEdition: new("INFINITE"), dedicatedEdition: new("INFINITE"), nodeCount: 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			projectID, clusterName := acc.ProjectIDExecutionWithFreeCluster(t, tc.nodeCount, 1)
			freeConfig := configTenant(t, projectID, clusterName, "Zone 1", freeInstanceSize)
			if tc.freeEdition != nil {
				freeConfig = strings.Replace(freeConfig, `cluster_type = "REPLICASET"`, fmt.Sprintf("cluster_type = \"REPLICASET\"\n database_edition = %q", *tc.freeEdition), 1)
			}
			dedicatedConfig := configDatabaseEditionUpgrade(projectID, clusterName, tc.dedicatedEdition, tc.nodeCount)
			steps := []resource.TestStep{}
			for _, stage := range []struct {
				edition *string
				config  string
			}{
				{config: freeConfig, edition: tc.freeEdition},
				{config: dedicatedConfig, edition: tc.dedicatedEdition},
			} {
				step := resource.TestStep{Config: stage.config}
				if stage.edition == nil {
					step.Check = resource.ComposeAggregateTestCheckFunc(
						acc.CheckExistsCluster(resourceName),
						resource.TestCheckNoResourceAttr(resourceName, "database_edition"),
						resource.TestCheckNoResourceAttr(dataSourceName, "database_edition"),
					)
					step.ConfigStateChecks = []statecheck.StateCheck{
						acc.PluralResultCheck(dataSourcePluralName, "name", knownvalue.StringExact(clusterName), map[string]knownvalue.Check{
							"database_edition": knownvalue.Null(),
						}),
					}
				} else {
					step.Check = checkDatabaseEdition(stage.edition, *stage.edition)
					step.ConfigStateChecks = pluralDatabaseEditionChecks(clusterName, stage.edition, *stage.edition)
				}
				steps = append(steps, step, acc.TestStepImportCluster(resourceName), resource.TestStep{Config: stage.config, PlanOnly: true})
			}
			if tc.dedicatedEdition != nil {
				changedEdition := "CORE"
				if *tc.dedicatedEdition == "CORE" {
					changedEdition = "INFINITE"
				}
				steps = append(steps,
					resource.TestStep{
						Config:      strings.Replace(dedicatedConfig, fmt.Sprintf("database_edition = %q", *tc.dedicatedEdition), fmt.Sprintf("database_edition = %q", changedEdition), 1),
						ExpectError: regexp.MustCompile("databaseEdition cannot be changed"),
					},
					resource.TestStep{Config: dedicatedConfig, PlanOnly: true},
				)
			}
			resource.ParallelTest(t, resource.TestCase{
				PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				CheckDestroy:             acc.CheckDestroyCluster,
				Steps:                    steps,
			})
		})
	}
}

func TestAccAdvancedCluster_flexDatabaseEditionUpgrade(t *testing.T) {
	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 2)
	flexConfig := configFlexCluster(t, projectID, clusterName, "AWS", "US_EAST_1", "Zone 1", "", false, nil)
	dedicatedConfig := configDatabaseEditionUpgrade(projectID, clusterName, new("INFINITE"), 2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             resource.ComposeAggregateTestCheckFunc(acc.CheckDestroyCluster, acc.CheckDestroyFlexCluster),
		Steps: []resource.TestStep{
			{
				Config: flexConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkFlexClusterConfig(projectID, clusterName, "AWS", "US_EAST_1", false, false),
					resource.TestCheckNoResourceAttr(resourceName, "database_edition"),
					resource.TestCheckNoResourceAttr(dataSourceName, "database_edition"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					acc.PluralResultCheck(dataSourcePluralName, "name", knownvalue.StringExact(clusterName), map[string]knownvalue.Check{
						"database_edition": knownvalue.Null(),
					}),
				},
			},
			acc.TestStepImportCluster(resourceName),
			{Config: flexConfig, PlanOnly: true},
			{
				Config:            dedicatedConfig,
				Check:             checkDatabaseEdition(new("INFINITE"), "INFINITE"),
				ConfigStateChecks: pluralDatabaseEditionChecks(clusterName, new("INFINITE"), "INFINITE"),
			},
			acc.TestStepImportCluster(resourceName),
			{Config: dedicatedConfig, PlanOnly: true},
			{
				Config:      strings.Replace(dedicatedConfig, `database_edition = "INFINITE"`, `database_edition = "CORE"`, 1),
				ExpectError: regexp.MustCompile("databaseEdition cannot be changed"),
			},
			{Config: dedicatedConfig, PlanOnly: true},
		},
	})
}

func configDatabaseEditionUpgrade(projectID, clusterName string, databaseEdition *string, nodeCount int) string {
	config := strings.Replace(configDatabaseEdition(projectID, clusterName, databaseEdition, nodeCount), "pit_enabled    = true", "", 1)
	if databaseEdition != nil && *databaseEdition == "INFINITE" {
		config = strings.Replace(config, `instance_size = "M10"`, `instance_size = "M30_GEN_2"`, 1)
	}
	return config
}
