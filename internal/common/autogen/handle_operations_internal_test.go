package autogen

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen/customtypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testClient(t *testing.T, handler http.Handler) *config.MongoDBClient {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	client, err := config.NewClient(&config.Credentials{BaseURL: srv.URL}, "test")
	require.NoError(t, err)
	return client
}

func jsonResponse(status int, body string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		fmt.Fprint(w, body)
	}
}

func atlasError(status int, detail string) http.HandlerFunc {
	return jsonResponse(status, fmt.Sprintf(`{"detail":%q,"error":%d,"errorCode":"TEST_ERROR"}`, detail, status))
}

func TestCallDelete(t *testing.T) {
	testCases := map[string]struct {
		handler             http.HandlerFunc
		expectedErrContains string
	}{
		"400 returns error":                {atlasError(http.StatusBadRequest, "cannot delete resource"), "cannot delete resource"},
		"401 returns error":                {atlasError(http.StatusUnauthorized, "unauthorized"), "unauthorized"},
		"409 returns error":                {atlasError(http.StatusConflict, "conflict"), "conflict"},
		"500 returns error":                {atlasError(http.StatusInternalServerError, "server error"), "server error"},
		"404 is tolerated already deleted": {atlasError(http.StatusNotFound, "not found"), ""},
		"204 succeeds":                     {jsonResponse(http.StatusNoContent, ""), ""},
		"200 with body succeeds":           {jsonResponse(http.StatusOK, `{"stateName":"DELETING"}`), ""},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			req := &HandleDeleteReq{
				Client:     testClient(t, tc.handler),
				CallParams: &config.APICallParams{RelativePath: "/api/test", Method: http.MethodDelete},
			}
			err := callDelete(context.Background(), req)
			if tc.expectedErrContains == "" {
				require.NoError(t, err)
				return
			}
			require.ErrorContains(t, err, tc.expectedErrContains)
		})
	}
}

func TestRefreshFunc(t *testing.T) {
	testCases := map[string]struct {
		handler             http.HandlerFunc
		hooks               any
		expectedState       string
		expectedErrContains string
	}{
		"500 returns error":                   {atlasError(http.StatusInternalServerError, "server error"), nil, "", "server error"},
		"404 is DELETED":                      {atlasError(http.StatusNotFound, "not found"), nil, retrystrategy.RetryStrategyDeletedState, ""},
		"200 returns state value":             {jsonResponse(http.StatusOK, `{"stateName":"IDLE"}`), nil, "IDLE", ""},
		"200 empty JSON with hook is DELETED": {jsonResponse(http.StatusOK, `{}`), emptyJSONAsNotFoundHook, retrystrategy.RetryStrategyDeletedState, ""},
		"200 empty JSON without hook errors":  {jsonResponse(http.StatusOK, `{}`), nil, "", "wait state attribute not found"},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			client := testClient(t, tc.handler)
			wait := &WaitReq{
				StateProperty: "stateName",
				PendingStates: []string{"IDLE"},
				TargetStates:  []string{retrystrategy.RetryStrategyDeletedState},
				CallParams: func(model any) *config.APICallParams {
					return &config.APICallParams{RelativePath: "/api/test", Method: http.MethodGet}
				},
			}
			_, state, err := refreshFunc(context.Background(), wait, client, nil, tc.hooks)()
			if tc.expectedErrContains != "" {
				require.ErrorContains(t, err, tc.expectedErrContains)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedState, state)
		})
	}
}

type testPostReadHook struct {
	postRead func(HandleReadReq, APICallResult) APICallResult
}

func (h *testPostReadHook) PostReadAPICall(req HandleReadReq, result APICallResult) APICallResult {
	return h.postRead(req, result)
}

// emptyJSONAsNotFoundHook replicates hooks of resources whose API returns an ok status code with an
// empty JSON body for missing resources, e.g. search deployment.
var emptyJSONAsNotFoundHook = &testPostReadHook{
	postRead: func(_ HandleReadReq, result APICallResult) APICallResult {
		if result.Err == nil && IsEmptyJSON(result.Body) {
			result.Err = ErrNotFound
		}
		return result
	},
}

