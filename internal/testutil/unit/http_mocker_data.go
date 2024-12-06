package unit

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

type statusText struct {
	Text               string
	ResponseIndex      int
	Status             int
	DuplicateResponses int
}

func (s statusText) MarshalYAML() (interface{}, error) {
	childNodes := []*yaml.Node{
		{Kind: yaml.ScalarNode, Value: "response_index"},
		{Kind: yaml.ScalarNode, Value: fmt.Sprintf("%d", s.ResponseIndex)},

		{Kind: yaml.ScalarNode, Value: "status"},
		{Kind: yaml.ScalarNode, Value: fmt.Sprintf("%d", s.Status)},
	}
	if s.DuplicateResponses > 0 {
		childNodes = append(childNodes,
			&yaml.Node{Kind: yaml.ScalarNode, Value: "duplicate_responses"},
			&yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("%d", s.DuplicateResponses)},
		)
	}
	childNodes = append(childNodes,
		&yaml.Node{Kind: yaml.ScalarNode, Value: "text"},
		&yaml.Node{Kind: yaml.ScalarNode, Value: s.Text, Tag: "!!str", Style: yaml.DoubleQuotedStyle},
	)
	return &yaml.Node{
		Kind:    yaml.MappingNode,
		Content: childNodes,
	}, nil
}

func (s *statusText) IncreaseDuplicateResponses() {
	s.DuplicateResponses++
}

type RequestInfo struct {
	Path      string       `yaml:"path"`
	Method    string       `yaml:"method"`
	Version   string       `yaml:"version"`
	Text      string       `yaml:"text"`
	Responses []statusText `yaml:"responses"`
}

// Custom marshaling is necessary to use `flow` style only on response fields (text and responses.*.text)
func (i RequestInfo) MarshalYAML() (interface{}, error) { //nolint:gocritic // Using a pointer method leads to inconsistent dump results
	responseNode := []*yaml.Node{}
	for _, response := range i.Responses {
		node, err := response.MarshalYAML()
		if err != nil {
			return nil, err
		}
		responseNode = append(responseNode, node.(*yaml.Node))
	}
	childNodes := []*yaml.Node{
		{Kind: yaml.ScalarNode, Value: "path"},
		{Kind: yaml.ScalarNode, Value: i.Path},
		{Kind: yaml.ScalarNode, Value: "method"},
		{Kind: yaml.ScalarNode, Value: i.Method},
		{Kind: yaml.ScalarNode, Value: "version"},
		{Kind: yaml.ScalarNode, Value: i.Version, Tag: "!!str", Style: yaml.SingleQuotedStyle},
		{Kind: yaml.ScalarNode, Value: "text"},
		{Kind: yaml.ScalarNode, Value: i.Text, Tag: "!!str", Style: yaml.DoubleQuotedStyle},
		{Kind: yaml.ScalarNode, Value: "responses"},
		{Kind: yaml.SequenceNode, Content: responseNode},
	}
	return &yaml.Node{
		Kind:    yaml.MappingNode,
		Style:   yaml.FoldedStyle,
		Content: childNodes,
	}, nil
}

func (i *RequestInfo) id() string {
	return fmt.Sprintf("%s_%s_%s_%s", i.Method, i.Path, i.Version, i.Text)
}

func (i *RequestInfo) idShort() string {
	return fmt.Sprintf("%s_%s_%s", i.Method, i.Path, i.Version)
}

func (i *RequestInfo) Match(method, urlPath, version string, vars map[string]string) bool {
	if i.Method != method {
		return false
	}
	selfPath := replaceVars(i.Path, vars)
	return selfPath == urlPath && i.Version == version
}

type stepRequests struct {
	DiffRequests     []RequestInfo `yaml:"diff_requests"`
	RequestResponses []RequestInfo `yaml:"request_responses"`
}

func (s *stepRequests) findRequest(request *RequestInfo) (*RequestInfo, bool) {
	for i := range s.RequestResponses {
		if s.RequestResponses[i].id() == request.id() {
			return &s.RequestResponses[i], true
		}
	}
	return nil, false
}

func (s *stepRequests) AddRequest(request *RequestInfo, isDiff bool) {
	if isDiff {
		s.DiffRequests = append(s.DiffRequests, *request)
	}
	existing, found := s.findRequest(request)
	if found {
		lastResponse := existing.Responses[len(existing.Responses)-1]
		newResponse := request.Responses[0]
		if lastResponse.Status == newResponse.Status && lastResponse.Text == newResponse.Text {
			existing.Responses[len(existing.Responses)-1].IncreaseDuplicateResponses()
		} else {
			existing.Responses = append(existing.Responses, request.Responses[0])
		}
	} else {
		s.RequestResponses = append(s.RequestResponses, *request)
	}
}

