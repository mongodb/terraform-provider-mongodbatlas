package databaseuserapi_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/databaseuser"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const resourceName = "mongodbatlas_database_user_api.test"

func TestAccDatabaseUserAPI_basic(t *testing.T) {
	var (
		groupID  = acc.ProjectIDExecution(t)
		username = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(groupID, username, "atlasAdmin", "First Key", "First value"),
				Check:  checkBasic(groupID, username, "atlasAdmin", "First Key", "First value"),
			},
			{
				Config: configBasic(groupID, username, "read", "Second Key", "Second value"),
				Check:  checkBasic(groupID, username, "read", "Second Key", "Second value"),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       importStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func configBasic(groupID, username, roleName, keyLabel, valueLabel string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_database_user_api" "test" {
			group_id         = %[1]q
			username           = %[2]q
			password           = "test-acc-password"
			database_name = "admin"

			roles = [{
				role_name     = %[3]q
				database_name = "admin"
			}]

			labels = [{
				key   = %[4]q
				value = %[5]q
			}]
		}
	`, groupID, username, roleName, keyLabel, valueLabel)
}

func checkBasic(groupID, username, roleName, keyLabel, valueLabel string) resource.TestCheckFunc {
	mapChecks := map[string]string{
		"group_id":          groupID,
		"username":          username,
		"password":          "test-acc-password",
		"database_name":     "admin",
		"labels.#":          "1",
		"labels.0.key":      keyLabel,
		"labels.0.value":    valueLabel,
		"roles.#":           "1",
		"roles.0.role_name": roleName,
	}
	checks := acc.AddAttrChecks(resourceName, nil, mapChecks)
	checks = append(checks, checkExists(resourceName))
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		groupID := rs.Primary.Attributes["group_id"]
		databaseName := rs.Primary.Attributes["database_name"]
		username := rs.Primary.Attributes["username"]
		if groupID != "" || databaseName != "" || username != "" {
			return fmt.Errorf("attributes not found for: %s", resourceName)
		}
		if _, _, err := acc.ConnV2().DatabaseUsersApi.GetDatabaseUser(context.Background(), groupID, databaseName, username).Execute(); err == nil {
			return nil
		}
		return fmt.Errorf("database user(%s-%s-%s) does not exist", groupID, databaseName, username)
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_database_user_api" {
			continue
		}
		groupID, username, databaseName, err := databaseuser.SplitDatabaseUserImportID(rs.Primary.ID)
		if err != nil {
			continue
		}
		_, _, err = acc.ConnV2().DatabaseUsersApi.GetDatabaseUser(context.Background(), groupID, databaseName, username).Execute()
		if err == nil {
			return fmt.Errorf("database user (%s) still exists", groupID)
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
		ids := conversion.DecodeStateID(rs.Primary.ID)
		return fmt.Sprintf("%s-%s-%s", ids["group_id"], ids["username"], ids["database_name"]), nil
	}
}
