package alertconfiguration_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/alertconfiguration"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"go.mongodb.org/atlas-sdk/v20231115007/admin"
)

func TestAccConfigRSAlertConfiguration_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_alert_configuration.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		alert        = &admin.GroupAlertsConfig{}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasicRS(orgID, projectName, true),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "2"),
				),
			},
			{
				Config: configBasicRS(orgID, projectName, false),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "2"),
				),
			},
		},
	})
}

func TestAccConfigRSAlertConfiguration_withEmptyMetricThresholdConfig(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithEmptyMetricThresholdConfig(orgID, projectName, true),
			},
		},
	})
}

func TestAccConfigRSAlertConfiguration_withEmptyMatcherMetricThresholdConfig(t *testing.T) {
	var (
		resourceName = "mongodbatlas_alert_configuration.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		alert        = &admin.GroupAlertsConfig{}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithEmptyMatcherMetricThresholdConfig(orgID, projectName, true),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "1"),
				),
			},
		},
	})
}
func TestAccConfigRSAlertConfiguration_withNotifications(t *testing.T) {
	var (
		resourceName = "mongodbatlas_alert_configuration.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		alert        = &admin.GroupAlertsConfig{}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithNotifications(orgID, projectName, true, true, false),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
			{
				Config: configWithNotifications(orgID, projectName, false, false, true),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
		},
	})
}

func TestAccConfigRSAlertConfiguration_withMatchers(t *testing.T) {
	var (
		resourceName = "mongodbatlas_alert_configuration.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		alert        = &admin.GroupAlertsConfig{}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithMatchers(orgID, projectName, true, false, true,
					map[string]interface{}{
						"fieldName": "TYPE_NAME",
						"operator":  "EQUALS",
						"value":     "SECONDARY",
					},
					map[string]interface{}{
						"fieldName": "TYPE_NAME",
						"operator":  "CONTAINS",
						"value":     "MONGOS",
					}),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
			{
				Config: configWithMatchers(orgID, projectName, false, true, false,
					map[string]interface{}{
						"fieldName": "TYPE_NAME",
						"operator":  "NOT_EQUALS",
						"value":     "SECONDARY",
					},
					map[string]interface{}{
						"fieldName": "HOSTNAME",
						"operator":  "EQUALS",
						"value":     "PRIMARY",
					}),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
		},
	})
}

func TestAccConfigRSAlertConfiguration_withMetricUpdated(t *testing.T) {
	var (
		resourceName = "mongodbatlas_alert_configuration.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		alert        = &admin.GroupAlertsConfig{}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithMetricUpdated(orgID, projectName, true, 99.0),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
			{
				Config: configWithMetricUpdated(orgID, projectName, false, 89.7),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
		},
	})
}

func TestAccConfigRSAlertConfiguration_withThresholdUpdated(t *testing.T) {
	var (
		resourceName = "mongodbatlas_alert_configuration.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		alert        = &admin.GroupAlertsConfig{}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithThresholdUpdated(orgID, projectName, true, 1),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
			{
				Config: configWithThresholdUpdated(orgID, projectName, false, 3),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       importStateProjectIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id", "updated", "matcher.0.field_name"},
			},
		},
	})
}

func TestAccConfigRSAlertConfiguration_withoutRoles(t *testing.T) {
	var (
		alert        = &admin.GroupAlertsConfig{}
		resourceName = "mongodbatlas_alert_configuration.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithoutRoles(orgID, projectName, true, 99.0),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
		},
	})
}

func TestAccConfigRSAlertConfiguration_withoutOptionalAttributes(t *testing.T) {
	var (
		alert        = &admin.GroupAlertsConfig{}
		resourceName = "mongodbatlas_alert_configuration.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithEmptyOptionalAttributes(orgID, projectName),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
		},
	})
}

func TestAccConfigRSAlertConfiguration_importBasic(t *testing.T) {
	var (
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		resourceName = "mongodbatlas_alert_configuration.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasicRS(orgID, projectName, true),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       importStateProjectIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id", "updated"},
			},
		},
	})
}

