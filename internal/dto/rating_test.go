package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateRatingRequestValidation(t *testing.T) {
	req := CreateRatingRequest{
		JobID:        "3f026267-f3ce-4b70-9ee5-6d82aa4f4f78",
		ContractorID: "de1178db-8326-4f3d-884d-4ea5e18e2a26",
		Score:        5,
	}
	assert.NoError(t, ValidateStruct(req))
}

func TestCreateRatingRequestInvalidScore(t *testing.T) {
	req := CreateRatingRequest{
		JobID:        "3f026267-f3ce-4b70-9ee5-6d82aa4f4f78",
		ContractorID: "de1178db-8326-4f3d-884d-4ea5e18e2a26",
		Score:        10,
	}
	err := ValidateStruct(req)
	assert.Error(t, err)
	assert.Contains(t, ValidationErrors(err), "score")
}
