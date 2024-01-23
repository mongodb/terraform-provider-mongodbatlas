package acc

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func CheckDestroyLDAPConfiguration(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_ldap_configuration" {
			continue
		}
		_, _, err := Conn().LDAPConfigurations.Get(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("ldapConfiguration (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}
