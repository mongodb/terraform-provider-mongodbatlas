package acc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func CheckDatabaseUserExists(resourceName string, dbUser *matlas.DatabaseUser) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := TestMongoDBClient.(*config.MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.Attributes["project_id"] == "" {
			return fmt.Errorf("no project_id is set")
		}

		if rs.Primary.Attributes["auth_database_name"] == "" {
			return fmt.Errorf("no auth_database_name is set")
		}

		if rs.Primary.Attributes["username"] == "" {
			return fmt.Errorf("no username is set")
		}

		authDB := rs.Primary.Attributes["auth_database_name"]
		projectID := rs.Primary.Attributes["project_id"]
		username := rs.Primary.Attributes["username"]

		if dbUserResp, _, err := conn.DatabaseUsers.Get(context.Background(), authDB, projectID, username); err == nil {
			*dbUser = *dbUserResp
			return nil
		}

		return fmt.Errorf("database user(%s-%s-%s) does not exist", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["username"], rs.Primary.Attributes["auth_database_name"])
	}
}

func CheckDatabaseUserAttributes(dbUser *matlas.DatabaseUser, username string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Printf("[DEBUG] difference dbUser.Username: %s , username : %s", dbUser.Username, username)
		if dbUser.Username != username {
			return fmt.Errorf("bad username: %s", dbUser.Username)
		}

		return nil
	}
}
