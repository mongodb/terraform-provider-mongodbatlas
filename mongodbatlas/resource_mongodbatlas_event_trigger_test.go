package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mwielbut/pointy"
	"go.mongodb.org/realm/realm"
)

func TestAccResourceMongoDBAtlasEventTriggerDatabase_basic(t *testing.T) {
	SkipTest(t)
	var (
		resourceName = "mongodbatlas_event_trigger.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		appID        = os.Getenv("MONGODB_REALM_APP_ID")
		eventResp    = realm.EventTrigger{}
	)
	event := realm.EventTriggerRequest{
		Name:       acctest.RandomWithPrefix("test-acc"),
		Type:       "DATABASE",
		FunctionID: os.Getenv("MONGODB_REALM_FUNCTION_ID"),
		Disabled:   pointy.Bool(false),
		Config: &realm.EventTriggerConfig{
			OperationTypes: []string{"INSERT", "UPDATE"},
			Database:       "sample_airbnb",
			Collection:     "listingsAndReviews",
			ServiceID:      os.Getenv("MONGODB_REALM_SERVICE_ID"),
			FullDocument:   pointy.Bool(false),
		},
	}
	eventUpdated := realm.EventTriggerRequest{
		Name:       acctest.RandomWithPrefix("test-acc"),
		Type:       "DATABASE",
		FunctionID: os.Getenv("MONGODB_REALM_FUNCTION_ID"),
		Disabled:   pointy.Bool(false),
		Config: &realm.EventTriggerConfig{
			OperationTypes: []string{"INSERT", "UPDATE", "DELETE"},
			Database:       "sample_airbnb",
			Collection:     "listingsAndReviews",
			ServiceID:      os.Getenv("MONGODB_REALM_SERVICE_ID"),
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasEventTriggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasEventTriggerDatabaseConfigDatabase(projectID, appID, `"INSERT", "UPDATE"`, &event, false, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEventTriggerExists(resourceName, &eventResp),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
			{
				Config: testAccMongoDBAtlasEventTriggerDatabaseConfigDatabase(projectID, appID, `"INSERT", "UPDATE", "DELETE"`, &eventUpdated, true, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEventTriggerExists(resourceName, &eventResp),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasEventTriggerImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceMongoDBAtlasEventTriggerDatabase_eventProccesor(t *testing.T) {
	SkipTest(t)
	var (
		resourceName = "mongodbatlas_event_trigger.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		awsAccountID = os.Getenv("AWS_ACCOUNT_ID")
		awsRegion    = os.Getenv("AWS_REGION")
		appID        = os.Getenv("MONGODB_REALM_APP_ID")
		eventResp    = realm.EventTrigger{}
	)
	event := realm.EventTriggerRequest{
		Name:       acctest.RandomWithPrefix("test-acc"),
		Type:       "DATABASE",
		FunctionID: os.Getenv("MONGODB_REALM_FUNCTION_ID"),
		Disabled:   pointy.Bool(false),
		Config: &realm.EventTriggerConfig{
			OperationTypes: []string{"INSERT", "UPDATE"},
			Database:       "sample_airbnb",
			Collection:     "listingsAndReviews",
			ServiceID:      os.Getenv("MONGODB_REALM_SERVICE_ID"),
			FullDocument:   pointy.Bool(false),
		},
	}
	eventUpdated := realm.EventTriggerRequest{
		Name:       acctest.RandomWithPrefix("test-acc"),
		Type:       "DATABASE",
		FunctionID: os.Getenv("MONGODB_REALM_FUNCTION_ID"),
		Disabled:   pointy.Bool(false),
		Config: &realm.EventTriggerConfig{
			OperationTypes: []string{"INSERT", "UPDATE", "DELETE"},
			Database:       "sample_airbnb",
			Collection:     "listingsAndReviews",
			ServiceID:      os.Getenv("MONGODB_REALM_SERVICE_ID"),
			FullDocument:   pointy.Bool(false),
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasEventTriggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasEventTriggerDatabaseConfigDatabaseEP(projectID, appID, `"INSERT", "UPDATE"`, awsAccountID, awsRegion, &event),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEventTriggerExists(resourceName, &eventResp),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
			{
				Config: testAccMongoDBAtlasEventTriggerDatabaseConfigDatabaseEP(projectID, appID, `"INSERT", "UPDATE", "DELETE"`, awsAccountID, awsRegion, &eventUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEventTriggerExists(resourceName, &eventResp),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasEventTriggerImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceMongoDBAtlasEventTriggerAuth_basic(t *testing.T) {
	SkipTest(t)
	var (
		resourceName = "mongodbatlas_event_trigger.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		appID        = os.Getenv("MONGODB_REALM_APP_ID")
		eventResp    = realm.EventTrigger{}
	)
	event := realm.EventTriggerRequest{
		Name:       acctest.RandomWithPrefix("test-acc"),
		Type:       "AUTHENTICATION",
		FunctionID: os.Getenv("MONGODB_REALM_FUNCTION_ID"),
		Disabled:   pointy.Bool(false),
		Config: &realm.EventTriggerConfig{
			OperationType: "LOGIN",
			Providers:     []string{"anon-user", "local-userpass"},
		},
	}
	eventUpdated := realm.EventTriggerRequest{
		Name:       acctest.RandomWithPrefix("test-acc"),
		Type:       "AUTHENTICATION",
		FunctionID: os.Getenv("MONGODB_REALM_FUNCTION_ID"),
		Disabled:   pointy.Bool(false),
		Config: &realm.EventTriggerConfig{
			OperationTypes: []string{"INSERT", "UPDATE", "DELETE"},
			OperationType:  "LOGIN",
			Providers:      []string{"anon-user", "local-userpass"},
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasEventTriggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasEventTriggerDatabaseConfigAuthentication(projectID, appID, `"anon-user", "local-userpass"`, &event),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEventTriggerExists(resourceName, &eventResp),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
			{
				Config: testAccMongoDBAtlasEventTriggerDatabaseConfigAuthentication(projectID, appID, `"anon-user", "local-userpass", "api-key"`, &eventUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEventTriggerExists(resourceName, &eventResp),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasEventTriggerImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceMongoDBAtlasEventTriggerAuth_eventProcessor(t *testing.T) {
	SkipTest(t)
	var (
		resourceName = "mongodbatlas_event_trigger.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		awsAccountID = os.Getenv("AWS_ACCOUNT_ID")
		awsRegion    = os.Getenv("AWS_REGION")
		appID        = os.Getenv("MONGODB_REALM_APP_ID")
		eventResp    = realm.EventTrigger{}
	)
	event := realm.EventTriggerRequest{
		Name:       acctest.RandomWithPrefix("test-acc"),
		Type:       "AUTHENTICATION",
		FunctionID: os.Getenv("MONGODB_REALM_FUNCTION_ID"),
		Disabled:   pointy.Bool(false),
		Config: &realm.EventTriggerConfig{
			OperationType: "LOGIN",
			Providers:     []string{"anon-user", "local-userpass"},
		},
	}
	eventUpdated := realm.EventTriggerRequest{
		Name:       acctest.RandomWithPrefix("test-acc"),
		Type:       "AUTHENTICATION",
		FunctionID: os.Getenv("MONGODB_REALM_FUNCTION_ID"),
		Disabled:   pointy.Bool(false),
		Config: &realm.EventTriggerConfig{
			OperationTypes: []string{"INSERT", "UPDATE", "DELETE"},
			OperationType:  "LOGIN",
			Providers:      []string{"anon-user", "local-userpass"},
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasEventTriggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasEventTriggerDatabaseConfigAuthenticationEP(projectID, appID, `"anon-user", "local-userpass"`, awsAccountID, awsRegion, &event),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEventTriggerExists(resourceName, &eventResp),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
			{
				Config: testAccMongoDBAtlasEventTriggerDatabaseConfigAuthenticationEP(projectID, appID, `"anon-user", "local-userpass", "api-key"`, awsAccountID, awsRegion, &eventUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEventTriggerExists(resourceName, &eventResp),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasEventTriggerImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceMongoDBAtlasEventTriggerSchedule_basic(t *testing.T) {
	SkipTest(t)
	var (
		resourceName = "mongodbatlas_event_trigger.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		appID        = os.Getenv("MONGODB_REALM_APP_ID")
		eventResp    = realm.EventTrigger{}
	)
	event := realm.EventTriggerRequest{
		Name:       acctest.RandomWithPrefix("test-acc"),
		Type:       "SCHEDULED",
		FunctionID: os.Getenv("MONGODB_REALM_FUNCTION_ID"),
		Disabled:   pointy.Bool(false),
		Config: &realm.EventTriggerConfig{
			Schedule: "* * * * *",
		},
	}
	eventUpdated := realm.EventTriggerRequest{
		Name:       acctest.RandomWithPrefix("test-acc"),
		Type:       "SCHEDULED",
		FunctionID: os.Getenv("MONGODB_REALM_FUNCTION_ID"),
		Disabled:   pointy.Bool(false),
		Config: &realm.EventTriggerConfig{
			Schedule: "* * * * *",
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasEventTriggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasEventTriggerDatabaseConfigSchedule(projectID, appID, &event),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEventTriggerExists(resourceName, &eventResp),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
			{
				Config: testAccMongoDBAtlasEventTriggerDatabaseConfigSchedule(projectID, appID, &eventUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEventTriggerExists(resourceName, &eventResp),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasEventTriggerImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceMongoDBAtlasEventTriggerSchedule_eventProcessor(t *testing.T) {
	SkipTest(t)
	var (
		resourceName = "mongodbatlas_event_trigger.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		awsAccountID = os.Getenv("AWS_ACCOUNT_ID")
		awsRegion    = os.Getenv("AWS_REGION")
		appID        = os.Getenv("MONGODB_REALM_APP_ID")
		eventResp    = realm.EventTrigger{}
	)
	event := realm.EventTriggerRequest{
		Name:       acctest.RandomWithPrefix("test-acc"),
		Type:       "SCHEDULED",
		FunctionID: os.Getenv("MONGODB_REALM_FUNCTION_ID"),
		Disabled:   pointy.Bool(false),
		Config: &realm.EventTriggerConfig{
			Schedule: "* * * * *",
		},
	}
	eventUpdated := realm.EventTriggerRequest{
		Name:       acctest.RandomWithPrefix("test-acc"),
		Type:       "SCHEDULED",
		FunctionID: os.Getenv("MONGODB_REALM_FUNCTION_ID"),
		Disabled:   pointy.Bool(false),
		Config: &realm.EventTriggerConfig{
			Schedule: "* * * * *",
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasEventTriggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasEventTriggerDatabaseConfigScheduleEP(projectID, appID, awsAccountID, awsRegion, &event),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEventTriggerExists(resourceName, &eventResp),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
			{
				Config: testAccMongoDBAtlasEventTriggerDatabaseConfigScheduleEP(projectID, appID, awsAccountID, awsRegion, &eventUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEventTriggerExists(resourceName, &eventResp),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasEventTriggerImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceMongoDBAtlasEventTriggerFunction_basic(t *testing.T) {
	SkipTest(t)
	var (
		resourceName = "mongodbatlas_event_trigger.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		appID        = os.Getenv("MONGODB_REALM_APP_ID")
		eventResp    = realm.EventTrigger{}
	)
	event := realm.EventTriggerRequest{
		Name:       acctest.RandomWithPrefix("test-acc"),
		Type:       "SCHEDULED",
		FunctionID: os.Getenv("MONGODB_REALM_FUNCTION_ID"),
		Disabled:   pointy.Bool(false),
		Config: &realm.EventTriggerConfig{
			Schedule: "0 8 * * *",
		},
	}
	eventUpdated := realm.EventTriggerRequest{
		Name:       acctest.RandomWithPrefix("test-acc"),
		Type:       "SCHEDULED",
		FunctionID: os.Getenv("MONGODB_REALM_FUNCTION_ID"),
		Disabled:   pointy.Bool(false),
		Config: &realm.EventTriggerConfig{
			Schedule: "0 8 * * *",
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasEventTriggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasEventTriggerFunctionConfig(projectID, appID, &event),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEventTriggerExists(resourceName, &eventResp),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
			{
				Config: testAccMongoDBAtlasEventTriggerFunctionConfig(projectID, appID, &eventUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEventTriggerExists(resourceName, &eventResp),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasEventTriggerImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMongoDBAtlasEventTriggerExists(resourceName string, eventTrigger *realm.EventTrigger) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := context.Background()
		conn, err := testAccProvider.Meta().(*MongoDBClient).GetRealmClient(ctx)
		if err != nil {
			return err
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		ids := decodeStateID(rs.Primary.ID)

		if ids["trigger_id"] == "" {
			return fmt.Errorf("no ID is set")
		}

		log.Printf("[DEBUG] trigger_id ID: %s", ids["trigger_id"])

		res, _, err := conn.EventTriggers.Get(ctx, ids["project_id"], ids["app_id"], ids["trigger_id"])
		if err == nil {
			*eventTrigger = *res
			return nil
		}

		return fmt.Errorf("cloudProviderSnapshot (%s) does not exist", rs.Primary.Attributes["snapshot_id"])
	}
}

func testAccCheckMongoDBAtlasEventTriggerDestroy(s *terraform.State) error {
	ctx := context.Background()
	conn, err := testAccProvider.Meta().(*MongoDBClient).GetRealmClient(ctx)
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_event_trigger" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)

		res, _, _ := conn.EventTriggers.Get(ctx, ids["project_id"], ids["app_id"], ids["trigger_id"])

		if res != nil {
			return fmt.Errorf("event trigger (%s) still exists", ids["trigger_id"])
		}
	}

	return nil
}

func testAccCheckMongoDBAtlasEventTriggerImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return fmt.Sprintf("%s--%s--%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["app_id"], rs.Primary.Attributes["trigger_id"]), nil
	}
}

func testAccMongoDBAtlasEventTriggerDatabaseConfigDatabase(projectID, appID, operationTypes string, eventTrigger *realm.EventTriggerRequest, fullDoc, fullDocBefore bool) string {
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

func testAccMongoDBAtlasEventTriggerDatabaseConfigDatabaseEP(projectID, appID, operationTypes, awsAccID, awsRegion string, eventTrigger *realm.EventTriggerRequest) string {
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

func testAccMongoDBAtlasEventTriggerDatabaseConfigAuthentication(projectID, appID, providers string, eventTrigger *realm.EventTriggerRequest) string {
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

func testAccMongoDBAtlasEventTriggerDatabaseConfigAuthenticationEP(projectID, appID, providers, awsAccID, awsRegion string, eventTrigger *realm.EventTriggerRequest) string {
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

func testAccMongoDBAtlasEventTriggerDatabaseConfigSchedule(projectID, appID string, eventTrigger *realm.EventTriggerRequest) string {
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

func testAccMongoDBAtlasEventTriggerDatabaseConfigScheduleEP(projectID, appID, awsAccID, awsRegion string, eventTrigger *realm.EventTriggerRequest) string {
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

func testAccMongoDBAtlasEventTriggerFunctionConfig(projectID, appID string, eventTrigger *realm.EventTriggerRequest) string {
	return fmt.Sprintf(`
resource "mongodbatlas_event_trigger" "test" {
  project_id      = %[1]q
  app_id          = %[2]q
  name            = %[3]q
  type            = %[4]q
  function_id     = %[5]q
  disabled        = %[6]t
  config_schedule = %[7]q
}
	`, projectID, appID, eventTrigger.Name, eventTrigger.Type, eventTrigger.FunctionID,
		*eventTrigger.Disabled, eventTrigger.Config.Schedule)
}
