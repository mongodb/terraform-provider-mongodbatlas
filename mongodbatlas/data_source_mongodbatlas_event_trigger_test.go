package mongodbatlas_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mwielbut/pointy"
	"go.mongodb.org/realm/realm"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestEventTrigger_basic(t *testing.T) {
	acc.SkipTestForCI(t)
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
			Unordered:      pointy.Bool(true),
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasEventTriggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataSourceEventTriggerConfig(projectID, appID, `"INSERT", "UPDATE"`, &event),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEventTriggerExists(resourceName, &eventResp),
				),
			},
		},
	})
}

func testAccMongoDBAtlasDataSourceEventTriggerConfig(projectID, appID, operationTypes string, eventTrigger *realm.EventTriggerRequest) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_event_trigger" "test" {
			project_id = %[1]q
			app_id = %[2]q
			name = %[3]q
			type = %[4]q
			function_id = %[5]q
			disabled = %[6]t
			unordered = %[7]t
			config_operation_types = [%s]
			config_database = %[8]q
			config_collection = %[9]q
			config_service_id = %[10]q
			config_match = "{\"updateDescription.updatedFields\":{\"status\":\"blocked\"}}"
		}

		data "mongodbatlas_event_trigger" "test" {
			project_id = mongodbatlas_event_trigger.test.project_id
			app_id = mongodbatlas_event_trigger.test.app_id
			trigger_id = mongodbatlas_event_trigger.test.id
		}
`, projectID, appID, eventTrigger.Name, eventTrigger.Type, eventTrigger.FunctionID, *eventTrigger.Disabled, *eventTrigger.Config.Unordered, operationTypes,
		eventTrigger.Config.Database, eventTrigger.Config.Collection,
		eventTrigger.Config.ServiceID)
}