func TestAccConfigRSAlertConfiguration_importIncorrectId(t *testing.T) {
	var (
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		resourceName = "mongodbatlas_alert_configuration.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasicRS(orgID, projectName, true),
			},
			{
				ResourceName:  resourceName,
				ImportState:   true,
				ImportStateId: "incorrect_id_without_project_id_and_dash",
				ExpectError:   regexp.MustCompile("import format error"),
			},
		},
	})
}

func TestAccConfigRSAlertConfiguration_importConfigNotifications(t *testing.T) {
	var (
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		resourceName = "mongodbatlas_alert_configuration.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithNotifications(orgID, projectName, true, true, false),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       importStateProjectIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id", "updated"},
			},
		},
	})
}

// dummy keys used for credential values in third party notifications
const dummy32CharKey = "11111111111111111111111111111111"
const dummy36CharKey = "11111111-1111-1111-1111-111111111111"

// used for testing notification that does not define interval_min attribute
func TestAccConfigRSAlertConfiguration_importPagerDuty(t *testing.T) {
	var (
		resourceName = "mongodbatlas_alert_configuration.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		serviceKey   = dummy32CharKey
		alert        = &admin.GroupAlertsConfig{}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithPagerDuty(orgID, projectName, serviceKey, true),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       importStateProjectIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"updated", "notification.0.service_key"}, // service key is not returned by api in import operation
			},
		},
	})
}

func TestAccConfigRSAlertConfiguration_updatePagerDutyWithNotifierId(t *testing.T) {
	var (
		resourceName = "mongodbatlas_alert_configuration.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		serviceKey   = dummy32CharKey
		notifierID   = "651dd9336afac13e1c112222"
		alert        = &admin.GroupAlertsConfig{}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithPagerDutyNotifierID(orgID, projectName, notifierID, 10, &serviceKey),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "notification.0.delay_min", "10"),
					resource.TestCheckResourceAttr(resourceName, "notification.0.service_key", serviceKey),
				),
			},
			{
				Config: configWithPagerDutyNotifierID(orgID, projectName, notifierID, 15, nil),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "notification.0.delay_min", "15"),
				),
			},
		},
	})
}

func TestAccConfigRSAlertConfiguration_withDataDog(t *testing.T) {
	var (
		resourceName = "mongodbatlas_alert_configuration.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		ddAPIKey     = dummy32CharKey
		ddRegion     = "US"
		alert        = &admin.GroupAlertsConfig{}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithDataDog(orgID, projectName, ddAPIKey, ddRegion, true),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
		},
	})
}

func TestAccConfigRSAlertConfiguration_withPagerDuty(t *testing.T) {
	var (
		resourceName = "mongodbatlas_alert_configuration.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		serviceKey   = dummy32CharKey
		alert        = &admin.GroupAlertsConfig{}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithPagerDuty(orgID, projectName, serviceKey, true),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
		},
	})
}

func TestAccConfigRSAlertConfiguration_withOpsGenie(t *testing.T) {
	var (
		resourceName = "mongodbatlas_alert_configuration.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		apiKey       = dummy36CharKey
		alert        = &admin.GroupAlertsConfig{}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithOpsGenie(orgID, projectName, apiKey, true),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
		},
	})
}

func TestAccConfigRSAlertConfiguration_withVictorOps(t *testing.T) {
	var (
		resourceName = "mongodbatlas_alert_configuration.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		apiKey       = dummy36CharKey
		alert        = &admin.GroupAlertsConfig{}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithVictorOps(orgID, projectName, apiKey, true),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
		},
	})
}