// notFoundSignalHook replicates hooks that signal a missing resource via ErrNotFound, e.g. service account secrets.
var notFoundSignalHook = &testPostReadHook{
	postRead: func(_ HandleReadReq, _ APICallResult) APICallResult {
		return APICallResult{Err: fmt.Errorf("secret not found in response: %w", ErrNotFound)}
	},
}

func assertDiagError(t *testing.T, diags *diag.Diagnostics, summary, detailContains string) {
	t.Helper()
	require.True(t, diags.HasError(), "expected a diagnostics error")
	if summary != "" {
		assert.Equal(t, summary, diags.Errors()[0].Summary())
	}
	if detailContains != "" {
		assert.Contains(t, diags.Errors()[0].Detail(), detailContains)
	}
}

func assertDiagsOK(t *testing.T, diags *diag.Diagnostics) {
	t.Helper()
	require.False(t, diags.HasError(), "unexpected diagnostics: %#v", diags.Errors())
}

func TestNotFound(t *testing.T) {
	respWithStatus := func(status int) *http.Response {
		return &http.Response{StatusCode: status}
	}
	testCases := map[string]struct {
		callResult APICallResult
		expected   bool
	}{
		"404 with error":                        {APICallResult{Resp: respWithStatus(http.StatusNotFound), Err: fmt.Errorf("404 not found")}, true},
		"hook sentinel with no response":        {APICallResult{Err: fmt.Errorf("secret not found: %w", ErrNotFound)}, true},
		"error with 2xx response":               {APICallResult{Resp: respWithStatus(http.StatusOK), Err: fmt.Errorf("decode failure")}, false},
		"transport error with no response":      {APICallResult{Err: fmt.Errorf("connection refused")}, false},
		"success with empty body":               {APICallResult{Resp: respWithStatus(http.StatusNoContent)}, false},
		"success with empty JSON body":          {APICallResult{Resp: respWithStatus(http.StatusOK), Body: []byte("{}")}, false},
		"success with body":                     {APICallResult{Resp: respWithStatus(http.StatusOK), Body: []byte(`{"a":1}`)}, false},
		"zero value result is not classifiable": {APICallResult{}, false},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, notFound(tc.callResult))
		})
	}
}

type testReadModel struct {
	Name types.String `tfsdk:"name" apiname:"name"`
}

func testReadRequest(t *testing.T, handler http.Handler) (HandleReadReq, *tfsdk.State, *diag.Diagnostics) {
	t.Helper()
	stateValue := tftypes.NewValue(
		tftypes.Object{AttributeTypes: map[string]tftypes.Type{"name": tftypes.String}},
		map[string]tftypes.Value{"name": tftypes.NewValue(tftypes.String, "original")},
	)
	state := &tfsdk.State{
		Raw:    stateValue,
		Schema: schema.Schema{Attributes: map[string]schema.Attribute{"name": schema.StringAttribute{}}},
	}
	diags := &diag.Diagnostics{}
	return HandleReadReq{
		State:      &testReadModel{Name: types.StringValue("original")},
		RespState:  state,
		Client:     testClient(t, handler),
		CallParams: &config.APICallParams{RelativePath: "/api/test", Method: http.MethodGet},
		RespDiags:  diags,
	}, state, diags
}

func TestHandleRead(t *testing.T) {
	testCases := map[string]struct {
		handler       http.HandlerFunc
		hooks         any
		wantErrDetail string
		wantName      string
		wantRemoved   bool
	}{
		"500 reports error and keeps state":                    {atlasError(http.StatusInternalServerError, "server error"), nil, "server error", "", false},
		"404 removes resource from state":                      {atlasError(http.StatusNotFound, "not found"), nil, "", "", true},
		"hook signaling not found removes resource from state": {jsonResponse(http.StatusOK, `{"secrets":[]}`), notFoundSignalHook, "", "", true},
		"200 sets state":                                       {jsonResponse(http.StatusOK, `{"name":"updated"}`), nil, "", "updated", false},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			req, state, diags := testReadRequest(t, tc.handler)
			req.Hooks = tc.hooks
			HandleRead(context.Background(), req)
			if tc.wantErrDetail != "" {
				assertDiagError(t, diags, "", tc.wantErrDetail)
				assert.False(t, state.Raw.IsNull(), "resource must not be removed from state on API error")
				return
			}
			assertDiagsOK(t, diags)
			assert.Equal(t, tc.wantRemoved, state.Raw.IsNull())
			if tc.wantName != "" {
				var model testReadModel
				require.False(t, state.Get(context.Background(), &model).HasError())
				assert.Equal(t, tc.wantName, model.Name.ValueString())
			}
		})
	}
}

