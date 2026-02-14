package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTAuthRejectsMissingHeader(t *testing.T) {
	t.Parallel()

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	handler := JWTAuth("test-secret")(next)
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.False(t, nextCalled)
}

func TestJWTAuthPassesValidTokenAndStoresClaims(t *testing.T) {
	t.Parallel()

	claims := JWTClaims{
		Sub:            "user-id",
		Email:          "user@example.com",
		Role:           "ADMIN",
		OrganizationID: "org-id",
	}

	token, err := GenerateToken("test-secret", claims)
	require.NoError(t, err)

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		current := CurrentUser(r.Context())
		require.NotNil(t, current)
		assert.Equal(t, "user-id", current.Sub)
		assert.Equal(t, "user@example.com", current.Email)
		assert.Equal(t, "ADMIN", current.Role)
		assert.Equal(t, "org-id", current.OrganizationID)
		w.WriteHeader(http.StatusOK)
	})

	handler := JWTAuth("test-secret")(next)
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.True(t, nextCalled)
	assert.Equal(t, http.StatusOK, w.Code)
}
