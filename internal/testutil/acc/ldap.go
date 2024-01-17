package acc

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func CheckDestroyLDAPConfiguration(s *terraform.State) error {
	conn := TestAccProviderSdkV2.Meta().(*config.MongoDBClient).Atlas

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_ldap_configuration" {
			continue
		}

		_, _, err := conn.LDAPConfigurations.Get(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("ldapConfiguration (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}
