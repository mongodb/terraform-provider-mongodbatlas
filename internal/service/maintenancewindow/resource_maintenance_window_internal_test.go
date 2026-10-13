// Package declaration is not `maintenancewindow_test` because it tests the unexported
// `newProtectedHours` helper, which cannot be reached through the exported API without
// a full Terraform diff (resourceCreate reads GetRawConfig).
package maintenancewindow

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestNewProtectedHoursKeepsZero(t *testing.T) {
	d := schema.TestResourceDataRaw(t, Resource().Schema, map[string]any{
		"protected_hours": []any{map[string]any{
			"start_hour_of_day": 12,
			"end_hour_of_day":   0,
		}},
	})

	protectedHours := newProtectedHours(d)
	require.NotNil(t, protectedHours, "protected_hours must be sent")
	require.NotNil(t, protectedHours.StartHourOfDay, "start_hour_of_day = 12 must be sent, not omitted")
	require.Equal(t, 12, *protectedHours.StartHourOfDay)
	require.NotNil(t, protectedHours.EndHourOfDay, "end_hour_of_day = 0 must be sent, not omitted")
	require.Equal(t, 0, *protectedHours.EndHourOfDay)
}
