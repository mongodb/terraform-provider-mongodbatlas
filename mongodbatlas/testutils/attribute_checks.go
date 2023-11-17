package testutils

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

func JSONEquals[T any](value T) resource.CheckResourceAttrWithFunc {
	return func(input string) error {
		var actual T
		if err := json.Unmarshal([]byte(input), &actual); err != nil {
			return fmt.Errorf("could not unmarshal json: %s", err)
		}

		if !reflect.DeepEqual(actual, value) {
			return fmt.Errorf("expected `%v`, got `%v`", value, actual)
		}

		return nil
	}
}
