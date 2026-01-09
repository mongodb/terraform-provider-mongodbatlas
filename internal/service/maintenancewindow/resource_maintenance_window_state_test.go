package maintenancewindow_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"go.mongodb.org/atlas-sdk/v20250312011/admin"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

// TestAccMaintenanceWindow_UpdateErrorStateNotCorrupted tests that when an API error occurs
// during update, the Terraform state is corrupted but a follow up PATCH still works (can go back to original state)
//
// This test reproduces the bug reported in HELP-87150 where:
// 1. User has a maintenance window configured (e.g., Saturday at 00:00)
// 2. User tries to update to a new time but API rejects (e.g., scheduled maintenance pending)
// 3. Despite the error, the Terraform state is updated with the new values
// 4. Subsequent applies pass, since no real PATCH is sent (API value matches config value).
func TestAccMaintenanceWindow_UpdateErrorStateNotCorrupted(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		// Initial valid configuration
		dayOfWeek           = 7
		hourOfDay           = 0
		validProtectedHours = &admin.ProtectedHours{
			StartHourOfDay: conversion.Pointer(9),
			EndHourOfDay:   conversion.Pointer(17),
		}
		// Invalid configuration: start == end should trigger API error
		invalidProtectedHours = &admin.ProtectedHours{
			StartHourOfDay: conversion.Pointer(10),
			EndHourOfDay:   conversion.Pointer(10), // Same as start - should be invalid
		}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with valid configuration
			{
				Config: configBasic(orgID, projectName, dayOfWeek, hourOfDay, validProtectedHours),
				Check:  checkBasic(dayOfWeek, hourOfDay, validProtectedHours),
			},
			// Step 2: Try to update with invalid protected_hours (same start/end) - expect error
			{
				Config:      configBasic(orgID, projectName, dayOfWeek, hourOfDay, invalidProtectedHours),
				ExpectError: regexp.MustCompile(`(?i)(protected.*hours|invalid|bad.*request)`),
			},
			// Step 3: Apply original config again - should have NO changes if state wasn't corrupted
			// If the bug exists, this step will try to update because state has invalid values
			{
				Config:   configBasic(orgID, projectName, dayOfWeek, hourOfDay, validProtectedHours),
				Check:    checkBasic(dayOfWeek, hourOfDay, validProtectedHours),
			},
		},
	})
}
