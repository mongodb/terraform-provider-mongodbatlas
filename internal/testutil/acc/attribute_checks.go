package acc

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	matchTimestamp = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:\d{2})$`)
	matchUsername  = regexp.MustCompile(`.*@mongodb\.com$`)
)

func MatchesExpression(expr string) resource.CheckResourceAttrWithFunc {
	return func(value string) error {
		matched, err := regexp.MatchString(expr, value)
		if err != nil {
			return err
		}
		if !matched {
			return fmt.Errorf("%s did not match expression %s", value, expr)
		}
		return nil
	}
}

// IsTimestamp checks if the value is a valid timestamp in RFC3339 format.
func IsTimestamp() resource.CheckResourceAttrWithFunc {
	return func(value string) error {
		matched, err := regexp.MatchString(matchTimestamp.String(), value)
		if err != nil {
			return err
		}
		if !matched {
			return fmt.Errorf("expected a timestamp, got %s", value)
		}
		return nil
	}
}

func IsUsername() resource.CheckResourceAttrWithFunc {
	return func(value string) error {
		matched, err := regexp.MatchString(matchUsername.String(), value)
		if err != nil {
			return err
		}
		if !matched {
			return fmt.Errorf("expected a username, got %s", value)
		}
		return nil
	}
}

func CIDRBlockExpression() resource.CheckResourceAttrWithFunc {
	return func(value string) error {
		_, _, err := net.ParseCIDR(value)
		return err
	}
}

func IntGreatThan(value int) resource.CheckResourceAttrWithFunc {
	return func(input string) error {
		inputInt, err := strconv.Atoi(input)
		if err != nil {
			return err
		}
		if inputInt <= value {
			return fmt.Errorf("%d is not greater than %d", inputInt, value)
		}
		return nil
	}
}

func JSONEquals(expected string) resource.CheckResourceAttrWithFunc {
	return func(input string) error {
		var expectedAny, inputAny any

		if err := json.Unmarshal([]byte(expected), &expectedAny); err != nil {
			return fmt.Errorf("could not unmarshal json: %s", err)
		}

		if err := json.Unmarshal([]byte(input), &inputAny); err != nil {
			return fmt.Errorf("could not unmarshal json: %s", err)
		}

		if !reflect.DeepEqual(expectedAny, inputAny) {
			return fmt.Errorf("expected `%v`, got `%v`", expected, input)
		}
		return nil
	}
}

// IsProjectNameOrID accepts a project id or name and checks if the input project id matches the expected project
func IsProjectNameOrID(expected string) resource.CheckResourceAttrWithFunc {
	return func(input string) error {
		projectID := expected
		if startNumber, _ := regexp.MatchString(`^\d`, expected); !startNumber {
			resp, _, _ := ConnV2().ProjectsApi.GetGroupByName(context.Background(), expected).Execute()
			projectID = resp.GetId()
			if projectID == "" {
				return fmt.Errorf("project not found %q", expected)
			}
		}
		if projectID != input {
			return fmt.Errorf("project expected %q but got %q", projectID, input)
		}
		return nil
	}
}

// CheckRSAndDS returns a check function that asserts a set of attributes (presence and values) in resource and data sources.
// If a plural data source name is provided, it will apply checks over first result
func CheckRSAndDS(resourceName string, dataSourceName, pluralDataSourceName *string, attrsSet []string, attrsMap map[string]string, extra ...resource.TestCheckFunc) resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{}
	checks = AddAttrChecks(resourceName, checks, attrsMap)
	checks = AddAttrSetChecks(resourceName, checks, attrsSet...)
	if dataSourceName != nil {
		checks = AddAttrChecks(*dataSourceName, checks, attrsMap)
		checks = AddAttrSetChecks(*dataSourceName, checks, attrsSet...)
	}
	if pluralDataSourceName != nil {
		checks = AddAttrChecksPrefix(*pluralDataSourceName, checks, attrsMap, "results.0")
		checks = AddAttrSetChecksPrefix(*pluralDataSourceName, checks, attrsSet, "results.0")
	}
	checks = append(checks, extra...)
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func AddAttrSetChecks(targetName string, checks []resource.TestCheckFunc, attrNames ...string) []resource.TestCheckFunc {
	newChecks := copyChecks(checks, attrNames)
	for _, attrName := range attrNames {
		newChecks = append(newChecks, resource.TestCheckResourceAttrSet(targetName, attrName))
	}
	return newChecks
}

func AddNoAttrSetChecks(targetName string, checks []resource.TestCheckFunc, attrNames ...string) []resource.TestCheckFunc {
	newChecks := copyChecks(checks, attrNames)
	for _, attrName := range attrNames {
		newChecks = append(newChecks, resource.TestCheckNoResourceAttr(targetName, attrName))
	}
	return newChecks
}

func AddAttrChecks(targetName string, checks []resource.TestCheckFunc, mapChecks map[string]string) []resource.TestCheckFunc {
	newChecks := copyChecks(checks, mapChecks)
	for key, value := range mapChecks {
		newChecks = append(newChecks, resource.TestCheckResourceAttr(targetName, key, value))
	}
	return newChecks
}

func AddAttrChecksPrefix(targetName string, checks []resource.TestCheckFunc, mapChecks map[string]string, prefix string, skipNames ...string) []resource.TestCheckFunc {
	newChecks := copyChecks(checks, mapChecks)
	prefix, _ = strings.CutSuffix(prefix, ".")
	for key, value := range mapChecks {
		if slices.Contains(skipNames, key) {
			continue
		}
		keyWithPrefix := fmt.Sprintf("%s.%s", prefix, key)
		newChecks = append(newChecks, resource.TestCheckResourceAttr(targetName, keyWithPrefix, value))
	}
	return newChecks
}

func AddAttrSetChecksPrefix(targetName string, checks []resource.TestCheckFunc, attrNames []string, prefix string) []resource.TestCheckFunc {
	newChecks := copyChecks(checks, attrNames)
	prefix, _ = strings.CutSuffix(prefix, ".")
	for _, key := range attrNames {
		keyWithPrefix := fmt.Sprintf("%s.%s", prefix, key)
		newChecks = append(newChecks, resource.TestCheckResourceAttrSet(targetName, keyWithPrefix))
	}
	return newChecks
}

// copyChecks helps to prevent the accidental modification of the existing slice
func copyChecks[T map[string]string | []string](checks []resource.TestCheckFunc, additionalChecks T) []resource.TestCheckFunc {
	newChecks := make([]resource.TestCheckFunc, len(checks), len(checks)+len(additionalChecks))
	copy(newChecks, checks)
	return newChecks
}