type RoundTrip struct {
	Variables  map[string]string
	Request    RequestInfo
	Response   statusText
	StepNumber int
}

func NewMockHTTPData(stepCount int) MockHTTPData {
	steps := make([]stepRequests, stepCount)
	return MockHTTPData{
		StepCount: stepCount,
		Steps:     steps,
		Variables: map[string]string{},
	}
}

type VariableChange struct {
	OldName  string
	NewName  string
	OldValue string
	NewValue string
}

type VariablesChangedError struct {
	Changes []VariableChange
}

func (e VariablesChangedError) Error() string {
	return fmt.Sprintf("variables changed: %v", e.Changes)
}

func (e VariablesChangedError) ChangedNamesMap() map[string]string {
	result := map[string]string{}
	for _, change := range e.Changes {
		result[change.OldName] = change.NewName
	}
	return result
}

func (e VariablesChangedError) ChangedValuesMap() map[string]string {
	result := map[string]string{}
	for _, change := range e.Changes {
		result[change.OldValue] = change.NewValue
	}
	return result
}

type MockHTTPData struct {
	Variables map[string]string `yaml:"variables"`
	Steps     []stepRequests    `yaml:"steps"`
	StepCount int               `yaml:"step_count"`
}

func (m *MockHTTPData) AddRoundtrip(t *testing.T, rt *RoundTrip, isDiff bool) error {
	t.Helper()
	rtVariables := rt.Variables
	err := m.UpdateVariables(t, rtVariables)
	if vce, ok := err.(*VariablesChangedError); ok {
		for _, change := range vce.Changes {
			delete(rtVariables, change.OldName)
			rtVariables[change.NewName] = change.NewValue
		}
	} else if err != nil {
		return err
	}
	normalizedPath := useVariables(rtVariables, rt.Request.Path)
	step := &m.Steps[rt.StepNumber-1]
	requestInfo := RequestInfo{
		Version: rt.Request.Version,
		Method:  rt.Request.Method,
		Path:    normalizedPath,
		Text:    useVariables(rtVariables, rt.Request.Text),
		Responses: []statusText{
			{
				Text:          useVariables(rtVariables, rt.Response.Text),
				Status:        rt.Response.Status,
				ResponseIndex: rt.Response.ResponseIndex,
			},
		},
	}
	step.AddRequest(&requestInfo, isDiff)
	return nil
}

func (m *MockHTTPData) UpdateVariables(t *testing.T, variables map[string]string) error {
	t.Helper()
	var missingValue []string
	for name, value := range variables {
		if value == "" {
			missingValue = append(missingValue, name)
		}
	}
	if len(missingValue) > 0 {
		sort.Strings(missingValue)
		return fmt.Errorf("missing values for variables: %v", missingValue)
	}
	changes := []VariableChange{}
	for name, value := range variables {
		oldValue, exists := m.Variables[name]
		if exists && oldValue != value {
			change, err := findVariableChange(t, name, m.Variables, oldValue, value)
			if err != nil {
				return err
			}
			changes = append(changes, *change)
			m.Variables[change.NewName] = change.NewValue
		} else {
			m.Variables[name] = value
		}
	}
	if len(changes) > 0 {
		return &VariablesChangedError{Changes: changes}
	}
	return nil
}

func findVariableChange(t *testing.T, name string, vars map[string]string, oldValue, newValue string) (*VariableChange, error) {
	t.Helper()
	for suffix := 2; suffix < 10; suffix++ {
		newName := fmt.Sprintf("%s%d", name, suffix)
		oldValue2, exists := vars[newName]
		if exists && oldValue2 != newValue {
			continue
		}
		if !exists {
			t.Logf("Adding variable %s to %s=%s", name, newName, newValue)
		}
		return &VariableChange{name, newName, oldValue, newValue}, nil
	}
	return nil, fmt.Errorf("too many variables with the same name and different values: %s", name)
}

func useVariables(vars map[string]string, text string) string {
	for key, value := range vars {
		text = strings.ReplaceAll(text, value, fmt.Sprintf("{%s}", key))
	}
	return text
}
