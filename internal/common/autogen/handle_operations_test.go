package autogen_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen/customtypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func testClient(t *testing.T, handler http.Handler) *config.MongoDBClient {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return newClient(t, srv.URL)
}

// unreachableClient points at a server that is already closed, so calls fail without an HTTP response.
func unreachableClient(t *testing.T) *config.MongoDBClient {
	t.Helper()
	srv := httptest.NewServer(http.NotFoundHandler())
	srv.Close()
	return newClient(t, srv.URL)
}

func newClient(t *testing.T, baseURL string) *config.MongoDBClient {
	t.Helper()
	client, err := config.NewClient(&config.Credentials{BaseURL: baseURL}, "test")
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

// deleteThenPoll answers the DELETE of the resource and then every wait poll, so a single handler can
// drive HandleDelete end to end.
func deleteThenPoll(deleteHandler, pollHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			deleteHandler(w, r)
			return
		}
		pollHandler(w, r)
	}
}

func testDeleteRequest(t *testing.T, handler http.Handler) (autogen.HandleDeleteReq, *diag.Diagnostics) {
	t.Helper()
	diags := &diag.Diagnostics{}
	return autogen.HandleDeleteReq{
		Client:     testClient(t, handler),
		State:      &testReadModel{},
		Diags:      diags,
		CallParams: &config.APICallParams{RelativePath: "/api/test", Method: http.MethodDelete},
	}, diags
}

func TestHandleDelete(t *testing.T) {
	testCases := map[string]struct {
		handler             http.HandlerFunc
		expectedErrContains string
	}{
		"400 reports error":                {atlasError(http.StatusBadRequest, "cannot delete resource"), "cannot delete resource"},
		"401 reports error":                {atlasError(http.StatusUnauthorized, "unauthorized"), "unauthorized"},
		"409 reports error":                {atlasError(http.StatusConflict, "conflict"), "conflict"},
		"500 reports error":                {atlasError(http.StatusInternalServerError, "server error"), "server error"},
		"404 is tolerated already deleted": {atlasError(http.StatusNotFound, "not found"), ""},
		"204 succeeds":                     {jsonResponse(http.StatusNoContent, ""), ""},
		"200 with body succeeds":           {jsonResponse(http.StatusOK, `{"stateName":"DELETING"}`), ""},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			req, diags := testDeleteRequest(t, tc.handler)
			autogen.HandleDelete(context.Background(), req)
			if tc.expectedErrContains == "" {
				assertDiagsOK(t, diags)
				return
			}
			assertDiagError(t, diags, "", tc.expectedErrContains)
		})
	}
}

func deleteWaitReq() *autogen.WaitReq {
	return &autogen.WaitReq{
		StateProperty: "stateName",
		PendingStates: []string{"DELETING"},
		TargetStates:  []string{retrystrategy.RetryStrategyDeletedState},
		Timeout:       10 * time.Second,
		CallParams: func(_ any) *config.APICallParams {
			return &config.APICallParams{RelativePath: "/api/test", Method: http.MethodGet}
		},
	}
}

// notFoundHook is the minimal hook needed to raise ErrNotFound. Hook behavior itself is tested in the
// packages that define the hooks; this only provides the sentinel, which no HTTP response can carry.
type notFoundHook struct{}

func (notFoundHook) PostReadAPICall(_ autogen.HandleReadReq, _ autogen.APICallResult) autogen.APICallResult {
	return autogen.APICallResult{Err: fmt.Errorf("resource missing from list response: %w", autogen.ErrNotFound)}
}

// TestErrNotFoundSentinel covers the not-found classification for APIs that cannot express it as an
// HTTP 404, at all three call sites of notFound.
func TestErrNotFoundSentinel(t *testing.T) {
	okResponse := jsonResponse(http.StatusOK, `{"name":"present"}`) // the sentinel comes from the hook, not the status code

	t.Run("resource read removes the resource from state", func(t *testing.T) {
		req, state, diags := testReadRequest(t, testClient(t, okResponse))
		req.Hooks = notFoundHook{}
		autogen.HandleRead(context.Background(), req)
		assertDiagsOK(t, diags)
		assert.True(t, state.Raw.IsNull(), "resource must be removed from state")
	})

	t.Run("data source read reports resource not found", func(t *testing.T) {
		req, _, diags := testReadRequest(t, testClient(t, okResponse))
		req.Hooks = notFoundHook{}
		autogen.HandleDataSourceRead(context.Background(), req)
		assertDiagError(t, diags, "Resource not found", "")
	})

	t.Run("delete wait poll completes the delete", func(t *testing.T) {
		req, diags := testDeleteRequest(t, deleteThenPoll(jsonResponse(http.StatusNoContent, ""), okResponse))
		req.Hooks = notFoundHook{}
		req.Wait = deleteWaitReq()
		autogen.HandleDelete(context.Background(), req)
		assertDiagsOK(t, diags)
	})
}

