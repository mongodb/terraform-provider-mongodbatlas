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
	t.Run("500 reports error and keeps state", func(t *testing.T) {
		req, state, diags := testReadRequest(t, atlasError(http.StatusInternalServerError, "server error"))
		HandleRead(context.Background(), req)
		require.True(t, diags.HasError(), "expected a diagnostics error")
		assert.Contains(t, diags.Errors()[0].Detail(), "server error")
		assert.False(t, state.Raw.IsNull(), "resource must not be removed from state on API error")
	})

	t.Run("404 removes resource from state", func(t *testing.T) {
		req, state, diags := testReadRequest(t, atlasError(http.StatusNotFound, "not found"))
		HandleRead(context.Background(), req)
		require.False(t, diags.HasError(), "unexpected diagnostics: %#v", diags.Errors())
		assert.True(t, state.Raw.IsNull(), "resource must be removed from state on 404")
	})

	t.Run("hook signaling not found removes resource from state", func(t *testing.T) {
		req, state, diags := testReadRequest(t, jsonResponse(http.StatusOK, `{"secrets":[]}`))
		req.Hooks = &testPostReadHook{
			postRead: func(_ HandleReadReq, result APICallResult) APICallResult {
				return APICallResult{Err: fmt.Errorf("secret not found in response: %w", ErrNotFound)}
			},
		}
		HandleRead(context.Background(), req)
		require.False(t, diags.HasError(), "unexpected diagnostics: %#v", diags.Errors())
		assert.True(t, state.Raw.IsNull(), "resource must be removed from state when a hook signals not found")
	})

	t.Run("200 sets state", func(t *testing.T) {
		req, state, diags := testReadRequest(t, jsonResponse(http.StatusOK, `{"name":"updated"}`))
		HandleRead(context.Background(), req)
		require.False(t, diags.HasError(), "unexpected diagnostics: %#v", diags.Errors())
		var model testReadModel
		require.False(t, state.Get(context.Background(), &model).HasError())
		assert.Equal(t, "updated", model.Name.ValueString())
	})
}

type testListResultModel struct {
	Name types.String `tfsdk:"name"`
}

type testListModel struct {
	Results customtypes.NestedListValue[testListResultModel] `tfsdk:"results" autogen:"omitjson"`
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
		req, _, diags := testListReadRequest(t, atlasError(http.StatusInternalServerError, "server error"))
		HandleDataSourceReadList(context.Background(), req)
		require.True(t, diags.HasError(), "expected a diagnostics error")
		assert.Contains(t, diags.Errors()[0].Detail(), "server error")
	})

	t.Run("200 empty JSON reports resource not found", func(t *testing.T) {
		req, _, diags := testListReadRequest(t, jsonResponse(http.StatusOK, `{}`))
		HandleDataSourceReadList(context.Background(), req)
		require.True(t, diags.HasError(), "expected a diagnostics error")
		assert.Contains(t, diags.Errors()[0].Detail(), "resource not found")
	})

	t.Run("200 with results sets state", func(t *testing.T) {
		req, state, diags := testListReadRequest(t, jsonResponse(http.StatusOK, `{"results":[{"name":"one"}],"totalCount":1}`))
		HandleDataSourceReadList(context.Background(), req)
		require.False(t, diags.HasError(), "unexpected diagnostics: %#v", diags.Errors())
		var model testListModel
		require.False(t, state.Get(context.Background(), &model).HasError())
		assert.Len(t, model.Results.Elements(), 1)
	})
}

func TestHandleDataSourceRead(t *testing.T) {
	t.Run("500 reports the API error instead of not found", func(t *testing.T) {
		req, _, diags := testReadRequest(t, atlasError(http.StatusInternalServerError, "server error"))
		HandleDataSourceRead(context.Background(), req)
		require.True(t, diags.HasError(), "expected a diagnostics error")
		assert.Contains(t, diags.Errors()[0].Detail(), "server error")
	})

	t.Run("404 reports resource not found", func(t *testing.T) {
		req, _, diags := testReadRequest(t, atlasError(http.StatusNotFound, "not found"))
		HandleDataSourceRead(context.Background(), req)
		require.True(t, diags.HasError(), "expected a diagnostics error")
		assert.Equal(t, "Resource not found", diags.Errors()[0].Summary())
	})
}
