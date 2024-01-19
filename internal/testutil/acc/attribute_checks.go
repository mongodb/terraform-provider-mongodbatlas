package acc

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"

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