func checkExists(resourceName string, alert *admin.GroupAlertsConfig) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		alertResp, _, err := acc.ConnV2().AlertConfigurationsApi.GetAlertConfiguration(context.Background(), ids[alertconfiguration.EncodedIDKeyProjectID], ids[alertconfiguration.EncodedIDKeyAlertID]).Execute()
		if err != nil {
			return fmt.Errorf("the Alert Configuration(%s) does not exist", ids[alertconfiguration.EncodedIDKeyAlertID])
		}
		alert = alertResp
		return nil
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_alert_configuration" {
			continue
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		alert, _, err := acc.Conn().AlertConfigurations.GetAnAlertConfig(context.Background(), ids[alertconfiguration.EncodedIDKeyProjectID], ids[alertconfiguration.EncodedIDKeyAlertID])
		if alert != nil {
			return fmt.Errorf("the Project Alert Configuration(%s) still exists %s", ids[alertconfiguration.EncodedIDKeyAlertID], err)
		}
	}
	return nil
}

func importStateProjectIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s-%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["alert_configuration_id"]), nil
	}
}

func configBasicRS(orgID, projectName string, enabled bool) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "test" {
	name   = %[2]q
	org_id = %[1]q
}

resource "mongodbatlas_alert_configuration" "test" {
  project_id = mongodbatlas_project.test.id
  event_type = "OUTSIDE_METRIC_THRESHOLD"
  enabled    = "%[3]t"

  notification {
    type_name     = "GROUP"
    interval_min  = 5
    delay_min     = 0
    sms_enabled   = false
    email_enabled = true
    roles = ["GROUP_DATA_ACCESS_READ_ONLY", "GROUP_CLUSTER_MANAGER", "GROUP_DATA_ACCESS_ADMIN"]
  }

  notification {
    type_name     = "ORG"
    interval_min  = 5
    delay_min     = 0
    sms_enabled   = true
    email_enabled = false
  }

  matcher {
    field_name = "HOSTNAME_AND_PORT"
    operator   = "EQUALS"
    value      = "SECONDARY"
  }

  metric_threshold_config {
    metric_name = "ASSERT_REGULAR"
    operator    = "LESS_THAN"
    threshold   = 99.0
    units       = "RAW"
    mode        = "AVERAGE"
  }
}
	`, orgID, projectName, enabled)
}

func configWithNotifications(orgID, projectName string, enabled, smsEnabled, emailEnabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_alert_configuration" "test" {
			project_id = mongodbatlas_project.test.id
			event_type = "NO_PRIMARY"
			enabled    = "%[3]t"

			notification {
				type_name     = "GROUP"
				interval_min  = 5
				delay_min     = 0
				sms_enabled   = %[4]t
				email_enabled = %[5]t
				roles = ["GROUP_DATA_ACCESS_READ_ONLY"]
			}

			notification {
				type_name     = "ORG"
				interval_min  = 5
				delay_min     = 1
				sms_enabled   = %[4]t
				email_enabled = %[5]t
			}
		}
	`, orgID, projectName, enabled, smsEnabled, emailEnabled)
}

func configWithMatchers(orgID, projectName string, enabled, smsEnabled, emailEnabled bool, m1, m2 map[string]interface{}) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_alert_configuration" "test" {
			project_id = mongodbatlas_project.test.id
			event_type = "HOST_DOWN"
			enabled    = "%[3]t"

			notification {
				type_name     = "GROUP"
				interval_min  = 5
				delay_min     = 0
				sms_enabled   = %[4]t
				email_enabled = %[5]t
				roles = ["GROUP_DATA_ACCESS_READ_ONLY", "GROUP_CLUSTER_MANAGER"]
			}

			matcher {
				field_name = %[6]q
				operator   = %[7]q
				value      = %[8]q
			}
			matcher {
				field_name = %[9]q
				operator   = %[10]q
				value      = %[11]q
			}
		}
	`, orgID, projectName, enabled, smsEnabled, emailEnabled,
		m1["fieldName"], m1["operator"], m1["value"],
		m2["fieldName"], m2["operator"], m2["value"])
}

func configWithMetricUpdated(orgID, projectName string, enabled bool, threshold float64) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_alert_configuration" "test" {
			project_id = mongodbatlas_project.test.id
			event_type = "OUTSIDE_METRIC_THRESHOLD"
			enabled    = "%[3]t"

			notification {
				type_name     = "GROUP"
				interval_min  = 5
				delay_min     = 0
				sms_enabled   = false
				email_enabled = true
				roles = ["GROUP_DATA_ACCESS_READ_ONLY"]
			}

			matcher {
				field_name = "HOSTNAME_AND_PORT"
				operator   = "EQUALS"
				value      = "SECONDARY"
			}

			metric_threshold_config {
				metric_name = "ASSERT_REGULAR"
				operator    = "LESS_THAN"
				threshold   = %[4]f
				units       = "RAW"
				mode        = "AVERAGE"
			}
		}
	`, orgID, projectName, enabled, threshold)
}