type testListResultModel struct {
	Name types.String `tfsdk:"name"`
}

type testListModel struct {
	Results customtypes.NestedListValue[testListResultModel] `tfsdk:"results" autogen:"omitjson,emptyjsonaslist"`
}

func testListReadRequest(t *testing.T, handler http.Handler) (HandleReadReq, *tfsdk.State, *diag.Diagnostics) {
	t.Helper()
	ctx := context.Background()
	testSchema := dsschema.Schema{Attributes: map[string]dsschema.Attribute{
		"results": dsschema.ListNestedAttribute{
			Computed:   true,
			CustomType: customtypes.NewNestedListType[testListResultModel](ctx),
			NestedObject: dsschema.NestedAttributeObject{
				Attributes: map[string]dsschema.Attribute{
					"name": dsschema.StringAttribute{Computed: true},
				},
			},
		},
	}}
	state := &tfsdk.State{
		Raw:    tftypes.NewValue(testSchema.Type().TerraformType(ctx), nil),
		Schema: testSchema,
	}
	diags := &diag.Diagnostics{}
	return HandleReadReq{
		State:      &testListModel{},
		RespState:  state,
		Client:     testClient(t, handler),
		CallParams: &config.APICallParams{RelativePath: "/api/test", Method: http.MethodGet},
		RespDiags:  diags,
	}, state, diags
}

func TestHandleDataSourceReadList(t *testing.T) {
	t.Run("500 reports the API error", func(t *testing.T) {
		req, _, diags := testReadRequest(t, atlasError(http.StatusInternalServerError, "server error")) // list schema not needed to assert the error
		HandleDataSourceReadList(context.Background(), req)
		assertDiagError(t, diags, "", "server error")
	})

	t.Run("200 empty JSON sets empty results", func(t *testing.T) {
		req, state, diags := testListReadRequest(t, jsonResponse(http.StatusOK, `{}`))
		HandleDataSourceReadList(context.Background(), req)
		assertDiagsOK(t, diags)
		var model testListModel
		require.False(t, state.Get(context.Background(), &model).HasError())
		assert.Empty(t, model.Results.Elements())
	})

	t.Run("200 empty page sets empty results list", func(t *testing.T) {
		req, state, diags := testListReadRequest(t, jsonResponse(http.StatusOK, `{"results":[],"totalCount":0}`))
		HandleDataSourceReadList(context.Background(), req)
		assertDiagsOK(t, diags)
		var model testListModel
		require.False(t, state.Get(context.Background(), &model).HasError())
		require.False(t, model.Results.IsNull())
		assert.Empty(t, model.Results.Elements())
	})

	t.Run("200 with results sets state", func(t *testing.T) {
		req, state, diags := testListReadRequest(t, jsonResponse(http.StatusOK, `{"results":[{"name":"one"}],"totalCount":1}`))
		HandleDataSourceReadList(context.Background(), req)
		assertDiagsOK(t, diags)
		var model testListModel
		require.False(t, state.Get(context.Background(), &model).HasError())
		assert.Len(t, model.Results.Elements(), 1)
	})
}

func TestHandleDataSourceRead(t *testing.T) {
	testCases := map[string]struct {
		handler       http.HandlerFunc
		hooks         any
		wantSummary   string
		wantErrDetail string
	}{
		"500 reports the API error instead of not found":      {atlasError(http.StatusInternalServerError, "server error"), nil, "", "server error"},
		"404 reports resource not found":                      {atlasError(http.StatusNotFound, "not found"), nil, "Resource not found", ""},
		"hook signaling not found reports resource not found": {jsonResponse(http.StatusOK, `{"secrets":[]}`), notFoundSignalHook, "Resource not found", ""},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			req, _, diags := testReadRequest(t, tc.handler)
			req.Hooks = tc.hooks
			HandleDataSourceRead(context.Background(), req)
			assertDiagError(t, diags, tc.wantSummary, tc.wantErrDetail)
		})
	}
}
