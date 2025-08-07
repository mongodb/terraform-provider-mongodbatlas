package maintenancewindow_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/spf13/cast"
	"go.mongodb.org/atlas-sdk/v20250312006/admin"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const resourceName = "mongodbatlas_maintenance_window.test"

var (
	defaultProtectedHours = &admin.ProtectedHours{
		StartHourOfDay: conversion.Pointer(9),
		EndHourOfDay:   conversion.Pointer(17),
	}
	updatedProtectedHours = &admin.ProtectedHours{
		StartHourOfDay: conversion.Pointer(10),
		EndHourOfDay:   conversion.Pointer(15),
	}
)

func TestAccConfigRSMaintenanceWindow_basic(t *testing.T) {
	var (
		orgID            = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName      = acc.RandomProjectName()
		dayOfWeek        = 7
		hourOfDay        = 0
		dayOfWeekUpdated = 4
		hourOfDayUpdated = 5
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, projectName, dayOfWeek, hourOfDay, defaultProtectedHours),
				Check:  checkBasic(dayOfWeek, hourOfDay, defaultProtectedHours),
			},
			{
				Config: configBasic(orgID, projectName, dayOfWeek, hourOfDayUpdated, updatedProtectedHours),
				Check:  checkBasic(dayOfWeek, hourOfDayUpdated, updatedProtectedHours),
			},
			{
				Config: configBasic(orgID, projectName, dayOfWeekUpdated, hourOfDay, nil),
				Check:  checkBasic(dayOfWeekUpdated, hourOfDay, nil),
			},
			{
				Config: configBasic(orgID, projectName, dayOfWeek, hourOfDay, defaultProtectedHours),
				Check:  checkBasic(dayOfWeek, hourOfDay, defaultProtectedHours),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccConfigRSMaintenanceWindow_autoDeferActivated(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		dayOfWeek   = 7
		hourOfDay   = 3
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithAutoDeferEnabled(orgID, projectName, dayOfWeek, hourOfDay),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "day_of_week", cast.ToString(dayOfWeek)),
					resource.TestCheckResourceAttr(resourceName, "hour_of_day", cast.ToString(hourOfDay)),
					resource.TestCheckResourceAttr(resourceName, "number_of_deferrals", "0"),
					resource.TestCheckResourceAttr(resourceName, "auto_defer_once_enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "time_zone_id"),
				),
			},
		},
	})
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		log.Printf("[DEBUG] projectID: %s", rs.Primary.ID)
		_, _, err := acc.ConnV2().MaintenanceWindowsApi.GetMaintenanceWindow(context.Background(), rs.Primary.ID).Execute()
		if err != nil {
			return fmt.Errorf("maintenance Window (%s) does not exist", rs.Primary.ID)
		}
		return nil
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_maintenance_window" {
			continue
		}
		projectID := rs.Primary.ID
		if projectID == "" {
			return fmt.Errorf("checkDestroy, no ID is set for: %s", resourceName)
		}
		maintenanceWindow, _, _ := acc.ConnV2().MaintenanceWindowsApi.GetMaintenanceWindow(context.Background(), projectID).Execute()
		// Check if it's back to default settings (day_of_week = 0 means it's been reset)
		if maintenanceWindow.GetDayOfWeek() != 0 {
			return fmt.Errorf("maintenance window for project (%s) was not properly reset to defaults", projectID)
		}
	}
	return nil
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}

func configBasic(orgID, projectName string, dayOfWeek, hourOfDay int, protectedHours *admin.ProtectedHours) string {
	protectedHoursStr := ""
	if protectedHours != nil {
		protectedHoursStr = fmt.Sprintf(`
			protected_hours {
				start_hour_of_day = %[1]d
				end_hour_of_day   = %[2]d
			}`, *protectedHours.StartHourOfDay, *protectedHours.EndHourOfDay)
	}

	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_maintenance_window" "test" {
			project_id  = mongodbatlas_project.test.id
			day_of_week = %[3]d
			hour_of_day = %[4]d
			%[5]s

		}`, orgID, projectName, dayOfWeek, hourOfDay, protectedHoursStr)
}

func configWithAutoDeferEnabled(orgID, projectName string, dayOfWeek, hourOfDay int) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_maintenance_window" "test" {
			project_id  = mongodbatlas_project.test.id
			day_of_week = %[3]d
			hour_of_day = %[4]d
			auto_defer_once_enabled = true
		}`, orgID, projectName, dayOfWeek, hourOfDay)
}

func checkBasic(dayOfWeek, hourOfDay int, protectedHours *admin.ProtectedHours) resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{
		checkExists(resourceName),
		resource.TestCheckResourceAttrSet(resourceName, "project_id"),
		resource.TestCheckResourceAttr(resourceName, "day_of_week", cast.ToString(dayOfWeek)),
		resource.TestCheckResourceAttr(resourceName, "hour_of_day", cast.ToString(hourOfDay)),
		resource.TestCheckResourceAttr(resourceName, "number_of_deferrals", "0"),
	}
	if protectedHours != nil {
		checks = append(checks,
			resource.TestCheckResourceAttr(resourceName, "protected_hours.0.start_hour_of_day", cast.ToString(*protectedHours.StartHourOfDay)),
			resource.TestCheckResourceAttr(resourceName, "protected_hours.0.end_hour_of_day", cast.ToString(*protectedHours.EndHourOfDay)),
		)
	} else {
		checks = append(checks, resource.TestCheckResourceAttr(resourceName, "protected_hours.#", "0"))
	}
	return resource.ComposeAggregateTestCheckFunc(checks...)
}
