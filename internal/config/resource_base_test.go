package config_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/stretchr/testify/assert"
)

func TestNoResourceInterfaceLoss(t *testing.T) {
	analyticsResource := config.AnalyticsResourceFunc(advancedcluster.Resource())()
	_, ok := analyticsResource.(resource.ResourceWithModifyPlan)
	assert.True(t, ok)
	_, ok = analyticsResource.(resource.ResourceWithUpgradeState)
	assert.True(t, ok)
	_, ok = analyticsResource.(resource.ResourceWithMoveState)
	assert.True(t, ok)
	_, ok = analyticsResource.(resource.ResourceWithUpgradeState)
	assert.True(t, ok)
	_, ok = analyticsResource.(resource.ResourceWithImportState)
	assert.True(t, ok)
}
