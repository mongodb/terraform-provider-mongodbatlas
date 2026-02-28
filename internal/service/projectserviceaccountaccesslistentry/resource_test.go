package projectserviceaccountaccesslistentry_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/hcl"
	"go.mongodb.org/atlas-sdk/v20250312014/admin"
)

const (
	resourceName         = "mongodbatlas_project_service_account_access_list_entry.test"
	dataSourceName       = "data.mongodbatlas_project_service_account_access_list_entry.test"
	dataSourcePluralName = "data.mongodbatlas_project_service_account_access_list_entries.test"
)

type testEntry struct {
	cidr string
	ip   string
}

func TestAccProjectServiceAccountAccessListEntry_singleEntry(t *testing.T) {
	var (
		projectID     = acc.ProjectIDExecution(t)
		name          = acc.RandomName()
		cidrEntries   = []testEntry{{cidr: "192.168.1.0/24"}}
		ipEntries     = []testEntry{{ip: "192.168.1.1"}}
		resourceName0 = resourceName + "_0"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{ // Create with cidr
				Config: configBasic(projectID, name, cidrEntries),
				Check:  checkBasic(cidrEntries),
			},
			{
				ResourceName:                         resourceName0,
				ImportStateIdFunc:                    importStateIDFunc(resourceName0),
				ImportStateVerifyIdentifierAttribute: "cidr_block",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
			{ // Replace with ip
				Config: configBasic(projectID, name, ipEntries),
				Check:  checkBasic(ipEntries),
			},
			{
				ResourceName:                         resourceName0,
				ImportStateIdFunc:                    importStateIDFunc(resourceName0),
				ImportStateVerifyIdentifierAttribute: "ip_address",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
		},
	})
}

func TestAccProjectServiceAccountAccessListEntry_multipleEntries(t *testing.T) {
	var (
		projectID     = acc.ProjectIDExecution(t)
		name          = acc.RandomName()
		entries       = []testEntry{{cidr: "100.200.30.4/32"}, {ip: "4.3.2.1"}, {cidr: "123.234.0.0/16"}}
		resourceName0 = resourceName + "_0"
		resourceName1 = resourceName + "_1"
		resourceName2 = resourceName + "_2"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, name, entries),
				Check:  checkBasic(entries),
			},
			{
				ResourceName:                         resourceName0,
				ImportStateIdFunc:                    importStateIDFunc(resourceName0),
				ImportStateVerifyIdentifierAttribute: "cidr_block",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
			{
				ResourceName:                         resourceName1,
				ImportStateIdFunc:                    importStateIDFunc(resourceName1),
				ImportStateVerifyIdentifierAttribute: "ip_address",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
			{
				ResourceName:                         resourceName2,
				ImportStateIdFunc:                    importStateIDFunc(resourceName2),
				ImportStateVerifyIdentifierAttribute: "cidr_block",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
		},
	})
}

func TestAccProjectServiceAccountAccessListEntry_errors(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config:      configError(testEntry{}),
				ExpectError: regexp.MustCompile("cidr_block or ip_address must be provided"),
			},
			{
				Config:      configError(testEntry{cidr: "invalid cidr"}),
				ExpectError: regexp.MustCompile("Attribute cidr_block string value must be defined as a valid cidr"),
			},
			{
				Config:      configError(testEntry{ip: "invalid ip"}),
				ExpectError: regexp.MustCompile("Attribute ip_address string value must be defined as a valid IP Address"),
			},
			{
				Config:      configError(testEntry{cidr: "192.168.1.0/24", ip: "192.168.1.1"}),
				ExpectError: regexp.MustCompile(`Attribute "ip_address" cannot be specified when "cidr_block" is specified`),
			},
		},
	})
}

func configBasic(projectID, name string, entries []testEntry) string {
	var entriesStr strings.Builder
	resourceNames := []string{}
	for i, entry := range entries {
		fmt.Fprintf(&entriesStr, `
			resource "mongodbatlas_project_service_account_access_list_entry" "test_%[1]d" {
				project_id = %[2]q
				client_id  = mongodbatlas_project_service_account.test.client_id
				%[3]s
			}

			data "mongodbatlas_project_service_account_access_list_entry" "test_%[1]d" {
				project_id = %[2]q
				client_id  = mongodbatlas_project_service_account.test.client_id
				%[3]s
				depends_on = [mongodbatlas_project_service_account_access_list_entry.test_%[1]d]
			}
		`, i, projectID, entry.hclStr())
		resourceNames = append(resourceNames, fmt.Sprintf("%s_%d", resourceName, i))
	}
	resourceNamesStr := hcl.StringSliceToHCL(resourceNames)

	return fmt.Sprintf(`
		resource "mongodbatlas_project_service_account" "test" {
			project_id                 = %[1]q
			name                       = %[2]q
			description                = "Acceptance Test SA for Project SA access list"
			roles                      = ["GROUP_READ_ONLY"]
			secret_expires_after_hours = 12
		}

		%[3]s

		data "mongodbatlas_project_service_account_access_list_entries" "test" {
			project_id = %[1]q
			client_id  = mongodbatlas_project_service_account.test.client_id
			depends_on = %[4]s
		}
	`, projectID, name, entriesStr.String(), resourceNamesStr)
}

