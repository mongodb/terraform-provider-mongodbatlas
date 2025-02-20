package update

import (
	"encoding/json"
	"reflect"
	"regexp"
	"slices"
	"strings"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/wI2L/jsondiff"
)

type attrPatchOperations struct {
	data                 map[string][]jsondiff.Operation
	ignoreInStateSuffix  []string
	ignoreInStatePrefix  []string
	includeInStateSuffix []string
	forceUpdateAttr      []string
}

func (m *attrPatchOperations) ignoreInStatePath(path string) bool {
	for _, include := range m.includeInStateSuffix {
		suffix := "/" + include
		if strings.HasSuffix(path, suffix) {
			return false
		}
	}
	for _, ignore := range m.ignoreInStateSuffix {
		suffix := "/" + ignore
		if strings.HasSuffix(path, suffix) {
			return true
		}
	}
	for _, ignore := range m.ignoreInStatePrefix {
		for _, part := range strings.Split(path, "/") {
			if ignore == part {
				return true
			}
		}
	}
	return false
}

func newAttrPatchOperations(patch jsondiff.Patch, options []PatchOptions) *attrPatchOperations {
	var (
		ignoreSuffixInState  []string
		ignorePrefixInState  []string
		includeSuffixInState []string
		forceUpdateAttr      []string
	)
	for _, option := range options {
		ignoreSuffixInState = append(ignoreSuffixInState, option.IgnoreInStateSuffix...)
		ignorePrefixInState = append(ignorePrefixInState, option.IgnoreInStatePrefix...)
		includeSuffixInState = append(includeSuffixInState, option.IncludeInStateSuffix...)
		forceUpdateAttr = append(forceUpdateAttr, option.ForceUpdateAttr...)
	}
	self := &attrPatchOperations{
		data:                 map[string][]jsondiff.Operation{},
		ignoreInStateSuffix:  ignoreSuffixInState,
		ignoreInStatePrefix:  ignorePrefixInState,
		includeInStateSuffix: includeSuffixInState,
		forceUpdateAttr:      forceUpdateAttr,
	}
	for _, op := range patch {
		if op.Path == "" {
			self.set("", &op)
		} else {
			rootPath := strings.Split(op.Path, "/")[1]
			self.set(rootPath, &op)
		}
	}
	return self
}

func (m *attrPatchOperations) set(attr string, value *jsondiff.Operation) {
	m.data[attr] = append(m.data[attr], *value)
}

func (m *attrPatchOperations) get(attr string) []jsondiff.Operation {
	return m.data[attr]
}

var changeOps = []string{jsondiff.OperationReplace, jsondiff.OperationAdd}
var digitRegex = regexp.MustCompile(`^\d+$`)

func indexRemoval(path string) bool {
	if strings.Count(path, "/") == 0 {
		return false
	}
	pathParts := strings.Split(path, "/")
	lastPart := pathParts[len(pathParts)-1]
	return digitRegex.MatchString(lastPart)
}

func (m *attrPatchOperations) hasChanged(attr string) bool {
	if slices.Contains(m.forceUpdateAttr, attr) {
		return true
	}
	for _, op := range m.get(attr) {
		if slices.Contains(changeOps, op.Type) {
			return true
		}
		if op.Type == jsondiff.OperationRemove && indexRemoval(op.Path) {
			return true
		}
	}
	return false
}

func (m *attrPatchOperations) ChangedAttributes() []string {
	attrs := []string{}
	for attr := range m.data {
		if m.hasChanged(attr) {
			attrs = append(attrs, attr)
		}
	}
	// There might be a case where there are no changes in m.data for the attributes in forceUpdateAttr
	for _, attr := range m.forceUpdateAttr {
		if !slices.Contains(attrs, attr) {
			attrs = append(attrs, attr)
		}
	}
	return attrs
}

