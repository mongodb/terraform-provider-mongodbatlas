package acc

import (
	"context"
	"errors"
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// PluralResultCheck creates a StateCheck that finds results.x by using x.resultAttributeID and then performs the checks.
// This is useful when tests are run in parallel and the index of the result is not known.
// Code is inspired by the example here: https://developer.hashicorp.com/terraform/plugin/testing/acceptance-tests/state-checks/custom
func PluralResultCheck(resourceAddress, resultAttributeID string, resultAttributeMatch knownvalue.Check, checks map[string]knownvalue.Check) statecheck.StateCheck {
	return expectPluralResultChecks{
		resourceAddress:      resourceAddress,
		resultAttributeID:    resultAttributeID,
		resultAttributeMatch: resultAttributeMatch,
		checks:               checks,
	}
}

type expectPluralResultChecks struct {
	resultAttributeMatch knownvalue.Check
	checks               map[string]knownvalue.Check
	resourceAddress      string
	resultAttributeID    string
}

func (e expectPluralResultChecks) CheckState(ctx context.Context, req statecheck.CheckStateRequest, resp *statecheck.CheckStateResponse) {
	var resource *tfjson.StateResource

	if req.State == nil {
		resp.Error = fmt.Errorf("state is nil")
	}

	if req.State.Values == nil {
		resp.Error = fmt.Errorf("state does not contain any state values")
	}

	if req.State.Values.RootModule == nil {
		resp.Error = fmt.Errorf("state does not contain a root module")
	}

	for _, r := range req.State.Values.RootModule.Resources {
		if e.resourceAddress == r.Address {
			resource = r

			break
		}
	}

	if resource == nil {
		resp.Error = fmt.Errorf("%s - Resource not found in state", e.resourceAddress)

		return
	}
	resultsPath := tfjsonpath.New("results")
	result, err := tfjsonpath.Traverse(resource.AttributeValues, resultsPath)

	if err != nil {
		resp.Error = err
		return
	}

	foundValue, err := findResultsMatch(result, e.resultAttributeID, e.resultAttributeMatch)
	if err != nil {
		resp.Error = err
		return
	}
	if err = doResultChecks(foundValue, e.checks); err != nil {
		resp.Error = err
		return
	}
}

func findResultsMatch(result any, resultAttributeID string, resultAttributeMatch knownvalue.Check) (any, error) {
	resultList, ok := result.([]any)
	if !ok {
		return nil, fmt.Errorf("results is not a list of map")
	}
	if len(resultList) == 0 {
		return nil, fmt.Errorf("results is empty")
	}
	for _, resultCandidate := range resultList {
		attrValue, err := tfjsonpath.Traverse(resultCandidate, PathFromString(resultAttributeID))
		if err != nil {
			return nil, fmt.Errorf("result does not contain %s, err=%w", resultAttributeID, err)
		}
		if err := resultAttributeMatch.CheckValue(attrValue); err == nil {
			return resultCandidate, nil
		}
	}
	return nil, fmt.Errorf("none of the results.*.%s matched", resultAttributeID)
}

func doResultChecks(foundValue any, checks map[string]knownvalue.Check) error {
	checkErrors := make([]error, 0)
	for key, check := range checks {
		if check == nil {
			continue
		}
		jsonPath := PathFromString(key)
		value, err := tfjsonpath.Traverse(foundValue, jsonPath)
		if err != nil {
			checkErrors = append(checkErrors, fmt.Errorf("result does not contain %s", key), err)
			continue
		}
		if err := check.CheckValue(value); err != nil {
			checkErrors = append(checkErrors, fmt.Errorf("result %s failed check: %s", key, err))
		}
	}
	if len(checkErrors) > 0 {
		return errors.Join(checkErrors...)
	}
	return nil
}