func configWithoutRoles(orgID, projectName string, enabled bool, threshold float64) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_alert_configuration" "test" {
			project_id = mongodbatlas_project.test.id
			event_type = "OUTSIDE_METRIC_THRESHOLD"
			enabled    = "%[3]t"

			notification {
				type_name     = "EMAIL"
				email_address = "mongodbatlas.testing@gmail.com"
				interval_min  = 5
				delay_min     = 0
				sms_enabled   = false
				email_enabled = false
			}

			matcher {
				field_name = "HOSTNAME_AND_PORT"
				operator   = "EQUALS"
				value      = "SECONDARY"
			}

			metric_threshold_config {
				metric_name = "ASSERT_REGULAR"
				operator    = "LESS_THAN"
				threshold   = %[4]f
				units       = "RAW"
				mode        = "AVERAGE"
			}
		}
	`, orgID, projectName, enabled, threshold)
}

func configWithThresholdUpdated(orgID, projectName string, enabled bool, threshold float64) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_alert_configuration" "test" {
			project_id = mongodbatlas_project.test.id
			event_type = "REPLICATION_OPLOG_WINDOW_RUNNING_OUT"
			enabled    = "%[3]t"

			notification {
				type_name     = "GROUP"
				interval_min  = 5
				delay_min     = 0
				sms_enabled   = false
				email_enabled = true
				roles = ["GROUP_DATA_ACCESS_READ_ONLY", "GROUP_CLUSTER_MANAGER"]
			}

			matcher {
				field_name = "REPLICA_SET_NAME"
				operator   = "EQUALS"
				value      = "SECONDARY"
			}

			threshold_config {
				operator    = "LESS_THAN"
				units       = "HOURS"
				threshold   = %[4]f
			}
		}
	`, orgID, projectName, enabled, threshold)
}

func configWithDataDog(orgID, projectName, dataDogAPIKey, dataDogRegion string, enabled bool) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "test" {
	name   = %[2]q
	org_id = %[1]q
}
resource "mongodbatlas_third_party_integration" "atlas_datadog" {
  project_id = mongodbatlas_project.test.id
  type = "DATADOG"
  api_key = "%[4]s"
  region = "%[5]s"
}

resource "mongodbatlas_alert_configuration" "test" {
  project_id = mongodbatlas_project.test.id
  event_type = "REPLICATION_OPLOG_WINDOW_RUNNING_OUT"
  enabled    = %[3]t

  notification {
    type_name     = "GROUP"
    interval_min  = 5
    delay_min     = 0
    sms_enabled   = false
    email_enabled = true
    roles         = ["GROUP_OWNER"]
  }

  notification {
    type_name = "DATADOG"
    datadog_api_key = mongodbatlas_third_party_integration.atlas_datadog.api_key
    datadog_region = mongodbatlas_third_party_integration.atlas_datadog.region
    interval_min  = 5
    delay_min     = 0
  }

  matcher {
    field_name = "REPLICA_SET_NAME"
    operator   = "EQUALS"
    value      = "SECONDARY"
  }

  threshold_config {
    operator    = "LESS_THAN"
    threshold   = 72
    units       = "HOURS"
  }
}
	`, orgID, projectName, enabled, dataDogAPIKey, dataDogRegion)
}

func configWithPagerDuty(orgID, projectName, serviceKey string, enabled bool) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "test" {
	name   = %[2]q
	org_id = %[1]q
}
resource "mongodbatlas_alert_configuration" "test" {
  project_id = mongodbatlas_project.test.id
  event_type = "NO_PRIMARY"
  enabled    = "%[4]t"

  notification {
    type_name    = "PAGER_DUTY"
    service_key  = %[3]q
    delay_min    = 0
  }
}
	`, orgID, projectName, serviceKey, enabled)
}

