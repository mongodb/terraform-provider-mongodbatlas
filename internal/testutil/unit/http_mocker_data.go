package unit

import (
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/hcl"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type statusText struct {
	Text               string `yaml:"text"`
	ResponseIndex      int    `yaml:"response_index"`
	Status             int    `yaml:"status"`
	DuplicateResponses int    `yaml:"duplicate_responses"`
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

func (i *RequestInfo) NormalizePath(reqURL *url.URL) string {
	queryVars := i.QueryVars()
	if len(queryVars) == 0 {
		return reqURL.Path
	}
	return removeQueryParamsAndTrim(reqURL.Path) + "?" + relevantQuery(queryVars, reqURL.Query())
}

func (i *RequestInfo) QueryVars() []string {
	selfURL, _ := url.Parse("http://localhost" + i.Path)
	query := selfURL.Query()
	queryVars := []string{}
	for key := range query {
		queryVars = append(queryVars, key)
	}
	return queryVars
}

func (i *RequestInfo) Match(t *testing.T, method, version string, reqURL *url.URL, usedVars map[string]string) bool {
	t.Helper()
	if i.Method != method || i.Version != version {
		return false
	}
	reqPath := i.NormalizePath(reqURL)
	if replaceVars(i.Path, usedVars) == reqPath {
		return true
	}
	apiPath := APISpecPath{Path: removeQueryParamsAndTrim(i.Path)}
	if !apiPath.Match(reqURL.Path) {
		return false
	}
	pathVars := apiPath.Variables(reqURL.Path)
	for name, value := range pathVars {
		oldValue, exists := usedVars[name]
		if !exists {
			t.Logf("Adding variable to %s=%s based on match of %s", name, value, i.Path)
		}
		if exists && oldValue != value {
			change, err := findVariableChange(t, name, usedVars, oldValue, value)
			if err != nil {
				t.Error(err)
				return false
			}
			usedVars[change.NewName] = change.NewValue
		} else {
			usedVars[name] = value
		}
	}
	return replaceVars(i.Path, usedVars) == reqPath
}

type stepRequests struct {
	Config           string        `yaml:"config,omitempty"`
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
			existing.Responses = append(existing.Responses, newResponse)
		}
	} else {
		s.RequestResponses = append(s.RequestResponses, *request)
	}
}

type RoundTrip struct {
	Variables   map[string]string
	QueryString string
	Request     RequestInfo
	Response    statusText
	StepNumber  int
}

func NewMockHTTPData(t *testing.T, stepCount int, tfConfigs []string) *MockHTTPData {
	steps := make([]stepRequests, stepCount)
	data := MockHTTPData{
		Steps:     steps,
		Variables: map[string]string{},
	}
	for i := range steps {
		tfConfig := tfConfigs[i]
		tfConfig = hcl.PrettyHCL(t, tfConfig)
		configVars := ExtractConfigVariables(t, tfConfig)
		err := data.UpdateVariables(t, configVars)
		if err == nil {
			continue
		}
		if _, ok := err.(*VariablesChangedError); ok {
			continue
		}
		require.NoError(t, err)
	}
	return &data
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
}

// Normalize happens after all data is captured, as a cluster.name might only be discovered as a variable in later steps
func (m *MockHTTPData) Normalize() {
	for i := range m.Steps {
		step := &m.Steps[i]
		for j := range step.RequestResponses {
			request := &step.RequestResponses[j]
			m.normalizeRequest(request)
		}
		for j := range step.DiffRequests {
			request := &step.DiffRequests[j]
			m.normalizeRequest(request)
		}
	}
}

func (m *MockHTTPData) normalizeRequest(request *RequestInfo) {
	request.Text = useVars(m.Variables, request.Text)
	for k := range request.Responses {
		response := &request.Responses[k]
		response.Text = useVars(m.Variables, response.Text)
	}
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
	normalizedPath := useVars(rtVariables, rt.Request.Path)
	if rt.QueryString != "" {
		normalizedPath += "?" + useVars(rtVariables, rt.QueryString)
	}
	if rt.StepNumber > len(m.Steps) {
		return fmt.Errorf("step number %d is out of bounds, are you re-running the same test case?", rt.StepNumber)
	}
	step := &m.Steps[rt.StepNumber-1]
	requestInfo := RequestInfo{
		Version: rt.Request.Version,
		Method:  rt.Request.Method,
		Path:    normalizedPath,
		Text:    useVars(rtVariables, rt.Request.Text),
		Responses: []statusText{
			{
				Text:          useVars(rtVariables, rt.Response.Text),
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

func useVars(vars map[string]string, text string) string {
	for key, value := range vars {
		replaceInRegex := regexp.MustCompile(fmt.Sprintf(`\W(%s)\W`, key))
		text = replaceInRegex.ReplaceAllString(text, fmt.Sprintf("{%s}", key))
		text = strings.ReplaceAll(text, value, fmt.Sprintf("{%s}", key))
	}
	return text
}
