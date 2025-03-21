package eventtrigger_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb-labs/go-client-mongodb-atlas-app-services/appservices"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccEventTrigger_basic(t *testing.T) {
	acc.SkipTestForCI(t) // needs a project configured for triggers

	var (
		resourceName = "mongodbatlas_event_trigger.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		appID        = os.Getenv("MONGODB_APP_SERVICES_APP_ID")
	)
	event := appservices.EventTriggerRequest{
		Name:       acc.RandomName(),
		Type:       "DATABASE",
		FunctionID: os.Getenv("MONGODB_APP_SERVICES_FUNCTION_ID"),
		Disabled:   conversion.Pointer(true),
		Config: &appservices.EventTriggerConfig{
			OperationTypes: []string{"INSERT", "UPDATE"},
			Database:       "sample_airbnb",
			Collection:     "listingsAndReviews",
			ServiceID:      os.Getenv("MONGODB_APP_SERVICES_SERVICE_ID"),
			FullDocument:   conversion.Pointer(false),
		},
	}
	eventUpdated := appservices.EventTriggerRequest{
		Name:       acc.RandomName(),
		Type:       "DATABASE",
		FunctionID: os.Getenv("MONGODB_APP_SERVICES_FUNCTION_ID"),
		Disabled:   conversion.Pointer(true),
		Config: &appservices.EventTriggerConfig{
			OperationTypes: []string{"INSERT", "UPDATE", "DELETE"},
			Database:       "sample_airbnb",
			Collection:     "listingsAndReviews",
			ServiceID:      os.Getenv("MONGODB_APP_SERVICES_SERVICE_ID"),
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configDatabaseTrigger(projectID, appID, `"INSERT", "UPDATE"`, &event, false, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "disabled", "true"),
				),
			},
			{
				Config: configDatabaseTrigger(projectID, appID, `"INSERT", "UPDATE", "DELETE"`, &eventUpdated, true, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "disabled", "true"),
				),
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

