package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateJobRequest_Valid(t *testing.T) {
	req := CreateJobRequest{
		Title:  "Самосвал на объект",
		Region: "Moscow",
	}
	err := ValidateStruct(req)
	assert.NoError(t, err)
}

func TestCreateJobRequest_MissingTitle(t *testing.T) {
	req := CreateJobRequest{
		Region: "Moscow",
	}
	err := ValidateStruct(req)
	assert.Error(t, err)
	errs := ValidationErrors(err)
	assert.Contains(t, errs, "title")
}

func TestCreateJobRequest_ShortTitle(t *testing.T) {
	req := CreateJobRequest{
		Title:  "ab",
		Region: "Moscow",
	}
	err := ValidateStruct(req)
	assert.Error(t, err)
	errs := ValidationErrors(err)
	assert.Contains(t, errs, "title")
}

func TestCreateJobRequest_MissingRegion(t *testing.T) {
	req := CreateJobRequest{
		Title: "Some job",
	}
	err := ValidateStruct(req)
	assert.Error(t, err)
	errs := ValidationErrors(err)
	assert.Contains(t, errs, "region")
}

func TestQueryJobsRequest_ValidStatus(t *testing.T) {
	status := "DRAFT"
	req := QueryJobsRequest{
		Status: &status,
		Page:   1,
		Limit:  20,
	}
	err := ValidateStruct(req)
	assert.NoError(t, err)
}

func TestQueryJobsRequest_InvalidStatus(t *testing.T) {
	status := "INVALID"
	req := QueryJobsRequest{
		Status: &status,
		Page:   1,
		Limit:  20,
	}
	err := ValidateStruct(req)
	assert.Error(t, err)
	errs := ValidationErrors(err)
	assert.Contains(t, errs, "status")
}
