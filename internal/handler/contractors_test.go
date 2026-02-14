package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsePaginationDefaults(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/contractors", nil)
	page, limit, err := parsePagination(req)
	require.NoError(t, err)
	assert.Equal(t, 1, page)
	assert.Equal(t, 20, limit)
}

func TestParsePaginationInvalid(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/contractors?page=abc", nil)
	_, _, err := parsePagination(req)
	assert.Error(t, err)
}

func TestParseQueryContractorsRequest(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/contractors?region=Moscow&type=IP&isAvailable=true&search=ivan&page=2&limit=50", nil)
	parsed, err := parseQueryContractorsRequest(req)
	require.NoError(t, err)

	require.NotNil(t, parsed.Region)
	assert.Equal(t, "Moscow", *parsed.Region)
	require.NotNil(t, parsed.Type)
	assert.Equal(t, "IP", *parsed.Type)
	require.NotNil(t, parsed.IsAvailable)
	assert.True(t, *parsed.IsAvailable)
	require.NotNil(t, parsed.Search)
	assert.Equal(t, "ivan", *parsed.Search)
	assert.Equal(t, 2, parsed.Page)
	assert.Equal(t, 50, parsed.Limit)
}

func TestParseQueryContractorsRequestInvalidBool(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/contractors?isAvailable=nope", nil)
	_, err := parseQueryContractorsRequest(req)
	assert.Error(t, err)
}