func (m *attrPatchOperations) StatePatch(attr string) jsondiff.Patch {
	patch := jsondiff.Patch{}
	var lastValue any
	for _, op := range m.get(attr) {
		if op.Type == jsondiff.OperationTest {
			lastValue = op.Value
		}
		if op.Type == jsondiff.OperationRemove && !indexRemoval(op.Path) {
			path := op.Path
			if m.ignoreInStatePath(path) {
				continue
			}
			patch = append(patch, jsondiff.Operation{
				Type:  jsondiff.OperationAdd,
				Value: lastValue,
				Path:  path,
			})
		}
	}
	return patch
}

func filterPatches(attr string, patches []jsondiff.Operation) jsondiff.Patch {
	newPatch := jsondiff.Patch{}
	for _, op := range patches {
		if attr == "" && op.Path == "" {
			newPatch = append(newPatch, op)
		} else if strings.HasPrefix(op.Path, "/"+attr) {
			newPatch = append(newPatch, op)
		}
	}
	return newPatch
}

func convertJSONDiffToJSONPatch(patch jsondiff.Patch) (jsonpatch.Patch, error) {
	patchKeyBytes, err := json.Marshal(patch)
	if err != nil {
		return nil, err
	}
	decodedPatch, err := jsonpatch.DecodePatch(patchKeyBytes)
	if err != nil {
		return nil, err
	}
	return decodedPatch, nil
}

// Current limitation if the field is set as part of a nested attribute in a map
type PatchOptions struct {
	IgnoreInStateSuffix  []string
	IgnoreInStatePrefix  []string
	IncludeInStateSuffix []string
	ForceUpdateAttr      []string
}

// PatchPayload uses the state and plan to changes to find the patch request, including changes only when:
// - The plan has replaced, added, or removed list values from the state
// Note that we intentionally do NOT include removed state values:
// - The state value is probably computed and not needed in the request
// However, for nested attributes, we MUST include some of the removed state values (e.g., `replication_spec[*].(id|zone_id)`)
// Assumptions:
// - Only Optional|Required attributes are set in the state|plan. `connection_strings` are not needed
// - --> Except specific computed attributes in nested_attributes (e.g., `replication_spec[*].(id|zone_id`)
// - The state and plan can be dumped to json
// How it works:
// 1. Use `jsondiff` to find the patch, aka. operations to go from state to plan
// 2. Groups the operations by attribute name
// 3. Filters the operations to only include replaced, added or removed list values
// 4. Adds nested "removed" values from the state to the request
// 5. Use `jsonpatch` to apply each attribute plan & state patch to an empty JSON object
// 6. Create a `patchReq` pointer with the final JSON object marshaled to `T` or return nil if there are no changes (`{}`)
func PatchPayload[T any, U any](state *T, plan *U, options ...PatchOptions) (*U, error) {
	if plan == nil {
		return nil, nil
	}
	statePlanPatch, err := jsondiff.Compare(state, plan, jsondiff.Invertible())
	if err != nil {
		return nil, err
	}
	attrOperations := newAttrPatchOperations(statePlanPatch, options)
	reqJSON := []byte(`{}`)

	addPatchToRequest := func(patchDiff jsondiff.Patch) error {
		if len(patchDiff) == 0 {
			return nil
		}
		patch, err := convertJSONDiffToJSONPatch(patchDiff)
		if err != nil {
			return err
		}
		reqJSON, err = patch.Apply(reqJSON)
		if err != nil {
			return err
		}
		return nil
	}

	patchReq := new(U)
	patchFromPlanDiff, err := jsondiff.Compare(patchReq, plan)
	if err != nil {
		return nil, err
	}
	for _, attr := range attrOperations.ChangedAttributes() {
		patchFromPlan := filterPatches(attr, patchFromPlanDiff)
		err = addPatchToRequest(patchFromPlan)
		if err != nil {
			return nil, err
		}
		patchFromState := attrOperations.StatePatch(attr)
		err = addPatchToRequest(patchFromState)
		if err != nil {
			return nil, err
		}
	}
	if string(reqJSON) == "{}" {
		return nil, nil
	}
	return patchReq, json.Unmarshal(reqJSON, patchReq)
}

func IsZeroValues[T any](last *T) bool {
	if last == nil {
		return true
	}
	empty := new(T)
	return reflect.DeepEqual(last, empty)
}
