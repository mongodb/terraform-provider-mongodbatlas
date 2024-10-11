package advancedcluster_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

// TODO: if acceptance test will be run in an existing CI group of resources, the name should include the group in the prefix followed by the name of the resource e.i. TestAccStreamRSStreamInstance_basic
// In addition, if acceptance test contains testing of both resource and data sources, the RS/DS can be omitted.
func TestAccAdvancedClusterRS_basic(t *testing.T) {

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
//		CheckDestroy:             checkDestroyAdvancedCluster,
		Steps: []resource.TestStep{ // TODO: verify updates and import in case of resources
//			{
//				Config: advancedClusterConfig(),
//				Check:  advancedClusterAttributeChecks(),
//			},
//          {
//				Config: advancedClusterConfig(),
//				Check:  advancedClusterAttributeChecks(),
//			},
//			{
//				Config:            advancedClusterConfig(),
//				ResourceName:      resourceName,
//				ImportStateIdFunc: checkAdvancedClusterImportStateIDFunc,
//				ImportState:       true,
//				ImportStateVerify: true,
			},
		},
	)
}