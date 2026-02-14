package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseQueryJobsRequest(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/jobs?status=DRAFT&region=Moscow&page=3&limit=15", nil)
	parsed, err := parseQueryJobsRequest(req)
	require.NoError(t, err)

	require.NotNil(t, parsed.Status)
	assert.Equal(t, "DRAFT", *parsed.Status)
	require.NotNil(t, parsed.Region)
	assert.Equal(t, "Moscow", *parsed.Region)
	assert.Equal(t, 3, parsed.Page)
	assert.Equal(t, 15, parsed.Limit)
}

func TestParseQueryJobsRequestDefaultPagination(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/jobs", nil)
	parsed, err := parseQueryJobsRequest(req)
	require.NoError(t, err)

	assert.Equal(t, 1, parsed.Page)
	assert.Equal(t, 20, parsed.Limit)
}

func TestParseQueryJobsRequestInvalidPagination(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/jobs?page=oops", nil)
	_, err := parseQueryJobsRequest(req)
	assert.Error(t, err)
}
