package todoacc

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func CheckDestroyProject(s *terraform.State) error {
	conn := TestMongoDBClient.(*config.MongoDBClient).Atlas

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_project" {
			continue
		}

		projectRes, _, _ := conn.Projects.GetOneProjectByName(context.Background(), rs.Primary.ID)
		if projectRes != nil {
			return fmt.Errorf("project (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func CheckDestroyTeam(s *terraform.State) error {
	conn := TestMongoDBClient.(*config.MongoDBClient).Atlas

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_teams" {
			continue
		}

		ids := config.DecodeStateID(rs.Primary.ID)

		orgID := ids["org_id"]
		id := ids["id"]

		// Try to find the team
		_, _, err := conn.Teams.Get(context.Background(), orgID, id)
		if err == nil {
			return fmt.Errorf("team (%s) still exists", id)
		}
	}

	return nil
}