func configError(entry testEntry) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project_service_account_access_list_entry" "test" {
			project_id = "000000000000000000000000"
			client_id  = "mdb_sa_id_000000000000000000000000"
			%s
		}
	`, entry.hclStr())
}

func checkBasic(entries []testEntry) resource.TestCheckFunc {
	// Check plural DS first result only when there is 1 entry
	var pluralDSName *string
	if len(entries) == 1 {
		pluralDSName = new(dataSourcePluralName)
	}

	attrsSet := []string{"client_id", "created_at", "request_count"}
	checks := []resource.TestCheckFunc{}
	for i, entry := range entries {
		resourceName := fmt.Sprintf("%s_%d", resourceName, i)
		dataSourceName := fmt.Sprintf("%s_%d", dataSourceName, i)
		checks = append(checks, acc.CheckRSAndDS(
			resourceName, new(dataSourceName), pluralDSName,
			attrsSet, entry.attrMap(),
			checkExists(resourceName),
		))
	}

	checks = append(checks, resource.TestCheckResourceAttr(dataSourcePluralName, "results.#", strconv.Itoa(len(entries))))

	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		projectID := rs.Primary.Attributes["project_id"]
		clientID := rs.Primary.Attributes["client_id"]
		cidrOrIP := getCidrOrIP(rs)

		if projectID == "" || clientID == "" || cidrOrIP == "" {
			return fmt.Errorf("checkExists, attributes not found for: %s", resourceName)
		}

		entry, err := getEntry(projectID, clientID, cidrOrIP)
		if entry == nil || err != nil {
			return fmt.Errorf("access list entry (%s/%s/%s) does not exist", projectID, clientID, cidrOrIP)
		}
		return nil
	}
}

func checkDestroy(s *terraform.State) error {
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	for name, rs := range s.RootModule().Resources {
		if name != resourceName {
			continue
		}

		projectID := rs.Primary.Attributes["project_id"]
		clientID := rs.Primary.Attributes["client_id"]
		cidrOrIP := getCidrOrIP(rs)

		if projectID == "" || clientID == "" || cidrOrIP == "" {
			return fmt.Errorf("checkExists, attributes not found for: %s", resourceName)
		}

		entry, _ := getEntry(projectID, clientID, cidrOrIP)
		if entry != nil {
			return fmt.Errorf("access list entry (%s/%s/%s) still exists", projectID, clientID, cidrOrIP)
		}

		// Delete the service account (project_service_account DELETE only removes the project assignment)
		_, _ = acc.ConnV2().ServiceAccountsApi.DeleteOrgServiceAccount(context.Background(), clientID, orgID).Execute()
		return nil
	}
	return nil
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		projectID := rs.Primary.Attributes["project_id"]
		clientID := rs.Primary.Attributes["client_id"]
		cidrOrIP := getCidrOrIP(rs)
		if projectID == "" || clientID == "" || cidrOrIP == "" {
			return "", fmt.Errorf("import, attributes not found for: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s/%s", projectID, clientID, cidrOrIP), nil
	}
}

func getCidrOrIP(rs *terraform.ResourceState) string {
	cidrOrIP := rs.Primary.Attributes["ip_address"]
	if cidrOrIP == "" {
		cidrOrIP = rs.Primary.Attributes["cidr_block"]
	}
	return cidrOrIP
}

func getEntry(projectID, clientID, cidrOrIP string) (*admin.ServiceAccountIPAccessListEntry, error) {
	res, _, err := acc.ConnV2().ServiceAccountsApi.ListAccessList(context.Background(), projectID, clientID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get access list: %w", err)
	}
	entries := res.GetResults()
	for i := range entries {
		entry := &entries[i]
		if entry.GetIpAddress() == cidrOrIP || entry.GetCidrBlock() == cidrOrIP {
			return entry, nil
		}
	}
	return nil, nil
}

func (e testEntry) hclStr() string {
	result := ""
	if e.cidr != "" {
		result += fmt.Sprintf("cidr_block = %q\n", e.cidr)
	}
	if e.ip != "" {
		result += fmt.Sprintf("ip_address = %q\n", e.ip)
	}
	return result
}

func (e testEntry) attrMap() map[string]string {
	result := map[string]string{}
	if e.cidr != "" {
		result["cidr_block"] = e.cidr
	}
	if e.ip != "" {
		result["ip_address"] = e.ip
	}
	return result
}
