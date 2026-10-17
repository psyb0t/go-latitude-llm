package latitude_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/psyb0t/go-latitude-llm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testAPIKey      = "EXAMPLE-DO-NOT-USE"
	testProjectName = "Example project"
	testProjectJSON = `{
		"id":"example-project-id",
		"organizationId":"example-organization-id",
		"name":"Example project",
		"slug":"example-project",
		"settings":null,
		"firstTraceAt":null,
		"deletedAt":null,
		"lastEditedAt":"2026-10-17T00:00:00Z",
		"createdAt":"2026-10-17T00:00:00Z",
		"updatedAt":"2026-10-17T00:00:00Z"
	}`
)

func newTestClient(
	t *testing.T,
	handler http.HandlerFunc,
) *latitude.ClientWithResponses {
	t.Helper()

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	auth := latitude.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+testAPIKey)

			return nil
		},
	)

	client, err := latitude.NewClientWithResponses(srv.URL, auth)
	require.NoError(t, err)

	return client
}

func assertRequest(
	t *testing.T,
	r *http.Request,
	method, path string,
) {
	t.Helper()
	assert.Equal(t, method, r.Method)
	assert.Equal(t, path, r.URL.Path)
	assert.Equal(t, "Bearer "+testAPIKey, r.Header.Get("Authorization"))
}

func writeJSON(
	t *testing.T,
	w http.ResponseWriter,
	status int,
	body string,
) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err := w.Write([]byte(body))
	assert.NoError(t, err)
}

func TestListProjects(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertRequest(t, r, http.MethodGet, "/v1/projects")
		assert.Empty(t, r.URL.RawQuery)
		writeJSON(t, w, http.StatusOK, `{
			"items":[`+testProjectJSON+`],
			"nextCursor":null,
			"hasMore":false
		}`)
	})

	res, err := client.ListProjectsWithResponse(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.NotNil(t, res.JSON200)
	require.Len(t, res.JSON200.Items, 1)
	assert.Equal(t, testProjectName, res.JSON200.Items[0].Name)
	assert.Nil(t, res.JSON200.Items[0].Settings)
	assert.Nil(t, res.JSON200.Items[0].FirstTraceAt)
	assert.Nil(t, res.JSON200.Items[0].DeletedAt)
	assert.Nil(t, res.JSON200.NextCursor)
	assert.False(t, res.JSON200.HasMore)
}

func TestCreateProject(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertRequest(t, r, http.MethodPost, "/v1/projects")
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var body latitude.CreateProjectBody

		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, testProjectName, body.Name)
		writeJSON(t, w, http.StatusCreated, testProjectJSON)
	})

	res, err := client.CreateProjectWithResponse(
		context.Background(),
		latitude.CreateProjectBody{Name: testProjectName},
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode())
	require.NotNil(t, res.JSON201)
	assert.Equal(t, "example-project-id", res.JSON201.Id)
	assert.Equal(t, testProjectName, res.JSON201.Name)
}

func TestGetAccountWithAPIKey(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertRequest(t, r, http.MethodGet, "/v1/account")
		writeJSON(t, w, http.StatusOK, `{
			"user":null,
			"role":null,
			"organization":{
				"id":"example-organization-id",
				"name":"Example organization",
				"slug":"example-organization"
			}
		}`)
	})

	res, err := client.GetAccountWithResponse(context.Background())
	require.NoError(t, err)
	require.NotNil(t, res.JSON200)
	assert.Nil(t, res.JSON200.User)
	assert.Nil(t, res.JSON200.Role)
	assert.Equal(t, "example-organization-id", res.JSON200.Organization.Id)
}

func TestListProjectsUnauthorized(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertRequest(t, r, http.MethodGet, "/v1/projects")
		writeJSON(t, w, http.StatusUnauthorized, `{"error":"Unauthorized"}`)
	})

	res, err := client.ListProjectsWithResponse(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode())
	assert.Nil(t, res.JSON200)
	require.NotNil(t, res.JSON401)
	assert.Equal(t, "Unauthorized", res.JSON401.Error)
}

func TestListProjectsMalformedJSON(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertRequest(t, r, http.MethodGet, "/v1/projects")
		writeJSON(t, w, http.StatusOK, `{"items":invalid}`)
	})

	res, err := client.ListProjectsWithResponse(context.Background())

	var syntaxError *json.SyntaxError

	require.ErrorAs(t, err, &syntaxError)
	assert.Nil(t, res)
}

func TestListProjectsCanceledContext(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(_ http.ResponseWriter, _ *http.Request) {
		t.Error("canceled request reached the server")
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	res, err := client.ListProjectsWithResponse(ctx)
	require.ErrorIs(t, err, context.Canceled)
	assert.Nil(t, res)
}

func TestCustomScoreFloat64AndNull(t *testing.T) {
	t.Parallel()

	var score latitude.CustomScoreResponse

	err := json.Unmarshal([]byte(`{
		"value":0.123456789012345,
		"sessionId":null,
		"traceId":null,
		"error":null,
		"draftedAt":null,
		"metadata":{"category":"example"}
	}`), &score)
	require.NoError(t, err)
	assert.InDelta(t, 0.123456789012345, score.Value, 0)
	assert.Nil(t, score.SessionId)
	assert.Nil(t, score.TraceId)
	assert.Nil(t, score.Error)
	assert.Nil(t, score.DraftedAt)
	assert.Equal(t, "example", score.Metadata["category"])
}
