package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"platform/internal/handler"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthCheck(t *testing.T) {
	r := New(Deps{
		AuthHandler:          handler.NewAuthHandler(nil),
		UsersHandler:         handler.NewUsersHandler(nil),
		OrganizationsHandler: handler.NewOrganizationsHandler(nil),
		ContractorsHandler:   handler.NewContractorsHandler(nil),
		JobsHandler:          handler.NewJobsHandler(nil),
		MatchingHandler:      handler.NewMatchingHandler(nil),
		ResponsesHandler:     handler.NewResponsesHandler(nil),
		AssignmentsHandler:   handler.NewAssignmentsHandler(nil),
		RatingsHandler:       handler.NewRatingsHandler(nil),
		JWTSecret:            "test-secret",
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, "ok", resp["status"])
}

func TestCORS(t *testing.T) {
	r := New(Deps{
		AuthHandler:          handler.NewAuthHandler(nil),
		UsersHandler:         handler.NewUsersHandler(nil),
		OrganizationsHandler: handler.NewOrganizationsHandler(nil),
		ContractorsHandler:   handler.NewContractorsHandler(nil),
		JobsHandler:          handler.NewJobsHandler(nil),
		MatchingHandler:      handler.NewMatchingHandler(nil),
		ResponsesHandler:     handler.NewResponsesHandler(nil),
		AssignmentsHandler:   handler.NewAssignmentsHandler(nil),
		RatingsHandler:       handler.NewRatingsHandler(nil),
		JWTSecret:            "test-secret",
	})

	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Authorization")
}