// TestHandleDeleteWait covers the wait polling of a delete: only a 404 confirms the resource is gone,
// and any polling error must abort the wait.
func TestHandleDeleteWait(t *testing.T) {
	testCases := map[string]struct {
		poll                http.HandlerFunc
		expectedErrContains string
	}{
		"404 poll completes the delete":            {atlasError(http.StatusNotFound, "not found"), ""},
		"poll reaching the target state completes": {jsonResponse(http.StatusOK, `{"stateName":"DELETED"}`), ""},
		"500 poll aborts the wait":                 {atlasError(http.StatusInternalServerError, "server error"), "server error"},
		"401 poll aborts the wait":                 {atlasError(http.StatusUnauthorized, "unauthorized"), "unauthorized"},
		"empty JSON poll aborts the wait":          {jsonResponse(http.StatusOK, `{}`), "wait state attribute not found"},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			req, diags := testDeleteRequest(t, deleteThenPoll(jsonResponse(http.StatusNoContent, ""), tc.poll))
			req.Wait = deleteWaitReq()
			autogen.HandleDelete(context.Background(), req)
			if tc.expectedErrContains == "" {
				assertDiagsOK(t, diags)
				return
			}
			assertDiagError(t, diags, "", tc.expectedErrContains)
		})
	}
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

type testReadModel struct {
	Name types.String `tfsdk:"name" apiname:"name"`
}

func testReadRequest(t *testing.T, client *config.MongoDBClient) (autogen.HandleReadReq, *tfsdk.State, *diag.Diagnostics) {
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
	return autogen.HandleReadReq{
		State:      &testReadModel{Name: types.StringValue("original")},
		RespState:  state,
		Client:     client,
		CallParams: &config.APICallParams{RelativePath: "/api/test", Method: http.MethodGet},
		RespDiags:  diags,
	}, state, diags
}

func TestHandleRead(t *testing.T) {
	testCases := map[string]struct {
		handler       http.HandlerFunc
		wantErrDetail string
		wantName      string
		unreachable   bool
		wantErr       bool
		wantRemoved   bool
	}{
		"400 reports error and keeps state":                   {handler: atlasError(http.StatusBadRequest, "bad request"), wantErr: true, wantErrDetail: "bad request"},
		"401 reports error and keeps state":                   {handler: atlasError(http.StatusUnauthorized, "unauthorized"), wantErr: true, wantErrDetail: "unauthorized"},
		"409 reports error and keeps state":                   {handler: atlasError(http.StatusConflict, "conflict"), wantErr: true, wantErrDetail: "conflict"},
		"500 reports error and keeps state":                   {handler: atlasError(http.StatusInternalServerError, "server error"), wantErr: true, wantErrDetail: "server error"},
		"transport error reports error and keeps state":       {unreachable: true, wantErr: true},
		"200 with invalid JSON reports error and keeps state": {handler: jsonResponse(http.StatusOK, `not-json`), wantErr: true},

		"404 removes resource from state": {handler: atlasError(http.StatusNotFound, "not found"), wantRemoved: true},
		"200 sets state":                  {handler: jsonResponse(http.StatusOK, `{"name":"updated"}`), wantName: "updated"},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			client := unreachableClient(t)
			if !tc.unreachable {
				client = testClient(t, tc.handler)
			}
			req, state, diags := testReadRequest(t, client)
			autogen.HandleRead(context.Background(), req)
			if tc.wantErr {
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
	Results customtypes.NestedListValue[testListResultModel] `tfsdk:"results" autogen:"omitjson"`
}

func testListReadRequest(t *testing.T, handler http.Handler) (autogen.HandleReadReq, *tfsdk.State, *diag.Diagnostics) {
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
	return autogen.HandleReadReq{
		State:      &testListModel{},
		RespState:  state,
		Client:     testClient(t, handler),
		CallParams: &config.APICallParams{RelativePath: "/api/test", Method: http.MethodGet},
		RespDiags:  diags,
	}, state, diags
}

func TestHandleDataSourceReadList(t *testing.T) {
	t.Run("500 reports the API error", func(t *testing.T) {
		req, _, diags := testReadRequest(t, testClient(t, atlasError(http.StatusInternalServerError, "server error"))) // list schema not needed to assert the error
		autogen.HandleDataSourceReadList(context.Background(), req)
		assertDiagError(t, diags, "", "server error")
	})

	t.Run("200 empty JSON sets empty results", func(t *testing.T) {
		req, state, diags := testListReadRequest(t, jsonResponse(http.StatusOK, `{}`))
		autogen.HandleDataSourceReadList(context.Background(), req)
		assertDiagsOK(t, diags)
		var model testListModel
		require.False(t, state.Get(context.Background(), &model).HasError())
		assert.Empty(t, model.Results.Elements())
	})

	t.Run("200 with results sets state", func(t *testing.T) {
		req, state, diags := testListReadRequest(t, jsonResponse(http.StatusOK, `{"results":[{"name":"one"}],"totalCount":1}`))
		autogen.HandleDataSourceReadList(context.Background(), req)
		assertDiagsOK(t, diags)
		var model testListModel
		require.False(t, state.Get(context.Background(), &model).HasError())
		assert.Len(t, model.Results.Elements(), 1)
	})
}

func TestHandleDataSourceRead(t *testing.T) {
	testCases := map[string]struct {
		handler       http.HandlerFunc
		wantSummary   string
		wantErrDetail string
	}{
		"400 reports the API error instead of not found": {atlasError(http.StatusBadRequest, "bad request"), "", "bad request"},
		"500 reports the API error instead of not found": {atlasError(http.StatusInternalServerError, "server error"), "", "server error"},
		"404 reports resource not found":                 {atlasError(http.StatusNotFound, "not found"), "Resource not found", ""},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			req, _, diags := testReadRequest(t, testClient(t, tc.handler))
			autogen.HandleDataSourceRead(context.Background(), req)
			assertDiagError(t, diags, tc.wantSummary, tc.wantErrDetail)
		})
	}
}
