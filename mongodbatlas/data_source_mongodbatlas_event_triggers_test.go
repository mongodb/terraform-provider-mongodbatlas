package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/mwielbut/pointy"
	"go.mongodb.org/realm/realm"
)

func TestAccDataSourceMongoDBAtlasEventTriggers_basic(t *testing.T) {
	var (
		//resourceName = "data.mongodbatlas_event_triggers.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		awsAccountID = os.Getenv("AWS_ACCOUNT_ID")
		awsRegion    = os.Getenv("AWS_REGION")
		appID        = "testing-edgar-utlvf"
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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasEventTriggersDataSourceConfig(projectID, appID, "expr", "something", awsAccountID, awsRegion, &event),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testAccMongoDBAtlasEventTriggersDataSourceConfig(projectID, appID, matchKey, matchValue, awsAccID, awsRegion string, eventTrigger *realm.EventTriggerRequest) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_event_triggers" "test" {
			project_id = mongodbatlas_event_trigger.test.project_id
			app_id = mongodbatlas_event_trigger.test.app_id
		}
`, testAccMongoDBAtlasEventTriggerConfig(projectID, appID, matchKey, matchValue, awsAccID, awsRegion, eventTrigger))
}
