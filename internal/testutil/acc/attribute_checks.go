package acc

import (
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

// copyChecks helps to prevent the accidental modification of the existing slice
func copyChecks[T map[string]string | []string](checks []resource.TestCheckFunc, additionalChecks T) []resource.TestCheckFunc {
	newChecks := make([]resource.TestCheckFunc, len(checks), len(checks)+len(additionalChecks))
	copy(newChecks, checks)
	return newChecks
}
