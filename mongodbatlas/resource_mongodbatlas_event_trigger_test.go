package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/mwielbut/pointy"
	"go.mongodb.org/realm/realm"
)

func TestAccResourceMongoDBAtlasEventTrigger_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_event_trigger.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		awsAccountID = os.Getenv("AWS_ACCOUNT_ID")
		awsRegion    = os.Getenv("AWS_REGION")
		appID        = "testing-edgar-utlvf"
		eventResp    = realm.EventTrigger{}
	)
	event := realm.EventTriggerRequest{
		Name:       acctest.RandomWithPrefix("test-acc"),
		Type:       "DATABASE",
		FunctionID: "1",
		Disabled:   pointy.Bool(false),
		Config: &realm.EventTriggerConfig{
			OperationTypes: []interface{}{"INSERT", "UPDATE"},
			OperationType:  "LOGIN",
			Providers:      "anon-user",
			Database:       "database",
			Collection:     "collection",
			ServiceID:      "1",
			Match: map[string]interface{}{
				"expr": "something",
			},
			FullDocument: pointy.Bool(false),
			Schedule:     "*",
		},
		EventProcessors: map[string]interface{}{
			"AWS_EVENTBRIDGE": map[string]interface{}{
				"type": "AWS_EVENTBRIDGE",
				"config": map[string]interface{}{
					"account_id": awsAccountID,
					"region":     awsRegion,
				},
			},
		},
	}
	eventUpdated := realm.EventTriggerRequest{
		Name:       acctest.RandomWithPrefix("test-acc-updated"),
		Type:       "DATABASE",
		FunctionID: "1",
		Disabled:   pointy.Bool(false),
		Config: &realm.EventTriggerConfig{
			OperationTypes: []interface{}{"INSERT", "UPDATE"},
			OperationType:  "LOGIN",
			Providers:      "anon-user",
			Database:       "database",
			Collection:     "collection",
			ServiceID:      "1",
			Match: map[string]interface{}{
				"expr": "something",
			},
			FullDocument: pointy.Bool(false),
			Schedule:     "*",
		},
		EventProcessors: map[string]interface{}{
			"AWS_EVENTBRIDGE": map[string]interface{}{
				"type": "AWS_EVENTBRIDGE",
				"config": map[string]interface{}{
					"account_id": awsAccountID,
					"region":     awsRegion,
				},
			},
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasEventTriggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasEventTriggerConfig(projectID, appID, "expr", "something", awsAccountID, awsRegion, &event),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEventTriggerExists(resourceName, &eventResp),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
			{
				Config: testAccMongoDBAtlasEventTriggerConfig(projectID, appID, "expr", "something", awsAccountID, awsRegion, &eventUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEventTriggerExists(resourceName, &eventResp),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
		},
	})
}

func testAccCheckMongoDBAtlasEventTriggerExists(resourceName string, eventTrigger *realm.EventTrigger) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Realm

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		ids := decodeStateID(rs.Primary.ID)

		if ids["trigger_id"] == "" {
			return fmt.Errorf("no ID is set")
		}

		log.Printf("[DEBUG] trigger_id ID: %s", ids["trigger_id"])

		res, _, err := conn.EventTriggers.Get(context.Background(), ids["project_id"], ids["app_id"], ids["trigger_id"])
		if err == nil {
			*eventTrigger = *res
			return nil
		}

		return fmt.Errorf("cloudProviderSnapshot (%s) does not exist", rs.Primary.Attributes["snapshot_id"])
	}
}

func testAccCheckMongoDBAtlasEventTriggerDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*MongoDBClient).Realm

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_event_trigger" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)

		res, _, _ := conn.EventTriggers.Get(context.Background(), ids["project_id"], ids["app_id"], ids["trigger_id"])

		if res != nil {
			return fmt.Errorf("event trigger (%s) still exists", ids["trigger_id"])
		}

	}

	return nil
}

func testAccMongoDBAtlasEventTriggerConfig(projectID, appID, matchKey, matchValue, awsAccID, awsRegion string, eventTrigger *realm.EventTriggerRequest) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_event_trigger" "test" {
			project_id = %[1]q
			app_id = %[2]q
			name = %[3]q
			type = %[4]q
			function_id = %[5]q
			disabled = %[6]t
			config_operation_types = ["INSERT", "UPDATE"]
			config_operation_type = %[7]q
			config_providers = %[8]q
			config_database = %[9]q
			config_collection = %[10]q
			config_service_id = %[11]q
			config_full_document = %[14]t
			config_schedule = %[15]q
			event_processors {
				aws_eventbridge {
					config_account_id = %[16]q
					config_region = %[17]q
				}
			}

		}
	`, projectID, appID, eventTrigger.Name, eventTrigger.Type, eventTrigger.FunctionID, *eventTrigger.Disabled,
		eventTrigger.Config.OperationType, eventTrigger.Config.Providers, eventTrigger.Config.Database, eventTrigger.Config.Collection,
		eventTrigger.Config.ServiceID, matchKey, matchValue, *eventTrigger.Config.FullDocument, eventTrigger.Config.Schedule,
		awsAccID, awsRegion)
}
