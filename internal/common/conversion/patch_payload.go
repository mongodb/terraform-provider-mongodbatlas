package conversion

import (
	"encoding/json"
	"slices"
	"strings"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/wI2L/jsondiff"
)

type attrPatchOperations struct {
	data map[string][]jsondiff.Operation
}

func newAttrPatchOperations(patch jsondiff.Patch) *attrPatchOperations {
	self := &attrPatchOperations{
		data: map[string][]jsondiff.Operation{},
	}
	for _, op := range patch {
		rootPath := strings.Split(op.Path, "/")[1]
		self.set(rootPath, &op)
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

func (m *attrPatchOperations) hasChanged(attr string) bool {
	for _, op := range m.get(attr) {
		if slices.Contains(changeOps, op.Type) {
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
	return attrs
}

func (m *attrPatchOperations) StatePatch(attr string) jsondiff.Patch {
	patch := jsondiff.Patch{}
	var lastValue any
	for _, op := range m.get(attr) {
		if op.Type == jsondiff.OperationTest {
			lastValue = op.Value
		}
		if op.Type == jsondiff.OperationRemove {
			patch = append(patch, jsondiff.Operation{
				Type:  jsondiff.OperationAdd,
				Value: lastValue,
				Path:  op.Path,
			})
		}
	}
	return patch
}

func filterPatches(attr string, patches []jsondiff.Operation) jsondiff.Patch {
	newPatch := jsondiff.Patch{}
	for _, op := range patches {
		if strings.HasPrefix(op.Path, "/"+attr) {
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

func PatchPayloadNoChanges[T any](state, plan, reqPatch *T) (bool, error) {
	statePlanPatch, err := jsondiff.Compare(state, plan, jsondiff.Invertible())
	if err != nil {
		return false, err
	}
	attrOperations := newAttrPatchOperations(statePlanPatch)
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

	patchFromPlanDiff, err := jsondiff.Compare(reqPatch, plan)
	if err != nil {
		return false, err
	}
	for _, attr := range attrOperations.ChangedAttributes() {
		patchFromPlan := filterPatches(attr, patchFromPlanDiff)
		err = addPatchToRequest(patchFromPlan)
		if err != nil {
			return false, err
		}
		patchFromState := attrOperations.StatePatch(attr)
		err = addPatchToRequest(patchFromState)
		if err != nil {
			return false, err
		}
	}
	if string(reqJSON) == "{}" {
		return true, nil
	}
	return false, json.Unmarshal(reqJSON, reqPatch)
}