func TestAccEventTrigger_databaseNoCollection(t *testing.T) {
	acc.SkipTestForCI(t) // needs a project configured for triggers

	var (
		resourceName = "mongodbatlas_event_trigger.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		appID        = os.Getenv("MONGODB_APP_SERVICES_APP_ID")
	)
	event := appservices.EventTriggerRequest{
		Name:       acc.RandomName(),
		Type:       "DATABASE",
		FunctionID: os.Getenv("MONGODB_APP_SERVICES_FUNCTION_ID"),
		Disabled:   conversion.Pointer(false),
		Config: &appservices.EventTriggerConfig{
			OperationTypes: []string{"INSERT", "UPDATE"},
			Database:       "sample_airbnb",
			Collection:     "listingsAndReviews",
			ServiceID:      os.Getenv("MONGODB_APP_SERVICES_SERVICE_ID"),
			FullDocument:   conversion.Pointer(false),
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configDatabaseNoCollectionTrigger(projectID, appID, `"INSERT", "UPDATE"`, &event),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "config_database", event.Config.Database),
					resource.TestCheckResourceAttr(resourceName, "config_collection", ""),
				),
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

func TestAccEventTrigger_databaseEventProccesor(t *testing.T) {
	acc.SkipTestForCI(t) // needs a project configured for triggers

	var (
		resourceName            = "mongodbatlas_event_trigger.test"
		projectID               = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		eventBridgeAwsAccountID = os.Getenv("AWS_EVENTBRIDGE_ACCOUNT_ID")
		eventBridgeAwsRegion    = conversion.MongoDBRegionToAWSRegion(os.Getenv("AWS_REGION"))
		appID                   = os.Getenv("MONGODB_APP_SERVICES_APP_ID")
	)
	event := appservices.EventTriggerRequest{
		Name:       acc.RandomName(),
		Type:       "DATABASE",
		FunctionID: os.Getenv("MONGODB_APP_SERVICES_FUNCTION_ID"),
		Disabled:   conversion.Pointer(false),
		Config: &appservices.EventTriggerConfig{
			OperationTypes: []string{"INSERT", "UPDATE"},
			Database:       "sample_airbnb",
			Collection:     "listingsAndReviews",
			ServiceID:      os.Getenv("MONGODB_APP_SERVICES_SERVICE_ID"),
			FullDocument:   conversion.Pointer(false),
		},
	}
	eventUpdated := appservices.EventTriggerRequest{
		Name:       acc.RandomName(),
		Type:       "DATABASE",
		FunctionID: os.Getenv("MONGODB_APP_SERVICES_FUNCTION_ID"),
		Disabled:   conversion.Pointer(false),
		Config: &appservices.EventTriggerConfig{
			OperationTypes: []string{"INSERT", "UPDATE", "DELETE"},
			Database:       "sample_airbnb",
			Collection:     "listingsAndReviews",
			ServiceID:      os.Getenv("MONGODB_APP_SERVICES_SERVICE_ID"),
			FullDocument:   conversion.Pointer(false),
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configDatabaseEPTrigger(projectID, appID, `"INSERT", "UPDATE"`, eventBridgeAwsAccountID, eventBridgeAwsRegion, &event),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
			{
				Config: configDatabaseEPTrigger(projectID, appID, `"INSERT", "UPDATE", "DELETE"`, eventBridgeAwsAccountID, eventBridgeAwsRegion, &eventUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
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

func TestAccEventTrigger_authBasic(t *testing.T) {
	acc.SkipTestForCI(t) // needs a project configured for triggers

	var (
		resourceName = "mongodbatlas_event_trigger.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		appID        = os.Getenv("MONGODB_APP_SERVICES_APP_ID")
	)
	event := appservices.EventTriggerRequest{
		Name:       acc.RandomName(),
		Type:       "AUTHENTICATION",
		FunctionID: os.Getenv("MONGODB_APP_SERVICES_FUNCTION_ID"),
		Disabled:   conversion.Pointer(false),
		Config: &appservices.EventTriggerConfig{
			OperationType: "LOGIN",
			Providers:     []string{"anon-user", "local-userpass"},
		},
	}
	eventUpdated := appservices.EventTriggerRequest{
		Name:       acc.RandomName(),
		Type:       "AUTHENTICATION",
		FunctionID: os.Getenv("MONGODB_APP_SERVICES_FUNCTION_ID"),
		Disabled:   conversion.Pointer(false),
		Config: &appservices.EventTriggerConfig{
			OperationTypes: []string{"INSERT", "UPDATE", "DELETE"},
			OperationType:  "LOGIN",
			Providers:      []string{"anon-user", "local-userpass"},
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configAuthenticationTrigger(projectID, appID, `"anon-user", "local-userpass"`, &event),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
			{
				Config: configAuthenticationTrigger(projectID, appID, `"anon-user", "local-userpass", "api-key"`, &eventUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
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

func TestAccEventTrigger_authEventProcessor(t *testing.T) {
	acc.SkipTestForCI(t) // needs a project configured for triggers

	var (
		resourceName            = "mongodbatlas_event_trigger.test"
		projectID               = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		eventBridgeAwsAccountID = os.Getenv("AWS_EVENTBRIDGE_ACCOUNT_ID")
		eventBridgeAwsRegion    = conversion.MongoDBRegionToAWSRegion(os.Getenv("AWS_REGION"))
		appID                   = os.Getenv("MONGODB_APP_SERVICES_APP_ID")
	)
	event := appservices.EventTriggerRequest{
		Name:       acc.RandomName(),
		Type:       "AUTHENTICATION",
		FunctionID: os.Getenv("MONGODB_APP_SERVICES_FUNCTION_ID"),
		Disabled:   conversion.Pointer(false),
		Config: &appservices.EventTriggerConfig{
			OperationType: "LOGIN",
			Providers:     []string{"anon-user", "local-userpass"},
		},
	}
	eventUpdated := appservices.EventTriggerRequest{
		Name:       acc.RandomName(),
		Type:       "AUTHENTICATION",
		FunctionID: os.Getenv("MONGODB_APP_SERVICES_FUNCTION_ID"),
		Disabled:   conversion.Pointer(false),
		Config: &appservices.EventTriggerConfig{
			OperationTypes: []string{"INSERT", "UPDATE", "DELETE"},
			OperationType:  "LOGIN",
			Providers:      []string{"anon-user", "local-userpass"},
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configAuthenticationEPTrigger(projectID, appID, `"anon-user", "local-userpass"`, eventBridgeAwsAccountID, eventBridgeAwsRegion, &event),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
			{
				Config: configAuthenticationEPTrigger(projectID, appID, `"anon-user", "local-userpass", "api-key"`, eventBridgeAwsAccountID, eventBridgeAwsRegion, &eventUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
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

func TestAccEventTrigger_scheduleBasic(t *testing.T) {
	acc.SkipTestForCI(t) // needs a project configured for triggers

	var (
		resourceName = "mongodbatlas_event_trigger.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		appID        = os.Getenv("MONGODB_APP_SERVICES_APP_ID")
	)
	event := appservices.EventTriggerRequest{
		Name:       acc.RandomName(),
		Type:       "SCHEDULED",
		FunctionID: os.Getenv("MONGODB_APP_SERVICES_FUNCTION_ID"),
		Disabled:   conversion.Pointer(false),
		Config: &appservices.EventTriggerConfig{
			Schedule: "* * * * *",
		},
	}
	eventUpdated := appservices.EventTriggerRequest{
		Name:       acc.RandomName(),
		Type:       "SCHEDULED",
		FunctionID: os.Getenv("MONGODB_APP_SERVICES_FUNCTION_ID"),
		Disabled:   conversion.Pointer(false),
		Config: &appservices.EventTriggerConfig{
			Schedule: "* * * * *",
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configScheduleTrigger(projectID, appID, &event),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
			{
				Config: configScheduleTrigger(projectID, appID, &eventUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
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

func TestAccEventTrigger_scheduleEventProcessor(t *testing.T) {
	acc.SkipTestForCI(t) // needs a project configured for triggers

	var (
		resourceName            = "mongodbatlas_event_trigger.test"
		projectID               = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		eventBridgeAwsAccountID = os.Getenv("AWS_EVENTBRIDGE_ACCOUNT_ID")
		eventBridgeAwsRegion    = conversion.MongoDBRegionToAWSRegion(os.Getenv("AWS_REGION"))
		appID                   = os.Getenv("MONGODB_APP_SERVICES_APP_ID")
	)
	event := appservices.EventTriggerRequest{
		Name:       acc.RandomName(),
		Type:       "SCHEDULED",
		FunctionID: os.Getenv("MONGODB_APP_SERVICES_FUNCTION_ID"),
		Disabled:   conversion.Pointer(false),
		Config: &appservices.EventTriggerConfig{
			Schedule: "* * * * *",
		},
	}
	eventUpdated := appservices.EventTriggerRequest{
		Name:       acc.RandomName(),
		Type:       "SCHEDULED",
		FunctionID: os.Getenv("MONGODB_APP_SERVICES_FUNCTION_ID"),
		Disabled:   conversion.Pointer(false),
		Config: &appservices.EventTriggerConfig{
			Schedule: "* * * * *",
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configScheduleEPTrigger(projectID, appID, eventBridgeAwsAccountID, eventBridgeAwsRegion, &event),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
			{
				Config: configScheduleEPTrigger(projectID, appID, eventBridgeAwsAccountID, eventBridgeAwsRegion, &eventUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
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

func TestAccEventTrigger_functionBasic(t *testing.T) {
	acc.SkipTestForCI(t) // needs a project configured for triggers

	var (
		resourceName = "mongodbatlas_event_trigger.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		appID        = os.Getenv("MONGODB_APP_SERVICES_APP_ID")
	)
	event := appservices.EventTriggerRequest{
		Name:       acc.RandomName(),
		Type:       "SCHEDULED",
		FunctionID: os.Getenv("MONGODB_APP_SERVICES_FUNCTION_ID"),
		Disabled:   conversion.Pointer(false),
		Config: &appservices.EventTriggerConfig{
			Schedule: "0 8 * * *",
		},
	}
	eventUpdated := appservices.EventTriggerRequest{
		Name:       acc.RandomName(),
		Type:       "SCHEDULED",
		FunctionID: os.Getenv("MONGODB_APP_SERVICES_FUNCTION_ID"),
		Disabled:   conversion.Pointer(false),
		Config: &appservices.EventTriggerConfig{
			Schedule: "0 8 * * *",
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configScheduleTrigger(projectID, appID, &event),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
			{
				Config: configScheduleTrigger(projectID, appID, &eventUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
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

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := context.Background()
		conn, err := acc.MongoDBClient.GetAppServicesClient(ctx)
		if err != nil {
			return err
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)

		if ids["trigger_id"] == "" {
			return fmt.Errorf("no ID is set")
		}

		log.Printf("[DEBUG] trigger_id ID: %s", ids["trigger_id"])

		_, _, err = conn.EventTriggers.Get(ctx, ids["project_id"], ids["app_id"], ids["trigger_id"])
		if err == nil {
			return nil
		}

		return fmt.Errorf("cloudProviderSnapshot (%s) does not exist", rs.Primary.Attributes["snapshot_id"])
	}
}

func checkDestroy(s *terraform.State) error {
	ctx := context.Background()
	conn, err := acc.MongoDBClient.GetAppServicesClient(ctx)
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_event_trigger" {
			continue
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)

		res, _, _ := conn.EventTriggers.Get(ctx, ids["project_id"], ids["app_id"], ids["trigger_id"])

		if res != nil {
			return fmt.Errorf("event trigger (%s) still exists", ids["trigger_id"])
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

		return fmt.Sprintf("%s--%s--%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["app_id"], rs.Primary.Attributes["trigger_id"]), nil
	}
}

func configDatabaseTrigger(projectID, appID, operationTypes string, eventTrigger *appservices.EventTriggerRequest, fullDoc, fullDocBefore bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_event_trigger" "test" {
			project_id = %[1]q
			app_id = %[2]q
			name = %[3]q
			type = %[4]q
			function_id = %[5]q
			disabled = %[6]t
			config_operation_types = [%s]
			config_database = %[8]q
			config_collection = %[9]q
			config_service_id = %[10]q
			config_full_document = %[11]t
			config_full_document_before = %[12]t
			config_match = <<-EOF
			{
			  "updateDescription.updatedFields": {
				"status": "blocked"
			  }
			}
			EOF		
}
	`, projectID, appID, eventTrigger.Name, eventTrigger.Type, eventTrigger.FunctionID, *eventTrigger.Disabled, operationTypes,
		eventTrigger.Config.Database, eventTrigger.Config.Collection,
		eventTrigger.Config.ServiceID, fullDoc, fullDocBefore)
}

func configDatabaseNoCollectionTrigger(projectID, appID, operationTypes string, eventTrigger *appservices.EventTriggerRequest) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_event_trigger" "test" {
			project_id = %[1]q
			app_id = %[2]q
			name = %[3]q
			type = %[4]q
			function_id = %[5]q
			disabled = %[6]t
			config_operation_types = [%[7]s]
			config_database = %[8]q
			config_service_id = %[9]q
			config_full_document = false
		}
	`, projectID, appID, eventTrigger.Name, eventTrigger.Type, eventTrigger.FunctionID, *eventTrigger.Disabled, operationTypes,
		eventTrigger.Config.Database, eventTrigger.Config.ServiceID)
}

func configDatabaseEPTrigger(projectID, appID, operationTypes, awsAccID, awsRegion string, eventTrigger *appservices.EventTriggerRequest) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_event_trigger" "test" {
			project_id = %[1]q
			app_id = %[2]q
			name = %[3]q
			type = %[4]q
			disabled = %[5]t
			config_operation_types = [%[6]s]
			config_database = %[7]q
			config_collection = %[8]q
			config_service_id = %[9]q
			config_match = "{\"updateDescription.updatedFields\":{\"status\":\"blocked\"}}"
			event_processors{
				aws_eventbridge{
					config_account_id = %[10]q
					config_region = %[11]q
				}

			}
		}
	`, projectID, appID, eventTrigger.Name, eventTrigger.Type, *eventTrigger.Disabled, operationTypes,
		eventTrigger.Config.Database, eventTrigger.Config.Collection,
		eventTrigger.Config.ServiceID, awsAccID, awsRegion)
}

func configAuthenticationTrigger(projectID, appID, providers string, eventTrigger *appservices.EventTriggerRequest) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_event_trigger" "test" {
			project_id = %[1]q
			app_id = %[2]q
			name = %[3]q
			type = %[4]q
			function_id = %[5]q
			disabled = %[6]t
			config_operation_type = %[7]q
			config_providers = [%[8]s]
		}
	`, projectID, appID, eventTrigger.Name, eventTrigger.Type, eventTrigger.FunctionID, *eventTrigger.Disabled,
		eventTrigger.Config.OperationType, providers)
}

func configAuthenticationEPTrigger(projectID, appID, providers, awsAccID, awsRegion string, eventTrigger *appservices.EventTriggerRequest) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_event_trigger" "test" {
			project_id = %[1]q
			app_id = %[2]q
			name = %[3]q
			type = %[4]q
			disabled = %[5]t
			config_operation_type = %[6]q
			config_providers = [%[7]s]
			event_processors{
				aws_eventbridge{
					config_account_id = %[8]q
					config_region = %[9]q
				}
			}
		}
	`, projectID, appID, eventTrigger.Name, eventTrigger.Type, *eventTrigger.Disabled, eventTrigger.Config.OperationType, providers,
		awsAccID, awsRegion)
}

func configScheduleTrigger(projectID, appID string, eventTrigger *appservices.EventTriggerRequest) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_event_trigger" "test" {
			project_id = %[1]q
			app_id = %[2]q
			name = %[3]q
			type = %[4]q
			function_id = %[5]q
			disabled = %[6]t
			config_schedule = %[7]q
		}
	`, projectID, appID, eventTrigger.Name, eventTrigger.Type, eventTrigger.FunctionID, *eventTrigger.Disabled,
		eventTrigger.Config.Schedule)
}

func configScheduleEPTrigger(projectID, appID, awsAccID, awsRegion string, eventTrigger *appservices.EventTriggerRequest) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_event_trigger" "test" {
			project_id = %[1]q
			app_id = %[2]q
			name = %[3]q
			type = %[4]q
			disabled = %[5]t
			config_schedule = %[6]q
			event_processors{
				aws_eventbridge{
					config_account_id = %[7]q
					config_region = %[8]q
				}
			}
		}
	`, projectID, appID, eventTrigger.Name, eventTrigger.Type, *eventTrigger.Disabled, eventTrigger.Config.Schedule,
		awsAccID, awsRegion)
}