func configWithPagerDutyNotifierID(orgID, projectName, notifierID string, delayMin int, serviceKey *string) string {
	var serviceKeyString string
	if serviceKey != nil {
		serviceKeyString = fmt.Sprintf(`service_key = %q`, *serviceKey)
	}
	return fmt.Sprintf(`
resource "mongodbatlas_project" "test" {
	name   = %[2]q
	org_id = %[1]q
}
resource "mongodbatlas_alert_configuration" "test" {
  project_id = mongodbatlas_project.test.id
  event_type = "NO_PRIMARY"
  enabled    = "true"

  notification {
    type_name    = "PAGER_DUTY"
    notifier_id  = %[3]q
	%[4]s
    delay_min    = %[5]d
  }
}
	`, orgID, projectName, notifierID, serviceKeyString, delayMin)
}

func configWithOpsGenie(orgID, projectName, apiKey string, enabled bool) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "test" {
	name   = %[2]q
	org_id = %[1]q
}
resource "mongodbatlas_alert_configuration" "test" {
  project_id = mongodbatlas_project.test.id
  event_type = "NO_PRIMARY"
  enabled    = "%[4]t"

  notification {
    type_name          = "OPS_GENIE"
    ops_genie_api_key  = %[3]q
    ops_genie_region   = "US"
    delay_min          = 0
  }
}
	`, orgID, projectName, apiKey, enabled)
}

func configWithVictorOps(orgID, projectName, apiKey string, enabled bool) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "test" {
	name   = %[2]q
	org_id = %[1]q
}
resource "mongodbatlas_alert_configuration" "test" {
  project_id = mongodbatlas_project.test.id
  event_type = "NO_PRIMARY"
  enabled    = "%[4]t"

  notification {
    type_name              = "VICTOR_OPS"
    victor_ops_api_key     = %[3]q
    victor_ops_routing_key = "testing"
    delay_min              = 0
  }
}
	`, orgID, projectName, apiKey, enabled)
}

func configWithEmptyMetricThresholdConfig(orgID, projectName string, enabled bool) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "test" {
	name   = %[2]q
	org_id = %[1]q
}

resource "mongodbatlas_alert_configuration" "test" {
  project_id = mongodbatlas_project.test.id
  event_type = "REPLICATION_OPLOG_WINDOW_RUNNING_OUT"
  enabled    = "%[3]t"

  notification {
    type_name     = "GROUP"
    interval_min  = 60
    delay_min     = 0
    sms_enabled   = true
    email_enabled = false
	roles         = ["GROUP_OWNER"]
  }

  threshold_config {
    operator    = "LESS_THAN"
    threshold   = 72
    units       = "HOURS"
  }

}
	`, orgID, projectName, enabled)
}

func configWithEmptyMatcherMetricThresholdConfig(orgID, projectName string, enabled bool) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "test" {
	name   = %[2]q
	org_id = %[1]q
}

resource "mongodbatlas_alert_configuration" "test" {
  project_id = mongodbatlas_project.test.id
  event_type = "CLUSTER_MONGOS_IS_MISSING"
  enabled    = "%[3]t"

  notification {
    type_name     = "GROUP"
    interval_min  = 60
    delay_min     = 0
    sms_enabled   = true
    email_enabled = false
	roles         = ["GROUP_OWNER"]
  }
}
	`, orgID, projectName, enabled)
}
